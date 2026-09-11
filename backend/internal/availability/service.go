package availability

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service coordinates lightweight session heartbeats, gap detection, and availability metrics.
type Service struct {
	pool          *pgxpool.Pool
	mu            sync.RWMutex
	lastLatency   int
	lastLatencyAt time.Time
	stopWorker    chan struct{}
}

// NewService constructs a new availability Service instance with the database connection pool.
func NewService(pool *pgxpool.Pool) *Service {
	s := &Service{
		pool:       pool,
		stopWorker: make(chan struct{}),
	}
	// Initial async probe for external latency
	go s.probeExternalLatency()
	return s
}

// Initialize performs startup gap detection using native Linux /proc telemetry.
// If previous last_seen_at exists and gap > 7 minutes, it estimates downtime and creates an incident record.
// Persists session metadata to availability_state without inserting a row for every tick.
func (s *Service) Initialize(ctx context.Context, uptimeSec int64) error {
	bootID := bootID()
	now := time.Now()
	if uptimeSec <= 0 {
		uptimeSec = uptimeSeconds()
	}
	bootTime := now.Add(-time.Duration(uptimeSec) * time.Second)

	var (
		stateID          int
		prevBootID       string
		prevStartedAt    time.Time
		prevLastSeenAt   time.Time
		firstMonitoredAt time.Time
		heartbeatCount   int64
	)

	err := s.pool.QueryRow(ctx, `
		SELECT id, boot_id, started_at, last_seen_at, first_monitored_at, heartbeat_count
		FROM availability_state
		WHERE id = 1
	`).Scan(&stateID, &prevBootID, &prevStartedAt, &prevLastSeenAt, &firstMonitoredAt, &heartbeatCount)

	if errors.Is(err, pgx.ErrNoRows) {
		// First run ever on this machine
		_, err = s.pool.Exec(ctx, `
			INSERT INTO availability_state (id, boot_id, started_at, last_seen_at, first_monitored_at, heartbeat_count, current_status)
			VALUES (1, $1, $2, $3, $4, 1, 'OPERATIONAL')
		`, bootID, bootTime, now, bootTime)
		if err != nil {
			return fmt.Errorf("failed to initialize availability_state: %w", err)
		}
		log.Printf("[Availability] Initialized monitoring baseline (boot_id=%s, bootTime=%s)", bootID, bootTime.Format(time.RFC3339))
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to query availability_state: %w", err)
	}

	// Gap detection: check if there was a meaningful downtime between previous heartbeat and now
	gapDuration := now.Sub(prevLastSeenAt)
	if gapDuration > 7*time.Minute {
		var incidentStart, incidentEnd time.Time
		var cause, details string

		if prevBootID != bootID {
			// System rebooted! Downtime is from prevLastSeenAt until current bootTime
			incidentStart = prevLastSeenAt
			if bootTime.After(prevLastSeenAt) {
				incidentEnd = bootTime
			} else {
				incidentEnd = now
			}
			cause = "SYSTEM_REBOOT"
			details = fmt.Sprintf("Host reboot detected via boot_id change (previous session active until %s)", prevLastSeenAt.Format("Jan 02, 15:04 MST"))
		} else {
			// Same boot ID, but service was stopped or suspended
			incidentStart = prevLastSeenAt
			incidentEnd = now
			cause = "SERVICE_SUSPENDED"
			details = fmt.Sprintf("Service downtime gap detected between heartbeats (%s to %s)", prevLastSeenAt.Format("15:04"), now.Format("15:04 MST"))
		}

		durSec := int64(incidentEnd.Sub(incidentStart).Seconds())
		if durSec >= 120 { // only record meaningful outages (>= 2 minutes)
			_, err = s.pool.Exec(ctx, `
				INSERT INTO availability_incidents (started_at, ended_at, duration_seconds, cause, details)
				VALUES ($1, $2, $3, $4, $5)
			`, incidentStart, incidentEnd, durSec, cause, details)
			if err != nil {
				log.Printf("[Availability] Warning: failed to record incident: %v", err)
			} else {
				log.Printf("[Availability] Recorded downtime incident: %s (%ds gap)", cause, durSec)
			}
		}
	}

	// Update existing state in-place with current session and last_seen
	_, err = s.pool.Exec(ctx, `
		UPDATE availability_state
		SET boot_id = $1, started_at = $2, last_seen_at = $3, heartbeat_count = heartbeat_count + 1, current_status = 'OPERATIONAL'
		WHERE id = 1
	`, bootID, bootTime, now)
	return err
}

// StartHeartbeatWorker runs a background ticker that performs lightweight in-place heartbeat updates.
func (s *Service) StartHeartbeatWorker(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.stopWorker:
				return
			case <-ticker.C:
				s.tickHeartbeat(ctx)
			}
		}
	}()
	log.Printf("[Availability] Lightweight heartbeat worker active (cadence: %v)", interval)
}

func (s *Service) tickHeartbeat(ctx context.Context) {
	now := time.Now()
	// Single in-place update — zero row spamming!
	_, err := s.pool.Exec(ctx, `
		UPDATE availability_state
		SET last_seen_at = $1, heartbeat_count = heartbeat_count + 1, current_status = 'OPERATIONAL'
		WHERE id = 1
	`, now)
	if err != nil {
		log.Printf("[Availability] Heartbeat update error: %v", err)
	}
}

// Stop cleanly terminates background workers.
func (s *Service) Stop() {
	close(s.stopWorker)
}

package availability

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ProjectAvailabilityStats contains aggregated availability metrics for a project.
type ProjectAvailabilityStats struct {
	UptimePercent       *float64  `json:"uptime_percent"`
	LatestResponseTimeMs int       `json:"latest_response_time_ms"`
	LastCheckedAt       *time.Time `json:"last_checked_at"`
	TotalChecks30d      int       `json:"total_checks_30d"`
	CurrentStatus       string    `json:"current_status"`
}

type projectProbeTarget struct {
	ID        string
	Name      string
	PublicURL string
	Status    string
}

// StartProjectCheckWorker starts the periodic HTTP health check worker for published projects.
func (s *Service) StartProjectCheckWorker(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	ticker := time.NewTicker(interval)

	// Run an initial check after a short delay on startup
	go func() {
		time.Sleep(10 * time.Second)
		s.ProbePublishedProjects(ctx)
	}()

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.stopWorker:
				return
			case <-ticker.C:
				s.ProbePublishedProjects(ctx)
			}
		}
	}()
	log.Printf("[Availability] Project health probe worker active (cadence: %v)", interval)
}

// ProbePublishedProjects executes bounded concurrent HTTP health checks against all published projects.
func (s *Service) ProbePublishedProjects(ctx context.Context) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, public_url, status
		FROM projects
		WHERE visibility = 'PUBLIC' AND public_url != '' AND status IN ('ONLINE', 'OFFLINE')
	`)
	if err != nil {
		log.Printf("[Availability] Error querying projects for probe: %v", err)
		return
	}
	defer rows.Close()

	targets := make([]projectProbeTarget, 0)
	for rows.Next() {
		var t projectProbeTarget
		if err := rows.Scan(&t.ID, &t.Name, &t.PublicURL, &t.Status); err == nil {
			targets = append(targets, t)
		}
	}

	if len(targets) == 0 {
		return
	}

	// Bounded concurrency worker pool (max 8 concurrent checks)
	concurrency := 8
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	for _, target := range targets {
		wg.Add(1)
		sem <- struct{}{}
		go func(t projectProbeTarget) {
			defer wg.Done()
			defer func() { <-sem }()
			s.probeSingleProject(ctx, client, t)
		}(target)
	}

	wg.Wait()
}

func (s *Service) probeSingleProject(ctx context.Context, client *http.Client, t projectProbeTarget) {
	start := time.Now()
	var statusCode int
	var isSuccess bool
	var responseTimeMs int

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.PublicURL, nil)
	if err != nil {
		isSuccess = false
	} else {
		req.Header.Set("User-Agent", "NgumpulHost-AvailabilityProbe/2.0")
		resp, err := client.Do(req)
		responseTimeMs = int(time.Since(start).Milliseconds())
		if responseTimeMs <= 0 {
			responseTimeMs = 1
		}

		if err == nil {
			statusCode = resp.StatusCode
			_ = resp.Body.Close()
			// HTTP 200 to 399 considered healthy
			isSuccess = (statusCode >= 200 && statusCode < 400)
		} else {
			isSuccess = false
			statusCode = 0
		}
	}

	// Query latest consecutive failures
	var prevFailures int
	_ = s.pool.QueryRow(ctx, `
		SELECT consecutive_failures
		FROM project_availability_checks
		WHERE project_id = $1
		ORDER BY checked_at DESC
		LIMIT 1
	`, t.ID).Scan(&prevFailures)

	consecutiveFailures := 0
	if !isSuccess {
		consecutiveFailures = prevFailures + 1
	}

	// Record probe record
	_, _ = s.pool.Exec(ctx, `
		INSERT INTO project_availability_checks (project_id, checked_at, status_code, response_time_ms, is_successful, consecutive_failures)
		VALUES ($1, NOW(), $2, $3, $4, $5)
	`, t.ID, statusCode, responseTimeMs, isSuccess, consecutiveFailures)

	// State transitions
	if !isSuccess && consecutiveFailures >= 2 && t.Status == "ONLINE" {
		// Mark OFFLINE after 2 consecutive failures
		_, err := s.pool.Exec(ctx, `
			UPDATE projects
			SET status = 'OFFLINE', updated_at = NOW()
			WHERE id = $1
		`, t.ID)
		if err == nil {
			log.Printf("[Availability] Project %s marked OFFLINE after %d consecutive failures", t.Name, consecutiveFailures)
			_, _ = s.pool.Exec(ctx, `
				INSERT INTO activities (project_id, type, metadata, visibility)
				VALUES ($1, 'PROJECT_OFFLINE', $2, 'PUBLIC')
			`, t.ID, map[string]string{
				"title":   fmt.Sprintf("%s is currently unavailable", t.Name),
				"message": fmt.Sprintf("Health check failed (HTTP %d, latency %dms)", statusCode, responseTimeMs),
			})
		}
	} else if isSuccess && t.Status == "OFFLINE" {
		// Immediate recovery to ONLINE
		_, err := s.pool.Exec(ctx, `
			UPDATE projects
			SET status = 'ONLINE', updated_at = NOW()
			WHERE id = $1
		`, t.ID)
		if err == nil {
			log.Printf("[Availability] Project %s recovered to ONLINE", t.Name)
			_, _ = s.pool.Exec(ctx, `
				INSERT INTO activities (project_id, type, metadata, visibility)
				VALUES ($1, 'PROJECT_ONLINE', $2, 'PUBLIC')
			`, t.ID, map[string]string{
				"title":   fmt.Sprintf("%s recovered and is online", t.Name),
				"message": fmt.Sprintf("Health check passed (HTTP %d, latency %dms)", statusCode, responseTimeMs),
			})
		}
	}
}

// GetProjectAvailabilityStats calculates 30-day uptime and latency for a specific project.
func (s *Service) GetProjectAvailabilityStats(ctx context.Context, projectID string) (ProjectAvailabilityStats, error) {
	var stats ProjectAvailabilityStats

	// Query latest check and current status
	err := s.pool.QueryRow(ctx, `
		SELECT p.status, c.checked_at, c.response_time_ms
		FROM projects p
		LEFT JOIN LATERAL (
			SELECT checked_at, response_time_ms
			FROM project_availability_checks
			WHERE project_id = p.id
			ORDER BY checked_at DESC
			LIMIT 1
		) c ON true
		WHERE p.id = $1
	`, projectID).Scan(&stats.CurrentStatus, &stats.LastCheckedAt, &stats.LatestResponseTimeMs)
	if err != nil {
		return stats, err
	}

	// 30-day checks summary
	var totalChecks, successfulChecks int
	err = s.pool.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE is_successful = true)
		FROM project_availability_checks
		WHERE project_id = $1 AND checked_at >= NOW() - INTERVAL '30 days'
	`, projectID).Scan(&totalChecks, &successfulChecks)
	if err == nil {
		stats.TotalChecks30d = totalChecks
		if totalChecks > 0 {
			pct := (float64(successfulChecks) / float64(totalChecks)) * 100.0
			stats.UptimePercent = &pct
		} else {
			// No checks recorded yet, default to 100% if project is ONLINE
			if stats.CurrentStatus == "ONLINE" {
				pct := 100.0
				stats.UptimePercent = &pct
			}
		}
	}

	return stats, nil
}

package availability

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ProjectAvailabilityBlock represents a bucketed slice of project availability.
type ProjectAvailabilityBlock struct {
	Index            int      `json:"index"`
	Label            string   `json:"label"`
	Date             string   `json:"date"`
	Status           string   `json:"status"` // "no_data" | "operational" | "partial" | "incident"
	UptimePercent    *float64 `json:"uptimePercent"`
	SuccessfulChecks int      `json:"successful_checks"`
	TotalChecks      int      `json:"total_checks"`
	Details          string   `json:"details"`
}

// ProjectRangeAvailability contains availability stats and blocks for a single range.
type ProjectRangeAvailability struct {
	RangeFilter      string                     `json:"range"` // "1d", "7d", "30d"
	UptimePercent    *float64                   `json:"uptime_percent"`
	TotalChecks      int                        `json:"total_checks"`
	SuccessfulChecks int                        `json:"successful_checks"`
	Blocks           []ProjectAvailabilityBlock `json:"blocks"`
}

// ProjectAvailabilityStats contains aggregated availability metrics for a project across multiple ranges.
type ProjectAvailabilityStats struct {
	UptimePercent        *float64                            `json:"uptime_percent"`
	LatestResponseTimeMs int                                 `json:"latest_response_time_ms"`
	LastCheckedAt        *time.Time                          `json:"last_checked_at"`
	TotalChecks30d       int                                 `json:"total_checks_30d"`
	CurrentStatus        string                              `json:"current_status"`
	CurrentAvailability  string                              `json:"current_availability"`
	AvailabilityReason   string                              `json:"availability_reason,omitempty"`
	Blocks               []ProjectAvailabilityBlock          `json:"blocks"`
	Ranges               map[string]ProjectRangeAvailability `json:"ranges"`
}

type projectProbeTarget struct {
	ID              string
	Name            string
	PublicURL       string
	Status          string
	Availability    string
	LifecycleStatus string
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

// ProbePublishedProjects executes bounded concurrent HTTP health checks against all active published projects.
func (s *Service) ProbePublishedProjects(ctx context.Context) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, public_url, status, availability, lifecycle_status
		FROM projects
		WHERE lifecycle_status = 'ACTIVE' AND visibility = 'PUBLIC' AND public_url != ''
	`)
	if err != nil {
		log.Printf("[Availability] Error querying projects for probe: %v", err)
		return
	}
	defer rows.Close()

	targets := make([]projectProbeTarget, 0)
	for rows.Next() {
		var t projectProbeTarget
		if err := rows.Scan(&t.ID, &t.Name, &t.PublicURL, &t.Status, &t.Availability, &t.LifecycleStatus); err == nil {
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
	var failureReason string

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.PublicURL, nil)
	if err != nil {
		isSuccess = false
		failureReason = "CONNECTION_FAILED"
	} else {
		req.Header.Set("User-Agent", "NgumpulHost-AvailabilityProbe/2.0")
		resp, doErr := client.Do(req)
		responseTimeMs = int(time.Since(start).Milliseconds())
		if responseTimeMs <= 0 {
			responseTimeMs = 1
		}

		if doErr == nil {
			statusCode = resp.StatusCode
			_ = resp.Body.Close()
			// HTTP 200 to 399 considered healthy
			isSuccess = (statusCode >= 200 && statusCode < 400)
			if !isSuccess {
				failureReason = "HTTP_ERROR"
			}
		} else {
			isSuccess = false
			statusCode = 0
			errStr := strings.ToLower(doErr.Error())
			if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
				failureReason = "TIMEOUT"
			} else if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "dns") {
				failureReason = "DNS_ERROR"
			} else if strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls") || strings.Contains(errStr, "x509") {
				failureReason = "TLS_ERROR"
			} else {
				failureReason = "CONNECTION_FAILED"
			}
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

	// State transitions: only update availability, NEVER touch lifecycle_status
	if !isSuccess && consecutiveFailures >= 2 && t.Availability != "UNREACHABLE" {
		// Mark UNREACHABLE after 2 consecutive failures
		_, updateErr := s.pool.Exec(ctx, `
			UPDATE projects
			SET availability = 'UNREACHABLE', availability_reason = $2, status = 'OFFLINE', updated_at = NOW()
			WHERE id = $1
		`, t.ID, failureReason)
		if updateErr == nil {
			log.Printf("[Availability] Project %s marked UNREACHABLE (%s) after %d consecutive failures", t.Name, failureReason, consecutiveFailures)
			_, _ = s.pool.Exec(ctx, `
				INSERT INTO activities (project_id, type, metadata, visibility)
				VALUES ($1, 'PROJECT_UNREACHABLE', $2, 'PUBLIC')
			`, t.ID, map[string]string{
				"title":   fmt.Sprintf("%s is unreachable", t.Name),
				"message": fmt.Sprintf("Health check failed (%s, HTTP %d, latency %dms)", failureReason, statusCode, responseTimeMs),
				"reason":  failureReason,
			})
		}
	} else if isSuccess && (t.Availability == "UNREACHABLE" || t.Availability == "UNKNOWN") {
		// Recovered back to REACHABLE
		_, updateErr := s.pool.Exec(ctx, `
			UPDATE projects
			SET availability = 'REACHABLE', availability_reason = '', status = 'ONLINE', updated_at = NOW()
			WHERE id = $1
		`, t.ID)
		if updateErr == nil && t.Availability == "UNREACHABLE" {
			log.Printf("[Availability] Project %s recovered to REACHABLE", t.Name)
			_, _ = s.pool.Exec(ctx, `
				INSERT INTO activities (project_id, type, metadata, visibility)
				VALUES ($1, 'PROJECT_REACHABLE', $2, 'PUBLIC')
			`, t.ID, map[string]string{
				"title":   fmt.Sprintf("%s recovered and is reachable", t.Name),
				"message": fmt.Sprintf("Health check passed (HTTP %d, latency %dms)", statusCode, responseTimeMs),
			})
		}
	}
}

type projectCheckEntry struct {
	checkedAt    time.Time
	isSuccessful bool
}

func buildProjectBuckets(
	rangeName string,
	blockCount int,
	blockDuration time.Duration,
	checks []projectCheckEntry,
	createdAt time.Time,
	now time.Time,
) ProjectRangeAvailability {
	blocks := make([]ProjectAvailabilityBlock, 0, blockCount)
	var rangeTotal, rangeSuccess int

	for i := blockCount - 1; i >= 0; i-- {
		var bStart, bEnd time.Time
		var label, dateStr string

		if rangeName == "1d" {
			bStart = now.Add(-time.Duration(i) * 30 * time.Minute).Truncate(30 * time.Minute)
			bEnd = bStart.Add(30 * time.Minute)
			label = bStart.Format("15:04")
			dateStr = bStart.Format("2006-01-02 15:04")
		} else if rangeName == "7d" {
			bStart = now.Add(-time.Duration(i*3) * time.Hour).Truncate(3 * time.Hour)
			bEnd = bStart.Add(3 * time.Hour)
			label = bStart.Format("Jan 02 15:04")
			dateStr = bStart.Format("2006-01-02 15:04")
		} else {
			bStart = now.AddDate(0, 0, -i).Truncate(24 * time.Hour)
			bEnd = bStart.Add(24 * time.Hour)
			label = bStart.Format("Jan 02")
			if i == 0 {
				label = "Today"
			}
			dateStr = bStart.Format("2006-01-02")
		}

		// Non-penalization rule: If block ended before the project was created, strictly mark as no_data
		if bEnd.Before(createdAt) {
			blocks = append(blocks, ProjectAvailabilityBlock{
				Index:            blockCount - 1 - i,
				Label:            label,
				Date:             dateStr,
				Status:           "no_data",
				UptimePercent:    nil,
				SuccessfulChecks: 0,
				TotalChecks:      0,
				Details:          "No recorded telemetry",
			})
			continue
		}

		// Count checks falling in [bStart, bEnd)
		var bTotal, bSuccess int
		for _, c := range checks {
			if (c.checkedAt.Equal(bStart) || c.checkedAt.After(bStart)) && c.checkedAt.Before(bEnd) {
				bTotal++
				if c.isSuccessful {
					bSuccess++
				}
			}
		}

		if bTotal == 0 {
			blocks = append(blocks, ProjectAvailabilityBlock{
				Index:            blockCount - 1 - i,
				Label:            label,
				Date:             dateStr,
				Status:           "no_data",
				UptimePercent:    nil,
				SuccessfulChecks: 0,
				TotalChecks:      0,
				Details:          "No recorded telemetry",
			})
			continue
		}

		rangeTotal += bTotal
		rangeSuccess += bSuccess

		pct := math.Round((float64(bSuccess)/float64(bTotal))*1000) / 10
		status := "operational"
		details := fmt.Sprintf("100.0%% availability (%d/%d checks)", bSuccess, bTotal)
		if bSuccess < bTotal {
			if pct >= 50 {
				status = "partial"
				details = fmt.Sprintf("%.1f%% availability (%d/%d checks)", pct, bSuccess, bTotal)
			} else {
				status = "incident"
				details = fmt.Sprintf("%.1f%% availability (%d/%d checks)", pct, bSuccess, bTotal)
			}
		}

		blocks = append(blocks, ProjectAvailabilityBlock{
			Index:            blockCount - 1 - i,
			Label:            label,
			Date:             dateStr,
			Status:           status,
			UptimePercent:    &pct,
			SuccessfulChecks: bSuccess,
			TotalChecks:      bTotal,
			Details:          details,
		})
	}

	var rangeUptime *float64
	if rangeTotal > 0 {
		pct := math.Round((float64(rangeSuccess)/float64(rangeTotal))*1000) / 10
		rangeUptime = &pct
	}

	return ProjectRangeAvailability{
		RangeFilter:      rangeName,
		UptimePercent:    rangeUptime,
		TotalChecks:      rangeTotal,
		SuccessfulChecks: rangeSuccess,
		Blocks:           blocks,
	}
}

// GetProjectAvailabilityStats calculates multi-range uptime and latency for a specific project.
func (s *Service) GetProjectAvailabilityStats(ctx context.Context, projectID string) (ProjectAvailabilityStats, error) {
	var stats ProjectAvailabilityStats
	stats.Ranges = make(map[string]ProjectRangeAvailability)

	var projectCreatedAt time.Time

	// Query latest check, current status, availability, and project creation baseline
	err := s.pool.QueryRow(ctx, `
		SELECT p.status, p.availability, p.availability_reason, p.created_at, c.checked_at, c.response_time_ms
		FROM projects p
		LEFT JOIN LATERAL (
			SELECT checked_at, response_time_ms
			FROM project_availability_checks
			WHERE project_id = p.id
			ORDER BY checked_at DESC
			LIMIT 1
		) c ON true
		WHERE p.id = $1
	`, projectID).Scan(&stats.CurrentStatus, &stats.CurrentAvailability, &stats.AvailabilityReason, &projectCreatedAt, &stats.LastCheckedAt, &stats.LatestResponseTimeMs)
	if err != nil {
		return stats, err
	}

	// Query 30-day recorded checks
	rows, err := s.pool.Query(ctx, `
		SELECT checked_at, is_successful
		FROM project_availability_checks
		WHERE project_id = $1 AND checked_at >= NOW() - INTERVAL '30 days'
		ORDER BY checked_at ASC
	`, projectID)
	var checks []projectCheckEntry
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var entry projectCheckEntry
			if err := rows.Scan(&entry.checkedAt, &entry.isSuccessful); err == nil {
				checks = append(checks, entry)
			}
		}
	}

	now := time.Now()

	// Calculate discrete observation blocks: 1d (48 bars), 7d (56 bars), 30d (30 bars)
	r1d := buildProjectBuckets("1d", 48, 30*time.Minute, checks, projectCreatedAt, now)
	r7d := buildProjectBuckets("7d", 56, 3*time.Hour, checks, projectCreatedAt, now)
	r30d := buildProjectBuckets("30d", 30, 24*time.Hour, checks, projectCreatedAt, now)

	stats.Ranges["1d"] = r1d
	stats.Ranges["7d"] = r7d
	stats.Ranges["30d"] = r30d

	// Baseline values reflect 30-day window
	stats.Blocks = r30d.Blocks
	stats.TotalChecks30d = r30d.TotalChecks
	stats.UptimePercent = r30d.UptimePercent
	return stats, nil
}

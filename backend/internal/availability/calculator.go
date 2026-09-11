package availability

import (
	"context"
	"fmt"
	"math"
	"time"
)

// CalculateAvailability computes mathematical uptime strictly over the observed period.
// Days prior to first_monitored_at are marked as "no_data", never falsely counted as downtime.
func (s *Service) CalculateAvailability(ctx context.Context, rangeFilter string) (*AvailabilityResponse, error) {
	if rangeFilter != "1d" && rangeFilter != "7d" {
		rangeFilter = "30d"
	}

	now := time.Now()
	var (
		bootID           string
		startedAt        time.Time
		lastSeenAt       time.Time
		firstMonitoredAt time.Time
		currentStatusDB  string
	)

	err := s.pool.QueryRow(ctx, `
		SELECT boot_id, started_at, last_seen_at, first_monitored_at, current_status
		FROM availability_state
		WHERE id = 1
	`).Scan(&bootID, &startedAt, &lastSeenAt, &firstMonitoredAt, &currentStatusDB)

	if err != nil {
		// Fallback if state is not yet initialized
		firstMonitoredAt = now
		startedAt = now
		lastSeenAt = now
		currentStatusDB = "OPERATIONAL"
	}

	uptimeSec := uptimeSeconds()
	if uptimeSec <= 0 && !startedAt.IsZero() {
		uptimeSec = int64(now.Sub(startedAt).Seconds())
	}

	currentStatus := "operational"
	if time.Since(lastSeenAt) > 12*time.Minute {
		currentStatus = "degraded"
	}
	if time.Since(lastSeenAt) > 30*time.Minute {
		currentStatus = "down"
	}

	// Determine time window boundaries
	var windowStart time.Time
	var blockCount int
	var isHourly bool

	switch rangeFilter {
	case "1d":
		windowStart = now.Add(-24 * time.Hour)
		blockCount = 24
		isHourly = true
	case "7d":
		windowStart = now.AddDate(0, 0, -7)
		blockCount = 7
		isHourly = false
	default: // 30d
		windowStart = now.AddDate(0, 0, -30)
		blockCount = 30
		isHourly = false
	}

	// Query incidents that intersect with windowStart to now
	rows, err := s.pool.Query(ctx, `
		SELECT id, started_at, ended_at, duration_seconds, cause, details
		FROM availability_incidents
		WHERE ended_at >= $1
		ORDER BY started_at DESC
	`, windowStart)

	var incidents []IncidentRecord
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var inc IncidentRecord
			if err := rows.Scan(&inc.ID, &inc.StartedAt, &inc.EndedAt, &inc.DurationSeconds, &inc.Cause, &inc.Details); err == nil {
				inc.DurationFormatted = formatDuration(inc.DurationSeconds)
				incidents = append(incidents, inc)
			}
		}
	}

	// Calculate blocks
	blocks := make([]Block, 0, blockCount)

	if isHourly {
		// 24 hourly blocks (23 hours ago to current hour)
		for i := blockCount - 1; i >= 0; i-- {
			hStart := now.Add(-time.Duration(i) * time.Hour).Truncate(time.Hour)
			hEnd := hStart.Add(time.Hour)
			label := hStart.Format("15:04")
			dateStr := hStart.Format("2006-01-02 15:04")

			// Check if this hour is completely unmonitored
			if hEnd.Before(firstMonitoredAt) {
				blocks = append(blocks, Block{
					Index:         blockCount - 1 - i,
					Label:         label,
					Date:          dateStr,
					Status:        "no_data",
					UptimePercent: nil,
					Details:       "No recorded telemetry",
				})
				continue
			}

			// Active monitored hour
			mStart := hStart
			if mStart.Before(firstMonitoredAt) {
				mStart = firstMonitoredAt
			}
			mEnd := hEnd
			if mEnd.After(now) {
				mEnd = now
			}

			mSec := mEnd.Sub(mStart).Seconds()
			if mSec <= 0 {
				mSec = 3600
			}

			var downSec float64
			for _, inc := range incidents {
				oStart := inc.StartedAt
				if oStart.Before(mStart) {
					oStart = mStart
				}
				oEnd := inc.EndedAt
				if oEnd.After(mEnd) {
					oEnd = mEnd
				}
				if oEnd.After(oStart) {
					downSec += oEnd.Sub(oStart).Seconds()
				}
			}

			uptimePct := math.Round(((mSec-downSec)/mSec)*1000) / 10
			if uptimePct < 0 {
				uptimePct = 0
			}
			if uptimePct > 100 {
				uptimePct = 100
			}

			status := "operational"
			details := fmt.Sprintf("%.1f%% operational", uptimePct)
			if downSec > 0 {
				if downSec >= mSec {
					status = "incident"
					details = "Confirmed downtime incident"
				} else {
					status = "partial"
					details = fmt.Sprintf("Partial outage (%s down)", formatDuration(int64(downSec)))
				}
			}

			blocks = append(blocks, Block{
				Index:         blockCount - 1 - i,
				Label:         label,
				Date:          dateStr,
				Status:        status,
				UptimePercent: &uptimePct,
				Details:       details,
			})
		}
	} else {
		// Daily blocks (7 or 30)
		for i := blockCount - 1; i >= 0; i-- {
			dStart := now.AddDate(0, 0, -i).Truncate(24 * time.Hour)
			dEnd := dStart.Add(24 * time.Hour)
			label := dStart.Format("Jan 02")
			if i == 0 {
				label = "Today"
			}
			dateStr := dStart.Format("2006-01-02")

			// Check if this day is completely unmonitored
			if dEnd.Before(firstMonitoredAt) {
				blocks = append(blocks, Block{
					Index:         blockCount - 1 - i,
					Label:         label,
					Date:          dateStr,
					Status:        "no_data",
					UptimePercent: nil,
					Details:       "No recorded telemetry",
				})
				continue
			}

			// Active monitored day
			mStart := dStart
			if mStart.Before(firstMonitoredAt) {
				mStart = firstMonitoredAt
			}
			mEnd := dEnd
			if mEnd.After(now) {
				mEnd = now
			}

			mSec := mEnd.Sub(mStart).Seconds()
			if mSec <= 0 {
				mSec = 86400
			}

			var downSec float64
			for _, inc := range incidents {
				oStart := inc.StartedAt
				if oStart.Before(mStart) {
					oStart = mStart
				}
				oEnd := inc.EndedAt
				if oEnd.After(mEnd) {
					oEnd = mEnd
				}
				if oEnd.After(oStart) {
					downSec += oEnd.Sub(oStart).Seconds()
				}
			}

			uptimePct := math.Round(((mSec-downSec)/mSec)*1000) / 10
			if uptimePct < 0 {
				uptimePct = 0
			}
			if uptimePct > 100 {
				uptimePct = 100
			}

			status := "operational"
			details := fmt.Sprintf("%.1f%% operational", uptimePct)
			if downSec > 0 {
				if downSec >= mSec {
					status = "incident"
					details = "Confirmed downtime incident"
				} else {
					status = "partial"
					details = fmt.Sprintf("Partial outage (%s down)", formatDuration(int64(downSec)))
				}
			}

			blocks = append(blocks, Block{
				Index:         blockCount - 1 - i,
				Label:         label,
				Date:          dateStr,
				Status:        status,
				UptimePercent: &uptimePct,
				Details:       details,
			})
		}
	}

	// Mathematical availability strictly over the recorded observation window
	effectiveStart := windowStart
	if firstMonitoredAt.After(effectiveStart) {
		effectiveStart = firstMonitoredAt
	}

	totalRecordedSec := now.Sub(effectiveStart).Seconds()
	recordedDuration := now.Sub(firstMonitoredAt)

	var recordedPeriod string
	var recordedDays int = int(recordedDuration.Hours() / 24)
	if recordedDays >= 1 {
		recordedPeriod = fmt.Sprintf("%d days recorded", recordedDays+1)
	} else {
		hrs := int(recordedDuration.Hours())
		if hrs >= 1 {
			recordedPeriod = fmt.Sprintf("%d hours recorded", hrs)
		} else {
			recordedPeriod = fmt.Sprintf("%d mins recorded", int(recordedDuration.Minutes()))
		}
	}

	hasSufficientData := totalRecordedSec >= 1800 // at least 30 minutes of observations
	var availPct *float64
	var availFormatted string

	if totalRecordedSec <= 0 {
		availFormatted = "No data"
	} else {
		var windowDownSec float64
		for _, inc := range incidents {
			oStart := inc.StartedAt
			if oStart.Before(effectiveStart) {
				oStart = effectiveStart
			}
			oEnd := inc.EndedAt
			if oEnd.After(now) {
				oEnd = now
			}
			if oEnd.After(oStart) {
				windowDownSec += oEnd.Sub(oStart).Seconds()
			}
		}

		calculated := ((totalRecordedSec - windowDownSec) / totalRecordedSec) * 100
		calculated = math.Round(calculated*100) / 100
		if calculated < 0 {
			calculated = 0
		}
		if calculated > 100 {
			calculated = 100
		}
		availPct = &calculated

		if !hasSufficientData {
			availFormatted = "Collecting data"
		} else if calculated == 100 {
			availFormatted = "100%"
		} else {
			availFormatted = fmt.Sprintf("%.2f%%", calculated)
		}
	}

	var totalDownSec int64
	for _, inc := range incidents {
		totalDownSec += inc.DurationSeconds
	}
	downFormatted := "0m"
	if totalDownSec > 0 {
		downFormatted = formatDuration(totalDownSec)
	}

	return &AvailabilityResponse{
		CurrentStatus:          currentStatus,
		AvailabilityPercent:    availPct,
		AvailabilityFormatted:  availFormatted,
		RecordedPeriod:         recordedPeriod,
		RecordedDays:           recordedDays + 1,
		HasSufficientData:      hasSufficientData,
		UptimeSeconds:          uptimeSec,
		UptimeFormatted:        formatUptime(uptimeSec),
		LastHeartbeatAt:        lastSeenAt,
		LastHeartbeatText:      formatRelativeTime(lastSeenAt),
		FirstMonitoredAt:       firstMonitoredAt,
		RangeFilter:            rangeFilter,
		TotalIncidents:         len(incidents),
		TotalDowntimeSeconds:   totalDownSec,
		TotalDowntimeFormatted: downFormatted,
		Blocks:                 blocks,
		Incidents:              incidents,
	}, nil
}

package availability

import "time"

// Block represents a single discrete slice of telemetry (daily or hourly).
type Block struct {
	Index         int      `json:"index"`
	Label         string   `json:"label"`         // e.g. "Today", "Sep 08", "14:00"
	Date          string   `json:"date"`          // e.g. "2026-09-11" or "2026-09-11 14:00"
	Status        string   `json:"status"`        // "no_data" | "operational" | "partial" | "incident"
	UptimePercent *float64 `json:"uptimePercent"` // nil when status == "no_data"
	Details       string   `json:"details"`
}

// IncidentRecord documents a confirmed downtime gap or state transition.
type IncidentRecord struct {
	ID                int64     `json:"id"`
	StartedAt         time.Time `json:"started_at"`
	EndedAt           time.Time `json:"ended_at"`
	DurationSeconds   int64     `json:"duration_seconds"`
	DurationFormatted string    `json:"duration_formatted"`
	Cause             string    `json:"cause"`
	Details           string    `json:"details"`
}

// AvailabilityResponse encapsulates the full, observation-bounded availability contract.
type AvailabilityResponse struct {
	CurrentStatus          string           `json:"current_status"`      // "operational" | "degraded" | "down"
	AvailabilityPercent    *float64        `json:"availability_percent"` // nil if no recorded data exists
	AvailabilityFormatted  string           `json:"availability_formatted"`
	RecordedPeriod         string           `json:"recorded_period"`     // e.g. "3 days recorded"
	RecordedDays           int              `json:"recorded_days"`
	HasSufficientData      bool             `json:"has_sufficient_data"`
	UptimeSeconds          int64            `json:"uptime_seconds"`
	UptimeFormatted        string           `json:"uptime_formatted"`
	LastHeartbeatAt        time.Time        `json:"last_heartbeat_at"`
	LastHeartbeatText      string           `json:"last_heartbeat_text"`
	FirstMonitoredAt       time.Time        `json:"first_monitored_at"`
	RangeFilter            string           `json:"range_filter"`        // "1d" | "7d" | "30d"
	TotalIncidents         int              `json:"total_incidents"`
	TotalDowntimeSeconds   int64            `json:"total_downtime_seconds"`
	TotalDowntimeFormatted string           `json:"total_downtime_formatted"`
	Blocks                 []Block          `json:"blocks"`
	Incidents              []IncidentRecord `json:"incidents"`
}

// ServiceHealth describes an application-level service's live probed state.
type ServiceHealth struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"` // "OPERATIONAL" | "DEGRADED" | "DOWN"
	LastCheckedAt  time.Time `json:"last_checked_at"`
	ResponseTimeMs int       `json:"response_time_ms"`
}

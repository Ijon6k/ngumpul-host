package visits

// DailyMetric represents page views and outbound visits for a single day.
type DailyMetric struct {
	Date           string `json:"date"`
	Views          int    `json:"views"`
	OutboundVisits int    `json:"outbound_visits"`
}

// ProjectVisitsResponse provides honest 30-day metrics for project owners.
type ProjectVisitsResponse struct {
	TotalViews30d    int           `json:"total_views_30d"`
	TotalOutbound30d int           `json:"total_outbound_30d"`
	DailyBreakdown   []DailyMetric `json:"daily_breakdown"`
}

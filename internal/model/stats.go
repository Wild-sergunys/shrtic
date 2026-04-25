package model

type StatItem struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
	Count      int     `json:"count"`
}

type LinkStats struct {
	TotalClicks int        `json:"total_clicks"`
	Browsers    []StatItem `json:"browsers"`
	Devices     []StatItem `json:"devices"`
	Countries   []StatItem `json:"countries"`
	Referrers   []StatItem `json:"referrers"`
}

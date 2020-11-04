package api

import "time"

type SummaryQueryResult struct {
	Page                  Page                   `json:"page"`
	Components            []Component            `json:"components"`
	Incidents             []Incident             `json:"incidents"`
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
	Status                Status                 `json:"status"`
}
type IncidentsQueryResult struct {
	Page      Page       `json:"page"`
	Incidents []Incident `json:"incidents"`
}
type ScheduledMaintenancesQueryResult struct {
	Page                  Page                   `json:"page"`
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
}
type ComponentsQueryResult struct {
	Page       Page        `json:"page"`
	Components []Component `json:"components"`
}
type StatusQueryResult struct {
	Page   Page   `json:"page"`
	Status Status `json:"status"`
}
type AffectedComponents struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
}
type ScheduledMaintenance struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Status          string           `json:"status"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	MonitoringAt    interface{}      `json:"monitoring_at"`
	ResolvedAt      time.Time        `json:"resolved_at"`
	Impact          string           `json:"impact"`
	Shortlink       string           `json:"shortlink"`
	StartedAt       time.Time        `json:"started_at"`
	PageID          string           `json:"page_id"`
	IncidentUpdates []IncidentUpdate `json:"incident_updates"`
	Components      []Component      `json:"components"`
	ScheduledFor    time.Time        `json:"scheduled_for"`
	ScheduledUntil  time.Time        `json:"scheduled_until"`
}
type Page struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	TimeZone  string    `json:"time_zone"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Component struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Status             string      `json:"status"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	Position           int         `json:"position"`
	Description        string      `json:"description"`
	Showcase           bool        `json:"showcase"`
	StartDate          interface{} `json:"start_date"`
	GroupID            interface{} `json:"group_id"`
	PageID             string      `json:"page_id"`
	Group              bool        `json:"group"`
	OnlyShowIfDegraded bool        `json:"only_show_if_degraded"`
}
type Status struct {
	Indicator   string `json:"indicator"`
	Description string `json:"description"`
}
type IncidentUpdate struct {
	ID                   string      `json:"id"`
	Status               string      `json:"status"`
	Body                 string      `json:"body"`
	IncidentID           string      `json:"incident_id"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
	DisplayAt            time.Time   `json:"display_at"`
	AffectedComponents   interface{} `json:"affected_components"`
	DeliverNotifications bool        `json:"deliver_notifications"`
	CustomTweet          interface{} `json:"custom_tweet"`
	TweetID              interface{} `json:"tweet_id"`
}
type Incident struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Status          string           `json:"status"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	MonitoringAt    interface{}      `json:"monitoring_at"`
	ResolvedAt      time.Time        `json:"resolved_at"`
	Impact          string           `json:"impact"`
	Shortlink       string           `json:"shortlink"`
	StartedAt       time.Time        `json:"started_at"`
	PageID          string           `json:"page_id"`
	IncidentUpdates []IncidentUpdate `json:"incident_updates"`
	Components      []Component      `json:"components"`
}

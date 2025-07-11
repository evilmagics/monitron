package models

import (
	"time"

	"github.com/google/uuid"
)

type OperationalPage struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Slug        string    `db:"slug" json:"slug"` // Unique slug for the page URL
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	IsPublic    bool      `db:"is_public" json:"is_public"` // True if public, false if private (requires auth)
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (p *OperationalPage) TableName() string {
	return "operational_pages"
}

type OperationalPageComponent struct {
	ID            uuid.UUID `db:"id" json:"id"`
	PageID        uuid.UUID `db:"page_id" json:"page_id"`
	ComponentType string    `db:"component_type" json:"component_type"` // "service" or "domain_ssl"
	ComponentID   uuid.UUID `db:"component_id" json:"component_id"`
	ComponentName string    `db:"component_name" json:"component_name"` // User-defined name for display
	DisplayOrder  int       `db:"display_order" json:"display_order"`
	Description   string    `db:"description" json:"description"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (c *OperationalPageComponent) TableName() string {
	return "operational_page_components"
}

// OperationalPageStats represents aggregated stats for an operational page
type OperationalPageStats struct {
	PageID              uuid.UUID `db:"page_id" json:"page_id"`
	OverallUptime       float64   `db:"overall_uptime" json:"overall_uptime"`
	IncidentsTotal      int       `db:"incidents_total" json:"incidents_total"`
	AverageResponseTime float64   `db:"average_response_time" json:"average_response_time"`
	UptimeHistory       string    `db:"uptime_history" json:"uptime_history"` // JSON string of 30-day history
	LastUpdated         time.Time `db:"last_updated" json:"last_updated"`
}

func (s *OperationalPageStats) TableName() string {
	return "operational_page_stats"
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type DomainSSL struct {
	ID               uuid.UUID `db:"id" json:"id"`
	Domain           string    `db:"domain" json:"domain"`
	WarningThreshold int       `db:"warning_threshold" json:"warning_threshold"` // Days before expiry to warn
	ExpiryThreshold  int       `db:"expiry_threshold" json:"expiry_threshold"`   // Days before expiry to consider expired
	CheckInterval    int       `db:"check_interval" json:"check_interval"`       // in days
	Label            string    `db:"label" json:"label"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`

	// Parsed SSL/Domain details
	CertificateDetail string    `db:"certificate_detail" json:"certificate_detail"` // JSON string
	Issuer            string    `db:"issuer" json:"issuer"`
	ValidFrom         time.Time `db:"valid_from" json:"valid_from"`
	ResolvedIP        string    `db:"resolved_ip" json:"resolved_ip"`
	Expiry            time.Time `db:"expiry" json:"expiry"`
	DaysLeft          int       `db:"days_left" json:"days_left"` // Calculated, not stored
}

func (DomainSSL) TableName() string {
	return "domain_ssl"
}

// DomainSSLStats represents the monitoring statistics for a domain/SSL
type DomainSSLStats struct {
	DomainSSLID         uuid.UUID `db:"domain_ssl_id" json:"domain_ssl_id"`
	ResponseTime        float64   `db:"response_time" json:"response_time"`
	Uptime              float64   `db:"uptime" json:"uptime"` // Percentage
	LastChecked         time.Time `db:"last_checked" json:"last_checked"`
	AverageResponseTime float64   `db:"average_response_time" json:"average_response_time"`
	IncidentTotal       int       `db:"incident_total" json:"incident_total"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}

func (DomainSSLStats) TableName() string {
	return "domain_ssl_stats"
}

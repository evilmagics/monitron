package models

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	ReportType string    `db:"report_type" json:"report_type"` // e.g., "instance_summary", "service_uptime"
	Format    string    `db:"format" json:"format"`     // e.g., "CSV", "PDF", "Excel"
	Status    string    `db:"status" json:"status"`     // e.g., "pending", "generating", "completed", "failed"
	GeneratedAt time.Time `db:"generated_at" json:"generated_at"`
	FilePath  string    `db:"file_path" json:"file_path"` // Path to the generated report file
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type LogEntry struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Level     string    `db:"level" json:"level"`     // e.g., "info", "warn", "error"
	Message   string    `db:"message" json:"message"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Service   string    `db:"service" json:"service"` // e.g., "monitron-server", "monitron-agent"
	RequestID string    `db:"request_id" json:"request_id"` // Optional: for tracing requests
}



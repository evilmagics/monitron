package models

import (
	"time"

	"github.com/google/uuid"
)

type Instance struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Name          string    `db:"name" json:"name" validate:"required"`
	Host          string    `db:"host" json:"host" validate:"required"`
	CheckInterval int       `db:"check_interval" json:"check_interval" validate:"required,min=1"`
	CheckTimeout  int       `db:"check_timeout" json:"check_timeout" validate:"required,min=1"`
	AgentPort     int       `db:"agent_port" json:"agent_port" validate:"required,min=1"`
	AgentAuth     string    `db:"agent_auth" json:"agent_auth"` // Encrypted agent authentication string
	Description   string    `db:"description" json:"description"`
	Label         string    `db:"label" json:"label"`
	Group         string    `db:"group" json:"group"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type InstanceStats struct {
	ID           uuid.UUID `db:"id" json:"id"`
	InstanceID   uuid.UUID `db:"instance_id" json:"instance_id"`
	CPUUsage     float64   `db:"cpu_usage" json:"cpu_usage"`
	MemoryUsage  float64   `db:"memory_usage" json:"memory_usage"`
	DiskUsage    float64   `db:"disk_usage" json:"disk_usage"`
	NetworkIOIn  float64   `db:"network_io_in" json:"network_io_in"`
	NetworkIOOut float64   `db:"network_io_out" json:"network_io_out"`
	Timestamp    time.Time `db:"timestamp" json:"timestamp"`
}

type InstanceMetric struct {
	ID          uuid.UUID `db:"id" json:"id"`
	InstanceID  uuid.UUID `db:"instance_id" json:"instance_id"`
	MetricName  string    `db:"metric_name" json:"metric_name"`
	MetricValue float64   `db:"metric_value" json:"metric_value"`
	Timestamp   time.Time `db:"timestamp" json:"timestamp"`
}

type DeviceInfo struct {
	ID           uuid.UUID `db:"id" json:"id"`
	InstanceID   uuid.UUID `db:"instance_id" json:"instance_id"`
	OS           string    `db:"os" json:"os"`
	Architecture string    `db:"architecture" json:"architecture"`
	Hostname     string    `db:"hostname" json:"hostname"`
	TotalCPU     int       `db:"total_cpu" json:"total_cpu"`
	TotalMemory  int       `db:"total_memory" json:"total_memory"`
	BootTime     time.Time `db:"boot_time" json:"boot_time"`
	Timestamp    time.Time `db:"timestamp" json:"timestamp"`
}

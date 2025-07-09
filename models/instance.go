package models

import (
	"time"

	"gorm.io/gorm"
)

type Instance struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	IPAddress   string         `json:"ip_address" gorm:"not null"`
	Port        int            `json:"port"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Metrics []InstanceMetric `json:"metrics,omitempty" gorm:"foreignKey:InstanceID"`
}

type InstanceMetric struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	InstanceID      uint      `json:"instance_id" gorm:"not null;index"`
	Timestamp       time.Time `json:"timestamp" gorm:"not null;index"`
	CPUUsage        *float64  `json:"cpu_usage"`
	MemoryUsage     *float64  `json:"memory_usage"`
	DiskUsage       *float64  `json:"disk_usage"`
	NetworkIOIn     *int64    `json:"network_io_in"`
	NetworkIOOut    *int64    `json:"network_io_out"`
	Uptime          *int64    `json:"uptime"`
	ProcessCount    *int      `json:"process_count"`
	LoadAverage1m   *float64  `json:"load_average_1m"`
	LoadAverage5m   *float64  `json:"load_average_5m"`
	LoadAverage15m  *float64  `json:"load_average_15m"`
	
	// Relationships
	Instance Instance `json:"instance,omitempty" gorm:"foreignKey:InstanceID"`
}


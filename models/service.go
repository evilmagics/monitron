package models

import (
	"time"

	"gorm.io/gorm"
)

type ServiceType string

const (
	ServiceTypeHTTP ServiceType = "http"
	ServiceTypeGRPC ServiceType = "grpc"
	ServiceTypeTCP  ServiceType = "tcp"
	ServiceTypeMQTT ServiceType = "mqtt"
)

type Service struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null"`
	Type            ServiceType    `json:"type" gorm:"not null"`
	Host            string         `json:"host" gorm:"not null"`
	Port            int            `json:"port"`
	IntervalSeconds int            `json:"interval_seconds" gorm:"default:60"`
	Description     string         `json:"description"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	HTTPConfig *ServiceHTTPConfig `json:"http_config,omitempty" gorm:"foreignKey:ServiceID"`
	GRPCConfig *ServiceGRPCConfig `json:"grpc_config,omitempty" gorm:"foreignKey:ServiceID"`
	TCPConfig  *ServiceTCPConfig  `json:"tcp_config,omitempty" gorm:"foreignKey:ServiceID"`
	MQTTConfig *ServiceMQTTConfig `json:"mqtt_config,omitempty" gorm:"foreignKey:ServiceID"`
	Checks     []ServiceCheck     `json:"checks,omitempty" gorm:"foreignKey:ServiceID"`
}

type ServiceHTTPConfig struct {
	ServiceID      uint   `json:"service_id" gorm:"primaryKey"`
	Path           string `json:"path" gorm:"default:/"`
	Method         string `json:"method" gorm:"default:GET"`
	ExpectedStatus int    `json:"expected_status" gorm:"default:200"`
	RegexMatch     string `json:"regex_match"`
	SSLCheck       bool   `json:"ssl_check" gorm:"default:true"`
	
	// Relationships
	Service Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

type ServiceGRPCConfig struct {
	ServiceID              uint   `json:"service_id" gorm:"primaryKey"`
	ServiceName            string `json:"service_name" gorm:"not null"`
	MethodName             string `json:"method_name" gorm:"not null"`
	RequestBody            string `json:"request_body"`
	ExpectedResponseRegex  string `json:"expected_response_regex"`
	
	// Relationships
	Service Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

type ServiceTCPConfig struct {
	ServiceID uint `json:"service_id" gorm:"primaryKey"`
	TimeoutMs int  `json:"timeout_ms" gorm:"default:5000"`
	
	// Relationships
	Service Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

type ServiceMQTTConfig struct {
	ServiceID uint   `json:"service_id" gorm:"primaryKey"`
	Topic     string `json:"topic" gorm:"not null"`
	QoS       int    `json:"qos" gorm:"default:0"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	
	// Relationships
	Service Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

type ServiceCheckStatus string

const (
	ServiceCheckStatusUp      ServiceCheckStatus = "up"
	ServiceCheckStatusDown    ServiceCheckStatus = "down"
	ServiceCheckStatusUnknown ServiceCheckStatus = "unknown"
)

type ServiceCheck struct {
	ID             uint               `json:"id" gorm:"primaryKey"`
	ServiceID      uint               `json:"service_id" gorm:"not null;index"`
	Timestamp      time.Time          `json:"timestamp" gorm:"not null;index"`
	Status         ServiceCheckStatus `json:"status" gorm:"not null"`
	ResponseTimeMs *int               `json:"response_time_ms"`
	ErrorMessage   string             `json:"error_message"`
	
	// Relationships
	Service Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}


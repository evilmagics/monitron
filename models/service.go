package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	APIType       string    `db:"api_type" json:"api_type"`             // HTTP API, gRPC, MQTT, TCP, DNS, Ping
	CheckInterval int       `db:"check_interval" json:"check_interval"` // in seconds
	Timeout       int       `db:"timeout" json:"timeout"`               // in seconds
	Description   string    `db:"description" json:"description"`
	Label         string    `db:"label" json:"label"`
	Group         string    `db:"group" json:"group"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`

	// Specific fields for HTTP API
	HTTPMethod         string `db:"http_method" json:"http_method"`
	HTTPHealthURL      string `db:"http_health_url" json:"http_health_url"`
	HTTPExpectedStatus int    `db:"http_expected_status" json:"http_expected_status"`

	// Specific fields for gRPC
	GRPCHost  string `db:"grpc_host" json:"grpc_host"`
	GRPCPort  int    `db:"grpc_port" json:"grpc_port"`
	GRPCAuth  string `db:"grpc_auth" json:"grpc_auth"`
	GRPCProto string `db:"grpc_proto" json:"grpc_proto"`

	// Specific fields for MQTT
	MQTTHost  string `db:"mqtt_host" json:"mqtt_host"`
	MQTTPort  int    `db:"mqtt_port" json:"mqtt_port"`
	MQTTQoS   int    `db:"mqtt_qos" json:"mqtt_qos"`
	MQTTTopic string `db:"mqtt_topic" json:"mqtt_topic"`
	MQTTAuth  string `db:"mqtt_auth" json:"mqtt_auth"`

	// Specific fields for TCP
	TCPHost string `db:"tcp_host" json:"tcp_host"`
	TCPPort int    `db:"tcp_port" json:"tcp_port"`

	// Specific fields for DNS
	DNSDomainName string `db:"dns_domain_name" json:"dns_domain_name"`

	// Specific fields for Ping
	PingHost string `db:"ping_host" json:"ping_host"`
}

func (Service) TableName() string {
	return "services"
}

// ServiceStats represents the monitoring statistics for a service
type ServiceStats struct {
	ServiceID           uuid.UUID `db:"service_id" json:"service_id"`
	ResponseTime        float64   `db:"response_time" json:"response_time"`
	Uptime              float64   `db:"uptime" json:"uptime"` // Percentage
	LastChecked         time.Time `db:"last_checked" json:"last_checked"`
	AverageResponseTime float64   `db:"average_response_time" json:"average_response_time"`
	IncidentTotal       int       `db:"incident_total" json:"incident_total"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}

func (ServiceStats) TableName() string {
	return "service_stats"
}

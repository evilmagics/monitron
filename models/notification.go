package models

import (
	"time"
)

type NotificationType string

const (
	NotificationTypeEmail    NotificationType = "email"
	NotificationTypeSlack    NotificationType = "slack"
	NotificationTypeTelegram NotificationType = "telegram"
	NotificationTypeWebhook  NotificationType = "webhook"
)

type Notification struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	UserID    uint             `json:"user_id" gorm:"not null"`
	Type      NotificationType `json:"type" gorm:"not null"`
	Config    string           `json:"config" gorm:"type:json"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	
	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type AlertSourceType string

const (
	AlertSourceTypeInstance AlertSourceType = "instance"
	AlertSourceTypeService  AlertSourceType = "service"
	AlertSourceTypeDNS      AlertSourceType = "dns"
)

type AlertStatus string

const (
	AlertStatusTriggered AlertStatus = "triggered"
	AlertStatusResolved  AlertStatus = "resolved"
)

type Alert struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	SourceType   AlertSourceType `json:"source_type" gorm:"not null"`
	SourceID     uint            `json:"source_id" gorm:"not null"`
	MetricName   string          `json:"metric_name"`
	Threshold    string          `json:"threshold"`
	CurrentValue string          `json:"current_value"`
	Status       AlertStatus     `json:"status" gorm:"not null"`
	TriggeredAt  time.Time       `json:"triggered_at" gorm:"not null"`
	ResolvedAt   *time.Time      `json:"resolved_at"`
	Message      string          `json:"message"`
}


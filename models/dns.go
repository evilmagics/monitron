package models

import (
	"time"

	"gorm.io/gorm"
)

type DNSRecordType string

const (
	DNSRecordTypeA     DNSRecordType = "A"
	DNSRecordTypeAAAA  DNSRecordType = "AAAA"
	DNSRecordTypeCNAME DNSRecordType = "CNAME"
	DNSRecordTypeMX    DNSRecordType = "MX"
	DNSRecordTypeNS    DNSRecordType = "NS"
	DNSRecordTypeTXT   DNSRecordType = "TXT"
	DNSRecordTypeSRV   DNSRecordType = "SRV"
	DNSRecordTypePTR   DNSRecordType = "PTR"
)

type DNSRecord struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null"`
	RecordType      DNSRecordType  `json:"record_type" gorm:"not null"`
	Target          string         `json:"target" gorm:"not null"`
	ExpectedIP      string         `json:"expected_ip"`
	DNSServer       string         `json:"dns_server"`
	IntervalSeconds int            `json:"interval_seconds" gorm:"default:300"`
	Description     string         `json:"description"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Checks []DNSCheck `json:"checks,omitempty" gorm:"foreignKey:DNSRecordID"`
}

type DNSCheckStatus string

const (
	DNSCheckStatusResolved DNSCheckStatus = "resolved"
	DNSCheckStatusFailed   DNSCheckStatus = "failed"
)

type DNSCheck struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	DNSRecordID      uint           `json:"dns_record_id" gorm:"not null;index"`
	Timestamp        time.Time      `json:"timestamp" gorm:"not null;index"`
	Status           DNSCheckStatus `json:"status" gorm:"not null"`
	ResolvedIP       string         `json:"resolved_ip"`
	ResolutionTimeMs *int           `json:"resolution_time_ms"`
	ErrorMessage     string         `json:"error_message"`
	
	// Relationships
	DNSRecord DNSRecord `json:"dns_record,omitempty" gorm:"foreignKey:DNSRecordID"`
}


package services

import (
	"monitoring-backend/models"
	"time"

	"gorm.io/gorm"
)

type DNSService struct {
	db *gorm.DB
}

func NewDNSService(db *gorm.DB) *DNSService {
	return &DNSService{db: db}
}

func (s *DNSService) GetDNSRecords() ([]models.DNSRecord, error) {
	var records []models.DNSRecord
	err := s.db.Find(&records).Error
	return records, err
}

func (s *DNSService) GetDNSRecordByID(id uint) (*models.DNSRecord, error) {
	var record models.DNSRecord
	err := s.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *DNSService) CreateDNSRecord(record *models.DNSRecord) error {
	return s.db.Create(record).Error
}

func (s *DNSService) UpdateDNSRecord(id uint, record *models.DNSRecord) error {
	return s.db.Model(&models.DNSRecord{}).Where("id = ?", id).Updates(record).Error
}

func (s *DNSService) DeleteDNSRecord(id uint) error {
	return s.db.Delete(&models.DNSRecord{}, id).Error
}

func (s *DNSService) SaveCheck(recordID uint, check *models.DNSCheck) error {
	check.DNSRecordID = recordID
	check.Timestamp = time.Now()
	return s.db.Create(check).Error
}

func (s *DNSService) GetDNSChecks(recordID uint, from, to time.Time) ([]models.DNSCheck, error) {
	var checks []models.DNSCheck
	query := s.db.Where("dns_record_id = ?", recordID)
	
	if !from.IsZero() {
		query = query.Where("timestamp >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("timestamp <= ?", to)
	}
	
	err := query.Order("timestamp DESC").Find(&checks).Error
	return checks, err
}

func (s *DNSService) GetLatestCheck(recordID uint) (*models.DNSCheck, error) {
	var check models.DNSCheck
	err := s.db.Where("dns_record_id = ?", recordID).Order("timestamp DESC").First(&check).Error
	if err != nil {
		return nil, err
	}
	return &check, nil
}


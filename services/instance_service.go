package services

import (
	"monitoring-backend/models"
	"time"

	"gorm.io/gorm"
)

type InstanceService struct {
	db *gorm.DB
}

func NewInstanceService(db *gorm.DB) *InstanceService {
	return &InstanceService{db: db}
}

func (s *InstanceService) GetInstances() ([]models.Instance, error) {
	var instances []models.Instance
	err := s.db.Find(&instances).Error
	return instances, err
}

func (s *InstanceService) GetInstanceByID(id uint) (*models.Instance, error) {
	var instance models.Instance
	err := s.db.First(&instance, id).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (s *InstanceService) CreateInstance(instance *models.Instance) error {
	return s.db.Create(instance).Error
}

func (s *InstanceService) UpdateInstance(id uint, instance *models.Instance) error {
	return s.db.Model(&models.Instance{}).Where("id = ?", id).Updates(instance).Error
}

func (s *InstanceService) DeleteInstance(id uint) error {
	return s.db.Delete(&models.Instance{}, id).Error
}

func (s *InstanceService) SaveMetrics(instanceID uint, metrics *models.InstanceMetric) error {
	metrics.InstanceID = instanceID
	metrics.Timestamp = time.Now()
	return s.db.Create(metrics).Error
}

func (s *InstanceService) GetInstanceMetrics(instanceID uint, from, to time.Time) ([]models.InstanceMetric, error) {
	var metrics []models.InstanceMetric
	query := s.db.Where("instance_id = ?", instanceID)
	
	if !from.IsZero() {
		query = query.Where("timestamp >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("timestamp <= ?", to)
	}
	
	err := query.Order("timestamp DESC").Find(&metrics).Error
	return metrics, err
}

func (s *InstanceService) GetLatestMetrics(instanceID uint) (*models.InstanceMetric, error) {
	var metric models.InstanceMetric
	err := s.db.Where("instance_id = ?", instanceID).Order("timestamp DESC").First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}


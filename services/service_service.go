package services

import (
	"monitoring-backend/models"
	"time"

	"gorm.io/gorm"
)

type ServiceService struct {
	db *gorm.DB
}

func NewServiceService(db *gorm.DB) *ServiceService {
	return &ServiceService{db: db}
}

func (s *ServiceService) GetServices() ([]models.Service, error) {
	var services []models.Service
	err := s.db.Preload("HTTPConfig").Preload("GRPCConfig").Preload("TCPConfig").Preload("MQTTConfig").Find(&services).Error
	return services, err
}

func (s *ServiceService) GetServiceByID(id uint) (*models.Service, error) {
	var service models.Service
	err := s.db.Preload("HTTPConfig").Preload("GRPCConfig").Preload("TCPConfig").Preload("MQTTConfig").First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *ServiceService) CreateService(service *models.Service) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(service).Error; err != nil {
			return err
		}

		// Create type-specific config
		switch service.Type {
		case models.ServiceTypeHTTP:
			if service.HTTPConfig != nil {
				service.HTTPConfig.ServiceID = service.ID
				return tx.Create(service.HTTPConfig).Error
			}
		case models.ServiceTypeGRPC:
			if service.GRPCConfig != nil {
				service.GRPCConfig.ServiceID = service.ID
				return tx.Create(service.GRPCConfig).Error
			}
		case models.ServiceTypeTCP:
			if service.TCPConfig != nil {
				service.TCPConfig.ServiceID = service.ID
				return tx.Create(service.TCPConfig).Error
			}
		case models.ServiceTypeMQTT:
			if service.MQTTConfig != nil {
				service.MQTTConfig.ServiceID = service.ID
				return tx.Create(service.MQTTConfig).Error
			}
		}
		return nil
	})
}

func (s *ServiceService) UpdateService(id uint, service *models.Service) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Service{}).Where("id = ?", id).Updates(service).Error; err != nil {
			return err
		}

		// Update type-specific config
		switch service.Type {
		case models.ServiceTypeHTTP:
			if service.HTTPConfig != nil {
				return tx.Model(&models.ServiceHTTPConfig{}).Where("service_id = ?", id).Updates(service.HTTPConfig).Error
			}
		case models.ServiceTypeGRPC:
			if service.GRPCConfig != nil {
				return tx.Model(&models.ServiceGRPCConfig{}).Where("service_id = ?", id).Updates(service.GRPCConfig).Error
			}
		case models.ServiceTypeTCP:
			if service.TCPConfig != nil {
				return tx.Model(&models.ServiceTCPConfig{}).Where("service_id = ?", id).Updates(service.TCPConfig).Error
			}
		case models.ServiceTypeMQTT:
			if service.MQTTConfig != nil {
				return tx.Model(&models.ServiceMQTTConfig{}).Where("service_id = ?", id).Updates(service.MQTTConfig).Error
			}
		}
		return nil
	})
}

func (s *ServiceService) DeleteService(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete type-specific configs
		tx.Delete(&models.ServiceHTTPConfig{}, "service_id = ?", id)
		tx.Delete(&models.ServiceGRPCConfig{}, "service_id = ?", id)
		tx.Delete(&models.ServiceTCPConfig{}, "service_id = ?", id)
		tx.Delete(&models.ServiceMQTTConfig{}, "service_id = ?", id)
		
		// Delete service
		return tx.Delete(&models.Service{}, id).Error
	})
}

func (s *ServiceService) SaveCheck(serviceID uint, check *models.ServiceCheck) error {
	check.ServiceID = serviceID
	check.Timestamp = time.Now()
	return s.db.Create(check).Error
}

func (s *ServiceService) GetServiceChecks(serviceID uint, from, to time.Time) ([]models.ServiceCheck, error) {
	var checks []models.ServiceCheck
	query := s.db.Where("service_id = ?", serviceID)
	
	if !from.IsZero() {
		query = query.Where("timestamp >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("timestamp <= ?", to)
	}
	
	err := query.Order("timestamp DESC").Find(&checks).Error
	return checks, err
}

func (s *ServiceService) GetLatestCheck(serviceID uint) (*models.ServiceCheck, error) {
	var check models.ServiceCheck
	err := s.db.Where("service_id = ?", serviceID).Order("timestamp DESC").First(&check).Error
	if err != nil {
		return nil, err
	}
	return &check, nil
}


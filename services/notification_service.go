package services

import (
	"encoding/json"
	"fmt"
	"monitoring-backend/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) SendAlert(alert *models.Alert) error {
	// Save alert to database
	if err := s.db.Create(alert).Error; err != nil {
		return err
	}

	// Get all notifications
	var notifications []models.Notification
	if err := s.db.Find(&notifications).Error; err != nil {
		return err
	}

	// Send notifications
	for _, notification := range notifications {
		go s.sendNotification(notification, alert)
	}

	return nil
}

func (s *NotificationService) sendNotification(notification models.Notification, alert *models.Alert) {
	message := fmt.Sprintf("Alert: %s - %s", alert.MetricName, alert.Message)
	
	switch notification.Type {
	case models.NotificationTypeEmail:
		s.sendEmailNotification(notification.Config, message)
	case models.NotificationTypeSlack:
		s.sendSlackNotification(notification.Config, message)
	case models.NotificationTypeTelegram:
		s.sendTelegramNotification(notification.Config, message)
	case models.NotificationTypeWebhook:
		s.sendWebhookNotification(notification.Config, message)
	}
}

func (s *NotificationService) sendEmailNotification(config, message string) {
	// TODO: Implement email notification
	fmt.Printf("Email notification: %s\n", message)
}

func (s *NotificationService) sendSlackNotification(config, message string) {
	// TODO: Implement Slack notification
	fmt.Printf("Slack notification: %s\n", message)
}

func (s *NotificationService) sendTelegramNotification(config, message string) {
	// TODO: Implement Telegram notification
	fmt.Printf("Telegram notification: %s\n", message)
}

func (s *NotificationService) sendWebhookNotification(config, message string) {
	// TODO: Implement webhook notification
	fmt.Printf("Webhook notification: %s\n", message)
}

func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	// Validate config is valid JSON
	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(notification.Config), &configMap); err != nil {
		return fmt.Errorf("invalid config JSON: %v", err)
	}
	
	return s.db.Create(notification).Error
}

func (s *NotificationService) GetNotifications(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := s.db.Where("user_id = ?", userID).Find(&notifications).Error
	return notifications, err
}

func (s *NotificationService) DeleteNotification(id uint) error {
	return s.db.Delete(&models.Notification{}, id).Error
}


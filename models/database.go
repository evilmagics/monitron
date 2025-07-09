package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Instance{},
		&InstanceMetric{},
		&Service{},
		&ServiceHTTPConfig{},
		&ServiceGRPCConfig{},
		&ServiceTCPConfig{},
		&ServiceMQTTConfig{},
		&ServiceCheck{},
		&DNSRecord{},
		&DNSCheck{},
		&Notification{},
		&Alert{},
	)
}


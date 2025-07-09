package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string         `json:"-" gorm:"not null;size:255"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:255"`
	Role         string         `json:"role" gorm:"default:user;size:50"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}


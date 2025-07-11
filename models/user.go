package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"password"` // Hashed password
	Role      string     `db:"role" json:"role"`         // e.g., "admin", "user"
	Status    string     `db:"status" json:"status"`     // e.g., "active", "inactive"
	LastLogin *time.Time `db:"last_login" json:"last_login"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// PasswordResetToken represents a token for password reset
type PasswordResetToken struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

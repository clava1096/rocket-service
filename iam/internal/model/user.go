package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID            `json:"id" db:"id"`
	Login               string               `json:"login" db:"login"`
	Email               string               `json:"email" db:"email"`
	PasswordHash        string               `json:"-" db:"password_hash"`
	NotificationMethods []NotificationMethod `json:"notification_methods" db:"notification_methods"`
	CreatedAt           time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at" db:"updated_at"`
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type RegisterRequest struct {
	Login               string
	Password            string
	Email               string
	NotificationMethods []NotificationMethod
}

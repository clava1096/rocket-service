package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID                uuid.UUID
	Login               string
	Email               string
	PasswordHash        string
	NotificationMethods []NotificationMethod
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}

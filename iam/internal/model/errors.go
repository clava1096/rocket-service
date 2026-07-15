package model

import "errors"

var (
	ErrUserEmailExists      = errors.New("this email already exists")
	ErrUserLoginExists      = errors.New("this login already exists")
	ErrUserAlreadyExists    = errors.New("this user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidSessionFormat = errors.New("invalid session format")
	ErrSessionNotFound      = errors.New("session not found")
)

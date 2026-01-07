package model

import "errors"

var (
	ErrNotFound    = errors.New("part not found")
	ErrWhileCreate = errors.New("error while creating model in mongodb")
)

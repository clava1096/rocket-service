package model

import "errors"

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrOrderNotPending = errors.New("order is not in pending state")
	ErrPartNotFound    = errors.New("one or more parts not found")
	ErrThisOrderExists = errors.New("this order already exists")
)

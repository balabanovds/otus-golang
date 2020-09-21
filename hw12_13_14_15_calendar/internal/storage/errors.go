package storage

import "errors"

var (
	ErrEventExists = errors.New("event already exists")
	ErrEvent404    = errors.New("event not found")
)

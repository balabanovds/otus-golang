package app

import "errors"

var (
	ErrBusyDate    = errors.New("selected date is busy")
	ErrWrongFormat = errors.New("event wrong format")
)

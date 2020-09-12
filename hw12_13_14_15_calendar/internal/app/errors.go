package app

import "errors"

var (
	ErrForbidden  = errors.New("forbidden action")
	ErrAppGeneral = errors.New("general app error")
)

package router

import "errors"

var (
	ErrWrongMethod = errors.New("wrong method used")
	ErrWrongPath   = errors.New("wrong path")
)

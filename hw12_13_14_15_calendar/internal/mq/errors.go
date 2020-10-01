package mq

import "errors"

var (
	ErrChannelNil = errors.New("channel not opened")
)

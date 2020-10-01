package server

import (
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"
)

type IServer interface {
	Starter
	utils.CloseStringer
}

type Starter interface {
	Start() error
}

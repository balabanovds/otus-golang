package logger

import (
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"go.uber.org/zap"
)

func New(config config.Logger, production bool) error {
	var cfg zap.Config

	if production {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	al := zap.NewAtomicLevel()
	err := al.UnmarshalText([]byte(config.Level))
	if err != nil {
		return err
	}

	cfg.Level.SetLevel(al.Level())

	cfg.OutputPaths = []string{"stderr"}
	if config.LogFile != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, config.LogFile)
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(l)

	return nil
}

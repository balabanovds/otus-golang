package main

import "go.uber.org/zap"

func configLogger(level, logFile string, production bool) error {
	var cfg zap.Config

	if production {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	al := zap.NewAtomicLevel()
	err := al.UnmarshalText([]byte(level))
	if err != nil {
		return err
	}

	cfg.Level.SetLevel(al.Level())

	cfg.OutputPaths = []string{"stderr"}
	if logFile != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, logFile)
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(l)

	return nil
}

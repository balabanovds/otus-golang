package main

import (
	"strings"

	//nolint:goimports
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

const (
	envPrefix = "CAL_"
)

type Config struct {
	Storage storage.Config `koanf:"storage"`
	Server  server.Config  `koanf:"server"`
	Logger  struct {
		Level   string `koanf:"level"`
		LogFile string `koanf:"log_file"`
	} `koanf:"logger"`
	Production bool `koanf:"production"`
}

func NewConfig(fileName string) (*Config, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider(fileName), toml.Parser()); err != nil {
		return nil, err
	}

	if err := k.Load(env.Provider(envPrefix, "_", envCallback), nil); err != nil {
		return nil, err
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func envCallback(s string) string {
	return strings.ToLower(strings.TrimPrefix(s, envPrefix))
}

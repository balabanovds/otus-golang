package config

import (
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

const (
	envPrefix = "CAL_"
)

type Config struct {
	Storage    Storage `koanf:"storage"`
	Server     Server  `koanf:"server"`
	Rmq        Rmq     `koanf:"rmq"`
	Logger     Logger  `koanf:"logger"`
	Production bool    `koanf:"production"`
}

func New(fileName string) (*Config, error) {
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

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
	filename string
	k        *koanf.Koanf
}

func New(filename string) *Config {
	return &Config{
		filename: filename,
		k:        koanf.New("."),
	}
}

func (c *Config) Unmarshal(cfg interface{}) error {
	if err := c.k.Load(file.Provider(c.filename), toml.Parser()); err != nil {
		return err
	}

	if err := c.k.Load(env.Provider(envPrefix, "_", envCallback), nil); err != nil {
		return err
	}

	if err := c.k.Unmarshal("", cfg); err != nil {
		return err
	}

	return nil
}

func envCallback(s string) string {
	return strings.ToLower(strings.TrimPrefix(s, envPrefix))
}

package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	cfg "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/logger"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
	pflag.Parse()
}

type config struct {
	Storage    cfg.Storage `koanf:"storage"`
	Rmq        cfg.Rmq     `koanf:"rmq"`
	Logger     cfg.Logger  `koanf:"logger"`
	Production bool        `koanf:"production"`
}

func main() {
	var c config
	err := cfg.New(configFile).Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	if c.Rmq.Login == "" {
		c.Rmq.Login = "guest"
	}

	if c.Rmq.Password == "" {
		c.Rmq.Password = "guest"
	}

	err = logger.New(c.Logger, c.Production)
	if err != nil {
		log.Fatalf("failed to configure logger: %v\n", err)
	}

	var st storage.IStorage

	if c.Storage.SQL {
		st = sqlstorage.New(c.Storage.Dsn)
	} else {
		st = memorystorage.New()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = st.Connect(ctx); err != nil {
		zap.L().Error("failed to connect to db", zap.Error(err))
		os.Exit(1)
	}

	var wg sync.WaitGroup
	utils.HandleGracefulShutdown(&wg, st)

	wg.Wait()
}
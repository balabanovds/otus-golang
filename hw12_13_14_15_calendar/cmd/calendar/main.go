package main

import (
	"context"
	"log"
	"os"
	"time"

	cfg "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/logger"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/grpcsrv"
	internalhttp "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory" //nolint:gci
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
	HTTP       cfg.HTTP    `koanf:"http"`
	GRPC       cfg.GRPC    `koanf:"grpc"`
	Logger     cfg.Logger  `koanf:"logger"`
	Production bool        `koanf:"production"`
}

func main() {
	var c config
	err := cfg.New(configFile).Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	l, err := logger.New(c.Logger, c.Production)
	if err != nil {
		log.Fatalf("failed to configure logger: %v\n", err)
	}
	defer func() {
		_ = l.Sync()
	}()

	var st storage.IStorage

	if c.Storage.SQL {
		st = sqlstorage.New(c.Storage)
	} else {
		st = memorystorage.New()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = st.Connect(ctx); err != nil {
		zap.L().Error("failed to connect to db", zap.Error(err))
		os.Exit(1)
	}

	calendar := app.New(st)

	grpc := grpcsrv.New(calendar, c.GRPC)
	httpsrv := internalhttp.New(calendar, c.HTTP)
	defer utils.Close(st, grpc, httpsrv)

	doneCh := make(chan struct{})

	go utils.HandleGracefulShutdown(st, grpc, httpsrv)

	if err := fireUp(httpsrv, grpc); err != nil {
		doneCh <- struct{}{}
	}

	<-doneCh
}

func fireUp(starters ...server.Starter) error {
	done := make(chan error)
	for _, st := range starters {
		go func(s server.Starter) {
			if err := s.Start(); err != nil {
				done <- err
			}
		}(st)
	}
	return <-done
}

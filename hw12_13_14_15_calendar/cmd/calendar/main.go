package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
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

func main() {
	config, err := config.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = logger.New(config.Logger, config.Production)
	if err != nil {
		log.Fatalf("failed to configure logger: %v\n", err)
	}

	var st storage.IStorage

	if config.Storage.SQL {
		st = sqlstorage.New(config.Storage.Dsn)
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

	var srv server.IServer
	if config.Server.Grpc {
		srv = grpcsrv.New(calendar, config.Server)
	} else {
		srv = internalhttp.New(calendar, config.Server)
	}

	var wg sync.WaitGroup
	go utils.HandleGracefulShutdown(&wg, st, srv)

	if err := srv.Start(); err != nil {
		os.Exit(1)
	}
	wg.Wait()
}

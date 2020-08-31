package main

import (
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory" //nolint:gci
	"github.com/spf13/pflag"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	pflag.Parse()

	if configFile == "" {
		pflag.Usage()
		os.Exit(1)
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = configLogger(config.Logger.Level, config.Logger.LogFile, config.Production)
	if err != nil {
		log.Fatalf("failed to configure logger: %v\n", err)
	}

	storage := memorystorage.New()
	calendar := app.New(storage)

	server := internalhttp.NewServer(calendar, config.Server)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)

		if err := server.Stop(); err != nil {
			zap.L().Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(); err != nil {
		zap.L().Error("failed to start http server: " + err.Error())
		os.Exit(1)
	}
}

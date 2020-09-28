package main

import (
	"context"
	"log"
	"os"
	"time"

	cfg "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/logger"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/mq/rabbitmq"
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
	Storage    cfg.Storage   `koanf:"storage"`
	Rmq        cfg.Rmq       `koanf:"rmq"`
	Scheduler  cfg.Scheduler `koanf:"scheduler"`
	Logger     cfg.Logger    `koanf:"logger"`
	Production bool          `koanf:"production"`
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

	ctx, cancelT := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelT()

	if err = st.Connect(ctx); err != nil {
		zap.L().Error("failed to connect to db", zap.Error(err))
		os.Exit(1)
	}

	pub, err := rabbitmq.NewPublisher(c.Rmq)
	if err != nil {
		zap.L().Error("failed to connect to amqp", zap.Error(err))
		os.Exit(1)
	}
	defer utils.Close(st, pub)

	ctx, cancelC := context.WithCancel(context.Background())
	defer cancelC()

	sch := new(pub, st, time.Duration(c.Scheduler.Interval)*time.Second)
	sch.run(ctx)
	doneCh := make(chan struct{})
	go utils.HandleGracefulShutdown(st, pub, sch)
	// go every(ctx, doneCh, time.Second, publishEvents)
	// go every(ctx, doneCh, time.Second, clearEvents)

	<-doneCh
}

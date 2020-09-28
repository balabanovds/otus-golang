package main

import (
	"context"
	"log"
	"os"

	cfg "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/logger"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/mq/rabbitmq"
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
	Rmq        cfg.Rmq    `koanf:"rmq"`
	Logger     cfg.Logger `koanf:"logger"`
	Production bool       `koanf:"production"`
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

	sub, err := rabbitmq.NewConsumer(c.Rmq)
	if err != nil {
		zap.L().Error("failed to connect to amqp", zap.Error(err))
		os.Exit(1)
	}
	defer utils.Close(sub)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go utils.HandleGracefulShutdown(sub)

	if err := sub.Channel().Open(); err != nil {
		zap.L().Error("open channel", zap.Error(err))
		os.Exit(1)
	}
	defer func() {
		if err := sub.Channel().Close(); err != nil {
			zap.L().Error("close channel", zap.Error(err))
		}
	}()

	msgsCh, err := sub.Consume(ctx)
	if err != nil {
		zap.L().Error("close channel", zap.Error(err))
		return
	}

	for msg := range msgsCh {
		if msg.Err != nil {
			zap.L().Error("err received in message", zap.Error(err))
			continue
		}
		send(msg.Data)
	}
}

func send(msg models.MQNotification) {
	zap.L().Info("sender: send message",
		zap.Int("event_id", msg.EventID),
		zap.Int("user_id", msg.UserID),
		zap.String("title", msg.Title),
		zap.Time("at", msg.Date),
	)
}

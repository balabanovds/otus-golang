package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/amqp"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type scheduler struct {
	pub      amqp.Publisher
	st       storage.IStorage
	interval time.Duration
	doneCh   chan struct{}
}

func new(pub amqp.Publisher, st storage.IStorage, interval time.Duration) *scheduler {
	return &scheduler{
		pub:      pub,
		st:       st,
		interval: interval,
		doneCh:   make(chan struct{}),
	}
}

func (s *scheduler) run(ctx context.Context) {
	logInfo("RUN", zap.Duration("interval", s.interval))
	go every(ctx, s.doneCh, s.interval, s.publishEvents)
	go every(ctx, s.doneCh, s.interval, s.clearEvents)
	<-s.doneCh
}

func (s *scheduler) Close() error {
	s.doneCh <- struct{}{}
	return nil
}

func (s *scheduler) String() string {
	return "scheduler"
}

func (s *scheduler) publishEvents(ctx context.Context, date time.Time) {
	if err := s.pub.Channel().Open(); err != nil {
		zap.L().Error("open channel", zap.Error(err))
		return
	}
	defer func() {
		if err := s.pub.Channel().Close(); err != nil {
			logErr("close channel", zap.Error(err))
		}
	}()

	logInfo("publish start")
	list := s.st.Events().ListByReminderBetweenDates(ctx, date, date.Add(s.interval))

	cntr := 0
	start := time.Now()
	for _, ev := range list {
		data, err := json.Marshal(models.NewMQNotification(ev))
		if err != nil {
			logErr("marshal notification", zap.Error(err))
			continue
		}
		if err := s.pub.Publish(ctx, data); err != nil {
			logErr("publish event", zap.Error(err))
			continue
		}
		cntr++
	}
	logInfo("publish report",
		zap.Int("total", len(list)),
		zap.Int("count", cntr),
		zap.Duration("duration", time.Since(start)),
	)
}

// clearEvents clears all events older that 1 Year since date.
func (s *scheduler) clearEvents(ctx context.Context, date time.Time) {
	logInfo("clear old events start")
	start := time.Now()
	date = date.AddDate(-1, 0, 0)

	list := s.st.Events().ListBeforeDate(ctx, date)
	for _, ev := range list {
		s.st.Events().Delete(ctx, ev.ID)
	}

	logInfo("clear old events report",
		zap.Int("total", len(list)),
		zap.Duration("duration", time.Since(start)),
	)
}

func logInfo(msg string, fields ...zapcore.Field) {
	zap.L().Info("scheduler", append([]zapcore.Field{zap.String("msg", msg)}, fields...)...)
}

func logErr(msg string, fields ...zapcore.Field) {
	zap.L().Error("scheduler", append([]zapcore.Field{zap.String("msg", msg)}, fields...)...)
}

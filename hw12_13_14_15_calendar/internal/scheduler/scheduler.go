package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/amqp"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type Scheduler struct {
	pub      amqp.Publisher
	st       storage.IStorage
	interval time.Duration
}

func New(pub amqp.Publisher, st storage.IStorage, interval time.Duration) *Scheduler {
	return &Scheduler{
		pub:      pub,
		st:       st,
		interval: interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	every(ctx, s.interval, s.publishEvents)
	every(ctx, s.interval, s.clearEvents)
}

func (s *Scheduler) publishEvents(ctx context.Context, date time.Time) {
	zap.L().Info("scheduler: publish start")
	list := s.st.Events().ListByReminderBetweenDates(ctx, date, date.Add(s.interval))

	cntr := 0
	start := time.Now()
	for _, ev := range list {
		data, err := json.Marshal(models.NewMQNotification(ev))
		if err != nil {
			zap.L().Error("scheduler: marshal notification", zap.Error(err))
			continue
		}
		if err := s.pub.Publish(data); err != nil {
			zap.L().Error("scheduler: publish event", zap.Error(err))
			continue
		}
		cntr++
	}
	zap.L().Info("scheduler: publish report",
		zap.Int("total", len(list)),
		zap.Int("count", cntr),
		zap.Duration("duration", time.Since(start)),
	)
}

// clearEvents clears all events older that 1 Year since date
func (s *Scheduler) clearEvents(ctx context.Context, date time.Time) {
	start := time.Now()
	date = date.AddDate(-1, 0, 0)

	list := s.st.Events().ListBeforeDate(ctx, date)
	for _, ev := range list {
		s.st.Events().Delete(ctx, ev.ID)
	}

	zap.L().Info("scheduler: clear old events report",
		zap.Int("total", len(list)),
		zap.Duration("duration", time.Since(start)),
	)
}

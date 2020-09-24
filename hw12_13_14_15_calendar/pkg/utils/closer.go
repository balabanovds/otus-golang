package utils

import (
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type CloseStringer interface {
	Close() error
	String() string
}

func Close(closers ...io.Closer) {
	for _, cl := range closers {
		if err := cl.Close(); err != nil {
			zap.L().Error("error to close", zap.Error(err))
		}
	}
}

// HandleGracefulShutdown monitor for signals to graceful shutdown closers.
func HandleGracefulShutdown(wg *sync.WaitGroup, closers ...CloseStringer) {
	wg.Add(1)
	defer wg.Done()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	s := <-signals
	zap.L().Info("closing services by signal", zap.String("signal", s.String()))
	signal.Stop(signals)

	cntr := 0
	start := time.Now()

	for _, cl := range closers {
		if err := cl.Close(); err != nil {
			zap.L().Error("faled to close",
				zap.String("service", cl.String()),
				zap.Error(err))
			continue
		}
		cntr++
	}

	zap.L().Info("close report",
		zap.Int("closed", cntr),
		zap.Int("total", len(closers)),
		zap.Duration("duration", time.Since(start)),
	)
}

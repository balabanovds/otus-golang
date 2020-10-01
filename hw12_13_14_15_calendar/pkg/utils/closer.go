package utils

import (
	"io"
	"os"
	"os/signal"
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
		if cl != nil {
			if err := cl.Close(); err != nil {
				zap.L().Error("error to close", zap.Error(err))
			}
		}
	}
}

// HandleGracefulShutdown monitor for signals to graceful shutdown closers.
func HandleGracefulShutdown(closers ...CloseStringer) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	s := <-signals
	zap.L().Info("closing services by signal", zap.String("signal", s.String()))
	signal.Stop(signals)

	cntr := 0
	start := time.Now()

	for _, cl := range closers {
		if err := cl.Close(); err != nil {
			zap.L().Error("failed to close",
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
	os.Exit(0)
}

func HandleSigterm(doneCh chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	s := <-signals
	zap.L().Info("closing services by signal", zap.String("signal", s.String()))
	signal.Stop(signals)
	doneCh <- struct{}{}
}

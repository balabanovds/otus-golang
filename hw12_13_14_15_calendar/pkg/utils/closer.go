package utils

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type CloseStringer interface {
	Close() error
	String() string
}

// HandleGracefulShutdown monitor for signals to graceful shutdown closers.
func HandleGracefulShutdown(wg *sync.WaitGroup, closers ...CloseStringer) {
	wg.Add(1)
	defer wg.Done()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	s := <-signals
	log.Printf("catched signal %s. closing services", s.String())
	signal.Stop(signals)

	cntr := 0
	start := time.Now()

	for _, cl := range closers {
		if err := cl.Close(); err != nil {
			log.Printf("failed to stop %s: %v", cl, err)
			continue
		}
		cntr++
	}

	log.Printf("gracefully closed %d/%d services. took %d ns",
		cntr, len(closers), time.Since(start))
}

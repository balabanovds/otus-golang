package memorystorage

import (
	"context"
	"sync"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type Storage struct {
	mu           sync.RWMutex
	data         map[int]models.Event
	eventStorage storage.IEventStorage
}

func New() storage.IStorage {
	return &Storage{
		data: make(map[int]models.Event),
	}
}

func (s *Storage) Events() storage.IEventStorage {
	if s.eventStorage == nil {
		s.eventStorage = newEventStorage(s)
	}
	return s.eventStorage
}

func (s *Storage) Connect(_ context.Context) error {
	zap.L().Info("connected to memory storage")
	return nil
}

func (s *Storage) Close() error {
	zap.L().Info("closing memory storage")
	return nil
}

func (s *Storage) String() string {
	return "memory storage"
}

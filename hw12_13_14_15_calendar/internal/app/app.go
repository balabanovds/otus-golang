package app

import (
	"context"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage storage.IStorage
}

func New(storage storage.IStorage) *App {
	return &App{storage}
}

func (a *App) CreateEvent(ctx context.Context, id string, title string) error {
	return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO

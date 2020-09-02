package app

import (
	"context"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage storage.Repo
}

func New(storage storage.Repo) *App {
	return &App{storage}
}

func (a *App) CreateEvent(ctx context.Context, id string, title string) error {
	_, err := a.storage.Events().Create(storage.Event{ID: id, Title: title})
	return err
}

// TODO

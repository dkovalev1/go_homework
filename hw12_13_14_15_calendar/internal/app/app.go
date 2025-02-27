package app

import (
	"context"
	"errors"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"  //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
)

var (
	ErrNotImpl  = errors.New("not implemented")
	ErrNotFound = errors.New("event not found")
)

type App struct { // TODO
	Log     *logger.Logger
	Storage Storage
}

type Storage interface { // TODO
	CreateEvent(storage.Event) error
	UpdateEvent(storage.Event) error
	DeleteEvent(storage.Event) error
	GetAllEventsDay(time.Time) ([]storage.Event, error)
	GetAllEventsWeek(time.Time) ([]storage.Event, error)
	GetAllEventsMonth(time.Time) ([]storage.Event, error)
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{
		Log:     logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(_ context.Context, id, title string) error {
	return a.Storage.CreateEvent(storage.Event{ID: id, Title: title})
}

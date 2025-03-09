package scheduler

import (
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"    //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/config" //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger" //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/rabbit" //nolint
)

type Scheduler struct {
	interval   time.Duration
	keepEvents time.Duration
	storage    app.Storage
	logger     *logger.Logger
	queue      *rabbit.Rabbit
}

func New(config *config.RabbitConf, storage app.Storage) *Scheduler {
	ret := &Scheduler{
		interval: config.Interval,
		storage:  storage,
	}

	return ret
}

func (s *Scheduler) Run(interruptChan <-chan struct{}) error {
	// Main loop
	for {
		s.performSendEvents()
		s.deleteOldEvents()

		select {
		case <-time.After(s.interval):
			// Perform scheduled tasks on the next step
		case <-interruptChan:
			// Got interrupt signal, exit loop if needed
			return nil
		}
	}
}

func (s *Scheduler) Start(interruptChan <-chan struct{}) error {
	go s.Run(interruptChan)
	return nil
}

func (s *Scheduler) performSendEvents() {
	// Perform sender tasks here
	events, err := s.storage.GetUpcomingEvents(time.Now())
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	for _, event := range events {
		// It could be more efficient to filter out
		// already sent events on the database level,
		// but for the sake of stricter program layer design will do it here
		if !event.NotificationSent {
			continue
		}

		notification := rabbit.NewNotification(event)
		err = s.queue.SendNotification(notification)
		if err != nil {
			s.logger.Error(err.Error())
		}
	}
}

func (s *Scheduler) deleteOldEvents() {
	// Perform deletion tasks here
	old := time.Now().Add(-s.keepEvents)
	err := s.storage.DeleteEventOlderThan(old)
	if err != nil {
		s.logger.Error(err.Error())
	}
}

package scheduler

import (
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/rabbit"
)

type Scheduler struct {
	interval   time.Duration
	keepEvents time.Duration
	storage    app.Storage
	logger     *logger.Logger
	queue      *rabbit.Rabbit
}

func New(config *config.Config) *Scheduler {

	ret := &Scheduler{
		interval: config.RabbitMQ.Interval,
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
			// Perform scheduled tasks here
		case <-interruptChan:
			// Handle interrupt signal and exit loop if needed
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

	events, err := s.storage.GetAllCurrentEvents()
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	for _, event := range events {
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

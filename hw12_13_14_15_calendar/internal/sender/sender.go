package sender

import (
	"fmt"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"                    //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/config"                 //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"                 //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/rabbit"                 //nolint
	sqlstorage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage/sql" //nolint
)

type Sender struct {
	config *config.Config
	logger *logger.Logger
	rabbit rabbit.IRabbit

	// used to store rabbit notificatrions in the database
	db app.Storage
}

func New(config *config.Config, logger *logger.Logger) *Sender {
	r := rabbit.NewRabbit(config.RabbitMQ.Connstr, logger)

	var storage app.Storage

	if config.Storage.Type == "SQL" {
		storage = sqlstorage.New(config.Storage.Connstr)
	}

	return &Sender{
		config: config,
		logger: logger,
		rabbit: r,
		db:     storage,
	}
}

func (s *Sender) Run(interruptChan <-chan struct{}) error {
	s.logger.Info("Sender started")
	err := s.rabbit.ReceiveNotifications(interruptChan, func(n *rabbit.Notification) {
		s.logger.Info("Received a message: " + n.EventID + " " + n.Title + " " + n.Time.String())
		s.reportMessage(n)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) Start(interruptChan <-chan struct{}) error {
	go s.Run(interruptChan)
	return nil
}

func (s *Sender) reportMessage(n *rabbit.Notification) {
	if s.db != nil {
		err := s.db.AddNotification(n.EventID, n.Title, n.Time, int(n.User))
		if err != nil {
			s.logger.Error(fmt.Errorf("could not save notification: %w", err).Error())
		} else {
			s.logger.Info("Stored notification " + n.EventID)
		}
	}
}

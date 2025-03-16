package sender

import (
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/config" //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger" //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/rabbit" //nolint
)

type Sender struct {
	config *config.Config
	logger *logger.Logger
	rabbit rabbit.IRabbit
}

func New(config *config.Config, logger *logger.Logger) *Sender {
	r := rabbit.NewRabbit(config.RabbitMQ.Connstr, logger)

	return &Sender{
		config: config,
		logger: logger,
		rabbit: r,
	}
}

func (s *Sender) Run(interruptChan <-chan struct{}) error {
	s.logger.Info("Sender started")
	err := s.rabbit.ReceiveNotifications(interruptChan, func(n *rabbit.Notification) {
		s.logger.Info("Received a message: " + n.EventID + " " + n.Title + " " + n.Time.String())
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

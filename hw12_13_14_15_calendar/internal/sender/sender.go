package sender

import (
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/rabbit"
)

type Sender struct {
	config *config.Config
	app    app.App
	logger logger.Logger
	rabbit rabbit.Rabbit
}

func New(config *config.Config) *Sender {
	return &Sender{
		config: config,
	}
}

func (s *Sender) Run(interruptChan <-chan struct{}) error {

}

func (s *Sender) Start(interruptChan <-chan struct{}) error {
	for {

	}
}

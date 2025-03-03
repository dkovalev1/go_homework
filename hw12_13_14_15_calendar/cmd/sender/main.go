package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/sender"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := config.NewConfig(configFile)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	interruptChan := make(chan struct{})

	sender := sender.New(config)
	sender.Start(interruptChan)

	<-ctx.Done()
	interruptChan <- struct{}{}
}

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"                          //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/config"                       //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"                       //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/scheduler"                    //nolint
	memorystorage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage/memory" //nolint
	sqlstorage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage/sql"       //nolint
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := config.NewConfig(configFile)

	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	switch config.Storage.Type {
	case "IM":
		storage = memorystorage.New()
	case "SQL":
		storage = sqlstorage.New(config.Storage.Connstr)
	default:
		logg.Error("unknown storage type: " + config.Storage.Type)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	interruptChan := make(chan struct{})

	scheduler := scheduler.New(&config.RabbitMQ, logg, storage)
	scheduler.Start(interruptChan)

	<-ctx.Done()
	interruptChan <- struct{}{}
}

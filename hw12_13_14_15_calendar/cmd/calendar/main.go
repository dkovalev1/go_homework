package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"                          //nolint
	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"                //nolint
	internalhttp "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/server/http"     //nolint
	memorystorage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage/memory" //nolint
	sqlstorage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage/sql"       //nolint
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
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

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

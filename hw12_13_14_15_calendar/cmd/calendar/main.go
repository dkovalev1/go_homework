package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"                          //nolint
	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"                //nolint
	internalgrpc "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/server/grpc"     //nolint
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

	server := internalhttp.NewServer(config.REST.Port, logg, calendar.Storage)
	grpcServer := internalgrpc.NewService(config.GRPC.Port, logg, calendar.Storage)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	if err := grpcServer.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		cancel()
		os.Exit(1)
	}

	<-ctx.Done()

	if err := server.Stop(ctx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}

	if err := grpcServer.Stop(); err != nil {
		logg.Error("failed to stop grpc server: " + err.Error())
	}
}

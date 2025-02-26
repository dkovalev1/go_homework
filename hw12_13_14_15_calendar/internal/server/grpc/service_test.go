package internalgrpc

import (
	"context"
	"log"
	"testing"
	"time"

	calendarpb "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/api"                        //nolint
	app "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"                      //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"                       //nolint
	memorystorage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage/memory" //nolint
	"github.com/stretchr/testify/require"                                                           //nolint
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type testContext struct {
	server  *CalendarService
	storage app.Storage
	client  calendarpb.CalendarClient
}

func setUp(t *testing.T) *testContext {
	t.Helper()

	testLogger := logger.New("debug")
	testStorage := memorystorage.New()

	testServer := NewService(testLogger, testStorage)

	err := testServer.Start()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := calendarpb.NewCalendarClient(conn)

	return &testContext{
		storage: testStorage,
		server:  testServer,
		client:  client,
	}
}

func tearDown(t *testing.T, ctx *testContext) {
	t.Helper()
	ctx.server.Stop()
}

func TestStorage(t *testing.T) {
	ctx := setUp(t)
	defer func() {
		tearDown(t, ctx)
	}()

	t.Run("CreateEvent", func(t *testing.T) {
		result, err := ctx.client.CreateEvent(context.Background(), &calendarpb.Event{
			ID:        "1",
			Title:     "titletest",
			StartTime: timestamppb.New(time.Now().Add(1 * time.Hour)),
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		// Check in storage that event was created
		events, err := ctx.storage.GetAllEventsDay(time.Now())
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0].ID, "1")
		require.Equal(t, events[0].Title, "titletest")
	})
}

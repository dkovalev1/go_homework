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

const testPort = 50051

func setUp(t *testing.T) *testContext {
	t.Helper()

	testLogger := logger.New("debug")
	testStorage := memorystorage.New()

	testServer := NewService(testPort, testLogger, testStorage)

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

func TestServiceGRPC(t *testing.T) {
	ctx := setUp(t)
	defer func() {
		tearDown(t, ctx)
	}()

	startTime := time.Now().Add(1 * time.Hour).UTC()
	t.Run("CreateEvent", func(t *testing.T) {
		result, err := ctx.client.CreateEvent(context.Background(), &calendarpb.Event{
			ID:        "1",
			Title:     "titletest",
			StartTime: timestamppb.New(startTime),
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		// Check in storage that event was created
		events, err := ctx.storage.GetAllEventsDay(time.Now())
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, "1", events[0].ID)
		require.Equal(t, "titletest", events[0].Title)
		require.Equal(t, startTime, events[0].StartTime)
	})
	t.Run("UpdateEvent", func(t *testing.T) {
		// create test event
		result, err := ctx.client.CreateEvent(context.Background(), &calendarpb.Event{
			ID:        "1",
			Title:     "titletest",
			StartTime: timestamppb.New(time.Now().Add(1 * time.Hour)),
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		result, err = ctx.client.UpdateEvent(context.Background(), &calendarpb.Event{
			ID:        "1",
			Title:     "Updated title",
			StartTime: timestamppb.New(time.Now().Add(1 * time.Hour)),
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		// Check in storage that event was updated
		events, err := ctx.storage.GetAllEventsDay(time.Now())
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, "1", events[0].ID)
		require.Equal(t, "Updated title", events[0].Title)
	})
	t.Run("DeleteEvent", func(t *testing.T) {
		// create test event
		result, err := ctx.client.CreateEvent(context.Background(), &calendarpb.Event{
			ID:        "1",
			Title:     "titletest",
			StartTime: timestamppb.New(time.Now().Add(1 * time.Hour)),
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		result, err = ctx.client.DeleteEvent(context.Background(), &calendarpb.EventId{
			Id: "1",
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		// Check in storage that event was deleted
		events, err := ctx.storage.GetAllEventsDay(time.Now())
		require.NoError(t, err)
		require.Equal(t, 0, len(events))
	})
	t.Run("GetAllEventsDay", func(t *testing.T) {
		startTime := time.Now().Add(1 * time.Hour)
		// create test event
		result, err := ctx.client.CreateEvent(context.Background(), &calendarpb.Event{
			ID:        "1",
			Title:     "titletest",
			StartTime: timestamppb.New(startTime),
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		events, err := ctx.client.GetAllEventsDay(context.Background(), &calendarpb.TimeSpec{
			Stamp: timestamppb.New(startTime),
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(events.Events))
		require.Equal(t, "1", events.Events[0].ID)
		require.Equal(t, "titletest", events.Events[0].Title)
	})

	t.Run("GetAllEventsWeek", func(t *testing.T) {
		// create test event
		result, err := ctx.client.CreateEvent(context.Background(), &calendarpb.Event{
			ID:        "1",
			Title:     "titletest",
			StartTime: timestamppb.New(time.Now().Add(1 * time.Hour)),
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		events, err := ctx.client.GetAllEventsWeek(context.Background(), &calendarpb.TimeSpec{
			Stamp: timestamppb.New(startTime),
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(events.Events))
		require.Equal(t, "1", events.Events[0].ID)
		require.Equal(t, "titletest", events.Events[0].Title)
	})

	t.Run("GetAllEventsMonth", func(t *testing.T) {
		// create test event
		result, err := ctx.client.CreateEvent(context.Background(), &calendarpb.Event{
			ID:        "1",
			Title:     "titletest",
			StartTime: timestamppb.New(time.Now().Add(1 * time.Hour)),
		})
		require.NoError(t, err)
		require.True(t, result.IsOk)

		events, err := ctx.client.GetAllEventsMonth(context.Background(), &calendarpb.TimeSpec{
			Stamp: timestamppb.New(startTime),
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(events.Events))
		require.Equal(t, "1", events.Events[0].ID)
		require.Equal(t, "titletest", events.Events[0].Title)
	})
}

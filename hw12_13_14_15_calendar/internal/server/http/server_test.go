package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"                          //nolint
	logger "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"                //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage"                      //nolint
	memorystorage "github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage/memory" //nolint
	"github.com/stretchr/testify/require"                                                           //nolint
	//nolint
	//nolint
	//nolint
)

type testContext struct {
	server  *Server
	ts      *httptest.Server
	storage app.Storage
	ctx     context.Context
}

const (
	testPort   = 50052
	eventID    = "testrest"
	eventTitle = "We will test the REST"
)

func setUp(t *testing.T) *testContext {
	t.Helper()

	testLogger := logger.New("debug")
	testStorage := memorystorage.New()

	testServer := NewServer(testPort, testLogger, testStorage)

	ts := httptest.NewServer(testServer.server.Handler)

	serverContext := context.Background()

	return &testContext{
		storage: testStorage,
		server:  testServer,
		ts:      ts,
		ctx:     serverContext,
	}
}

func tearDown(t *testing.T, ctx *testContext) {
	t.Helper()
	ctx.ts.Close()
	// ctx.server.Stop(ctx.ctx)
}

func makeTestEvent(t *testing.T, eventID, eventTitle string, eventStart time.Time) []byte {
	t.Helper()
	event := &storage.Event{
		ID:        eventID,
		Title:     eventTitle,
		StartTime: eventStart,
	}

	b, err := json.Marshal(event)
	require.NoError(t, err)
	return b
}

func testCreateEvent(t *testing.T, testC *testContext, eventStart time.Time) {
	t.Helper()
	client := &http.Client{}
	b := makeTestEvent(t, eventID, eventTitle, eventStart)
	url := fmt.Sprintf("%s/createevent", testC.ts.URL)

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(b))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, "200 OK", resp.Status)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, body, []byte("ok\n"))
}

func TestServiceREST(t *testing.T) {
	eventStart := time.Now()
	client := &http.Client{}

	t.Run("CreateEvent", func(t *testing.T) {
		// To provide tests isolation by the cost of performance, create a storage for each test
		testC := setUp(t)
		defer func() {
			tearDown(t, testC)
		}()

		testCreateEvent(t, testC, eventStart)

		// Check in storage that event was created
		events, err := testC.storage.GetAllEventsDay(time.Now().Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, eventID, events[0].ID)
		require.Equal(t, eventTitle, events[0].Title)
		require.True(t, eventStart.Equal(events[0].StartTime))
	})
	t.Run("UpdateEvent", func(t *testing.T) {
		testC := setUp(t)
		defer func() {
			tearDown(t, testC)
		}()
		// create test event
		testCreateEvent(t, testC, eventStart)

		// update test event
		url := fmt.Sprintf("%s/updateevent", testC.ts.URL)
		b := makeTestEvent(t, eventID, "updated title", eventStart)
		req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(b))
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, "200 OK", resp.Status)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, body, []byte("ok\n"))

		// Check in storage that event was updated
		events, err := testC.storage.GetAllEventsDay(time.Now().Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0].ID, eventID)
		require.Equal(t, events[0].Title, "updated title")
	})
	t.Run("DeleteEvent", func(t *testing.T) {
		testC := setUp(t)
		defer func() {
			tearDown(t, testC)
		}()
		// create test event
		testCreateEvent(t, testC, eventStart)

		// delete test event
		url := fmt.Sprintf("%s/deleteevent", testC.ts.URL)
		b := makeTestEvent(t, eventID, "", time.Time{})
		req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(b))
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, "200 OK", resp.Status)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, body, []byte("ok\n"))

		// Check in storage that event was deleted
		events, err := testC.storage.GetAllEventsDay(time.Now().Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Equal(t, 0, len(events))
	})
	t.Run("GetAllEventsDay", func(t *testing.T) {
		testC := setUp(t)
		defer func() {
			tearDown(t, testC)
		}()

		testCreateEvent(t, testC, eventStart)

		events, err := testC.storage.GetAllEventsDay(time.Now().Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0].ID, eventID)
		require.Equal(t, events[0].Title, eventTitle)
		require.True(t, eventStart.Equal(events[0].StartTime))
	})
	t.Run("GetAllEventsWeek", func(t *testing.T) {
		testC := setUp(t)
		defer func() {
			tearDown(t, testC)
		}()

		testCreateEvent(t, testC, eventStart)

		events, err := testC.storage.GetAllEventsWeek(time.Now().Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0].ID, eventID)
		require.Equal(t, events[0].Title, eventTitle)
		require.True(t, eventStart.Equal(events[0].StartTime))
	})
	t.Run("GetAllEventsMonth", func(t *testing.T) {
		testC := setUp(t)
		defer func() {
			tearDown(t, testC)
		}()
		testCreateEvent(t, testC, eventStart)

		events, err := testC.storage.GetAllEventsMonth(time.Now().Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0].ID, eventID)
		require.Equal(t, events[0].Title, eventTitle)
		require.True(t, eventStart.Equal(events[0].StartTime))
	})
}

package sqlstorage

import (
	"fmt"
	"testing"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"     //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/jmoiron/sqlx"                                                  //nolint
	"github.com/stretchr/testify/require"                                      //nolint
)

const (
	connstr     = "host=localhost user=calendar password=calendar dbname=calendardb sslmode=disable"
	testEventID = "test"
	title       = "titletest"
)

type testContext struct {
	con *sqlx.DB
}

func setUp(t *testing.T) *testContext {
	t.Helper()
	con, err := sqlx.Connect("postgres", connstr)
	if err != nil {
		t.Fatal(err)
	}

	return &testContext{
		con: con,
	}
}

func tearDown(t *testing.T, ctx *testContext) {
	t.Helper()
	_, err := ctx.con.Exec("DELETE FROM event")
	require.NoError(t, err)
	ctx.con.Close()
}

func TestStorage(t *testing.T) { //nolint:funlen
	t.Run("CreateEvent", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		s := New(connstr)

		now := time.Now()

		err := s.CreateEvent(storage.Event{
			ID:         testEventID,
			Title:      title,
			StartTime:  now,
			NotifyTime: time.Hour,
		})
		require.NoError(t, err)

		// Do check
		type dbEvent struct {
			ID         string
			Title      string
			StartTime  time.Time
			NotifyTime int
		}

		var events []dbEvent
		err = context.con.Select(&events, "SELECT id,title,starttime,notifytime FROM event")
		require.NoError(t, err)
		require.Len(t, events, 1)

		require.Equal(t, testEventID, events[0].ID)
		require.Equal(t, title, events[0].Title)
		delta := now.Sub(events[0].StartTime)

		// db can not support nanoseconds, so lets consider precision is 1 mcs
		require.True(t, delta < time.Microsecond)
		require.Equal(t, 3600, events[0].NotifyTime)
	})

	t.Run("CreateEvents", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		s := New(connstr)

		title := "titletest"
		start := time.Now()

		const nEvents = 2000
		for i := range nEvents {
			id := fmt.Sprintf("%s_%d", testEventID, i)
			start = start.Add(10 * time.Minute)
			err := s.CreateEvent(storage.Event{
				ID:         id,
				Title:      title,
				StartTime:  start,
				NotifyTime: time.Hour,
			})
			require.NoError(t, err)
		}

		// Do check

		var count []int
		err := context.con.Select(&count, "SELECT COUNT(*) FROM event")
		require.NoError(t, err)
		require.Equal(t, nEvents, count[0])

		type dbEvent struct {
			ID         string
			Title      string
			StartTime  time.Time
			NotifyTime int
		}
		var events []dbEvent
		err = context.con.Select(&events, "SELECT id,title,starttime,notifytime FROM event")
		require.NoError(t, err)
		require.Equal(t, nEvents, len(events))
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		s := New(connstr)

		title := "titletest"
		updatedTitle := "updatedtitle"

		err := s.CreateEvent(storage.Event{
			ID:    testEventID,
			Title: title,
		})
		require.NoError(t, err)

		err = s.UpdateEvent(storage.Event{
			ID:    testEventID,
			Title: updatedTitle,
		})
		require.NoError(t, err)

		// Do check
		var events []storage.Event
		err = context.con.Select(&events, "SELECT id,title FROM event")
		require.NoError(t, err)
		require.Len(t, events, 1)

		require.Equal(t, testEventID, events[0].ID)
		require.Equal(t, updatedTitle, events[0].Title)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		s := New(connstr)

		err := s.CreateEvent(storage.Event{
			ID:    testEventID,
			Title: "titletest",
		})
		require.NoError(t, err)

		err = s.DeleteEvent(testEventID)
		require.NoError(t, err)

		var events []storage.Event
		err = context.con.Select(&events, "SELECT id,title FROM event")

		require.NoError(t, err)
		require.Len(t, events, 0)
		require.Empty(t, events)

		err = s.DeleteEvent(testEventID)
		require.Error(t, err)
		require.ErrorIs(t, err, app.ErrNotFound)
	})

	t.Run("MarkEventAsNotificationSent", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		s := New(connstr)

		err := s.CreateEvent(storage.Event{
			ID:    testEventID,
			Title: "titletest",
		})
		require.NoError(t, err)

		err = s.MarkEventAsNotificationSent(testEventID)
		require.NoError(t, err)

		// Do check
		var events []storage.Event
		err = context.con.Select(&events, "SELECT id,notification_sent FROM event")
		require.NoError(t, err)
		require.Len(t, events, 1)

		require.Equal(t, testEventID, events[0].ID)
		require.True(t, events[0].NotificationSent)

		err = s.MarkEventAsNotificationSent(testEventID)
		require.NoError(t, err)

		err = s.MarkEventAsNotificationSent("nonexist")
		require.Error(t, err)
		require.ErrorIs(t, err, app.ErrNotFound)
	})

	t.Run("DeleteEventOlderThan", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		s := New(connstr)

		// Create test events
		times := []time.Time{
			time.Now().Add(-48 * time.Hour),
			time.Now().Add(-24 * time.Hour),
			time.Now().Add(24 * time.Hour),
		}

		for i, startTime := range times {
			err := s.CreateEvent(storage.Event{
				ID:        fmt.Sprintf("test%d", i),
				Title:     fmt.Sprintf("titletest%d", i),
				StartTime: startTime,
			})
			require.NoError(t, err)
		}

		// Delete events older than 23 hours
		err := s.DeleteEventOlderThan(time.Now().Add(-23 * time.Hour))
		require.NoError(t, err)

		// Check remaining events
		var events []storage.Event
		err = context.con.Select(&events, "SELECT id, starttime FROM event")
		require.NoError(t, err)
		require.Len(t, events, 1)

		// Ensure the remaining event is the future one
		require.Equal(t, "test2", events[0].ID)
		require.True(t, events[0].StartTime.After(time.Now()))
	})

	t.Run("GetUpcomingEvents", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		s := New(connstr)

		// Create test events
		times := []time.Time{
			time.Now().Add(-24 * time.Hour),
			time.Now().Add(1 * time.Hour),
			time.Now().Add(2 * time.Hour),
			time.Now().Add(24 * time.Hour),
		}

		for i, startTime := range times {
			err := s.CreateEvent(storage.Event{
				ID:         fmt.Sprintf("test%d", i),
				Title:      fmt.Sprintf("titletest%d", i),
				StartTime:  startTime,
				NotifyTime: 3 * time.Hour,
			})
			require.NoError(t, err)
		}

		upcomingEvents, err := s.GetUpcomingEvents(time.Now())
		require.NoError(t, err)
		require.Len(t, upcomingEvents, 2)

		require.Equal(t, "test1", upcomingEvents[0].ID)
		require.Equal(t, "test2", upcomingEvents[1].ID)
	})
}

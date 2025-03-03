package sqlstorage

import (
	"testing"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"     //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/jmoiron/sqlx"                                                  //nolint
	"github.com/stretchr/testify/require"                                      //nolint
)

const (
	connstr     = "host=localhost user=calendar password=calendar dbname=calendardb sslmode=disable"
	testEventID = "test"
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

func TestStorage(t *testing.T) {
	t.Run("CreateEvent", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		s := New(connstr)

		title := "titletest"

		err := s.CreateEvent(storage.Event{
			ID:    testEventID,
			Title: title,
		})
		require.NoError(t, err)

		// Do check
		var events []storage.Event
		err = context.con.Select(&events, "SELECT id,title FROM event")
		require.NoError(t, err)
		require.Len(t, events, 1)

		require.Equal(t, testEventID, events[0].ID)
		require.Equal(t, title, events[0].Title)
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
}

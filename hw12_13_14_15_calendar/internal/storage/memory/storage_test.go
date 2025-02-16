package memorystorage

import (
	"fmt"
	"testing"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/stretchr/testify/require"                                      //nolint
)

func createTestEvents(s *StorageIM, times []time.Time) error {
	for i, st := range times {
		err := s.CreateEvent(storage.Event{
			ID:        fmt.Sprintf("test%d", i),
			Title:     fmt.Sprintf("titletest%d", i),
			StartTime: st,
			Duration:  time.Hour,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func checkEvents(t *testing.T, start time.Time, expectedLen int,
	getEventsFunc func(time.Time) ([]storage.Event, error),
) {
	t.Helper()
	events, err := getEventsFunc(start)
	require.NoError(t, err)
	require.Len(t, events, expectedLen)
}

const testEventID = "test"

func TestStorage(t *testing.T) {
	t.Run("CreateEvent", func(t *testing.T) {
		s := New()

		title := "titletest"

		err := s.CreateEvent(storage.Event{
			ID:    testEventID,
			Title: title,
		})
		require.NoError(t, err)

		value, ok := s.events[testEventID]
		require.True(t, ok)
		require.Equal(t, testEventID, value.ID)
		require.Equal(t, title, value.Title)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		s := New()

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
		require.Equal(t, updatedTitle, s.events[testEventID].Title)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		s := New()

		err := s.CreateEvent(storage.Event{
			ID:    testEventID,
			Title: "titletest",
		})
		require.NoError(t, err)

		err = s.DeleteEvent(storage.Event{
			ID: testEventID,
		})
		require.NoError(t, err)
		require.Empty(t, s.events)

		err = s.DeleteEvent(storage.Event{
			ID: testEventID,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotFound)
		require.Empty(t, s.events)
	})

	t.Run("GetAllEventsDay", func(t *testing.T) {
		s := New()

		times := []time.Time{
			time.Now(),
			time.Now().Add(time.Hour),
			time.Now().Add(time.Hour * 2),
			time.Now().Add(time.Hour * 3),
			time.Now().Add(time.Hour * 4),
		}

		err := createTestEvents(s, times)
		require.NoError(t, err)

		checkEvents(t, times[0], 5, s.GetAllEventsDay)
		checkEvents(t, times[1], 4, s.GetAllEventsDay)
	})

	t.Run("GetAllEventsWeek", func(t *testing.T) {
		s := New()

		times := []time.Time{
			time.Now(),
			time.Now().Add(time.Hour * 24),
			time.Now().Add(time.Hour * 24 * 5),
			time.Now().Add(time.Hour * 24 * 6),
			time.Now().Add(time.Hour * 24 * 7),
			time.Now().Add(time.Hour * 24 * 8),
		}

		err := createTestEvents(s, times)
		require.NoError(t, err)

		checkEvents(t, times[0], 4, s.GetAllEventsWeek)
		checkEvents(t, times[1], 4, s.GetAllEventsWeek)
	})

	t.Run("GetAllEventsMonth", func(t *testing.T) {
		s := New()

		times := []time.Time{
			time.Now(),
			time.Now().Add(time.Hour * 24 * 10),
			time.Now().Add(time.Hour * 24 * 20),
			time.Now().Add(time.Hour * 24 * 25),
			time.Now().Add(time.Hour * 24 * 30),
			time.Now().Add(time.Hour * 24 * 35),
		}

		err := createTestEvents(s, times)
		require.NoError(t, err)

		checkEvents(t, times[0], 4, s.GetAllEventsMonth)
		checkEvents(t, times[1], 5, s.GetAllEventsMonth)
	})
}

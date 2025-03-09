package memorystorage

import (
	"fmt"
	"testing"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"     //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/stretchr/testify/require"                                      //nolint
)

func createTestEvents(s *StorageIM, times []time.Time) error {
	for i, st := range times {
		err := s.CreateEvent(storage.Event{
			ID:         fmt.Sprintf("test%d", i),
			Title:      fmt.Sprintf("titletest%d", i),
			StartTime:  st,
			Duration:   time.Hour,
			NotifyTime: time.Hour * 2,
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

func TestStorage(t *testing.T) { //nolint:funlen
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

		err = s.DeleteEvent(testEventID)
		require.NoError(t, err)
		require.Empty(t, s.events)

		err = s.DeleteEvent(testEventID)
		require.Error(t, err)
		require.ErrorIs(t, err, app.ErrNotFound)
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

	t.Run("MarkEventAsNotificationSent", func(t *testing.T) {
		s := New()

		err := s.CreateEvent(storage.Event{
			ID:    testEventID,
			Title: "titletest",
		})
		require.NoError(t, err)

		err = s.MarkEventAsNotificationSent(testEventID)
		require.NoError(t, err)
		require.True(t, s.events[testEventID].NotificationSent)

		// Test for non-existent event
		err = s.MarkEventAsNotificationSent("nonexistent")
		require.Error(t, err)
		require.ErrorIs(t, err, app.ErrNotFound)
	})

	t.Run("DeleteEventOlderThan", func(t *testing.T) {
		s := New()

		times := []time.Time{
			time.Now(),
			time.Now().Add(time.Hour * 24 * -10),
			time.Now().Add(time.Hour * 24 * -20),
			time.Now().Add(time.Hour * 24 * -25),
			time.Now().Add(time.Hour * 24 * -30),
			time.Now().Add(time.Hour * 24 * -35),
			time.Now().Add(time.Hour * 24 * 5),
			time.Now().Add(time.Hour * 24 * 10),
		}

		err := createTestEvents(s, times)
		require.NoError(t, err)

		err = s.DeleteEventOlderThan(time.Now().Add(time.Hour * 24 * -40))
		require.NoError(t, err)
		require.Len(t, s.events, 8)

		err = s.DeleteEventOlderThan(time.Now().Add(time.Hour * 24 * -21))
		require.NoError(t, err)
		require.Len(t, s.events, 5)

		err = s.DeleteEventOlderThan(time.Now())
		require.NoError(t, err)
		require.Len(t, s.events, 2)
	})

	t.Run("GetUpcomingEvents", func(t *testing.T) {
		s := New()

		times := []time.Time{
			time.Now().Add(time.Hour * -3),
			time.Now().Add(time.Hour * -1),
			time.Now().Add(time.Hour * 1),
			time.Now().Add(time.Hour * 3),
		}

		err := createTestEvents(s, times)
		require.NoError(t, err)

		checkEvents(t, time.Now(), 1, s.GetUpcomingEvents)
		checkEvents(t, time.Now().Add(time.Hour*2), 1, s.GetUpcomingEvents)
		checkEvents(t, time.Now().Add(time.Hour*4), 0, s.GetUpcomingEvents)
	})
}

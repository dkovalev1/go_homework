package memorystorage

import (
	"sync"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"     //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
)

type StorageIM struct {
	mu sync.RWMutex

	events map[string]storage.Event
}

func (s *StorageIM) MarkEventAsNotificationSent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	event, ok := s.events[id]
	if !ok {
		return app.ErrNotFound
	}
	event.NotificationSent = true
	s.events[id] = event

	return nil
}

func (s *StorageIM) DeleteEventOlderThan(t time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id := range s.events {
		if s.events[id].StartTime.Before(t) {
			delete(s.events, id)
		}
	}

	return nil
}

func (s *StorageIM) GetUpcomingEvents(now time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var ret []storage.Event
	for _, v := range s.events {
		notify := v.StartTime.Add(-v.NotifyTime)

		if v.StartTime.After(now) &&
			notify.Before(now) {
			// but sending not noti
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func New() *StorageIM {
	return &StorageIM{
		mu:     sync.RWMutex{},
		events: make(map[string]storage.Event),
	}
}

func (s *StorageIM) CreateEvent(base storage.Event) error {
	ev := storage.Event{
		ID:          base.ID,
		Title:       base.Title,
		StartTime:   base.StartTime,
		Duration:    base.Duration,
		Description: base.Description,
		UserID:      base.UserID,
		NotifyTime:  base.NotifyTime,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[ev.ID] = ev

	return nil
}

func (s *StorageIM) UpdateEvent(ev storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[ev.ID] = ev

	return nil
}

func (s *StorageIM) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return app.ErrNotFound
	}
	delete(s.events, id)

	return nil
}

func (s *StorageIM) getAllEventsByPeriod(start, end time.Time) ([]storage.Event, error) {
	ret := make([]storage.Event, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.events {
		if v.StartTime.Equal(start) || (v.StartTime.After(start) && v.StartTime.Before(end)) {
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func (s *StorageIM) GetAllEventsDay(start time.Time) ([]storage.Event, error) {
	end := start.Add(time.Hour * 24)
	return s.getAllEventsByPeriod(start, end)
}

func (s *StorageIM) GetAllEventsWeek(start time.Time) ([]storage.Event, error) {
	end := start.Add(time.Hour * 24 * 7)
	return s.getAllEventsByPeriod(start, end)
}

func (s *StorageIM) GetAllEventsMonth(start time.Time) ([]storage.Event, error) {
	end := start.Add(time.Hour * 24 * 30)
	return s.getAllEventsByPeriod(start, end)
}

func (s *StorageIM) AddNotification(_, _ string, _ time.Time, _ int) error {
	return app.ErrNotImpl
}

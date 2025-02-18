package sqlstorage

import (
	"log"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/app"     //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/jmoiron/sqlx"                                                  //nolint
	_ "github.com/lib/pq"                                                      //nolint
)

/*
var schema = `
CREATE TABLE event (
	id          text PRIMARY KEY,
	title       string,
	starttime   time,
	duration    int, -- in seconds
	description string
	userid      int64
	NotifyTime  int -- in seconds
);
)`
*/

type StorageSQL struct { // TODO
	db *sqlx.DB
}

func New(connstr string) *StorageSQL {
	db, err := sqlx.Connect("postgres", connstr)
	if err != nil {
		log.Fatalln(err)
	}

	return &StorageSQL{
		db: db,
	}
}

func (s *StorageSQL) CreateEvent(event storage.Event) error {
	_, err := s.db.Exec("INSERT INTO event VALUES ($1, $2, $3, $4, $5, $6, $7)",
		event.ID, event.Title, event.StartTime,
		event.Duration, event.Description,
		event.UserID, event.NotifyTime)

	return err
}

func (s *StorageSQL) UpdateEvent(event storage.Event) error {
	_, err := s.db.Exec(`
UPDATE event SET 
	title=$2, 
	starttime=$3, 
	duration=$4, 
	description=$5, 
	userid=$6, 
	notifytime=$7 
WHERE id=$1`,
		event.ID, event.Title, event.StartTime, event.Duration, event.Description, event.UserID, event.NotifyTime)

	return err
}

func (s *StorageSQL) DeleteEvent(event storage.Event) error {
	_, err := s.db.Exec("DELETE FROM event WHERE id=$1", event.ID)

	return err
}

func (s *StorageSQL) GetAllEventsDay(day time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.AddDate(0, 0, 1)

	err := s.db.Select(&events, `
SELECT id, title, starttime, duration, description, userid, notifytime
FROM event
WHERE starttime >= $1 AND starttime < $2`,
		start, end)
	if err != nil {
		return nil, err
	}

	return events, err
}

func (s *StorageSQL) GetAllEventsWeek(_ time.Time) ([]storage.Event, error) {
	return nil, app.ErrNotImpl
}

func (s *StorageSQL) GetAllEventsMonth(_ time.Time) ([]storage.Event, error) {
	return nil, app.ErrNotImpl
}

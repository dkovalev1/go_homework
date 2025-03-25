package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/jackc/pgx/v5"                                                  //nolint
	. "github.com/onsi/ginkgo/v2"                                              //nolint
	. "github.com/onsi/gomega"                                                 //nolint
)

const (
	// restURL    = "http://calendar:8080/event"
	restURL    = "http://localhost:8080/event"
	eventID    = "testrest"
	eventTitle = "We will test the REST"
	connstr    = "host=localhost user=calendar password=calendar dbname=calendardb sslmode=disable"
	dsn        = "postgres://calendar:calendar@localhost:5432/calendardb?sslmode=disable"
)

func cleanupDatabase() {
	conn, err := pgx.Connect(context.Background(), dsn)
	Expect(err).NotTo(HaveOccurred())

	ctx := context.Background()

	defer conn.Close(ctx)
	// Ensure the connection is actually working
	err = conn.Ping(ctx)
	Expect(err).NotTo(HaveOccurred())

	_, err = conn.Exec(ctx, "DELETE FROM event")
	Expect(err).NotTo(HaveOccurred())

	_, err = conn.Exec(ctx, "DELETE FROM notification")
	Expect(err).NotTo(HaveOccurred())
}

func getDBEventIDs() []string {
	conn, err := pgx.Connect(context.Background(), dsn)
	Expect(err).NotTo(HaveOccurred())

	ctx := context.Background()

	defer conn.Close(ctx)
	// Ensure the connection is actually working
	err = conn.Ping(ctx)
	Expect(err).NotTo(HaveOccurred())

	rows, err := conn.Query(ctx, "SELECT id FROM event")
	Expect(err).NotTo(HaveOccurred())
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		Expect(err).NotTo(HaveOccurred())
		ids = append(ids, id)
	}
	return ids
}

type Notification struct {
	eventID string
	title   string
	start   time.Time
	user    int64
}

func getDBNotifications() (ret []Notification) {
	conn, err := pgx.Connect(context.Background(), dsn)
	Expect(err).NotTo(HaveOccurred())

	ctx := context.Background()

	defer conn.Close(ctx)
	// Ensure the connection is actually working
	err = conn.Ping(ctx)
	Expect(err).NotTo(HaveOccurred())

	rows, err := conn.Query(ctx, "SELECT eventId, title, start, userid FROM notification")

	Expect(err).NotTo(HaveOccurred())
	defer rows.Close()

	for rows.Next() {
		var n Notification
		err := rows.Scan(&n.eventID, &n.title, &n.start, &n.user)
		Expect(err).NotTo(HaveOccurred())
		ret = append(ret, n)
	}

	return
}

func waitForNotifications() bool {
	for {
		notifications := getDBNotifications()

		if len(notifications) > 0 {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
}

type Event struct {
	ID               string        `db:"id"`
	Title            string        `db:"title"`
	StartTime        time.Time     `db:"starttime"`
	Duration         time.Duration `db:"duration"`
	Description      string        `db:"description"`
	UserID           int64         `db:"userid"`
	NotifyTime       time.Duration `db:"notifytime"`
	NotificationSent bool          `db:"notification_sent"`
}

func makeTestEvent(eventID, eventTitle string, eventStart time.Time) []byte {
	event := Event{
		ID:         eventID,
		Title:      eventTitle,
		StartTime:  eventStart,
		NotifyTime: time.Hour * 3,
	}

	b, err := json.Marshal(event)
	Expect(err).NotTo(HaveOccurred())
	return b
}

func testCreateEvent(eventStart time.Time, shouldFail bool) {
	client := &http.Client{}
	b := makeTestEvent(eventID, eventTitle, eventStart)
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "PUT", restURL, bytes.NewBuffer(b))
	Expect(err).NotTo(HaveOccurred())

	resp, err := client.Do(req)
	Expect(err).NotTo(HaveOccurred())

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())

	bodystr := string(body)
	if shouldFail {
		Expect(bodystr).To(Equal("pq: duplicate key value violates unique constraint \"event_pkey\""))
		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	} else {
		Expect(bodystr).To(Equal("ok\n"))
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(resp.Status).To(Equal("200 OK"))
	}
}

func testCreateEvents(eventStart time.Time, count int) {
	client := &http.Client{}

	dt := 3*time.Hour + 2*time.Minute

	ctx := context.Background()

	start := eventStart
	actual := 0

	for i := range count {
		start = start.Add(dt)
		b := makeTestEvent(fmt.Sprintf("%s_%d", eventID, i), fmt.Sprintf("%s_%d", eventTitle, i), start)

		req, err := http.NewRequestWithContext(ctx, "PUT", restURL, bytes.NewBuffer(b))
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(resp.Status).To(Equal("200 OK"))

		body, err := io.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(Equal([]byte("ok\n")))
		actual++
	}
	Expect(actual).To(Equal(count))
}

func testGetEvents(interval string) []storage.Event {
	req, err := http.NewRequestWithContext(context.Background(), "GET", restURL+"/"+interval, nil)
	Expect(err).NotTo(HaveOccurred())

	client := &http.Client{}
	resp, err := client.Do(req)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	Expect(resp.Status).To(Equal("200 OK"))
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(err).NotTo(HaveOccurred())

	var events []storage.Event
	err = json.Unmarshal(body, &events)
	Expect(err).NotTo(HaveOccurred())

	return events
}

var _ = Describe("Integration tests", func() {
	now := time.Now()
	bod := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	BeforeEach(func() {
		cleanupDatabase()
	})

	AfterEach(func() {
		cleanupDatabase()
	})

	Describe("Sending event", func() {
		It("can accept event", func(_ SpecContext) {
			testCreateEvent(time.Now().Add(time.Hour), false)

			ids := getDBEventIDs()
			Expect(ids).To(ContainElement(eventID))
		}, SpecTimeout(time.Second*20))

		It("Reject duplicate event", func(_ SpecContext) {
			testCreateEvent(time.Now().Add(time.Hour), false)
			testCreateEvent(time.Now().Add(time.Hour*2), true)

			ids := getDBEventIDs()
			Expect(len(ids)).Should(Equal(1))
			Expect(ids).To(ContainElement(eventID))
		}, SpecTimeout(time.Second*20))
	})

	Describe("get events", func() {
		It("gets events for 1 day", func(_ SpecContext) {
			testCreateEvents(bod, 35)

			ids := getDBEventIDs()
			Expect(len(ids)).Should(Equal(35))

			// Check that events are in database and can be received by the REST
			events := testGetEvents("day")

			Expect(len(events)).Should(Equal(8))
		}, SpecTimeout(time.Second*20))

		It("gets events for 1 week", func(_ SpecContext) {
			testCreateEvents(bod, 60)

			ids := getDBEventIDs()

			Expect(len(ids)).Should(Equal(60))

			// Check that events are in database and can be received by the REST
			events := testGetEvents("week")

			Expect(len(events)).Should(Equal(56))
		}, SpecTimeout(time.Second*300))

		It("gets events for 1 month", func(_ SpecContext) {
			testCreateEvents(bod, 300)

			ids := getDBEventIDs()

			Expect(len(ids)).Should(Equal(300))

			// Check that events are in database and can be received by the REST
			events := testGetEvents("month")

			Expect(len(events)).Should(Equal(246))
		}, SpecTimeout(time.Second*300))
	})

	Describe("sender", func() {
		It("can send notification", func(_ SpecContext) {
			testCreateEvents(bod, 10)

			ids := getDBEventIDs()

			Expect(len(ids)).Should(Equal(10))

			// Check that events are in database and can be received by the REST
			Expect(waitForNotifications()).To(BeTrue())
		}, SpecTimeout(time.Second*30))
	})
})

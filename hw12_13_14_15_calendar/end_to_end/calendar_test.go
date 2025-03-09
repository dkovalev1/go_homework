package end_to_end_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	//nolint
	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/gomega"
)

var (
	// calendarCmd *exec.Cmd
	calendar  = "./bin/calendar"
	scheduler = "./bin/scheduler"
	sender    = "./bin/sender"
)

const (
	testPort   = 8080
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

func testCreateEvent(eventStart time.Time) {
	client := &http.Client{}
	b := makeTestEvent(eventID, eventTitle, eventStart)
	url := fmt.Sprintf("http://localhost:%d/event", testPort)

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(b))
	Expect(err).NotTo(HaveOccurred())

	resp, err := client.Do(req)
	Expect(err).NotTo(HaveOccurred())

	defer resp.Body.Close()

	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(resp.Status).To(Equal("200 OK"))

	body, err := io.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(body).To(Equal([]byte("ok\n")))
}

func checkStart(program, config, waitFor string) *gexec.Session {
	cmd := exec.Command(program, config)
	cmd.Dir = ".."
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gbytes.Say(waitFor))
	return session
}

func checkStop(session *gexec.Session) {
	session.Interrupt()
	Eventually(session).Should(gexec.Exit(0))
}

var _ = Describe("End to End plugin tests", func() {
	Describe("Start", func() {
		It("should start calendar", func(ctx SpecContext) {
			session := checkStart(calendar, "--config=configs/config.toml", `\[INFO\] starting grpc server on \[::\]:8081`)
			time.Sleep(time.Millisecond * 200)
			checkStop(session)
		}, SpecTimeout(time.Second*3))

		It("should start scheduler", func(ctx SpecContext) {
			session := checkStart(scheduler, "--config=configs/scheduler.toml", `\[INFO\] Scheduler started`)
			time.Sleep(time.Millisecond * 200)
			checkStop(session)
		}, SpecTimeout(time.Second*3))

		It("should start sender", func(ctx SpecContext) {
			session := checkStart(sender, "--config=configs/sender.toml", `\[INFO\] Sender started`)
			time.Sleep(time.Millisecond * 200)
			checkStop(session)
		}, SpecTimeout(time.Second*3))
	})

	Describe("Sending event", func() {
		var calendarSession *gexec.Session
		var schedulerSession *gexec.Session
		var senderSession *gexec.Session

		BeforeEach(func(ctx SpecContext) {

			cleanupDatabase()

			calendarSession = checkStart(calendar, "--config=configs/config.toml", `\[INFO\] starting grpc server on \[::\]:8081`)
			schedulerSession = checkStart(scheduler, "--config=configs/scheduler.toml", `\[INFO\] Scheduler started`)
			senderSession = checkStart(sender, "--config=configs/sender.toml", `\[INFO\] Sender started`)
		})
		AfterEach(func(ctx SpecContext) {
			checkStop(calendarSession)
			checkStop(senderSession)
			checkStop(schedulerSession)

			cleanupDatabase()
		})

		It("can start all services", func(ctx SpecContext) {
			Expect(calendarSession).NotTo(BeNil())
			Expect(schedulerSession).NotTo(BeNil())
			Expect(senderSession).NotTo(BeNil())
		}, SpecTimeout(time.Second*3))

		It("can accept event", func(ctx SpecContext) {
			testCreateEvent(time.Now().Add(time.Hour))

			Eventually(calendarSession).Should(gbytes.Say(`PUT /event`))
			Eventually(schedulerSession).WithTimeout(time.Second * 10).Should(gbytes.Say(`\[INFO\] Sent a message: testrest`))
			Eventually(senderSession).WithTimeout(time.Second * 10).Should(gbytes.Say(`\[INFO\] Received a message: testrest`))
		}, SpecTimeout(time.Second*20))
	})
})

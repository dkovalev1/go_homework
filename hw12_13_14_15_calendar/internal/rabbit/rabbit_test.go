package rabbit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger" //nolint
	"github.com/google/uuid"                                                  //nolint
	"github.com/streadway/amqp"                                               //nolint
	"github.com/stretchr/testify/require"                                     //nolint
)

const connstr = "amqp://calendar:calendarmq@localhost:5672/"

type testContext struct {
	logger *logger.Logger

	conn *amqp.Connection
	ch   *amqp.Channel

	// Not used yet
	// queue *amqp.Queue
}

func (c *testContext) cleanQueue() {
	c.ch.QueuePurge(QueueName, false)
}

func setUp(t *testing.T) *testContext {
	t.Helper()

	ret := &testContext{
		logger: logger.New("debug"),
	}

	// Open a channel and declare a queue
	conn, err := amqp.Dial(connstr) // "amqp://calendar:calendarmq@localhost:5672/"
	require.NoError(t, err)

	ret.conn = conn

	channel, err := conn.Channel()
	require.NoError(t, err)

	ret.ch = channel

	ret.cleanQueue()

	return ret
}

func tearDown(t *testing.T, ctx *testContext) {
	t.Helper()
	ctx.cleanQueue()
	ctx.ch.Close()
	ctx.conn.Close()
}

func TestRabbit(t *testing.T) {
	t.Run("Connect", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		r := NewRabbit(connstr, context.logger)
		defer r.Close()

		require.NotNil(t, r)
		rabbit := r.(*Rabbit)

		require.NotNil(t, rabbit.conn)
		require.NotNil(t, rabbit.ch)
		require.NotNil(t, rabbit.queue)

		require.Equal(t, QueueName, rabbit.queue.Name)
	})

	t.Run("SendNotification", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		r := NewRabbit(connstr, context.logger)
		defer r.Close()

		require.NotNil(t, r)

		evid := uuid.New().String()
		evTime := time.Now()
		notification := &Notification{
			EventID: evid,
			Title:   "event1",
			Time:    evTime,
			User:    1,
		}

		err := r.SendNotification(notification)
		require.NoError(t, err)

		rabbit := r.(*Rabbit)
		require.NotNil(t, rabbit.conn)

		msgs, err := context.ch.Consume(QueueName, "checker", false, false, false, true, nil)
		require.NoError(t, err)

		timer := time.NewTimer(time.Second * 2)

		select {
		case <-timer.C:
			require.Fail(t, "No message in the queue")
		case x, ok := <-msgs:
			require.True(t, ok)
			require.NotEmpty(t, x.Body)

			var recv Notification
			require.NoError(t, json.Unmarshal(x.Body, &recv))

			require.Equal(t, notification.EventID, recv.EventID)
			require.Equal(t, notification.Title, recv.Title)

			// JSON could lost precision, so approx compare
			dt := recv.Time.Sub(notification.Time)
			require.LessOrEqual(t, dt, time.Millisecond)

			require.Equal(t, notification.User, recv.User)
		}
	})

	t.Run("ReceiveNotification", func(t *testing.T) {
		context := setUp(t)
		defer func() {
			tearDown(t, context)
		}()

		r := NewRabbit(connstr, context.logger)
		defer r.Close()

		require.NotNil(t, r)
		rabbit := r.(*Rabbit)
		require.NotNil(t, rabbit.conn)

		evid := uuid.New().String()
		evTime := time.Now()

		notification := &Notification{
			EventID: evid,
			Title:   "event2",
			Time:    evTime,
			User:    2,
		}

		body, err := notification.body()
		require.NoError(t, err)

		queue, err := context.ch.QueueInspect(QueueName)
		require.NoError(t, err)
		require.NotNil(t, queue)

		require.Equal(t, QueueName, queue.Name)
		require.Equal(t, 0, queue.Messages)
		require.Equal(t, 0, queue.Consumers)

		err = context.ch.Publish(
			"",        // exchange)
			QueueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/json",
				Body:        body,
			})
		require.NoError(t, err)

		finish := make(chan struct{})

		notificationChan := make(chan *Notification, 1)

		go func() {
			r.ReceiveNotifications(finish, func(recv *Notification) {
				require.NotNil(t, recv)
				require.Equal(t, notification.EventID, recv.EventID)
				require.Equal(t, notification.Title, recv.Title)

				// JSON could lost precision, so approx compare
				dt := recv.Time.Sub(notification.Time)
				require.LessOrEqual(t, dt, time.Millisecond)
				require.Equal(t, notification.User, recv.User)

				notificationChan <- recv
				finish <- struct{}{}
			})
		}()

		hasNotification := false

		select {
		case recv := <-notificationChan:
			require.NotNil(t, recv)
			require.Equal(t, notification.EventID, recv.EventID)
			hasNotification = true
		case <-finish:
		case <-time.After(time.Second * 20):
			require.Fail(t, "Timeout closing receiver")
		}

		require.True(t, hasNotification)
	})
}

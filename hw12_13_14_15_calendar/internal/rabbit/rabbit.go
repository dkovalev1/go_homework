package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	//nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/logger"  //nolint
	"github.com/dkovalev1/go_homework/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/streadway/amqp"                                                //nolint
)

const (
	QueueName = "calendar-queue"
)

type IRabbit interface {
	Close() error
	SendNotification(*Notification) error
	ReceiveNotifications(interruptChan <-chan struct{}, handler func(*Notification)) error
}

type Rabbit struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	queue  *amqp.Queue
	logger *logger.Logger
}

type Notification struct {
	EventID string
	Title   string
	Time    time.Time
	User    int64
}

func (n *Notification) body() ([]byte, error) {
	b, err := json.Marshal(n)
	return b, err
}

func NewNotification(event storage.Event) *Notification {
	return &Notification{
		EventID: event.ID,
		Title:   event.Description,
		Time:    event.StartTime,
		User:    event.UserID,
	}
}

func NewRabbit(connstr string, logger *logger.Logger) IRabbit {
	ret := &Rabbit{
		logger: logger,
	}
	if err := ret.connect(connstr); err != nil {
		log.Fatal(err)
	}
	return ret
}

func (r *Rabbit) connect(connstr string) error {
	conn, err := amqp.Dial(connstr) // "amqp://calendar:calendarmq@localhost:5672/"
	if err != nil {
		r.logger.Error("Failed to connect to RabbitMQ: " + err.Error())
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return err
	}

	// We create a Queue to send the message to.
	q, err := ch.QueueDeclare(
		QueueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	r.queue = &q
	r.ch = ch
	r.conn = conn

	return nil
}

func (r *Rabbit) SendNotification(notification *Notification) error {
	// We set the payload for the message.
	body, err := notification.body()
	if err != nil {
		return err
	}
	err = r.ch.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        body,
		})
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to publish a message: %v", err))
	}

	return err
}

func (r *Rabbit) ReceiveNotifications(interruptChan <-chan struct{}, handler func(*Notification)) error {
	msgs, err := r.ch.Consume(r.queue.Name, "calendar", true, false, false, false, nil)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to register a consumer: %v", err))
		return err
	}

	for {
		select {
		case <-interruptChan:
			return nil
		case msg := <-msgs:
			var notification Notification
			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				r.logger.Error(err.Error())
				continue
			}
			handler(&notification)
		}
	}
}

func (r *Rabbit) Close() error {
	r.ch.Close()
	return r.conn.Close()
}

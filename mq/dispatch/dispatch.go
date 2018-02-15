package dispatch

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

// A Dispatcher dispatches messages to the message queue
type Dispatcher interface {
	DispatchMessage(message interface{}) (err error)
}

// AMQPDispatcher implements the Dispatcher interface
type AMQPDispatcher struct {
	channel       *amqp.Channel
	queueName     string
	mandatorySend bool
}

// DispatchMessage function of AMQPDispatcher
func (d *AMQPDispatcher) DispatchMessage(message interface{}) (err error) {
	log.Println("Dispatching message to queue", d.queueName)
	body, err := json.Marshal(message)
	if err == nil {
		err = d.channel.Publish(
			"",              // exchange
			d.queueName,     // routing key
			d.mandatorySend, // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Println("Failed to dispatch message, err", err)
		}
	} else {
		log.Println("Failed to marshal:", err, "message", message)
	}
	return
}

// NewDispatcher returns a new Dispatcher for the given connection and queue
func NewDispatcher(conn *amqp.Connection, queueName string) Dispatcher {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("Failed to declare a queue", err)
	}

	return &AMQPDispatcher{
		channel:       ch,
		queueName:     q.Name,
		mandatorySend: false,
	}
}

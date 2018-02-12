package mq_helper

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/seagullbird/headr-common/config"
	"reflect"
)

type Receiver interface {
	RegisterListener(queueName string, listener Listener)
}

type AMQPReceiver struct {
	ch 				*amqp.Channel
	registration	map[string]Listener
}

// Listener is a function that takes action when an event is received.
type Listener func(delivery amqp.Delivery)

func (r *AMQPReceiver) RegisterListener(queueName string, listener Listener) {
	if l, ok := r.registration[queueName]; ok {
		log.Println("Listener already registered", reflect.ValueOf(l))
		return
	}
	r.registration[queueName] = listener
	log.Println("New Listener registered, queue", queueName)

	q, _ := r.ch.QueueDeclare(
		queueName, 		// name
		false,          // durable
		false,		// delete when usused
		false,			// exclusive
		false,			// no-wait
		nil,				// arguments
	)
	qIn, _ := r.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	// Start listening
	go func() {
		for d := range qIn {
			listener(d)
		}
	}()
}

func NewReceiver() Receiver {
	uri := amqp.URI{
		Scheme:   "amqp",
		Host:     config.MQSERVERNAME,
		Port:     5672,
		Username: "user",
		Password: config.MQSERVERPWD,
		Vhost:    "/",
	}
	conn, err := amqp.Dial(uri.String())

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open AMQP channel", err)
	}

	return &AMQPReceiver{
		ch: ch,
		registration: make(map[string]Listener),
	}
}

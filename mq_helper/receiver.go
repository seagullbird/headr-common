package mq_helper

import (
	"github.com/go-kit/kit/log"
	"github.com/streadway/amqp"
	"github.com/seagullbird/headr-common/config"
	"reflect"
)

type Receiver interface {
	RegisterListener(queueName string, event interface{}, listener Listener)
}

type AMQPReceiver struct {
	ch 				*amqp.Channel
	logger 			log.Logger
	registration	map[string]Listener
}

// Listener is a function that takes action when an event is received.
type Listener func(delivery amqp.Delivery, logger log.Logger)

func (r *AMQPReceiver) RegisterListener(queueName string, listener Listener) {
	if l, ok := r.registration[queueName]; ok {
		r.logger.Log("Listener already registered", reflect.ValueOf(l))
		return
	}
	r.registration[queueName] = listener

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
			listener(d, r.logger)
		}
	}()
}

func NewReceiver(logger log.Logger) *AMQPReceiver {
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
		logger.Log("Failed to open AMQP channel", err)
	}

	return &AMQPReceiver{
		ch: ch,
		logger:	logger,
		registration: make(map[string]Listener),
	}
}
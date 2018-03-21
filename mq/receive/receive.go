package receive

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/seagullbird/headr-common/mq/client"
	"github.com/streadway/amqp"
	"sync"
)

// A Receiver receives Messages from the message queue,
// and consumes them with registered listener functions
type Receiver interface {
	RegisterListener(queueName string, listener Listener) error
}

// ErrQueueAlreadyRegistered makes sure only each queue has only one Listener
var ErrQueueAlreadyRegistered = errors.New("this queue is already registered by another Listener")

// AMQPReceiver implements the Receiver interface
type AMQPReceiver struct {
	client.Client
	ch           *amqp.Channel
	registration map[string]Listener
	logger       log.Logger
	// mux is for mutual exclusion of listener goroutines
	mux sync.Mutex
}

// Listener is a function that takes action when an event is received.
type Listener func(delivery amqp.Delivery)

// RegisterListener register one Listener for the given queue and start a goroutine listening that queue,
// each arrived message from that queue will be consumed by the registered Listener,
// each consuming is mutual exclusive
func (r *AMQPReceiver) RegisterListener(queueName string, listener Listener) error {
	if _, ok := r.registration[queueName]; ok {
		r.logger.Log("errror_desc", "queue already registered", "queue_name", queueName)
		return ErrQueueAlreadyRegistered
	}
	r.registration[queueName] = listener
	r.logger.Log("info", "New Listener registered", "queue_name", queueName)

	q, _ := r.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
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
			r.mux.Lock()
			listener(d)
			r.mux.Unlock()
		}
	}()
	return nil
}

// NewReceiver returns a new Receiver for the given connection
func NewReceiver(c client.Client, logger log.Logger) (Receiver, error) {
	if err := c.Connect(); err != nil {
		logger.Log("errror_desc", "Failed to connect to RabbitMQ", "error", err)
		return nil, err
	}
	ch, err := c.Connection().Channel()
	if err != nil {
		logger.Log("error_desc", "Failed to open a channel", "error", err)
		return nil, err
	}

	receiver := &AMQPReceiver{
		Client:       c,
		ch:           ch,
		registration: make(map[string]Listener),
		logger:       logger,
	}

	go func() {
		for {
			if <-receiver.Connection().NotifyClose(make(chan *amqp.Error)) != nil {
				receiver.logger.Log("[warning]", fmt.Sprintf("MQ connection closing: %v", err))
				// Connection broken, try reconnecting
				for {
					retryTime := 1
					if err := receiver.Reconnect(retryTime); err != nil {
						receiver.logger.Log("error_desc", "Reconnection failed", "retrytime", retryTime)
						retryTime++
					} else {
						receiver.logger.Log("[info]", "Reconnection succeeded")
						receiver.ch, err = receiver.Connection().Channel()
						if err != nil {
							receiver.logger.Log("error_desc", "Failed to open a channel", "error", err)
						}
						// TODO: deal with reconnecting success but open channel fail???
						break
					}
				}
			}
		}
	}()

	return receiver, nil
}

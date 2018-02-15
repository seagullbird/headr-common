package receive

import (
	"github.com/streadway/amqp"
	"log"
	"reflect"
	"sync"
)

// A Receiver receives Messages from the message queue,
// and consumes them with registered listener functions
type Receiver interface {
	RegisterListener(queueName string, listener Listener)
}

// AMQPReceiver implements the Receiver interface
type AMQPReceiver struct {
	ch           *amqp.Channel
	registration map[string]Listener
	// mux is for mutual exclusion of listener goroutines
	mux sync.Mutex
}

// Listener is a function that takes action when an event is received.
type Listener func(delivery amqp.Delivery)

// RegisterListener register one Listener for the given queue and start a goroutine listening that queue,
// each arrived message from that queue will be consumed by the registered Listener,
// each consuming is mutual exclusive
func (r *AMQPReceiver) RegisterListener(queueName string, listener Listener) {
	if l, ok := r.registration[queueName]; ok {
		log.Println("Listener already registered", reflect.ValueOf(l))
		return
	}
	r.registration[queueName] = listener
	log.Println("New Listener registered, queue", queueName)

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
}

// NewReceiver returns a new Receiver for the given connection
func NewReceiver(conn *amqp.Connection) Receiver {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open AMQP channel", err)
	}

	return &AMQPReceiver{
		ch:           ch,
		registration: make(map[string]Listener),
	}
}

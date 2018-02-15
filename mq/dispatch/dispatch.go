package dispatch

import (
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/streadway/amqp"
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
	logger        log.Logger
}

// DispatchMessage function of AMQPDispatcher
func (d *AMQPDispatcher) DispatchMessage(message interface{}) (err error) {
	d.logger.Log("info", "Dispatching message to queue", "queue_name", d.queueName)
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
			d.logger.Log("error_desc", "Failed to dispatch message", "error", err)
		}
	} else {
		d.logger.Log("error_desc", "Failed to marshal", "error", err, "message", message)
	}
	return
}

// NewDispatcher returns a new Dispatcher for the given connection and queue
func NewDispatcher(conn *amqp.Connection, queueName string, logger log.Logger) (Dispatcher, error) {
	ch, err := conn.Channel()
	if err != nil {
		logger.Log("error_desc", "Failed to open a channel", "error", err)
		return nil, err
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
		logger.Log("error_desc", "Failed to declare a queue", "error", err)
		return nil, err
	}

	return &AMQPDispatcher{
		channel:       ch,
		queueName:     q.Name,
		mandatorySend: false,
		logger:        logger,
	}, nil
}

package dispatch

//go:generate mockgen -destination=./mock/mock_dispatch.go -package=mock github.com/seagullbird/headr-common/mq/dispatch Dispatcher

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/seagullbird/headr-common/mq/client"
	"github.com/streadway/amqp"
)

// A Dispatcher dispatches messages to the message queue
type Dispatcher interface {
	DispatchMessage(queueName string, message interface{}) (err error)
}

// AMQPDispatcher implements the Dispatcher interface
type AMQPDispatcher struct {
	client.Client
	channel       *amqp.Channel
	mandatorySend bool
	logger        log.Logger
}

// DispatchMessage function of AMQPDispatcher
func (d *AMQPDispatcher) DispatchMessage(queueName string, message interface{}) (err error) {
	d.logger.Log("info", "Dispatching message to queue", "queue_name", queueName)
	body, err := json.Marshal(message)
	if err != nil {
		d.logger.Log("error_desc", "Failed to marshal", "error", err, "message", message)
		return
	}
	_, err = d.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		d.logger.Log("error_desc", "Failed to declare a queue", "error", err)
		return err
	}

	err = d.channel.Publish(
		"",              // exchange
		queueName,       // routing key
		d.mandatorySend, // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		d.logger.Log("error_desc", "Failed to dispatch message", "error", err)
	}

	return
}

// NewDispatcher returns a new Dispatcher for the given connection and queue
func NewDispatcher(c client.Client, logger log.Logger) (Dispatcher, error) {
	if err := c.Connect(); err != nil {
		logger.Log("errror_desc", "Failed to connect to RabbitMQ", "error", err)
		return nil, err
	}
	ch, err := c.Connection().Channel()
	if err != nil {
		logger.Log("error_desc", "Failed to open a channel", "error", err)
		return nil, err
	}

	dispatcher := &AMQPDispatcher{
		Client:        c,
		channel:       ch,
		mandatorySend: false,
		logger:        logger,
	}

	go func() {
		for {
			if err := <-dispatcher.Connection().NotifyClose(make(chan *amqp.Error)); err != nil {
				dispatcher.logger.Log("[warning]", fmt.Sprintf("MQ connection closing: %v", err))
				// Connection broken, try reconnecting
				for {
					retryTime := 1
					if err := dispatcher.Reconnect(retryTime); err != nil {
						dispatcher.logger.Log("error_desc", "Reconnection failed", "retrytime", retryTime)
						retryTime++
					} else {
						dispatcher.logger.Log("[info]", "Reconnection succeeded")
						dispatcher.channel, err = dispatcher.Connection().Channel()
						if err != nil {
							dispatcher.logger.Log("error_desc", "Failed to open a channel", "error", err)
						}
						// TODO: deal with reconnecting success but open channel fail???
						break
					}
				}
			}
		}
	}()

	return dispatcher, nil
}

// FakeDispatcher is Only used in tests
type FakeDispatcher struct{}

// DispatchMessage function of FakeDispatcher
func (d FakeDispatcher) DispatchMessage(queueName string, message interface{}) (err error) {
	return nil
}

package mq_helper

import (
	"github.com/streadway/amqp"
	"github.com/go-kit/kit/log"
	"encoding/json"
	"github.com/seagullbird/headr-common/config"
)

type Dispatcher interface {
	DispatchMessage(message interface{}) (err error)
}

type AMQPDispatcher struct {
	channel       	*amqp.Channel
	queueName     	string
	mandatorySend 	bool
	logger 			log.Logger
}

func (d *AMQPDispatcher) DispatchMessage(message interface{}) (err error) {
	d.logger.Log("Dispatching message to queue", d.queueName)
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
			d.logger.Log("Failed to dispatch message, err", err)
		}
	} else {
		d.logger.Log("Failed to marshal err", err, "message", message)
	}
	return
}

func NewDispatcher(queueName string, logger log.Logger) Dispatcher {
	uri := amqp.URI{
		Scheme:   "amqp",
		Host:     config.MQSERVERNAME,
		Port:     5672,
		Username: "user",
		Password: config.MQSERVERPWD,
		Vhost:    "/",
	}
	conn, err := amqp.Dial(uri.String())
	if err != nil {
		logger.Log("Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Log( "Failed to open a channel", err)
	}

	q, err := ch.QueueDeclare(
		queueName, 			// name
		false,		// durable
		false,	// delete when unused
		false,		// exclusive
		false,		// no-wait
		nil,			// arguments
	)
	if err != nil {
		logger.Log( "Failed to declare a queue", err)
	}

	return &AMQPDispatcher{
		channel: ch,
		queueName: q.Name,
		mandatorySend: false,
		logger: logger,
	}
}

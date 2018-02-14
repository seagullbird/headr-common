package mq_helper

import (
	"github.com/seagullbird/headr-common/config"
	"github.com/streadway/amqp"
)

func MakeConn() (*amqp.Connection, error) {
	uri := amqp.URI{
		Scheme:   "amqp",
		Host:     config.MQSERVERNAME,
		Port:     5672,
		Username: "user",
		Password: config.MQSERVERPWD,
		Vhost:    "/",
	}
	return amqp.Dial(uri.String())
}

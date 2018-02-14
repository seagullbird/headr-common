package mq

import (
	"github.com/streadway/amqp"
)

func MakeConn(servername, username, passwd string) (*amqp.Connection, error) {
	uri := amqp.URI{
		Scheme:   "amqp",
		Host:     servername,
		Port:     5672,
		Username: username,
		Password: passwd,
		Vhost:    "/",
	}
	return amqp.Dial(uri.String())
}

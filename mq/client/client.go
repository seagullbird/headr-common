package client

import (
	"github.com/streadway/amqp"
	"math/rand"
	"time"
)

// Client represents a basic rabbitmq client that connects and closes.
type Client interface {
	Connect() error
	Reconnect(retryTime int) error
	Close()
	Connection() *amqp.Connection
}

type rabbitmqClient struct {
	servername string
	username   string
	password   string
	conn       *amqp.Connection
}

// New returns a rabbitmqClient instance.
func New(servername, username, password string) Client {
	return &rabbitmqClient{
		servername: servername,
		username:   username,
		password:   password,
	}
}

func (c *rabbitmqClient) Connect() error {
	uri := amqp.URI{
		Scheme:   "amqp",
		Host:     c.servername,
		Port:     5672,
		Username: c.username,
		Password: c.password,
		Vhost:    "/",
	}
	var err error
	c.conn, err = amqp.Dial(uri.String())
	return err
}

func (c *rabbitmqClient) Reconnect(retryTime int) error {
	c.Close()
	time.Sleep(time.Duration(15+rand.Intn(60)+2*retryTime) * time.Second)
	return c.Connect()
}

func (c *rabbitmqClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *rabbitmqClient) Connection() *amqp.Connection {
	return c.conn
}

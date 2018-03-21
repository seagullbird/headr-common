package client

import (
	"github.com/streadway/amqp"
	"os"
	"testing"
)

func TestRabbitmqClient(t *testing.T) {
	var (
		servername = os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
		username   = "guest"
		passwd     = "guest"
	)
	c := New(servername, username, passwd)
	if err := c.Connect(); err != nil {
		t.Fatal(err)
	}

	receiver := make(chan *amqp.Error)
	c.Connection().NotifyClose(receiver)
	// Notice: c.Close() won't send any error into receiver coz this close is treat as a wanted action.
	// Only unexpected connection broken error will be received at receiver. There is a difference here.
	// I am calling c.NotifyClose() here only to test the function works well.
	c.Close()
	if err := c.Reconnect(1); err != nil {
		t.Fatal(err)
	}
}

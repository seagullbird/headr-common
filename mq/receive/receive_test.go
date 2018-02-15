package receive_test

import (
	"encoding/json"
	"fmt"
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/streadway/amqp"
	"os"
	"testing"
)

func TestReceiveMessage(t *testing.T) {
	var (
		servername = os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
		username   = "guest"
		passwd     = "guest"
	)
	conn, err := mq.MakeConn(servername, username, passwd)
	if err != nil {
		t.Fatal("Cannot connection to RabbitMQ", err)
	}

	receiver, err := receive.NewReceiver(conn)
	if err != nil {
		t.Fatal("Failed to create receiver", err)
	}

	err = receiver.RegisterListener("dispatch_test", func(delivery amqp.Delivery) {
		var event mq.ExampleEvent
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			panic(err)
		}
		fmt.Printf("Received new event: %s", event)
	})

	if err != nil {
		t.Fatal("Failed to register Listener to queue dispatch_test", err)
	}

	// test duplicate registration
	err = receiver.RegisterListener("dispatch_test", func(delivery amqp.Delivery) {
		fmt.Printf("I should not be run.")
	})
	if err != receive.ErrQueueAlreadyRegistered {
		t.Fatal("Duplicate Listener to the same queue", err)
	}
}

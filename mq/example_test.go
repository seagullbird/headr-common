package mq_test

import (
	"encoding/json"
	"fmt"
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/dispatch"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/streadway/amqp"
	"os"
	"time"
)

// Example event
type exampleTestEvent struct {
	name string `json:"name"`
}

func (e exampleTestEvent) String() string {
	return fmt.Sprintf("exampleTestEvent, name=%s", e.name)
}

// Example listener
func exampleListener(delivery amqp.Delivery) {
	var event exampleTestEvent
	if err := json.Unmarshal(delivery.Body, &event); err != nil {
		panic(err)
	}
	fmt.Printf("Received new event: %s", event)
}

func Example() {
	// Make connection to rabbitmq server
	var (
		servername = os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
		username   = "guest"
		passwd     = "guest"
	)
	conn, err := mq.MakeConn(servername, username, passwd)
	if err != nil {
		panic(err)
	}
	// New dispatcher
	dispatcher := dispatch.NewDispatcher(conn, "example_test")
	// New receiver
	receiver := receive.NewReceiver(conn)

	// Dispatch a message
	msg := exampleTestEvent{
		name: "example-message",
	}
	err = dispatcher.DispatchMessage(msg)
	if err != nil {
		panic(err)
	}

	// Register a message listener
	receiver.RegisterListener("example_test", exampleListener)

	// Wait for the message to be consumed
	time.Sleep(3 * time.Second)

	// Output:
	// Received new event: exampleTestEvent, name=example-message
}

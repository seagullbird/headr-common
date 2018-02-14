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
type ExampleTestEvent struct {
	Message string `json:"Message"`
}

func (e ExampleTestEvent) String() string {
	return fmt.Sprintf("ExampleTestEvent, Message=%s", e.Message)
}

// Example listener
func exampleListener(delivery amqp.Delivery) {
	var event ExampleTestEvent
	if err := json.Unmarshal(delivery.Body, &event); err != nil {
		panic(err)
	}
	fmt.Printf("Received new event: %s", event)
}

func Example() {
	var (
		servername = os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
		username   = "guest"
		passwd     = "guest"
	)

	// New dispatcher
	// Make connection to rabbitmq server
	dConn, err := mq.MakeConn(servername, username, passwd)
	if err != nil {
		panic(err)
	}
	dispatcher := dispatch.NewDispatcher(dConn, "example_test")

	// Dispatch a Message
	msg := ExampleTestEvent{
		Message: "example-message",
	}
	err = dispatcher.DispatchMessage(msg)
	if err != nil {
		panic(err)
	}

	// Wait for the Message to be produced
	time.Sleep(time.Second)

	// New receiver
	// Make connection to rabbitmq server
	rConn, err := mq.MakeConn(servername, username, passwd)
	if err != nil {
		panic(err)
	}
	receiver := receive.NewReceiver(rConn)

	// Register a Message listener
	receiver.RegisterListener("example_test", exampleListener)

	// Wait for the Message to be consumed
	time.Sleep(time.Second)

	// Output:
	// Received new event: ExampleTestEvent, Message=example-message
}

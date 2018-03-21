package mq_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/client"
	"github.com/seagullbird/headr-common/mq/dispatch"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/streadway/amqp"
	"os"
	"time"
)

// Example listener
func exampleListener(delivery amqp.Delivery) {
	var event mq.ExampleEvent
	if err := json.Unmarshal(delivery.Body, &event); err != nil {
		panic(err)
	}
	fmt.Printf("Received new event: %s", event)
}

func makeExampleListener(logger log.Logger) receive.Listener {
	return func(delivery amqp.Delivery) {
		var event mq.ExampleEvent
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			panic(err)
		}
		logger.Log("info", "received new event", "event", event)
	}
}

func Example() {
	logger := log.NewLogfmtLogger(os.Stdout)

	var (
		servername = os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
		username   = "guest"
		passwd     = "guest"
	)

	// New dispatcher
	dispatcher, err := dispatch.NewDispatcher(client.New(servername, username, passwd), logger)
	if err != nil {
		panic(err)
	}

	// Dispatch a Message
	msg := mq.ExampleEvent{
		Message: "example-message",
	}
	err = dispatcher.DispatchMessage("example_test", msg)
	if err != nil {
		panic(err)
	}

	// Wait for the Message to be produced
	time.Sleep(time.Second)

	// New receiver
	receiver, err := receive.NewReceiver(client.New(servername, username, passwd), logger)
	if err != nil {
		panic(err)
	}

	// Register a Message listener
	receiver.RegisterListener("example_test", makeExampleListener(logger))

	// Wait for the Message to be consumed
	time.Sleep(time.Second)

	// Output:
	// info="Dispatching message to queue" queue_name=example_test
	// info="New Listener registered" queue_name=example_test
	// info="received new event" event="ExampleTestEvent, Message=example-message"
}

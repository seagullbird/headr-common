package mq_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/seagullbird/headr-common/mq"
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
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

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
	dispatcher, err := dispatch.NewDispatcher(dConn, logger)
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
	// Make connection to rabbitmq server
	rConn, err := mq.MakeConn(servername, username, passwd)
	if err != nil {
		panic(err)
	}

	receiver, err := receive.NewReceiver(rConn, logger)
	if err != nil {
		panic(err)
	}

	// Register a Message listener
	receiver.RegisterListener("example_test", makeExampleListener(logger))

	// Wait for the Message to be consumed
	time.Sleep(time.Second)

	// Output:
	// caller=dispatch.go:23 info="Dispatching message to queue" queue_name=example_test
	// caller=receive.go:40 info="New Listener registered" queue_name=example_test
	// caller=example_test.go:30 info="received new event" event="ExampleTestEvent, Message=example-message"
}

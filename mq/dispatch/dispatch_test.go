package dispatch_test

import (
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/dispatch"
	"os"
	"testing"
)

func TestDispatchMessage(t *testing.T) {
	var (
		servername = os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
		username   = "guest"
		passwd     = "guest"
	)
	conn, err := mq.MakeConn(servername, username, passwd)
	if err != nil {
		t.Fatal("Cannot connection to RabbitMQ", err)
	}

	dispatcher, err := dispatch.NewDispatcher(conn, "dispatch_test")
	if err != nil {
		t.Fatal("Cannot create dispatcher", err)
	}

	event := mq.ExampleEvent{
		Message: "dispatch-test",
	}
	if err := dispatcher.DispatchMessage(event); err != nil {
		t.Fatal("Failed to dispatch message", err)
	}
}

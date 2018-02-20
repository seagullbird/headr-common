package dispatch_test

import (
	"github.com/go-kit/kit/log"
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/dispatch"
	"os"
	"testing"
)

func TestDispatchMessage(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		servername = os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
		username   = "guest"
		passwd     = "guest"
	)
	conn, err := mq.MakeConn(servername, username, passwd)
	if err != nil {
		t.Fatal("Cannot connection to RabbitMQ", err)
	}

	dispatcher, err := dispatch.NewDispatcher(conn, logger)
	if err != nil {
		t.Fatal("Cannot create dispatcher", err)
	}

	event := mq.ExampleEvent{
		Message: "dispatch-test",
	}
	if err := dispatcher.DispatchMessage("dispatch_test", event); err != nil {
		t.Fatal("Failed to dispatch message", err)
	}
}

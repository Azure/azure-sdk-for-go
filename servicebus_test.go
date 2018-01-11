package servicebus

import (
	"testing"

	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"pack.ag/amqp"
	"time"
)

func TestServiceBus(t *testing.T) {
	connStr := os.Getenv("AZURE_SERVICE_BUS_CONN_STR") // `Endpoint=sb://XXXX.servicebus.windows.net/;SharedAccessKeyName=XXXX;SharedAccessKey=XXXX`
	sb, err := New(connStr)
	defer sb.Close()
	assert.Nil(t, err)

	catcher := make(chan string)
	sb.Receive("myQueue", func(c context.Context, message *amqp.Message) error {
		catcher <- string(message.Data)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	go func() {
		sb.Send(ctx, "myQueue", &amqp.Message{
			Data: []byte("Hello World!"),
			Properties: &amqp.MessageProperties{
				MessageID: "1",
			},
		})
	}()

	select {
	case msg := <-catcher:
		assert.Equal(t, msg, "Hello World!")
	case <- ctx.Done():
		t.Log("timed out")
		t.Fail()
	}

}

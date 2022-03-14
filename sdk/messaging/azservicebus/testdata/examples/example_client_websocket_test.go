package examples_test

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/joho/godotenv"
	"nhooyr.io/websocket"
)

func Example_usingWebsockets() {
	const exampleQueueName = "service-bus-examples"
	connectionString, cleanup := exampleSetup(exampleQueueName)
	defer cleanup()

	// NOTE: If you'd like to authenticate via Azure Active Directory look at
	// the `azservicebus.NewClient` function instead.
	client, err := azservicebus.NewClientFromConnectionString(connectionString, &azservicebus.ClientOptions{
		NewWebSocketConn: func(ctx context.Context, args azservicebus.NewWebSocketConnArgs) (net.Conn, error) {
			opts := &websocket.DialOptions{Subprotocols: []string{"amqp"}}
			wssConn, _, err := websocket.Dial(ctx, args.Host, opts)

			if err != nil {
				return nil, err
			}

			return websocket.NetConn(context.TODO(), wssConn, websocket.MessageBinary), nil
		},
	})

	if err != nil {
		panic(err)
	}

	defer func() { _ = client.Close(context.TODO()) }()

	receiver, err := client.NewReceiverForQueue(exampleQueueName, nil)

	if err != nil {
		panic(err)
	}

	defer func() { _ = receiver.Close(context.TODO()) }()

	sender, err := client.NewSender(exampleQueueName, nil)

	if err != nil {
		panic(err)
	}

	err = sender.SendMessage(context.TODO(), &azservicebus.Message{
		Body: []byte("hello world"),
	})

	if err != nil {
		panic(err)
	}

	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)

	if err != nil {
		panic(err)
	}

	if len(messages) != 1 {
		panic("No messages in queue")
	}

	body, err := messages[0].Body()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Body: \"%s\"", string(body))

	err = receiver.CompleteMessage(context.TODO(), messages[0])

	if err != nil {
		panic(err)
	}

	// Output:
	// Body: "hello world"
}

func exampleSetup(queueName string) (string, func()) {
	_ = godotenv.Load()
	connectionString := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	if connectionString == "" {
		panic("No connection string given, can't run example")
	}

	adminClient, err := admin.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		panic(err)
	}

	_, err = adminClient.CreateQueue(context.TODO(), queueName, nil, nil)

	if err != nil {
		panic(err)
	}

	return connectionString, func() {
		_, _ = adminClient.DeleteQueue(context.TODO(), queueName, nil)
	}
}

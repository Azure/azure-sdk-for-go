// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/coder/websocket"
)

func Example_usingWebsocketsAndProxies() {
	eventHubNamespace := os.Getenv("EVENTHUB_NAMESPACE") // <ex: myeventhubnamespace.servicebus.windows.net>
	eventHubName := os.Getenv("EVENTHUB_NAME")

	if eventHubName == "" || eventHubNamespace == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	// You can use an HTTP proxy, with websockets, by setting the appropriate HTTP(s)_PROXY
	// variable in your environment, as described in the https://pkg.go.dev/net/http#ProxyFromEnvironment
	// function.
	//
	// A proxy is NOT required to use websockets.
	newWebSocketConnFn := func(ctx context.Context, args azeventhubs.WebSocketConnParams) (net.Conn, error) {
		opts := &websocket.DialOptions{Subprotocols: []string{"amqp"}}
		wssConn, _, err := websocket.Dial(ctx, args.Host, opts)

		if err != nil {
			return nil, err
		}

		return websocket.NetConn(ctx, wssConn, websocket.MessageBinary), nil
	}

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	consumerClient, err = azeventhubs.NewConsumerClient(eventHubNamespace, eventHubName, azeventhubs.DefaultConsumerGroup, defaultAzureCred, &azeventhubs.ConsumerClientOptions{
		NewWebSocketConn: newWebSocketConnFn,
	})

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// NOTE: For users of `coder/websocket` there's an open discussion here:
	//   https://github.com/coder/websocket/issues/520
	//
	// An error ("failed to read frame header: EOF") can be returned when the
	// websocket connection is closed. This error will be returned from the
	// `ConsumerClient.Close` or `ProducerClient.Close` functions and can be
	// ignored, as the websocket "close handshake" has already completed.
	defer func() { _ = consumerClient.Close(context.TODO()) }()
}

var _ any // (ignore, used for docs)

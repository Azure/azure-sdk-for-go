// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"net"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"nhooyr.io/websocket"
)

func ExampleNewClient_usingWebsocketsAndProxies() {
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

	// NOTE: If you'd like to authenticate via Azure Active Directory use the
	// the `NewProducerClient` or `NewConsumerClient` functions.
	consumerClient, err = azeventhubs.NewConsumerClientFromConnectionString("connection-string", "event-hub", azeventhubs.DefaultConsumerGroup, &azeventhubs.ConsumerClientOptions{
		NewWebSocketConn: newWebSocketConnFn,
	})

	if err != nil {
		panic(err)
	}

	// NOTE: For users of `nhooyr.io/websocket` there's an open discussion here:
	//   https://github.com/nhooyr/websocket/discussions/380
	//
	// An error ("failed to read frame header: EOF") can be returned when the
	// websocket connection is closed. This error will be returned from the
	// `ConsumerClient.Close` or `ProducerClient.Close` functions and can be
	// ignored, as the websocket "close handshake" has already completed.
	defer consumerClient.Close(context.TODO())
}

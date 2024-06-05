// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"net"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"nhooyr.io/websocket"
)

func ExampleNewClient_usingWebsocketsAndProxies() {
	// You can use an HTTP proxy, with websockets, by setting the appropriate HTTP(s)_PROXY
	// variable in your environment, as described in the https://pkg.go.dev/net/http#ProxyFromEnvironment
	// function.
	//
	// A proxy is NOT required to use websockets.
	newWebSocketConnFn := func(ctx context.Context, args azservicebus.NewWebSocketConnArgs) (net.Conn, error) {
		opts := &websocket.DialOptions{Subprotocols: []string{"amqp"}}
		wssConn, _, err := websocket.Dial(ctx, args.Host, opts)

		if err != nil {
			return nil, err
		}

		return websocket.NetConn(ctx, wssConn, websocket.MessageBinary), nil
	}

	// NOTE: If you'd like to authenticate using a Service Bus connection string
	// look at `NewClientFromConnectionString` instead.
	client, err = azservicebus.NewClient(endpoint, tokenCredential, &azservicebus.ClientOptions{
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
	// `Client.Close` function and can be ignored, as the websocket "close handshake"
	// has already completed.
	defer client.Close(context.TODO())
}

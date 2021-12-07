// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"net"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"nhooyr.io/websocket"
)

func ExampleNewClient() {
	// NOTE: If you'd like to authenticate using a Service Bus connection string
	// look at `NewClientFromConnectionString` instead.

	// For more information about the DefaultAzureCredential:
	// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential
	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	client, err = azservicebus.NewClient("<ex: myservicebus.servicebus.windows.net>", credential, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleNewClientFromConnectionString() {
	// NOTE: If you'd like to authenticate via Azure Active Directory look at
	// the `NewClient` function instead.

	client, err = azservicebus.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleNewClient_usingWebsockets() {
	// NOTE: If you'd like to authenticate via Azure Active Directory look at
	// the `NewClient` function instead.
	client, err = azservicebus.NewClientFromConnectionString(connectionString, &azservicebus.ClientOptions{
		NewWebSocketConn: func(ctx context.Context, args azservicebus.NewWebSocketConnArgs) (net.Conn, error) {
			opts := &websocket.DialOptions{Subprotocols: []string{"amqp"}}
			wssConn, _, err := websocket.Dial(ctx, args.Host, opts)

			if err != nil {
				return nil, err
			}

			return websocket.NetConn(context.Background(), wssConn, websocket.MessageBinary), nil
		},
	})

	if err != nil {
		panic(err)
	}
}

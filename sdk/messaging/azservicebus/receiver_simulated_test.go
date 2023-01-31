// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

func TestReceiver_Simulated(t *testing.T) {
	md := emulation.NewMockData(t, nil)

	client, err := newClientImpl(clientCreds{
		connectionString: "Endpoint=sb://example.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=DEADBEEF",
	}, struct {
		Client *ClientOptions
		NS     []internal.NamespaceOption
	}{
		NS: []internal.NamespaceOption{
			internal.NamespaceWithNewClientFn(md.NewConnection),
		},
	})

	require.NoError(t, err)
	require.NotNil(t, client)

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, sender)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	test.EnableStdoutLogging()

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	require.Equal(t, 1, len(md.Events.GetOpenConns()))
	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver each own 3 links apiece ($mgmt, actual link)")

	err = client.Close(context.Background())
	require.NoError(t, err)

	emulation.RequireNoLeaks(t, md.Events)
}

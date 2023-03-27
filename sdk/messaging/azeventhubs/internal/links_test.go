// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/stretchr/testify/require"
)

func TestLinksCBSLinkStillOpen(t *testing.T) {
	// we're not going to use this client for these tests.
	testParams := test.GetConnectionParamsForTest(t)
	ns, err := NewNamespace(NamespaceWithConnectionString(testParams.ConnectionString))
	require.NoError(t, err)

	defer func() { _ = ns.Close(context.Background(), true) }()

	session, oldConnID, err := ns.NewAMQPSession(context.Background())
	require.NoError(t, err)

	// opening a Sender to the $cbs endpoint. This endpoint can only be opened by a single
	// sender/receiver pair in a connection.
	_, err = session.NewSender(context.Background(), "$cbs", nil)
	require.NoError(t, err)

	newLinkFn := func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (AMQPSenderCloser, error) {
		return session.NewSender(ctx, entityPath, &amqp.SenderOptions{
			SettlementMode:              to.Ptr(amqp.SenderSettleModeMixed),
			RequestedReceiverSettleMode: to.Ptr(amqp.ReceiverSettleModeFirst),
		})
	}

	formatEntityPath := func(partitionID string) string {
		return fmt.Sprintf("%s/Partitions/%s", testParams.EventHubName, partitionID)
	}

	links := NewLinks(ns, fmt.Sprintf("%s/$management", testParams.EventHubName), formatEntityPath, newLinkFn)

	var lwid LinkWithID[AMQPSenderCloser]

	err = links.Retry(context.Background(), exported.EventConn, "test", "0", exported.RetryOptions{
		RetryDelay:    -1,
		MaxRetryDelay: time.Millisecond,
	}, func(ctx context.Context, innerLWID LinkWithID[AMQPSenderCloser]) error {
		lwid = innerLWID
		return nil
	})
	require.NoError(t, err)

	defer func() {
		err := links.Close(context.Background())
		require.NoError(t, err)
	}()

	require.NoError(t, err)
	require.Equal(t, oldConnID+1, lwid.ConnID, "Connection gets incremented since it had to be reset")
}

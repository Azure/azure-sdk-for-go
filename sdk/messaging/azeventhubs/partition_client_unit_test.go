// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/stretchr/testify/require"
)

func TestUnit_PartitionClient_PrefetchOff(t *testing.T) {
	ns := &internal.FakeNSForPartClient{
		Receiver: &internal.FakeAMQPReceiver{
			Messages: []*amqp.Message{
				{}, {}, {},
			},
		},
	}

	client, err := newPartitionClient(partitionClientArgs{
		namespace: ns,
	}, &PartitionClientOptions{
		Prefetch: -1,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	events, err := client.ReceiveEvents(ctx, 3, nil)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	require.Equal(t, []uint32{uint32(3)}, ns.Receiver.IssuedCredit, "Non-prefetch scenarios will issue credit at the time of request")
	require.Equal(t, uint32(0), ns.Receiver.ActiveCredits, "All messages should have been received")
	require.True(t, ns.Receiver.ManualCreditsSetFromOptions)
}

func TestUnit_PartitionClient_PrefetchOff_CreditLimits(t *testing.T) {
	ns := &internal.FakeNSForPartClient{
		Receiver: &internal.FakeAMQPReceiver{
			Messages: fakeMessages(int(defaultMaxCreditSize)),
		},
	}

	client, err := newPartitionClient(partitionClientArgs{
		namespace: ns,
	}, &PartitionClientOptions{
		Prefetch: -1,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// can't request over the max
	events, err := client.ReceiveEvents(ctx, int(defaultMaxCreditSize+1), nil)
	require.EqualError(t, err, "count cannot exceed 5000")
	require.Empty(t, events)

	// can't request negative/0 credits
	events, err = client.ReceiveEvents(ctx, 0, nil)
	require.EqualError(t, err, "count should be greater than 0")
	require.Empty(t, events)

	events, err = client.ReceiveEvents(ctx, -1, nil)
	require.EqualError(t, err, "count should be greater than 0")
	require.Empty(t, events)

	// can request the max
	events, err = client.ReceiveEvents(ctx, int(defaultMaxCreditSize), nil)
	require.NoError(t, err)
	require.NotEmpty(t, events)
}

func TestUnit_PartitionClient_PrefetchOffOnlyBackfillsCredits(t *testing.T) {
	testData := []struct {
		Name    string
		Initial uint32
		Issued  []uint32
	}{
		{"Need some more credits", 2, []uint32{uint32(1)}},
		{"No more credits needed", 3, nil},
	}

	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {
			ns := &internal.FakeNSForPartClient{
				Receiver: &internal.FakeAMQPReceiver{
					Messages: []*amqp.Message{{}, {}, {}},
				},
			}

			client, err := newPartitionClient(partitionClientArgs{
				namespace: ns,
			}, &PartitionClientOptions{
				Prefetch: -1,
			})
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// we're going to artifically make it seem like we already have credits issued
			// this makes it so the next call will issue just enough credits to get it to
			// what we requested.
			ns.Receiver.ActiveCredits = td.Initial

			events, err := client.ReceiveEvents(ctx, 3, nil)
			require.NoError(t, err)
			require.NotEmpty(t, events)

			require.Equal(t, td.Issued, ns.Receiver.IssuedCredit, "Only issue credits to backfill missing credits")
			require.Equal(t, uint32(0), ns.Receiver.ActiveCredits, "All messages should have been received")
			require.True(t, ns.Receiver.ManualCreditsSetFromOptions)
		})
	}
}

func TestUnit_PartitionClient_PrefetchOn(t *testing.T) {
	testData := []struct {
		options        *PartitionClientOptions
		initialCredits uint32
	}{
		{nil, 300}, //  (300 credits is the default prefetch)
		{&PartitionClientOptions{Prefetch: 101}, 101},
	}

	for _, td := range testData {
		ns := &internal.FakeNSForPartClient{
			Receiver: &internal.FakeAMQPReceiver{
				Messages: []*amqp.Message{{}, {}, {}},
			},
		}

		client, err := newPartitionClient(partitionClientArgs{
			namespace: ns,
		}, td.options)
		require.NoError(t, err)

		events, err := client.ReceiveEvents(context.Background(), 3, nil)
		require.NoError(t, err)
		require.NotEmpty(t, events)

		require.Equal(t, td.initialCredits, ns.Receiver.CreditsSetFromOptions, "All messages should have been received")
		require.Nil(t, ns.Receiver.IssuedCredit, "prefetching doesn't manually issue credits")

		require.Equal(t, uint32(td.initialCredits-3), ns.Receiver.ActiveCredits, "All messages should have been received")
	}
}

func TestUnit_PartitionClient_PrefetchLimit(t *testing.T) {
	newPartitionClient := func(prefetch int32) (*PartitionClient, error) {
		ns := &internal.FakeNSForPartClient{
			Receiver: &internal.FakeAMQPReceiver{
				Messages: fakeMessages(int(defaultMaxCreditSize) + 1),
			},
		}

		client, err := newPartitionClient(partitionClientArgs{namespace: ns}, &PartitionClientOptions{
			Prefetch: prefetch,
		})

		return client, err
	}

	t.Run("max allowed credits is defaultMaxCreditSize", func(t *testing.T) {
		client, err := newPartitionClient(int32(defaultMaxCreditSize))
		require.NoError(t, err)
		require.NotNil(t, client)

		test.RequireClose(t, client)
	})

	t.Run("can't request zero or negative credits", func(t *testing.T) {
		client, err := newPartitionClient(int32(defaultMaxCreditSize))
		require.NoError(t, err)
		require.NotNil(t, client)

		events, err := client.ReceiveEvents(context.Background(), 0, nil)
		require.EqualError(t, err, "count should be greater than 0")
		require.Empty(t, events)

		events, err = client.ReceiveEvents(context.Background(), -1, nil)
		require.EqualError(t, err, "count should be greater than 0")
		require.Empty(t, events)

		test.RequireClose(t, client)
	})

	t.Run("can receive more than defaultMaxCreditSize in prefetch mode", func(t *testing.T) {
		client, err := newPartitionClient(0)
		require.NoError(t, err)
		require.NotNil(t, client)

		// if you're using prefetch it's fine to ask for more than the `defaultMaxCreditSize`
		// since it doesn't actually cause to request more credits than is allowed.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		events, err := client.ReceiveEvents(ctx, int(defaultMaxCreditSize+1), nil)
		require.NoError(t, err)
		require.NotEmpty(t, events)

		test.RequireClose(t, client)
	})

	t.Run("can't set a option.Prefetch > defaultMaxCreditSize", func(t *testing.T) {
		// and you can't create a PartitionClient that uses more credits than allowed.
		client, err := newPartitionClient(int32(defaultMaxCreditSize) + 1)
		require.EqualError(t, err, fmt.Sprintf("options.Prefetch cannot exceed %d", defaultMaxCreditSize))
		require.Nil(t, client)
	})
}

func fakeMessages(count int) []*amqp.Message {
	var messages []*amqp.Message

	for i := 0; i < count; i++ {
		messages = append(messages, &amqp.Message{})
	}

	return messages
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/stretchr/testify/require"
)

func TestConsumerClient_Recovery(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	// Uncomment to see the entire recovery playbook run.
	test.EnableStdoutLogging()

	tokenCred, err := credential.New(nil)
	require.NoError(t, err)

	// Overview:
	// 1. Send one event per partition
	// 2. Receive one event per partition. This'll ensure the links are live.
	// 3. Grub into the client to get access to it's connection and shut it off.
	// 4. Try again, everything should recover.
	producerClient, err := NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, tokenCred, nil)
	require.NoError(t, err)

	ehProps, err := producerClient.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	// trim the partition list down so the test executes in resonable time.
	ehProps.PartitionIDs = ehProps.PartitionIDs[0:3] // min for testing is 3 partitions anyways

	type sendResult struct {
		PartitionID  string
		OffsetBefore int64
	}

	sendResults := make([]sendResult, len(ehProps.PartitionIDs))
	wg := sync.WaitGroup{}

	log.Printf("1. sending 2 events to %d partitions", len(ehProps.PartitionIDs))

	for i, pid := range ehProps.PartitionIDs {
		wg.Add(1)

		go func(i int, pid string) {
			defer wg.Done()

			partProps, err := producerClient.GetPartitionProperties(context.Background(), pid, nil)
			require.NoError(t, err)
			require.Equal(t, pid, partProps.PartitionID)

			t.Logf("[%s] Starting props %#v", pid, partProps)

			batch, err := producerClient.NewEventDataBatch(context.Background(), &EventDataBatchOptions{
				PartitionID: &pid,
			})
			require.NoError(t, err)

			require.NoError(t, batch.AddEventData(&EventData{
				Body: []byte(fmt.Sprintf("event 1 for partition %s", pid)),
			}, nil))

			require.NoError(t, batch.AddEventData(&EventData{
				Body: []byte(fmt.Sprintf("event 2 for partition %s", pid)),
			}, nil))

			err = producerClient.SendEventDataBatch(context.Background(), batch, nil)
			require.NoError(t, err)

			afterPartProps, err := producerClient.GetPartitionProperties(context.Background(), pid, nil)
			require.NoError(t, err)
			require.Equal(t, pid, afterPartProps.PartitionID)

			t.Logf("[%s] After props %#v", pid, afterPartProps)

			require.Equalf(t, int64(2), afterPartProps.LastEnqueuedSequenceNumber-partProps.LastEnqueuedSequenceNumber, "Expected only 2 messages in partition %s", pid)

			sendResults[i] = sendResult{PartitionID: pid, OffsetBefore: partProps.LastEnqueuedOffset}
		}(i, pid)
	}

	wg.Wait()

	test.RequireClose(t, producerClient)

	// now we'll receive an event (so we know each partition client is alive)
	// each partition actually has two offsets.
	consumerClient, err := NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, DefaultConsumerGroup, tokenCred, nil)
	require.NoError(t, err)

	partitionClients := make([]*PartitionClient, len(sendResults))

	log.Printf("2. receiving the first event for each partition")

	for i, sr := range sendResults {
		wg.Add(1)

		go func(i int, sr sendResult) {
			defer wg.Done()

			partClient, err := consumerClient.NewPartitionClient(sr.PartitionID, &PartitionClientOptions{
				StartPosition: StartPosition{Inclusive: false, Offset: &sr.OffsetBefore},
				Prefetch:      -1,
			})
			require.NoError(t, err)

			partitionClients[i] = partClient

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			events, err := partClient.ReceiveEvents(ctx, 1, nil)
			require.NoError(t, err)
			require.EqualValues(t, 1, len(events))

			t.Logf("[%s] Received seq:%d, offset:%d", sr.PartitionID, events[0].SequenceNumber, events[0].Offset)

			require.Equal(t, fmt.Sprintf("event 1 for partition %s", sr.PartitionID), string(events[0].Body))
		}(i, sr)
	}

	wg.Wait()

	defer test.RequireClose(t, consumerClient)

	log.Printf("3. closing internal connection (non-permanently), which will force recovery for each partition client so they can read the next event")

	// now we'll close the internal connection, simulating a connection break
	require.NoError(t, consumerClient.namespace.Close(context.Background(), false))

	var best int64

	log.Printf("4. try to read the second event, which force clients to recover")

	// and try to receive the second event for each client
	for i, pc := range partitionClients {
		wg.Add(1)

		go func(i int, pc *PartitionClient) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			events, err := pc.ReceiveEvents(ctx, 1, nil)
			require.NoError(t, err)
			require.EqualValues(t, 1, len(events))
			require.Equal(t, fmt.Sprintf("event 2 for partition %s", sendResults[i].PartitionID), string(events[0].Body))

			atomic.AddInt64(&best, 1)
		}(i, pc)
	}

	wg.Wait()
	require.Equal(t, int64(len(ehProps.PartitionIDs)), best)
}

func TestConsumerClient_RecoveryLink(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	// Uncomment to see the entire recovery playbook run.
	test.EnableStdoutLogging()

	tokenCred, err := credential.New(nil)
	require.NoError(t, err)

	// Overview:
	// 1. Send one event per partition
	// 2. Receive one event per partition. This'll ensure the links are live.
	// 3. Grub into the client to get access to it's connection and shut it off.
	// 4. Try again, everything should recover.
	producerClient, err := NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, tokenCred, nil)
	require.NoError(t, err)

	ehProps, err := producerClient.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	// trim the partition list down so the test executes in resonable time.
	ehProps.PartitionIDs = ehProps.PartitionIDs[0:3] // min for testing is 3 partitions anyways

	type sendResult struct {
		PartitionID  string
		OffsetBefore int64
	}

	sendResults := make([]sendResult, len(ehProps.PartitionIDs))
	wg := sync.WaitGroup{}

	log.Printf("== 1. sending 2 events to %d partitions ==", len(ehProps.PartitionIDs))

	for i, pid := range ehProps.PartitionIDs {
		wg.Add(1)

		go func(i int, pid string) {
			defer wg.Done()

			partProps, err := producerClient.GetPartitionProperties(context.Background(), pid, nil)
			require.NoError(t, err)

			batch, err := producerClient.NewEventDataBatch(context.Background(), &EventDataBatchOptions{
				PartitionID: &pid,
			})
			require.NoError(t, err)

			require.NoError(t, batch.AddEventData(&EventData{
				Body: []byte(fmt.Sprintf("event 1 for partition %s", pid)),
			}, nil))

			require.NoError(t, batch.AddEventData(&EventData{
				Body: []byte(fmt.Sprintf("event 2 for partition %s", pid)),
			}, nil))

			err = producerClient.SendEventDataBatch(context.Background(), batch, nil)
			require.NoError(t, err)

			sendResults[i] = sendResult{PartitionID: pid, OffsetBefore: partProps.LastEnqueuedOffset}
		}(i, pid)
	}

	wg.Wait()

	test.RequireClose(t, producerClient)

	// now we'll receive an event (so we know each partition client is alive)
	// each partition actually has two offsets.
	consumerClient, err := NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, DefaultConsumerGroup, tokenCred, nil)
	require.NoError(t, err)

	partitionClients := make([]*PartitionClient, len(sendResults))

	log.Printf("== 2. receiving the first event for each partition == ")

	for i, sr := range sendResults {
		wg.Add(1)

		go func(i int, sr sendResult) {
			defer wg.Done()

			partClient, err := consumerClient.NewPartitionClient(sr.PartitionID, &PartitionClientOptions{
				StartPosition: StartPosition{Inclusive: false, Offset: &sr.OffsetBefore},
				Prefetch:      -1,
			})
			require.NoError(t, err)

			partitionClients[i] = partClient

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			events, err := partClient.ReceiveEvents(ctx, 1, nil)
			require.NoError(t, err)
			require.EqualValues(t, 1, len(events))
			require.Equal(t, fmt.Sprintf("event 1 for partition %s", sr.PartitionID), string(events[0].Body))
		}(i, sr)
	}

	wg.Wait()

	defer test.RequireClose(t, consumerClient)

	var best int64

	log.Printf("== 3. Closing links, but leaving connection intact ==")

	for i, pc := range partitionClients {
		links := pc.links.(*internal.Links[amqpwrap.AMQPReceiverCloser])
		lwid, err := links.GetLink(context.Background(), sendResults[i].PartitionID)
		require.NoError(t, err)
		require.NoError(t, lwid.Link().Close(context.Background()))
	}

	log.Printf("== 4. try to read the second event, which force clients to recover ==")

	// and try to receive the second event for each client
	for i, pc := range partitionClients {
		wg.Add(1)

		go func(i int, pc *PartitionClient) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			events, err := pc.ReceiveEvents(ctx, 1, nil)
			require.NoError(t, err)
			require.EqualValues(t, 1, len(events))
			require.Equal(t, fmt.Sprintf("event 2 for partition %s", sendResults[i].PartitionID), string(events[0].Body))

			atomic.AddInt64(&best, 1)
		}(i, pc)
	}

	wg.Wait()
	require.Equal(t, int64(len(ehProps.PartitionIDs)), best)
}

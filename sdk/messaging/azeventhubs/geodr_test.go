// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/stretchr/testify/require"
)

func TestConsumerClient_GeoReplication(t *testing.T) {
	// this test just needs a single partition to test
	var partitionID = "0"

	testParams := test.GetConnectionParamsForTest(t)

	if testParams.GeoDRNamespace == "" || testParams.GeoDRHubName == "" || testParams.GeoDRStorageEndpoint == "" {
		t.Skipf("Skipping GeoDR test, EVENTHUBS_GEODR_NAMESPACE or EVENTHUBS_GEODR_HUBNAME or EVENTHUBS_GEODR_CHECKPOINTSTORE_STORAGE_ENDPOINT was not set")
	}

	propsBeforeTest := func() azeventhubs.PartitionProperties {
		producer, err := azeventhubs.NewProducerClient(testParams.GeoDRNamespace, testParams.GeoDRHubName, testParams.Cred, nil)
		require.NoError(t, err)

		defer test.RequireClose(t, producer)

		props, err := producer.GetEventHubProperties(context.Background(), nil)
		require.NoError(t, err)
		require.True(t, props.GeoReplicationEnabled)

		propsBeforeTest, err := producer.GetPartitionProperties(context.Background(), "0", nil)
		require.NoError(t, err)

		// This is what the partition properties look like, with geo-replication enabled
		// (note, this example event hub was empty):
		// {
		//     BeginningSequenceNumber:-1,
		// 	   EventHubName:"ehrp2",
		// 	   IsEmpty:true,
		// 	   LastEnqueuedOffset:"2:-1:-1",
		// 	   LastEnqueuedOn:time.Date(1, time.January, 1, 0, 0, 0, 0, time.Local),
		//     LastEnqueuedSequenceNumber:-1,
		//     PartitionID:"0"
		// }

		t.Logf("LastEnqueuedOffset: %#v, LastEnqueuedSequenceNumber: %#v", propsBeforeTest.LastEnqueuedOffset, propsBeforeTest.LastEnqueuedSequenceNumber)

		// we send a couple of events so the processor tests, that can't be started inclusive, will still have something
		// predictable to retrieve.
		batch, err := producer.NewEventDataBatch(context.Background(), &azeventhubs.EventDataBatchOptions{
			PartitionID: &partitionID,
		})
		require.NoError(t, err)

		// the Event Hub is re-used, so this is the first in our sent messages, but not necessarily
		// the actual first message in the partition.
		err = batch.AddEventData(&azeventhubs.EventData{Body: []byte("1")}, nil)
		require.NoError(t, err)

		err = batch.AddEventData(&azeventhubs.EventData{Body: []byte("2")}, nil)
		require.NoError(t, err)

		// even if the event hub is re-used, this is still the last message for sure.
		err = batch.AddEventData(&azeventhubs.EventData{Body: []byte("3")}, nil)
		require.NoError(t, err)

		err = producer.SendEventDataBatch(context.Background(), batch, nil)
		require.NoError(t, err)

		return propsBeforeTest
	}()

	earliestEvent, ourFirstEvent := func() (*azeventhubs.ReceivedEventData, *azeventhubs.ReceivedEventData) {
		consumer, err := azeventhubs.NewConsumerClient(testParams.GeoDRNamespace, testParams.GeoDRHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
		require.NoError(t, err)

		pc, err := consumer.NewPartitionClient(partitionID, &azeventhubs.PartitionClientOptions{
			StartPosition: azeventhubs.StartPosition{
				Earliest:  to.Ptr(true),
				Inclusive: true,
			},
		})
		require.NoError(t, err)
		defer test.RequireClose(t, pc)

		earliestEvents, err := pc.ReceiveEvents(context.Background(), 1, nil)
		require.NoError(t, err)
		require.NotEmpty(t, earliestEvents)

		// Now let's start at the point just before we started sending events. We want those offsets as well.
		pc, err = consumer.NewPartitionClient(partitionID, &azeventhubs.PartitionClientOptions{
			StartPosition: azeventhubs.StartPosition{
				Offset:    &propsBeforeTest.LastEnqueuedOffset,
				Inclusive: false,
			},
		})
		require.NoError(t, err)
		defer test.RequireClose(t, pc)

		ourEvents, err := pc.ReceiveEvents(context.Background(), 3, nil)
		require.NoError(t, err)
		require.NotEmpty(t, ourEvents)

		return earliestEvents[0], ourEvents[0]
	}()

	runTest := func(t *testing.T, proc *azeventhubs.Processor) *azeventhubs.ReceivedEventData {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch := make(chan struct{})

		go func() {
			defer close(ch)
			err := proc.Run(ctx)
			require.NoError(t, err)
		}()

		var event *azeventhubs.ReceivedEventData

		for {
			pc := proc.NextPartitionClient(context.Background())

			if pc.PartitionID() == "0" {
				events, err := pc.ReceiveEvents(context.Background(), 1, nil)
				require.NoError(t, err)
				require.NotEmpty(t, events)

				event = events[0]

				cancel()
				_ = pc.Close(context.Background())
				break
			} else {
				_ = pc.Close(context.Background())
			}
		}

		require.NotNil(t, event)
		<-ch
		return event
	}

	t.Run("ProcessorWithLegacyOffset", func(t *testing.T) {
		setup := func(t *testing.T) *processorTestData {
			td := setupProcessorTest(t, true)

			err = td.CheckpointStore.SetCheckpoint(context.Background(), azeventhubs.Checkpoint{
				ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
				FullyQualifiedNamespace: testParams.GeoDRNamespace,
				EventHubName:            testParams.GeoDRHubName,
				PartitionID:             "0",
				// this is invalid - you can't use old offsets with a new GeoDR-enabled Event Hub once it's
				// been promoted.
				Offset:         to.Ptr("0"),
				SequenceNumber: to.Ptr(int64(1)),
			}, nil)
			require.NoError(t, err)

			return td
		}

		td := setup(t)
		proc := td.Create(nil)
		event := runTest(t, proc)

		// Here's what happens here:
		// 1. Processor loads up checkpoint, which contains a "legacy" offset (ie, just an integer)
		// 2. It attempts to create a consumer using that offset, which Event Hubs rejects, with a GeoDR related error
		// 3. We then fallback to opening up the start of the partition instead (ie: earliest)
		require.Equal(t, earliestEvent, event)
	})

	t.Run("Processor", func(t *testing.T) {
		setup := func(t *testing.T) *processorTestData {
			td := setupProcessorTest(t, true)

			err = td.CheckpointStore.SetCheckpoint(context.Background(), azeventhubs.Checkpoint{
				ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
				FullyQualifiedNamespace: testParams.GeoDRNamespace,
				EventHubName:            testParams.GeoDRHubName,
				PartitionID:             "0",
				// Checkpoints always point to the last event received, so we will receive the event just after ourFirstEvent
				Offset:         &ourFirstEvent.Offset,
				SequenceNumber: &ourFirstEvent.SequenceNumber,
			}, nil)
			require.NoError(t, err)

			return td
		}

		td := setup(t)
		proc := td.Create(nil)
		event := runTest(t, proc)

		require.Equal(t, "2", string(event.Body))
	})

	t.Run("StartWithOffsetFromGetPartitionProperties", func(t *testing.T) {
		cc, err := azeventhubs.NewConsumerClient(testParams.GeoDRNamespace, testParams.GeoDRHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
		require.NoError(t, err)
		defer test.RequireClose(t, cc)

		_, err = strconv.ParseInt(propsBeforeTest.LastEnqueuedOffset, 10, 64)
		require.Error(t, err, "offsets are no longer just integers")

		pc, err := cc.NewPartitionClient("0", &azeventhubs.PartitionClientOptions{
			StartPosition: azeventhubs.StartPosition{
				Offset: &propsBeforeTest.LastEnqueuedOffset,
			},
		})
		require.NoError(t, err)
		defer test.RequireClose(t, pc)

		events, err := pc.ReceiveEvents(context.Background(), 1, nil)
		require.NoError(t, err)
		require.NotEmpty(t, events)
	})
}

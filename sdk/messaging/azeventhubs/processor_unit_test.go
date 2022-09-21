// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
	"github.com/stretchr/testify/require"
)

func TestUnit_Processor_loadBalancing(t *testing.T) {
	cps := newCheckpointStoreForTest()
	firstProcessor := newProcessorForTest(t, "first-processor", cps)
	newAddressForPartition := func(partitionID string) CheckpointStoreAddress {
		return CheckpointStoreAddress{
			ConsumerGroup:           "consumer-group",
			EventHubName:            "event-hub",
			FullyQualifiedNamespace: "fqdn",
			PartitionID:             partitionID,
		}
	}
	require.Equal(t, ProcessorStrategyBalanced, firstProcessor.lb.strategy)

	allPartitionIDs := []string{"1", "100", "1001"}
	lbinfo, err := firstProcessor.lb.getAvailablePartitions(context.Background(), allPartitionIDs)
	require.NoError(t, err)

	// this is a completely empty checkpoint store so nobody owns any partitions yet
	// which means that we get to claim them all
	require.Empty(t, lbinfo.aboveMax)
	require.Empty(t, lbinfo.current)
	require.False(t, lbinfo.extraPartitionPossible)
	require.Equal(t, 3, lbinfo.maxAllowed, "only 1 possible owner (us), so we're allowed all the available partitions")

	expectedOwnerships := []Ownership{
		{
			CheckpointStoreAddress: newAddressForPartition("1"),
			OwnershipData: OwnershipData{
				OwnerID: "first-processor",
			}},
		{
			CheckpointStoreAddress: newAddressForPartition("100"),
			OwnershipData: OwnershipData{
				OwnerID: "first-processor",
			}},
		{
			CheckpointStoreAddress: newAddressForPartition("1001"),
			OwnershipData: OwnershipData{
				OwnerID: "first-processor",
			}},
	}

	require.Equal(t, expectedOwnerships, lbinfo.unownedOrExpired)

	// getAvailablePartitions doesn't mutate the checkpoint store.
	lbinfo, err = firstProcessor.lb.getAvailablePartitions(context.Background(), allPartitionIDs)
	require.NoError(t, err)
	require.Equal(t, expectedOwnerships, lbinfo.unownedOrExpired)

	// the balanced strategy claims one new partition per round, until balanced.
	// we'll do more in-depth testing in other tests, but this is just a basic
	// run through.
	firstProcessorOwnerships, err := firstProcessor.lb.LoadBalance(context.Background(), allPartitionIDs)
	require.NoError(t, err)

	expectedLoadBalancingOwnership := updateDynamicData(t, firstProcessorOwnerships[0], Ownership{
		CheckpointStoreAddress: newAddressForPartition("1001"),
		OwnershipData: OwnershipData{
			OwnerID: "first-processor",
		},
	}, allPartitionIDs)
	require.Equal(t, []Ownership{expectedLoadBalancingOwnership}, firstProcessorOwnerships)

	// at this point this is our state:
	// 3 total partitions ("1", "100", "1001")
	// 1 of those partitions is owned by our client ("first-processor")
	// 2 are still available.

	secondProcessor := newProcessorForTest(t, "second-processor", cps)

	// when we ask for available partitions we take into account the owners that are
	// present in the checkpoint store and ourselves, since we're about to try to claim
	// some partitions. So now it has to divide 3 partitions amongst two Processors.
	lbinfo, err = secondProcessor.lb.getAvailablePartitions(context.Background(), allPartitionIDs)
	require.NoError(t, err)
	require.Empty(t, lbinfo.aboveMax)
	require.Empty(t, lbinfo.current)
	require.True(t, lbinfo.extraPartitionPossible, "divvying 3 partitions amongst 2 processors")
	require.Equal(t, 2, lbinfo.maxAllowed, "now we're divvying up 3 partitions between 2 processors. At _most_ you can have min+1")

	// there are two available partition ownerships - we should be getting one of them.
	newProcessorOwnerships, err := secondProcessor.lb.LoadBalance(context.Background(), allPartitionIDs)
	require.NoError(t, err)

	newExpectedLoadBalancingOwnership := updateDynamicData(t, newProcessorOwnerships[0], Ownership{
		CheckpointStoreAddress: newAddressForPartition("1001"),
		OwnershipData: OwnershipData{
			OwnerID: "second-processor",
		},
	}, allPartitionIDs)

	require.Equal(t, []Ownership{newExpectedLoadBalancingOwnership}, newProcessorOwnerships)
	require.NotEqual(t, newExpectedLoadBalancingOwnership.PartitionID, expectedLoadBalancingOwnership.PartitionID, "partitions should not be assigned twice")

	//
	// now let's assign out the last partition - we'll pick a winner here and just use the second processor, but either one can technically claim it (or even attempt to at the same time!)
	//

	secondProcessorOwnershipsForLastPartition, err := secondProcessor.lb.LoadBalance(context.Background(), allPartitionIDs)
	require.NoError(t, err)

	require.Equal(t, 2, len(secondProcessorOwnershipsForLastPartition))

	// no overlap in partition assignments
	for _, o := range secondProcessorOwnershipsForLastPartition {
		require.NotEqual(t, firstProcessorOwnerships[0].PartitionID, o.PartitionID)
	}

	// and if we try to claim now with the first, it won't get anything new since
	// a. we're in a balanced state (all extra partitions assigned)
	// b. it has the minimum.
	time.Sleep(100 * time.Millisecond) // give a little gap so our last modified time is definitely greater
	afterBalanceOwnerships, err := firstProcessor.lb.LoadBalance(context.Background(), allPartitionIDs)
	require.NoError(t, err)
	require.Equal(t, 1, len(afterBalanceOwnerships))
	require.Equal(t, firstProcessorOwnerships[0].PartitionID, afterBalanceOwnerships[0].PartitionID)
	require.NotEqual(t, firstProcessorOwnerships[0].ETag, afterBalanceOwnerships[0].ETag, "ownership (etag) also gets updated each time we load balance")
	require.Greater(t, afterBalanceOwnerships[0].LastModifiedTime, firstProcessorOwnerships[0].LastModifiedTime, "ownership (last modified time) also gets updated each time we load balance")
}

func TestUnit_Processor_Run(t *testing.T) {
	cps := newCheckpointStoreForTest()

	processor, err := newProcessorImpl(simpleFakeConsumerClient(), cps, &NewProcessorOptions{
		PartitionExpirationDuration: time.Hour,
	})

	require.NoError(t, err)

	procCtx, procCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer procCancel()

	partClientValue := atomic.Value{}

	go func() {
		partitionClient := processor.NextPartitionClient(context.Background())
		partClientValue.Store(partitionClient)
		procCancel()
	}()

	time.Sleep(time.Second)
	require.NoError(t, processor.Run(procCtx))

	partitionClient := partClientValue.Load().(*ProcessorPartitionClient)
	require.NotNil(t, partitionClient)
	require.Equal(t, "a", partitionClient.partitionID)
}

func TestUnit_Processor_Run_singleConsumerPerPartition(t *testing.T) {
	cps := newCheckpointStoreForTest()
	ehProps := EventHubProperties{
		PartitionIDs: []string{"a"},
	}

	partitionClientsCreated := 0

	cc := &fakeConsumerClient{
		details: consumerClientDetails{
			ConsumerGroup:           "consumer-group",
			EventHubName:            "event-hub",
			FullyQualifiedNamespace: "fqdn",
			ClientID:                "my-client-id",
		},
		getEventHubPropertiesResult: ehProps,
		newPartitionClientFn: func(partitionID string, options *NewPartitionClientOptions) (*PartitionClient, error) {
			partitionClientsCreated++
			return newFakePartitionClient(partitionID, ""), nil
		},
	}

	processor, err := newProcessorImpl(cc, cps, &NewProcessorOptions{
		PartitionExpirationDuration: time.Hour,
	})
	require.NoError(t, err)

	consumersSyncMap := &sync.Map{}

	// to make the test easier (and less dependent on timing) we're calling through to the
	// pieces of the runImpl function
	_, err = processor.initNextClientsCh(context.Background())
	require.NoError(t, err)
	require.Empty(t, processor.nextClients)
	require.Equal(t, len(ehProps.PartitionIDs), cap(processor.nextClients))

	// the first dispatch - we have a single partition available ("a") and it gets assigned
	err = processor.dispatch(context.Background(), ehProps, consumersSyncMap)
	require.NoError(t, err)
	require.Equal(t, 1, len(processor.nextClients), "the client we created is ready to get picked up by NextPartitionClient()")

	consumers := syncMapToNormalMap(consumersSyncMap)
	origPartClient := consumers["a"]
	require.Equal(t, "a", origPartClient.partitionID)
	require.Equal(t, 1, partitionClientsCreated)

	// pull the client from the channel - it should be for the "a" partition
	procClient := processor.NextPartitionClient(context.Background())
	require.Equal(t, "a", procClient.partitionID)
	require.Empty(t, processor.nextClients)

	// the second dispatch - we reaffirm our ownership of "a" _but_ since we're already processing it no new
	// client is returned.
	err = processor.dispatch(context.Background(), ehProps, consumersSyncMap)
	require.NoError(t, err)

	// make sure we didn't create any new clients since we're already actively subscribed.
	consumers = syncMapToNormalMap(consumersSyncMap)
	afterSecondDispatchPartClient := consumers["a"]
	require.Equal(t, "a", afterSecondDispatchPartClient.partitionID)
	require.Same(t, origPartClient, afterSecondDispatchPartClient, "the client in our map is still the active one from before")
}

func TestUnit_Processor_Run_startPosition(t *testing.T) {
	cps := newCheckpointStoreForTest()

	err := cps.UpdateCheckpoint(context.Background(), Checkpoint{
		CheckpointStoreAddress: CheckpointStoreAddress{
			ConsumerGroup:           "consumer-group",
			EventHubName:            "event-hub",
			FullyQualifiedNamespace: "fqdn",
			PartitionID:             "a",
		},
		CheckpointData: CheckpointData{
			SequenceNumber: to.Ptr[int64](202),
		},
	}, nil)
	require.NoError(t, err)

	fakeConsumerClient := simpleFakeConsumerClient()

	fakeConsumerClient.newPartitionClientFn = func(partitionID string, options *NewPartitionClientOptions) (*PartitionClient, error) {
		offsetExpr, err := getOffsetExpression(options.StartPosition)
		require.NoError(t, err)

		return newFakePartitionClient(partitionID, offsetExpr), nil
	}

	processor, err := newProcessorImpl(fakeConsumerClient, cps, &NewProcessorOptions{
		PartitionExpirationDuration: time.Hour,
	})
	require.NoError(t, err)

	ehProps, err := processor.initNextClientsCh(context.Background())
	require.NoError(t, err)

	consumers := sync.Map{}
	err = processor.dispatch(context.Background(), ehProps, &consumers)
	require.NoError(t, err)

	checkpoints, err := cps.ListCheckpoints(context.Background(),
		processor.consumerClientDetails.FullyQualifiedNamespace,
		processor.consumerClientDetails.EventHubName,
		processor.consumerClientDetails.ConsumerGroup, nil)
	require.NoError(t, err)
	require.Equal(t, int64(202), *checkpoints[0].SequenceNumber)

	partClient := processor.NextPartitionClient(context.Background())
	require.Equal(t, "amqp.annotation.x-opt-sequence-number > '202'", partClient.innerClient.offsetExpression)

	err = partClient.UpdateCheckpoint(context.Background(), &ReceivedEventData{
		SequenceNumber: 405,
	})
	require.NoError(t, err)
	checkpoints, err = cps.ListCheckpoints(context.Background(),
		processor.consumerClientDetails.FullyQualifiedNamespace,
		processor.consumerClientDetails.EventHubName,
		processor.consumerClientDetails.ConsumerGroup, nil)
	require.NoError(t, err)
	require.Equal(t, 1, len(checkpoints))
	require.Equal(t, int64(405), *checkpoints[0].SequenceNumber)
}

func syncMapToNormalMap(src *sync.Map) map[string]*ProcessorPartitionClient {
	dest := map[string]*ProcessorPartitionClient{}

	src.Range(func(key, value any) bool {
		dest[key.(string)] = value.(*ProcessorPartitionClient)
		return true
	})

	return dest
}

func TestUnit_Processor_Run_cancellation(t *testing.T) {
	cps := newCheckpointStoreForTest()

	processor, err := newProcessorImpl(&fakeConsumerClient{
		details: consumerClientDetails{
			ConsumerGroup:           "consumer-group",
			EventHubName:            "event-hub",
			FullyQualifiedNamespace: "fqdn",
			ClientID:                "my-client-id",
		},
	}, cps, &NewProcessorOptions{
		PartitionExpirationDuration: time.Hour,
	})

	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// note that the cancellation here doesn't cause an error.
	err = processor.Run(ctx)
	require.NoError(t, err)
}

// updateDynamicData updates the passed in `expected` Ownership with any fields that are
// dynamically or randomly chosen. It returns the updated value.
func updateDynamicData(t *testing.T, src Ownership, expected Ownership, allPartitionIDs []string) Ownership {
	// these fields are dynamic (current time, randomly generated etag and randomly chosen partition ID) so we'll just copy them over so we can easily compare, after we validate they're
	// not bogus.
	require.NotEmpty(t, src.ETag)
	expected.ETag = src.ETag

	require.NotEqual(t, time.Time{}, src.LastModifiedTime)
	expected.LastModifiedTime = src.LastModifiedTime

	require.Contains(t, allPartitionIDs, src.PartitionID, "partition ID is randomly chosen but in our domain of partitions")
	expected.PartitionID = src.PartitionID

	return expected
}

func newProcessorForTest(t *testing.T, clientID string, cps CheckpointStore) *Processor {
	processor, err := newProcessorImpl(&fakeConsumerClient{
		details: consumerClientDetails{
			ConsumerGroup:           "consumer-group",
			EventHubName:            "event-hub",
			FullyQualifiedNamespace: "fqdn",
			ClientID:                clientID,
		},
	}, cps, &NewProcessorOptions{
		PartitionExpirationDuration: time.Hour,
	})
	require.NoError(t, err)
	return processor
}

type fakeConsumerClient struct {
	details consumerClientDetails

	getEventHubPropertiesResult EventHubProperties
	getEventHubPropertiesErr    error

	partitionClients     map[string]newMockPartitionClientResult
	newPartitionClientFn func(partitionID string, options *NewPartitionClientOptions) (*PartitionClient, error)
}

type newMockPartitionClientResult struct {
	client *PartitionClient
	err    error
}

func (cc *fakeConsumerClient) GetEventHubProperties(ctx context.Context, options *GetEventHubPropertiesOptions) (EventHubProperties, error) {
	return cc.getEventHubPropertiesResult, cc.getEventHubPropertiesErr
}

func (cc *fakeConsumerClient) NewPartitionClient(partitionID string, options *NewPartitionClientOptions) (*PartitionClient, error) {
	if cc.newPartitionClientFn != nil {
		return cc.newPartitionClientFn(partitionID, options)
	}

	if cc.partitionClients == nil {
		panic("bad test, no partition clients defined")
	}

	value, exists := cc.partitionClients[partitionID]

	if !exists {
		panic(fmt.Sprintf("bad test, partition client needed for partition %s but didn't exist in test map", partitionID))
	}

	return value.client, value.err
}

func (cc *fakeConsumerClient) getDetails() consumerClientDetails {
	return cc.details
}

func simpleFakeConsumerClient() *fakeConsumerClient {
	return &fakeConsumerClient{
		details: consumerClientDetails{
			ConsumerGroup:           "consumer-group",
			EventHubName:            "event-hub",
			FullyQualifiedNamespace: "fqdn",
			ClientID:                "my-client-id",
		},
		getEventHubPropertiesResult: EventHubProperties{
			PartitionIDs: []string{"a"},
		},
		partitionClients: map[string]newMockPartitionClientResult{
			"a": {
				client: newFakePartitionClient("a", ""),
				err:    nil,
			},
		},
	}
}

type fakeLinksForPartitionClient struct {
	internal.LinksForPartitionClient[amqpwrap.AMQPReceiverCloser]
}

func (fc *fakeLinksForPartitionClient) Retry(ctx context.Context, eventName log.Event, operation string, partitionID string, retryOptions exported.RetryOptions, fn func(ctx context.Context, lwid internal.LinkWithID[amqpwrap.AMQPReceiverCloser]) error) error {
	return nil
}

func (fc *fakeLinksForPartitionClient) Close(ctx context.Context) error {
	return nil
}

func newFakePartitionClient(partitionID string, offsetExpr string) *PartitionClient {
	return &PartitionClient{
		partitionID:      partitionID,
		offsetExpression: offsetExpr,
		links:            &fakeLinksForPartitionClient{},
	}
}

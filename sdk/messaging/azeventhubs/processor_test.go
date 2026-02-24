// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
)

func TestProcessor_PartitionsAreReqlinquished(t *testing.T) {
	res := mustCreateProcessorForTest(t, TestProcessorArgs{
		Prefix: "loadbalance",
		ProcessorOptions: &azeventhubs.ProcessorOptions{
			LoadBalancingStrategy: azeventhubs.ProcessorStrategyGreedy,
		},
	})

	hubProps, err := res.Consumer.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	ctx, stopProcessor := context.WithCancel(context.Background())
	defer stopProcessor()
	processorClosed := make(chan struct{})

	go func() {
		err := res.Processor.Run(ctx)
		require.NoError(t, err)
		close(processorClosed)
	}()

	// we expect to own all the partitions so we'll just wait until they're all claimed.
	for i := 0; i < len(hubProps.PartitionIDs); i++ {
		_ = res.Processor.NextPartitionClient(context.Background())
	}

	stopProcessor()
	<-processorClosed

	requireAllOwnershipsRelinquished(t, res)
}

func TestProcessor_Balanced(t *testing.T) {
	testWithLoadBalancer(t, azeventhubs.ProcessorStrategyBalanced)
}

func TestProcessor_Balanced_AcquisitionOnly(t *testing.T) {
	testPartitionAcquisition(t, azeventhubs.ProcessorStrategyBalanced)
}

func TestProcessor_Greedy_AcquisitionOnly(t *testing.T) {
	testPartitionAcquisition(t, azeventhubs.ProcessorStrategyGreedy)
}

func TestProcessor_Greedy(t *testing.T) {
	testWithLoadBalancer(t, azeventhubs.ProcessorStrategyGreedy)
}

func TestProcessor_Contention(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	containerName := test.RandomString("proctest", 10)
	cc, err := container.NewClient(test.URLJoinPaths(testParams.StorageEndpoint, containerName), testParams.Cred, nil)
	require.NoError(t, err)

	_, err = cc.Create(context.Background(), nil)
	require.NoError(t, err)

	defer func() {
		t.Logf("Deleting storage container")
		_, err = cc.Delete(context.Background(), nil)
		require.NoError(t, err)
	}()

	log.Printf("Producer client created")
	producerClient, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, nil)
	require.NoError(t, err)

	defer func() {
		err := producerClient.Close(context.Background())
		require.NoError(t, err)
	}()

	ehProps, err := producerClient.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	checkpointStore, err := checkpoints.NewBlobStore(cc, nil)
	require.NoError(t, err)

	type testData struct {
		name           string
		consumerClient *azeventhubs.ConsumerClient
		processor      *azeventhubs.Processor

		ctx    context.Context
		cancel context.CancelFunc
		closed chan struct{}
	}

	var processors []testData

	const numConsumers = 3

	// create a few consumer clients and processors.
	for i := 0; i < numConsumers; i++ {
		log.Printf("Consumer client %d created", i)

		consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
		require.NoError(t, err)

		// warm up the connection itself.
		_, err = consumerClient.GetEventHubProperties(context.Background(), nil)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())

		processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)
		require.NoError(t, err)

		processors = append(processors, testData{
			name:           fmt.Sprintf("ID%d", i),
			consumerClient: consumerClient,
			processor:      processor,
			ctx:            ctx,
			cancel:         cancel,
			closed:         make(chan struct{}),
		})
	}

	defer func() {
		// cancel all the processors
		for _, ps := range processors {
			ps.cancel()
			<-ps.closed
		}
	}()

	wg := sync.WaitGroup{}

	for _, client := range processors {
		wg.Add(1)

		go func(procStuff testData) {
			defer wg.Done()

			defer func() {
				err := procStuff.consumerClient.Close(context.Background())
				require.NoError(t, err)
			}()

			go func() {
				defer close(procStuff.closed)
				err := procStuff.processor.Run(procStuff.ctx)
				require.NoError(t, err)
			}()

			// we'll keep debouncing a timer for a bit - if we go 1 minute without any changes
			// to our ownership we can consider things settled.
			nextCtx, cancelNext := context.WithCancel(context.Background())
			defer cancelNext()

			// arbitrary interval, we just want to give enough time that things seem balanced.
			const idleInterval = 10 * time.Second
			active := time.AfterFunc(idleInterval, cancelNext)

			for {
				partitionClient := procStuff.processor.NextPartitionClient(nextCtx)

				if partitionClient == nil {
					break
				}

				t.Logf("%s claimed partition %s", procStuff.name, partitionClient.PartitionID())

				printOwnerships(context.Background(), t, checkpointStore, testParams, ehProps.PartitionIDs, numConsumers)

				active.Reset(time.Minute)
			}

			t.Logf("%s hasn't received a new partition in %s", procStuff.name, idleInterval)
		}(client)
	}

	wg.Wait()

	// were the partitions properly distributed?
	ownerships, err := checkpointStore.ListOwnership(context.Background(), testParams.EventHubNamespace, testParams.EventHubName, "$Default", nil)
	require.NoError(t, err)
	require.Equal(t, len(ehProps.PartitionIDs), len(ownerships))

	printOwnerships(context.Background(), t, checkpointStore, testParams, ehProps.PartitionIDs, len(ehProps.PartitionIDs))

	// check that our ownerships balanced properly
	maxAllowed := len(ehProps.PartitionIDs) / 3

	if len(ehProps.PartitionIDs)%3 > 0 {
		maxAllowed++
	}

	owners := map[string]int{}

	for _, o := range ownerships {
		owners[o.OwnerID]++
	}

	for o, numOwned := range owners {
		require.LessOrEqualf(t, numOwned, maxAllowed, "Owner %s should own max or less partitions", o)
	}
}

func TestProcessor_RunNotConcurrent(t *testing.T) {
	newProc := setupProcessorTest(t, false).Create
	proc := newProc(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	var errors int32

	wg.Add(2)

	go func() {
		defer wg.Done()

		if err := proc.Run(ctx); err != nil {
			require.EqualError(t, err, "the Processor is currently running, concurrent calls to Run() are not allowed")
			atomic.AddInt32(&errors, 1)
		}
	}()

	go func() {
		defer wg.Done()

		if err := proc.Run(ctx); err != nil {
			require.EqualError(t, err, "the Processor is currently running, concurrent calls to Run() are not allowed")
			atomic.AddInt32(&errors, 1)
		}
	}()

	wg.Wait()
	require.Equal(t, int32(1), errors)
}

func TestProcessor_Run(t *testing.T) {
	t.Run("VariationsWithoutCallingNextPartitionClient", func(t *testing.T) {
		newProc := setupProcessorTest(t, false).Create

		processor := newProc(&azeventhubs.ProcessorOptions{
			UpdateInterval:        time.Millisecond,
			LoadBalancingStrategy: azeventhubs.ProcessorStrategyBalanced,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// start the processor up and let it run for a few seconds
		err = processor.Run(ctx)
		require.NoError(t, err)

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// try to start it again (it'll fail)
		err = processor.Run(ctx)
		require.EqualError(t, err, "the Processor has been stopped. Create a new instance to start processing again")

		// and it takes precedence over the context
		err = processor.Run(ctx)
		require.EqualError(t, err, "the Processor has been stopped. Create a new instance to start processing again")
	})

	t.Run("NextPartitionClient quits properly when Run is stopped.", func(t *testing.T) {
		newProc := setupProcessorTest(t, false).Create

		processor := newProc(&azeventhubs.ProcessorOptions{
			UpdateInterval:        time.Millisecond,
			LoadBalancingStrategy: azeventhubs.ProcessorStrategyBalanced,
		})

		runProcessor := func(firstRun bool) {
			counter := newNextPartitionCounter(t, processor)
			go counter.Func()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := processor.Run(ctx)

			if firstRun {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, "the Processor has been stopped. Create a new instance to start processing again")
			}

			// now don't explicitly cancel the previous NextPartitionClient() run.
			t.Logf("Waiting for NextPartitionClient goroutine to stop")
			select {
			case <-time.After(time.Minute):
				require.Fail(t, "old goroutine didn't close after a minute")
			case <-counter.Closed:
			}

			if firstRun {
				require.NotZero(t, *counter.Called, "we should be assigned some partitions")
			} else {
				require.Zero(t, *counter.Called, "we should be assigned some partitions")
			}
		}

		runProcessor(true)
		runProcessor(false)
		runProcessor(false)
	})

	t.Run("Run stopped quickly", func(t *testing.T) {
		newProc := setupProcessorTest(t, false).Create

		processor := newProc(&azeventhubs.ProcessorOptions{
			UpdateInterval:        time.Millisecond,
			LoadBalancingStrategy: azeventhubs.ProcessorStrategyBalanced,
		})

		runProcessor := func(firstRun bool) {
			counter := newNextPartitionCounter(t, processor)
			go counter.Func()

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()

			err := processor.Run(ctx)

			if firstRun {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, "the Processor has been stopped. Create a new instance to start processing again")
			}

			// now don't explicitly cancel the previous NextPartitionClient() run.
			t.Logf("Waiting for NextPartitionClient goroutine to stop")
			select {
			case <-time.After(time.Minute):
				require.Fail(t, "old goroutine didn't close after a minute")
			case <-counter.Closed:
			}
		}

		runProcessor(true)
		runProcessor(false)
		runProcessor(false)
	})
}

func newNextPartitionCounter(t *testing.T, processor *azeventhubs.Processor) struct {
	Func   func()
	Closed chan struct{}
	Called *int
} {
	closed := make(chan struct{})
	got := 0

	return struct {
		Func   func()
		Closed chan struct{}
		Called *int
	}{
		Func: func() {
			defer close(closed)

			for {
				pc := processor.NextPartitionClient(context.Background())

				if pc == nil {
					t.Logf("partition counter: got nil, exiting loop")
					break
				}

				got++
				err = pc.Close(context.Background())
				require.NoError(t, err)
			}
		},
		Closed: closed,
		Called: &got,
	}
}

type processorTestData struct {
	Create          func(*azeventhubs.ProcessorOptions) *azeventhubs.Processor
	CheckpointStore azeventhubs.CheckpointStore
}

func setupProcessorTest(t *testing.T, geoDR bool) *processorTestData {
	testParams := test.GetConnectionParamsForTest(t)

	containerName := test.RandomString("proctest", 10)

	storageEndpoint := testParams.StorageEndpoint

	if geoDR {
		storageEndpoint = testParams.GeoDRStorageEndpoint
	}

	cc, err := container.NewClient(test.URLJoinPaths(storageEndpoint, containerName), testParams.Cred, nil)

	require.NoError(t, err)

	t.Logf("Creating storage container %s", containerName)
	_, err = cc.Create(context.Background(), nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		t.Logf("Deleting storage container")
		_, err = cc.Delete(context.Background(), nil)
		require.NoError(t, err)
	})

	// Create the checkpoint store
	// NOTE: the container must exist before the checkpoint store can be used.
	t.Logf("Checkpoint store created")
	checkpointStore, err := checkpoints.NewBlobStore(cc, nil)
	require.NoError(t, err)

	namespace, eventHubName := testParams.EventHubNamespace, testParams.EventHubName

	if geoDR {
		namespace, eventHubName = testParams.GeoDRNamespace, testParams.GeoDRHubName
	}

	t.Logf("Consumer client created")
	consumerClient, err := azeventhubs.NewConsumerClient(namespace, eventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := consumerClient.Close(context.Background())
		require.NoError(t, err)
	})

	mustCreateProcessor := func(options *azeventhubs.ProcessorOptions) *azeventhubs.Processor {
		t.Logf("Processor created")
		processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, options)
		require.NoError(t, err)

		return processor
	}

	return &processorTestData{
		Create:          mustCreateProcessor,
		CheckpointStore: checkpointStore,
	}
}

func testPartitionAcquisition(t *testing.T, loadBalancerStrategy azeventhubs.ProcessorStrategy) {
	testParams := test.GetConnectionParamsForTest(t)

	containerName := test.RandomString("proctest", 10)
	cc, err := container.NewClient(test.URLJoinPaths(testParams.StorageEndpoint, containerName), testParams.Cred, nil)
	require.NoError(t, err)

	t.Logf("Creating storage container %s", containerName)
	_, err = cc.Create(context.Background(), nil)
	require.NoError(t, err)

	defer func() {
		t.Logf("Deleting storage container")
		_, err = cc.Delete(context.Background(), nil)
		require.NoError(t, err)
	}()

	// Create the checkpoint store
	// NOTE: the container must exist before the checkpoint store can be used.
	t.Logf("Checkpoint store created")
	checkpointStore, err := checkpoints.NewBlobStore(cc, nil)
	require.NoError(t, err)

	t.Logf("Consumer client created")
	consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
	require.NoError(t, err)

	t.Logf("Processor created")
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, &azeventhubs.ProcessorOptions{
		UpdateInterval:        time.Millisecond,
		LoadBalancingStrategy: loadBalancerStrategy,
	})
	require.NoError(t, err)

	runCtx, cancelRun := context.WithCancel(context.TODO())
	defer cancelRun()

	processorClosed := make(chan struct{})

	// customer launches load balancer in a goroutine, and it continually runs
	// until they cancel the context. There is no Close() function on the Processor()
	go func() {
		defer close(processorClosed)

		t.Logf("Starting processor in separate goroutine")
		err := processor.Run(runCtx)
		require.NoError(t, err)
	}()

	// get the connection warmed up
	ehProps, err := consumerClient.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	partitionsAcquired := map[string]bool{}

	// acquire all the partitions
	for i := 0; i < len(ehProps.PartitionIDs); i++ {
		t.Logf("Waiting for next partition client")
		partitionClient := processor.NextPartitionClient(runCtx)
		require.False(t, partitionsAcquired[partitionClient.PartitionID()], "No previous client for %s", partitionClient.PartitionID())
	}

	// close all the clients.
	t.Logf("All partitions acquired and tested. Closing processor...")
	cancelRun()

	<-processorClosed
}

type countingCheckpointStore struct {
	*checkpoints.BlobStore
	listed int
}

func (c *countingCheckpointStore) ListCheckpoints(ctx context.Context, fullyQualifiedNamespace string, eventHubName string, consumerGroup string, options *azeventhubs.ListCheckpointsOptions) ([]azeventhubs.Checkpoint, error) {
	c.listed++
	return c.BlobStore.ListCheckpoints(ctx, fullyQualifiedNamespace, eventHubName, consumerGroup, options)
}

func testWithLoadBalancer(t *testing.T, loadBalancerStrategy azeventhubs.ProcessorStrategy) {
	testParams := test.GetConnectionParamsForTest(t)

	containerName := test.RandomString("proctest", 10)
	cc, err := container.NewClient(test.URLJoinPaths(testParams.StorageEndpoint, containerName), testParams.Cred, nil)
	require.NoError(t, err)

	t.Logf("Creating storage container %s", containerName)
	_, err = cc.Create(context.Background(), nil)
	require.NoError(t, err)

	defer func() {
		t.Logf("Deleting storage container")
		_, err = cc.Delete(context.Background(), nil)
		require.NoError(t, err)
	}()

	// Create the checkpoint store
	// NOTE: the container must exist before the checkpoint store can be used.
	t.Logf("Checkpoint store created")
	var checkpointStore *countingCheckpointStore
	{
		inner, err := checkpoints.NewBlobStore(cc, nil)
		require.NoError(t, err)

		checkpointStore = &countingCheckpointStore{BlobStore: inner, listed: 0}
	}

	t.Logf("Consumer client created")
	consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
	require.NoError(t, err)

	t.Logf("Processor created")
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, &azeventhubs.ProcessorOptions{
		UpdateInterval:        time.Millisecond,
		LoadBalancingStrategy: loadBalancerStrategy,
	})
	require.NoError(t, err)

	// get the connection warmed up
	ehProps, err := consumerClient.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	producerClient, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, nil)
	require.NoError(t, err)

	defer func() {
		err := producerClient.Close(context.Background())
		require.NoError(t, err)
	}()

	runCtx, cancelRun := context.WithCancel(context.TODO())
	defer cancelRun()

	go func() {
		defer cancelRun()

		wg := sync.WaitGroup{}
		partitionsAcquired := map[string]bool{}

		// acquire all the partitions
		for i := 0; i < len(ehProps.PartitionIDs); i++ {
			t.Logf("Waiting for next partition client")
			partitionClient := processor.NextPartitionClient(runCtx)

			wg.Add(1)

			require.False(t, partitionsAcquired[partitionClient.PartitionID()], "No previous client for %s", partitionClient.PartitionID())

			go func(client *azeventhubs.ProcessorPartitionClient) {
				defer wg.Done()
				err := processEventsForTest(t, producerClient, client)
				require.NoError(t, err)
			}(partitionClient)
		}

		wg.Wait()

		// close all the clients.
		t.Logf("All partitions acquired and tested. Closing processor...")
	}()

	t.Logf("Starting processor in separate goroutine")
	err = processor.Run(runCtx)
	require.NoError(t, err)

	if loadBalancerStrategy == azeventhubs.ProcessorStrategyGreedy {
		// with greedy we do one dispatch, so all 'x' number of partition clients will use/reuse
		// the one ListCheckpoints() call's results.
		require.Equal(t, 1, checkpointStore.listed)
	} else {
		// with balanced, we do a dispatch round _per_ partition, so we'll see 'x' number of requests
		require.Equal(t, len(ehProps.PartitionIDs), checkpointStore.listed)
	}
}

func processEventsForTest(t *testing.T, producerClient *azeventhubs.ProducerClient, partitionClient *azeventhubs.ProcessorPartitionClient) error {
	// initialize any resources needed to process the partition
	// This is the equivalent to PartitionOpen
	t.Logf("goroutine started for partition %s", partitionClient.PartitionID())

	const expectedEventsCount = 1000
	const batchSize = 1000
	require.Zero(t, expectedEventsCount%batchSize, "keep the math simple - even # of messages for each batch")

	// start producing events. We'll give the consumer client a moment, just to ensure
	// it's actually started up the link.
	go func() {
		time.Sleep(10 * time.Second)

		ctr := 0

		for i := 0; i < expectedEventsCount/batchSize; i++ {
			pid := partitionClient.PartitionID()
			batch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.EventDataBatchOptions{
				PartitionID: &pid,
			})
			require.NoError(t, err)

			for j := 0; j < batchSize; j++ {
				err := batch.AddEventData(&azeventhubs.EventData{
					Body: []byte(fmt.Sprintf("[%s:%d] Message", partitionClient.PartitionID(), ctr)),
				}, nil)
				require.NoError(t, err)
				ctr++
			}

			err = producerClient.SendEventDataBatch(context.Background(), batch, nil)
			require.NoError(t, err)
		}
	}()

	var allEvents []*azeventhubs.ReceivedEventData

	for {
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), 3*time.Second)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			if eventHubError := (*azeventhubs.Error)(nil); errors.As(err, &eventHubError) && eventHubError.Code == exported.ErrorCodeOwnershipLost {
				fmt.Printf("Partition %s was stolen\n", partitionClient.PartitionID())
			}

			return err
		}

		if len(events) != 0 {
			t.Logf("Processing %d event(s) for partition %s", len(events), partitionClient.PartitionID())

			allEvents = append(allEvents, events...)

			// Update the checkpoint with the last event received. If the processor is restarted
			// it will resume from this point in the partition.

			t.Logf("Updating checkpoint for partition %s", partitionClient.PartitionID())

			if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1], nil); err != nil {
				return err
			}

			if len(allEvents) == expectedEventsCount {
				t.Logf("! All events acquired for partition %s, ending...", partitionClient.PartitionID())
				return nil
			}
		}
	}
}

func printOwnerships(ctx context.Context, t *testing.T, cps azeventhubs.CheckpointStore, testParams test.ConnectionParamsForTest, partitionIDs []string, expectedConsumers int) {
	max := len(partitionIDs) / expectedConsumers

	if len(partitionIDs)%expectedConsumers > 0 {
		max++
	}

	// dump out the state of the checkpoint store so we can see how things are shaking out.
	ownerships, err := cps.ListOwnership(ctx, testParams.EventHubNamespace, testParams.EventHubName, "$Default", nil)

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return
	}

	require.NoError(t, err)

	owners := map[string][]string{}
	ownedPartitions := map[string]bool{}

	for _, o := range ownerships {
		trimmedID := o.OwnerID[0:4]
		owners[trimmedID] = append(owners[trimmedID], o.PartitionID)
		ownedPartitions[o.PartitionID] = true
	}

	sort.Strings(partitionIDs)

	var unowned []string

	for _, partID := range partitionIDs {
		if !ownedPartitions[partID] {
			unowned = append(unowned, partID)
		}
	}

	sb := strings.Builder{}

	for o, parts := range owners {
		sort.Strings(parts)
		sb.WriteString(fmt.Sprintf("  [%s (%d)] %s\n", o, len(parts), strings.Join(parts, ",")))
	}

	sb.WriteString(fmt.Sprintf("  Unowned (%d): %s\n", len(unowned), strings.Join(unowned, ",")))

	sort.Strings(partitionIDs)

	t.Logf("\nOwnerships (partitions: %d/%d), unique owners: %d, max can own: %d\n%s\n",
		len(ownedPartitions),
		len(partitionIDs),
		len(owners),
		max,
		sb.String())
}

type TestProcessorArgs struct {
	Prefix           string
	ProcessorOptions *azeventhubs.ProcessorOptions
	ConsumerOptions  *azeventhubs.ConsumerClientOptions
}

type TestProcessorResult struct {
	ContainerName   string
	TestParams      test.ConnectionParamsForTest
	ContainerClient *container.Client
	CheckpointStore azeventhubs.CheckpointStore
	Processor       *azeventhubs.Processor
	Consumer        *azeventhubs.ConsumerClient
}

func mustCreateProcessorForTest(t *testing.T, args TestProcessorArgs) TestProcessorResult {
	require.NotEmpty(t, args.Prefix)

	testParams := test.GetConnectionParamsForTest(t)

	containerName := test.RandomString(args.Prefix, 10)
	cc, err := container.NewClient(test.URLJoinPaths(testParams.StorageEndpoint, containerName), testParams.Cred, nil)
	require.NoError(t, err)

	t.Logf("Creating storage container %s", containerName)
	_, err = cc.Create(context.Background(), nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		t.Logf("Deleting storage container")
		_, err = cc.Delete(context.Background(), nil)
		require.NoError(t, err)
	})

	// Create the checkpoint store
	// NOTE: the container must exist before the checkpoint store can be used.
	t.Logf("Checkpoint store created")
	checkpointStore, err := checkpoints.NewBlobStore(cc, nil)
	require.NoError(t, err)

	t.Logf("Consumer client created")
	consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, args.ConsumerOptions)
	require.NoError(t, err)

	t.Cleanup(func() { test.RequireClose(t, consumerClient) })

	t.Logf("Processor created")
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, args.ProcessorOptions)
	require.NoError(t, err)

	return TestProcessorResult{
		CheckpointStore: checkpointStore,
		ContainerClient: cc,
		ContainerName:   containerName,
		TestParams:      testParams,
		Consumer:        consumerClient,
		Processor:       processor,
	}
}

func requireAllOwnershipsRelinquished(t *testing.T, res TestProcessorResult) {
	// now check that the ownerships exist but were all cleared out.
	ownerships, err := res.CheckpointStore.ListOwnership(context.Background(), res.TestParams.EventHubNamespace, res.TestParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
	require.NoError(t, err)

	for _, o := range ownerships {
		require.Empty(t, o.OwnerID)
	}
}

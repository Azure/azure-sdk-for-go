// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

func ProcessorStressTester(ctx context.Context) error {
	test, err := newProcessorStressTest(os.Args[2:])

	if err != nil {
		return err
	}

	defer test.Close()

	return test.Run(ctx)
}

type processorStressTest struct {
	*stressTestData

	containerName  string
	numProcessors  int
	eventsPerRound int
	rounds         int64

	prefetch     int32
	sleepAfterFn func()

	checkpointStore azeventhubs.CheckpointStore
}

func newProcessorStressTest(args []string) (*processorStressTest, error) {
	fs := flag.NewFlagSet("infiniteprocessor", flag.ContinueOnError)

	numProcessors := fs.Int("processors", 1, "Number of processors to run, concurrently")
	eventsPerRound := fs.Int("send", 5000, "Number of events to send per round")
	rounds := fs.Int64("rounds", 100, "Number of rounds. -1 means math.MaxInt64")
	prefetch := fs.Int("prefetch", 0, "Number of events to set for the prefetch. Negative numbers disable prefetch altogether. 0 uses the default for the package.")
	enableVerboseLoggingFn := addVerboseLoggingFlag(fs, nil)
	sleepAfterFn := addSleepAfterFlag(fs)

	if err := fs.Parse(args); err != nil {
		fs.PrintDefaults()
		return nil, err
	}

	enableVerboseLoggingFn()

	if *rounds == -1 {
		*rounds = math.MaxInt64
	}

	testData, err := newStressTestData("infiniteprocessor", map[string]string{
		"Processors":     fmt.Sprintf("%d", numProcessors),
		"EventsPerRound": fmt.Sprintf("%d", eventsPerRound),
		"Rounds":         fmt.Sprintf("%d", rounds),
		"Prefetch":       fmt.Sprintf("%d", *prefetch),
	})

	if err != nil {
		return nil, err
	}

	containerName := testData.runID

	storageEndpoint := test.URLJoinPaths(testData.StorageEndpoint, containerName)

	containerClient, err := container.NewClient(storageEndpoint, testData.Cred, nil)

	if err != nil {
		return nil, err
	}

	blobStore, err := checkpoints.NewBlobStore(containerClient, nil)

	if err != nil {
		return nil, err
	}

	return &processorStressTest{
		stressTestData:  testData,
		containerName:   containerName,
		numProcessors:   *numProcessors,
		eventsPerRound:  *eventsPerRound,
		rounds:          *rounds,
		checkpointStore: blobStore,
		prefetch:        int32(*prefetch),
		sleepAfterFn:    sleepAfterFn,
	}, nil
}

func (inf *processorStressTest) Run(ctx context.Context) error {
	log.Printf("======= Starting infinite processing test\n  %d processors\n  %d events sent per round\n  container name %s =======",
		inf.numProcessors,
		inf.eventsPerRound,
		inf.containerName)

	defer inf.sleepAfterFn()

	checkpoints, err := initCheckpointStore(ctx, inf.containerName, inf.stressTestData)

	if err != nil {
		return err
	}

	// start up the processors - they'll stay alive for the entire test.
	for i := 0; i < inf.numProcessors; i++ {
		cc, proc, err := inf.newProcessorForTest()

		if err != nil {
			return err
		}

		shortConsumerID := string(cc.InstanceID()[0:5])

		go func() {
			for {
				partClient := proc.NextPartitionClient(ctx)

				if partClient == nil {
					break
				}

				logger := func(format string, v ...any) {
					msg := fmt.Sprintf(format, v...)
					log.Printf("[c(%s), p(%s)]: %s", shortConsumerID, partClient.PartitionID(), msg)
				}

				go func() {
					if err := inf.receiveForever(ctx, partClient, logger, inf.eventsPerRound); err != nil {
						inf.TC.TrackException(err)
						panic(err)
					}
				}()
			}
		}()

		go func() {
			if err := proc.Run(ctx); err != nil {
				inf.TC.TrackException(err)
				panic(err)
			}
		}()
	}

	// this is the main driver for the entire test - we send, wait for the events to all be
	// accounted for, and then send again.
	producerClient, err := azeventhubs.NewProducerClient(inf.Namespace, inf.HubName, inf.Cred, nil)

	if err != nil {
		return err
	}

	defer func() { _ = producerClient.Close(context.Background()) }()

	for round := int64(0); round < inf.rounds; round++ {
		log.Printf("===== [BEGIN] Round %d/%d ===== ", round, inf.rounds)

		start := time.Now()

		endPositionsCh := make(chan azeventhubs.PartitionProperties, len(checkpoints))

		wg := sync.WaitGroup{}

		for _, cp := range checkpoints {
			wg.Add(1)

			go func(partID string) {
				defer wg.Done()
				_, ep, err := sendEventsToPartition(ctx, sendEventsToPartitionArgs{
					client:        producerClient,
					partitionID:   partID,
					messageLimit:  inf.eventsPerRound,
					testData:      inf.stressTestData,
					numExtraBytes: 1024,
				})

				if err != nil {
					inf.TC.TrackException(err)
					panic(err)
				}

				endPositionsCh <- ep
			}(cp.PartitionID)
		}

		wg.Wait()
		log.Printf("Done sending events...")
		close(endPositionsCh)

		endPositions := channelToSortedSlice(endPositionsCh, func(a, b azeventhubs.PartitionProperties) bool {
			aAsInt, _ := strconv.ParseInt(a.PartitionID, 10, 64)
			bAsInt, _ := strconv.ParseInt(b.PartitionID, 10, 64)

			return aAsInt < bAsInt
		})

		// start checking the checkpoint store to see how far along we are, and when
		// we're at the end.
		for {
			var elapsed = time.Since(start) / time.Second
			header := fmt.Sprintf("round %d, elapsed %d seconds", round, elapsed)
			output, done, err := inf.report(ctx, header, endPositions)

			if err != nil {
				log.Printf("Failed to check if partitions were balanced: %s", err.Error())
				inf.TC.TrackException(err)
			}

			if done {
				log.Printf("%s", output)
				log.Printf("!!! DONE, all partitions fully received and checkpointed.")
				break
			} else {
				log.Printf("%s", output)
			}

			<-time.After(5 * time.Second)
		}

		log.Printf("===== [END] Round %d ===== ", round)
	}

	return nil
}

func (inf *processorStressTest) receiveForever(ctx context.Context, partClient *azeventhubs.ProcessorPartitionClient, logger logf, eventsPerRound int) error {
	defer func() {
		logger("Closing")

		err := partClient.Close(context.Background())

		if err != nil {
			inf.TC.TrackException(err)
			logger("Failed when closing client: %s", err.Error())
		}
	}()

	logger("Starting receive loop")

	batchSize := int(math.Min(float64(eventsPerRound), 100))

	for {
		receiveCtx, cancelReceive := context.WithCancel(ctx)
		events, err := partClient.ReceiveEvents(receiveCtx, batchSize, nil)
		cancelReceive()

		if errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil {
			// this is fine - it just means we ran out of time waiting for events.
			// This'll happen periodically in between tests when there are no messages.
			inf.TC.TrackMetricWithProps(MetricDeadlineExceeded, 1.0, map[string]string{
				"PartitionID": partClient.PartitionID(),
			})
			continue
		}

		if ehErr := (*azeventhubs.Error)(nil); errors.As(err, &ehErr) && ehErr.Code == azeventhubs.ErrorCodeOwnershipLost {
			// this can happen as partitions are rebalanced between processors - Event Hubs
			// actually detaches us with this error.
			inf.TC.TrackMetricWithProps(MetricNameOwnershipLost, 1.0, map[string]string{
				"PartitionID": partClient.PartitionID(),
			})
			logger("Ownership lost")
			break
		}

		if err != nil {
			logger("Fatal error from ReceiveEvents: %s", err)
			inf.TC.TrackException(err)
			panic(err)
		}

		if len(events) > 0 {
			// we're okay, let's update our checkpoint
			if err := partClient.UpdateCheckpoint(ctx, events[len(events)-1], nil); err != nil {
				logger("Fatal error updating checkpoint: %s", err)
				inf.TC.TrackException(err)
				panic(err)
			}

			inf.TC.TrackMetricWithProps(MetricNameReceived, float64(len(events)), map[string]string{
				"PartitionID": partClient.PartitionID(),
			})
		}
	}

	return nil
}

func (inf *processorStressTest) Close() {
	inf.stressTestData.Close()
}

func (inf *processorStressTest) report(ctx context.Context, header string, endPositions []azeventhubs.PartitionProperties) (string, bool, error) {
	ownerships, err := inf.checkpointStore.ListOwnership(ctx, inf.Namespace, inf.HubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		return "", false, err
	}

	checkpoints, err := inf.checkpointStore.ListCheckpoints(ctx, inf.Namespace, inf.HubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		return "", false, err
	}

	ownershipMap := sliceToMap(ownerships, func(v azeventhubs.Ownership) string { return v.PartitionID })
	checkpointsMap := sliceToMap(checkpoints, func(v azeventhubs.Checkpoint) string { return v.PartitionID })

	stats := strings.Builder{}

	stats.WriteString(fmt.Sprintf("=== Stats (%s) ===\n", header))

	done := 0

	// iterate over all partitions, noting if they are unowned, how far we've gotten into the partitions, etc...
	for _, endProps := range endPositions {
		owner := "none"
		lastUpdate := ""

		o, exists := ownershipMap[endProps.PartitionID]

		if exists {
			owner = string(o.OwnerID[0:5])
			lastUpdate = o.LastModifiedTime.Format(time.RFC3339)
		}

		cp, exists := checkpointsMap[endProps.PartitionID]

		remaining := int64(-1)

		if exists {
			remaining = endProps.LastEnqueuedSequenceNumber - *cp.SequenceNumber
		}

		if remaining == 0 {
			done++
		}

		stats.WriteString(fmt.Sprintf("  [%s] o:%s (last: %s), remaining: %d/%d\n", endProps.PartitionID, owner, lastUpdate, remaining, inf.eventsPerRound))
	}

	return stats.String(), done == len(endPositions), nil
}

func sliceToMap[T any](values []T, key func(v T) string) map[string]T {
	m := map[string]T{}

	for _, v := range values {
		m[key(v)] = v
	}

	return m
}

func (inf *processorStressTest) newProcessorForTest() (*azeventhubs.ConsumerClient, *azeventhubs.Processor, error) {
	storageEndpoint := test.URLJoinPaths(inf.StorageEndpoint, inf.containerName)
	containerClient, err := container.NewClient(storageEndpoint, inf.Cred, nil)

	if err != nil {
		return nil, nil, err
	}

	cps, err := checkpoints.NewBlobStore(containerClient, nil)

	if err != nil {
		return nil, nil, err
	}

	cc, err := azeventhubs.NewConsumerClient(inf.Namespace, inf.HubName, azeventhubs.DefaultConsumerGroup, inf.Cred, nil)

	if err != nil {
		return nil, nil, err
	}

	processor, err := azeventhubs.NewProcessor(cc, cps, &azeventhubs.ProcessorOptions{
		Prefetch: inf.prefetch,
	})

	if err != nil {
		return nil, nil, err
	}

	return cc, processor, nil
}

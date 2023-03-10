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
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

/*
customMetrics
| where name == "received"
| where customDimensions["TestRunId"] == "Run-1663984812911993945"
| project timestamp, expected=toint(customDimensions["size"]), actual=valueMax
// We do batching on the client, so the "binning by x interval" is already happening.
//| summarize by bin(timestamp, 1s), valueMax
| render timechart with (title="no prefetch, exp 10000, time 1s")
*/

func getBatchTesterParams(args []string) (batchTesterParams, error) {
	params := batchTesterParams{}

	fs := flag.NewFlagSet("batch", flag.ContinueOnError)

	var batchDurationStr string

	// NOTE: these values aren't particularly special, but they do try to create a reasonable default
	// test just to make sure everything is working.
	//
	// Look in ../templates/deploy-job.yaml for some of the other parameter variations we use in stress/longevity
	// testing.
	fs.IntVar(&params.numToSend, "send", 1000000, "Number of events to send.")
	fs.IntVar(&params.batchSize, "receive", 1000, "Size to request each time we call ReceiveEvents(). Higher batch sizes will require higher amounts of memory for this test.")
	fs.StringVar(&batchDurationStr, "timeout", "60s", "Time to wait for each batch (ie: 1m, 30s, etc..)")
	prefetch := fs.Int("prefetch", 0, "Number of events to set for the prefetch. Negative numbers disable prefetch altogether. 0 uses the default for the package.")

	fs.Int64Var(&params.rounds, "rounds", 100, "Number of rounds to run with these parameters. -1 means math.MaxInt64")
	fs.IntVar(&params.paddingBytes, "padding", 1024, "Extra number of bytes to add into each message body")
	fs.StringVar(&params.partitionID, "partition", "0", "Partition ID to send and receive events to")
	fs.IntVar(&params.maxDeadlineExceeded, "maxtimeouts", 10, "Number of consecutive receive timeouts allowed before quitting")
	fs.BoolVar(&params.enableVerboseLogging, "verbose", false, "enable verbose azure sdk logging")
	sleepAfterFn := addSleepAfterFlag(fs)

	if err := fs.Parse(os.Args[2:]); err != nil {
		fs.PrintDefaults()
		return batchTesterParams{}, err
	}

	params.prefetch = int32(*prefetch)

	if params.rounds == -1 {
		params.rounds = math.MaxInt64
	}

	batchDuration, err := time.ParseDuration(batchDurationStr)

	if err != nil {
		fs.PrintDefaults()
		return batchTesterParams{}, err
	}

	params.batchDuration = batchDuration
	params.sleepAfterFn = sleepAfterFn

	return params, nil
}

// BatchStressTester sends a limited number of events and then consumes
// that set of events over and over to see what we get with different wait times.
func BatchStressTester(ctx context.Context) error {
	params, err := getBatchTesterParams(os.Args[2:])

	if err != nil {
		return err
	}

	defer params.sleepAfterFn()

	testData, err := newStressTestData("batch", params.enableVerboseLogging, map[string]string{
		"BatchDuration":       params.batchDuration.String(),
		"BatchSize":           fmt.Sprintf("%d", params.batchSize),
		"NumToSend":           fmt.Sprintf("%d", params.numToSend),
		"PaddingBytes":        fmt.Sprintf("%d", params.paddingBytes),
		"PartitionId":         params.partitionID,
		"Prefetch":            fmt.Sprintf("%d", params.prefetch),
		"Rounds":              fmt.Sprintf("%d", params.rounds),
		"Verbose":             fmt.Sprintf("%t", params.enableVerboseLogging),
		"MaxDeadlineExceeded": fmt.Sprintf("%d", params.maxDeadlineExceeded),
	})

	if err != nil {
		return err
	}

	defer testData.Close()

	log.Printf("Starting test with: batch size %d, wait time %s, prefetch: %d", params.batchSize, params.batchDuration, params.prefetch)

	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(testData.ConnectionString, testData.HubName, nil)

	if err != nil {
		return err
	}

	// we're going to read (and re-read these events over and over in our tests)
	log.Printf("Sending messages to partition %s", params.partitionID)

	sp, ep, err := sendEventsToPartition(context.Background(), sendEventsToPartitionArgs{
		client:        producerClient,
		partitionID:   params.partitionID,
		messageLimit:  params.numToSend,
		numExtraBytes: params.paddingBytes,
		testData:      testData,
	})

	closeOrPanic(producerClient)

	if err != nil {
		log.Fatalf("Failed to send events to partition %s: %s", params.partitionID, err)
	}

	log.Printf("Starting receive tests for partition %s", params.partitionID)
	log.Printf("  Start position: %#v\nEnd position: %#v", sp, ep)

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(testData.ConnectionString, testData.HubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		return err
	}

	defer closeOrPanic(consumerClient)

	// warm up the connection
	if _, err := consumerClient.GetEventHubProperties(ctx, nil); err != nil {
		log.Fatalf("Failed to warm up connection for consumer client: %s", err.Error())
	}

	for i := int64(0); i < params.rounds; i++ {
		if err := consumeForBatchTester(context.Background(), i, consumerClient, sp, params, testData); err != nil {
			log.Fatalf("Failed running round %d: %s", i, err.Error())
		}
	}

	log.Printf("Finished, check TestRunId = %s", testData.runID)
	return nil
}

type batchTesterParams struct {
	numToSend            int
	paddingBytes         int
	partitionID          string
	batchSize            int
	batchDuration        time.Duration
	rounds               int64
	prefetch             int32
	maxDeadlineExceeded  int
	enableVerboseLogging bool
	sleepAfterFn         func()
}

func consumeForBatchTester(ctx context.Context, round int64, cc *azeventhubs.ConsumerClient, sp azeventhubs.StartPosition, params batchTesterParams, testData *stressTestData) error {
	partClient, err := cc.NewPartitionClient(params.partitionID, &azeventhubs.PartitionClientOptions{
		StartPosition: sp,
		Prefetch:      params.prefetch,
	})

	if err != nil {
		return fmt.Errorf("failed to create partition client: %w", err)
	}

	defer closeOrPanic(partClient)

	log.Printf("[r:%d/%d,p:%s] Starting to receive messages from partition", round, params.rounds, params.partitionID)
	defer log.Printf("[r:%d/%d,p:%s] Done receiving messages from partition", round, params.rounds, params.partitionID)

	total := 0
	numCancels := 0
	const cancelLimit = 5

	analyzeErrorFn := func(err error) error {
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				// track these, we can use it as a proxy for "network was slow" or similar.
				testData.TC.TrackMetric(MetricDeadlineExceeded, float64(1), nil)
				numCancels++

				if numCancels >= cancelLimit {
					return fmt.Errorf("cancellation errors were received %d times in a row. Stopping test as this indicates a problem", numCancels)
				}
			} else {
				return fmt.Errorf("received %d/%d, but then got err: %w", total, params.numToSend, err)
			}
		}

		return nil
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), params.batchDuration)
		events, err := partClient.ReceiveEvents(ctx, params.batchSize, nil)
		cancel()

		if err := analyzeErrorFn(err); err != nil {
			panic(err)
		}

		testData.TC.TrackMetric(MetricReceived, float64(len(events)), nil)
		total += len(events)

		if total >= params.numToSend {
			log.Printf("[r:%d/%d,p:%s] All messages received (%d/%d)", round, params.rounds, params.partitionID, total, params.numToSend)
			break
		}
	}

	return nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/eh/stress/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/joho/godotenv"
)

const (
	endProperty       = "End"
	partitionProperty = "PartitionID"
	numProperty       = "Number"
)

type stressTestData struct {
	name  string
	runID string
	TC    *shared.TelemetryClientWrapper[Metric, Event]

	Namespace       string
	HubName         string
	StorageEndpoint string

	Cred azcore.TokenCredential
}

func (td *stressTestData) Close() {
	td.TC.TrackEvent(EventEnd)
}

type logf func(format string, v ...any)

func newStressTestData(name string, baggage map[string]string) (*stressTestData, error) {
	td := &stressTestData{
		name:  name,
		runID: fmt.Sprintf("%s-%d", name, time.Now().UnixNano()),
	}

	envFilePath := "../../../.env"

	if os.Getenv("ENV_FILE") != "" {
		envFilePath = os.Getenv("ENV_FILE")
	}

	if err := godotenv.Load(envFilePath); err != nil {
		return nil, err
	}

	var missing []string

	variables := map[string]*string{
		"EVENTHUB_NAMESPACE":               &td.Namespace,
		"EVENTHUB_NAME_STRESS":             &td.HubName,
		"CHECKPOINTSTORE_STORAGE_ENDPOINT": &td.StorageEndpoint,
	}

	for name, dest := range variables {
		val := os.Getenv(name)

		if val == "" {
			missing = append(missing, name)
		}

		*dest = val
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing environment variables (%s)", strings.Join(missing, ","))
	}

	td.TC = shared.NewTelemetryClientWrapper[Metric, Event]()

	// NOTE: this isn't run in the live testing pipelines, only within stress testing
	// so you shouldn't use the test credential.
	var err error
	td.Cred, err = azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, err
	}

	if td.TC.Context().CommonProperties == nil {
		td.TC.Context().CommonProperties = map[string]string{}
	}

	td.TC.Context().CommonProperties["TestRunId"] = td.runID
	td.TC.Context().CommonProperties["Scenario"] = td.name

	log.Printf("Name: %s, TestRunID: %s", td.name, td.runID)

	startBaggage := map[string]string{
		"Namespace": td.Namespace,
		"HubName":   td.HubName,
	}

	for k, v := range baggage {
		startBaggage[k] = v
	}

	td.TC.TrackEventWithProps(EventStart, startBaggage)

	return td, nil
}

type sendEventsToPartitionArgs struct {
	// required arguments
	client       *azeventhubs.ProducerClient
	partitionID  string
	messageLimit int

	testData *stressTestData

	// the number of extra bytes to add to the message - this helps with
	// testing conditions that require transfer times to not be instantaneous.
	// This is optional.
	numExtraBytes int
}

func sendEventsToPartition(ctx context.Context, args sendEventsToPartitionArgs) (azeventhubs.StartPosition, azeventhubs.PartitionProperties, error) {
	log.Printf("[BEGIN] Sending %d messages to partition ID %s, with messages of size %db", args.messageLimit, args.partitionID, args.numExtraBytes)

	beforeSendProps, err := args.client.GetPartitionProperties(ctx, args.partitionID, nil)

	if err != nil {
		return azeventhubs.StartPosition{}, azeventhubs.PartitionProperties{}, err
	}

	extraBytes := make([]byte, args.numExtraBytes)

	batch, err := args.client.NewEventDataBatch(context.Background(), &azeventhubs.EventDataBatchOptions{
		PartitionID: &args.partitionID,
	})

	if err != nil {
		return azeventhubs.StartPosition{}, azeventhubs.PartitionProperties{}, err
	}

	sendFn := func() error {
		if err := args.client.SendEventDataBatch(context.Background(), batch, nil); err != nil {
			return err
		}

		args.testData.TC.TrackMetricWithProps(MetricNameSent, float64(batch.NumEvents()), map[string]string{
			"PartitionID": args.partitionID,
		})

		return nil
	}

	for i := 0; i < args.messageLimit; i++ {
		ed := &azeventhubs.EventData{
			Body: extraBytes,
			Properties: map[string]any{
				numProperty:       i,
				partitionProperty: args.partitionID,
			},
		}

		if i == (args.messageLimit - 1) {
			addEndProperty(ed, int64(args.messageLimit))
		}

		err := batch.AddEventData(ed, nil)

		if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
			if batch.NumEvents() == 0 {
				return azeventhubs.StartPosition{}, azeventhubs.PartitionProperties{}, errors.New("single event was too large to fit into batch")
			}

			if err := sendFn(); err != nil {
				return azeventhubs.StartPosition{}, azeventhubs.PartitionProperties{}, err
			}

			tempBatch, err := args.client.NewEventDataBatch(context.Background(), &azeventhubs.EventDataBatchOptions{
				PartitionID: &args.partitionID,
			})

			if err != nil {
				return azeventhubs.StartPosition{}, azeventhubs.PartitionProperties{}, err
			}

			batch = tempBatch
			i-- // retry adding the same message
		} else if err != nil {
			return azeventhubs.StartPosition{}, azeventhubs.PartitionProperties{}, err
		}
	}

	if batch.NumEvents() > 0 {
		if err := sendFn(); err != nil {
			return azeventhubs.StartPosition{}, azeventhubs.PartitionProperties{}, err
		}
	}

	endProps, err := args.client.GetPartitionProperties(ctx, args.partitionID, nil)

	if err != nil {
		return azeventhubs.StartPosition{}, azeventhubs.PartitionProperties{}, err
	}

	sp := azeventhubs.StartPosition{
		Inclusive: false,
	}

	if beforeSendProps.IsEmpty {
		log.Printf("Partition %s is empty, starting sequence at 0 (not inclusive)", args.partitionID)
		sp.Earliest = to.Ptr(true)
	} else {
		log.Printf("Partition %s is NOT empty, starting sequence at %d (not inclusive)", args.partitionID, beforeSendProps.LastEnqueuedSequenceNumber)
		sp.SequenceNumber = &beforeSendProps.LastEnqueuedSequenceNumber
	}

	log.Printf("[END] Sending %d messages to partition ID %s, with messages of size %db", args.messageLimit, args.partitionID, args.numExtraBytes)

	return sp, endProps, nil
}

// initCheckpointStore creates the blob container and creates checkpoints for
// every partition so the next Processor will start from the end.
//
// Returns the checkpoints we updated, sorted by partition ID.
func initCheckpointStore(ctx context.Context, containerName string, testData *stressTestData) ([]azeventhubs.Checkpoint, error) {
	// create the container first - it shouldn't already exist
	storageEndpoint := test.URLJoinPaths(testData.StorageEndpoint, containerName)

	cc, err := container.NewClient(storageEndpoint, testData.Cred, nil)

	if err != nil {
		return nil, err
	}

	if _, err := cc.Create(ctx, nil); err != nil {
		return nil, err
	}

	cps, err := checkpoints.NewBlobStore(cc, nil)

	if err != nil {
		return nil, err
	}

	// now grab the current state of the partitions so, when the test starts up, we
	// don't read in any old data.
	producerClient, err := azeventhubs.NewProducerClient(testData.Namespace, testData.HubName, testData.Cred, nil)

	if err != nil {
		return nil, err
	}

	defer func() { _ = producerClient.Close(ctx) }()

	hubProps, err := producerClient.GetEventHubProperties(ctx, nil)

	if err != nil {
		return nil, err
	}

	var updatedCheckpoints []azeventhubs.Checkpoint

	sort.Strings(hubProps.PartitionIDs)

	for _, partitionID := range hubProps.PartitionIDs {
		partProps, err := producerClient.GetPartitionProperties(ctx, partitionID, nil)

		if err != nil {
			return nil, err
		}

		newCheckpoint := azeventhubs.Checkpoint{
			ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
			EventHubName:            testData.HubName,
			FullyQualifiedNamespace: testData.Namespace,
			PartitionID:             partitionID,
		}

		if partProps.IsEmpty {
			newCheckpoint.Offset = to.Ptr("-1")
			newCheckpoint.SequenceNumber = to.Ptr[int64](0)
		} else {
			newCheckpoint.Offset = &partProps.LastEnqueuedOffset
			newCheckpoint.SequenceNumber = &partProps.LastEnqueuedSequenceNumber
		}

		if err = cps.SetCheckpoint(ctx, newCheckpoint, nil); err != nil {
			return nil, err
		}

		updatedCheckpoints = append(updatedCheckpoints, newCheckpoint)
	}

	return updatedCheckpoints, nil
}

func addEndProperty(ed *azeventhubs.EventData, expectedCount int64) {
	ed.Properties[endProperty] = expectedCount
}

func channelToSortedSlice[T any](ch chan T, less func(i, j T) bool) []T {
	var values []T

	for v := range ch {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		return less(values[i], values[j])
	})
	return values
}

func closeOrPanic(closeable interface {
	Close(ctx context.Context) error
}) {
	if err := closeable.Close(context.Background()); err != nil {
		// TODO: there's an interesting thing happening here when I close out the connection
		// where it sometimes complains about it being idle. This is "ok" but I'd like to see
		// why EH's behavior seems different than expected.
		// Issue: https://github.com/Azure/azure-sdk-for-go/issues/19220

		var eherr *azeventhubs.Error
		if errors.As(err, &eherr) && eherr.Code == azeventhubs.ErrorCodeConnectionLost {
			// for now we'll say this is okay - it didn't interfere with the core operation
			// of the test.
			return
		}

		panic(err)
	}
}

func addSleepAfterFlag(fs *flag.FlagSet) func() {
	var durationStr string
	fs.StringVar(&durationStr, "sleepAfter", "0m", "Time to sleep after test completes")

	return func() {
		sleepAfter, err := time.ParseDuration(durationStr)

		if err != nil {
			log.Printf("Invalid sleepAfter duration given: %s", sleepAfter)
			return
		}

		if sleepAfter > 0 {
			log.Printf("Sleeping for %s", sleepAfter)
			time.Sleep(sleepAfter)
			log.Printf("Done sleeping for %s", sleepAfter)
		}
	}
}

func addVerboseLoggingFlag(fs *flag.FlagSet, customLogFn func(verbose string, e azlog.Event, s string)) func() {
	verbose := fs.String("v", "", "Enable verbose SDK logging. Valid values are test or sdk or all")

	logFn := func(e azlog.Event, s string) {
		log.Printf("[%s] %s", e, s)
	}

	if customLogFn != nil {
		logFn = func(e azlog.Event, s string) {
			customLogFn(*verbose, e, s)
		}
	}

	return func() {
		switch *verbose {
		case "":
		case "test":
			azlog.SetEvents(EventBalanceTest)
			azlog.SetListener(logFn)
		case "sdk":
			azlog.SetEvents(EventBalanceTest, azeventhubs.EventConsumer, azeventhubs.EventProducer)
			azlog.SetListener(logFn)
		case "all":
			azlog.SetListener(logFn)
		default:
			fmt.Printf("%s is not a valid logging value. Valid values are test or sdk or all", *verbose)
		}
	}
}

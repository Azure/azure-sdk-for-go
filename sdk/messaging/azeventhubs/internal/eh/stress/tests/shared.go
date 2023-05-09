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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/joho/godotenv"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

const (
	endProperty       = "End"
	partitionProperty = "PartitionID"
	numProperty       = "Number"
)

// metric names
const (
	// standard to all tests
	MetricSent          = "Sent"
	MetricReceived      = "Received"
	MetricOwnershipLost = "OwnershipLost"

	// go specific
	MetricDeadlineExceeded = "DeadlineExceeded"
)

type stressTestData struct {
	name  string
	runID string
	TC    telemetryClient

	ConnectionString        string
	Namespace               string
	HubName                 string
	StorageConnectionString string
}

func (td *stressTestData) Close() {
	td.TC.TrackEvent("end", nil)
	td.TC.Channel().Flush()
	<-td.TC.Channel().Close()
}

type logf func(format string, v ...any)

func newStressTestData(name string, verbose bool, baggage map[string]string) (*stressTestData, error) {
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
		"EVENTHUB_CONNECTION_STRING":                &td.ConnectionString,
		"EVENTHUB_NAME":                             &td.HubName,
		"CHECKPOINTSTORE_STORAGE_CONNECTION_STRING": &td.StorageConnectionString,
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

	if verbose {
		enableVerboseLogging()
	}

	tc, err := loadAppInsights()

	if err != nil {
		return nil, err
	}

	td.TC = telemetryClient{tc}

	if td.TC.Context().CommonProperties == nil {
		td.TC.Context().CommonProperties = map[string]string{}
	}

	td.TC.Context().CommonProperties["TestRunId"] = td.runID
	td.TC.Context().CommonProperties["Scenario"] = td.name

	log.Printf("Name: %s, TestRunID: %s", td.name, td.runID)

	props, err := exported.ParseConnectionString(td.ConnectionString)

	if err != nil {
		return nil, err
	}

	td.Namespace = props.FullyQualifiedNamespace

	startBaggage := map[string]string{
		"Namespace": td.Namespace,
		"HubName":   td.HubName,
	}

	for k, v := range baggage {
		startBaggage[k] = v
	}

	td.TC.TrackEvent("start", startBaggage)

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

		args.testData.TC.TrackMetric(MetricSent, float64(batch.NumEvents()), map[string]string{
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
	cc, err := container.NewClientFromConnectionString(testData.StorageConnectionString, containerName, nil)

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
	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(testData.ConnectionString, testData.HubName, nil)

	if err != nil {
		return nil, err
	}

	defer producerClient.Close(ctx)

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
			newCheckpoint.Offset = to.Ptr[int64](-1)
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

func loadAppInsights() (appinsights.TelemetryClient, error) {
	aiKey := os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY")

	if aiKey == "" {
		return nil, errors.New("missing APPINSIGHTS_INSTRUMENTATIONKEY environment variable")
	}

	config := appinsights.NewTelemetryConfiguration(aiKey)
	config.MaxBatchInterval = 5 * time.Second
	return appinsights.NewTelemetryClientFromConfig(config), nil
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

func enableVerboseLogging() {
	//azlog.SetEvents(azeventhubs.EventAuth, azeventhubs.EventConn, azeventhubs.EventConsumer)
	azlog.SetListener(func(e azlog.Event, s string) {
		log.Printf("[%s] %s", e, s)
	})
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

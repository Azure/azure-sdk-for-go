// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/conn"
	"github.com/joho/godotenv"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

type stressTestData struct {
	name  string
	runID string
	TC    appinsights.TelemetryClient

	ConnectionString        string
	Namespace               string
	HubName                 string
	StorageConnectionString string
}

func (td *stressTestData) Close() {
	td.TC.TrackEvent("end")
	td.TC.Channel().Flush()
	<-td.TC.Channel().Close()
}

func newStressTestData(name string) (*stressTestData, error) {
	td := &stressTestData{
		name:  name,
		runID: fmt.Sprintf("Run-%d", time.Now().UnixNano()),
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

	telemetryClient, err := loadAppInsights()

	if err != nil {
		return nil, err
	}

	td.TC = telemetryClient

	if telemetryClient.Context().CommonProperties == nil {
		telemetryClient.Context().CommonProperties = map[string]string{}
	}

	telemetryClient.Context().CommonProperties["TestRunId"] = td.runID
	telemetryClient.Context().CommonProperties["Scenario"] = td.name

	log.Printf("Name: %s, TestRunID: %s", td.name, td.runID)

	parsedConn, err := conn.ParsedConnectionFromStr(td.ConnectionString)

	if err != nil {
		return nil, err
	}

	td.Namespace = parsedConn.Namespace

	startEvent := appinsights.NewEventTelemetry("start")
	startEvent.Properties = map[string]string{
		"Namespace": td.Namespace,
		"HubName":   td.HubName,
	}

	telemetryClient.Track(startEvent)
	return td, nil
}

func sendEventsToPartition(ctx context.Context, producerClient *azeventhubs.ProducerClient, partitionID string, messageLimit int, numExtraBytes int) (azeventhubs.StartPosition, error) {
	log.Printf("Sending %d messages to partition ID %s, with messages of size %db", messageLimit, partitionID, numExtraBytes)

	beforeSendProps, err := producerClient.GetPartitionProperties(ctx, partitionID, nil)

	if err != nil {
		return azeventhubs.StartPosition{}, err
	}

	extraBytes := make([]byte, numExtraBytes)

	batch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
		PartitionID: &partitionID,
	})

	if err != nil {
		return azeventhubs.StartPosition{}, err
	}

	for i := 0; i < messageLimit; i++ {
		ed := &azeventhubs.EventData{
			Body: extraBytes,
			Properties: map[string]interface{}{
				numProperty:       i,
				partitionProperty: partitionID,
			},
		}

		if i == (messageLimit - 1) {
			ed.Properties[endProperty] = messageLimit
		}

		err := batch.AddEventData(ed, nil)

		if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
			if batch.NumMessages() == 0 {
				return azeventhubs.StartPosition{}, errors.New("single event was too large to fit into batch")
			}

			log.Printf("[%s] Sending batch with %d messages", partitionID, batch.NumMessages())
			if err := producerClient.SendEventBatch(context.Background(), batch, nil); err != nil {
				return azeventhubs.StartPosition{}, err
			}

			tempBatch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
				PartitionID: &partitionID,
			})

			if err != nil {
				return azeventhubs.StartPosition{}, err
			}

			batch = tempBatch
			i-- // retry adding the same message
		} else if err != nil {
			return azeventhubs.StartPosition{}, err
		}
	}

	if batch.NumMessages() > 0 {
		log.Printf("[%s] Sending last batch with %d messages", partitionID, batch.NumMessages())
		if err := producerClient.SendEventBatch(ctx, batch, nil); err != nil {
			return azeventhubs.StartPosition{}, err
		}
	}

	afterSendProps, err := producerClient.GetPartitionProperties(context.Background(), partitionID, nil)

	if err != nil {
		return azeventhubs.StartPosition{}, err
	}

	log.Printf("[%s] Sending is complete, sequence number diff: %d", partitionID, afterSendProps.LastEnqueuedSequenceNumber-beforeSendProps.LastEnqueuedSequenceNumber)

	sp := azeventhubs.StartPosition{
		Inclusive: false,
	}

	if beforeSendProps.IsEmpty {
		log.Printf("Partition %s is empty, starting sequence at 0 (not inclusive)", partitionID)
		sp.Earliest = to.Ptr(true)
	} else {
		log.Printf("Partition %s is NOT empty, starting sequence at %d (not inclusive)", partitionID, beforeSendProps.LastEnqueuedSequenceNumber)
		sp.SequenceNumber = &beforeSendProps.LastEnqueuedSequenceNumber
	}

	return sp, nil
}

func initCheckpointStore(ctx context.Context, client propertiesClient, baseAddress azeventhubs.CheckpointStoreAddress, cps azeventhubs.CheckpointStore) error {
	hubProps, err := client.GetEventHubProperties(ctx, nil)

	if err != nil {
		return err
	}

	for _, partitionID := range hubProps.PartitionIDs {
		partProps, err := client.GetPartitionProperties(ctx, partitionID, nil)

		if err != nil {
			return err
		}

		newAddress := baseAddress
		newAddress.PartitionID = partitionID

		var checkpointData azeventhubs.CheckpointData

		if partProps.IsEmpty {
			checkpointData.Offset = to.Ptr[int64](-1)
			checkpointData.SequenceNumber = to.Ptr[int64](0)
		} else {
			checkpointData.Offset = &partProps.LastEnqueuedOffset
			checkpointData.SequenceNumber = &partProps.LastEnqueuedSequenceNumber
		}

		err = cps.UpdateCheckpoint(ctx, azeventhubs.Checkpoint{
			CheckpointStoreAddress: newAddress,
			CheckpointData:         checkpointData,
		}, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

type propertiesClient interface {
	GetEventHubProperties(ctx context.Context, options *azeventhubs.GetEventHubPropertiesOptions) (azeventhubs.EventHubProperties, error)
	GetPartitionProperties(ctx context.Context, partitionID string, options *azeventhubs.GetPartitionPropertiesOptions) (azeventhubs.PartitionProperties, error)
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

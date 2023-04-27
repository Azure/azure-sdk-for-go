// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// This tool lets you consume events from a single partition using the ProducerClient.
// The PartitionClient does not do checkpointing and can only consume from a single
// partition at a time. Look at the "ehprocessor" tool, which uses the Processor.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

func main() {
	if err := partitionCmd(os.Args[:]); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func printPartitionExamples() {
	fmt.Fprintf(os.Stderr, "\n"+
		"Examples for partition:\n"+
		"  # Consume from after the latest event on partition \"partitionid\"\n"+
		"  ehpartition -namespace <your event hub namespace>. + servicebus.windows.net -eventhub tests -partition \"partitionid\"\n"+
		"\n"+
		"  # Consume including the latest event on partition \"partitionid\"\n"+
		"  ehpartition -namespace <your event hub namespace>. + servicebus.windows.net -eventhub tests -partition \"partitionid\" -start \"@latest\" -inclusive\n"+
		"\n"+
		"  # Consume from the beginning of partition \"partitionid\"\n"+
		"  ehpartition -namespace <your event hub namespace>. + servicebus.windows.net -eventhub tests -partition \"partitionid\" -start \"@earliest\"\n")
}

// partitionCmd handles receiving from a single partition using a PartitionClient
func partitionCmd(commandLineArgs []string) error {
	fs := flag.NewFlagSet("partition", flag.ContinueOnError)

	eventHubNamespace := fs.String("namespace", "", "The fully qualified hostname of your Event Hub namespace (ex: <your event hub>.servicebus.windows.net)")
	eventHubName := fs.String("eventhub", "", "The name of your Event Hub")
	eventHubConsumerGroup := fs.String("consumergroup", azeventhubs.DefaultConsumerGroup, "The Event Hub consumer group used by your application")
	eventHubOwnerLevel := fs.Int64("ownerlevel", -1, "The owner level of your consumer")
	partitionID := fs.String("partition", "", "Partition ID to receive events from")

	startPositionStr := fs.String("start", "@latest", "Start position: @latest or @earliest or o:<offset> or s:<sequence number>")
	startInclusive := fs.Bool("inclusive", false, "Include the event pointed to by the start position")

	maxBatchWaitTime := fs.Duration("wait", 30*time.Second, "Max wait time for events, per batch")
	maxBatchSize := fs.Int("count", 1, "Maximum number of events to receive, per batch")

	if err := fs.Parse(commandLineArgs); err != nil {
		printPartitionExamples()
		return err
	}

	if *eventHubName == "" || *eventHubNamespace == "" || *eventHubConsumerGroup == "" || *partitionID == "" {
		fs.PrintDefaults()
		printPartitionExamples()
		return errors.New("missing command line arguments")
	}

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return err
	}

	startPosition, startPosDesc, err := calculateStartPosition(*startPositionStr, *startInclusive)

	if err != nil {
		return err
	}

	// Using an owner level lets you control exclusivity when consuming a partition.
	//
	// See the PartitionClientOptions.OwnerLevel field for more details: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#PartitionClientOptions
	ownerLevelDesc := fmt.Sprintf("%d", *eventHubOwnerLevel)

	if *eventHubOwnerLevel == -1 {
		eventHubOwnerLevel = nil
		ownerLevelDesc = "<none>"
	}

	cc, err := azeventhubs.NewConsumerClient(*eventHubNamespace, *eventHubName, *eventHubConsumerGroup, defaultAzureCred, nil)

	if err != nil {
		return err
	}

	defer cc.Close(context.Background())

	pc, err := cc.NewPartitionClient(*partitionID, &azeventhubs.PartitionClientOptions{
		StartPosition: startPosition,
		OwnerLevel:    eventHubOwnerLevel,
	})

	if err != nil {
		return err
	}

	log.Printf("Processing events from partition %s, %s, owner level: %s", *partitionID, startPosDesc, ownerLevelDesc)
	processPartition(context.Background(), pc, *partitionID, *maxBatchWaitTime, *maxBatchSize)
	return nil
}

func calculateStartPosition(startPositionStr string, startInclusive bool) (azeventhubs.StartPosition, string, error) {
	startPosition := azeventhubs.StartPosition{
		Inclusive: startInclusive,
	}

	startPosDesc := fmt.Sprintf("Inclusive: %t", startInclusive)

	if strings.HasPrefix(startPositionStr, "s:") {
		v, err := strconv.ParseInt((startPositionStr)[2:], 10, 64)

		if err != nil {
			return azeventhubs.StartPosition{}, "", fmt.Errorf("'%s' is an invalid start position", startPositionStr)
		}

		startPosDesc = fmt.Sprintf("sequence number: %d, %s", v, startPosDesc)
		startPosition.SequenceNumber = &v
	} else if strings.HasPrefix(startPositionStr, "o:") {
		v, err := strconv.ParseInt((startPositionStr)[2:], 10, 64)

		if err != nil {
			return azeventhubs.StartPosition{}, "", fmt.Errorf("'%s' is an invalid start position", startPositionStr)
		}

		startPosDesc = fmt.Sprintf("offset: %d, %s", v, startPosDesc)
		startPosition.Offset = &v
	} else if startPositionStr == "@earliest" {
		startPosDesc = "earliest, " + startPosDesc
		startPosition.Earliest = to.Ptr(true)
	} else if startPositionStr == "@latest" {
		startPosDesc = "latest, " + startPosDesc
		startPosition.Latest = to.Ptr(true)
	} else {
		return azeventhubs.StartPosition{}, "", fmt.Errorf("'%s' is an invalid start position", startPositionStr)
	}
	return startPosition, startPosDesc, nil
}

func processPartition(ctx context.Context, pc *azeventhubs.PartitionClient, partitionID string, eventHubMaxTime time.Duration, eventHubMaxSize int) {
	defer pc.Close(ctx)

	for {
		receiveCtx, cancelReceive := context.WithTimeout(ctx, eventHubMaxTime)
		events, err := pc.ReceiveEvents(receiveCtx, eventHubMaxSize, nil)
		cancelReceive()

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			if ctx.Err() != nil { // parent cancelled
				break
			}

			// timing out without any events is fine. Continue receiving...
			continue
		} else if err != nil {
			log.Printf("ERROR while processing partition %q: %s", partitionID, err)
			break
		}

		if err := printEventsAsJSON(partitionID, events); err != nil {
			log.Printf("ERROR: %s", err)
			break
		}
	}
}

func printEventsAsJSON(partitionID string, events []*azeventhubs.ReceivedEventData) error {
	for _, evt := range events {
		var bodyBytes []int

		for _, b := range evt.Body {
			bodyBytes = append(bodyBytes, int(b))
		}

		// pick out some of the common fields
		jsonBytes, err := json.Marshal(struct {
			PartitionID    string
			MessageID      any
			BodyAsString   string
			Body           []int
			SequenceNumber int64
			Offset         int64
		}{partitionID, evt.MessageID, string(evt.Body), bodyBytes, evt.SequenceNumber, evt.Offset})

		if err != nil {
			return fmt.Errorf("Failed to marshal received event with message ID %v: %s", evt.MessageID, err.Error())
		}

		fmt.Printf("%s\n", string(jsonBytes))
	}

	return nil
}

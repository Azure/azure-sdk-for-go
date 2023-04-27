// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// This tool shows how to send events to an Event Hub, targeting a partition ID/partition key or allowing
// Event Hubs to choose the destination partition.
//
// For more information about partitioning see: https://learn.microsoft.com/azure/event-hubs/event-hubs-features#partitions

package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

func main() {
	if err := produceEventsTool(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func printProduceEventsExamples() {
	fmt.Fprintf(os.Stderr, "Examples:\n"+
		"  # Send a single event to partition with ID \"partitionid\" from STDIN\n"+
		"  echo hello | ehproducer -namespace your-event-hub-namespace.servicebus.windows.net -eventhub tests -partition \"partitionid\"\n"+
		"\n"+
		"  # Send a single event to partition with ID \"partitionid\" from a file\n"+
		"  ehproducer -namespace your-event-hub-namespace.servicebus.windows.net -eventhub tests -partition \"partitionid\" < samplemessage.txt\n"+

		"\n"+
		"  # Send multiple events to partition with ID \"partitionid\" from a file\n"+
		"  ehproducer -namespace your-event-hub-namespace.servicebus.windows.net -eventhub testing -partition \"partitionid\" < file_with_one_message_per_line.txt\n",
	)
}

func produceEventsTool() error {
	fs := flag.NewFlagSet("ehproducer", flag.ContinueOnError)

	eventHubNamespace := fs.String("namespace", "", "The fully qualified hostname of your Event Hub namespace (ex: <your event hub>.servicebus.windows.net)")
	eventHubName := fs.String("eventhub", "", "The name of your Event Hub")
	partitionKey := fs.String("partitionkey", "", "Partition key for events we send.")
	partitionID := fs.String("partition", "", "Partition ID to send events to. By default, allows Event Hubs to assign a partition")
	readMultiple := fs.Bool("multiple", false, "Whether each line of STDIN should be treated as a separate event, or if all the lines should be joined and sent as a single event")

	verbose := fs.Bool("v", false, "Enable Azure SDK verbose logging")

	if err := fs.Parse(os.Args[1:]); err != nil {
		printProduceEventsExamples()
		return err
	}

	if *eventHubNamespace == "" || *eventHubName == "" && (*partitionKey == "" || *partitionID == "") {
		fs.PrintDefaults()
		printProduceEventsExamples()
		return errors.New("Missing command line arguments")
	}

	if *verbose {
		azlog.SetEvents(azeventhubs.EventConsumer, azeventhubs.EventConn, azeventhubs.EventAuth, azeventhubs.EventProducer)
		azlog.SetListener(func(e azlog.Event, s string) {
			log.Printf("[%s] %s", e, s)
		})
	}

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return err
	}

	producerClient, err := azeventhubs.NewProducerClient(*eventHubNamespace, *eventHubName, defaultAzureCred, nil)

	if err != nil {
		return err
	}

	defer producerClient.Close(context.Background())

	batchOptions := &azeventhubs.EventDataBatchOptions{}

	if *partitionKey != "" {
		batchOptions.PartitionKey = partitionKey
	}

	if *partitionID != "" {
		batchOptions.PartitionID = partitionID
	}

	batch, err := producerClient.NewEventDataBatch(context.Background(), batchOptions)

	if err != nil {
		return err
	}

	if err := readEventsFromStdin(*readMultiple, batch); err != nil {
		return err
	}

	if err := producerClient.SendEventDataBatch(context.Background(), batch, nil); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Sent %d events, %d bytes\n", batch.NumEvents(), batch.NumBytes())
	return nil
}

func readEventsFromStdin(readMultiple bool, batch *azeventhubs.EventDataBatch) error {
	if readMultiple {
		fmt.Fprintf(os.Stderr, "Reading multiple events from stdin, one per line\n(type CTRL+d to send)...\n")
		scanner := bufio.NewScanner(os.Stdin)

		// This is a very simplified approach and will fail if the amount of messages exceeds the maximum
		// allowed size. For an example of how to handle this see this example:
		// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#example-package-ProducingEventsUsingProducerClient
		for scanner.Scan() {
			if err := batch.AddEventData(&azeventhubs.EventData{
				Body: scanner.Bytes(),
			}, nil); err != nil {
				return err
			}
		}

		return scanner.Err()
	} else {
		fmt.Fprintf(os.Stderr, "Reading a single event from stdin\n(type CTRL+d to send)...\n")
		bytes, err := io.ReadAll(os.Stdin)

		if err != nil {
			return err
		}

		return batch.AddEventData(&azeventhubs.EventData{
			Body: bytes,
		}, nil)
	}
}

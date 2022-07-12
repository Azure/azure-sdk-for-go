// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/persist"
)

func main() {
	ctx := context.Background()

	fp, err := persist.NewFilePersister(os.Getenv("EVENTHUB_FILEPERSIST_DIRECTORY"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	output, err := NewBatchWriter(fp, os.Stdout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	hub, err := eventhub.NewHubFromEnvironment(eventhub.HubWithOffsetPersistence(output))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer hub.Close(ctx)

	partitionId := os.Getenv("EVENTHUB_PARTITIONID")
	if partitionId == "" {
		parts := strings.SplitN(os.Getenv("HOSTNAME"), "-", 2)
		if len(parts) == 2 {
			partitionId = parts[1]
		} else {
			fmt.Println("EVENTHUB_PARTITIONID environment variable must be set")
			os.Exit(1)
		}
	}

	consumerGroup := os.Getenv("EVENTHUB_CONSUMERGROUP")
	if consumerGroup == "" {
		consumerGroup = "$Default"
	}

	_, err = hub.Receive(ctx, partitionId, output.HandleEvent, eventhub.ReceiveWithConsumerGroup(consumerGroup), eventhub.ReceiveWithPrefetchCount(20000))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("FATAL: ", err)
	}
}

func ExampleHub_helloWorld() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	connStr := os.Getenv("EVENTHUB_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable EVENTHUB_CONNECTION_STRING not set")
		return
	}

	hubManager, err := eventhub.NewHubManagerFromConnectionString(connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	hubEntity, err := ensureHub(ctx, hubManager, "ExampleHub_helloWorld")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a client to communicate with EventHub
	hub, err := eventhub.NewHubFromConnectionString(connStr + ";EntityPath=" + hubEntity.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = hub.Send(ctx, eventhub.NewEventFromString("Hello World!"))
	if err != nil {
		fmt.Println(err)
		return
	}

	exit := make(chan struct{})
	handler := func(ctx context.Context, event *eventhub.Event) error {
		text := string(event.Data)
		fmt.Println(text)
		exit <- struct{}{}
		return nil
	}

	for _, partitionID := range *hubEntity.PartitionIDs {
		_, err = hub.Receive(ctx, partitionID, handler)

		if err != nil {
			panic(err)
		}
	}

	// wait for the first handler to get called with "Hello World!"
	select {
	case <-exit:
		// test completed
	case <-ctx.Done():
		// test timed out
	}
	err = hub.Close(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExampleHub_webSocket() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	connStr := os.Getenv("EVENTHUB_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable EVENTHUB_CONNECTION_STRING not set")
		return
	}

	hubManager, err := eventhub.NewHubManagerFromConnectionString(connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	hubEntity, err := ensureHub(ctx, hubManager, "ExampleHub_helloWorld")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a client to communicate with EventHub
	hub, err := eventhub.NewHubFromConnectionString(connStr+";EntityPath="+hubEntity.Name, eventhub.HubWithWebSocketConnection())
	if err != nil {
		fmt.Println(err)
		return
	}

	err = hub.Send(ctx, eventhub.NewEventFromString("this message was sent and received via web socket!!"))
	if err != nil {
		fmt.Println(err)
		return
	}

	exit := make(chan struct{})
	handler := func(ctx context.Context, event *eventhub.Event) error {
		text := string(event.Data)
		fmt.Println(text)
		exit <- struct{}{}
		return nil
	}

	for _, partitionID := range *hubEntity.PartitionIDs {
		_, err = hub.Receive(ctx, partitionID, handler)

		if err != nil {
			panic(err)
		}
	}

	// wait for the first handler to get called with "Hello World!"
	select {
	case <-exit:
		// test completed
	case <-ctx.Done():
		// test timed out
	}
	err = hub.Close(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ensureHub(ctx context.Context, em *eventhub.HubManager, name string, opts ...eventhub.HubManagementOption) (*eventhub.HubEntity, error) {
	_, err := em.Get(ctx, name)
	if err == nil {
		_ = em.Delete(ctx, name)
	}

	he, err := em.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return he, nil
}

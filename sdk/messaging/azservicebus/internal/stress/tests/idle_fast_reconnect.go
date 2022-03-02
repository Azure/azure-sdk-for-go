// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func IdleFastReconnect(remainingArgs []string) {
	sc := shared.MustCreateStressContext("IdleFastReconnect")

	topicName := fmt.Sprintf("topic-%s", sc.Nano)

	startEvent := appinsights.NewEventTelemetry("Start")
	startEvent.Properties["Topic"] = topicName
	sc.Track(startEvent)

	defer sc.End()

	ac, err := admin.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("Failed to create a topic manager", err)

	_, err = ac.CreateTopic(context.Background(), topicName, nil, nil)
	sc.PanicOnError("Failed to create topic", err)

	defer func() { _, _ = ac.DeleteTopic(context.Background(), topicName, nil) }()

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, &azservicebus.ClientOptions{
		RetryOptions: azservicebus.RetryOptions{
			// NOTE: we'll _never_ use this timeout - the idle detach below
			// should use the "quick" reconnect.
			RetryDelay: time.Hour,
		},
	})

	if err != nil {
		panic(err)
	}

	sender, err := client.NewSender(topicName, nil)

	if err != nil {
		panic(err)
	}

	log.Printf("Sending first message to make sure connection is open")

	err = sender.SendMessage(context.Background(), &azservicebus.Message{
		Body: []byte("hello"),
	})

	if err != nil {
		log.Printf("%#v", err)
		panic(err)
	}

	// let the link idle out
	log.Printf("Sleeping for 11 minutes to trigger the idle link detaching")
	time.Sleep(11 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sender.SendMessage(ctx, &azservicebus.Message{
		Body: []byte("hello"),
	})

	if err != nil {
		log.Printf("%#v", err)
		panic(err)
	}

	log.Printf("Quicker reconnect worked")
}

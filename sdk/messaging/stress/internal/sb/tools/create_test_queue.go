// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tools

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
)

// CreateTempQueue creates a queue for quick experimentation that deletes itself when it's
// idle for a configurable interval.
func CreateTempQueue(remainingArgs []string) int {
	fs := flag.NewFlagSet("tempqueue", flag.ExitOnError)

	clientCreator := shared.AddAuthFlags(fs)
	durationStr := fs.String("duration", "10m", "Amount of time until the queue deletes itself due to inactivity")

	if err := fs.Parse(remainingArgs); err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		fs.PrintDefaults()
		return 1
	}

	_ = shared.LoadEnvironment()

	_, adminClient, err := clientCreator()

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return 1
	}

	duration, err := time.ParseDuration(*durationStr)

	if err != nil {
		fmt.Printf("ERROR: failed to parse duration %s", err.Error())
		return 1
	}

	autoDeleteOnIdle := utils.DurationToStringPtr(&duration)

	queueName := strings.ToLower(fmt.Sprintf("tempqueue-%X", time.Now().UnixNano()))

	_, err = adminClient.CreateQueue(context.Background(), queueName, &admin.CreateQueueOptions{
		Properties: &admin.QueueProperties{
			AutoDeleteOnIdle: autoDeleteOnIdle,
		},
	})

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return 1
	}

	fmt.Printf("Created queue '%s', will delete in %s\n", queueName, *durationStr)
	return 0
}

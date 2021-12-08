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
)

// CreateTempQueue creates a queue for quick experimentation that deletes itself when it's
// idle for a configurable interval.
func CreateTempQueue(remainingArgs []string) int {
	fs := flag.NewFlagSet("tempqueue", flag.ExitOnError)

	clientCreator := shared.AddAuthFlags(fs)
	duration := fs.String("duration", "10m", "Amount of time until the queue deletes itself due to inactivity")

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

	autoDeleteOnIdle, err := time.ParseDuration(*duration)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return 1
	}

	queueName := strings.ToLower(fmt.Sprintf("tempqueue-%X", time.Now().UnixNano()))

	_, err = adminClient.CreateQueue(context.Background(), queueName, &admin.QueueProperties{
		AutoDeleteOnIdle: &autoDeleteOnIdle,
	}, nil)

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return 1
	}

	fmt.Printf("Created queue '%s', will delete in %s\n", queueName, *duration)
	return 0
}

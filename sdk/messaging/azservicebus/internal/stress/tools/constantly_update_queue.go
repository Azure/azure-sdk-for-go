// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tools

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func ConstantlyUpdateQueue(remainingArgs []string) int {
	fs := flag.NewFlagSet("constantupdate", flag.ExitOnError)

	queueName := flag.String("name", "", "Name of the queue to update")
	interval := fs.String("interval", "5s", "How often to update the entity")
	clientCreator := shared.AddAuthFlags(fs)

	if err := fs.Parse(remainingArgs); err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		fs.PrintDefaults()
		return 1
	}

	duration, err := time.ParseDuration(*interval)

	if err != nil {
		fmt.Printf("ERROR: failed to parse interval: %s", err.Error())
		fs.PrintDefaults()
		return 1
	}

	_, adminClient, err := clientCreator()

	if err != nil {
		fs.PrintDefaults()
		return 1
	}

	if *queueName == "" {
		fmt.Printf("ERROR: No queue name given\n")
		fs.PrintDefaults()
		os.Exit(1)
	}

	ctx, cancel := shared.NewCtrlCContext()
	defer cancel()

	if err := shared.ConstantlyUpdateQueue(ctx, adminClient, *queueName, duration); err != nil {
		fmt.Printf("ERROR: %s", err)
	}

	return 0
}

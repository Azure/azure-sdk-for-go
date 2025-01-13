// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"fmt"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

func Example_enableLogging() {
	// print log output to stdout
	azlog.SetListener(printLoggedEvent)

	// pick the set of events to log
	azlog.SetEvents(
		azeventhubs.EventConn,
		azeventhubs.EventAuth,
		azeventhubs.EventProducer,
		azeventhubs.EventConsumer,
	)

	fmt.Printf("Logging enabled\n")
}

func printLoggedEvent(event azlog.Event, s string) {
	fmt.Printf("[%s] %s\n", event, s)
}

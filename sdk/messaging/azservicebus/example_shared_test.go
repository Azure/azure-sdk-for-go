// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

// Helper functions to make the main examples compilable

func exitOnError(message string, err error) {
	// these errors/failures are expected since the example
	// code, at the moment, can't run.
	if err == nil {
		return
	}

	log.Panicf("(error in example): %s: %s", message, err.Error())
}

// these just make it so our examples don't have to have a bunch of extra declarations
// for unrelated entities.
var connectionString string

var client *azservicebus.Client
var sender *azservicebus.Sender
var receiver *azservicebus.Receiver

var messages []*azservicebus.ReceivedMessage
var err error

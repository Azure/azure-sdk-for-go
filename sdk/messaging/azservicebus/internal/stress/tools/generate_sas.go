// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tools

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func GenerateSas(remainingArgs []string) int {
	fs := flag.NewFlagSet("generatesas", flag.ExitOnError)

	csVarName := fs.String("varname", "SERVICEBUS_CONNECTION_STRING", "The environment variable with a Service Bus connection string")
	durationStr := fs.String("duration", "1h", "Duration that the SAS key will be valid for")

	_ = shared.LoadEnvironment()

	// parse the connection string.
	connectionString := os.Getenv(*csVarName)

	if connectionString == "" {
		fmt.Printf("No connection string found in environment variable '%s'", *csVarName)
		fs.PrintDefaults()
		return 1
	}

	duration, err := time.ParseDuration(*durationStr)

	if err != nil {
		fmt.Printf("Duration '%s' was invalid: %s", *durationStr, err.Error())
		fs.PrintDefaults()
		return 1
	}

	cs, err := sas.CreateConnectionStringWithSASUsingExpiry(connectionString, time.Now().UTC().Add(duration))

	if err != nil {
		fmt.Printf("Failed to generate a connection with string with SAS: %s", err.Error())
		return 1
	}

	fmt.Println(cs)
	return 0
}

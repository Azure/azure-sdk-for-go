// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/tests"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/tools"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Bad command line\n")
		fmt.Printf("Usage: stress (tests|tools)\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "tests":
		tests.Run(os.Args[2:])
	case "tools":
		tools.Run(os.Args[2:])
	default:
		fmt.Printf("Usage: stress (tests|tools)\n")
		os.Exit(1)
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"fmt"
	"os"

	"stress/internal/servicebus/tests"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Printf("Bad command line\n")
		fmt.Printf("Usage: stress\n")
		os.Exit(1)
	}

	tests.Run(os.Args[1:])
}

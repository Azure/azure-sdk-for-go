// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/eh/stress/tests"
)

func main() {
	tests := []struct {
		name string
		fn   func(ctx context.Context) error
	}{
		{name: "finite", fn: tests.FiniteSendAndReceiveTest},
		{name: "finiteprocessor", fn: tests.FiniteProcessorTest},
		{name: "infinite", fn: tests.InfiniteProcessorTest},
	}

	sort.Slice(tests, func(i, j int) bool {
		return tests[i].name < tests[j].name
	})

	if len(os.Args) != 2 {
		fmt.Printf("Usage: stress <scenario-name>\n")

		fmt.Printf("Scenarios:\n")

		for _, test := range tests {
			fmt.Printf("  %s\n", test.name)
		}

		os.Exit(1)
	}

	testName := os.Args[1]

	for _, test := range tests {
		if test.name == testName {
			//azlog.SetEvents(azeventhubs.EventAuth, azeventhubs.EventConn, azeventhubs.EventConsumer)
			azlog.SetListener(func(e azlog.Event, s string) {
				//log.Printf("[%s] %s", e, s)
			})
			rand.Seed(time.Now().UnixNano())

			if err := test.fn(context.Background()); err != nil {
				fmt.Printf("ERROR: %s\n", err)
				os.Exit(1)
			}

			os.Exit(0)
		}
	}

	os.Exit(1)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

// Simple query to view some of the stats reported by these stress tests.
//
// customMetrics | where customDimensions["TestRunId"] == "169C8700A767E314"
// | project timestamp, name, valueCount
// | summarize Sum=sum(valueCount) by bin(timestamp, 30s), name
// | render timechart

func Run(remainingArgs []string) {
	// turn on some simple stderr diagnostics
	// tracer := utils.NewSimpleTracer(map[string]bool{
	// 	// tracing.SpanRecover:        true,
	// 	//tracing.SpanNegotiateClaim: true,
	// 	// tracing.SpanRecoverClient:  true,
	// 	// tracing.SpanRecoverLink:    true,
	// }, nil)

	// tab.Register(tracer)

	allTests := map[string]func(args []string){
		"constantDetach":           ConstantDetachment,
		"constantDetachmentSender": ConstantDetachmentSender,
		"emptySessions":            EmptySessions,
		"finitePeeks":              FinitePeeks,
		"finiteSendAndReceive":     FiniteSendAndReceiveTest,
		"finiteSessions":           FiniteSessions,
		"idleFastReconnect":        IdleFastReconnect,
		"infiniteSendAndReceive":   InfiniteSendAndReceiveRun,
		"longRunningRenewLock":     LongRunningRenewLockTest,
		"mostlyIdleReceiver":       MostlyIdleReceiver,
		"openCloseMeasurements":    OpenCloseMeasurements,
		"rapidOpenClose":           RapidOpenCloseTest,
		"receiveCancellation":      ReceiveCancellation,
		"sendAndReceiveDrain":      SendAndReceiveDrain,
	}

	if len(remainingArgs) == 0 {
		printUsageAndQuit(allTests)
	}

	testFn, validTestName := allTests[remainingArgs[0]]

	if !validTestName {
		printUsageAndQuit(allTests)
	}

	if err := shared.LoadEnvironment(); err != nil {
		log.Fatalf("Failed to load .env file: %s", err.Error())
	}

	testFn(remainingArgs[1:])
}

func printUsageAndQuit(allTests map[string]func(args []string)) {
	fmt.Printf("Missing/invalid test name\n")

	var names []string

	for k := range allTests {
		names = append(names, k)
	}

	sort.Strings(names)

	fmt.Printf("Usage: stress test <stress test name, listed below> \n  %s\n", strings.Join(names, "\n  "))
	os.Exit(1)
}

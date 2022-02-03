// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

// StressContext holds onto some common useful state for stress tests, including some simple stats tracking,
// a telemetry client and a context that represents the lifetime of the test itself (and will be cancelled if the user
// quits out of the stress)
type StressContext struct {
	appinsights.TelemetryClient
	*statsPrinter
	context.Context

	// TestRunID represents the test run and can be used to tie into other container metrics generated within the test cluster.
	TestRunID string

	// Nano is the nanoseconds start time for the stress test run
	Nano string

	// ConnectionString represents the value of the environment variable SERVICEBUS_CONNECTION_STRING.
	ConnectionString string

	logMessages chan string

	cancel context.CancelFunc
}

func MustCreateStressContext(testName string) *StressContext {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	if cs == "" {
		log.Fatalf("missing SERVICEBUS_CONNECTION_STRING environment variable")
	}

	aiKey := os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY")

	if aiKey == "" {
		log.Fatalf("missing APPINSIGHTS_INSTRUMENTATIONKEY environment variable")
	}

	config := appinsights.NewTelemetryConfiguration(aiKey)
	config.MaxBatchInterval = 5 * time.Second
	telemetryClient := appinsights.NewTelemetryClientFromConfig(config)

	testRunID := strings.ToLower(fmt.Sprintf("%X", time.Now().UnixNano()))

	telemetryClient.Context().CommonProperties = map[string]string{
		"Test":      testName,
		"TestRunId": testRunID,
	}

	ctx, cancel := NewCtrlCContext()

	azlog.SetEvents("azsb.Conn", "azsb.Auth", "azsb.Retry", "azsb.Mgmt")

	logMessages := make(chan string, 10000)

	go func() {
	PrintLoop:
		for {
			select {
			case <-ctx.Done():
				break PrintLoop
			case msg := <-logMessages:
				fmt.Println(msg)
			}
		}
	}()

	azlog.SetListener(func(e azlog.Event, msg string) {
		logMessages <- fmt.Sprintf("%s %10s %s", time.Now().Format(time.RFC3339), e, msg)
	})

	return &StressContext{
		TestRunID:        testRunID,
		Nano:             testRunID, // the same for now
		ConnectionString: cs,
		TelemetryClient:  telemetryClient,
		statsPrinter:     newStatsPrinter(ctx, testName, 5*time.Second, telemetryClient),
		logMessages:      logMessages,
		Context:          ctx,
		cancel:           cancel,
	}
}

func (sc *StressContext) Start(entityName string, attributes map[string]string) {
	startEvent := appinsights.NewEventTelemetry("Start")
	startEvent.Properties = map[string]string{
		"Entity": entityName,
	}

	for k, v := range attributes {
		startEvent.Properties[k] = v
	}

	log.Printf("Start: %#v", startEvent.Properties)

	sc.Track(startEvent)
}

func (sc *StressContext) End() {
	log.Printf("Stopping and flushing telemetry")

	sc.cancel()

	sc.TrackEvent("End")

	sc.Channel().Flush()
	<-sc.Channel().Close()

	time.Sleep(5 * time.Second)

	// dump out any remaining log messages
PrintLoop:
	for {
		select {
		case msg := <-sc.logMessages:
			fmt.Println(msg)
		default:
			break PrintLoop
		}
	}

	// dump out the last stats.
	sc.PrintStats()

	log.Printf("Done")
}

// PanicOnError logs, sends telemetry and then closes on error
func (tracker *StressContext) PanicOnError(message string, err error) {
	tracker.LogIfFailed(message, err, nil)

	if err != nil {
		panic(err)
	}
}

func (tracker *StressContext) Assert(condition bool, message string) {
	tracker.LogIfFailed(message, nil, nil)

	if !condition {
		panic(message)
	}
}

func (sc *StressContext) LogIfFailed(message string, err error, stats *Stats) {
	if err != nil {
		log.Printf("Error: %s: %#v, %T", message, err, err)

		if stats != nil {
			atomic.AddInt32(&stats.Errors, 1)
		}

		et := appinsights.NewExceptionTelemetry(err)
		et.Properties["Reason"] = message
		sc.Track(et)
	}
}

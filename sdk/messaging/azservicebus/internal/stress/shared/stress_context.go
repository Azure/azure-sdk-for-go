// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

// StressContext holds onto some common useful state for stress tests, including some simple stats tracking,
// a telemetry client and a context that represents the lifetime of the test itself (and will be cancelled if the user
// quits out of the stress)
type StressContext struct {
	TC      appinsights.TelemetryClient
	Context context.Context

	// TestRunID represents the test run and can be used to tie into other container metrics generated within the test cluster.
	TestRunID string

	// Nano is the nanoseconds start time for the stress test run
	Nano string

	// ConnectionString represents the value of the environment variable SERVICEBUS_CONNECTION_STRING.
	ConnectionString string

	logMessages chan string

	cancel context.CancelFunc
}

// TrackDuration tracks durations (as a metric), using the initial call to TrackDuration as the start. The duration is
// ended when you call the returned function.
// TrackDuration respects any included baggage in the context.
func TrackDuration(ctx context.Context, tc appinsights.TelemetryClient, name Metric) func(map[string]string) {
	start := time.Now()

	return func(attrs map[string]string) {
		duration := time.Since(start) / time.Millisecond
		TrackMetric(ctx, tc, name, float64(duration), attrs)
	}
}

// TrackMetric tracks metric and respects any included baggage in the context.
func TrackMetric(ctx context.Context, tc appinsights.TelemetryClient, name Metric, value float64, attrs map[string]string) {
	tc.Track(&appinsights.MetricTelemetry{
		Name:  string(name),
		Value: value,
		BaseTelemetry: appinsights.BaseTelemetry{
			Properties: UpdateBaggage(ctx, attrs),
		},
	})
}

// TrackError tracks an error (using the AppInsights exceptions table).
// TrackError respects any included baggage in the context.
//
// NOTE: this function does not consider context cancellations/deadlines as errors.
func TrackError(ctx context.Context, tc appinsights.TelemetryClient, err error) {
	// track all errors except for cancellation errors - the caller can take care of those since
	// they're the only one that knows if it's a true error or just normal behavior.
	if err != nil && !isCancelError(err) {
		log.Printf("Error: %#v, %T", err, err)

		ext := appinsights.NewExceptionTelemetry(err)
		ext.BaseTelemetry.Properties = UpdateBaggage(ctx, nil)

		tc.Track(ext)
	}
}

type StressContextOptions struct {
	// Duration is the amount of time the stress test should run before
	// the StressContext.Context expires.
	Duration time.Duration
}

func MustCreateStressContext(testName string, options *StressContextOptions) *StressContext {
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

	log.Printf("Common properties\n:%#v", telemetryClient.Context().CommonProperties)

	ctx, cancel := NewCtrlCContext()

	if options != nil && options.Duration > 0 {
		ctx, cancel = context.WithTimeout(ctx, options.Duration)
	}

	azlog.SetEvents(azservicebus.EventSender, azservicebus.EventReceiver, azservicebus.EventConn)

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

	// A little while back I had issues because the base image didn't include ca-certificates, so the SSL
	// cert for appinsights couldn't be validated. Uncommenting this will show you potential issues the
	// appinsights client has when attempting to upload telemetry.
	//
	// appinsights.NewDiagnosticsMessageListener(func(msg string) error {
	// 	fmt.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
	// 	return nil
	// })

	return &StressContext{
		TestRunID:        testRunID,
		Nano:             testRunID, // the same for now
		ConnectionString: cs,
		TC:               telemetryClient,
		// you could always change the interval here. A minute feels like often enough
		// to know things are running, while not so often that you end up flooding logging
		// with duplicate information.
		logMessages: logMessages,
		Context:     ctx,
		cancel:      cancel,
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
	sc.TC.Track(startEvent)
}

func (sc *StressContext) End() {
	log.Printf("Stopping and flushing telemetry")

	sc.cancel()

	sc.TC.TrackEvent("End")

	sc.TC.Channel().Flush()
	<-sc.TC.Channel().Close()

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

	log.Printf("Done")
}

// PanicOnError logs, sends telemetry and then closes on error
func (tracker *StressContext) PanicOnError(message string, err error) {
	tracker.LogIfFailed(message, err)

	if err != nil {
		panic(err)
	}
}

func (tracker *StressContext) Failf(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	tracker.LogIfFailed(err.Error(), err)
	panic(err)
}

func (tracker *StressContext) NoError(err error) {
	if err == nil {
		return
	}

	tracker.LogIfFailed(err.Error(), err)
	panic(err)
}

func (tracker *StressContext) NoErrorf(err error, format string, args ...any) {
	if err == nil {
		return
	}

	msg := fmt.Sprintf(format, args...)
	tracker.LogIfFailed(fmt.Sprintf("%s: %s", msg, err.Error()), err)
	panic(err)
}

func (tracker *StressContext) Assert(condition bool, message string) {
	tracker.LogIfFailed(message, nil)

	if !condition {
		panic(message)
	}
}

func (tracker *StressContext) Equal(val1 any, val2 any) {
	if val1 != val2 {
		panic(fmt.Errorf("Expected %v, got %v", val1, val2))
	}
}

func (tracker *StressContext) Nil(val1 any) {
	if val1 == nil {
		panic("value was not nil")
	}
}

func (sc *StressContext) LogIfFailed(message string, err error) {
	if err != nil {
		log.Printf("Error: %s: %#v, %T", message, err, err)
	}
}

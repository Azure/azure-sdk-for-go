// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

// OpenCloseMeasurements tests that we are able to consistently open and close our links and connections
// in a timely way. This test doesn't immediately fail, it's primary purpose is just to provide historical
// and measurable data on our performance.
//
// The origin of this test was a bug we found in go-amqp where, if the frame was too small, it wouldn't parse
// and return it until we received more data (any data, just so long as it caused our total local buffer to exceed
// 8 bytes) [PR#320]
//
// PR#320: https://github.com/Azure/go-amqp/pull/320
func OpenCloseMeasurements(remainingArgs []string) {
	type testArgs struct {
		SleepDuration time.Duration
		MessageCount  int
		BodySize      int
	}

	fs := flag.NewFlagSet("args", flag.PanicOnError)
	numRounds := fs.Int("rounds", 10, "The number of rounds of sends and closes to run")
	_ = fs.Parse(remainingArgs)

	fn := func(args testArgs) {
		sc := shared.MustCreateStressContext("OpenCloseMeasurements", &shared.StressContextOptions{
			CommonBaggage: map[string]string{
				"SleepDuration": args.SleepDuration.String(),
				"MessageCount":  strconv.FormatInt(int64(args.MessageCount), 10),
				"BodySize":      strconv.FormatInt(int64(args.BodySize), 10),
			},
			EmitStartEvent: true,
		})

		defer sc.End()

		queueName := fmt.Sprintf("OpenCloseMeasurements-%s", sc.Nano)
		_ = shared.MustCreateAutoDeletingQueue(sc, queueName, &admin.QueueProperties{})

		client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
		sc.PanicOnError("failed to create client", err)

		trackingSender, err := shared.NewTrackingSender(sc.TC, client, queueName, nil)
		sc.PanicOnError("failed to create sender", err)

		log.Printf("Sending message to warm up connection and links.")

		body := make([]byte, args.BodySize)

		for i := 0; i < args.MessageCount; i++ {
			err = trackingSender.SendMessage(context.Background(), &azservicebus.Message{
				Body: body,
			}, nil)
			sc.NoErrorf(err, "failed to send message %d", i)
		}

		log.Printf("Sleeping for %s, done at %s...", args.SleepDuration, time.Now().Add(args.SleepDuration))
		time.Sleep(args.SleepDuration)

		log.Printf("Done sleeping, now attempting to close link")
		// the error is reported for now to metrics - not going to kill this as we have a bug where
		// the "detach because idle" error comes back from Close() right now.

		start := time.Now()
		max := 10 * time.Second
		_ = trackingSender.Close(context.Background())

		if time.Since(start) > max {
			sc.PanicOnError("Slow close", fmt.Errorf("Took longer than %s", max))
		}
	}

	// some simple cases
	testCases := []testArgs{
		{1 * time.Minute, 1, 10},
		{5 * time.Minute, 100, 100},
		{5 * time.Minute, 100, 10000},
		{5 * time.Minute, 1, 10},
		{10*time.Minute + 30*time.Second, 1, 10},
		{11 * time.Minute, 1, 10},
	}

	for i := 0; i < *numRounds; i++ {
		wg := sync.WaitGroup{}

		for _, args := range testCases {
			wg.Add(1)

			go func(args testArgs) {
				defer wg.Done()
				fn(args)
			}(args)
		}

		wg.Wait()
	}
}

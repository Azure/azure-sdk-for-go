// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

// EmptySessions attempts to get the next available session from a session-enabled subscription that has no
// available sessions. This makes sure the server-enforced timeout works properly in our code over a long
// period of time.
func EmptySessions(remainingArgs []string) {
	// Queries:
	// customMetrics | where customDimensions['TestRunId'] == '17402ab9210e0340' and name == "SessionTimeoutMS" | project timestamp, seconds=valueMax/1000 | render timechart

	params := struct {
		numAttempts int
		rounds      int
	}{}

	fs := flag.NewFlagSet("emptysessions", flag.PanicOnError)

	fs.IntVar(&params.numAttempts, "sessions", 2000, "Number of attempts to get a session")
	fs.IntVar(&params.rounds, "rounds", 100, "Number of rounds to run with these parameters. -1 means math.MaxInt64")

	topicName := strings.ToLower(fmt.Sprintf("topic-%X", time.Now().UnixNano()))
	log.Printf("Creating topic %s", topicName)

	sc := shared.MustCreateStressContext("emptysessions")
	defer sc.End()

	cleanup := shared.MustCreateSubscriptions(sc, topicName, []string{"sub1"}, &shared.MustCreateSubscriptionsOptions{
		Subscription: &admin.CreateSubscriptionOptions{
			Properties: &admin.SubscriptionProperties{
				RequiresSession: to.Ptr(true),
			},
		},
	})
	defer cleanup()

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, &azservicebus.ClientOptions{
		RetryOptions: azservicebus.RetryOptions{
			MaxRetries: -1,
		},
	})
	sc.NoError(err)

	defer func() {
		sc.NoError(client.Close(context.Background()))
	}()

	for round := 0; round < params.rounds; round++ {
		for i := 0; i < params.numAttempts; i++ {
			// the default "wait for next session" timeout is basically a minute. If we exceed that then something is broken.
			// we give it a little bit of time in case the connection needed to be started/recovered as well.
			acceptCtx, cancelAccept := context.WithTimeout(context.Background(), 4*time.Minute)

			start := time.Now()

			sessionReceiver, err := client.AcceptNextSessionForSubscription(acceptCtx, topicName, "sub1", nil)
			cancelAccept()
			sc.Nil(sessionReceiver)

			endMS := time.Since(start) / time.Millisecond

			sc.TrackMetric(string(MetricNameSessionTimeoutMS), float64(endMS))

			// the error should indicate that we timed out waiting for a new session
			if sbErr := (*azservicebus.Error)(nil); errors.As(err, &sbErr) {
				switch sbErr.Code {
				case azservicebus.CodeConnectionLost:
					// these are okay, we'll just let it recover on the next call.
					sc.TrackMetric(string(MetricNameConnectionLost), 1)
				case azservicebus.CodeTimeout:
					// great! this is what we expect to get - the service-side timeout fires off since there's
					// no session available.
				default:
					sc.Failf("Unexpected service code '%s'", sbErr.Code)
				}
			} else if errors.Is(err, context.DeadlineExceeded) {
				// this means we gave the service a max timeout (passed in via
				// the link attach request) and there was some issue
				sc.Failf("Deadline exceeded, session timeout took too long (%dms)", endMS)
			} else if err != nil {
				sc.PanicOnError("A non-timeout error occurred", err)
			}
		}
	}
}

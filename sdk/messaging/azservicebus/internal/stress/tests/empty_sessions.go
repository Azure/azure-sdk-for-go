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

func EmptySessions(remainingArgs []string) {
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
			// barrel through retryable failures - we're specifically trying to see if we ever stop
			// receiving the standard "no session within time limit" error.
			MaxRetries: 10,
		},
	})
	sc.NoError(err)

	defer func() {
		sc.NoError(client.Close(context.Background()))
	}()

	for round := 0; round < params.rounds; round++ {
		for i := 0; i < params.numAttempts; i++ {
			// the default "wait for next session" timeout is basically a minute. If we exceed that then something is broken.
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			sessionReceiver, err := client.AcceptNextSessionForSubscription(ctx, topicName, "sub1", nil)
			cancel()
			sc.Nil(sessionReceiver)

			// the error should indicate that we timed out waiting for a new session
			if sbErr := (*azservicebus.Error)(nil); errors.As(err, &sbErr) {
				sc.Equal(azservicebus.CodeTimeout, sbErr.Code)
			} else if err != nil {
				sc.PanicOnError("A non-timeout error occurred", err)
			}
		}
	}
}

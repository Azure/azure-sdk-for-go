// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

type finiteSessionsArgs struct {
	numSessions int
	rounds      int
}

func FiniteSessions(remainingArgs []string) {
	// NOTE: these values aren't particularly special, but they do try to create a reasonable default
	// test just to make sure everything is working.
	//
	// Look in ../templates/deploy-job.yaml for some of the other parameter variations we use in stress/longevity
	// testing.
	fs := flag.NewFlagSet("FiniteSessions", flag.ContinueOnError)

	params := finiteSessionsArgs{}

	fs.IntVar(&params.numSessions, "sessions", 2000, "Number of sessions to test")
	fs.IntVar(&params.rounds, "rounds", 300, "Number of rounds to run with these parameters. -1 means math.MaxInt64")

	sc := shared.MustCreateStressContext("FiniteSessions", nil)
	defer sc.End()

	topicName := strings.ToLower(fmt.Sprintf("topic-%X", time.Now().UnixNano()))

	sc.Start(topicName, map[string]string{
		"NumSessions": fmt.Sprintf("%d", params.numSessions),
		"Rounds":      fmt.Sprintf("%d", params.rounds),
	})

	log.Printf("Creating topic %s", topicName)

	cleanup := shared.MustCreateSubscriptions(sc, topicName, []string{"sub1"}, &shared.MustCreateSubscriptionsOptions{
		Subscription: &admin.CreateSubscriptionOptions{
			Properties: &admin.SubscriptionProperties{
				RequiresSession: to.Ptr(true),
			},
		},
	})
	defer cleanup()

	client, err := azservicebus.NewClient(sc.Endpoint, sc.Cred, nil)
	sc.NoError(err)

	sender, err := shared.NewTrackingSender(sc.TC, client, topicName, nil)
	sc.NoError(err)

	defer sender.Close(sc.Context)

	for round := 0; round < int(params.rounds); round++ {
		var sessionReceivers []*azservicebus.SessionReceiver

		wg := sync.WaitGroup{}

		for i := 0; i < params.numSessions; i++ {
			sessionID := fmt.Sprintf("%d:%d", round, i)

			err = sender.SendMessage(sc.Context, &azservicebus.Message{
				SessionID: &sessionID,
			}, nil)
			sc.NoError(err)
			shared.TrackMetric(sc.Context, sc.TC, shared.MetricMessagesSent, 1, map[string]string{
				"SessionID": sessionID,
			})

			sessionReceiver, err := client.AcceptNextSessionForSubscription(sc.Context, topicName, "sub1", nil)
			sc.NoError(err)
			shared.TrackMetric(sc.Context, sc.TC, shared.MetricSessionAccept, 1, nil)

			// one of the things mentioned in the customer issue - they keep the session receivers
			// alive for a long time.
			sessionReceivers = append(sessionReceivers, sessionReceiver)

			wg.Add(1)

			go func() {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(sc.Context, time.Minute)

				messages, err := sessionReceiver.ReceiveMessages(ctx, 2, nil)
				cancel()

				sc.NoError(err)
				sc.Equal(1, len(messages))
				sc.Equal(sessionID, *messages[0].SessionID)

				shared.TrackMetric(ctx, sc.TC, shared.MetricMessageReceived, float64(len(messages)), nil)

				trackDuration := shared.TrackDuration(ctx, sc.TC, shared.MetricSettlementRequestDuration)
				sc.NoError(sessionReceiver.CompleteMessage(sc.Context, messages[0], nil))
				trackDuration(nil)
			}()
		}

		wg.Wait()

		for _, receiver := range sessionReceivers {
			err = receiver.Close(sc.Context)
			sc.NoErrorf(err, "No errors when session receiver is closed")
		}
	}
}

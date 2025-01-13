// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func ReceiveCancellation(remainingArgs []string) {
	const rounds = 2000

	sc := shared.MustCreateStressContext("FinitePeeks", nil)
	defer sc.End()

	queueName := fmt.Sprintf("finite-peeks-%s", sc.Nano)
	sc.Start(queueName, map[string]string{
		"Rounds": fmt.Sprintf("%d", rounds),
	})

	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClient(sc.Endpoint, sc.Cred, nil)
	sc.PanicOnError("failed to create client", err)

	for i := 0; i < rounds; i += 100 {
		func() {
			receiver, err := shared.NewTrackingReceiverForQueue(sc.TC, client, queueName, nil)
			sc.PanicOnError("failed to create receiver", err)

			defer receiver.Close(context.Background())

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i)*time.Millisecond)
			defer cancel()

			// from a cold receiver link
			_, err = receiver.ReceiveMessages(ctx, 95, nil)

			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				sc.PanicOnError("failed to receive messages (1)", err)
			}

			shared.TrackMetric(sc.Context, sc.TC, shared.MetricStressSuccessfulCancels, float64(1), map[string]string{
				"Type": "cold",
			})

			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(i)*time.Millisecond)
			defer cancel()

			// and one more time, now that the link has been warmed up
			_, err = receiver.ReceiveMessages(ctx, 95, nil)

			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				sc.PanicOnError("failed to receive messages (2)", err)
			}

			shared.TrackMetric(sc.Context, sc.TC, shared.MetricStressSuccessfulCancels, float64(1), map[string]string{
				"Type": "warm",
			})
		}()
	}
}

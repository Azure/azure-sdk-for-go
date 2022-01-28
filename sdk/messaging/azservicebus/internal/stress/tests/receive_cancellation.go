// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func ReceiveCancellation(remainingArgs []string) {
	sc := shared.MustCreateStressContext("FinitePeeks")
	defer sc.Done()

	queueName := fmt.Sprintf("finite-peeks-%s", sc.Nano)
	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	for i := 0; i < 2000; i += 100 {
		func() {
			receiver, err := client.NewReceiverForQueue(queueName, nil)
			sc.PanicOnError("failed to create receiver", err)

			defer receiver.Close(context.Background())

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i)*time.Millisecond)
			defer cancel()

			// from a cold receiver link
			_, err = receiver.ReceiveMessages(ctx, 95, nil)
			sc.PanicOnError("failed to receive messages (1)", err)

			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(i)*time.Millisecond)
			defer cancel()

			// and one more time, now that the link has been warmed up
			_, err = receiver.ReceiveMessages(ctx, 95, nil)
			sc.PanicOnError("failed to receive messages (2)", err)
		}()
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func LongRunningRenewLockTest(remainingArgs []string) {
	sc := shared.MustCreateStressContext("LongRunningRenewLockTest", nil)
	defer sc.End()

	queueName := fmt.Sprintf("renew-lock-test-%s", sc.Nano)
	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create admin.Client", err)

	sender, err := shared.NewTrackingSender(sc.TC, client, queueName, nil)
	sc.PanicOnError("failed to create Sender", err)

	err = sender.SendMessage(context.Background(), &azservicebus.Message{
		Body: []byte("ping"),
	}, nil)
	sc.PanicOnError("failed to send message", err)

	receiver, err := shared.NewTrackingReceiverForQueue(sc.TC, client, queueName, nil)
	sc.PanicOnError("failed to create receiver", err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	sc.PanicOnError("failed to receive messages", err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()

	now := time.Now()
	i := 0

	for {
		if i != 0 && i%50 == 0 {
			log.Printf("Renewed %d times, for %d minutes", i, time.Since(now)/time.Minute)
		}

		i++
		err := receiver.RenewMessageLock(ctx, messages[0], nil)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Cancellation/deadline exceeded. Can stop, will complete message now")

			err = receiver.CompleteMessage(context.Background(), messages[0], nil)
			sc.PanicOnError("failed to complete message", err)
			break
		}

		sc.PanicOnError("failed to renew message lock", err)

		time.Sleep(5 * time.Second)
	}
}

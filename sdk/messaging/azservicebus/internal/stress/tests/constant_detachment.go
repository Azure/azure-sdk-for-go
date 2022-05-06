// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func ConstantDetachment(remainingArgs []string) {
	sc := shared.MustCreateStressContext("ConstantDetachment")
	defer sc.End()

	queueName := fmt.Sprintf("detach-tester-%s", sc.Nano)

	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	stats := sc.NewStat("stats")

	sender, err := client.NewSender(queueName, nil)
	sc.PanicOnError("create a sender", err)

	const maxMessages = 10000

	shared.MustGenerateMessages(sc, sender, maxMessages, 1024, stats)

	// now attempt to receive messages, while constant detaching is happening
	receiver, err := client.NewReceiverForQueue(queueName, &azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModeReceiveAndDelete,
	})
	sc.PanicOnError("failed to create receiver", err)

	adminClient, err := admin.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create admin client", err)

	go func() {
		// keep updating the definition of the queue, which will cause the service to detach us.
		err := shared.ConstantlyUpdateQueue(sc.Context, adminClient, queueName, 30*time.Second)

		if errors.Is(err, context.Canceled) {
			return
		}

		sc.PanicOnError("constantly updated queue failed", err)
	}()

InfiniteLoop:
	for stats.Received != maxMessages {
		select {
		case <-sc.Done():
			break InfiniteLoop
		default:
		}

		messages, err := receiver.ReceiveMessages(sc.Context, 10, nil)

		if err != nil {
			sc.LogIfFailed("receive failed, continuing", err, stats)
			continue
		}

		stats.AddReceived(int32(len(messages)))
	}

	if stats.Received != maxMessages {
		sc.PanicOnError("constantdetach failed", fmt.Errorf("not all messages received (got %d, wanted %d)", stats.Received, maxMessages))
	}
}

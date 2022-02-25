// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func ConstantDetachment(remainingArgs []string) {
	sc := shared.MustCreateStressContext("ConstantDetachment")

	queueName := fmt.Sprintf("detach-tester-%s", sc.Nano)

	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	senderStats := sc.NewStat("sender")
	receiverStats := sc.NewStat("receiver")

	sender, err := client.NewSender(queueName, nil)
	sc.PanicOnError("create a sender", err)

	const maxMessages = 5000

	shared.MustGenerateMessages(sc, sender, maxMessages, 1024, senderStats)

	// now attempt to receive messages, while constant detaching is happening
	receiver, err := client.NewReceiverForQueue(queueName, &azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModeReceiveAndDelete,
	})
	sc.PanicOnError("failed to create receiver", err)

	adminClient, err := admin.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create admin client", err)

	go func() {
		err := shared.ConstantlyUpdateQueue(sc.Context, adminClient, queueName, 30*time.Second)
		sc.PanicOnError("constantly updated queue failed", err)
	}()

InfiniteLoop:
	for receiverStats.Received != maxMessages {
		select {
		case <-sc.Done():
			break InfiniteLoop
		default:
		}

		log.Printf("Waiting for message")
		messages, err := receiver.ReceiveMessages(sc.Context, 1, nil)

		if err != nil {
			sc.LogIfFailed("receive failed, continuing", err, receiverStats)
			continue
		}

		receiverStats.AddReceived(int32(len(messages)))

		// TODO: this easily proves that the 410 error (where a lock is lost) is a problem for amqp-common since it just continually retries, wasting time since
		// TODO: lock lost is _FATAL_.
		// This is covered here: https://github.com/Azure/azure-sdk-for-go/issues/16088
		//TODO: if you run this with receiveAndDelete mode our recovery is _perfect_.
		// for _, msg := range messages {
		// 	receiver.CompleteMessage(sc.Context, msg)

		// 	if sc.LogOnError("failed to complete message", err, receiverStats) == nil {
		// 		receiverStats.AddReceived(1)
		// 	}
		// }
	}
}

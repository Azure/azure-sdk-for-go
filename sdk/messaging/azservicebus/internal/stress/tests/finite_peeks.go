// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func FinitePeeks(remainingArgs []string) {
	sc := shared.MustCreateStressContext("FinitePeeks")
	defer sc.Done()

	queueName := fmt.Sprintf("finite-peeks-%s", sc.Nano)
	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	receiverStats := sc.NewStat("peeks")

	sender, err := client.NewSender(queueName, nil)
	sc.PanicOnError("failed to create sender", err)

	err = sender.SendMessage(sc.Context, &azservicebus.Message{
		Body: []byte("peekable message"),
	})
	sc.PanicOnError("failed to send message", err)

	_ = sender.Close(sc.Context)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	sc.PanicOnError("failed to create receiver", err)

	// receiving here just guarantees the message has arrived and is available (sometimes
	// there's a slight delay)
	receiveCtx, cancel := context.WithTimeout(sc.Context, time.Minute)
	defer cancel()

	tmp, err := receiver.ReceiveMessages(receiveCtx, 1, nil)
	sc.PanicOnError("failed to receive messages", err)
	sc.Assert(len(tmp) == 1, "message was never available")

	// return the message back from whence it came.
	sc.PanicOnError("failed to abandon message",
		receiver.AbandonMessage(sc.Context, tmp[0], nil))

	for i := 0; i < 10000; i++ {
		log.Printf("Sleeping for 1 second before iteration %d", i)
		time.Sleep(time.Second)

		seqNum := int64(0)

		messages, err := receiver.PeekMessages(sc.Context, 1, &azservicebus.PeekMessagesOptions{
			FromSequenceNumber: &seqNum,
		})
		sc.PanicOnError("failed to peek messages", err)
		sc.Assert(len(messages) == 1, "no messages returned in peek")

		receiverStats.AddReceived(int32(1))
	}
}

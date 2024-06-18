// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func ConstantDetachmentSender(remainingArgs []string) {
	sc := shared.MustCreateStressContext("ConstantDetachmentSender", nil)
	defer sc.End()

	adminClient, err := admin.NewClient(sc.Endpoint, sc.Cred, nil)
	sc.PanicOnError("failed to create admin client", err)

	wg := sync.WaitGroup{}
	wg.Add(2)

	const numToSend = 100

	go func() {
		defer wg.Done()
		queueName, sender := createDetachResources(sc, "send")

		for i := 0; i < numToSend; i++ {
			err := shared.ForceQueueDetach(sc.Context, adminClient, queueName)
			sc.PanicOnError("failed updating queue", err)

			err = sender.SendMessage(sc.Context, &azservicebus.Message{
				Body: []byte(fmt.Sprintf("test body %d", i)),
			}, nil)
			sc.PanicOnError("failed to send message", err)
		}

		checkMessages(sc, queueName, numToSend)
	}()

	go func() {
		defer wg.Done()
		queueName, sender := createDetachResources(sc, "sendBatch")

		for i := 0; i < numToSend; i++ {
			err := shared.ForceQueueDetach(sc.Context, adminClient, queueName)
			sc.PanicOnError("failed updating queue", err)

			batch, err := sender.NewMessageBatch(sc.Context, nil)
			sc.PanicOnError("failed to create message batch", err)

			err = batch.AddMessage(&azservicebus.Message{
				Body: []byte(fmt.Sprintf("batch test body %d", i)),
			}, nil)
			sc.PanicOnError("failed to add message", err)

			sc.PanicOnError("failed updating queue", err)

			err = sender.SendMessageBatch(sc.Context, batch, nil)
			sc.PanicOnError("failed to send message batch", err)
		}

		checkMessages(sc, queueName, numToSend)
	}()

	wg.Wait()
}

func createDetachResources(sc *shared.StressContext, name string) (string, *shared.TrackingSender) {
	queueName := fmt.Sprintf("detach_%s-%s", name, sc.Nano)

	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClient(sc.Endpoint, sc.Cred, nil)
	sc.PanicOnError("failed to create client", err)

	sender, err := shared.NewTrackingSender(sc.TC, client, queueName, nil)
	sc.PanicOnError("failed to create a sender", err)

	return queueName, sender
}

func checkMessages(sc *shared.StressContext, queueName string, numSent int) {
	client, err := azservicebus.NewClient(sc.Endpoint, sc.Cred, nil)
	sc.PanicOnError("failed to create client", err)

	receiver, err := shared.NewTrackingReceiverForQueue(sc.TC, client, queueName, &azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModeReceiveAndDelete,
	})
	sc.PanicOnError("failed to create receiver", err)

	defer func() { _ = receiver.Close(sc.Context) }()

	var all []*azservicebus.ReceivedMessage

	for {
		ctx, cancel := context.WithTimeout(sc.Context, 60*time.Second)
		defer cancel()

		messages, err := receiver.ReceiveMessages(ctx, numSent, nil)

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			// test is over or we received nothing (ie, we're probably done)
			break
		}

		sc.PanicOnError("failed to receive messages", err)

		if len(messages) == 0 {
			// probably done
			break
		}

		all = append(all, messages...)
	}

	if numSent != len(all) {
		sc.PanicOnError(fmt.Sprintf("Incorrect # of messages. Expected %d, got %d", numSent, len(all)), errors.New("Bad"))
	}
}

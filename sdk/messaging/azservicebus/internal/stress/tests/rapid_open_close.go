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

func RapidOpenCloseTest(remainingArgs []string) {
	sc := shared.MustCreateStressContext("RapidOpenCloseTest")
	queueName := fmt.Sprintf("rapid_open_close-%X", time.Now().UnixNano())

	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	for round := 0; round < 100; round++ {
		func() {
			log.Printf("[%d] Open/Close", round)
			client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
			sc.PanicOnError("failed to create client", err)

			defer func() {
				err = client.Close(context.Background())
				sc.PanicOnError("failed to close client", err)
			}()

			for i := 0; i < 1000; i++ {
				sender, err := client.NewSender(queueName, nil)
				sc.PanicOnError("failed to create sender", err)

				err = sender.SendMessage(context.Background(), &azservicebus.Message{
					Body: []byte("ping"),
				}, nil)
				sc.PanicOnError("failed to send message", err)

				err = sender.Close(sc.Context)
				sc.PanicOnError("failed to close client", err)

				receiver, err := client.NewReceiverForQueue(queueName, nil)
				sc.NoError(err)

				messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
				sc.NoError(err)
				sc.Equal(1, len(messages))

				err = receiver.Close(context.Background())
				sc.NoError(err)
			}
		}()
	}
}

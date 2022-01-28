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

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	sc.PanicOnError("failed to create receiver", err)

	for i := 0; i < 10000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Nanosecond)
		defer cancel()
		_, err := receiver.ReceiveMessages(ctx, 95, nil)
		sc.PanicOnError("failed to receive messages", err)
	}
}

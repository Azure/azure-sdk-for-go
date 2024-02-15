// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package benchmarks

import (
	"context"
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

// BackupSettlementLeak checks that, when we use backup settlement, that we're not
// leaking memory. This came up in a couple of issues for a customer:
// - https://github.com/Azure/azure-sdk-for-go/issues/22318
// - https://github.com/Azure/azure-sdk-for-go/issues/22157
//
// The use case for backup settlement is for when the original link we've received
// on has gone offline, so we need to settle via the management$ link instead. However,
// the underlying go-amqp link is tracking several bits of state for the message which
// will never get cleared. Since that receiver was dead it was going to get garbage
// collected anyways, so this was non-issue.
//
// This customer's use case was slightly different - they were completing on a separate
// receiver even when the original receiving link was still alive. This means the memory
// leak is just accumulating and never gets garbage collected since there's no trigger
// to know when to clear out any tracking state for the message.
func BenchmarkBackupSettlementLeakWhileOldReceiverStillAlive(b *testing.B) {
	b.StopTimer()

	sc := shared.MustCreateStressContext("BenchmarkBackupSettlementLeak", nil)
	defer sc.End()

	sent := int64(100000)

	client, queueName := mustInitBenchmarkBackupSettlementLeak(sc, b, int(sent))

	oldReceiver, err := client.NewReceiverForQueue(queueName, nil)
	sc.NoError(err)

	newReceiver, err := client.NewReceiverForQueue(queueName, nil)
	sc.NoError(err)

	b.StartTimer()

	var completed int64
	expected := maxDeliveryCount * int64(sent)

	for completed < expected {
		// receive from the old receiver and...
		receiveCtx, cancel := context.WithTimeout(context.Background(), time.Minute)

		messages, err := oldReceiver.ReceiveMessages(receiveCtx, int(math.Min(float64(expected-completed), 5000)), &azservicebus.ReceiveMessagesOptions{
			// not super scientific - mostly just want to get slightly fuller batches
			TimeAfterFirstMessage: 30 * time.Second,
		})
		cancel()
		sc.NoError(err)

		wg := sync.WaitGroup{}
		wg.Add(len(messages))

		// ...completing on another receiver
		for _, m := range messages {
			m := m

			go func() {
				defer wg.Done()

				// abandon it so we see the message a few times (until it's deadlettered after 10 tries)
				err := newReceiver.AbandonMessage(context.Background(), m, nil)
				sc.NoError(err)
				atomic.AddInt64(&completed, 1)
			}()
		}

		wg.Wait()

		b.Logf("Settled %d/%d", completed, sent)
	}

	log.Printf("Forcing garbage collection\n")
	runtime.GC()
	log.Printf("Done with collection\n")
	time.Sleep(1 * time.Minute)
}

func mustInitBenchmarkBackupSettlementLeak(sc *shared.StressContext, b *testing.B, numToSend int) (*azservicebus.Client, string) {
	queueName := fmt.Sprintf("backup-settlement-tester-%s", sc.Nano)
	shared.MustCreateAutoDeletingQueue(sc, queueName, &admin.QueueProperties{
		MaxDeliveryCount: to.Ptr[int32](maxDeliveryCount),
	})

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	sender, err := shared.NewTrackingSender(sc.TC, client, queueName, nil)
	sc.PanicOnError("create a sender", err)

	shared.MustGenerateMessages(sc, sender, numToSend, 0)

	return client, queueName
}

const maxDeliveryCount = 20

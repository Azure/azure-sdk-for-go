// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package eph

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/devigned/tab"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
)

type (
	leasedReceiver struct {
		handle    *eventhub.ListenerHandle
		processor *EventProcessorHost
		lease     LeaseMarker
		done      func()
	}
)

func newLeasedReceiver(processor *EventProcessorHost, lease LeaseMarker) *leasedReceiver {
	return &leasedReceiver{
		processor: processor,
		lease:     lease,
	}
}

func (lr *leasedReceiver) Run(ctx context.Context) error {
	span, ctx := lr.startConsumerSpanFromContext(ctx, "eph.leasedReceiver.Run")
	defer span.End()

	partitionID := lr.lease.GetPartitionID()
	epoch := lr.lease.GetEpoch()
	lr.dlog(ctx, "running...")

	renewLeaseCtx, cancelRenewLease := context.WithCancel(context.Background())
	lr.done = cancelRenewLease

	go lr.periodicallyRenewLease(renewLeaseCtx)

	opts := []eventhub.ReceiveOption{eventhub.ReceiveWithEpoch(epoch)}
	if lr.processor.consumerGroup != "" {
		opts = append(opts, eventhub.ReceiveWithConsumerGroup(lr.processor.consumerGroup))
	}

	handle, err := lr.processor.client.Receive(ctx, partitionID, lr.processor.compositeHandlers(), opts...)
	if err != nil {
		return err
	}
	lr.handle = handle
	lr.listenForClose()
	return nil
}

func (lr *leasedReceiver) Close(ctx context.Context) error {
	span, ctx := lr.startConsumerSpanFromContext(ctx, "eph.leasedReceiver.Close")
	defer span.End()

	if lr.done != nil {
		lr.done()
	}

	if lr.handle != nil {
		return lr.handle.Close(ctx)
	}

	return nil
}

func (lr *leasedReceiver) listenForClose() {
	go func() {
		<-lr.handle.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		span, ctx := lr.startConsumerSpanFromContext(ctx, "eph.leasedReceiver.listenForClose")
		defer span.End()
		err := lr.processor.scheduler.stopReceiver(ctx, lr.lease)
		if err != nil {
			tab.For(ctx).Error(err)
		}
	}()
}

func (lr *leasedReceiver) periodicallyRenewLease(ctx context.Context) {
	span, ctx := lr.startConsumerSpanFromContext(ctx, "eph.leasedReceiver.periodicallyRenewLease")
	defer span.End()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			skew := time.Duration(rand.Intn(1000)-500) * time.Millisecond
			time.Sleep(DefaultLeaseRenewalInterval + skew)
			err := lr.tryRenew(ctx)
			if err != nil {
				tab.For(ctx).Error(err)
				_ = lr.processor.scheduler.stopReceiver(ctx, lr.lease)
			}
		}
	}
}

func (lr *leasedReceiver) tryRenew(ctx context.Context) error {
	span, ctx := lr.startConsumerSpanFromContext(ctx, "eph.leasedReceiver.tryRenew")
	defer span.End()

	lease, ok, err := lr.processor.leaser.RenewLease(ctx, lr.lease.GetPartitionID())
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}
	if !ok {
		err = errors.New("can't renew lease")
		tab.For(ctx).Error(err)
		return err
	}
	lr.dlog(ctx, "lease renewed")
	lr.lease = lease
	return nil
}

func (lr *leasedReceiver) dlog(ctx context.Context, msg string) {
	name := lr.processor.name
	partitionID := lr.lease.GetPartitionID()
	epoch := lr.lease.GetEpoch()
	tab.For(ctx).Debug(fmt.Sprintf("eph %q, partition %q, epoch %d: "+msg, name, partitionID, epoch))
}

func (lr *leasedReceiver) startConsumerSpanFromContext(ctx context.Context, operationName string) (tab.Spanner, context.Context) {
	span, ctx := startConsumerSpanFromContext(ctx, operationName)
	span.AddAttributes(
		tab.StringAttribute("eph.id", lr.processor.name),
		tab.StringAttribute(partitionIDTag, lr.lease.GetPartitionID()),
		tab.Int64Attribute(epochTag, lr.lease.GetEpoch()),
	)
	return span, ctx
}

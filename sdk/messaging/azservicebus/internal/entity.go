// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"sync"
	"time"

	"github.com/devigned/tab"
)

type (
	entity struct {
		Name           string
		managementPath string
		namespace      *Namespace
		rpcClient      *rpcClient
		rpcClientMu    sync.RWMutex
	}

	sendingEntity struct {
		*entity
	}

	receivingEntity struct {
		renewMessageLockMutex sync.Mutex
		*entity
	}

	sendAndReceiveEntity struct {
		*entity
		*sendingEntity
		*receivingEntity
	}
)

func newEntity(name string, managementPath string, ns *Namespace) *entity {
	return &entity{
		Name:           name,
		managementPath: managementPath,
		namespace:      ns,
	}
}

func newReceivingEntity(e *entity) *receivingEntity {
	return &receivingEntity{
		entity: e,
	}
}

func newSendingEntity(e *entity) *sendingEntity {
	return &sendingEntity{
		entity: e,
	}
}

func newSendAndReceiveEntity(entity *entity) *sendAndReceiveEntity {
	return &sendAndReceiveEntity{
		entity:          entity,
		receivingEntity: newReceivingEntity(entity),
		sendingEntity:   newSendingEntity(entity),
	}
}

func (e *entity) GetRPCClient(ctx context.Context) (*rpcClient, error) {
	if err := e.ensureRPCClient(ctx); err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	return e.rpcClient, nil
}

// ManagementPath is the relative uri to address the entity's management functionality
func (e *entity) ManagementPath() string {
	return e.managementPath
}

func (e *entity) Namespace() *Namespace {
	return e.namespace
}

func (e *entity) getEntity() *entity {
	return e
}

// ReceiveDeferred will receive and handle a set of deferred messages
//
// When a queue or subscription client receives a message that it is willing to process, but for which processing is
// not currently possible due to special circumstances inside of the application, it has the option of "deferring"
// retrieval of the message to a later point. The message remains in the queue or subscription, but it is set aside.
//
// Deferral is a feature specifically created for workflow processing scenarios. Workflow frameworks may require certain
// operations to be processed in a particular order, and may have to postpone processing of some received messages
// until prescribed prior work that is informed by other messages has been completed.
//
// A simple illustrative example is an order processing sequence in which a payment notification from an external
// payment provider appears in a system before the matching purchase order has been propagated from the store front
// to the fulfillment system. In that case, the fulfillment system might defer processing the payment notification
// until there is an order with which to associate it. In rendezvous scenarios, where messages from different sources
// drive a workflow forward, the real-time execution order may indeed be correct, but the messages reflecting the
// outcomes may arrive out of order.
//
// Ultimately, deferral aids in reordering messages from the arrival order into an order in which they can be
// processed, while leaving those messages safely in the message store for which processing needs to be postponed.
func (re *receivingEntity) ReceiveDeferred(ctx context.Context, handler Handler, sequenceNumbers ...int64) error {
	ctx, span := re.startSpanFromContext(ctx, "sb.receivingEntity.ReceiveDeferred")
	defer span.End()

	return re.ReceiveDeferredWithMode(ctx, handler, PeekLockMode, sequenceNumbers...)
}

// ReceiveDeferredWithMode will receive and handle a set of deferred messages
//
// When a queue or subscription client receives a message that it is willing to process, but for which processing is
// not currently possible due to special circumstances inside of the application, it has the option of "deferring"
// retrieval of the message to a later point. The message remains in the queue or subscription, but it is set aside.
//
// Deferral is a feature specifically created for workflow processing scenarios. Workflow frameworks may require certain
// operations to be processed in a particular order, and may have to postpone processing of some received messages
// until prescribed prior work that is informed by other messages has been completed.
//
// A simple illustrative example is an order processing sequence in which a payment notification from an external
// payment provider appears in a system before the matching purchase order has been propagated from the store front
// to the fulfillment system. In that case, the fulfillment system might defer processing the payment notification
// until there is an order with which to associate it. In rendezvous scenarios, where messages from different sources
// drive a workflow forward, the real-time execution order may indeed be correct, but the messages reflecting the
// outcomes may arrive out of order.
//
// Ultimately, deferral aids in reordering messages from the arrival order into an order in which they can be
// processed, while leaving those messages safely in the message store for which processing needs to be postponed.
func (re *receivingEntity) ReceiveDeferredWithMode(ctx context.Context, handler Handler, mode ReceiveMode, sequenceNumbers ...int64) error {
	ctx, span := re.startSpanFromContext(ctx, "sb.receivingEntity.ReceiveDeferred")
	defer span.End()

	rpcClient, err := re.entity.GetRPCClient(ctx)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}

	messages, err := rpcClient.ReceiveDeferred(ctx, mode, sequenceNumbers...)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}

	for _, msg := range messages {
		if err := handler.Handle(ctx, msg); err != nil {
			tab.For(ctx).Error(err)
			return err
		}
	}
	return nil
}

// RenewLocks renews the locks on messages provided
func (re *receivingEntity) RenewLocks(ctx context.Context, messages ...*Message) error {
	ctx, span := re.startSpanFromContext(ctx, "sb.receivingEntity.RenewLocks")
	defer span.End()

	re.renewMessageLockMutex.Lock()
	defer re.renewMessageLockMutex.Unlock()

	client, err := re.entity.GetRPCClient(ctx)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}
	return client.RenewLocks(ctx, messages...)
}

// ScheduleAt will send a batch of messages to a Queue, schedule them to be enqueued, and return the sequence numbers
// that can be used to cancel each message.
func (se *sendingEntity) ScheduleAt(ctx context.Context, enqueueTime time.Time, messages ...*Message) ([]int64, error) {
	ctx, span := se.startSpanFromContext(ctx, "sb.sendingEntity.ScheduleAt")
	defer span.End()

	client, err := se.entity.GetRPCClient(ctx)
	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	return client.ScheduleAt(ctx, enqueueTime, messages...)
}

// CancelScheduled allows for removal of messages that have been handed to the Service Bus broker for later delivery,
// but have not yet ben enqueued.
func (se *sendingEntity) CancelScheduled(ctx context.Context, seq ...int64) error {
	ctx, span := se.startSpanFromContext(ctx, "sb.sendingEntity.CancelScheduled")
	defer span.End()

	client, err := se.entity.GetRPCClient(ctx)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}

	return client.CancelScheduled(ctx, seq...)
}

func (e *entity) ensureRPCClient(ctx context.Context) error {
	ctx, span := e.startSpanFromContext(ctx, "sb.entity.ensureRPCClient")
	defer span.End()

	e.rpcClientMu.Lock()
	defer e.rpcClientMu.Unlock()

	if e.rpcClient != nil {
		return nil
	}

	client, err := newRPCClient(ctx, e)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}

	e.rpcClient = client
	return nil
}

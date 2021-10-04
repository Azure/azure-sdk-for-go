// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/go-amqp"
)

var errReceiveAndDeleteReceiver = errors.New("messages that are received in receiveAndDelete mode are not settlable")

type settler interface {
	CompleteMessage(ctx context.Context, message *ReceivedMessage) error
	AbandonMessage(ctx context.Context, message *ReceivedMessage) error
	DeferMessage(ctx context.Context, message *ReceivedMessage) error
	DeadLetterMessage(ctx context.Context, message *ReceivedMessage, options *DeadLetterOptions) error
}

type messageSettler struct {
	links                  internal.AMQPLinks
	onlyDoBackupSettlement bool
	baseRetrier            internal.Retrier
}

func newMessageSettler(links internal.AMQPLinks, baseRetrier internal.Retrier) settler {
	return &messageSettler{
		links:       links,
		baseRetrier: baseRetrier,
	}
}

func (s *messageSettler) useManagementLink(m *ReceivedMessage, linkRevision uint64) bool {
	return s.onlyDoBackupSettlement ||
		m.deferred ||
		m.linkRevision != linkRevision
}

func (s *messageSettler) settleWithRetries(ctx context.Context, message *ReceivedMessage, settleFn func(receiver internal.AMQPReceiver, mgmt internal.MgmtClient, linkRevision uint64) error) error {
	if s == nil {
		return errReceiveAndDeleteReceiver
	}

	retrier := s.baseRetrier.Copy()
	var lastErr error

	for retrier.Try(ctx) {
		var receiver internal.AMQPReceiver
		var mgmt internal.MgmtClient
		var linkRevision uint64

		_, receiver, mgmt, linkRevision, lastErr = s.links.Get(ctx)

		if lastErr != nil {
			_ = s.links.RecoverIfNeeded(ctx, linkRevision, lastErr)
			continue
		}

		lastErr := settleFn(receiver, mgmt, linkRevision)

		if lastErr != nil {
			_ = s.links.RecoverIfNeeded(ctx, linkRevision, lastErr)
			continue
		}

		break
	}

	return lastErr
}

// CompleteMessage completes a message, deleting it from the queue or subscription.
func (s *messageSettler) CompleteMessage(ctx context.Context, message *ReceivedMessage) error {
	return s.settleWithRetries(ctx, message, func(receiver internal.AMQPReceiver, mgmt internal.MgmtClient, linkRevision uint64) error {
		if s.useManagementLink(message, linkRevision) {
			return mgmt.SendDisposition(ctx, bytesToAMQPUUID(message.LockToken), internal.Disposition{Status: internal.CompletedDisposition})
		} else {
			return receiver.AcceptMessage(ctx, message.rawAMQPMessage)
		}
	})
}

// AbandonMessage will cause a message to be returned to the queue or subscription.
// This will increment its delivery count, and potentially cause it to be dead lettered
// depending on your queue or subscription's configuration.
func (s *messageSettler) AbandonMessage(ctx context.Context, message *ReceivedMessage) error {
	return s.settleWithRetries(ctx, message, func(receiver internal.AMQPReceiver, mgmt internal.MgmtClient, linkRevision uint64) error {
		if s.useManagementLink(message, linkRevision) {
			d := internal.Disposition{
				Status: internal.AbandonedDisposition,
			}
			return mgmt.SendDisposition(ctx, bytesToAMQPUUID(message.LockToken), d)
		}

		return receiver.ModifyMessage(ctx, message.rawAMQPMessage, false, false, nil)
	})
}

// DeferMessage will cause a message to be deferred. Deferred messages
// can be received using `Receiver.ReceiveDeferredMessages`.
func (s *messageSettler) DeferMessage(ctx context.Context, message *ReceivedMessage) error {
	return s.settleWithRetries(ctx, message, func(receiver internal.AMQPReceiver, mgmt internal.MgmtClient, linkRevision uint64) error {
		if s.useManagementLink(message, linkRevision) {
			d := internal.Disposition{
				Status: internal.DeferredDisposition,
			}
			return mgmt.SendDisposition(ctx, bytesToAMQPUUID(message.LockToken), d)
		}

		return receiver.ModifyMessage(ctx, message.rawAMQPMessage, false, true, nil)
	})
}

// DeadLetterOptions describe the reason and error description for dead lettering
// a message using the `Receiver.DeadLetterMessage()`
type DeadLetterOptions struct {
	// ErrorDescription that caused the dead lettering of the message.
	ErrorDescription *string

	// Reason for dead lettering the message.
	Reason *string

	// PropertiesToModify specifies properties to modify in the message when it is dead lettered.
	PropertiesToModify map[string]interface{}
}

// DeadLetterMessage settles a message by moving it to the dead letter queue for a
// queue or subscription. To receive these messages create a receiver with `Client.NewReceiver()`
// using the `SubQueue` option.
func (s *messageSettler) DeadLetterMessage(ctx context.Context, message *ReceivedMessage, options *DeadLetterOptions) error {
	return s.settleWithRetries(ctx, message, func(receiver internal.AMQPReceiver, mgmt internal.MgmtClient, linkRevision uint64) error {
		reason := ""
		description := ""

		if options != nil {
			if options.Reason != nil {
				reason = *options.Reason
			}

			if options.ErrorDescription != nil {
				description = *options.ErrorDescription
			}
		}

		if s.useManagementLink(message, linkRevision) {
			d := internal.Disposition{
				Status:                internal.SuspendedDisposition,
				DeadLetterDescription: &description,
				DeadLetterReason:      &reason,
			}
			return mgmt.SendDisposition(ctx, bytesToAMQPUUID(message.LockToken), d)
		}

		info := map[string]interface{}{
			"DeadLetterReason":           reason,
			"DeadLetterErrorDescription": description,
		}

		if options != nil && options.PropertiesToModify != nil {
			for key, val := range options.PropertiesToModify {
				info[key] = val
			}
		}

		amqpErr := amqp.Error{
			Condition: "com.microsoft:dead-letter",
			Info:      info,
		}

		return receiver.RejectMessage(ctx, message.rawAMQPMessage, &amqpErr)
	})
}

func bytesToAMQPUUID(bytes [16]byte) *amqp.UUID {
	uuid := amqp.UUID(bytes)
	return &uuid
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/go-amqp"
	"github.com/Azure/go-autorest/autorest/to"
)

var errReceiveAndDeleteReceiver = errors.New("messages that are received in receiveAndDelete mode are not settlable")

type messageSettler struct {
	links                  internal.AMQPLinks
	onlyDoBackupSettlement bool
}

// CompleteMessage completes a message, deleting it from the queue or subscription.
func (s *messageSettler) CompleteMessage(ctx context.Context, message *ReceivedMessage) error {
	if s == nil {
		return errReceiveAndDeleteReceiver
	}

	_, _, mgmt, linkRevision, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	// complete
	if s.onlyDoBackupSettlement || message.linkRevision != linkRevision {
		return mgmt.SendDisposition(ctx, message.LockToken, internal.Disposition{Status: internal.CompletedDisposition})
	}

	return message.RawAMQPMessage.Accept(ctx)
}

// AbandonMessage will cause a message to be returned to the queue or subscription.
// This will increment its delivery count, and potentially cause it to be dead lettered
// depending on your queue or subscription's configuration.
func (s *messageSettler) AbandonMessage(ctx context.Context, message *ReceivedMessage) error {
	if s == nil {
		return errReceiveAndDeleteReceiver
	}

	_, _, mgmt, linkRevision, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	if s.onlyDoBackupSettlement || linkRevision != message.linkRevision {
		// abandon
		d := internal.Disposition{
			Status: internal.AbandonedDisposition,
		}
		return mgmt.SendDisposition(ctx, message.LockToken, d)
	}

	return message.RawAMQPMessage.Modify(ctx, false, false, nil)
}

// DeferMessage will cause a message to be deferred. Deferred messages
// can be received using `Receiver.ReceiveDeferredMessages`.
func (s *messageSettler) DeferMessage(ctx context.Context, message *ReceivedMessage) error {
	if s == nil {
		return errReceiveAndDeleteReceiver
	}

	// TODO: this looks totally silly, but the settlement methods are moving
	// to 'receiver' in the next go-amqp rev _and_ we need to figure out how to do
	// backup 'defer' settlement.
	_, _, _, _, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	return message.RawAMQPMessage.Modify(ctx, true, true, nil)
}

type DeadLetterOptions struct {
	errorDescription   *string
	reason             *string
	propertiesToModify map[string]interface{}
}

type DeadLetterOption func(options *DeadLetterOptions) error

func DeadLetterWithErrorDescription(description string) DeadLetterOption {
	return func(options *DeadLetterOptions) error {
		options.errorDescription = &description
		return nil
	}
}

func DeadLetterWithReason(reason string) DeadLetterOption {
	return func(options *DeadLetterOptions) error {
		options.reason = &reason
		return nil
	}
}

func DeadLetterWithPropertiesToModify(propertiesToModify map[string]interface{}) DeadLetterOption {
	return func(options *DeadLetterOptions) error {
		options.propertiesToModify = propertiesToModify
		return nil
	}
}

// DeadLetterMessage settles a message by moving it to the dead letter queue for a
// queue or subscription. To receive these messages create a receiver with `Client.NewReceiver()`
// using the `ReceiverWithSubQueue()` option.
func (s *messageSettler) DeadLetterMessage(ctx context.Context, message *ReceivedMessage, options ...DeadLetterOption) error {
	if s == nil {
		return errReceiveAndDeleteReceiver
	}

	deadLetterOptions := &DeadLetterOptions{}

	for _, opt := range options {
		if err := opt(deadLetterOptions); err != nil {
			return err
		}
	}

	_, _, mgmt, linkRevision, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	if s.onlyDoBackupSettlement || linkRevision != message.linkRevision {
		d := internal.Disposition{
			Status:                internal.SuspendedDisposition,
			DeadLetterDescription: to.StringPtr(err.Error()),
			DeadLetterReason:      to.StringPtr("amqp:error"),
		}
		return mgmt.SendDisposition(ctx, message.LockToken, d)
	}

	info := map[string]interface{}{
		"DeadLetterReason":           deadLetterOptions.reason,
		"DeadLetterErrorDescription": deadLetterOptions.errorDescription,
	}

	if deadLetterOptions.propertiesToModify != nil {
		for key, val := range deadLetterOptions.propertiesToModify {
			info[key] = val
		}
	}

	amqpErr := amqp.Error{
		Condition: "com.microsoft:dead-letter",
		Info:      info,
	}

	return message.RawAMQPMessage.Reject(ctx, &amqpErr)
}

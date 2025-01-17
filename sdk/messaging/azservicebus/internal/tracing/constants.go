// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"

type SpanName string

const (
	NegotiateClaimSpanName SpanName = "Namespace.NegotiateClaim"

	SendSpanName            SpanName = "Sender.SendMessage"
	SendBatchSpanName       SpanName = "Sender.SendMessageBatch"
	ScheduleSpanName        SpanName = "Sender.ScheduleMessages"
	CancelScheduledSpanName SpanName = "Sender.CancelScheduledMessages"

	ReceiveSpanName          SpanName = "Receiver.ReceiveMessages"
	PeekSpanName             SpanName = "Receiver.PeekMessages"
	ReceiveDeferredSpanName  SpanName = "Receiver.ReceiveDeferredMessages"
	RenewMessageLockSpanName SpanName = "Receiver.RenewMessageLock"

	CompleteSpanName   SpanName = "Receiver.CompleteMessage"
	AbandonSpanName    SpanName = "Receiver.AbandonMessage"
	DeferSpanName      SpanName = "Receiver.DeferMessage"
	DeadLetterSpanName SpanName = "Receiver.DeadLetterMessage"

	AcceptSessionSpanName    SpanName = "SessionReceiver.AcceptSession"
	GetSessionStateSpanName  SpanName = "SessionReceiver.GetSessionState"
	SetSessionStateSpanName  SpanName = "SessionReceiver.SetSessionState"
	RenewSessionLockSpanName SpanName = "SessionReceiver.RenewSessionLock"
)

type SpanKind = tracing.SpanKind

const (
	SpanKindInternal tracing.SpanKind = tracing.SpanKindInternal
	SpanKindClient   tracing.SpanKind = tracing.SpanKindClient
	SpanKindProducer tracing.SpanKind = tracing.SpanKindProducer
	SpanKindConsumer tracing.SpanKind = tracing.SpanKindConsumer
)

// OTel-specific messaging attributes
const (
	ServerAddress     = "server.address"
	MessagingSystem   = "messaging.system"
	OperationName     = "messaging.operation.name"
	BatchMessageCount = "messaging.batch.message_count"
	DestinationName   = "messaging.destination.name"
	SubscriptionName  = "messaging.destination.subscription.name"
	OperationType     = "messaging.operation.type"
	DispositionStatus = "messaging.servicebus.disposition_status"
	DeliveryCount     = "messaging.servicebus.message.delivery_count"
	ConversationID    = "messaging.message.conversation_id"
	MessageID         = "messaging.message.id"
	EnqueuedTime      = "messaging.servicebus.message.enqueued_time"
)

type MessagingOperationType string

const (
	SendOperationType    MessagingOperationType = "send"
	ReceiveOperationType MessagingOperationType = "receive"
	SettleOperationType  MessagingOperationType = "settle"
	SessionOperationType MessagingOperationType = "session"
)

type MessagingOperationName string

const (
	SendOperationName            MessagingOperationName = "send"
	ScheduleOperationName        MessagingOperationName = "schedule"
	CancelScheduledOperationName MessagingOperationName = "cancel_scheduled"

	ReceiveOperationName          MessagingOperationName = "receive"
	PeekOperationName             MessagingOperationName = "peek"
	ReceiveDeferredOperationName  MessagingOperationName = "receive_deferred"
	RenewMessageLockOperationName MessagingOperationName = "renew_message_lock"

	AbandonOperationName    MessagingOperationName = "abandon"
	CompleteOperationName   MessagingOperationName = "complete"
	DeferOperationName      MessagingOperationName = "defer"
	DeadLetterOperationName MessagingOperationName = "dead_letter"

	AcceptSessionOperationName    MessagingOperationName = "accept_session"
	GetSessionStateOperationName  MessagingOperationName = "get_session_state"
	SetSessionStateOperationName  MessagingOperationName = "set_session_state"
	RenewSessionLockOperationName MessagingOperationName = "renew_session_lock"
)

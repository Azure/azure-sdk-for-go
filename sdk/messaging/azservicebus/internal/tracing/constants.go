// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"

type SpanKind = tracing.SpanKind

const (
	SpanKindInternal = tracing.SpanKindInternal
	SpanKindClient   = tracing.SpanKindClient
	SpanKindProducer = tracing.SpanKindProducer
	SpanKindConsumer = tracing.SpanKindConsumer
)

type SpanContext = tracing.SpanContext

const (
	SpanStatusUnset = tracing.SpanStatusUnset
	SpanStatusError = tracing.SpanStatusError
	SpanStatusOK    = tracing.SpanStatusOK
)

const (
	AttrServerAddress     = "server.address"
	AttrMessagingSystem   = "messaging.system"
	AttrOperationName     = "messaging.operation.name"
	AttrBatchMessageCount = "messaging.batch.message_count"
	AttrDestinationName   = "messaging.destination.name"
	AttrSubscriptionName  = "messaging.destination.subscription.name"
	AttrOperationType     = "messaging.operation.type"
	AttrDispositionStatus = "messaging.servicebus.disposition_status"
	AttrDeliveryCount     = "messaging.servicebus.message.delivery_count"
	AttrConversationID    = "messaging.message.conversation_id"
	AttrMessageID         = "messaging.message.id"
	AttrEnqueuedTime      = "messaging.servicebus.message.enqueued_time"
	AttrErrorType         = "error.type"
)

type MessagingOperationType string

const (
	CreateOperationType  MessagingOperationType = "create"
	SendOperationType    MessagingOperationType = "send"
	ReceiveOperationType MessagingOperationType = "receive"
	SettleOperationType  MessagingOperationType = "settle"
)

type MessagingOperationName string

const (
	CreateOperationName          MessagingOperationName = "create"
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

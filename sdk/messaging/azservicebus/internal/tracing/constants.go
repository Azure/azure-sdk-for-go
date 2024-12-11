// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

// OTel-specific messaging attributes
const (
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

	ServerAddress = "server.address"
	ServerPort    = "server.port"
)

type MessagingOperationType string

const (
	SendOperationType    MessagingOperationType = "send"
	ReceiveOperationType MessagingOperationType = "receive"
	SettleOperationType  MessagingOperationType = "settle"
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
	DeadLetterOperationName MessagingOperationName = "deadletter"
	DeleteOperationName     MessagingOperationName = "delete"
)

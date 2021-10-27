// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"time"
)

// QueueProperties represents the static properties of the queue.
type QueueProperties struct {
	// Name of the queue relative to the namespace base address.
	Name string

	// LockDuration - The duration a message is locked when using the PeekLock receive mode.
	// Default is 1 minute.
	LockDuration *time.Duration

	// MaxSizeInMegabytes - The maximum size of the queue in megabytes, which is the size of memory
	// allocated for the queue.
	// Default is 1024.
	MaxSizeInMegabytes *int32

	// RequiresDuplicateDetection - A value indicating if this queue requires duplicate detection.
	RequiresDuplicateDetection *bool

	// RequiresSession indicates whether the queue supports the concept of sessions.
	// Sessionful-messages follow FIFO ordering.
	// Default is false.
	RequiresSession *bool

	// DefaultMessageTimeToLive is the duration after which the message expires, starting from when
	// the message is sent to Service Bus. This is the default value used when TimeToLive is not
	// set on a message itself.
	DefaultMessageTimeToLive *time.Duration

	// DeadLetteringOnMessageExpiration indicates whether this queue has dead letter
	// support when a message expires.
	DeadLetteringOnMessageExpiration *bool

	// DuplicateDetectionHistoryTimeWindow is the duration of duplicate detection history.
	// Default value is 10 minutes.
	DuplicateDetectionHistoryTimeWindow *time.Duration

	// MaxDeliveryCount is the maximum amount of times a message can be delivered before it is automatically
	// sent to the dead letter queue.
	// Default value is 10.
	MaxDeliveryCount *int32

	// EnableBatchedOperations indicates whether server-side batched operations are enabled.
	EnableBatchedOperations *bool

	// The current status of the queue.
	Status *EntityStatus

	// AutoDeleteOnIdle is the idle interval after which the queue is automatically deleted.
	AutoDeleteOnIdle *time.Duration

	// Indicates whether the queue is to be partitioned across multiple message brokers.
	EnablePartitioning *bool

	// ForwardTo is the name of the recipient entity to which all the messages sent to the queue
	// are forwarded to.
	ForwardTo *string

	// ForwardDeadLetteredMessagesTo - absolute URI of the entity to forward dead letter messages
	ForwardDeadLetteredMessagesTo *string
}

// QueueRuntimeProperties represent dynamic properties of a queue, such as the ActiveMessageCount.
type QueueRuntimeProperties struct {
	// Name is the name of the queue.
	Name string

	// SizeInBytes - The size of the queue, in bytes.
	SizeInBytes int64

	// CreatedAt is when the entity was created.
	CreatedAt time.Time

	// UpdatedAt is when the entity was last updated.
	UpdatedAt time.Time

	// AccessedAt is when the entity was last updated.
	AccessedAt time.Time

	// TotalMessageCount is the number of messages in the queue.
	TotalMessageCount int64

	// ActiveMessageCount is the number of active messages in the entity.
	ActiveMessageCount int32

	// DeadLetterMessageCount is the number of dead-lettered messages in the entity.
	DeadLetterMessageCount int32

	// ScheduledMessageCount is the number of messages that are scheduled to be enqueued.
	ScheduledMessageCount int32

	// TransferDeadLetterMessageCount is the number of messages transfer-messages which are dead-lettered
	// into transfer-dead-letter subqueue.
	TransferDeadLetterMessageCount int32

	// TransferMessageCount is the number of messages which are yet to be transferred/forwarded to destination entity.
	TransferMessageCount int32
}

// EntityStatus represents the current status of the entity.
type EntityStatus string

const (
	// EntityStatusActive indicates an entity can be used for sending and receiving.
	EntityStatusActive EntityStatus = "Active"
	// EntityStatusDisabled indicates an entity cannot be used for sending or receiving.
	EntityStatusDisabled EntityStatus = "Disabled"
	// EntityStatusSendDisabled indicates that an entity cannot be used for sending.
	EntityStatusSendDisabled EntityStatus = "SendDisabled"
	// EntityStatusReceiveDisabled indicates that an entity cannot be used for receiving.
	EntityStatusReceiveDisabled EntityStatus = "ReceiveDisabled"
)

// AccessRights represents the rights for an authorization rule.
type AccessRights int

const (
	// AccessRightsManage allows management of entities.
	AccessRightsManage AccessRights = 0
	// AccessRightsSend allows sending to entities.
	AccessRightsSend AccessRights = 1
	// AccessRightsListen allows listening to entities.
	AccessRightsListen AccessRights = 2
)

// EntityAvailabilityStatus is the availability status of the entity.
type EntityAvailabilityStatus string

const (
	// EntityAvailabilityStatusAvailable indicates the entity is available.
	EntityAvailabilityStatusAvailable EntityAvailabilityStatus = "Available"

	EntityAvailabilityStatusLimited EntityAvailabilityStatus = "Limited"

	// EntityAvailabilityStatusRenaming indicates the entity is being renamed.
	EntityAvailabilityStatusRenaming EntityAvailabilityStatus = "Renaming"

	// EntityAvailabilityStatusRestoring indicates the entity is being restored.
	EntityAvailabilityStatusRestoring EntityAvailabilityStatus = "Restoring"

	EntityAvailabilityStatusUnknown EntityAvailabilityStatus = "Unknown"
)

// AuthorizationRule encompasses access rights, metadata and a key for authentication.
type AuthorizationRule struct {
	ClaimType    string
	Rights       []AccessRights
	KeyName      string
	CreatedTime  time.Time
	ModifiedTime time.Time
}

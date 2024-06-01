// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package aznamespaces

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"

// AcknowledgeResult - The result of the Acknowledge operation.
type AcknowledgeResult struct {
	// REQUIRED; Array of FailedLockToken for failed cloud events. Each FailedLockToken includes the lock token along with the
	// related error information (namely, the error code and description).
	FailedLockTokens []FailedLockToken

	// REQUIRED; Array of lock tokens for the successfully acknowledged cloud events.
	SucceededLockTokens []string
}

// BrokerProperties - Properties of the Event Broker operation.
type BrokerProperties struct {
	// REQUIRED; The attempt count for delivering the event.
	DeliveryCount *int32

	// REQUIRED; The token of the lock on the event.
	LockToken *string
}

// Error - The error object.
type Error struct {
	// REQUIRED; One of a server-defined set of error codes.
	Code *string

	// REQUIRED; A human-readable representation of the error.
	Message *string
}

// FailedLockToken - Failed LockToken information.
type FailedLockToken struct {
	// REQUIRED; Error information of the failed operation result for the lock token in the request.
	Error *Error

	// REQUIRED; The lock token of an entry in the request.
	LockToken *string
}

// PublishResult - The result of the Publish operation.
type PublishResult struct {
}

// ReceiveDetails - Receive operation details per Cloud Event.
type ReceiveDetails struct {
	// REQUIRED; The Event Broker details.
	BrokerProperties *BrokerProperties

	// REQUIRED; Cloud Event details.
	Event messaging.CloudEvent
}

// ReceiveResult - Details of the Receive operation response.
type ReceiveResult struct {
	// REQUIRED; Array of receive responses, one per cloud event.
	Details []ReceiveDetails
}

// RejectResult - The result of the Reject operation.
type RejectResult struct {
	// REQUIRED; Array of FailedLockToken for failed cloud events. Each FailedLockToken includes the lock token along with the
	// related error information (namely, the error code and description).
	FailedLockTokens []FailedLockToken

	// REQUIRED; Array of lock tokens for the successfully rejected cloud events.
	SucceededLockTokens []string
}

// ReleaseResult - The result of the Release operation.
type ReleaseResult struct {
	// REQUIRED; Array of FailedLockToken for failed cloud events. Each FailedLockToken includes the lock token along with the
	// related error information (namely, the error code and description).
	FailedLockTokens []FailedLockToken

	// REQUIRED; Array of lock tokens for the successfully released cloud events.
	SucceededLockTokens []string
}

// RenewLocksResult - The result of the RenewLock operation.
type RenewLocksResult struct {
	// REQUIRED; Array of FailedLockToken for failed cloud events. Each FailedLockToken includes the lock token along with the
	// related error information (namely, the error code and description).
	FailedLockTokens []FailedLockToken

	// REQUIRED; Array of lock tokens for the successfully renewed locks.
	SucceededLockTokens []string
}

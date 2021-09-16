// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

const vendorPrefix = "com.microsoft:"

// Operations
const (
	lockRenewalOperationName   = vendorPrefix + "renew-lock"
	peekMessageOperationID     = vendorPrefix + "peek-message"
	scheduleMessageOperationID = vendorPrefix + "schedule-message"
	cancelScheduledOperationID = vendorPrefix + "cancel-scheduled-message"
)

// Field Descriptions
const (
	operationFieldName     = "operation"
	lockTokensFieldName    = "lock-tokens"
	serverTimeoutFieldName = vendorPrefix + "server-timeout"
	associatedLinkName     = "associated-link-name"
)

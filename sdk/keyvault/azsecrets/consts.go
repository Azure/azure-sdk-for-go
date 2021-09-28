//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets/internal"

type DeletionRecoveryLevel string

const (
	// CustomizedRecoverable - Denotes a vault state in which deletion is recoverable without the possibility for immediate and permanent
	// deletion (i.e. purge when 7<= SoftDeleteRetentionInDays < 90).This level guarantees the recoverability of the deleted entity during the retention interval
	// and while the subscription is still available.
	CustomizedRecoverable DeletionRecoveryLevel = "CustomizedRecoverable"

	// CustomizedRecoverableProtectedSubscription - Denotes a vault and subscription state in which deletion is recoverable, immediate
	// and permanent deletion (i.e. purge) is not permitted, and in which the subscription itself cannot be permanently canceled when 7<= SoftDeleteRetentionInDays
	// < 90. This level guarantees the recoverability of the deleted entity during the retention interval, and also reflects the fact that the subscription
	// itself cannot be cancelled.
	CustomizedRecoverableProtectedSubscription DeletionRecoveryLevel = "CustomizedRecoverable+ProtectedSubscription"

	// CustomizedRecoverablePurgeable - Denotes a vault state in which deletion is recoverable, and which also permits immediate and permanent
	// deletion (i.e. purge when 7<= SoftDeleteRetentionInDays < 90). This level guarantees the recoverability of the deleted entity during the retention interval,
	// unless a Purge operation is requested, or the subscription is cancelled.
	CustomizedRecoverablePurgeable DeletionRecoveryLevel = "CustomizedRecoverable+Purgeable"

	// Purgeable - Denotes a vault state in which deletion is an irreversible operation, without the possibility for recovery. This level
	// corresponds to no protection being available against a Delete operation; the data is irretrievably lost upon accepting a Delete operation at the entity
	// level or higher (vault, resource group, subscription etc.)
	Purgeable DeletionRecoveryLevel = "Purgeable"

	// Recoverable - Denotes a vault state in which deletion is recoverable without the possibility for immediate and permanent deletion
	// (i.e. purge). This level guarantees the recoverability of the deleted entity during the retention interval(90 days) and while the subscription is still
	// available. System wil permanently delete it after 90 days, if not recovered
	Recoverable DeletionRecoveryLevel = "Recoverable"

	// RecoverableProtectedSubscription - Denotes a vault and subscription state in which deletion is recoverable within retention interval
	// (90 days), immediate and permanent deletion (i.e. purge) is not permitted, and in which the subscription itself cannot be permanently canceled. System
	// wil permanently delete it after 90 days, if not recovered
	RecoverableProtectedSubscription DeletionRecoveryLevel = "Recoverable+ProtectedSubscription"

	// RecoverablePurgeable - Denotes a vault state in which deletion is recoverable, and which also permits immediate and permanent deletion
	// (i.e. purge). This level guarantees the recoverability of the deleted entity during the retention interval (90 days), unless a Purge operation is requested,
	// or the subscription is cancelled. System wil permanently delete it after 90 days, if not recovered
	RecoverablePurgeable DeletionRecoveryLevel = "Recoverable+Purgeable"
)

func deletionRecoveryLevelFromGenerated(i internal.DeletionRecoveryLevel) DeletionRecoveryLevel {
	if i == internal.DeletionRecoveryLevelCustomizedRecoverable {
		return CustomizedRecoverable
	} else if i == internal.DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription {
		return CustomizedRecoverableProtectedSubscription
	} else if i == internal.DeletionRecoveryLevelCustomizedRecoverablePurgeable {
		return CustomizedRecoverablePurgeable
	} else if i == internal.DeletionRecoveryLevelPurgeable {
		return Purgeable
	} else if i == internal.DeletionRecoveryLevelRecoverable {
		return Recoverable
	} else if i == internal.DeletionRecoveryLevelRecoverableProtectedSubscription {
		return RecoverableProtectedSubscription
	} else {
		return RecoverablePurgeable
	}
}

func (d DeletionRecoveryLevel) toGenerated() internal.DeletionRecoveryLevel {
	if d == CustomizedRecoverable {
		return internal.DeletionRecoveryLevelCustomizedRecoverable
	} else if d == CustomizedRecoverableProtectedSubscription {
		return internal.DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription
	} else if d == CustomizedRecoverablePurgeable {
		return internal.DeletionRecoveryLevelCustomizedRecoverablePurgeable
	} else if d == Purgeable {
		return internal.DeletionRecoveryLevelPurgeable
	} else if d == Recoverable {
		return internal.DeletionRecoveryLevelRecoverable
	} else if d == RecoverableProtectedSubscription {
		return internal.DeletionRecoveryLevelRecoverableProtectedSubscription
	} else {
		return internal.DeletionRecoveryLevelRecoverablePurgeable
	}
}

// ToPtr returns a *DeletionRecoveryLevel pointing to the current value.
func (c DeletionRecoveryLevel) ToPtr() *DeletionRecoveryLevel {
	return &c
}

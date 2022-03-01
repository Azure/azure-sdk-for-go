//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets/internal"

type DeletionRecoveryLevel string

const (
	// DeletionRecoveryLevelCustomizedRecoverable - Denotes a vault state in which deletion is recoverable without the possibility for immediate and permanent
	// deletion (i.e. purge when 7<= SoftDeleteRetentionInDays < 90).This level guarantees the recoverability of the deleted entity during the retention interval
	// and while the subscription is still available.
	DeletionRecoveryLevelCustomizedRecoverable DeletionRecoveryLevel = "CustomizedRecoverable"

	// DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription - Denotes a vault and subscription state in which deletion is recoverable, immediate
	// and permanent deletion (i.e. purge) is not permitted, and in which the subscription itself cannot be permanently canceled when 7<= SoftDeleteRetentionInDays
	// < 90. This level guarantees the recoverability of the deleted entity during the retention interval, and also reflects the fact that the subscription
	// itself cannot be cancelled.
	DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription DeletionRecoveryLevel = "CustomizedRecoverable+ProtectedSubscription"

	// DeletionRecoveryLevelCustomizedRecoverablePurgeable - Denotes a vault state in which deletion is recoverable, and which also permits immediate and permanent
	// deletion (i.e. purge when 7<= SoftDeleteRetentionInDays < 90). This level guarantees the recoverability of the deleted entity during the retention interval,
	// unless a Purge operation is requested, or the subscription is cancelled.
	DeletionRecoveryLevelCustomizedRecoverablePurgeable DeletionRecoveryLevel = "CustomizedRecoverable+Purgeable"

	// DeletionRecoveryLevelPurgeable - Denotes a vault state in which deletion is an irreversible operation, without the possibility for recovery. This level
	// corresponds to no protection being available against a Delete operation; the data is irretrievably lost upon accepting a Delete operation at the entity
	// level or higher (vault, resource group, subscription etc.)
	DeletionRecoveryLevelPurgeable DeletionRecoveryLevel = "Purgeable"

	// DeletionRecoveryLevelRecoverable - Denotes a vault state in which deletion is recoverable without the possibility for immediate and permanent deletion
	// (i.e. purge). This level guarantees the recoverability of the deleted entity during the retention interval(90 days) and while the subscription is still
	// available. System wil permanently delete it after 90 days, if not recovered
	DeletionRecoveryLevelRecoverable DeletionRecoveryLevel = "Recoverable"

	// DeletionRecoveryLevelRecoverableProtectedSubscription - Denotes a vault and subscription state in which deletion is recoverable within retention interval
	// (90 days), immediate and permanent deletion (i.e. purge) is not permitted, and in which the subscription itself cannot be permanently canceled. System
	// wil permanently delete it after 90 days, if not recovered
	DeletionRecoveryLevelRecoverableProtectedSubscription DeletionRecoveryLevel = "Recoverable+ProtectedSubscription"

	// DeletionRecoveryLevelRecoverablePurgeable - Denotes a vault state in which deletion is recoverable, and which also permits immediate and permanent deletion
	// (i.e. purge). This level guarantees the recoverability of the deleted entity during the retention interval (90 days), unless a Purge operation is requested,
	// or the subscription is cancelled. System wil permanently delete it after 90 days, if not recovered
	DeletionRecoveryLevelRecoverablePurgeable DeletionRecoveryLevel = "Recoverable+Purgeable"
)

func deletionRecoveryLevelFromGenerated(i internal.DeletionRecoveryLevel) DeletionRecoveryLevel {
	if i == internal.DeletionRecoveryLevelCustomizedRecoverable {
		return DeletionRecoveryLevelCustomizedRecoverable
	} else if i == internal.DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription {
		return DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription
	} else if i == internal.DeletionRecoveryLevelCustomizedRecoverablePurgeable {
		return DeletionRecoveryLevelCustomizedRecoverablePurgeable
	} else if i == internal.DeletionRecoveryLevelPurgeable {
		return DeletionRecoveryLevelPurgeable
	} else if i == internal.DeletionRecoveryLevelRecoverable {
		return DeletionRecoveryLevelRecoverable
	} else if i == internal.DeletionRecoveryLevelRecoverableProtectedSubscription {
		return DeletionRecoveryLevelRecoverableProtectedSubscription
	} else {
		return DeletionRecoveryLevelRecoverablePurgeable
	}
}

func (d *DeletionRecoveryLevel) toGenerated() *internal.DeletionRecoveryLevel {
	if d == nil {
		return nil
	}
	if *d == DeletionRecoveryLevelCustomizedRecoverable {
		i := internal.DeletionRecoveryLevelCustomizedRecoverable
		return i.ToPtr()
	} else if *d == DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription {
		i := internal.DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription
		return i.ToPtr()
	} else if *d == DeletionRecoveryLevelCustomizedRecoverablePurgeable {
		i := internal.DeletionRecoveryLevelCustomizedRecoverablePurgeable
		return i.ToPtr()
	} else if *d == DeletionRecoveryLevelPurgeable {
		i := internal.DeletionRecoveryLevelPurgeable
		return i.ToPtr()
	} else if *d == DeletionRecoveryLevelRecoverable {
		i := internal.DeletionRecoveryLevelRecoverable
		return i.ToPtr()
	} else if *d == DeletionRecoveryLevelRecoverableProtectedSubscription {
		i := internal.DeletionRecoveryLevelRecoverableProtectedSubscription
		return i.ToPtr()
	} else {
		i := internal.DeletionRecoveryLevelRecoverablePurgeable
		return i.ToPtr()
	}
}

// ToPtr returns a *DeletionRecoveryLevel pointing to the current value.
func (c DeletionRecoveryLevel) ToPtr() *DeletionRecoveryLevel {
	return &c
}

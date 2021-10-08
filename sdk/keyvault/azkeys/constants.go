//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

// DeletionRecoveryLevel - Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains 'Purgeable', the
// certificate can be permanently deleted by a privileged user; otherwise,
// only the system can purge the certificate, at the end of the retention interval.
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

func (d DeletionRecoveryLevel) ToPtr() *DeletionRecoveryLevel {
	return &d
}

// JSONWebKeyCurveName - Elliptic curve name. For valid values, see JsonWebKeyCurveName.
type JSONWebKeyCurveName string

const (
	// P256 - The NIST P-256 elliptic curve, AKA SECG curve SECP256R1.
	P256 JSONWebKeyCurveName = "P-256"

	// P256K - The SECG SECP256K1 elliptic curve.
	P256K JSONWebKeyCurveName = "P-256K"

	// P384 - The NIST P-384 elliptic curve, AKA SECG curve SECP384R1.
	P384 JSONWebKeyCurveName = "P-384"

	// P521 - The NIST P-521 elliptic curve, AKA SECG curve SECP521R1.
	eP521 JSONWebKeyCurveName = "P-521"
)

// ToPtr returns a *JSONWebKeyCurveName pointing to the current value.
func (c JSONWebKeyCurveName) ToPtr() *JSONWebKeyCurveName {
	return &c
}

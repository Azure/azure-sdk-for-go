//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

// ActionType - The type of the action.
type ActionType string

const (
	ActionTypeEmailContacts ActionType = "EmailContacts"
	ActionTypeAutoRenew     ActionType = "AutoRenew"
)

// ToPtr returns a *ActionType pointing to the current value.
func (c ActionType) ToPtr() *ActionType {
	return &c
}

// DeletionRecoveryLevel - Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains 'Purgeable', the
// certificate can be permanently deleted by a privileged user; otherwise,
// only the system can purge the certificate, at the end of the retention interval.
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

// ToPtr returns a *DeletionRecoveryLevel pointing to the current value.
func (c DeletionRecoveryLevel) ToPtr() *DeletionRecoveryLevel {
	return &c
}

// JSONWebKeyCurveName - Elliptic curve name. For valid values, see JsonWebKeyCurveName.
type JSONWebKeyCurveName string

const (
	JSONWebKeyCurveNameP256  JSONWebKeyCurveName = "P-256"
	JSONWebKeyCurveNameP256K JSONWebKeyCurveName = "P-256K"
	JSONWebKeyCurveNameP384  JSONWebKeyCurveName = "P-384"
	JSONWebKeyCurveNameP521  JSONWebKeyCurveName = "P-521"
)

// ToPtr returns a *JSONWebKeyCurveName pointing to the current value.
func (c JSONWebKeyCurveName) ToPtr() *JSONWebKeyCurveName {
	return &c
}

// JSONWebKeyType - The type of key pair to be used for the certificate.
type JSONWebKeyType string

const (
	JSONWebKeyTypeEC     JSONWebKeyType = "EC"
	JSONWebKeyTypeECHSM  JSONWebKeyType = "EC-HSM"
	JSONWebKeyTypeOct    JSONWebKeyType = "oct"
	JSONWebKeyTypeOctHSM JSONWebKeyType = "oct-HSM"
	JSONWebKeyTypeRSA    JSONWebKeyType = "RSA"
	JSONWebKeyTypeRSAHSM JSONWebKeyType = "RSA-HSM"
)

// ToPtr returns a *JSONWebKeyType pointing to the current value.
func (c JSONWebKeyType) ToPtr() *JSONWebKeyType {
	return &c
}

type KeyUsageType string

const (
	KeyUsageTypeCRLSign          KeyUsageType = "cRLSign"
	KeyUsageTypeDataEncipherment KeyUsageType = "dataEncipherment"
	KeyUsageTypeDecipherOnly     KeyUsageType = "decipherOnly"
	KeyUsageTypeDigitalSignature KeyUsageType = "digitalSignature"
	KeyUsageTypeEncipherOnly     KeyUsageType = "encipherOnly"
	KeyUsageTypeKeyAgreement     KeyUsageType = "keyAgreement"
	KeyUsageTypeKeyCertSign      KeyUsageType = "keyCertSign"
	KeyUsageTypeKeyEncipherment  KeyUsageType = "keyEncipherment"
	KeyUsageTypeNonRepudiation   KeyUsageType = "nonRepudiation"
)

// ToPtr returns a *KeyUsageType pointing to the current value.
func (c KeyUsageType) ToPtr() *KeyUsageType {
	return &c
}

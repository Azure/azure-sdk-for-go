//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

// CertificatePolicyAction - The type of the action.
type CertificatePolicyAction string

const (
	CertificatePolicyActionEmailContacts CertificatePolicyAction = "EmailContacts"
	CertificatePolicyActionAutoRenew     CertificatePolicyAction = "AutoRenew"
)

// ToPtr returns a *ActionType pointing to the current value.
func (c CertificatePolicyAction) ToPtr() *CertificatePolicyAction {
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

// CertificateKeyCurveName - Elliptic curve name. For valid values, see CertificateKeyCurveName.
type CertificateKeyCurveName string

const (
	CertificateKeyCurveNameP256  CertificateKeyCurveName = "P-256"
	CertificateKeyCurveNameP256K CertificateKeyCurveName = "P-256K"
	CertificateKeyCurveNameP384  CertificateKeyCurveName = "P-384"
	CertificateKeyCurveNameP521  CertificateKeyCurveName = "P-521"
)

// ToPtr returns a *CertificateKeyCurveName pointing to the current value.
func (c CertificateKeyCurveName) ToPtr() *CertificateKeyCurveName {
	return &c
}

// CertificateKeyType - The type of key pair to be used for the certificate.
type CertificateKeyType string

const (
	CertificateKeyTypeEC     CertificateKeyType = "EC"
	CertificateKeyTypeECHSM  CertificateKeyType = "EC-HSM"
	CertificateKeyTypeOct    CertificateKeyType = "oct"
	CertificateKeyTypeOctHSM CertificateKeyType = "oct-HSM"
	CertificateKeyTypeRSA    CertificateKeyType = "RSA"
	CertificateKeyTypeRSAHSM CertificateKeyType = "RSA-HSM"
)

// ToPtr returns a *CertificateKeyType pointing to the current value.
func (c CertificateKeyType) ToPtr() *CertificateKeyType {
	return &c
}

type CerificateKeyUsage string

const (
	CertificateKeyUsageCRLSign          CerificateKeyUsage = "cRLSign"
	CertificateKeyUsageDataEncipherment CerificateKeyUsage = "dataEncipherment"
	CertificateKeyUsageDecipherOnly     CerificateKeyUsage = "decipherOnly"
	CertificateKeyUsageDigitalSignature CerificateKeyUsage = "digitalSignature"
	CertificateKeyUsageEncipherOnly     CerificateKeyUsage = "encipherOnly"
	CertificateKeyUsageKeyAgreement     CerificateKeyUsage = "keyAgreement"
	CertificateKeyUsageKeyCertSign      CerificateKeyUsage = "keyCertSign"
	CertificateKeyUsageKeyEncipherment  CerificateKeyUsage = "keyEncipherment"
	CertificateKeyUsageNonRepudiation   CerificateKeyUsage = "nonRepudiation"
)

// ToPtr returns a *CertificateKeyUsage pointing to the current value.
func (c CerificateKeyUsage) ToPtr() *CerificateKeyUsage {
	return &c
}

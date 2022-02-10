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

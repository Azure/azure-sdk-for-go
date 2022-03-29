//go:build go1.18
// +build go1.18

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

// PossibleCertificatePolicyActionValues returns a slice of all possible CertificatePolicyAction values.
func PossibleCertificatePolicyActionValues() []CertificatePolicyAction {
	return []CertificatePolicyAction{
		CertificatePolicyActionEmailContacts,
		CertificatePolicyActionAutoRenew,
	}
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

// PossibleCertificateKeyCurveNameValues returns a slice of all possible CertificateKeyCurveName values.
func PossibleCertificateKeyCurveNameValues() []CertificateKeyCurveName {
	return []CertificateKeyCurveName{
		CertificateKeyCurveNameP256,
		CertificateKeyCurveNameP256K,
		CertificateKeyCurveNameP384,
		CertificateKeyCurveNameP521,
	}
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

// PossibleCertificateKeyTypeValues returns a slice of all possible CertificateKeyType values.
func PossibleCertificateKeyTypeValues() []CertificateKeyType {
	return []CertificateKeyType{
		CertificateKeyTypeEC,
		CertificateKeyTypeECHSM,
		CertificateKeyTypeOct,
		CertificateKeyTypeOctHSM,
		CertificateKeyTypeRSA,
		CertificateKeyTypeRSAHSM,
	}
}

// CertificateKeyUsage is the key usage for a certificate
type CertificateKeyUsage string

const (
	CertificateKeyUsageCRLSign          CertificateKeyUsage = "cRLSign"
	CertificateKeyUsageDataEncipherment CertificateKeyUsage = "dataEncipherment"
	CertificateKeyUsageDecipherOnly     CertificateKeyUsage = "decipherOnly"
	CertificateKeyUsageDigitalSignature CertificateKeyUsage = "digitalSignature"
	CertificateKeyUsageEncipherOnly     CertificateKeyUsage = "encipherOnly"
	CertificateKeyUsageKeyAgreement     CertificateKeyUsage = "keyAgreement"
	CertificateKeyUsageKeyCertSign      CertificateKeyUsage = "keyCertSign"
	CertificateKeyUsageKeyEncipherment  CertificateKeyUsage = "keyEncipherment"
	CertificateKeyUsageNonRepudiation   CertificateKeyUsage = "nonRepudiation"
)

// ToPtr returns a *CertificateKeyUsage pointing to the current value.
func (c CertificateKeyUsage) ToPtr() *CertificateKeyUsage {
	return &c
}

// PossibleCertificateKeyUsageValues returns a slice of all possible CertificateKeyUsage values.
func PossibleCertificateKeyUsageValues() []CertificateKeyUsage {
	return []CertificateKeyUsage{
		CertificateKeyUsageCRLSign,
		CertificateKeyUsageDataEncipherment,
		CertificateKeyUsageDecipherOnly,
		CertificateKeyUsageDigitalSignature,
		CertificateKeyUsageEncipherOnly,
		CertificateKeyUsageKeyAgreement,
		CertificateKeyUsageKeyCertSign,
		CertificateKeyUsageKeyEncipherment,
		CertificateKeyUsageNonRepudiation,
	}
}

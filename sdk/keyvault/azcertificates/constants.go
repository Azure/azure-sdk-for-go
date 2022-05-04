//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

// PolicyAction - The type of the action.
type PolicyAction string

const (
	PolicyActionEmailContacts PolicyAction = "EmailContacts"
	PolicyActionAutoRenew     PolicyAction = "AutoRenew"
)

// PossiblePolicyActionValues returns a slice of all possible CertificatePolicyAction values.
func PossiblePolicyActionValues() []PolicyAction {
	return []PolicyAction{
		PolicyActionEmailContacts,
		PolicyActionAutoRenew,
	}
}

// KeyCurveName - Elliptic curve name. For valid values, see KeyCurveName.
type KeyCurveName string

const (
	KeyCurveNameP256  KeyCurveName = "P-256"
	KeyCurveNameP256K KeyCurveName = "P-256K"
	KeyCurveNameP384  KeyCurveName = "P-384"
	KeyCurveNameP521  KeyCurveName = "P-521"
)

// PossibleKeyCurveNameValues returns a slice of all possible CertificateKeyCurveName values.
func PossibleKeyCurveNameValues() []KeyCurveName {
	return []KeyCurveName{
		KeyCurveNameP256,
		KeyCurveNameP256K,
		KeyCurveNameP384,
		KeyCurveNameP521,
	}
}

// KeyType - The type of key pair to be used for the certificate.
type KeyType string

const (
	KeyTypeEC     KeyType = "EC"
	KeyTypeECHSM  KeyType = "EC-HSM"
	KeyTypeRSA    KeyType = "RSA"
	KeyTypeRSAHSM KeyType = "RSA-HSM"
)

// PossibleKeyTypeValues returns a slice of all possible CertificateKeyType values.
func PossibleKeyTypeValues() []KeyType {
	return []KeyType{
		KeyTypeEC,
		KeyTypeECHSM,
		KeyTypeRSA,
		KeyTypeRSAHSM,
	}
}

// KeyUsage is the key usage for a certificate
type KeyUsage string

const (
	KeyUsageCRLSign          KeyUsage = "cRLSign"
	KeyUsageDataEncipherment KeyUsage = "dataEncipherment"
	KeyUsageDecipherOnly     KeyUsage = "decipherOnly"
	KeyUsageDigitalSignature KeyUsage = "digitalSignature"
	KeyUsageEncipherOnly     KeyUsage = "encipherOnly"
	KeyUsageKeyAgreement     KeyUsage = "keyAgreement"
	KeyUsageKeyCertSign      KeyUsage = "keyCertSign"
	KeyUsageKeyEncipherment  KeyUsage = "keyEncipherment"
	KeyUsageNonRepudiation   KeyUsage = "nonRepudiation"
)

// PossibleKeyUsageValues returns a slice of all possible CertificateKeyUsage values.
func PossibleKeyUsageValues() []KeyUsage {
	return []KeyUsage{
		KeyUsageCRLSign,
		KeyUsageDataEncipherment,
		KeyUsageDecipherOnly,
		KeyUsageDigitalSignature,
		KeyUsageEncipherOnly,
		KeyUsageKeyAgreement,
		KeyUsageKeyCertSign,
		KeyUsageKeyEncipherment,
		KeyUsageNonRepudiation,
	}
}

// CertificateContentType is the content type of the certificate.
type CertificateContentType string

const (
	CertificateContentTypePKCS12 CertificateContentType = "application/x-pkcs12"
	CertificateContentTypePEM    CertificateContentType = "application/x-pem-file"
)

// PossibleCertificateContentTypeValues returns a slice of all possible CertificateContentType values.
func PossibleCertificateContentTypeValues() []CertificateContentType {
	return []CertificateContentType{
		CertificateContentTypePEM,
		CertificateContentTypePKCS12,
	}
}

// WellKnownIssuerNames names you can use when creating a certificate policy
type WellKnownIssuerNames string

const (
	WellKnownIssuerNamesSelf    WellKnownIssuerNames = "Self"
	WellKnownIssuerNamesUnknown WellKnownIssuerNames = "Unknown"
)

// PossibleWellKnownIssuerNamesValues returns a slice of all possible WellKnownIssuer values.
func PossibleWellKnownIssuerNamesValues() []WellKnownIssuerNames {
	return []WellKnownIssuerNames{
		WellKnownIssuerNamesSelf,
		WellKnownIssuerNamesUnknown,
	}
}

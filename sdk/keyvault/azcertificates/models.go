//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates/internal/generated"
)

// Action - The action that will be executed.
type Action struct {
	// The type of the action.
	ActionType *ActionType `json:"action_type,omitempty"`
}

func (a *Action) toGenerated() *generated.Action {
	if a == nil {
		return nil
	}

	return &generated.Action{
		ActionType: (*generated.ActionType)(a.ActionType),
	}
}

// Attributes - The object attributes managed by the KeyVault service.
type Attributes struct {
	// Determines whether the object is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Expiry date in UTC.
	Expires *time.Time `json:"exp,omitempty"`

	// Not before date in UTC.
	NotBefore *time.Time `json:"nbf,omitempty"`

	// READ-ONLY; Creation time in UTC.
	Created *time.Time `json:"created,omitempty" azure:"ro"`

	// READ-ONLY; Last updated time in UTC.
	Updated *time.Time `json:"updated,omitempty" azure:"ro"`
}

// CertificateAttributes - The certificate management attributes.
type CertificateAttributes struct {
	Attributes
	// READ-ONLY; softDelete data retention days. Value should be >=7 and <=90 when softDelete enabled, otherwise 0.
	RecoverableDays *int32 `json:"recoverableDays,omitempty" azure:"ro"`

	// READ-ONLY; Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains 'Purgeable', the certificate
	// can be permanently deleted by a privileged user; otherwise,
	// only the system can purge the certificate, at the end of the retention interval.
	RecoveryLevel *DeletionRecoveryLevel `json:"recoveryLevel,omitempty" azure:"ro"`
}

func (c *CertificateAttributes) toGenerated() *generated.CertificateAttributes {
	if c == nil {
		return &generated.CertificateAttributes{}
	}

	return &generated.CertificateAttributes{
		RecoverableDays: c.RecoverableDays,
		RecoveryLevel:   (*generated.DeletionRecoveryLevel)(c.RecoveryLevel),
		Attributes: generated.Attributes{
			Enabled:   c.Enabled,
			Expires:   c.Expires,
			NotBefore: c.NotBefore,
			Created:   c.Created,
			Updated:   c.Updated,
		},
	}
}

// CertificateOperation - A certificate operation is returned in case of asynchronous requests.
type CertificateOperation struct {
	// Indicates if cancellation was requested on the certificate operation.
	CancellationRequested *bool `json:"cancellation_requested,omitempty"`

	// The certificate signing request (CSR) that is being used in the certificate operation.
	Csr []byte `json:"csr,omitempty"`

	// Error encountered, if any, during the certificate operation.
	Error *generated.Error `json:"error,omitempty"`

	// Parameters for the issuer of the X509 component of a certificate.
	IssuerParameters *generated.IssuerParameters `json:"issuer,omitempty"`

	// Identifier for the certificate operation.
	RequestID *string `json:"request_id,omitempty"`

	// Status of the certificate operation.
	Status *string `json:"status,omitempty"`

	// The status details of the certificate operation.
	StatusDetails *string `json:"status_details,omitempty"`

	// Location which contains the result of the certificate operation.
	Target *string `json:"target,omitempty"`

	// READ-ONLY; The certificate id.
	ID *string `json:"id,omitempty" azure:"ro"`
}

// CertificatePolicy - Management policy for a certificate.
type CertificatePolicy struct {
	// The certificate attributes.
	Attributes *CertificateAttributes `json:"attributes,omitempty"`

	// Parameters for the issuer of the X509 component of a certificate.
	IssuerParameters *IssuerParameters `json:"issuer,omitempty"`

	// Properties of the key backing a certificate.
	KeyProperties *KeyProperties `json:"key_props,omitempty"`

	// Actions that will be performed by Key Vault over the lifetime of a certificate.
	LifetimeActions []*LifetimeAction `json:"lifetime_actions,omitempty"`

	// Properties of the secret backing a certificate.
	SecretProperties *SecretProperties `json:"secret_props,omitempty"`

	// Properties of the X509 component of a certificate.
	X509CertificateProperties *X509CertificateProperties `json:"x509_props,omitempty"`
}

func (c CertificatePolicy) toGeneratedCertificateCreateParameters() *generated.CertificatePolicy {
	var lifetimeActions []*generated.LifetimeAction
	if len(c.LifetimeActions) > 0 {
		for _, la := range c.LifetimeActions {
			lifetimeActions = append(lifetimeActions, la.toGenerated())
		}
	}

	return &generated.CertificatePolicy{
		Attributes:                c.Attributes.toGenerated(),
		IssuerParameters:          c.IssuerParameters.toGenerated(),
		KeyProperties:             c.KeyProperties.toGenerated(),
		SecretProperties:          c.SecretProperties.toGenerated(),
		LifetimeActions:           lifetimeActions,
		X509CertificateProperties: c.X509CertificateProperties.toGenerated(),
	}
}

// IssuerParameters - Parameters for the issuer of the X509 component of a certificate.
type IssuerParameters struct {
	// Indicates if the certificates generated under this policy should be published to certificate transparency logs.
	CertificateTransparency *bool `json:"cert_transparency,omitempty"`

	// Certificate type as supported by the provider (optional); for example 'OV-SSL', 'EV-SSL'
	CertificateType *string `json:"cty,omitempty"`

	// Name of the referenced issuer object or reserved names; for example, 'Self' or 'Unknown'.
	Name *string `json:"name,omitempty"`
}

func (i *IssuerParameters) toGenerated() *generated.IssuerParameters {
	if i == nil {
		return &generated.IssuerParameters{}
	}

	return &generated.IssuerParameters{
		CertificateTransparency: i.CertificateTransparency,
		CertificateType:         i.CertificateType,
		Name:                    i.Name,
	}
}

// KeyProperties - Properties of the key pair backing a certificate.
type KeyProperties struct {
	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	Curve *JSONWebKeyCurveName `json:"crv,omitempty"`

	// Indicates if the private key can be exported.
	Exportable *bool `json:"exportable,omitempty"`

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// The type of key pair to be used for the certificate.
	KeyType *JSONWebKeyType `json:"kty,omitempty"`

	// Indicates if the same key pair will be used on certificate renewal.
	ReuseKey *bool `json:"reuse_key,omitempty"`
}

func (k *KeyProperties) toGenerated() *generated.KeyProperties {
	if k == nil {
		return nil
	}

	return &generated.KeyProperties{
		Curve:      (*generated.JSONWebKeyCurveName)(k.Curve),
		Exportable: k.Exportable,
		KeySize:    k.KeySize,
		KeyType:    (*generated.JSONWebKeyType)(k.KeyType),
		ReuseKey:   k.ReuseKey,
	}
}

// LifetimeAction - Action and its trigger that will be performed by Key Vault over the lifetime of a certificate.
type LifetimeAction struct {
	// The action that will be executed.
	Action *Action `json:"action,omitempty"`

	// The condition that will execute the action.
	Trigger *Trigger `json:"trigger,omitempty"`
}

func (l LifetimeAction) toGenerated() *generated.LifetimeAction {
	return &generated.LifetimeAction{
		Action:  l.Action.toGenerated(),
		Trigger: l.Trigger.toGenerated(),
	}
}

// SecretProperties - Properties of the key backing a certificate.
type SecretProperties struct {
	// The media type (MIME type).
	ContentType *string `json:"contentType,omitempty"`
}

func (s *SecretProperties) toGenerated() *generated.SecretProperties {
	if s == nil {
		return nil
	}

	return &generated.SecretProperties{
		ContentType: s.ContentType,
	}
}

// SubjectAlternativeNames - The subject alternate names of a X509 object.
type SubjectAlternativeNames struct {
	// Domain names.
	DNSNames []*string `json:"dns_names,omitempty"`

	// Email addresses.
	Emails []*string `json:"emails,omitempty"`

	// User principal names.
	Upns []*string `json:"upns,omitempty"`
}

func (s *SubjectAlternativeNames) toGenerated() *generated.SubjectAlternativeNames {
	if s == nil {
		return &generated.SubjectAlternativeNames{}
	}

	return &generated.SubjectAlternativeNames{
		DNSNames: s.DNSNames,
		Emails:   s.Emails,
		Upns:     s.Upns,
	}
}

// Trigger - A condition to be satisfied for an action to be executed.
type Trigger struct {
	// Days before expiry to attempt renewal. Value should be between 1 and validityinmonths multiplied by 27. If validityinmonths is 36, then value should
	// be between 1 and 972 (36 * 27).
	DaysBeforeExpiry *int32 `json:"days_before_expiry,omitempty"`

	// Percentage of lifetime at which to trigger. Value should be between 1 and 99.
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

func (t *Trigger) toGenerated() *generated.Trigger {
	if t == nil {
		return nil
	}

	return &generated.Trigger{
		DaysBeforeExpiry:   t.DaysBeforeExpiry,
		LifetimePercentage: t.LifetimePercentage,
	}
}

// X509CertificateProperties - Properties of the X509 component of a certificate.
type X509CertificateProperties struct {
	// The enhanced key usage.
	Ekus []*string `json:"ekus,omitempty"`

	// List of key usages.
	KeyUsage []*KeyUsageType `json:"key_usage,omitempty"`

	// The subject name. Should be a valid X509 distinguished Name.
	Subject *string `json:"subject,omitempty"`

	// The subject alternative names.
	SubjectAlternativeNames *SubjectAlternativeNames `json:"sans,omitempty"`

	// The duration that the certificate is valid in months.
	ValidityInMonths *int32 `json:"validity_months,omitempty"`
}

func (x *X509CertificateProperties) toGenerated() *generated.X509CertificateProperties {
	if x == nil {
		return &generated.X509CertificateProperties{}
	}

	var keyUsage []*generated.KeyUsageType
	for _, k := range x.KeyUsage {
		keyUsage = append(keyUsage, (*generated.KeyUsageType)(k))
	}

	return &generated.X509CertificateProperties{
		Ekus:                    x.Ekus,
		KeyUsage:                keyUsage,
		Subject:                 x.Subject,
		SubjectAlternativeNames: x.SubjectAlternativeNames.toGenerated(),
		ValidityInMonths:        x.ValidityInMonths,
	}
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// AdministratorContact - Details of the organization administrator of the certificate issuer.
type AdministratorContact struct {
	// Email address.
	Email *string

	// First name.
	FirstName *string

	// Last name.
	LastName *string

	// Phone number.
	Phone *string
}

// Properties - The certificate management properties.
type Properties struct {
	// READ-ONLY; Creation time in UTC.
	CreatedOn *time.Time

	// Determines whether the object is enabled.
	Enabled *bool

	// Expiry date in UTC.
	ExpiresOn *time.Time

	// Not before date in UTC.
	NotBefore *time.Time

	// READ-ONLY; softDelete data retention days. Value should be >=7 and <=90 when softDelete enabled, otherwise 0.
	RecoverableDays *int32

	// READ-ONLY; Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains
	// 'Purgeable', the certificate can be permanently deleted by a privileged user; otherwise,
	// only the system can purge the certificate, at the end of the retention interval.
	RecoveryLevel *string

	// READ-ONLY; Last updated time in UTC.
	UpdatedOn *time.Time

	// READ-ONLY; The ID of the certificate
	ID *string

	// READ-ONLY; The name of the certificate
	Name *string

	// Application specific metadata in the form of key-value pairs
	Tags map[string]*string

	// READ-ONLY; The vault URL for the certificate
	VaultURL *string

	// READ-ONLY; The version for the certificate
	Version *string

	// Thumbprint of the certificate.
	X509Thumbprint []byte
}

func (c *Properties) toGenerated() *generated.CertificateAttributes {
	if c == nil {
		return nil
	}

	return &generated.CertificateAttributes{
		Created:         c.CreatedOn,
		Enabled:         c.Enabled,
		Expires:         c.ExpiresOn,
		NotBefore:       c.NotBefore,
		RecoverableDays: c.RecoverableDays,
		RecoveryLevel:   (*generated.DeletionRecoveryLevel)(c.RecoveryLevel),
		Updated:         c.UpdatedOn,
	}
}

func propertiesFromGenerated(g *generated.CertificateAttributes, tags map[string]*string, id *string, thumbprint []byte) *Properties {
	if g == nil {
		return nil
	}

	vaulURL, name, version := shared.ParseID(id)

	return &Properties{
		Enabled:         g.Enabled,
		ExpiresOn:       g.Expires,
		NotBefore:       g.NotBefore,
		CreatedOn:       g.Created,
		UpdatedOn:       g.Updated,
		RecoverableDays: g.RecoverableDays,
		RecoveryLevel:   (*string)(g.RecoveryLevel),
		Tags:            tags,
		ID:              id,
		Name:            name,
		VaultURL:        vaulURL,
		Version:         version,
		X509Thumbprint:  thumbprint,
	}
}

// Certificate - A certificate bundle consists of a certificate (X509) plus its properties.
type Certificate struct {
	// The certificate attributes.
	Properties *Properties

	// CER contents of x509 certificate.
	CER []byte

	// READ-ONLY; The certificate id.
	ID *string

	// READ-ONLY; The name of the certificate
	Name *string

	// READ-ONLY; The key ID.
	KeyID *string

	// READ-ONLY; The secret ID.
	SecretID *string
}

// CertificateWithPolicy - A certificate bundle consists of a certificate (X509) with a policy, and its properties.
type CertificateWithPolicy struct {
	// The certificate attributes.
	Properties *Properties

	// CER contents of x509 certificate.
	CER []byte

	// The content type of the secret. eg. 'application/x-pem-file' or 'application/x-pkcs12',
	ContentType *string

	// READ-ONLY; The certificate id.
	ID *string

	// READ-ONLY; The key ID.
	KeyID *string

	// READ-ONLY; The management policy.
	Policy *Policy

	// READ-ONLY; The secret ID.
	SecretID *string
}

// UnmarshalJSON implements the json.Unmarshaler interface for the CertificateWithPolicy type.
func (c *CertificateWithPolicy) UnmarshalJSON(data []byte) error {
	var g generated.CertificateBundle
	err := json.Unmarshal(data, &g)
	if err != nil {
		return err
	}
	c.Properties = propertiesFromGenerated(g.Attributes, g.Tags, g.ID, g.X509Thumbprint)
	c.CER = g.Cer
	c.ContentType = g.ContentType
	c.ID = g.ID
	c.KeyID = g.Kid
	c.Policy = certificatePolicyFromGenerated(g.Policy)
	c.SecretID = g.Sid
	return nil
}

func certificateFromGenerated(g *generated.CertificateBundle) Certificate {
	if g == nil {
		return Certificate{}
	}

	_, name, _ := shared.ParseID(g.ID)
	return Certificate{
		Properties: propertiesFromGenerated(g.Attributes, g.Tags, g.ID, g.X509Thumbprint),
		CER:        g.Cer,
		ID:         g.ID,
		Name:       name,
		KeyID:      g.Kid,
		SecretID:   g.Sid,
	}
}

// CertificateOperationError - The key vault server error.
type CertificateOperationError struct {
	// READ-ONLY; The error code.
	Code *string

	// READ-ONLY; The key vault server error.
	innerError *CertificateOperationError

	// READ-ONLY; The error message.
	message *string
}

// Error returns the error string detailing why the Certificate Operation failed.
func (c *CertificateOperationError) Error() string {
	if c == nil {
		return ""
	}
	marshalled, err := json.Marshal(*c)
	if err != nil {
		return fmt.Sprintf("could not turn operation error into a string: %v", err)
	}
	return string(marshalled)
}

func certificateErrorFromGenerated(g *generated.Error) *CertificateOperationError {
	if g == nil {
		return nil
	}

	return &CertificateOperationError{
		Code:       g.Code,
		message:    g.Message,
		innerError: certificateErrorFromGenerated(g.InnerError),
	}
}

// IssuerItem - The certificate issuer item containing certificate issuer metadata.
type IssuerItem struct {
	// Certificate Identifier.
	ID *string

	// The issuer provider.
	Provider *string
}

func certificateIssuerItemFromGenerated(g *generated.CertificateIssuerItem) *IssuerItem {
	if g == nil {
		return nil
	}

	return &IssuerItem{ID: g.ID, Provider: g.Provider}
}

// CertificateItem - The certificate item containing certificate metadata.
type CertificateItem struct {
	// The certificate management attributes.
	Properties *Properties

	// Certificate identifier.
	ID *string
}

// Operation - A certificate operation is returned in case of asynchronous requests.
type Operation struct {
	// Indicates if cancellation was requested on the certificate operation.
	CancellationRequested *bool

	// The certificate signing request (CSR) that is being used in the certificate operation.
	CSR []byte

	// Error encountered, if any, during the certificate operation.
	Error *CertificateOperationError

	// Parameters for the issuer of the X509 component of a certificate.
	IssuerParameters *IssuerParameters

	// Identifier for the certificate operation.
	RequestID *string

	// Status of the certificate operation.
	Status *string

	// The status details of the certificate operation.
	StatusDetails *string

	// Location which contains the result of the certificate operation.
	Target *string

	// READ-ONLY; The certificate id.
	ID *string
}

func certificateOperationFromGenerated(g generated.CertificateOperation) Operation {
	return Operation{
		CancellationRequested: g.CancellationRequested,
		CSR:                   g.Csr,
		Error:                 certificateErrorFromGenerated(g.Error),
		IssuerParameters:      issuerParametersFromGenerated(g.IssuerParameters),
		RequestID:             g.RequestID,
		Status:                g.Status,
		StatusDetails:         g.StatusDetails,
		Target:                g.Target,
		ID:                    g.ID,
	}
}

// Policy - Management policy for a certificate.
type Policy struct {
	// The certificate properties.
	Properties *Properties

	// Parameters for the issuer of the X509 component of a certificate.
	IssuerParameters *IssuerParameters

	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	KeyCurveName *KeyCurveName

	// Indicates if the private key can be exported.
	Exportable *bool

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32

	// The type of key pair to be used for the certificate.
	KeyType *KeyType

	// Indicates if the same key pair will be used on certificate renewal.
	ReuseKey *bool

	// Actions that will be performed by Key Vault over the lifetime of a certificate.
	LifetimeActions []*LifetimeAction

	// ContentType of the downloaded certificate
	ContentType *CertificateContentType

	// Properties of the X509 component of a certificate.
	X509Properties *X509CertificateProperties
}

// NewDefaultCertificatePolicy returns a Policy with IssuerName "Self" and Subject "CN=DefaultPolicy"
func NewDefaultCertificatePolicy() Policy {
	return Policy{
		IssuerParameters: &IssuerParameters{IssuerName: (*string)(to.Ptr(WellKnownIssuerNamesSelf))},
		X509Properties:   &X509CertificateProperties{Subject: to.Ptr("CN=DefaultPolicy")},
	}
}

func (c *Policy) toGeneratedCertificateCreateParameters() *generated.CertificatePolicy {
	if c == nil {
		return nil
	}
	var la []*generated.LifetimeAction
	for _, l := range c.LifetimeActions {
		la = append(la, l.toGenerated())
	}

	var keyProps *generated.KeyProperties
	if c.KeyCurveName != nil || c.Exportable != nil || c.KeySize != nil || c.KeyType != nil || c.ReuseKey != nil {
		keyProps = &generated.KeyProperties{
			Curve:      (*generated.JSONWebKeyCurveName)(c.KeyCurveName),
			Exportable: c.Exportable,
			KeySize:    c.KeySize,
			KeyType:    (*generated.JSONWebKeyType)(c.KeyType),
			ReuseKey:   c.ReuseKey,
		}
	}

	return &generated.CertificatePolicy{
		Attributes:                c.Properties.toGenerated(),
		IssuerParameters:          c.IssuerParameters.toGenerated(),
		KeyProperties:             keyProps,
		SecretProperties:          &generated.SecretProperties{ContentType: (*string)(c.ContentType)},
		LifetimeActions:           la,
		X509CertificateProperties: c.X509Properties.toGenerated(),
	}
}

func certificatePolicyFromGenerated(g *generated.CertificatePolicy) *Policy {
	if g == nil {
		return nil
	}

	var la []*LifetimeAction
	for _, l := range g.LifetimeActions {
		la = append(la, lifetimeActionFromGenerated(l))
	}

	c := &Policy{}
	if g.KeyProperties != nil {
		c.KeyCurveName = (*KeyCurveName)(g.KeyProperties.Curve)
		c.Exportable = g.KeyProperties.Exportable
		c.KeySize = g.KeyProperties.KeySize
		c.KeyType = (*KeyType)(g.KeyProperties.KeyType)
		c.ReuseKey = g.KeyProperties.ReuseKey
	}

	c.Properties = propertiesFromGenerated(g.Attributes, nil, nil, nil)
	c.IssuerParameters = issuerParametersFromGenerated(g.IssuerParameters)
	c.LifetimeActions = la
	c.ContentType = (*CertificateContentType)(g.SecretProperties.ContentType)
	c.X509Properties = x509CertificatePropertiesFromGenerated(g.X509CertificateProperties)
	return c
}

// Contact - The contact information for the vault certificates.
type Contact struct {
	// Email address.
	Email *string

	// Name.
	Name *string

	// Phone number.
	Phone *string
}

func contactListFromGenerated(g []*generated.Contact) []*Contact {
	var ret []*Contact
	for _, c := range g {
		ret = append(ret, &Contact{
			Email: c.EmailAddress,
			Name:  c.Name,
			Phone: c.Phone,
		})
	}
	return ret
}

// Contacts - The contacts for the vault certificates.
type Contacts struct {
	// The contact list for the vault certificates.
	ContactList []*Contact

	// READ-ONLY; Identifier for the contacts collection.
	ID *string
}

func (c *Contacts) toGenerated() generated.Contacts {
	if c == nil {
		return generated.Contacts{}
	}

	var contacts []*generated.Contact
	for _, contact := range c.ContactList {
		contacts = append(contacts, &generated.Contact{
			EmailAddress: contact.Email,
			Name:         contact.Name,
			Phone:        contact.Phone,
		})
	}

	return generated.Contacts{
		ID:          c.ID,
		ContactList: contacts,
	}
}

// DeletedCertificate consists of its previous id, attributes and its tags, as well as information on when it will be purged.
type DeletedCertificate struct {
	// The certificate properties.
	Properties *Properties

	// CER contents of x509 certificate.
	CER []byte

	// The content type of the secret. eg. 'application/x-pem-file' or 'application/x-pkcs12',
	ContentType *string

	// The url of the recovery object, used to identify and recover the deleted certificate.
	RecoveryID *string

	// READ-ONLY; The time when the certificate was deleted, in UTC
	DeletedOn *time.Time

	// READ-ONLY; The certificate id.
	ID *string

	// READ-ONLY; The name of the certificate
	Name *string

	// READ-ONLY; The key ID.
	KeyID *string

	// READ-ONLY; The management policy.
	Policy *Policy

	// READ-ONLY; The time when the certificate is scheduled to be purged, in UTC
	ScheduledPurgeDate *time.Time

	// READ-ONLY; The secret ID.
	SecretID *string
}

// DeletedCertificateItem - The deleted certificate item containing metadata about the deleted certificate.
type DeletedCertificateItem struct {
	// The certificate management properties.
	Properties *Properties

	// Certificate identifier.
	ID *string

	// READ-ONLY; The name of the certificate
	Name *string

	// The url of the recovery object, used to identify and recover the deleted certificate.
	RecoveryID *string

	// READ-ONLY; The time when the certificate was deleted, in UTC
	DeletedOn *time.Time

	// READ-ONLY; The time when the certificate is scheduled to be purged, in UTC
	ScheduledPurgeDate *time.Time
}

// Issuer - The issuer for Key Vault certificate.
type Issuer struct {
	// Determines whether the issuer is enabled.
	Enabled *bool

	// READ-ONLY; Creation time in UTC.
	CreatedOn *time.Time

	// READ-ONLY; Last updated time in UTC.
	UpdatedOn *time.Time

	// The credentials to be used for the issuer.
	Credentials *IssuerCredentials

	// Details of the organization administrator.
	AdministratorContacts []*AdministratorContact

	// Id of the organization.
	OrganizationID *string

	// The issuer provider.
	Provider *string

	// READ-ONLY; Identifier for the issuer object.
	ID *string

	// READ-ONLY; Name is the name of the issuer
	Name *string
}

// IssuerCredentials - The credentials to be used for the certificate issuer.
type IssuerCredentials struct {
	// The user name/account name/account id.
	AccountID *string

	// The password/secret/account key.
	Password *string
}

func issuerCredentialsFromGenerated(g *generated.IssuerCredentials) *IssuerCredentials {
	if g == nil {
		return nil
	}

	return &IssuerCredentials{AccountID: g.AccountID, Password: g.Password}
}

func (i *IssuerCredentials) toGenerated() *generated.IssuerCredentials {
	if i == nil {
		return nil
	}

	return &generated.IssuerCredentials{
		AccountID: i.AccountID,
		Password:  i.Password,
	}
}

// IssuerParameters - Parameters for the issuer of the X509 component of a certificate.
type IssuerParameters struct {
	// Indicates if the certificates generated under this policy should be published to certificate transparency logs.
	CertificateTransparency *bool

	// Certificate type as supported by the provider (optional); for example 'OV-SSL', 'EV-SSL'
	CertificateType *string

	// IssuerName of the referenced issuer object or reserved names; for example, 'Self' or 'Unknown'.
	IssuerName *string
}

func (i *IssuerParameters) toGenerated() *generated.IssuerParameters {
	if i == nil {
		return &generated.IssuerParameters{}
	}

	return &generated.IssuerParameters{
		CertificateTransparency: i.CertificateTransparency,
		CertificateType:         i.CertificateType,
		Name:                    i.IssuerName,
	}
}

func issuerParametersFromGenerated(g *generated.IssuerParameters) *IssuerParameters {
	if g == nil {
		return nil
	}

	return &IssuerParameters{
		CertificateTransparency: g.CertificateTransparency,
		CertificateType:         g.CertificateType,
		IssuerName:              g.Name,
	}
}

// LifetimeAction - Action and its trigger that will be performed by Key Vault over the lifetime of a certificate.
type LifetimeAction struct {
	// The action that will be executed.
	Action *PolicyAction
	// Days before expiry to attempt renewal. Value should be between 1 and validityinmonths multiplied by 27. If validityinmonths is 36, then value should
	// be between 1 and 972 (36 * 27).
	DaysBeforeExpiry *int32

	// Percentage of lifetime at which to trigger. Value should be between 1 and 99.
	LifetimePercentage *int32
}

func (l LifetimeAction) toGenerated() *generated.LifetimeAction {
	return &generated.LifetimeAction{
		Action: &generated.Action{
			ActionType: (*generated.ActionType)(l.Action),
		},
		Trigger: &generated.Trigger{
			DaysBeforeExpiry:   l.DaysBeforeExpiry,
			LifetimePercentage: l.LifetimePercentage,
		},
	}
}

func lifetimeActionFromGenerated(g *generated.LifetimeAction) *LifetimeAction {
	if g == nil {
		return nil
	}

	return &LifetimeAction{
		Action:             (*PolicyAction)(g.Action.ActionType),
		DaysBeforeExpiry:   g.Trigger.DaysBeforeExpiry,
		LifetimePercentage: g.Trigger.LifetimePercentage,
	}
}

// SubjectAlternativeNames - The subject alternate names of a X509 object.
type SubjectAlternativeNames struct {
	// Domain names.
	DNSNames []*string

	// Email addresses.
	Emails []*string

	// User principal names.
	UserPrincipalNames []*string
}

func (s *SubjectAlternativeNames) toGenerated() *generated.SubjectAlternativeNames {
	if s == nil {
		return &generated.SubjectAlternativeNames{}
	}

	return &generated.SubjectAlternativeNames{
		DNSNames: s.DNSNames,
		Emails:   s.Emails,
		Upns:     s.UserPrincipalNames,
	}
}

func subjectAlternativeNamesFromGenerated(g *generated.SubjectAlternativeNames) *SubjectAlternativeNames {
	if g == nil {
		return nil
	}

	return &SubjectAlternativeNames{
		DNSNames:           g.DNSNames,
		Emails:             g.Emails,
		UserPrincipalNames: g.Upns,
	}
}

// X509CertificateProperties - Properties of the X509 component of a certificate.
type X509CertificateProperties struct {
	// The enhanced key usage.
	EnhancedKeyUsages []*string

	// List of key usages.
	KeyUsages []*KeyUsage

	// The subject name. Should be a valid X509 distinguished Name.
	Subject *string

	// The subject alternative names.
	SubjectAlternativeNames *SubjectAlternativeNames

	// The duration that the certificate is valid in months.
	ValidityInMonths *int32
}

func (x *X509CertificateProperties) toGenerated() *generated.X509CertificateProperties {
	if x == nil {
		return &generated.X509CertificateProperties{}
	}

	var keyUsage []*generated.KeyUsageType
	for _, k := range x.KeyUsages {
		keyUsage = append(keyUsage, (*generated.KeyUsageType)(k))
	}

	return &generated.X509CertificateProperties{
		Ekus:                    x.EnhancedKeyUsages,
		KeyUsage:                keyUsage,
		Subject:                 x.Subject,
		SubjectAlternativeNames: x.SubjectAlternativeNames.toGenerated(),
		ValidityInMonths:        x.ValidityInMonths,
	}
}

func x509CertificatePropertiesFromGenerated(g *generated.X509CertificateProperties) *X509CertificateProperties {
	if g == nil {
		return nil
	}

	var ku []*KeyUsage
	for _, k := range g.KeyUsage {
		ku = append(ku, (*KeyUsage)(k))
	}

	return &X509CertificateProperties{
		EnhancedKeyUsages:       g.Ekus,
		Subject:                 g.Subject,
		KeyUsages:               ku,
		SubjectAlternativeNames: subjectAlternativeNamesFromGenerated(g.SubjectAlternativeNames),
		ValidityInMonths:        g.ValidityInMonths,
	}
}

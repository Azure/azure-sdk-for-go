//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Properties - The properties of a key managed by the key vault service.
type Properties struct {
	// READ-ONLY; Creation time in UTC.
	CreatedOn *time.Time `json:"created,omitempty" azure:"ro"`

	// Determines whether the object is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Expiry date in UTC.
	ExpiresOn *time.Time `json:"exp,omitempty"`

	// Indicates if the private key can be exported.
	Exportable *bool `json:"exportable,omitempty"`

	// ID identifies the key
	ID *string

	// READ-ONLY; True if the key's lifetime is managed by key vault. If this is a key backing a certificate, then managed will be true.
	Managed *bool `json:"managed,omitempty" azure:"ro"`

	// Name of the key
	Name *string

	// Not before date in UTC.
	NotBefore *time.Time `json:"nbf,omitempty"`

	// READ-ONLY; softDelete data retention days. Value should be >=7 and <=90 when softDelete enabled, otherwise 0.
	RecoverableDays *int32 `json:"recoverableDays,omitempty" azure:"ro"`

	// READ-ONLY; Reflects the deletion recovery level currently in effect for keys in the current vault. If it contains 'Purgeable'
	// the key can be permanently deleted by a privileged user; otherwise, only the system
	// can purge the key, at the end of the retention interval.
	RecoveryLevel *string `json:"recoveryLevel,omitempty" azure:"ro"`

	// Tags contain application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Last updated time in UTC.
	UpdatedOn *time.Time `json:"updated,omitempty" azure:"ro"`

	// VaultURL for the key
	VaultURL *string

	// Version of the key
	Version *string
}

// converts a KeyAttributes to *generated.KeyAttributes
func (k *Properties) toGenerated() *generated.KeyAttributes {
	if k == nil {
		return nil
	}
	return &generated.KeyAttributes{
		RecoverableDays: k.RecoverableDays,
		RecoveryLevel:   toGeneratedDeletionRecoveryLevel(k.RecoveryLevel),
		Enabled:         k.Enabled,
		Expires:         k.ExpiresOn,
		NotBefore:       k.NotBefore,
		Created:         k.CreatedOn,
		Updated:         k.UpdatedOn,
		Exportable:      k.Exportable,
	}
}

// converts *generated.KeyAttributes to *KeyAttributes
func keyPropertiesFromGenerated(i *generated.KeyAttributes, id *string, name *string, version *string, managed *bool, vaultURL *string, tags map[string]*string) *Properties {
	if i == nil {
		return &Properties{}
	}

	return &Properties{
		CreatedOn:       i.Created,
		Enabled:         i.Enabled,
		ExpiresOn:       i.Expires,
		Exportable:      i.Exportable,
		ID:              id,
		Managed:         managed,
		Name:            name,
		NotBefore:       i.NotBefore,
		RecoverableDays: i.RecoverableDays,
		RecoveryLevel:   to.Ptr(string(*i.RecoveryLevel)),
		Tags:            tags,
		UpdatedOn:       i.Updated,
		VaultURL:        vaultURL,
		Version:         version,
	}
}

// Key - A Key consists of a WebKey plus its attributes.
type Key struct {
	// The key management properties.
	Properties *Properties `json:"attributes,omitempty"`

	// The Json web key.
	JSONWebKey *JSONWebKey `json:"key,omitempty"`

	// The policy rules under which the key can be exported.
	ReleasePolicy *ReleasePolicy `json:"release_policy,omitempty"`

	// ID identifies the key
	ID *string

	// READ-ONLY; The name of the key
	Name *string
}

// convert the options to generated.KeyUpdateParameters struct
func (k Key) toKeyUpdateParameters() generated.KeyUpdateParameters {
	var attribs *generated.KeyAttributes
	if k.Properties != nil {
		attribs = k.Properties.toGenerated()
	}

	var tags map[string]*string
	if k.Properties != nil && k.Properties.Tags != nil {
		tags = k.Properties.Tags
	}

	return generated.KeyUpdateParameters{
		KeyAttributes: attribs,
		ReleasePolicy: k.ReleasePolicy.toGenerated(),
		Tags:          tags,
	}
}

// JSONWebKey - As of http://tools.ietf.org/html/draft-ietf-jose-json-web-key-18
type JSONWebKey struct {
	// Elliptic curve name. For valid values, see PossibleCurveNameValues.
	Crv *CurveName `json:"crv,omitempty"`

	// RSA private exponent, or the D component of an EC private key.
	D []byte `json:"d,omitempty"`

	// RSA private key parameter.
	DP []byte `json:"dp,omitempty"`

	// RSA private key parameter.
	DQ []byte `json:"dq,omitempty"`

	// RSA public exponent.
	E []byte `json:"e,omitempty"`

	// Symmetric key.
	K      []byte       `json:"k,omitempty"`
	KeyOps []*Operation `json:"key_ops,omitempty"`

	// ID identifies the key
	ID *string `json:"kid,omitempty"`

	// JsonWebKey Key Type (kty), as defined in https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-40.
	KeyType *KeyType `json:"kty,omitempty"`

	// RSA modulus.
	N []byte `json:"n,omitempty"`

	// RSA secret prime.
	P []byte `json:"p,omitempty"`

	// RSA secret prime, with p < q.
	Q []byte `json:"q,omitempty"`

	// RSA private key parameter.
	QI []byte `json:"qi,omitempty"`

	// Protected Key, used with 'Bring Your Own Key'.
	T []byte `json:"key_hsm,omitempty"`

	// X component of an EC public key.
	X []byte `json:"x,omitempty"`

	// Y component of an EC public key.
	Y []byte `json:"y,omitempty"`
}

// converts generated.JSONWebKey to publicly exposed version
func jsonWebKeyFromGenerated(i *generated.JSONWebKey) *JSONWebKey {
	if i == nil {
		return &JSONWebKey{}
	}

	ops := make([]*Operation, len(i.KeyOps))
	for j, op := range i.KeyOps {
		ops[j] = (*Operation)(op)
	}

	return &JSONWebKey{
		Crv:     (*CurveName)(i.Crv),
		D:       i.D,
		DP:      i.DP,
		DQ:      i.DQ,
		E:       i.E,
		K:       i.K,
		KeyOps:  ops,
		ID:      i.Kid,
		KeyType: (*KeyType)(i.Kty),
		N:       i.N,
		P:       i.P,
		Q:       i.Q,
		QI:      i.QI,
		T:       i.T,
		X:       i.X,
		Y:       i.Y,
	}
}

// converts JSONWebKey to *generated.JSONWebKey
func (j JSONWebKey) toGenerated() *generated.JSONWebKey {
	ops := make([]*string, len(j.KeyOps))
	for i, op := range j.KeyOps {
		ops[i] = (*string)(op)
	}
	return &generated.JSONWebKey{
		Crv:    (*generated.JSONWebKeyCurveName)(j.Crv),
		D:      j.D,
		DP:     j.DP,
		DQ:     j.DQ,
		E:      j.E,
		K:      j.K,
		KeyOps: ops,
		Kid:    j.ID,
		Kty:    (*generated.JSONWebKeyType)(j.KeyType),
		N:      j.N,
		P:      j.P,
		Q:      j.Q,
		QI:     j.QI,
		T:      j.T,
		X:      j.X,
		Y:      j.Y,
	}
}

// KeyItem - The key item containing key metadata.
type KeyItem struct {
	// The key management properties.
	Properties *Properties `json:"attributes,omitempty"`

	// ID identifies the key
	ID *string `json:"kid,omitempty"`

	// Name of the key
	Name *string
}

// convert *generated.KeyItem to *KeyItem
func keyItemFromGenerated(i *generated.KeyItem) *KeyItem {
	if i == nil {
		return nil
	}

	_, name, _ := shared.ParseID(i.Kid)
	return &KeyItem{
		Properties: keyPropertiesFromGenerated(i.Attributes, nil, nil, nil, i.Managed, nil, i.Tags),
		ID:         i.Kid,
		Name:       name,
	}
}

// DeletedKey - A DeletedKey consisting of a WebKey plus its Attributes and deletion info
type DeletedKey struct {
	// The key management properties.
	Properties *Properties `json:"attributes,omitempty"`

	// The Json web key.
	Key *JSONWebKey `json:"key,omitempty"`

	// The url of the recovery object, used to identify and recover the deleted key.
	RecoveryID *string `json:"recoveryId,omitempty"`

	// The policy rules under which the key can be exported.
	ReleasePolicy *ReleasePolicy `json:"release_policy,omitempty"`

	// READ-ONLY; The time when the key was deleted, in UTC
	DeletedOn *time.Time `json:"deletedDate,omitempty" azure:"ro"`

	// READ-ONLY; The time when the key is scheduled to be purged, in UTC
	ScheduledPurgeDate *time.Time `json:"scheduledPurgeDate,omitempty" azure:"ro"`
}

// DeletedKeyItem - The deleted key item containing the deleted key metadata and information about deletion.
type DeletedKeyItem struct {
	// The key management attributes.
	Properties *Properties `json:"attributes,omitempty"`

	// ID identifies the key
	ID *string `json:"kid,omitempty"`

	// Name of the key
	Name *string

	// The url of the recovery object, used to identify and recover the deleted key.
	RecoveryID *string `json:"recoveryId,omitempty"`

	// READ-ONLY; The time when the key was deleted, in UTC
	DeletedOn *time.Time `json:"deletedDate,omitempty" azure:"ro"`

	// READ-ONLY; The time when the key is scheduled to be purged, in UTC
	ScheduledPurgeDate *time.Time `json:"scheduledPurgeDate,omitempty" azure:"ro"`
}

// convert *generated.DeletedKeyItem to *DeletedKeyItem
func deletedKeyItemFromGenerated(i *generated.DeletedKeyItem) *DeletedKeyItem {
	if i == nil {
		return nil
	}

	vaultURL, name, version := shared.ParseID(i.Kid)
	return &DeletedKeyItem{
		RecoveryID:         i.RecoveryID,
		DeletedOn:          i.DeletedDate,
		ScheduledPurgeDate: i.ScheduledPurgeDate,
		Properties:         keyPropertiesFromGenerated(i.Attributes, i.Kid, name, version, i.Managed, vaultURL, i.Tags),
		ID:                 i.Kid,
		Name:               name,
	}
}

// ReleasePolicy represents the release policy for a Key Vault Key. For more information regarding
// the release policy grammar for Azure Key Vault, please refer to:
//  - https://aka.ms/policygrammarkeys for Azure Key Vault release policy grammar.
//  - https://aka.ms/policygrammarhsm for Azure Managed HSM release policy grammar.
type ReleasePolicy struct {
	// Content type and version of key release policy
	ContentType *string `json:"contentType,omitempty"`

	// Blob encoding the policy rules under which the key can be released.
	EncodedPolicy []byte `json:"data,omitempty"`

	// Defines the mutability state of the policy. Once marked immutable, this flag cannot be reset and the policy cannot be changed
	// under any circumstances.
	Immutable *bool `json:"immutable,omitempty"`
}

func (k *ReleasePolicy) toGenerated() *generated.KeyReleasePolicy {
	if k == nil {
		return nil
	}

	return &generated.KeyReleasePolicy{
		ContentType:   k.ContentType,
		EncodedPolicy: k.EncodedPolicy,
		Immutable:     k.Immutable,
	}
}

func keyReleasePolicyFromGenerated(i *generated.KeyReleasePolicy) *ReleasePolicy {
	if i == nil {
		return nil
	}
	return &ReleasePolicy{
		ContentType:   i.ContentType,
		EncodedPolicy: i.EncodedPolicy,
		Immutable:     i.Immutable,
	}
}

// RotationPolicy - Management policy for a key.
type RotationPolicy struct {
	// The key rotation policy attributes.
	Attributes *RotationPolicyAttributes `json:"attributes,omitempty"`

	// Actions that will be performed by Key Vault over the lifetime of a key. For preview, lifetimeActions can only have two items at maximum: one for rotate,
	// one for notify. Notification time would be
	// default to 30 days before expiry and it is not configurable.
	LifetimeActions []*LifetimeActions `json:"lifetimeActions,omitempty"`

	// READ-ONLY; The key policy id.
	ID *string `json:"id,omitempty" azure:"ro"`
}

func (u RotationPolicy) toGenerated() generated.KeyRotationPolicy {
	var attribs *generated.KeyRotationPolicyAttributes
	if u.Attributes != nil {
		attribs = u.Attributes.toGenerated()
	}
	la := make([]*generated.LifetimeActions, len(u.LifetimeActions))
	for i, l := range u.LifetimeActions {
		la[i] = l.toGenerated()
	}

	return generated.KeyRotationPolicy{
		ID:              u.ID,
		LifetimeActions: la,
		Attributes:      attribs,
	}
}

// RotationPolicyAttributes - The key rotation policy attributes.
type RotationPolicyAttributes struct {
	// ExpiresIn will be applied on the new key version. It should be at least 28 days.
	// It will be in ISO 8601 Format. Examples: 90 days: P90D, 3 months: P3M, 48 hours: PT48H, 1 year and 10 days: P1Y10D
	// It should be at least 28 days.
	ExpiresIn *string `json:"expiryTime,omitempty"`

	// READ-ONLY; The key rotation policy created time in UTC.
	CreatedOn *time.Time `json:"created,omitempty" azure:"ro"`

	// READ-ONLY; The key rotation policy's last updated time in UTC.
	UpdatedOn *time.Time `json:"updated,omitempty" azure:"ro"`
}

func (k RotationPolicyAttributes) toGenerated() *generated.KeyRotationPolicyAttributes {
	return &generated.KeyRotationPolicyAttributes{
		ExpiryTime: k.ExpiresIn,
		Created:    k.CreatedOn,
		Updated:    k.UpdatedOn,
	}
}

// LifetimeActions - Action and its trigger that will be performed by Key Vault over the lifetime of a key.
type LifetimeActions struct {
	// The action that will be executed.
	Action *LifetimeActionsType `json:"action,omitempty"`

	// The condition that will execute the action.
	Trigger *LifetimeActionsTrigger `json:"trigger,omitempty"`
}

func (l *LifetimeActions) toGenerated() *generated.LifetimeActions {
	if l == nil {
		return nil
	}
	return &generated.LifetimeActions{
		Action: &generated.LifetimeActionsType{
			Type: (*generated.ActionType)(l.Action.Type),
		},
		Trigger: &generated.LifetimeActionsTrigger{
			TimeAfterCreate:  l.Trigger.TimeAfterCreate,
			TimeBeforeExpiry: l.Trigger.TimeBeforeExpiry,
		},
	}
}

func lifetimeActionsFromGenerated(i *generated.LifetimeActions) *LifetimeActions {
	if i == nil {
		return nil
	}
	return &LifetimeActions{
		Trigger: &LifetimeActionsTrigger{
			TimeAfterCreate:  i.Trigger.TimeAfterCreate,
			TimeBeforeExpiry: i.Trigger.TimeBeforeExpiry,
		},
		Action: &LifetimeActionsType{
			Type: (*RotationAction)(i.Action.Type),
		},
	}
}

// LifetimeActionsType - The action that will be executed.
type LifetimeActionsType struct {
	// The type of the action.
	Type *RotationAction `json:"type,omitempty"`
}

// LifetimeActionsTrigger - A condition to be satisfied for an action to be executed.
type LifetimeActionsTrigger struct {
	// Time after creation to attempt to rotate. It only applies to rotate. It will be in ISO 8601 duration format. Example: 90 days : "P90D"
	TimeAfterCreate *string `json:"timeAfterCreate,omitempty"`

	// Time before expiry to attempt to rotate or notify. It will be in ISO 8601 duration format. Example: 90 days : "P90D"
	TimeBeforeExpiry *string `json:"timeBeforeExpiry,omitempty"`
}

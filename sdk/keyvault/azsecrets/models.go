//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets/internal"
)

// DeletedSecret consists of the previous ID, attributes, tags, and information on when it will be purged.
type DeletedSecret struct {
	// The secret management attributes.
	Properties *Properties `json:"attributes,omitempty"`

	// The content type of the secret.
	ContentType *string `json:"contentType,omitempty"`

	// The secret id.
	ID *string `json:"id,omitempty"`

	// The url of the recovery object, used to identify and recover the deleted secret.
	RecoveryID *string `json:"recoveryId,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`

	// The secret value.
	Value *string `json:"value,omitempty"`

	// READ-ONLY; The time when the secret was deleted, in UTC
	DeletedOn *time.Time `json:"deletedDate,omitempty" azure:"ro"`

	// READ-ONLY; If this is a secret backing a KV certificate, then this field specifies the corresponding key backing the KV
	// certificate.
	KeyID *string `json:"kid,omitempty" azure:"ro"`

	// READ-ONLY; True if the secret's lifetime is managed by key vault. If this is a secret backing a certificate, then managed
	// will be true.
	Managed *bool `json:"managed,omitempty" azure:"ro"`

	// READ-ONLY; The time when the secret is scheduled to be purged, in UTC
	ScheduledPurgeDate *time.Time `json:"scheduledPurgeDate,omitempty" azure:"ro"`
}

// Secret - A secret consisting of a value, id and its attributes.
type Secret struct {
	// The secret management attributes.
	Properties *Properties `json:"attributes,omitempty"`

	// The content type of the secret.
	ContentType *string `json:"contentType,omitempty"`

	// The secret id.
	ID *string `json:"id,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`

	// The secret value.
	Value *string `json:"value,omitempty"`

	// READ-ONLY; If this is a secret backing a KV certificate, then this field specifies the corresponding key backing the KV certificate.
	KeyID *string `json:"kid,omitempty" azure:"ro"`

	// READ-ONLY; True if the secret's lifetime is managed by key vault. If this is a secret backing a certificate, then managed will be true.
	Managed *bool `json:"managed,omitempty" azure:"ro"`
}

// Properties - The secret management properties.
type Properties struct {
	// Determines whether the object is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Expiry date in UTC.
	ExpiresOn *time.Time `json:"exp,omitempty"`

	// Not before date in UTC.
	NotBefore *time.Time `json:"nbf,omitempty"`

	// READ-ONLY; Creation time in UTC.
	CreatedOn *time.Time `json:"created,omitempty" azure:"ro"`

	// READ-ONLY; Last updated time in UTC.
	UpdatedOn *time.Time `json:"updated,omitempty" azure:"ro"`

	// READ-ONLY; softDelete data retention days. Value should be >=7 and <=90 when softDelete enabled, otherwise 0.
	RecoverableDays *int32 `json:"recoverableDays,omitempty" azure:"ro"`

	// READ-ONLY; Reflects the deletion recovery level currently in effect for secrets in the current vault. If it contains 'Purgeable', the secret can be permanently
	// deleted by a privileged user; otherwise, only the
	// system can purge the secret, at the end of the retention interval.
	RecoveryLevel *DeletionRecoveryLevel `json:"recoveryLevel,omitempty" azure:"ro"`
}

func (s *Properties) toGenerated() *internal.SecretAttributes {
	if s == nil {
		return nil
	}
	return &internal.SecretAttributes{
		RecoverableDays: s.RecoverableDays,
		RecoveryLevel:   s.RecoveryLevel.toGenerated(),
		Enabled:         s.Enabled,
		Expires:         s.ExpiresOn,
		NotBefore:       s.NotBefore,
		Created:         s.CreatedOn,
		Updated:         s.UpdatedOn,
	}
}

// create a SecretAttributes object from an internal.SecretAttributes object
func secretPropertiesFromGenerated(i *internal.SecretAttributes) *Properties {
	if i == nil {
		return nil
	}
	return &Properties{
		RecoverableDays: i.RecoverableDays,
		RecoveryLevel:   deletionRecoveryLevelFromGenerated(*i.RecoveryLevel).ToPtr(),
		Enabled:         i.Enabled,
		ExpiresOn:       i.Expires,
		NotBefore:       i.NotBefore,
		CreatedOn:       i.Created,
		UpdatedOn:       i.Updated,
	}
}

// SecretItem contains secret metadata.
type SecretItem struct {
	// The secret management attributes.
	Properties *Properties `json:"attributes,omitempty"`

	// Type of the secret value such as a password.
	ContentType *string `json:"contentType,omitempty"`

	// Secret identifier.
	ID *string `json:"id,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`

	// READ-ONLY; True if the secret's lifetime is managed by key vault. If this is a key backing a certificate, then managed will be true.
	Managed *bool `json:"managed,omitempty" azure:"ro"`
}

// create a SecretItem from the internal.SecretItem model
func secretItemFromGenerated(i *internal.SecretItem) SecretItem {
	if i == nil {
		return SecretItem{}
	}

	return SecretItem{
		Properties:  secretPropertiesFromGenerated(i.Attributes),
		ContentType: i.ContentType,
		ID:          i.ID,
		Tags:        convertPtrMap(i.Tags),
		Managed:     i.Managed,
	}
}

// DeletedSecretItem - The deleted secret item containing metadata about the deleted secret.
type DeletedSecretItem struct {
	// The secret management attributes.
	Properties *Properties `json:"attributes,omitempty"`

	// Type of the secret value such as a password.
	ContentType *string `json:"contentType,omitempty"`

	// Secret identifier.
	ID *string `json:"id,omitempty"`

	// The url of the recovery object, used to identify and recover the deleted secret.
	RecoveryID *string `json:"recoveryId,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`

	// READ-ONLY; The time when the secret was deleted, in UTC
	DeletedDate *time.Time `json:"deletedDate,omitempty" azure:"ro"`

	// READ-ONLY; True if the secret's lifetime is managed by key vault. If this is a key backing a certificate, then managed
	// will be true.
	Managed *bool `json:"managed,omitempty" azure:"ro"`

	// READ-ONLY; The time when the secret is scheduled to be purged, in UTC
	ScheduledPurgeDate *time.Time `json:"scheduledPurgeDate,omitempty" azure:"ro"`
}

func deletedSecretItemFromGenerated(i *internal.DeletedSecretItem) DeletedSecretItem {
	if i == nil {
		return DeletedSecretItem{}
	}

	return DeletedSecretItem{
		Properties:         secretPropertiesFromGenerated(i.Attributes),
		ContentType:        i.ContentType,
		ID:                 i.ID,
		RecoveryID:         i.RecoveryID,
		Tags:               convertPtrMap(i.Tags),
		DeletedDate:        i.DeletedDate,
		Managed:            i.Managed,
		ScheduledPurgeDate: i.ScheduledPurgeDate,
	}
}

func convertPtrMap(m map[string]*string) map[string]string {
	if m == nil {
		return nil
	}

	ret := map[string]string{}
	for key, val := range m {
		ret[key] = *val
	}

	return ret
}

func convertToGeneratedMap(m map[string]string) map[string]*string {
	if m == nil {
		return nil
	}

	ret := map[string]*string{}
	for key, val := range m {
		ret[key] = &val
	}

	return ret
}

//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package convert

import (
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/models"
)

// converts a KeyAttributes to *generated.KeyAttributes
func KeyAttributesToGenerated(k *models.KeyAttributes) *generated.KeyAttributes {
	return &generated.KeyAttributes{
		RecoverableDays: k.RecoverableDays,
		RecoveryLevel:   k.RecoveryLevel,
		Enabled:         k.Enabled,
		Expires:         k.Expires,
		NotBefore:       k.NotBefore,
		Created:         k.Created,
		Updated:         k.Updated,
	}
}

// converts *generated.KeyAttributes to *KeyAttributes
func KeyAttributesFromGenerated(i *generated.KeyAttributes) *models.KeyAttributes {
	if i == nil {
		return &models.KeyAttributes{}
	}

	return &models.KeyAttributes{
		RecoverableDays: i.RecoverableDays,
		RecoveryLevel:   i.RecoveryLevel,
		Attributes: models.Attributes{
			Enabled:   i.Enabled,
			Expires:   i.Expires,
			NotBefore: i.NotBefore,
			Created:   i.Created,
			Updated:   i.Updated,
		},
	}
}

// converts generated.JSONWebKey to publicly exposed version
func JSONWebKeyFromGenerated(i *generated.JSONWebKey) *models.JSONWebKey {
	if i == nil {
		return &models.JSONWebKey{}
	}

	return &models.JSONWebKey{
		Crv:     (*models.JSONWebKeyCurveName)(i.Crv),
		D:       i.D,
		DP:      i.DP,
		DQ:      i.DQ,
		E:       i.E,
		K:       i.K,
		KeyOps:  i.KeyOps,
		ID:      i.Kid,
		KeyType: (*models.KeyType)(i.Kty),
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
func JSONWebKeyToGenerated(j models.JSONWebKey) *generated.JSONWebKey {
	return &generated.JSONWebKey{
		Crv:    (*generated.JSONWebKeyCurveName)(j.Crv),
		D:      j.D,
		DP:     j.DP,
		DQ:     j.DQ,
		E:      j.E,
		K:      j.K,
		KeyOps: j.KeyOps,
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

// convert KeyType to *generated.JSONWebKeyType
func KeyTypeToGenerated(j models.KeyType) *generated.JSONWebKeyType {
	return generated.JSONWebKeyType(j).ToPtr()
}

// convert *generated.KeyItem to *KeyItem
func KeyItemFromGenerated(i *generated.KeyItem) *models.KeyItem {
	if i == nil {
		return nil
	}

	return &models.KeyItem{
		Attributes: KeyAttributesFromGenerated(i.Attributes),
		KID:        i.Kid,
		Tags:       FromGeneratedMap(i.Tags),
		Managed:    i.Managed,
	}
}

// convert *generated.DeletedKeyItem to *DeletedKeyItem
func DeletedKeyItemFromGenerated(i *generated.DeletedKeyItem) *models.DeletedKeyItem {
	if i == nil {
		return nil
	}

	return &models.DeletedKeyItem{
		RecoveryID:         i.RecoveryID,
		DeletedDate:        i.DeletedDate,
		ScheduledPurgeDate: i.ScheduledPurgeDate,
		KeyItem: models.KeyItem{
			Attributes: &models.KeyAttributes{
				Attributes: models.Attributes{
					Enabled:   i.Attributes.Enabled,
					Expires:   i.Attributes.Expires,
					NotBefore: i.Attributes.NotBefore,
					Created:   i.Attributes.Created,
					Updated:   i.Attributes.Updated,
				},
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   i.Attributes.RecoveryLevel,
			},
			KID:     i.Kid,
			Tags:    FromGeneratedMap(i.Tags),
			Managed: i.Managed,
		},
	}
}

func KeyReleasePolicyFromGenerated(i *generated.KeyReleasePolicy) *models.KeyReleasePolicy {
	if i == nil {
		return nil
	}
	return &models.KeyReleasePolicy{
		ContentType: i.ContentType,
		Data:        i.Data,
	}
}

func KeyRotationPolicyAttributesToGenerated(k *models.KeyRotationPolicyAttributes) *generated.KeyRotationPolicyAttributes {
	return &generated.KeyRotationPolicyAttributes{
		ExpiryTime: k.ExpiryTime,
		Created:    k.Created,
		Updated:    k.Updated,
	}
}

func LifetimeActionsToGenerated(l *models.LifetimeActions) *generated.LifetimeActions {
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

func LifetimeActionsFromGenerated(i *generated.LifetimeActions) *models.LifetimeActions {
	if i == nil {
		return nil
	}
	return &models.LifetimeActions{
		Trigger: &models.LifetimeActionsTrigger{
			TimeAfterCreate:  i.Trigger.TimeAfterCreate,
			TimeBeforeExpiry: i.Trigger.TimeBeforeExpiry,
		},
		Action: &models.LifetimeActionsType{
			Type: (*models.ActionType)(i.Action.Type),
		},
	}
}

func ToGeneratedMap(m map[string]string) map[string]*string {
	if m == nil {
		return nil
	}

	ret := make(map[string]*string)
	for k, v := range m {
		ret[k] = &v
	}
	return ret
}

func FromGeneratedMap(m map[string]*string) map[string]string {
	if m == nil {
		return nil
	}

	ret := make(map[string]string)
	for k, v := range m {
		ret[k] = *v
	}
	return ret
}

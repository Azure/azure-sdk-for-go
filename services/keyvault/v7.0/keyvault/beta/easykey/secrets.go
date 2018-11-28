package easykey

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
)

// Secret describes secret information stored in KeyVault.
type Secret struct {
	// ID is the ID information associated with this secret.
	ID ID
	// Value is the value stored by the secret.
	Value string
	// Attr are attributes describing the secret.
	Attr SecretAttr
}

// version implements versioner.
func (s Secret) version() string {
	return s.ID.Version
}

func (s Secret) toBundle() keyvault.SecretBundle {
	id := s.ID.String()
	return keyvault.SecretBundle{
		Value:      &s.Value,
		ID:         &id,
		Attributes: s.Attr.toBundle(),
	}
}

// SecretAttr are attributes associated with this secret.
type SecretAttr struct {
	// RecoveryLevel the level of recovery for this password when deleted.  See the description of
	// DeletionRecoveryLevel above.
	RecoveryLevel DeletionRecoveryLevel
	// Enabled indicates if the secret is currently enabled.
	Enabled bool
	// Created indicates the time the secret was created in UTC. If set to the zero value, it indicates
	// this was not set.
	Created time.Time
	// NotBefore indicate that the key isn't valid before this time in UTC. If set to the zero value, it indicates
	// this was not set.
	NotBefore time.Time
	// Updated indicates the last time the secret was updated in UTC. If set to the zero value, it indicates
	// this was not set.
	Updated time.Time
}

// fromSecretAttributes converts from the internal secrets to this more garbage collector friendly version.
// It drops all the response data that isn't useful.
func (s *SecretAttr) fromSecretAttributes(sa keyvault.SecretAttributes) {
	s.RecoveryLevel = sa.RecoveryLevel
	s.Enabled = *sa.Enabled
	if sa.NotBefore != nil {
		s.NotBefore = time.Time(*sa.NotBefore)
	}
	if sa.Created != nil {
		s.Created = time.Time(*sa.Created)
	}
	if sa.Updated != nil {
		s.Updated = time.Time(*sa.Updated)
	}
}

func (s SecretAttr) toBundle() *keyvault.SecretAttributes {
	sa := &keyvault.SecretAttributes{
		RecoveryLevel: s.RecoveryLevel,
		Enabled:       &s.Enabled,
	}

	if s.NotBefore.IsZero() {
		z := date.NewUnixTimeFromSeconds(0.0)
		sa.NotBefore = &z
	} else {
		d := date.UnixTime(s.NotBefore)
		sa.NotBefore = &d
	}
	if s.Created.IsZero() {
		z := date.NewUnixTimeFromSeconds(0.0)
		sa.Created = &z
	} else {
		d := date.UnixTime(s.Created)
		sa.Created = &d
	}
	if s.Updated.IsZero() {
		z := date.NewUnixTimeFromSeconds(0.0)
		sa.Updated = &z
	} else {
		d := date.UnixTime(s.Updated)
		sa.Updated = &d
	}
	return sa
}

// Secret returns a secret from the vault that is defined by key at version.
// version can be set to constant LatestVersion to retreive the current version of the key.
func (k *KeyVault) Secret(ctx context.Context, key, version string) (Secret, error) {
	secret, err := k.vault.GetSecret(ctx, k.vaultURL, key, version)
	if err != nil {
		if strings.Contains(err.Error(), "StatusCode=403") {
			return Secret{}, fmt.Errorf("access denied to keyvault(%s) for key(%s), often due to not having an access policy set for the user in AZURE_CLIENT_ID or missing environmental variable AZURE_CLIENT_SECRET: %s", k.vaultName, key, err)
		}
		return Secret{}, fmt.Errorf("could not connect to vault %s and retrieve key %q: %s", k.vaultURL, key, err)
	}

	if secret.Value == nil || *secret.Value == "" {
		return Secret{}, fmt.Errorf("vault %q with secretKey %q contained an empty string or was nil", k.vaultName, key)
	}

	attr := SecretAttr{}
	attr.fromSecretAttributes(*secret.Attributes)

	u, _ := url.Parse(*secret.ID)
	id, err := urlToID(u)
	if err != nil {
		return Secret{}, err
	}
	s := Secret{
		ID:    id,
		Value: *secret.Value,
		Attr:  attr,
	}

	return s, nil
}

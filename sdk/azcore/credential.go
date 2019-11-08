// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"encoding/json"
	"time"
)

// Credential represents any credential type.
type Credential interface {
	// Policy returns a policy that wil request the underlying credential type with
	// the specified set of scopes (if applicable) and apply it to the HTTP request.
	Policy(scopes ...string) Policy
}

// CredentialFunc is a type that implements the Credential interface.
// Use this type when implementing a stateless credential as a first-class function.
type CredentialFunc func(scopes ...string) Policy

// Policy implements the Credential interface on CredentialFunc.
func (cf CredentialFunc) Policy(scopes ...string) Policy {
	return cf(scopes...)
}

// TokenCredential represents a credential capable of providing an OAuth token.
type TokenCredential interface {
	// GetToken requests an access token for the specified set of scopes.
	GetToken(ctx context.Context, scopes []string) (*AccessToken, error)
}

// AccessToken represents an Azure service bearer access token with expiry information.
type AccessToken struct {
	Token     string      `json:"access_token"`
	ExpiresIn json.Number `json:"expires_in"`
	ExpiresOn time.Time
}

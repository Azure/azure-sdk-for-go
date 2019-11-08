// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"encoding/json"
	"time"
)

// CredentialPolicyOptions contains various options used to create a credential policy.
type CredentialPolicyOptions struct {
	// Scopes is the list of OAuth2 authentication scopes used when requesting a token.
	// This field is ignored for other forms of authentication (e.g. shared key).
	Scopes []string
}

// Credential represents any credential type.
type Credential interface {
	// Policy returns a policy that requests the credential and applies it to the HTTP request.
	Policy(options CredentialPolicyOptions) Policy
}

// CredentialFunc is a type that implements the Credential interface.
// Use this type when implementing a stateless credential as a first-class function.
type CredentialFunc func(options CredentialPolicyOptions) Policy

// Policy implements the Credential interface on CredentialFunc.
func (cf CredentialFunc) Policy(options CredentialPolicyOptions) Policy {
	return cf(options)
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

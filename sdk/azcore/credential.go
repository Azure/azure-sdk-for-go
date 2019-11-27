// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"encoding/json"
	"net/url"
	"time"
)

// AuthenticationPolicyOptions contains various options used to create a credential policy.
type AuthenticationPolicyOptions struct {
	// Scopes is the list of OAuth2 authentication scopes used when requesting a token.
	// This field is ignored for other forms of authentication (e.g. shared key).
	Scopes        []string
	AuthorityHost url.URL
}

// Credential represents any credential type.
type Credential interface {
	// AuthenticationPolicy returns a policy that requests the credential and applies it to the HTTP request.
	AuthenticationPolicy(options AuthenticationPolicyOptions) Policy
}

// credentialFunc is a type that implements the Credential interface.
// Use this type when implementing a stateless credential as a first-class function.
type credentialFunc func(options AuthenticationPolicyOptions) Policy

// AuthenticationPolicy implements the Credential interface on credentialFunc.
func (cf credentialFunc) AuthenticationPolicy(options AuthenticationPolicyOptions) Policy {
	return cf(options)
}

// TokenCredential represents a credential capable of providing an OAuth token.
type TokenCredential interface {
	Credential
	// GetToken requests an access token for the specified set of scopes.
	GetToken(ctx context.Context, options TokenRequestOptions) (*AccessToken, error)
}

// AccessToken represents an Azure service bearer access token with expiry information.
type AccessToken struct {
	Token     string      `json:"access_token"`
	ExpiresIn json.Number `json:"expires_in"`
	ExpiresOn time.Time
}

// TokenRequestOptions contain specific parameter that may be used by credentials types when attempting to get a token
type TokenRequestOptions struct {
	Scopes []string
}

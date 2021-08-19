//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"time"
)

// AuthenticationOptions contains various options used to create a credential policy.
type AuthenticationOptions struct {
	// TokenRequest is a TokenRequestOptions that includes a scopes field which contains
	// the list of OAuth2 authentication scopes used when requesting a token.
	// This field is ignored for other forms of authentication (e.g. shared key).
	TokenRequest TokenRequestOptions
	// AuxiliaryTenants contains a list of additional tenant IDs to be used to authenticate
	// in cross-tenant applications.
	AuxiliaryTenants []string
}

// Credential represents any credential type.
type Credential interface {
	// AuthenticationPolicy returns a policy that requests the credential and applies it to the HTTP request.
	NewAuthenticationPolicy(options AuthenticationOptions) Policy
}

// credentialFunc is a type that implements the Credential interface.
// Use this type when implementing a stateless credential as a first-class function.
type credentialFunc func(options AuthenticationOptions) Policy

// AuthenticationPolicy implements the Credential interface on credentialFunc.
func (cf credentialFunc) NewAuthenticationPolicy(options AuthenticationOptions) Policy {
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
	Token     string
	ExpiresOn time.Time
}

// TokenRequestOptions contain specific parameter that may be used by credentials types when attempting to get a token.
type TokenRequestOptions struct {
	// Scopes contains the list of permission scopes required for the token.
	Scopes []string
	// TenantID contains the tenant ID to use in a multi-tenant authentication scenario, if TenantID is set
	// it will override the tenant ID that was added at credential creation time.
	TenantID string
}

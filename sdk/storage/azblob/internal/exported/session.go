// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"net/http"
	"time"
)

// SessionMode specifies how session-based authentication is handled.
type SessionMode string

const (
	// SessionModeDefault is the default mode where sessions are disabled.
	SessionModeDefault SessionMode = ""
	// SessionModeDisabled explicitly disables session-based authentication.
	SessionModeDisabled SessionMode = "disabled"
	// SessionModeEnabled enables session-based authentication.
	SessionModeEnabled SessionMode = "enabled"
)

// PossibleSessionModeValues returns a slice of possible values for SessionMode.
func PossibleSessionModeValues() []SessionMode {
	return []SessionMode{
		SessionModeDefault,
		SessionModeDisabled,
		SessionModeEnabled,
	}
}

type SessionCredential struct {
	token  string
	key    string
	expiry time.Time
	// fallback indicates that session creation failed and the caller should use bearer token
	// authentication instead. This is stored as a field rather than returned as an error because
	// temporal.Resource only caches successful (non-error) results. Returning a non-error fallback
	// value allows the decision to be cached for the duration of expiry, avoiding repeated
	// session creation attempts when the service indicates the feature is unavailable.
	fallback bool
}

// NewSessionCredential creates a SessionCredential with the given token, key, and expiry.
func NewSessionCredential(token, key string, expiry time.Time) SessionCredential {
	return SessionCredential{token: token, key: key, expiry: expiry}
}

// NewSessionCredentialFallback creates a SessionCredential that indicates fallback to bearer token auth.
func NewSessionCredentialFallback(expiry time.Time) SessionCredential {
	return SessionCredential{fallback: true, expiry: expiry}
}

// Token returns the session token.
func (s SessionCredential) Token() string { return s.token }

// Key returns the session key.
func (s SessionCredential) Key() string { return s.key }

// Expiry returns the session expiry time.
func (s SessionCredential) Expiry() time.Time { return s.expiry }

// Fallback returns true if the caller should fall back to bearer token authentication.
func (s SessionCredential) Fallback() bool { return s.fallback }

type SessionContext struct {
	ContainerName string
}

type SessionProvider interface {
	GetSession(ctx context.Context, sessionCtx SessionContext) (SessionCredential, error)
	InvalidateSession(sessionCtx SessionContext, current SessionCredential) (err error)
	IsRequestEligible(req *http.Request) bool
}

// SessionOptions configures session-based authentication behavior.
type SessionOptions struct {
	// Mode specifies the session authentication mode.
	Mode SessionMode

	// AccountName is the optional storage account name.
	AccountName string

	// SessionProvider is the optional session cache. If nil, a default client-scoped cache is used.
	SessionProvider SessionProvider
}

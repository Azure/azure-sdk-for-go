// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
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
	// fallback indicates that session creation is unavailable and the caller should use bearer
	// token authentication instead. This is represented as a successful (non-error) value so the
	// decision can be cached for the duration of expiry, avoiding repeated session creation
	// attempts while the service indicates the feature is unavailable.
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

// SessionProvider supplies and caches session credentials for outgoing requests.
// Implementations must be safe for concurrent use.
type SessionProvider interface {
	// GetSession returns a session credential for the given request, acquiring or refreshing
	// one as needed. Implementations derive the scope (e.g. the container name) from the
	// request URL. A credential whose Fallback reports true indicates the caller should use
	// bearer token authentication instead.
	GetSession(req *http.Request) (SessionCredential, error)

	// InvalidateSession discards the cached session for the given request's scope, but only if
	// the cached credential is still the one described by current. If it has already been
	// refreshed, the call is a no-op. A new session is acquired on the next call to GetSession.
	InvalidateSession(req *http.Request, current SessionCredential) (err error)

	// IsRequestEligible reports whether the request can be authenticated with a session.
	IsRequestEligible(req *http.Request) bool
}

// SessionOptions configures session-based authentication behavior.
type SessionOptions struct {
	// Mode specifies the session authentication mode.
	Mode SessionMode

	// AccountName is the optional storage account name.
	AccountName string

	// Provider is the optional session cache. If nil, a default client-scoped cache is used.
	Provider SessionProvider
}

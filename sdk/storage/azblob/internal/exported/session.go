// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/internal/autorefresh"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/internal/locality"
)

// SessionMode specifies whether eligible requests are authenticated with a session.
type SessionMode string

const (
	// SessionModeAuto is the zero value, and therefore the default when no value is specified.
	// The client library decides whether sessions are used, and that decision may change in a
	// future release. Currently, SessionModeAuto resolves to SessionModeDisabled.
	SessionModeAuto SessionMode = ""

	// SessionModeDisabled authenticates every request with a bearer token.
	SessionModeDisabled SessionMode = "Disabled"

	// SessionModeEnabled authenticates eligible requests with a session, created and cached per
	// container. It requires the storage account name: when it can be determined from neither
	// SessionOptions.AccountName nor the client's URL, client construction fails.
	SessionModeEnabled SessionMode = "Enabled"
)

// PossibleSessionModeValues returns the possible values for SessionMode.
func PossibleSessionModeValues() []SessionMode {
	return []SessionMode{
		SessionModeAuto,
		SessionModeDisabled,
		SessionModeEnabled,
	}
}

// ResolveSessionMode resolves SessionModeAuto to the mode it currently stands for.
func ResolveSessionMode(mode SessionMode) SessionMode {
	if mode == SessionModeAuto {
		return SessionModeDisabled
	}
	return mode
}

// SessionOptions configures session authentication. Sessions apply only to clients
// authenticated with an azcore.TokenCredential; clients using a shared key, a SAS or no
// credential ignore these options.
type SessionOptions struct {
	// Mode specifies whether sessions are used. The default, SessionModeAuto, currently resolves
	// to SessionModeDisabled.
	Mode SessionMode

	// AccountName is the storage account name used to sign requests with the session key. When
	// empty, it is derived from the client's URL. Set it when the URL is a custom domain from
	// which the account name can't be derived.
	// When the account name can be determined from neither this field nor the URL, client
	// construction fails if Mode is SessionModeEnabled; otherwise the client authenticates every
	// request with a bearer token.
	AccountName string

	// Provider is an optional session provider that owns the session cache, created with
	// NewContainerSessionProvider. Clients configured with the same provider share its cached
	// sessions, so a session is created once per container no matter how many clients use it,
	// and the cache outlives any one of those clients. The provider must target the same blob
	// service endpoint as the clients it is used with.
	// When nil, the client creates a provider of its own, shared only with the clients derived
	// from it (for example the container and blob clients created from a service client).
	Provider SessionProvider
}

// SessionProvider provides and caches the sessions used to authenticate eligible requests. It
// can be shared across clients through SessionOptions.Provider.
//
// SessionProvider can't be implemented outside this module. Use NewContainerSessionProvider to
// create one.
type SessionProvider interface {
	// getSession returns the cached session for the request's scope, acquiring one if needed.
	// The result may be a fallback marker meaning the request should use a bearer token.
	getSession(req *http.Request) (sessionCredential, error)

	// invalidateSession discards the cached session for the request's scope, but only when it is
	// still used. The next eligible request acquires a new session.
	invalidateSession(req *http.Request, used sessionCredential)

	// isRequestEligible reports whether the request can be authenticated with a session.
	isRequestEligible(req *http.Request) bool
}

// sessionCredential is a session returned by Create Session, or a marker that sessions are
// unavailable and requests should fall back to bearer token authentication.
type sessionCredential struct {
	token string
	key   string
	// fallback marks a cached "sessions are unavailable" decision. It is cached like a session
	// so that the service isn't asked again until it expires.
	fallback bool
}

const (
	// sessionRefreshBuffer is how long before a session expires that a background refresh starts.
	sessionRefreshBuffer = 30 * time.Second

	// sessionBackgroundAcquireTimeout bounds a background session refresh.
	sessionBackgroundAcquireTimeout = 30 * time.Second

	// sessionFallbackCooldown is how long a fallback-to-bearer decision is cached after Create
	// Session fails with a 5xx, a 403, or a 400 FeatureNotEnabled.
	sessionFallbackCooldown = 5 * time.Minute

	featureNotEnabled = "FeatureNotEnabled"
)

// ContainerSessionProvider is a SessionProvider that creates sessions with an
// azcore.TokenCredential and caches one per container. It applies only to clients authenticated
// with a token credential.
//
// Use NewContainerSessionProvider to create one. The zero value is not usable.
type ContainerSessionProvider struct {
	client *generated.ServiceClient
	now    func() time.Time

	// caches maps a container name to the cache holding its session. Entries are never removed,
	// so a provider retains one entry per distinct container it has been used with.
	caches sync.Map // map[string]*autorefresh.Cache[sessionCredential]
}

// NewContainerSessionProviderFromClient returns a ContainerSessionProvider that creates sessions
// through client, which must target a blob service endpoint and must not itself use sessions.
// It is the constructor the public NewContainerSessionProvider and the clients use once they
// have built a session-free pipeline.
func NewContainerSessionProviderFromClient(client *generated.ServiceClient) *ContainerSessionProvider {
	return &ContainerSessionProvider{client: client, now: time.Now}
}

func (p *ContainerSessionProvider) getSession(req *http.Request) (sessionCredential, error) {
	cache, err := p.cacheFor(req)
	if err != nil {
		return sessionCredential{}, err
	}
	// The request may be routed to a layout endpoint, but its session must be created at the
	// account endpoint.
	return cache.Get(locality.WithoutEndpoint(req.Context()))
}

func (p *ContainerSessionProvider) invalidateSession(req *http.Request, used sessionCredential) {
	cache, err := p.cacheFor(req)
	if err != nil {
		return
	}
	cache.InvalidateIf(func(current sessionCredential) bool {
		return current.token == used.token
	})
}

// isRequestEligible reports whether req is a Get Blob request, the only operation a session
// authenticates: a GET addressed to a blob, with neither a comp nor a restype query parameter
// (whatever their values), and not requesting a structured message body. Every other request,
// including Get Blob Properties, Put Blob, Put Block, Put Block List, Get Layout and every
// container- or service-level operation, is authenticated with the bearer token.
func (p *ContainerSessionProvider) isRequestEligible(req *http.Request) bool {
	if req == nil || req.URL == nil || req.Method != http.MethodGet {
		return false
	}
	if _, blob, err := shared.GetContainerAndBlobName(req.URL); err != nil || blob == "" {
		return false
	}
	query := req.URL.Query()
	if query.Has("comp") || query.Has("restype") {
		return false
	}
	if shared.HeaderValue(req.Header, shared.HeaderXmsStructuredBody) != "" {
		return false
	}
	return true
}

// cacheFor returns the session cache for the request's container, creating it on first use.
func (p *ContainerSessionProvider) cacheFor(req *http.Request) (*autorefresh.Cache[sessionCredential], error) {
	if req == nil {
		return nil, errors.New("a request is required to determine the session container")
	}
	container, _, err := shared.GetContainerAndBlobName(req.URL)
	if err != nil {
		return nil, err
	}
	// container names are case-insensitive to the service
	key := strings.ToLower(container)
	if cache, ok := p.caches.Load(key); ok {
		return cache.(*autorefresh.Cache[sessionCredential]), nil
	}
	// Creating a cache has no side effects, so losing the race to store one only discards an
	// unused cache; it never causes an extra Create Session.
	cache := autorefresh.New(p.acquireFunc(container), &autorefresh.Options{
		BackgroundAcquireTimeout: sessionBackgroundAcquireTimeout,
		Now:                      p.now,
	})
	actual, _ := p.caches.LoadOrStore(key, cache)
	return actual.(*autorefresh.Cache[sessionCredential]), nil
}

// acquireFunc returns the function that creates a session for container. A Create Session
// failure that indicates sessions are unavailable (any 5xx, a 403, or a 400 with error code
// FeatureNotEnabled) becomes a fallback marker cached for five minutes, during which this
// container's requests use bearer tokens without contacting the service again. Other failures
// are returned to the caller.
func (p *ContainerSessionProvider) acquireFunc(container string) autorefresh.AcquireFunc[sessionCredential] {
	containerClient := generated.NewContainerClient(runtime.JoinPaths(p.client.Endpoint(), container), p.client.InternalClient())
	return func(ctx context.Context) (autorefresh.Entry[sessionCredential], error) {
		resp, err := containerClient.CreateSession(ctx, generated.CreateSessionConfiguration{AuthenticationType: to.Ptr(generated.AuthenticationTypeHmac)}, nil)
		if err != nil {
			if isSessionFallbackError(err) {
				expires := p.now().Add(sessionFallbackCooldown)
				// no refresh before expiry: the decision holds for the whole cooldown
				return autorefresh.Entry[sessionCredential]{Value: sessionCredential{fallback: true}, ExpiresOn: expires, RefreshOn: expires}, nil
			}
			return autorefresh.Entry[sessionCredential]{}, err
		}
		if resp.Expiration == nil || resp.Credentials == nil || resp.Credentials.SessionToken == nil || resp.Credentials.SessionKey == nil {
			return autorefresh.Entry[sessionCredential]{}, errors.New("the Create Session response is missing the session credentials or their expiration")
		}
		return autorefresh.Entry[sessionCredential]{
			Value:     sessionCredential{token: *resp.Credentials.SessionToken, key: *resp.Credentials.SessionKey},
			ExpiresOn: *resp.Expiration,
			RefreshOn: resp.Expiration.Add(-sessionRefreshBuffer),
		}, nil
	}
}

// isSessionFallbackError reports whether a Create Session failure means sessions are
// unavailable rather than that the request was wrong.
func isSessionFallbackError(err error) bool {
	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return false
	}
	return respErr.StatusCode >= http.StatusInternalServerError ||
		respErr.StatusCode == http.StatusForbidden ||
		(respErr.StatusCode == http.StatusBadRequest && strings.EqualFold(respErr.ErrorCode, featureNotEnabled))
}

// ValidateSessionMode returns an error for a mode that isn't one of the SessionMode constants.
func ValidateSessionMode(mode SessionMode) error {
	switch mode {
	case SessionModeAuto, SessionModeDisabled, SessionModeEnabled:
		return nil
	}
	return fmt.Errorf("unsupported session mode %q", mode)
}

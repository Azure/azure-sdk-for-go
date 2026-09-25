// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

const featureNotEnabled = "FeatureNotEnabled"

// cooldown duration applied when the service indicates sessions are unavailable.
const transientFailureCooldown = 5 * time.Minute

// containerSessionProvider implements SessionProvider for container-scoped token credential sessions.
type containerSessionProvider struct {
	genClient *generated.ServiceClient

	// sessions maps container names to their cached session credential.
	sessions sync.Map // map[string]*temporal.Resource[exported.SessionCredential, context.Context]

	// invalidateMu serializes invalidation so that the "is this still the cached session" check
	// and the expiry that follows it can't interleave with another invalidation. Without it,
	// every request that was in flight with a rejected session could discard the replacement
	// acquired by the previous one. Invalidation is a rare, off the happy path operation, so a
	// single lock for all containers is cheap.
	invalidateMu sync.Mutex
}

// NewContainerSessionProvider creates a SessionProvider that manages container-scoped sessions
// using the provided token credential.
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - storageURL - the URL of the storage account e.g. https://<account>.blob.core.windows.net/
//   - options - client options; pass nil to accept the default values
//
// The provider caches one session per container for its lifetime, so its memory use grows with
// the number of distinct containers it's used with. Sharing a single provider across an
// application that accesses a large or unbounded set of containers will retain an entry per
// container.
func NewContainerSessionProvider(cred azcore.TokenCredential, storageURL string, options *ClientOptions) (exported.SessionProvider, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	svcURL, err := shared.GetServiceURL(storageURL)
	if err != nil {
		return nil, err
	}

	// Create a service client with sessions disabled to avoid recursive session creation.
	svcOpts := *options
	svcOpts.Session.Mode = exported.SessionModeDisabled
	azClient, err := GetAzClient(svcURL, cred, nil, &svcOpts)
	if err != nil {
		return nil, err
	}

	return &containerSessionProvider{
		genClient: generated.NewServiceClient(svcURL, azClient),
	}, nil
}

// GetSession returns a cached session credential for the request's container, refreshing it if expired.
func (p *containerSessionProvider) GetSession(req *http.Request) (exported.SessionCredential, error) {
	resource, err := p.resourceForRequest(req)
	if err != nil {
		return exported.SessionCredential{}, err
	}
	return resource.Get(req.Context())
}

// InvalidateSession discards the cached session for the request's container so that a new one is
// acquired on the next call to GetSession.
//
// The cached session is only discarded when it is still the one the caller was rejected with, so
// that the many requests that were in flight with a rejected session cause a single replacement
// rather than one each.
func (p *containerSessionProvider) InvalidateSession(req *http.Request, reqCred exported.SessionCredential) error {
	resource, err := p.resourceForRequest(req)
	if err != nil {
		return err
	}

	// the check and the expiry that follows it must be atomic with respect to other invalidations
	p.invalidateMu.Lock()
	defer p.invalidateMu.Unlock()

	// A caller reporting a session that is no longer cached has nothing to discard: the session
	// was already replaced, either by an eager refresh or by an earlier invalidation.
	currCred, _ := resource.Get(req.Context())
	if currCred.Token() == reqCred.Token() {
		resource.Expire()
	}
	return nil
}

// Query parameter values identifying the blob operations a session can authenticate.
const (
	compBlock     = "block"
	compBlockList = "blocklist"
)

// IsRequestEligible returns true if the request is eligible for session-based authentication.
//
// Eligibility is deliberately narrow: a request must be blob-level and must match one of the
// operations the service accepts a session for. Everything else, including every other verb and
// every other comp value, continues to authenticate with the bearer token.
//
// The eligible operations are Get Blob, Get Blob Properties, Put Blob, Put Block and
// Put Block List. Each is identified by the request shape the generated clients produce rather
// than by the verb alone, so that operations sharing a verb, such as Copy Blob or
// Put Block From URL, are not swept in.
func (p *containerSessionProvider) IsRequestEligible(req *http.Request) bool {
	if req == nil || req.URL == nil {
		return false
	}

	// Session auth is not supported for structured message requests
	if shared.HeaderValue(req.Header, shared.HeaderXmsStructuredBody) != "" {
		return false
	}

	query := req.URL.Query()

	// restype scopes a request to the container or the account, neither of which is a
	// blob-level operation a container session covers.
	if query.Get("restype") != "" {
		return false
	}

	// A session is scoped to a container, and only blob-level requests are eligible.
	if _, blob, err := shared.GetContainerAndBlobName(req.URL); err != nil || blob == "" {
		return false
	}

	comp := query.Get("comp")

	switch req.Method {
	case http.MethodGet:
		// Get Blob is the comp-less GET. Every other GET is a sub-resource read
		// (comp=tags, comp=blocklist, comp=layout, ...) and stays on the bearer token.
		return comp == ""

	case http.MethodHead:
		// Get Blob Properties is the comp-less HEAD. It is currently the only HEAD the blob
		// API defines, but the comp check keeps a later one from becoming eligible silently.
		return comp == ""

	case http.MethodPut:
		// The *FromURL operations (Copy Blob, Put Blob From URL, Put Block From URL) share a
		// URL shape with the operations below but are distinct REST operations that are not
		// approved for session auth. The service tells them apart by the copy-source header,
		// so that is what excludes them here.
		if shared.HeaderValue(req.Header, shared.HeaderXmsCopySource) != "" {
			return false
		}
		switch comp {
		case compBlock:
			// Put Block. The generated client also sets blockid, but comp=block is what
			// identifies the operation.
			return true
		case compBlockList:
			// Put Block List.
			return true
		case "":
			// Put Blob. It is the only comp-less PUT that creates blob content, and the
			// generated clients identify it with x-ms-blob-type; requiring that header keeps
			// a comp-less PUT added later from becoming eligible by default.
			return shared.HeaderValue(req.Header, shared.HeaderXmsBlobType) != ""
		}
	}

	return false
}

// resourceForRequest resolves the container for the request and returns its session resource.
func (p *containerSessionProvider) resourceForRequest(req *http.Request) (*temporal.Resource[exported.SessionCredential, context.Context], error) {
	if req == nil {
		return nil, errors.New("request URL is required to determine the session container")
	}
	containerName, _, err := shared.GetContainerAndBlobName(req.URL)
	if err != nil {
		return nil, err
	}
	return p.getOrCreateResource(containerName), nil
}

// getOrCreateResource returns the temporal.Resource for the given container, creating one if needed.
func (p *containerSessionProvider) getOrCreateResource(containerName string) *temporal.Resource[exported.SessionCredential, context.Context] {
	if v, ok := p.sessions.Load(containerName); ok {
		return v.(*temporal.Resource[exported.SessionCredential, context.Context])
	}

	containerClient := p.getContainerClient(containerName)
	resource := temporal.NewResourceWithOptions(acquireSession(containerClient), temporal.ResourceOptions[exported.SessionCredential, context.Context]{
		ShouldRefresh: shouldRefreshSession,
	})
	actual, _ := p.sessions.LoadOrStore(containerName, resource)
	return actual.(*temporal.Resource[exported.SessionCredential, context.Context])
}

func (p *containerSessionProvider) getContainerClient(containerName string) *generated.ContainerClient {
	containerURL := runtime.JoinPaths(p.genClient.Endpoint(), containerName)
	return generated.NewContainerClient(containerURL, p.genClient.InternalClient())
}

// acquireSession is the function called by temporal.Resource to create a new session.
// When the service indicates that session creation is unavailable, a fallback credential is
// returned so the decision is cached for the duration of its expiry rather than retried per request.
func acquireSession(client *generated.ContainerClient) func(context.Context) (exported.SessionCredential, time.Time, error) {
	return func(ctx context.Context) (creds exported.SessionCredential, expiry time.Time, err error) {
		resp, err := client.CreateSession(ctx, generated.CreateSessionConfiguration{AuthenticationType: to.Ptr(generated.AuthenticationTypeHMAC)}, nil)
		// Fall back to using bearer token if session is unable to be created
		if err != nil {
			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				switch {
				case respErr.StatusCode >= 500,
					respErr.StatusCode == http.StatusBadRequest && respErr.ErrorCode == featureNotEnabled,
					respErr.StatusCode == http.StatusForbidden:
					errorExpiry := time.Now().Add(transientFailureCooldown)
					return exported.NewSessionCredentialFallback(errorExpiry), errorExpiry, nil
				}
			}
			return creds, expiry, err
		}

		if resp.Expiration != nil {
			expiry = *resp.Expiration
		}
		var token, key string
		if resp.Credentials != nil {
			if resp.Credentials.SessionToken != nil {
				token = *resp.Credentials.SessionToken
			}
			if resp.Credentials.SessionKey != nil {
				key = *resp.Credentials.SessionKey
			}
		}

		return exported.NewSessionCredential(token, key, expiry), expiry, nil
	}
}

func shouldRefreshSession(resource exported.SessionCredential, _ context.Context) bool {
	if resource.Fallback() {
		// A fallback credential is a cached "sessions are unavailable" decision. Refreshing it
		// early would contact the service before the decision was due to be reconsidered; let it
		// expire instead.
		return false
	}
	return resource.Expiry().Add(-30 * time.Second).Before(time.Now())
}

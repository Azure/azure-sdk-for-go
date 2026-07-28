// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

const featureNotEnabled = "FeatureNotEnabled"

// containerSessionProvider implements SessionProvider for container-scoped token credential sessions.
type containerSessionProvider struct {
	cred      azcore.TokenCredential
	svcURL    string
	options   *ClientOptions
	genClient *generated.ServiceClient

	// sessions maps container names to their temporal.Resource for session credentials.
	sessions sync.Map // map[string]*temporal.Resource[exported.SessionCredential, context.Context]
}

// NewContainerSessionProvider creates a SessionProvider that manages container-scoped sessions
// using the provided token credential.
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - storageURL - the URL of the storage account e.g. https://<account>.blob.core.windows.net/
//   - options - client options; pass nil to accept the default values
func NewContainerSessionProvider(cred azcore.TokenCredential, storageURL string, options *ClientOptions) (SessionProvider, error) {
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
	azClient, err := base.GetAzClient(cred, nil, (*base.ClientOptions)(&svcOpts))
	if err != nil {
		return nil, err
	}

	return &containerSessionProvider{
		cred:      cred,
		svcURL:    svcURL,
		options:   options,
		genClient: generated.NewServiceClient(svcURL, azClient),
	}, nil
}

// GetSession returns a cached session credential for the given container, refreshing it if expired.
func (p *containerSessionProvider) GetSession(ctx context.Context, sessionCtx SessionContext) (SessionCredential, error) {
	resource := p.getOrCreateResource(sessionCtx.ContainerName)
	return resource.Get(ctx)
}

// InvalidateSession marks the session for the given container as expired, forcing a refresh on the next call.
// It only invalidates if the stored session matches current to avoid invalidating an already-refreshed session.
func (p *containerSessionProvider) InvalidateSession(sessionCtx SessionContext, current SessionCredential) error {
	if v, ok := p.sessions.Load(sessionCtx.ContainerName); ok {
		resource := v.(*temporal.Resource[exported.SessionCredential, context.Context])
		stored, err := resource.Get(context.Background())
		if err == nil && stored.Token() == current.Token() {
			resource.Expire()
		}
	}
	return nil
}

// IsRequestEligible returns true if the request is eligible for session-based authentication.
func (p *containerSessionProvider) IsRequestEligible(req *http.Request) bool {
	if req == nil || req.Method != http.MethodGet {
		return false
	}

	u := req.URL
	if u == nil {
		return false
	}

	// Session auth is not supported for requests with comp query parameter
	if u.Query().Get("comp") != "" {
		return false
	}

	// Path format: /<container>/<blob>
	path := strings.TrimPrefix(u.Path, "/")
	if path == "" {
		return false
	}

	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return false
	}

	return true
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

// acquireSession is the function called by temporal.Resource to create a new session.
func acquireSession(client *generated.ContainerClient) func(context.Context) (exported.SessionCredential, time.Time, error) {
	return func(ctx context.Context) (creds exported.SessionCredential, expiry time.Time, err error) {
		resp, err := client.CreateSession(ctx, generated.CreateSessionConfiguration{AuthenticationType: to.Ptr(generated.AuthenticationTypeHMAC)}, nil)
		// Fall back to using bearer token if session is unable to be created
		if err != nil {
			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				if respErr.StatusCode >= 500 {
					errorExpiry := time.Now().Add(1 * time.Minute)
					return exported.NewSessionCredentialFallback(errorExpiry), errorExpiry, nil
				}
				errorExpiry := time.Now().Add(24 * time.Hour)
				if respErr.StatusCode == http.StatusBadRequest && respErr.ErrorCode == featureNotEnabled {
					return exported.NewSessionCredentialFallback(errorExpiry), errorExpiry, nil
				}
				if respErr.StatusCode == http.StatusForbidden {
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
	return resource.Expiry().Add(-30 * time.Second).Before(time.Now())
}

func (p *containerSessionProvider) getContainerClient(containerName string) *generated.ContainerClient {
	containerURL := runtime.JoinPaths(p.genClient.Endpoint(), containerName)
	return generated.NewContainerClient(containerURL, p.genClient.InternalClient())
}

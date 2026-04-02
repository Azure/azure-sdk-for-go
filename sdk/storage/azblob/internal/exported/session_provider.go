// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// ErrFallbackToBearer indicates that the container does not support sessions
// and the caller should fall back to bearer token authentication.
var ErrFallbackToBearer = errors.New("container does not support sessions, fallback to bearer token")

type SessionState struct {
	client *generated.ContainerClient
	ctx    context.Context
}

// SingleContainerProvider caches a session for a single container using a temporal resource.
// It is safe for concurrent use.
type SingleContainerProvider struct {
	client *generated.ContainerClient
	mu     sync.RWMutex
	// mutex protects access to containerName and resource to ensure thread safety
	// when we extend to support multiple containers in the future, we can change this to a map of containerName to resource
	containerName string
	resource      *temporal.Resource[SessionCredentials, SessionState]
}

// NewSingleContainerProvider creates a new SingleContainerProvider instance with the specified client.
func NewSingleContainerProvider(client *generated.ServiceClient, containerName string) *SingleContainerProvider {
	containerURL := runtime.JoinPaths(client.Endpoint(), containerName)
	cc := generated.NewContainerClient(containerURL, client.InternalClient())

	return &SingleContainerProvider{
		client:        cc,
		containerName: containerName,
		resource:      temporal.NewResource(acquire),
	}
}

// acquire is the function called by temporal.Resource to create a new session.
func acquire(state SessionState) (creds generated.SessionCredentials, expiry time.Time, err error) {
	resp, err := state.client.CreateSession(state.ctx, generated.CreateSessionConfiguration{AuthenticationType: to.Ptr(generated.AuthenticationTypeHMAC)}, nil)
	// Fall back to using bearer token if session is unable to be created
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			if respErr.StatusCode >= 500 {
				return creds, expiry, ErrFallbackToBearer
			}
			if respErr.StatusCode == http.StatusBadRequest && respErr.ErrorCode == "FeatureNotEnabled" {
				return creds, expiry, ErrFallbackToBearer
			}
			if respErr.StatusCode == http.StatusForbidden {
				return creds, expiry, ErrFallbackToBearer
			}
		}
		return creds, expiry, err
	}

	if resp.Expiration != nil {
		expiry = *resp.Expiration
	}
	if resp.Credentials != nil {
		creds = *resp.Credentials
	}

	return creds, expiry, err
}

func (sm *SingleContainerProvider) GetSessionCredentials(ctx context.Context, containerName string) (SessionCredentials, error) {
	sm.mu.Lock()

	// If container name matches, get session
	if sm.containerName == containerName {
		sm.mu.Unlock()
		return sm.resource.Get(SessionState{
			client: sm.client,
			ctx:    ctx,
		})
	}

	// If container name is set and does not match, return error to fall back to bearer token
	sm.mu.Unlock()
	return SessionCredentials{}, ErrFallbackToBearer

}

func (sm *SingleContainerProvider) ExpireSessionCredentials(containerName string) {
	sm.mu.Lock()

	// If container name is set and matches, refresh session
	if sm.containerName == containerName {
		sm.mu.Unlock()
		sm.resource.Expire()
		return
	}

	// If container name is set and does not match, return error to fall back to bearer token
	sm.mu.Unlock()
	return
}

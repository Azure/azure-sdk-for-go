// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

type SessionProvider interface {
	GetSessionCredentials(ctx context.Context, containerName string) (SessionCredentials, error)
	ExpireSessionCredentials(containerName string)
}

// ErrFallbackToBearer indicates that the container does not support sessions
// and the caller should fall back to bearer token authentication.
var ErrFallbackToBearer = errors.New("container does not support sessions, fallback to bearer token")

type sessionState struct {
	client *generated.ContainerClient
	ctx    context.Context
}

// acquire is the function called by temporal.Resource to create a new session.
func acquire(state sessionState) (creds generated.SessionCredentials, expiry time.Time, err error) {
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

// SingleContainerProvider caches a session for a single container using a temporal resource.
// It is safe for concurrent use.
type SingleContainerProvider struct {
	client        *generated.ContainerClient
	containerName string
	resource      *temporal.Resource[SessionCredentials, sessionState]
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

func (sm *SingleContainerProvider) GetSessionCredentials(ctx context.Context, containerName string) (SessionCredentials, error) {
	// If container name matches, get session
	if sm.containerName == containerName {
		return sm.resource.Get(sessionState{
			client: sm.client,
			ctx:    ctx,
		})
	}

	// If container name does not match, return error to fall back to bearer token
	return SessionCredentials{}, ErrFallbackToBearer
}

func (sm *SingleContainerProvider) ExpireSessionCredentials(containerName string) {
	// If container name is set and matches, expire session
	if sm.containerName == containerName {
		sm.resource.Expire()
	}
}

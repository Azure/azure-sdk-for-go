// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package session

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
)

// ErrFallbackToBearer indicates that the container does not support sessions
// and the caller should fall back to bearer token authentication.
var ErrFallbackToBearer = errors.New("container does not support sessions, fallback to bearer token")

type State struct {
	client        *Client
	ctx           context.Context
	containerName string
}

// Manager caches a session for a single container using a temporal resource.
// It is safe for concurrent use.
type Manager struct {
	client *Client
	mu     sync.RWMutex
	// mutex protects access to containerName and resource to ensure thread safety
	// when we extend to support multiple containers in the future, we can change this to a map of containerName to resource
	containerName *string
	resource      *temporal.Resource[ContainerCreateSessionResponse, State]
}

// NewManager creates a new Manager instance with the specified client.
func NewManager(client *Client) *Manager {
	sm := &Manager{
		client: client,
	}

	sm.resource = temporal.NewResource(acquire)
	return sm
}

// acquire is the function called by temporal.Resource to create a new session.
func acquire(state State) (resp ContainerCreateSessionResponse, expiry time.Time, err error) {
	resp, err = state.client.ContainerCreateSession(state.ctx, state.containerName)
	// Fall back to using bearer token if session is unable to be created
	if err != nil {
		return resp, expiry, fmt.Errorf("%w: %v", ErrFallbackToBearer, err)
	}

	if resp.Expiration != nil {
		expiry = *resp.Expiration
	}

	return resp, expiry, err
}

// Get returns a valid session for the specified container.
// If the cached session is for a different container or has expired, a new session is acquired.
// Returns ErrFallbackToBearer if the container does not support sessions.
func (sm *Manager) Get(ctx context.Context, containerName string) (ContainerCreateSessionResponse, error) {
	sm.mu.Lock()

	// If container name not set, set container name, get session
	// If container name is set and matches, get session
	if sm.containerName == nil || *sm.containerName == containerName {
		sm.containerName = &containerName
		sm.mu.Unlock()
		return sm.resource.Get(State{
			client:        sm.client,
			ctx:           ctx,
			containerName: containerName,
		})
	}

	// If container name is set and does not match, return error to fall back to bearer token
	sm.mu.Unlock()
	return ContainerCreateSessionResponse{}, ErrFallbackToBearer

}

// Refresh forces acquisition of a new session for the specified container,
// invalidating any cached session.
func (sm *Manager) Refresh(ctx context.Context, containerName string) (ContainerCreateSessionResponse, error) {
	sm.mu.Lock()

	// If container name not set, return error - must call Get first
	if sm.containerName == nil {
		sm.mu.Unlock()
		return ContainerCreateSessionResponse{}, errors.New("session not initialized, call Get first")
	}

	// If container name is set and matches, refresh session
	if *sm.containerName == containerName {
		sm.mu.Unlock()
		sm.resource.Expire()
		return sm.resource.Get(State{
			client:        sm.client,
			ctx:           ctx,
			containerName: containerName,
		})
	}

	// If container name is set and does not match, return error to fall back to bearer token
	sm.mu.Unlock()
	return ContainerCreateSessionResponse{}, ErrFallbackToBearer
}

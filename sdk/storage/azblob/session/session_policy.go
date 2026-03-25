// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package session

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type sessionContextKey struct{}

// SessionContext holds session-related information passed via context.
type Context struct {
	ContainerURL string
}

// WithSessionContext adds session information to the context.
func WithSessionContext(ctx context.Context, containerURL string) context.Context {
	return context.WithValue(ctx, sessionContextKey{}, &Context{ContainerURL: containerURL})
}

// getSessionContainer extracts the session container from the context, if it exists. Returns nil if no session information is found.
// if container is "", this indicates a session token request for the account rather than a specific container
func getSessionContainer(ctx context.Context) *string {
	if v := ctx.Value(sessionContextKey{}); v != nil {
		return v.(*string)
	}
	return nil
}

type Provider interface {
	GetSessionToken(containerURL string) (string, error)
	RefreshSessionToken(containerURL string) (string, error)
}
type Policy struct {
	manager           *Manager
	bearerTokenPolicy policy.Policy
}

func NewPolicy(manager *Manager, bearerTokenPolicy policy.Policy) *Policy {
	return &Policy{
		manager:           manager,
		bearerTokenPolicy: bearerTokenPolicy,
	}
}

// Do implements the policy.Policy interface.
func (p *Policy) Do(req *policy.Request) (*http.Response, error) {
	// Look at request URL - if its a get blob, extract the container name and get a session token for that container
	if req.Raw().Method == http.MethodGet {

	}

	// Fall back to bearer token policy
	return p.bearerTokenPolicy.Do(req)
}

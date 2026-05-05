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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

const sessionUnavailable = "SessionOperationsTemporarilyUnavailable"
const sessionSchemeNotSupported = "Authentication scheme Session is not supported."

// errFallbackToBearer is a sentinel error indicating that session-based authentication
// is unavailable and the request should fall back to bearer token authentication.
var errFallbackToBearer = errors.New("session unavailable, falling back to bearer token authentication")

type sessionPolicy struct {
	bearerTokenPolicy  policy.Policy
	opts               SessionOptions
	oauthServiceClient *generated.ServiceClient

	// resources maps container names to their temporal.Resource for session credentials.
	resources sync.Map // map[string]*temporal.Resource[sessionCredentials, context.Context]
}

func NewSessionPolicy(opts SessionOptions, bearerTokenPolicy policy.Policy, oauthServiceClient *generated.ServiceClient) (policy.Policy, error) {
	if opts.Mode == SessionModeDisabled || opts.Mode == SessionModeDefault {
		return bearerTokenPolicy, nil
	} else if opts.Mode == SessionModeEnabled {
		if opts.AccountName == "" {
			return nil, errors.New("account name is required for enabled mode")
		}
		sessionPl := &sessionPolicy{
			bearerTokenPolicy:  bearerTokenPolicy,
			opts:               opts,
			oauthServiceClient: oauthServiceClient,
		}
		return sessionPl, nil
	} else {
		return nil, fmt.Errorf("unsupported session mode %v", opts.Mode)
	}
}

// getOrCreateResource returns the temporal.Resource for the given container name,
// creating a new one if it doesn't already exist in the map.
func (p *sessionPolicy) getOrCreateResource(containerName string) *temporal.Resource[sessionCredentials, context.Context] {
	if v, ok := p.resources.Load(containerName); ok {
		return v.(*temporal.Resource[sessionCredentials, context.Context])
	}
	cc := getContainerClient(p.oauthServiceClient, containerName)
	res := temporal.NewResourceWithOptions(acquireSession(cc), temporal.ResourceOptions[sessionCredentials, context.Context]{
		ShouldRefresh: shouldRefreshSession,
	})
	actual, _ := p.resources.LoadOrStore(containerName, res)
	return actual.(*temporal.Resource[sessionCredentials, context.Context])
}

func (p *sessionPolicy) Do(req *policy.Request) (*http.Response, error) {
	containerName, ok := supportsSession(req.Raw())
	if !ok {
		return p.bearerTokenPolicy.Do(req)
	}

	resp, err := p.doWithSession(req, containerName)
	if errors.Is(err, errFallbackToBearer) {
		// rewind the request body before falling back to bearer token authentication,
		// as it may have been consumed by a prior call to req.Next().
		if rwErr := req.RewindBody(); rwErr != nil {
			return nil, rwErr
		}
		return p.bearerTokenPolicy.Do(req)
	}
	return resp, err
}

// doWithSession attempts to authenticate the request using session credentials.
// It applies the session auth header, sends the request, and handles any session-specific errors.
func (p *sessionPolicy) doWithSession(req *policy.Request, containerName string) (*http.Response, error) {
	resource := p.getOrCreateResource(containerName)
	resp, err := p.applySessionReq(req, resource)
	if err == nil {
		return resp, nil
	}
	return p.handleSessionError(req, resp, err, resource)
}

// handleSessionError inspects the error from a session-authenticated request and determines
// whether to fall back to bearer token auth, retry with a new session, or return the error.
func (p *sessionPolicy) handleSessionError(req *policy.Request, resp *http.Response, err error, resource *temporal.Resource[sessionCredentials, context.Context]) (*http.Response, error) {
	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return resp, err
	}

	if (respErr.StatusCode == http.StatusServiceUnavailable && respErr.ErrorCode == sessionUnavailable) ||
		(respErr.StatusCode == http.StatusForbidden && strings.Contains(respErr.Error(), sessionSchemeNotSupported)) {
		// drain the failed response to avoid leaking the connection
		runtime.Drain(resp)
		return nil, errFallbackToBearer
	}

	if respErr.StatusCode == http.StatusUnauthorized {
		// drain the failed response to avoid leaking the connection
		runtime.Drain(resp)

		// rewind the request body before retrying
		if err := req.RewindBody(); err != nil {
			return nil, err
		}

		// retry with new session
		resource.Expire()
		return p.applySessionReq(req, resource)
	}

	return resp, err
}

// applySessionReq signs the request with session credentials and sends it.
func (p *sessionPolicy) applySessionReq(req *policy.Request, resource *temporal.Resource[sessionCredentials, context.Context]) (*http.Response, error) {
	sessionCreds, err := resource.Get(req.Raw().Context())
	if err != nil {
		return nil, err
	}
	if sessionCreds.fallback {
		return nil, errFallbackToBearer
	}

	cred, err := NewSharedKeyCredential(p.opts.AccountName, sessionCreds.key)
	if err != nil {
		return nil, err
	}

	// always set a fresh date so the signature matches the current time, including on retries
	req.Raw().Header.Set(shared.HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))

	stringToSign, err := cred.buildStringToSign(req.Raw())
	if err != nil {
		return nil, err
	}
	signature, err := cred.computeHMACSHA256(stringToSign)
	if err != nil {
		return nil, err
	}
	authHeader := "Session " + sessionCreds.token + ":" + signature
	req.Raw().Header.Set(shared.HeaderAuthorization, authHeader)

	return req.Next()
}

// supportsSession checks if the request can use session-based authentication.
// Currently limited to Get Blob requests (GET method on blob URLs without comp query param).
// Returns the container name and true if session can be used, empty string and false otherwise.
func supportsSession(req *http.Request) (containerName string, ok bool) {
	// Only GET requests are supported for sessions
	if req.Method != http.MethodGet {
		return "", false
	}

	u := req.URL
	if u == nil {
		return "", false
	}

	// Session auth is not supported for requests with comp query parameter
	if u.Query().Get("comp") != "" {
		return "", false
	}

	// Path format: /<container>/<blob>
	path := strings.TrimPrefix(u.Path, "/")
	if path == "" {
		return "", false
	}

	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}

	return parts[0], true
}

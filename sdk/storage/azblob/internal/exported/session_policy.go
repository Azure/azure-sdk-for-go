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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

const SessionExpiring = "session_expiring"
const SessionRevoking = "session_revoking"

type SessionCredentials = generated.SessionCredentials

type Provider interface {
	GetSessionCredentials(ctx context.Context, containerName string) (SessionCredentials, error)
	ExpireSessionCredentials(containerName string)
}

type SessionPolicy struct {
	bearerTokenPolicy policy.Policy
	opts              SessionOptions
	provider          Provider
	refreshMu         sync.Mutex
}

func NewSessionPolicy(opts SessionOptions, bearerTokenPolicy policy.Policy, oauthServiceClient *generated.ServiceClient) (policy.Policy, error) {
	var provider Provider
	switch opts.Mode {
	case SessionModeSingleContainer:
		if opts.AccountName == "" {
			return nil, errors.New("account name is required for singlecontainer mode")
		}
		if opts.ContainerName == "" {
			return nil, errors.New("container name is required for singlecontainer mode")
		}
		provider = NewSingleContainerProvider(oauthServiceClient, opts.ContainerName)
	default:
		return nil, fmt.Errorf("unsupported session mode %v", opts.Mode)
	}

	return &SessionPolicy{
		bearerTokenPolicy: bearerTokenPolicy,
		opts:              opts,
		provider:          provider,
	}, nil
}

func (p *SessionPolicy) Do(req *policy.Request) (*http.Response, error) {
	containerName, ok := canUseSession(req.Raw())
	if !ok {
		return p.bearerTokenPolicy.Do(req)
	}

	resp, err := p.doWithSession(req, containerName)
	if err != nil && errors.Is(err, ErrFallbackToBearer) {
		return p.bearerTokenPolicy.Do(req)
	}
	return resp, err
}

func (p *SessionPolicy) doWithSession(req *policy.Request, containerName string) (*http.Response, error) {
	sessionCreds, err := p.provider.GetSessionCredentials(req.Raw().Context(), containerName)
	if err != nil {
		return nil, err
	}

	resp, err := p.applySessionReq(req, sessionCreds)
	if err == nil {
		p.handleSessionRefresh(resp, containerName)
		return resp, nil
	}

	return p.handleSessionError(req, resp, err, containerName)
}

func (p *SessionPolicy) handleSessionRefresh(resp *http.Response, containerName string) {
	authInfo := getHeader(shared.HeaderXmsAuthInfo, resp.Header)
	if strings.Contains(authInfo, SessionExpiring) || strings.Contains(authInfo, SessionRevoking) {
		// Use TryLock to ensure only one goroutine attempts refresh at a time
		if !p.refreshMu.TryLock() {
			return
		}
		go func() {
			defer p.refreshMu.Unlock()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			_, _ = p.provider.GetSessionCredentials(ctx, containerName)
		}()
	}
}

func (p *SessionPolicy) handleSessionError(req *policy.Request, resp *http.Response, err error, containerName string) (*http.Response, error) {
	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return resp, err
	}

	if resp == nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusServiceUnavailable && respErr.ErrorCode == "SessionOperationsTemporarilyUnavailable" {
		return nil, ErrFallbackToBearer
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return p.retryWithNewSession(req, containerName)
	}

	return resp, err
}

func (p *SessionPolicy) retryWithNewSession(req *policy.Request, containerName string) (*http.Response, error) {
	p.provider.ExpireSessionCredentials(containerName)
	sessionCreds, err := p.provider.GetSessionCredentials(req.Raw().Context(), containerName)
	if err != nil {
		if errors.Is(err, ErrFallbackToBearer) {
			return nil, ErrFallbackToBearer
		}
		return nil, err
	}
	return p.applySessionReq(req, sessionCreds)
}

func (p *SessionPolicy) applySessionReq(req *policy.Request, sessionCreds SessionCredentials) (*http.Response, error) {
	var key, token string
	if sessionCreds.SessionKey != nil {
		key = *sessionCreds.SessionKey
	}
	if sessionCreds.SessionToken != nil {
		token = *sessionCreds.SessionToken
	}
	cred, err := NewSharedKeyCredential(p.opts.AccountName, key)
	if err != nil {
		return nil, err
	}

	if d := getHeader(shared.HeaderXmsDate, req.Raw().Header); d == "" {
		req.Raw().Header.Set(shared.HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))
	}
	stringToSign, err := cred.buildStringToSign(req.Raw())
	if err != nil {
		return nil, err
	}
	signature, err := cred.computeHMACSHA256(stringToSign)
	if err != nil {
		return nil, err
	}
	authHeader := strings.Join([]string{"Session ", token, ":", signature}, "")
	req.Raw().Header.Set(shared.HeaderAuthorization, authHeader)

	return req.Next()
}

// canUseSession checks if the request can use session-based authentication.
// Currently limited to Get Blob requests (GET method on blob URLs without comp query param).
// Returns the container name and true if session can be used, empty string and false otherwise.
func canUseSession(req *http.Request) (containerName string, ok bool) {
	// Only GET requests are supported for sessions
	if req.Method != http.MethodGet {
		return "", false
	}

	u := req.URL
	if u == nil {
		return "", false
	}

	// If there's a 'comp' query param, it's not a Get Blob request
	// (e.g., comp=tags, comp=metadata, comp=properties)
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

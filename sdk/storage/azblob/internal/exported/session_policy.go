// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

type sessionPolicy struct {
	bearerTokenPolicy policy.Policy
	provider          SessionProvider
	accountName       string
}

func NewSessionPolicy(accountName string, provider SessionProvider, bearerTokenPolicy policy.Policy) policy.Policy {
	// If we get here, assumption is opts.Mode = Enabled
	return &sessionPolicy{
		accountName:       accountName,
		provider:          provider,
		bearerTokenPolicy: bearerTokenPolicy,
	}
}

func (p *sessionPolicy) Do(req *policy.Request) (*http.Response, error) {
	if !p.provider.IsRequestEligible(req.Raw()) {
		return p.bearerTokenPolicy.Do(req)
	}

	sessionCreds, err := p.provider.GetSession(req.Raw())
	if err != nil {
		return nil, err
	}

	if !sessionCreds.Fallback() {
		resp, err := p.applySessionReq(req, sessionCreds)
		// a 401 means the service rejected the session; the pipeline surfaces the response
		// without an error, so the status code is what's checked here
		if err != nil || resp.StatusCode != http.StatusUnauthorized {
			return resp, err
		}

		// The session was rejected, so discard it; a new one is acquired on the next eligible
		// request. This request falls back to bearer token authentication.
		// drain the failed response to avoid leaking the connection
		runtime.Drain(resp)

		if invErr := p.provider.InvalidateSession(req.Raw(), sessionCreds); invErr != nil {
			return nil, invErr
		}

		// rewind the request body before falling back to bearer token authentication,
		// as it may have been consumed by the prior call to req.Next().
		if rwErr := req.RewindBody(); rwErr != nil {
			return nil, rwErr
		}
	}

	return p.bearerTokenPolicy.Do(req)
}

// applySessionReq signs the request with the given session credentials and sends it.
func (p *sessionPolicy) applySessionReq(req *policy.Request, sessionCreds SessionCredential) (*http.Response, error) {
	cred, err := NewSharedKeyCredential(p.accountName, sessionCreds.Key())
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
	authHeader := "Session " + sessionCreds.Token() + ":" + signature
	req.Raw().Header.Set(shared.HeaderAuthorization, authHeader)

	return req.Next()
}

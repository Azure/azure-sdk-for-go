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

// sessionPolicy occupies the authentication slot of a token credential pipeline whose sessions
// resolve to enabled. Requests the provider deems eligible are signed with a session; every
// other request, and any request whose session can't be used, goes to the bearer token policy.
type sessionPolicy struct {
	bearerTokenPolicy policy.Policy
	provider          SessionProvider
	accountName       string
}

// NewSessionPolicy returns the authentication policy for a token credential client with sessions
// enabled. accountName signs session requests; bearerTokenPolicy authenticates everything else.
func NewSessionPolicy(accountName string, provider SessionProvider, bearerTokenPolicy policy.Policy) policy.Policy {
	return &sessionPolicy{
		accountName:       accountName,
		provider:          provider,
		bearerTokenPolicy: bearerTokenPolicy,
	}
}

func (p *sessionPolicy) Do(req *policy.Request) (*http.Response, error) {
	if !p.provider.isRequestEligible(req.Raw()) {
		return p.bearerTokenPolicy.Do(req)
	}
	session, err := p.provider.getSession(req.Raw())
	if err != nil {
		return nil, err
	}
	if session.fallback {
		return p.bearerTokenPolicy.Do(req)
	}

	resp, err := p.sendWithSession(req, session)
	// The pipeline returns a 401 as a response rather than an error, so the status code is what
	// tells a rejected session apart.
	if err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	// The service rejected the session. Free the connection, discard the session if it is still
	// the cached one so the next eligible request creates a new one, and send this request once
	// more with a bearer token. This request doesn't wait for a new session.
	runtime.Drain(resp)
	p.provider.invalidateSession(req.Raw(), session)
	if err := req.RewindBody(); err != nil {
		return nil, err
	}
	// Hand the bearer token policy the request as it was before the session was applied. It
	// sets Authorization itself, and bearer requests don't carry the signing date.
	req.Raw().Header.Del(shared.HeaderAuthorization)
	req.Raw().Header.Del(shared.HeaderXmsDate)
	return p.bearerTokenPolicy.Do(req)
}

// sendWithSession signs the request with the session key using the shared key scheme, sets
// "Authorization: Session <token>:<signature>", and sends it.
func (p *sessionPolicy) sendWithSession(req *policy.Request, session sessionCredential) (*http.Response, error) {
	cred, err := NewSharedKeyCredential(p.accountName, session.key)
	if err != nil {
		return nil, err
	}
	// a fresh date for every attempt keeps the signature current on retries
	req.Raw().Header.Set(shared.HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))
	stringToSign, err := cred.buildStringToSign(req.Raw())
	if err != nil {
		return nil, err
	}
	signature, err := cred.computeHMACSHA256(stringToSign)
	if err != nil {
		return nil, err
	}
	req.Raw().Header.Set(shared.HeaderAuthorization, "Session "+session.token+":"+signature)
	return req.Next()
}

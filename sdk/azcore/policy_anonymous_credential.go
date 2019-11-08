// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import "context"

var singletonAnonymousCredential = AnonymousCredential{}

// NewAnonymousCredential creates an anonymous credential for use with HTTP(S) requests that read public resource
// or for use with Shared Access Signatures (SAS).
func NewAnonymousCredential() *AnonymousCredential { return &singletonAnonymousCredential }

// AnonymousCredential represent an anonymous credential.
type AnonymousCredential struct{}

// Do implements the credential's policy interface.
func (p *AnonymousCredential) Do(ctx context.Context, req *Request) (*Response, error) {
	// For anonymous credentials, this is effectively a no-op except for calling the next Policy
	return req.Do(ctx)
}

// marker satisfies the Credential interface making Credential policies "special"
func (p *AnonymousCredential) marker() {}

//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import "net/http"

func anonCredAuthPolicyFunc(AuthenticationOptions) Policy {
	return policyFunc(anonCredPolicyFunc)
}

func anonCredPolicyFunc(req *Request) (*http.Response, error) {
	return req.Next()
}

// NewAnonymousCredential is for use with HTTP(S) requests that read public resource
// or for use with Shared Access Signatures (SAS).
func NewAnonymousCredential() Credential {
	return credentialFunc(anonCredAuthPolicyFunc)
}

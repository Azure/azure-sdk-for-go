//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// KeyCredential

type KeyCredential struct {
	APIKey string
}

// APIKeyPolicy authorizes requests with an API key acquired from a KeyCredential.
type APIKeyPolicy struct {
	header string
	cred   KeyCredential
}

// NewAPIKeyPolicy creates a policy object that authorizes requests with an API Key.
// cred: a KeyCredential implementation.
func NewAPIKeyPolicy(cred KeyCredential, header string) *APIKeyPolicy {
	return &APIKeyPolicy{
		header: header,
		cred:   cred,
	}
}

// Do returns a function which authorizes req with a token from the policy's credential
func (b *APIKeyPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Set(b.header, b.cred.APIKey)
	return req.Next()
}

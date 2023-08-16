//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// KeyCredential is used when doing APIKey-based authentication.
type KeyCredential struct {
	// apiKey is the api key for the client.
	apiKey string
}

// NewKeyCredential creates a KeyCredential containing an API key for
// either Azure OpenAI or OpenAI.
func NewKeyCredential(apiKey string) (KeyCredential, error) {
	return KeyCredential{apiKey: apiKey}, nil
}

// apiKeyPolicy authorizes requests with an API key acquired from a KeyCredential.
type apiKeyPolicy struct {
	header string
	cred   KeyCredential
}

// newAPIKeyPolicy creates a policy object that authorizes requests with an API Key.
// cred: a KeyCredential implementation.
func newAPIKeyPolicy(cred KeyCredential, header string) *apiKeyPolicy {
	return &apiKeyPolicy{
		header: header,
		cred:   cred,
	}
}

// Do returns a function which authorizes req with a token from the policy's credential
func (b *apiKeyPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Set(b.header, b.cred.apiKey)
	return req.Next()
}

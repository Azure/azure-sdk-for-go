//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azanomalydetector

// this file contains handwritten additions to the generated code

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	clientName = "azanomalydetector.Client"
	apiVersion = "v1.1"
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

// Clients

// UnivariateClientOptions contains optional settings for UnivariateClient.
type UnivariateClientOptions struct {
	azcore.ClientOptions
}

// NewUnivariateClient creates a client that accesses Azure Monitor metrics data.
func NewUnivariateClientWithKeyCredential(endpoint string, credential KeyCredential, options *UnivariateClientOptions) (*UnivariateClient, error) {
	if options == nil {
		options = &UnivariateClientOptions{}
	}
	authPolicy := NewAPIKeyPolicy(credential, "Ocp-Apim-Subscription-Key")
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &UnivariateClient{endpoint: endpoint + "/anomalydetector/" + apiVersion, internal: azcoreClient}, nil
}

// MultivariateClientOptions contains optional settings for MultivariateClient.
type MultivariateClientOptions struct {
	azcore.ClientOptions
}

// NewMultivariateClient creates a client that accesses Azure Monitor logs data.
func NewMultivariateClientWithKeyCredential(endpoint string, credential KeyCredential, options *MultivariateClientOptions) (*MultivariateClient, error) {
	if options == nil {
		options = &MultivariateClientOptions{}
	}
	authPolicy := NewAPIKeyPolicy(credential, "Ocp-Apim-Subscription-Key")
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &MultivariateClient{endpoint: endpoint + "/anomalydetector/" + apiVersion, internal: azcoreClient}, nil
}

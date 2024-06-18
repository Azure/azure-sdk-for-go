//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	clientName = "azopenaiassistants.Client"
	tokenScope = "https://cognitiveservices.azure.com/.default"
)

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a new instance of Client that connects to an Azure OpenAI endpoint.
//   - endpoint - Azure OpenAI service endpoint, for example: https://{your-resource-name}.openai.azure.com
//   - credential - used to authorize requests. Usually a credential from [github.com/Azure/azure-sdk-for-go/sdk/azidentity].
//   - options - client options, pass nil to accept the default values.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{tokenScope}, nil)
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy, &azureOpenAIPolicy{}}}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		endpoint: endpoint,
		internal: azcoreClient,
		cd: clientData{
			azure: true,
		},
	}, nil
}

// NewClientWithKeyCredential creates a new instance of Client that connects to an Azure OpenAI endpoint.
//   - endpoint - Azure OpenAI service endpoint, for example: https://{your-resource-name}.openai.azure.com
//   - credential - used to authorize requests with an API Key credential
//   - options - client options, pass nil to accept the default values.
func NewClientWithKeyCredential(endpoint string, credential *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewKeyCredentialPolicy(credential, "api-key", nil)
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy, &azureOpenAIPolicy{}}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &Client{
		endpoint: endpoint,
		internal: azcoreClient,
		cd: clientData{
			azure: true,
		},
	}, nil
}

// NewClientForOpenAI creates a new instance of Client which connects to the public OpenAI endpoint.
//   - endpoint - OpenAI service endpoint, for example: https://api.openai.com/v1
//   - credential - used to authorize requests with an API Key credential
//   - options - client options, pass nil to accept the default values.
func NewClientForOpenAI(endpoint string, credential *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	kp := runtime.NewKeyCredentialPolicy(credential, "authorization", &runtime.KeyCredentialPolicyOptions{
		Prefix: "Bearer ",
	})

	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{
		PerRetry: []policy.Policy{kp, &openAIPolicy{}},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		endpoint: endpoint,
		internal: azcoreClient,
	}, nil
}

// openAIPolicy is an internal pipeline policy to remove the api-version query parameter
type openAIPolicy struct{}

// Do returns a function which adapts a request to target OpenAI.
// Specifically, it removes the api-version query parameter.
func (b *openAIPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Add("OpenAI-Beta", "assistants=v2")
	return req.Next()
}

type azureOpenAIPolicy struct{}

func (b *azureOpenAIPolicy) Do(req *policy.Request) (*http.Response, error) {
	reqQP := req.Raw().URL.Query()

	// TODO: there used to be a constant with this value but it doesn't appear to be generated anymore.
	reqQP.Set("api-version", string("2024-05-01-preview"))

	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Add("OpenAI-Beta", "assistants=v2")
	return req.Next()
}

func (client *Client) formatURL(path string) string {
	if client.cd.azure {
		return runtime.JoinPaths("openai", path)
	}

	return path
}

type clientData struct {
	azure bool
}

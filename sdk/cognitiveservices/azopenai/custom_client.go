//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

// this file contains handwritten additions to the generated code

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	clientName = "azopenai.Client"
	apiVersion = "2023-03-15-preview"
	tokenScope = "https://cognitiveservices.azure.com/.default"
)

// Clients

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a new instance of Client with the specified values.
//   - endpoint - Azure OpenAI service endpoint, for example: https://{your-resource-name}.openai.azure.com
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - deploymentID - the deployment ID of the model to query
//   - options - client options, pass nil to accept the default values.
func NewClient(endpoint string, credential azcore.TokenCredential, deploymentID string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{tokenScope}, nil)
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	fullEndpoint := formatAzureOpenAIURL(endpoint, deploymentID)
	return &Client{endpoint: fullEndpoint, internal: azcoreClient}, nil
}

func formatAzureOpenAIURL(endpoint, deploymentID string) string {
	escapedDeplID := url.PathEscape(deploymentID)
	return runtime.JoinPaths(endpoint, "openai", "deployments", escapedDeplID)
}

// NewClientWithKeyCredential creates a new instance of Client with the specified values.
//   - endpoint - Azure OpenAI service endpoint, for example: https://{your-resource-name}.openai.azure.com
//   - credential - used to authorize requests with an API Key credential
//   - deploymentID - the deployment ID of the model to query
//   - options - client options, pass nil to accept the default values.
func NewClientWithKeyCredential(endpoint string, credential KeyCredential, deploymentID string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := newAPIKeyPolicy(credential, "api-key")
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	fullEndpoint := formatAzureOpenAIURL(endpoint, deploymentID)
	return &Client{endpoint: fullEndpoint, internal: azcoreClient}, nil
}

// NewClientForOpenAI creates a new instance of Client with the specified values.
//   - endpoint - OpenAI service endpoint, for example: https://api.openai.com/v1
//   - credential - used to authorize requests with an API Key credential
//   - options - client options, pass nil to accept the default values.
func NewClientForOpenAI(endpoint string, credential KeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	openAIPolicy := newOpenAIPolicy(credential)
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{openAIPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{endpoint: endpoint, internal: azcoreClient}, nil
}

// openAIPolicy is an internal pipeline policy to remove the api-version query parameter
type openAIPolicy struct {
	cred KeyCredential
}

// newOpenAIPolicy creates a new instance of openAIPolicy.
// cred: a KeyCredential implementation.
func newOpenAIPolicy(cred KeyCredential) *openAIPolicy {
	return &openAIPolicy{cred: cred}
}

// Do returns a function which adapts a request to target OpenAI.
// Specifically, it removes the api-version query parameter.
func (b *openAIPolicy) Do(req *policy.Request) (*http.Response, error) {
	q := req.Raw().URL.Query()
	q.Del("api-version")
	req.Raw().Header.Set("authorization", "Bearer "+b.cred.APIKey)
	return req.Next()
}

// Methods that return streaming response

type streamCompletionsOptions struct {
	CompletionsOptions
	Stream bool `json:"stream"`
}

func (o streamCompletionsOptions) MarshalJSON() ([]byte, error) {
	bytes, err := o.CompletionsOptions.MarshalJSON()
	if err != nil {
		return nil, err
	}
	objectMap := make(map[string]any)
	err = json.Unmarshal(bytes, &objectMap)
	if err != nil {
		return nil, err
	}
	objectMap["stream"] = o.Stream
	return json.Marshal(objectMap)
}

// GetCompletionsStream - Return the completions for a given prompt as a sequence of events.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - GetCompletionsOptions contains the optional parameters for the Client.GetCompletions method.
func (client *Client) GetCompletionsStream(ctx context.Context, body CompletionsOptions, options *GetCompletionsStreamOptions) (GetCompletionsStreamResponse, error) {
	req, err := client.getCompletionsCreateRequest(ctx, CompletionsOptions{}, &GetCompletionsOptions{})

	if err != nil {
		return GetCompletionsStreamResponse{}, err
	}

	if err := runtime.MarshalAsJSON(req, streamCompletionsOptions{body, true}); err != nil {
		return GetCompletionsStreamResponse{}, err
	}

	runtime.SkipBodyDownload(req)

	resp, err := client.internal.Pipeline().Do(req)

	if err != nil {
		return GetCompletionsStreamResponse{}, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return GetCompletionsStreamResponse{}, runtime.NewResponseError(resp)
	}

	return GetCompletionsStreamResponse{
		Events: newEventReader[Completions](resp.Body),
	}, nil
}

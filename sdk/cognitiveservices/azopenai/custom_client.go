//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Package azopenai Azure OpenAI Service provides access to OpenAI's powerful language models including the GPT-4,
// GPT-35-Turbo, and Embeddings model series, as well as image generation using DALL-E.
//
// The [Client] in this package can be used with Azure OpenAI or OpenAI.
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
	tokenScope = "https://cognitiveservices.azure.com/.default"
)

// Clients

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
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azcoreClient,
		clientData: clientData{
			endpoint: endpoint,
			azure:    true,
		},
	}, nil
}

// NewClientWithKeyCredential creates a new instance of Client that connects to an Azure OpenAI endpoint.
//   - endpoint - Azure OpenAI service endpoint, for example: https://{your-resource-name}.openai.azure.com
//   - credential - used to authorize requests with an API Key credential
//   - options - client options, pass nil to accept the default values.
func NewClientWithKeyCredential(endpoint string, credential KeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := newAPIKeyPolicy(credential, "api-key")
	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azcoreClient,
		clientData: clientData{
			endpoint: endpoint,
			azure:    true,
		},
	}, nil
}

// NewClientForOpenAI creates a new instance of Client which connects to the public OpenAI endpoint.
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

	return &Client{
		internal: azcoreClient,
		clientData: clientData{
			endpoint: endpoint,
			azure:    false,
		},
	}, nil
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
	req.Raw().Header.Set("authorization", "Bearer "+b.cred.apiKey)
	return req.Next()
}

// Methods that return streaming response
type streamCompletionsOptions struct {
	// we strip out the 'stream' field from the options exposed to the customer so
	// now we need to add it back in.
	any
	Stream bool `json:"stream"`
}

func (o streamCompletionsOptions) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(o.any)

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
	req, err := client.getCompletionsCreateRequest(ctx, body, &GetCompletionsOptions{})

	if err != nil {
		return GetCompletionsStreamResponse{}, err
	}

	if err := runtime.MarshalAsJSON(req, streamCompletionsOptions{
		any:    body,
		Stream: true,
	}); err != nil {
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
		CompletionsStream: newEventReader[Completions](resp.Body),
	}, nil
}

// GetChatCompletionsStream - Return the chat completions for a given prompt as a sequence of events.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - GetCompletionsOptions contains the optional parameters for the Client.GetCompletions method.
func (client *Client) GetChatCompletionsStream(ctx context.Context, body ChatCompletionsOptions, options *GetChatCompletionsStreamOptions) (GetChatCompletionsStreamResponse, error) {
	req, err := client.getChatCompletionsCreateRequest(ctx, body, &GetChatCompletionsOptions{})

	if err != nil {
		return GetChatCompletionsStreamResponse{}, err
	}

	if err := runtime.MarshalAsJSON(req, streamCompletionsOptions{
		any:    body,
		Stream: true,
	}); err != nil {
		return GetChatCompletionsStreamResponse{}, err
	}

	runtime.SkipBodyDownload(req)

	resp, err := client.internal.Pipeline().Do(req)

	if err != nil {
		return GetChatCompletionsStreamResponse{}, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return GetChatCompletionsStreamResponse{}, runtime.NewResponseError(resp)
	}

	return GetChatCompletionsStreamResponse{
		ChatCompletionsStream: newEventReader[ChatCompletions](resp.Body),
	}, nil
}

func (client *Client) formatURL(path string, deploymentID string) string {
	switch path {
	// https://learn.microsoft.com/en-us/azure/cognitive-services/openai/reference#image-generation
	case "/images/generations:submit":
		return runtime.JoinPaths(client.endpoint, "openai", path)
	default:
		if client.azure {
			escapedDeplID := url.PathEscape(deploymentID)
			return runtime.JoinPaths(client.endpoint, "openai", "deployments", escapedDeplID, path)
		}

		return runtime.JoinPaths(client.endpoint, path)
	}
}

func (client *Client) newError(resp *http.Response) error {
	return newContentFilterResponseError(resp)
}

type clientData struct {
	endpoint string
	azure    bool
}

func getDeploymentID[T ChatCompletionsOptions | CompletionsOptions | EmbeddingsOptions | ImageGenerationOptions](v T) string {
	switch a := any(v).(type) {
	case ChatCompletionsOptions:
		return a.DeploymentID
	case CompletionsOptions:
		return a.DeploymentID
	case EmbeddingsOptions:
		return a.DeploymentID
	case ImageGenerationOptions:
		return ""
	default:
		return ""
	}
}

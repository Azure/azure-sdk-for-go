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
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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

const apiVersion = "2025-01-01-preview"

// NewClient creates a new instance of Client that connects to an Azure OpenAI endpoint.
//   - endpoint - Azure OpenAI service endpoint, for example: https://{your-resource-name}.openai.azure.com
//   - credential - used to authorize requests. Usually a credential from [github.com/Azure/azure-sdk-for-go/sdk/azidentity].
//   - options - client options, pass nil to accept the default values.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{tokenScope}, &policy.BearerTokenOptions{
		InsecureAllowCredentialWithHTTP: allowInsecure(options),
	})

	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy, tempAPIVersionPolicy{}},
	}, &options.ClientOptions)

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
func NewClientWithKeyCredential(endpoint string, credential *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewKeyCredentialPolicy(credential, "api-key", &runtime.KeyCredentialPolicyOptions{
		InsecureAllowCredentialWithHTTP: allowInsecure(options),
	})

	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy, tempAPIVersionPolicy{}},
	}, &options.ClientOptions)
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
func NewClientForOpenAI(endpoint string, credential *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	kp := runtime.NewKeyCredentialPolicy(credential, "authorization", &runtime.KeyCredentialPolicyOptions{
		Prefix:                          "Bearer ",
		InsecureAllowCredentialWithHTTP: allowInsecure(options),
	})

	azcoreClient, err := azcore.NewClient(clientName, version, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			kp,
			newOpenAIPolicy(),
		},
	}, &options.ClientOptions)

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
type openAIPolicy struct{}

// newOpenAIPolicy creates a new instance of openAIPolicy.
func newOpenAIPolicy() *openAIPolicy {
	return &openAIPolicy{}
}

// Do returns a function which adapts a request to target OpenAI.
// Specifically, it removes the api-version query parameter.
func (b *openAIPolicy) Do(req *policy.Request) (*http.Response, error) {
	q := req.Raw().URL.Query()
	q.Del("api-version")
	return req.Next()
}

func (options *CompletionsOptions) toWireType() completionsOptions {
	return completionsOptions{
		Stream:           to.Ptr(false),
		StreamOptions:    nil,
		Prompt:           options.Prompt,
		BestOf:           options.BestOf,
		MaxTokens:        options.MaxTokens,
		Temperature:      options.Temperature,
		TopP:             options.TopP,
		FrequencyPenalty: options.FrequencyPenalty,
		PresencePenalty:  options.PresencePenalty,
		Stop:             options.Stop,
		Echo:             options.Echo,
		LogitBias:        options.LogitBias,
		LogProbs:         options.LogProbs,
		DeploymentName:   options.DeploymentName,
		N:                options.N,
		Seed:             options.Seed,
		Suffix:           options.Suffix,
		User:             options.User,
	}
}

func (options *CompletionsStreamOptions) toWireType() completionsOptions {
	return completionsOptions{
		Stream:           to.Ptr(true),
		StreamOptions:    options.StreamOptions,
		Prompt:           options.Prompt,
		BestOf:           options.BestOf,
		MaxTokens:        options.MaxTokens,
		Temperature:      options.Temperature,
		TopP:             options.TopP,
		FrequencyPenalty: options.FrequencyPenalty,
		PresencePenalty:  options.PresencePenalty,
		Stop:             options.Stop,
		Echo:             options.Echo,
		LogitBias:        options.LogitBias,
		LogProbs:         options.LogProbs,
		DeploymentName:   options.DeploymentName,
		N:                options.N,
		Seed:             options.Seed,
		Suffix:           options.Suffix,
		User:             options.User,
	}
}

func (options *ChatCompletionsOptions) toWireType() chatCompletionsOptions {
	return chatCompletionsOptions{
		Stream:                 to.Ptr(false),
		StreamOptions:          nil,
		Audio:                  options.Audio,
		Metadata:               options.Metadata,
		Modalities:             options.Modalities,
		Prediction:             options.Prediction,
		ReasoningEffort:        options.ReasoningEffort,
		Store:                  options.Store,
		UserSecurityContext:    options.UserSecurityContext,
		Messages:               options.Messages,
		AzureExtensionsOptions: options.AzureExtensionsOptions,
		Enhancements:           options.Enhancements,
		FrequencyPenalty:       options.FrequencyPenalty,
		FunctionCall:           options.FunctionCall,
		Functions:              options.Functions,
		LogitBias:              options.LogitBias,
		LogProbs:               options.LogProbs,
		MaxCompletionTokens:    options.MaxCompletionTokens,
		MaxTokens:              options.MaxTokens,
		DeploymentName:         options.DeploymentName,
		N:                      options.N,
		ParallelToolCalls:      options.ParallelToolCalls,
		PresencePenalty:        options.PresencePenalty,
		ResponseFormat:         options.ResponseFormat,
		Seed:                   options.Seed,
		Stop:                   options.Stop,
		Temperature:            options.Temperature,
		ToolChoice:             options.ToolChoice,
		Tools:                  options.Tools,
		TopLogProbs:            options.TopLogProbs,
		TopP:                   options.TopP,
		User:                   options.User,
	}
}

func (options *ChatCompletionsStreamOptions) toWireType() chatCompletionsOptions {
	return chatCompletionsOptions{
		Stream:                 to.Ptr(true),
		StreamOptions:          options.StreamOptions,
		Audio:                  options.Audio,
		Metadata:               options.Metadata,
		Modalities:             options.Modalities,
		Prediction:             options.Prediction,
		ReasoningEffort:        options.ReasoningEffort,
		Store:                  options.Store,
		UserSecurityContext:    options.UserSecurityContext,
		Messages:               options.Messages,
		AzureExtensionsOptions: options.AzureExtensionsOptions,
		Enhancements:           options.Enhancements,
		FrequencyPenalty:       options.FrequencyPenalty,
		FunctionCall:           options.FunctionCall,
		Functions:              options.Functions,
		LogitBias:              options.LogitBias,
		LogProbs:               options.LogProbs,
		MaxCompletionTokens:    options.MaxCompletionTokens,
		MaxTokens:              options.MaxTokens,
		DeploymentName:         options.DeploymentName,
		N:                      options.N,
		ParallelToolCalls:      options.ParallelToolCalls,
		PresencePenalty:        options.PresencePenalty,
		ResponseFormat:         options.ResponseFormat,
		Seed:                   options.Seed,
		Stop:                   options.Stop,
		Temperature:            options.Temperature,
		ToolChoice:             options.ToolChoice,
		Tools:                  options.Tools,
		TopLogProbs:            options.TopLogProbs,
		TopP:                   options.TopP,
		User:                   options.User,
	}
}

// GetCompletionsStream - Return the completions for a given prompt as a sequence of events.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - GetCompletionsOptions contains the optional parameters for the Client.GetCompletions method.
func (client *Client) GetCompletionsStream(ctx context.Context, body CompletionsStreamOptions, options *GetCompletionsStreamOptions) (GetCompletionsStreamResponse, error) {
	req, err := client.getCompletionsCreateRequest(ctx, body.toWireType(), &GetCompletionsOptions{})
	if err != nil {
		return GetCompletionsStreamResponse{}, err
	}

	runtime.SkipBodyDownload(req)

	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GetCompletionsStreamResponse{}, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK) {
		_ = resp.Body.Close()
		return GetCompletionsStreamResponse{}, runtime.NewResponseError(resp)
	}

	return GetCompletionsStreamResponse{
		CompletionsStream: newEventReader[Completions](resp.Body),
	}, nil
}

// GetCompletions - Gets completions for the provided input prompts. Completions support a wide variety of tasks and generate
// text that continues from or "completes" provided prompt data.
// If the operation fails it returns an *azcore.ResponseError type.
//   - body - The configuration information for a completions request. Completions support a wide variety of tasks and generate
//     text that continues from or "completes" provided prompt data.
//   - options - GetCompletionsOptions contains the optional parameters for the Client.getCompletions method.
func (client *Client) GetCompletions(ctx context.Context, body CompletionsOptions, options *GetCompletionsOptions) (GetCompletionsResponse, error) {
	return client.getCompletions(ctx, body.toWireType(), options)
}

// GetChatCompletionsStream - Return the chat completions for a given prompt as a sequence of events.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - GetCompletionsOptions contains the optional parameters for the Client.GetCompletions method.
func (client *Client) GetChatCompletionsStream(ctx context.Context, body ChatCompletionsStreamOptions, options *GetChatCompletionsStreamOptions) (GetChatCompletionsStreamResponse, error) {
	req, err := client.getChatCompletionsCreateRequest(ctx, body.toWireType(), &GetChatCompletionsOptions{})
	if err != nil {
		return GetChatCompletionsStreamResponse{}, err
	}

	runtime.SkipBodyDownload(req)

	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GetChatCompletionsStreamResponse{}, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK) {
		_ = resp.Body.Close()
		return GetChatCompletionsStreamResponse{}, runtime.NewResponseError(resp)
	}

	return GetChatCompletionsStreamResponse{
		ChatCompletionsStream: newEventReader[ChatCompletions](resp.Body),
	}, nil
}

// GetChatCompletions - Gets chat completions for the provided chat messages. Completions support a wide variety of tasks
// and generate text that continues from or "completes" provided prompt data.
// If the operation fails it returns an *azcore.ResponseError type.
//   - body - The configuration information for a chat completions request. Completions support a wide variety of tasks and generate
//     text that continues from or "completes" provided prompt data.
//   - options - GetChatCompletionsOptions contains the optional parameters for the Client.getChatCompletions method.
func (client *Client) GetChatCompletions(ctx context.Context, body ChatCompletionsOptions, options *GetChatCompletionsOptions) (GetChatCompletionsResponse, error) {
	return client.getChatCompletions(ctx, body.toWireType(), options)
}

func allowInsecure(options *ClientOptions) bool {
	return options != nil && options.InsecureAllowCredentialWithHTTP
}

// NOTE: This is a workaround for an emitter issue, see: https://github.com/Azure/azure-sdk-for-go/issues/23417
type tempAPIVersionPolicy struct{}

func (tavp tempAPIVersionPolicy) Do(req *policy.Request) (*http.Response, error) {
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req.Next()
}

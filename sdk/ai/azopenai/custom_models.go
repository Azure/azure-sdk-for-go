//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

// Models for methods that return streaming response

// GetCompletionsStreamOptions contains the optional parameters for the [Client.GetCompletionsStream] method.
type GetCompletionsStreamOptions struct {
	// placeholder for future optional parameters
}

// GetCompletionsStreamResponse is the response from [Client.GetCompletionsStream].
type GetCompletionsStreamResponse struct {
	// CompletionsStream returns the stream of completions. Token limits and other settings may limit the number of completions returned by the service.
	CompletionsStream *EventReader[Completions]
}

// GetChatCompletionsStreamOptions contains the optional parameters for the [Client.GetChatCompletionsStream] method.
type GetChatCompletionsStreamOptions struct {
	// placeholder for future optional parameters
}

// GetChatCompletionsStreamResponse is the response from [Client.GetChatCompletionsStream].
type GetChatCompletionsStreamResponse struct {
	// ChatCompletionsStream returns the stream of completions. Token limits and other settings may limit the number of chat completions returned by the service.
	ChatCompletionsStream *EventReader[ChatCompletions]
}

// ImageGenerationsDataItem contains the results of image generation.
//
// The field that's set will be based on [ImageGenerationOptions.ResponseFormat] and
// are mutually exclusive.
type ImageGenerationsDataItem struct {
	// Base64Data is set to image data, encoded as a base64 string, if [ImageGenerationOptions.ResponseFormat]
	// was set to [ImageGenerationResponseFormatB64JSON].
	Base64Data *string `json:"b64_json"`

	// URL is the address of a generated image if [ImageGenerationOptions.ResponseFormat] was set
	// to [ImageGenerationResponseFormatURL].
	URL *string `json:"url"`
}

// AzureChatExtensionOptions provides Azure specific options to extend ChatCompletions.
type AzureChatExtensionOptions struct {
	// Extensions is a slice of extensions to the chat completions endpoint, like Azure Cognitive Search.
	Extensions []AzureChatExtensionConfiguration
}

// Error implements the error interface for type Error.
// Note that the message contents are not contractual and can change over time.
func (e *Error) Error() string {
	if e.message == nil {
		return ""
	}

	return *e.message
}

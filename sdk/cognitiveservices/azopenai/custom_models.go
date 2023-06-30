//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import "encoding/json"

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

// ImageGenerationsDataItem holds the result of image generation.
type ImageGenerationsDataItem struct {
	// Result will be:
	// - [ImageLocation] if [ImageGenerationOptions.ResponseFormat] was set to [ImageGenerationResponseFormatURL]
	// - [ImagePayload] if [ImageGenerationOptions.ResponseFormat] was set to [ImageGenerationResponseFormatB64JSON]
	Result any
}

// UnmarshalJSON implements the json.Unmarshaler interface for [ImageGenerationsDataItem].
func (di *ImageGenerationsDataItem) UnmarshalJSON(data []byte) error {
	var v *struct {
		Base64JSON *string `json:"b64_json"`
		URL        *string `json:"url"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// which payload do we have?
	if v.Base64JSON != nil {
		di.Result = ImagePayload{B64JSON: v.Base64JSON}
	} else {
		di.Result = ImageLocation{URL: v.URL}
	}

	return nil
}

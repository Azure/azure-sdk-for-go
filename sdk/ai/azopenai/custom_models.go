//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

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

// ContentFilterResponseError is an error as a result of a request being filtered.
type ContentFilterResponseError struct {
	azcore.ResponseError

	// ContentFilterResults contains Information about the content filtering category, if it has been detected.
	ContentFilterResults *ContentFilterResults
}

// Unwrap returns the inner error for this error.
func (e *ContentFilterResponseError) Unwrap() error {
	return &e.ResponseError
}

func newContentFilterResponseError(resp *http.Response) error {
	respErr := runtime.NewResponseError(resp).(*azcore.ResponseError)

	if respErr.ErrorCode != "content_filter" {
		return respErr
	}

	body, err := runtime.Payload(resp)

	if err != nil {
		return err
	}

	var envelope *struct {
		Error struct {
			InnerError struct {
				FilterResult *ContentFilterResults `json:"content_filter_result"`
			} `json:"innererror"`
		}
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return err
	}

	return &ContentFilterResponseError{
		ResponseError:        *respErr,
		ContentFilterResults: envelope.Error.InnerError.FilterResult,
	}
}

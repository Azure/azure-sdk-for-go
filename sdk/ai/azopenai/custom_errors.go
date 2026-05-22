// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/openai/openai-go/v3"
)

// ContentFilterError can be extracted from an openai.Error using [ExtractContentFilterError].
type ContentFilterError struct {
	OpenAIError *openai.Error
	ContentFilterResultDetailsForPrompt
}

// Error implements the error interface for type ContentFilterError.
func (c *ContentFilterError) Error() string {
	return c.OpenAIError.Error()
}

// Unwrap returns the inner error for this error.
func (c *ContentFilterError) Unwrap() error {
	return c.OpenAIError
}

// ExtractContentFilterError checks the error to see if it contains content filtering
// information. If so it'll assign the resulting information to *contentFilterErr,
// similar to errors.As().
//
// Prompt filtering information will be present if you see an error message similar to
// this: 'The response was filtered due to the prompt triggering'.
// (NOTE: error message is for illustrative purposes, and can change).
//
// Usage looks like this:
//
//	resp, err := chatCompletionsService.New(args)
//
//	var contentFilterErr *azopenai.ContentFilterError
//
//	if openai.ExtractContentFilterError(err, &contentFilterErr) {
//		// contentFilterErr.Hate, contentFilterErr.SelfHarm, contentFilterErr.Sexual or contentFilterErr.Violence
//		// contain information about why content was flagged.
//	}
func ExtractContentFilterError(err error, contentFilterErr **ContentFilterError) bool {
	// This is for a very specific case - when Azure rejects a request, outright, because
	// it violates a content filtering rule. In that case you get a StatusBadRequest, and the
	// underlying response contains a payload with the content filtering details.

	var openaiErr *openai.Error

	if !errors.As(err, &openaiErr) {
		return false
	}

	if openaiErr.Response != nil && openaiErr.Response.StatusCode != http.StatusBadRequest {
		return false
	}

	body, origErr := runtime.Payload(openaiErr.Response)

	if origErr != nil {
		return false
	}

	var envelope *struct {
		Error struct {
			Param      string `json:"prompt"`
			Message    string `json:"message"`
			Code       string `json:"code"`
			Status     int    `json:"status"`
			InnerError struct {
				Code                 string                              `json:"code"`
				ContentFilterResults ContentFilterResultDetailsForPrompt `json:"content_filter_result"`
			} `json:"innererror"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return false
	}

	if envelope.Error.Code != "content_filter" {
		return false
	}

	*contentFilterErr = &ContentFilterError{
		OpenAIError:                         openaiErr,
		ContentFilterResultDetailsForPrompt: envelope.Error.InnerError.ContentFilterResults,
	}

	return true
}

// NonRetriable is a marker method, indicating the request failure is terminal.
func (c *ContentFilterError) NonRetriable() {}

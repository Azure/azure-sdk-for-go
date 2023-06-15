//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

// Models for methods that return streaming response

// GetCompletionsStreamOptions contains the optional parameters for the Client.GetCompletions method.
type GetCompletionsStreamOptions struct {
	// placeholder for future optional parameters
}

// GetCompletionsStreamResponse is the response from [GetCompletionsStream].
type GetCompletionsStreamResponse struct {
	// CompletionsStream returns the stream of completions. Token limits and other settings may limit the number of completions returned by the service.
	CompletionsStream *EventReader[Completions]
}

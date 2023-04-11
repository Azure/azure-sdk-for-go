//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

// Models for methods that return streaming response

// ClientGetCompletionsStreamOptions contains the optional parameters for the Client.GetCompletions method.
type ClientGetCompletionsStreamOptions struct {
	// placeholder for future optional parameters
}

type CompletionEventsResponse struct {
	// REQUIRED; An EventReader to obtain the streaming completions choices associated with this completions response.
	// Generally, n choices are generated per provided prompt with a default value of 1. Token limits and other settings
	// may limit the number of choices generated.
	Events *EventReader[Completions]
}

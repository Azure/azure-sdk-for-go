// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"encoding/json"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/respjson"
)

//
// ChatCompletions (non-streaming)
//

// ChatCompletion wraps an [openai.ChatCompletion], allowing access to Azure specific properties.
type ChatCompletion openai.ChatCompletion

// ChatCompletionChoice wraps an [openai.ChatCompletionChoice], allowing access to Azure specific properties.
type ChatCompletionChoice openai.ChatCompletionChoice

// ChatCompletionMessage wraps an [openai.ChatCompletionMessage], allowing access to Azure specific properties.
type ChatCompletionMessage openai.ChatCompletionMessage

//
// Completions (streaming)
//

// ChatCompletionChunk wraps an [openai.ChatCompletionChunk], allowing access to Azure specific properties.
type ChatCompletionChunk openai.ChatCompletionChunk

// ChatCompletionChunkChoiceDelta wraps an [openai.ChatCompletionChunkChoiceDelta], allowing access to Azure specific properties.
type ChatCompletionChunkChoiceDelta openai.ChatCompletionChunkChoiceDelta

//
// Completions (streaming and non-streaming)
//

// Completion wraps an [openai.Completion], allowing access to Azure specific properties.
type Completion openai.Completion

// CompletionChoice wraps an [openai.CompletionChoice], allowing access to Azure specific properties.
type CompletionChoice openai.CompletionChoice

// PromptFilterResults contains content filtering results for zero or more prompts in the request.
func (c ChatCompletion) PromptFilterResults() ([]ContentFilterResultsForPrompt, error) {
	return unmarshalField[[]ContentFilterResultsForPrompt](c.JSON.ExtraFields["prompt_filter_results"])
}

// ContentFilterResults contains content filtering information for this choice.
func (c ChatCompletionChoice) ContentFilterResults() (*ContentFilterResultsForChoice, error) {
	return unmarshalField[*ContentFilterResultsForChoice](c.JSON.ExtraFields["content_filter_results"])
}

// Context contains additional context information available when Azure OpenAI chat extensions are involved
// in the generation of a corresponding chat completions response.
func (c ChatCompletionMessage) Context() (*AzureChatExtensionsMessageContext, error) {
	return unmarshalField[*AzureChatExtensionsMessageContext](c.JSON.ExtraFields["context"])
}

// PromptFilterResults contains content filtering results for zero or more prompts in the request. In a streaming request,
// results for different prompts may arrive at different times or in different orders.
func (c ChatCompletionChunk) PromptFilterResults() ([]ContentFilterResultsForPrompt, error) {
	return unmarshalField[[]ContentFilterResultsForPrompt](c.JSON.ExtraFields["prompt_filter_results"])
}

// Context contains additional context information available when Azure OpenAI chat extensions are involved
// in the generation of a corresponding chat completions response.
func (c ChatCompletionChunkChoiceDelta) Context() (*AzureChatExtensionsMessageContext, error) {
	return unmarshalField[*AzureChatExtensionsMessageContext](c.JSON.ExtraFields["context"])
}

// PromptFilterResults contains content filtering results for zero or more prompts in the request.
func (c Completion) PromptFilterResults() ([]ContentFilterResultsForPrompt, error) {
	return unmarshalField[[]ContentFilterResultsForPrompt](c.JSON.ExtraFields["prompt_filter_results"])
}

// ContentFilterResults contains content filtering information for this choice.
func (c CompletionChoice) ContentFilterResults() (*ContentFilterResultsForChoice, error) {
	return unmarshalField[*ContentFilterResultsForChoice](c.JSON.ExtraFields["content_filter_results"])
}

// unmarshalField is a generic way for us to unmarshal our 'extra' fields.
func unmarshalField[T any](field respjson.Field) (T, error) {
	var zero T

	raw := field.Raw()
	if len(raw) == 0 {
		return zero, nil
	}

	var obj *T

	if err := json.Unmarshal([]byte(field.Raw()), &obj); err != nil {
		return zero, err
	}

	return *obj, nil
}

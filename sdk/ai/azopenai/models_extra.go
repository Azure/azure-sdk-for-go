// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azopenai

/**
 * This file is required for the generated code to be valid. The difference between the files
 * with a custom_ prefix and ones like this with an _extra suffix is that the _extra files are
 * not modified by the customization scripts, so they can be run safely. Files with the custom_
 * would change if they weren't ignored, so we have to keep them separate.
 */

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ChatCompletionsOptionsFunctionCall - Controls how the model responds to function calls. "none" means the model does not
// call a function, and responds to the end-user. "auto" means the model can pick between an end-user or calling a
// function. Specifying a particular function via {"name": "my_function"} forces the model to call that function. "none" is
// the default when no functions are present. "auto" is the default if functions
// are present.
type ChatCompletionsOptionsFunctionCall struct {
	// IsFunction is true if Value refers to a function name.
	IsFunction bool

	// Value is one of:
	// - "auto", meaning the model can pick between an end-user or calling a function
	// - "none", meaning the model does not call a function,
	// - name of a function, in which case [IsFunction] should be set to true.
	Value *string
}

// MarshalJSON implements the json.Marshaller interface for type ChatCompletionsOptionsFunctionCall.
func (c ChatCompletionsOptionsFunctionCall) MarshalJSON() ([]byte, error) {
	if c.IsFunction {
		if c.Value == nil {
			return nil, errors.New("the Value should be the function name to call, not nil")
		}

		return json.Marshal(map[string]string{"name": *c.Value})
	}

	return json.Marshal(c.Value)
}

// ChatCompletionsToolChoice controls which tool is used for this ChatCompletions call.
// You can choose between:
// - [ChatCompletionsToolChoiceAuto] means the model can pick between generating a message or calling a function.
// - [ChatCompletionsToolChoiceNone] means the model will not call a function and instead generates a message
// - Use the [NewChatCompletionsToolChoice] function to specify a specific tool.
type ChatCompletionsToolChoice struct {
	value any
}

var (
	// ChatCompletionsToolChoiceAuto means the model can pick between generating a message or calling a function.
	ChatCompletionsToolChoiceAuto *ChatCompletionsToolChoice = &ChatCompletionsToolChoice{value: "auto"}

	// ChatCompletionsToolChoiceNone means the model will not call a function and instead generates a message.
	ChatCompletionsToolChoiceNone *ChatCompletionsToolChoice = &ChatCompletionsToolChoice{value: "none"}
)

// NewChatCompletionsToolChoice creates a ChatCompletionsToolChoice for a specific tool.
func NewChatCompletionsToolChoice[T ChatCompletionsToolChoiceFunction](v T) *ChatCompletionsToolChoice {
	return &ChatCompletionsToolChoice{value: v}
}

// ChatCompletionsToolChoiceFunction can be used to force the model to call a particular function.
type ChatCompletionsToolChoiceFunction struct {
	// Name is the name of the function to call.
	Name string
}

// MarshalJSON implements the json.Marshaller interface for type ChatCompletionsToolChoiceFunction.
func (tf ChatCompletionsToolChoiceFunction) MarshalJSON() ([]byte, error) {
	type jsonInnerFunc struct {
		Name string `json:"name"`
	}

	type jsonFormat struct {
		Type     string        `json:"type"`
		Function jsonInnerFunc `json:"function"`
	}

	return json.Marshal(jsonFormat{
		Type:     "function",
		Function: jsonInnerFunc(tf),
	})
}

// MarshalJSON implements the json.Marshaller interface for type ChatCompletionsToolChoice.
func (tc ChatCompletionsToolChoice) MarshalJSON() ([]byte, error) {
	return json.Marshal(tc.value)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ChatCompletionsToolChoice.
func (tc *ChatCompletionsToolChoice) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &tc.value)
}

// ChatRequestAssistantMessageContent represents the content for a [azopenai.ChatRequestAssistantMessage].
// NOTE: This should be created using [azopenai.NewChatRequestAssistantMessageContent]
type ChatRequestAssistantMessageContent struct {
	value any
}

// NewChatRequestAssistantMessageContent creates a [azopenai.ChatRequestAssistantMessageContent].
func NewChatRequestAssistantMessageContent[T []ChatMessageRefusalContentItem | []ChatMessageTextContentItem | string](value T) *ChatRequestAssistantMessageContent {
	switch any(value).(type) {
	case []ChatMessageRefusalContentItem:
		return &ChatRequestAssistantMessageContent{value: value}
	case []ChatMessageTextContentItem:
		return &ChatRequestAssistantMessageContent{value: value}
	case string:
		return &ChatRequestAssistantMessageContent{value: value}
	default:
		panic(fmt.Sprintf("Invalid type %T for ChatRequestAssistantMessageContent", value))
	}
}

// MarshalJSON implements the json.Marshaller interface for type ChatRequestAssistantMessageContent.
func (c ChatRequestAssistantMessageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// ChatRequestSystemMessageContent contains the content for a [ChatRequestSystemMessage].
// NOTE: This should be created using [azopenai.NewChatRequestSystemMessageContent]
type ChatRequestSystemMessageContent struct {
	value any
}

// MarshalJSON implements the json.Marshaller interface for type ChatRequestSystemMessageContent.
func (c ChatRequestSystemMessageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// NewChatRequestSystemMessageContent creates a [azopenai.ChatRequestSystemMessageContent].
func NewChatRequestSystemMessageContent[T []ChatMessageTextContentItem | string](value T) *ChatRequestSystemMessageContent {
	switch any(value).(type) {
	case []ChatMessageTextContentItem:
		return &ChatRequestSystemMessageContent{value: value}
	case string:
		return &ChatRequestSystemMessageContent{value: value}
	default:
		panic(fmt.Sprintf("Invalid type %T for ChatRequestSystemMessageContent", value))
	}
}

// ChatRequestDeveloperMessageContent contains the content for a [ChatRequestDeveloperMessage].
// NOTE: This should be created using [azopenai.NewChatRequestDeveloperMessageContent]
type ChatRequestDeveloperMessageContent struct {
	value any
}

// MarshalJSON implements the json.Marshaller interface for type ChatRequestSystemMessageContent.
func (c ChatRequestDeveloperMessageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// NewChatRequestDeveloperMessageContent creates a [azopenai.ChatRequestDeveloperMessageContent].
func NewChatRequestDeveloperMessageContent[T []ChatMessageTextContentItem | string](value T) *ChatRequestDeveloperMessageContent {
	switch any(value).(type) {
	case []ChatMessageTextContentItem:
		return &ChatRequestDeveloperMessageContent{value: value}
	case string:
		return &ChatRequestDeveloperMessageContent{value: value}
	default:
		panic(fmt.Sprintf("Invalid type %T for ChatRequestDeveloperMessageContent", value))
	}
}

// ChatRequestToolMessageContent contains the content for a [ChatRequestToolMessage].
// NOTE: This should be created using [azopenai.NewChatRequestToolMessageContent]
type ChatRequestToolMessageContent struct {
	value any
}

// MarshalJSON implements the json.Marshaller interface for type ChatRequestToolMessageContent.
func (c ChatRequestToolMessageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// NewChatRequestToolMessageContent creates a [azopenai.ChatRequestToolMessageContent].
func NewChatRequestToolMessageContent[T []ChatMessageTextContentItem | string](value T) *ChatRequestToolMessageContent {
	switch any(value).(type) {
	case []ChatMessageTextContentItem:
		return &ChatRequestToolMessageContent{value: value}
	case string:
		return &ChatRequestToolMessageContent{value: value}
	default:
		panic(fmt.Sprintf("Invalid type %T for ChatRequestToolMessageContent", value))
	}
}

// ChatRequestUserMessageContent contains the user prompt - either as a single string
// or as a []ChatCompletionRequestMessageContentPart, enabling images and text as input.
//
// NOTE: This should be created using [azopenai.NewChatRequestUserMessageContent]
type ChatRequestUserMessageContent struct {
	value any
}

// NewChatRequestUserMessageContent creates a [azopenai.ChatRequestUserMessageContent].
func NewChatRequestUserMessageContent[T string | []ChatCompletionRequestMessageContentPartClassification](v T) *ChatRequestUserMessageContent {
	switch actualV := any(v).(type) {
	case string:
		return &ChatRequestUserMessageContent{value: &actualV}
	case []ChatCompletionRequestMessageContentPartClassification:
		return &ChatRequestUserMessageContent{value: actualV}
	}
	return &ChatRequestUserMessageContent{}
}

// MarshalJSON implements the json.Marshaller interface for type Error.
func (c ChatRequestUserMessageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ChatRequestUserMessageContent.
func (c *ChatRequestUserMessageContent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.value)
}

// EmbeddingItem - Representation of a single embeddings relatedness comparison.
type EmbeddingItem struct {
	// List of embeddings value for the input prompt. These represent a measurement of the vector-based relatedness
	// of the provided input when when [EmbeddingEncodingFormatFloat] is specified.
	Embedding []float32

	// EmbeddingBase64 represents the Embeddings when [EmbeddingEncodingFormatBase64] is specified.
	EmbeddingBase64 string

	// REQUIRED; Index of the prompt to which the EmbeddingItem corresponds.
	Index *int32

	// The object type which is always 'embedding'.
	Object string
}

func deserializeEmbeddingsArray(msg json.RawMessage, embeddingItem *EmbeddingItem) error {
	if len(msg) == 0 {
		return nil
	}

	if msg[0] == '"' && len(msg) > 2 && msg[len(msg)-1] == '"' {
		var s = string(msg)
		embeddingItem.EmbeddingBase64 = s[1 : len(s)-1]
		return nil
	}

	return json.Unmarshal(msg, &embeddingItem.Embedding)
}

// MongoDBChatExtensionParametersEmbeddingDependency contains the embedding dependency for the [MongoDBChatExtensionParameters].
// NOTE: This should be created using [azopenai.NewMongoDBChatExtensionParametersEmbeddingDependency]
type MongoDBChatExtensionParametersEmbeddingDependency struct {
	value any
}

// NewMongoDBChatExtensionParametersEmbeddingDependency creates a [azopenai.MongoDBChatExtensionParametersEmbeddingDependency].
func NewMongoDBChatExtensionParametersEmbeddingDependency[T OnYourDataDeploymentNameVectorizationSource | OnYourDataEndpointVectorizationSource](value T) *MongoDBChatExtensionParametersEmbeddingDependency {
	switch any(value).(type) {
	case OnYourDataDeploymentNameVectorizationSource:
		return &MongoDBChatExtensionParametersEmbeddingDependency{value: value}
	case OnYourDataEndpointVectorizationSource:
		return &MongoDBChatExtensionParametersEmbeddingDependency{value: value}
	default:
		panic(fmt.Sprintf("Invalid type %T for MongoDBChatExtensionParametersEmbeddingDependency", value))
	}
}

// MarshalJSON implements the json.Marshaller interface for type MongoDBChatExtensionParametersEmbeddingDependency.
func (c MongoDBChatExtensionParametersEmbeddingDependency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// PredictionContentContent contains the content for a [PredictionContent].
// NOTE: This should be created using [azopenai.NewPredictionContentContent]
type PredictionContentContent struct {
	value any
}

// MarshalJSON implements the json.Marshaller interface for type PredictionContentContent.
func (c PredictionContentContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// NewPredictionContentContent creates a [azopenai.PredictionContentContent].
func NewPredictionContentContent[T []ChatMessageTextContentItem | string](value T) *PredictionContentContent {
	switch any(value).(type) {
	case []ChatMessageTextContentItem:
		return &PredictionContentContent{value: value}
	case string:
		return &PredictionContentContent{value: value}
	default:
		panic(fmt.Sprintf("Invalid type %T for PredictionContentContent", value))
	}
}

// ContentFilterResponseError is an error as a result of a request being filtered.
type ContentFilterResponseError struct {
	azcore.ResponseError

	// ContentFilterResults contains Information about the content filtering category, if it has been detected.
	ContentFilterResults *ContentFilterResults
}

// ContentFilterResults are the content filtering results for a [ContentFilterResponseError].
type ContentFilterResults struct {
	// Describes language attacks or uses that include pejorative or discriminatory language with reference to a person or identity
	// group on the basis of certain differentiating attributes of these groups
	// including but not limited to race, ethnicity, nationality, gender identity and expression, sexual orientation, religion,
	// immigration status, ability status, personal appearance, and body size.
	Hate *ContentFilterResult `json:"hate"`

	// Describes language related to physical actions intended to purposely hurt, injure, or damage one’s body, or kill oneself.
	SelfHarm *ContentFilterResult `json:"self_harm"`

	// Describes language related to anatomical organs and genitals, romantic relationships, acts portrayed in erotic or affectionate
	// terms, physical sexual acts, including those portrayed as an assault or a
	// forced sexual violent act against one’s will, prostitution, pornography, and abuse.
	Sexual *ContentFilterResult `json:"sexual"`

	// Describes language related to physical actions intended to hurt, injure, damage, or kill someone or something; describes
	// weapons, etc.
	Violence *ContentFilterResult `json:"violence"`
}

// Unwrap returns the inner error for this error.
func (e *ContentFilterResponseError) Unwrap() error {
	return &e.ResponseError
}

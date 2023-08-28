// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/stretchr/testify/require"
)

type Params struct {
	Type       string                   `json:"type"`
	Properties map[string]ParamProperty `json:"properties"`
	Required   []string                 `json:"required,omitempty"`
}

type ParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

func TestGetChatCompletions_usingFunctions(t *testing.T) {
	// https://platform.openai.com/docs/guides/gpt/function-calling

	t.Run("OpenAI", func(t *testing.T) {
		chatClient := newOpenAIClientForTest(t)
		testChatCompletionsFunctions(t, chatClient, openAI)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		chatClient := newAzureOpenAIClientForTest(t, azureOpenAI)
		testChatCompletionsFunctions(t, chatClient, azureOpenAI)
	})
}

func testChatCompletionsFunctions(t *testing.T, chatClient *azopenai.Client, tv testVars) {
	body := azopenai.ChatCompletionsOptions{
		DeploymentID: tv.ChatCompletions,
		Messages: []azopenai.ChatMessage{
			{
				Role:    to.Ptr(azopenai.ChatRoleUser),
				Content: to.Ptr("What's the weather like in Boston, MA, in celsius?"),
			},
		},
		FunctionCall: &azopenai.ChatCompletionsOptionsFunctionCall{
			Value: to.Ptr("auto"),
		},
		Functions: []azopenai.FunctionDefinition{
			{
				Name:        to.Ptr("get_current_weather"),
				Description: to.Ptr("Get the current weather in a given location"),
				Parameters: Params{
					Required: []string{"location"},
					Type:     "object",
					Properties: map[string]ParamProperty{
						"location": {
							Type:        "string",
							Description: "The city and state, e.g. San Francisco, CA",
						},
						"unit": {
							Type: "string",
							Enum: []string{"celsius", "fahrenheit"},
						},
					},
				},
			},
		},
		Temperature: to.Ptr[float32](0.0),
	}

	resp, err := chatClient.GetChatCompletions(context.Background(), body, nil)
	require.NoError(t, err)

	funcCall := resp.ChatCompletions.Choices[0].Message.FunctionCall

	require.Equal(t, "get_current_weather", *funcCall.Name)

	type location struct {
		Location string `json:"location"`
		Unit     string `json:"unit"`
	}

	var funcParams *location
	err = json.Unmarshal([]byte(*funcCall.Arguments), &funcParams)
	require.NoError(t, err)

	require.Equal(t, location{Location: "Boston, MA", Unit: "celsius"}, *funcParams)
}

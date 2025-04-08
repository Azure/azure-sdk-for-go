// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestGetChatCompletions_usingFunctions(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	// https://platform.openai.com/docs/guides/gpt/function-calling

	parametersJSON, err := json.Marshal(map[string]any{
		"required": []string{"location"},
		"type":     "object",
		"properties": map[string]any{
			"location": map[string]any{
				"type":        "string",
				"description": "The city and state, e.g. San Francisco, CA",
			},
			"unit": map[string]any{
				"type": "string",
				"enum": []string{"celsius", "fahrenheit"},
			},
		},
	})
	require.NoError(t, err)

	testFn := func(t *testing.T, chatClient *azopenai.Client, deploymentName string, toolChoice *azopenai.ChatCompletionsToolChoice) {
		body := azopenai.ChatCompletionsOptions{
			DeploymentName: &deploymentName,
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestAssistantMessage{
					Content: azopenai.NewChatRequestAssistantMessageContent("What's the weather like in Boston, MA, in celsius?"),
				},
			},
			Tools: []azopenai.ChatCompletionsToolDefinitionClassification{
				&azopenai.ChatCompletionsFunctionToolDefinition{
					Function: &azopenai.ChatCompletionsFunctionToolDefinitionFunction{
						Name:        to.Ptr("get_current_weather"),
						Description: to.Ptr("Get the current weather in a given location"),
						Parameters:  parametersJSON,
					},
				},
			},
			ToolChoice:  toolChoice,
			Temperature: to.Ptr[float32](0.0),
		}

		resp, err := chatClient.GetChatCompletions(context.Background(), body, nil)
		customRequireNoError(t, err, true)

		funcCall := resp.Choices[0].Message.ToolCalls[0].(*azopenai.ChatCompletionsFunctionToolCall).Function

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

	useSpecificTool := azopenai.NewChatCompletionsToolChoice(
		azopenai.ChatCompletionsToolChoiceFunction{Name: "get_current_weather"},
	)

	t.Run("AzureOpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, azureOpenAI.ChatCompletions.Endpoint)

		testData := []struct {
			Model      string
			ToolChoice *azopenai.ChatCompletionsToolChoice
		}{
			// all of these variants use the tool provided - auto just also works since we did provide
			// a tool reference and ask a question to use it.
			{Model: azureOpenAI.ChatCompletions.Model, ToolChoice: nil},
			{Model: azureOpenAI.ChatCompletions.Model, ToolChoice: azopenai.ChatCompletionsToolChoiceAuto},
			{Model: azureOpenAI.ChatCompletions.Model, ToolChoice: useSpecificTool},
		}

		for _, td := range testData {
			testFn(t, chatClient, td.Model, td.ToolChoice)
		}
	})

	t.Run("OpenAI", func(t *testing.T) {
		testData := []struct {
			EPM        endpointWithModel
			ToolChoice *azopenai.ChatCompletionsToolChoice
		}{
			// all of these variants use the tool provided - auto just also works since we did provide
			// a tool reference and ask a question to use it.
			{EPM: openAI.ChatCompletions, ToolChoice: nil},
			{EPM: openAI.ChatCompletions, ToolChoice: azopenai.ChatCompletionsToolChoiceAuto},
			{EPM: openAI.ChatCompletionsLegacyFunctions, ToolChoice: useSpecificTool},
		}

		for i, td := range testData {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				chatClient := newTestClient(t, td.EPM.Endpoint)
				testFn(t, chatClient, td.EPM.Model, td.ToolChoice)
			})
		}
	})
}

func TestGetChatCompletions_usingFunctions_legacy(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	parametersJSON, err := json.Marshal(map[string]any{
		"required": []string{"location"},
		"type":     "object",
		"properties": map[string]any{
			"location": map[string]any{
				"type":        "string",
				"description": "The city and state, e.g. San Francisco, CA",
			},
			"unit": map[string]any{
				"type": "string",
				"enum": []string{"celsius", "fahrenheit"},
			},
		},
	})
	require.NoError(t, err)

	testFn := func(t *testing.T, epm endpointWithModel) {
		client := newTestClient(t, epm.Endpoint)

		body := azopenai.ChatCompletionsOptions{
			DeploymentName: &epm.Model,
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestAssistantMessage{
					Content: azopenai.NewChatRequestAssistantMessageContent("What's the weather like in Boston, MA, in celsius?"),
				},
			},
			FunctionCall: &azopenai.ChatCompletionsOptionsFunctionCall{
				Value: to.Ptr("auto"),
			},
			Functions: []azopenai.FunctionDefinition{
				{
					Name:        to.Ptr("get_current_weather"),
					Description: to.Ptr("Get the current weather in a given location"),
					Parameters:  parametersJSON,
				},
			},
			Temperature: to.Ptr[float32](0.0),
		}

		resp, err := client.GetChatCompletions(context.Background(), body, nil)
		customRequireNoError(t, err, true)

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

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, azureOpenAI.ChatCompletionsLegacyFunctions)
	})

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, openAI.ChatCompletions)
	})

	t.Run("OpenAI.LegacyFunctions", func(t *testing.T) {
		testFn(t, openAI.ChatCompletionsLegacyFunctions)

	})
}

func TestGetChatCompletions_usingFunctions_streaming(t *testing.T) {
	parametersJSON, err := json.Marshal(map[string]any{
		"required": []string{"location"},
		"type":     "object",
		"properties": map[string]any{
			"location": map[string]any{
				"type":        "string",
				"description": "The city and state, e.g. San Francisco, CA",
			},
			"unit": map[string]any{
				"type": "string",
				"enum": []string{"celsius", "fahrenheit"},
			},
		},
	})
	require.NoError(t, err)

	testFn := func(t *testing.T, epm endpointWithModel) {
		body := azopenai.ChatCompletionsStreamOptions{
			DeploymentName: &epm.Model,
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestAssistantMessage{
					Content: azopenai.NewChatRequestAssistantMessageContent("What's the weather like in Boston, MA, in celsius?"),
				},
			},
			Tools: []azopenai.ChatCompletionsToolDefinitionClassification{
				&azopenai.ChatCompletionsFunctionToolDefinition{
					Function: &azopenai.ChatCompletionsFunctionToolDefinitionFunction{
						Name:        to.Ptr("get_current_weather"),
						Description: to.Ptr("Get the current weather in a given location"),
						Parameters:  parametersJSON,
					},
				},
			},
			Temperature: to.Ptr[float32](0.0),
		}

		chatClient := newTestClient(t, epm.Endpoint)

		resp, err := chatClient.GetChatCompletionsStream(context.Background(), body, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp)

		defer func() {
			err := resp.ChatCompletionsStream.Close()
			require.NoError(t, err)
		}()

		// these results are way trickier than they should be, but we have to accumulate across
		// multiple fields to get a full result.

		funcCall := &azopenai.FunctionCall{
			Arguments: to.Ptr(""),
			Name:      to.Ptr(""),
		}

		for {
			streamResp, err := resp.ChatCompletionsStream.Read()
			require.NoError(t, err)

			if len(streamResp.Choices) == 0 {
				// there are prompt filter results.
				require.NotEmpty(t, streamResp.PromptFilterResults)
				continue
			}

			if streamResp.Choices[0].FinishReason != nil {
				break
			}

			var functionToolCall *azopenai.ChatCompletionsFunctionToolCall = streamResp.Choices[0].Delta.ToolCalls[0].(*azopenai.ChatCompletionsFunctionToolCall)
			require.NotEmpty(t, functionToolCall.Function)

			if functionToolCall.Function.Arguments != nil {
				*funcCall.Arguments += *functionToolCall.Function.Arguments
			}

			if functionToolCall.Function.Name != nil {
				*funcCall.Name += *functionToolCall.Function.Name
			}
		}

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

	// https://platform.openai.com/docs/guides/gpt/function-calling

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, azureOpenAI.ChatCompletions)
	})

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, openAI.ChatCompletions)
	})
}

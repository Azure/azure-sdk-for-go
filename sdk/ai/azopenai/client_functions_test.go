// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
	"github.com/stretchr/testify/require"
)

var weatherFuncTool = []openai.ChatCompletionToolUnionParam{{
	OfFunction: &openai.ChatCompletionFunctionToolParam{
		Function: shared.FunctionDefinitionParam{
			Name:        "get_current_weather",
			Description: openai.String("Get the current weather in a given location"),
			Parameters: openai.FunctionParameters{
				"required": []string{"location"},
				"type":     "object",
				"properties": map[string]interface{}{
					"location": map[string]string{
						"type":        "string",
						"description": "The city and state, e.g. San Francisco, CA",
					},
					"unit": map[string]interface{}{
						"type": "string",
						"enum": []string{"celsius", "fahrenheit"},
					},
				},
			},
		},
	},
}}

func TestGetChatCompletions_usingFunctions(t *testing.T) {
	// https://platform.openai.com/docs/guides/gpt/function-calling

	testFn := func(t *testing.T, chatClient *openai.Client, deploymentName string, toolChoice *openai.ChatCompletionToolChoiceOptionUnionParam) {
		body := openai.ChatCompletionNewParams{
			Model: openai.ChatModel(deploymentName),
			Messages: []openai.ChatCompletionMessageParamUnion{{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String("What's the weather like in Boston, MA, in celsius?"),
					},
				},
			}},
			Tools:       weatherFuncTool,
			Temperature: openai.Float(0.0),
		}
		if toolChoice != nil {
			body.ToolChoice = *toolChoice
		}

		resp, err := chatClient.Chat.Completions.New(context.Background(), body)
		require.NoError(t, err)

		funcCall := resp.Choices[0].Message.ToolCalls[0]

		if recording.GetRecordMode() == recording.PlaybackMode {
			require.Equal(t, "Sanitized", funcCall.Function.Name)
		} else {
			require.Equal(t, "get_current_weather", funcCall.Function.Name)
		}

		type location struct {
			Location string `json:"location"`
			Unit     string `json:"unit"`
		}

		var funcParams *location
		err = json.Unmarshal([]byte(funcCall.Function.Arguments), &funcParams)
		require.NoError(t, err)

		require.Equal(t, location{Location: "Boston, MA", Unit: "celsius"}, *funcParams)
	}

	chatClient := newStainlessTestClientWithAzureURL(t, azureOpenAI.ChatCompletions.Endpoint)

	testData := []struct {
		Model      string
		ToolChoice *openai.ChatCompletionToolChoiceOptionUnionParam
	}{
		// all of these variants use the tool provided - auto just also works since we did provide
		// a tool reference and ask a question to use it.
		{Model: azureOpenAI.ChatCompletions.Model, ToolChoice: nil},
		{Model: azureOpenAI.ChatCompletions.Model, ToolChoice: &openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: openai.String("auto"),
		}},
		{Model: azureOpenAI.ChatCompletions.Model, ToolChoice: &openai.ChatCompletionToolChoiceOptionUnionParam{
			OfFunctionToolChoice: &openai.ChatCompletionNamedToolChoiceParam{
				Function: openai.ChatCompletionNamedToolChoiceFunctionParam{
					Name: "get_current_weather",
				},
			},
		}},
	}

	for _, td := range testData {
		testFn(t, &chatClient, td.Model, td.ToolChoice)
	}
}

func TestGetChatCompletions_usingFunctions_streaming(t *testing.T) {
	body := openai.ChatCompletionNewParams{
		Model: openai.ChatModel(azureOpenAI.ChatCompletions.Model),
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String("What's the weather like in Boston, MA, in celsius?"),
				},
			},
		}},
		Tools:       weatherFuncTool,
		Temperature: openai.Float(0.0),
	}

	chatClient := newStainlessTestClientWithAzureURL(t, azureOpenAI.ChatCompletions.Endpoint)

	stream := chatClient.Chat.Completions.NewStreaming(context.Background(), body)

	defer func() {
		err := stream.Close()
		require.NoError(t, err)
	}()

	// these results are way trickier than they should be, but we have to accumulate across
	// multiple fields to get a full result.

	funcCall := &struct {
		Arguments *string
		Name      *string
	}{
		Arguments: to.Ptr(""),
		Name:      to.Ptr(""),
	}

	for stream.Next() {
		chunk := stream.Current()

		if len(chunk.Choices) == 0 {
			azureChunk := azopenai.ChatCompletionChunk(chunk)

			promptFilterResults, err := azureChunk.PromptFilterResults()
			require.NoError(t, err)

			// there are prompt filter results.
			require.NotEmpty(t, promptFilterResults)
			continue
		}

		if chunk.Choices[0].FinishReason != "" {
			require.Equal(t, "tool_calls", chunk.Choices[0].FinishReason)
			continue
		}

		functionToolCall := chunk.Choices[0].Delta.ToolCalls[0]

		require.NotEmpty(t, functionToolCall.Function)

		*funcCall.Arguments += functionToolCall.Function.Arguments
		*funcCall.Name += functionToolCall.Function.Name
	}

	require.NoError(t, stream.Err())
	require.Equal(t, "get_current_weather", *funcCall.Name)

	type location struct {
		Location string `json:"location"`
		Unit     string `json:"unit"`
	}

	var funcParams *location
	err := json.Unmarshal([]byte(*funcCall.Arguments), &funcParams)
	require.NoError(t, err)

	require.Equal(t, location{Location: "Boston, MA", Unit: "celsius"}, *funcParams)
}

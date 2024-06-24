//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestStreaming_CreateThreadAndRunStream(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		asstID := *createResp.ID

		resp, err := client.CreateThreadAndRunStream(context.Background(), azopenaiassistants.CreateAndRunThreadBody{
			AssistantID:    &asstID,
			DeploymentName: &assistantsModel,
			Instructions:   to.Ptr("You're a helpful assistant that provides answers on questions about beetles."),
			// TODO: if you set this you can trigger the incomplete message
			//MaxCompletionTokens: to.Ptr[int32](100),
			Thread: &azopenaiassistants.CreateThreadBody{
				Messages: []azopenaiassistants.CreateMessageBody{
					{
						Content: to.Ptr("Tell me about the origins of the humble beetle in two sentences."),
						Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
					},
				},
			},
		}, nil)
		requireNoErr(t, azure, err)

		defer func() {
			err = resp.Stream.Close()
			require.NoError(t, err)
		}()

		processStream(t, azure, asstID, streamScenarioThreadAndRun, resp.Stream)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestStreaming_CreateRunStream(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			Instructions: "You're a helpful assistant that provides answers on questions about beetles.",
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		asstID := *createResp.ID

		createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.CreateThreadBody{
			Messages: []azopenaiassistants.CreateMessageBody{
				{
					Content: to.Ptr("Tell me about the origins of the humble beetle in two sentences."),
					Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
				},
			},
		}, nil)
		requireNoErr(t, azure, err)

		createRunResp, err := client.CreateRunStream(context.Background(), *createThreadResp.ID, azopenaiassistants.CreateRunBody{
			AssistantID: &asstID,
		}, nil)
		requireNoErr(t, azure, err)

		defer func() {
			err = createRunResp.Stream.Close()
			require.NoError(t, err)
		}()

		processStream(t, azure, asstID, streamScenarioRun, createRunResp.Stream)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestStreaming_SubmitToolOutputsAndRunStream(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createAssistantResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		// create our thread and a run for that thread with our question "What's the weather like in Boston, MA, in celsius?"
		var threadID string
		{
			createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.CreateThreadBody{}, nil)
			requireNoErr(t, azure, err)

			t.Cleanup(func() {
				_, err := client.DeleteThread(context.Background(), *createThreadResp.ID, nil)
				requireNoErr(t, azure, err)
			})

			threadID = *createThreadResp.ID

			_, err = client.CreateMessage(context.Background(), threadID, azopenaiassistants.CreateMessageBody{
				Content: to.Ptr("What's the weather like in Boston, MA, in celsius?"),
				Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
			}, nil)
			requireNoErr(t, azure, err)

			// run the thread
			createRunResp, err := client.CreateRunStream(context.Background(), *createThreadResp.ID, azopenaiassistants.CreateRunBody{
				AssistantID:  createAssistantResp.ID,
				Instructions: to.Ptr("Use functions to answer questions, when possible."),
			}, nil)
			requireNoErr(t, azure, err)

			defer func() {
				err = createRunResp.Stream.Close()
				require.NoError(t, err)
			}()

			threadRun := processStream(t, azure, *createAssistantResp.ID, streamScenarioRun, createRunResp.Stream)
			require.NotNil(t, threadRun)

			submitToolOutputsWithStreaming(t, client, threadRun, azure)
		}
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func submitToolOutputsWithStreaming(t *testing.T, client *azopenaiassistants.Client, lastResp *azopenaiassistants.ThreadRun, azure bool) {
	var weatherFuncArgs *WeatherFuncArgs
	var funcToolCall *azopenaiassistants.RequiredFunctionToolCall
	{
		// the arguments we need to run the next thing are inside of the run
		submitToolsAction, ok := lastResp.RequiredAction.(*azopenaiassistants.SubmitToolOutputsAction)
		require.True(t, ok)

		tmpFuncToolCall, ok := submitToolsAction.SubmitToolOutputs.ToolCalls[0].(*azopenaiassistants.RequiredFunctionToolCall)
		require.True(t, ok)
		require.Equal(t, "get_current_weather", *tmpFuncToolCall.Function.Name)

		err := json.Unmarshal([]byte(*tmpFuncToolCall.Function.Arguments), &weatherFuncArgs)
		require.NoError(t, err)

		funcToolCall = tmpFuncToolCall
	}

	// now call the "function" (OpenAI is just providing the arguments)
	{
		require.Equal(t, "Boston, MA", weatherFuncArgs.Location)
		require.Equal(t, "celsius", weatherFuncArgs.Unit)

		resp, err := client.SubmitToolOutputsToRunStream(context.Background(), *lastResp.ThreadID, *lastResp.ID, azopenaiassistants.SubmitToolOutputsToRunBody{
			ToolOutputs: []azopenaiassistants.ToolOutput{
				{
					Output:     to.Ptr("0C"),
					ToolCallID: funcToolCall.ID,
				},
			},
		}, nil)
		requireNoErr(t, azure, err)

		defer func() {
			err = resp.Stream.Close()
			require.NoError(t, err)
		}()

		lastResp := processStream(t, azure, *lastResp.AssistantID, streamScenarioTool, resp.Stream)
		require.Nil(t, lastResp)
	}
}

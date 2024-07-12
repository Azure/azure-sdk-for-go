//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var assistantsModel = "gpt-4-1106-preview"

func Test_UsingIdentity(t *testing.T) {
	if os.Getenv("USE_TOKEN_CREDS") != "true" || recording.GetRecordMode() != recording.LiveMode {
		t.Skip("WARNING: Not testing token credentials")
	}

	testFn := func(t *testing.T) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure:       true,
				UseIdentity: true,
			},
		})

		found := false

		pager := client.NewListAssistantsPager(&azopenaiassistants.ListAssistantsOptions{
			Limit: to.Ptr(int32(100)),
		})

		// let's find our assistant in the list
	PagingLoop:
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			requireNoErr(t, true, err)

			for _, a := range page.Data {
				name := "<none>"

				if a.Name != nil {
					name = *a.Name
				}

				fmt.Printf("[%s] %s\n", *a.ID, name)

				if *a.ID == *createResp.ID {
					found = true
					break PagingLoop
				}
			}
		}

		require.True(t, found)
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t)
	})
}

func TestAssistantCreationAndListing(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		found := false

		pager := client.NewListAssistantsPager(&azopenaiassistants.ListAssistantsOptions{
			Limit: to.Ptr(int32(100)),
		})

		// let's find our assistant in the list
	PagingLoop:
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			requireNoErr(t, azure, err)

			for _, a := range page.Data {
				name := "<none>"

				if a.Name != nil {
					name = *a.Name
				}

				fmt.Printf("[%s] %s\n", *a.ID, name)

				if *a.ID == *createResp.ID {
					found = true
					break PagingLoop
				}
			}
		}

		require.True(t, found)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestAssistantMessages(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	testFn := func(t *testing.T, azure bool) {
		client := mustGetClient(t, newClientArgs{
			Azure: azure,
		})

		threadResp, err := client.CreateThread(context.Background(), azopenaiassistants.CreateThreadBody{}, nil)
		requireNoErr(t, azure, err)

		defer func() {
			_, err := client.DeleteThread(context.Background(), *threadResp.ID, nil)
			requireNoErr(t, azure, err)
		}()

		threadID := threadResp.ID

		uploadResp, err := client.UploadFile(context.Background(), bytes.NewReader([]byte("hello world")), azopenaiassistants.FilePurposeAssistants, &azopenaiassistants.UploadFileOptions{
			Filename: getFileName(t, "txt"),
		})
		requireNoErr(t, azure, err)

		defer func() {
			_, err := client.DeleteFile(context.Background(), *uploadResp.ID, nil)
			requireNoErr(t, azure, err)
		}()

		messageResp, err := client.CreateMessage(context.Background(), *threadID, azopenaiassistants.CreateMessageBody{
			Content: to.Ptr("How many ears does a dog usually have?"),
			Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
			Attachments: []azopenaiassistants.MessageAttachment{
				{
					FileID: uploadResp.ID,
					Tools: []azopenaiassistants.MessageAttachmentToolDefinition{
						{CodeInterpreterToolDefinition: &azopenaiassistants.CodeInterpreterToolDefinition{}},
						{FileSearchToolDefinition: &azopenaiassistants.FileSearchToolDefinition{}},
					},
				},
			},
		}, nil)
		requireNoErr(t, azure, err)

		attachmentTools := messageResp.Attachments[0].Tools

		// just trying to keep a consistent ordering of the tools for our checks.
		if attachmentTools[0].CodeInterpreterToolDefinition == nil {
			attachmentTools[0], attachmentTools[1] = attachmentTools[1], attachmentTools[0]
		}

		require.Equal(t, azopenaiassistants.CodeInterpreterToolDefinition{
			Type: to.Ptr("code_interpreter"),
		}, *attachmentTools[0].CodeInterpreterToolDefinition)
		require.Nil(t, attachmentTools[0].FileSearchToolDefinition)

		require.Equal(t, azopenaiassistants.FileSearchToolDefinition{
			Type: to.Ptr("file_search"),
		}, *attachmentTools[1].FileSearchToolDefinition)
		require.Nil(t, attachmentTools[1].CodeInterpreterToolDefinition)

		messageID := messageResp.ID

		getMessageResp, err := client.GetMessage(context.Background(), *threadID, *messageID, nil)
		requireNoErr(t, azure, err)

		require.Equal(t, "How many ears does a dog usually have?", *getMessageResp.Content[0].(*azopenaiassistants.MessageTextContent).Text.Value)

		getMessageFileResp, err := client.GetFile(context.Background(), *uploadResp.ID, nil)
		requireNoErr(t, azure, err)

		require.Equal(t, "file", *getMessageFileResp.Object)

		{
			listFilesResp, err := client.ListFiles(context.Background(), &azopenaiassistants.ListFilesOptions{
				Purpose: to.Ptr(azopenaiassistants.FilePurposeAssistants),
			})
			requireNoErr(t, azure, err)

			found := false

			for _, file := range listFilesResp.Data {
				if *file.ID == *getMessageFileResp.ID {
					found = true
					break
				}
			}

			require.True(t, found)
		}
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func skipRecordingsCantMatchRoutesTestHack(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("skipping due to issue where recordings never match. Issue #22839. Also #22869")
	}
}

func TestAssistantConversationLoop(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createAssistantResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.CreateThreadBody{}, nil)
		requireNoErr(t, azure, err)

		t.Cleanup(func() {
			_, err := client.DeleteThread(context.Background(), *createThreadResp.ID, nil)
			requireNoErr(t, azure, err)
		})

		threadID := *createThreadResp.ID

		convoIdx := 0
		responses := []string{
			"What is the y-intercept for y=x+4?",
			"That answer was nice, thank you. Was my question clear?",
		}

		var convo func(threadMessages []azopenaiassistants.ThreadMessage) []string = func(threadMessages []azopenaiassistants.ThreadMessage) []string {
			// we have a few scripted interactions, just to test how the run loop works.
			defer func() { convoIdx++ }()

			for tmIndex, tm := range threadMessages {
				for contentIndex, content := range tm.Content {
					t.Logf("[ASSISTANT:%d,%d] %s", tmIndex, contentIndex, stringize(content))
				}
			}

			if convoIdx >= len(responses) {
				return nil
			}

			return []string{responses[convoIdx]}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		runAssistant := func(ctx context.Context) {
			var lastResponses []azopenaiassistants.ThreadMessage
			var lastMessageID *string

			for {
				convoResponses := convo(lastResponses)

				if convoResponses == nil {
					break
				}

				for _, msg := range convoResponses {
					// now, let's actually ask it some questions
					createMessageResp, err := client.CreateMessage(ctx, threadID, azopenaiassistants.CreateMessageBody{
						Content: &msg,
						Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
					}, nil)
					requireNoErr(t, azure, err)
					require.NotEmpty(t, createMessageResp)

					t.Logf("[ME] %s", msg)

					lastMessageID = createMessageResp.ID
				}

				// run the thread
				createRunResp, err := client.CreateRun(context.Background(), *createThreadResp.ID, azopenaiassistants.CreateRunBody{
					AssistantID:  createAssistantResp.ID,
					Instructions: to.Ptr("This user is known to be sad, please be kind"),
					// (getting an error with Azure OpenAI on this one)
					// {
					// 	"error": {
					// 	  "message": "1 validation error for Request\nbody -> additional_instructions\n  extra fields not permitted (type=value_error.extra)",
					// 	  "type": "invalid_request_error",
					// 	  "param": null,
					// 	  "code": null
					// 	}
					//   }
					//AdditionalInstructions:
				}, nil)
				requireNoErr(t, azure, err)

				runID := *createRunResp.ID

				lastGetRunResp := pollForTests(t, context.Background(), client, threadID, runID, azure)
				require.Equal(t, azopenaiassistants.RunStatusCompleted, *lastGetRunResp.Status)

				// grab any messages that occurred after our last known message
				lastResponses, err = getLatestMessages(context.Background(), client, threadID, lastMessageID)
				requireNoErr(t, azure, err)
			}
		}

		runAssistant(ctx)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestAssistantRequiredAction(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createAssistantResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		// create our thread and a run for that thread with our question "What's the weather like in Boston, MA, in celsius?"
		var lastMessageID, runID, threadID string
		{
			createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.CreateThreadBody{}, nil)
			requireNoErr(t, azure, err)

			t.Cleanup(func() {
				_, err := client.DeleteThread(context.Background(), *createThreadResp.ID, nil)
				requireNoErr(t, azure, err)
			})

			threadID = *createThreadResp.ID

			msgResponse, err := client.CreateMessage(context.Background(), threadID, azopenaiassistants.CreateMessageBody{
				Content: to.Ptr("What's the weather like in Boston, MA, in celsius?"),
				Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
			}, nil)
			requireNoErr(t, azure, err)

			lastMessageID = *msgResponse.ID

			// run the thread
			createRunResp, err := client.CreateRun(context.Background(), *createThreadResp.ID, azopenaiassistants.CreateRunBody{
				AssistantID:  createAssistantResp.ID,
				Instructions: to.Ptr("Use functions to answer questions, when possible."),
			}, nil)
			requireNoErr(t, azure, err)

			runID = *createRunResp.ID
		}

		lastResp := pollForTests(t, context.Background(), client, threadID, runID, azure)
		require.Equal(t, azopenaiassistants.RunStatusRequiresAction, *lastResp.Status, "More action would be required since we need to run the action and feed it's inputs back in")

		// extract the tool we need to call and the arguments we need to call it with
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

			// submit our outputs from evaluating the tool
			submitToolOutputResp, err := client.SubmitToolOutputsToRun(context.Background(), threadID, runID, azopenaiassistants.SubmitToolOutputsToRunBody{
				ToolOutputs: []azopenaiassistants.ToolOutput{
					{
						Output:     to.Ptr("0C"),
						ToolCallID: funcToolCall.ID,
					},
				},
			}, nil)
			requireNoErr(t, azure, err)
			require.NotEmpty(t, submitToolOutputResp)
		}

		// the run will restart now, we just need to wait until it finishes
		lastResp = pollForTests(t, context.Background(), client, threadID, runID, azure)

		// note our status is Completed now, instead of RunStatusRequiresAction
		require.Equal(t, azopenaiassistants.RunStatusCompleted, *lastResp.Status, "Run should complete now that we've submitted tool outputs")

		latestMessages, err := getLatestMessages(context.Background(), client, threadID, &lastMessageID)
		requireNoErr(t, azure, err)

		require.NotEmpty(t, latestMessages)

		// Prints something like:
		// [ASSISTANT] <id>: Text response: The current weather in Boston, MA, measured in Celsius, is 0°C.
		err = printAssistantMessages(context.Background(), client, latestMessages)
		require.NoError(t, err)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestNewListRunsPager(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		assistantID := *createResp.ID

		runs := map[string]bool{}

		threadAndRunResp, err := client.CreateThreadAndRun(context.Background(), azopenaiassistants.CreateAndRunThreadBody{
			AssistantID:    &assistantID,
			Instructions:   to.Ptr("You're a helpful assistant, but refuse to speak about cats"),
			DeploymentName: &assistantsModel,
			Thread: &azopenaiassistants.CreateThreadBody{
				Messages: []azopenaiassistants.CreateMessageBody{
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("How many ears do cats have?")},
				},
			},
		}, nil)
		requireNoErr(t, azure, err)

		runID := *threadAndRunResp.ID
		threadID := *threadAndRunResp.ThreadID
		runs[runID] = true

		lastRun := pollForTests(t, context.Background(), client, threadID, runID, azure)
		require.Equal(t, azopenaiassistants.RunStatusCompleted, *lastRun.Status)

		run2Resp, err := client.CreateRun(context.Background(), threadID, azopenaiassistants.CreateRunBody{
			AssistantID: &assistantID,
		}, nil)
		requireNoErr(t, azure, err)
		runs[*run2Resp.ID] = true

		lastRun = pollForTests(t, context.Background(), client, threadID, runID, azure)
		require.Equal(t, azopenaiassistants.RunStatusCompleted, *lastRun.Status)

		pager := client.NewListRunsPager(threadID, &azopenaiassistants.ListRunsOptions{
			Limit: to.Ptr[int32](1),
		})

		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, run := range page.Data {
				require.True(t, runs[*run.ID])
				delete(runs, *run.ID)
			}
		}

		require.Empty(t, runs)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestNewListRunStepsPager(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createAsstResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		createThreadAndRunResp, err := client.CreateThreadAndRun(context.Background(), azopenaiassistants.CreateAndRunThreadBody{
			AssistantID:    createAsstResp.ID,
			DeploymentName: &assistantsModel,
			Instructions:   to.Ptr("You are a mysterious assistant"),
			Thread: &azopenaiassistants.CreateThreadBody{
				Messages: []azopenaiassistants.CreateMessageBody{
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("First, message A")},
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("Next, message B")},
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("And then, message C")},
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("And lastly, message D")},
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("What should come next?")},
				},
			},
		}, nil)
		requireNoErr(t, azure, err)

		threadID := *createThreadAndRunResp.ThreadID
		runID := *createThreadAndRunResp.ID

		lastRun := pollForTests(t, context.Background(), client, threadID, runID, azure)
		require.Equal(t, azopenaiassistants.RunStatusCompleted, *lastRun.Status)

		// Run steps are described here:
		// https://platform.openai.com/docs/assistants/how-it-works/runs-and-run-steps
		//
		// The gist is that each time something is added to the thread from the assistant or a tool you
		// get a run step. In this particular run I'm seeing messages indicating the assistant is attempting
		// to answer the question.
		pager := client.NewListRunStepsPager(threadID, runID, nil)

		gotResponse := false

		for pager.More() {
			page, err := pager.NextPage(context.Background())
			requireNoErr(t, azure, err)

			for _, runStep := range page.Data {
				require.Equal(t, azopenaiassistants.RunStepStatusCompleted, *runStep.Status)

				// little sanity check - yes, we can re-read the same step with the ID.
				rereadRunStep, err := client.GetRunStep(context.Background(), threadID, runID, *runStep.ID, nil)
				require.NoError(t, err)
				require.Equal(t, *runStep.ID, *rereadRunStep.ID)

				stepDetails := runStep.StepDetails.(*azopenaiassistants.RunStepMessageCreationDetails)
				messageResp, err := client.GetMessage(context.Background(), threadID, *stepDetails.MessageCreation.MessageID, nil)
				require.NoError(t, err)

				if *messageResp.Role == azopenaiassistants.MessageRoleAssistant {
					body := *messageResp.Content[0].(*azopenaiassistants.MessageTextContent).Text.Value
					fmt.Printf("Assistant response: %s\n", body)
					gotResponse = true
				}
			}
		}

		require.True(t, gotResponse)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestFiles(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	testFn := func(t *testing.T, azure bool) {
		client := mustGetClient(t, newClientArgs{
			Azure: azure,
		})

		textBytes := []byte("test text")
		expectedLen := int32(len(textBytes))

		uploadResp, err := client.UploadFile(context.Background(), bytes.NewReader(textBytes), azopenaiassistants.FilePurposeAssistants, &azopenaiassistants.UploadFileOptions{
			Filename: getFileName(t, "txt"),
		})
		requireNoErr(t, azure, err)
		require.Equal(t, expectedLen, *uploadResp.Bytes)

		defer func() {
			_, err := client.DeleteFile(context.Background(), *uploadResp.ID, nil)
			requireNoErr(t, azure, err)
		}()

		getFileResp, err := client.GetFile(context.Background(), *uploadResp.ID, nil)
		requireNoErr(t, azure, err)

		require.Equal(t, expectedLen, *getFileResp.Bytes)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func pollForTests(t *testing.T, ctx context.Context, client *azopenaiassistants.Client, threadID string, runID string, azure bool) azopenaiassistants.GetRunResponse {
	resp, err := pollUntilRunEnds(ctx, client, threadID, runID)
	requireSuccessfulPolling(t, azure, resp, err)
	return resp
}

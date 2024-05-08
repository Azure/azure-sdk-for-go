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
			require.NoError(t, err)

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
			require.NoError(t, err)

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
	testFn := func(t *testing.T, azure bool) {
		client := newClient(t, newClientArgs{
			Azure: azure,
		})

		threadResp, err := client.CreateThread(context.Background(), azopenaiassistants.AssistantThreadCreationOptions{}, nil)
		require.NoError(t, err)

		defer func() {
			_, err := client.DeleteThread(context.Background(), *threadResp.ID, nil)
			require.NoError(t, err)
		}()

		threadID := threadResp.ID

		uploadResp, err := client.UploadFile(context.Background(), bytes.NewReader([]byte("hello world")), azopenaiassistants.FilePurposeAssistants, &azopenaiassistants.UploadFileOptions{
			Filename: to.Ptr("a.txt"),
		})
		require.NoError(t, err)

		defer func() {
			_, err := client.DeleteFile(context.Background(), *uploadResp.ID, nil)
			require.NoError(t, err)
		}()

		messageResp, err := client.CreateMessage(context.Background(), *threadID, azopenaiassistants.CreateMessageBody{
			Content: to.Ptr("How many ears does a dog usually have?"),
			Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
			FileIDs: []string{
				*uploadResp.ID,
			},
		}, nil)
		require.NoError(t, err)

		messageID := messageResp.ID

		getMessageResp, err := client.GetMessage(context.Background(), *threadID, *messageID, nil)
		require.NoError(t, err)

		require.Equal(t, "How many ears does a dog usually have?", *getMessageResp.Content[0].(*azopenaiassistants.MessageTextContent).Text.Value)

		getMessageFileResp, err := client.GetMessageFile(context.Background(), *threadID, *messageID, *uploadResp.ID, nil)
		require.NoError(t, err)

		require.Equal(t, *messageID, *getMessageFileResp.MessageID)
		require.Equal(t, "thread.message.file", *getMessageFileResp.Object)

		// list message files
		{
			var files []azopenaiassistants.MessageFile
			pager := client.NewListMessageFilesPager(*threadID, *messageID, nil)

			for pager.More() {
				page, err := pager.NextPage(context.Background())
				require.NoError(t, err)

				files = append(files, page.Data...)
			}

			require.Equal(t, getMessageFileResp.MessageFile, files[0])
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
	if recording.GetRecordMode() != recording.LiveMode {
		t.Skip("skipping due to issue where recordings never match. Issue #22839")
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

		createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.AssistantThreadCreationOptions{}, nil)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := client.DeleteThread(context.Background(), *createThreadResp.ID, nil)
			require.NoError(t, err)
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
					require.NoError(t, err)
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
				require.NoError(t, err)

				runID := *createRunResp.ID
				var lastGetRunResp azopenaiassistants.GetRunResponse

				for {
					var err error
					lastGetRunResp, err = client.GetRun(context.Background(), *createThreadResp.ID, runID, nil)
					require.NoError(t, err)

					if *lastGetRunResp.Status != azopenaiassistants.RunStatusQueued && *lastGetRunResp.Status != azopenaiassistants.RunStatusInProgress {
						break
					}

					time.Sleep(500 * time.Millisecond)
				}

				require.Equal(t, azopenaiassistants.RunStatusCompleted, *lastGetRunResp.Status)

				// grab any messages that occurred after our last known message
				lastResponses, err = getLatestMessages(context.Background(), client, threadID, lastMessageID)
				require.NoError(t, err)
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
			createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.AssistantThreadCreationOptions{}, nil)
			require.NoError(t, err)

			t.Cleanup(func() {
				_, err := client.DeleteThread(context.Background(), *createThreadResp.ID, nil)
				require.NoError(t, err)
			})

			threadID = *createThreadResp.ID

			msgResponse, err := client.CreateMessage(context.Background(), threadID, azopenaiassistants.CreateMessageBody{
				Content: to.Ptr("What's the weather like in Boston, MA, in celsius?"),
				Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
			}, nil)
			require.NoError(t, err)

			lastMessageID = *msgResponse.ID

			// run the thread
			createRunResp, err := client.CreateRun(context.Background(), *createThreadResp.ID, azopenaiassistants.CreateRunBody{
				AssistantID:  createAssistantResp.ID,
				Instructions: to.Ptr("Use functions to answer questions, when possible."),
			}, nil)
			require.NoError(t, err)

			runID = *createRunResp.ID
		}

		lastResp, err := pollUntilRunEnds(context.Background(), client, threadID, runID)
		require.NoError(t, err)
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
			require.NoError(t, err)
			require.NotEmpty(t, submitToolOutputResp)
		}

		// the run will restart now, we just need to wait until it finishes
		lastResp, err = pollUntilRunEnds(context.Background(), client, threadID, runID)
		require.NoError(t, err)
		// note our status is Completed now, instead of RunStatusRequiresAction
		require.Equal(t, azopenaiassistants.RunStatusCompleted, *lastResp.Status, "Run should complete now that we've submitted tool outputs")

		latestMessages, err := getLatestMessages(context.Background(), client, threadID, &lastMessageID)
		require.NoError(t, err)

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

func TestNewAssistantFilesPager(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		var createdIDs []string

		for i := 0; i < 5; i++ {
			createAsstFileResp, err := client.CreateAssistantFile(context.Background(), *createResp.ID, azopenaiassistants.CreateAssistantFileBody{
				FileID: mustUploadFile(t, client, "hello world").ID,
			}, nil)
			require.NoError(t, err)
			require.NotEmpty(t, createAsstFileResp)

			createdIDs = append(createdIDs, *createAsstFileResp.ID)
		}

		for _, sortOrder := range []azopenaiassistants.ListSortOrder{azopenaiassistants.ListSortOrderAscending, azopenaiassistants.ListSortOrderDescending} {
			t.Run("with sort order "+string(sortOrder), func(t *testing.T) {
				m := map[string]bool{}

				var first *azopenaiassistants.AssistantFile
				var last *azopenaiassistants.AssistantFile

				for _, id := range createdIDs {
					m[id] = true
				}

				pager := client.NewListAssistantFilesPager(*createResp.ID, &azopenaiassistants.ListAssistantFilesOptions{
					Limit: to.Ptr[int32](1),
					Order: &sortOrder,
				})

				for pager.More() {
					page, err := pager.NextPage(context.Background())
					require.NoError(t, err)

					for _, item := range page.Data {
						require.Contains(t, m, *item.ID)
						delete(m, *item.ID) // catch if we got the same file twice somehow.

						if first == nil {
							first = &item
						}

						last = &item
					}
				}

				require.Empty(t, m)

				if sortOrder == azopenaiassistants.ListSortOrderAscending {
					require.Greater(t, last.CreatedAt.Sub(*first.CreatedAt), time.Duration(0))
				} else {
					require.Greater(t, first.CreatedAt.Sub(*last.CreatedAt), time.Duration(0))
				}
			})
		}
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

		threadAndRunResp, err := client.CreateThreadAndRun(context.Background(), azopenaiassistants.CreateAndRunThreadOptions{
			AssistantID:    &assistantID,
			Instructions:   to.Ptr("You're a helpful assistant, but refuse to speak about cats"),
			DeploymentName: &assistantsModel,
			Thread: &azopenaiassistants.AssistantThreadCreationOptions{
				Messages: []azopenaiassistants.ThreadInitializationMessage{
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("How many ears do cats have?")},
				},
			},
		}, nil)
		require.NoError(t, err)

		runID := *threadAndRunResp.ID
		threadID := *threadAndRunResp.ThreadID
		runs[runID] = true

		lastRun, err := pollUntilRunEnds(context.Background(), client, threadID, runID)
		require.NoError(t, err)
		require.Equal(t, azopenaiassistants.RunStatusCompleted, *lastRun.Status)

		run2Resp, err := client.CreateRun(context.Background(), threadID, azopenaiassistants.CreateRunBody{
			AssistantID: &assistantID,
		}, nil)
		require.NoError(t, err)
		runs[*run2Resp.ID] = true

		lastRun, err = pollUntilRunEnds(context.Background(), client, threadID, runID)
		require.NoError(t, err)
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

		createThreadAndRunResp, err := client.CreateThreadAndRun(context.Background(), azopenaiassistants.CreateAndRunThreadOptions{
			AssistantID:    createAsstResp.ID,
			DeploymentName: &assistantsModel,
			Instructions:   to.Ptr("You are a mysterious assistant"),
			Thread: &azopenaiassistants.AssistantThreadCreationOptions{
				Messages: []azopenaiassistants.ThreadInitializationMessage{
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("First, message A")},
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("Next, message B")},
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("And then, message C")},
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("And lastly, message D")},
					{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("What should come next?")},
				},
			},
		}, nil)
		require.NoError(t, err)

		threadID := *createThreadAndRunResp.ThreadID
		runID := *createThreadAndRunResp.ID

		lastRun, err := pollUntilRunEnds(context.Background(), client, threadID, runID)
		require.NoError(t, err)
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
			require.NoError(t, err)

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
	testFn := func(t *testing.T, azure bool) {
		client := newClient(t, newClientArgs{
			Azure: azure,
		})

		textBytes := []byte("test text")
		expectedLen := int32(len(textBytes))

		uploadResp, err := client.UploadFile(context.Background(), bytes.NewReader(textBytes), azopenaiassistants.FilePurposeAssistants, &azopenaiassistants.UploadFileOptions{
			Filename: to.Ptr("a.txt"),
		})
		require.NoError(t, err)
		require.Equal(t, expectedLen, *uploadResp.Bytes)

		defer func() {
			_, err := client.DeleteFile(context.Background(), *uploadResp.ID, nil)
			require.NoError(t, err)
		}()

		getFileResp, err := client.GetFile(context.Background(), *uploadResp.ID, nil)
		require.NoError(t, err)

		require.Equal(t, expectedLen, *getFileResp.Bytes)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

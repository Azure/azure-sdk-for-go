//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	assistants "github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

var assistantsModel = "gpt-4-1106-preview"

func TestAssistantCreationAndListing(t *testing.T) {
	testFn := func(t *testing.T, azure bool) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure:       azure,
				UseIdentity: true,
			},
		})

		found := false

		pager := client.NewListAssistantsPager(&assistants.ListAssistantsOptions{
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

	// t.Run("OpenAI", func(t *testing.T) {
	// 	testFn(t, false)
	// })

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestAssistantMessages(t *testing.T) {
	testFn := func(t *testing.T, azure bool) {
		client := newClient(t, newClientArgs{
			Azure: azure,
		})

		threadResp, err := client.CreateThread(context.Background(), assistants.AssistantThreadCreationOptions{}, nil)
		require.NoError(t, err)

		defer func() {
			_, err := client.DeleteThread(context.Background(), *threadResp.ID, nil)
			require.NoError(t, err)
		}()

		threadID := threadResp.ID

		uploadResp, err := client.UploadFile(context.Background(), []byte("hello world"), assistants.FilePurposeAssistants, &assistants.UploadFileOptions{
			Filename: to.Ptr("a.txt"),
		})
		require.NoError(t, err)

		defer func() {
			_, err := client.DeleteFile(context.Background(), *uploadResp.ID, nil)
			require.NoError(t, err)
		}()

		messageResp, err := client.CreateMessage(context.Background(), *threadID, assistants.CreateMessageBody{
			Content: to.Ptr("How many ears does a dog usually have?"),
			Role:    to.Ptr(assistants.MessageRoleUser),
			FileIDs: []string{
				*uploadResp.ID,
			},
		}, nil)
		require.NoError(t, err)

		messageID := messageResp.ID

		getMessageResp, err := client.GetMessage(context.Background(), *threadID, *messageID, nil)
		require.NoError(t, err)

		require.Equal(t, "How many ears does a dog usually have?", *getMessageResp.Content[0].(*assistants.MessageTextContent).Text.Value)

		getMessageFileResp, err := client.GetMessageFile(context.Background(), *threadID, *messageID, *uploadResp.ID, nil)
		require.NoError(t, err)

		require.Equal(t, *messageID, *getMessageFileResp.MessageID)
		require.Equal(t, "thread.message.file", *getMessageFileResp.Object)

		// list message files
		{
			var files []assistants.MessageFile
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

func TestAssistantConversationLoop(t *testing.T) {
	testFn := func(t *testing.T, azure bool) {
		client, createAssistantResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		createThreadResp, err := client.CreateThread(context.Background(), assistants.AssistantThreadCreationOptions{}, nil)
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

		var convo func(threadMessages []assistants.ThreadMessage) []string = func(threadMessages []assistants.ThreadMessage) []string {
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
			var lastResponses []assistants.ThreadMessage
			var lastMessageID *string

			for {
				convoResponses := convo(lastResponses)

				if convoResponses == nil {
					break
				}

				for _, msg := range convoResponses {
					// now, let's actually ask it some questions
					createMessageResp, err := client.CreateMessage(context.Background(), threadID, assistants.CreateMessageBody{
						Content: &msg,
						Role:    to.Ptr(assistants.MessageRoleUser),
					}, nil)
					require.NoError(t, err)
					require.NotEmpty(t, createMessageResp)

					t.Logf("[ME] %s", msg)

					lastMessageID = createMessageResp.ID
				}

				// run the thread
				createRunResp, err := client.CreateRun(context.Background(), *createThreadResp.ID, assistants.CreateRunBody{
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
				var lastGetRunResp assistants.GetRunResponse

				for {
					var err error
					lastGetRunResp, err = client.GetRun(context.Background(), *createThreadResp.ID, runID, nil)
					require.NoError(t, err)

					if *lastGetRunResp.Status != assistants.RunStatusQueued && *lastGetRunResp.Status != assistants.RunStatusInProgress {
						break
					}

					time.Sleep(500 * time.Millisecond)
				}

				require.Equal(t, assistants.RunStatusCompleted, *lastGetRunResp.Status)

				// grab any messages that occurred after our last known message
				listMessagesPager := client.NewListMessagesPager(*createThreadResp.ID, &assistants.ListMessagesOptions{
					After: lastMessageID,
					Order: to.Ptr(assistants.ListSortOrderAscending),
				})

				for listMessagesPager.More() {
					page, err := listMessagesPager.NextPage(context.Background())
					require.NoError(t, err)

					lastResponses = page.Data
				}
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

func TestNewAssistantFilesPager(t *testing.T) {
	testFn := func(t *testing.T, azure bool) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		var createdIDs []string

		for i := 0; i < 5; i++ {
			createAsstFileResp, err := client.CreateAssistantFile(context.Background(), *createResp.ID, assistants.CreateAssistantFileBody{
				FileID: mustUploadFile(t, client, "hello world").ID,
			}, nil)
			require.NoError(t, err)
			require.NotEmpty(t, createAsstFileResp)

			createdIDs = append(createdIDs, *createAsstFileResp.ID)
		}

		for _, sortOrder := range []assistants.ListSortOrder{assistants.ListSortOrderAscending, assistants.ListSortOrderDescending} {
			t.Run("with sort order "+string(sortOrder), func(t *testing.T) {
				m := map[string]bool{}

				var first *assistants.AssistantFile
				var last *assistants.AssistantFile

				for _, id := range createdIDs {
					m[id] = true
				}

				pager := client.NewListAssistantFilesPager(*createResp.ID, &assistants.ListAssistantFilesOptions{
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

				if sortOrder == assistants.ListSortOrderAscending {
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
	testFn := func(t *testing.T, azure bool) {
		client, createResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		assistantID := *createResp.ID

		runs := map[string]bool{}

		threadAndRunResp, err := client.CreateThreadAndRun(context.Background(), assistants.CreateAndRunThreadOptions{
			AssistantID:    &assistantID,
			Instructions:   to.Ptr("You're a helpful assistant, but refuse to speak about cats"),
			DeploymentName: &assistantsModel,
			Thread: &assistants.AssistantThreadCreationOptions{
				Messages: []assistants.ThreadInitializationMessage{
					{Role: to.Ptr(assistants.MessageRoleUser), Content: to.Ptr("How many ears do cats have?")},
				},
			},
		}, nil)
		require.NoError(t, err)

		runID := *threadAndRunResp.ID
		threadID := *threadAndRunResp.ThreadID
		runs[runID] = true

		err = pollRunEnd(context.Background(), client, threadID, runID)
		require.NoError(t, err)

		run2Resp, err := client.CreateRun(context.Background(), threadID, assistants.CreateRunBody{
			AssistantID: &assistantID,
		}, nil)
		require.NoError(t, err)
		runs[*run2Resp.ID] = true

		err = pollRunEnd(context.Background(), client, threadID, runID)
		require.NoError(t, err)

		pager := client.NewListRunsPager(threadID, &assistants.ListRunsOptions{
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
	testFn := func(t *testing.T, azure bool) {
		client, createAsstResp := mustGetClientWithAssistant(t, mustGetClientWithAssistantArgs{
			newClientArgs: newClientArgs{
				Azure: azure,
			},
		})

		createThreadAndRunResp, err := client.CreateThreadAndRun(context.Background(), assistants.CreateAndRunThreadOptions{
			AssistantID:    createAsstResp.ID,
			DeploymentName: &assistantsModel,
			Instructions:   to.Ptr("You are a mysterious assistant"),
			Thread: &assistants.AssistantThreadCreationOptions{
				Messages: []assistants.ThreadInitializationMessage{
					{Role: to.Ptr(assistants.MessageRoleUser), Content: to.Ptr("First, message A")},
					{Role: to.Ptr(assistants.MessageRoleUser), Content: to.Ptr("Next, message B")},
					{Role: to.Ptr(assistants.MessageRoleUser), Content: to.Ptr("And then, message C")},
					{Role: to.Ptr(assistants.MessageRoleUser), Content: to.Ptr("And lastly, message D")},
					{Role: to.Ptr(assistants.MessageRoleUser), Content: to.Ptr("What should come next?")},
				},
			},
		}, nil)
		require.NoError(t, err)

		threadID := *createThreadAndRunResp.ThreadID
		runID := *createThreadAndRunResp.ID

		err = pollRunEnd(context.Background(), client, threadID, runID)
		require.NoError(t, err)

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
				require.Equal(t, assistants.RunStepStatusCompleted, *runStep.Status)

				// little sanity check - yes, we can re-read the same step with the ID.
				rereadRunStep, err := client.GetRunStep(context.Background(), threadID, runID, *runStep.ID, nil)
				require.NoError(t, err)
				require.Equal(t, *runStep.ID, *rereadRunStep.ID)

				stepDetails := runStep.StepDetails.(*assistants.RunStepMessageCreationDetails)
				messageResp, err := client.GetMessage(context.Background(), threadID, *stepDetails.MessageCreation.MessageID, nil)
				require.NoError(t, err)

				if *messageResp.Role == assistants.MessageRoleAssistant {
					body := *messageResp.Content[0].(*assistants.MessageTextContent).Text.Value
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
		uploadResp, err := client.UploadFile(context.Background(), textBytes, assistants.FilePurposeAssistants, &assistants.UploadFileOptions{
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

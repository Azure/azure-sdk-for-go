//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/stretchr/testify/require"
)

type streamScenario string

const (
	streamScenarioThreadAndRun streamScenario = "threadAndRun"
	streamScenarioRun          streamScenario = "run"
	streamScenarioTool         streamScenario = "tool"
)

// processStream checks that all the streaming events arrive and work the way you'd expect.
//
//   - threadAlreadyCreated - if you started your stream from CreateRun() you won't see the "thread created"
//     event because _you_ already created the thread. This only fires off in the all-in-one CreateThreadAndRun.
func processStream(t *testing.T, azure bool, asstID string, scenario streamScenario, stream *azopenaiassistants.EventReader[azopenaiassistants.StreamEvent]) *azopenaiassistants.ThreadRun {
	type threadStats struct {
		MessageContent string
		EventCounts    map[azopenaiassistants.AssistantStreamEvent]int
	}

	thread := threadStats{
		EventCounts: map[azopenaiassistants.AssistantStreamEvent]int{},
	}

	t.Logf("(%s) Stream has started", scenario)

	var toolOutputsRequiredRet *azopenaiassistants.ThreadRun

	for {
		event, err := stream.Read()

		if errors.Is(err, io.EOF) {
			t.Logf("(%s) Stream has ended normally", scenario)
			break
		}

		require.NoError(t, err)

		thread.EventCounts[event.Reason]++

		switch event.Reason {
		case azopenaiassistants.AssistantStreamEventThreadCreated:
			v := requireType[*azopenaiassistants.AssistantThread](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)

		case azopenaiassistants.AssistantStreamEventThreadRunCreated:
			v := requireType[*azopenaiassistants.ThreadRun](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)
			require.Equal(t, asstID, *v.AssistantID)

		case azopenaiassistants.AssistantStreamEventThreadRunQueued:
			v := requireType[*azopenaiassistants.ThreadRun](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)
			require.Equal(t, asstID, *v.AssistantID)
			require.Equal(t, azopenaiassistants.RunStatusQueued, *v.Status)

		case azopenaiassistants.AssistantStreamEventThreadRunInProgress:
			v := requireType[*azopenaiassistants.ThreadRun](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)
			require.Equal(t, asstID, *v.AssistantID)
			require.Equal(t, azopenaiassistants.RunStatusInProgress, *v.Status)

		case azopenaiassistants.AssistantStreamEventThreadMessageCreated:
			v := requireType[*azopenaiassistants.ThreadMessage](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)
			require.Equal(t, azopenaiassistants.MessageRoleAssistant, *v.Role)

		case azopenaiassistants.AssistantStreamEventThreadMessageCompleted:
			v := requireType[*azopenaiassistants.ThreadMessage](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)
			require.Equal(t, azopenaiassistants.MessageStatusCompleted, *v.Status)
			require.Equal(t, azopenaiassistants.MessageRoleAssistant, *v.Role)

		case azopenaiassistants.AssistantStreamEventThreadMessageInProgress:
			v := requireType[*azopenaiassistants.ThreadMessage](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)
			require.Equal(t, azopenaiassistants.MessageStatusInProgress, *v.Status)
			require.Equal(t, azopenaiassistants.MessageRoleAssistant, *v.Role)

		case azopenaiassistants.AssistantStreamEventThreadMessageDelta:
			v := requireType[*azopenaiassistants.MessageDeltaChunk](t, event)
			require.NotEmpty(t, *v.Delta)

			for _, c := range v.Delta.Content {
				switch actualContent := c.(type) {
				case *azopenaiassistants.MessageDeltaImageFileContent:
					require.Fail(t, "Answer included image content, which is not expected")
				case *azopenaiassistants.MessageDeltaTextContentObject:
					thread.MessageContent += *actualContent.Text.Value
				}
			}

		case azopenaiassistants.AssistantStreamEventThreadRunCompleted:
			v := requireType[*azopenaiassistants.ThreadRun](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)

			require.Equal(t, asstID, *v.AssistantID)

		case azopenaiassistants.AssistantStreamEventThreadRunFailed:
			v := requireType[*azopenaiassistants.ThreadRun](t, event)

			if azure {
				skipifThrottled(t, v.LastError)
			}

			jsonText, err := json.MarshalIndent(v, "  ", "  ")

			if err != nil {
				require.Failf(t, "ThreadRun failed", "%#v", v)
			} else {
				require.Failf(t, "ThreadRun failed", "%s", string(jsonText))
			}

		case azopenaiassistants.AssistantStreamEventThreadRunRequiresAction:
			v := requireType[*azopenaiassistants.ThreadRun](t, event)
			t.Logf("(%s) %s", event.Reason, *v.ID)
			t.Logf("We need to run tool outputs, stream is going to end.")

			// we'll need this in order to properly submit tool outputs.
			toolOutputsRequiredRet = v

		case azopenaiassistants.AssistantStreamEventThreadRunStepDelta:
			v := requireType[*azopenaiassistants.RunStepDeltaChunk](t, event)
			switch details := v.Delta.StepDetails.(type) {
			case *azopenaiassistants.RunStepDeltaMessageCreation:
				t.Logf("(%s): %s, message created: %s", event.Reason, *v.ID, *details.MessageCreation.MessageID)
			case *azopenaiassistants.RunStepDeltaToolCallObject:
				t.Logf("(%s): %s", event.Reason, *v.ID)

				for _, toolCallClassification := range details.ToolCalls {
					switch tc := toolCallClassification.(type) {
					case *azopenaiassistants.RunStepDeltaFileSearchToolCall:
						t.Logf("  toolCall[%d]: id: %s, type: %s", *tc.Index, *tc.ID, *tc.Type)
					case *azopenaiassistants.RunStepDeltaFunctionToolCall:
						funcName := getValue(tc.Function.Name, "")
						funcArgs := getValue(tc.Function.Arguments, "")
						id := getValue(tc.ID, "(unknown)")

						t.Logf("  toolCall[%d]: id: %s, type: %s, func: %s, args: %s",
							*tc.Index,
							id,
							*tc.Type,
							funcName, funcArgs)
					case *azopenaiassistants.RunStepDeltaToolCall:
						t.Logf("  toolCall[%d]: id: %s, type: %s", *tc.Index, *tc.ID, *tc.Type)
					}
				}
			}

		case azopenaiassistants.AssistantStreamEventThreadRunStepCompleted,
			azopenaiassistants.AssistantStreamEventThreadRunStepCreated,
			azopenaiassistants.AssistantStreamEventThreadRunStepInProgress:
			v := requireType[*azopenaiassistants.RunStep](t, event)
			t.Logf("(%s): %s", event.Reason, *v.ID)
		case azopenaiassistants.AssistantStreamEventError,
			azopenaiassistants.AssistantStreamEventThreadRunStepExpired,
			azopenaiassistants.AssistantStreamEventThreadRunStepFailed,
			azopenaiassistants.AssistantStreamEventThreadRunStepCancelled,
			azopenaiassistants.AssistantStreamEventThreadRunCancelled,
			azopenaiassistants.AssistantStreamEventThreadMessageIncomplete,
			azopenaiassistants.AssistantStreamEventThreadRunCancelling,
			azopenaiassistants.AssistantStreamEventThreadRunExpired:
			require.Failf(t, "Failure", "kind %s should not happen in this test", string(event.Reason))
		case azopenaiassistants.AssistantStreamEventDone:
			// this is handled by the EventReader, causing it to return io.EOF
			require.Failf(t, "Unhandled kind", "kind %s should not happen in this test", string(event.Reason))
		default:
			require.Failf(t, "Unhandled kind", "kind %s should not happen in this test", string(event.Reason))
		}
	}

	switch scenario {
	case streamScenarioThreadAndRun:
		require.Equal(t, 1, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadCreated])
		require.Equal(t, 1, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadRunCreated])
	case streamScenarioRun:
		// thread was already created, so we don't see that event here.
		require.Zero(t, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadCreated])
	case streamScenarioTool:
		require.Equal(t, 0, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadCreated])
		require.Equal(t, 0, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadRunCreated])
	}

	statsJSON, err := json.MarshalIndent(thread, "  ", "  ")
	require.NoError(t, err)

	t.Logf("Thread stats(%s): %s", scenario, statsJSON)

	if toolOutputsRequiredRet == nil {
		// this is a run that has started and completed.
		require.NotZero(t, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadMessageCreated])
		require.NotZero(t, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadMessageDelta])
		require.NotZero(t, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadMessageInProgress])
		require.NotZero(t, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadMessageCompleted])
		require.Equal(t, 1, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadRunCompleted])

		// each type that was created should also be completed.
		require.Equalf(t,
			thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadMessageCompleted],
			thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadMessageCreated], "messages created == messages completed for %s", scenario)
	} else {
		// the run isn't complete, it's "paused" while we wait for tool outputs to be
		// submitted.
		require.Equal(t, 0, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadRunCompleted])
		require.NotZero(t, thread.EventCounts[azopenaiassistants.AssistantStreamEventThreadRunStepInProgress], "a step will still be running since we're waiting for tool outputs")
	}

	return toolOutputsRequiredRet
}

func skipifThrottled(t *testing.T, lastErr *azopenaiassistants.ThreadRunLastError) {
	if lastErr != nil && lastErr.Code != nil && *lastErr.Code == "rate_limit_exceeded" {
		msg := "(no message)"

		if lastErr.Message != nil {
			msg = *lastErr.Message
		}

		t.Skipf("Assistant test is being throttled (%s): %s", *lastErr.Code, msg)
	}
}

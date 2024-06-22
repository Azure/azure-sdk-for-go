//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"encoding/json"
	"fmt"
	"io"
)

// StreamEvent contains an event from an Assistants API stream.
type StreamEvent struct {
	// Reason identifies the reason this event was sent.
	Reason AssistantStreamEvent

	// Event is the payload for the StreamEvent.
	Event StreamEventDataClassification
}

// StreamEventData is the common data for all stream events.
type StreamEventData struct{}

// StreamEventDataClassification provides polymorphic access to related types.
// Call the interface's GetStreamEventDataContent() method to access the common type.
// Use a type switch to determine the concrete type. The possible types are:
// - [*AssistantThread]
// - [*MessageDeltaChunk]
// - [*RunStep]
// - [*RunStepDeltaChunk]
// - [*ThreadMessage]
// - [*ThreadRun]
type StreamEventDataClassification interface {
	// GetStreamEventDataContent returns the StreamEventData content of the underlying type.
	GetStreamEventDataContent() *StreamEventData
}

func unmarshalStreamEvent(eventName string, data []byte) (StreamEvent, error) {
	kind := AssistantStreamEvent(eventName)

	switch kind {
	case AssistantStreamEventDone: // AssistantStreamEventDone - Event sent when the stream is done.
		// This should never get hit - EventReader takes care of processing this event for us.
		return StreamEvent{}, io.EOF
	case AssistantStreamEventError: // AssistantStreamEventError - Event sent when an error occurs, such as an internal server error or a timeout.
		return StreamEvent{}, fmt.Errorf("error occurred while streaming: %s", string(data))
	case AssistantStreamEventThreadCreated:
		var v *AssistantThread
		if err := json.Unmarshal(data, &v); err != nil {
			return StreamEvent{}, err
		}

		return StreamEvent{
			Event:  v,
			Reason: kind,
		}, nil
	case AssistantStreamEventThreadMessageCompleted,
		AssistantStreamEventThreadMessageCreated,
		AssistantStreamEventThreadMessageInProgress,
		AssistantStreamEventThreadMessageIncomplete:
		var v *ThreadMessage
		if err := json.Unmarshal(data, &v); err != nil {
			return StreamEvent{}, err
		}

		return StreamEvent{
			Event:  v,
			Reason: kind,
		}, nil

	case AssistantStreamEventThreadMessageDelta:
		var v *MessageDeltaChunk
		if err := json.Unmarshal(data, &v); err != nil {
			return StreamEvent{}, err
		}

		return StreamEvent{
			Event:  v,
			Reason: kind,
		}, nil

	case AssistantStreamEventThreadRunCancelled,
		AssistantStreamEventThreadRunCancelling,
		AssistantStreamEventThreadRunCompleted,
		AssistantStreamEventThreadRunCreated,
		AssistantStreamEventThreadRunExpired,
		AssistantStreamEventThreadRunFailed,
		AssistantStreamEventThreadRunInProgress,
		AssistantStreamEventThreadRunQueued,
		AssistantStreamEventThreadRunRequiresAction:
		var v *ThreadRun
		if err := json.Unmarshal(data, &v); err != nil {
			return StreamEvent{}, err
		}

		return StreamEvent{
			Event:  v,
			Reason: kind,
		}, nil

	case AssistantStreamEventThreadRunStepCancelled, AssistantStreamEventThreadRunStepCompleted, AssistantStreamEventThreadRunStepCreated, AssistantStreamEventThreadRunStepExpired, AssistantStreamEventThreadRunStepFailed, AssistantStreamEventThreadRunStepInProgress:
		var v *RunStep
		if err := json.Unmarshal(data, &v); err != nil {
			return StreamEvent{}, err
		}

		return StreamEvent{
			Event:  v,
			Reason: kind,
		}, nil

	case AssistantStreamEventThreadRunStepDelta:
		var v *RunStepDeltaChunk
		if err := json.Unmarshal(data, &v); err != nil {
			return StreamEvent{}, err
		}

		return StreamEvent{
			Event:  v,
			Reason: kind,
		}, nil

	default:
		return StreamEvent{}, fmt.Errorf("unhandled kind %s", kind)
	}
}

// GetStreamEventDataContent returns the common data of the underlying type.
func (v *AssistantThread) GetStreamEventDataContent() *StreamEventData {
	return &StreamEventData{}
}

// GetStreamEventDataContent returns the common data of the underlying type.
func (v *MessageDeltaChunk) GetStreamEventDataContent() *StreamEventData {
	return &StreamEventData{}
}

// GetStreamEventDataContent returns the common data of the underlying type.
func (v *RunStep) GetStreamEventDataContent() *StreamEventData {
	return &StreamEventData{}
}

// GetStreamEventDataContent returns the common data of the underlying type.
func (v *RunStepDeltaChunk) GetStreamEventDataContent() *StreamEventData {
	return &StreamEventData{}
}

// GetStreamEventDataContent returns the common data of the underlying type.
func (v *ThreadMessage) GetStreamEventDataContent() *StreamEventData {
	return &StreamEventData{}
}

// GetStreamEventDataContent returns the common data of the underlying type.
func (v *ThreadRun) GetStreamEventDataContent() *StreamEventData {
	return &StreamEventData{}
}

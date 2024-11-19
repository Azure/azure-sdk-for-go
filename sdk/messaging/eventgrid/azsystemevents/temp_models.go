//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

import (
	"encoding/json"
	"fmt"
)

// NOTE: there appears to be a bug where, when a type is used as a base class in TypeSpec, we automatically trim it _but_ there are some
// cases, like these events, where they're not just vestigial, they are legitimate models that should be generated.
// Filed as https://github.com/Azure/autorest.go/issues/1466. Once this is fixed, we can delete these temporarily copied models.

const (
	TypeMediaJobOutputStateChange = "Microsoft.Media.JobOutputStateChange" // maps to MediaJobOutputStateChangeEventData
	TypeMediaJobStateChange       = "Microsoft.Media.JobStateChange"       // maps to MediaJobStateChangeEventData
)

// MediaJobOutputStateChangeEventData - Schema of the Data property of an EventGridEvent for a
// Microsoft.Media.JobOutputStateChange event.
type MediaJobOutputStateChangeEventData struct {
	// REQUIRED; Gets the Job correlation data.
	JobCorrelationData map[string]*string

	// REQUIRED; Gets the output.
	Output MediaJobOutputClassification

	// REQUIRED; The previous state of the Job.
	PreviousState *MediaJobState
}

// MediaJobStateChangeEventData - Schema of the Data property of an EventGridEvent for a
// Microsoft.Media.JobStateChange event.
type MediaJobStateChangeEventData struct {
	// REQUIRED; Gets the Job correlation data.
	CorrelationData map[string]*string

	// REQUIRED; The previous state of the Job.
	PreviousState *MediaJobState

	// REQUIRED; The new state of the Job.
	State *MediaJobState
}

// MarshalJSON implements the json.Marshaller interface for type MediaJobOutputStateChangeEventData.
func (m MediaJobOutputStateChangeEventData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "jobCorrelationData", m.JobCorrelationData)
	populate(objectMap, "output", m.Output)
	populate(objectMap, "previousState", m.PreviousState)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type MediaJobOutputStateChangeEventData.
func (m *MediaJobOutputStateChangeEventData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", m, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "jobCorrelationData":
			err = unpopulate(val, "JobCorrelationData", &m.JobCorrelationData)
			delete(rawMsg, key)
		case "output":
			m.Output, err = unmarshalMediaJobOutputClassification(val)
			delete(rawMsg, key)
		case "previousState":
			err = unpopulate(val, "PreviousState", &m.PreviousState)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", m, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type MediaJobStateChangeEventData.
func (m MediaJobStateChangeEventData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "correlationData", m.CorrelationData)
	populate(objectMap, "previousState", m.PreviousState)
	populate(objectMap, "state", m.State)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type MediaJobStateChangeEventData.
func (m *MediaJobStateChangeEventData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", m, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "correlationData":
			err = unpopulate(val, "CorrelationData", &m.CorrelationData)
			delete(rawMsg, key)
		case "previousState":
			err = unpopulate(val, "PreviousState", &m.PreviousState)
			delete(rawMsg, key)
		case "state":
			err = unpopulate(val, "State", &m.State)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", m, err)
		}
	}
	return nil
}

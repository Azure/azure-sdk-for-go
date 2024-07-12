//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

import (
	"encoding/json"
	"fmt"
	"time"
)

// EventGridEvent - Properties of an event published to an Event Grid topic using the EventGrid Schema.
type EventGridEvent struct {
	// REQUIRED; Event data specific to the event type.
	Data any

	// REQUIRED; The schema version of the data object.
	DataVersion *string

	// REQUIRED; The time (in UTC) the event was generated.
	EventTime *time.Time

	// REQUIRED; The type of the event that occurred.
	EventType *string

	// REQUIRED; An unique identifier for the event.
	ID *string

	// REQUIRED; A resource path relative to the topic path.
	Subject *string

	// The resource path of the event source.
	Topic *string

	// READ-ONLY; The schema version of the event metadata.
	MetadataVersion *string
}

// MarshalJSON implements the json.Marshaller interface for type EventGridEvent.
func (e EventGridEvent) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populateAny(objectMap, "data", e.Data)
	populate(objectMap, "dataVersion", e.DataVersion)
	populateDateTimeRFC3339(objectMap, "eventTime", e.EventTime)
	populate(objectMap, "eventType", e.EventType)
	populate(objectMap, "id", e.ID)
	populate(objectMap, "metadataVersion", e.MetadataVersion)
	populate(objectMap, "subject", e.Subject)
	populate(objectMap, "topic", e.Topic)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type EventGridEvent.
func (e *EventGridEvent) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %w", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "data":
			e.Data = []byte(val)
			delete(rawMsg, key)
		case "dataVersion":
			err = unpopulate(val, "DataVersion", &e.DataVersion)
			delete(rawMsg, key)
		case "eventTime":
			err = unpopulateDateTimeRFC3339(val, "EventTime", &e.EventTime)
			delete(rawMsg, key)
		case "eventType":
			err = unpopulate(val, "EventType", &e.EventType)
			delete(rawMsg, key)
		case "id":
			err = unpopulate(val, "ID", &e.ID)
			delete(rawMsg, key)
		case "metadataVersion":
			err = unpopulate(val, "MetadataVersion", &e.MetadataVersion)
			delete(rawMsg, key)
		case "subject":
			err = unpopulate(val, "Subject", &e.Subject)
			delete(rawMsg, key)
		case "topic":
			err = unpopulate(val, "Topic", &e.Topic)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %w", e, err)
		}
	}
	return nil
}

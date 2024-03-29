//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azeventgrid

import "time"

// Event - Properties of an event published to an Event Grid topic using the EventGrid Schema.
type Event struct {
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

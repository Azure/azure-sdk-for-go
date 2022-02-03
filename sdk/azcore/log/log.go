//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package log provides functionality for configuring logging facilities.
package log

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// Event is used to group entries.  Each group can be toggled on or off.
type Event = log.Event

const (
	// EventRequest entries contain information about HTTP requests.
	// This includes information like the URL, query parameters, and headers.
	EventRequest = log.EventRequest

	// EventResponse entries contain information about HTTP responses.
	// This includes information like the HTTP status code, headers, and request URL.
	EventResponse = log.EventResponse

	// EventRetryPolicy entries contain information specific to the retry policy in use.
	EventRetryPolicy = log.EventRetryPolicy

	// EventLRO entries contain information specific to long-running operations.
	// This includes information like polling location, operation state and sleep intervals.
	EventLRO = log.EventLRO
)

// SetEvents is used to control which events are written to
// the log.  By default all log events are writen.
func SetEvents(cls ...Event) {
	log.SetEvents(cls...)
}

// SetListener will set the Logger to write to the specified Listener.
func SetListener(lst func(Event, string)) {
	log.SetListener(lst)
}

// for testing purposes
func resetEvents() {
	log.TestResetEvents()
}

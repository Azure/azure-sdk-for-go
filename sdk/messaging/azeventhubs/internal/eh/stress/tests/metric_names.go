// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

type Metric string

const (
	// standard to all tests
	MetricNameSent          Metric = "Sent"
	MetricNameReceived      Metric = "Received"
	MetricNameOwnershipLost Metric = "OwnershipLost"

	// go specific
	MetricDeadlineExceeded Metric = "DeadlineExceeded"
)

type Event string

const (
	EventUnbalanced Event = "Unbalanced"
	EventBalanced   Event = "Balanced"
	EventEnd        Event = "end"
	EventStart      Event = "start"
)

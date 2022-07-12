// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package persist

import (
	"time"
)

const (
	// StartOfStream is a constant defined to represent the start of a partition stream in EventHub.
	StartOfStream = "-1"

	// EndOfStream is a constant defined to represent the current end of a partition stream in EventHub.
	// This can be used as an offset argument in receiver creation to start receiving from the latest
	// event, instead of a specific offset or point in time.
	EndOfStream = "@latest"
)

type (
	// Checkpoint is the information needed to determine the last message processed
	Checkpoint struct {
		Offset         string    `json:"offset"`
		SequenceNumber int64     `json:"sequenceNumber"`
		EnqueueTime    time.Time `json:"enqueueTime"`
	}
)

// NewCheckpointFromStartOfStream returns a checkpoint for the start of the stream
func NewCheckpointFromStartOfStream() Checkpoint {
	return Checkpoint{
		Offset: StartOfStream,
	}
}

// NewCheckpointFromEndOfStream returns a checkpoint for the end of the stream
func NewCheckpointFromEndOfStream() Checkpoint {
	return Checkpoint{
		Offset: EndOfStream,
	}
}

// NewCheckpoint contains the information needed to checkpoint Event Hub progress
func NewCheckpoint(offset string, sequence int64, enqueueTime time.Time) Checkpoint {
	return Checkpoint{
		Offset:         offset,
		SequenceNumber: sequence,
		EnqueueTime:    enqueueTime,
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package persist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryPersister(t *testing.T) {
	p := NewMemoryPersister()

	// read a checkpoint before it's been saved will default to "beginning of stream"
	actualCheckpoint, err := p.Read("namespace", "name", "consumerGroup", "0")
	assert.NoError(t, err, "Not an error to get a checkpoint for a partition for the first time")
	assert.Equal(t, NewCheckpointFromStartOfStream(), actualCheckpoint, "If we've never stored a checkpoint, default to beginning of stream")

	now := time.Now()

	// now write one and read it back
	err = p.Write("namespace", "name", "consumerGroup", "0", Checkpoint{
		Offset:         "100",
		SequenceNumber: 2,
		EnqueueTime:    now,
	})

	assert.NoError(t, err)

	actualCheckpoint, err = p.Read("namespace", "name", "consumerGroup", "0")
	assert.NoError(t, err)

	assert.Equal(t, Checkpoint{
		Offset:         "100",
		SequenceNumber: 2,
		EnqueueTime:    now,
	}, actualCheckpoint)
}

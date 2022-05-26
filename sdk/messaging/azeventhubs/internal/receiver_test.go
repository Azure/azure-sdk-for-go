// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/persist"
	"github.com/stretchr/testify/assert"
)

func TestGetOffsetExpression(t *testing.T) {
	expr := getOffsetExpression(persist.NewCheckpointFromStartOfStream())
	assert.EqualValues(t, "amqp.annotation.x-opt-offset > '-1'", expr)

	expr = getOffsetExpression(persist.NewCheckpointFromEndOfStream())
	assert.EqualValues(t, "amqp.annotation.x-opt-offset > '@latest'", expr)

	expr = getOffsetExpression(persist.Checkpoint{Offset: "100"})
	assert.EqualValues(t, "amqp.annotation.x-opt-offset > '100'", expr)

	// offset wins - the time is ignored if they've specified an offset in the checkpoint.t
	now, err := time.Parse(time.RFC3339, "1975-04-04T01:02:03Z")
	assert.NoError(t, err)

	// now's ignored here - the offset will win.
	checkpoint := persist.NewCheckpoint("100", 1, now)
	expr = getOffsetExpression(checkpoint)
	assert.EqualValues(t, "amqp.annotation.x-opt-offset > '100'", expr)

	// no offset this time, date will be used
	checkpoint = persist.NewCheckpoint("", 1, now)
	expr = getOffsetExpression(checkpoint)
	assert.EqualValues(t, "amqp.annotation.x-opt-enqueued-time > '165805323000'", expr)
}

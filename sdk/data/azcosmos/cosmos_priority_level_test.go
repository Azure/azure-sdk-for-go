// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPriorityLevelValues(t *testing.T) {
	values := PriorityLevelValues()
	require.Len(t, values, 2, "expected 2 priority levels")
	require.Equal(t, PriorityLevelHigh, values[0], "expected first value to be High")
	require.Equal(t, PriorityLevelLow, values[1], "expected second value to be Low")
}

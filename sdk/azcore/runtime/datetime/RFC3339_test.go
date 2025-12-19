// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRFC3339(t *testing.T) {
	originalTime := time.Date(2023, time.June, 15, 14, 30, 45, 0, time.UTC)
	dt := RFC3339(originalTime)
	result := dt.String()
	require.NotEmpty(t, result)

	jsonBytes, err := dt.MarshalJSON()
	require.NoError(t, err)
	var dt2 RFC3339
	err = dt2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(dt2))

	textBytes, err := dt.MarshalText()
	require.NoError(t, err)
	var dt3 RFC3339
	err = dt3.UnmarshalText(textBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(dt3))
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnix(t *testing.T) {
	originalTime := time.Date(2023, time.June, 15, 14, 30, 45, 0, time.Local)
	tu := Unix(originalTime)
	result := tu.String()
	expected := fmt.Sprintf("%d", originalTime.Unix())
	require.Equal(t, expected, result)

	jsonBytes, err := tu.MarshalJSON()
	require.NoError(t, err)
	var tu2 Unix
	err = tu2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(tu2))
}

func TestUnix_Invalid(t *testing.T) {
	var tu Unix
	err := tu.UnmarshalJSON([]byte("not-a-number"))
	require.Error(t, err)
}

func TestUnix_Epoch(t *testing.T) {
	tu := Unix(time.Unix(0, 0).UTC())
	result := tu.String()
	require.Equal(t, "0", result)
}

func TestUnix_NegativeTimestamp(t *testing.T) {
	beforeEpoch := time.Date(1969, time.December, 31, 23, 59, 59, 0, time.UTC)
	tu := Unix(beforeEpoch)

	jsonBytes, err := tu.MarshalJSON()
	require.NoError(t, err)

	var tu2 Unix
	err = tu2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, beforeEpoch.Unix(), time.Time(tu2).Unix())
}

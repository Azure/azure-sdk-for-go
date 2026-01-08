// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/datetime"
	"github.com/stretchr/testify/require"
)

func TestUnix(t *testing.T) {
	originalTime := time.Date(2023, time.June, 15, 14, 30, 45, 0, time.Local)
	tu := datetime.Unix(originalTime)
	result := tu.String()
	expected := fmt.Sprintf("%d", originalTime.Unix())
	require.Equal(t, expected, result)

	jsonBytes, err := tu.MarshalJSON()
	require.NoError(t, err)
	var tu2 datetime.Unix
	err = tu2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(tu2))
}

func TestUnix_Invalid(t *testing.T) {
	var tu datetime.Unix
	err := tu.UnmarshalJSON([]byte("not-a-number"))
	require.Error(t, err)
}

func TestUnix_MarshalText(t *testing.T) {
	originalTime := time.Date(2022, time.December, 25, 10, 0, 0, 0, time.Local)
	tu := datetime.Unix(originalTime)
	textBytes, err := tu.MarshalText()
	require.NoError(t, err)
	expected := fmt.Sprintf("%d", originalTime.Unix())
	require.Equal(t, expected, string(textBytes))
}

func TestUnix_UnmarshalText(t *testing.T) {
	originalTime := time.Date(2021, time.November, 5, 8, 15, 30, 0, time.Local)
	text := fmt.Sprintf("%d", originalTime.Unix())
	var tu datetime.Unix
	err := tu.UnmarshalText([]byte(text))
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(tu))
}

func TestUnix_UnmarshalText_Empty(t *testing.T) {
	var tu datetime.Unix
	err := tu.UnmarshalText([]byte(""))
	require.NoError(t, err)
	require.Zero(t, tu)
}

func TestUnix_Epoch(t *testing.T) {
	tu := datetime.Unix(time.Unix(0, 0).UTC())
	result := tu.String()
	require.Equal(t, "0", result)
}

func TestUnix_NegativeTimestamp(t *testing.T) {
	beforeEpoch := time.Date(1969, time.December, 31, 23, 59, 59, 0, time.UTC)
	tu := datetime.Unix(beforeEpoch)

	jsonBytes, err := tu.MarshalJSON()
	require.NoError(t, err)

	var tu2 datetime.Unix
	err = tu2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, beforeEpoch.Unix(), time.Time(tu2).Unix())
}

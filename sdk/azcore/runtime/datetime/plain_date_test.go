// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime_test

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/datetime"
	"github.com/stretchr/testify/require"
)

func TestPlainDate(t *testing.T) {
	originalDate := time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC)
	pd := datetime.PlainDate(originalDate)

	jsonBytes, err := pd.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"2023-01-15"`, string(jsonBytes))

	var pd2 datetime.PlainDate
	err = pd2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)

	require.Equal(t, time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC), time.Time(pd2))
}

func TestPlainDate_TimeIgnored(t *testing.T) {
	date1 := time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2023, time.January, 15, 23, 59, 59, 999999999, time.UTC)

	pd1 := datetime.PlainDate(date1)
	pd2 := datetime.PlainDate(date2)

	result1, _ := pd1.MarshalJSON()
	result2, _ := pd2.MarshalJSON()

	require.Equal(t, `"2023-01-15"`, string(result1))
	require.Equal(t, `"2023-01-15"`, string(result2))
	require.Equal(t, result1, result2)
}

func TestPlainDate_UnmarshalJSON_Invalid_Format(t *testing.T) {
	var pd datetime.PlainDate
	err := pd.UnmarshalJSON([]byte("2023/01/15"))
	require.Error(t, err)
}

func TestPlainDate_UnmarshalJSON_Invalid_Month(t *testing.T) {
	var pd datetime.PlainDate
	err := pd.UnmarshalJSON([]byte("2023-13-01"))
	require.Error(t, err)
}

func TestPlainDate_UnmarshalJSON_Invalid_Day(t *testing.T) {
	var pd datetime.PlainDate
	err := pd.UnmarshalJSON([]byte("2023-01-32"))
	require.Error(t, err)
}

func TestPlainDate_UnmarshalJSON_Empty(t *testing.T) {
	var pd datetime.PlainDate
	err := pd.UnmarshalJSON([]byte(""))
	require.Error(t, err)
}

func TestPlainDate_UnmarshalJSON_Null(t *testing.T) {
	var pd datetime.PlainDate
	err := pd.UnmarshalJSON([]byte("null"))
	require.NoError(t, err)
	require.Zero(t, pd)
}

func TestPlainDate_UnmarshalJSON_PartialDate(t *testing.T) {
	var pd datetime.PlainDate
	err := pd.UnmarshalJSON([]byte("2023-01"))
	require.Error(t, err)
}

func TestPlainDate_MarshalText(t *testing.T) {
	plainDate := datetime.PlainDate(time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC))
	textBytes, err := plainDate.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "2023-01-15", string(textBytes))
}

func TestPlainDate_UnmarshalText(t *testing.T) {
	var plainDate datetime.PlainDate
	err := plainDate.UnmarshalText([]byte("2023-01-15"))
	require.NoError(t, err)
	require.Equal(t, time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC), time.Time(plainDate))
}

func TestPlainDate_UnmarshalText_Nil(t *testing.T) {
	var plainDate datetime.PlainDate
	err := plainDate.UnmarshalText(nil)
	require.NoError(t, err)
	require.Zero(t, plainDate)
}

func TestPlainDate_UnmarshalText_Empty(t *testing.T) {
	var plainDate datetime.PlainDate
	err := plainDate.UnmarshalText([]byte(""))
	require.NoError(t, err)
	require.Zero(t, plainDate)
}

func TestPlainDate_String(t *testing.T) {
	plainDate := datetime.PlainDate(time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC))
	require.Equal(t, "2023-01-15", plainDate.String())
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPlainDate(t *testing.T) {
	originalDate := time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC)
	pd := PlainDate(originalDate)

	jsonBytes, err := pd.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"2023-01-15"`, string(jsonBytes))

	var pd2 PlainDate
	err = pd2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)

	require.Equal(t, time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC), time.Time(pd2))
}

func TestPlainDate_TimeIgnored(t *testing.T) {
	date1 := time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2023, time.January, 15, 23, 59, 59, 999999999, time.UTC)

	pd1 := PlainDate(date1)
	pd2 := PlainDate(date2)

	result1, _ := pd1.MarshalJSON()
	result2, _ := pd2.MarshalJSON()

	require.Equal(t, `"2023-01-15"`, string(result1))
	require.Equal(t, `"2023-01-15"`, string(result2))
	require.Equal(t, result1, result2)
}

func TestPlainDate_UnmarshalJSON_Invalid_Format(t *testing.T) {
	var pd PlainDate
	err := pd.UnmarshalJSON([]byte("2023/01/15"))
	require.Error(t, err)
}

func TestPlainDate_UnmarshalJSON_Invalid_Month(t *testing.T) {
	var pd PlainDate
	err := pd.UnmarshalJSON([]byte("2023-13-01"))
	require.Error(t, err)
}

func TestPlainDate_UnmarshalJSON_Invalid_Day(t *testing.T) {
	var pd PlainDate
	err := pd.UnmarshalJSON([]byte("2023-01-32"))
	require.Error(t, err)
}

func TestPlainDate_UnmarshalJSON_Empty(t *testing.T) {
	var pd PlainDate
	err := pd.UnmarshalJSON([]byte(""))
	require.Error(t, err)
}

func TestPlainDate_UnmarshalJSON_PartialDate(t *testing.T) {
	var pd PlainDate
	err := pd.UnmarshalJSON([]byte("2023-01"))
	require.Error(t, err)
}

func TestPlainDate_String(t *testing.T) {
	plainDate := PlainDate(time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC))
	require.Equal(t, "2023-01-15", plainDate.String())
}

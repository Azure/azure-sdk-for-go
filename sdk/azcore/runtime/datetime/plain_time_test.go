// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime_test

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/datetime"
	"github.com/stretchr/testify/require"
)

func TestPlainTime(t *testing.T) {
	originalTime := time.Date(2023, time.June, 15, 10, 30, 45, 0, time.UTC)
	pt := datetime.PlainTime(originalTime)
	result := pt.String()
	require.Equal(t, "10:30:45", result)

	jsonBytes, err := pt.MarshalJSON()
	require.NoError(t, err)
	var pt2 datetime.PlainTime
	err = pt2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime.Hour(), time.Time(pt2).Hour())
	require.Equal(t, originalTime.Minute(), time.Time(pt2).Minute())
	require.Equal(t, originalTime.Second(), time.Time(pt2).Second())

	textBytes, err := pt.MarshalText()
	require.NoError(t, err)
	var pt3 datetime.PlainTime
	err = pt3.UnmarshalText(textBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime.Hour(), time.Time(pt3).Hour())
	require.Equal(t, originalTime.Minute(), time.Time(pt3).Minute())
	require.Equal(t, originalTime.Second(), time.Time(pt3).Second())
}

func TestPlainTime_UnmarshalText_Empty(t *testing.T) {
	var pt datetime.PlainTime
	err := pt.UnmarshalText([]byte(""))
	require.NoError(t, err)
	require.Zero(t, pt)
}

func TestPlainTime_UnmarshalText_Nil(t *testing.T) {
	var pt datetime.PlainTime
	err := pt.UnmarshalText(nil)
	require.NoError(t, err)
	require.Zero(t, pt)
}

func TestPlainTime_Various(t *testing.T) {
	tests := []struct {
		name   string
		hour   int
		minute int
		second int
	}{
		{"Midnight", 0, 0, 0},
		{"Noon", 12, 0, 0},
		{"Afternoon", 14, 30, 45},
		{"End of day", 23, 59, 59},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalTime := time.Date(2023, time.January, 15, tt.hour, tt.minute, tt.second, 0, time.UTC)
			pt := datetime.PlainTime(originalTime)

			textBytes, err := pt.MarshalText()
			require.NoError(t, err)

			var pt2 datetime.PlainTime
			err = pt2.UnmarshalText(textBytes)
			require.NoError(t, err)

			require.Equal(t, tt.hour, time.Time(pt2).Hour())
			require.Equal(t, tt.minute, time.Time(pt2).Minute())
			require.Equal(t, tt.second, time.Time(pt2).Second())
		})
	}
}

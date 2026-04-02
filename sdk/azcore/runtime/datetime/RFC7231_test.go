// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime_test

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/datetime"
	"github.com/stretchr/testify/require"
)

func TestRFC7231(t *testing.T) {
	originalTime := time.Date(2023, time.June, 15, 14, 30, 45, 0, time.UTC)
	dt := datetime.RFC7231(originalTime)

	// String should always use GMT
	result := dt.String()
	require.Equal(t, "Thu, 15 Jun 2023 14:30:45 GMT", result)

	// MarshalJSON round-trip
	jsonBytes, err := dt.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"Thu, 15 Jun 2023 14:30:45 GMT"`, string(jsonBytes))
	var dt2 datetime.RFC7231
	err = dt2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(dt2))

	// MarshalText round-trip
	textBytes, err := dt.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "Thu, 15 Jun 2023 14:30:45 GMT", string(textBytes))
	var dt3 datetime.RFC7231
	err = dt3.UnmarshalText(textBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(dt3))
}

func TestRFC7231MarshalConvertsToGMT(t *testing.T) {
	// a time in a non-GMT zone should be converted to GMT on marshal
	est := time.FixedZone("EST", -5*60*60)
	originalTime := time.Date(2023, time.June, 15, 9, 30, 45, 0, est) // 09:30 EST = 14:30 UTC/GMT
	dt := datetime.RFC7231(originalTime)

	result := dt.String()
	require.Equal(t, "Thu, 15 Jun 2023 14:30:45 GMT", result)

	jsonBytes, err := dt.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"Thu, 15 Jun 2023 14:30:45 GMT"`, string(jsonBytes))

	textBytes, err := dt.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "Thu, 15 Jun 2023 14:30:45 GMT", string(textBytes))
}

func TestRFC7231UnmarshalJSONError(t *testing.T) {
	var dt datetime.RFC7231
	err := dt.UnmarshalJSON([]byte(`"not a valid date"`))
	require.Error(t, err)
}

func TestRFC7231UnmarshalTextError(t *testing.T) {
	var dt datetime.RFC7231
	err := dt.UnmarshalText([]byte("not a valid date"))
	require.Error(t, err)
}

func TestRFC7231_empty(t *testing.T) {
	tt := datetime.RFC7231{}
	require.NoError(t, xml.Unmarshal([]byte("<RFC7231/>"), &tt))
	require.Zero(t, tt)
}

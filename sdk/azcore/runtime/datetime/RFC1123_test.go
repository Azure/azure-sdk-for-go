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

func TestRFC1123(t *testing.T) {
	originalTime := time.Date(2023, time.June, 15, 14, 30, 45, 0, time.UTC)
	dt := datetime.RFC1123(originalTime)
	result := dt.String()
	require.Equal(t, "Thu, 15 Jun 2023 14:30:45 UTC", result)

	jsonBytes, err := dt.MarshalJSON()
	require.NoError(t, err)
	var dt2 datetime.RFC1123
	err = dt2.UnmarshalJSON(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(dt2))

	textBytes, err := dt.MarshalText()
	require.NoError(t, err)
	var dt3 datetime.RFC1123
	err = dt3.UnmarshalText(textBytes)
	require.NoError(t, err)
	require.Equal(t, originalTime, time.Time(dt3))
}

func TestRFC1123_empty(t *testing.T) {
	tt := datetime.RFC1123{}
	require.NoError(t, xml.Unmarshal([]byte("<RFC1123/>"), &tt))
	require.Zero(t, tt)
}

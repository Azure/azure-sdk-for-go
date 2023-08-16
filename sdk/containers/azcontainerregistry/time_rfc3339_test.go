//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_timeRFC3339_MarshalJSON(t *testing.T) {
	t1 := timeRFC3339(time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC))
	b, err := t1.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"2023-05-06T10:23:15.123456789Z"`, string(b))
}

func Test_timeRFC3339_MarshalText(t *testing.T) {
	t1 := timeRFC3339(time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC))
	b, err := t1.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "2023-05-06T10:23:15.123456789Z", string(b))
}

func Test_timeRFC3339_UnmarshalJSON(t *testing.T) {
	var t1 timeRFC3339
	err := t1.UnmarshalJSON([]byte(`"2023-05-06T10:23:15.123456789Z"`))
	require.NoError(t, err)
	require.Equal(t, timeRFC3339(time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC)), t1)
}

func Test_timeRFC3339_UnmarshalText(t *testing.T) {
	var t1 timeRFC3339
	err := t1.UnmarshalText([]byte("2023-05-06T10:23:15.123456789Z"))
	require.NoError(t, err)
	require.Equal(t, timeRFC3339(time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC)), t1)
}

func Test_populateTimeRFC3339(t *testing.T) {
	m := map[string]any{}
	populateTimeRFC3339(m, "test", nil)
	require.Equal(t, map[string]any{}, m)
	populateTimeRFC3339(m, "test", azcore.NullValue[*time.Time]())
	require.Equal(t, map[string]any{"test": nil}, m)
	t1 := time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC)
	populateTimeRFC3339(m, "test", &t1)
	tt1 := timeRFC3339(t1)
	require.Equal(t, map[string]any{"test": &tt1}, m)
}

func Test_unpopulateTimeRFC3339(t *testing.T) {
	var t1 *time.Time
	var data json.RawMessage
	err := unpopulateTimeRFC3339(data, "test", &t1)
	require.NoError(t, err)
	var tt1 *time.Time
	require.Equal(t, tt1, t1)
	data = json.RawMessage("null")
	err = unpopulateTimeRFC3339(data, "test", &t1)
	require.NoError(t, err)
	require.Equal(t, tt1, t1)
	data = json.RawMessage("wrong value")
	err = unpopulateTimeRFC3339(data, "test", &t1)
	require.Error(t, err)
	data = json.RawMessage(`"2023-05-06T10:23:15.123456789Z"`)
	err = unpopulateTimeRFC3339(data, "test", &t1)
	require.NoError(t, err)
	tt2 := time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC)
	require.Equal(t, &tt2, t1)
}

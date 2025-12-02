// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func Test_dateTimeRFC3339_MarshalJSON(t *testing.T) {
	t1 := dateTimeRFC3339(time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC))
	b, err := t1.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"2023-05-06T10:23:15.123456789Z"`, string(b))
}

func Test_dateTimeRFC3339_MarshalText(t *testing.T) {
	t1 := dateTimeRFC3339(time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC))
	b, err := t1.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "2023-05-06T10:23:15.123456789Z", string(b))
}

func Test_dateTimeRFC3339_UnmarshalJSON(t *testing.T) {
	var t1 dateTimeRFC3339
	err := t1.UnmarshalJSON([]byte(`"2023-05-06T10:23:15.123456789Z"`))
	require.NoError(t, err)
	require.Equal(t, dateTimeRFC3339(time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC)), t1)
}

func Test_dateTimeRFC3339_UnmarshalText(t *testing.T) {
	var t1 dateTimeRFC3339
	err := t1.UnmarshalText([]byte("2023-05-06T10:23:15.123456789Z"))
	require.NoError(t, err)
	require.Equal(t, dateTimeRFC3339(time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC)), t1)
}

func Test_populateDateTimeRFC3339(t *testing.T) {
	m := map[string]any{}
	populateDateTimeRFC3339(m, "test", nil)
	require.Equal(t, map[string]any{}, m)
	populateDateTimeRFC3339(m, "test", azcore.NullValue[*time.Time]())
	require.Equal(t, map[string]any{"test": nil}, m)
	t1 := time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC)
	populateDateTimeRFC3339(m, "test", &t1)
	tt1 := dateTimeRFC3339(t1)
	require.Equal(t, map[string]any{"test": &tt1}, m)
}

func Test_unpopulateDateTimeRFC3339(t *testing.T) {
	var t1 *time.Time
	var data json.RawMessage
	err := unpopulateDateTimeRFC3339(data, "test", &t1)
	require.NoError(t, err)
	var tt1 *time.Time
	require.Equal(t, tt1, t1)
	data = json.RawMessage("null")
	err = unpopulateDateTimeRFC3339(data, "test", &t1)
	require.NoError(t, err)
	require.Equal(t, tt1, t1)
	data = json.RawMessage("wrong value")
	err = unpopulateDateTimeRFC3339(data, "test", &t1)
	require.Error(t, err)
	data = json.RawMessage(`"2023-05-06T10:23:15.123456789Z"`)
	err = unpopulateDateTimeRFC3339(data, "test", &t1)
	require.NoError(t, err)
	tt2 := time.Date(2023, 5, 6, 10, 23, 15, 123456789, time.UTC)
	require.Equal(t, &tt2, t1)
}

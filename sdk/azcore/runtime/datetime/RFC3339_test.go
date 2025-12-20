// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"encoding/json"
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type timeInfoJSON struct {
	StartTime *RFC3339
}

type timeInfoXML struct {
	Expiry *RFC3339
}

type timeInfoSlice struct {
	Expiry *[]RFC3339
}

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

func TestRFC3339_WithSpace_JSON(t *testing.T) {
	dst := timeInfoJSON{}
	require.NoError(t, json.Unmarshal([]byte(`{"startTime":"2024-01-18 14:18:54Z"}`), &dst))
	require.NotNil(t, dst.StartTime)
	require.WithinDuration(t, time.Date(2024, 1, 18, 14, 18, 54, 0, time.UTC), time.Time(*dst.StartTime), 0)

	dst = timeInfoJSON{}
	require.NoError(t, json.Unmarshal([]byte(`{"startTime":"2024-01-18 14:18:54.123Z"}`), &dst))
	require.NotNil(t, dst.StartTime)
	require.WithinDuration(t, time.Date(2024, 1, 18, 14, 18, 54, 123000000, time.UTC), time.Time(*dst.StartTime), 0)

	dst = timeInfoJSON{}
	require.NoError(t, json.Unmarshal([]byte(`{"startTime":"2024-01-18 14:18:54"}`), &dst))
	require.NotNil(t, dst.StartTime)
	require.WithinDuration(t, time.Date(2024, 1, 18, 14, 18, 54, 0, time.UTC), time.Time(*dst.StartTime), 0)

	dst = timeInfoJSON{}
	require.NoError(t, json.Unmarshal([]byte(`{"startTime":"2024-01-18 14:18:54.123"}`), &dst))
	require.NotNil(t, dst.StartTime)
	require.WithinDuration(t, time.Date(2024, 1, 18, 14, 18, 54, 123000000, time.UTC), time.Time(*dst.StartTime), 0)
}

func TestRFC3339_WithSpace_XML(t *testing.T) {
	dst := timeInfoXML{}
	require.NoError(t, xml.Unmarshal([]byte(`<timeInfo><Expiry>2024-01-18 14:18:54Z</Expiry></timeInfo>`), &dst))
	require.NotNil(t, dst.Expiry)
	require.WithinDuration(t, time.Date(2024, 1, 18, 14, 18, 54, 0, time.UTC), time.Time(*dst.Expiry), 0)

	dst = timeInfoXML{}
	require.NoError(t, xml.Unmarshal([]byte(`<timeInfo><Expiry>2024-01-18 14:18:54.123Z</Expiry></timeInfo>`), &dst))
	require.NotNil(t, dst.Expiry)
	require.WithinDuration(t, time.Date(2024, 1, 18, 14, 18, 54, 123000000, time.UTC), time.Time(*dst.Expiry), 0)

	dst = timeInfoXML{}
	require.NoError(t, xml.Unmarshal([]byte(`<timeInfo><Expiry>2024-01-18 14:18:54</Expiry></timeInfo>`), &dst))
	require.NotNil(t, dst.Expiry)
	require.WithinDuration(t, time.Date(2024, 1, 18, 14, 18, 54, 0, time.UTC), time.Time(*dst.Expiry), 0)

	dst = timeInfoXML{}
	require.NoError(t, xml.Unmarshal([]byte(`<timeInfo><Expiry>2024-01-18 14:18:54.123</Expiry></timeInfo>`), &dst))
	require.NotNil(t, dst.Expiry)
	require.WithinDuration(t, time.Date(2024, 1, 18, 14, 18, 54, 123000000, time.UTC), time.Time(*dst.Expiry), 0)
}

func TestEmptyAndNullTime(t *testing.T) {
	dst := timeInfoSlice{}
	require.NoError(t, json.Unmarshal([]byte("{}"), &dst))
	require.Nil(t, dst.Expiry)
	require.NoError(t, json.Unmarshal([]byte(`{"interval": null}`), &dst))
	require.Nil(t, dst.Expiry)
}

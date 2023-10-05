// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package messaging

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestCloudEvent_Minimum(t *testing.T) {
	e, err := NewCloudEvent("source", "eventType", nil, nil)
	require.NoError(t, err)

	require.NotEmpty(t, e)

	require.NotEmpty(t, e.ID)
	require.GreaterOrEqual(t, time.Since(*e.Time), time.Duration(0))

	require.Equal(t, CloudEvent{
		ID:          e.ID,
		Source:      "source",
		SpecVersion: "1.0",
		Time:        e.Time,
		Type:        "eventType",
	}, e)

	actualCE := roundTrip(t, e)

	require.Equal(t, CloudEvent{
		ID:          e.ID,
		Source:      "source",
		SpecVersion: "1.0",
		Time:        e.Time,
		Type:        "eventType",
	}, *actualCE)
}

func TestCloudEventDefaultToTimeNowUTC(t *testing.T) {
	ce, err := NewCloudEvent("source", "type", nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, ce.Time)
}

func TestCloudEventJSONData(t *testing.T) {
	data := map[string]string{
		"randomData": "hello",
	}
	ce, err := NewCloudEvent("source", "type", data, nil)
	require.NoError(t, err)
	require.Equal(t, data, ce.Data)

	// The types change here because the map is ultimately treated as
	// a JSON object, which means the type ends up being map[string]any{}
	// when deserialized.
	actualCE := roundTrip(t, ce)

	var dest *map[string]string
	require.NoError(t, json.Unmarshal(actualCE.Data.([]byte), &dest))

	require.Equal(t, data, *dest)
}

func TestCloudEventUnmarshalFull(t *testing.T) {
	tm, err := time.Parse(time.RFC3339, "2023-06-16T02:54:01Z")
	require.NoError(t, err)

	ce, err := NewCloudEvent("source", "type", []byte{1, 2, 3}, &CloudEventOptions{
		DataContentType: to.Ptr("data content type"),
		DataSchema:      to.Ptr("microsoft.com/dataschema"),
		Extensions: map[string]any{
			"extstr":  "extstring",
			"extnum":  float64(1),
			"extbool": true,
			"exturi":  "http://microsoft.com",
		},
		Subject: to.Ptr("subject"),
		Time:    &tm,
	})
	require.NoError(t, err)
	require.NotEmpty(t, ce.ID)
	require.NotEmpty(t, ce.Time)

	actualCE := roundTrip(t, ce)

	require.NotEmpty(t, actualCE.ID)
	require.Equal(t, &CloudEvent{
		ID:              ce.ID,
		Source:          "source",
		Subject:         to.Ptr("subject"),
		SpecVersion:     "1.0",
		Time:            &tm,
		Type:            "type",
		DataSchema:      to.Ptr("microsoft.com/dataschema"),
		DataContentType: to.Ptr("data content type"),
		Data:            []byte{1, 2, 3},
		Extensions: map[string]any{
			"extstr":  "extstring",
			"extnum":  float64(1),
			"extbool": true,
			"exturi":  "http://microsoft.com",
		},
	}, actualCE)
}

func TestCloudEventUnmarshalFull_InteropWithPython(t *testing.T) {
	// this event is a Python serialized CloudEvent
	text, err := ioutil.ReadFile("testdata/cloudevent_binary_with_extensions.json")
	require.NoError(t, err)

	var ce *CloudEvent

	err = json.Unmarshal(text, &ce)
	require.NoError(t, err)

	tm, err := time.Parse(time.RFC3339, "2023-06-16T02:54:01.470515Z")
	require.NoError(t, err)

	require.Equal(t, &CloudEvent{
		ID:              "2de93014-a793-4170-88f4-1ef74002dfc9",
		Source:          "source",
		Subject:         to.Ptr("subject"),
		SpecVersion:     "1.0",
		Time:            &tm,
		Type:            "type",
		DataSchema:      to.Ptr("microsoft.com/dataschema"),
		DataContentType: to.Ptr("data content type"),
		Data:            []byte{1, 2, 3},
		Extensions: map[string]any{
			"extstr":  "extstring",
			"extnum":  float64(1),
			"extbool": true,
			"exturi":  "http://microsoft.com",
		},
	}, ce)
}

func TestCloudEventUnmarshalRequiredFieldsOnly(t *testing.T) {
	text, err := ioutil.ReadFile("testdata/cloudevent_required_only.json")
	require.NoError(t, err)

	var ce *CloudEvent

	err = json.Unmarshal(text, &ce)
	require.NoError(t, err)

	require.Equal(t, &CloudEvent{
		ID:          "2de93014-a793-4170-88f4-1ef74002dfc9",
		Source:      "source",
		SpecVersion: "1.0",
		Type:        "type",
	}, ce)
}

func TestCloudEventUnmarshalInvalidEvents(t *testing.T) {
	var ce *CloudEvent

	err := json.Unmarshal([]byte("{}"), &ce)
	require.EqualError(t, err, "required field 'id' was not present, or was empty")

	err = json.Unmarshal([]byte(`{"id": "hello"}`), &ce)
	require.EqualError(t, err, "required field 'source' was not present, or was empty")

	err = json.Unmarshal([]byte(`{"id": "hello", "source": "hello"}`), &ce)
	require.EqualError(t, err, "required field 'specversion' was not present, or was empty")

	err = json.Unmarshal([]byte(`{"id": "hello", "source": "hello", "specversion": "1.0"}`), &ce)
	require.EqualError(t, err, "required field 'type' was not present, or was empty")

	err = json.Unmarshal([]byte("invalid-json"), &ce)
	require.EqualError(t, err, "invalid character 'i' looking for beginning of value")

	err = json.Unmarshal([]byte("[]"), &ce)
	require.EqualError(t, err, "json: cannot unmarshal array into Go value of type map[string]json.RawMessage")

	err = json.Unmarshal([]byte(`{"id":100}`), &ce)
	require.EqualError(t, err, `failed to deserialize "id": json: cannot unmarshal number into Go value of type string`)

	err = json.Unmarshal([]byte(`{"data_base64": 1}`), &ce)
	require.EqualError(t, err, `failed to deserialize "data_base64": json: cannot unmarshal number into Go value of type string`)

	err = json.Unmarshal([]byte(`{"data_base64": "not-base-64"}`), &ce)
	require.EqualError(t, err, `failed to deserialize "data_base64": illegal base64 data at input byte 3`)

	err = json.Unmarshal([]byte(`{"time": 100}`), &ce)
	require.EqualError(t, err, `failed to deserialize "time": json: cannot unmarshal number into Go value of type string`)

	err = json.Unmarshal([]byte(`{"time": "not an RFC timestamp"}`), &ce)
	require.EqualError(t, err, `failed to deserialize "time": parsing time "not an RFC timestamp" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "not an RFC timestamp" as "2006"`)
}

func TestGetValue(t *testing.T) {
	var s string
	require.NoError(t, getValue("k", "hello", &s))
	require.Equal(t, "hello", s)

	// this doesn't work because we assume the [T] here is *string
	// and that's not what the rawValue would be.
	var ps *string
	require.EqualError(t, getValue("k", "hello", &ps), `field "k" is a string, but should be *string`)
}

func TestInvalidCloudEvent(t *testing.T) {
	ce, err := NewCloudEvent("", "eventType", nil, nil)
	require.Empty(t, ce)
	require.EqualError(t, err, "source cannot be empty")

	ce, err = NewCloudEvent("source", "", nil, nil)
	require.Empty(t, ce)
	require.EqualError(t, err, "eventType cannot be empty")
}

func roundTrip(t *testing.T, ce CloudEvent) *CloudEvent {
	bytes, err := json.Marshal(ce)
	require.NoError(t, err)

	var dest *CloudEvent
	err = json.Unmarshal(bytes, &dest)
	require.NoError(t, err)

	return dest
}

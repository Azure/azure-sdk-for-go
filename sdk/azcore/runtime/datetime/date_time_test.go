//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/datetime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSuccess(t *testing.T) {
	actual, err := datetime.Parse(datetime.FormatDateOnly, "2024-07-17")
	require.NoError(t, err)
	assert.WithinDuration(t, time.Date(2024, 7, 17, 0, 0, 0, 0, time.UTC), actual, 0)

	actual, err = datetime.Parse(datetime.FormatRFC1123, "Wed, 17 Jul 2024 13:15:00 UTC")
	require.NoError(t, err)
	assert.WithinDuration(t, time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC), actual, 0)

	for _, rfc3339Variant := range []struct {
		Value    string
		Expected time.Time
	}{
		{
			Value:    "2024-07-17T13:15:00Z",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC),
		},
		{
			Value:    "2024-07-17t13:15:00Z",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC),
		},
		{
			Value:    "2024-07-17 13:15:00Z",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC),
		},
		{
			Value:    "2024-07-17T13:15:00",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC),
		},
		{
			Value:    "2024-07-17t13:15:00",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC),
		},
		{
			Value:    "2024-07-17 13:15:00",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC),
		},
		{
			Value:    "2024-07-17T13:15:00.12345Z",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 123450000, time.UTC),
		},
		{
			Value:    "2024-07-17t13:15:00.12345Z",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 123450000, time.UTC),
		},
		{
			Value:    "2024-07-17 13:15:00.12345Z",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 123450000, time.UTC),
		},
		{
			Value:    "2024-07-17T13:15:00.12345",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 123450000, time.UTC),
		},
		{
			Value:    "2024-07-17t13:15:00.12345",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 123450000, time.UTC),
		},
		{
			Value:    "2024-07-17 13:15:00.12345",
			Expected: time.Date(2024, 7, 17, 13, 15, 0, 123450000, time.UTC),
		},
	} {
		t.Run(rfc3339Variant.Value, func(t *testing.T) {
			actual, err = datetime.Parse(datetime.FormatRFC3339, rfc3339Variant.Value)
			require.NoError(t, err)
			assert.WithinDuration(t, rfc3339Variant.Expected, actual, 0)
		})
	}

	actual, err = datetime.Parse(datetime.FormatRFC7231, "Wed, 17 Jul 2024 13:15:00 GMT")
	require.NoError(t, err)
	assert.WithinDuration(t, time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC), actual, 0)

	actual, err = datetime.Parse(datetime.FormatTimeOnly, "13:15:00")
	require.NoError(t, err)
	// note that we use 1 for month and day as, when omitted, are initialized
	// to 1 during parsing. see the docs for time.Parse for why this is.
	assert.WithinDuration(t, time.Date(0, 1, 1, 13, 15, 0, 0, time.UTC), actual, 0)

	actual, err = datetime.Parse(datetime.FormatTimeOnly, "13:15:00.12345")
	require.NoError(t, err)
	// note that we use 1 for month and day as, when omitted, are initialized
	// to 1 during parsing. see the docs for time.Parse for why this is.
	assert.WithinDuration(t, time.Date(0, 1, 1, 13, 15, 0, 123450000, time.UTC), actual, 0)

	actual, err = datetime.Parse(datetime.FormatUnix, "1721222100")
	require.NoError(t, err)
	assert.WithinDuration(t, time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC), actual, 0)
}

func TestParseError(t *testing.T) {
	actual, err := datetime.Parse(datetime.FormatDateOnly, "2024-07-17T13:15:00Z")
	assert.Error(t, err)
	assert.Zero(t, actual)

	actual, err = datetime.Parse(datetime.FormatRFC1123, "2024-07-17")
	assert.Error(t, err)
	assert.Zero(t, actual)

	actual, err = datetime.Parse(datetime.FormatRFC3339, "Wed, 17 Jul 2024 13:15:00 GMT")
	assert.Error(t, err)
	assert.Zero(t, actual)

	actual, err = datetime.Parse(datetime.FormatRFC7231, "2024-07-17T13:15:00Z")
	assert.Error(t, err)
	assert.Zero(t, actual)

	actual, err = datetime.Parse(datetime.FormatTimeOnly, "2024-07-17")
	assert.Error(t, err)
	assert.Zero(t, actual)

	actual, err = datetime.Parse(datetime.FormatUnix, "2024-07-17T13:15:00Z")
	assert.Error(t, err)
	assert.Zero(t, actual)

	actual, err = datetime.Parse(123, "2024-07-17T13:15:00Z")
	assert.Error(t, err)
	assert.Zero(t, actual)
}

func TestUTC(t *testing.T) {
	tt := datetime.UTC(nil)
	assert.Nil(t, tt)

	theTime := time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC)
	assert.EqualValues(t, &theTime, datetime.UTC(&theTime))

	pdt, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	theTime = time.Date(2024, 7, 17, 13, 15, 0, 0, pdt)
	assert.EqualValues(t, to.Ptr(time.Date(2024, 7, 17, 20, 15, 0, 0, time.UTC)), datetime.UTC(&theTime))
}

func TestDateOnly(t *testing.T) {
	const (
		dateValue     = "2024-07-17"
		jsonDateValue = "\"" + dateValue + "\""
		xmlDateValue  = "<DateOnly>" + dateValue + "</DateOnly>"
	)

	src := datetime.DateOnly(time.Date(2024, 7, 17, 0, 0, 0, 0, time.UTC))
	data, err := json.Marshal(src)
	require.NoError(t, err)
	assert.EqualValues(t, jsonDateValue, string(data))
	data, err = xml.Marshal(src)
	require.NoError(t, err)
	assert.EqualValues(t, xmlDateValue, string(data))
	assert.EqualValues(t, dateValue, src.String())

	// including time components should have no effect for marshal/string
	src = datetime.DateOnly(time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC))
	data, err = json.Marshal(src)
	require.NoError(t, err)
	assert.EqualValues(t, jsonDateValue, string(data))
	assert.EqualValues(t, dateValue, src.String())

	dst := datetime.DateOnly{}
	require.NoError(t, json.Unmarshal([]byte(jsonDateValue), &dst))
	assert.EqualValues(t, dateValue, dst.String())
	dst = datetime.DateOnly{}
	require.NoError(t, xml.Unmarshal([]byte(xmlDateValue), &dst))
	assert.EqualValues(t, dateValue, dst.String())

	dst = datetime.DateOnly{}
	require.NoError(t, json.Unmarshal([]byte("null"), &dst))
	assert.Zero(t, dst)
	dst = datetime.DateOnly{}
	require.NoError(t, xml.Unmarshal([]byte("<DateOnly/>"), &dst))
	assert.Zero(t, dst)

	dst = datetime.DateOnly{}
	require.Error(t, json.Unmarshal([]byte(`"invalid-format"`), &dst))
	assert.Zero(t, dst)
}

func TestRFC1123(t *testing.T) {
	gmt, err := time.LoadLocation("GMT")
	require.NoError(t, err)

	for _, rfc1123Variant := range []struct {
		Src      datetime.RFC1123
		Expected string
		JSON     string
		XML      string
	}{
		{
			Src:      datetime.RFC1123(time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC)),
			Expected: "Wed, 17 Jul 2024 13:15:00 UTC",
			JSON:     `"Wed, 17 Jul 2024 13:15:00 UTC"`,
			XML:      "<RFC1123>Wed, 17 Jul 2024 13:15:00 UTC</RFC1123>",
		},
		{
			Src:      datetime.RFC1123(time.Date(2024, 7, 17, 13, 15, 0, 0, gmt)),
			Expected: "Wed, 17 Jul 2024 13:15:00 GMT",
			JSON:     `"Wed, 17 Jul 2024 13:15:00 GMT"`,
			XML:      "<RFC1123>Wed, 17 Jul 2024 13:15:00 GMT</RFC1123>",
		},
	} {
		t.Run(rfc1123Variant.Expected, func(t *testing.T) {
			data, err := json.Marshal(rfc1123Variant.Src)
			require.NoError(t, err)
			assert.EqualValues(t, rfc1123Variant.JSON, string(data))
			data, err = xml.Marshal(rfc1123Variant.Src)
			require.NoError(t, err)
			assert.EqualValues(t, rfc1123Variant.XML, string(data))
			assert.EqualValues(t, rfc1123Variant.Expected, rfc1123Variant.Src.String())
		})
	}

	const (
		rfc1123ValueUTC = "Wed, 17 Jul 2024 13:15:00 UTC"
		rfc1123ValueGMT = "Wed, 17 Jul 2024 13:15:00 GMT"
	)

	for _, rfc1123Variant := range []struct {
		Value    string
		Expected string
	}{
		{
			Value:    "Wed, 17 Jul 2024 13:15:00 UTC",
			Expected: rfc1123ValueUTC,
		},
		{
			Value:    "Wed, 17 Jul 2024 13:15:00 GMT",
			Expected: rfc1123ValueGMT,
		},
	} {
		t.Run(rfc1123Variant.Value, func(t *testing.T) {
			jsonRFC1123UnmarshalValue := fmt.Sprintf("\"%s\"", rfc1123Variant.Value)
			xmlRFC1123UnmarshalValue := fmt.Sprintf("<RFC1123>%s</RFC1123>", rfc1123Variant.Value)

			dst := datetime.RFC1123{}
			require.NoError(t, json.Unmarshal([]byte(jsonRFC1123UnmarshalValue), &dst))
			assert.EqualValues(t, rfc1123Variant.Expected, dst.String())
			dst = datetime.RFC1123{}
			require.NoError(t, xml.Unmarshal([]byte(xmlRFC1123UnmarshalValue), &dst))
			assert.EqualValues(t, rfc1123Variant.Expected, dst.String())
		})
	}

	dst := datetime.RFC1123{}
	require.NoError(t, json.Unmarshal([]byte("null"), &dst))
	assert.Zero(t, dst)
	dst = datetime.RFC1123{}
	require.NoError(t, xml.Unmarshal([]byte("<RFC1123/>"), &dst))
	assert.Zero(t, dst)

	dst = datetime.RFC1123{}
	require.Error(t, json.Unmarshal([]byte(`"invalid-format"`), &dst))
	assert.Zero(t, dst)
}

func TestRFC3339(t *testing.T) {
	for _, rfc3339Variant := range []struct {
		Src      datetime.RFC3339
		Expected string
		JSON     string
		XML      string
	}{
		{
			Src:      datetime.RFC3339(time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC)),
			Expected: "2024-07-17T13:15:00Z",
			JSON:     `"2024-07-17T13:15:00Z"`,
			XML:      "<RFC3339>2024-07-17T13:15:00Z</RFC3339>",
		},
		{
			Src:      datetime.RFC3339(time.Date(2024, 7, 17, 13, 15, 0, 123450000, time.UTC)),
			Expected: "2024-07-17T13:15:00.12345Z",
			JSON:     `"2024-07-17T13:15:00.12345Z"`,
			XML:      "<RFC3339>2024-07-17T13:15:00.12345Z</RFC3339>",
		},
	} {
		t.Run(rfc3339Variant.Expected, func(t *testing.T) {
			data, err := json.Marshal(rfc3339Variant.Src)
			require.NoError(t, err)
			assert.EqualValues(t, rfc3339Variant.JSON, string(data))
			data, err = xml.Marshal(rfc3339Variant.Src)
			require.NoError(t, err)
			assert.EqualValues(t, rfc3339Variant.XML, string(data))
			assert.EqualValues(t, rfc3339Variant.Expected, rfc3339Variant.Src.String())
		})
	}

	const (
		rfc3339Value     = "2024-07-17T13:15:00Z"
		rfc3339ValueNano = "2024-07-17T13:15:00.12345Z"
	)

	for _, rfc3339Variant := range []struct {
		Value    string
		Expected string
	}{
		{
			Value:    "2024-07-17T13:15:00Z",
			Expected: rfc3339Value,
		},
		{
			Value:    "2024-07-17t13:15:00Z",
			Expected: rfc3339Value,
		},
		{
			Value:    "2024-07-17 13:15:00Z",
			Expected: rfc3339Value,
		},
		{
			Value:    "2024-07-17T13:15:00",
			Expected: rfc3339Value,
		},
		{
			Value:    "2024-07-17t13:15:00",
			Expected: rfc3339Value,
		},
		{
			Value:    "2024-07-17 13:15:00",
			Expected: rfc3339Value,
		},
		{
			Value:    "2024-07-17T13:15:00.12345Z",
			Expected: rfc3339ValueNano,
		},
		{
			Value:    "2024-07-17t13:15:00.12345Z",
			Expected: rfc3339ValueNano,
		},
		{
			Value:    "2024-07-17 13:15:00.12345Z",
			Expected: rfc3339ValueNano,
		},
		{
			Value:    "2024-07-17T13:15:00.12345",
			Expected: rfc3339ValueNano,
		},
		{
			Value:    "2024-07-17t13:15:00.12345",
			Expected: rfc3339ValueNano,
		},
		{
			Value:    "2024-07-17 13:15:00.12345",
			Expected: rfc3339ValueNano,
		},
	} {
		t.Run(rfc3339Variant.Value, func(t *testing.T) {
			jsonRFC3339UnmarshalValue := fmt.Sprintf("\"%s\"", rfc3339Variant.Value)
			xmlRFC3339UnmarshalValue := fmt.Sprintf("<RFC3339>%s</RFC3339>", rfc3339Variant.Value)

			dst := datetime.RFC3339{}
			require.NoError(t, json.Unmarshal([]byte(jsonRFC3339UnmarshalValue), &dst))
			assert.EqualValues(t, rfc3339Variant.Expected, dst.String())
			dst = datetime.RFC3339{}
			require.NoError(t, xml.Unmarshal([]byte(xmlRFC3339UnmarshalValue), &dst))
			assert.EqualValues(t, rfc3339Variant.Expected, dst.String())
		})
	}

	dst := datetime.RFC3339{}
	require.NoError(t, json.Unmarshal([]byte("null"), &dst))
	assert.Zero(t, dst)
	dst = datetime.RFC3339{}
	require.NoError(t, xml.Unmarshal([]byte("<RFC3339/>"), &dst))
	assert.Zero(t, dst)

	dst = datetime.RFC3339{}
	require.Error(t, json.Unmarshal([]byte(`"invalid-format"`), &dst))
	assert.Zero(t, dst)
}

func TestRFC7231(t *testing.T) {
	const (
		rfc7231Value  = "Wed, 17 Jul 2024 13:15:00 GMT"
		jsonDateValue = "\"" + rfc7231Value + "\""
		xmlDateValue  = "<RFC7231>" + rfc7231Value + "</RFC7231>"
	)

	src := datetime.RFC7231(time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC))
	data, err := json.Marshal(src)
	require.NoError(t, err)
	assert.EqualValues(t, jsonDateValue, string(data))
	data, err = xml.Marshal(src)
	require.NoError(t, err)
	assert.EqualValues(t, xmlDateValue, string(data))
	assert.EqualValues(t, rfc7231Value, src.String())

	dst := datetime.RFC7231{}
	require.NoError(t, json.Unmarshal([]byte(jsonDateValue), &dst))
	assert.EqualValues(t, rfc7231Value, dst.String())
	dst = datetime.RFC7231{}
	require.NoError(t, xml.Unmarshal([]byte(xmlDateValue), &dst))
	assert.EqualValues(t, rfc7231Value, dst.String())

	dst = datetime.RFC7231{}
	require.NoError(t, json.Unmarshal([]byte("null"), &dst))
	assert.Zero(t, dst)
	dst = datetime.RFC7231{}
	require.NoError(t, xml.Unmarshal([]byte("<RFC7231/>"), &dst))
	assert.Zero(t, dst)

	dst = datetime.RFC7231{}
	require.Error(t, json.Unmarshal([]byte(`"invalid-format"`), &dst))
	assert.Zero(t, dst)
}

func TestTimeOnly(t *testing.T) {
	for _, timeVariant := range []struct {
		Src      datetime.TimeOnly
		Expected string
		JSON     string
		XML      string
	}{
		{
			Src:      datetime.TimeOnly(time.Date(0, 0, 0, 13, 15, 0, 0, time.UTC)),
			Expected: "13:15:00",
			JSON:     `"13:15:00"`,
			XML:      "<TimeOnly>13:15:00</TimeOnly>",
		},
		{
			Src:      datetime.TimeOnly(time.Date(0, 0, 0, 13, 15, 0, 123450000, time.UTC)),
			Expected: "13:15:00.12345",
			JSON:     `"13:15:00.12345"`,
			XML:      "<TimeOnly>13:15:00.12345</TimeOnly>",
		},
		{
			// including date components should have no effect for marshal/string
			Src:      datetime.TimeOnly(time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC)),
			Expected: "13:15:00",
			JSON:     `"13:15:00"`,
			XML:      "<TimeOnly>13:15:00</TimeOnly>",
		},
	} {
		t.Run(timeVariant.Expected, func(t *testing.T) {
			data, err := json.Marshal(timeVariant.Src)
			require.NoError(t, err)
			assert.EqualValues(t, timeVariant.JSON, string(data))
			data, err = xml.Marshal(timeVariant.Src)
			require.NoError(t, err)
			assert.EqualValues(t, timeVariant.XML, string(data))
			assert.EqualValues(t, timeVariant.Expected, timeVariant.Src.String())
		})
	}

	const (
		timeOnlyValue     = "13:15:00"
		timeOnlyValueNano = "13:15:00.12345"
	)

	for _, timeVariant := range []struct {
		Value    string
		Expected string
	}{
		{
			Value:    timeOnlyValue,
			Expected: timeOnlyValue,
		},
		{
			Value:    timeOnlyValueNano,
			Expected: timeOnlyValueNano,
		},
	} {
		t.Run(timeVariant.Expected, func(t *testing.T) {
			jsonTimeOnlyUnmarshalValue := fmt.Sprintf("\"%s\"", timeVariant.Value)
			xmlTimeOnlyUnmarshalValue := fmt.Sprintf("<TimeOnly>%s</TimeOnly>", timeVariant.Value)

			dst := datetime.TimeOnly{}
			require.NoError(t, json.Unmarshal([]byte(jsonTimeOnlyUnmarshalValue), &dst))
			assert.EqualValues(t, timeVariant.Expected, dst.String())
			dst = datetime.TimeOnly{}
			require.NoError(t, xml.Unmarshal([]byte(xmlTimeOnlyUnmarshalValue), &dst))
			assert.EqualValues(t, timeVariant.Expected, dst.String())
		})
	}

	dst := datetime.TimeOnly{}
	require.NoError(t, json.Unmarshal([]byte("null"), &dst))
	assert.Zero(t, dst)
	dst = datetime.TimeOnly{}
	require.NoError(t, xml.Unmarshal([]byte("<TimeOnly/>"), &dst))
	assert.Zero(t, dst)

	dst = datetime.TimeOnly{}
	require.Error(t, json.Unmarshal([]byte(`"invalid-format"`), &dst))
	assert.Zero(t, dst)
}

func TestUnix(t *testing.T) {
	const unixValue = 1721222100
	jsonDateValue := fmt.Sprintf("%d", unixValue)
	xmlDateValue := fmt.Sprintf("<Unix>%d</Unix>", unixValue)

	src := datetime.Unix(time.Date(2024, 7, 17, 13, 15, 0, 0, time.UTC))
	data, err := json.Marshal(src)
	require.NoError(t, err)
	assert.EqualValues(t, jsonDateValue, string(data))
	data, err = xml.Marshal(src)
	require.NoError(t, err)
	assert.EqualValues(t, xmlDateValue, string(data))
	assert.EqualValues(t, jsonDateValue, src.String())

	dst := datetime.Unix{}
	require.NoError(t, json.Unmarshal([]byte(jsonDateValue), &dst))
	assert.EqualValues(t, jsonDateValue, dst.String())
	dst = datetime.Unix{}
	require.NoError(t, xml.Unmarshal([]byte(xmlDateValue), &dst))
	assert.EqualValues(t, jsonDateValue, dst.String())

	dst = datetime.Unix{}
	require.NoError(t, json.Unmarshal([]byte("null"), &dst))
	assert.Zero(t, dst)
	dst = datetime.Unix{}
	require.NoError(t, xml.Unmarshal([]byte("<Unix/>"), &dst))
	assert.Zero(t, dst)

	dst = datetime.Unix{}
	require.Error(t, json.Unmarshal([]byte(`"invalid-format"`), &dst))
	assert.Zero(t, dst)
}

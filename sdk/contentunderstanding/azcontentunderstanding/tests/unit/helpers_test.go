// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentRanges(t *testing.T) {
	for _, test := range []struct {
		name, want string
		build      func() (sdk.ContentRange, error)
	}{
		{"page", "5", func() (sdk.ContentRange, error) { return sdk.Page(5) }},
		{"pages", "1-3", func() (sdk.ContentRange, error) { return sdk.Pages(1, 3) }},
		{"pages-from", "9-", func() (sdk.ContentRange, error) { return sdk.PagesFrom(9) }},
		{"time", "1200-3651", func() (sdk.ContentRange, error) { return sdk.TimeRange(1200*time.Millisecond, 3651*time.Millisecond) }},
		{"time-from", "0-", func() (sdk.ContentRange, error) { return sdk.TimeRangeFrom(0) }},
		{"combined", "1-3,5,9-", func() (sdk.ContentRange, error) { return sdk.CombineRanges("1-3", "5", "9-") }},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.build()
			require.NoError(t, err)
			assert.Equal(t, test.want, got.String())
			assert.Equal(t, test.want, *got.Ptr())
		})
	}
	for _, build := range []func() (sdk.ContentRange, error){
		func() (sdk.ContentRange, error) { return sdk.Page(0) },
		func() (sdk.ContentRange, error) { return sdk.Pages(3, 2) },
		func() (sdk.ContentRange, error) { return sdk.PagesFrom(-1) },
		func() (sdk.ContentRange, error) { return sdk.TimeRange(-1, 0) },
		func() (sdk.ContentRange, error) { return sdk.TimeRange(time.Second, 0) },
		func() (sdk.ContentRange, error) { return sdk.TimeRangeFrom(-1) },
		func() (sdk.ContentRange, error) { return sdk.CombineRanges() },
	} {
		_, err := build()
		require.Error(t, err)
	}
	assert.Equal(t, "custom", sdk.ContentRange("custom").String())
}

func TestContentSources(t *testing.T) {
	raw := "D(1,0.5,1,2.5,1,2.5,3);AV(1200,100,200,50,60);D(2);AV(0)"
	got, err := sdk.ParseContentSources(raw)
	require.NoError(t, err)
	require.Len(t, got, 4)
	assert.Equal(t, 1, *got[0].PageNumber)
	assert.Equal(t, []sdk.Point{{X: .5, Y: 1}, {X: 2.5, Y: 1}, {X: 2.5, Y: 3}}, got[0].Polygon)
	assert.Equal(t, &sdk.BoundingBox{X: .5, Y: 1, Width: 2, Height: 2}, got[0].BoundingBox)
	assert.Equal(t, 1200*time.Millisecond, *got[1].Time)
	assert.Equal(t, &sdk.BoundingBox{X: 100, Y: 200, Width: 50, Height: 60}, got[1].BoundingBox)
	assert.Nil(t, got[2].Polygon)
	assert.Nil(t, got[2].BoundingBox)
	assert.Equal(t, time.Duration(0), *got[3].Time)
	back, err := sdk.FormatContentSources(got)
	require.NoError(t, err)
	assert.Equal(t, raw, back)
	for _, invalid := range []string{"", "D(0)", "D(x)", "D(1,0,1)", "D(1,0,1,2,3,4,5,6)", "D(1,NaN,0,1,1,2,2)", "AV(x)", "AV(1,2)", "AV(2147483648)", "X(1)", "D(1"} {
		_, err := sdk.ParseContentSources(invalid)
		require.Error(t, err, invalid)
	}
	spaced, err := sdk.ParseContentSources(" ; D(3); ; ")
	require.NoError(t, err)
	require.Len(t, spaced, 1)
	_, err = sdk.FormatContentSources(nil)
	require.Error(t, err)
	field := &sdk.StringField{Source: to.Ptr("D(1)")}
	sources, err := sdk.GetFieldSources(field)
	require.NoError(t, err)
	require.Len(t, sources, 1)
	sources, err = sdk.GetFieldSources(&sdk.StringField{})
	require.NoError(t, err)
	assert.Nil(t, sources)
	_, err = sdk.GetFieldSources(nil)
	require.Error(t, err)
}

func capturedResponse(payload string) *http.Response {
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(payload))}
}

func TestCapturedAnalysisHelpers(t *testing.T) {
	raw := `{"id":"op1","status":"Succeeded","usage":{"documentPagesMinimal":0,"audioHours":0.25,"tokens":{"model":12}},"result":{"contents":[{"kind":"audioVisual","startTimeMs":0,"endTimeMs":3651,"KeyFrameTimesMs":[1,1200],"metadata":{"author":"someone"}}]},"future":true}`
	response := capturedResponse(raw)
	id, err := sdk.GetOperationID(response)
	require.NoError(t, err)
	assert.Equal(t, "op1", id)
	usage, err := sdk.GetUsageDetails(response)
	require.NoError(t, err)
	require.NotNil(t, usage)
	assert.Equal(t, int32(0), *usage.DocumentPagesMinimal)
	assert.Equal(t, float32(.25), *usage.AudioHours)
	result, err := sdk.GetAnalysisResult(response)
	require.NoError(t, err)
	require.Len(t, result.Contents, 1)
	av, ok := result.Contents[0].(*sdk.AudioVisualContent)
	require.True(t, ok)
	times, err := av.KeyFrameTimes()
	require.NoError(t, err)
	assert.Equal(t, []time.Duration{time.Millisecond, 1200 * time.Millisecond}, times)
	end, err := av.EndTime()
	require.NoError(t, err)
	assert.Equal(t, 3651*time.Millisecond, end)
	original, err := sdk.GetRawJSON(response)
	require.NoError(t, err)
	assert.JSONEq(t, raw, string(original))
	assert.Contains(t, string(original), `"KeyFrameTimesMs"`)
	assert.NotContains(t, string(original), `"keyFrameTimesMs"`)
	response.Header.Set("Operation-Location", "https://example.test/contentunderstanding/analyzerResults/initial-id?api-version=2025-11-01")
	id, err = sdk.GetOperationID(response)
	require.NoError(t, err)
	assert.Equal(t, "initial-id", id)
	for _, payload := range []string{`{"usage":null}`, `{}`} {
		u, err := sdk.GetUsageDetails(capturedResponse(payload))
		require.NoError(t, err)
		assert.Nil(t, u)
	}
	for _, payload := range []string{`not json`, `{"usage":"bad"}`, `{"usage":{"tokens":{"model":"bad"}}}`} {
		_, err := sdk.GetUsageDetails(capturedResponse(payload))
		require.Error(t, err)
	}
	_, err = sdk.GetOperationID(capturedResponse(`{}`))
	require.Error(t, err)
	_, err = sdk.GetRawJSON(nil)
	require.Error(t, err)
	_, err = sdk.GetAnalysisResult(capturedResponse(`{"status":"Failed"}`))
	require.Error(t, err)
}

func TestMediaTimes(t *testing.T) {
	av := &sdk.AudioVisualContent{StartTimeMs: to.Ptr(int64(0)), EndTimeMs: to.Ptr(int64(3651)), CameraShotTimesMs: []*int64{to.Ptr(int64(1234))}}
	start, err := av.StartTime()
	require.NoError(t, err)
	assert.Zero(t, start)
	shots, err := av.CameraShotTimes()
	require.NoError(t, err)
	assert.Equal(t, []time.Duration{1234 * time.Millisecond}, shots)
	empty, err := av.KeyFrameTimes()
	require.NoError(t, err)
	assert.Empty(t, empty)
	assert.NotNil(t, empty)
	for _, getter := range []func() (time.Duration, error){
		(&sdk.AudioVisualContentSegment{StartTimeMs: to.Ptr(int64(1234))}).StartTime,
		(&sdk.AudioVisualContentSegment{EndTimeMs: to.Ptr(int64(1234))}).EndTime,
		(&sdk.TranscriptPhrase{StartTimeMs: to.Ptr(int64(1234))}).StartTime,
		(&sdk.TranscriptPhrase{EndTimeMs: to.Ptr(int64(1234))}).EndTime,
		(&sdk.TranscriptWord{StartTimeMs: to.Ptr(int64(1234))}).StartTime,
		(&sdk.TranscriptWord{EndTimeMs: to.Ptr(int64(1234))}).EndTime,
	} {
		got, err := getter()
		require.NoError(t, err)
		assert.Equal(t, 1234*time.Millisecond, got)
	}
	_, err = sdk.DurationFromMilliseconds(math.MaxInt64)
	require.Error(t, err)
	_, err = (&sdk.AudioVisualContent{}).StartTime()
	require.Error(t, err)
}

func TestKeyConstructorAndCheckedInline(t *testing.T) {
	calls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		assert.Equal(t, "key", r.Header.Get("Ocp-Apim-Subscription-Key"))
		assert.NotEmpty(t, r.Header.Get("x-ms-client-request-id"))
		_, _ = io.WriteString(w, `{"status":"Failed","error":{"code":"InvalidContent"}}`)
	}))
	defer server.Close()
	client, err := sdk.NewClientWithKeyCredential(server.URL, azcore.NewKeyCredential("key"), &sdk.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: server.Client()}})
	require.NoError(t, err)
	_, err = client.AnalyzeInline(t.Context(), "a", []sdk.AnalysisInput{{URL: to.Ptr("https://example.test/file")}}, "utf16", nil)
	var serviceErr *azcore.ResponseError
	require.ErrorAs(t, err, &serviceErr)
	assert.Equal(t, 200, serviceErr.StatusCode)
	assert.Equal(t, "InvalidContent", serviceErr.ErrorCode)
	assert.Equal(t, 1, calls, "failed inline analysis must not resubmit a billable request")
	_, err = sdk.NewClientWithKeyCredential("http://example.test", azcore.NewKeyCredential("key"), nil)
	require.Error(t, err)
	_, err = sdk.NewClient(server.URL, nil, nil)
	require.Error(t, err)
}

func TestRawJSONPreservesInteger(t *testing.T) {
	raw, err := sdk.GetRawJSON(capturedResponse(`{"value":9007199254740993}`))
	require.NoError(t, err)
	var value struct {
		Value int64 `json:"value"`
	}
	require.NoError(t, json.Unmarshal(raw, &value))
	assert.Equal(t, int64(9007199254740993), value.Value)
}

func TestKeyframeAliasUsesFirstNonNullServiceValue(t *testing.T) {
	for _, fields := range []string{
		`"KeyFrameTimesMs":[100],"keyFrameTimesMs":[200]`,
		`"keyFrameTimesMs":[100],"KeyFrameTimesMs":[200]`,
		`"keyFrameTimesMs":null,"KeyFrameTimesMs":[100]`,
	} {
		t.Run(fields, func(t *testing.T) {
			raw := `{"status":"Succeeded","result":{"contents":[{"kind":"audioVisual",` + fields + `}]}}`
			result, err := sdk.GetAnalysisResult(capturedResponse(raw))
			require.NoError(t, err)
			content, ok := result.Contents[0].(*sdk.AudioVisualContent)
			require.True(t, ok)
			times, err := content.KeyFrameTimes()
			require.NoError(t, err)
			assert.Equal(t, []time.Duration{100 * time.Millisecond}, times)
		})
	}

}

func TestBinaryContentTypeDefault(t *testing.T) {
	client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
		_, _ = io.WriteString(w, succeeded)
	})
	_, err := client.AnalyzeBinaryInline(t.Context(), "a", streamingBody("input"), "utf16", "", nil)
	require.NoError(t, err)
}

func TestKeyCredentialRotation(t *testing.T) {
	var keys []string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys = append(keys, r.Header.Get("Ocp-Apim-Subscription-Key"))
		_, _ = io.WriteString(w, `{"modelDeployments":{}}`)
	}))
	defer server.Close()
	credential := azcore.NewKeyCredential("first-key")
	client, err := sdk.NewClientWithKeyCredential(server.URL, credential, &sdk.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: server.Client()},
	})
	require.NoError(t, err)
	_, err = client.GetDefaults(t.Context(), nil)
	require.NoError(t, err)
	credential.Update("second-key")
	_, err = client.GetDefaults(t.Context(), nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"first-key", "second-key"}, keys)
}

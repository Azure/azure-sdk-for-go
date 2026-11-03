// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdditionalOperations(t *testing.T) {
	client, _ := testClient(t, "2025-11-01", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPatch:
			assert.Equal(t, "/contentunderstanding/analyzers/a", r.URL.Path)
			assert.Equal(t, "application/merge-patch+json", r.Header.Get("Content-Type"))
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"description":"updated"}`, string(body))
			_, _ = io.WriteString(w, `{"analyzerId":"a","description":"updated"}`)
		case r.Method == http.MethodDelete:
			assert.Equal(t, "/contentunderstanding/analyzerResults/op1", r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
		default:
			assert.Equal(t, "/contentunderstanding/analyzerResults/op1/files/keyframes/1200", r.URL.Path)
			w.Header().Set("Content-Type", "image/jpeg")
			_, _ = w.Write([]byte{0xff, 0xd8, 0xff, 0xd9})
		}
	})
	result, err := client.UpdateAnalyzer(t.Context(), "a", sdk.ContentAnalyzer{Description: to.Ptr("updated")}, nil)
	require.NoError(t, err)
	assert.Equal(t, "updated", *result.Description)
	file, err := client.GetResultFile(t.Context(), "op1", "keyframes/1200", nil)
	require.NoError(t, err)
	require.NotNil(t, file.Body)
	bytes, err := io.ReadAll(file.Body)
	require.NoError(t, err)
	require.NoError(t, file.Body.Close())
	assert.Equal(t, []byte{0xff, 0xd8, 0xff, 0xd9}, bytes)
	require.NotNil(t, file.ContentType)
	assert.Equal(t, "image/jpeg", *file.ContentType)
	_, err = client.DeleteResult(t.Context(), "op1", nil)
	require.NoError(t, err)
}

func TestNativeOperationResumeDoesNotResubmit(t *testing.T) {
	var server *httptest.Server
	submissions := 0
	polls := 0
	client, srv := testClient(t, "2025-11-01", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			submissions++
			w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzerResults/resumable")
			w.WriteHeader(http.StatusAccepted)
			_, _ = io.WriteString(w, `{"status":"Running"}`)
		} else {
			polls++
			_, _ = io.WriteString(w, succeeded)
		}
	})
	server = srv
	poller, err := client.BeginAnalyze(t.Context(), "a", []sdk.AnalysisInput{{URL: to.Ptr("https://example.test/file")}}, "utf16", nil)
	require.NoError(t, err)
	token, err := poller.ResumeToken()
	require.NoError(t, err)
	require.NotEmpty(t, token)
	resumed, err := client.BeginAnalyze(t.Context(), "a", nil, "utf16", &sdk.ClientBeginAnalyzeOptions{ResumeToken: token})
	require.NoError(t, err)
	result, err := resumed.PollUntilDone(t.Context(), &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
	require.NoError(t, err)
	assertResult(t, result.AnalysisResult)
	assert.Equal(t, 1, submissions)
	assert.Equal(t, 1, polls)
}

func TestServiceErrorsRemainTyped(t *testing.T) {
	for _, status := range []int{400, 401, 403, 404, 409, 429} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			client, _ := testClient(t, "2025-11-01", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = io.WriteString(w, `{"error":{"code":"ServiceError","message":"test failure"}}`)
			})
			_, err := client.GetAnalyzer(t.Context(), "a", nil)
			var serviceErr *azcore.ResponseError
			require.ErrorAs(t, err, &serviceErr)
			assert.Equal(t, status, serviceErr.StatusCode)
			assert.Equal(t, "ServiceError", serviceErr.ErrorCode)
		})
	}
}

func TestPaginationAndPollingCannotLeakCredentials(t *testing.T) {
	client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"value":[],"nextLink":"https://other.test/contentunderstanding/analyzers"}`)
	})
	pager := client.NewListAnalyzersPager(nil)
	_, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	_, err = pager.NextPage(t.Context())
	require.ErrorContains(t, err, "resource endpoint")
	client, _ = testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Operation-Location", "https://other.test/poll")
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, `{"status":"Running"}`)
	})
	poller, err := client.BeginAnalyze(t.Context(), "a", []sdk.AnalysisInput{{URL: to.Ptr("https://example.test/file")}}, "utf16", nil)
	require.NoError(t, err)
	_, err = poller.Poll(t.Context())
	require.ErrorContains(t, err, "resource endpoint")
}

func TestInlineMalformedAndMissingResults(t *testing.T) {
	for _, body := range []string{`null`, `{}`, `{"status":"Succeeded"}`, `{"status":"Succeeded","result":null}`, `{"status":`} {
		t.Run(body, func(t *testing.T) {
			client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, body) })
			_, err := client.AnalyzeBinaryInline(t.Context(), "a", streamingBody("input"), "utf16", "application/pdf", nil)
			require.Error(t, err)
		})
	}
}

func streamingBody(value string) io.ReadSeekCloser {
	return streaming.NopCloser(strings.NewReader(value))
}

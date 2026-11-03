// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawAnalyzerDefinitionsPreservePatchValues(t *testing.T) {
	body := json.RawMessage(`{"description":"raw definition","extra":{"enabled":false,"value":9007199254740993},"remove":null}`)
	client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(body), string(data))
		if r.Method == http.MethodPut {
			assert.Equal(t, "false", r.URL.Query().Get("allowReplace"))
			w.Header().Set("Operation-Location", "https://example.test/contentunderstanding/analyzers/a/operations/id")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"analyzerId":"a","status":"creating"}`)
		} else {
			assert.Equal(t, "application/merge-patch+json", r.Header.Get("Content-Type"))
			_, _ = io.WriteString(w, `{"analyzerId":"a","description":"raw definition"}`)
		}
	})
	_, err := client.BeginCreateAnalyzerJSON(t.Context(), "a", body, &sdk.ClientBeginCreateAnalyzerOptions{AllowReplace: to.Ptr(false)})
	require.NoError(t, err)
	result, err := client.UpdateAnalyzerJSON(t.Context(), "a", body, nil)
	require.NoError(t, err)
	assert.Equal(t, "raw definition", *result.Description)
	for _, invalid := range []string{`[]`, `null`, `not-json`, ``} {
		_, err := client.BeginCreateAnalyzerJSON(t.Context(), "a", json.RawMessage(invalid), nil)
		require.Error(t, err)
	}
}

func TestRawAnalysisPreservesRequest(t *testing.T) {
	body := json.RawMessage(`{"inputs":[{"url":"https://example.test/file?sig=a%2Bb"}],"modelDeployments":{"gpt-5.2":"model"},"future":false}`)
	for _, inline := range []bool{false, true} {
		t.Run(map[bool]string{false: "lro", true: "inline"}[inline], func(t *testing.T) {
			client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, string(body), string(data))
				assert.Equal(t, "utf16", r.URL.Query().Get("stringEncoding"))
				if inline {
					_, _ = io.WriteString(w, succeeded)
					return
				}
				w.Header().Set("Operation-Location", "https://example.test/contentunderstanding/analyzerResults/id")
				w.WriteHeader(http.StatusAccepted)
				_, _ = io.WriteString(w, `{"status":"Running"}`)
			})
			if inline {
				response, err := client.AnalyzeInlineJSON(t.Context(), "a", body, "utf16", nil)
				require.NoError(t, err)
				require.NotNil(t, response.Result)
			} else {
				_, err := client.BeginAnalyzeJSON(t.Context(), "a", body, "utf16", nil)
				require.NoError(t, err)
			}
		})
	}
}

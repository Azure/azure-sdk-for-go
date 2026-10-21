//go:build live

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const liveSampleURL = "https://raw.githubusercontent.com/Azure-Samples/azure-ai-content-understanding-assets/main/document/invoice.pdf"

func liveClient(t *testing.T, apiVersion string) *Client {
	t.Helper()
	endpoint := os.Getenv("CONTENTUNDERSTANDING_ENDPOINT")
	require.NotEmpty(t, endpoint, "live tests require CONTENTUNDERSTANDING_ENDPOINT")
	credential, err := azidentity.NewAzureCLICredential(nil)
	require.NoError(t, err)
	client, err := NewClient(endpoint, credential, &ClientOptions{APIVersion: apiVersion})
	require.NoError(t, err)
	return client
}

func liveContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Minute)
	t.Cleanup(cancel)
	return ctx
}

func assertLiveAnalysis(t *testing.T, result AnalyzeResponse) {
	t.Helper()
	assert.Equal(t, "Succeeded", result.Status)
	require.NotEmpty(t, result.RawJSON)
	var projected struct {
		AnalyzerID string `json:"analyzerId"`
		Contents   []struct {
			Markdown string                     `json:"markdown"`
			Fields   map[string]json.RawMessage `json:"fields"`
		} `json:"contents"`
	}
	require.NoError(t, json.Unmarshal(result.Result, &projected))
	require.NotEmpty(t, projected.Contents)
	assert.NotEmpty(t, projected.Contents[0].Markdown)
	require.NotEmpty(t, result.Usage, "usage must survive completed analysis")
	assert.NotEqual(t, "null", string(result.Usage))
}

func TestLiveAnalyze(t *testing.T) {
	ctx := liveContext(t)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, liveSampleURL, nil)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)
	data, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	for _, apiVersion := range []string{APIVersionGA, APIVersionPreview} {
		t.Run(apiVersion, func(t *testing.T) {
			client := liveClient(t, apiVersion)
			defaults, err := client.GetDefaults(ctx)
			require.NoError(t, err)
			assert.NotEmpty(t, defaults.ModelDeployments)
			t.Run("binary", func(t *testing.T) {
				poller, err := client.BeginAnalyzeBinary(ctx, "prebuilt-layout", streaming.NopCloser(bytes.NewReader(data)), "application/pdf")
				require.NoError(t, err)
				result, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
				require.NoError(t, err)
				assertLiveAnalysis(t, result)
			})
			t.Run("url", func(t *testing.T) {
				poller, err := client.BeginAnalyze(ctx, "prebuilt-layout", AnalyzeRequest{Inputs: []AnalysisInput{{URL: liveSampleURL}}})
				require.NoError(t, err)
				result, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
				require.NoError(t, err)
				assertLiveAnalysis(t, result)
			})
			t.Run("base64", func(t *testing.T) {
				poller, err := client.BeginAnalyze(ctx, "prebuilt-layout", AnalyzeRequest{
					Inputs: []AnalysisInput{{Data: data, Name: "invoice.pdf", MIMEType: "application/pdf"}},
				})
				require.NoError(t, err)
				result, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
				require.NoError(t, err)
				assertLiveAnalysis(t, result)
			})
			if apiVersion == APIVersionGA {
				t.Run("schema-suggestion", func(t *testing.T) {
					poller, err := client.BeginAnalyze(ctx, "prebuilt-documentFieldSchema", AnalyzeRequest{
						Inputs: []AnalysisInput{{Data: data, Name: "invoice.pdf", MIMEType: "application/pdf"}},
					})
					require.NoError(t, err)
					result, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
					require.NoError(t, err)
					var projected struct {
						Contents []struct {
							Fields map[string]struct {
								ValueJSON json.RawMessage `json:"valueJson"`
							} `json:"fields"`
						} `json:"contents"`
					}
					require.NoError(t, json.Unmarshal(result.Result, &projected))
					require.NotEmpty(t, projected.Contents)
					require.Contains(t, projected.Contents[0].Fields, "schema")
					assert.NotEmpty(t, projected.Contents[0].Fields["schema"].ValueJSON)
				})
			}
			if apiVersion == APIVersionPreview {
				t.Run("inline-binary", func(t *testing.T) {
					result, err := client.AnalyzeBinaryInline(ctx, "prebuilt-layout", streaming.NopCloser(bytes.NewReader(data)), "application/pdf")
					require.NoError(t, err)
					assertLiveAnalysis(t, result)
				})
				t.Run("inline-url", func(t *testing.T) {
					result, err := client.AnalyzeInline(ctx, "prebuilt-layout", AnalyzeRequest{Inputs: []AnalysisInput{{URL: liveSampleURL}}})
					require.NoError(t, err)
					assertLiveAnalysis(t, result)
				})
				t.Run("invalid-inline-input", func(t *testing.T) {
					_, err := client.AnalyzeInline(ctx, "prebuilt-layout", AnalyzeRequest{
						Inputs: []AnalysisInput{{URL: "https://example.invalid/no-input.pdf"}},
					})
					require.Error(t, err)
					var responseErr *azcore.ResponseError
					require.ErrorAs(t, err, &responseErr)
				})
			}
		})
	}
}

func TestLiveAnalyzerLifecycle(t *testing.T) {
	ctx := liveContext(t)
	client := liveClient(t, APIVersionGA)
	id := fmt.Sprintf("go_prototype_%d", time.Now().UnixNano())
	copyID := id + "_copy"
	definition := json.RawMessage(`{"baseAnalyzerId":"prebuilt-document","description":"Go SDK prototype live test","models":{"completion":"gpt-5.2"},"config":{"enableLayout":true},"fieldSchema":{"fields":{"Vendor":{"type":"string","method":"extract","description":"Vendor name on the invoice"}}}}`)
	// IDs belong exclusively to this test. Cleanup runs even after a failed assertion.
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		for _, analyzerID := range []string{copyID, id} {
			err := client.DeleteAnalyzer(cleanupCtx, analyzerID)
			var responseErr *azcore.ResponseError
			if err != nil && !(errors.As(err, &responseErr) && responseErr.StatusCode == 404) {
				t.Errorf("cleanup failed for test analyzer: %v", err)
			}
		}
	})
	poller, err := client.BeginCreateAnalyzer(ctx, id, definition)
	require.NoError(t, err)
	created, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	assert.Equal(t, id, created.AnalyzerID)
	got, err := client.GetAnalyzer(ctx, id)
	require.NoError(t, err)
	assert.Contains(t, string(got.RawJSON), `"Vendor"`)
	t.Run("reject-same-resource-authorization", func(t *testing.T) {
		resourceID := os.Getenv("CONTENTUNDERSTANDING_RESOURCE_ID")
		require.NotEmpty(t, resourceID, "authorization test requires the test resource ARM ID")
		region := os.Getenv("CONTENTUNDERSTANDING_REGION")
		require.NotEmpty(t, region, "authorization test requires the test resource region")
		_, err := client.GrantCopyAuthorization(ctx, id, GrantCopyAuthorizationRequest{
			TargetAzureResourceID: resourceID, TargetRegion: region,
		})
		var responseErr *azcore.ResponseError
		require.ErrorAs(t, err, &responseErr)
		assert.Equal(t, 400, responseErr.StatusCode)
		assert.Equal(t, "InvalidRequest", responseErr.ErrorCode)
		assert.Contains(t, err.Error(), "same as the source resource")
	})
	list := client.NewListAnalyzersPager()
	found := false
	for list.More() {
		page, err := list.NextPage(ctx)
		require.NoError(t, err)
		for _, item := range page.Value {
			if item.AnalyzerID == id {
				found = true
			}
		}
	}
	require.True(t, found, "created analyzer must be visible in paged listing")
	_, err = client.BeginCreateAnalyzer(ctx, id, definition)
	require.Error(t, err)
	var conflict *azcore.ResponseError
	require.ErrorAs(t, err, &conflict)
	assert.Equal(t, 409, conflict.StatusCode)
	analysis, err := client.BeginAnalyze(ctx, id, AnalyzeRequest{Inputs: []AnalysisInput{{URL: liveSampleURL}}})
	require.NoError(t, err)
	result, err := analysis.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	assertLiveAnalysis(t, result)
	var projected struct {
		Contents []struct {
			Fields map[string]json.RawMessage `json:"fields"`
		} `json:"contents"`
	}
	require.NoError(t, json.Unmarshal(result.Result, &projected))
	require.NotEmpty(t, projected.Contents)
	assert.Contains(t, projected.Contents[0].Fields, "Vendor")
	copier, err := client.BeginCopyAnalyzer(ctx, copyID, CopyAnalyzerRequest{SourceAnalyzerID: id})
	require.NoError(t, err)
	copied, err := copier.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	assert.Equal(t, copyID, copied.AnalyzerID)
	require.NoError(t, client.DeleteAnalyzer(ctx, copyID))
	_, err = client.GetAnalyzer(ctx, copyID)
	var missing *azcore.ResponseError
	require.ErrorAs(t, err, &missing)
	assert.Equal(t, 404, missing.StatusCode)
}

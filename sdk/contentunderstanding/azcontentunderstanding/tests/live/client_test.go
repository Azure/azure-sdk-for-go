//go:build live

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleURL = "https://raw.githubusercontent.com/Azure-Samples/azure-ai-content-understanding-assets/main/document/invoice.pdf"

func liveClient(t *testing.T, version string) *sdk.Client {
	t.Helper()
	endpoint := os.Getenv("CONTENTUNDERSTANDING_ENDPOINT")
	require.NotEmpty(t, endpoint, "live tests require CONTENTUNDERSTANDING_ENDPOINT")
	credential, err := azidentity.NewAzureCLICredential(nil)
	require.NoError(t, err)
	client, err := sdk.NewClient(endpoint, credential, &sdk.ClientOptions{ClientOptions: azcore.ClientOptions{APIVersion: version}})
	require.NoError(t, err)
	return client
}

func assertAnalysis(t *testing.T, result sdk.AnalysisResult, response *http.Response) {
	t.Helper()
	require.NotEmpty(t, result.Contents)
	require.NotNil(t, result.Contents[0].GetAnalysisContent().Markdown)
	assert.NotEmpty(t, *result.Contents[0].GetAnalysisContent().Markdown)
	require.NotNil(t, response)
	payload, err := runtime.Payload(response)
	require.NoError(t, err)
	var envelope sdk.ContentAnalyzerAnalyzeOperationStatus
	require.NoError(t, json.Unmarshal(payload, &envelope))
	require.NotNil(t, envelope.Usage, "capture must retain top-level usage omitted from the generated LRO result")
	require.NotNil(t, envelope.Result)
}

func TestLiveGeneratedAnalyze(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Minute)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sampleURL, nil)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, 200, response.StatusCode)
	data, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	for _, version := range []string{"2025-11-01", "2026-06-01-preview"} {
		t.Run(version, func(t *testing.T) {
			client := liveClient(t, version)
			defaults, err := client.GetDefaults(ctx, nil)
			require.NoError(t, err)
			require.NotEmpty(t, defaults.ModelDeployments)
			for _, mode := range []string{"binary", "url", "base64"} {
				t.Run(mode, func(t *testing.T) {
					var raw *http.Response
					capture := policy.WithCaptureResponse(ctx, &raw)
					if mode == "binary" {
						poller, err := client.BeginAnalyzeBinary(capture, "prebuilt-layout", streaming.NopCloser(bytes.NewReader(data)), "codePoint", "", nil)
						require.NoError(t, err)
						result, err := poller.PollUntilDone(capture, &runtime.PollUntilDoneOptions{Frequency: time.Second})
						require.NoError(t, err)
						assertAnalysis(t, result.AnalysisResult, raw)
					} else {
						input := sdk.AnalysisInput{URL: to.Ptr(sampleURL)}
						if mode == "base64" {
							input = sdk.AnalysisInput{Data: data, Name: to.Ptr("invoice.pdf"), MimeType: to.Ptr("application/pdf")}
						}
						poller, err := client.BeginAnalyze(capture, "prebuilt-layout", []sdk.AnalysisInput{input}, "codePoint", nil)
						require.NoError(t, err)
						result, err := poller.PollUntilDone(capture, &runtime.PollUntilDoneOptions{Frequency: time.Second})
						require.NoError(t, err)
						assertAnalysis(t, result.AnalysisResult, raw)
					}
				})
			}
			if version == "2025-11-01" {
				t.Run("schema-suggestion", func(t *testing.T) {
					poller, err := client.BeginAnalyze(ctx, "prebuilt-documentFieldSchema", []sdk.AnalysisInput{{
						Data: data, Name: to.Ptr("invoice.pdf"), MimeType: to.Ptr("application/pdf"),
					}}, "codePoint", nil)
					require.NoError(t, err)
					result, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
					require.NoError(t, err)
					require.NotEmpty(t, result.Contents)
					field, ok := result.Contents[0].GetAnalysisContent().Fields["schema"].(*sdk.JSONField)
					require.True(t, ok)
					var schema map[string]json.RawMessage
					require.NoError(t, json.Unmarshal(field.ValueJSON, &schema))
					assert.NotEmpty(t, schema)
				})
			}
			if version == "2026-06-01-preview" {
				t.Run("inline-binary", func(t *testing.T) {
					result, err := client.AnalyzeBinaryInline(ctx, "prebuilt-layout", streaming.NopCloser(bytes.NewReader(data)), "codePoint", "application/pdf", nil)
					require.NoError(t, err)
					require.NotNil(t, result.Status)
					assert.Equal(t, sdk.OperationStateSucceeded, *result.Status)
					require.NotNil(t, result.Result)
					require.NotEmpty(t, result.Result.Contents)
					require.NotNil(t, result.Usage)
				})
				t.Run("inline-url", func(t *testing.T) {
					result, err := client.AnalyzeInline(ctx, "prebuilt-layout", []sdk.AnalysisInput{{URL: to.Ptr(sampleURL)}}, "codePoint", nil)
					require.NoError(t, err)
					require.NotNil(t, result.Status)
					assert.Equal(t, sdk.OperationStateSucceeded, *result.Status)
					require.NotNil(t, result.Result)
					require.NotNil(t, result.Usage)
				})
				t.Run("invalid-inline-url", func(t *testing.T) {
					result, err := client.AnalyzeInline(ctx, "prebuilt-layout", []sdk.AnalysisInput{{
						URL: to.Ptr("https://example.invalid/no-input.pdf"),
					}}, "codePoint", nil)
					if err != nil {
						var responseErr *azcore.ResponseError
						require.ErrorAs(t, err, &responseErr)
						return
					}
					require.NotNil(t, result.Status)
					assert.Equal(t, sdk.OperationStateFailed, *result.Status)
					assert.Nil(t, result.Result)
				})
			}
		})
	}
}

func TestLiveGeneratedAnalyzerLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Minute)
	defer cancel()
	client := liveClient(t, "2025-11-01")
	id := fmt.Sprintf("go_emitter_%d", time.Now().UnixNano())
	copyID := id + "_copy"
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		for _, analyzerID := range []string{copyID, id} {
			_, err := client.DeleteAnalyzer(cleanupCtx, analyzerID, nil)
			var responseErr *azcore.ResponseError
			if err != nil && !(errors.As(err, &responseErr) && responseErr.StatusCode == 404) {
				t.Errorf("owned analyzer cleanup failed: %v", err)
			}
		}
	})
	var definition sdk.ContentAnalyzer
	require.NoError(t, json.Unmarshal([]byte(`{"baseAnalyzerId":"prebuilt-document","description":"Go emitter local test","models":{"completion":"gpt-5.2"},"config":{"enableLayout":true},"fieldSchema":{"fields":{"Vendor":{"type":"string","method":"extract","description":"Vendor name on the invoice"}}}}`), &definition))
	poller, err := client.BeginCreateAnalyzer(ctx, id, definition, &sdk.ClientBeginCreateAnalyzerOptions{AllowReplace: to.Ptr(false)})
	require.NoError(t, err)
	created, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	require.NotNil(t, created.AnalyzerID)
	assert.Equal(t, id, *created.AnalyzerID)
	got, err := client.GetAnalyzer(ctx, id, nil)
	require.NoError(t, err)
	require.NotNil(t, got.FieldSchema)
	assert.Contains(t, got.FieldSchema.Fields, "Vendor")
	pager := client.NewListAnalyzersPager(nil)
	found := false
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		for _, item := range page.Value {
			if item.AnalyzerID != nil && *item.AnalyzerID == id {
				found = true
			}
		}
	}
	assert.True(t, found)
	_, err = client.BeginCreateAnalyzer(ctx, id, definition, &sdk.ClientBeginCreateAnalyzerOptions{AllowReplace: to.Ptr(false)})
	var conflict *azcore.ResponseError
	require.ErrorAs(t, err, &conflict)
	assert.Equal(t, 409, conflict.StatusCode)
	var raw *http.Response
	capture := policy.WithCaptureResponse(ctx, &raw)
	analysis, err := client.BeginAnalyze(capture, id, []sdk.AnalysisInput{{URL: to.Ptr(sampleURL)}}, "codePoint", nil)
	require.NoError(t, err)
	analyzed, err := analysis.PollUntilDone(capture, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	assertAnalysis(t, analyzed.AnalysisResult, raw)
	require.NotEmpty(t, analyzed.Contents)
	assert.Contains(t, analyzed.Contents[0].GetAnalysisContent().Fields, "Vendor")
	copier, err := client.BeginCopyAnalyzer(ctx, copyID, id, &sdk.ClientBeginCopyAnalyzerOptions{AllowReplace: to.Ptr(false)})
	require.NoError(t, err)
	copied, err := copier.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	require.NotNil(t, copied.AnalyzerID, "generated copy LRO must project its final analyzer")
	assert.Equal(t, copyID, *copied.AnalyzerID)
	_, err = client.DeleteAnalyzer(ctx, copyID, nil)
	require.NoError(t, err)
	_, err = client.GetAnalyzer(ctx, copyID, nil)
	var missing *azcore.ResponseError
	require.ErrorAs(t, err, &missing)
	assert.Equal(t, 404, missing.StatusCode)
}

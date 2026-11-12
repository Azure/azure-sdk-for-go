//go:build live

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sdkSampleBase = "https://raw.githubusercontent.com/Azure/azure-sdk-for-net/" +
	"fe8ee421b95316cb301c3342112f496fb9dab90c/" +
	"sdk/contentunderstanding/Azure.AI.ContentUnderstanding/tests/samples/SampleFiles/"

func fullAnalyze(t *testing.T, client *sdk.Client, analyzerID string, input sdk.AnalysisInput) (*sdk.AnalysisResult, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Minute)
	t.Cleanup(cancel)
	var initial *http.Response
	start := policy.WithCaptureResponse(ctx, &initial)
	poller, err := client.BeginAnalyze(start, analyzerID, []sdk.AnalysisInput{input}, "utf16", nil)
	require.NoError(t, err)
	id, err := sdk.GetOperationID(initial)
	require.NoError(t, err)
	token, err := poller.ResumeToken()
	require.NoError(t, err)
	require.NotEmpty(t, token)
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		_, err := client.DeleteResult(cleanup, id, nil)
		assert.NoError(t, err, "test-owned result cleanup")
	})
	resumed, err := client.BeginAnalyze(ctx, analyzerID, nil, "utf16", &sdk.ClientBeginAnalyzeOptions{ResumeToken: token})
	require.NoError(t, err)
	var final *http.Response
	wait := policy.WithCaptureResponse(ctx, &final)
	_, err = resumed.PollUntilDone(wait, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	usage, err := sdk.GetUsageDetails(final)
	require.NoError(t, err)
	require.NotNil(t, usage)
	raw, err := sdk.GetRawJSON(final)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"result"`)
	result, err := sdk.GetAnalysisResult(final)
	require.NoError(t, err)
	require.NotEmpty(t, result.Contents)
	return result, id
}

func fullAnalyzer(t *testing.T, client *sdk.Client, definition sdk.ContentAnalyzer) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Minute)
	defer cancel()
	id := fmt.Sprintf("go_full_%d", time.Now().UnixNano())
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		_, err := client.DeleteAnalyzer(cleanup, id, nil)
		assert.NoError(t, err, "test-owned analyzer cleanup")
	})
	body, err := json.Marshal(definition)
	require.NoError(t, err)
	poller, err := client.BeginCreateAnalyzerJSON(ctx, id, body, &sdk.ClientBeginCreateAnalyzerOptions{AllowReplace: to.Ptr(false)})
	require.NoError(t, err)
	_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	return id
}

func TestLiveFullDocumentWorkflows(t *testing.T) {
	client := liveClient(t, sdk.APIVersionPreview)
	t.Run("invoice-grounding-usage-formatting", func(t *testing.T) {
		result, _ := fullAnalyze(t, client, "prebuilt-invoice", sdk.AnalysisInput{URL: to.Ptr(sampleURL)})
		grounded := false
		for _, content := range result.Contents {
			for _, field := range content.GetAnalysisContent().Fields {
				sources, err := sdk.GetFieldSources(field)
				require.NoError(t, err)
				if len(sources) > 0 {
					grounded = true
					assert.NotNil(t, sources[0].PageNumber)
				}
			}
		}
		assert.True(t, grounded, "invoice should contain at least one grounded field")
		text, err := sdk.ToLLMInput(result, map[string]any{"source": "public invoice"}, nil)
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(text, "---\n"))
		assert.Contains(t, text, "fields:")
		assert.Contains(t, text, "customMetadata:")
	})
	t.Run("range", func(t *testing.T) {
		page, err := sdk.Page(1)
		require.NoError(t, err)
		result, _ := fullAnalyze(t, client, "prebuilt-layout", sdk.AnalysisInput{URL: to.Ptr(sampleURL), ContentRange: page.Ptr()})
		document, ok := result.Contents[0].(*sdk.DocumentContent)
		require.True(t, ok)
		require.Len(t, document.Pages, 1)
		assert.Equal(t, int32(1), *document.StartPageNumber)
		text, err := sdk.ToLLMInput(result, nil, nil)
		require.NoError(t, err)
		assert.Contains(t, text, "InputPageNumber: 1")
	})
	t.Run("raw-json-analysis", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Minute)
		defer cancel()
		body, err := json.Marshal(map[string]any{"inputs": []map[string]string{{"url": sampleURL}}})
		require.NoError(t, err)
		var initial *http.Response
		poller, err := client.BeginAnalyzeJSON(policy.WithCaptureResponse(ctx, &initial), "prebuilt-layout", body, "utf16", nil)
		require.NoError(t, err)
		id, err := sdk.GetOperationID(initial)
		require.NoError(t, err)
		t.Cleanup(func() {
			cleanup, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := client.DeleteResult(cleanup, id, nil)
			assert.NoError(t, err)
		})
		var raw *http.Response
		_, err = poller.PollUntilDone(policy.WithCaptureResponse(ctx, &raw), &runtime.PollUntilDoneOptions{Frequency: time.Second})
		require.NoError(t, err)
		result, err := sdk.GetAnalysisResult(raw)
		require.NoError(t, err)
		require.NotEmpty(t, result.Contents)
		inline, err := client.AnalyzeInlineJSON(ctx, "prebuilt-layout", body, "utf16", nil)
		require.NoError(t, err)
		require.NotNil(t, inline.Result)
		require.NotNil(t, inline.Usage)
	})
	for _, file := range []string{"sample_metadata.pdf", "sample_metadata.docx"} {
		t.Run(file, func(t *testing.T) {
			result, _ := fullAnalyze(t, client, "prebuilt-layout", sdk.AnalysisInput{URL: to.Ptr(sdkSampleBase + file)})
			metadata := result.Contents[0].GetAnalysisContent().Metadata
			assert.NotEmpty(t, metadata)
			text, err := sdk.ToLLMInput(result, nil, nil)
			require.NoError(t, err)
			assert.Contains(t, text, "metadata:")
		})
	}
	t.Run("signatures", func(t *testing.T) {
		result, _ := fullAnalyze(t, client, "prebuilt-layout", sdk.AnalysisInput{URL: to.Ptr(sdkSampleBase + "sample_signature.png")})
		document, ok := result.Contents[0].(*sdk.DocumentContent)
		require.True(t, ok)
		require.NotEmpty(t, document.Signatures)
	})
}

func TestLiveFullPreviewAnalyzers(t *testing.T) {
	client := liveClient(t, sdk.APIVersionPreview)
	t.Run("semantic-chunks", func(t *testing.T) {
		id := fullAnalyzer(t, client, sdk.ContentAnalyzer{
			BaseAnalyzerID: to.Ptr("prebuilt-document"), Models: map[string]*string{"completion": to.Ptr("gpt-5.2")},
			Config: &sdk.ContentAnalyzerConfig{EnableLayout: to.Ptr(true), ReturnDetails: to.Ptr(true),
				ChunkingStrategy: &sdk.SemanticChunkingStrategy{MaxTokens: to.Ptr(int32(300))}},
		})
		result, _ := fullAnalyze(t, client, id, sdk.AnalysisInput{URL: to.Ptr(sampleURL)})
		document, ok := result.Contents[0].(*sdk.DocumentContent)
		require.True(t, ok)
		require.NotEmpty(t, document.Chunks)
		require.NotNil(t, document.Markdown)
		for _, chunk := range document.Chunks {
			for _, span := range chunk.Spans {
				text, err := sdk.SliceContentSpan(*document.Markdown, *span, "utf16")
				require.NoError(t, err)
				assert.NotEmpty(t, text)
			}
		}
		updated, err := client.UpdateAnalyzer(t.Context(), id, sdk.ContentAnalyzer{Description: to.Ptr("updated by owned live test")}, nil)
		require.NoError(t, err)
		assert.Equal(t, "updated by owned live test", *updated.Description)
		rawUpdated, err := client.UpdateAnalyzerJSON(t.Context(), id, json.RawMessage(`{"description":"raw merge patch from owned live test"}`), nil)
		require.NoError(t, err)
		assert.Equal(t, "raw merge patch from owned live test", *rawUpdated.Description)
	})
	t.Run("agentic", func(t *testing.T) {
		id := fullAnalyzer(t, client, sdk.ContentAnalyzer{
			BaseAnalyzerID: to.Ptr("prebuilt-document"), Models: map[string]*string{"completion": to.Ptr("gpt-5.2")},
			Config: &sdk.ContentAnalyzerConfig{Workflow: to.Ptr(sdk.ContentAnalyzerWorkflowAgentic), ReturnDetails: to.Ptr(true)},
			FieldSchema: &sdk.ContentFieldSchema{Fields: map[string]*sdk.ContentFieldDefinition{
				"Vendor": {Type: to.Ptr(sdk.ContentFieldTypeString), Description: to.Ptr("Name of the invoice vendor")},
			}},
		})
		result, _ := fullAnalyze(t, client, id, sdk.AnalysisInput{URL: to.Ptr(sampleURL)})
		assert.Contains(t, result.Contents[0].GetAnalysisContent().Fields, "Vendor")
	})
	t.Run("in-page-classification", func(t *testing.T) {
		id := fullAnalyzer(t, client, sdk.ContentAnalyzer{
			BaseAnalyzerID: to.Ptr("prebuilt-document"), Models: map[string]*string{"completion": to.Ptr("gpt-5.2")},
			Config: &sdk.ContentAnalyzerConfig{ReturnDetails: to.Ptr(true), EnableSegment: to.Ptr(true), AllowInPageSegments: to.Ptr(true),
				ContentCategories: map[string]*sdk.ContentCategoryDefinition{
					"Invoice":       {Description: to.Ptr("An invoice requesting payment for goods or services, with line items, totals and payment terms.")},
					"BankStatement": {Description: to.Ptr("A bank account statement listing balances, deposits, withdrawals, fees and transactions.")},
				}},
		})
		result, _ := fullAnalyze(t, client, id, sdk.AnalysisInput{URL: to.Ptr(sdkSampleBase + "mixed_financial_docs_in_page.pdf")})
		document, ok := result.Contents[0].(*sdk.DocumentContent)
		require.True(t, ok)
		require.GreaterOrEqual(t, len(document.Segments), 2)
		assert.Equal(t, *document.StartPageNumber, *document.EndPageNumber, "sample must exercise multiple segments within one page")
		text, err := sdk.ToLLMInput(result, nil, nil)
		require.NoError(t, err)
		assert.Contains(t, text, "category:")
	})
}

func TestLiveFullMediaAndResultFile(t *testing.T) {
	client := liveClient(t, sdk.APIVersionPreview)
	t.Run("video-keyframes", func(t *testing.T) {
		selected, err := sdk.TimeRange(0, 5*time.Second)
		require.NoError(t, err)
		result, id := fullAnalyze(t, client, "prebuilt-videoSearch", sdk.AnalysisInput{
			URL:          to.Ptr("https://raw.githubusercontent.com/Azure-Samples/azure-ai-content-understanding-assets/main/videos/sdk_samples/FlightSimulator.mp4"),
			ContentRange: selected.Ptr(),
		})
		video, ok := result.Contents[0].(*sdk.AudioVisualContent)
		require.True(t, ok)
		times, err := video.KeyFrameTimes()
		require.NoError(t, err)
		require.NotEmpty(t, times)
		end, err := video.EndTime()
		require.NoError(t, err)
		assert.Greater(t, end, time.Duration(0))
		file, err := client.GetResultFile(t.Context(), id, fmt.Sprintf("keyframes/%d", times[0].Milliseconds()), nil)
		require.NoError(t, err)
		require.NotNil(t, file.Body)
		data, err := io.ReadAll(file.Body)
		require.NoError(t, err)
		require.NoError(t, file.Body.Close())
		assert.NotEmpty(t, data)
		require.NotNil(t, file.ContentType)
		assert.True(t, strings.HasPrefix(*file.ContentType, "image/"))
		_, err = sdk.ToLLMInput(result, nil, nil)
		require.NoError(t, err)
	})
	t.Run("audio", func(t *testing.T) {
		selected, err := sdk.TimeRange(0, 5*time.Second)
		require.NoError(t, err)
		result, _ := fullAnalyze(t, client, "prebuilt-audioSearch", sdk.AnalysisInput{
			URL:          to.Ptr("https://raw.githubusercontent.com/Azure-Samples/azure-ai-content-understanding-assets/main/audio/callCenterRecording.mp3"),
			ContentRange: selected.Ptr(),
		})
		audio, ok := result.Contents[0].(*sdk.AudioVisualContent)
		require.True(t, ok)
		require.NotEmpty(t, audio.TranscriptPhrases)
		text, err := sdk.ToLLMInput(result, nil, nil)
		require.NoError(t, err)
		assert.NotEmpty(t, text)
	})
}

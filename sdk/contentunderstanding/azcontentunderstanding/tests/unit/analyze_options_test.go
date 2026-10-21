// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeOptionTriStateAndRequestID(t *testing.T) {
	for _, mode := range []string{"json", "binary", "json-inline", "binary-inline"} {
		for _, state := range []string{"omitted", "false", "true"} {
			t.Run(mode+"/"+state, func(t *testing.T) {
				var truncation *bool
				if state != "omitted" {
					truncation = to.Ptr(state == "true")
				}
				requestID := "d20e4222-5bb8-4adf-bcb5-ad43776b5a71"
				inline := mode == "json-inline" || mode == "binary-inline"
				client, _ := testClient(t, sdk.APIVersionPreview, func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, requestID, r.Header.Get("x-ms-client-request-id"))
					assert.Equal(t, "geography", r.URL.Query().Get("processingLocation"))
					if truncation == nil {
						assert.False(t, r.URL.Query().Has("allowInputTruncation"))
					} else {
						assert.Equal(t, strconv.FormatBool(*truncation), r.URL.Query().Get("allowInputTruncation"))
					}
					if mode == "binary" || mode == "binary-inline" {
						assert.Equal(t, "1-3", r.URL.Query().Get("range"))
					}
					if inline {
						_, _ = io.WriteString(w, succeeded)
						return
					}
					w.Header().Set("Operation-Location", "https://example.test/contentunderstanding/analyzerResults/op")
					w.WriteHeader(http.StatusAccepted)
					_, _ = io.WriteString(w, `{"status":"Running"}`)
				})
				ctx := policy.WithHTTPHeader(t.Context(), http.Header{"X-Ms-Client-Request-Id": {requestID}})
				inputs := []sdk.AnalysisInput{{URL: to.Ptr("https://example.test/file")}}
				var err error
				switch mode {
				case "json":
					_, err = client.BeginAnalyze(ctx, "a", inputs, "utf16", &sdk.ClientBeginAnalyzeOptions{
						AllowInputTruncation: truncation, ProcessingLocation: to.Ptr(sdk.ProcessingLocationGeography),
					})
				case "binary":
					_, err = client.BeginAnalyzeBinary(ctx, "a", streamingBody("input"), "utf16", "application/pdf", &sdk.ClientBeginAnalyzeBinaryOptions{
						AllowInputTruncation: truncation, ProcessingLocation: to.Ptr(sdk.ProcessingLocationGeography), Range: to.Ptr("1-3"),
					})
				case "json-inline":
					_, err = client.AnalyzeInline(ctx, "a", inputs, "utf16", &sdk.ClientAnalyzeInlineOptions{
						AllowInputTruncation: truncation, ProcessingLocation: to.Ptr(sdk.ProcessingLocationGeography),
					})
				case "binary-inline":
					_, err = client.AnalyzeBinaryInline(ctx, "a", streamingBody("input"), "utf16", "application/pdf", &sdk.ClientAnalyzeBinaryInlineOptions{
						AllowInputTruncation: truncation, ProcessingLocation: to.Ptr(sdk.ProcessingLocationGeography), Range: to.Ptr("1-3"),
					})
				}
				require.NoError(t, err)
			})
		}
	}
}

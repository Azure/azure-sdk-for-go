// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type credential struct{}

func (credential) GetToken(_ context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "test-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func testClient(t *testing.T, version string, handler http.HandlerFunc) (*sdk.Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewTLSServer(handler)
	t.Cleanup(server.Close)
	client, err := sdk.NewClient(server.URL, credential{}, &sdk.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			APIVersion: version, Transport: server.Client(),
			Retry: policy.RetryOptions{MaxRetries: -1},
		},
	})
	require.NoError(t, err)
	return client, server
}

const succeeded = `{"id":"op1","status":"Succeeded","usage":{"documentPagesMinimal":0,"tokens":{"model":10}},"result":{"analyzerId":"prebuilt-layout","stringEncoding":"codePoint","contents":[{"kind":"document","markdown":"你好 🌍","fields":{"Amount":{"type":"integer","valueInteger":9007199254740993}}}],"future":true},"unknown":{"retained":true}}`

func assertResult(t *testing.T, result sdk.AnalysisResult) {
	t.Helper()
	require.Len(t, result.Contents, 1)
	content, ok := result.Contents[0].(*sdk.DocumentContent)
	require.True(t, ok)
	require.NotNil(t, content.Markdown)
	assert.Equal(t, "你好 🌍", *content.Markdown)
	amount, ok := content.Fields["Amount"].(*sdk.IntegerField)
	require.True(t, ok)
	require.NotNil(t, amount.ValueInteger)
	assert.Equal(t, int64(9007199254740993), *amount.ValueInteger)
}

func TestGeneratedAnalysisPaths(t *testing.T) {
	for _, binary := range []bool{false, true} {
		for _, inline := range []bool{false, true} {
			name := "json"
			if binary {
				name = "binary"
			}
			if inline {
				name += "-inline"
			}
			t.Run(name, func(t *testing.T) {
				var server *httptest.Server
				calls := 0
				client, srv := testClient(t, "2026-06-01-preview", func(w http.ResponseWriter, r *http.Request) {
					calls++
					assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
					assert.Equal(t, "2026-06-01-preview", r.URL.Query().Get("api-version"))
					w.Header().Set("Content-Type", "application/json")
					if r.Method == http.MethodGet {
						_, _ = io.WriteString(w, succeeded)
						return
					}
					action := ":analyze"
					if binary {
						action += "Binary"
					}
					if inline {
						action += "Inline"
					}
					assert.Equal(t, "/contentunderstanding/analyzers/prebuilt-layout"+action, r.URL.Path)
					assert.Equal(t, "codePoint", r.URL.Query().Get("stringEncoding"))
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					if binary {
						assert.Equal(t, "input", string(body))
						assert.Equal(t, "application/pdf", r.Header.Get("Content-Type"))
					} else {
						assert.JSONEq(t, `{"inputs":[{"url":"https://storage.test/in.pdf?sig=a%2Bb%3D&sp=r"}],"modelDeployments":null}`, string(body))
					}
					if inline {
						_, _ = io.WriteString(w, succeeded)
					} else {
						w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzerResults/op1?api-version=2026-06-01-preview")
						w.WriteHeader(http.StatusAccepted)
						_, _ = io.WriteString(w, `{"status":"Running"}`)
					}
				})
				server = srv
				var raw *http.Response
				ctx := policy.WithCaptureResponse(t.Context(), &raw)
				inputs := []sdk.AnalysisInput{{URL: to.Ptr("https://storage.test/in.pdf?sig=a%2Bb%3D&sp=r")}}
				if inline {
					if binary {
						response, err := client.AnalyzeBinaryInline(ctx, "prebuilt-layout", streaming.NopCloser(strings.NewReader("input")), "codePoint", "application/pdf", nil)
						require.NoError(t, err)
						require.NotNil(t, response.Result)
						assertResult(t, *response.Result)
						require.NotNil(t, response.Usage)
					} else {
						response, err := client.AnalyzeInline(ctx, "prebuilt-layout", inputs, "codePoint", nil)
						require.NoError(t, err)
						require.NotNil(t, response.Result)
						assertResult(t, *response.Result)
						require.NotNil(t, response.Usage)
					}
				} else if binary {
					poller, err := client.BeginAnalyzeBinary(ctx, "prebuilt-layout", streaming.NopCloser(strings.NewReader("input")), "codePoint", "application/pdf", nil)
					require.NoError(t, err)
					result, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
					require.NoError(t, err)
					assertResult(t, result.AnalysisResult)
				} else {
					poller, err := client.BeginAnalyze(ctx, "prebuilt-layout", inputs, "codePoint", nil)
					require.NoError(t, err)
					result, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
					require.NoError(t, err)
					assertResult(t, result.AnalysisResult)
				}
				require.NotNil(t, raw)
				payload, err := runtime.Payload(raw)
				require.NoError(t, err)
				assert.JSONEq(t, succeeded, string(payload))
				var envelope sdk.ContentAnalyzerAnalyzeOperationStatus
				require.NoError(t, json.Unmarshal(payload, &envelope))
				require.NotNil(t, envelope.Usage)
				assert.Equal(t, int32(0), *envelope.Usage.DocumentPagesMinimal)
				if inline {
					assert.Equal(t, 1, calls)
				} else {
					assert.Equal(t, 2, calls)
				}
			})
		}
	}
}

func TestGeneratedBase64Input(t *testing.T) {
	client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"inputs":[{"data":"aGVsbG8=","name":"sample.txt","mimeType":"text/plain"}],"modelDeployments":null}`, string(body))
		_, _ = io.WriteString(w, succeeded)
	})
	_, err := client.AnalyzeInline(t.Context(), "prebuilt-layout", []sdk.AnalysisInput{{
		Data: []byte("hello"), Name: to.Ptr("sample.txt"), MimeType: to.Ptr("text/plain"),
	}}, "codePoint", nil)
	require.NoError(t, err)
}

func TestGeneratedDefaults(t *testing.T) {
	client, _ := testClient(t, "2025-11-01", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2025-11-01", r.URL.Query().Get("api-version"))
		if r.Method == http.MethodPatch {
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Equal(t, "application/merge-patch+json", r.Header.Get("Content-Type"))
			assert.JSONEq(t, `{"modelDeployments":{"gpt-5.2":"deployment","old":null}}`, string(body))
		}
		_, _ = io.WriteString(w, `{"modelDeployments":{"gpt-5.2":"deployment"}}`)
	})
	defaults, err := client.GetDefaults(t.Context(), nil)
	require.NoError(t, err)
	assert.Equal(t, "deployment", *defaults.ModelDeployments["gpt-5.2"])
	_, err = client.UpdateDefaults(t.Context(), &sdk.ClientUpdateDefaultsOptions{
		ModelDeployments: &sdk.RecordMergePatchUpdate{
			AdditionalProperties: map[string]*string{"gpt-5.2": to.Ptr("deployment"), "old": nil},
		},
	})
	require.NoError(t, err)
}

func TestGeneratedPagingAndLifecycle(t *testing.T) {
	var server *httptest.Server
	client, srv := testClient(t, "2025-11-01", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut:
			assert.Equal(t, "false", r.URL.Query().Get("allowReplace"))
			w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzers/a/operations/op1")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"analyzerId":"a","status":"creating"}`)
		case strings.Contains(r.URL.Path, "/operations/"):
			_, _ = io.WriteString(w, `{"status":"Succeeded","result":{"analyzerId":"a","status":"ready"}}`)
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case strings.HasSuffix(r.URL.Path, "/a"):
			_, _ = io.WriteString(w, `{"analyzerId":"a","status":"ready"}`)
		case r.URL.Query().Get("skip") == "1":
			_, _ = io.WriteString(w, `{"value":[{"analyzerId":"b"}]}`)
		default:
			_, _ = io.WriteString(w, `{"value":[{"analyzerId":"a"}],"nextLink":"`+server.URL+`/contentunderstanding/analyzers?skip=1"}`)
		}
	})
	server = srv
	poller, err := client.BeginCreateAnalyzer(t.Context(), "a", sdk.ContentAnalyzer{BaseAnalyzerID: to.Ptr("prebuilt-document")}, &sdk.ClientBeginCreateAnalyzerOptions{AllowReplace: to.Ptr(false)})
	require.NoError(t, err)
	result, err := poller.PollUntilDone(t.Context(), &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
	require.NoError(t, err)
	require.NotNil(t, result.AnalyzerID)
	assert.Equal(t, "a", *result.AnalyzerID)
	pager := client.NewListAnalyzersPager(nil)
	var ids []string
	for pager.More() {
		page, err := pager.NextPage(t.Context())
		require.NoError(t, err)
		for _, item := range page.Value {
			ids = append(ids, *item.AnalyzerID)
		}
	}
	assert.Equal(t, []string{"a", "b"}, ids)
	_, err = client.DeleteAnalyzer(t.Context(), "a", nil)
	require.NoError(t, err)
}

func TestGeneratedCopyRequestContracts(t *testing.T) {
	client, _ := testClient(t, "2025-11-01", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		if strings.HasSuffix(r.URL.Path, ":grantCopyAuthorization") {
			assert.JSONEq(t, `{"targetAzureResourceId":"/subscriptions/target","targetRegion":"westus"}`, string(body))
			_, _ = io.WriteString(w, `{"source":"a","targetAzureResourceId":"/subscriptions/target","expiresAt":"2026-10-21T23:00:00Z"}`)
		} else {
			assert.JSONEq(t, `{"sourceAnalyzerId":"a","sourceAzureResourceId":"/subscriptions/source","sourceRegion":"eastus"}`, string(body))
			w.WriteHeader(http.StatusConflict)
			_, _ = io.WriteString(w, `{"error":{"code":"Conflict","message":"already exists"}}`)
		}
	})
	grant, err := client.GrantCopyAuthorization(t.Context(), "a", "/subscriptions/target", &sdk.ClientGrantCopyAuthorizationOptions{TargetRegion: to.Ptr("westus")})
	require.NoError(t, err)
	assert.Equal(t, "/subscriptions/target", *grant.TargetAzureResourceID)
	_, err = client.BeginCopyAnalyzer(t.Context(), "b", "a", &sdk.ClientBeginCopyAnalyzerOptions{
		SourceAzureResourceID: to.Ptr("/subscriptions/source"), SourceRegion: to.Ptr("eastus"), AllowReplace: to.Ptr(false),
	})
	var responseErr *azcore.ResponseError
	require.ErrorAs(t, err, &responseErr)
	assert.Equal(t, 409, responseErr.StatusCode)
}

func TestInlineFailedEnvelopeReturnsError(t *testing.T) {
	calls := 0
	client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		calls++
		_, _ = io.WriteString(w, `{"status":"Failed","error":{"code":"InvalidContent"}}`)
	})
	_, err := client.AnalyzeInline(t.Context(), "a", []sdk.AnalysisInput{{URL: to.Ptr("https://example.test/in.pdf")}}, "codePoint", nil)
	var responseErr *azcore.ResponseError
	require.ErrorAs(t, err, &responseErr)
	assert.Equal(t, "InvalidContent", responseErr.ErrorCode)
	assert.Equal(t, 200, responseErr.StatusCode)
	assert.Equal(t, 1, calls)
}

func TestGeneratedLROFailure(t *testing.T) {
	var server *httptest.Server
	client, srv := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzerResults/op1")
			w.WriteHeader(http.StatusAccepted)
			_, _ = io.WriteString(w, `{"status":"Running"}`)
			return
		}
		_, _ = io.WriteString(w, `{"status":"Failed","error":{"code":"AnalysisFailed","message":"invalid input"}}`)
	})
	server = srv
	poller, err := client.BeginAnalyze(t.Context(), "a", []sdk.AnalysisInput{{URL: to.Ptr("https://example.test/in.pdf")}}, "codePoint", nil)
	require.NoError(t, err)
	_, err = poller.PollUntilDone(t.Context(), &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
	var responseErr *azcore.ResponseError
	require.ErrorAs(t, err, &responseErr)
	assert.Equal(t, "AnalysisFailed", responseErr.ErrorCode)
}

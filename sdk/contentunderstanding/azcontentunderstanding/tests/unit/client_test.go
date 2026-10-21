// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T, apiVersion string, handler http.HandlerFunc) (*sdk.Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewTLSServer(handler)
	t.Cleanup(server.Close)
	client, err := sdk.NewClientWithKeyCredential(server.URL, azcore.NewKeyCredential("test-key"), &sdk.ClientOptions{
		APIVersion: apiVersion, ClientOptions: azcore.ClientOptions{
			Transport: server.Client(), Retry: policy.RetryOptions{MaxRetries: -1},
		},
	})
	require.NoError(t, err)
	return client, server
}

const succeededAnalysis = `{"id":"op1","status":"Succeeded","usage":{"documentPagesMinimal":0,"tokens":{"large":9007199254740993}},"result":{"analyzerId":"prebuilt-layout","stringEncoding":"codePoint","contents":[{"kind":"document","markdown":"你好 🌍","fields":{"Amount":{"type":"integer","valueInteger":9007199254740993}}}],"future":true},"unknown":{"retained":true}}`

func TestAnalysisPaths(t *testing.T) {
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
				var calls atomic.Int32
				var server *httptest.Server
				client, srv := testClient(t, sdk.APIVersionPreview, func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					assert.Equal(t, "test-key", r.Header.Get("Ocp-Apim-Subscription-Key"))
					assert.Equal(t, sdk.APIVersionPreview, r.URL.Query().Get("api-version"))
					w.Header().Set("Content-Type", "application/json")
					if r.Method == http.MethodGet {
						assert.Equal(t, "/contentunderstanding/analyzerResults/op1", r.URL.Path)
						_, _ = io.WriteString(w, succeededAnalysis)
						return
					}
					assert.Equal(t, http.MethodPost, r.Method)
					action := ":analyze"
					if binary {
						action += "Binary"
					}
					if inline {
						action += "Inline"
					}
					assert.Equal(t, "/contentunderstanding/analyzers/prebuilt-layout"+action, r.URL.Path)
					assert.Equal(t, "codePoint", r.URL.Query().Get("stringEncoding"))
					data, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					if binary {
						assert.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
						assert.Equal(t, []byte("input"), data)
					} else {
						assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
						assert.JSONEq(t, `{"inputs":[{"url":"https://storage.test/a.pdf?sig=a%2Bb%3D&sp=r"}]}`, string(data))
					}
					if inline {
						_, _ = io.WriteString(w, succeededAnalysis)
						return
					}
					w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzerResults/op1?api-version="+sdk.APIVersionPreview)
					w.WriteHeader(http.StatusAccepted)
					_, _ = io.WriteString(w, `{"status":"Running"}`)
				})
				server = srv
				var got sdk.AnalyzeResponse
				var err error
				input := sdk.AnalyzeRequest{Inputs: []sdk.AnalysisInput{{URL: "https://storage.test/a.pdf?sig=a%2Bb%3D&sp=r"}}}
				if inline {
					if binary {
						got, err = client.AnalyzeBinaryInline(t.Context(), "prebuilt-layout", streamingBody("input"), "")
					} else {
						got, err = client.AnalyzeInline(t.Context(), "prebuilt-layout", input)
					}
				} else {
					var poller *runtime.Poller[sdk.AnalyzeResponse]
					if binary {
						poller, err = client.BeginAnalyzeBinary(t.Context(), "prebuilt-layout", streamingBody("input"), "")
					} else {
						poller, err = client.BeginAnalyze(t.Context(), "prebuilt-layout", input)
					}
					require.NoError(t, err)
					got, err = poller.PollUntilDone(t.Context(), &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
				}
				require.NoError(t, err)
				assert.JSONEq(t, succeededAnalysis, string(got.RawJSON))
				assert.Contains(t, string(got.Result), "你好 🌍")
				assert.Contains(t, string(got.Usage), "9007199254740993")
				if inline {
					assert.Equal(t, int32(1), calls.Load())
				} else {
					assert.Equal(t, int32(2), calls.Load())
				}
			})
		}
	}
}

func streamingBody(value string) io.ReadSeekCloser {
	return streaming.NopCloser(strings.NewReader(value))
}

func TestBase64Input(t *testing.T) {
	client, _ := testClient(t, sdk.APIVersionPreview, func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"inputs":[{"data":"aGVsbG8=","name":"sample.txt","mimeType":"text/plain"}]}`, string(data))
		_, _ = io.WriteString(w, succeededAnalysis)
	})
	_, err := client.AnalyzeInline(t.Context(), "prebuilt-layout", sdk.AnalyzeRequest{
		Inputs: []sdk.AnalysisInput{{Data: []byte("hello"), Name: "sample.txt", MIMEType: "text/plain"}},
	})
	require.NoError(t, err)
}

func TestInlineErrorsAndValidation(t *testing.T) {
	for _, payload := range []string{
		`{"status":"Failed","error":{"code":"InvalidContent","message":"invalid content"}}`,
		`{"status":"Canceled"}`, `{}`, `null`, `{"status":"Succeeded","result":null}`,
		`{"status":"Succeeded","result":[]}`, `{"status":`,
	} {
		t.Run(payload, func(t *testing.T) {
			client, _ := testClient(t, sdk.APIVersionPreview, func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, payload)
			})
			_, err := client.AnalyzeInline(t.Context(), "a", sdk.AnalyzeRequest{Inputs: []sdk.AnalysisInput{{URL: "https://example.test/a"}}})
			require.Error(t, err)
			if strings.Contains(payload, "InvalidContent") {
				var responseErr *azcore.ResponseError
				require.ErrorAs(t, err, &responseErr)
				assert.Equal(t, "InvalidContent", responseErr.ErrorCode)
				assert.Equal(t, 200, responseErr.StatusCode)
			}
		})
	}
	client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) { t.Error("invalid request reached transport") })
	_, err := client.AnalyzeInline(t.Context(), "a", sdk.AnalyzeRequest{})
	require.ErrorContains(t, err, "preview")
	for _, body := range []sdk.AnalyzeRequest{
		{}, {Inputs: []sdk.AnalysisInput{{}}}, {Inputs: []sdk.AnalysisInput{{URL: "https://example.test", Data: []byte("both")}}},
	} {
		_, err := client.BeginAnalyze(t.Context(), "a", body)
		require.Error(t, err)
	}
	_, err = client.BeginAnalyzeBinary(t.Context(), "a", nil, "")
	require.Error(t, err)
}

func TestDefaultsAndServiceErrors(t *testing.T) {
	var count int
	client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		count++
		assert.Equal(t, sdk.APIVersionGA, r.URL.Query().Get("api-version"))
		assert.Equal(t, "/contentunderstanding/defaults", r.URL.Path)
		switch count {
		case 1:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"error":{"code":"InvalidRequest","message":"DefaultsNotSet","innererror":{"code":"DefaultsNotSet"}}}`)
		case 2:
			assert.Equal(t, "PATCH", r.Method)
			assert.Equal(t, "application/merge-patch+json", r.Header.Get("Content-Type"))
			data, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"modelDeployments":{"old":null,"gpt-5.2":"deployment"}}`, string(data))
			_, _ = io.WriteString(w, `{"modelDeployments":{"gpt-5.2":"deployment"}}`)
		case 3:
			_, _ = io.WriteString(w, `{"modelDeployments":{"gpt-5.2":"deployment"}}`)
		}
	})
	_, err := client.GetDefaults(t.Context())
	var responseErr *azcore.ResponseError
	require.ErrorAs(t, err, &responseErr)
	assert.Equal(t, http.StatusBadRequest, responseErr.StatusCode)
	value := "deployment"
	got, err := client.UpdateDefaults(t.Context(), map[string]*string{"old": nil, "gpt-5.2": &value})
	require.NoError(t, err)
	assert.Equal(t, value, got.ModelDeployments["gpt-5.2"])
	got, err = client.GetDefaults(t.Context())
	require.NoError(t, err)
	assert.Equal(t, value, got.ModelDeployments["gpt-5.2"])
}

func TestAnalyzerLifecycleAndPaging(t *testing.T) {
	var server *httptest.Server
	definition := json.RawMessage(`{"baseAnalyzerId":"prebuilt-document","future":{"enabled":false}}`)
	client, srv := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case r.Method == http.MethodPut:
			assert.Equal(t, "false", r.URL.Query().Get("allowReplace"))
			data, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, string(definition), string(data))
			w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzers/a/operations/1")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"analyzerId":"a","status":"Creating"}`)
		case strings.Contains(path, "/operations/"):
			_, _ = io.WriteString(w, `{"status":"Succeeded","result":{"analyzerId":"a","status":"ready"}}`)
		case path == "/contentunderstanding/analyzers/a":
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			_, _ = io.WriteString(w, `{"analyzerId":"a","status":"ready","future":{"preserve":true}}`)
		case r.URL.Query().Get("skip") == "1":
			_, _ = io.WriteString(w, `{"value":[{"analyzerId":"b"}]}`)
		default:
			_, _ = io.WriteString(w, `{"value":[{"analyzerId":"a"}],"nextLink":"`+server.URL+`/contentunderstanding/analyzers?skip=1&api-version=2025-11-01"}`)
		}
	})
	server = srv
	poller, err := client.BeginCreateAnalyzer(t.Context(), "a", definition)
	require.NoError(t, err)
	got, err := poller.PollUntilDone(t.Context(), &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
	require.NoError(t, err)
	assert.Equal(t, "a", got.AnalyzerID)
	got, err = client.GetAnalyzer(t.Context(), "a")
	require.NoError(t, err)
	assert.Contains(t, string(got.RawJSON), `"preserve":true`)
	pager := client.NewListAnalyzersPager()
	var ids []string
	for pager.More() {
		page, err := pager.NextPage(t.Context())
		require.NoError(t, err)
		for _, item := range page.Value {
			ids = append(ids, item.AnalyzerID)
		}
	}
	assert.Equal(t, []string{"a", "b"}, ids)
	require.NoError(t, client.DeleteAnalyzer(t.Context(), "a"))
}

func TestCopyRequests(t *testing.T) {
	var server *httptest.Server
	var grantSeen bool
	client, srv := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, ":grantCopyAuthorization"):
			grantSeen = true
			data, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"targetAzureResourceId":"/subscriptions/target","targetRegion":"westus"}`, string(data))
			_, _ = io.WriteString(w, `{"source":"source","targetAzureResourceId":"/subscriptions/target","expiresAt":"2026-10-21T23:00:00Z"}`)
		case strings.HasSuffix(r.URL.Path, ":copy"):
			assert.True(t, grantSeen)
			assert.Equal(t, "false", r.URL.Query().Get("allowReplace"))
			data, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"sourceAnalyzerId":"a","sourceAzureResourceId":"/subscriptions/source","sourceRegion":"eastus"}`, string(data))
			w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzers/b/operations/1")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"analyzerId":"b","status":"Creating"}`)
		default:
			_, _ = io.WriteString(w, `{"status":"Succeeded","result":{"analyzerId":"b","status":"ready"}}`)
		}
	})
	server = srv
	grant, err := client.GrantCopyAuthorization(t.Context(), "a", sdk.GrantCopyAuthorizationRequest{
		TargetAzureResourceID: "/subscriptions/target", TargetRegion: "westus",
	})
	require.NoError(t, err)
	assert.Equal(t, "/subscriptions/target", grant.TargetAzureResourceID)
	poller, err := client.BeginCopyAnalyzer(t.Context(), "b", sdk.CopyAnalyzerRequest{
		SourceAnalyzerID: "a", SourceAzureResourceID: "/subscriptions/source", SourceRegion: "eastus",
	})
	require.NoError(t, err)
	result, err := poller.PollUntilDone(t.Context(), &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
	require.NoError(t, err)
	assert.Equal(t, "b", result.AnalyzerID)
}

func TestFailedAndCanceledPolling(t *testing.T) {
	for _, state := range []string{"Failed", "Canceled"} {
		t.Run(state, func(t *testing.T) {
			var server *httptest.Server
			client, srv := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzerResults/op?api-version=2025-11-01")
					w.WriteHeader(http.StatusAccepted)
					_, _ = io.WriteString(w, `{"status":"Running"}`)
					return
				}
				_, _ = io.WriteString(w, `{"status":"`+state+`","error":{"code":"AnalysisFailed","message":"failed"}}`)
			})
			server = srv
			poller, err := client.BeginAnalyze(t.Context(), "a", sdk.AnalyzeRequest{Inputs: []sdk.AnalysisInput{{URL: "https://example.test/a"}}})
			require.NoError(t, err)
			_, err = poller.PollUntilDone(t.Context(), &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
			require.Error(t, err)
		})
	}
}

func TestContinuationOriginAndCancellation(t *testing.T) {
	client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"value":[],"nextLink":"https://other.test/next"}`)
	})
	pager := client.NewListAnalyzersPager()
	_, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	_, err = pager.NextPage(t.Context())
	require.ErrorContains(t, err, "resource endpoint")
	client, _ = testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Operation-Location", "https://other.test/poll")
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, `{"status":"Running"}`)
	})
	poller, err := client.BeginAnalyze(t.Context(), "a", sdk.AnalyzeRequest{Inputs: []sdk.AnalysisInput{{URL: "https://example.test/a"}}})
	require.NoError(t, err)
	_, err = poller.Poll(t.Context())
	require.ErrorContains(t, err, "resource endpoint")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = client.GetDefaults(ctx)
	require.True(t, errors.Is(err, context.Canceled))
}

type testCredential struct {
	scopes []string
}

func (c *testCredential) GetToken(_ context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.scopes = options.Scopes
	return azcore.AccessToken{Token: "fake-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestBearerAuthAndConstructorValidation(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer fake-token", r.Header.Get("Authorization"))
		_, _ = io.WriteString(w, `{"modelDeployments":{}}`)
	}))
	defer server.Close()
	credential := &testCredential{}
	client, err := sdk.NewClient(server.URL, credential, &sdk.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: server.Client()}})
	require.NoError(t, err)
	_, err = client.GetDefaults(t.Context())
	require.NoError(t, err)
	assert.Equal(t, []string{"https://cognitiveservices.azure.com/.default"}, credential.scopes)
	for _, endpoint := range []string{"", "http://example.test", "https://user:pass@example.test", "https://example.test/path", "https://example.test?key=x"} {
		_, err := sdk.NewClientWithKeyCredential(endpoint, azcore.NewKeyCredential("key"), nil)
		require.Error(t, err)
	}
	_, err = sdk.NewClient(server.URL, nil, nil)
	require.Error(t, err)
	_, err = sdk.NewClientWithKeyCredential(server.URL, nil, nil)
	require.Error(t, err)
}

func TestMalformedResponse(t *testing.T) {
	for _, payload := range []string{`{"invalid"`, `null`, `[]`, ``} {
		t.Run(payload, func(t *testing.T) {
			client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.Copy(w, bytes.NewBufferString(payload))
			})
			_, err := client.GetDefaults(t.Context())
			require.Error(t, err)
		})
	}
}

func TestServiceErrorCodes(t *testing.T) {
	for _, status := range []int{401, 403, 404, 409} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			client, _ := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = io.WriteString(w, `{"error":{"code":"ServiceCode","message":"service error"}}`)
			})
			_, err := client.GetAnalyzer(t.Context(), "a")
			var responseErr *azcore.ResponseError
			require.ErrorAs(t, err, &responseErr)
			assert.Equal(t, status, responseErr.StatusCode)
			assert.Equal(t, "ServiceCode", responseErr.ErrorCode)
		})
	}
}

func TestRetryResendsBodyAndRetainsRequestID(t *testing.T) {
	var calls int
	var requestID string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "input", string(body))
		assert.NotEmpty(t, r.Header.Get("x-ms-client-request-id"))
		if calls == 1 {
			requestID = r.Header.Get("x-ms-client-request-id")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = io.WriteString(w, `{"error":{"code":"Unavailable"}}`)
			return
		}
		assert.Equal(t, requestID, r.Header.Get("x-ms-client-request-id"))
		_, _ = io.WriteString(w, succeededAnalysis)
	}))
	defer server.Close()
	client, err := sdk.NewClientWithKeyCredential(server.URL, azcore.NewKeyCredential("key"), &sdk.ClientOptions{
		APIVersion: sdk.APIVersionPreview,
		ClientOptions: azcore.ClientOptions{Transport: server.Client(), Retry: policy.RetryOptions{
			MaxRetries: 1, RetryDelay: time.Millisecond, MaxRetryDelay: time.Millisecond,
		}},
	})
	require.NoError(t, err)
	_, err = client.AnalyzeBinaryInline(t.Context(), "a", streamingBody("input"), "")
	require.NoError(t, err)
	assert.Equal(t, 2, calls)
}

func TestPollingHonorsRetryAfterAndDeadline(t *testing.T) {
	var calls atomic.Int32
	var server *httptest.Server
	client, srv := testClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Operation-Location", server.URL+"/contentunderstanding/analyzerResults/op")
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, `{"status":"Running"}`)
	})
	server = srv
	poller, err := client.BeginAnalyze(t.Context(), "a", sdk.AnalyzeRequest{Inputs: []sdk.AnalysisInput{{URL: "https://example.test/a"}}})
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Millisecond)
	defer cancel()
	_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, int32(1), calls.Load(), "Retry-After must delay polling until the deadline expires")
}

func TestAPIVersionOverridesApplyToContinuations(t *testing.T) {
	for _, explicit := range []bool{false, true} {
		t.Run(map[bool]string{false: "embedded", true: "explicit"}[explicit], func(t *testing.T) {
			var server *httptest.Server
			var calls int
			server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				assert.Equal(t, sdk.APIVersionPreview, r.URL.Query().Get("api-version"))
				if calls == 1 {
					_, _ = io.WriteString(w, `{"value":[],"nextLink":"`+server.URL+`/contentunderstanding/analyzers?api-version=2025-11-01"}`)
				} else {
					_, _ = io.WriteString(w, `{"value":[]}`)
				}
			}))
			defer server.Close()
			options := &sdk.ClientOptions{ClientOptions: azcore.ClientOptions{
				Transport: server.Client(), APIVersion: sdk.APIVersionPreview,
			}}
			if explicit {
				options.ClientOptions.APIVersion = sdk.APIVersionGA
				options.APIVersion = sdk.APIVersionPreview
			}
			client, err := sdk.NewClientWithKeyCredential(server.URL, azcore.NewKeyCredential("key"), options)
			require.NoError(t, err)
			pager := client.NewListAnalyzersPager()
			for pager.More() {
				_, err := pager.NextPage(t.Context())
				require.NoError(t, err)
			}
			assert.Equal(t, 2, calls)
			_, err = client.AnalyzeInline(t.Context(), "a", sdk.AnalyzeRequest{})
			require.ErrorContains(t, err, "at least one", "preview version must not be rejected as GA")
		})
	}
}

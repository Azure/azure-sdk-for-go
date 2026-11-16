// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armsovereign

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type regressionTransport func(*http.Request) (*http.Response, error)

func (f regressionTransport) Do(req *http.Request) (*http.Response, error) { return f(req) }

func regressionResponse(req *http.Request, status int, body string) *http.Response {
	return &http.Response{Request: req, StatusCode: status,
		Header: http.Header{"Content-Type": {"application/json"}},
		Body:   io.NopCloser(strings.NewReader(body))}
}

func TestRemediatePublicPoller(t *testing.T) {
	for _, tc := range []struct {
		name, terminal, wantID, wantError string
		immediate, resume, location       bool
	}{
		{name: "terminal envelope", terminal: `{"status":"Succeeded","properties":{"remediationId":"terminal-id"}}`, wantID: "terminal-id"},
		{name: "resumed envelope", terminal: `{"status":"Succeeded","properties":{"remediationId":"terminal-id"}}`, wantID: "terminal-id", resume: true},
		{name: "flat", terminal: `{"status":"Succeeded","remediationId":"flat-id"}`, wantID: "flat-id"},
		{name: "immediate flat", terminal: `{"remediationId":"flat-id"}`, wantID: "flat-id", immediate: true},
		{name: "immediate envelope", terminal: `{"status":"Succeeded","properties":{"remediationId":"terminal-id"}}`, wantID: "terminal-id", immediate: true},
		{name: "location final", terminal: `{"remediationId":"location-id"}`, wantID: "location-id", location: true},
		{name: "failed", terminal: `{"status":"Failed","error":{"code":"BackendFailed","message":"not remediated"}}`, wantError: "BackendFailed"},
		{name: "cancelled", terminal: `{"status":"Canceled","error":{"code":"Canceled","message":"cancelled"}}`, wantError: "Canceled", resume: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			posts, polls, finals := 0, 0, 0
			transport := regressionTransport(func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodPost {
					posts++
					if req.URL.Query().Get("api-version") != "2026-10-31-preview" {
						t.Fatalf("unexpected API version: %s", req.URL)
					}
					if tc.immediate {
						return regressionResponse(req, 200, tc.terminal), nil
					}
					resp := regressionResponse(req, 202, `{"remediationId":"initial-id"}`)
					resp.Header.Set("Azure-AsyncOperation", "https://management.azure.com/operations/job?api-version=2026-10-31-preview")
					resp.Header.Set("Retry-After", "0")
					if tc.location {
						resp.Header.Set("Location", "https://management.azure.com/results/job?api-version=2026-10-31-preview")
					}
					return resp, nil
				}
				if strings.Contains(req.URL.Path, "/results/") {
					finals++
					return regressionResponse(req, 200, tc.terminal), nil
				}
				polls++
				if polls == 1 {
					return regressionResponse(req, 200, `{"status":"Running"}`), nil
				}
				body := tc.terminal
				if tc.location {
					body = `{"status":"Succeeded"}`
				}
				return regressionResponse(req, 200, body), nil
			})
			client, err := NewViewOperationsClient("subscription", nil, &arm.ClientOptions{
				ClientOptions: policy.ClientOptions{Transport: transport, InsecureAllowCredentialWithHTTP: true,
					Retry: policy.RetryOptions{MaxRetries: -1}},
			})
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			poller, err := client.BeginRemediate(ctx, "rg", "view", ViewRemediateContent{}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !tc.immediate {
				if _, err = poller.Poll(ctx); err != nil {
					t.Fatal(err)
				}
				if poller.Done() {
					t.Fatal("Running was treated as terminal")
				}
			}
			if tc.resume {
				token, err := poller.ResumeToken()
				if err != nil {
					t.Fatal(err)
				}
				poller, err = client.BeginRemediate(ctx, "rg", "view", ViewRemediateContent{},
					&ViewOperationsClientBeginRemediateOptions{ResumeToken: token})
				if err != nil {
					t.Fatal(err)
				}
			}
			result, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
			if tc.wantError != "" {
				var responseError *azcore.ResponseError
				if !errors.As(err, &responseError) || responseError.ErrorCode != tc.wantError {
					t.Fatalf("expected %s, got %v", tc.wantError, err)
				}
			} else if err != nil || result.RemediationID == nil || *result.RemediationID != tc.wantID {
				t.Fatalf("public result ID=%v, error=%v; want %q", result.RemediationID, err, tc.wantID)
			}
			if posts != 1 || (!tc.immediate && polls != 2) || (tc.location && finals != 1) {
				t.Fatalf("unexpected requests: posts=%d polls=%d finals=%d", posts, polls, finals)
			}
		})
	}
}

func TestEvaluatedAtTypedRoundTrip(t *testing.T) {
	for _, timestamp := range []string{
		"2026-11-08T20:19:50.7989851Z",
		"2026-11-08T20:19:50.1234567+05:30",
		"2026-11-08T20:19:50.0000001Z",
		"2026-11-08T20:19:50.123456789Z",
	} {
		t.Run(timestamp, func(t *testing.T) {
			var details ViewComplianceDetailsResult
			if err := json.Unmarshal([]byte(`{"resources":[{"dimensionDetails":[{"evaluatedAt":"`+timestamp+`"}]}]}`), &details); err != nil {
				t.Fatal(err)
			}

			value := details.Resources[0].DimensionDetails[0].EvaluatedAt
			request := ViewComplianceResourceDetailsContent{EvaluatedAt: value}
			body, err := json.Marshal(request)
			if err != nil {
				t.Fatal(err)
			}
			var wire struct{ EvaluatedAt string }
			if err := json.Unmarshal(body, &wire); err != nil {
				t.Fatal(err)
			}
			before, _ := time.Parse(time.RFC3339Nano, timestamp)
			after, err := time.Parse(time.RFC3339Nano, wire.EvaluatedAt)
			if err != nil || !before.Equal(after) || before.Nanosecond() != after.Nanosecond() {
				t.Fatalf("timestamp lost precision: %s -> %s (%v)", timestamp, body, err)
			}
		})
	}
}

func TestRemediateResponseShapes(t *testing.T) {
	for _, tc := range []struct {
		name, body, wantID string
		wantError          bool
	}{
		{name: "flat wins", body: `{"remediationId":"flat","properties":{"remediationId":"nested"}}`, wantID: "flat"},
		{name: "missing stays missing", body: `{"status":"Succeeded","properties":{}}`},
		{name: "null", body: `{"properties":null}`},
		{name: "invalid nested ID", body: `{"properties":{"remediationId":42}}`, wantError: true},
		{name: "invalid JSON", body: `{`, wantError: true},
		{name: "invalid properties", body: `{"properties":[]}`, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := ViewOperationsClientRemediateResponse{ViewRemediateResult: ViewRemediateResult{RemediationID: ptr("stale")}}
			err := json.Unmarshal([]byte(tc.body), &r)
			if (err != nil) != tc.wantError {
				t.Fatalf("error=%v; want error=%v", err, tc.wantError)
			}
			if err == nil && ((tc.wantID == "" && r.RemediationID != nil) ||
				(tc.wantID != "" && (r.RemediationID == nil || *r.RemediationID != tc.wantID))) {
				t.Fatalf("ID=%v; want %q", r.RemediationID, tc.wantID)
			}
		})
	}
}

func ptr(value string) *string { return &value }

func TestListContinuationProviderError(t *testing.T) {
	calls := 0
	client, err := NewViewOperationsClient("subscription", nil, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Retry: policy.RetryOptions{MaxRetries: -1},
			Transport: regressionTransport(func(req *http.Request) (*http.Response, error) {
				calls++
				if calls == 1 {
					return regressionResponse(req, 200, `{"value":[{"id":"/views/one"}],"nextLink":"https://management.azure.com/subscriptions/subscription/providers/Microsoft.Sovereign/sovereignViews?api-version=2026-10-31-preview&$skiptoken=a%2Bb%2Fc%3D"}`), nil
				}
				if req.URL.Query().Get("$skiptoken") != "a+b/c=" {
					t.Fatalf("continuation changed: %s", req.URL)
				}
				return regressionResponse(req, 502, `{"error":{"code":"ProviderError","message":"page two failed"}}`), nil
			})},
	})
	if err != nil {
		t.Fatal(err)
	}
	pager := client.NewListBySubscriptionPager(nil)
	if _, err := pager.NextPage(context.Background()); err != nil {
		t.Fatal(err)
	}
	_, err = pager.NextPage(context.Background())
	var responseError *azcore.ResponseError
	if !errors.As(err, &responseError) || responseError.StatusCode != 502 ||
		responseError.ErrorCode != "ProviderError" || calls != 2 {
		t.Fatalf("continuation error was not preserved: %v (calls=%d)", err, calls)
	}
}

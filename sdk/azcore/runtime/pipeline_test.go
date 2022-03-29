//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

// policy that tracks the number of times it was invoked
type countingPolicy struct {
	count int
}

func (p *countingPolicy) Do(req *policy.Request) (*http.Response, error) {
	p.count++
	return req.Next()
}

func TestNewPipelineTelemetry(t *testing.T) {
	for _, disabled := range []bool{true, false} {
		name := "enabled"
		if disabled {
			name = "disabled"
		}
		t.Run(name, func(t *testing.T) {
			srv, close := mock.NewServer()
			defer close()
			srv.AppendResponse()
			opt := policy.ClientOptions{Telemetry: policy.TelemetryOptions{Disabled: disabled}, Transport: srv}
			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			module := "test"
			version := "v1.2.3"
			resp, err := NewPipeline(module, version, PipelineOptions{}, &opt).Do(req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			ua := resp.Request.Header.Get("User-Agent")
			if (!disabled && !strings.HasPrefix(ua, fmt.Sprintf("azsdk-go-%s/%s", module, version))) || (disabled && ua != "") {
				t.Fatalf("Unexpected User-Agent %s", ua)
			}
		})
	}
}

func TestNewPipelineCustomTelemetry(t *testing.T) {
	const appID = "something"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse()
	opts := policy.ClientOptions{Transport: srv, Telemetry: policy.TelemetryOptions{ApplicationID: appID}}
	if opts.Telemetry.ApplicationID != appID {
		t.Fatalf("telemetry was modified: %s", opts.Telemetry.ApplicationID)
	}
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := NewPipeline("armtest", "v1.2.3", PipelineOptions{}, &opts).Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if ua := resp.Request.Header.Get("User-Agent"); !strings.HasPrefix(ua, appID+" "+"azsdk-go-armtest/v1.2.3") {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}

func TestNewPipelineCustomPolicies(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	opts := policy.ClientOptions{Transport: srv, Retry: policy.RetryOptions{RetryDelay: time.Microsecond, MaxRetries: 1}}
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	perCallPolicy := &countingPolicy{}
	perRetryPolicy := &countingPolicy{}
	pl := NewPipeline("",
		"",
		PipelineOptions{PerCall: []pipeline.Policy{perCallPolicy}, PerRetry: []pipeline.Policy{perRetryPolicy}},
		&opts,
	)
	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if perCallPolicy.count != 1 {
		t.Fatalf("Per call policy received %d requests instead of the expected 1", perCallPolicy.count)
	}
	if perRetryPolicy.count != 2 {
		t.Fatalf("Per call policy received %d requests instead of the expected 2", perRetryPolicy.count)
	}
}

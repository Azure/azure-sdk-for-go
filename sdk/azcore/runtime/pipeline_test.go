// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

// policy that tracks the number of times it was invoked
type countingPolicy struct {
	count    int
	callback func()
}

func TestPipelineSanitizer(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	endpoint := srv.URL() + "?client=visible&pipeline=visible&sig=secret"
	expectedErr := &url.Error{Op: "Get", URL: endpoint, Err: errors.New("failed")}
	srv.SetError(expectedErr)
	messages := map[log.Event]string{}
	log.SetListener(func(event log.Event, message string) { messages[event] += message })
	t.Cleanup(func() { log.SetListener(nil) })
	options := policy.ClientOptions{
		Transport: srv,
		Retry:     policy.RetryOptions{MaxRetries: -1},
		Logging: policy.LogOptions{
			AllowedHeaders:     []string{"Client-Header"},
			AllowedQueryParams: []string{"client"},
		},
	}
	pl := NewPipeline("module", "v1.0.0", PipelineOptions{
		AllowedHeaders:         []string{"Pipeline-Header"},
		AllowedQueryParameters: []string{"pipeline"},
	}, &options)
	s := exported.PipelineSanitizer(pl)
	require.NotNil(t, s)
	require.Same(t, s, newRetryPolicy(&options.Retry, s).(*retryPolicy).sanitizer)
	require.Same(t, s, newLogPolicy(&options.Logging, s).(*logPolicy).sanitizer)
	require.Same(t, s, newHTTPTracePolicy(s).(*httpTracePolicy).sanitizer)
	req, err := NewRequest(context.Background(), http.MethodGet, endpoint)
	require.NoError(t, err)
	req.Raw().Header.Set("Client-Header", "client-value")
	req.Raw().Header.Set("Pipeline-Header", "pipeline-value")
	req.Raw().Header.Set("Secret-Header", "header-secret")
	_, err = pl.Do(req)
	require.Same(t, expectedErr, err)
	for _, event := range []log.Event{log.EventRetryPolicy, log.EventRequest, log.EventResponse} {
		require.Contains(t, messages[event], "client=visible&pipeline=visible&sig=REDACTED")
		require.NotContains(t, messages[event], "secret")
	}
	require.Contains(t, messages[log.EventRequest], "Client-Header: client-value")
	require.Contains(t, messages[log.EventRequest], "Pipeline-Header: pipeline-value")
	require.Equal(t, []string{"client"}, options.Logging.AllowedQueryParams)
	require.Equal(t, []string{"Client-Header"}, options.Logging.AllowedHeaders)
}

func (p *countingPolicy) Do(req *policy.Request) (*http.Response, error) {
	p.count++
	if p.callback != nil {
		p.callback()
	}
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
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)

	// NewPipeline should place policies from ClientOptions (i.e. application-specified policies)
	// after policies from PipelineOptions (i.e. client default policies)
	defaultPerCallPolicy := &countingPolicy{}
	defaultPerRetryPolicy := &countingPolicy{}
	customPerCallPolicy := &countingPolicy{}
	customPerCallPolicy.callback = func() {
		require.Equal(t, 1, defaultPerCallPolicy.count)
	}
	customPerRetryPolicy := &countingPolicy{}
	customPerRetryPolicy.callback = func() {
		require.Equal(t, 1, defaultPerCallPolicy.count)
		require.Equal(t, 1, customPerCallPolicy.count)
		require.GreaterOrEqual(t, defaultPerRetryPolicy.count, 1)
	}

	pl := NewPipeline("",
		"",
		PipelineOptions{PerCall: []policy.Policy{defaultPerCallPolicy}, PerRetry: []policy.Policy{defaultPerRetryPolicy}},
		&policy.ClientOptions{
			Transport:        srv,
			Retry:            policy.RetryOptions{RetryDelay: time.Microsecond, MaxRetries: 1},
			PerCallPolicies:  []policy.Policy{customPerCallPolicy},
			PerRetryPolicies: []policy.Policy{customPerRetryPolicy},
		},
	)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, defaultPerCallPolicy.count)
	require.Equal(t, 1, customPerCallPolicy.count)
	require.Equal(t, 2, defaultPerRetryPolicy.count)
	require.Equal(t, 2, customPerRetryPolicy.count)
}

func TestPipelineDoConcurrent(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()

	pl := NewPipeline("TestPipelineDoConcurrent", shared.Version, PipelineOptions{}, nil)

	plErr := make(chan error, 1)
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			if err != nil {
				// test bug
				panic(err)
			}
			_, err = pl.Do(req)
			if err != nil {
				select {
				case plErr <- err:
					// set error
				default:
					// pending error
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	select {
	case err := <-plErr:
		t.Fatal(err)
	default:
		// no error
	}
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPipelineWithAPIVersion(t *testing.T) {
	version := "42"
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse()
	pl, err := NewPipeline("...", "...", mockCredential{}, azruntime.PipelineOptions{}, &armpolicy.ClientOptions{
		ClientOptions: policy.ClientOptions{
			APIVersion: version,
			Transport:  srv,
		},
	})
	require.NoError(t, err)
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	res, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, version, res.Request.URL.Query().Get(string("api-version")))
}

func TestNewPipelineWithOptions(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse()
	opt := armpolicy.ClientOptions{}
	opt.Transport = srv
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	pl, err := NewPipeline("armtest", "v1.2.3", mockCredential{}, azruntime.PipelineOptions{}, &opt)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if ua := resp.Request.Header.Get("User-Agent"); !strings.HasPrefix(ua, "azsdk-go-armtest/v1.2.3") {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}

func TestNewPipelineWithCustomTelemetry(t *testing.T) {
	const myTelemetry = "something"
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse()
	opt := armpolicy.ClientOptions{}
	opt.Transport = srv
	opt.Telemetry.ApplicationID = myTelemetry
	if opt.Telemetry.ApplicationID != myTelemetry {
		t.Fatalf("telemetry was modified: %s", opt.Telemetry.ApplicationID)
	}
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	pl, err := NewPipeline("armtest", "v1.2.3", mockCredential{}, azruntime.PipelineOptions{}, &opt)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if ua := resp.Request.Header.Get("User-Agent"); !strings.HasPrefix(ua, myTelemetry+" "+"azsdk-go-armtest/v1.2.3") {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}

func TestDisableAutoRPRegistration(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	// initial response that RP is unregistered
	srv.SetResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp1)))
	opts := &armpolicy.ClientOptions{DisableRPRegistration: true, ClientOptions: policy.ClientOptions{Transport: srv}}
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// log only RP registration
	log.SetEvents(LogRPRegistration)
	defer func() {
		// reset logging
		log.SetEvents()
	}()
	logEntries := 0
	log.SetListener(func(cls log.Event, msg string) {
		logEntries++
	})
	pl, err := NewPipeline("armtest", "v1.2.3", mockCredential{}, azruntime.PipelineOptions{}, opts)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status code %d:", resp.StatusCode)
	}
	// shouldn't be any log entries
	if logEntries != 0 {
		t.Fatalf("expected 0 log entries, got %d", logEntries)
	}
}

// policy that tracks the number of times it was invoked
type countingPolicy struct {
	count int
}

func (p *countingPolicy) Do(req *policy.Request) (*http.Response, error) {
	p.count++
	return req.Next()
}

func TestPipelineWithCustomPolicies(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	// initial response is a failure to trigger retry
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	perCallPolicy := countingPolicy{}
	perRetryPolicy := countingPolicy{}
	opts := &armpolicy.ClientOptions{
		DisableRPRegistration: true,
		ClientOptions: policy.ClientOptions{
			PerCallPolicies:  []policy.Policy{&perCallPolicy},
			PerRetryPolicies: []policy.Policy{&perRetryPolicy},
			Retry:            policy.RetryOptions{RetryDelay: time.Microsecond},
			Transport:        srv,
		},
	}
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	pl, err := NewPipeline("armtest", "v1.2.3", mockCredential{}, azruntime.PipelineOptions{}, opts)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if perCallPolicy.count != 1 {
		t.Fatalf("unexpected per call policy count %d", perCallPolicy.count)
	}
	if perRetryPolicy.count != 2 {
		t.Fatalf("unexpected per retry policy count %d", perRetryPolicy.count)
	}
}

func TestPipelineAudience(t *testing.T) {
	for _, c := range []cloud.Configuration{cloud.AzureChina, cloud.AzureGovernment, cloud.AzurePublic} {
		srv, close := mock.NewTLSServer()
		defer close()
		srv.AppendResponse(mock.WithStatusCode(200))
		opts := &armpolicy.ClientOptions{}
		opts.Cloud = c
		opts.Transport = srv
		audience := opts.Cloud.Services[cloud.ResourceManager].Audience
		if !strings.HasPrefix(audience, "https://") {
			t.Fatal("unexpected audience " + audience)
		}
		getTokenCalled := false
		cred := mockCredential{getTokenImpl: func(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
			getTokenCalled = true
			if n := len(options.Scopes); n != 1 {
				t.Fatalf("expected 1 scope, got %d", n)
			}
			if options.Scopes[0] != audience+"/.default" {
				t.Fatalf(`unexpected scope "%s"`, options.Scopes[0])
			}
			return azcore.AccessToken{Token: "...", ExpiresOn: time.Now().Add(time.Hour)}, nil
		}}
		req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
		if err != nil {
			t.Fatal(err)
		}
		pl, err := NewPipeline("test", "v0.1.0", cred, azruntime.PipelineOptions{}, opts)
		if err != nil {
			t.Fatal(err)
		}
		_, err = pl.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if !getTokenCalled {
			t.Fatal("mock credential's GetToken method wasn't called")
		}
	}
}

func TestPipelineWithIncompleteCloudConfig(t *testing.T) {
	partialConfigs := []cloud.Configuration{
		{Services: map[cloud.ServiceName]cloud.ServiceConfiguration{"...": {Endpoint: "..."}}},
		{Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {Audience: "..."},
		}},
		{Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {Endpoint: "http://localhost"},
		}},
	}
	for _, c := range partialConfigs {
		opts := &armpolicy.ClientOptions{}
		opts.Cloud = c
		_, err := NewPipeline("test", "v0.1.0", mockCredential{}, azruntime.PipelineOptions{}, opts)
		if err == nil {
			t.Fatal("expected an error")
		}
	}
}

func TestPipelineDoConcurrent(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse()

	pl, err := NewPipeline(
		"TestPipelineDoConcurrent",
		shared.Version,
		mockCredential{},
		azruntime.PipelineOptions{},
		&armpolicy.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Transport: srv,
			},
		},
	)
	require.NoError(t, err)

	plErr := make(chan error, 1)
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
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

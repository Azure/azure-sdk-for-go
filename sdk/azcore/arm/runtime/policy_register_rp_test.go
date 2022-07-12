//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const rpUnregisteredResp = `{
	"error":{
		"code":"MissingSubscriptionRegistration",
		"message":"The subscription registration is in 'Unregistered' state. The subscription must be registered to use namespace 'Microsoft.Storage'. See https://aka.ms/rps-not-found for how to register subscriptions.",
		"details":[{
				"code":"MissingSubscriptionRegistration",
				"target":"Microsoft.Storage",
				"message":"The subscription registration is in 'Unregistered' state. The subscription must be registered to use namespace 'Microsoft.Storage'. See https://aka.ms/rps-not-found for how to register subscriptions."
			}
		]
	}
}`

// some content was omitted here as it's not relevant
const rpRegisteringResp = `{
    "id": "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage",
    "namespace": "Microsoft.Storage",
    "registrationState": "Registering",
    "registrationPolicy": "RegistrationRequired"
}`

// some content was omitted here as it's not relevant
const rpRegisteredResp = `{
    "id": "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage",
    "namespace": "Microsoft.Storage",
    "registrationState": "Registered",
    "registrationPolicy": "RegistrationRequired"
}`

const requestEndpoint = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/fakeResourceGroupo/providers/Microsoft.Storage/storageAccounts/fakeAccountName"

func newTestRPRegistrationPipeline(t *testing.T, srv *mock.Server) runtime.Pipeline {
	opts := testRPRegistrationOptions(srv)
	rp, err := NewRPRegistrationPolicy(mockCredential{}, testRPRegistrationOptions(srv))
	if err != nil {
		t.Fatal(err)
	}
	return runtime.NewPipeline("test", "v0.1.0", runtime.PipelineOptions{PerCall: []azpolicy.Policy{rp}}, &opts.ClientOptions)
}

func testRPRegistrationOptions(srv *mock.Server) *armpolicy.RegistrationOptions {
	def := armpolicy.RegistrationOptions{}
	def.Cloud = cloud.Configuration{
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {Endpoint: srv.URL(), Audience: srv.URL()},
		}}
	def.Transport = srv
	def.PollingDelay = 100 * time.Millisecond
	def.PollingDuration = 1 * time.Second
	return &def
}

func TestRPRegistrationPolicySuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	// polling responses to Register() and Get(), in progress
	srv.RepeatResponse(5, mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteringResp)))
	// polling response, successful registration
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteredResp)))
	// response for original request (different status code than any of the other responses)
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	pl := newTestRPRegistrationPipeline(t, srv)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, runtime.JoinPaths(srv.URL(), requestEndpoint))
	if err != nil {
		t.Fatal(err)
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
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code %d:", resp.StatusCode)
	}
	if resp.Request.URL.Path != requestEndpoint {
		t.Fatalf("unexpected path in response %s", resp.Request.URL.Path)
	}
	// should be four entries
	// 1st is for start
	// 2nd is for first response to get state
	// 3rd is when state transitions to success
	// 4th is for end
	if logEntries != 4 {
		t.Fatalf("expected 4 log entries, got %d", logEntries)
	}
}

func TestRPRegistrationPolicyNA(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// response indicates no RP registration is required, policy does nothing
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	pl := newTestRPRegistrationPipeline(t, srv)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	// log only RP registration
	log.SetEvents(LogRPRegistration)
	defer func() {
		// reset logging
		log.SetEvents()
	}()
	log.SetListener(func(cls log.Event, msg string) {
		t.Fatalf("unexpected log entry %s: %s", cls, msg)
	})
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestRPRegistrationPolicy409Other(t *testing.T) {
	const failedResp = `{
		"error":{
			"code":"CannotDoTheThing",
			"message":"Something failed in your API call.",
			"details":[{
					"code":"ThisIsForTesting",
					"message":"This is fake."
				}
			]
		}
	}`
	srv, close := mock.NewServer()
	defer close()
	// test getting a 409 but not due to registration required
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(failedResp)))
	pl := newTestRPRegistrationPipeline(t, srv)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	// log only RP registration
	log.SetEvents(LogRPRegistration)
	defer func() {
		// reset logging
		log.SetEvents()
	}()
	log.SetListener(func(cls log.Event, msg string) {
		t.Fatalf("unexpected log entry %s: %s", cls, msg)
	})
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestRPRegistrationPolicyTimesOut(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	// polling responses to Register() and Get(), in progress but slow
	// tests registration takes too long, times out
	srv.RepeatResponse(10, mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteringResp)), mock.WithSlowResponse(400*time.Millisecond))
	pl := newTestRPRegistrationPipeline(t, srv)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, runtime.JoinPaths(srv.URL(), requestEndpoint))
	if err != nil {
		t.Fatal(err)
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
	resp, err := pl.Do(req)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	// should be three entries
	// 1st is for start
	// 2nd is for first response to get state
	// 3rd is the deadline exceeded error
	if logEntries != 3 {
		t.Fatalf("expected 3 log entries, got %d", logEntries)
	}
	// we should get the response from the original request
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestRPRegistrationPolicyExceedsAttempts(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// add a cycle of unregistered->registered so that we keep retrying and hit the cap
	for i := 0; i < 4; i++ {
		// initial response that RP is unregistered
		srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
		// polling responses to Register() and Get(), in progress
		srv.RepeatResponse(2, mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteringResp)))
		// polling response, successful registration
		srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteredResp)))
	}
	pl := newTestRPRegistrationPipeline(t, srv)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, runtime.JoinPaths(srv.URL(), requestEndpoint))
	if err != nil {
		t.Fatal(err)
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
	resp, err := pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if !strings.HasPrefix(err.Error(), "exceeded attempts to register Microsoft.Storage") {
		t.Fatalf("unexpected error message %s", err.Error())
	}
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status code %d:", resp.StatusCode)
	}
	if resp.Request.URL.Path != requestEndpoint {
		t.Fatalf("unexpected path in response %s", resp.Request.URL.Path)
	}
	// should be 4 entries for each attempt, total 12 entries
	// 1st is for start
	// 2nd is for first response to get state
	// 3rd is when state transitions to success
	// 4th is for end
	if logEntries != 12 {
		t.Fatalf("expected 12 log entries, got %d", logEntries)
	}
}

// test cancelling registration
func TestRPRegistrationPolicyCanCancel(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	// polling responses to Register() and Get(), in progress but slow so we have time to cancel
	srv.RepeatResponse(10, mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteringResp)), mock.WithSlowResponse(300*time.Millisecond))
	opts := armpolicy.RegistrationOptions{}
	opts.Transport = srv
	pl := newTestRPRegistrationPipeline(t, srv)
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

	wg := &sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	var resp *http.Response
	var err error
	go func() {
		defer wg.Done()
		// create request and start pipeline
		var req *azpolicy.Request
		req, err = runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(srv.URL(), requestEndpoint))
		if err != nil {
			return
		}
		resp, err = pl.Do(req)
	}()

	// wait for a bit then cancel the operation
	time.Sleep(500 * time.Millisecond)
	cancel()
	wg.Wait()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected Canceled error, got %v", err)
	}
	// there should be 1 or 2 entries depending on the timing
	if logEntries == 0 {
		t.Fatal("didn't get any log entries")
	}
	// should have original response
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestRPRegistrationPolicyDisabled(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	ops := testRPRegistrationOptions(srv)
	ops.MaxAttempts = -1
	rp, err := NewRPRegistrationPolicy(mockCredential{}, ops)
	if err != nil {
		t.Fatal(err)
	}
	pl := runtime.NewPipeline("test", "v0.1.0", runtime.PipelineOptions{PerCall: []azpolicy.Policy{rp}}, nil)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, runtime.JoinPaths(srv.URL(), requestEndpoint))
	if err != nil {
		t.Fatal(err)
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

func TestRPRegistrationPolicyAudience(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	// polling responses to Register() and Get(), in progress
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteringResp)))
	// polling response, successful registration
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteredResp)))
	// response for original request (different status code than any of the other responses)
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))

	audience := "audience"
	conf := cloud.Configuration{
		ActiveDirectoryAuthorityHost: srv.URL(),
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {Audience: audience, Endpoint: srv.URL()},
		},
	}
	getTokenCalled := false
	cred := mockCredential{getTokenImpl: func(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
		getTokenCalled = true
		if n := len(options.Scopes); n != 1 {
			t.Fatalf("expected 1 scope, got %d", n)
		}
		if options.Scopes[0] != audience+"/.default" {
			t.Fatalf(`unexpected scope "%s"`, options.Scopes[0])
		}
		return azcore.AccessToken{Token: "...", ExpiresOn: time.Now().Add(time.Hour)}, nil
	}}
	opts := azpolicy.ClientOptions{Cloud: conf, Transport: srv}
	rp, err := NewRPRegistrationPolicy(cred, &armpolicy.RegistrationOptions{ClientOptions: opts})
	if err != nil {
		t.Fatal(err)
	}
	pl := runtime.NewPipeline("test", "v0.1.0", runtime.PipelineOptions{PerCall: []azpolicy.Policy{rp}}, &azpolicy.ClientOptions{Transport: srv})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL()+requestEndpoint)
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

func TestRPRegistrationPolicyWithIncompleteCloudConfig(t *testing.T) {
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
		opts := azpolicy.ClientOptions{Cloud: c}
		_, err := NewRPRegistrationPolicy(mockCredential{}, &armpolicy.RegistrationOptions{ClientOptions: opts})
		if err == nil {
			t.Fatal("expected an error")
		}
	}
}

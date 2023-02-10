//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
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

const rpEnvsInSubExceeded = `{
	"code": "MaxNumberOfRegionalEnvironmentsInSubExceeded",
	"message": "The subscription '00000000-0000-0000-0000-000000000000' cannot have more than 1 environments in East US."
}`

const requestEndpoint = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/fakeResourceGroupo/providers/Microsoft.Storage/storageAccounts/fakeAccountName"

const fakeAPIBody = "success"

func newTestRPRegistrationPipeline(t *testing.T, srv *mock.Server, opts *armpolicy.RegistrationOptions) runtime.Pipeline {
	if opts == nil {
		opts = testRPRegistrationOptions(srv)
	}
	rp, err := NewRPRegistrationPolicy(mockCredential{}, opts)
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
	// response for original request
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(fakeAPIBody)))
	client := newFakeClient(t, srv, nil)
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
	resp, err := client.FakeAPI(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, fakeAPIBody, resp.Result)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, requestEndpoint, resp.Path)
	// should be four entries
	// 1st is for start
	// 2nd is for first response to get state
	// 3rd is when state transitions to success
	// 4th is for end
	require.EqualValues(t, 4, logEntries)
}

func TestRPRegistrationPolicyNA(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// response indicates no RP registration is required, policy does nothing
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(fakeAPIBody)))
	client := newFakeClient(t, srv, nil)
	// log only RP registration
	log.SetEvents(LogRPRegistration)
	defer func() {
		// reset logging
		log.SetEvents()
	}()
	log.SetListener(func(cls log.Event, msg string) {
		t.Fatalf("unexpected log entry %s: %s", cls, msg)
	})
	resp, err := client.FakeAPI(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, fakeAPIBody, resp.Result)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, requestEndpoint, resp.Path)
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
	client := newFakeClient(t, srv, nil)
	// log only RP registration
	log.SetEvents(LogRPRegistration)
	defer func() {
		// reset logging
		log.SetEvents()
	}()
	log.SetListener(func(cls log.Event, msg string) {
		t.Fatalf("unexpected log entry %s: %s", cls, msg)
	})
	resp, err := client.FakeAPI(context.Background())
	require.Error(t, err)
	require.Zero(t, resp)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.EqualValues(t, "CannotDoTheThing", respErr.ErrorCode)
}

func TestRPRegistrationPolicyTimesOut(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	// polling responses to Register() and Get(), in progress but slow
	// tests registration takes too long, times out
	srv.RepeatResponse(10, mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteringResp)), mock.WithSlowResponse(400*time.Millisecond))
	client := newFakeClient(t, srv, nil)
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
	resp, err := client.FakeAPI(context.Background())
	require.Error(t, err)
	require.Zero(t, resp)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	// should be three entries
	// 1st is for start
	// 2nd is for first response to get state
	// 3rd is the deadline exceeded error
	require.EqualValues(t, 3, logEntries)
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
	client := newFakeClient(t, srv, nil)
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
	resp, err := client.FakeAPI(context.Background())
	require.Error(t, err)
	require.Zero(t, resp)
	require.Contains(t, err.Error(), "exceeded attempts to register Microsoft.Storage")
	// should be 4 entries for each attempt, total 12 entries
	// 1st is for start
	// 2nd is for first response to get state
	// 3rd is when state transitions to success
	// 4th is for end
	require.EqualValues(t, 12, logEntries)
}

// test cancelling registration
func TestRPRegistrationPolicyCanCancel(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	// polling responses to Register() and Get(), in progress but slow so we have time to cancel
	srv.RepeatResponse(10, mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(rpRegisteringResp)), mock.WithSlowResponse(300*time.Millisecond))
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
	var resp fakeResponse
	var err error
	go func() {
		defer wg.Done()
		client := newFakeClient(t, srv, nil)
		resp, err = client.FakeAPI(ctx)
	}()

	// wait for a bit then cancel the operation
	time.Sleep(500 * time.Millisecond)
	cancel()
	wg.Wait()
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, resp)
	// there should be 1 or 2 entries depending on the timing
	require.NotZero(t, logEntries)
}

func TestRPRegistrationPolicyDisabled(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	ops := testRPRegistrationOptions(srv)
	ops.MaxAttempts = -1
	client := newFakeClient(t, srv, ops)
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
	resp, err := client.FakeAPI(context.Background())
	require.Error(t, err)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.EqualValues(t, "MissingSubscriptionRegistration", respErr.ErrorCode)
	require.Zero(t, resp)
	// shouldn't be any log entries
	require.Zero(t, logEntries)
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
		require.Len(t, options.Scopes, 1)
		require.EqualValues(t, audience+"/.default", options.Scopes[0])
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
		require.Error(t, err)
	}
}

func TestRPRegistrationPolicyEnvironmentsInSubExceeded(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// test getting a 409 due to exceeded environments in a subscription
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpEnvsInSubExceeded)))
	client := newFakeClient(t, srv, nil)
	// log only RP registration
	log.SetEvents(LogRPRegistration)
	logEntries := 0
	log.SetListener(func(cls log.Event, msg string) {
		logEntries++
	})
	defer func() {
		// reset logging
		log.SetEvents()
	}()
	resp, err := client.FakeAPI(context.Background())
	require.Error(t, err)
	require.Zero(t, resp)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.EqualValues(t, "MaxNumberOfRegionalEnvironmentsInSubExceeded", respErr.ErrorCode)
	require.Contains(t, err.Error(), "cannot have more than 1 environments")
	require.EqualValues(t, 0, logEntries)
}

type fakeClient struct {
	ep string
	pl runtime.Pipeline
}

func newFakeClient(t *testing.T, srv *mock.Server, opts *armpolicy.RegistrationOptions) *fakeClient {
	return &fakeClient{ep: srv.URL(), pl: newTestRPRegistrationPipeline(t, srv, opts)}
}

type fakeResponse struct {
	Result     string
	StatusCode int
	Path       string
}

// FakeAPI returns fakeResponse with Result "success" on a HTTP 200.
func (f *fakeClient) FakeAPI(ctx context.Context) (fakeResponse, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(f.ep, requestEndpoint))
	if err != nil {
		return fakeResponse{}, err
	}
	resp, err := f.pl.Do(req)
	if err != nil {
		return fakeResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return fakeResponse{}, runtime.NewResponseError(resp)
	}
	body, err := runtime.Payload(resp)
	if err != nil {
		return fakeResponse{}, err
	}
	return fakeResponse{Result: string(body), StatusCode: resp.StatusCode, Path: resp.Request.URL.Path}, nil
}

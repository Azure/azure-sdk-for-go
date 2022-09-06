// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"strings"

	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	tokenValue                = "***"
	accessTokenRespSuccess    = `{"access_token": "` + tokenValue + `", "expires_in": 3600}`
	accessTokenRespShortLived = `{"access_token": "` + tokenValue + `", "expires_in": 0}`
	scope                     = "scope"
)

type mockCredential struct {
	getTokenImpl func(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error)
}

func (mc mockCredential) GetToken(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
	if mc.getTokenImpl != nil {
		return mc.getTokenImpl(ctx, options)
	}
	return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func (mc mockCredential) Do(req *azpolicy.Request) (*http.Response, error) {
	return nil, nil
}

func newTestPipeline(opts *azpolicy.ClientOptions) runtime.Pipeline {
	return runtime.NewPipeline("testmodule", "v0.1.0", runtime.PipelineOptions{}, opts)
}

func defaultTestPipeline(srv azpolicy.Transporter, scope string) (runtime.Pipeline, error) {
	retryOpts := azpolicy.RetryOptions{
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    time.Millisecond,
	}
	return NewPipeline(
		"testmodule",
		"v0.1.0",
		mockCredential{},
		runtime.PipelineOptions{},
		&arm.ClientOptions{
			ClientOptions: azpolicy.ClientOptions{
				Retry:     retryOpts,
				Transport: srv,
			},
		})
}

func TestBearerPolicy_SuccessGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	pipeline, err := defaultTestPipeline(srv, scope)
	if err != nil {
		t.Fatal(err)
	}
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
	const expectedToken = shared.BearerTokenPrefix + tokenValue
	if token := resp.Request.Header.Get(shared.HeaderAuthorization); token != expectedToken {
		t.Fatalf("expected token '%s', got '%s'", expectedToken, token)
	}
}

func TestBearerPolicy_CredentialFailGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	expectedErr := errors.New("oops")
	failCredential := mockCredential{}
	failCredential.getTokenImpl = func(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
		return azcore.AccessToken{}, expectedErr
	}
	b := NewBearerTokenPolicy(failCredential, nil)
	pipeline := newTestPipeline(&azpolicy.ClientOptions{
		Transport: srv,
		Retry: azpolicy.RetryOptions{
			RetryDelay: 10 * time.Millisecond,
		},
		PerRetryPolicies: []azpolicy.Policy{b},
	})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err != expectedErr {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestBearerTokenPolicy_TokenExpired(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespShortLived)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	pipeline, err := defaultTestPipeline(srv, scope)
	if err != nil {
		t.Fatal(err)
	}
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestBearerPolicy_GetTokenFailsNoDeadlock(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	retryOpts := azpolicy.RetryOptions{
		// use a negative try timeout to trigger a deadline exceeded error causing GetToken() to fail
		TryTimeout:    -1 * time.Nanosecond,
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    50 * time.Millisecond,
		MaxRetries:    3,
	}
	b := NewBearerTokenPolicy(mockCredential{}, nil)
	pipeline := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, Retry: retryOpts, PerRetryPolicies: []azpolicy.Policy{b}})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestBearerTokenWithAuxiliaryTenants(t *testing.T) {
	t.Skip("this feature isn't implemented yet")
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse()
	retryOpts := azpolicy.RetryOptions{
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    50 * time.Millisecond,
	}
	b := NewBearerTokenPolicy(
		mockCredential{},
		&armpolicy.BearerTokenOptions{
			Scopes: []string{scope},
			//AuxiliaryTenants: []string{"tenant1", "tenant2", "tenant3"},
		},
	)
	pipeline := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, Retry: retryOpts, PerRetryPolicies: []azpolicy.Policy{b}})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	expectedHeader := strings.Repeat(shared.BearerTokenPrefix+tokenValue+", ", 3)
	expectedHeader = expectedHeader[:len(expectedHeader)-2]
	if auxH := resp.Request.Header.Get(shared.HeaderAuxiliaryAuthorization); auxH != expectedHeader {
		t.Fatalf("unexpected auxiliary authorization header %s", auxH)
	}
}

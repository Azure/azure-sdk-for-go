// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"strings"

	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
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
	if !options.EnableCAE {
		return azcore.AccessToken{}, errors.New("ARM clients should set EnableCAE to true")
	}
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

func defaultTestPipeline(srv azpolicy.Transporter) (runtime.Pipeline, error) {
	retryOpts := azpolicy.RetryOptions{
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    time.Millisecond,
	}
	return NewPipeline(
		"testmodule",
		"v0.1.0",
		mockCredential{},
		runtime.PipelineOptions{},
		&armpolicy.ClientOptions{
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
	pipeline, err := defaultTestPipeline(srv)
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
	expectedErr := "oops"
	failCredential := mockCredential{}
	failCredential.getTokenImpl = func(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
		return azcore.AccessToken{}, errors.New(expectedErr)
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
	require.NoError(t, err)
	resp, err := pipeline.Do(req)
	require.EqualError(t, err, expectedErr)
	require.Nil(t, resp)
	require.Implements(t, (*errorinfo.NonRetriable)(nil), err)
}

func TestBearerTokenPolicy_TokenExpired(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespShortLived)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	pipeline, err := defaultTestPipeline(srv)
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

func TestAuxiliaryTenants(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	primary := "primary"
	auxTenants := []string{"aux1", "aux2", "aux3"}
	expectCache := false
	b := NewBearerTokenPolicy(
		mockCredential{
			// getTokenImpl returns a token whose value equals the requested tenant so the test can validate how the policy handles tenants
			// i.e., primary tenant token goes in Authorization header and aux tenant tokens go in x-ms-authorization-auxiliary
			getTokenImpl: func(ctx context.Context, options azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
				require.False(t, expectCache, "client should have used a cached token instead of requesting another")
				tenant := primary
				if options.TenantID != "" {
					tenant = options.TenantID
				}
				return azcore.AccessToken{Token: tenant, ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
			},
		},
		&armpolicy.BearerTokenOptions{AuxiliaryTenants: auxTenants, Scopes: []string{scope}},
	)
	pipeline := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, PerRetryPolicies: []azpolicy.Policy{b}})
	expected := strings.Split(shared.BearerTokenPrefix+strings.Join(auxTenants, ","+shared.BearerTokenPrefix), ",")
	for i := 0; i < 3; i++ {
		if i == 1 {
			// policy should have a cached token after the first iteration
			expectCache = true
		}
		req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
		require.NoError(t, err)
		resp, err := pipeline.Do(req)
		require.NoError(t, err)
		require.Equal(t, shared.BearerTokenPrefix+primary, resp.Request.Header.Get(shared.HeaderAuthorization), "Authorization header must contain primary tenant token")
		actual := strings.Split(resp.Request.Header.Get(headerAuxiliaryAuthorization), ", ")
		// auxiliary tokens may appear in arbitrary order
		require.ElementsMatch(t, expected, actual)
	}
}

func TestBearerTokenPolicyRequiresHTTPS(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	b := NewBearerTokenPolicy(mockCredential{}, nil)
	pl := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, PerRetryPolicies: []azpolicy.Policy{b}})
	req, err := runtime.NewRequest(context.Background(), "GET", srv.URL())
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.Error(t, err)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
}

func TestBearerTokenPolicyAllowHTTP(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	b := NewBearerTokenPolicy(mockCredential{}, &armpolicy.BearerTokenOptions{
		InsecureAllowCredentialWithHTTP: true,
	})
	pl := newTestPipeline(&azpolicy.ClientOptions{Transport: srv, PerRetryPolicies: []azpolicy.Policy{b}})
	req, err := runtime.NewRequest(context.Background(), "GET", srv.URL())
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
}

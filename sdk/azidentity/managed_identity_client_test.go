//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

type userAgentValidatingPolicy struct {
	t     *testing.T
	appID string
}

func (p userAgentValidatingPolicy) Do(req *policy.Request) (*http.Response, error) {
	expected := "azsdk-go-" + component + "/" + version
	if p.appID != "" {
		expected = p.appID + " " + expected
	}
	if ua := req.Raw().Header.Get("User-Agent"); !strings.HasPrefix(ua, expected) {
		p.t.Fatalf("unexpected User-Agent %s", ua)
	}
	return req.Next()
}

func TestIMDSEndpointParse(t *testing.T) {
	_, err := url.Parse(imdsEndpoint)
	if err != nil {
		t.Fatalf("Failed to parse the IMDS endpoint: %v", err)
	}
}

func TestManagedIdentityClient_UserAgent(t *testing.T) {
	options := ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &mockSTS{}, PerCallPolicies: []policy.Policy{userAgentValidatingPolicy{t: t}},
		},
	}
	client, err := newManagedIdentityClient(&options)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.authenticate(context.Background(), nil, []string{liveTestScope})
	if err != nil {
		t.Fatal(err)
	}
}

func TestManagedIdentityClient_ApplicationID(t *testing.T) {
	appID := "customvalue"
	options := ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &mockSTS{}, PerCallPolicies: []policy.Policy{userAgentValidatingPolicy{t: t, appID: appID}},
		},
	}
	options.Telemetry.ApplicationID = appID
	client, err := newManagedIdentityClient(&options)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.authenticate(context.Background(), nil, []string{liveTestScope})
	if err != nil {
		t.Fatal(err)
	}
}

func TestManagedIdentityClient_IMDSErrors(t *testing.T) {
	for _, test := range []struct {
		body, desc string
		code       int
	}{
		{
			desc: "No identity assigned",
			code: http.StatusBadRequest,
			body: `{"error":"invalid_request","error_description":"Identity not found"}`,
		},
		{
			desc: "Docker Desktop",
			code: http.StatusForbidden,
			body: "connecting to 169.254.169.254:80: connecting to 169.254.169.254:80: dial tcp 169.254.169.254:80: connectex: A socket operation was attempted to an unreachable network.",
		},
		{
			desc: "Docker Desktop",
			code: http.StatusForbidden,
			body: "connecting to 169.254.169.254:80: connecting to 169.254.169.254:80: dial tcp 169.254.169.254:80: connectex: A socket operation was attempted to an unreachable host.",
		},
		{
			desc: "Azure Container Instances",
			code: http.StatusBadRequest,
			body: "Required metadata header not specified or not correct",
		},
	} {
		t.Run(fmt.Sprint(test.code), func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.SetResponse(mock.WithBody([]byte(test.body)), mock.WithStatusCode(test.code))
			client, err := newManagedIdentityClient(&ManagedIdentityCredentialOptions{
				ClientOptions: azcore.ClientOptions{Transport: srv},
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = client.authenticate(context.Background(), nil, testTRO.Scopes)
			if err == nil {
				t.Fatal("expected an error")
			}
			if actual := err.Error(); !strings.Contains(actual, test.body) {
				t.Fatalf("expected response body in error, got %q", actual)
			}
			var unavailableErr credentialUnavailable
			if !errors.As(err, &unavailableErr) {
				t.Fatalf("expected %T, got %T", unavailableErr, err)
			}
		})
	}
}

func TestManagedIdentityClient_IMDSProbeReturnsUnavailableError(t *testing.T) {
	for _, test := range []struct {
		desc string
		res  [][]mock.ResponseOption
	}{
		{
			"Azure Container Instance",
			[][]mock.ResponseOption{
				{
					mock.WithBody([]byte("Required metadata header not specified or not correct")),
					mock.WithStatusCode(http.StatusBadRequest),
				},
				{mock.WithBody([]byte("error")), mock.WithStatusCode(http.StatusBadRequest)},
			},
		},
		{
			"404",
			[][]mock.ResponseOption{
				{mock.WithStatusCode(http.StatusNotFound)},
				{mock.WithStatusCode(http.StatusNotFound)},
			},
		},
		{
			"non-JSON token response",
			[][]mock.ResponseOption{
				{mock.WithBody([]byte("Required metadata header not specified or not correct"))},
				{mock.WithBody([]byte("not json"))},
			},
		},
		{
			"no token in response",
			[][]mock.ResponseOption{
				{mock.WithBody([]byte("Required metadata header not specified or not correct"))},
				{mock.WithBody([]byte(`{"error": "no token here"}`)), mock.WithStatusCode(http.StatusBadRequest)},
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			for _, res := range test.res {
				srv.AppendResponse(res...)
			}
			c, err := newManagedIdentityClient(&ManagedIdentityCredentialOptions{
				ClientOptions: azcore.ClientOptions{
					// disabling retries simplifies the 404 test (404 is retriable for IMDS)
					Retry:     policy.RetryOptions{MaxRetries: -1},
					Transport: srv,
				},
				dac: true,
			})
			require.NoError(t, err)
			_, err = c.authenticate(ctx, nil, testTRO.Scopes)
			var cu credentialUnavailable
			require.ErrorAs(t, err, &cu)

			srv.AppendResponse(
				mock.WithBody(accessTokenRespSuccess),
				mock.WithPredicate(func(r *http.Request) bool {
					if _, ok := r.Header["Metadata"]; !ok {
						t.Error("client shouldn't send another probe after receiving a response")
					}
					return true
				}),
			)
			srv.AppendResponse()
			tk, err := c.authenticate(ctx, nil, testTRO.Scopes)
			require.NoError(t, err)
			require.Equal(t, tokenValue, tk.Token)
		})
	}
}

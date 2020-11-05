// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

type mockTokenCred struct{}

func (mockTokenCred) AuthenticationPolicy(azcore.AuthenticationPolicyOptions) azcore.Policy {
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		return req.Next()
	})
}

func (mockTokenCred) GetToken(context.Context, azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{
		Token:     "abc123",
		ExpiresOn: time.Now().Add(1 * time.Hour),
	}, nil
}

func TestNewDefaultConnection(t *testing.T) {
	opt := DefaultConnectionOptions()
	con := NewDefaultConnection(mockTokenCred{}, &opt)
	if ep := con.Endpoint(); ep != DefaultEndpoint {
		t.Fatalf("unexpected endpoint %s", ep)
	}
}

func TestNewConnection(t *testing.T) {
	const customEndpoint = "https://contoso.com/fake/endpoint"
	con := NewConnection(customEndpoint, mockTokenCred{}, nil)
	if ep := con.Endpoint(); ep != customEndpoint {
		t.Fatalf("unexpected endpoint %s", ep)
	}
}

func TestNewConnectionWithPipeline(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse()
	p := getPipeline(srv)
	con := NewConnectionWithPipeline(srv.URL(), p)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := con.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

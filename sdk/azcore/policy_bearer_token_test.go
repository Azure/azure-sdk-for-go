// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

type dummyCredential struct{}

func (c *dummyCredential) GetToken(ctx context.Context, options TokenRequestOptions) (*AccessToken, error) {
	return &AccessToken{
		Token:     "success_token",
		ExpiresOn: time.Date(2021, 06, 25, 3, 20, 0, 0, time.UTC),
	}, nil
}

func (c *dummyCredential) AuthenticationPolicy(options AuthenticationPolicyOptions) Policy {
	return NewBearerTokenPolicy(c, options)
}

func TestBearerTokenPolicyHTTPFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	cred := &dummyCredential{}
	pl := NewPipeline(srv, NewBearerTokenPolicy(cred, AuthenticationPolicyOptions{}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = pl.Do(req)
	if err == nil {
		t.Fatalf("expected an error but did not receive one")
	}
}

func TestBearerTokenPolicy(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	cred := &dummyCredential{}
	pl := NewPipeline(srv, NewBearerTokenPolicy(cred, AuthenticationPolicyOptions{}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(req.Header, resp.Request.Header) {
		t.Fatal("unexpected modification to request headers")
	}
	if resp.Request.Header.Get(HeaderAuthorization) != fmt.Sprintf("Bearer %s", "success_token") {
		t.Fatalf("unexpected value in Authorization header: %v", resp.Request.Header.Get(HeaderAuthorization))
	}
}

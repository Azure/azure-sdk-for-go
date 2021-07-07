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

// defaultTokenProcessor is used for the common case where only one token is requested directly from the
// credential
type defaultTokenProcessor struct{}

func (t *defaultTokenProcessor) IsZeroOrExpired(eo map[string]time.Time) bool {
	// the default case will only provide one token
	// check if the token's expieration time has passed or is uninitialized
	for _, t := range eo {
		return t.IsZero() || t.Before(time.Now())
	}
	// if a nil or empty map is passed then return true to perform the initial token request
	return true
}

func (t *defaultTokenProcessor) ShouldRefresh(eo map[string]time.Time) bool {
	// the default case will check that the token's expiration time is within two minutes
	// if it is it will signal a refresh for the token.
	const window = 2 * time.Minute
	for _, t := range eo {
		return t.Add(-window).Before(time.Now())
	}
	// TODO check if this should instead return an error?
	return false
}

func (t *defaultTokenProcessor) RefreshOrGet(ctx context.Context, cred TokenCredential, opts TokenRequestOptions) (string, error) {
	tk, err := cred.GetToken(ctx, opts)
	if err != nil {
		return "", err
	}
	return bearerTokenPrefix + tk.Token, nil
}

func (t *defaultTokenProcessor) Header() string {
	return HeaderAuthorization
}

func (c *dummyCredential) AuthenticationPolicy(options AuthenticationPolicyOptions) Policy {
	return NewTokenRefreshPolicy(c, &defaultTokenProcessor{}, options)
}

func TestTokenRefreshPolicyHTTPFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	cred := &dummyCredential{}
	pl := NewPipeline(srv, NewTokenRefreshPolicy(cred, &defaultTokenProcessor{}, AuthenticationPolicyOptions{}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = pl.Do(req)
	if err == nil {
		t.Fatalf("expected an error but did not receive one")
	}
}

func TestTokenRefreshPolicy(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	cred := &dummyCredential{}
	pl := NewPipeline(srv, NewTokenRefreshPolicy(cred, &defaultTokenProcessor{}, AuthenticationPolicyOptions{}))
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

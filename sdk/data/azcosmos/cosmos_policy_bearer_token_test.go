// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestConvertBearerToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := bearerTokenVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&mockAuthPolicy{}, &cosmosBearerTokenPolicy{}, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.authHeaderContent != "type=aad&ver=1.0&sig=this is a test token" {
		t.Fatalf("Expected auth header content to be 'type=aad&ver=1.0&sig=this is a test token', got %s", verifier.authHeaderContent)
	}
}

type bearerTokenVerify struct {
	authHeaderContent string
}

func (p *bearerTokenVerify) Do(req *policy.Request) (*http.Response, error) {
	p.authHeaderContent = req.Raw().Header.Get(headerAuthorization)

	return req.Next()
}

type mockAuthPolicy struct{}

func (p *mockAuthPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Set(headerAuthorization, "Bearer this is a test token")

	return req.Next()
}

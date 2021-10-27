// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestCosmosErrorOnEmptyResponse(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(404))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", []policy.Policy{}, []policy.Policy{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)

	cError := newCosmosError(resp)
	if cError.Error() != "response contained no body" {
		t.Errorf("Expected response contained no body, but got %v", cError)
	}
}

func TestCosmosErrorOnNonJsonBody(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody([]byte("This is not JSON")),
		mock.WithStatusCode(404))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", []policy.Policy{}, []policy.Policy{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)

	cError := newCosmosError(resp)
	if cError.Error() != "This is not JSON" {
		t.Errorf("Expected This is not JSON, but got %v", cError)
	}
}

func TestCosmosErrorOnJsonBody(t *testing.T) {
	someError := &cosmosErrorResponse{
		Code: "SomeCode",
	}

	jsonString, err := json.Marshal(someError)
	if err != nil {
		t.Fatal(err)
	}

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithStatusCode(404))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", []policy.Policy{}, []policy.Policy{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)

	cError := newCosmosError(resp)
	asError := cError.(*cosmosError)
	if asError.ErrorCode() != someError.Code {
		t.Errorf("Expected %v, but got %v", someError.Code, asError.ErrorCode())
	}

	if asError.StatusCode() != 404 {
		t.Errorf("Expected 404 Not Found, but got %v", asError.StatusCode())
	}

	if asError.Error() != string(jsonString) {
		t.Errorf("Expected %v, but got %v", string(jsonString), asError.Error())
	}
}

// +build !emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestCosmosErrorOnEmptyResponse(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(404))

	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azcore.NewPipeline(srv)
	resp, _ := pl.Do(req)

	cError := newCosmosError(resp)
	if cError.Error() != "" {
		t.Errorf("Expected empty error, but got %v", cError)
	}
}

func TestCosmosErrorOnNonJsonBody(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody([]byte("This is not JSON")),
		mock.WithStatusCode(404))

	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azcore.NewPipeline(srv)
	resp, _ := pl.Do(req)

	cError := newCosmosError(resp)
	if cError.Error() != "This is not JSON" {
		t.Errorf("Expected This is not JSON, but got %v", cError)
	}
}

func TestCosmosErrorOnJsonBody(t *testing.T) {
	someError := &cosmosError{
		Code:    "SomeCode",
		Message: "SomeMessage",
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

	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azcore.NewPipeline(srv)
	resp, _ := pl.Do(req)

	cError := newCosmosError(resp)
	asError := cError.(*cosmosError)
	if asError.Code != someError.Code {
		t.Errorf("Expected %v, but got %v", someError.Code, asError.Code)
	}
	if asError.Message != someError.Message {
		t.Errorf("Expected %v, but got %v", someError.Message, asError.Message)
	}
}

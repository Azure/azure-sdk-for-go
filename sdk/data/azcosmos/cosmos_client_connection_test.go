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

func TestEnsureErrorIsGeneratedOnResponse(t *testing.T) {
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

	pl := azruntime.NewPipeline(srv)
	connection := &cosmosClientConnection{endpoint: srv.URL(), Pipeline: pl}
	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}
	_, err = connection.sendGetRequest("/", context.Background(), operationContext, &ReadContainerOptions{}, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	asError := err.(*cosmosError)
	if asError.ErrorCode() != someError.Code {
		t.Errorf("Expected %v, but got %v", someError.Code, asError.ErrorCode())
	}

	if err.Error() != asError.Error() {
		t.Errorf("Expected %v, but got %v", err.Error(), asError.Error())
	}
}

func TestEnsureErrorIsNotGeneratedOnResponse(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))

	pl := azruntime.NewPipeline(srv)
	connection := &cosmosClientConnection{endpoint: srv.URL(), Pipeline: pl}
	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}
	_, err := connection.sendGetRequest("/", context.Background(), operationContext, &ReadContainerOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRequestEnricherIsCalled(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))

	pl := azruntime.NewPipeline(srv)
	connection := &cosmosClientConnection{endpoint: srv.URL(), Pipeline: pl}
	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	addHeader := func(r *policy.Request) {
		r.Raw().Header.Add("my-header", "12345")
	}

	req, err := connection.createRequest("/", context.Background(), http.MethodGet, operationContext, &ReadContainerOptions{}, addHeader)
	if err != nil {
		t.Fatal(err)
	}

	if req.Raw().Header.Get("my-header") != "12345" {
		t.Errorf("Expected %v, but got %v", "12345", req.Raw().Header.Get("my-header"))
	}
}

func TestNoOptionsIsCalled(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))

	pl := azruntime.NewPipeline(srv)
	connection := &cosmosClientConnection{endpoint: srv.URL(), Pipeline: pl}
	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	_, err := connection.createRequest("/", context.Background(), http.MethodGet, operationContext, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

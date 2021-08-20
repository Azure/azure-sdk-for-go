// +build !emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

func TestEnsureErrorIsGeneratedOnResponse(t *testing.T) {
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

	pl := azcore.NewPipeline(srv)
	connection := &cosmosClientConnection{endpoint: srv.URL(), Pipeline: pl}
	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}
	_, err = connection.sendGetRequest("/", context.Background(), operationContext, &CosmosContainerRequestOptions{})
	if err == nil {
		t.Fatal("Expected error")
	}

	asRuntimeError := err.(*runtime.ResponseError)
	asError := asRuntimeError.Unwrap().(*cosmosError)
	if asError.Code != someError.Code {
		t.Errorf("Expected %v, but got %v", someError.Code, asError.Code)
	}
	if asError.Message != someError.Message {
		t.Errorf("Expected %v, but got %v", someError.Message, asError.Message)
	}
}

func TestEnsureErrorIsNotGeneratedOnResponse(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))

	pl := azcore.NewPipeline(srv)
	connection := &cosmosClientConnection{endpoint: srv.URL(), Pipeline: pl}
	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}
	_, err := connection.sendGetRequest("/", context.Background(), operationContext, &CosmosContainerRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}
	_, err = client.sendGetRequest("/", context.Background(), operationContext, &ReadContainerOptions{}, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	asError := err.(*azcore.ResponseError)
	if asError.ErrorCode != someError.Code {
		t.Errorf("Expected %v, but got %v", someError.Code, asError.ErrorCode)
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

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}
	_, err := client.sendGetRequest("/", context.Background(), operationContext, &ReadContainerOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRequestEnricherIsCalled(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	addHeader := func(r *policy.Request) {
		r.Raw().Header.Add("my-header", "12345")
	}

	req, err := client.createRequest("/", context.Background(), http.MethodGet, operationContext, &ReadContainerOptions{}, addHeader)
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

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	_, err := client.createRequest("/", context.Background(), http.MethodGet, operationContext, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAttachContent(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	body := map[string]string{
		"foo": "bar",
	}

	marshalled, _ := json.Marshal(body)

	// Using the interface{}
	req, err := client.createRequest("/", context.Background(), http.MethodGet, operationContext, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = client.attachContent(body, req)
	if err != nil {
		t.Fatal(err)
	}

	readBody, _ := ioutil.ReadAll(req.Body())

	if string(readBody) != string(marshalled) {
		t.Errorf("Expected %v, but got %v", string(marshalled), string(readBody))
	}

	// Using the raw []byte
	req, err = client.createRequest("/", context.Background(), http.MethodGet, operationContext, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = client.attachContent(marshalled, req)
	if err != nil {
		t.Fatal(err)
	}

	readBody, _ = ioutil.ReadAll(req.Body())

	if string(readBody) != string(marshalled) {
		t.Errorf("Expected %v, but got %v", string(marshalled), string(readBody))
	}
}

func TestCreateRequest(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	req, err := client.createRequest("/", context.Background(), http.MethodGet, operationContext, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if req.Raw().URL.String() != srv.URL()+"/" {
		t.Errorf("Expected %v, but got %v", srv.URL()+"/", req.Raw().URL.String())
	}

	if req.Raw().Method != http.MethodGet {
		t.Errorf("Expected %v, but got %v", http.MethodGet, req.Raw().Method)
	}

	if req.Raw().Header.Get(headerXmsDate) == "" {
		t.Errorf("Expected %v, but got %v", "", req.Raw().Header.Get(headerXmsDate))
	}

	if req.Raw().Header.Get(headerXmsVersion) != "2020-11-05" {
		t.Errorf("Expected %v, but got %v", "2020-11-05", req.Raw().Header.Get(headerXmsVersion))
	}

	opValue := pipelineRequestOptions{}
	if !req.OperationValue(&opValue) {
		t.Error("Expected to find operation value")
	}
}

func TestSendDelete(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))
	verifier := pipelineVerifier{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	_, err := client.sendDeleteRequest("/", context.Background(), operationContext, &DeleteDatabaseOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.method != http.MethodDelete {
		t.Errorf("Expected %v, but got %v", http.MethodDelete, verifier.method)
	}
}

func TestSendGet(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))
	verifier := pipelineVerifier{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	_, err := client.sendGetRequest("/", context.Background(), operationContext, &DeleteDatabaseOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.method != http.MethodGet {
		t.Errorf("Expected %v, but got %v", http.MethodGet, verifier.method)
	}
}

func TestSendPut(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))
	verifier := pipelineVerifier{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	body := map[string]string{
		"foo": "bar",
	}

	marshalled, _ := json.Marshal(body)

	_, err := client.sendPutRequest("/", context.Background(), body, operationContext, &DeleteDatabaseOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.method != http.MethodPut {
		t.Errorf("Expected %v, but got %v", http.MethodPut, verifier.method)
	}

	if verifier.body != string(marshalled) {
		t.Errorf("Expected %v, but got %v", string(marshalled), verifier.body)
	}
}

func TestSendPost(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))
	verifier := pipelineVerifier{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	body := map[string]string{
		"foo": "bar",
	}

	marshalled, _ := json.Marshal(body)

	_, err := client.sendPostRequest("/", context.Background(), body, operationContext, &DeleteDatabaseOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.method != http.MethodPost {
		t.Errorf("Expected %v, but got %v", http.MethodPost, verifier.method)
	}

	if verifier.body != string(marshalled) {
		t.Errorf("Expected %v, but got %v", string(marshalled), verifier.body)
	}
}

func TestSendQuery(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))
	verifier := pipelineVerifier{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	_, err := client.sendQueryRequest("/", context.Background(), "SELECT * FROM c", operationContext, &DeleteDatabaseOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.method != http.MethodPost {
		t.Errorf("Expected %v, but got %v", http.MethodPost, verifier.method)
	}

	if verifier.isQuery != true {
		t.Errorf("Expected %v, but got %v", true, verifier.isQuery)
	}

	if verifier.contentType != cosmosHeaderValuesQuery {
		t.Errorf("Expected %v, but got %v", cosmosHeaderValuesQuery, verifier.contentType)
	}

	if verifier.body != "{\"query\":\"SELECT * FROM c\"}" {
		t.Errorf("Expected %v, but got %v", "{\"query\":\"SELECT * FROM c\"}", verifier.body)
	}
}

type pipelineVerifier struct {
	method      string
	body        string
	contentType string
	isQuery     bool
}

func (p *pipelineVerifier) Do(req *policy.Request) (*http.Response, error) {
	p.method = req.Raw().Method
	if req.Body() != nil {
		readBody, _ := ioutil.ReadAll(req.Body())
		p.body = string(readBody)
	}
	p.contentType = req.Raw().Header.Get(headerContentType)
	p.isQuery = req.Raw().Method == http.MethodPost && req.Raw().Header.Get(cosmosHeaderQuery) == "True"
	return req.Next()
}

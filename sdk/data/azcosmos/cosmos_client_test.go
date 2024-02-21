// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestNewClientFromConnStrReturnErrorOnWrongDelimiter(t *testing.T) {
	invalidStr := "invalid_connection_string"
	_, err := NewClientFromConnectionString(invalidStr, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	expected := "failed parsing connection string due to it not consist of two parts separated by ';'"
	actual := err.Error()
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func TestNewClientFromConnStrReturnErrorOnWrongAccEnpoint(t *testing.T) {
	invalidStr := "invalid_str;AccountKey=dG9fYmFzZV82NA=="
	_, err := NewClientFromConnectionString(invalidStr, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	expected := "failed parsing connection string due to unmatched key value separated by '='"
	actual := err.Error()
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func TestNewClientFromConnStrReturnErrorOnWrongAccKey(t *testing.T) {
	invalidStr := "AccountEndpoint=http://127.0.0.1:80;invalid_str"
	_, err := NewClientFromConnectionString(invalidStr, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	expected := "failed parsing connection string due to unmatched key value separated by '='"
	actual := err.Error()
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func TestNewClientFromConnStrSuccess(t *testing.T) {
	connStr := "AccountEndpoint=http://127.0.0.1:80;AccountKey=dG9fYmFzZV82NA==;"
	client, err := NewClientFromConnectionString(connStr, nil)
	if err != nil {
		t.Fatal(err)
	}

	actualEnpoint := client.endpoint
	expectedEndpoint := "http://127.0.0.1:80"
	if actualEnpoint != expectedEndpoint {
		t.Errorf("Expected %v, but got %v", expectedEndpoint, actualEnpoint)
	}
}

func TestEnsureErrorIsGeneratedOnResponse(t *testing.T) {
	someError := map[string]string{"Code": "SomeCode"}

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
	if asError.ErrorCode != "404 Not Found" {
		t.Errorf("Expected %v, but got %v", "404 Not Found", asError.ErrorCode)
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

	readBody, _ := io.ReadAll(req.Body())

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

	readBody, _ = io.ReadAll(req.Body())

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

	if req.Raw().Header.Get(headerXmsVersion) != apiVersion {
		t.Errorf("Expected %v, but got %v", apiVersion, req.Raw().Header.Get(headerXmsVersion))
	}

	if req.Raw().Header.Get(cosmosHeaderSDKSupportedCapabilities) != supportedCapabilitiesHeaderValue {
		t.Errorf("Expected %v, but got %v", supportedCapabilitiesHeaderValue, req.Raw().Header.Get(cosmosHeaderSDKSupportedCapabilities))
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

	if verifier.requests[0].method != http.MethodDelete {
		t.Errorf("Expected %v, but got %v", http.MethodDelete, verifier.requests[0].method)
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

	if verifier.requests[0].method != http.MethodGet {
		t.Errorf("Expected %v, but got %v", http.MethodGet, verifier.requests[0].method)
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

	if verifier.requests[0].method != http.MethodPut {
		t.Errorf("Expected %v, but got %v", http.MethodPut, verifier.requests[0].method)
	}

	if verifier.requests[0].body != string(marshalled) {
		t.Errorf("Expected %v, but got %v", string(marshalled), verifier.requests[0].body)
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

	if verifier.requests[0].method != http.MethodPost {
		t.Errorf("Expected %v, but got %v", http.MethodPost, verifier.requests[0].method)
	}

	if verifier.requests[0].body != string(marshalled) {
		t.Errorf("Expected %v, but got %v", string(marshalled), verifier.requests[0].body)
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

	_, err := client.sendQueryRequest("/", context.Background(), "SELECT * FROM c", []QueryParameter{}, operationContext, &DeleteDatabaseOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.requests[0].method != http.MethodPost {
		t.Errorf("Expected %v, but got %v", http.MethodPost, verifier.requests[0].method)
	}

	if verifier.requests[0].isQuery != true {
		t.Errorf("Expected %v, but got %v", true, verifier.requests[0].isQuery)
	}

	if verifier.requests[0].contentType != cosmosHeaderValuesQuery {
		t.Errorf("Expected %v, but got %v", cosmosHeaderValuesQuery, verifier.requests[0].contentType)
	}

	if verifier.requests[0].body != "{\"query\":\"SELECT * FROM c\"}" {
		t.Errorf("Expected %v, but got %v", "{\"query\":\"SELECT * FROM c\"}", verifier.requests[0].body)
	}
}

func TestSendQueryWithParameters(t *testing.T) {
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

	parameters := []QueryParameter{
		{"@id", "1"},
		{"@status", "enabled"},
	}

	_, err := client.sendQueryRequest("/", context.Background(), "SELECT * FROM c WHERE c.id = @id and c.status = @status", parameters, operationContext, &DeleteDatabaseOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.requests[0].method != http.MethodPost {
		t.Errorf("Expected %v, but got %v", http.MethodPost, verifier.requests[0].method)
	}

	if verifier.requests[0].isQuery != true {
		t.Errorf("Expected %v, but got %v", true, verifier.requests[0].isQuery)
	}

	if verifier.requests[0].contentType != cosmosHeaderValuesQuery {
		t.Errorf("Expected %v, but got %v", cosmosHeaderValuesQuery, verifier.requests[0].contentType)
	}

	expectedSerializedQuery := "{\"query\":\"SELECT * FROM c WHERE c.id = @id and c.status = @status\",\"parameters\":[{\"name\":\"@id\",\"value\":\"1\"},{\"name\":\"@status\",\"value\":\"enabled\"}]}"

	if verifier.requests[0].body != expectedSerializedQuery {
		t.Errorf("Expected %v, but got %v", expectedSerializedQuery, verifier.requests[0].body)
	}
}

func TestSendBatch(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(200))
	verifier := pipelineVerifier{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDocument,
		resourceAddress: "",
	}

	batch := TransactionalBatch{}
	batch.partitionKey = NewPartitionKeyString("foo")

	body := map[string]string{
		"foo": "bar",
	}

	itemMarshall, _ := json.Marshal(body)

	batch.CreateItem(itemMarshall, nil)
	batch.ReadItem("someId", nil)

	marshalled, err := json.Marshal(batch.operations)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.sendBatchRequest(context.Background(), "/", batch.operations, operationContext, &TransactionalBatchOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.requests[0].method != http.MethodPost {
		t.Errorf("Expected %v, but got %v", http.MethodPost, verifier.requests[0].method)
	}

	if verifier.requests[0].body != string(marshalled) {
		t.Errorf("Expected %v, but got %v", string(marshalled), verifier.requests[0].body)
	}
}

func TestSendPatch(t *testing.T) {
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

	body := map[string]any{
		"condition": "from c where c.Address.ZipCode ='98101' ",
		"operations": []struct {
			Op    string `json:"op"`
			Path  string `json:"path"`
			Value any    `json:"value"`
		}{
			{
				Op:    "replace",
				Path:  "/Address/ZipCode",
				Value: 98107,
			},
		},
	}

	marshalled, _ := json.Marshal(body)

	_, err := client.sendPatchRequest("/", context.Background(), body, operationContext, &ItemOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if verifier.requests[0].method != http.MethodPatch {
		t.Errorf("Expected %v, but got %v", http.MethodPost, verifier.requests[0].method)
	}

	if verifier.requests[0].body != string(marshalled) {
		t.Errorf("Expected %v, but got %v", string(marshalled), verifier.requests[0].body)
	}
}

func TestCreateScopeFromEndpoint(t *testing.T) {
	url := "https://foo.documents.azure.com:443/"
	scope, err := createScopeFromEndpoint(url)
	if err != nil {
		t.Fatal(err)
	}

	if scope[0] != "https://foo.documents.azure.com/.default" {
		t.Errorf("Expected %v, but got %v", "https://foo.documents.azure.com/.default", scope[0])
	}

	if len(scope) != 1 {
		t.Errorf("Expected %v, but got %v", 1, len(scope))
	}
}

func TestQueryDatabases(t *testing.T) {
	jsonStringpage1 := []byte(`{"Databases":[{"id":"doc1"},{"id":"doc2"}]}`)
	jsonStringpage2 := []byte(`{"Databases":[{"id":"doc3"},{"id":"doc4"},{"id":"doc5"}]}`)

	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(
		mock.WithBody(jsonStringpage1),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithHeader(cosmosHeaderContinuationToken, "someContinuationToken"),
		mock.WithStatusCode(200))
	srv.AppendResponse(
		mock.WithBody(jsonStringpage2),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}

	receivedIds := []string{}
	queryPager := client.NewQueryDatabasesPager("select * from c", nil)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query items: %v", err)
		}

		for _, dbs := range queryResponse.Databases {
			receivedIds = append(receivedIds, dbs.ID)
		}

		if queryPager.More() && queryResponse.ContinuationToken != "someContinuationToken" {
			t.Errorf("Expected ContinuationToken to be %s, but got %s", "someContinuationToken", queryResponse.ContinuationToken)
		}

		if queryResponse.ActivityID == "" {
			t.Fatal("Activity id was not returned")
		}

		if queryResponse.RequestCharge == 0 {
			t.Fatal("Request charge was not returned")
		}
	}

	for i := 0; i < 5; i++ {
		if receivedIds[i] != "doc"+strconv.Itoa(i+1) {
			t.Fatalf("Expected id %d, got %s", i, receivedIds[i])
		}
	}

	if len(verifier.requests) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(verifier.requests))
	}

	for index, request := range verifier.requests {
		if request.method != http.MethodPost {
			t.Errorf("Expected method to be %s, but got %s", http.MethodPost, request.method)
		}

		if request.url.RequestURI() != "/dbs" {
			t.Errorf("Expected url to be %s, but got %s", "/dbs", request.url.RequestURI())
		}

		if !request.isQuery {
			t.Errorf("Expected request to be a query, but it was not")
		}

		if request.body != "{\"query\":\"select * from c\"}" {
			t.Errorf("Expected %v, but got %v", "{\"query\":\"select * from c\"}", request.body)
		}

		if request.contentType != cosmosHeaderValuesQuery {
			t.Errorf("Expected %v, but got %v", cosmosHeaderValuesQuery, request.contentType)
		}

		if index == 0 && request.headers.Get(cosmosHeaderContinuationToken) != "" {
			t.Errorf("Expected ContinuationToken to be %s, but got %s", "", request.headers.Get(cosmosHeaderContinuationToken))
		}

		if index == 1 && request.headers.Get(cosmosHeaderContinuationToken) != "someContinuationToken" {
			t.Errorf("Expected ContinuationToken to be %s, but got %s", "someContinuationToken", request.headers.Get(cosmosHeaderContinuationToken))
		}
	}
}

type pipelineVerifier struct {
	requests []pipelineVerifierRequest
}

type pipelineVerifierRequest struct {
	method      string
	body        string
	contentType string
	isQuery     bool
	url         *url.URL
	headers     http.Header
}

func (p *pipelineVerifier) Do(req *policy.Request) (*http.Response, error) {
	pr := pipelineVerifierRequest{}
	pr.method = req.Raw().Method
	pr.url = req.Raw().URL
	if req.Body() != nil {
		readBody, _ := io.ReadAll(req.Body())
		pr.body = string(readBody)
	}
	pr.contentType = req.Raw().Header.Get(headerContentType)
	pr.headers = req.Raw().Header
	pr.isQuery = req.Raw().Method == http.MethodPost && req.Raw().Header.Get(cosmosHeaderQuery) == "True"
	p.requests = append(p.requests, pr)
	return req.Next()
}

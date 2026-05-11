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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

func TestAddContentHeaderDefaultOnWriteOperation(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
		resourceType:     resourceTypeDocument,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !verifier.isEnableContentResponseOnWriteHeaderSet {
		t.Fatalf("expected content response header to be set")
	}
}

func TestAddContentHeaderDefaultOnNonDocumentWriteOperation(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
		resourceType:     resourceTypeCollection,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.isEnableContentResponseOnWriteHeaderSet {
		t.Fatalf("expected content response header not to be set")
	}
}

func TestAddContentHeaderDefaultOnReadOperation(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: false,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.isEnableContentResponseOnWriteHeaderSet {
		t.Fatalf("expected content response header to not be set")
	}
}

func TestAddContentHeaderOnWriteOperation(t *testing.T) {
	headerPolicy := &headerPolicies{
		enableContentResponseOnWrite: true,
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
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

	if verifier.isEnableContentResponseOnWriteHeaderSet {
		t.Fatalf("expected content response header to not be set")
	}
}

func TestAddContentHeaderOnWriteOperationWithOverride(t *testing.T) {
	headerPolicy := &headerPolicies{
		enableContentResponseOnWrite: true,
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	contentOverride := false
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
		resourceType:     resourceTypeDocument,
		headerOptionsOverride: &headerOptionsOverride{
			enableContentResponseOnWrite: &contentOverride,
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !verifier.isEnableContentResponseOnWriteHeaderSet {
		t.Fatalf("expected content response header to be set")
	}
}

func TestAddContentHeaderDefaultOnWriteOperationWithOverride(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	contentOverride := true
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
		headerOptionsOverride: &headerOptionsOverride{
			enableContentResponseOnWrite: &contentOverride,
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.isEnableContentResponseOnWriteHeaderSet {
		t.Fatalf("expected content response header to not be set")
	}
}

func TestAddPartitionKeyHeader(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())

	partitionKey := NewPartitionKeyString("some string")
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
		headerOptionsOverride: &headerOptionsOverride{
			partitionKey: &partitionKey,
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.isPartitionKeyHeaderSet != "[\"some string\"]" {
		t.Fatalf("expected pk header to be set")
	}
}

func TestAddCorrelatedActivityIdHeader(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())

	correlatedActivityId, _ := uuid.New()
	req.SetOperationValue(pipelineRequestOptions{
		headerOptionsOverride: &headerOptionsOverride{
			correlatedActivityId: &correlatedActivityId,
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.isCorrelatedActivityIdSet != correlatedActivityId.String() {
		t.Fatalf("expected correlatedActivityId header to be set")
	}
}

type headerPoliciesVerify struct {
	isEnableContentResponseOnWriteHeaderSet bool
	isPartitionKeyHeaderSet                 string
	isCorrelatedActivityIdSet               string
	priorityLevelHeaderValue                string
	throughputBucketHeaderValue             string
}

func (p *headerPoliciesVerify) Do(req *policy.Request) (*http.Response, error) {
	p.isEnableContentResponseOnWriteHeaderSet = req.Raw().Header.Get(cosmosHeaderPrefer) != ""
	p.isPartitionKeyHeaderSet = req.Raw().Header.Get(cosmosHeaderPartitionKey)
	p.isCorrelatedActivityIdSet = req.Raw().Header.Get(cosmosHeaderCorrelatedActivityId)
	p.priorityLevelHeaderValue = req.Raw().Header.Get(cosmosHeaderPriorityLevel)
	p.throughputBucketHeaderValue = req.Raw().Header.Get(cosmosHeaderThroughputBucket)

	return req.Next()
}

func TestPriorityLevelHeaderFromClientDefault(t *testing.T) {
	priority := PriorityLevelLow
	headerPolicy := &headerPolicies{
		priorityLevel: &priority,
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.SetOperationValue(pipelineRequestOptions{})

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.priorityLevelHeaderValue != "Low" {
		t.Fatalf("expected priority level header to be Low but got %v", verifier.priorityLevelHeaderValue)
	}
}

func TestPriorityLevelHeaderRequestOverridesClient(t *testing.T) {
	clientPriority := PriorityLevelHigh
	headerPolicy := &headerPolicies{
		priorityLevel: &clientPriority,
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requestPriority := PriorityLevelLow
	req.SetOperationValue(pipelineRequestOptions{
		headerOptionsOverride: &headerOptionsOverride{
			priorityLevel: &requestPriority,
		},
	})

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.priorityLevelHeaderValue != "Low" {
		t.Fatalf("expected priority level header to be Low (request override) but got %v", verifier.priorityLevelHeaderValue)
	}
}

func TestPriorityLevelHeaderNotSetWhenNil(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.SetOperationValue(pipelineRequestOptions{})

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.priorityLevelHeaderValue != "" {
		t.Fatalf("expected no priority level header but got %v", verifier.priorityLevelHeaderValue)
	}
}

func TestThroughputBucketHeaderFromClientDefault(t *testing.T) {
	bucket := int32(3)
	headerPolicy := &headerPolicies{
		throughputBucket: &bucket,
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.SetOperationValue(pipelineRequestOptions{})

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.throughputBucketHeaderValue != "3" {
		t.Fatalf("expected throughput bucket header to be 3 but got %v", verifier.throughputBucketHeaderValue)
	}
}

func TestThroughputBucketHeaderRequestOverridesClient(t *testing.T) {
	clientBucket := int32(1)
	headerPolicy := &headerPolicies{
		throughputBucket: &clientBucket,
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requestBucket := int32(5)
	req.SetOperationValue(pipelineRequestOptions{
		headerOptionsOverride: &headerOptionsOverride{
			throughputBucket: &requestBucket,
		},
	})

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.throughputBucketHeaderValue != "5" {
		t.Fatalf("expected throughput bucket header to be 5 (request override) but got %v", verifier.throughputBucketHeaderValue)
	}
}

func TestThroughputBucketHeaderNotSetWhenNil(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.SetOperationValue(pipelineRequestOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier.throughputBucketHeaderValue != "" {
		t.Fatalf("expected no throughput bucket header but got %v", verifier.throughputBucketHeaderValue)
	}
}

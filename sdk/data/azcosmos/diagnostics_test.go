// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestDiagnosticsSummaryIncludesRetriedGatewayCalls(t *testing.T) {
	rootTrace := newRootTrace("test_operation")
	requestTrace := rootTrace.StartChild(traceDatumKeyTransportRequest)
	clientStats := newClientSideRequestStatisticsTraceDatum(time.Now().UTC(), requestTrace)
	requestTrace.AddDatum(traceDatumKeyClientSideRequestStats, clientStats)
	requestTrace.AddDatum(traceDatumKeyPointOperationStatistics, pointOperationStatisticsTraceDatum{
		ActivityID:      "activity-2",
		ResponseTimeUTC: time.Now().UTC(),
		StatusCode:      http.StatusOK,
		RequestCharge:   2.5,
		RequestURI:      "https://example.com/dbs/test",
		BELatencyInMs:   "12",
	})

	req503, err := http.NewRequest(http.MethodGet, "https://example.com/dbs/test", nil)
	require.NoError(t, err)
	resp503 := &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Header: http.Header{
			cosmosHeaderActivityId:    []string{"activity-1"},
			cosmosHeaderRequestCharge: []string{"1.5"},
		},
		Request: req503,
	}
	clientStats.recordHTTPResponse(time.Now().Add(-20*time.Millisecond), resp503, resourceTypeDatabase, "local")

	req200, err := http.NewRequest(http.MethodGet, "https://example.com/dbs/test", nil)
	require.NoError(t, err)
	resp200 := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			cosmosHeaderActivityId:    []string{"activity-2"},
			cosmosHeaderRequestCharge: []string{"2.5"},
			cosmosHeaderSessionToken:  []string{"response-session"},
		},
		Request: req200,
	}
	clientStats.recordHTTPResponse(time.Now().Add(-10*time.Millisecond), resp200, resourceTypeDatabase, "local")

	requestTrace.End()
	rootTrace.End()

	diagnosticsPayload := newDiagnostics(requestTrace).String()
	require.NotEmpty(t, diagnosticsPayload)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(diagnosticsPayload), &parsed))

	summary := parsed["Summary"].(map[string]any)
	gatewayCalls := summary["GatewayCalls"].(map[string]any)
	require.Equal(t, float64(1), gatewayCalls["(200, 0)"])
	require.Equal(t, float64(1), gatewayCalls["(503, 0)"])
	require.Equal(t, float64(1), summary["RegionsContacted"])

	children := parsed["children"].([]any)
	require.Len(t, children, 1)

	child := children[0].(map[string]any)
	data := child["data"].(map[string]any)
	require.Contains(t, data, traceDatumKeyClientSideRequestStats)
	require.Contains(t, data, traceDatumKeyPointOperationStatistics)
}

func TestDiagnosticsFromErrorReturnsResponseDiagnostics(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()

	srv.SetResponse(
		mock.WithStatusCode(http.StatusNotFound),
		mock.WithHeader(cosmosHeaderActivityId, "activity-404"),
		mock.WithHeader(cosmosHeaderRequestCharge, "3.25"),
	)

	client := newTestDiagnosticsClient(t, srv, false)
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	_, err := client.sendGetRequest("/", context.Background(), operationContext, &ReadContainerOptions{}, nil)
	require.Error(t, err)

	var responseErr *azcore.ResponseError
	require.True(t, errors.As(err, &responseErr))

	diagnostics, ok := DiagnosticsFromError(err)
	require.True(t, ok)
	require.NotEmpty(t, diagnostics.String())

	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(diagnostics.String()), &parsed))

	data := parsed["data"].(map[string]any)
	pointStats := data[traceDatumKeyPointOperationStatistics].(map[string]any)
	require.Equal(t, float64(http.StatusNotFound), pointStats["StatusCode"])
}

func TestQueryResponseDiagnosticsIncludeQueryMetrics(t *testing.T) {
	queryResponseRaw := map[string][]map[string]string{
		"Documents": {
			{"id": "id1"},
		},
	}

	jsonString, err := json.Marshal(queryResponseRaw)
	require.NoError(t, err)

	srv, close := mock.NewTLSServer()
	defer close()

	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithStatusCode(http.StatusOK),
		mock.WithHeader(cosmosHeaderActivityId, "query-activity"),
		mock.WithHeader(cosmosHeaderRequestCharge, "5.5"),
		mock.WithHeader(cosmosHeaderQueryMetrics, "retrievedDocumentCount=1"),
	)

	client := newTestDiagnosticsClient(t, srv, false)
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDocument,
		resourceAddress: "dbs/test/colls/test",
	}

	resp, err := client.sendQueryRequest("/", context.Background(), "SELECT * FROM c", nil, operationContext, &QueryOptions{}, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	parsedResponse, err := newQueryResponse(resp)
	require.NoError(t, err)
	require.NotNil(t, parsedResponse.QueryMetrics)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(parsedResponse.Diagnostics.String()), &parsed))

	data := parsed["data"].(map[string]any)
	require.Equal(t, "retrievedDocumentCount=1", data[traceDatumKeyQueryMetrics])
}

func newTestDiagnosticsClient(t *testing.T, transport policy.Transporter, withRetryPolicy bool) *Client {
	t.Helper()

	serverURL := transport.(*mock.Server).URL()
	endpointURL, err := url.Parse(serverURL)
	require.NoError(t, err)

	gem := &globalEndpointManager{
		preferredLocations: []string{"local"},
		locationCache:      newLocationCache([]string{"local"}, *endpointURL, true),
	}

	pipelineOptions := azruntime.PipelineOptions{}
	if withRetryPolicy {
		pipelineOptions.PerRetry = []policy.Policy{
			&clientRetryPolicy{gem: gem},
		}
	}

	internalClient, err := azcore.NewClient("azcosmostest", "v1.0.0", pipelineOptions, &policy.ClientOptions{Transport: transport})
	require.NoError(t, err)

	return &Client{
		endpoint:    serverURL,
		endpointUrl: endpointURL,
		internal:    internalClient,
		gem:         gem,
	}
}

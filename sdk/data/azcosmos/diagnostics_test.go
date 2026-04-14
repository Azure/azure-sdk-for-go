// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
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

func TestDiagnosticsStringRendersSampleGatewayJSON(t *testing.T) {
	rootTrace := newFixedTrace("sample_operation", fixedDiagnosticsTime(0), fixedDiagnosticsTime(100*time.Millisecond), nil)
	requestTrace := newFixedTrace(traceDatumKeyTransportRequest, fixedDiagnosticsTime(5*time.Millisecond), fixedDiagnosticsTime(95*time.Millisecond), rootTrace)

	requestTrace.AddDatum(traceDatumKeyClientSideRequestStats, &clientSideRequestStatisticsTraceDatum{
		trace:               requestTrace,
		requestStartTimeUTC: fixedDiagnosticsTime(5 * time.Millisecond),
		requestEndTimeUTC:   timePtr(fixedDiagnosticsTime(24 * time.Millisecond)),
		regionsContacted: []contactedRegion{
			{name: "local", uri: "https://example.com/dbs/test"},
		},
		regionsContactedByURI: map[string]struct{}{
			"https://example.com/dbs/test": {},
		},
		httpResponseStatistics: []httpResponseStatistics{
			{
				startTimeUTC:   fixedDiagnosticsTime(10 * time.Millisecond),
				duration:       3 * time.Millisecond,
				requestURI:     "https://example.com/dbs/test",
				resourceType:   resourceTypeDatabase,
				httpMethod:     http.MethodGet,
				activityID:     "activity-1",
				statusCode:     http.StatusServiceUnavailable,
				statusCodeText: "ServiceUnavailable",
				reasonPhrase:   "Service Unavailable",
			},
			{
				startTimeUTC:   fixedDiagnosticsTime(20 * time.Millisecond),
				duration:       4 * time.Millisecond,
				requestURI:     "https://example.com/dbs/test",
				resourceType:   resourceTypeDatabase,
				httpMethod:     http.MethodGet,
				activityID:     "activity-2",
				statusCode:     http.StatusOK,
				statusCodeText: "OK",
				reasonPhrase:   "OK",
			},
		},
	})
	requestTrace.AddDatum(traceDatumKeyPointOperationStatistics, pointOperationStatisticsTraceDatum{
		ActivityID:      "activity-2",
		ResponseTimeUTC: fixedDiagnosticsTime(25 * time.Millisecond),
		StatusCode:      http.StatusOK,
		RequestCharge:   2.5,
		RequestURI:      "https://example.com/dbs/test",
		BELatencyInMs:   "12",
	})
	rootTrace.summary.incrementFailedCount()

	diagnostics := newDiagnostics(rootTrace)
	require.Equal(t, 1, diagnostics.FailedRequestCount())
	requireRenderedDiagnosticsJSON(t, `
{
  "Summary": {
    "RegionsContacted": 1,
    "GatewayCalls": {
      "(200, 0)": 1,
      "(503, 0)": 1
    }
  },
  "name": "sample_operation",
  "start datetime": "2024-01-02T03:04:05.000Z",
  "duration in milliseconds": 100,
  "children": [
    {
      "name": "Microsoft.Azure.Documents.ServerStoreModel Transport Request",
      "duration in milliseconds": 90,
      "data": {
        "Client Side Request Stats": {
          "Id": "AggregatedClientSideRequestStatistics",
          "ContactedReplicas": [],
          "RegionsContacted": [
            "https://example.com/dbs/test"
          ],
          "FailedReplicas": [],
          "AddressResolutionStatistics": [],
          "StoreResponseStatistics": [],
          "HttpResponseStats": [
            {
              "StartTimeUTC": "2024-01-02T03:04:05.0100000Z",
              "DurationInMs": 3,
              "RequestUri": "https://example.com/dbs/test",
              "ResourceType": "Database",
              "HttpMethod": "GET",
              "ActivityId": "activity-1",
              "StatusCode": "ServiceUnavailable",
              "ReasonPhrase": "Service Unavailable"
            },
            {
              "StartTimeUTC": "2024-01-02T03:04:05.0200000Z",
              "DurationInMs": 4,
              "RequestUri": "https://example.com/dbs/test",
              "ResourceType": "Database",
              "HttpMethod": "GET",
              "ActivityId": "activity-2",
              "StatusCode": "OK"
            }
          ]
        },
        "PointOperationStatisticsTraceDatum": {
          "Id": "PointOperationStatistics",
          "ActivityId": "activity-2",
          "ResponseTimeUtc": "2024-01-02T03:04:05.0250000Z",
          "StatusCode": 200,
          "SubStatusCode": 0,
          "RequestCharge": 2.5,
          "RequestUri": "https://example.com/dbs/test",
          "ErrorMessage": null,
          "RequestSessionToken": null,
          "ResponseSessionToken": null,
          "BELatencyInMs": "12"
        }
      }
    }
  ]
}
`, diagnostics.String())
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

func TestDiagnosticsFromErrorRendersSampleJSON(t *testing.T) {
	requestTrace := newFixedTrace(traceDatumKeyTransportRequest, fixedDiagnosticsTime(0), fixedDiagnosticsTime(40*time.Millisecond), nil)
	requestTrace.AddDatum(traceDatumKeyPointOperationStatistics, pointOperationStatisticsTraceDatum{
		ActivityID:      "missing-activity",
		ResponseTimeUTC: fixedDiagnosticsTime(30 * time.Millisecond),
		StatusCode:      http.StatusNotFound,
		RequestCharge:   3.25,
		RequestURI:      "https://example.com/dbs/test",
		ErrorMessage:    "resource missing",
	})

	req, err := http.NewRequest(http.MethodGet, "https://example.com/dbs/test", nil)
	require.NoError(t, err)
	req = req.WithContext(withRequestDiagnosticsState(context.Background(), &requestDiagnosticsState{
		requestTrace: requestTrace,
	}))

	diagnostics, ok := DiagnosticsFromError(&azcore.ResponseError{
		StatusCode:  http.StatusNotFound,
		RawResponse: &http.Response{StatusCode: http.StatusNotFound, Request: req},
	})
	require.True(t, ok)
	requireRenderedDiagnosticsJSON(t, `
{
  "Summary": {},
  "name": "Microsoft.Azure.Documents.ServerStoreModel Transport Request",
  "start datetime": "2024-01-02T03:04:05.000Z",
  "duration in milliseconds": 40,
  "data": {
    "PointOperationStatisticsTraceDatum": {
      "Id": "PointOperationStatistics",
      "ActivityId": "missing-activity",
      "ResponseTimeUtc": "2024-01-02T03:04:05.0300000Z",
      "StatusCode": 404,
      "SubStatusCode": 0,
      "RequestCharge": 3.25,
      "RequestUri": "https://example.com/dbs/test",
      "ErrorMessage": "resource missing",
      "RequestSessionToken": null,
      "ResponseSessionToken": null,
      "BELatencyInMs": null
    }
  }
}
`, diagnostics.String())
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

func TestDiagnosticsStringRendersSampleQueryMetricsJSON(t *testing.T) {
	rootTrace := newFixedTrace("sample_query", fixedDiagnosticsTime(0), fixedDiagnosticsTime(50*time.Millisecond), nil)
	rootTrace.AddDatum(traceDatumKeyPointOperationStatistics, pointOperationStatisticsTraceDatum{
		ActivityID:           "query-activity",
		ResponseTimeUTC:      fixedDiagnosticsTime(20 * time.Millisecond),
		StatusCode:           http.StatusOK,
		RequestCharge:        5.5,
		RequestURI:           "https://example.com/dbs/test/colls/test",
		RequestSessionToken:  "request-session",
		ResponseSessionToken: "response-session",
		BELatencyInMs:        "7",
	})
	rootTrace.AddDatum(traceDatumKeyQueryMetrics, "retrievedDocumentCount=1")

	diagnostics := newDiagnostics(rootTrace)
	requireRenderedDiagnosticsJSON(t, `
{
  "Summary": {},
  "name": "sample_query",
  "start datetime": "2024-01-02T03:04:05.000Z",
  "duration in milliseconds": 50,
  "data": {
    "PointOperationStatisticsTraceDatum": {
      "Id": "PointOperationStatistics",
      "ActivityId": "query-activity",
      "ResponseTimeUtc": "2024-01-02T03:04:05.0200000Z",
      "StatusCode": 200,
      "SubStatusCode": 0,
      "RequestCharge": 5.5,
      "RequestUri": "https://example.com/dbs/test/colls/test",
      "ErrorMessage": null,
      "RequestSessionToken": "request-session",
      "ResponseSessionToken": "response-session",
      "BELatencyInMs": "7"
    },
    "Query Metrics": "retrievedDocumentCount=1"
  }
}
`, diagnostics.String())
}

func TestDiagnosticsFromErrorReturnsQueryParseDiagnostics(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()

	srv.SetResponse(
		mock.WithBody([]byte(`{"Documents":[`)),
		mock.WithStatusCode(http.StatusOK),
		mock.WithHeader(cosmosHeaderActivityId, "query-activity"),
		mock.WithHeader(cosmosHeaderRequestCharge, "5.5"),
	)

	client := newTestDiagnosticsClient(t, srv, false)
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDocument,
		resourceAddress: "dbs/test/colls/test",
	}

	resp, err := client.sendQueryRequest("/", context.Background(), "SELECT * FROM c", nil, operationContext, &QueryOptions{}, nil)
	require.NoError(t, err)

	_, err = newQueryResponse(resp)
	require.Error(t, err)

	diagnostics, ok := DiagnosticsFromError(err)
	require.True(t, ok)
	require.NotEmpty(t, diagnostics.String())

	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(diagnostics.String()), &parsed))

	data := parsed["data"].(map[string]any)
	pointStats := data[traceDatumKeyPointOperationStatistics].(map[string]any)
	require.Equal(t, float64(http.StatusOK), pointStats["StatusCode"])
	require.Equal(t, float64(5.5), pointStats["RequestCharge"])
}

func TestReadManyItemsDiagnosticsAreStableAfterReturn(t *testing.T) {
	containerBody, err := json.Marshal(ContainerProperties{
		ID: "testcontainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	})
	require.NoError(t, err)

	rangeBody, err := json.Marshal(partitionKeyRangeResponse{
		PartitionKeyRanges: []partitionKeyRange{
			{
				ID:           "0",
				MinInclusive: "",
				MaxExclusive: "",
			},
		},
		Count: 1,
	})
	require.NoError(t, err)

	queryBody, err := json.Marshal(map[string][]map[string]string{
		"Documents": {
			{"id": "item1", "pk": "pk1"},
		},
	})
	require.NoError(t, err)

	srv, close := mock.NewTLSServer()
	defer close()

	srv.AppendResponse(
		mock.WithBody(containerBody),
		mock.WithStatusCode(http.StatusOK),
		mock.WithHeader(cosmosHeaderActivityId, "container-read"),
	)
	srv.AppendResponse(
		mock.WithBody(rangeBody),
		mock.WithStatusCode(http.StatusOK),
		mock.WithHeader(cosmosHeaderActivityId, "range-read"),
	)
	srv.AppendResponse(
		mock.WithBody(queryBody),
		mock.WithStatusCode(http.StatusOK),
		mock.WithHeader(cosmosHeaderActivityId, "query-read"),
		mock.WithHeader(cosmosHeaderRequestCharge, "1.5"),
	)

	container := newTestDiagnosticsContainer(t, srv, false)

	resp, err := container.ReadManyItems(context.Background(), []ItemIdentity{
		{
			ID:           "item1",
			PartitionKey: NewPartitionKeyString("pk1"),
		},
	}, nil)
	require.NoError(t, err)
	require.Len(t, resp.Items, 1)

	firstDuration := resp.Diagnostics.ClientElapsedTime()
	firstPayload := resp.Diagnostics.String()
	require.Greater(t, firstDuration, time.Duration(0))
	require.NotEmpty(t, firstPayload)

	time.Sleep(20 * time.Millisecond)

	secondDuration := resp.Diagnostics.ClientElapsedTime()
	secondPayload := resp.Diagnostics.String()
	require.Equal(t, firstDuration, secondDuration)
	require.Equal(t, firstPayload, secondPayload)
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

func newTestDiagnosticsContainer(t *testing.T, transport policy.Transporter, withRetryPolicy bool) *ContainerClient {
	t.Helper()

	client := newTestDiagnosticsClient(t, transport, withRetryPolicy)

	database, err := newDatabase("testdb", client)
	require.NoError(t, err)

	container, err := database.NewContainer("testcontainer")
	require.NoError(t, err)

	return container
}

func fixedDiagnosticsTime(offset time.Duration) time.Time {
	return time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC).Add(offset)
}

func newFixedTrace(name string, start, end time.Time, parent *trace) *trace {
	endUTC := end.UTC()
	fixedTrace := &trace{
		name:      name,
		startTime: start.UTC(),
		endTime:   &endUTC,
		parent:    parent,
		children:  []*trace{},
	}

	if parent != nil {
		fixedTrace.summary = parent.summary
		parent.AddChild(fixedTrace)
	} else {
		fixedTrace.summary = &traceSummary{}
	}

	return fixedTrace
}

func timePtr(value time.Time) *time.Time {
	value = value.UTC()
	return &value
}

func requireRenderedDiagnosticsJSON(t *testing.T, expected string, actual string) {
	t.Helper()

	var compact bytes.Buffer
	require.NoError(t, json.Compact(&compact, []byte(expected)))
	require.Equal(t, compact.String(), actual)
}

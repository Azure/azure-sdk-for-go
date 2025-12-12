// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/stretchr/testify/require"
)

func TestMetricsClient(t *testing.T) {
	client, err := azquery.NewMetricsClient(credential, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	c := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "https://...",
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Audience: "",
				Endpoint: "",
			},
		},
	}
	opts := azcore.ClientOptions{Cloud: c}
	cloudClient, err := azquery.NewMetricsClient(credential, &azquery.MetricsClientOptions{ClientOptions: opts})
	require.Error(t, err)
	require.Equal(t, err.Error(), "provided Cloud field is missing Azure Monitor Metrics configuration")
	require.Nil(t, cloudClient)
}

func TestQueryResource_BasicQuerySuccess(t *testing.T) {
	client := startMetricsTest(t)
	timespan := azquery.TimeInterval("PT12H")
	res, err := client.QueryResource(context.Background(), resourceURI,
		&azquery.MetricsClientQueryResourceOptions{
			Timespan:        to.Ptr(timespan),
			Interval:        to.Ptr("PT1M"),
			Aggregation:     to.SliceOfPtrs(azquery.AggregationTypeAverage, azquery.AggregationTypeCount),
			OrderBy:         to.Ptr("Average asc"),
			MetricNamespace: to.Ptr("Microsoft.AppConfiguration/configurationStores"),
		})
	require.NoError(t, err)
	require.NotNil(t, res.Timespan)
	require.Equal(t, *res.Value[0].ErrorCode, "Success")
	require.Equal(t, *res.Namespace, "Microsoft.AppConfiguration/configurationStores")

	testSerde(t, &res)
	testSerde(t, res.Value[0])
	testSerde(t, res.Value[0].Name)
	testSerde(t, res.Value[0].TimeSeries[0])
}

func TestQueryResource_BasicQueryFailure(t *testing.T) {
	client := startMetricsTest(t)
	invalidResourceURI := "123"
	var httpErr *azcore.ResponseError

	res, err := client.QueryResource(context.Background(), invalidResourceURI, nil)

	require.Error(t, err)
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "MissingSubscription")
	require.Equal(t, httpErr.StatusCode, 404)
	require.Nil(t, res.Timespan)
	require.Nil(t, res.Value)
	require.Nil(t, res.Cost)
	require.Nil(t, res.Interval)
	require.Nil(t, res.Namespace)
	require.Nil(t, res.ResourceRegion)

	testSerde(t, &res)
}

func TestNewListDefinitionsPager_Success(t *testing.T) {
	client := startMetricsTest(t)

	pager := client.NewListDefinitionsPager(resourceURI, nil)

	// test if first page is valid
	if pager.More() {
		res, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res.Value)
		testSerde(t, &res.MetricDefinitionCollection)
	} else {
		t.Fatal("no response")
	}

}

func TestNewListDefinitionsPager_Failure(t *testing.T) {
	client := startMetricsTest(t)

	pager := client.NewListDefinitionsPager(resourceURI, nil)

	// test if first page is valid
	if pager.More() {
		res, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res.Value)
		testSerde(t, &res.MetricDefinitionCollection)
	} else {
		t.Fatal("no response")
	}
}

func TestNewListNamespacesPager_Success(t *testing.T) {
	client := startMetricsTest(t)

	pager := client.NewListNamespacesPager(resourceURI, &azquery.MetricsClientListNamespacesOptions{})

	// test if first page is valid
	if pager.More() {
		res, err := pager.NextPage(context.Background())
		if err != nil {
			t.Fatalf("failed to advance page: %v", err)
		}
		if res.Value == nil {
			t.Fatal("expected a response")
		}
		testSerde(t, &res.MetricNamespaceCollection)
	} else {
		t.Fatal("no response")
	}

}

func TestNewListNamespacesPager_Failure(t *testing.T) {
	client := startMetricsTest(t)
	invalidResourceURI := "123"
	var httpErr *azcore.ResponseError

	pager := client.NewListNamespacesPager(invalidResourceURI, nil)
	if pager.More() {
		res, err := pager.NextPage(context.Background())
		require.Error(t, err)
		require.ErrorAs(t, err, &httpErr)
		require.NotEqual(t, 200, httpErr.StatusCode)
		require.Nil(t, res.Value)
	} else {
		t.Fatal("no response")
	}

}

func TestMetricsAPIVersion(t *testing.T) {
	apiVersion := "2023-10-01"
	var requireVersion = func(t *testing.T) func(req *http.Request) bool {
		return func(r *http.Request) bool {
			version := r.URL.Query().Get("api-version")
			require.Equal(t, version, apiVersion)
			return true
		}
	}
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithStatusCode(200),
		mock.WithPredicate(requireVersion(t)),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	opts := &azquery.MetricsClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:  srv,
			APIVersion: apiVersion,
		},
	}
	client, err := azquery.NewMetricsClient(&azcred.Fake{}, opts)
	require.NoError(t, err)
	_, err = client.QueryResource(context.Background(), resourceURI, nil)
	require.NoError(t, err)
}

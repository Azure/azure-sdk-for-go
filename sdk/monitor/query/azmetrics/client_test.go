//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azmetrics_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics"
	"github.com/stretchr/testify/require"
)

func TestQueryResources_Pass(t *testing.T) {
	client := startTest(t)
	metricName := "HttpIncomingRequestCount"
	resourceIDList := azmetrics.ResourceIDList{ResourceIDs: []string{resourceURI}}

	res, err := client.QueryResources(
		context.Background(),
		subscriptionID,
		"Microsoft.AppConfiguration/configurationStores",
		[]string{"HttpIncomingRequestCount"},
		resourceIDList,
		&azmetrics.QueryResourcesOptions{
			Aggregation: to.Ptr("average"),
			StartTime:   to.Ptr("P1D"),
			Interval:    to.Ptr("PT1H"),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, resource := range res.Values {
		for _, metric := range resource.Values {
			require.Equal(t, metricName, *metric.Name.Value)
			for _, timeSeriesElement := range metric.TimeSeries {
				require.NotNil(t, timeSeriesElement)
			}
		}
	}

	testSerde(t, &res)
	testSerde(t, &resourceIDList)
}

func TestQueryResources_Fail(t *testing.T) {
	client := startTest(t)

	res, err := client.QueryResources(
		context.Background(),
		"fakesubscriptionID",
		"Microsoft.AppConfiguration/configurationStores",
		[]string{"HttpIncomingRequestCount"},
		azmetrics.ResourceIDList{ResourceIDs: []string{resourceURI}},
		nil,
	)
	require.Error(t, err)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.NotEqual(t, 200, httpErr.StatusCode)
	require.Nil(t, res.Values)

	testSerde(t, &res)
}

func TestAPIVersion(t *testing.T) {
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
	opts := &azmetrics.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:  srv,
			APIVersion: apiVersion,
		},
	}
	client, err := azmetrics.NewClient(endpoint, &azcred.Fake{}, opts)
	require.NoError(t, err)
	_, err = client.QueryResources(context.Background(), subscriptionID, "testing", []string{"test"}, azmetrics.ResourceIDList{}, nil)
	require.NoError(t, err)
}

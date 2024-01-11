//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/stretchr/testify/require"
)

func TestQueryBatch_Metrics(t *testing.T) {
	client := startMetricsBatchTest(t)
	metricName := "HttpIncomingRequestCount"
	resourceIDList := azquery.ResourceIDList{ResourceIDs: to.SliceOfPtrs(resourceURI)}

	res, err := client.QueryBatch(
		context.Background(),
		subscriptionID,
		"Microsoft.AppConfiguration/configurationStores",
		[]string{"HttpIncomingRequestCount"},
		resourceIDList,
		&azquery.MetricsBatchClientQueryBatchOptions{
			Aggregation: to.Ptr("average"),
			StartTime:   to.Ptr("2023-11-15"),
			EndTime:     to.Ptr("2023-11-16"),
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

func TestQueryBatch_MetricsFailure(t *testing.T) {
	client := startMetricsBatchTest(t)

	res, err := client.QueryBatch(
		context.Background(),
		"fakesubscriptionID",
		"Microsoft.AppConfiguration/configurationStores",
		[]string{"HttpIncomingRequestCount"},
		azquery.ResourceIDList{ResourceIDs: to.SliceOfPtrs(resourceURI)},
		nil,
	)
	require.Error(t, err)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "InvalidSubscriptionId")
	require.Equal(t, httpErr.StatusCode, 400)
	require.Nil(t, res.Values)

	testSerde(t, &res)
}

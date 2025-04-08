// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics/fake"
	"github.com/stretchr/testify/require"
)

var (
	resourceID1 = "test"
	resourceID2 = "test2"
)

func getServer() fake.Server {
	return fake.Server{
		QueryResources: func(ctx context.Context, subscriptionID string, metricNamespace string, metricNames []string, batchRequest azmetrics.ResourceIDList, options *azmetrics.QueryResourcesOptions) (resp azfake.Responder[azmetrics.QueryResourcesResponse], errResp azfake.ErrorResponder) {
			metricsResp := azmetrics.QueryResourcesResponse{
				MetricResults: azmetrics.MetricResults{
					Values: []azmetrics.MetricData{
						{
							ResourceID: &batchRequest.ResourceIDs[0],
						},
						{
							ResourceID: &batchRequest.ResourceIDs[1],
						},
					},
				},
			}
			resp.SetResponse(http.StatusOK, metricsResp, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()
	client, err := azmetrics.NewClient("https://fake.metrics.monitor.azure.com", &azfake.TokenCredential{}, &azmetrics.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	res, err := client.QueryResources(context.TODO(), "fake-sub-id", "fake-namespace", []string{"test"}, azmetrics.ResourceIDList{ResourceIDs: []string{resourceID1, resourceID2}}, &azmetrics.QueryResourcesOptions{})
	require.NoError(t, err)
	require.Len(t, res.Values, 2)
	require.Equal(t, resourceID1, *res.Values[0].ResourceID)
	require.Equal(t, resourceID2, *res.Values[1].ResourceID)
}

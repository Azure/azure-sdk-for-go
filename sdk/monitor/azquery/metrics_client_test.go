//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

const resourceURI1 = "/subscriptions/faa080af-c1d8-40ad-9cce-e1a450ca5b57/resourceGroups/ripark/providers/Microsoft.Cache/Redis/ripark"

func getMetricsClient(t *testing.T) *azquery.MetricsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	return azquery.NewMetricsClient(cred, nil)
}

func TestQueryResource_BasicQuerySuccess(t *testing.T) {
	client := getMetricsClient(t)
	res, err := client.QueryResource(context.Background(), resourceURI1, nil)
	if err != nil {
		t.Fatal("error")
	}
	if res.Response.Timespan == nil {
		t.Fatal("error")
	}
}

func TestNewListMetricDefinitionsPager_Success(t *testing.T) {
	//client := getMetricsClient(t)

	//pager := client.NewListMetricDefinitionsPager(resourceURI1, nil)
	//_ = res

}

func TestNewListMetricNamespacesPager_Success(t *testing.T) {
	//client := getMetricsClient(t)

}

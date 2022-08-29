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

func TestQueryResource_BasicQuerySuccess(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	client := azquery.NewMetricsClient(cred, nil)
	resourceURI := "/subscriptions/faa080af-c1d8-40ad-9cce-e1a450ca5b57/resourceGroups/srnagar-azuresdkgroup/providers/Microsoft.CognitiveServices/accounts/srnagara-textanalytics"
	res, err := client.QueryResource(context.Background(), resourceURI, nil)
	if err != nil {
		t.Fatal("error")
	}
	if res.Response.Timespan == nil {
		t.Fatal("error")
	}
}

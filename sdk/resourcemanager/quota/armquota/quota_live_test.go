//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armquota_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota/v2"
	"github.com/stretchr/testify/suite"
)

type QuotaTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *QuotaTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *QuotaTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestQuotaTestSuite(t *testing.T) {
	suite.Run(t, new(QuotaTestSuite))
}

// Microsoft.Quota/quotas/{resourceName}
func (testsuite *QuotaTestSuite) TestQuota() {
	var err error

	// From step Quota_List
	fmt.Println("Call operation: Quota_List")
	client, err := armquota.NewClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientNewListPager := client.NewListPager("subscriptions/"+testsuite.subscriptionId+"/providers/Microsoft.Network/locations/eastus", nil)
	for clientNewListPager.More() {
		_, err := clientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Quota_Get
	fmt.Println("Call operation: Quota_Get")
	_, err = client.Get(testsuite.ctx, "MinPublicIpInterNetworkPrefixLength", "subscriptions/"+testsuite.subscriptionId+"/providers/Microsoft.Network/locations/eastus", nil)
	testsuite.Require().NoError(err)

	// From step QuotaRequestStatus_List
	fmt.Println("Call operation: QuotaRequestStatus_List")
	requestStatusClient, err := armquota.NewRequestStatusClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	requestStatusClientNewListPager := requestStatusClient.NewListPager("subscriptions/"+testsuite.subscriptionId+"/providers/Microsoft.Network/locations/eastus", &armquota.RequestStatusClientListOptions{Filter: nil,
		Top:       nil,
		Skiptoken: nil,
	})
	for requestStatusClientNewListPager.More() {
		_, err := requestStatusClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		break
	}
}

// Microsoft.Quota/operations
func (testsuite *QuotaTestSuite) TestQuotaOperation() {
	var err error
	// From step QuotaOperation_List
	fmt.Println("Call operation: QuotaOperation_List")
	operationClient, err := armquota.NewOperationClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationClientNewListPager := operationClient.NewListPager(nil)
	for operationClientNewListPager.More() {
		_, err := operationClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Quota/usages/{resourceName}
func (testsuite *QuotaTestSuite) TestUsages() {
	var err error
	// From step Usages_List
	fmt.Println("Call operation: Usages_List")
	usagesClient, err := armquota.NewUsagesClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	usagesClientNewListPager := usagesClient.NewListPager("subscriptions/"+testsuite.subscriptionId+"/providers/Microsoft.Network/locations/eastus", nil)
	for usagesClientNewListPager.More() {
		_, err = usagesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

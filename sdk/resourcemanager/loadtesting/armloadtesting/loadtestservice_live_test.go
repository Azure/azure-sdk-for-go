// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armloadtesting_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/loadtesting/armloadtesting"
	"github.com/stretchr/testify/suite"
)

type LoadtestserviceTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	loadTestName      string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *LoadtestserviceTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.loadTestName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "loadtest", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *LoadtestserviceTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestLoadtestserviceTestSuite(t *testing.T) {
	suite.Run(t, new(LoadtestserviceTestSuite))
}

// Microsoft.LoadTestService/loadTests/{loadTestName}
func (testsuite *LoadtestserviceTestSuite) TestLoadTests() {
	var err error
	// From step LoadTests_CreateOrUpdate
	fmt.Println("Call operation: LoadTests_CreateOrUpdate")
	loadTestsClient, err := armloadtesting.NewLoadTestsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadTestsClientCreateOrUpdateResponsePoller, err := loadTestsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadTestName, armloadtesting.LoadTestResource{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"Team": to.Ptr("Dev Exp"),
		},
		Properties: &armloadtesting.LoadTestProperties{
			Description: to.Ptr("This is new load test resource"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, loadTestsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LoadTests_ListBySubscription
	fmt.Println("Call operation: LoadTests_ListBySubscription")
	loadTestsClientNewListBySubscriptionPager := loadTestsClient.NewListBySubscriptionPager(nil)
	for loadTestsClientNewListBySubscriptionPager.More() {
		_, err := loadTestsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LoadTests_Get
	fmt.Println("Call operation: LoadTests_Get")
	_, err = loadTestsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadTestName, nil)
	testsuite.Require().NoError(err)

	// From step LoadTests_ListOutboundNetworkDependenciesEndpoints
	fmt.Println("Call operation: LoadTests_ListOutboundNetworkDependenciesEndpoints")
	loadTestsClientNewListOutboundNetworkDependenciesEndpointsPager := loadTestsClient.NewListOutboundNetworkDependenciesEndpointsPager(testsuite.resourceGroupName, testsuite.loadTestName, nil)
	for loadTestsClientNewListOutboundNetworkDependenciesEndpointsPager.More() {
		_, err := loadTestsClientNewListOutboundNetworkDependenciesEndpointsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LoadTests_ListByResourceGroup
	fmt.Println("Call operation: LoadTests_ListByResourceGroup")
	loadTestsClientNewListByResourceGroupPager := loadTestsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for loadTestsClientNewListByResourceGroupPager.More() {
		_, err := loadTestsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LoadTests_Update
	fmt.Println("Call operation: LoadTests_Update")
	loadTestsClientUpdateResponsePoller, err := loadTestsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadTestName, armloadtesting.LoadTestResourcePatchRequestBody{
		Tags: map[string]*string{
			"Division": to.Ptr("LT"),
			"Team":     to.Ptr("Dev Exp"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, loadTestsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LoadTests_Delete
	fmt.Println("Call operation: LoadTests_Delete")
	loadTestsClientDeleteResponsePoller, err := loadTestsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadTestName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, loadTestsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.LoadTestService/operations
func (testsuite *LoadtestserviceTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armloadtesting.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.LoadTestService/locations/{location}/quotas/{quotaBucketName}
func (testsuite *LoadtestserviceTestSuite) TestQuotas() {
	var quotaBucketName string
	var err error
	// From step Quotas_List
	fmt.Println("Call operation: Quotas_List")
	quotasClient, err := armloadtesting.NewQuotasClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	quotasClientNewListPager := quotasClient.NewListPager(testsuite.location, nil)
	for quotasClientNewListPager.More() {
		nextResult, err := quotasClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		quotaBucketName = *nextResult.Value[0].Name
		break
	}

	// From step Quotas_Get
	fmt.Println("Call operation: Quotas_Get")
	_, err = quotasClient.Get(testsuite.ctx, testsuite.location, quotaBucketName, nil)
	testsuite.Require().NoError(err)

	// From step Quotas_CheckAvailability
	fmt.Println("Call operation: Quotas_CheckAvailability")
	_, err = quotasClient.CheckAvailability(testsuite.ctx, testsuite.location, quotaBucketName, armloadtesting.QuotaBucketRequest{
		Properties: &armloadtesting.QuotaBucketRequestProperties{
			CurrentQuota: to.Ptr[int32](40),
			CurrentUsage: to.Ptr[int32](20),
			Dimensions: &armloadtesting.QuotaBucketRequestPropertiesDimensions{
				Location:       to.Ptr(testsuite.location),
				SubscriptionID: to.Ptr(testsuite.subscriptionId),
			},
			NewQuota: to.Ptr[int32](50),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

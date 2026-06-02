// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmonitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/stretchr/testify/suite"
)

type DatacollectionendpointsTestSuite struct {
	suite.Suite

	ctx                        context.Context
	cred                       azcore.TokenCredential
	options                    *arm.ClientOptions
	dataCollectionEndpointName string
	location                   string
	resourceGroupName          string
	subscriptionId             string
}

func (testsuite *DatacollectionendpointsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.dataCollectionEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "datacollectionendpointna", 30, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *DatacollectionendpointsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDatacollectionendpointsTestSuite(t *testing.T) {
	suite.Run(t, new(DatacollectionendpointsTestSuite))
}

// Microsoft.Insights/dataCollectionEndpoints
func (testsuite *DatacollectionendpointsTestSuite) TestDatacollectionendpoint() {
	var err error
	// From step DataCollectionEndpoints_Create
	fmt.Println("Call operation: DataCollectionEndpoints_Create")
	dataCollectionEndpointsClient, err := armmonitor.NewDataCollectionEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = dataCollectionEndpointsClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.dataCollectionEndpointName, &armmonitor.DataCollectionEndpointsClientCreateOptions{
		Body: &armmonitor.DataCollectionEndpointResource{
			Location: to.Ptr(testsuite.location),
			Properties: &armmonitor.DataCollectionEndpointResourceProperties{
				NetworkACLs: &armmonitor.DataCollectionEndpointNetworkACLs{
					PublicNetworkAccess: to.Ptr(armmonitor.KnownPublicNetworkAccessOptionsEnabled),
				},
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step DataCollectionEndpoints_ListBySubscription
	fmt.Println("Call operation: DataCollectionEndpoints_ListBySubscription")
	dataCollectionEndpointsClientNewListBySubscriptionPager := dataCollectionEndpointsClient.NewListBySubscriptionPager(nil)
	for dataCollectionEndpointsClientNewListBySubscriptionPager.More() {
		_, err := dataCollectionEndpointsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataCollectionEndpoints_Get
	fmt.Println("Call operation: DataCollectionEndpoints_Get")
	_, err = dataCollectionEndpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.dataCollectionEndpointName, nil)
	testsuite.Require().NoError(err)

	// From step DataCollectionEndpoints_ListByResourceGroup
	fmt.Println("Call operation: DataCollectionEndpoints_ListByResourceGroup")
	dataCollectionEndpointsClientNewListByResourceGroupPager := dataCollectionEndpointsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for dataCollectionEndpointsClientNewListByResourceGroupPager.More() {
		_, err := dataCollectionEndpointsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataCollectionEndpoints_Update
	fmt.Println("Call operation: DataCollectionEndpoints_Update")
	_, err = dataCollectionEndpointsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.dataCollectionEndpointName, &armmonitor.DataCollectionEndpointsClientUpdateOptions{
		Body: &armmonitor.ResourceForUpdate{
			Tags: map[string]*string{
				"tag1": to.Ptr("A"),
				"tag2": to.Ptr("B"),
				"tag3": to.Ptr("C"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step DataCollectionRuleAssociations_ListByDataCollectionEndpoint
	fmt.Println("Call operation: DataCollectionRuleAssociations_ListByDataCollectionEndpoint")
	dataCollectionRuleAssociationsClient, err := armmonitor.NewDataCollectionRuleAssociationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dataCollectionRuleAssociationsClientNewListByDataCollectionEndpointPager := dataCollectionRuleAssociationsClient.NewListByDataCollectionEndpointPager(testsuite.resourceGroupName, testsuite.dataCollectionEndpointName, nil)
	for dataCollectionRuleAssociationsClientNewListByDataCollectionEndpointPager.More() {
		_, err := dataCollectionRuleAssociationsClientNewListByDataCollectionEndpointPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataCollectionEndpoints_Delete
	fmt.Println("Call operation: DataCollectionEndpoints_Delete")
	_, err = dataCollectionEndpointsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.dataCollectionEndpointName, nil)
	testsuite.Require().NoError(err)
}

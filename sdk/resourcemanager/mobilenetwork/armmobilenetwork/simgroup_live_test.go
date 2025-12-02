// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mobilenetwork/armmobilenetwork/v4"
	"github.com/stretchr/testify/suite"
)

type SimGroupTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	mobileNetworkId   string
	mobileNetworkName string
	simGroupName      string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SimGroupTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.simGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "simgroup", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SimGroupTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSimGroupTestSuite(t *testing.T) {
	suite.Run(t, new(SimGroupTestSuite))
}

func (testsuite *SimGroupTestSuite) Prepare() {
	var err error
	// From step MobileNetworks_CreateOrUpdate
	fmt.Println("Call operation: MobileNetworks_CreateOrUpdate")
	mobileNetworksClient, err := armmobilenetwork.NewMobileNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	mobileNetworksClientCreateOrUpdateResponsePoller, err := mobileNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, armmobilenetwork.MobileNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.PropertiesFormat{
			PublicLandMobileNetworkIdentifier: &armmobilenetwork.PlmnID{
				Mcc: to.Ptr("001"),
				Mnc: to.Ptr("01"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var mobileNetworksClientCreateOrUpdateResponse *armmobilenetwork.MobileNetworksClientCreateOrUpdateResponse
	mobileNetworksClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, mobileNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.mobileNetworkId = *mobileNetworksClientCreateOrUpdateResponse.ID
}

// Microsoft.MobileNetwork/simGroups/{simGroupName}
func (testsuite *SimGroupTestSuite) TestSimGroups() {
	var err error
	// From step SimGroups_CreateOrUpdate
	fmt.Println("Call operation: SimGroups_CreateOrUpdate")
	simGroupsClient, err := armmobilenetwork.NewSimGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	simGroupsClientCreateOrUpdateResponsePoller, err := simGroupsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, armmobilenetwork.SimGroup{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.SimGroupPropertiesFormat{
			MobileNetwork: &armmobilenetwork.ResourceID{
				ID: to.Ptr(testsuite.mobileNetworkId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simGroupsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step SimGroups_ListBySubscription
	fmt.Println("Call operation: SimGroups_ListBySubscription")
	simGroupsClientNewListBySubscriptionPager := simGroupsClient.NewListBySubscriptionPager(nil)
	for simGroupsClientNewListBySubscriptionPager.More() {
		_, err := simGroupsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SimGroups_Get
	fmt.Println("Call operation: SimGroups_Get")
	_, err = simGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, nil)
	testsuite.Require().NoError(err)

	// From step SimGroups_ListByResourceGroup
	fmt.Println("Call operation: SimGroups_ListByResourceGroup")
	simGroupsClientNewListByResourceGroupPager := simGroupsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for simGroupsClientNewListByResourceGroupPager.More() {
		_, err := simGroupsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SimGroups_UpdateTags
	fmt.Println("Call operation: SimGroups_UpdateTags")
	_, err = simGroupsClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, armmobilenetwork.IdentityAndTagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step SimGroups_Delete
	fmt.Println("Call operation: SimGroups_Delete")
	simGroupsClientDeleteResponsePoller, err := simGroupsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simGroupsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

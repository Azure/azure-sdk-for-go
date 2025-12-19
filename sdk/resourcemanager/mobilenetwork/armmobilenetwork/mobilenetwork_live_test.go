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

type MobileNetworkTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	mobileNetworkName string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *MobileNetworkTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *MobileNetworkTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestMobileNetworkTestSuite(t *testing.T) {
	suite.Run(t, new(MobileNetworkTestSuite))
}

// Microsoft.MobileNetwork/mobileNetworks/{mobileNetworkName}
func (testsuite *MobileNetworkTestSuite) TestMobileNetworks() {
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
	_, err = testutil.PollForTest(testsuite.ctx, mobileNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step MobileNetworks_ListBySubscription
	fmt.Println("Call operation: MobileNetworks_ListBySubscription")
	mobileNetworksClientNewListBySubscriptionPager := mobileNetworksClient.NewListBySubscriptionPager(nil)
	for mobileNetworksClientNewListBySubscriptionPager.More() {
		_, err := mobileNetworksClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MobileNetworks_ListByResourceGroup
	fmt.Println("Call operation: MobileNetworks_ListByResourceGroup")
	mobileNetworksClientNewListByResourceGroupPager := mobileNetworksClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for mobileNetworksClientNewListByResourceGroupPager.More() {
		_, err := mobileNetworksClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MobileNetworks_Get
	fmt.Println("Call operation: MobileNetworks_Get")
	_, err = mobileNetworksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, nil)
	testsuite.Require().NoError(err)

	// From step MobileNetworks_UpdateTags
	fmt.Println("Call operation: MobileNetworks_UpdateTags")
	_, err = mobileNetworksClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, armmobilenetwork.IdentityAndTagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step MobileNetworks_Delete
	fmt.Println("Call operation: MobileNetworks_Delete")
	mobileNetworksClientDeleteResponsePoller, err := mobileNetworksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mobileNetworksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

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

type DataNetworkTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	dataNetworkName   string
	mobileNetworkName string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *DataNetworkTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.dataNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "datanetw", 14, false)
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DataNetworkTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDataNetworkTestSuite(t *testing.T) {
	suite.Run(t, new(DataNetworkTestSuite))
}

func (testsuite *DataNetworkTestSuite) Prepare() {
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
}

// Microsoft.MobileNetwork/mobileNetworks/{mobileNetworkName}/dataNetworks/{dataNetworkName}
func (testsuite *DataNetworkTestSuite) TestDataNetworks() {
	var err error
	// From step DataNetworks_CreateOrUpdate
	fmt.Println("Call operation: DataNetworks_CreateOrUpdate")
	dataNetworksClient, err := armmobilenetwork.NewDataNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dataNetworksClientCreateOrUpdateResponsePoller, err := dataNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.dataNetworkName, armmobilenetwork.DataNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.DataNetworkPropertiesFormat{
			Description: to.Ptr("myFavouriteDataNetwork"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dataNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DataNetworks_ListByMobileNetwork
	fmt.Println("Call operation: DataNetworks_ListByMobileNetwork")
	dataNetworksClientNewListByMobileNetworkPager := dataNetworksClient.NewListByMobileNetworkPager(testsuite.resourceGroupName, testsuite.mobileNetworkName, nil)
	for dataNetworksClientNewListByMobileNetworkPager.More() {
		_, err := dataNetworksClientNewListByMobileNetworkPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataNetworks_Get
	fmt.Println("Call operation: DataNetworks_Get")
	_, err = dataNetworksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.dataNetworkName, nil)
	testsuite.Require().NoError(err)

	// From step DataNetworks_UpdateTags
	fmt.Println("Call operation: DataNetworks_UpdateTags")
	_, err = dataNetworksClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.dataNetworkName, armmobilenetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step DataNetworks_Delete
	fmt.Println("Call operation: DataNetworks_Delete")
	dataNetworksClientDeleteResponsePoller, err := dataNetworksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.dataNetworkName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dataNetworksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

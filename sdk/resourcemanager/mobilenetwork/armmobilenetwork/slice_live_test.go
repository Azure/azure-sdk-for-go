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

type SliceTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	mobileNetworkName string
	sliceName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SliceTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.sliceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "slicenam", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SliceTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSliceTestSuite(t *testing.T) {
	suite.Run(t, new(SliceTestSuite))
}

func (testsuite *SliceTestSuite) Prepare() {
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

// Microsoft.MobileNetwork/mobileNetworks/{mobileNetworkName}/slices/{sliceName}
func (testsuite *SliceTestSuite) TestSlices() {
	var err error
	// From step Slices_CreateOrUpdate
	fmt.Println("Call operation: Slices_CreateOrUpdate")
	slicesClient, err := armmobilenetwork.NewSlicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	slicesClientCreateOrUpdateResponsePoller, err := slicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.sliceName, armmobilenetwork.Slice{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.SlicePropertiesFormat{
			Description: to.Ptr("myFavouriteSlice"),
			Snssai: &armmobilenetwork.Snssai{
				Sd:  to.Ptr("1abcde"),
				Sst: to.Ptr[int32](1),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, slicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Slices_ListByMobileNetwork
	fmt.Println("Call operation: Slices_ListByMobileNetwork")
	slicesClientNewListByMobileNetworkPager := slicesClient.NewListByMobileNetworkPager(testsuite.resourceGroupName, testsuite.mobileNetworkName, nil)
	for slicesClientNewListByMobileNetworkPager.More() {
		_, err := slicesClientNewListByMobileNetworkPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Slices_Get
	fmt.Println("Call operation: Slices_Get")
	_, err = slicesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.sliceName, nil)
	testsuite.Require().NoError(err)

	// From step Slices_UpdateTags
	fmt.Println("Call operation: Slices_UpdateTags")
	_, err = slicesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.sliceName, armmobilenetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Slices_Delete
	fmt.Println("Call operation: Slices_Delete")
	slicesClientDeleteResponsePoller, err := slicesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.sliceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, slicesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

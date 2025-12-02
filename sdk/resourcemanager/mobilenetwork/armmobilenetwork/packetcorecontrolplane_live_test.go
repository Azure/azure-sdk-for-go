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

type PacketCoreControlPlaneTestSuite struct {
	suite.Suite

	ctx                        context.Context
	cred                       azcore.TokenCredential
	options                    *arm.ClientOptions
	armEndpoint                string
	mobileNetworkName          string
	packetCoreControlPlaneName string
	siteName                   string
	sitesId                    string
	location                   string
	resourceGroupName          string
	subscriptionId             string
}

func (testsuite *PacketCoreControlPlaneTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.packetCoreControlPlaneName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "packetco", 14, false)
	testsuite.siteName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sitename", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *PacketCoreControlPlaneTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPacketCoreControlPlaneTestSuite(t *testing.T) {
	suite.Run(t, new(PacketCoreControlPlaneTestSuite))
}

func (testsuite *PacketCoreControlPlaneTestSuite) Prepare() {
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

	// From step Sites_CreateOrUpdate
	fmt.Println("Call operation: Sites_CreateOrUpdate")
	sitesClient, err := armmobilenetwork.NewSitesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sitesClientCreateOrUpdateResponsePoller, err := sitesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.siteName, armmobilenetwork.Site{
		Location:   to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.SitePropertiesFormat{},
	}, nil)
	testsuite.Require().NoError(err)
	var sitesClientCreateOrUpdateResponse *armmobilenetwork.SitesClientCreateOrUpdateResponse
	sitesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, sitesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.sitesId = *sitesClientCreateOrUpdateResponse.ID
}

// Microsoft.MobileNetwork/packetCoreControlPlanes/{packetCoreControlPlaneName}
func (testsuite *PacketCoreControlPlaneTestSuite) TestPacketCoreControlPlanes() {
	var err error
	// From step PacketCoreControlPlanes_CreateOrUpdate
	fmt.Println("Call operation: PacketCoreControlPlanes_CreateOrUpdate")
	packetCoreControlPlanesClient, err := armmobilenetwork.NewPacketCoreControlPlanesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	packetCoreControlPlanesClientCreateOrUpdateResponsePoller, err := packetCoreControlPlanesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, armmobilenetwork.PacketCoreControlPlane{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.PacketCoreControlPlanePropertiesFormat{
			ControlPlaneAccessInterface: &armmobilenetwork.InterfaceProperties{
				Name:        to.Ptr("N2"),
				IPv4Address: to.Ptr("2.4.0.1"),
			},
			CoreNetworkTechnology: to.Ptr(armmobilenetwork.CoreNetworkTypeFiveGC),
			Installation: &armmobilenetwork.Installation{
				DesiredState: to.Ptr(armmobilenetwork.DesiredInstallationStateInstalled),
			},
			LocalDiagnosticsAccess: &armmobilenetwork.LocalDiagnosticsAccessConfiguration{
				AuthenticationType: to.Ptr(armmobilenetwork.AuthenticationTypeAAD),
			},
			Platform: &armmobilenetwork.PlatformConfiguration{
				Type: to.Ptr(armmobilenetwork.PlatformTypeAKSHCI),
			},
			Sites: []*armmobilenetwork.SiteResourceID{
				{
					ID: to.Ptr(testsuite.sitesId),
				}},
			SKU:   to.Ptr(armmobilenetwork.BillingSKUG0),
			UeMtu: to.Ptr[int32](1600),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, packetCoreControlPlanesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PacketCoreControlPlanes_ListBySubscription
	fmt.Println("Call operation: PacketCoreControlPlanes_ListBySubscription")
	packetCoreControlPlanesClientNewListBySubscriptionPager := packetCoreControlPlanesClient.NewListBySubscriptionPager(nil)
	for packetCoreControlPlanesClientNewListBySubscriptionPager.More() {
		_, err := packetCoreControlPlanesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PacketCoreControlPlanes_Get
	fmt.Println("Call operation: PacketCoreControlPlanes_Get")
	_, err = packetCoreControlPlanesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, nil)
	testsuite.Require().NoError(err)

	// From step PacketCoreControlPlanes_ListByResourceGroup
	fmt.Println("Call operation: PacketCoreControlPlanes_ListByResourceGroup")
	packetCoreControlPlanesClientNewListByResourceGroupPager := packetCoreControlPlanesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for packetCoreControlPlanesClientNewListByResourceGroupPager.More() {
		_, err := packetCoreControlPlanesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PacketCoreControlPlanes_Delete
	fmt.Println("Call operation: PacketCoreControlPlanes_Delete")
	packetCoreControlPlanesClientDeleteResponsePoller, err := packetCoreControlPlanesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, packetCoreControlPlanesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.MobileNetwork/packetCoreControlPlaneVersions
func (testsuite *PacketCoreControlPlaneTestSuite) TestPacketCoreControlPlaneVersions() {
	var err error
	// From step PacketCoreControlPlaneVersions_List
	fmt.Println("Call operation: PacketCoreControlPlaneVersions_List")
	packetCoreControlPlaneVersionsClient, err := armmobilenetwork.NewPacketCoreControlPlaneVersionsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	packetCoreControlPlaneVersionsClientNewListPager := packetCoreControlPlaneVersionsClient.NewListPager(nil)
	for packetCoreControlPlaneVersionsClientNewListPager.More() {
		_, err := packetCoreControlPlaneVersionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PacketCoreControlPlaneVersions_Get
	fmt.Println("Call operation: PacketCoreControlPlaneVersions_Get")
	_, err = packetCoreControlPlaneVersionsClient.Get(testsuite.ctx, "PMN-4-11-1", nil)
	testsuite.Require().NoError(err)

	// From step PacketCoreControlPlaneVersions_ListBySubscription
	fmt.Println("Call operation: PacketCoreControlPlaneVersions_ListBySubscription")
	packetCoreControlPlaneVersionsClientNewListBySubscriptionPager := packetCoreControlPlaneVersionsClient.NewListBySubscriptionPager(testsuite.subscriptionId, nil)
	for packetCoreControlPlaneVersionsClientNewListBySubscriptionPager.More() {
		_, err := packetCoreControlPlaneVersionsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

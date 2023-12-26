//go:build go1.18
// +build go1.18

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mobilenetwork/armmobilenetwork/v3"
	"github.com/stretchr/testify/suite"
)

type AttachedDataNetworkTestSuite struct {
	suite.Suite

	ctx                        context.Context
	cred                       azcore.TokenCredential
	options                    *arm.ClientOptions
	armEndpoint                string
	attachedDataNetworkName    string
	mobileNetworkName          string
	packetCoreControlPlaneName string
	packetCoreDataPlaneName    string
	siteName                   string
	sitesId                    string
	location                   string
	resourceGroupName          string
	subscriptionId             string
}

func (testsuite *AttachedDataNetworkTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/mobilenetwork/armmobilenetwork/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.attachedDataNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "attached", 14, false)
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.packetCoreControlPlaneName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "packetco", 14, false)
	testsuite.packetCoreDataPlaneName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "packetco", 14, false)
	testsuite.siteName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sitename", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AttachedDataNetworkTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAttachedDataNetworkTestSuite(t *testing.T) {
	suite.Run(t, new(AttachedDataNetworkTestSuite))
}

func (testsuite *AttachedDataNetworkTestSuite) Prepare() {
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

	// From step PacketCoreDataPlanes_CreateOrUpdate
	fmt.Println("Call operation: PacketCoreDataPlanes_CreateOrUpdate")
	packetCoreDataPlanesClient, err := armmobilenetwork.NewPacketCoreDataPlanesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	packetCoreDataPlanesClientCreateOrUpdateResponsePoller, err := packetCoreDataPlanesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, testsuite.packetCoreDataPlaneName, armmobilenetwork.PacketCoreDataPlane{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.PacketCoreDataPlanePropertiesFormat{
			UserPlaneAccessInterface: &armmobilenetwork.InterfaceProperties{
				Name: to.Ptr("N3"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, packetCoreDataPlanesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DataNetworks_CreateOrUpdate
	fmt.Println("Call operation: DataNetworks_CreateOrUpdate")
	dataNetworksClient, err := armmobilenetwork.NewDataNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dataNetworksClientCreateOrUpdateResponsePoller, err := dataNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.attachedDataNetworkName, armmobilenetwork.DataNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.DataNetworkPropertiesFormat{
			Description: to.Ptr("myFavouriteDataNetwork"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dataNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.MobileNetwork/packetCoreControlPlanes/{packetCoreControlPlaneName}/packetCoreDataPlanes/{packetCoreDataPlaneName}/attachedDataNetworks/{attachedDataNetworkName}
func (testsuite *AttachedDataNetworkTestSuite) TestAttachedDataNetworks() {
	var err error
	// From step AttachedDataNetworks_CreateOrUpdate
	fmt.Println("Call operation: AttachedDataNetworks_CreateOrUpdate")
	attachedDataNetworksClient, err := armmobilenetwork.NewAttachedDataNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	attachedDataNetworksClientCreateOrUpdateResponsePoller, err := attachedDataNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, testsuite.packetCoreDataPlaneName, testsuite.attachedDataNetworkName, armmobilenetwork.AttachedDataNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.AttachedDataNetworkPropertiesFormat{
			DNSAddresses: []*string{
				to.Ptr("1.1.1.1")},
			NaptConfiguration: &armmobilenetwork.NaptConfiguration{
				Enabled:       to.Ptr(armmobilenetwork.NaptEnabledEnabled),
				PinholeLimits: to.Ptr[int32](65536),
				PinholeTimeouts: &armmobilenetwork.PinholeTimeouts{
					Icmp: to.Ptr[int32](30),
					TCP:  to.Ptr[int32](180),
					UDP:  to.Ptr[int32](30),
				},
				PortRange: &armmobilenetwork.PortRange{
					MaxPort: to.Ptr[int32](49999),
					MinPort: to.Ptr[int32](1024),
				},
				PortReuseHoldTime: &armmobilenetwork.PortReuseHoldTimes{
					TCP: to.Ptr[int32](120),
					UDP: to.Ptr[int32](60),
				},
			},
			UserEquipmentAddressPoolPrefix: []*string{
				to.Ptr("2.2.0.0/16")},
			UserEquipmentStaticAddressPoolPrefix: []*string{
				to.Ptr("2.4.0.0/16")},
			UserPlaneDataInterface: &armmobilenetwork.InterfaceProperties{
				Name: to.Ptr("N6"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, attachedDataNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AttachedDataNetworks_ListByPacketCoreDataPlane
	fmt.Println("Call operation: AttachedDataNetworks_ListByPacketCoreDataPlane")
	attachedDataNetworksClientNewListByPacketCoreDataPlanePager := attachedDataNetworksClient.NewListByPacketCoreDataPlanePager(testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, testsuite.packetCoreDataPlaneName, nil)
	for attachedDataNetworksClientNewListByPacketCoreDataPlanePager.More() {
		_, err := attachedDataNetworksClientNewListByPacketCoreDataPlanePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AttachedDataNetworks_Get
	fmt.Println("Call operation: AttachedDataNetworks_Get")
	_, err = attachedDataNetworksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, testsuite.packetCoreDataPlaneName, testsuite.attachedDataNetworkName, nil)
	testsuite.Require().NoError(err)

	// From step AttachedDataNetworks_UpdateTags
	fmt.Println("Call operation: AttachedDataNetworks_UpdateTags")
	_, err = attachedDataNetworksClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, testsuite.packetCoreDataPlaneName, testsuite.attachedDataNetworkName, armmobilenetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AttachedDataNetworks_Delete
	fmt.Println("Call operation: AttachedDataNetworks_Delete")
	attachedDataNetworksClientDeleteResponsePoller, err := attachedDataNetworksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.packetCoreControlPlaneName, testsuite.packetCoreDataPlaneName, testsuite.attachedDataNetworkName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, attachedDataNetworksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

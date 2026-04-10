// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armprivatedns_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type PrivatednsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	globalLocation    string
	privateZoneName   string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *PrivatednsTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.globalLocation = "Global"
	testsuite.privateZoneName = "scenario_privatezone.com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	testutil.StartRecording(testsuite.T(), pathToPackage)
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *PrivatednsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPrivatednsTestSuite(t *testing.T) {
	suite.Run(t, new(PrivatednsTestSuite))
}

func (testsuite *PrivatednsTestSuite) Prepare() {
	var err error
	// From step PrivateZone_Create
	fmt.Println("Call operation: PrivateZones_Create")
	privateZonesClient, err := armprivatedns.NewPrivateZonesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateZonesClientCreateOrUpdateResponsePoller, err := privateZonesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, armprivatedns.PrivateZone{
		Location: to.Ptr(testsuite.globalLocation),
	}, &armprivatedns.PrivateZonesClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateZonesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/privateDnsZones
func (testsuite *PrivatednsTestSuite) TestPrivatezone() {
	var err error
	// From step PrivateZone_Update
	fmt.Println("Call operation: PrivateZones_Update")
	privateZonesClient, err := armprivatedns.NewPrivateZonesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateZonesClientUpdateResponsePoller, err := privateZonesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, armprivatedns.PrivateZone{
		Tags: map[string]*string{
			"key2": to.Ptr("value2"),
		},
	}, &armprivatedns.PrivateZonesClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateZonesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateZone_Get
	fmt.Println("Call operation: PrivateZones_Get")
	_, err = privateZonesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateZone_ListBySubscription
	fmt.Println("Call operation: PrivateZones_List")
	privateZonesClientNewListPager := privateZonesClient.NewListPager(&armprivatedns.PrivateZonesClientListOptions{Top: nil})
	for privateZonesClientNewListPager.More() {
		_, err := privateZonesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateZone_ListByResourceGroup
	fmt.Println("Call operation: PrivateZones_ListByResourceGroup")
	privateZonesClientNewListByResourceGroupPager := privateZonesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armprivatedns.PrivateZonesClientListByResourceGroupOptions{Top: nil})
	for privateZonesClientNewListByResourceGroupPager.More() {
		_, err := privateZonesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Network/privateDnsZones/{recordType}/{relativeRecordSetName}
func (testsuite *PrivatednsTestSuite) TestRecordset() {
	relativeRecordSetName := "scenario_recordset"
	var err error
	// From step RecordSet_Create
	fmt.Println("Call operation: RecordSets_CreateOrUpdate")
	recordSetsClient, err := armprivatedns.NewRecordSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = recordSetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, armprivatedns.RecordTypeA, relativeRecordSetName, armprivatedns.RecordSet{
		Properties: &armprivatedns.RecordSetProperties{
			ARecords: []*armprivatedns.ARecord{
				{
					IPv4Address: to.Ptr("1.2.3.4"),
				}},
			Metadata: map[string]*string{
				"key1": to.Ptr("value1"),
			},
			TTL: to.Ptr[int64](3600),
		},
	}, &armprivatedns.RecordSetsClientCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)

	// From step RecordSet_Update
	fmt.Println("Call operation: RecordSets_Update")
	_, err = recordSetsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, armprivatedns.RecordTypeA, relativeRecordSetName, armprivatedns.RecordSet{
		Properties: &armprivatedns.RecordSetProperties{
			Metadata: map[string]*string{
				"key2": to.Ptr("value2"),
			},
		},
	}, &armprivatedns.RecordSetsClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step RecordSet_Get
	fmt.Println("Call operation: RecordSets_Get")
	_, err = recordSetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, armprivatedns.RecordTypeA, relativeRecordSetName, nil)
	testsuite.Require().NoError(err)

	// From step RecordSet_ListByType
	fmt.Println("Call operation: RecordSets_ListByType")
	recordSetsClientNewListByTypePager := recordSetsClient.NewListByTypePager(testsuite.resourceGroupName, testsuite.privateZoneName, armprivatedns.RecordTypeA, &armprivatedns.RecordSetsClientListByTypeOptions{Top: nil,
		Recordsetnamesuffix: nil,
	})
	for recordSetsClientNewListByTypePager.More() {
		_, err := recordSetsClientNewListByTypePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RecordSet_List
	fmt.Println("Call operation: RecordSets_List")
	recordSetsClientNewListPager := recordSetsClient.NewListPager(testsuite.resourceGroupName, testsuite.privateZoneName, &armprivatedns.RecordSetsClientListOptions{Top: nil,
		Recordsetnamesuffix: nil,
	})
	for recordSetsClientNewListPager.More() {
		_, err := recordSetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RecordSet_Delete
	fmt.Println("Call operation: RecordSets_Delete")
	_, err = recordSetsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, armprivatedns.RecordTypeA, relativeRecordSetName, &armprivatedns.RecordSetsClientDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.Network/privateDnsZones/virtualNetworkLinks
func (testsuite *PrivatednsTestSuite) TestVirtualnetworklink() {
	var virtaulNetworkId string
	virtualNetworkLinkName := "scenario_vnetlink"
	var err error
	// From step NetworkAndSubnet_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"virtaulNetworkId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks', parameters('name'))]",
			},
		},
		"parameters": map[string]interface{}{
			"name": map[string]interface{}{
				"type":         "string",
				"defaultValue": "pn-network",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('name')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"location":   "eastus",
				"properties": map[string]interface{}{
					"addressSpace": map[string]interface{}{
						"addressPrefixes": []interface{}{
							"10.0.0.0/16",
						},
					},
					"subnets": []interface{}{
						map[string]interface{}{
							"name": "test-1",
							"properties": map[string]interface{}{
								"addressPrefix": "10.0.0.0/24",
							},
						},
					},
				},
				"tags": map[string]interface{}{},
			},
		},
	}
	params := map[string]interface{}{}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "NetworkAndSubnet_Create", &deployment)
	testsuite.Require().NoError(err)
	virtaulNetworkId = deploymentExtend.Properties.Outputs.(map[string]interface{})["virtaulNetworkId"].(map[string]interface{})["value"].(string)

	// From step VirtualNetworkLink_Create
	fmt.Println("Call operation: VirtualNetworkLinks_CreateOrUpdate")
	virtualNetworkLinksClient, err := armprivatedns.NewVirtualNetworkLinksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworkLinksClientCreateOrUpdateResponsePoller, err := virtualNetworkLinksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, virtualNetworkLinkName, armprivatedns.VirtualNetworkLink{
		Location: to.Ptr(testsuite.globalLocation),
		Properties: &armprivatedns.VirtualNetworkLinkProperties{
			RegistrationEnabled: to.Ptr(false),
			VirtualNetwork: &armprivatedns.SubResource{
				ID: to.Ptr(virtaulNetworkId),
			},
		},
	}, &armprivatedns.VirtualNetworkLinksClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkLinksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkLink_Update
	fmt.Println("Call operation: VirtualNetworkLinks_Update")
	virtualNetworkLinksClientUpdateResponsePoller, err := virtualNetworkLinksClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, virtualNetworkLinkName, armprivatedns.VirtualNetworkLink{
		Tags: map[string]*string{
			"key2": to.Ptr("value2"),
		},
		Properties: &armprivatedns.VirtualNetworkLinkProperties{
			RegistrationEnabled: to.Ptr(true),
		},
	}, &armprivatedns.VirtualNetworkLinksClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkLinksClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkLink_Get
	fmt.Println("Call operation: VirtualNetworkLinks_Get")
	_, err = virtualNetworkLinksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, virtualNetworkLinkName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkLink_List
	fmt.Println("Call operation: VirtualNetworkLinks_List")
	virtualNetworkLinksClientNewListPager := virtualNetworkLinksClient.NewListPager(testsuite.resourceGroupName, testsuite.privateZoneName, &armprivatedns.VirtualNetworkLinksClientListOptions{Top: nil})
	for virtualNetworkLinksClientNewListPager.More() {
		_, err := virtualNetworkLinksClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworkLink_Delete
	fmt.Println("Call operation: VirtualNetworkLinks_Delete")
	virtualNetworkLinksClientDeleteResponsePoller, err := virtualNetworkLinksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.privateZoneName, virtualNetworkLinkName, &armprivatedns.VirtualNetworkLinksClientBeginDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkLinksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

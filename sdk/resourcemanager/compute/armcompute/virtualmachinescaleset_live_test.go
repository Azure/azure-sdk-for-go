// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type VirtualMachineScaleSetTestSuite struct {
	suite.Suite

	ctx                      context.Context
	cred                     azcore.TokenCredential
	options                  *arm.ClientOptions
	adminUsername            string
	subnetId                 string
	virtualNetworkSubnetName string
	vmScaleSetName           string
	vmssExtensionName        string
	adminPassword            string
	location                 string
	resourceGroupName        string
	subscriptionId           string
}

func (testsuite *VirtualMachineScaleSetTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.adminUsername, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmuserna", 14, true)
	testsuite.virtualNetworkSubnetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmssvnetna", 16, false)
	testsuite.vmScaleSetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmscaleset", 16, false)
	testsuite.vmssExtensionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmssextens", 16, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *VirtualMachineScaleSetTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestVirtualMachineScaleSetTestSuite(t *testing.T) {
	suite.Run(t, new(VirtualMachineScaleSetTestSuite))
}

func (testsuite *VirtualMachineScaleSetTestSuite) Prepare() {
	var err error
	// From step Create_NetworkAndSubnet
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"subnetId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworkSubnetName'), 'default')]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworkSubnetName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.virtualNetworkSubnetName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworkSubnetName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"subnets": []any{
						map[string]any{
							"name": "default",
							"properties": map[string]any{
								"addressPrefix": "10.0.0.0/24",
							},
						},
					},
				},
			},
		},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_NetworkAndSubnet", &deployment)
	testsuite.Require().NoError(err)
	testsuite.subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step VirtualMachineScaleSets_CreateOrUpdate
	fmt.Println("Call operation: VirtualMachineScaleSets_CreateOrUpdate")
	virtualMachineScaleSetsClient, err := armcompute.NewVirtualMachineScaleSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachineScaleSetsClientCreateOrUpdateResponsePoller, err := virtualMachineScaleSetsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, armcompute.VirtualMachineScaleSet{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.VirtualMachineScaleSetProperties{
			Overprovision: to.Ptr(true),
			UpgradePolicy: &armcompute.UpgradePolicy{
				Mode: to.Ptr(armcompute.UpgradeModeManual),
			},
			VirtualMachineProfile: &armcompute.VirtualMachineScaleSetVMProfile{
				NetworkProfile: &armcompute.VirtualMachineScaleSetNetworkProfile{
					NetworkInterfaceConfigurations: []*armcompute.VirtualMachineScaleSetNetworkConfiguration{
						{
							Name: to.Ptr(testsuite.vmScaleSetName),
							Properties: &armcompute.VirtualMachineScaleSetNetworkConfigurationProperties{
								EnableIPForwarding: to.Ptr(true),
								IPConfigurations: []*armcompute.VirtualMachineScaleSetIPConfiguration{
									{
										Name: to.Ptr(testsuite.vmScaleSetName),
										Properties: &armcompute.VirtualMachineScaleSetIPConfigurationProperties{
											Subnet: &armcompute.APIEntityReference{
												ID: to.Ptr(testsuite.subnetId),
											},
										},
									}},
								Primary: to.Ptr(true),
							},
						}},
				},
				OSProfile: &armcompute.VirtualMachineScaleSetOSProfile{
					AdminPassword:      to.Ptr(testsuite.adminPassword),
					AdminUsername:      to.Ptr(testsuite.adminUsername),
					ComputerNamePrefix: to.Ptr("vmss"),
				},
				StorageProfile: &armcompute.VirtualMachineScaleSetStorageProfile{
					ImageReference: &armcompute.ImageReference{
						Offer:     to.Ptr("WindowsServer"),
						Publisher: to.Ptr("MicrosoftWindowsServer"),
						SKU:       to.Ptr("2016-Datacenter"),
						Version:   to.Ptr("latest"),
					},
					OSDisk: &armcompute.VirtualMachineScaleSetOSDisk{
						Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
						CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
						ManagedDisk: &armcompute.VirtualMachineScaleSetManagedDiskParameters{
							StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS),
						},
					},
				},
			},
		},
		SKU: &armcompute.SKU{
			Name:     to.Ptr("Standard_D1_v2"),
			Capacity: to.Ptr[int64](3),
			Tier:     to.Ptr("Standard"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}
func (testsuite *VirtualMachineScaleSetTestSuite) TestVirtualMachineScaleSets() {
	var err error
	// From step VirtualMachineScaleSets_ListByLocation
	fmt.Println("Call operation: VirtualMachineScaleSets_ListByLocation")
	virtualMachineScaleSetsClient, err := armcompute.NewVirtualMachineScaleSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachineScaleSetsClientNewListByLocationPager := virtualMachineScaleSetsClient.NewListByLocationPager(testsuite.location, nil)
	for virtualMachineScaleSetsClientNewListByLocationPager.More() {
		_, err := virtualMachineScaleSetsClientNewListByLocationPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineScaleSets_GetInstanceView
	fmt.Println("Call operation: VirtualMachineScaleSets_GetInstanceView")
	_, err = virtualMachineScaleSetsClient.GetInstanceView(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_GetOSUpgradeHistory
	fmt.Println("Call operation: VirtualMachineScaleSets_GetOSUpgradeHistory")
	virtualMachineScaleSetsClientNewGetOSUpgradeHistoryPager := virtualMachineScaleSetsClient.NewGetOSUpgradeHistoryPager(testsuite.resourceGroupName, testsuite.vmScaleSetName, nil)
	for virtualMachineScaleSetsClientNewGetOSUpgradeHistoryPager.More() {
		_, err := virtualMachineScaleSetsClientNewGetOSUpgradeHistoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineScaleSets_ListAll
	fmt.Println("Call operation: VirtualMachineScaleSets_ListAll")
	virtualMachineScaleSetsClientNewListAllPager := virtualMachineScaleSetsClient.NewListAllPager(nil)
	for virtualMachineScaleSetsClientNewListAllPager.More() {
		_, err := virtualMachineScaleSetsClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineScaleSets_Get
	fmt.Println("Call operation: VirtualMachineScaleSets_Get")
	_, err = virtualMachineScaleSetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_List
	fmt.Println("Call operation: VirtualMachineScaleSets_List")
	virtualMachineScaleSetsClientNewListPager := virtualMachineScaleSetsClient.NewListPager(testsuite.resourceGroupName, nil)
	for virtualMachineScaleSetsClientNewListPager.More() {
		_, err := virtualMachineScaleSetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineScaleSets_ListSkus
	fmt.Println("Call operation: VirtualMachineScaleSets_ListSKUs")
	virtualMachineScaleSetsClientNewListSKUsPager := virtualMachineScaleSetsClient.NewListSKUsPager(testsuite.resourceGroupName, testsuite.vmScaleSetName, nil)
	for virtualMachineScaleSetsClientNewListSKUsPager.More() {
		_, err := virtualMachineScaleSetsClientNewListSKUsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineScaleSets_Update
	fmt.Println("Call operation: VirtualMachineScaleSets_Update")
	virtualMachineScaleSetsClientUpdateResponsePoller, err := virtualMachineScaleSetsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, armcompute.VirtualMachineScaleSetUpdate{}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_Redeploy
	fmt.Println("Call operation: VirtualMachineScaleSets_Redeploy")
	virtualMachineScaleSetsClientRedeployResponsePoller, err := virtualMachineScaleSetsClient.BeginRedeploy(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientBeginRedeployOptions{VMInstanceIDs: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientRedeployResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_Deallocate
	fmt.Println("Call operation: VirtualMachineScaleSets_Deallocate")
	virtualMachineScaleSetsClientDeallocateResponsePoller, err := virtualMachineScaleSetsClient.BeginDeallocate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientBeginDeallocateOptions{VMInstanceIDs: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientDeallocateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_Start
	fmt.Println("Call operation: VirtualMachineScaleSets_Start")
	virtualMachineScaleSetsClientStartResponsePoller, err := virtualMachineScaleSetsClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientBeginStartOptions{VMInstanceIDs: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientStartResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_Reimage
	fmt.Println("Call operation: VirtualMachineScaleSets_Reimage")
	virtualMachineScaleSetsClientReimageResponsePoller, err := virtualMachineScaleSetsClient.BeginReimage(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientBeginReimageOptions{VMScaleSetReimageInput: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientReimageResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_Restart
	fmt.Println("Call operation: VirtualMachineScaleSets_Restart")
	virtualMachineScaleSetsClientRestartResponsePoller, err := virtualMachineScaleSetsClient.BeginRestart(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientBeginRestartOptions{VMInstanceIDs: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientRestartResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_ReimageAll
	fmt.Println("Call operation: VirtualMachineScaleSets_ReimageAll")
	virtualMachineScaleSetsClientReimageAllResponsePoller, err := virtualMachineScaleSetsClient.BeginReimageAll(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientBeginReimageAllOptions{VMInstanceIDs: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientReimageAllResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/virtualMachines/{instanceId}
func (testsuite *VirtualMachineScaleSetTestSuite) TestVirtualMachineScaleSetVMs() {
	instanceId := "0"
	var err error
	// From step VirtualMachineScaleSetVMs_List
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_List")
	virtualMachineScaleSetVMsClient, err := armcompute.NewVirtualMachineScaleSetVMsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachineScaleSetVMsClientNewListPager := virtualMachineScaleSetVMsClient.NewListPager(testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetVMsClientListOptions{Filter: nil,
		Select: nil,
		Expand: nil,
	})
	for virtualMachineScaleSetVMsClientNewListPager.More() {
		_, err := virtualMachineScaleSetVMsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineScaleSetVMs_Get
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_Get")
	_, err = virtualMachineScaleSetVMsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, &armcompute.VirtualMachineScaleSetVMsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_GetInstanceView
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_GetInstanceView")
	_, err = virtualMachineScaleSetVMsClient.GetInstanceView(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_Redeploy
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_Redeploy")
	virtualMachineScaleSetVMsClientRedeployResponsePoller, err := virtualMachineScaleSetVMsClient.BeginRedeploy(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMsClientRedeployResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_Start
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_Start")
	virtualMachineScaleSetVMsClientStartResponsePoller, err := virtualMachineScaleSetVMsClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMsClientStartResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_Restart
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_Restart")
	virtualMachineScaleSetVMsClientRestartResponsePoller, err := virtualMachineScaleSetVMsClient.BeginRestart(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMsClientRestartResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_Deallocate
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_Deallocate")
	virtualMachineScaleSetVMsClientDeallocateResponsePoller, err := virtualMachineScaleSetVMsClient.BeginDeallocate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMsClientDeallocateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_Reimage
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_Reimage")
	virtualMachineScaleSetVMsClientReimageResponsePoller, err := virtualMachineScaleSetVMsClient.BeginReimage(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, &armcompute.VirtualMachineScaleSetVMsClientBeginReimageOptions{VMScaleSetVMReimageInput: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMsClientReimageResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_PowerOff
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_PowerOff")
	virtualMachineScaleSetVMsClientPowerOffResponsePoller, err := virtualMachineScaleSetVMsClient.BeginPowerOff(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, &armcompute.VirtualMachineScaleSetVMsClientBeginPowerOffOptions{SkipShutdown: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMsClientPowerOffResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_Delete
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_Delete")
	virtualMachineScaleSetVMsClientDeleteResponsePoller, err := virtualMachineScaleSetVMsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, instanceId, &armcompute.VirtualMachineScaleSetVMsClientBeginDeleteOptions{ForceDeletion: to.Ptr(true)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *VirtualMachineScaleSetTestSuite) Cleanup() {
	var err error
	// From step VirtualMachineScaleSets_PowerOff
	fmt.Println("Call operation: VirtualMachineScaleSets_PowerOff")
	virtualMachineScaleSetsClient, err := armcompute.NewVirtualMachineScaleSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachineScaleSetsClientPowerOffResponsePoller, err := virtualMachineScaleSetsClient.BeginPowerOff(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientBeginPowerOffOptions{SkipShutdown: nil,
		VMInstanceIDs: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientPowerOffResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSets_Delete
	fmt.Println("Call operation: VirtualMachineScaleSets_Delete")
	virtualMachineScaleSetsClientDeleteResponsePoller, err := virtualMachineScaleSetsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, &armcompute.VirtualMachineScaleSetsClientBeginDeleteOptions{ForceDeletion: to.Ptr(true)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

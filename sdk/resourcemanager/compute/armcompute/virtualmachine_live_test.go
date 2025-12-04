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

type VirtualMachineTestSuite struct {
	suite.Suite

	ctx                  context.Context
	cred                 azcore.TokenCredential
	options              *arm.ClientOptions
	adminUsername        string
	networkInterfaceId   string
	networkInterfaceName string
	virtualNetworksName  string
	vmName               string
	adminPassword        string
	location             string
	resourceGroupName    string
	subscriptionId       string
}

func (testsuite *VirtualMachineTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.adminUsername, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmuserna", 14, true)
	testsuite.networkInterfaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmnic", 11, false)
	testsuite.virtualNetworksName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmvnet", 12, false)
	testsuite.vmName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmname", 12, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *VirtualMachineTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestVirtualMachineTestSuite(t *testing.T) {
	suite.Run(t, new(VirtualMachineTestSuite))
}

func (testsuite *VirtualMachineTestSuite) Prepare() {
	var err error
	// From step Create_NetworkInterface
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"networkInterfaceId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/networkInterfaces', parameters('networkInterfaceName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.networkInterfaceName,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.virtualNetworksName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
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
			map[string]any{
				"name":       "[parameters('networkInterfaceName')]",
				"type":       "Microsoft.Network/networkInterfaces",
				"apiVersion": "2021-08-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"ipConfigurations": []any{
						map[string]any{
							"name": "Ipv4config",
							"properties": map[string]any{
								"subnet": map[string]any{
									"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
								},
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_NetworkInterface", &deployment)
	testsuite.Require().NoError(err)
	testsuite.networkInterfaceId = deploymentExtend.Properties.Outputs.(map[string]interface{})["networkInterfaceId"].(map[string]interface{})["value"].(string)

	// From step VirtualMachines_CreateOrUpdate
	fmt.Println("Call operation: VirtualMachines_CreateOrUpdate")
	virtualMachinesClient, err := armcompute.NewVirtualMachinesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachinesClientCreateOrUpdateResponsePoller, err := virtualMachinesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, armcompute.VirtualMachine{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.VirtualMachineProperties{
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(armcompute.VirtualMachineSizeTypesStandardD1V2),
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: to.Ptr(testsuite.networkInterfaceId),
						Properties: &armcompute.NetworkInterfaceReferenceProperties{
							Primary: to.Ptr(true),
						},
					}},
			},
			OSProfile: &armcompute.OSProfile{
				AdminPassword: to.Ptr(testsuite.adminPassword),
				AdminUsername: to.Ptr(testsuite.adminUsername),
				ComputerName:  to.Ptr(testsuite.vmName),
			},
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					Offer:     to.Ptr("WindowsServer"),
					Publisher: to.Ptr("MicrosoftWindowsServer"),
					SKU:       to.Ptr("2016-Datacenter"),
					Version:   to.Ptr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr(testsuite.vmName + "osdisk"),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS),
					},
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/virtualMachines/{vmName}
func (testsuite *VirtualMachineTestSuite) TestVirtualMachines() {
	var err error
	// From step VirtualMachines_ListByLocation
	fmt.Println("Call operation: VirtualMachines_ListByLocation")
	virtualMachinesClient, err := armcompute.NewVirtualMachinesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachinesClientNewListByLocationPager := virtualMachinesClient.NewListByLocationPager(testsuite.location, nil)
	for virtualMachinesClientNewListByLocationPager.More() {
		_, err := virtualMachinesClientNewListByLocationPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachines_InstanceView
	fmt.Println("Call operation: VirtualMachines_InstanceView")
	_, err = virtualMachinesClient.InstanceView(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_ListAll
	fmt.Println("Call operation: VirtualMachines_ListAll")
	virtualMachinesClientNewListAllPager := virtualMachinesClient.NewListAllPager(&armcompute.VirtualMachinesClientListAllOptions{StatusOnly: nil,
		Filter: nil,
	})
	for virtualMachinesClientNewListAllPager.More() {
		_, err := virtualMachinesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachines_List
	fmt.Println("Call operation: VirtualMachines_List")
	virtualMachinesClientNewListPager := virtualMachinesClient.NewListPager(testsuite.resourceGroupName, &armcompute.VirtualMachinesClientListOptions{Filter: nil})
	for virtualMachinesClientNewListPager.More() {
		_, err := virtualMachinesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachines_Get
	fmt.Println("Call operation: VirtualMachines_Get")
	_, err = virtualMachinesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, &armcompute.VirtualMachinesClientGetOptions{Expand: to.Ptr(armcompute.InstanceViewTypesUserData)})
	testsuite.Require().NoError(err)

	// From step VirtualMachines_ListAvailableSizes
	fmt.Println("Call operation: VirtualMachines_ListAvailableSizes")
	virtualMachinesClientNewListAvailableSizesPager := virtualMachinesClient.NewListAvailableSizesPager(testsuite.resourceGroupName, testsuite.vmName, nil)
	for virtualMachinesClientNewListAvailableSizesPager.More() {
		_, err := virtualMachinesClientNewListAvailableSizesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachines_Update
	fmt.Println("Call operation: VirtualMachines_Update")
	virtualMachinesClientUpdateResponsePoller, err := virtualMachinesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, armcompute.VirtualMachineUpdate{
		Tags: map[string]*string{
			"virtaulMachine": to.Ptr("vmupdate"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_InstallPatches
	fmt.Println("Call operation: VirtualMachines_InstallPatches")
	virtualMachinesClientInstallPatchesResponsePoller, err := virtualMachinesClient.BeginInstallPatches(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, armcompute.VirtualMachineInstallPatchesParameters{
		MaximumDuration: to.Ptr("PT4H"),
		RebootSetting:   to.Ptr(armcompute.VMGuestPatchRebootSettingIfRequired),
		WindowsParameters: &armcompute.WindowsParameters{
			ClassificationsToInclude: []*armcompute.VMGuestPatchClassificationWindows{
				to.Ptr(armcompute.VMGuestPatchClassificationWindowsCritical),
				to.Ptr(armcompute.VMGuestPatchClassificationWindowsSecurity)},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientInstallPatchesResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Deallocate
	fmt.Println("Call operation: VirtualMachines_Deallocate")
	virtualMachinesClientDeallocateResponsePoller, err := virtualMachinesClient.BeginDeallocate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, &armcompute.VirtualMachinesClientBeginDeallocateOptions{Hibernate: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientDeallocateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Start
	fmt.Println("Call operation: VirtualMachines_Start")
	virtualMachinesClientStartResponsePoller, err := virtualMachinesClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientStartResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_AssessPatches
	fmt.Println("Call operation: VirtualMachines_AssessPatches")
	virtualMachinesClientAssessPatchesResponsePoller, err := virtualMachinesClient.BeginAssessPatches(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientAssessPatchesResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Restart
	fmt.Println("Call operation: VirtualMachines_Restart")
	virtualMachinesClientRestartResponsePoller, err := virtualMachinesClient.BeginRestart(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientRestartResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Reapply
	fmt.Println("Call operation: VirtualMachines_Reapply")
	virtualMachinesClientReapplyResponsePoller, err := virtualMachinesClient.BeginReapply(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientReapplyResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_RunCommand
	fmt.Println("Call operation: VirtualMachines_RunCommand")
	virtualMachinesClientRunCommandResponsePoller, err := virtualMachinesClient.BeginRunCommand(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, armcompute.RunCommandInput{
		CommandID: to.Ptr("RunPowerShellScript"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientRunCommandResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Redeploy
	fmt.Println("Call operation: VirtualMachines_Redeploy")
	virtualMachinesClientRedeployResponsePoller, err := virtualMachinesClient.BeginRedeploy(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientRedeployResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_PowerOff
	fmt.Println("Call operation: VirtualMachines_PowerOff")
	virtualMachinesClientPowerOffResponsePoller, err := virtualMachinesClient.BeginPowerOff(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, &armcompute.VirtualMachinesClientBeginPowerOffOptions{SkipShutdown: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientPowerOffResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *VirtualMachineTestSuite) Cleanup() {
	var err error
	// From step VirtualMachines_Delete
	fmt.Println("Call operation: VirtualMachines_Delete")
	virtualMachinesClient, err := armcompute.NewVirtualMachinesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachinesClientDeleteResponsePoller, err := virtualMachinesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, &armcompute.VirtualMachinesClientBeginDeleteOptions{ForceDeletion: to.Ptr(true)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

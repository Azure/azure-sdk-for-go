//go:build go1.18
// +build go1.18

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

type RunCommandTestSuite struct {
	suite.Suite

	ctx                      context.Context
	cred                     azcore.TokenCredential
	options                  *arm.ClientOptions
	adminUsername            string
	networkInterfaceId       string
	networkInterfaceName     string
	subnetId                 string
	virtualNetworkSubnetName string
	virtualNetworksName      string
	vmName                   string
	vmScaleSetName           string
	vmRunCommandName         string
	vmssRunCommandName       string
	adminPassword            string
	location                 string
	resourceGroupName        string
	subscriptionId           string
}

func (testsuite *RunCommandTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.adminUsername, _ = recording.GenerateAlphaNumericID(testsuite.T(), "rc", 8, true)
	testsuite.networkInterfaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmnicrc", 13, false)
	testsuite.virtualNetworkSubnetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmssvnetnarc", 18, false)
	testsuite.virtualNetworksName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmvnetrc", 14, false)
	testsuite.vmName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmcommand", 15, false)
	testsuite.vmScaleSetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmscalesetcommand", 23, false)
	testsuite.vmRunCommandName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmruncommand", 18, false)
	testsuite.vmssRunCommandName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmssruncommand", 20, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *RunCommandTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestRunCommandTestSuite(t *testing.T) {
	suite.Run(t, new(RunCommandTestSuite))
}

func (testsuite *RunCommandTestSuite) Prepare() {
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

	// From step Create_NetworkAndSubnet
	template = map[string]any{
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
	deployment = armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_NetworkAndSubnet", &deployment)
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

// Microsoft.Compute/virtualMachines/{vmName}/runCommands/{runCommandName}
func (testsuite *RunCommandTestSuite) TestVirtualMachineRunCommands() {
	var err error
	// From step VirtualMachineRunCommands_CreateOrUpdate
	fmt.Println("Call operation: VirtualMachineRunCommands_CreateOrUpdate")
	virtualMachineRunCommandsClient, err := armcompute.NewVirtualMachineRunCommandsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachineRunCommandsClientCreateOrUpdateResponsePoller, err := virtualMachineRunCommandsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, testsuite.vmRunCommandName, armcompute.VirtualMachineRunCommand{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.VirtualMachineRunCommandProperties{
			AsyncExecution: to.Ptr(false),
			Parameters: []*armcompute.RunCommandInputParameter{
				{
					Name:  to.Ptr("param1"),
					Value: to.Ptr("value1"),
				},
				{
					Name:  to.Ptr("param2"),
					Value: to.Ptr("value2"),
				}},
			RunAsPassword: to.Ptr("<runAsPassword>"),
			RunAsUser:     to.Ptr("user1"),
			Source: &armcompute.VirtualMachineRunCommandScriptSource{
				Script: to.Ptr("Write-Host Hello World!"),
			},
			TimeoutInSeconds: to.Ptr[int32](3600),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineRunCommandsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineRunCommands_List
	fmt.Println("Call operation: VirtualMachineRunCommands_List")
	virtualMachineRunCommandsClientNewListPager := virtualMachineRunCommandsClient.NewListPager(testsuite.location, nil)
	for virtualMachineRunCommandsClientNewListPager.More() {
		_, err := virtualMachineRunCommandsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineRunCommands_Get
	fmt.Println("Call operation: VirtualMachineRunCommands_Get")
	_, err = virtualMachineRunCommandsClient.Get(testsuite.ctx, testsuite.location, "RunPowerShellScript", nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachineRunCommands_ListByVirtualMachine
	fmt.Println("Call operation: VirtualMachineRunCommands_ListByVirtualMachine")
	virtualMachineRunCommandsClientNewListByVirtualMachinePager := virtualMachineRunCommandsClient.NewListByVirtualMachinePager(testsuite.resourceGroupName, testsuite.vmName, &armcompute.VirtualMachineRunCommandsClientListByVirtualMachineOptions{Expand: nil})
	for virtualMachineRunCommandsClientNewListByVirtualMachinePager.More() {
		_, err := virtualMachineRunCommandsClientNewListByVirtualMachinePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineRunCommands_GetByVirtualMachine
	fmt.Println("Call operation: VirtualMachineRunCommands_GetByVirtualMachine")
	_, err = virtualMachineRunCommandsClient.GetByVirtualMachine(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, testsuite.vmRunCommandName, &armcompute.VirtualMachineRunCommandsClientGetByVirtualMachineOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step VirtualMachineRunCommands_Update
	fmt.Println("Call operation: VirtualMachineRunCommands_Update")
	virtualMachineRunCommandsClientUpdateResponsePoller, err := virtualMachineRunCommandsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, testsuite.vmRunCommandName, armcompute.VirtualMachineRunCommandUpdate{
		Properties: &armcompute.VirtualMachineRunCommandProperties{
			Source: &armcompute.VirtualMachineRunCommandScriptSource{
				Script: to.Ptr("Write-Host Script Source Updated!"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineRunCommandsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_RunCommand
	fmt.Println("Call operation: VirtualMachines_RunCommand")
	virtualMachinesClient, err := armcompute.NewVirtualMachinesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachinesClientRunCommandResponsePoller, err := virtualMachinesClient.BeginRunCommand(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, armcompute.RunCommandInput{
		CommandID: to.Ptr("RunPowerShellScript"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientRunCommandResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineRunCommands_Delete
	fmt.Println("Call operation: VirtualMachineRunCommands_Delete")
	virtualMachineRunCommandsClientDeleteResponsePoller, err := virtualMachineRunCommandsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmName, testsuite.vmRunCommandName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineRunCommandsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/virtualMachines/{instanceId}/runCommands/{runCommandName}
func (testsuite *RunCommandTestSuite) TestVirtualMachineScaleSetVmRunCommands() {
	var err error
	// From step VirtualMachineScaleSetVMRunCommands_CreateOrUpdate
	fmt.Println("Call operation: VirtualMachineScaleSetVMRunCommands_CreateOrUpdate")
	virtualMachineScaleSetVMRunCommandsClient, err := armcompute.NewVirtualMachineScaleSetVMRunCommandsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachineScaleSetVMRunCommandsClientCreateOrUpdateResponsePoller, err := virtualMachineScaleSetVMRunCommandsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, "0", testsuite.vmssRunCommandName, armcompute.VirtualMachineRunCommand{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.VirtualMachineRunCommandProperties{
			AsyncExecution: to.Ptr(false),
			Parameters: []*armcompute.RunCommandInputParameter{
				{
					Name:  to.Ptr("param1"),
					Value: to.Ptr("value1"),
				},
				{
					Name:  to.Ptr("param2"),
					Value: to.Ptr("value2"),
				}},
			RunAsPassword: to.Ptr("<runAsPassword>"),
			RunAsUser:     to.Ptr("user1"),
			Source: &armcompute.VirtualMachineRunCommandScriptSource{
				Script: to.Ptr("Write-Host Hello World!"),
			},
			TimeoutInSeconds: to.Ptr[int32](3600),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMRunCommandsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMRunCommands_List
	fmt.Println("Call operation: VirtualMachineScaleSetVMRunCommands_List")
	virtualMachineScaleSetVMRunCommandsClientNewListPager := virtualMachineScaleSetVMRunCommandsClient.NewListPager(testsuite.resourceGroupName, testsuite.vmScaleSetName, "0", &armcompute.VirtualMachineScaleSetVMRunCommandsClientListOptions{Expand: nil})
	for virtualMachineScaleSetVMRunCommandsClientNewListPager.More() {
		_, err := virtualMachineScaleSetVMRunCommandsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualMachineScaleSetVMRunCommands_Get
	fmt.Println("Call operation: VirtualMachineScaleSetVMRunCommands_Get")
	_, err = virtualMachineScaleSetVMRunCommandsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, "0", testsuite.vmssRunCommandName, &armcompute.VirtualMachineScaleSetVMRunCommandsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMRunCommands_Update
	fmt.Println("Call operation: VirtualMachineScaleSetVMRunCommands_Update")
	virtualMachineScaleSetVMRunCommandsClientUpdateResponsePoller, err := virtualMachineScaleSetVMRunCommandsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, "0", testsuite.vmssRunCommandName, armcompute.VirtualMachineRunCommandUpdate{
		Properties: &armcompute.VirtualMachineRunCommandProperties{
			Source: &armcompute.VirtualMachineRunCommandScriptSource{
				Script: to.Ptr("Write-Host Script Source Updated!"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMRunCommandsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMs_RunCommand
	fmt.Println("Call operation: VirtualMachineScaleSetVMs_RunCommand")
	virtualMachineScaleSetVMsClient, err := armcompute.NewVirtualMachineScaleSetVMsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachineScaleSetVMsClientRunCommandResponsePoller, err := virtualMachineScaleSetVMsClient.BeginRunCommand(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, "0", armcompute.RunCommandInput{
		CommandID: to.Ptr("RunPowerShellScript"),
		Script: []*string{
			to.Ptr("Write-Host Hello World!")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMsClientRunCommandResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachineScaleSetVMRunCommands_Delete
	fmt.Println("Call operation: VirtualMachineScaleSetVMRunCommands_Delete")
	virtualMachineScaleSetVMRunCommandsClientDeleteResponsePoller, err := virtualMachineScaleSetVMRunCommandsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.vmScaleSetName, "0", testsuite.vmssRunCommandName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachineScaleSetVMRunCommandsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

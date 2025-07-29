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

type RestorePointTestSuite struct {
	suite.Suite

	ctx                        context.Context
	cred                       azcore.TokenCredential
	options                    *arm.ClientOptions
	adminUsername              string
	networkInterfaceId         string
	networkInterfaceName       string
	restorePointCollectionName string
	restorePointName           string
	virtaulMachineId           string
	virtualNetworksName        string
	vmName                     string
	adminPassword              string
	location                   string
	resourceGroupName          string
	subscriptionId             string
}

func (testsuite *RestorePointTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.adminUsername, _ = recording.GenerateAlphaNumericID(testsuite.T(), "rp", 8, true)
	testsuite.networkInterfaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmnicrp", 13, false)
	testsuite.restorePointCollectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "restorepoi", 16, false)
	testsuite.restorePointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "restorepoi", 16, false)
	testsuite.virtualNetworksName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmvnetrp", 14, false)
	testsuite.vmName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vmnamerp", 14, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *RestorePointTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestRestorePointTestSuite(t *testing.T) {
	suite.Run(t, new(RestorePointTestSuite))
}

func (testsuite *RestorePointTestSuite) Prepare() {
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
	var virtualMachinesClientCreateOrUpdateResponse *armcompute.VirtualMachinesClientCreateOrUpdateResponse
	virtualMachinesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.virtaulMachineId = *virtualMachinesClientCreateOrUpdateResponse.ID

	// From step RestorePointCollections_CreateOrUpdate
	fmt.Println("Call operation: RestorePointCollections_CreateOrUpdate")
	restorePointCollectionsClient, err := armcompute.NewRestorePointCollectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = restorePointCollectionsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.restorePointCollectionName, armcompute.RestorePointCollection{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"myTag1": to.Ptr("tagValue1"),
		},
		Properties: &armcompute.RestorePointCollectionProperties{
			Source: &armcompute.RestorePointCollectionSourceProperties{
				ID: to.Ptr(testsuite.virtaulMachineId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/restorePointCollections
func (testsuite *RestorePointTestSuite) TestRestorePointCollections() {
	var err error
	// From step RestorePointCollections_ListAll
	fmt.Println("Call operation: RestorePointCollections_ListAll")
	restorePointCollectionsClient, err := armcompute.NewRestorePointCollectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	restorePointCollectionsClientNewListAllPager := restorePointCollectionsClient.NewListAllPager(nil)
	for restorePointCollectionsClientNewListAllPager.More() {
		_, err := restorePointCollectionsClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RestorePointCollections_List
	fmt.Println("Call operation: RestorePointCollections_List")
	restorePointCollectionsClientNewListPager := restorePointCollectionsClient.NewListPager(testsuite.resourceGroupName, nil)
	for restorePointCollectionsClientNewListPager.More() {
		_, err := restorePointCollectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RestorePointCollections_Get
	fmt.Println("Call operation: RestorePointCollections_Get")
	_, err = restorePointCollectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.restorePointCollectionName, &armcompute.RestorePointCollectionsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/restorePointCollections/restorePoints
func (testsuite *RestorePointTestSuite) TestRestorePoints() {
	var err error
	// From step RestorePoints_Create
	fmt.Println("Call operation: RestorePoints_Create")
	restorePointsClient, err := armcompute.NewRestorePointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	restorePointsClientCreateResponsePoller, err := restorePointsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.restorePointCollectionName, testsuite.restorePointName, armcompute.RestorePoint{}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, restorePointsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step RestorePoints_Get
	fmt.Println("Call operation: RestorePoints_Get")
	_, err = restorePointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.restorePointCollectionName, testsuite.restorePointName, &armcompute.RestorePointsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step RestorePoints_Delete
	fmt.Println("Call operation: RestorePoints_Delete")
	restorePointsClientDeleteResponsePoller, err := restorePointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.restorePointCollectionName, testsuite.restorePointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, restorePointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *RestorePointTestSuite) Cleanup() {
	var err error
	// From step RestorePointCollections_Update
	restorePointCollectionsClient, err := armcompute.NewRestorePointCollectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	fmt.Println("Call operation: RestorePointCollections_Update")
	_, err = restorePointCollectionsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.restorePointCollectionName, armcompute.RestorePointCollectionUpdate{}, nil)
	testsuite.Require().NoError(err)
	// From step RestorePointCollections_Delete
	fmt.Println("Call operation: RestorePointCollections_Delete")
	restorePointCollectionsClientDeleteResponsePoller, err := restorePointCollectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.restorePointCollectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, restorePointCollectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache/v4"
	"github.com/stretchr/testify/suite"
)

type AmlfilesystemTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	amlFilesystemName string
	armEndpoint       string
	subnetId          string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *AmlfilesystemTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.amlFilesystemName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "amlfiles", 14, false)
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AmlfilesystemTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestAmlfilesystemTestSuite(t *testing.T) {
	suite.Run(t, new(AmlfilesystemTestSuite))
}

func (testsuite *AmlfilesystemTestSuite) Prepare() {
	var err error
	// From step Create_VirtualNetwork
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"subnetId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "storagecachevnet",
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
				"tags": map[string]any{},
			},
		},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_VirtualNetwork", &deployment)
	testsuite.Require().NoError(err)
	testsuite.subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)
}

// Microsoft.StorageCache/amlFilesystems/{amlFilesystemName}
func (testsuite *AmlfilesystemTestSuite) TestAmlFilesystems() {
	var err error
	// From step amlFilesystems_CreateOrUpdate
	fmt.Println("Call operation: amlFilesystems_CreateOrUpdate")
	amlFilesystemsClient, err := armstoragecache.NewAmlFilesystemsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	amlFilesystemsClientCreateOrUpdateResponsePoller, err := amlFilesystemsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.amlFilesystemName, armstoragecache.AmlFilesystem{
		Location: to.Ptr(testsuite.location),
		Properties: &armstoragecache.AmlFilesystemProperties{
			FilesystemSubnet: to.Ptr(testsuite.subnetId),
			MaintenanceWindow: &armstoragecache.AmlFilesystemPropertiesMaintenanceWindow{
				DayOfWeek:    to.Ptr(armstoragecache.MaintenanceDayOfWeekTypeMonday),
				TimeOfDayUTC: to.Ptr("23:25"),
			},
			StorageCapacityTiB: to.Ptr[float32](16),
		},
		SKU: &armstoragecache.SKUName{
			Name: to.Ptr("AMLFS-Durable-Premium-125"),
		},
		Zones: []*string{
			to.Ptr("1")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, amlFilesystemsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step amlFilesystems_List
	fmt.Println("Call operation: amlFilesystems_List")
	amlFilesystemsClientNewListPager := amlFilesystemsClient.NewListPager(nil)
	for amlFilesystemsClientNewListPager.More() {
		_, err := amlFilesystemsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step amlFilesystems_Get
	fmt.Println("Call operation: amlFilesystems_Get")
	_, err = amlFilesystemsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.amlFilesystemName, nil)
	testsuite.Require().NoError(err)

	// From step amlFilesystems_ListByResourceGroup
	fmt.Println("Call operation: amlFilesystems_ListByResourceGroup")
	amlFilesystemsClientNewListByResourceGroupPager := amlFilesystemsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for amlFilesystemsClientNewListByResourceGroupPager.More() {
		_, err := amlFilesystemsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step amlFilesystems_Update
	fmt.Println("Call operation: amlFilesystems_Update")
	amlFilesystemsClientUpdateResponsePoller, err := amlFilesystemsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.amlFilesystemName, armstoragecache.AmlFilesystemUpdate{
		Tags: map[string]*string{
			"Dept": to.Ptr("ContosoAds"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, amlFilesystemsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step checkAmlFSSubnets
	fmt.Println("Call operation: checkAmlFSSubnets")
	managementClient, err := armstoragecache.NewManagementClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managementClient.CheckAmlFSSubnets(testsuite.ctx, &armstoragecache.ManagementClientCheckAmlFSSubnetsOptions{
		AmlFilesystemSubnetInfo: &armstoragecache.AmlFilesystemSubnetInfo{
			Location:           to.Ptr(testsuite.location),
			FilesystemSubnet:   to.Ptr(testsuite.subnetId),
			StorageCapacityTiB: to.Ptr[float32](16),
			SKU: &armstoragecache.SKUName{
				Name: to.Ptr("AMLFS-Durable-Premium-125"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step getRequiredAmlFSSubnetsSize
	fmt.Println("Call operation: getRequiredAmlFSSubnetsSize")
	_, err = managementClient.GetRequiredAmlFSSubnetsSize(testsuite.ctx, &armstoragecache.ManagementClientGetRequiredAmlFSSubnetsSizeOptions{
		RequiredAMLFilesystemSubnetsSizeInfo: &armstoragecache.RequiredAmlFilesystemSubnetsSizeInfo{
			StorageCapacityTiB: to.Ptr[float32](16),
			SKU: &armstoragecache.SKUName{
				Name: to.Ptr("AMLFS-Durable-Premium-125"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step amlFilesystems_Delete
	fmt.Println("Call operation: amlFilesystems_Delete")
	amlFilesystemsClientDeleteResponsePoller, err := amlFilesystemsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.amlFilesystemName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, amlFilesystemsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

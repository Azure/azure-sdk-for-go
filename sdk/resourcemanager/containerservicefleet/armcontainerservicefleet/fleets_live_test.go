// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcontainerservicefleet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservicefleet/armcontainerservicefleet/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type FleetsTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	armEndpoint         string
	fleetMemberName     string
	fleetName           string
	managedClustersName string
	updateRunName       string
	updateStrategyName  string
	azureClientId       string
	azureClientSecret   string
	location            string
	resourceGroupName   string
	subscriptionId      string
}

func (testsuite *FleetsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.fleetMemberName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "fleetmem", 14, true)
	testsuite.fleetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "fleetnam", 14, true)
	testsuite.managedClustersName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "askcluster", 16, true)
	testsuite.updateRunName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "updateru", 14, true)
	testsuite.updateStrategyName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "updatest", 14, true)
	testsuite.azureClientId = recording.GetEnvVariable("AZURE_CLIENT_ID", "")
	testsuite.azureClientSecret = recording.GetEnvVariable("AZURE_CLIENT_SECRET", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *FleetsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestFleetsTestSuite(t *testing.T) {
	suite.Run(t, new(FleetsTestSuite))
}

func (testsuite *FleetsTestSuite) Prepare() {
	var err error
	// From step Fleets_CreateOrUpdate
	fmt.Println("Call operation: Fleets_CreateOrUpdate")
	fleetsClient, err := armcontainerservicefleet.NewFleetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	fleetsClientCreateOrUpdateResponsePoller, err := fleetsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, armcontainerservicefleet.Fleet{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"archv2": to.Ptr(""),
			"tier":   to.Ptr("production"),
		},
		Properties: &armcontainerservicefleet.FleetProperties{},
	}, &armcontainerservicefleet.FleetsClientBeginCreateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, fleetsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerService/fleets/{fleetName}
func (testsuite *FleetsTestSuite) TestFleets() {
	var err error
	// From step Fleets_ListBySubscription
	fmt.Println("Call operation: Fleets_ListBySubscription")
	fleetsClient, err := armcontainerservicefleet.NewFleetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	fleetsClientNewListBySubscriptionPager := fleetsClient.NewListBySubscriptionPager(nil)
	for fleetsClientNewListBySubscriptionPager.More() {
		_, err := fleetsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Fleets_Get
	fmt.Println("Call operation: Fleets_Get")
	_, err = fleetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, nil)
	testsuite.Require().NoError(err)

	// From step Fleets_ListByResourceGroup
	fmt.Println("Call operation: Fleets_ListByResourceGroup")
	fleetsClientNewListByResourceGroupPager := fleetsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for fleetsClientNewListByResourceGroupPager.More() {
		_, err := fleetsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Fleets_Update
	fmt.Println("Call operation: Fleets_Update")
	fleetsClientUpdateResponsePoller, err := fleetsClient.BeginUpdateAsync(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, armcontainerservicefleet.FleetPatch{
		Tags: map[string]*string{
			"env":  to.Ptr("prod"),
			"tier": to.Ptr("secure"),
		},
	}, &armcontainerservicefleet.FleetsClientBeginUpdateAsyncOptions{IfMatch: to.Ptr("dfjkwelr7384")})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, fleetsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerService/fleets/{fleetName}/members/{fleetMemberName}
func (testsuite *FleetsTestSuite) TestFleetMembers() {
	var managedClusterId string
	var err error
	// From step Create_ManageCluster
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"managedClusterId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.ContainerService/managedClusters', parameters('managedClustersName'))]",
			},
		},
		"parameters": map[string]any{
			"azureClientId": map[string]any{
				"type":         "securestring",
				"defaultValue": testsuite.azureClientId,
			},
			"azureClientSecret": map[string]any{
				"type":         "securestring",
				"defaultValue": testsuite.azureClientSecret,
			},
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"managedClustersName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.managedClustersName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('managedClustersName')]",
				"type":       "Microsoft.ContainerService/managedClusters",
				"apiVersion": "2023-01-02-preview",
				"identity": map[string]any{
					"type": "SystemAssigned",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"agentPoolProfiles": []any{
						map[string]any{
							"name":                "nodepool1",
							"type":                "VirtualMachineScaleSets",
							"count":               float64(3),
							"enableAutoScaling":   true,
							"enableFIPS":          false,
							"enableNodePublicIP":  true,
							"kubeletDiskType":     "OS",
							"maxCount":            float64(10),
							"maxPods":             float64(110),
							"minCount":            float64(1),
							"mode":                "System",
							"orchestratorVersion": "1.26.6",
							"osDiskSizeGB":        float64(128),
							"osDiskType":          "Managed",
							"osSKU":               "Ubuntu",
							"osType":              "Linux",
							"powerState": map[string]any{
								"code": "Running",
							},
							"upgradeSettings": map[string]any{},
							"vmSize":          "Standard_DS2_v2",
						},
					},
					"dnsPrefix":         "dnsprefix1",
					"kubernetesVersion": "1.26.6",
					"servicePrincipalProfile": map[string]any{
						"clientId": "[parameters('azureClientId')]",
						"secret":   "[parameters('azureClientSecret')]",
					},
				},
				"sku": map[string]any{
					"name": "Basic",
					"tier": "Free",
				},
			},
			map[string]any{
				"name":       "[concat(parameters('managedClustersName'), '/nodepool1')]",
				"type":       "Microsoft.ContainerService/managedClusters/agentPools",
				"apiVersion": "2023-01-02-preview",
				"dependsOn": []any{
					"[resourceId('Microsoft.ContainerService/managedClusters', parameters('managedClustersName'))]",
				},
				"properties": map[string]any{
					"type":                "VirtualMachineScaleSets",
					"count":               float64(3),
					"enableAutoScaling":   true,
					"enableFIPS":          false,
					"enableNodePublicIP":  true,
					"kubeletDiskType":     "OS",
					"maxCount":            float64(10),
					"maxPods":             float64(110),
					"minCount":            float64(1),
					"mode":                "System",
					"orchestratorVersion": "1.26.6",
					"osDiskSizeGB":        float64(128),
					"osDiskType":          "Managed",
					"osSKU":               "Ubuntu",
					"osType":              "Linux",
					"powerState": map[string]any{
						"code": "Running",
					},
					"upgradeSettings": map[string]any{},
					"vmSize":          "Standard_DS2_v2",
				},
			},
		},
		"variables": map[string]any{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_ManageCluster", &deployment)
	testsuite.Require().NoError(err)
	managedClusterId = deploymentExtend.Properties.Outputs.(map[string]interface{})["managedClusterId"].(map[string]interface{})["value"].(string)

	// From step FleetMembers_Create
	fmt.Println("Call operation: FleetMembers_Create")
	fleetMembersClient, err := armcontainerservicefleet.NewFleetMembersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	fleetMembersClientCreateResponsePoller, err := fleetMembersClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, testsuite.fleetMemberName, armcontainerservicefleet.FleetMember{
		Properties: &armcontainerservicefleet.FleetMemberProperties{
			ClusterResourceID: to.Ptr(managedClusterId),
		},
	}, &armcontainerservicefleet.FleetMembersClientBeginCreateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, fleetMembersClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step FleetMembers_ListByFleet
	fmt.Println("Call operation: FleetMembers_ListByFleet")
	fleetMembersClientNewListByFleetPager := fleetMembersClient.NewListByFleetPager(testsuite.resourceGroupName, testsuite.fleetName, nil)
	for fleetMembersClientNewListByFleetPager.More() {
		_, err := fleetMembersClientNewListByFleetPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FleetMembers_Get
	fmt.Println("Call operation: FleetMembers_Get")
	_, err = fleetMembersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, testsuite.fleetMemberName, nil)
	testsuite.Require().NoError(err)

	// From step FleetMembers_Update
	fmt.Println("Call operation: FleetMembers_Update")
	fleetMembersClientUpdateResponsePoller, err := fleetMembersClient.BeginUpdateAsync(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, testsuite.fleetMemberName, armcontainerservicefleet.FleetMemberUpdate{
		Properties: &armcontainerservicefleet.FleetMemberUpdateProperties{
			Group: to.Ptr("staging"),
		},
	}, &armcontainerservicefleet.FleetMembersClientBeginUpdateAsyncOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, fleetMembersClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerService/fleets/{fleetName}/updateStrategies/{updateStrategyName}
func (testsuite *FleetsTestSuite) TestFleetUpdateStrategies() {
	var err error
	// From step FleetUpdateStrategies_CreateOrUpdate
	fmt.Println("Call operation: FleetUpdateStrategies_CreateOrUpdate")
	fleetUpdateStrategiesClient, err := armcontainerservicefleet.NewFleetUpdateStrategiesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	fleetUpdateStrategiesClientCreateOrUpdateResponsePoller, err := fleetUpdateStrategiesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, testsuite.updateStrategyName, armcontainerservicefleet.FleetUpdateStrategy{
		Properties: &armcontainerservicefleet.FleetUpdateStrategyProperties{
			Strategy: &armcontainerservicefleet.UpdateRunStrategy{
				Stages: []*armcontainerservicefleet.UpdateStage{
					{
						Name:                    to.Ptr("stage1"),
						AfterStageWaitInSeconds: to.Ptr[int32](3600),
						Groups: []*armcontainerservicefleet.UpdateGroup{
							{
								Name: to.Ptr("group-a"),
							}},
					}},
			},
		},
	}, &armcontainerservicefleet.FleetUpdateStrategiesClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, fleetUpdateStrategiesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step FleetUpdateStrategies_ListByFleet
	fmt.Println("Call operation: FleetUpdateStrategies_ListByFleet")
	fleetUpdateStrategiesClientNewListByFleetPager := fleetUpdateStrategiesClient.NewListByFleetPager(testsuite.resourceGroupName, testsuite.fleetName, nil)
	for fleetUpdateStrategiesClientNewListByFleetPager.More() {
		_, err := fleetUpdateStrategiesClientNewListByFleetPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FleetUpdateStrategies_Get
	fmt.Println("Call operation: FleetUpdateStrategies_Get")
	_, err = fleetUpdateStrategiesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, testsuite.updateStrategyName, nil)
	testsuite.Require().NoError(err)

	// From step FleetUpdateStrategies_Delete
	fmt.Println("Call operation: FleetUpdateStrategies_Delete")
	fleetUpdateStrategiesClientDeleteResponsePoller, err := fleetUpdateStrategiesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.fleetName, testsuite.updateStrategyName, &armcontainerservicefleet.FleetUpdateStrategiesClientBeginDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, fleetUpdateStrategiesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

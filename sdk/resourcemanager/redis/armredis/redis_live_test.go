// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armredis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redis/armredis/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	name                string
	privateEndpointName string
	redisId             string
	location            string
	resourceGroupName   string
	subnetId            string
	subscriptionId      string
}

func (testsuite *RedisTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.name, _ = recording.GenerateAlphaNumericID(testsuite.T(), "redisna", 13, false)
	testsuite.privateEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "redisprivateendpoint", 26, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subnetId = recording.GetEnvVariable("SUBNET_ID", "")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *RedisTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

func (testsuite *RedisTestSuite) Prepare() {
	var err error
	// From step NetworkSubnet_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"subnetId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'redissubnet')]",
			},
		},
		"parameters": map[string]interface{}{
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworksName": map[string]interface{}{
				"type":         "string",
				"defaultValue": "redisvnet",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"addressSpace": map[string]interface{}{
						"addressPrefixes": []interface{}{
							"10.0.0.0/16",
						},
					},
					"subnets": []interface{}{
						map[string]interface{}{
							"name": "redissubnet",
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
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "NetworkSubnet_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step Redis_Create
	fmt.Println("Call operation: Redis_Create")
	client, err := armredis.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientCreateResponsePoller, err := client.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, armredis.CreateParameters{
		Location: to.Ptr(testsuite.location),
		Properties: &armredis.CreateProperties{
			EnableNonSSLPort:  to.Ptr(true),
			MinimumTLSVersion: to.Ptr(armredis.TLSVersionOne2),
			RedisConfiguration: &armredis.CommonPropertiesRedisConfiguration{
				MaxmemoryPolicy: to.Ptr("allkeys-lru"),
			},
			ReplicasPerPrimary: to.Ptr[int32](2),
			ShardCount:         to.Ptr[int32](2),
			SKU: &armredis.SKU{
				Name:     to.Ptr(armredis.SKUNamePremium),
				Capacity: to.Ptr[int32](1),
				Family:   to.Ptr(armredis.SKUFamilyP),
			},
		},
		Zones: []*string{
			to.Ptr("1")},
	}, nil)
	testsuite.Require().NoError(err)
	var clientCreateResponse *armredis.ClientCreateResponse
	clientCreateResponse, err = testutil.PollForTest(testsuite.ctx, clientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.redisId = *clientCreateResponse.ID
}

// Microsoft.Cache/redis
func (testsuite *RedisTestSuite) TestRedis() {
	var err error
	// From step Redis_CheckNameAvailability
	fmt.Println("Call operation: Redis_CheckNameAvailability")
	client, err := armredis.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.CheckNameAvailability(testsuite.ctx, armredis.CheckNameAvailabilityParameters{
		Name: to.Ptr("cacheName"),
		Type: to.Ptr("Microsoft.Cache/Redis"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Redis_ListBySubscription
	fmt.Println("Call operation: Redis_ListBySubscription")
	clientNewListBySubscriptionPager := client.NewListBySubscriptionPager(nil)
	for clientNewListBySubscriptionPager.More() {
		_, err := clientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Redis_ListByResourceGroup
	fmt.Println("Call operation: Redis_ListByResourceGroup")
	clientNewListByResourceGroupPager := client.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for clientNewListByResourceGroupPager.More() {
		_, err := clientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Redis_Get
	fmt.Println("Call operation: Redis_Get")
	_, err = client.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, nil)
	testsuite.Require().NoError(err)

	// From step Redis_Update
	fmt.Println("Call operation: Redis_Update")
	clientUpdateResponsePoller, err := client.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, armredis.UpdateParameters{
		Properties: &armredis.UpdateProperties{
			EnableNonSSLPort:   to.Ptr(true),
			ReplicasPerPrimary: to.Ptr[int32](2),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Redis_RegenerateKey
	fmt.Println("Call operation: Redis_RegenerateKey")
	_, err = client.RegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, armredis.RegenerateKeyParameters{
		KeyType: to.Ptr(armredis.RedisKeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Redis_ListKeys
	fmt.Println("Call operation: Redis_ListKeys")
	_, err = client.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, nil)
	testsuite.Require().NoError(err)

	// From step Redis_ForceReboot
	fmt.Println("Call operation: Redis_ForceReboot")
	_, err = client.ForceReboot(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, armredis.RebootParameters{
		Ports: []*int32{
			to.Ptr[int32](13000),
			to.Ptr[int32](15001)},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Cache/redis/firewallRules
func (testsuite *RedisTestSuite) TestFirewallRule() {
	cacheName := testsuite.name
	var err error
	// From step FirewallRules_CreateOrUpdate
	fmt.Println("Call operation: FirewallRules_CreateOrUpdate")
	firewallRulesClient, err := armredis.NewFirewallRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = firewallRulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, cacheName, "rule1", armredis.FirewallRule{
		Properties: &armredis.FirewallRuleProperties{
			EndIP:   to.Ptr("10.0.1.4"),
			StartIP: to.Ptr("10.0.1.1"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step FirewallRules_List
	fmt.Println("Call operation: FirewallRules_List")
	firewallRulesClientNewListPager := firewallRulesClient.NewListPager(testsuite.resourceGroupName, cacheName, nil)
	for firewallRulesClientNewListPager.More() {
		_, err := firewallRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FirewallRules_Get
	fmt.Println("Call operation: FirewallRules_Get")
	_, err = firewallRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, cacheName, "rule1", nil)
	testsuite.Require().NoError(err)

	// From step FirewallRules_Delete
	fmt.Println("Call operation: FirewallRules_Delete")
	_, err = firewallRulesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, cacheName, "rule1", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Cache/redis/patchSchedules
func (testsuite *RedisTestSuite) TestPatchSchedule() {
	cacheName := testsuite.name
	var err error
	// From step PatchSchedules_CreateOrUpdate
	fmt.Println("Call operation: PatchSchedules_CreateOrUpdate")
	patchSchedulesClient, err := armredis.NewPatchSchedulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = patchSchedulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, armredis.DefaultNameDefault, armredis.PatchSchedule{
		Properties: &armredis.ScheduleEntries{
			ScheduleEntries: []*armredis.ScheduleEntry{
				{
					DayOfWeek:         to.Ptr(armredis.DayOfWeekMonday),
					MaintenanceWindow: to.Ptr("PT5H"),
					StartHourUTC:      to.Ptr[int32](12),
				},
				{
					DayOfWeek:    to.Ptr(armredis.DayOfWeekTuesday),
					StartHourUTC: to.Ptr[int32](12),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PatchSchedules_ListByRedisResource
	fmt.Println("Call operation: PatchSchedules_ListByRedisResource")
	patchSchedulesClientNewListByRedisResourcePager := patchSchedulesClient.NewListByRedisResourcePager(testsuite.resourceGroupName, cacheName, nil)
	for patchSchedulesClientNewListByRedisResourcePager.More() {
		_, err := patchSchedulesClientNewListByRedisResourcePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PatchSchedules_Get
	fmt.Println("Call operation: PatchSchedules_Get")
	_, err = patchSchedulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, armredis.DefaultNameDefault, nil)
	testsuite.Require().NoError(err)

	// From step PatchSchedules_Delete
	fmt.Println("Call operation: PatchSchedules_Delete")
	_, err = patchSchedulesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, armredis.DefaultNameDefault, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Cache/operations
func (testsuite *RedisTestSuite) TestOperation() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armredis.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Cache/redis/privateEndpointConnections
func (testsuite *RedisTestSuite) TestPrivateEndpointConnections() {
	cacheName := testsuite.name
	var privateEndpointConnectionName string
	var err error
	// From step PrivateEndpoint_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]interface{}{
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"privateEndpointName": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.privateEndpointName,
			},
			"redisId": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.redisId,
			},
			"subnetId": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.subnetId,
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[concat(parameters('privateEndpointName'), '-nic')]",
				"type":       "Microsoft.Network/networkInterfaces",
				"apiVersion": "2020-11-01",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"dnsSettings": map[string]interface{}{
						"dnsServers": []interface{}{},
					},
					"enableIPForwarding": false,
					"ipConfigurations": []interface{}{
						map[string]interface{}{
							"name": "privateEndpointIpConfig",
							"properties": map[string]interface{}{
								"primary":                   true,
								"privateIPAddress":          "10.0.0.4",
								"privateIPAddressVersion":   "IPv4",
								"privateIPAllocationMethod": "Dynamic",
								"subnet": map[string]interface{}{
									"id": "[parameters('subnetId')]",
								},
							},
						},
					},
				},
			},
			map[string]interface{}{
				"name":       "[parameters('privateEndpointName')]",
				"type":       "Microsoft.Network/privateEndpoints",
				"apiVersion": "2020-11-01",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"customDnsConfigs":                    []interface{}{},
					"manualPrivateLinkServiceConnections": []interface{}{},
					"privateLinkServiceConnections": []interface{}{
						map[string]interface{}{
							"name": "[parameters('privateEndpointName')]",
							"properties": map[string]interface{}{
								"groupIds": []interface{}{
									"redisCache",
								},
								"privateLinkServiceConnectionState": map[string]interface{}{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('redisId')]",
							},
						},
					},
					"subnet": map[string]interface{}{
						"id": "[parameters('subnetId')]",
					},
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "PrivateEndpoint_Create", &deployment)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_List
	fmt.Println("Call operation: PrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armredis.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, cacheName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Put
	fmt.Println("Call operation: PrivateEndpointConnections_Put")
	privateEndpointConnectionsClientPutResponsePoller, err := privateEndpointConnectionsClient.BeginPut(testsuite.ctx, testsuite.resourceGroupName, cacheName, privateEndpointConnectionName, armredis.PrivateEndpointConnection{
		Properties: &armredis.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armredis.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Auto-Approved"),
				Status:      to.Ptr(armredis.PrivateEndpointServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientPutResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, cacheName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_ListByRedisCache
	fmt.Println("Call operation: PrivateLinkResources_ListByRedisCache")
	privateLinkResourcesClient, err := armredis.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListByRedisCachePager := privateLinkResourcesClient.NewListByRedisCachePager(testsuite.resourceGroupName, cacheName, nil)
	for privateLinkResourcesClientNewListByRedisCachePager.More() {
		_, err := privateLinkResourcesClientNewListByRedisCachePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	_, err = privateEndpointConnectionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, cacheName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *RedisTestSuite) Cleanup() {
	var err error
	// From step Redis_Delete
	fmt.Println("Call operation: Redis_Delete")
	client, err := armredis.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientDeleteResponsePoller, err := client.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.name, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

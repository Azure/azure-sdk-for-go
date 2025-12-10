// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armredisenterprise_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type RedisenterpriseTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	clusterName       string
	databaseName      string
	redisEnterpriseId string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *RedisenterpriseTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.clusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clustern", 14, false)
	testsuite.databaseName = "default"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *RedisenterpriseTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestRedisenterpriseTestSuite(t *testing.T) {
	suite.Run(t, new(RedisenterpriseTestSuite))
}

func (testsuite *RedisenterpriseTestSuite) Prepare() {
	var err error
	// From step RedisEnterprise_Create
	fmt.Println("Call operation: RedisEnterprise_Create")
	client, err := armredisenterprise.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientCreateResponsePoller, err := client.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armredisenterprise.Cluster{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
		},
		Identity: &armredisenterprise.ManagedServiceIdentity{
			Type: to.Ptr(armredisenterprise.ManagedServiceIdentityTypeSystemAssigned),
		},
		Properties: &armredisenterprise.ClusterCreateProperties{
			MinimumTLSVersion: to.Ptr(armredisenterprise.TLSVersionOne2),
		},
		SKU: &armredisenterprise.SKU{
			Name:     to.Ptr(armredisenterprise.SKUNameEnterpriseFlashF300),
			Capacity: to.Ptr[int32](3),
		},
		Zones: []*string{
			to.Ptr("1"),
			to.Ptr("2"),
			to.Ptr("3"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var clientCreateResponse *armredisenterprise.ClientCreateResponse
	clientCreateResponse, err = testutil.PollForTest(testsuite.ctx, clientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.redisEnterpriseId = *clientCreateResponse.ID
}

// Microsoft.Cache/redisEnterprise/{clusterName}
func (testsuite *RedisenterpriseTestSuite) TestRedisEnterprise() {
	var err error
	// From step RedisEnterprise_List
	fmt.Println("Call operation: RedisEnterprise_List")
	client, err := armredisenterprise.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientNewListPager := client.NewListPager(nil)
	for clientNewListPager.More() {
		_, err := clientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RedisEnterprise_ListByResourceGroup
	fmt.Println("Call operation: RedisEnterprise_ListByResourceGroup")
	clientNewListByResourceGroupPager := client.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for clientNewListByResourceGroupPager.More() {
		_, err := clientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RedisEnterprise_Get
	fmt.Println("Call operation: RedisEnterprise_Get")
	_, err = client.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)

	// From step RedisEnterprise_Update
	fmt.Println("Call operation: RedisEnterprise_Update")
	clientUpdateResponsePoller, err := client.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armredisenterprise.ClusterUpdate{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cache/redisEnterprise/{clusterName}/databases/{databaseName}
func (testsuite *RedisenterpriseTestSuite) TestDatabases() {
	var err error
	// From step Databases_Create
	fmt.Println("Call operation: Databases_Create")
	databasesClient, err := armredisenterprise.NewDatabasesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databasesClientCreateResponsePoller, err := databasesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, armredisenterprise.Database{
		Properties: &armredisenterprise.DatabaseCreateProperties{
			ClientProtocol:   to.Ptr(armredisenterprise.ProtocolEncrypted),
			ClusteringPolicy: to.Ptr(armredisenterprise.ClusteringPolicyOSSCluster),
			EvictionPolicy:   to.Ptr(armredisenterprise.EvictionPolicyNoEviction),
			Persistence: &armredisenterprise.Persistence{
				AofEnabled: to.Ptr(false),
				RdbEnabled: to.Ptr(false),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	resumeToken, err := databasesClientCreateResponsePoller.ResumeToken()
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	operationId := regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}").FindAllString(resumeToken, -1)[1]

	// From step OperationsStatus_Get
	fmt.Println("Call operation: OperationsStatus_Get")
	operationsStatusClient, err := armredisenterprise.NewOperationsStatusClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = operationsStatusClient.Get(testsuite.ctx, testsuite.location, operationId, nil)
	testsuite.Require().NoError(err)

	// From step Databases_ListByCluster
	fmt.Println("Call operation: Databases_ListByCluster")
	databasesClientNewListByClusterPager := databasesClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for databasesClientNewListByClusterPager.More() {
		_, err := databasesClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Databases_Get
	fmt.Println("Call operation: Databases_Get")
	_, err = databasesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step Databases_RegenerateKey
	fmt.Println("Call operation: Databases_RegenerateKey")
	databasesClientRegenerateKeyResponsePoller, err := databasesClient.BeginRegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, armredisenterprise.RegenerateKeyParameters{
		KeyType: to.Ptr(armredisenterprise.AccessKeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientRegenerateKeyResponsePoller)
	testsuite.Require().NoError(err)

	// From step Databases_ListKeys
	fmt.Println("Call operation: Databases_ListKeys")
	_, err = databasesClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step Databases_Delete
	fmt.Println("Call operation: Databases_Delete")
	databasesClientDeleteResponsePoller, err := databasesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cache/redisEnterprise/{clusterName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *RedisenterpriseTestSuite) TestPrivateEndpointConnections() {
	var privateEndpointConnectionName string
	var err error
	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "epredisenter-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "epredisenter",
			},
			"redisEnterpriseId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.redisEnterpriseId,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "epredisentervnet",
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2020-11-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"enableDdosProtection": false,
					"subnets": []any{
						map[string]any{
							"name": "default",
							"properties": map[string]any{
								"addressPrefix":                     "10.0.0.0/24",
								"delegations":                       []any{},
								"privateEndpointNetworkPolicies":    "Disabled",
								"privateLinkServiceNetworkPolicies": "Enabled",
							},
						},
					},
					"virtualNetworkPeerings": []any{},
				},
			},
			map[string]any{
				"name":       "[parameters('networkInterfaceName')]",
				"type":       "Microsoft.Network/networkInterfaces",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"dnsSettings": map[string]any{
						"dnsServers": []any{},
					},
					"enableIPForwarding": false,
					"ipConfigurations": []any{
						map[string]any{
							"name": "privateEndpointIpConfig",
							"properties": map[string]any{
								"primary":                   true,
								"privateIPAddress":          "10.0.0.4",
								"privateIPAddressVersion":   "IPv4",
								"privateIPAllocationMethod": "Dynamic",
								"subnet": map[string]any{
									"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
								},
							},
						},
					},
				},
			},
			map[string]any{
				"name":       "[parameters('privateEndpointName')]",
				"type":       "Microsoft.Network/privateEndpoints",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"customDnsConfigs":                    []any{},
					"manualPrivateLinkServiceConnections": []any{},
					"privateLinkServiceConnections": []any{
						map[string]any{
							"name": "[parameters('privateEndpointName')]",
							"properties": map[string]any{
								"groupIds": []any{
									"redisEnterprise",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('redisEnterpriseId')]",
							},
						},
					},
					"subnet": map[string]any{
						"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
					},
				},
			},
			map[string]any{
				"name":       "[concat(parameters('virtualNetworksName'), '/default')]",
				"type":       "Microsoft.Network/virtualNetworks/subnets",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
				},
				"properties": map[string]any{
					"addressPrefix":                     "10.0.0.0/24",
					"delegations":                       []any{},
					"privateEndpointNetworkPolicies":    "Disabled",
					"privateLinkServiceNetworkPolicies": "Enabled",
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
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_List
	fmt.Println("Call operation: PrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armredisenterprise.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Put
	fmt.Println("Call operation: PrivateEndpointConnections_Put")
	privateEndpointConnectionsClientPutResponsePoller, err := privateEndpointConnectionsClient.BeginPut(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, privateEndpointConnectionName, armredisenterprise.PrivateEndpointConnection{
		Properties: &armredisenterprise.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armredisenterprise.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Auto-Approved"),
				Status:      to.Ptr(armredisenterprise.PrivateEndpointServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientPutResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_ListByCluster
	fmt.Println("Call operation: PrivateLinkResources_ListByCluster")
	privateLinkResourcesClient, err := armredisenterprise.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListByClusterPager := privateLinkResourcesClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for privateLinkResourcesClientNewListByClusterPager.More() {
		_, err := privateLinkResourcesClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cache/operations
func (testsuite *RedisenterpriseTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armredisenterprise.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *RedisenterpriseTestSuite) Cleanup() {
	var err error
	// From step RedisEnterprise_Delete
	fmt.Println("Call operation: RedisEnterprise_Delete")
	client, err := armredisenterprise.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientDeleteResponsePoller, err := client.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

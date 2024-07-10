// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmongocluster_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mongocluster/armmongocluster"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type MongoClusterTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	clientFactory     *armmongocluster.ClientFactory
	armEndpoint       string
	firewallRuleName  string
	mongoClusterId    string
	mongoClusterName  string
	adminPassword     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *MongoClusterTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/mongocluster/armmongocluster/testdata")

	var err error
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.firewallRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "firewall", 14, false)
	testsuite.mongoClusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mongoclu", 14, true)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "Sanitized")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.clientFactory, err = armmongocluster.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *MongoClusterTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestMongoClusterTestSuite(t *testing.T) {
	suite.Run(t, new(MongoClusterTestSuite))
}

func (testsuite *MongoClusterTestSuite) Prepare() {
	var err error
	// From step MongoClusters_CreateOrUpdate
	fmt.Println("Call operation: MongoClusters_CreateOrUpdate")
	mongoClustersClient := testsuite.clientFactory.NewMongoClustersClient()
	mongoClustersClientCreateOrUpdateResponsePoller, err := mongoClustersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, armmongocluster.MongoCluster{
		Location: to.Ptr(testsuite.location),
		Properties: &armmongocluster.Properties{
			AdministratorLogin:         to.Ptr("mongoAdmin"),
			AdministratorLoginPassword: to.Ptr(testsuite.adminPassword),
			NodeGroupSpecs: []*armmongocluster.NodeGroupSpec{
				{
					DiskSizeGB: to.Ptr[int64](128),
					EnableHa:   to.Ptr(true),
					Kind:       to.Ptr(armmongocluster.NodeKindShard),
					NodeCount:  to.Ptr[int32](1),
					SKU:        to.Ptr("M30"),
				},
			},
			ServerVersion: to.Ptr("5.0"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var mongoClustersClientCreateOrUpdateResponse *armmongocluster.MongoClustersClientCreateOrUpdateResponse
	mongoClustersClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, mongoClustersClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.mongoClusterId = *mongoClustersClientCreateOrUpdateResponse.ID
}

// Microsoft.DocumentDB/mongoClusters/{mongoClusterName}
func (testsuite *MongoClusterTestSuite) TestMongoClusters() {
	var err error
	// From step MongoClusters_CheckNameAvailability
	fmt.Println("Call operation: MongoClusters_CheckNameAvailability")
	mongoClustersClient := testsuite.clientFactory.NewMongoClustersClient()
	_, err = mongoClustersClient.CheckNameAvailability(testsuite.ctx, testsuite.location, armmongocluster.CheckNameAvailabilityRequest{
		Name: to.Ptr("newmongocluster"),
		Type: to.Ptr("Microsoft.DocumentDB/mongoClusters"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step MongoClusters_List
	fmt.Println("Call operation: MongoClusters_List")
	mongoClustersClientNewListPager := mongoClustersClient.NewListPager(nil)
	for mongoClustersClientNewListPager.More() {
		_, err := mongoClustersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MongoClusters_ListByResourceGroup
	fmt.Println("Call operation: MongoClusters_ListByResourceGroup")
	mongoClustersClientNewListByResourceGroupPager := mongoClustersClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for mongoClustersClientNewListByResourceGroupPager.More() {
		_, err := mongoClustersClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MongoClusters_Get
	fmt.Println("Call operation: MongoClusters_Get")
	_, err = mongoClustersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, nil)
	testsuite.Require().NoError(err)

	// From step MongoClusters_Update
	fmt.Println("Call operation: MongoClusters_Update")
	mongoClustersClientUpdateResponsePoller, err := mongoClustersClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, armmongocluster.Update{
		Tags: map[string]*string{
			"mongoClusterTag": to.Ptr("mongoClusterTagValue"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoClustersClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step MongoClusters_ListConnectionStrings
	fmt.Println("Call operation: MongoClusters_ListConnectionStrings")
	_, err = mongoClustersClient.ListConnectionStrings(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/mongoClusters/{mongoClusterName}/firewallRules/{firewallRuleName}
func (testsuite *MongoClusterTestSuite) TestFirewallRules() {
	var err error
	// From step FirewallRules_CreateOrUpdate
	fmt.Println("Call operation: FirewallRules_CreateOrUpdate")
	firewallRulesClient := testsuite.clientFactory.NewFirewallRulesClient()
	firewallRulesClientCreateOrUpdateResponsePoller, err := firewallRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, testsuite.firewallRuleName, armmongocluster.FirewallRule{
		Properties: &armmongocluster.FirewallRuleProperties{
			EndIPAddress:   to.Ptr("255.255.255.255"),
			StartIPAddress: to.Ptr("0.0.0.0"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallRulesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step FirewallRules_ListByMongoCluster
	fmt.Println("Call operation: FirewallRules_ListByMongoCluster")
	firewallRulesClientNewListByMongoClusterPager := firewallRulesClient.NewListByMongoClusterPager(testsuite.resourceGroupName, testsuite.mongoClusterName, nil)
	for firewallRulesClientNewListByMongoClusterPager.More() {
		_, err := firewallRulesClientNewListByMongoClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FirewallRules_Get
	fmt.Println("Call operation: FirewallRules_Get")
	_, err = firewallRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, testsuite.firewallRuleName, nil)
	testsuite.Require().NoError(err)

	// From step FirewallRules_Delete
	fmt.Println("Call operation: FirewallRules_Delete")
	firewallRulesClientDeleteResponsePoller, err := firewallRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, testsuite.firewallRuleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/mongoClusters/{mongoClusterName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *MongoClusterTestSuite) TestPrivateEndpointConnections() {
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
			"mongoClusterId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.mongoClusterId,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "epmongocluster-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "epmongocluster",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "mongoclustervnet",
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
									"MongoCluster",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('mongoClusterId')]",
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

	// From step PrivateEndpointConnections_ListByMongoCluster
	fmt.Println("Call operation: PrivateEndpointConnections_ListByMongoCluster")
	privateEndpointConnectionsClient := testsuite.clientFactory.NewPrivateEndpointConnectionsClient()
	privateEndpointConnectionsClientNewListByMongoClusterPager := privateEndpointConnectionsClient.NewListByMongoClusterPager(testsuite.resourceGroupName, testsuite.mongoClusterName, nil)
	for privateEndpointConnectionsClientNewListByMongoClusterPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListByMongoClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Create
	fmt.Println("Call operation: PrivateEndpointConnections_Create")
	privateEndpointConnectionsClientCreateResponsePoller, err := privateEndpointConnectionsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, privateEndpointConnectionName, armmongocluster.PrivateEndpointConnectionResource{
		Properties: &armmongocluster.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armmongocluster.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Auto-Approved"),
				Status:      to.Ptr(armmongocluster.PrivateEndpointServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinks_ListByMongoCluster
	fmt.Println("Call operation: PrivateLinks_ListByMongoCluster")
	privateLinksClient := testsuite.clientFactory.NewPrivateLinksClient()
	privateLinksClientNewListByMongoClusterPager := privateLinksClient.NewListByMongoClusterPager(testsuite.resourceGroupName, testsuite.mongoClusterName, nil)
	for privateLinksClientNewListByMongoClusterPager.More() {
		_, err := privateLinksClientNewListByMongoClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/operations
func (testsuite *MongoClusterTestSuite) TestOperation() {
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient := testsuite.clientFactory.NewOperationsClient()
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *MongoClusterTestSuite) Cleanup() {
	var err error
	// From step MongoClusters_Delete
	fmt.Println("Call operation: MongoClusters_Delete")
	mongoClustersClient := testsuite.clientFactory.NewMongoClustersClient()
	mongoClustersClientDeleteResponsePoller, err := mongoClustersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.mongoClusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoClustersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

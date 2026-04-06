// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcosmosforpostgresql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmosforpostgresql/armcosmosforpostgresql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type PostgresqlhscTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	clusterId         string
	clusterName       string
	firewallRuleName  string
	roleName          string
	adminPassword     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *PostgresqlhscTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.clusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clustern", 14, true)
	testsuite.firewallRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "firewall", 14, false)
	testsuite.roleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "rolename", 14, true)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *PostgresqlhscTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPostgresqlhscTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresqlhscTestSuite))
}

func (testsuite *PostgresqlhscTestSuite) Prepare() {
	var err error
	// From step Clusters_CheckNameAvailability
	fmt.Println("Call operation: Clusters_CheckNameAvailability")
	clustersClient, err := armcosmosforpostgresql.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = clustersClient.CheckNameAvailability(testsuite.ctx, armcosmosforpostgresql.NameAvailabilityRequest{
		Name: to.Ptr("name1"),
		Type: to.Ptr("Microsoft.DBforPostgreSQL/serverGroupsv2"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Clusters_Create
	fmt.Println("Call operation: Clusters_Create")
	clustersClientCreateResponsePoller, err := clustersClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armcosmosforpostgresql.Cluster{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmosforpostgresql.ClusterProperties{
			AdministratorLoginPassword:      to.Ptr(testsuite.adminPassword),
			CitusVersion:                    to.Ptr("11.1"),
			CoordinatorEnablePublicIPAccess: to.Ptr(true),
			CoordinatorServerEdition:        to.Ptr("GeneralPurpose"),
			CoordinatorStorageQuotaInMb:     to.Ptr[int32](524288),
			CoordinatorVCores:               to.Ptr[int32](4),
			EnableHa:                        to.Ptr(true),
			EnableShardsOnCoordinator:       to.Ptr(false),
			NodeCount:                       to.Ptr[int32](3),
			NodeEnablePublicIPAccess:        to.Ptr(false),
			NodeServerEdition:               to.Ptr("MemoryOptimized"),
			NodeStorageQuotaInMb:            to.Ptr[int32](524288),
			NodeVCores:                      to.Ptr[int32](8),
			PostgresqlVersion:               to.Ptr("15"),
			PreferredPrimaryZone:            to.Ptr("1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var clustersClientCreateResponse *armcosmosforpostgresql.ClustersClientCreateResponse
	clustersClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, clustersClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.clusterId = *clustersClientCreateResponse.ID
}

// Microsoft.DBforPostgreSQL/serverGroupsv2/{clusterName}
func (testsuite *PostgresqlhscTestSuite) TestClusters() {
	var err error
	// From step Clusters_List
	fmt.Println("Call operation: Clusters_List")
	clustersClient, err := armcosmosforpostgresql.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientNewListPager := clustersClient.NewListPager(nil)
	for clustersClientNewListPager.More() {
		_, err := clustersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_ListByResourceGroup
	fmt.Println("Call operation: Clusters_ListByResourceGroup")
	clustersClientNewListByResourceGroupPager := clustersClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for clustersClientNewListByResourceGroupPager.More() {
		_, err := clustersClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_Get
	fmt.Println("Call operation: Clusters_Get")
	_, err = clustersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)

	// From step Clusters_Update
	fmt.Println("Call operation: Clusters_Update")
	clustersClientUpdateResponsePoller, err := clustersClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armcosmosforpostgresql.ClusterForUpdate{}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_Stop
	fmt.Println("Call operation: Clusters_Stop")
	clustersClientStopResponsePoller, err := clustersClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_Start
	fmt.Println("Call operation: Clusters_Start")
	clustersClientStartResponsePoller, err := clustersClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientStartResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_Restart
	fmt.Println("Call operation: Clusters_Restart")
	clustersClientRestartResponsePoller, err := clustersClient.BeginRestart(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientRestartResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/serverGroupsv2/{clusterName}/firewallRules/{firewallRuleName}
func (testsuite *PostgresqlhscTestSuite) TestFirewallRules() {
	var err error
	// From step FirewallRules_CreateOrUpdate
	fmt.Println("Call operation: FirewallRules_CreateOrUpdate")
	firewallRulesClient, err := armcosmosforpostgresql.NewFirewallRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	firewallRulesClientCreateOrUpdateResponsePoller, err := firewallRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.firewallRuleName, armcosmosforpostgresql.FirewallRule{
		Properties: &armcosmosforpostgresql.FirewallRuleProperties{
			EndIPAddress:   to.Ptr("255.255.255.255"),
			StartIPAddress: to.Ptr("0.0.0.0"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallRulesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step FirewallRules_ListByCluster
	fmt.Println("Call operation: FirewallRules_ListByCluster")
	firewallRulesClientNewListByClusterPager := firewallRulesClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for firewallRulesClientNewListByClusterPager.More() {
		_, err := firewallRulesClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FirewallRules_Get
	fmt.Println("Call operation: FirewallRules_Get")
	_, err = firewallRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.firewallRuleName, nil)
	testsuite.Require().NoError(err)

	// From step FirewallRules_Delete
	fmt.Println("Call operation: FirewallRules_Delete")
	firewallRulesClientDeleteResponsePoller, err := firewallRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.firewallRuleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/serverGroupsv2/{clusterName}/roles/{roleName}
func (testsuite *PostgresqlhscTestSuite) TestRoles() {
	var err error
	// From step Roles_Create
	fmt.Println("Call operation: Roles_Create")
	rolesClient, err := armcosmosforpostgresql.NewRolesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	rolesClientCreateResponsePoller, err := rolesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.roleName, armcosmosforpostgresql.Role{
		Properties: &armcosmosforpostgresql.RoleProperties{
			Password: to.Ptr(testsuite.adminPassword),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, rolesClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Roles_ListByCluster
	fmt.Println("Call operation: Roles_ListByCluster")
	rolesClientNewListByClusterPager := rolesClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for rolesClientNewListByClusterPager.More() {
		_, err := rolesClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Roles_Get
	fmt.Println("Call operation: Roles_Get")
	_, err = rolesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.roleName, nil)
	testsuite.Require().NoError(err)

	// From step Roles_Delete
	fmt.Println("Call operation: Roles_Delete")
	rolesClientDeleteResponsePoller, err := rolesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.roleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, rolesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/serverGroupsv2/{clusterName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *PostgresqlhscTestSuite) TestPrivateEndpointConnections() {
	var err error
	var privateEndpointConnectionName string
	var privateLinkResourceName string
	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]any{
			"clusterId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.clusterId,
			},
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "endpointpgsqlhsc-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "endpointpgsqlhsc",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "pgsqlhscvnet",
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
									"coordinator",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('clusterId')]",
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

	// From step PrivateEndpointConnections_ListByCluster
	fmt.Println("Call operation: PrivateEndpointConnections_ListByCluster")
	privateEndpointConnectionsClient, err := armcosmosforpostgresql.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListByClusterPager := privateEndpointConnectionsClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for privateEndpointConnectionsClientNewListByClusterPager.More() {
		privateEndpointConnectionsClientListByClusterResponse, err := privateEndpointConnectionsClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		privateEndpointConnectionName = *privateEndpointConnectionsClientListByClusterResponse.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_CreateOrUpdate
	fmt.Println("Call operation: PrivateEndpointConnections_CreateOrUpdate")
	privateEndpointConnectionsClientCreateOrUpdateResponsePoller, err := privateEndpointConnectionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, privateEndpointConnectionName, armcosmosforpostgresql.PrivateEndpointConnection{
		Properties: &armcosmosforpostgresql.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armcosmosforpostgresql.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Approved by johndoe@contoso.com"),
				Status:      to.Ptr(armcosmosforpostgresql.PrivateEndpointServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_ListByCluster
	fmt.Println("Call operation: PrivateLinkResources_ListByCluster")
	privateLinkResourcesClient, err := armcosmosforpostgresql.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListByClusterPager := privateLinkResourcesClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for privateLinkResourcesClientNewListByClusterPager.More() {
		privateLinkResourcesClientListByClusterResponse, err := privateLinkResourcesClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		privateLinkResourceName = *privateLinkResourcesClientListByClusterResponse.Value[0].Name
		break
	}

	// From step PrivateLinkResources_Get
	fmt.Println("Call operation: PrivateLinkResources_Get")
	_, err = privateLinkResourcesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, privateLinkResourceName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/serverGroupsv2/{clusterName}/coordinatorConfigurations/{configurationName}
func (testsuite *PostgresqlhscTestSuite) TestConfigurations() {
	var err error
	var configurationName = "array_nulls"
	// From step Configurations_UpdateOnCoordinator
	fmt.Println("Call operation: Configurations_UpdateOnCoordinator")
	configurationsClient, err := armcosmosforpostgresql.NewConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configurationsClientUpdateOnCoordinatorResponsePoller, err := configurationsClient.BeginUpdateOnCoordinator(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, configurationName, armcosmosforpostgresql.ServerConfiguration{
		Properties: &armcosmosforpostgresql.ServerConfigurationProperties{
			Value: to.Ptr("on"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationsClientUpdateOnCoordinatorResponsePoller)
	testsuite.Require().NoError(err)

	// From step Configurations_ListByCluster
	fmt.Println("Call operation: Configurations_ListByCluster")
	configurationsClientNewListByClusterPager := configurationsClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for configurationsClientNewListByClusterPager.More() {
		_, err := configurationsClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Configurations_Get
	fmt.Println("Call operation: Configurations_Get")
	_, err = configurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, configurationName, nil)
	testsuite.Require().NoError(err)

	// From step Configurations_GetCoordinator
	fmt.Println("Call operation: Configurations_GetCoordinator")
	_, err = configurationsClient.GetCoordinator(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, configurationName, nil)
	testsuite.Require().NoError(err)

	// From step Configurations_UpdateOnNode
	fmt.Println("Call operation: Configurations_UpdateOnNode")
	configurationsClientUpdateOnNodeResponsePoller, err := configurationsClient.BeginUpdateOnNode(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, configurationName, armcosmosforpostgresql.ServerConfiguration{
		Properties: &armcosmosforpostgresql.ServerConfigurationProperties{
			Value: to.Ptr("off"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationsClientUpdateOnNodeResponsePoller)
	testsuite.Require().NoError(err)

	// From step Configurations_GetNode
	fmt.Println("Call operation: Configurations_GetNode")
	_, err = configurationsClient.GetNode(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, configurationName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/operations
func (testsuite *PostgresqlhscTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armcosmosforpostgresql.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforPostgreSQL/serverGroupsv2/{clusterName}/servers
func (testsuite *PostgresqlhscTestSuite) TestServers() {
	var err error
	var serverName string
	// From step Servers_ListByCluster
	fmt.Println("Call operation: Servers_ListByCluster")
	serversClient, err := armcosmosforpostgresql.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClientNewListByClusterPager := serversClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for serversClientNewListByClusterPager.More() {
		serversClientListByClusterResponse, err := serversClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		serverName = *serversClientListByClusterResponse.Value[0].Name
		break
	}

	// From step Servers_Get
	fmt.Println("Call operation: Servers_Get")
	_, err = serversClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, serverName, nil)
	testsuite.Require().NoError(err)

	// From step Configurations_ListByServer
	fmt.Println("Call operation: Configurations_ListByServer")
	configurationsClient, err := armcosmosforpostgresql.NewConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configurationsClientNewListByServerPager := configurationsClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.clusterName, serverName, nil)
	for configurationsClientNewListByServerPager.More() {
		_, err := configurationsClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *PostgresqlhscTestSuite) Cleanup() {
	var err error
	// From step Clusters_Delete
	fmt.Println("Call operation: Clusters_Delete")
	clustersClient, err := armcosmosforpostgresql.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientDeleteResponsePoller, err := clustersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

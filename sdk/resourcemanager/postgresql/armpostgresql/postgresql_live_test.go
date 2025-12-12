// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armpostgresql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type PostgresqlTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	serverName        string
	adminPassword     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *PostgresqlTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serverna", 14, true)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *PostgresqlTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPostgresqlTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresqlTestSuite))
}

func (testsuite *PostgresqlTestSuite) Prepare() {
	var err error
	// From step Servers_Create
	fmt.Println("Call operation: Servers_Create")
	serversClient, err := armpostgresql.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClientCreateResponsePoller, err := serversClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armpostgresql.ServerForCreate{
		Location: to.Ptr(testsuite.location),
		Properties: &armpostgresql.ServerPropertiesForDefaultCreate{
			CreateMode:        to.Ptr(armpostgresql.CreateModeDefault),
			MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS12),
			SSLEnforcement:    to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
			StorageProfile: &armpostgresql.StorageProfile{
				BackupRetentionDays: to.Ptr[int32](7),
				GeoRedundantBackup:  to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
				StorageMB:           to.Ptr[int32](128000),
			},
			AdministratorLogin:         to.Ptr("cloudsa"),
			AdministratorLoginPassword: to.Ptr(testsuite.adminPassword),
		},
		SKU: &armpostgresql.SKU{
			Name:   to.Ptr("GP_Gen5_8"),
			Family: to.Ptr("Gen5"),
			Tier:   to.Ptr(armpostgresql.SKUTierGeneralPurpose),
		},
		Tags: map[string]*string{
			"ElasticServer": to.Ptr("1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientCreateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/{serverName}
func (testsuite *PostgresqlTestSuite) TestServers() {
	var err error
	// From step Servers_List
	fmt.Println("Call operation: Servers_List")
	serversClient, err := armpostgresql.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClientNewListPager := serversClient.NewListPager(nil)
	for serversClientNewListPager.More() {
		_, err := serversClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Servers_ListByResourceGroup
	fmt.Println("Call operation: Servers_ListByResourceGroup")
	serversClientNewListByResourceGroupPager := serversClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for serversClientNewListByResourceGroupPager.More() {
		_, err := serversClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Servers_Get
	fmt.Println("Call operation: Servers_Get")
	_, err = serversClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)

	// From step Servers_Update
	fmt.Println("Call operation: Servers_Update")
	serversClientUpdateResponsePoller, err := serversClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armpostgresql.ServerUpdateParameters{
		Properties: &armpostgresql.ServerUpdateParametersProperties{
			AdministratorLoginPassword: to.Ptr(testsuite.adminPassword),
			MinimalTLSVersion:          to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS12),
			SSLEnforcement:             to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Servers_Restart
	fmt.Println("Call operation: Servers_Restart")
	serversClientRestartResponsePoller, err := serversClient.BeginRestart(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientRestartResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/{serverName}/firewallRules/{firewallRuleName}
func (testsuite *PostgresqlTestSuite) TestFirewallRules() {
	var err error
	// From step FirewallRules_CreateOrUpdate
	fmt.Println("Call operation: FirewallRules_CreateOrUpdate")
	firewallRulesClient, err := armpostgresql.NewFirewallRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	firewallRulesClientCreateOrUpdateResponsePoller, err := firewallRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "rule1", armpostgresql.FirewallRule{
		Properties: &armpostgresql.FirewallRuleProperties{
			EndIPAddress:   to.Ptr("255.255.255.255"),
			StartIPAddress: to.Ptr("0.0.0.0"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallRulesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step FirewallRules_ListByServer
	fmt.Println("Call operation: FirewallRules_ListByServer")
	firewallRulesClientNewListByServerPager := firewallRulesClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for firewallRulesClientNewListByServerPager.More() {
		_, err := firewallRulesClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FirewallRules_Get
	fmt.Println("Call operation: FirewallRules_Get")
	_, err = firewallRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "rule1", nil)
	testsuite.Require().NoError(err)

	// From step FirewallRules_Delete
	fmt.Println("Call operation: FirewallRules_Delete")
	firewallRulesClientDeleteResponsePoller, err := firewallRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "rule1", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/{serverName}/virtualNetworkRules/{virtualNetworkRuleName}
func (testsuite *PostgresqlTestSuite) TestVirtualNetworkRules() {
	var subnetId string
	var err error
	// From step VirtualNetwork_Create
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"subnetId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), parameters('subnetName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"subnetName": map[string]any{
				"type":         "string",
				"defaultValue": "pgsubnet",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "pgvnet",
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
							"name": "[parameters('subnetName')]",
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "VirtualNetwork_Create", &deployment)
	testsuite.Require().NoError(err)
	subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step VirtualNetworkRules_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworkRules_CreateOrUpdate")
	virtualNetworkRulesClient, err := armpostgresql.NewVirtualNetworkRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworkRulesClientCreateOrUpdateResponsePoller, err := virtualNetworkRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "vnet-firewall-rule", armpostgresql.VirtualNetworkRule{
		Properties: &armpostgresql.VirtualNetworkRuleProperties{
			IgnoreMissingVnetServiceEndpoint: to.Ptr(true),
			VirtualNetworkSubnetID:           to.Ptr(subnetId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkRulesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkRules_ListByServer
	fmt.Println("Call operation: VirtualNetworkRules_ListByServer")
	virtualNetworkRulesClientNewListByServerPager := virtualNetworkRulesClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for virtualNetworkRulesClientNewListByServerPager.More() {
		_, err := virtualNetworkRulesClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworkRules_Get
	fmt.Println("Call operation: VirtualNetworkRules_Get")
	_, err = virtualNetworkRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "vnet-firewall-rule", nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkRules_Delete
	fmt.Println("Call operation: VirtualNetworkRules_Delete")
	virtualNetworkRulesClientDeleteResponsePoller, err := virtualNetworkRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "vnet-firewall-rule", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/{serverName}/databases/{databaseName}
func (testsuite *PostgresqlTestSuite) TestDatabases() {
	var err error
	// From step Databases_CreateOrUpdate
	fmt.Println("Call operation: Databases_CreateOrUpdate")
	databasesClient, err := armpostgresql.NewDatabasesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databasesClientCreateOrUpdateResponsePoller, err := databasesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "db1", armpostgresql.Database{
		Properties: &armpostgresql.DatabaseProperties{
			Charset:   to.Ptr("UTF8"),
			Collation: to.Ptr("English_United States.1252"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Databases_ListByServer
	fmt.Println("Call operation: Databases_ListByServer")
	databasesClientNewListByServerPager := databasesClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for databasesClientNewListByServerPager.More() {
		_, err := databasesClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Databases_Get
	fmt.Println("Call operation: Databases_Get")
	_, err = databasesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "db1", nil)
	testsuite.Require().NoError(err)

	// From step Databases_Delete
	fmt.Println("Call operation: Databases_Delete")
	databasesClientDeleteResponsePoller, err := databasesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "db1", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/{serverName}/configurations/{configurationName}
func (testsuite *PostgresqlTestSuite) TestConfigurations() {
	var err error
	// From step Configurations_CreateOrUpdate
	fmt.Println("Call operation: Configurations_CreateOrUpdate")
	configurationsClient, err := armpostgresql.NewConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configurationsClientCreateOrUpdateResponsePoller, err := configurationsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "array_nulls", armpostgresql.Configuration{
		Properties: &armpostgresql.ConfigurationProperties{
			Source: to.Ptr("user-override"),
			Value:  to.Ptr("off"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Configurations_ListByServer
	fmt.Println("Call operation: Configurations_ListByServer")
	configurationsClientNewListByServerPager := configurationsClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for configurationsClientNewListByServerPager.More() {
		_, err := configurationsClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Configurations_Get
	fmt.Println("Call operation: Configurations_Get")
	_, err = configurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "array_nulls", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/{serverName}/administrators/activeDirectory
func (testsuite *PostgresqlTestSuite) TestServerAdministrators() {
	var err error
	// From step ServerAdministrators_CreateOrUpdate
	fmt.Println("Call operation: ServerAdministrators_CreateOrUpdate")
	serverAdministratorsClient, err := armpostgresql.NewServerAdministratorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serverAdministratorsClientCreateOrUpdateResponsePoller, err := serverAdministratorsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armpostgresql.ServerAdministratorResource{
		Properties: &armpostgresql.ServerAdministratorProperties{
			AdministratorType: to.Ptr("ActiveDirectory"),
			Login:             to.Ptr("bob@contoso.com"),
			Sid:               to.Ptr("c6b82b90-a647-49cb-8a62-0d2d3cb7ac7c"),
			TenantID:          to.Ptr("c6b82b90-a647-49cb-8a62-0d2d3cb7ac7c"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serverAdministratorsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ServerAdministrators_List
	fmt.Println("Call operation: ServerAdministrators_List")
	serverAdministratorsClientNewListPager := serverAdministratorsClient.NewListPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for serverAdministratorsClientNewListPager.More() {
		_, err := serverAdministratorsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ServerAdministrators_Get
	fmt.Println("Call operation: ServerAdministrators_Get")
	_, err = serverAdministratorsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)

	// From step ServerAdministrators_Delete
	fmt.Println("Call operation: ServerAdministrators_Delete")
	serverAdministratorsClientDeleteResponsePoller, err := serverAdministratorsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serverAdministratorsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/replicas
func (testsuite *PostgresqlTestSuite) TestReplicas() {
	var err error
	// From step Replicas_ListByServer
	fmt.Println("Call operation: Replicas_ListByServer")
	replicasClient, err := armpostgresql.NewReplicasClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	replicasClientNewListByServerPager := replicasClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for replicasClientNewListByServerPager.More() {
		_, err := replicasClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforPostgreSQL/servers/updateConfigurations
func (testsuite *PostgresqlTestSuite) TestServerParameters() {
	var err error
	// From step ServerParameters_ListUpdateConfigurations
	fmt.Println("Call operation: ServerParameters_ListUpdateConfigurations")
	serverParametersClient, err := armpostgresql.NewServerParametersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serverParametersClientListUpdateConfigurationsResponsePoller, err := serverParametersClient.BeginListUpdateConfigurations(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armpostgresql.ConfigurationListResult{
		Value: []*armpostgresql.Configuration{
			{
				Name: to.Ptr("array_nulls"),
				Properties: &armpostgresql.ConfigurationProperties{
					Value: to.Ptr("on"),
				},
			},
			{
				Name: to.Ptr("backslash_quote"),
				Properties: &armpostgresql.ConfigurationProperties{
					Value: to.Ptr("on"),
				},
			}},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serverParametersClientListUpdateConfigurationsResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/logFiles
func (testsuite *PostgresqlTestSuite) TestLogFiles() {
	var err error
	// From step LogFiles_ListByServer
	fmt.Println("Call operation: LogFiles_ListByServer")
	logFilesClient, err := armpostgresql.NewLogFilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	logFilesClientNewListByServerPager := logFilesClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for logFilesClientNewListByServerPager.More() {
		_, err := logFilesClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforPostgreSQL/servers/recoverableServers
func (testsuite *PostgresqlTestSuite) TestRecoverableServers() {
	var err error
	// From step RecoverableServers_Get
	fmt.Println("Call operation: RecoverableServers_Get")
	recoverableServersClient, err := armpostgresql.NewRecoverableServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = recoverableServersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/servers/performanceTiers
func (testsuite *PostgresqlTestSuite) TestServerBasedPerformanceTier() {
	var err error
	// From step ServerBasedPerformanceTier_List
	fmt.Println("Call operation: ServerBasedPerformanceTier_List")
	serverBasedPerformanceTierClient, err := armpostgresql.NewServerBasedPerformanceTierClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serverBasedPerformanceTierClientNewListPager := serverBasedPerformanceTierClient.NewListPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for serverBasedPerformanceTierClientNewListPager.More() {
		_, err := serverBasedPerformanceTierClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforPostgreSQL/locations/performanceTiers
func (testsuite *PostgresqlTestSuite) TestLocationBasedPerformanceTier() {
	var err error
	// From step LocationBasedPerformanceTier_List
	fmt.Println("Call operation: LocationBasedPerformanceTier_List")
	locationBasedPerformanceTierClient, err := armpostgresql.NewLocationBasedPerformanceTierClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	locationBasedPerformanceTierClientNewListPager := locationBasedPerformanceTierClient.NewListPager("WestUS", nil)
	for locationBasedPerformanceTierClientNewListPager.More() {
		_, err := locationBasedPerformanceTierClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforPostgreSQL/checkNameAvailability
func (testsuite *PostgresqlTestSuite) TestCheckNameAvailability() {
	var err error
	// From step CheckNameAvailability_Execute
	fmt.Println("Call operation: CheckNameAvailability_Execute")
	checkNameAvailabilityClient, err := armpostgresql.NewCheckNameAvailabilityClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = checkNameAvailabilityClient.Execute(testsuite.ctx, armpostgresql.NameAvailabilityRequest{
		Name: to.Ptr("name1"),
		Type: to.Ptr("Microsoft.DBforPostgreSQL"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforPostgreSQL/operations
func (testsuite *PostgresqlTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armpostgresql.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = operationsClient.List(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *PostgresqlTestSuite) Cleanup() {
	var err error
	// From step Servers_Delete
	fmt.Println("Call operation: Servers_Delete")
	serversClient, err := armpostgresql.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClientDeleteResponsePoller, err := serversClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmariadb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mariadb/armmariadb"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type MariadbTestSuite struct {
	suite.Suite

	ctx                    context.Context
	cred                   azcore.TokenCredential
	options                *arm.ClientOptions
	armEndpoint            string
	databaseName           string
	firewallRuleName       string
	serverName             string
	virtualNetworkRuleName string
	adminPassword          string
	location               string
	resourceGroupName      string
	subscriptionId         string
}

func (testsuite *MariadbTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.databaseName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "database", 14, false)
	testsuite.firewallRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "firewall", 14, false)
	testsuite.serverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serverna", 14, false)
	testsuite.virtualNetworkRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualn", 14, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *MariadbTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestMariadbTestSuite(t *testing.T) {
	suite.Run(t, new(MariadbTestSuite))
}

func (testsuite *MariadbTestSuite) Prepare() {
	var err error
	// From step Servers_Create
	fmt.Println("Call operation: Servers_Create")
	serversClient, err := armmariadb.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClientCreateResponsePoller, err := serversClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armmariadb.ServerForCreate{
		Location: to.Ptr(testsuite.location),
		Properties: &armmariadb.ServerPropertiesForDefaultCreate{
			CreateMode:        to.Ptr(armmariadb.CreateModeDefault),
			MinimalTLSVersion: to.Ptr(armmariadb.MinimalTLSVersionEnumTLS12),
			SSLEnforcement:    to.Ptr(armmariadb.SSLEnforcementEnumEnabled),
			StorageProfile: &armmariadb.StorageProfile{
				BackupRetentionDays: to.Ptr[int32](7),
				GeoRedundantBackup:  to.Ptr(armmariadb.GeoRedundantBackupEnabled),
				StorageMB:           to.Ptr[int32](128000),
			},
			AdministratorLogin:         to.Ptr("cloudsa"),
			AdministratorLoginPassword: to.Ptr(testsuite.adminPassword),
		},
		SKU: &armmariadb.SKU{
			Name:     to.Ptr("GP_Gen5_2"),
			Capacity: to.Ptr[int32](2),
			Family:   to.Ptr("Gen5"),
			Tier:     to.Ptr(armmariadb.SKUTierGeneralPurpose),
		},
		Tags: map[string]*string{
			"ElasticServer": to.Ptr("1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientCreateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}
func (testsuite *MariadbTestSuite) TestServers() {
	var err error
	// From step Servers_List
	fmt.Println("Call operation: Servers_List")
	serversClient, err := armmariadb.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
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
	serversClientUpdateResponsePoller, err := serversClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armmariadb.ServerUpdateParameters{
		Properties: &armmariadb.ServerUpdateParametersProperties{
			AdministratorLoginPassword: to.Ptr("<administratorLoginPassword>"),
			SSLEnforcement:             to.Ptr(armmariadb.SSLEnforcementEnumDisabled),
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

	// From step Servers_Stop
	fmt.Println("Call operation: Servers_Stop")
	serversClientStopResponsePoller, err := serversClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step Servers_Start
	fmt.Println("Call operation: Servers_Start")
	serversClientStartResponsePoller, err := serversClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientStartResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/replicas
func (testsuite *MariadbTestSuite) TestReplicas() {
	var err error
	// From step Replicas_ListByServer
	fmt.Println("Call operation: Replicas_ListByServer")
	replicasClient, err := armmariadb.NewReplicasClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	replicasClientNewListByServerPager := replicasClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for replicasClientNewListByServerPager.More() {
		_, err := replicasClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforMariaDB/servers/{serverName}/firewallRules/{firewallRuleName}
func (testsuite *MariadbTestSuite) TestFirewallRules() {
	var err error
	// From step FirewallRules_CreateOrUpdate
	fmt.Println("Call operation: FirewallRules_CreateOrUpdate")
	firewallRulesClient, err := armmariadb.NewFirewallRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	firewallRulesClientCreateOrUpdateResponsePoller, err := firewallRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.firewallRuleName, armmariadb.FirewallRule{
		Properties: &armmariadb.FirewallRuleProperties{
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
	_, err = firewallRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.firewallRuleName, nil)
	testsuite.Require().NoError(err)

	// From step FirewallRules_Delete
	fmt.Println("Call operation: FirewallRules_Delete")
	firewallRulesClientDeleteResponsePoller, err := firewallRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.firewallRuleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/virtualNetworkRules/{virtualNetworkRuleName}
func (testsuite *MariadbTestSuite) TestVirtualNetworkRules() {
	var subnetId string
	var err error
	// From step Create_Subnet
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
				"defaultValue": "mariadbvnet",
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
								"serviceEndpoints": []any{
									map[string]any{
										"service": "Microsoft.Sql",
									},
								},
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_Subnet", &deployment)
	testsuite.Require().NoError(err)
	subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step VirtualNetworkRules_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworkRules_CreateOrUpdate")
	virtualNetworkRulesClient, err := armmariadb.NewVirtualNetworkRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworkRulesClientCreateOrUpdateResponsePoller, err := virtualNetworkRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.virtualNetworkRuleName, armmariadb.VirtualNetworkRule{
		Properties: &armmariadb.VirtualNetworkRuleProperties{
			IgnoreMissingVnetServiceEndpoint: to.Ptr(false),
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
	_, err = virtualNetworkRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.virtualNetworkRuleName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkRules_Delete
	fmt.Println("Call operation: VirtualNetworkRules_Delete")
	virtualNetworkRulesClientDeleteResponsePoller, err := virtualNetworkRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.virtualNetworkRuleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/databases/{databaseName}
func (testsuite *MariadbTestSuite) TestDatabases() {
	var err error
	// From step Databases_CreateOrUpdate
	fmt.Println("Call operation: Databases_CreateOrUpdate")
	databasesClient, err := armmariadb.NewDatabasesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databasesClientCreateOrUpdateResponsePoller, err := databasesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.databaseName, armmariadb.Database{
		Properties: &armmariadb.DatabaseProperties{
			Charset:   to.Ptr("utf8"),
			Collation: to.Ptr("utf8_general_ci"),
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
	_, err = databasesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step Databases_Delete
	fmt.Println("Call operation: Databases_Delete")
	databasesClientDeleteResponsePoller, err := databasesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/configurations/{configurationName}
func (testsuite *MariadbTestSuite) TestConfigurations() {
	var err error
	// From step Configurations_CreateOrUpdate
	fmt.Println("Call operation: Configurations_CreateOrUpdate")
	configurationsClient, err := armmariadb.NewConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configurationsClientCreateOrUpdateResponsePoller, err := configurationsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "event_scheduler", armmariadb.Configuration{
		Properties: &armmariadb.ConfigurationProperties{
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
	_, err = configurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "event_scheduler", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/updateConfigurations
func (testsuite *MariadbTestSuite) TestServerParameters() {
	var err error
	// From step ServerParameters_ListUpdateConfigurations
	fmt.Println("Call operation: ServerParameters_ListUpdateConfigurations")
	serverParametersClient, err := armmariadb.NewServerParametersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serverParametersClientListUpdateConfigurationsResponsePoller, err := serverParametersClient.BeginListUpdateConfigurations(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armmariadb.ConfigurationListResult{
		Value: []*armmariadb.Configuration{
			{
				Name: to.Ptr("event_scheduler"),
				Properties: &armmariadb.ConfigurationProperties{
					Value: to.Ptr("OFF"),
				},
			}},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serverParametersClientListUpdateConfigurationsResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/logFiles
func (testsuite *MariadbTestSuite) TestLogFiles() {
	var err error
	// From step LogFiles_ListByServer
	fmt.Println("Call operation: LogFiles_ListByServer")
	logFilesClient, err := armmariadb.NewLogFilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	logFilesClientNewListByServerPager := logFilesClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for logFilesClientNewListByServerPager.More() {
		_, err := logFilesClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforMariaDB/servers/{serverName}/recoverableServers
func (testsuite *MariadbTestSuite) TestRecoverableServers() {
	var err error
	// From step RecoverableServers_Get
	fmt.Println("Call operation: RecoverableServers_Get")
	recoverableServersClient, err := armmariadb.NewRecoverableServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = recoverableServersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/performanceTiers
func (testsuite *MariadbTestSuite) TestServerBasedPerformanceTier() {
	var err error
	// From step ServerBasedPerformanceTier_List
	fmt.Println("Call operation: ServerBasedPerformanceTier_List")
	serverBasedPerformanceTierClient, err := armmariadb.NewServerBasedPerformanceTierClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serverBasedPerformanceTierClientNewListPager := serverBasedPerformanceTierClient.NewListPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for serverBasedPerformanceTierClientNewListPager.More() {
		_, err := serverBasedPerformanceTierClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforMariaDB/locations/{locationName}/performanceTiers
func (testsuite *MariadbTestSuite) TestLocationBasedPerformanceTier() {
	var err error
	// From step LocationBasedPerformanceTier_List
	fmt.Println("Call operation: LocationBasedPerformanceTier_List")
	locationBasedPerformanceTierClient, err := armmariadb.NewLocationBasedPerformanceTierClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	locationBasedPerformanceTierClientNewListPager := locationBasedPerformanceTierClient.NewListPager("WestUS", nil)
	for locationBasedPerformanceTierClientNewListPager.More() {
		_, err := locationBasedPerformanceTierClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DBforMariaDB/checkNameAvailability
func (testsuite *MariadbTestSuite) TestCheckNameAvailability() {
	var err error
	// From step CheckNameAvailability_Execute
	fmt.Println("Call operation: CheckNameAvailability_Execute")
	checkNameAvailabilityClient, err := armmariadb.NewCheckNameAvailabilityClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = checkNameAvailabilityClient.Execute(testsuite.ctx, armmariadb.NameAvailabilityRequest{
		Name: to.Ptr("name1"),
		Type: to.Ptr("Microsoft.DBforMariaDB"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/operations
func (testsuite *MariadbTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armmariadb.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = operationsClient.List(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *MariadbTestSuite) Cleanup() {
	var err error
	// From step Servers_Delete
	fmt.Println("Call operation: Servers_Delete")
	serversClient, err := armmariadb.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClientDeleteResponsePoller, err := serversClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

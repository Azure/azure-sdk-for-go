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
	serversClientCreateResponsePoller, err := serversClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armpostgresql.Server{
		Location: to.Ptr(testsuite.location),
		Properties: &armpostgresql.ServerProperties{
			AdministratorLogin:         to.Ptr("examplelogin"),
			AdministratorLoginPassword: to.Ptr("examplepassword"),
			Backup: &armpostgresql.Backup{
				BackupRetentionDays: to.Ptr[int32](7),
				GeoRedundantBackup:  to.Ptr(armpostgresql.GeographicallyRedundantBackupDisabled),
			},
			Cluster: &armpostgresql.Cluster{
				ClusterSize:         to.Ptr[int32](2),
				DefaultDatabaseName: to.Ptr("clusterdb"),
			},
			CreateMode: to.Ptr(armpostgresql.CreateModeCreate),
			HighAvailability: &armpostgresql.HighAvailability{
				Mode: to.Ptr(armpostgresql.FlexibleServerHighAvailabilityModeDisabled),
			},
			Network: &armpostgresql.Network{
				PublicNetworkAccess: to.Ptr(armpostgresql.ServerPublicNetworkAccessStateDisabled),
			},
			Storage: &armpostgresql.Storage{
				AutoGrow:      to.Ptr(armpostgresql.StorageAutoGrowDisabled),
				StorageSizeGB: to.Ptr[int32](256),
				Tier:          to.Ptr(armpostgresql.AzureManagedDiskPerformanceTierP15),
			},
			Version: to.Ptr(armpostgresql.PostgresMajorVersion16),
		},
		SKU: &armpostgresql.SKU{
			Name: to.Ptr("Standard_D4ds_v5"),
			Tier: to.Ptr(armpostgresql.SKUTierGeneralPurpose),
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
	serversClientUpdateResponsePoller, err := serversClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armpostgresql.ServerForPatch{
		Properties: &armpostgresql.ServerPropertiesForPatch{
			Replica: &armpostgresql.Replica{
				PromoteMode:   to.Ptr(armpostgresql.ReadReplicaPromoteModeStandalone),
				PromoteOption: to.Ptr(armpostgresql.ReadReplicaPromoteOptionForced),
			},
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

// Microsoft.DBforPostgreSQL/servers/{serverName}/databases/{databaseName}
func (testsuite *PostgresqlTestSuite) TestDatabases() {
	var err error
	// From step Databases_CreateOrUpdate
	fmt.Println("Call operation: Databases_CreateOrUpdate")
	databasesClient, err := armpostgresql.NewDatabasesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databasesClientCreateOrUpdateResponsePoller, err := databasesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "db1", armpostgresql.Database{
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

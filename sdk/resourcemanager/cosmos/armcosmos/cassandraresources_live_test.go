// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type CassandraResourcesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	accountName       string
	keyspaceName      string
	tableName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *CassandraResourcesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.keyspaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "keyspace", 14, false)
	testsuite.tableName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "tablenam", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *CassandraResourcesTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestCassandraResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(CassandraResourcesTestSuite))
}

func (testsuite *CassandraResourcesTestSuite) Prepare() {
	var err error
	// From step DatabaseAccounts_CreateOrUpdate
	fmt.Println("Call operation: DatabaseAccounts_CreateOrUpdate")
	databaseAccountsClient, err := armcosmos.NewDatabaseAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databaseAccountsClientCreateOrUpdateResponsePoller, err := databaseAccountsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.DatabaseAccountCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Properties: &armcosmos.DatabaseAccountCreateUpdateProperties{
			Capabilities: []*armcosmos.Capability{
				{
					Name: to.Ptr("EnableCassandra"),
				}},
			CreateMode:               to.Ptr(armcosmos.CreateModeDefault),
			DatabaseAccountOfferType: to.Ptr("Standard"),
			Locations: []*armcosmos.Location{
				{
					FailoverPriority: to.Ptr[int32](0),
					IsZoneRedundant:  to.Ptr(false),
					LocationName:     to.Ptr(testsuite.location),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databaseAccountsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/cassandraKeyspaces/{keyspaceName}
func (testsuite *CassandraResourcesTestSuite) TestCassandraKeyspace() {
	var err error
	// From step CassandraResources_CreateUpdateCassandraKeyspace
	fmt.Println("Call operation: CassandraResources_CreateUpdateCassandraKeyspace")
	cassandraResourcesClient, err := armcosmos.NewCassandraResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	cassandraResourcesClientCreateUpdateCassandraKeyspaceResponsePoller, err := cassandraResourcesClient.BeginCreateUpdateCassandraKeyspace(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, armcosmos.CassandraKeyspaceCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.CassandraKeyspaceCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{
				Throughput: to.Ptr[int32](2000),
			},
			Resource: &armcosmos.CassandraKeyspaceResource{
				ID: to.Ptr(testsuite.keyspaceName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientCreateUpdateCassandraKeyspaceResponsePoller)
	testsuite.Require().NoError(err)

	// From step CassandraResources_ListCassandraKeyspaces
	fmt.Println("Call operation: CassandraResources_ListCassandraKeyspaces")
	cassandraResourcesClientNewListCassandraKeyspacesPager := cassandraResourcesClient.NewListCassandraKeyspacesPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for cassandraResourcesClientNewListCassandraKeyspacesPager.More() {
		_, err := cassandraResourcesClientNewListCassandraKeyspacesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step CassandraResources_GetCassandraKeyspaceThroughput
	fmt.Println("Call operation: CassandraResources_GetCassandraKeyspaceThroughput")
	_, err = cassandraResourcesClient.GetCassandraKeyspaceThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, nil)
	testsuite.Require().NoError(err)

	// From step CassandraResources_GetCassandraKeyspace
	fmt.Println("Call operation: CassandraResources_GetCassandraKeyspace")
	_, err = cassandraResourcesClient.GetCassandraKeyspace(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, nil)
	testsuite.Require().NoError(err)

	// From step CassandraResources_MigrateCassandraKeyspaceToAutoscale
	fmt.Println("Call operation: CassandraResources_MigrateCassandraKeyspaceToAutoscale")
	cassandraResourcesClientMigrateCassandraKeyspaceToAutoscaleResponsePoller, err := cassandraResourcesClient.BeginMigrateCassandraKeyspaceToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientMigrateCassandraKeyspaceToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step CassandraResources_MigrateCassandraKeyspaceToManualThroughput
	fmt.Println("Call operation: CassandraResources_MigrateCassandraKeyspaceToManualThroughput")
	cassandraResourcesClientMigrateCassandraKeyspaceToManualThroughputResponsePoller, err := cassandraResourcesClient.BeginMigrateCassandraKeyspaceToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientMigrateCassandraKeyspaceToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step CassandraResources_UpdateCassandraKeyspaceThroughput
	fmt.Println("Call operation: CassandraResources_UpdateCassandraKeyspaceThroughput")
	cassandraResourcesClientUpdateCassandraKeyspaceThroughputResponsePoller, err := cassandraResourcesClient.BeginUpdateCassandraKeyspaceThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientUpdateCassandraKeyspaceThroughputResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/cassandraKeyspaces/{keyspaceName}/tables/{tableName}
func (testsuite *CassandraResourcesTestSuite) TestCassandraTable() {
	var err error
	// From step CassandraResources_CreateUpdateCassandraTable
	fmt.Println("Call operation: CassandraResources_CreateUpdateCassandraTable")
	cassandraResourcesClient, err := armcosmos.NewCassandraResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	cassandraResourcesClientCreateUpdateCassandraTableResponsePoller, err := cassandraResourcesClient.BeginCreateUpdateCassandraTable(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, testsuite.tableName, armcosmos.CassandraTableCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.CassandraTableCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{
				Throughput: to.Ptr[int32](2000),
			},
			Resource: &armcosmos.CassandraTableResource{
				Schema: &armcosmos.CassandraSchema{
					Columns: []*armcosmos.Column{
						{
							Name: to.Ptr("columnA"),
							Type: to.Ptr("Ascii"),
						}},
					PartitionKeys: []*armcosmos.CassandraPartitionKey{
						{
							Name: to.Ptr("columnA"),
						}},
				},
				DefaultTTL: to.Ptr[int32](100),
				ID:         to.Ptr(testsuite.tableName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientCreateUpdateCassandraTableResponsePoller)
	testsuite.Require().NoError(err)

	// From step CassandraResources_GetCassandraTable
	fmt.Println("Call operation: CassandraResources_GetCassandraTable")
	_, err = cassandraResourcesClient.GetCassandraTable(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)

	// From step CassandraResources_ListCassandraTables
	fmt.Println("Call operation: CassandraResources_ListCassandraTables")
	cassandraResourcesClientNewListCassandraTablesPager := cassandraResourcesClient.NewListCassandraTablesPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, nil)
	for cassandraResourcesClientNewListCassandraTablesPager.More() {
		_, err := cassandraResourcesClientNewListCassandraTablesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step CassandraResources_GetCassandraTableThroughput
	fmt.Println("Call operation: CassandraResources_GetCassandraTableThroughput")
	_, err = cassandraResourcesClient.GetCassandraTableThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)

	// From step CassandraResources_MigrateCassandraTableToAutoscale
	fmt.Println("Call operation: CassandraResources_MigrateCassandraTableToAutoscale")
	cassandraResourcesClientMigrateCassandraTableToAutoscaleResponsePoller, err := cassandraResourcesClient.BeginMigrateCassandraTableToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientMigrateCassandraTableToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step CassandraResources_MigrateCassandraTableToManualThroughput
	fmt.Println("Call operation: CassandraResources_MigrateCassandraTableToManualThroughput")
	cassandraResourcesClientMigrateCassandraTableToManualThroughputResponsePoller, err := cassandraResourcesClient.BeginMigrateCassandraTableToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientMigrateCassandraTableToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step CassandraResources_UpdateCassandraTableThroughput
	fmt.Println("Call operation: CassandraResources_UpdateCassandraTableThroughput")
	cassandraResourcesClientUpdateCassandraTableThroughputResponsePoller, err := cassandraResourcesClient.BeginUpdateCassandraTableThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, testsuite.tableName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientUpdateCassandraTableThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step CassandraResources_DeleteCassandraTable
	fmt.Println("Call operation: CassandraResources_DeleteCassandraTable")
	cassandraResourcesClientDeleteCassandraTableResponsePoller, err := cassandraResourcesClient.BeginDeleteCassandraTable(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientDeleteCassandraTableResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *CassandraResourcesTestSuite) Cleanup() {
	var err error
	// From step CassandraResources_DeleteCassandraKeyspace
	fmt.Println("Call operation: CassandraResources_DeleteCassandraKeyspace")
	cassandraResourcesClient, err := armcosmos.NewCassandraResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	cassandraResourcesClientDeleteCassandraKeyspaceResponsePoller, err := cassandraResourcesClient.BeginDeleteCassandraKeyspace(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.keyspaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cassandraResourcesClientDeleteCassandraKeyspaceResponsePoller)
	testsuite.Require().NoError(err)
}

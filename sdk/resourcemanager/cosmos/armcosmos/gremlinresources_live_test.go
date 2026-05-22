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

type GremlinResourcesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	accountName       string
	databaseName      string
	graphName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *GremlinResourcesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.databaseName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "gremlindb", 15, false)
	testsuite.graphName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "graphnam", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *GremlinResourcesTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestGremlinResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(GremlinResourcesTestSuite))
}

func (testsuite *GremlinResourcesTestSuite) Prepare() {
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
					Name: to.Ptr("EnableGremlin"),
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

	// From step GremlinResources_CreateUpdateGremlinDatabase
	fmt.Println("Call operation: GremlinResources_CreateUpdateGremlinDatabase")
	gremlinResourcesClient, err := armcosmos.NewGremlinResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	gremlinResourcesClientCreateUpdateGremlinDatabaseResponsePoller, err := gremlinResourcesClient.BeginCreateUpdateGremlinDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, armcosmos.GremlinDatabaseCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.GremlinDatabaseCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{
				Throughput: to.Ptr[int32](2000),
			},
			Resource: &armcosmos.GremlinDatabaseResource{
				ID: to.Ptr(testsuite.databaseName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientCreateUpdateGremlinDatabaseResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/gremlinDatabases/{databaseName}
func (testsuite *GremlinResourcesTestSuite) TestGremlinDatabase() {
	var err error
	// From step GremlinResources_ListGremlinDatabases
	fmt.Println("Call operation: GremlinResources_ListGremlinDatabases")
	gremlinResourcesClient, err := armcosmos.NewGremlinResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	gremlinResourcesClientNewListGremlinDatabasesPager := gremlinResourcesClient.NewListGremlinDatabasesPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for gremlinResourcesClientNewListGremlinDatabasesPager.More() {
		_, err := gremlinResourcesClientNewListGremlinDatabasesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step GremlinResources_GetGremlinDatabaseThroughput
	fmt.Println("Call operation: GremlinResources_GetGremlinDatabaseThroughput")
	_, err = gremlinResourcesClient.GetGremlinDatabaseThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step GremlinResources_GetGremlinDatabase
	fmt.Println("Call operation: GremlinResources_GetGremlinDatabase")
	_, err = gremlinResourcesClient.GetGremlinDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step GremlinResources_MigrateGremlinDatabaseToAutoscale
	fmt.Println("Call operation: GremlinResources_MigrateGremlinDatabaseToAutoscale")
	gremlinResourcesClientMigrateGremlinDatabaseToAutoscaleResponsePoller, err := gremlinResourcesClient.BeginMigrateGremlinDatabaseToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientMigrateGremlinDatabaseToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step GremlinResources_MigrateGremlinDatabaseToManualThroughput
	fmt.Println("Call operation: GremlinResources_MigrateGremlinDatabaseToManualThroughput")
	gremlinResourcesClientMigrateGremlinDatabaseToManualThroughputResponsePoller, err := gremlinResourcesClient.BeginMigrateGremlinDatabaseToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientMigrateGremlinDatabaseToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step GremlinResources_UpdateGremlinDatabaseThroughput
	fmt.Println("Call operation: GremlinResources_UpdateGremlinDatabaseThroughput")
	gremlinResourcesClientUpdateGremlinDatabaseThroughputResponsePoller, err := gremlinResourcesClient.BeginUpdateGremlinDatabaseThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientUpdateGremlinDatabaseThroughputResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/gremlinDatabases/{databaseName}/graphs/{graphName}
func (testsuite *GremlinResourcesTestSuite) TestGremlinGraph() {
	var err error
	// From step GremlinResources_CreateUpdateGremlinGraph
	fmt.Println("Call operation: GremlinResources_CreateUpdateGremlinGraph")
	gremlinResourcesClient, err := armcosmos.NewGremlinResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	gremlinResourcesClientCreateUpdateGremlinGraphResponsePoller, err := gremlinResourcesClient.BeginCreateUpdateGremlinGraph(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.graphName, armcosmos.GremlinGraphCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.GremlinGraphCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{
				Throughput: to.Ptr[int32](2000),
			},
			Resource: &armcosmos.GremlinGraphResource{
				ConflictResolutionPolicy: &armcosmos.ConflictResolutionPolicy{
					ConflictResolutionPath: to.Ptr("/path"),
					Mode:                   to.Ptr(armcosmos.ConflictResolutionModeLastWriterWins),
				},
				DefaultTTL: to.Ptr[int32](100),
				ID:         to.Ptr(testsuite.graphName),
				IndexingPolicy: &armcosmos.IndexingPolicy{
					Automatic:     to.Ptr(true),
					ExcludedPaths: []*armcosmos.ExcludedPath{},
					IncludedPaths: []*armcosmos.IncludedPath{
						{
							Path: to.Ptr("/*"),
							Indexes: []*armcosmos.Indexes{
								{
									DataType:  to.Ptr(armcosmos.DataTypeString),
									Kind:      to.Ptr(armcosmos.IndexKindRange),
									Precision: to.Ptr[int32](-1),
								},
								{
									DataType:  to.Ptr(armcosmos.DataTypeNumber),
									Kind:      to.Ptr(armcosmos.IndexKindRange),
									Precision: to.Ptr[int32](-1),
								}},
						}},
					IndexingMode: to.Ptr(armcosmos.IndexingModeConsistent),
				},
				PartitionKey: &armcosmos.ContainerPartitionKey{
					Kind: to.Ptr(armcosmos.PartitionKindHash),
					Paths: []*string{
						to.Ptr("/AccountNumber")},
				},
				UniqueKeyPolicy: &armcosmos.UniqueKeyPolicy{
					UniqueKeys: []*armcosmos.UniqueKey{
						{
							Paths: []*string{
								to.Ptr("/testPath")},
						}},
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientCreateUpdateGremlinGraphResponsePoller)
	testsuite.Require().NoError(err)

	// From step GremlinResources_GetGremlinGraph
	fmt.Println("Call operation: GremlinResources_GetGremlinGraph")
	_, err = gremlinResourcesClient.GetGremlinGraph(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.graphName, nil)
	testsuite.Require().NoError(err)

	// From step GremlinResources_GetGremlinGraphThroughput
	fmt.Println("Call operation: GremlinResources_GetGremlinGraphThroughput")
	_, err = gremlinResourcesClient.GetGremlinGraphThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.graphName, nil)
	testsuite.Require().NoError(err)

	// From step GremlinResources_ListGremlinGraphs
	fmt.Println("Call operation: GremlinResources_ListGremlinGraphs")
	gremlinResourcesClientNewListGremlinGraphsPager := gremlinResourcesClient.NewListGremlinGraphsPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	for gremlinResourcesClientNewListGremlinGraphsPager.More() {
		_, err := gremlinResourcesClientNewListGremlinGraphsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step GremlinResources_MigrateGremlinGraphToAutoscale
	fmt.Println("Call operation: GremlinResources_MigrateGremlinGraphToAutoscale")
	gremlinResourcesClientMigrateGremlinGraphToAutoscaleResponsePoller, err := gremlinResourcesClient.BeginMigrateGremlinGraphToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.graphName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientMigrateGremlinGraphToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step GremlinResources_MigrateGremlinGraphToManualThroughput
	fmt.Println("Call operation: GremlinResources_MigrateGremlinGraphToManualThroughput")
	gremlinResourcesClientMigrateGremlinGraphToManualThroughputResponsePoller, err := gremlinResourcesClient.BeginMigrateGremlinGraphToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.graphName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientMigrateGremlinGraphToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step GremlinResources_UpdateGremlinGraphThroughput
	fmt.Println("Call operation: GremlinResources_UpdateGremlinGraphThroughput")
	gremlinResourcesClientUpdateGremlinGraphThroughputResponsePoller, err := gremlinResourcesClient.BeginUpdateGremlinGraphThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.graphName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientUpdateGremlinGraphThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step GremlinResources_DeleteGremlinGraph
	fmt.Println("Call operation: GremlinResources_DeleteGremlinGraph")
	gremlinResourcesClientDeleteGremlinGraphResponsePoller, err := gremlinResourcesClient.BeginDeleteGremlinGraph(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.graphName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientDeleteGremlinGraphResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *GremlinResourcesTestSuite) Cleanup() {
	var err error
	// From step GremlinResources_DeleteGremlinDatabase
	fmt.Println("Call operation: GremlinResources_DeleteGremlinDatabase")
	gremlinResourcesClient, err := armcosmos.NewGremlinResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	gremlinResourcesClientDeleteGremlinDatabaseResponsePoller, err := gremlinResourcesClient.BeginDeleteGremlinDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gremlinResourcesClientDeleteGremlinDatabaseResponsePoller)
	testsuite.Require().NoError(err)
}

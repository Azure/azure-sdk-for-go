//go:build go1.18
// +build go1.18

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

type MongoDbResourcesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	accountName       string
	collectionName    string
	databaseName      string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *MongoDbResourcesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.collectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "collecti", 14, false)
	testsuite.databaseName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mongodb", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *MongoDbResourcesTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestMongoDbResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(MongoDbResourcesTestSuite))
}

func (testsuite *MongoDbResourcesTestSuite) Prepare() {
	var err error
	// From step DatabaseAccounts_CreateOrUpdate
	fmt.Println("Call operation: DatabaseAccounts_CreateOrUpdate")
	databaseAccountsClient, err := armcosmos.NewDatabaseAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databaseAccountsClientCreateOrUpdateResponsePoller, err := databaseAccountsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.DatabaseAccountCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Kind:     to.Ptr(armcosmos.DatabaseAccountKindMongoDB),
		Properties: &armcosmos.DatabaseAccountCreateUpdateProperties{
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

	// From step MongoDBResources_CreateUpdateMongoDBDatabase
	fmt.Println("Call operation: MongoDBResources_CreateUpdateMongoDBDatabase")
	mongoDBResourcesClient, err := armcosmos.NewMongoDBResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	mongoDBResourcesClientCreateUpdateMongoDBDatabaseResponsePoller, err := mongoDBResourcesClient.BeginCreateUpdateMongoDBDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, armcosmos.MongoDBDatabaseCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.MongoDBDatabaseCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{
				Throughput: to.Ptr[int32](2000),
			},
			Resource: &armcosmos.MongoDBDatabaseResource{
				ID: to.Ptr(testsuite.databaseName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientCreateUpdateMongoDBDatabaseResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/mongodbDatabases/{databaseName}
func (testsuite *MongoDbResourcesTestSuite) TestMongoDbDatabase() {
	var err error
	// From step MongoDBResources_ListMongoDBDatabases
	fmt.Println("Call operation: MongoDBResources_ListMongoDBDatabases")
	mongoDBResourcesClient, err := armcosmos.NewMongoDBResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	mongoDBResourcesClientNewListMongoDBDatabasesPager := mongoDBResourcesClient.NewListMongoDBDatabasesPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for mongoDBResourcesClientNewListMongoDBDatabasesPager.More() {
		_, err := mongoDBResourcesClientNewListMongoDBDatabasesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MongoDBResources_GetMongoDBDatabase
	fmt.Println("Call operation: MongoDBResources_GetMongoDBDatabase")
	_, err = mongoDBResourcesClient.GetMongoDBDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_GetMongoDBDatabaseThroughput
	fmt.Println("Call operation: MongoDBResources_GetMongoDBDatabaseThroughput")
	_, err = mongoDBResourcesClient.GetMongoDBDatabaseThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_MigrateMongoDBDatabaseToAutoscale
	fmt.Println("Call operation: MongoDBResources_MigrateMongoDBDatabaseToAutoscale")
	mongoDBResourcesClientMigrateMongoDBDatabaseToAutoscaleResponsePoller, err := mongoDBResourcesClient.BeginMigrateMongoDBDatabaseToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientMigrateMongoDBDatabaseToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_MigrateMongoDBDatabaseToManualThroughput
	fmt.Println("Call operation: MongoDBResources_MigrateMongoDBDatabaseToManualThroughput")
	mongoDBResourcesClientMigrateMongoDBDatabaseToManualThroughputResponsePoller, err := mongoDBResourcesClient.BeginMigrateMongoDBDatabaseToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientMigrateMongoDBDatabaseToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_UpdateMongoDBDatabaseThroughput
	fmt.Println("Call operation: MongoDBResources_UpdateMongoDBDatabaseThroughput")
	mongoDBResourcesClientUpdateMongoDBDatabaseThroughputResponsePoller, err := mongoDBResourcesClient.BeginUpdateMongoDBDatabaseThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientUpdateMongoDBDatabaseThroughputResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/mongodbDatabases/{databaseName}/collections/{collectionName}
func (testsuite *MongoDbResourcesTestSuite) TestMongoDbCollection() {
	var err error
	// From step MongoDBResources_CreateUpdateMongoDBCollection
	fmt.Println("Call operation: MongoDBResources_CreateUpdateMongoDBCollection")
	mongoDBResourcesClient, err := armcosmos.NewMongoDBResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	mongoDBResourcesClientCreateUpdateMongoDBCollectionResponsePoller, err := mongoDBResourcesClient.BeginCreateUpdateMongoDBCollection(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.collectionName, armcosmos.MongoDBCollectionCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.MongoDBCollectionCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{
				Throughput: to.Ptr[int32](2000),
			},
			Resource: &armcosmos.MongoDBCollectionResource{
				ID: to.Ptr(testsuite.collectionName),
				ShardKey: map[string]*string{
					"testKey": to.Ptr("Hash"),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientCreateUpdateMongoDBCollectionResponsePoller)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_ListMongoDBCollections
	fmt.Println("Call operation: MongoDBResources_ListMongoDBCollections")
	mongoDBResourcesClientNewListMongoDBCollectionsPager := mongoDBResourcesClient.NewListMongoDBCollectionsPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	for mongoDBResourcesClientNewListMongoDBCollectionsPager.More() {
		_, err := mongoDBResourcesClientNewListMongoDBCollectionsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MongoDBResources_GetMongoDBCollection
	fmt.Println("Call operation: MongoDBResources_GetMongoDBCollection")
	_, err = mongoDBResourcesClient.GetMongoDBCollection(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.collectionName, nil)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_GetMongoDBCollectionThroughput
	fmt.Println("Call operation: MongoDBResources_GetMongoDBCollectionThroughput")
	_, err = mongoDBResourcesClient.GetMongoDBCollectionThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.collectionName, nil)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_MigrateMongoDBCollectionToAutoscale
	fmt.Println("Call operation: MongoDBResources_MigrateMongoDBCollectionToAutoscale")
	mongoDBResourcesClientMigrateMongoDBCollectionToAutoscaleResponsePoller, err := mongoDBResourcesClient.BeginMigrateMongoDBCollectionToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.collectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientMigrateMongoDBCollectionToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_MigrateMongoDBCollectionToManualThroughput
	fmt.Println("Call operation: MongoDBResources_MigrateMongoDBCollectionToManualThroughput")
	mongoDBResourcesClientMigrateMongoDBCollectionToManualThroughputResponsePoller, err := mongoDBResourcesClient.BeginMigrateMongoDBCollectionToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.collectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientMigrateMongoDBCollectionToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_UpdateMongoDBCollectionThroughput
	fmt.Println("Call operation: MongoDBResources_UpdateMongoDBCollectionThroughput")
	mongoDBResourcesClientUpdateMongoDBCollectionThroughputResponsePoller, err := mongoDBResourcesClient.BeginUpdateMongoDBCollectionThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.collectionName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientUpdateMongoDBCollectionThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step MongoDBResources_DeleteMongoDBCollection
	fmt.Println("Call operation: MongoDBResources_DeleteMongoDBCollection")
	mongoDBResourcesClientDeleteMongoDBCollectionResponsePoller, err := mongoDBResourcesClient.BeginDeleteMongoDBCollection(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.collectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientDeleteMongoDBCollectionResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *MongoDbResourcesTestSuite) Cleanup() {
	var err error
	// From step MongoDBResources_DeleteMongoDBDatabase
	fmt.Println("Call operation: MongoDBResources_DeleteMongoDBDatabase")
	mongoDBResourcesClient, err := armcosmos.NewMongoDBResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	mongoDBResourcesClientDeleteMongoDBDatabaseResponsePoller, err := mongoDBResourcesClient.BeginDeleteMongoDBDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mongoDBResourcesClientDeleteMongoDBDatabaseResponsePoller)
	testsuite.Require().NoError(err)
}

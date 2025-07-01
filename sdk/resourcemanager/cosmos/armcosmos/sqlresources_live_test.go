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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type SqlResourcesTestSuite struct {
	suite.Suite

	ctx                     context.Context
	cred                    azcore.TokenCredential
	options                 *arm.ClientOptions
	accountName             string
	containerName           string
	databaseName            string
	triggerName             string
	storedProcedureName     string
	userDefinedFunctionName string
	clientEncryptionKeyName string
	location                string
	resourceGroupName       string
	subscriptionId          string
}

func (testsuite *SqlResourcesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.containerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sqlcontaine", 17, false)
	testsuite.databaseName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sqldb", 11, false)
	testsuite.triggerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "triggern", 14, false)
	testsuite.storedProcedureName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storedpr", 14, true)
	testsuite.userDefinedFunctionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "userdefi", 14, false)
	testsuite.clientEncryptionKeyName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clienten", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SqlResourcesTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestSqlResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(SqlResourcesTestSuite))
}

func (testsuite *SqlResourcesTestSuite) Prepare() {
	var err error
	// From step DatabaseAccounts_CreateOrUpdate
	fmt.Println("Call operation: DatabaseAccounts_CreateOrUpdate")
	databaseAccountsClient, err := armcosmos.NewDatabaseAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databaseAccountsClientCreateOrUpdateResponsePoller, err := databaseAccountsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.DatabaseAccountCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
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

	// From step SqlResources_CreateUpdateSqlDatabase
	fmt.Println("Call operation: SQLResources_CreateUpdateSQLDatabase")
	sQLResourcesClient, err := armcosmos.NewSQLResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sQLResourcesClientCreateUpdateSQLDatabaseResponsePoller, err := sQLResourcesClient.BeginCreateUpdateSQLDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, armcosmos.SQLDatabaseCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.SQLDatabaseCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{
				Throughput: to.Ptr[int32](2000),
			},
			Resource: &armcosmos.SQLDatabaseResource{
				ID: to.Ptr(testsuite.databaseName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientCreateUpdateSQLDatabaseResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_CreateUpdateSqlContainer
	fmt.Println("Call operation: SQLResources_CreateUpdateSQLContainer")
	sQLResourcesClientCreateUpdateSQLContainerResponsePoller, err := sQLResourcesClient.BeginCreateUpdateSQLContainer(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, armcosmos.SQLContainerCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.SQLContainerCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{
				Throughput: to.Ptr[int32](2000),
			},
			Resource: &armcosmos.SQLContainerResource{
				ConflictResolutionPolicy: &armcosmos.ConflictResolutionPolicy{
					ConflictResolutionPath: to.Ptr("/path"),
					Mode:                   to.Ptr(armcosmos.ConflictResolutionModeLastWriterWins),
				},
				ID: to.Ptr(testsuite.containerName),
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
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientCreateUpdateSQLContainerResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/sqlDatabases/{databaseName}
func (testsuite *SqlResourcesTestSuite) TestSqlDatabase() {
	var err error
	// From step SqlResources_ListSqlDatabases
	fmt.Println("Call operation: SQLResources_ListSQLDatabases")
	sQLResourcesClient, err := armcosmos.NewSQLResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sQLResourcesClientNewListSQLDatabasesPager := sQLResourcesClient.NewListSQLDatabasesPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for sQLResourcesClientNewListSQLDatabasesPager.More() {
		_, err := sQLResourcesClientNewListSQLDatabasesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SqlResources_GetSqlDatabase
	fmt.Println("Call operation: SQLResources_GetSQLDatabase")
	_, err = sQLResourcesClient.GetSQLDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step SqlResources_GetSqlDatabaseThroughput
	fmt.Println("Call operation: SQLResources_GetSQLDatabaseThroughput")
	_, err = sQLResourcesClient.GetSQLDatabaseThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step SqlResources_UpdateSqlDatabaseThroughput
	fmt.Println("Call operation: SQLResources_UpdateSQLDatabaseThroughput")
	sQLResourcesClientUpdateSQLDatabaseThroughputResponsePoller, err := sQLResourcesClient.BeginUpdateSQLDatabaseThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientUpdateSQLDatabaseThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_MigrateSqlDatabaseToAutoscale
	fmt.Println("Call operation: SQLResources_MigrateSQLDatabaseToAutoscale")
	sQLResourcesClientMigrateSQLDatabaseToAutoscaleResponsePoller, err := sQLResourcesClient.BeginMigrateSQLDatabaseToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientMigrateSQLDatabaseToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_MigrateSqlDatabaseToManualThroughput
	fmt.Println("Call operation: SQLResources_MigrateSQLDatabaseToManualThroughput")
	sQLResourcesClientMigrateSQLDatabaseToManualThroughputResponsePoller, err := sQLResourcesClient.BeginMigrateSQLDatabaseToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientMigrateSQLDatabaseToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/sqlDatabases/{databaseName}/containers/{containerName}
func (testsuite *SqlResourcesTestSuite) TestSqlContainer() {
	var err error
	// From step SqlResources_GetSqlContainer
	fmt.Println("Call operation: SQLResources_GetSQLContainer")
	sQLResourcesClient, err := armcosmos.NewSQLResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = sQLResourcesClient.GetSQLContainer(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, nil)
	testsuite.Require().NoError(err)

	// From step SqlResources_ListSqlContainers
	fmt.Println("Call operation: SQLResources_ListSQLContainers")
	sQLResourcesClientNewListSQLContainersPager := sQLResourcesClient.NewListSQLContainersPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	for sQLResourcesClientNewListSQLContainersPager.More() {
		_, err := sQLResourcesClientNewListSQLContainersPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SqlResources_GetSqlContainerThroughput
	fmt.Println("Call operation: SQLResources_GetSQLContainerThroughput")
	_, err = sQLResourcesClient.GetSQLContainerThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, nil)
	testsuite.Require().NoError(err)

	// From step SqlResources_UpdateSqlContainerThroughput
	fmt.Println("Call operation: SQLResources_UpdateSQLContainerThroughput")
	sQLResourcesClientUpdateSQLContainerThroughputResponsePoller, err := sQLResourcesClient.BeginUpdateSQLContainerThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientUpdateSQLContainerThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_MigrateSqlContainerToAutoscale
	fmt.Println("Call operation: SQLResources_MigrateSQLContainerToAutoscale")
	sQLResourcesClientMigrateSQLContainerToAutoscaleResponsePoller, err := sQLResourcesClient.BeginMigrateSQLContainerToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientMigrateSQLContainerToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_MigrateSqlContainerToManualThroughput
	fmt.Println("Call operation: SQLResources_MigrateSQLContainerToManualThroughput")
	sQLResourcesClientMigrateSQLContainerToManualThroughputResponsePoller, err := sQLResourcesClient.BeginMigrateSQLContainerToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientMigrateSQLContainerToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/sqlDatabases/{databaseName}/containers/{containerName}/triggers/{triggerName}
func (testsuite *SqlResourcesTestSuite) TestSqlTrigger() {
	var err error
	// From step SqlResources_CreateUpdateSqlTrigger
	fmt.Println("Call operation: SQLResources_CreateUpdateSQLTrigger")
	sQLResourcesClient, err := armcosmos.NewSQLResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sQLResourcesClientCreateUpdateSQLTriggerResponsePoller, err := sQLResourcesClient.BeginCreateUpdateSQLTrigger(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.triggerName, armcosmos.SQLTriggerCreateUpdateParameters{
		Properties: &armcosmos.SQLTriggerCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{},
			Resource: &armcosmos.SQLTriggerResource{
				Body:             to.Ptr("body"),
				ID:               to.Ptr(testsuite.triggerName),
				TriggerOperation: to.Ptr(armcosmos.TriggerOperationAll),
				TriggerType:      to.Ptr(armcosmos.TriggerTypePre),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientCreateUpdateSQLTriggerResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_GetSqlTrigger
	fmt.Println("Call operation: SQLResources_GetSQLTrigger")
	_, err = sQLResourcesClient.GetSQLTrigger(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.triggerName, nil)
	testsuite.Require().NoError(err)

	// From step SqlResources_ListSqlTriggers
	fmt.Println("Call operation: SQLResources_ListSQLTriggers")
	sQLResourcesClientNewListSQLTriggersPager := sQLResourcesClient.NewListSQLTriggersPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, nil)
	for sQLResourcesClientNewListSQLTriggersPager.More() {
		_, err := sQLResourcesClientNewListSQLTriggersPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SqlResources_DeleteSqlTrigger
	fmt.Println("Call operation: SQLResources_DeleteSQLTrigger")
	sQLResourcesClientDeleteSQLTriggerResponsePoller, err := sQLResourcesClient.BeginDeleteSQLTrigger(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.triggerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientDeleteSQLTriggerResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/sqlDatabases/{databaseName}/containers/{containerName}/storedProcedures/{storedProcedureName}
func (testsuite *SqlResourcesTestSuite) TestSqlStoredProcedure() {
	var err error
	// From step SqlResources_CreateUpdateSqlStoredProcedure
	fmt.Println("Call operation: SQLResources_CreateUpdateSQLStoredProcedure")
	sQLResourcesClient, err := armcosmos.NewSQLResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sQLResourcesClientCreateUpdateSQLStoredProcedureResponsePoller, err := sQLResourcesClient.BeginCreateUpdateSQLStoredProcedure(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.storedProcedureName, armcosmos.SQLStoredProcedureCreateUpdateParameters{
		Properties: &armcosmos.SQLStoredProcedureCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{},
			Resource: &armcosmos.SQLStoredProcedureResource{
				Body: to.Ptr("body"),
				ID:   to.Ptr(testsuite.storedProcedureName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientCreateUpdateSQLStoredProcedureResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_GetSqlStoredProcedure
	fmt.Println("Call operation: SQLResources_GetSQLStoredProcedure")
	_, err = sQLResourcesClient.GetSQLStoredProcedure(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.storedProcedureName, nil)
	testsuite.Require().NoError(err)

	// From step SqlResources_ListSqlStoredProcedures
	fmt.Println("Call operation: SQLResources_ListSQLStoredProcedures")
	sQLResourcesClientNewListSQLStoredProceduresPager := sQLResourcesClient.NewListSQLStoredProceduresPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, nil)
	for sQLResourcesClientNewListSQLStoredProceduresPager.More() {
		_, err := sQLResourcesClientNewListSQLStoredProceduresPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SqlResources_DeleteSqlStoredProcedure
	fmt.Println("Call operation: SQLResources_DeleteSQLStoredProcedure")
	sQLResourcesClientDeleteSQLStoredProcedureResponsePoller, err := sQLResourcesClient.BeginDeleteSQLStoredProcedure(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.storedProcedureName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientDeleteSQLStoredProcedureResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/sqlDatabases/{databaseName}/containers/{containerName}/userDefinedFunctions/{userDefinedFunctionName}
func (testsuite *SqlResourcesTestSuite) TestSqlUserDefinedFunction() {
	var err error
	// From step SqlResources_CreateUpdateSqlUserDefinedFunction
	fmt.Println("Call operation: SQLResources_CreateUpdateSQLUserDefinedFunction")
	sQLResourcesClient, err := armcosmos.NewSQLResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sQLResourcesClientCreateUpdateSQLUserDefinedFunctionResponsePoller, err := sQLResourcesClient.BeginCreateUpdateSQLUserDefinedFunction(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.userDefinedFunctionName, armcosmos.SQLUserDefinedFunctionCreateUpdateParameters{
		Properties: &armcosmos.SQLUserDefinedFunctionCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{},
			Resource: &armcosmos.SQLUserDefinedFunctionResource{
				Body: to.Ptr("body"),
				ID:   to.Ptr(testsuite.userDefinedFunctionName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientCreateUpdateSQLUserDefinedFunctionResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_GetSqlUserDefinedFunction
	fmt.Println("Call operation: SQLResources_GetSQLUserDefinedFunction")
	_, err = sQLResourcesClient.GetSQLUserDefinedFunction(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.userDefinedFunctionName, nil)
	testsuite.Require().NoError(err)

	// From step SqlResources_ListSqlUserDefinedFunctions
	fmt.Println("Call operation: SQLResources_ListSQLUserDefinedFunctions")
	sQLResourcesClientNewListSQLUserDefinedFunctionsPager := sQLResourcesClient.NewListSQLUserDefinedFunctionsPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, nil)
	for sQLResourcesClientNewListSQLUserDefinedFunctionsPager.More() {
		_, err := sQLResourcesClientNewListSQLUserDefinedFunctionsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SqlResources_DeleteSqlUserDefinedFunction
	fmt.Println("Call operation: SQLResources_DeleteSQLUserDefinedFunction")
	sQLResourcesClientDeleteSQLUserDefinedFunctionResponsePoller, err := sQLResourcesClient.BeginDeleteSQLUserDefinedFunction(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, testsuite.userDefinedFunctionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientDeleteSQLUserDefinedFunctionResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}/sqlDatabases/{databaseName}/clientEncryptionKeys/{clientEncryptionKeyName}
func (testsuite *SqlResourcesTestSuite) TestSqlClientEncryptionKey() {
	var err error
	// From step SqlResources_CreateUpdateClientEncryptionKey
	fmt.Println("Call operation: SQLResources_CreateUpdateClientEncryptionKey")
	sQLResourcesClient, err := armcosmos.NewSQLResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sQLResourcesClientCreateUpdateClientEncryptionKeyResponsePoller, err := sQLResourcesClient.BeginCreateUpdateClientEncryptionKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.clientEncryptionKeyName, armcosmos.ClientEncryptionKeyCreateUpdateParameters{
		Properties: &armcosmos.ClientEncryptionKeyCreateUpdateProperties{
			Resource: &armcosmos.ClientEncryptionKeyResource{
				EncryptionAlgorithm: to.Ptr("AEAD_AES_256_CBC_HMAC_SHA256"),
				ID:                  to.Ptr(testsuite.clientEncryptionKeyName),
				KeyWrapMetadata: &armcosmos.KeyWrapMetadata{
					Name:      to.Ptr("customerManagedKey"),
					Type:      to.Ptr("AzureKeyVault"),
					Algorithm: to.Ptr("RSA-OAEP"),
					Value:     to.Ptr("AzureKeyVault Key URL"),
				},
				WrappedDataEncryptionKey: []byte("U3dhZ2dlciByb2Nrcw=="),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientCreateUpdateClientEncryptionKeyResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_GetClientEncryptionKey
	fmt.Println("Call operation: SQLResources_GetClientEncryptionKey")
	_, err = sQLResourcesClient.GetClientEncryptionKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.clientEncryptionKeyName, nil)
	testsuite.Require().NoError(err)

	// From step SqlResources_ListClientEncryptionKeys
	fmt.Println("Call operation: SQLResources_ListClientEncryptionKeys")
	sQLResourcesClientNewListClientEncryptionKeysPager := sQLResourcesClient.NewListClientEncryptionKeysPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	for sQLResourcesClientNewListClientEncryptionKeysPager.More() {
		_, err := sQLResourcesClientNewListClientEncryptionKeysPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *SqlResourcesTestSuite) Cleanup() {
	var err error
	// From step SqlResources_DeleteSqlContainer
	fmt.Println("Call operation: SQLResources_DeleteSQLContainer")
	sQLResourcesClient, err := armcosmos.NewSQLResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sQLResourcesClientDeleteSQLContainerResponsePoller, err := sQLResourcesClient.BeginDeleteSQLContainer(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, testsuite.containerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientDeleteSQLContainerResponsePoller)
	testsuite.Require().NoError(err)

	// From step SqlResources_DeleteSqlDatabase
	fmt.Println("Call operation: SQLResources_DeleteSQLDatabase")
	sQLResourcesClientDeleteSQLDatabaseResponsePoller, err := sQLResourcesClient.BeginDeleteSQLDatabase(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sQLResourcesClientDeleteSQLDatabaseResponsePoller)
	testsuite.Require().NoError(err)
}

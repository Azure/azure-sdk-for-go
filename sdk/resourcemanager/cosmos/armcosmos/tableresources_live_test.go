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

type TableResourcesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	accountName       string
	tableName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *TableResourcesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.tableName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "tablenam", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *TableResourcesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestTableResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(TableResourcesTestSuite))
}

func (testsuite *TableResourcesTestSuite) Prepare() {
	var err error
	// From step DatabaseAccounts_CreateOrUpdate
	fmt.Println("Call operation: DatabaseAccounts_CreateOrUpdate")
	databaseAccountsClient, err := armcosmos.NewDatabaseAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databaseAccountsClientCreateOrUpdateResponsePoller, err := databaseAccountsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.DatabaseAccountCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Kind:     to.Ptr(armcosmos.DatabaseAccountKindGlobalDocumentDB),
		Properties: &armcosmos.DatabaseAccountCreateUpdateProperties{
			Capabilities: []*armcosmos.Capability{
				{
					Name: to.Ptr("EnableTable"),
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

// Microsoft.DocumentDB/databaseAccounts/{accountName}/tables/{tableName}
func (testsuite *TableResourcesTestSuite) TestTableResources() {
	var err error
	// From step TableResources_CreateUpdateTable
	fmt.Println("Call operation: TableResources_CreateUpdateTable")
	tableResourcesClient, err := armcosmos.NewTableResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tableResourcesClientCreateUpdateTableResponsePoller, err := tableResourcesClient.BeginCreateUpdateTable(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.tableName, armcosmos.TableCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.TableCreateUpdateProperties{
			Options: &armcosmos.CreateUpdateOptions{},
			Resource: &armcosmos.TableResource{
				ID: to.Ptr(testsuite.tableName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tableResourcesClientCreateUpdateTableResponsePoller)
	testsuite.Require().NoError(err)

	// From step TableResources_ListTables
	fmt.Println("Call operation: TableResources_ListTables")
	tableResourcesClientNewListTablesPager := tableResourcesClient.NewListTablesPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for tableResourcesClientNewListTablesPager.More() {
		_, err := tableResourcesClientNewListTablesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step TableResources_GetTableThroughput
	fmt.Println("Call operation: TableResources_GetTableThroughput")
	_, err = tableResourcesClient.GetTableThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)

	// From step TableResources_GetTable
	fmt.Println("Call operation: TableResources_GetTable")
	_, err = tableResourcesClient.GetTable(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)

	// From step TableResources_MigrateTableToAutoscale
	fmt.Println("Call operation: TableResources_MigrateTableToAutoscale")
	tableResourcesClientMigrateTableToAutoscaleResponsePoller, err := tableResourcesClient.BeginMigrateTableToAutoscale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tableResourcesClientMigrateTableToAutoscaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step TableResources_MigrateTableToManualThroughput
	fmt.Println("Call operation: TableResources_MigrateTableToManualThroughput")
	tableResourcesClientMigrateTableToManualThroughputResponsePoller, err := tableResourcesClient.BeginMigrateTableToManualThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tableResourcesClientMigrateTableToManualThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step TableResources_UpdateTableThroughput
	fmt.Println("Call operation: TableResources_UpdateTableThroughput")
	tableResourcesClientUpdateTableThroughputResponsePoller, err := tableResourcesClient.BeginUpdateTableThroughput(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.tableName, armcosmos.ThroughputSettingsUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmos.ThroughputSettingsUpdateProperties{
			Resource: &armcosmos.ThroughputSettingsResource{
				Throughput: to.Ptr[int32](400),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tableResourcesClientUpdateTableThroughputResponsePoller)
	testsuite.Require().NoError(err)

	// From step TableResources_DeleteTable
	fmt.Println("Call operation: TableResources_DeleteTable")
	tableResourcesClientDeleteTableResponsePoller, err := tableResourcesClient.BeginDeleteTable(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.tableName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tableResourcesClientDeleteTableResponsePoller)
	testsuite.Require().NoError(err)
}

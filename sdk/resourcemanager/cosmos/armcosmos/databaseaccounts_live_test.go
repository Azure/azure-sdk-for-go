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

type DatabaseAccountsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	accountName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *DatabaseAccountsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *DatabaseAccountsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestDatabaseAccountsTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseAccountsTestSuite))
}

// Microsoft.DocumentDB/databaseAccounts/{accountName}
func (testsuite *DatabaseAccountsTestSuite) TestDatabaseAccounts() {
	var err error
	// From step DatabaseAccounts_CheckNameExists
	fmt.Println("Call operation: DatabaseAccounts_CheckNameExists")
	databaseAccountsClient, err := armcosmos.NewDatabaseAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = databaseAccountsClient.CheckNameExists(testsuite.ctx, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// fmt.Println("Call operation: DatabaseAccounts_CreateOrUpdate")
	_, err = databaseAccountsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.DatabaseAccountCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Properties: &armcosmos.DatabaseAccountCreateUpdateProperties{
			CreateMode:               to.Ptr(armcosmos.CreateModeDefault),
			DatabaseAccountOfferType: to.Ptr("Standard"),
			Locations: []*armcosmos.Location{
				{
					FailoverPriority: to.Ptr[int32](0),
					LocationName:     to.Ptr("westus"),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_List
	fmt.Println("Call operation: DatabaseAccounts_List")
	databaseAccountsClientNewListPager := databaseAccountsClient.NewListPager(nil)
	for databaseAccountsClientNewListPager.More() {
		_, err := databaseAccountsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DatabaseAccounts_ListByResourceGroup
	fmt.Println("Call operation: DatabaseAccounts_ListByResourceGroup")
	databaseAccountsClientNewListByResourceGroupPager := databaseAccountsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for databaseAccountsClientNewListByResourceGroupPager.More() {
		_, err := databaseAccountsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DatabaseAccounts_Get
	fmt.Println("Call operation: DatabaseAccounts_Get")
	_, err = databaseAccountsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step RestorableDatabaseAccounts_ListByLocation
	fmt.Println("Call operation: RestorableDatabaseAccounts_ListByLocation")
	restorableDatabaseAccountsClient, err := armcosmos.NewRestorableDatabaseAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	restorableDatabaseAccountsClientNewListByLocationPager := restorableDatabaseAccountsClient.NewListByLocationPager(testsuite.location, nil)
	for restorableDatabaseAccountsClientNewListByLocationPager.More() {
		_, err := restorableDatabaseAccountsClientNewListByLocationPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RestorableDatabaseAccounts_List
	fmt.Println("Call operation: RestorableDatabaseAccounts_List")
	restorableDatabaseAccountsClientNewListPager := restorableDatabaseAccountsClient.NewListPager(nil)
	for restorableDatabaseAccountsClientNewListPager.More() {
		_, err := restorableDatabaseAccountsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

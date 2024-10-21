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

func TestDatabaseAccountsTestSuite(t *testing.T) {
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

	// From step DatabaseAccounts_CreateOrUpdate
	fmt.Println("Call operation: DatabaseAccounts_CreateOrUpdate")
	databaseAccountsClientCreateOrUpdateResponsePoller, err := databaseAccountsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.DatabaseAccountCreateUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Properties: &armcosmos.DatabaseAccountCreateUpdateProperties{
			CreateMode:               to.Ptr(armcosmos.CreateModeDefault),
			DatabaseAccountOfferType: to.Ptr("Standard"),
			Locations: []*armcosmos.Location{
				{
					FailoverPriority: to.Ptr[int32](2),
					LocationName:     to.Ptr("southcentralus"),
				},
				{
					FailoverPriority: to.Ptr[int32](1),
					LocationName:     to.Ptr("eastus"),
				},
				{
					FailoverPriority: to.Ptr[int32](0),
					LocationName:     to.Ptr("westus"),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databaseAccountsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_List
	fmt.Println("Call operation: DatabaseAccounts_List")
	databaseAccountsClientNewListPager := databaseAccountsClient.NewListPager(nil)
	for databaseAccountsClientNewListPager.More() {
		_, err := databaseAccountsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DatabaseAccounts_ListUsages
	fmt.Println("Call operation: DatabaseAccounts_ListUsages")
	databaseAccountsClientNewListUsagesPager := databaseAccountsClient.NewListUsagesPager(testsuite.resourceGroupName, testsuite.accountName, &armcosmos.DatabaseAccountsClientListUsagesOptions{Filter: to.Ptr("")})
	for databaseAccountsClientNewListUsagesPager.More() {
		_, err := databaseAccountsClientNewListUsagesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DatabaseAccounts_GetReadOnlyKeys
	fmt.Println("Call operation: DatabaseAccounts_GetReadOnlyKeys")
	_, err = databaseAccountsClient.GetReadOnlyKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_ListByResourceGroup
	fmt.Println("Call operation: DatabaseAccounts_ListByResourceGroup")
	databaseAccountsClientNewListByResourceGroupPager := databaseAccountsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for databaseAccountsClientNewListByResourceGroupPager.More() {
		_, err := databaseAccountsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DatabaseAccounts_ListMetricDefinitions
	fmt.Println("Call operation: DatabaseAccounts_ListMetricDefinitions")
	databaseAccountsClientNewListMetricDefinitionsPager := databaseAccountsClient.NewListMetricDefinitionsPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for databaseAccountsClientNewListMetricDefinitionsPager.More() {
		_, err := databaseAccountsClientNewListMetricDefinitionsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DatabaseAccounts_Get
	fmt.Println("Call operation: DatabaseAccounts_Get")
	_, err = databaseAccountsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_Update
	fmt.Println("Call operation: DatabaseAccounts_Update")
	databaseAccountsClientUpdateResponsePoller, err := databaseAccountsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.DatabaseAccountUpdateParameters{
		Tags: map[string]*string{
			"dept": to.Ptr("finance"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databaseAccountsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_RegenerateKey
	fmt.Println("Call operation: DatabaseAccounts_RegenerateKey")
	databaseAccountsClientRegenerateKeyResponsePoller, err := databaseAccountsClient.BeginRegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.DatabaseAccountRegenerateKeyParameters{
		KeyKind: to.Ptr(armcosmos.KeyKindPrimary),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databaseAccountsClientRegenerateKeyResponsePoller)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_ListReadOnlyKeys
	fmt.Println("Call operation: DatabaseAccounts_ListReadOnlyKeys")
	_, err = databaseAccountsClient.ListReadOnlyKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_ListConnectionStrings
	fmt.Println("Call operation: DatabaseAccounts_ListConnectionStrings")
	_, err = databaseAccountsClient.ListConnectionStrings(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_ListKeys
	fmt.Println("Call operation: DatabaseAccounts_ListKeys")
	_, err = databaseAccountsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step DatabaseAccounts_FailoverPriorityChange
	fmt.Println("Call operation: DatabaseAccounts_FailoverPriorityChange")
	databaseAccountsClientFailoverPriorityChangeResponsePoller, err := databaseAccountsClient.BeginFailoverPriorityChange(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armcosmos.FailoverPolicies{
		FailoverPolicies: []*armcosmos.FailoverPolicy{
			{
				FailoverPriority: to.Ptr[int32](0),
				LocationName:     to.Ptr("eastus"),
			},
			{
				FailoverPriority: to.Ptr[int32](2),
				LocationName:     to.Ptr("southcentralus"),
			},
			{
				FailoverPriority: to.Ptr[int32](1),
				LocationName:     to.Ptr("westus"),
			}},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databaseAccountsClientFailoverPriorityChangeResponsePoller)
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

	// From step DatabaseAccounts_Delete
	fmt.Println("Call operation: DatabaseAccounts_Delete")
	databaseAccountsClientDeleteResponsePoller, err := databaseAccountsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databaseAccountsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

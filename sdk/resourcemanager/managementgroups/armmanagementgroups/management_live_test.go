//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmanagementgroups_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups/v2"
	"github.com/stretchr/testify/suite"
)

type ManagementTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	groupId			string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ManagementTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.groupId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "groupid", 13, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *ManagementTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestManagementTestSuite(t *testing.T) {
	suite.Run(t, new(ManagementTestSuite))
}

// Microsoft.Management/managementGroups/{groupId}
func (testsuite *ManagementTestSuite) TestManagementGroups() {
	var err error
	// From step CheckNameAvailability
	fmt.Println("Call operation: CheckNameAvailability")
	aPIClient, err := armmanagementgroups.NewAPIClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIClient.CheckNameAvailability(testsuite.ctx, armmanagementgroups.CheckNameAvailabilityRequest{
		Name:	to.Ptr(testsuite.groupId),
		Type:	to.Ptr("Microsoft.Management/managementGroups"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step ManagementGroups_CreateOrUpdate
	fmt.Println("Call operation: ManagementGroups_CreateOrUpdate")
	client, err := armmanagementgroups.NewClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientCreateOrUpdateResponsePoller, err := client.BeginCreateOrUpdate(testsuite.ctx, testsuite.groupId, armmanagementgroups.CreateManagementGroupRequest{
		Properties: &armmanagementgroups.CreateManagementGroupProperties{
			DisplayName: to.Ptr(testsuite.groupId),
		},
	}, &armmanagementgroups.ClientBeginCreateOrUpdateOptions{CacheControl: to.Ptr("no-cache")})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ManagementGroups_List
	fmt.Println("Call operation: ManagementGroups_List")
	clientNewListPager := client.NewListPager(&armmanagementgroups.ClientListOptions{CacheControl: to.Ptr("no-cache"),
		Skiptoken:	nil,
	})
	for clientNewListPager.More() {
		_, err := clientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ManagementGroups_GetDescendants
	fmt.Println("Call operation: ManagementGroups_GetDescendants")
	clientNewGetDescendantsPager := client.NewGetDescendantsPager(testsuite.groupId, &armmanagementgroups.ClientGetDescendantsOptions{Skiptoken: nil,
		Top:	nil,
	})
	for clientNewGetDescendantsPager.More() {
		_, err := clientNewGetDescendantsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ManagementGroups_Get
	fmt.Println("Call operation: ManagementGroups_Get")
	_, err = client.Get(testsuite.ctx, testsuite.groupId, &armmanagementgroups.ClientGetOptions{Expand: nil,
		Recurse:	nil,
		Filter:		nil,
		CacheControl:	to.Ptr("no-cache"),
	})
	testsuite.Require().NoError(err)

	// From step ManagementGroups_Update
	fmt.Println("Call operation: ManagementGroups_Update")
	_, err = client.Update(testsuite.ctx, testsuite.groupId, armmanagementgroups.PatchManagementGroupRequest{}, &armmanagementgroups.ClientUpdateOptions{CacheControl: to.Ptr("no-cache")})
	testsuite.Require().NoError(err)

	// From step ManagementGroups_Delete
	fmt.Println("Call operation: ManagementGroups_Delete")
	clientDeleteResponsePoller, err := client.BeginDelete(testsuite.ctx, testsuite.groupId, &armmanagementgroups.ClientBeginDeleteOptions{CacheControl: to.Ptr("no-cache")})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Management/getEntities
func (testsuite *ManagementTestSuite) TestEntities() {
	var err error
	// From step Entities_List
	fmt.Println("Call operation: Entities_List")
	entitiesClient, err := armmanagementgroups.NewEntitiesClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	entitiesClientNewListPager := entitiesClient.NewListPager(&armmanagementgroups.EntitiesClientListOptions{Skiptoken: nil,
		Skip:		nil,
		Top:		nil,
		Select:		nil,
		Search:		nil,
		Filter:		nil,
		View:		nil,
		GroupName:	nil,
		CacheControl:	nil,
	})
	for entitiesClientNewListPager.More() {
		_, err := entitiesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Management/operations
func (testsuite *ManagementTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armmanagementgroups.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Management/startTenantBackfill
func (testsuite *ManagementTestSuite) TestStartTenantBackfill() {
	var err error
	// From step StartTenantBackfill
	fmt.Println("Call operation: StartTenantBackfill")
	aPIClient, err := armmanagementgroups.NewAPIClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIClient.StartTenantBackfill(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Management/tenantBackfillStatus
func (testsuite *ManagementTestSuite) TestTenantBackfillStatus() {
	var err error
	// From step TenantBackfillStatus
	fmt.Println("Call operation: TenantBackfillStatus")
	aPIClient, err := armmanagementgroups.NewAPIClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIClient.TenantBackfillStatus(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}

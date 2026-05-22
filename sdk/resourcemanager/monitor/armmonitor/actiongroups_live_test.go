// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmonitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/stretchr/testify/suite"
)

type ActiongroupsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	actionGroupName   string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ActiongroupsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.actionGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "actiongroupna", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *ActiongroupsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestActiongroupsTestSuite(t *testing.T) {
	suite.Run(t, new(ActiongroupsTestSuite))
}

// Microsoft.Insights/actionGroups
func (testsuite *ActiongroupsTestSuite) TestActiongroups() {
	var err error
	// From step ActionGroups_Create
	fmt.Println("Call operation: ActionGroups_CreateOrUpdate")
	actionGroupsClient, err := armmonitor.NewActionGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = actionGroupsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.actionGroupName, armmonitor.ActionGroupResource{
		Location: to.Ptr("Global"),
		Properties: &armmonitor.ActionGroup{
			EmailReceivers: []*armmonitor.EmailReceiver{
				{
					Name:                 to.Ptr("John Doe's email"),
					EmailAddress:         to.Ptr("johndoe@eamil.com"),
					UseCommonAlertSchema: to.Ptr(false),
				}},
			Enabled:        to.Ptr(true),
			GroupShortName: to.Ptr("sample"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ActionGroups_Get
	fmt.Println("Call operation: ActionGroups_Get")
	_, err = actionGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.actionGroupName, nil)
	testsuite.Require().NoError(err)

	// From step ActionGroups_Update
	fmt.Println("Call operation: ActionGroups_Update")
	_, err = actionGroupsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.actionGroupName, armmonitor.ActionGroupPatchBody{
		Properties: &armmonitor.ActionGroupPatch{
			Enabled: to.Ptr(false),
		},
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ActionGroups_ListBySubscriptionId
	fmt.Println("Call operation: ActionGroups_ListBySubscriptionId")
	actionGroupsClientNewListBySubscriptionIDPager := actionGroupsClient.NewListBySubscriptionIDPager(nil)
	for actionGroupsClientNewListBySubscriptionIDPager.More() {
		_, err := actionGroupsClientNewListBySubscriptionIDPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ActionGroups_ListByResourceGroup
	fmt.Println("Call operation: ActionGroups_ListByResourceGroup")
	actionGroupsClientNewListByResourceGroupPager := actionGroupsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for actionGroupsClientNewListByResourceGroupPager.More() {
		_, err := actionGroupsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ActionGroups_Delete
	fmt.Println("Call operation: ActionGroups_Delete")
	_, err = actionGroupsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.actionGroupName, nil)
	testsuite.Require().NoError(err)
}

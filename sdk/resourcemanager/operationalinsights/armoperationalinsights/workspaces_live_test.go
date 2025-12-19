// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armoperationalinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	"github.com/stretchr/testify/suite"
)

type WorkspacesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	workspaceName     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *WorkspacesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.workspaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "oiautorestws", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *WorkspacesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestWorkspacesTestSuite(t *testing.T) {
	suite.Run(t, new(WorkspacesTestSuite))
}

// Microsoft.OperationalInsights/workspaces
func (testsuite *WorkspacesTestSuite) TestWorkspace() {
	var err error
	// From step Workspaces_CreateOrUpdate
	fmt.Println("Call operation: Workspaces_CreateOrUpdate")
	workspacesClient, err := armoperationalinsights.NewWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	workspacesClientCreateOrUpdateResponsePoller, err := workspacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, armoperationalinsights.Workspace{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("val1"),
		},
		Properties: &armoperationalinsights.WorkspaceProperties{
			RetentionInDays: to.Ptr[int32](30),
			SKU: &armoperationalinsights.WorkspaceSKU{
				Name: to.Ptr(armoperationalinsights.WorkspaceSKUNameEnumPerGB2018),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, workspacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Workspaces_List
	fmt.Println("Call operation: Workspaces_List")
	workspacesClientNewListPager := workspacesClient.NewListPager(nil)
	for workspacesClientNewListPager.More() {
		_, err := workspacesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Workspaces_Get
	fmt.Println("Call operation: Workspaces_Get")
	_, err = workspacesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, nil)
	testsuite.Require().NoError(err)

	// From step Workspaces_ListByResourceGroup
	fmt.Println("Call operation: Workspaces_ListByResourceGroup")
	workspacesClientNewListByResourceGroupPager := workspacesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for workspacesClientNewListByResourceGroupPager.More() {
		_, err := workspacesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Workspaces_Update
	fmt.Println("Call operation: Workspaces_Update")
	_, err = workspacesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, armoperationalinsights.WorkspacePatch{
		Properties: &armoperationalinsights.WorkspaceProperties{
			RetentionInDays: to.Ptr[int32](30),
			SKU: &armoperationalinsights.WorkspaceSKU{
				Name: to.Ptr(armoperationalinsights.WorkspaceSKUNameEnumPerGB2018),
			},
			WorkspaceCapping: &armoperationalinsights.WorkspaceCapping{
				DailyQuotaGb: to.Ptr[float64](-1),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Workspaces_Delete
	fmt.Println("Call operation: Workspaces_Delete")
	workspacesClientDeleteResponsePoller, err := workspacesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, &armoperationalinsights.WorkspacesClientBeginDeleteOptions{Force: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, workspacesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step DeletedWorkspaces_List
	fmt.Println("Call operation: DeletedWorkspaces_List")
	deletedWorkspacesClient, err := armoperationalinsights.NewDeletedWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	deletedWorkspacesClientNewListPager := deletedWorkspacesClient.NewListPager(nil)
	for deletedWorkspacesClientNewListPager.More() {
		_, err := deletedWorkspacesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DeletedWorkspaces_ListByResourceGroup
	fmt.Println("Call operation: DeletedWorkspaces_ListByResourceGroup")
	deletedWorkspacesClientNewListByResourceGroupPager := deletedWorkspacesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for deletedWorkspacesClientNewListByResourceGroupPager.More() {
		_, err := deletedWorkspacesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.OperationalInsights/operations
func (testsuite *WorkspacesTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armoperationalinsights.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

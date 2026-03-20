// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type WorkspacesTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	workspaceName     string
}

func (testsuite *WorkspacesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	// Add EUAP redirect policy
	euapOptions := GetEUAPClientOptions()
	testsuite.options.PerCallPolicies = append(testsuite.options.PerCallPolicies, euapOptions.PerCallPolicies...)

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "discovery-test-rg")
	testsuite.workspaceName = "test-workspace"

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	fmt.Println("Created resource group:", testsuite.resourceGroupName)
}

func (testsuite *WorkspacesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestWorkspacesTestSuite(t *testing.T) {
	suite.Run(t, new(WorkspacesTestSuite))
}

// Test listing workspaces by subscription
func (testsuite *WorkspacesTestSuite) TestWorkspacesListBySubscription() {
	fmt.Println("Call operation: Workspaces_ListBySubscription")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewWorkspacesClient().NewListBySubscriptionPager(nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

// Test listing workspaces by resource group
func (testsuite *WorkspacesTestSuite) TestWorkspacesListByResourceGroup() {
	fmt.Println("Call operation: Workspaces_ListByResourceGroup")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewWorkspacesClient().NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

// Test workspace CRUD operations
func (testsuite *WorkspacesTestSuite) SkipTestWorkspacesCRUD() {
	fmt.Println("Call operation: Workspaces_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	workspacesClient := clientFactory.NewWorkspacesClient()

	// Create workspace
	poller, err := workspacesClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.workspaceName,
		armdiscovery.Workspace{
			Location:   to.Ptr(testsuite.location),
			Properties: &armdiscovery.WorkspaceProperties{},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	workspace, err := poller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(workspace.ID)
	fmt.Println("Created workspace:", *workspace.Name)

	// Get workspace
	fmt.Println("Call operation: Workspaces_Get")
	getResp, err := workspacesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(testsuite.workspaceName, *getResp.Name)

	// Update workspace (add tags)
	fmt.Println("Call operation: Workspaces_Update")
	updatePoller, err := workspacesClient.BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.workspaceName,
		armdiscovery.Workspace{
			Location: to.Ptr(testsuite.location),
			Tags: map[string]*string{
				"environment": to.Ptr("test"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	updateResp, err := updatePoller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(updateResp.Tags)

	// Delete workspace
	fmt.Println("Call operation: Workspaces_Delete")
	delPoller, err := workspacesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, nil)
	testsuite.Require().NoError(err)
	_, err = delPoller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	fmt.Println("Deleted workspace:", testsuite.workspaceName)
}

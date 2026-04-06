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

type SavedSearchesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	workspaceName     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SavedSearchesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.workspaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "oisavesearch", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SavedSearchesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestSavedSearchesTestSuite(t *testing.T) {
	suite.Run(t, new(SavedSearchesTestSuite))
}

func (testsuite *SavedSearchesTestSuite) Prepare() {
	var err error
	// From step Workspaces_Create
	fmt.Println("Call operation: Workspaces_CreateOrUpdate")
	workspacesClient, err := armoperationalinsights.NewWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	workspacesClientCreateOrUpdateResponsePoller, err := workspacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, armoperationalinsights.Workspace{
		Location: to.Ptr(testsuite.location),
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
}

// Microsoft.OperationalInsights/workspaces/savedSearches
func (testsuite *SavedSearchesTestSuite) TestSavedSearch() {
	var err error
	// From step SavedSearches_CreateOrUpdate
	fmt.Println("Call operation: SavedSearches_CreateOrUpdate")
	savedSearchesClient, err := armoperationalinsights.NewSavedSearchesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = savedSearchesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "00000000-0000-0000-0000-00000000000", armoperationalinsights.SavedSearch{
		Properties: &armoperationalinsights.SavedSearchProperties{
			Category:           to.Ptr("Saved Search Test Category"),
			DisplayName:        to.Ptr("Create or Update Saved Search Test"),
			FunctionAlias:      to.Ptr("heartbeat_func"),
			FunctionParameters: to.Ptr("a:int=1"),
			Query:              to.Ptr("Heartbeat | summarize Count() by Computer | take a"),
			Tags: []*armoperationalinsights.Tag{
				{
					Name:  to.Ptr("Group"),
					Value: to.Ptr("Computer"),
				}},
			Version: to.Ptr[int64](2),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step SavedSearches_ListByWorkspace
	fmt.Println("Call operation: SavedSearches_ListByWorkspace")
	_, err = savedSearchesClient.ListByWorkspace(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, nil)
	testsuite.Require().NoError(err)

	// From step SavedSearches_Get
	fmt.Println("Call operation: SavedSearches_Get")
	_, err = savedSearchesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "00000000-0000-0000-0000-00000000000", nil)
	testsuite.Require().NoError(err)

	// From step SavedSearches_Delete
	fmt.Println("Call operation: SavedSearches_Delete")
	_, err = savedSearchesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "00000000-0000-0000-0000-00000000000", nil)
	testsuite.Require().NoError(err)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ProjectsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	workspaceName     string
	projectName       string
}

func (testsuite *ProjectsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.workspaceName = "test-wrksp-go01"
	testsuite.projectName = "test-project"
}

func (testsuite *ProjectsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestProjectsTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectsTestSuite))
}

// Test listing projects in a workspace
func (testsuite *ProjectsTestSuite) TestProjectsListByWorkspace() {
	fmt.Println("Call operation: Projects_ListByWorkspace")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewProjectsClient().NewListByWorkspacePager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

// Test project CRUD operations
func (testsuite *ProjectsTestSuite) SkipTestProjectsCRUD() {
	fmt.Println("Call operation: Projects_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	projectsClient := clientFactory.NewProjectsClient()

	// Create project
	poller, err := projectsClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.workspaceName,
		testsuite.projectName,
		armdiscovery.Project{
			Properties: &armdiscovery.ProjectProperties{},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	project, err := poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(project.ID)
	fmt.Println("Created project:", *project.Name)

	// Get project
	fmt.Println("Call operation: Projects_Get")
	getResp, err := projectsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, testsuite.projectName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(testsuite.projectName, *getResp.Name)

	// Update project
	fmt.Println("Call operation: Projects_Update")
	updatePoller, err := projectsClient.BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.workspaceName,
		testsuite.projectName,
		armdiscovery.Project{
			Tags: map[string]*string{
				"environment": to.Ptr("test"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	updateResp, err := updatePoller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(updateResp.ID)

	// Delete project
	fmt.Println("Call operation: Projects_Delete")
	delPoller, err := projectsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, testsuite.projectName, nil)
	testsuite.Require().NoError(err)
	_, err = delPoller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	fmt.Println("Deleted project:", testsuite.projectName)
}

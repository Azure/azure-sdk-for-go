// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// ProjectResourceScenarioSuite mirrors .NET ProjectResourceTests.
type ProjectResourceScenarioSuite struct {
	scenarioBaseSuite

	moverName   string
	projectName string
}

func TestProjectResourceScenarioSuite(t *testing.T) {
	suite.Run(t, new(ProjectResourceScenarioSuite))
}

func (s *ProjectResourceScenarioSuite) SetupSuite() {
	s.setupBase()
	s.moverName = s.generateName("stomover")
	s.projectName = s.generateName("project")
}
func (s *ProjectResourceScenarioSuite) TearDownSuite() { s.teardownBase() }

// TestProjectGetUpdateDelete mirrors .NET ProjectResourceTests.GetUpdateDeleteTest. Creates a project,
// verifies get-equivalence, PATCH-updates the description, then deletes it and confirms gone via 404.
func (s *ProjectResourceScenarioSuite) TestProjectGetUpdateDelete() {
	s.createStorageMover(s.moverName, nil, "")
	created := s.createProject(s.moverName, s.projectName, "")

	client, err := armstoragemover.NewProjectsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.moverName, s.projectName, nil)
	s.Require().NoError(err)
	s.Equal(*created.Name, *getResp.Name)
	s.Equal(*created.ID, *getResp.ID)

	patchResp, err := client.Update(s.ctx, s.resourceGroupName, s.moverName, s.projectName, armstoragemover.ProjectUpdateParameters{
		Properties: &armstoragemover.ProjectUpdateProperties{
			Description: to.Ptr("This is an updated project"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal("This is an updated project", *patchResp.Properties.Description)

	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.moverName, s.projectName, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)

	_, err = client.Get(s.ctx, s.resourceGroupName, s.moverName, s.projectName, nil)
	s.expectResponseError(err)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// ProjectCollectionScenarioSuite mirrors .NET ProjectCollectionTests.
type ProjectCollectionScenarioSuite struct {
	scenarioBaseSuite

	moverName   string
	projectName string
}

func TestProjectCollectionScenarioSuite(t *testing.T) {
	suite.Run(t, new(ProjectCollectionScenarioSuite))
}

func (s *ProjectCollectionScenarioSuite) SetupSuite() {
	s.setupBase()
	s.moverName = s.generateName("stomover")
	s.projectName = s.generateName("project")
}
func (s *ProjectCollectionScenarioSuite) TearDownSuite() { s.teardownBase() }

// TestProjectCreateGetExist mirrors .NET ProjectCollectionTests.CrateGetExistTest (sic — the .NET
// method name is intentionally misspelled). Verifies the project resource type string is
// "microsoft.storagemover/storagemovers/projects" (case-insensitive, since the RP returns the
// canonical-cased form).
func (s *ProjectCollectionScenarioSuite) TestProjectCreateGetExist() {
	s.createStorageMover(s.moverName, nil, "")
	created := s.createProject(s.moverName, s.projectName, "")
	s.Equal(s.projectName, *created.Name)
	s.Nil(created.Properties.Description)
	s.Require().NotNil(created.Type)
	s.Equal("microsoft.storagemover/storagemovers/projects", strings.ToLower(*created.Type))

	client, err := armstoragemover.NewProjectsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.moverName, s.projectName, nil)
	s.Require().NoError(err)
	s.Equal(s.projectName, *getResp.Name)
	s.Nil(getResp.Properties.Description)
	s.Equal("microsoft.storagemover/storagemovers/projects", strings.ToLower(*getResp.Type))

	count := 0
	pager := client.NewListPager(s.resourceGroupName, s.moverName, nil)
	for pager.More() {
		page, err := pager.NextPage(s.ctx)
		s.Require().NoError(err)
		count += len(page.Value)
	}
	s.GreaterOrEqual(count, 1)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// JobRunScenarioSuite mirrors .NET JobRunTests. Without a registered agent, no JobRun is ever
// produced for a job definition. We exercise the read paths anyway: list returns empty, and Get on
// any name returns a non-success response — matching the Python port's combined assertion strategy.
type JobRunScenarioSuite struct {
	scenarioBaseSuite

	storageMoverName string
	projectName      string
	jobDefName       string
	nfsEndpointName  string
	blobEndpointName string
}

func TestJobRunScenarioSuite(t *testing.T) {
	suite.Run(t, new(JobRunScenarioSuite))
}

func (s *JobRunScenarioSuite) SetupSuite() {
	s.setupBase()
	s.storageMoverName = s.generateName("stomover")
	s.projectName = s.generateName("project")
	s.nfsEndpointName = s.generateName("nfsep")
	s.blobEndpointName = s.generateName("blobep")
	s.jobDefName = s.generateName("jobdef")
	s.createStorageMover(s.storageMoverName, nil, "")
	s.createProject(s.storageMoverName, s.projectName, "")
	s.createNfsEndpoint(s.storageMoverName, s.nfsEndpointName, "")
	s.createBlobEndpoint(s.storageMoverName, s.blobEndpointName, "")

	jdClient, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	_, err = jdClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, armstoragemover.JobDefinition{
		Properties: &armstoragemover.JobDefinitionProperties{
			CopyMode:   to.Ptr(armstoragemover.CopyModeAdditive),
			SourceName: to.Ptr(s.nfsEndpointName),
			TargetName: to.Ptr(s.blobEndpointName),
		},
	}, nil)
	s.Require().NoError(err)
}

func (s *JobRunScenarioSuite) TearDownSuite() { s.teardownBase() }

// TestJobRunGetExist mirrors .NET JobRunTests.GetExistTest. With no agent registered, list returns
// no JobRuns and Get on a synthetic name returns a non-success response. We assert both behaviors
// as a combined test.
func (s *JobRunScenarioSuite) TestJobRunGetExist() {
	client, err := armstoragemover.NewJobRunsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	count := 0
	pager := client.NewListPager(s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, nil)
	for pager.More() {
		page, err := pager.NextPage(s.ctx)
		s.Require().NoError(err)
		count += len(page.Value)
	}
	s.Equal(0, count, "expected no JobRuns when no agent is registered")

	_, err = client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, "6e8c0dfe-821a-427d-8d11-a9ed7f1c9c13", nil)
	s.expectResponseError(err)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// JobDefinitionJobRunScenarioSuite mirrors .NET JobDefinitionJobRunTests.
type JobDefinitionJobRunScenarioSuite struct {
	scenarioBaseSuite

	storageMoverName string
	projectName      string
	nfsEndpointName  string
	blobEndpointName string
	jobDefName       string
}

func TestJobDefinitionJobRunScenarioSuite(t *testing.T) {
	suite.Run(t, new(JobDefinitionJobRunScenarioSuite))
}

func (s *JobDefinitionJobRunScenarioSuite) SetupSuite() {
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
}

func (s *JobDefinitionJobRunScenarioSuite) TearDownSuite() { s.teardownBase() }

// TestJobDefinitionJobRun mirrors .NET JobDefinitionJobRunTests.JobDefinitionJobRunTest. Creates a
// job definition, asserts core properties, lists, get-equivalence, then attempts StartJob/StopJob;
// these are expected to fail because no agent is registered to execute them.
func (s *JobDefinitionJobRunScenarioSuite) TestJobDefinitionJobRun() {
	client, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, armstoragemover.JobDefinition{
		Properties: &armstoragemover.JobDefinitionProperties{
			CopyMode:   to.Ptr(armstoragemover.CopyModeAdditive),
			SourceName: to.Ptr(s.nfsEndpointName),
			TargetName: to.Ptr(s.blobEndpointName),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.jobDefName, *createResp.Name)
	s.Equal(s.blobEndpointName, *createResp.Properties.TargetName)
	s.Equal(s.nfsEndpointName, *createResp.Properties.SourceName)
	s.Equal(armstoragemover.CopyModeAdditive, *createResp.Properties.CopyMode)

	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, nil)
	s.Require().NoError(err)
	s.Equal(s.jobDefName, *getResp.Name)

	count := 0
	pager := client.NewListPager(s.resourceGroupName, s.storageMoverName, s.projectName, nil)
	for pager.More() {
		page, err := pager.NextPage(s.ctx)
		s.Require().NoError(err)
		count += len(page.Value)
	}
	s.GreaterOrEqual(count, 1)

	// Get-equivalence (two successive Gets return identical fields).
	getResp2, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, nil)
	s.Require().NoError(err)
	s.Equal(*getResp.Name, *getResp2.Name)
	s.Equal(*getResp.Properties.TargetName, *getResp2.Properties.TargetName)
	s.Equal(*getResp.Properties.SourceName, *getResp2.Properties.SourceName)
	s.Equal(*getResp.ID, *getResp2.ID)

	// StartJob / StopJob require a registered agent. Both calls are expected to fail.
	_, err = client.StartJob(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, nil)
	s.expectResponseError(err)

	_, err = client.StopJob(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, nil)
	s.expectResponseError(err)
}

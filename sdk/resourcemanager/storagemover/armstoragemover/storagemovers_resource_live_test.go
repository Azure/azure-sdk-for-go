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

// StorageMoverResourceScenarioSuite mirrors .NET StorageMoverResourceTests. It provisions a single
// shared Storage Mover and Project in SetupSuite that read-only tests can reuse; tests that mutate
// or delete state create their own instances using names also pre-generated in SetupSuite.
type StorageMoverResourceScenarioSuite struct {
	scenarioBaseSuite

	storageMoverName string
	projectName      string
	endpointName     string
	mutableMoverName string
}

func TestStorageMoverResourceScenarioSuite(t *testing.T) {
	suite.Run(t, new(StorageMoverResourceScenarioSuite))
}

func (s *StorageMoverResourceScenarioSuite) SetupSuite() {
	s.setupBase()
	s.storageMoverName = s.generateName("stomover")
	s.projectName = s.generateName("project")
	s.endpointName = s.generateName("endpoint")
	s.mutableMoverName = s.generateName("stomover")

	s.createStorageMover(s.storageMoverName, map[string]*string{"tag1": to.Ptr("value1")}, "Shared storage mover for resource scenario tests")
	s.createProject(s.storageMoverName, s.projectName, "Shared project for resource scenario tests")
	s.createBlobEndpoint(s.storageMoverName, s.endpointName, "")
}

func (s *StorageMoverResourceScenarioSuite) TearDownSuite() { s.teardownBase() }

// TestGetStorageMover mirrors .NET StorageMoverResourceTests.GetStorageMoverTest. Verifies that two
// successive Get calls return identical resource metadata.
func (s *StorageMoverResourceScenarioSuite) TestGetStorageMover() {
	client, err := armstoragemover.NewStorageMoversClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	first, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, nil)
	s.Require().NoError(err)
	second, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, nil)
	s.Require().NoError(err)

	s.Equal(*first.Name, *second.Name)
	s.Equal(*first.ID, *second.ID)
	s.Equal(*first.Location, *second.Location)
	s.Equal(*first.Type, *second.Type)
	s.Equal(first.Tags, second.Tags)
}

// TestGetStorageMoverAgent mirrors .NET StorageMoverResourceTests.GetStorageMoverAgentTest. Skipped:
// see AgentScenarioSuite.SetupSuite for the rationale.
func (s *StorageMoverResourceScenarioSuite) TestGetStorageMoverAgent() {
	s.T().Skip("requires a registered Hybrid Compute agent; see AgentScenarioSuite for rationale")
}

// TestGetStorageMoverEndpoint mirrors .NET StorageMoverResourceTests.GetStorageMoverEndpointTest.
// Verifies a get on a previously-created blob container endpoint returns the expected fields.
func (s *StorageMoverResourceScenarioSuite) TestGetStorageMoverEndpoint() {
	client, err := armstoragemover.NewEndpointsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	resp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.endpointName, nil)
	s.Require().NoError(err)
	s.Equal(s.endpointName, *resp.Name)
	props, ok := resp.Properties.(*armstoragemover.AzureStorageBlobContainerEndpointProperties)
	s.Require().True(ok, "expected AzureStorageBlobContainerEndpointProperties, got %T", resp.Properties)
	s.NotNil(props.EndpointType)
	s.Equal(armstoragemover.EndpointTypeAzureStorageBlobContainer, *props.EndpointType)
}

// TestGetStorageMoverProject mirrors .NET StorageMoverResourceTests.GetStorageMoverProjectTest.
func (s *StorageMoverResourceScenarioSuite) TestGetStorageMoverProject() {
	client, err := armstoragemover.NewProjectsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	resp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, nil)
	s.Require().NoError(err)
	s.Equal(s.projectName, *resp.Name)
}

// TestStorageMoverUpdateTagsDelete mirrors .NET StorageMoverResourceTests.UpdateAddSetRemoveTagDeletTest.
// Creates a private mover, exercises a description PATCH and three tag map mutations using
// CreateOrUpdate (Go SDK does not provide AddTag/SetTags/RemoveTag helpers), then deletes the mover
// and confirms it is gone via a 404 on Get.
func (s *StorageMoverResourceScenarioSuite) TestStorageMoverUpdateTagsDelete() {
	client, err := armstoragemover.NewStorageMoversClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.mutableMoverName, armstoragemover.StorageMover{
		Location: to.Ptr(s.location),
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.mutableMoverName, *createResp.Name)
	s.Equal(s.location, *createResp.Location)

	// PATCH description.
	patchResp, err := client.Update(s.ctx, s.resourceGroupName, s.mutableMoverName, armstoragemover.UpdateParameters{
		Properties: &armstoragemover.UpdateProperties{
			Description: to.Ptr("This is an updated storage mover"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal("This is an updated storage mover", *patchResp.Properties.Description)

	// Add tag1 -> CreateOrUpdate replaces our tags with {tag1: val1}. The subscription may inject
	// additional tags via Azure Policy (e.g., "Mover", "team"), so we assert on the tag we set
	// rather than the total tag count.
	addResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.mutableMoverName, armstoragemover.StorageMover{
		Location: to.Ptr(s.location),
		Tags:     map[string]*string{"tag1": to.Ptr("val1")},
	}, nil)
	s.Require().NoError(err)
	s.Require().NotNil(addResp.Tags["tag1"])
	s.Equal("val1", *addResp.Tags["tag1"])

	// Set tags to {tag2, tag3}.
	setResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.mutableMoverName, armstoragemover.StorageMover{
		Location: to.Ptr(s.location),
		Tags: map[string]*string{
			"tag2": to.Ptr("val2"),
			"tag3": to.Ptr("val3"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Require().NotNil(setResp.Tags["tag2"])
	s.Equal("val2", *setResp.Tags["tag2"])
	s.Require().NotNil(setResp.Tags["tag3"])
	s.Equal("val3", *setResp.Tags["tag3"])

	// Remove tag2.
	removeResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.mutableMoverName, armstoragemover.StorageMover{
		Location: to.Ptr(s.location),
		Tags:     map[string]*string{"tag3": to.Ptr("val3")},
	}, nil)
	s.Require().NoError(err)
	s.Require().NotNil(removeResp.Tags["tag3"])
	s.Equal("val3", *removeResp.Tags["tag3"])
	_, hasTag2 := removeResp.Tags["tag2"]
	s.False(hasTag2, "tag2 should have been removed")

	// Delete and confirm.
	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.mutableMoverName, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)

	_, err = client.Get(s.ctx, s.resourceGroupName, s.mutableMoverName, nil)
	s.expectResponseError(err)
}

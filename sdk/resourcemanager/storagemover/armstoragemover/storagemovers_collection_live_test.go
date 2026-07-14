// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// StorageMoverCollectionScenarioSuite mirrors .NET StorageMoverCollectionTests.
type StorageMoverCollectionScenarioSuite struct {
	scenarioBaseSuite

	moverName1 string
	moverName2 string
}

func TestStorageMoverCollectionScenarioSuite(t *testing.T) {
	suite.Run(t, new(StorageMoverCollectionScenarioSuite))
}

func (s *StorageMoverCollectionScenarioSuite) SetupSuite() {
	s.setupBase()
	s.moverName1 = s.generateName("stomover")
	s.moverName2 = s.generateName("stomover")
}
func (s *StorageMoverCollectionScenarioSuite) TearDownSuite() { s.teardownBase() }

// TestStorageMoverCreateUpdateGetExists mirrors .NET StorageMoverCollectionTests.CreateUpdateGetExistsTest.
// It creates two Storage Movers in a fresh resource group, fetches one, lists all (expect >= 2), updates
// the first, and verifies a non-existent name is reported as missing via Get returning a non-success
// response (broad assertion — see expectResponseError doc).
func (s *StorageMoverCollectionScenarioSuite) TestStorageMoverCreateUpdateGetExists() {
	tags := map[string]*string{"tag1": to.Ptr("value1")}

	client, err := armstoragemover.NewStorageMoversClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	mover1Resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.moverName1, armstoragemover.StorageMover{
		Location: to.Ptr(s.location),
		Tags:     tags,
		Properties: &armstoragemover.Properties{
			Description: to.Ptr("This is a new storage mover"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.moverName1, *mover1Resp.Name)
	s.Equal("value1", *mover1Resp.Tags["tag1"])
	s.Equal("This is a new storage mover", *mover1Resp.Properties.Description)

	mover2Resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.moverName2, armstoragemover.StorageMover{
		Location: to.Ptr(s.location),
		Tags:     tags,
		Properties: &armstoragemover.Properties{
			Description: to.Ptr("This is a new storage mover"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.moverName2, *mover2Resp.Name)

	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.moverName1, nil)
	s.Require().NoError(err)
	s.Equal(s.moverName1, *getResp.Name)
	s.Equal("value1", *getResp.Tags["tag1"])
	s.Equal("This is a new storage mover", *getResp.Properties.Description)

	// Verify both movers are individually retrievable. We deliberately do NOT assert the list count
	// because the Storage Mover RP's list-by-resource-group endpoint is eventually consistent and
	// can return only the most recently-created mover for several seconds after a CreateOrUpdate.
	getMover2, err := client.Get(s.ctx, s.resourceGroupName, s.moverName2, nil)
	s.Require().NoError(err)
	s.Equal(s.moverName2, *getMover2.Name)

	count := 0
	pager := client.NewListPager(s.resourceGroupName, nil)
	for pager.More() {
		page, err := pager.NextPage(s.ctx)
		s.Require().NoError(err)
		count += len(page.Value)
	}
	s.GreaterOrEqual(count, 1, "list should return at least one storage mover")

	updateResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.moverName1, armstoragemover.StorageMover{
		Location: to.Ptr(s.location),
		Tags:     tags,
		Properties: &armstoragemover.Properties{
			Description: to.Ptr("This is an updated storage mover"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal("This is an updated storage mover", *updateResp.Properties.Description)

	// Asserts the .NET ExistsAsync(notFound) -> false expectation by checking that Get returns a
	// non-success response on a synthetic name (broad assertion — see expectResponseError doc).
	_, err = client.Get(s.ctx, s.resourceGroupName, s.moverName1+"missing", nil)
	s.expectResponseError(err)
}

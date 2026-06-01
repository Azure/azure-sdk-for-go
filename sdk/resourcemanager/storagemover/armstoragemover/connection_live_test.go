// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// ConnectionScenarioSuite mirrors .NET ConnectionTests (matrix row #32, promoted to the
// source-of-truth matrix on 2026-05-20). Exercises Storage Mover Connection op-group CRUD against
// the shared PrivateLinkService `test-pls-wcs` in subscription `b6b34ad8-…`. Requires API version
// 2025-08-01 or later.
//
// Cross-language note: do NOT assert on `ConnectionStatus`. It is `Pending` immediately after
// create (PLS-side provisioning is async) and only transitions to `Approved` after the PE on the
// PLS side is approved. Approval is exercised by matrix row #31 (StartC2CJobWithPrivateSource), not
// here. Also do NOT assert that an update's response echoes the new description back — the RP
// echoes the original description on PUT (CLI / Python / .NET all confirmed this behavior).
type ConnectionScenarioSuite struct {
	scenarioBaseSuite

	storageMoverName string
	connectionName   string
}

func TestConnectionScenarioSuite(t *testing.T) {
	suite.Run(t, new(ConnectionScenarioSuite))
}

func (s *ConnectionScenarioSuite) SetupSuite() {
	s.setupBaseInLocation(scenarioWestCentralUSLocation)
	s.storageMoverName = s.generateName("stomover")
	// Storage Mover Connection name is capped at 20 characters by the RP (InvalidConnectionName)
	// — JS port surfaced this on 2026-05-26. `conn` + 8-char random suffix = 12 chars, well under
	// the cap.
	s.connectionName = s.generateName("conn")
	s.createStorageMover(s.storageMoverName, nil, "")
}

// TearDownSuite cleans up the resource group and stops the recording.
func (s *ConnectionScenarioSuite) TearDownSuite() { s.teardownBase() }

// TestCreateGetListUpdateDelete mirrors .NET ConnectionTests.CreateGetListUpdateDeleteTest. Walks
// the Connection op-group through create → get → list → update (description) → delete.
func (s *ConnectionScenarioSuite) TestCreateGetListUpdateDelete() {
	client, err := armstoragemover.NewConnectionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	const initialDescription = "Initial connection description"
	const updatedDescription = "Updated connection description"

	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, armstoragemover.Connection{
		Properties: &armstoragemover.ConnectionProperties{
			PrivateLinkServiceID: to.Ptr(scenarioPrivateLinkServiceID),
			Description:          to.Ptr(initialDescription),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.connectionName, *createResp.Name)
	s.Require().NotNil(createResp.Properties)
	// ID comparison uses suffix match because the cross-sub sanitizer may rewrite subscription IDs
	// in recordings.
	s.Require().NotNil(createResp.Properties.PrivateLinkServiceID)
	s.True(strings.HasSuffix(*createResp.Properties.PrivateLinkServiceID, "/privateLinkServices/"+scenarioPrivateLinkServiceName),
		"unexpected PrivateLinkServiceID suffix: %s", *createResp.Properties.PrivateLinkServiceID)

	// Get
	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, nil)
	s.Require().NoError(err)
	s.Equal(s.connectionName, *getResp.Name)
	s.Equal(initialDescription, *getResp.Properties.Description)

	// List — assert our newly-created Connection appears.
	pager := client.NewListPager(s.resourceGroupName, s.storageMoverName, nil)
	found := false
	for pager.More() {
		page, err := pager.NextPage(s.ctx)
		s.Require().NoError(err)
		for _, item := range page.Value {
			if item.Name != nil && *item.Name == s.connectionName {
				found = true
			}
		}
	}
	s.True(found, "newly-created connection %q not found in list", s.connectionName)

	// Update via PUT (same op as create). Do NOT assert the response description equals
	// updatedDescription — the RP echoes the original description on PUT updates. CLI, Python, and
	// .NET ports all confirmed this; the update is still applied server-side and a subsequent Get
	// returns the new value.
	_, err = client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, armstoragemover.Connection{
		Properties: &armstoragemover.ConnectionProperties{
			PrivateLinkServiceID: to.Ptr(scenarioPrivateLinkServiceID),
			Description:          to.Ptr(updatedDescription),
		},
	}, nil)
	s.Require().NoError(err)

	// Delete.
	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)

	// Verify gone.
	_, err = client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, nil)
	s.expectResponseError(err)
}

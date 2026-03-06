// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type PrivateEndpointsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	workspaceName     string
	bookshelfName     string
	peConnectionName  string
}

func (testsuite *PrivateEndpointsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	// Add EUAP redirect policy
	euapOptions := GetEUAPClientOptions()
	testsuite.options.PerCallPolicies = append(testsuite.options.PerCallPolicies, euapOptions.PerCallPolicies...)

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.workspaceName = "wrksptest44"
	testsuite.bookshelfName = "test-bookshelf"
	testsuite.peConnectionName = "test-pe-connection"
}

func (testsuite *PrivateEndpointsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestPrivateEndpointsTestSuite(t *testing.T) {
	suite.Run(t, new(PrivateEndpointsTestSuite))
}

// ============ Workspace Private Endpoint Connection Tests ============

// Test listing workspace private endpoint connections
func (testsuite *PrivateEndpointsTestSuite) SkipTestWorkspacePrivateEndpointConnectionsListByWorkspace() {
	fmt.Println("Call operation: WorkspacePrivateEndpointConnections_ListByWorkspace")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewWorkspacePrivateEndpointConnectionsClient().NewListByWorkspacePager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a workspace private endpoint connection
func (testsuite *PrivateEndpointsTestSuite) SkipTestWorkspacePrivateEndpointConnectionsGet() {
	fmt.Println("Call operation: WorkspacePrivateEndpointConnections_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewWorkspacePrivateEndpointConnectionsClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, testsuite.peConnectionName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a workspace private endpoint connection
func (testsuite *PrivateEndpointsTestSuite) SkipTestWorkspacePrivateEndpointConnectionsCreateOrUpdate() {
	fmt.Println("Call operation: WorkspacePrivateEndpointConnections_CreateOrUpdate")
	// Requires private endpoint configuration and network setup
}

// Test deleting a workspace private endpoint connection
func (testsuite *PrivateEndpointsTestSuite) SkipTestWorkspacePrivateEndpointConnectionsDelete() {
	fmt.Println("Call operation: WorkspacePrivateEndpointConnections_Delete")
	// Requires existing private endpoint connection
}

// ============ Workspace Private Link Resource Tests ============

// Test listing workspace private link resources
func (testsuite *PrivateEndpointsTestSuite) SkipTestWorkspacePrivateLinkResourcesListByWorkspace() {
	fmt.Println("Call operation: WorkspacePrivateLinkResources_ListByWorkspace")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewWorkspacePrivateLinkResourcesClient().NewListByWorkspacePager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a workspace private link resource
func (testsuite *PrivateEndpointsTestSuite) SkipTestWorkspacePrivateLinkResourcesGet() {
	fmt.Println("Call operation: WorkspacePrivateLinkResources_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewWorkspacePrivateLinkResourcesClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "workspace", nil)
	testsuite.Require().NoError(err)
}

// ============ Bookshelf Private Endpoint Connection Tests ============

// Test listing bookshelf private endpoint connections
func (testsuite *PrivateEndpointsTestSuite) SkipTestBookshelfPrivateEndpointConnectionsListByBookshelf() {
	fmt.Println("Call operation: BookshelfPrivateEndpointConnections_ListByBookshelf")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewBookshelfPrivateEndpointConnectionsClient().NewListByBookshelfPager(testsuite.resourceGroupName, testsuite.bookshelfName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a bookshelf private endpoint connection
func (testsuite *PrivateEndpointsTestSuite) SkipTestBookshelfPrivateEndpointConnectionsGet() {
	fmt.Println("Call operation: BookshelfPrivateEndpointConnections_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewBookshelfPrivateEndpointConnectionsClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.bookshelfName, testsuite.peConnectionName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a bookshelf private endpoint connection
func (testsuite *PrivateEndpointsTestSuite) SkipTestBookshelfPrivateEndpointConnectionsCreateOrUpdate() {
	fmt.Println("Call operation: BookshelfPrivateEndpointConnections_CreateOrUpdate")
	// Requires private endpoint configuration and network setup
}

// Test deleting a bookshelf private endpoint connection
func (testsuite *PrivateEndpointsTestSuite) SkipTestBookshelfPrivateEndpointConnectionsDelete() {
	fmt.Println("Call operation: BookshelfPrivateEndpointConnections_Delete")
	// Requires existing private endpoint connection
}

// ============ Bookshelf Private Link Resource Tests ============

// Test listing bookshelf private link resources
func (testsuite *PrivateEndpointsTestSuite) SkipTestBookshelfPrivateLinkResourcesListByBookshelf() {
	fmt.Println("Call operation: BookshelfPrivateLinkResources_ListByBookshelf")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewBookshelfPrivateLinkResourcesClient().NewListByBookshelfPager(testsuite.resourceGroupName, testsuite.bookshelfName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a bookshelf private link resource
func (testsuite *PrivateEndpointsTestSuite) SkipTestBookshelfPrivateLinkResourcesGet() {
	fmt.Println("Call operation: BookshelfPrivateLinkResources_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewBookshelfPrivateLinkResourcesClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.bookshelfName, "bookshelf", nil)
	testsuite.Require().NoError(err)
}

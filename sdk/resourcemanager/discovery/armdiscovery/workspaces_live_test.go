// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type WorkspacesTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	workspaceName     string
}

func (testsuite *WorkspacesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.workspaceName = "test-wrksp-go01"
}

func (testsuite *WorkspacesTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestWorkspacesTestSuite(t *testing.T) {
	suite.Run(t, new(WorkspacesTestSuite))
}

// Test listing workspaces by subscription
func (testsuite *WorkspacesTestSuite) TestWorkspacesListBySubscription() {
	fmt.Println("Call operation: Workspaces_ListBySubscription")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewWorkspacesClient().NewListBySubscriptionPager(nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

// Test listing workspaces by resource group
func (testsuite *WorkspacesTestSuite) TestWorkspacesListByResourceGroup() {
	fmt.Println("Call operation: Workspaces_ListByResourceGroup")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewWorkspacesClient().NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

// Test creating a workspace (skipped - CreateOrUpdate not yet recorded)
func (testsuite *WorkspacesTestSuite) SkipTestWorkspacesCreateOrUpdate() {
	fmt.Println("Call operation: Workspaces_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	miID := "/subscriptions/" + testsuite.subscriptionId + "/resourcegroups/" + testsuite.resourceGroupName + "/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity"
	subnetDefault := "/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/virtualNetworks/newapiv/subnets/default"
	subnetDefault2 := "/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/virtualNetworks/newapiv/subnets/default2"
	subnetDefault3 := "/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/virtualNetworks/newapiv/subnets/default3"
	logAnalyticsClusterID := "/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.OperationalInsights/clusters/mycluse"

	workspacesClient := clientFactory.NewWorkspacesClient()
	poller, err := workspacesClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.workspaceName,
		armdiscovery.Workspace{
			Location: to.Ptr(testsuite.location),
			Properties: &armdiscovery.WorkspaceProperties{
				WorkspaceIdentity: &armdiscovery.Identity{
					ID: to.Ptr(miID),
				},
				AgentSubnetID:          to.Ptr(subnetDefault3),
				PrivateEndpointSubnetID: to.Ptr(subnetDefault),
				WorkspaceSubnetID:      to.Ptr(subnetDefault2),
				CustomerManagedKeys:    to.Ptr(armdiscovery.CustomerManagedKeysEnabled),
				KeyVaultProperties: &armdiscovery.KeyVaultProperties{
					KeyName:     to.Ptr("discoverykey"),
					KeyVaultURI: to.Ptr("https://newapik.vault.azure.net/"),
					KeyVersion:  to.Ptr("2c9db3cf55d247b4a1c1831fbbdad906"),
				},
				LogAnalyticsClusterID: to.Ptr(logAnalyticsClusterID),
				PublicNetworkAccess:   to.Ptr(armdiscovery.PublicNetworkAccessDisabled),
				SupercomputerIDs:      []*string{},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	workspace, err := poller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(workspace.ID)
	fmt.Println("Created workspace:", *workspace.Name)
}

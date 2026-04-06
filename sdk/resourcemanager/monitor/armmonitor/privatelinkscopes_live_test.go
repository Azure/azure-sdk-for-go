// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmonitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/stretchr/testify/suite"
)

type PrivatelinkscopesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	scopeName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *PrivatelinkscopesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.scopeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "linkscopena", 17, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *PrivatelinkscopesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPrivatelinkscopesTestSuite(t *testing.T) {
	suite.Run(t, new(PrivatelinkscopesTestSuite))
}

// microsoft.insights/privateLinkScopes
func (testsuite *PrivatelinkscopesTestSuite) TestPrivatelinkscope() {
	var err error
	// From step PrivateLinkScopes_Create
	fmt.Println("Call operation: PrivateLinkScopes_CreateOrUpdate")
	privateLinkScopesClient, err := armmonitor.NewPrivateLinkScopesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = privateLinkScopesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.scopeName, armmonitor.AzureMonitorPrivateLinkScope{
		Location: to.Ptr("Global"),
		Properties: &armmonitor.AzureMonitorPrivateLinkScopeProperties{
			AccessModeSettings: &armmonitor.AccessModeSettings{
				QueryAccessMode:     to.Ptr(armmonitor.AccessModeOpen),
				IngestionAccessMode: to.Ptr(armmonitor.AccessModePrivateOnly),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkScopes_List
	fmt.Println("Call operation: PrivateLinkScopes_List")
	privateLinkScopesClientNewListPager := privateLinkScopesClient.NewListPager(nil)
	for privateLinkScopesClientNewListPager.More() {
		_, err := privateLinkScopesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateLinkScopes_ListByResourceGroup
	fmt.Println("Call operation: PrivateLinkScopes_ListByResourceGroup")
	privateLinkScopesClientNewListByResourceGroupPager := privateLinkScopesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for privateLinkScopesClientNewListByResourceGroupPager.More() {
		_, err := privateLinkScopesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateLinkScopes_Get
	fmt.Println("Call operation: PrivateLinkScopes_Get")
	_, err = privateLinkScopesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.scopeName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkScopes_UpdateTags
	fmt.Println("Call operation: PrivateLinkScopes_UpdateTags")
	_, err = privateLinkScopesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.scopeName, armmonitor.TagsResource{
		Tags: map[string]*string{
			"Tag1": to.Ptr("Value1"),
			"Tag2": to.Ptr("Value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkScopes_Delete
	fmt.Println("Call operation: PrivateLinkScopes_Delete")
	privateLinkScopesClientDeleteResponsePoller, err := privateLinkScopesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.scopeName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateLinkScopesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

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

type ToolsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	toolName          string
}

func (testsuite *ToolsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	// Add EUAP redirect policy
	euapOptions := GetEUAPClientOptions()
	testsuite.options.PerCallPolicies = append(testsuite.options.PerCallPolicies, euapOptions.PerCallPolicies...)

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.toolName = "test-tool"
}

func (testsuite *ToolsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestToolsTestSuite(t *testing.T) {
	suite.Run(t, new(ToolsTestSuite))
}

// Test listing tools by subscription
func (testsuite *ToolsTestSuite) SkipTestToolsListBySubscription() {
	fmt.Println("Call operation: Tools_ListBySubscription")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewToolsClient().NewListBySubscriptionPager(nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test listing tools by resource group
func (testsuite *ToolsTestSuite) SkipTestToolsListByResourceGroup() {
	fmt.Println("Call operation: Tools_ListByResourceGroup")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewToolsClient().NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a tool
func (testsuite *ToolsTestSuite) SkipTestToolsGet() {
	fmt.Println("Call operation: Tools_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewToolsClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.toolName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a tool
func (testsuite *ToolsTestSuite) SkipTestToolsCreateOrUpdate() {
	fmt.Println("Call operation: Tools_CreateOrUpdate")
	// Requires proper tool configuration
}

// Test updating a tool
func (testsuite *ToolsTestSuite) SkipTestToolsUpdate() {
	fmt.Println("Call operation: Tools_Update")
	// Requires existing tool
}

// Test deleting a tool
func (testsuite *ToolsTestSuite) SkipTestToolsDelete() {
	fmt.Println("Call operation: Tools_Delete")
	// Requires existing tool
}

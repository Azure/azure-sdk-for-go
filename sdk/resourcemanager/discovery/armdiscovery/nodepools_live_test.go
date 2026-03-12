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

type NodePoolsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	supercomputerName string
	nodePoolName      string
}

func (testsuite *NodePoolsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	// Add EUAP redirect policy
	euapOptions := GetEUAPClientOptions()
	testsuite.options.PerCallPolicies = append(testsuite.options.PerCallPolicies, euapOptions.PerCallPolicies...)

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "olawal")
	testsuite.supercomputerName = "test-supercomputer"
	testsuite.nodePoolName = "test-nodepool"
}

func (testsuite *NodePoolsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestNodePoolsTestSuite(t *testing.T) {
	suite.Run(t, new(NodePoolsTestSuite))
}

// Test listing node pools by supercomputer
func (testsuite *NodePoolsTestSuite) SkipTestNodePoolsListBySupercomputer() {
	fmt.Println("Call operation: NodePools_ListBySupercomputer")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewNodePoolsClient().NewListBySupercomputerPager(testsuite.resourceGroupName, testsuite.supercomputerName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a node pool
func (testsuite *NodePoolsTestSuite) SkipTestNodePoolsGet() {
	fmt.Println("Call operation: NodePools_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewNodePoolsClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.supercomputerName, testsuite.nodePoolName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a node pool
func (testsuite *NodePoolsTestSuite) SkipTestNodePoolsCreateOrUpdate() {
	fmt.Println("Call operation: NodePools_CreateOrUpdate")
	// Requires supercomputer with proper configuration
}

// Test updating a node pool
func (testsuite *NodePoolsTestSuite) SkipTestNodePoolsUpdate() {
	fmt.Println("Call operation: NodePools_Update")
	// Requires existing node pool
}

// Test deleting a node pool
func (testsuite *NodePoolsTestSuite) SkipTestNodePoolsDelete() {
	fmt.Println("Call operation: NodePools_Delete")
	// Requires existing node pool
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/suite"
)

type SubnetsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *SubnetsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/network/armnetwork/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *SubnetsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSubnetsClient(t *testing.T) {
	suite.Run(t, new(SubnetsClientTestSuite))
}

func (testsuite *SubnetsClientTestSuite) TestSubnetsCRUD() {
	// create virtual network
	vnClient, err := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vnName := "go-test-vn"
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		armnetwork.VirtualNetwork{
			Location: to.Ptr(testsuite.location),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.Ptr("10.1.0.0/16"),
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	vnResp, err := testutil.PollForTest(testsuite.ctx, vnPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vnName, *vnResp.Name)

	// create subnet
	subClient, err := armnetwork.NewSubnetsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	subName := "go-test-subnet"
	subPoller, err := subClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.Ptr("10.1.10.0/24"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	subResp, err := testutil.PollForTest(testsuite.ctx, subPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(subName, *subResp.Name)

	// get subnet
	getResp, err := subClient.Get(testsuite.ctx, testsuite.resourceGroupName, vnName, subName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(subName, *getResp.Name)

	// list subnet
	listPager := subClient.NewListPager(testsuite.resourceGroupName, vnName, nil)
	testsuite.Require().Equal(true, listPager.More())

	// delete subnet
	delPoller, err := subClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, vnName, subName, nil)
	testsuite.Require().NoError(err)
	delResp, err := testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)
	//testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
	_ = delResp
}

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

type VirtualNetworksClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *VirtualNetworksClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/network/armnetwork/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *VirtualNetworksClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestVirtualNetworksClient(t *testing.T) {
	suite.Run(t, new(VirtualNetworksClientTestSuite))
}

func (testsuite *VirtualNetworksClientTestSuite) TestVirtualMachineCRUD() {
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

	//virtual network update tags
	tagResp, err := vnClient.UpdateTags(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		armnetwork.TagsObject{
			Tags: map[string]*string{
				"tag1": to.Ptr("value1"),
				"tag2": to.Ptr("value2"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("value1", *tagResp.Tags["tag1"])

	// get virtual network
	vnResp2, err := vnClient.Get(testsuite.ctx, testsuite.resourceGroupName, vnName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vnName, *vnResp2.Name)

	//virtual network list
	listPager := vnClient.NewListPager(testsuite.resourceGroupName, nil)
	testsuite.Require().Equal(true, listPager.More())

	//virtual network delete
	delPoller, err := vnClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, vnName, nil)
	testsuite.Require().NoError(err)
	delResp, err := testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)
	//testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
	_ = delResp
}

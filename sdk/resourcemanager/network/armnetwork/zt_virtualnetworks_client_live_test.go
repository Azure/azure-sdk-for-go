//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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
	vnClient := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	vnName := "go-test-vn"
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		armnetwork.VirtualNetwork{
			Location: to.StringPtr(testsuite.location),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var vnResp armnetwork.VirtualNetworksClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vnPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if vnPoller.Poller.Done() {
				vnResp, err = vnPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		vnResp, err = vnPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(vnName, *vnResp.Name)

	//virtual network update tags
	tagResp, err := vnClient.UpdateTags(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		armnetwork.TagsObject{
			Tags: map[string]*string{
				"tag1": to.StringPtr("value1"),
				"tag2": to.StringPtr("value2"),
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
	listPager := vnClient.List(testsuite.resourceGroupName, nil)
	testsuite.Require().Equal(true, listPager.NextPage(context.Background()))

	//virtual network delete
	delPoller, err := vnClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, vnName, nil)
	testsuite.Require().NoError(err)
	var delResp armnetwork.VirtualNetworksClientDeleteResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = delPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if delPoller.Poller.Done() {
				delResp, err = delPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		delResp, err = delPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}

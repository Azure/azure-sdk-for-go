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
	testsuite.resourceGroupName = *testutil.CreateResourceGroup(testsuite.T(), testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location).Name
}

func (testsuite *SubnetsClientTestSuite) TearDownSuite() {
	testutil.DeleteResourceGroup(testsuite.T(), testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testutil.StopRecording(testsuite.T())
}

func TestSubnetsClient(t *testing.T) {
	suite.Run(t, new(SubnetsClientTestSuite))
}

func (testsuite *SubnetsClientTestSuite) TestSubnetsCRUD() {
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

	// create subnet
	subClient := armnetwork.NewSubnetsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	subName := "go-test-subnet"
	subPoller, err := subClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.1.10.0/24"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var subResp armnetwork.SubnetsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = subPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if subPoller.Poller.Done() {
				subResp, err = subPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		subResp, err = subPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(subName, *subResp.Name)

	/* need registered for feature Microsoft.Network/AllowPrepareNetworkPoliciesAction
	// prepare network policy
	pnpPoller, err := subClient.BeginPrepareNetworkPolicies(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		subName,
		armnetwork.PrepareNetworkPoliciesRequest{
			ServiceName: to.StringPtr("Microsoft.Sql/managedInstances"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var pnpResp armnetwork.SubnetsClientPrepareNetworkPoliciesResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = pnpPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if pnpPoller.Poller.Done() {
				pnpResp, err = pnpPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		pnpResp, err = pnpPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(200, pnpResp.RawResponse.StatusCode)

	// unprepare network policy
	unPnpPoller, err := subClient.BeginUnprepareNetworkPolicies(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		subName,
		armnetwork.UnprepareNetworkPoliciesRequest{
			ServiceName: to.StringPtr("Microsoft.Sql/managedInstances"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var unPnpResp armnetwork.SubnetsClientUnprepareNetworkPoliciesResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = unPnpPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if unPnpPoller.Poller.Done() {
				unPnpResp, err = unPnpPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		unPnpResp, err = unPnpPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(200, unPnpResp.RawResponse.StatusCode)
	*/

	// get subnet
	getResp, err := subClient.Get(testsuite.ctx, testsuite.resourceGroupName, vnName, subName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(subName, *getResp.Name)

	// list subnet
	listPager := subClient.List(testsuite.resourceGroupName, vnName, nil)
	testsuite.Require().Equal(true, listPager.NextPage(context.Background()))

	// delete subnet
	delPoller, err := subClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, vnName, subName, nil)
	testsuite.Require().NoError(err)
	var delResp armnetwork.SubnetsClientDeleteResponse
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

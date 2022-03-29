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

type IPGroupsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *IPGroupsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/network/armnetwork/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *IPGroupsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestIPGroupsClient(t *testing.T) {
	suite.Run(t, new(IPGroupsClientTestSuite))
}

func (testsuite *IPGroupsClientTestSuite) TestIPGroupsCRUD() {
	// create ip group
	ipgClient := armnetwork.NewIPGroupsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	ipgName := "go-test-ipg"
	ipgPoller, err := ipgClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		ipgName,
		armnetwork.IPGroup{
			Location: to.StringPtr(testsuite.location),
			Properties: &armnetwork.IPGroupPropertiesFormat{
				IPAddresses: []*string{
					to.StringPtr("13.64.39.16/32"),
					to.StringPtr("40.74.146.80/31"),
					to.StringPtr("40.74.147.32/28"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var resp armnetwork.IPGroupsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = ipgPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if ipgPoller.Poller.Done() {
				resp, err = ipgPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		resp, err = ipgPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(ipgName, *resp.Name)

	// update
	updateResp, err := ipgClient.UpdateGroups(
		testsuite.ctx,
		testsuite.resourceGroupName,
		ipgName,
		armnetwork.TagsObject{
			Tags: map[string]*string{
				"test": to.StringPtr("live"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("live", *updateResp.Tags["test"])

	// get ip group
	getResp, err := ipgClient.Get(testsuite.ctx, testsuite.resourceGroupName, ipgName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(ipgName, *getResp.Name)

	// list ip group
	listPager := ipgClient.List(nil)
	testsuite.Require().True(listPager.NextPage(context.Background()))

	// delete ip group
	delPoller, err := ipgClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, ipgName, nil)
	testsuite.Require().NoError(err)
	var delResp armnetwork.IPGroupsClientDeleteResponse
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

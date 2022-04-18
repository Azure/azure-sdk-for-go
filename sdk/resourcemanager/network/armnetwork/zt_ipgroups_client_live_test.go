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
	ipgClient, err := armnetwork.NewIPGroupsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	ipgName := "go-test-ipg"
	ipgPoller, err := ipgClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		ipgName,
		armnetwork.IPGroup{
			Location: to.Ptr(testsuite.location),
			Properties: &armnetwork.IPGroupPropertiesFormat{
				IPAddresses: []*string{
					to.Ptr("13.64.39.16/32"),
					to.Ptr("40.74.146.80/31"),
					to.Ptr("40.74.147.32/28"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	resp, err := testutil.PollForTest(testsuite.ctx, ipgPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(ipgName, *resp.Name)

	// update
	updateResp, err := ipgClient.UpdateGroups(
		testsuite.ctx,
		testsuite.resourceGroupName,
		ipgName,
		armnetwork.TagsObject{
			Tags: map[string]*string{
				"test": to.Ptr("live"),
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
	listPager := ipgClient.NewListPager(nil)
	testsuite.Require().True(listPager.More())

	// delete ip group
	delPoller, err := ipgClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, ipgName, nil)
	testsuite.Require().NoError(err)
	delResp, err := testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)
	//testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
	_ = delResp
}

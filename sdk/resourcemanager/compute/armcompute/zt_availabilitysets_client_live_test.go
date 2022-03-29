//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type AvailabilitySetsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *AvailabilitySetsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/compute/armcompute/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *AvailabilitySetsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAvailabilitySetsClient(t *testing.T) {
	suite.Run(t, new(AvailabilitySetsClientTestSuite))
}

func (testsuite *AvailabilitySetsClientTestSuite) TestAvailabilitySetsCRUD() {
	// create availability sets
	client := armcompute.NewAvailabilitySetsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	name := "go-test-availability"
	resp, err := client.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		name,
		armcompute.AvailabilitySet{
			Location: to.StringPtr("westus"),
			SKU: &armcompute.SKU{
				Name: to.StringPtr(string(armcompute.AvailabilitySetSKUTypesAligned)),
			},
			Properties: &armcompute.AvailabilitySetProperties{
				PlatformFaultDomainCount:  to.Int32Ptr(1),
				PlatformUpdateDomainCount: to.Int32Ptr(1),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*resp.Name, name)

	// get
	getResp, err := client.Get(testsuite.ctx, testsuite.resourceGroupName, name, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*getResp.Name, name)

	// list
	listPager := client.List(testsuite.resourceGroupName, nil)
	testsuite.Require().Equal(listPager.Err(), nil)

	// list available size
	listResp, err := client.ListAvailableSizes(testsuite.ctx, testsuite.resourceGroupName, name, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(listResp.RawResponse.StatusCode, 200)

	// list by subscription
	listBySubscription := client.ListBySubscription(nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(true, listBySubscription.NextPage(context.Background()))

	// update
	updateResp, err := client.Update(
		testsuite.ctx,
		testsuite.resourceGroupName,
		name,
		armcompute.AvailabilitySetUpdate{
			Tags: map[string]*string{
				"tag": to.StringPtr("value"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(name, *updateResp.Name)

	// delete
	delResp, err := client.Delete(testsuite.ctx, testsuite.resourceGroupName, name, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}

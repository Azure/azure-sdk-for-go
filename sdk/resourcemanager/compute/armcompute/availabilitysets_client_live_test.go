//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
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
	fmt.Println("Call operation: AvailabilitySets_CreateOrUpdate")
	client, err := armcompute.NewAvailabilitySetsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	name := "go-test-availability"
	resp, err := client.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		name,
		armcompute.AvailabilitySet{
			Location: to.Ptr("westus"),
			SKU: &armcompute.SKU{
				Name: to.Ptr(string(armcompute.AvailabilitySetSKUTypesAligned)),
			},
			Properties: &armcompute.AvailabilitySetProperties{
				PlatformFaultDomainCount:  to.Ptr[int32](1),
				PlatformUpdateDomainCount: to.Ptr[int32](1),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*resp.Name, name)

	// get
	fmt.Println("Call operation: AvailabilitySets_Get")
	getResp, err := client.Get(testsuite.ctx, testsuite.resourceGroupName, name, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*getResp.Name, name)

	// list
	fmt.Println("Call operation: AvailabilitySets_List")
	listPager := client.NewListPager(testsuite.resourceGroupName, nil)
	testsuite.Require().True(listPager.More())

	// list available size
	fmt.Println("Call operation: AvailabilitySets_ListAvailableSize")
	listResp := client.NewListAvailableSizesPager(testsuite.resourceGroupName, name, nil)
	testsuite.Require().True(listResp.More())

	// list by subscription
	fmt.Println("Call operation: AvailabilitySets_ListBySubscription")
	listBySubscription := client.NewListBySubscriptionPager(nil)
	testsuite.Require().True(listBySubscription.More())

	// update
	fmt.Println("Call operation: AvailabilitySets_Update")
	updateResp, err := client.Update(
		testsuite.ctx,
		testsuite.resourceGroupName,
		name,
		armcompute.AvailabilitySetUpdate{
			Tags: map[string]*string{
				"tag": to.Ptr("value"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(name, *updateResp.Name)

	// delete
	fmt.Println("Call operation: AvailabilitySets_Delete")
	_, err = client.Delete(testsuite.ctx, testsuite.resourceGroupName, name, nil)
	testsuite.Require().NoError(err)
}

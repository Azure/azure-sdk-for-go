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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type CapacityReservationTestSuite struct {
	suite.Suite

	ctx                          context.Context
	cred                         azcore.TokenCredential
	options                      *arm.ClientOptions
	capacityReservationGroupName string
	location                     string
	resourceGroupName            string
	subscriptionId               string
}

func (testsuite *CapacityReservationTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.capacityReservationGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "capacityre", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *CapacityReservationTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestCapacityReservationTestSuite(t *testing.T) {
	suite.Run(t, new(CapacityReservationTestSuite))
}

func (testsuite *CapacityReservationTestSuite) Prepare() {
	var err error
	// From step CapacityReservationGroups_CreateOrUpdate
	fmt.Println("Call operation: CapacityReservationGroups_CreateOrUpdate")
	capacityReservationGroupsClient, err := armcompute.NewCapacityReservationGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = capacityReservationGroupsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.capacityReservationGroupName, armcompute.CapacityReservationGroup{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"department": to.Ptr("finance"),
		},
		Zones: []*string{
			to.Ptr("1"),
			to.Ptr("2")},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/capacityReservationGroups
func (testsuite *CapacityReservationTestSuite) TestCapacityReservationGroups() {
	var err error
	// From step CapacityReservationGroups_ListBySubscription
	fmt.Println("Call operation: CapacityReservationGroups_ListBySubscription")
	capacityReservationGroupsClient, err := armcompute.NewCapacityReservationGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	capacityReservationGroupsClientNewListBySubscriptionPager := capacityReservationGroupsClient.NewListBySubscriptionPager(&armcompute.CapacityReservationGroupsClientListBySubscriptionOptions{Expand: to.Ptr(armcompute.ExpandTypesForGetCapacityReservationGroupsVirtualMachinesRef)})
	for capacityReservationGroupsClientNewListBySubscriptionPager.More() {
		_, err := capacityReservationGroupsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step CapacityReservationGroups_ListByResourceGroup
	fmt.Println("Call operation: CapacityReservationGroups_ListByResourceGroup")
	capacityReservationGroupsClientNewListByResourceGroupPager := capacityReservationGroupsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armcompute.CapacityReservationGroupsClientListByResourceGroupOptions{Expand: to.Ptr(armcompute.ExpandTypesForGetCapacityReservationGroupsVirtualMachinesRef)})
	for capacityReservationGroupsClientNewListByResourceGroupPager.More() {
		_, err := capacityReservationGroupsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step CapacityReservationGroups_Get
	fmt.Println("Call operation: CapacityReservationGroups_Get")
	_, err = capacityReservationGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.capacityReservationGroupName, &armcompute.CapacityReservationGroupsClientGetOptions{Expand: to.Ptr(armcompute.CapacityReservationGroupInstanceViewTypesInstanceView)})
	testsuite.Require().NoError(err)

	// From step CapacityReservationGroups_Update
	fmt.Println("Call operation: CapacityReservationGroups_Update")
	_, err = capacityReservationGroupsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.capacityReservationGroupName, armcompute.CapacityReservationGroupUpdate{}, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *CapacityReservationTestSuite) Cleanup() {
	var err error
	// From step CapacityReservationGroups_Delete
	fmt.Println("Call operation: CapacityReservationGroups_Delete")
	capacityReservationGroupsClient, err := armcompute.NewCapacityReservationGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = capacityReservationGroupsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.capacityReservationGroupName, nil)
	testsuite.Require().NoError(err)
}

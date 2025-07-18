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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ProximityPlacementGroupTestSuite struct {
	suite.Suite

	ctx                         context.Context
	cred                        azcore.TokenCredential
	options                     *arm.ClientOptions
	location                    string
	proximityPlacementGroupName string
	resourceGroupName           string
	subscriptionId              string
}

func (testsuite *ProximityPlacementGroupTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.proximityPlacementGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "proximityp", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *ProximityPlacementGroupTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestProximityPlacementGroupTestSuite(t *testing.T) {
	suite.Run(t, new(ProximityPlacementGroupTestSuite))
}

// Microsoft.Compute/proximityPlacementGroups
func (testsuite *ProximityPlacementGroupTestSuite) TestProximityPlacementGroups() {
	var err error
	// From step ProximityPlacementGroups_CreateOrUpdate
	fmt.Println("Call operation: ProximityPlacementGroups_CreateOrUpdate")
	proximityPlacementGroupsClient, err := armcompute.NewProximityPlacementGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = proximityPlacementGroupsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.proximityPlacementGroupName, armcompute.ProximityPlacementGroup{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.ProximityPlacementGroupProperties{
			Intent: &armcompute.ProximityPlacementGroupPropertiesIntent{
				VMSizes: []*string{
					to.Ptr("Basic_A0"),
					to.Ptr("Basic_A2")},
			},
			ProximityPlacementGroupType: to.Ptr(armcompute.ProximityPlacementGroupTypeStandard),
		},
		Zones: []*string{
			to.Ptr("1")},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ProximityPlacementGroups_ListBySubscription
	fmt.Println("Call operation: ProximityPlacementGroups_ListBySubscription")
	proximityPlacementGroupsClientNewListBySubscriptionPager := proximityPlacementGroupsClient.NewListBySubscriptionPager(nil)
	for proximityPlacementGroupsClientNewListBySubscriptionPager.More() {
		_, err := proximityPlacementGroupsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ProximityPlacementGroups_ListByResourceGroup
	fmt.Println("Call operation: ProximityPlacementGroups_ListByResourceGroup")
	proximityPlacementGroupsClientNewListByResourceGroupPager := proximityPlacementGroupsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for proximityPlacementGroupsClientNewListByResourceGroupPager.More() {
		_, err := proximityPlacementGroupsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ProximityPlacementGroups_Get
	fmt.Println("Call operation: ProximityPlacementGroups_Get")
	_, err = proximityPlacementGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.proximityPlacementGroupName, &armcompute.ProximityPlacementGroupsClientGetOptions{IncludeColocationStatus: nil})
	testsuite.Require().NoError(err)

	// From step ProximityPlacementGroups_Update
	fmt.Println("Call operation: ProximityPlacementGroups_Update")
	_, err = proximityPlacementGroupsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.proximityPlacementGroupName, armcompute.ProximityPlacementGroupUpdate{
		Tags: map[string]*string{
			"additionalProp1": to.Ptr("string"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ProximityPlacementGroups_Delete
	fmt.Println("Call operation: ProximityPlacementGroups_Delete")
	_, err = proximityPlacementGroupsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.proximityPlacementGroupName, nil)
	testsuite.Require().NoError(err)
}

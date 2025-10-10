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

type AvailabilitySetTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	location            string
	resourceGroupName   string
	subscriptionId      string
	availabilitySetName string
}

func (testsuite *AvailabilitySetTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.availabilitySetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "availabili", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *AvailabilitySetTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAvailabilitySetTestSuite(t *testing.T) {
	suite.Run(t, new(AvailabilitySetTestSuite))
}

// Microsoft.Compute/availabilitySets
func (testsuite *AvailabilitySetTestSuite) TestAvailabilitySets() {
	var err error
	// From step AvailabilitySets_CreateOrUpdate
	fmt.Println("Call operation: AvailabilitySets_CreateOrUpdate")
	availabilitySetsClient, err := armcompute.NewAvailabilitySetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = availabilitySetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.availabilitySetName, armcompute.AvailabilitySet{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.AvailabilitySetProperties{
			PlatformFaultDomainCount:  to.Ptr[int32](2),
			PlatformUpdateDomainCount: to.Ptr[int32](20),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AvailabilitySets_ListBySubscription
	fmt.Println("Call operation: AvailabilitySets_ListBySubscription")
	availabilitySetsClientNewListBySubscriptionPager := availabilitySetsClient.NewListBySubscriptionPager(&armcompute.AvailabilitySetsClientListBySubscriptionOptions{Expand: nil})
	for availabilitySetsClientNewListBySubscriptionPager.More() {
		_, err := availabilitySetsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AvailabilitySets_List
	fmt.Println("Call operation: AvailabilitySets_List")
	availabilitySetsClientNewListPager := availabilitySetsClient.NewListPager(testsuite.resourceGroupName, nil)
	for availabilitySetsClientNewListPager.More() {
		_, err := availabilitySetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AvailabilitySets_Get
	fmt.Println("Call operation: AvailabilitySets_Get")
	_, err = availabilitySetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.availabilitySetName, nil)
	testsuite.Require().NoError(err)

	// From step AvailabilitySets_ListAvailableSizes
	fmt.Println("Call operation: AvailabilitySets_ListAvailableSizes")
	availabilitySetsClientNewListAvailableSizesPager := availabilitySetsClient.NewListAvailableSizesPager(testsuite.resourceGroupName, testsuite.availabilitySetName, nil)
	for availabilitySetsClientNewListAvailableSizesPager.More() {
		_, err := availabilitySetsClientNewListAvailableSizesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AvailabilitySets_Update
	fmt.Println("Call operation: AvailabilitySets_Update")
	_, err = availabilitySetsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.availabilitySetName, armcompute.AvailabilitySetUpdate{}, nil)
	testsuite.Require().NoError(err)

	// From step AvailabilitySets_Delete
	fmt.Println("Call operation: AvailabilitySets_Delete")
	_, err = availabilitySetsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.availabilitySetName, nil)
	testsuite.Require().NoError(err)
}

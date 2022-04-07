//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armsubscriptions_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SubscriptionsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *SubscriptionsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/resources/armsubscriptions/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *SubscriptionsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSubscriptionsClient(t *testing.T) {
	suite.Run(t, new(SubscriptionsClientTestSuite))
}

func (testsuite *SubscriptionsClientTestSuite) TestSubscriptionsCRUD() {
	// get
	subscriptionsClient := armsubscriptions.NewClient(testsuite.cred, testsuite.options)
	resp, err := subscriptionsClient.Get(testsuite.ctx, testsuite.subscriptionID, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(testsuite.subscriptionID, *resp.SubscriptionID)

	// list
	list := subscriptionsClient.List(nil)
	testsuite.Require().NoError(list.Err())

	// list locations
	listLocations, err := subscriptionsClient.ListLocations(testsuite.ctx, testsuite.subscriptionID, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Greater(len(listLocations.Value), 0)

	// check resource
	subscriptionClient := armsubscriptions.NewSubscriptionClient(testsuite.cred, testsuite.options)
	resourceName, err := subscriptionClient.CheckResourceName(context.Background(), &armsubscriptions.SubscriptionClientCheckResourceNameOptions{
		ResourceNameDefinition: &armsubscriptions.ResourceName{
			Name: to.StringPtr("go-test-subnet"),
			Type: to.StringPtr("Microsoft.Compute"),
		},
	})
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(armsubscriptions.ResourceNameStatusAllowed, *resourceName.Status)
}

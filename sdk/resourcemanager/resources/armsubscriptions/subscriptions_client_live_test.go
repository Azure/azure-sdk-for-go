// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armsubscriptions_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/stretchr/testify/suite"
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
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), pathToPackage)
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
	fmt.Println("Call operation: Subscriptions_Get")
	subscriptionsClient, err := armsubscriptions.NewClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = subscriptionsClient.Get(testsuite.ctx, testsuite.subscriptionID, nil)
	testsuite.Require().NoError(err)

	// list
	fmt.Println("Call operation: Subscriptions_List")
	list := subscriptionsClient.NewListPager(nil)
	testsuite.Require().True(list.More())

	// list locations
	fmt.Println("Call operation: Subscriptions_ListLocations")
	listLocations := subscriptionsClient.NewListLocationsPager(testsuite.subscriptionID, nil)
	testsuite.Require().True(listLocations.More())

	// check resource
	fmt.Println("Call operation: Subscriptions_CheckResourceName")
	subscriptionClient, err := armsubscriptions.NewSubscriptionClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceName, err := subscriptionClient.CheckResourceName(context.Background(), &armsubscriptions.SubscriptionClientCheckResourceNameOptions{
		ResourceNameDefinition: &armsubscriptions.ResourceName{
			Name: to.Ptr("go-test-subnet"),
			Type: to.Ptr("Microsoft.Compute"),
		},
	})
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(armsubscriptions.ResourceNameStatusAllowed, *resourceName.Status)
}

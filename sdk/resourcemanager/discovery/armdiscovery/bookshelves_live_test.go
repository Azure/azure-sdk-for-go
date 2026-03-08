// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type BookshelvesTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	bookshelfName     string
}

func (testsuite *BookshelvesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.bookshelfName = "test-bookshelf-go01"
}

func (testsuite *BookshelvesTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestBookshelvesTestSuite(t *testing.T) {
	suite.Run(t, new(BookshelvesTestSuite))
}

// Test listing bookshelves by subscription
func (testsuite *BookshelvesTestSuite) TestBookshelvesListBySubscription() {
	fmt.Println("Call operation: Bookshelves_ListBySubscription")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewBookshelvesClient().NewListBySubscriptionPager(nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

// Test listing bookshelves by resource group
func (testsuite *BookshelvesTestSuite) TestBookshelvesListByResourceGroup() {
	fmt.Println("Call operation: Bookshelves_ListByResourceGroup")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewBookshelvesClient().NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

// Test creating a bookshelf
func (testsuite *BookshelvesTestSuite) TestBookshelvesCreateOrUpdate() {
	fmt.Println("Call operation: Bookshelves_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	subnetDefault := "/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/virtualNetworks/newapiv/subnets/default"
	subnetDefault2 := "/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/virtualNetworks/newapiv/subnets/default2"

	bookshelvesClient := clientFactory.NewBookshelvesClient()
	poller, err := bookshelvesClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.bookshelfName,
		armdiscovery.Bookshelf{
			Location: to.Ptr(testsuite.location),
			Properties: &armdiscovery.BookshelfProperties{
				PrivateEndpointSubnetID: to.Ptr(subnetDefault),
				SearchSubnetID:          to.Ptr(subnetDefault2),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	bookshelf, err := poller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(bookshelf.ID)
	fmt.Println("Created bookshelf:", *bookshelf.Name)
}

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

	// Add EUAP redirect policy
	euapOptions := GetEUAPClientOptions()
	testsuite.options.PerCallPolicies = append(testsuite.options.PerCallPolicies, euapOptions.PerCallPolicies...)

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "discovery-test-rg")
	testsuite.bookshelfName = "test-bookshelf"

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	fmt.Println("Created resource group:", testsuite.resourceGroupName)
}

func (testsuite *BookshelvesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
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

// Test bookshelf CRUD operations
func (testsuite *BookshelvesTestSuite) SkipTestBookshelvesCRUD() {
	fmt.Println("Call operation: Bookshelves_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	bookshelvesClient := clientFactory.NewBookshelvesClient()

	// Create bookshelf
	poller, err := bookshelvesClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.bookshelfName,
		armdiscovery.Bookshelf{
			Location:   to.Ptr(testsuite.location),
			Properties: &armdiscovery.BookshelfProperties{},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	bookshelf, err := poller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(bookshelf.ID)
	fmt.Println("Created bookshelf:", *bookshelf.Name)

	// Get bookshelf
	fmt.Println("Call operation: Bookshelves_Get")
	getResp, err := bookshelvesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.bookshelfName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(testsuite.bookshelfName, *getResp.Name)

	// Update bookshelf
	fmt.Println("Call operation: Bookshelves_Update")
	updatePoller, err := bookshelvesClient.BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.bookshelfName,
		armdiscovery.Bookshelf{
			Location: to.Ptr(testsuite.location),
			Tags: map[string]*string{
				"environment": to.Ptr("test"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	updateResp, err := updatePoller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(updateResp.ID)

	// Delete bookshelf
	fmt.Println("Call operation: Bookshelves_Delete")
	delPoller, err := bookshelvesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.bookshelfName, nil)
	testsuite.Require().NoError(err)
	_, err = delPoller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	fmt.Println("Deleted bookshelf:", testsuite.bookshelfName)
}

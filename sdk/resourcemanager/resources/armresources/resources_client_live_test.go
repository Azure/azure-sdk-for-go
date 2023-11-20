//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type ResourcesClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *ResourcesClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/resources/armresources/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *ResourcesClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestResourcesClient(t *testing.T) {
	suite.Run(t, new(ResourcesClientTestSuite))
}

func (testsuite *ResourcesClientTestSuite) TestResourcesCRUD() {
	// check existence resource
	fmt.Println("Call operation: Resources_CheckExistence")
	resourcesClient, err := armresources.NewClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceName := "go-test-resource"
	check, err := resourcesClient.CheckExistence(
		testsuite.ctx,
		testsuite.resourceGroupName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().False(check.Success)

	// check existence resource by id
	fmt.Println("Call operation: Resources_CheckExistenceByID")
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID, "{guid}", testsuite.subscriptionID)
	resourceID = strings.ReplaceAll(resourceID, "{resourcegroupname}", testsuite.resourceGroupName)
	resourceID = strings.ReplaceAll(resourceID, "{resourceprovidernamespace}", "Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID, "{resourcetype}", "availabilitySets")
	resourceID = strings.ReplaceAll(resourceID, "{resourcename}", resourceName)
	checkByID, err := resourcesClient.CheckExistenceByID(testsuite.ctx, resourceID, "2021-07-01", nil)
	testsuite.Require().NoError(err)
	testsuite.Require().False(checkByID.Success)

	// create resource
	fmt.Println("Call operation: Resources_CreateOrUpdate")
	createPoller, err := resourcesClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Location: to.Ptr("eastus"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	resp, err := testutil.PollForTest(testsuite.ctx, createPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(resourceName, *resp.Name)

	// create resource by id
	fmt.Println("Call operation: Resources_CreateOrUpdateByID")
	createByIDPoller, err := resourcesClient.BeginCreateOrUpdateByID(
		testsuite.ctx,
		resourceID,
		"2021-07-01",
		armresources.GenericResource{
			Location: to.Ptr("eastus"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	createByIDResp, err := testutil.PollForTest(testsuite.ctx, createByIDPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(resourceName, *createByIDResp.Name)

	// update resources
	fmt.Println("Call operation: Resources_Update")
	updatePoller, err := resourcesClient.BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2019-07-01",
		armresources.GenericResource{
			Tags: map[string]*string{
				"tag1": to.Ptr("value1"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	updateResp, err := testutil.PollForTest(testsuite.ctx, updatePoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("value1", *updateResp.Tags["tag1"])

	// update resource by id
	fmt.Println("Call operation: Resources_UpdateByID")
	updateByIDPoller, err := resourcesClient.BeginUpdateByID(
		testsuite.ctx,
		resourceID,
		"2019-07-01",
		armresources.GenericResource{
			Tags: map[string]*string{
				"key2": to.Ptr("value2"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	updateByIDResp, err := testutil.PollForTest(testsuite.ctx, updateByIDPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("value2", *updateByIDResp.Tags["key2"])

	// get resource
	fmt.Println("Call operation: Resources_Get")
	getResp, err := resourcesClient.Get(
		testsuite.ctx,
		testsuite.resourceGroupName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2019-07-01",
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(resourceName, *getResp.Name)

	// get resource by id
	fmt.Println("Call operation: Resources_GetByID")
	getByIDResp, err := resourcesClient.GetByID(
		testsuite.ctx,
		resourceID,
		"2019-07-01",
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(resourceName, *getByIDResp.Name)

	// list resource
	fmt.Println("Call operation: Resources_List")
	listPager := resourcesClient.NewListPager(nil)
	testsuite.Require().NoError(err)
	testsuite.Require().True(listPager.More())

	// list resource by resource group
	fmt.Println("Call operation: Resources_ListByResourceGroup")
	listByResourceGroup := resourcesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().True(listByResourceGroup.More())

	// delete resource
	fmt.Println("Call operation: Resources_Delete")
	delPoller, err := resourcesClient.BeginDelete(
		testsuite.ctx,
		testsuite.resourceGroupName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2019-07-01",
		nil,
	)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)

	// delete resources by id
	fmt.Println("Call operation: Resources_DeleteByID")
	delByIDPoller, err := resourcesClient.BeginDeleteByID(
		testsuite.ctx,
		resourceID,
		"2019-07-01",
		nil,
	)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, delByIDPoller)
	testsuite.Require().NoError(err)

	// // clean new resource group
	// delPollerResp, err := rgClient.BeginDelete(context.Background(), rgName, nil)
	// testsuite.Require().NoError(err)
	// _, err = testutil.PollForTest(testsuite.ctx, delPollerResp)
	// testsuite.Require().NoError(err)
}

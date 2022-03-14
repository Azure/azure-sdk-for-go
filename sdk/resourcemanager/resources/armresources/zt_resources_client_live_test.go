//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
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
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
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
	resourcesClient := armresources.NewClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
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
	createPoller, err := resourcesClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Location: to.StringPtr("westus"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var resp armresources.ClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = createPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if createPoller.Poller.Done() {
				resp, err = createPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		resp, err = createPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(resourceName, *resp.Name)

	// create resource by id
	createByIDPoller, err := resourcesClient.BeginCreateOrUpdateByID(
		testsuite.ctx,
		resourceID,
		"2021-07-01",
		armresources.GenericResource{
			Location: to.StringPtr("westus"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var createByIDResp armresources.ClientCreateOrUpdateByIDResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = createByIDPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if createByIDPoller.Poller.Done() {
				createByIDResp, err = createByIDPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		createByIDResp, err = createByIDPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(resourceName, *createByIDResp.Name)

	// update resources
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
				"tag1": to.StringPtr("value1"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var updateResp armresources.ClientUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = updatePoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if updatePoller.Poller.Done() {
				updateResp, err = updatePoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		updateResp, err = updatePoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal("value1", *updateResp.Tags["tag1"])

	// update resource by id
	updateByIDPoller, err := resourcesClient.BeginUpdateByID(
		testsuite.ctx,
		resourceID,
		"2019-07-01",
		armresources.GenericResource{
			Tags: map[string]*string{
				"key2": to.StringPtr("value2"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var updateByIDResp armresources.ClientUpdateByIDResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = updateByIDPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if updateByIDPoller.Poller.Done() {
				updateByIDResp, err = updateByIDPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		updateByIDResp, err = updateByIDPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal("value2", *updateByIDResp.Tags["key2"])

	// get resource
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
	getByIDResp, err := resourcesClient.GetByID(
		testsuite.ctx,
		resourceID,
		"2019-07-01",
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(resourceName, *getByIDResp.Name)

	// list resource
	listPager := resourcesClient.List(nil)
	testsuite.Require().NoError(err)
	testsuite.Require().True(listPager.NextPage(testsuite.ctx))

	// list resource by resource group
	listByResourceGroup := resourcesClient.ListByResourceGroup(testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().True(listByResourceGroup.NextPage(testsuite.ctx))

	// create resource group
	rgName := "go-test-rg"
	rgClient := armresources.NewResourceGroupsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	newRg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(rgName, *newRg.Name)

	// validate move resources
	validatePoller, err := resourcesClient.BeginValidateMoveResources(
		testsuite.ctx,
		testsuite.resourceGroupName,
		armresources.MoveInfo{
			Resources: []*string{
				to.StringPtr(*resp.ID),
			},
			TargetResourceGroup: to.StringPtr(*newRg.ID),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var validateResp armresources.ClientValidateMoveResourcesResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = validatePoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if validatePoller.Poller.Done() {
				validateResp, err = validatePoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		validateResp, err = validatePoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(204, validateResp.RawResponse.StatusCode)

	// move resources
	movePoller, err := resourcesClient.BeginMoveResources(
		testsuite.ctx,
		testsuite.resourceGroupName,
		armresources.MoveInfo{
			Resources: []*string{
				to.StringPtr(*resp.ID),
			},
			TargetResourceGroup: to.StringPtr(*newRg.ID),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var moveResp armresources.ClientMoveResourcesResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = movePoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if movePoller.Poller.Done() {
				moveResp, err = movePoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		moveResp, err = movePoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(204, moveResp.RawResponse.StatusCode)

	// delete resource
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
	var delResource armresources.ClientDeleteResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = delPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if delPoller.Poller.Done() {
				delResource, err = delPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		delResource, err = delPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(204, delResource.RawResponse.StatusCode)

	// delete resources by id
	delByIDPoller, err := resourcesClient.BeginDeleteByID(
		testsuite.ctx,
		resourceID,
		"2019-07-01",
		nil,
	)
	testsuite.Require().NoError(err)
	var delByIDResp armresources.ClientDeleteByIDResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = delByIDPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if delByIDPoller.Poller.Done() {
				delByIDResp, err = delByIDPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		delByIDResp, err = delByIDPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(204, delByIDResp.RawResponse.StatusCode)

	// clean new resource group
	delPollerResp, err := rgClient.BeginDelete(context.Background(), rgName, nil)
	testsuite.Require().NoError(err)
	var delResp armresources.ResourceGroupsClientDeleteResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = delPollerResp.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if delPollerResp.Poller.Done() {
				delResp, err = delPollerResp.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		delResp, err = delPollerResp.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}

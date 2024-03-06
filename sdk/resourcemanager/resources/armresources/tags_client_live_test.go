//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type TagsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *TagsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/resources/armresources/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *TagsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestTagsClient(t *testing.T) {
	suite.Run(t, new(TagsClientTestSuite))
}

func (testsuite *TagsClientTestSuite) TestTagsCRUD() {
	// create tag
	fmt.Println("Call operation: Tags_CreateOrUpdate")
	tagsClient, err := armresources.NewTagsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tagName := "go-test-tags"
	resp, err := tagsClient.CreateOrUpdate(testsuite.ctx, tagName, &armresources.TagsClientCreateOrUpdateOptions{})
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(tagName, *resp.TagName)

	// create tag value
	fmt.Println("Call operation: Tags_CreateOrUpdateValue")
	valueName := "go-test-value"
	valueResp, err := tagsClient.CreateOrUpdateValue(testsuite.ctx, tagName, valueName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(valueName, *valueResp.TagValue.TagValue)

	// list
	fmt.Println("Call operation: Tags_List")
	listPager := tagsClient.NewListPager(nil)
	testsuite.Require().True(listPager.More())

	// delete tag value
	fmt.Println("Call operation: Tags_DeleteValue")
	_, err = tagsClient.DeleteValue(testsuite.ctx, tagName, valueName, nil)
	testsuite.Require().NoError(err)

	// delete tag
	fmt.Println("Call operation: Tags_Delete")
	_, err = tagsClient.Delete(testsuite.ctx, tagName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *TagsClientTestSuite) TestTagsCRUDAtScope() {
	// create scope
	fmt.Println("Call operation: Tags_CreateOrUpdateAtScope")
	tagsClient, err := armresources.NewTagsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scopeName := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v", testsuite.subscriptionID, testsuite.resourceGroupName)
	resp, err := tagsClient.CreateOrUpdateAtScope(
		testsuite.ctx,
		scopeName,
		armresources.TagsResource{
			Properties: &armresources.Tags{
				Tags: map[string]*string{
					"tagKey1": to.Ptr("tagValue1"),
					"tagKey2": to.Ptr("tagValue2"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("default", *resp.Name)

	// update at scope
	fmt.Println("Call operation: Tags_UpdateAtScope")
	updateResp, err := tagsClient.UpdateAtScope(
		testsuite.ctx,
		scopeName,
		armresources.TagsPatchResource{
			Operation: to.Ptr(armresources.TagsPatchOperationDelete),
			Properties: &armresources.Tags{
				Tags: map[string]*string{
					"tagKey1": to.Ptr("tagKey1"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("default", *updateResp.Name)

	// get at scopes
	fmt.Println("Call operation: Tags_GetAtScope")
	getResp, err := tagsClient.GetAtScope(testsuite.ctx, scopeName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("default", *getResp.Name)

	// delete at scope
	fmt.Println("Call operation: Tags_DeleteAtScope")
	_, err = tagsClient.DeleteAtScope(testsuite.ctx, scopeName, nil)
	testsuite.Require().NoError(err)
}

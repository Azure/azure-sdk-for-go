//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
	"testing"
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
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
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
	tagsClient := armresources.NewTagsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	tagName := "go-test-tags"
	resp, err := tagsClient.CreateOrUpdate(testsuite.ctx, tagName, &armresources.TagsClientCreateOrUpdateOptions{})
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(tagName, *resp.TagName)

	// create tag value
	valueName := "go-test-value"
	valueResp, err := tagsClient.CreateOrUpdateValue(testsuite.ctx, tagName, valueName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(valueName, *valueResp.TagValue.TagValue)

	// list
	listPager := tagsClient.List(nil)
	testsuite.Require().NoError(listPager.Err())
	testsuite.Require().True(listPager.NextPage(testsuite.ctx))

	// delete tag value
	delResp, err := tagsClient.DeleteValue(testsuite.ctx, tagName, valueName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)

	// delete tag
	delTag, err := tagsClient.Delete(testsuite.ctx, tagName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delTag.RawResponse.StatusCode)
}

func (testsuite *TagsClientTestSuite) TestTagsCRUDAtScope() {
	// create scope
	tagsClient := armresources.NewTagsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	scopeName := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v", testsuite.subscriptionID, testsuite.resourceGroupName)
	resp, err := tagsClient.CreateOrUpdateAtScope(
		testsuite.ctx,
		scopeName,
		armresources.TagsResource{
			Properties: &armresources.Tags{
				Tags: map[string]*string{
					"tagKey1": to.StringPtr("tagValue1"),
					"tagKey2": to.StringPtr("tagValue2"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("default", *resp.Name)

	// update at scope
	updateResp, err := tagsClient.UpdateAtScope(
		testsuite.ctx,
		scopeName,
		armresources.TagsPatchResource{
			Operation: armresources.TagsPatchOperationDelete.ToPtr(),
			Properties: &armresources.Tags{
				Tags: map[string]*string{
					"tagKey1": to.StringPtr("tagKey1"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("default", *updateResp.Name)

	// get at scopes
	getResp, err := tagsClient.GetAtScope(testsuite.ctx, scopeName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("default", *getResp.Name)

	// delete at scope
	delResp, err := tagsClient.DeleteAtScope(testsuite.ctx, scopeName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}

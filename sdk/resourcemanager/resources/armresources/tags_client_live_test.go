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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
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
	testutil.StartRecording(testsuite.T(), pathToPackage)
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
	_, err = tagsClient.CreateOrUpdate(testsuite.ctx, tagName, &armresources.TagsClientCreateOrUpdateOptions{})
	testsuite.Require().NoError(err)

	// create tag value
	fmt.Println("Call operation: Tags_CreateOrUpdateValue")
	valueName := "go-test-value"
	_, err = tagsClient.CreateOrUpdateValue(testsuite.ctx, tagName, valueName, nil)
	testsuite.Require().NoError(err)

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
	pollerResp0, err := tagsClient.BeginCreateOrUpdateAtScope(
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
	_, err = testutil.PollForTest(testsuite.ctx, pollerResp0)
	testsuite.Require().NoError(err)
	// update at scope
	fmt.Println("Call operation: Tags_UpdateAtScope")
	pollerResp, err := tagsClient.BeginUpdateAtScope(
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
	_, err = testutil.PollForTest(testsuite.ctx, pollerResp)
	testsuite.Require().NoError(err)

	// get at scopes
	fmt.Println("Call operation: Tags_GetAtScope")
	_, err = tagsClient.GetAtScope(testsuite.ctx, scopeName, nil)
	testsuite.Require().NoError(err)

	// delete at scope
	fmt.Println("Call operation: Tags_DeleteAtScope")
	pollerResp2, err := tagsClient.BeginDeleteAtScope(testsuite.ctx, scopeName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, pollerResp2)
	testsuite.Require().NoError(err)
}

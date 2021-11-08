//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTagsClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create tag
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	tagName, err := createRandomName(t, "tag")
	require.NoError(t, err)
	resp, err := tagsClient.CreateOrUpdate(ctx, tagName, &armresources.TagsCreateOrUpdateOptions{})
	require.NoError(t, err)
	require.Equal(t, *resp.TagName, tagName)
}

func TestTagsClient_CreateOrUpdateValue(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create tag
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	tagName, err := createRandomName(t, "tag")
	require.NoError(t, err)
	resp, err := tagsClient.CreateOrUpdate(ctx, tagName, &armresources.TagsCreateOrUpdateOptions{})
	require.NoError(t, err)
	require.Equal(t, *resp.TagName, tagName)

	// create tag value
	valueName, err := createRandomName(t, "value")
	valueResp, err := tagsClient.CreateOrUpdateValue(ctx, tagName, valueName, nil)
	require.NoError(t, err)
	require.Equal(t, *valueResp.TagValue.TagValue, valueName)
}

func TestTagsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create tag
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	tagName, err := createRandomName(t, "tag")
	require.NoError(t, err)
	resp, err := tagsClient.CreateOrUpdate(ctx, tagName, &armresources.TagsCreateOrUpdateOptions{})
	require.NoError(t, err)
	require.Equal(t, *resp.TagName, tagName)

	listPager := tagsClient.List(nil)
	require.Equal(t, true, listPager.NextPage(ctx))
}

func TestTagsClient_DeleteValue(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create tag
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	tagName, err := createRandomName(t, "tag")
	require.NoError(t, err)
	resp, err := tagsClient.CreateOrUpdate(ctx, tagName, &armresources.TagsCreateOrUpdateOptions{})
	require.NoError(t, err)
	require.Equal(t, *resp.TagName, tagName)

	// create tag value
	tagValue, err := createRandomName(t, "value")
	valueResp, err := tagsClient.CreateOrUpdateValue(ctx, tagName, tagValue, nil)
	require.NoError(t, err)
	require.Equal(t, *valueResp.TagValue.TagValue, tagValue)

	// delete tag value
	delResp, err := tagsClient.DeleteValue(ctx, tagName, tagValue, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)
}

func TestTagsClient_Delete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create tag
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	tagName, err := createRandomName(t, "tag")
	require.NoError(t, err)
	resp, err := tagsClient.CreateOrUpdate(ctx, tagName, &armresources.TagsCreateOrUpdateOptions{})
	require.NoError(t, err)
	require.Equal(t, *resp.TagName, tagName)

	delResp, err := tagsClient.Delete(ctx, tagName, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)
}

func TestTagsClient_CreateOrUpdateAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createVMSS","westus2")
	defer clean()

	// create scope
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	scope := *rg.ID
	resp, err := tagsClient.CreateOrUpdateAtScope(
		ctx,
		scope,
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
	require.NoError(t, err)
	require.Equal(t, "default", *resp.Name)
}

func TestTagsClient_GetAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createVMSS","westus2")
	defer clean()

	// create scope
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	scope := *rg.ID
	resp, err := tagsClient.CreateOrUpdateAtScope(
		ctx,
		scope,
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
	require.NoError(t, err)
	require.Equal(t, "tagValue1", *resp.Properties.Tags["tagKey1"])

	// get at scopes
	getResp, err := tagsClient.GetAtScope(ctx, scope, nil)
	require.NoError(t, err)
	require.Equal(t, "default", *getResp.Name)
}

func TestTagsClient_UpdateAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createVMSS","westus2")
	defer clean()

	// create scope
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	scope := *rg.ID
	resp, err := tagsClient.CreateOrUpdateAtScope(
		ctx,
		scope,
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
	require.NoError(t, err)
	require.Equal(t, "tagValue1", *resp.Properties.Tags["tagKey1"])

	// update at scope
	updateResp, err := tagsClient.UpdateAtScope(
		ctx,
		scope,
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
	require.NoError(t, err)
	require.Equal(t, "default", *updateResp.Name)
}

func TestTagsClient_DeleteAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createVMSS","westus2")
	defer clean()

	// create scope
	tagsClient := armresources.NewTagsClient(subscriptionID, cred, opt)
	scope := *rg.ID
	resp, err := tagsClient.CreateOrUpdateAtScope(
		ctx,
		scope,
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
	require.NoError(t, err)
	require.Equal(t, "tagValue1", *resp.Properties.Tags["tagKey1"])

	// delete at scope
	delResp, err := tagsClient.DeleteAtScope(ctx, scope, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)
}

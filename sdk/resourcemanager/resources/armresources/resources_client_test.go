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
	"strings"
	"testing"
	"time"
)

func TestResourcesClient_CheckExistence(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"check","westus")
	defer clean()
	rgName := *rg.Name

	// check existence resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	resp,err := resourcesClient.CheckExistence(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		nil,
		)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestResourcesClient_CheckExistenceByID(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"check2","westus")
	defer clean()
	rgName := *rg.Name

	// check existence resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID,"{guid}",subscriptionID)
	resourceID = strings.ReplaceAll(resourceID,"{resourcegroupname}",rgName)
	resourceID = strings.ReplaceAll(resourceID,"{resourceprovidernamespace}","Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID,"{resourcetype}","availabilitySets")
	resourceID = strings.ReplaceAll(resourceID,"{resourcename}",resourceName)
	resp,err := resourcesClient.CheckExistenceByID(ctx, resourceID, "2021-07-01", nil)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestResourcesClient_BeginCreateOrUpdateByID(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createByID","westus")
	defer clean()
	rgName := *rg.Name

	// create resource by id
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID,"{guid}",subscriptionID)
	resourceID = strings.ReplaceAll(resourceID,"{resourcegroupname}",rgName)
	resourceID = strings.ReplaceAll(resourceID,"{resourceprovidernamespace}","Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID,"{resourcetype}","availabilitySets")
	resourceID = strings.ReplaceAll(resourceID,"{resourcename}",resourceName)
	pollerResp,err := resourcesClient.BeginCreateOrUpdateByID(
		ctx,
		resourceID,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
		)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)
}

func TestResourcesClient_BeginCreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
	defer clean()
	rgName := *rg.Name

	// create resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	pollerResp,err := resourcesClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)
}

func TestResourcesClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"get","westus")
	defer clean()
	rgName := *rg.Name

	// create resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	pollerResp,err := resourcesClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	// get resource
	getResp,err := resourcesClient.Get(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2019-07-01",
		nil,
		)
	require.NoError(t, err)
	require.Equal(t, resourceName,*getResp.Name)
}

func TestResourcesClient_GetByID(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"getByID","westus")
	defer clean()
	rgName := *rg.Name

	// create resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	pollerResp,err := resourcesClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	// get resource
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID,"{guid}",subscriptionID)
	resourceID = strings.ReplaceAll(resourceID,"{resourcegroupname}",rgName)
	resourceID = strings.ReplaceAll(resourceID,"{resourceprovidernamespace}","Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID,"{resourcetype}","availabilitySets")
	resourceID = strings.ReplaceAll(resourceID,"{resourcename}",resourceName)
	getResp,err := resourcesClient.GetByID(
		ctx,
		resourceID,
		"2019-07-01",
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, resourceName,*getResp.Name)
}

func TestResourcesClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"list","westus")
	defer clean()
	rgName := *rg.Name

	// create resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	pollerResp,err := resourcesClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	// list resource
	listPager := resourcesClient.List(nil)
	require.NoError(t, err)
	require.True(t, listPager.NextPage(ctx))
}

func TestResourcesClient_ListByResourceGroup(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"listBy","westus")
	defer clean()
	rgName := *rg.Name

	// create resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	pollerResp,err := resourcesClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	// list resource
	listPager := resourcesClient.ListByResourceGroup(rgName,nil)
	require.NoError(t, err)
	require.True(t, listPager.NextPage(ctx))
}

func  TestResourcesClient_BeginValidateMoveResources(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create new resource group
	newRg,clean2 := createResourceGroup(t,cred,opt,subscriptionID,"move2","westus")
	defer clean2()

	// create resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	pollerResp,err := resourcesClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	movePoller,err := resourcesClient.BeginValidateMoveResources(
		ctx,
		rgName,
		armresources.ResourcesMoveInfo{
			Resources: []*string{
				to.StringPtr(*resp.ID),
			},
			TargetResourceGroup: to.StringPtr(*newRg.ID),
		},
		nil,
		)
	require.NoError(t, err)
	moveResp,err := movePoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, 204,moveResp.RawResponse.StatusCode)
}

func TestResourcesClient_BeginMoveResources(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create new resource group
	newRg,clean2 := createResourceGroup(t,cred,opt,subscriptionID,"move2","westus")
	defer clean2()

	// create resource
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	pollerResp,err := resourcesClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	movePoller,err := resourcesClient.BeginMoveResources(
		ctx,
		rgName,
		armresources.ResourcesMoveInfo{
			Resources: []*string{
				to.StringPtr(*resp.ID),
			},
			TargetResourceGroup: to.StringPtr(*newRg.ID),
		},
		nil,
	)
	require.NoError(t, err)
	moveResp,err := movePoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, 204,moveResp.RawResponse.StatusCode)
}

func TestResourcesClient_BeginUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createByID","westus")
	defer clean()
	rgName := *rg.Name

	// create resource by id
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID,"{guid}",subscriptionID)
	resourceID = strings.ReplaceAll(resourceID,"{resourcegroupname}",rgName)
	resourceID = strings.ReplaceAll(resourceID,"{resourceprovidernamespace}","Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID,"{resourcetype}","availabilitySets")
	resourceID = strings.ReplaceAll(resourceID,"{resourcename}",resourceName)
	pollerResp,err := resourcesClient.BeginCreateOrUpdateByID(
		ctx,
		resourceID,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	updatePoller,err := resourcesClient.BeginUpdate(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2019-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Tags: map[string]*string{
					"tag1": to.StringPtr("value1"),
				},
			},
		},
		nil,
		)
	require.NoError(t, err)
	updateResp,err := updatePoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, resourceName,*updateResp.Name)
}

func TestResourcesClient_BeginUpdateByID(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createByID","westus")
	defer clean()
	rgName := *rg.Name

	// create resource by id
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID,"{guid}",subscriptionID)
	resourceID = strings.ReplaceAll(resourceID,"{resourcegroupname}",rgName)
	resourceID = strings.ReplaceAll(resourceID,"{resourceprovidernamespace}","Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID,"{resourcetype}","availabilitySets")
	resourceID = strings.ReplaceAll(resourceID,"{resourcename}",resourceName)
	pollerResp,err := resourcesClient.BeginCreateOrUpdateByID(
		ctx,
		resourceID,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	updatePoller,err := resourcesClient.BeginUpdateByID(
		ctx,
		resourceID,
		"2019-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Tags: map[string]*string{
					"tag1": to.StringPtr("value1"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	updateResp,err := updatePoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, resourceName,*updateResp.Name)
}

func TestResourcesClient_BeginDelete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createByID","westus")
	defer clean()
	rgName := *rg.Name

	// create resource by id
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID,"{guid}",subscriptionID)
	resourceID = strings.ReplaceAll(resourceID,"{resourcegroupname}",rgName)
	resourceID = strings.ReplaceAll(resourceID,"{resourceprovidernamespace}","Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID,"{resourcetype}","availabilitySets")
	resourceID = strings.ReplaceAll(resourceID,"{resourcename}",resourceName)
	pollerResp,err := resourcesClient.BeginCreateOrUpdateByID(
		ctx,
		resourceID,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	delPoller,err := resourcesClient.BeginDelete(
		ctx,
		rgName,
		"Microsoft.Compute",
		"",
		"availabilitySets",
		resourceName,
		"2019-07-01",
		nil,
		)
	delResp,err := delPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, 200,delResp.RawResponse.StatusCode)
}

func TestResourcesClient_BeginDeleteByID(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"createByID","westus")
	defer clean()
	rgName := *rg.Name

	// create resource by id
	resourcesClient := armresources.NewResourcesClient(subscriptionID,cred,opt)
	resourceName,err := createRandomName(t,"rs")
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID,"{guid}",subscriptionID)
	resourceID = strings.ReplaceAll(resourceID,"{resourcegroupname}",rgName)
	resourceID = strings.ReplaceAll(resourceID,"{resourceprovidernamespace}","Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID,"{resourcetype}","availabilitySets")
	resourceID = strings.ReplaceAll(resourceID,"{resourcename}",resourceName)
	pollerResp,err := resourcesClient.BeginCreateOrUpdateByID(
		ctx,
		resourceID,
		"2021-07-01",
		armresources.GenericResource{
			Resource: armresources.Resource{
				Location: to.StringPtr("westus"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, resourceName,*resp.Name)

	delPoller,err := resourcesClient.BeginDeleteByID(
		ctx,
		resourceID,
		"2019-07-01",
		nil,
	)
	delResp,err := delPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, 200,delResp.RawResponse.StatusCode)
}
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
	"time"
)

func TestResourceGroupsClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rgName, err := createRandomName(t, "createRG")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID,cred,opt)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	require.NoError(t, err)
	require.Equal(t, *rg.Name,rgName)

	// clean resource group
	cleanup(t,rgClient,rgName)
}

func TestResourceGroupsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rgName, err := createRandomName(t, "getRG")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID,cred,opt)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	require.NoError(t, err)
	require.Equal(t, *rg.Name,rgName)
	defer cleanup(t,rgClient,rgName)

	// get resource group
	getResp,err := rgClient.Get(context.Background(),rgName,nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Name,rgName)
}

func TestResourceGroupsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rgName, err := createRandomName(t, "listRG")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID,cred,opt)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	require.NoError(t, err)
	require.Equal(t, *rg.Name,rgName)
	defer cleanup(t,rgClient,rgName)

	// list resource group
	listPager := rgClient.List(nil)
	require.Equal(t,listPager.NextPage(context.Background()),true)
}

func TestResourceGroupsClient_Update(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rgName, err := createRandomName(t, "updateRG")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID,cred,opt)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	require.NoError(t, err)
	require.Equal(t, *rg.Name,rgName)
	defer cleanup(t,rgClient,rgName)

	// update resource group
	resp,err := rgClient.Update(context.Background(),rgName,armresources.ResourceGroupPatchable{
		Tags: map[string]*string{
			"key": to.StringPtr("value"),
		},
	},nil)
	require.NoError(t, err)
	require.Equal(t, *resp.Name,rgName)
}

func TestResourceGroupsClient_BeginDelete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rgName, err := createRandomName(t, "deleteRG")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID,cred,opt)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	require.NoError(t, err)
	require.Equal(t, *rg.Name,rgName)

	// clean resource group
	pollerResp,err := rgClient.BeginDelete(context.Background(),rgName,nil)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(context.Background(),10*time.Second)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode,200)
}

func TestResourceGroupsClient_CheckExistence(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rgName, err := createRandomName(t, "checkRG")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID,cred,opt)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	require.NoError(t, err)
	require.Equal(t, *rg.Name,rgName)
	defer cleanup(t,rgClient,rgName)

	// check existence resource group
	resp,err := rgClient.CheckExistence(context.Background(),rgName,nil)
	require.NoError(t, err)
	require.Equal(t, resp.Success,false)
}

func TestResourceGroupsClient_BeginExportTemplate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rgName, err := createRandomName(t, "templateRG")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID,cred,opt)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	require.NoError(t, err)
	require.Equal(t, *rg.Name,rgName)
	defer cleanup(t,rgClient,rgName)

	// export template resource group
	pollerResp,err := rgClient.BeginExportTemplate(context.Background(),rgName,armresources.ExportTemplateRequest{
		Resources: []*string{
			to.StringPtr("*"),
		},
	},nil)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(context.Background(),10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, resp.Template)
}
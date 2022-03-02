//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armtemplatespecs_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armtemplatespecs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	require.NotNil(t, rg)
	resourceGroupName := *rg.Name

	// create template spec
	templateSpecName, err := createRandomName(t, "ts")
	require.NoError(t, err)
	templateSpecsClient := armtemplatespecs.NewClient(subscriptionID, cred, opt)
	resp, err := templateSpecsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		templateSpecName,
		armtemplatespecs.TemplateSpec{
			Location: to.StringPtr(location),
			Properties: &armtemplatespecs.TemplateSpecProperties{
				Description: to.StringPtr("template spec properties description."),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, templateSpecName, *resp.Name)

	// get template spec
	getResp, err := templateSpecsClient.Get(ctx, resourceGroupName, templateSpecName, nil)
	require.NoError(t, err)
	require.Equal(t, templateSpecName, *getResp.Name)

	// update template spec
	updateResp, err := templateSpecsClient.Update(ctx, resourceGroupName, templateSpecName, &armtemplatespecs.ClientUpdateOptions{
		TemplateSpec: &armtemplatespecs.TemplateSpecUpdateModel{
			Tags: map[string]*string{
				"test": to.StringPtr("recording"),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "recording", *updateResp.Tags["test"])

	// list template spec by resource group
	pager := templateSpecsClient.ListByResourceGroup(resourceGroupName, nil)
	require.NoError(t, pager.Err())

	// list template spec by subscription
	pager2 := templateSpecsClient.ListBySubscription(nil)
	require.NoError(t, pager2.Err())

	// delete template spec
	delResp, err := templateSpecsClient.Delete(ctx, resourceGroupName, templateSpecName, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)
}

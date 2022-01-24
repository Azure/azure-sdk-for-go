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

func TestTemplateSpecVersionsClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create template spec
	templateSpecName, err := createRandomName(t, "ts")
	require.NoError(t, err)
	templateSpecsClient := armtemplatespecs.NewClient(subscriptionID, cred, opt)
	resp, err := templateSpecsClient.CreateOrUpdate(
		ctx,
		rgName,
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

	// create template version
	templateSpecVersion, err := createRandomName(t, "version")
	require.NoError(t, err)
	templateSpecVersionsClient := armtemplatespecs.NewTemplateSpecVersionsClient(subscriptionID, cred, opt)
	vResp, err := templateSpecVersionsClient.CreateOrUpdate(
		ctx,
		rgName,
		templateSpecName,
		templateSpecVersion,
		armtemplatespecs.TemplateSpecVersion{
			Location: to.StringPtr(location),
			Properties: &armtemplatespecs.TemplateSpecVersionProperties{
				Description: to.StringPtr("<description>"),
				MainTemplate: map[string]interface{}{
					"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
					"contentVersion": "1.0.0.0",
					"parameters":     map[string]interface{}{},
					"resources":      []interface{}{},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, vResp)
}

func TestTemplateSpecVersionsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create template spec
	templateSpecName, err := createRandomName(t, "ts")
	require.NoError(t, err)
	templateSpecsClient := armtemplatespecs.NewClient(subscriptionID, cred, opt)
	resp, err := templateSpecsClient.CreateOrUpdate(
		ctx,
		rgName,
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

	// create template version
	templateSpecVersion, err := createRandomName(t, "version")
	require.NoError(t, err)
	templateSpecVersionsClient := armtemplatespecs.NewTemplateSpecVersionsClient(subscriptionID, cred, opt)
	vResp, err := templateSpecVersionsClient.CreateOrUpdate(
		ctx,
		rgName,
		templateSpecName,
		templateSpecVersion,
		armtemplatespecs.TemplateSpecVersion{
			Location: to.StringPtr(location),
			Properties: &armtemplatespecs.TemplateSpecVersionProperties{
				Description: to.StringPtr("<description>"),
				MainTemplate: map[string]interface{}{
					"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
					"contentVersion": "1.0.0.0",
					"parameters":     map[string]interface{}{},
					"resources":      []interface{}{},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, vResp)

	// get
	getResp, err := templateSpecVersionsClient.Get(ctx, rgName, templateSpecName, templateSpecVersion, nil)
	require.NoError(t, err)
	require.Equal(t, templateSpecVersion, *getResp.Name)
}

func TestTemplateSpecVersionsClient_Update(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create template spec
	templateSpecName, err := createRandomName(t, "ts")
	require.NoError(t, err)
	templateSpecsClient := armtemplatespecs.NewClient(subscriptionID, cred, opt)
	resp, err := templateSpecsClient.CreateOrUpdate(
		ctx,
		rgName,
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

	// create template version
	templateSpecVersion, err := createRandomName(t, "version")
	require.NoError(t, err)
	templateSpecVersionsClient := armtemplatespecs.NewTemplateSpecVersionsClient(subscriptionID, cred, opt)
	vResp, err := templateSpecVersionsClient.CreateOrUpdate(
		ctx,
		rgName,
		templateSpecName,
		templateSpecVersion,
		armtemplatespecs.TemplateSpecVersion{
			Location: to.StringPtr(location),
			Properties: &armtemplatespecs.TemplateSpecVersionProperties{
				Description: to.StringPtr("<description>"),
				MainTemplate: map[string]interface{}{
					"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
					"contentVersion": "1.0.0.0",
					"parameters":     map[string]interface{}{},
					"resources":      []interface{}{},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, vResp)

	// update
	updateResp, err := templateSpecVersionsClient.Update(ctx, rgName, templateSpecName, templateSpecVersion, &armtemplatespecs.TemplateSpecVersionsClientUpdateOptions{
		TemplateSpecVersionUpdateModel: &armtemplatespecs.TemplateSpecVersionUpdateModel{
			Tags: map[string]*string{
				"test": to.StringPtr("recording"),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "recording", *updateResp.Tags["test"])
}

func TestTemplateSpecVersionsClient_Delete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create template spec
	templateSpecName, err := createRandomName(t, "ts")
	require.NoError(t, err)
	templateSpecsClient := armtemplatespecs.NewClient(subscriptionID, cred, opt)
	resp, err := templateSpecsClient.CreateOrUpdate(
		ctx,
		rgName,
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

	// create template version
	templateSpecVersion, err := createRandomName(t, "version")
	require.NoError(t, err)
	templateSpecVersionsClient := armtemplatespecs.NewTemplateSpecVersionsClient(subscriptionID, cred, opt)
	vResp, err := templateSpecVersionsClient.CreateOrUpdate(
		ctx,
		rgName,
		templateSpecName,
		templateSpecVersion,
		armtemplatespecs.TemplateSpecVersion{
			Location: to.StringPtr(location),
			Properties: &armtemplatespecs.TemplateSpecVersionProperties{
				Description: to.StringPtr("<description>"),
				MainTemplate: map[string]interface{}{
					"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
					"contentVersion": "1.0.0.0",
					"parameters":     map[string]interface{}{},
					"resources":      []interface{}{},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, vResp)

	// delete
	delResp, err := templateSpecVersionsClient.Delete(ctx, rgName, templateSpecName, templateSpecVersion, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)
}

func TestTemplateSpecVersionsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create template spec
	templateSpecName, err := createRandomName(t, "ts")
	require.NoError(t, err)
	templateSpecsClient := armtemplatespecs.NewClient(subscriptionID, cred, opt)
	resp, err := templateSpecsClient.CreateOrUpdate(
		ctx,
		rgName,
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

	// create template version
	templateSpecVersion, err := createRandomName(t, "version")
	require.NoError(t, err)
	templateSpecVersionsClient := armtemplatespecs.NewTemplateSpecVersionsClient(subscriptionID, cred, opt)
	vResp, err := templateSpecVersionsClient.CreateOrUpdate(
		ctx,
		rgName,
		templateSpecName,
		templateSpecVersion,
		armtemplatespecs.TemplateSpecVersion{
			Location: to.StringPtr(location),
			Properties: &armtemplatespecs.TemplateSpecVersionProperties{
				Description: to.StringPtr("<description>"),
				MainTemplate: map[string]interface{}{
					"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
					"contentVersion": "1.0.0.0",
					"parameters":     map[string]interface{}{},
					"resources":      []interface{}{},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, vResp)

	// delete
	pager := templateSpecVersionsClient.List(rgName, templateSpecName, nil)
	require.NoError(t, pager.Err())
}

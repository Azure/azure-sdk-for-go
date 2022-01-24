//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDeploymentOperationsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "list", "westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID, cred, opt)
	deploymentName, err := createRandomName(t, "rs")
	temp, err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:     armresources.DeploymentModeIncremental.ToPtr(),
				Template: temp,
				Parameters: map[string]interface{}{
					"location": map[string]string{
						"value": "West US",
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := pollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, deploymentName, *resp.Name)

	// list deployment operations
	deploymentOperationsClient := armresources.NewDeploymentOperationsClient(subscriptionID,cred,opt)
	listPager := deploymentOperationsClient.List(rgName,deploymentName,nil)
	require.True(t,listPager.NextPage(ctx))
}

func TestDeploymentOperationsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "get", "westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID, cred, opt)
	deploymentName, err := createRandomName(t, "rs")
	temp, err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:     armresources.DeploymentModeIncremental.ToPtr(),
				Template: temp,
				Parameters: map[string]interface{}{
					"location": map[string]string{
						"value": "West US",
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := pollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, deploymentName, *resp.Name)

	// list deployment operations
	deploymentOperationsClient := armresources.NewDeploymentOperationsClient(subscriptionID,cred,opt)
	listPager := deploymentOperationsClient.List(rgName,deploymentName,nil)
	require.True(t,listPager.NextPage(ctx))
	operationID := *listPager.PageResponse().Value[0].ID

	// get deployment operation
	getResp,err := deploymentOperationsClient.Get(ctx,rgName,deploymentName,operationID,nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.OperationID,operationID)
}

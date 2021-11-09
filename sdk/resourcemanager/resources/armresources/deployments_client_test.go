//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDeploymentsClient_CheckExistence(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"check","westus")
	defer clean()
	rgName := *rg.Name

	// check existence deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	resp,err := deploymentsClient.CheckExistence(ctx,rgName,deploymentName,nil)
	require.NoError(t, err)
	require.False(t, resp.Success)
}
var template = `
{
 "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
 "contentVersion": "1.0.0.0",
 "parameters": {
   "location": {
     "type": "string",
     "allowedValues": [
       "East US",
       "West US",
       "West Europe",
       "East Asia",
       "South East Asia"
     ],
     "metadata": {
       "description": "Location to deploy to"
     }
   }
 },
 "resources": [
   {
     "type": "Microsoft.Compute/availabilitySets",
     "name": "availabilitySet1",
     "apiVersion": "2019-07-01",
     "location": "[parameters('location')]",
     "properties": {}
   }
 ],
 "outputs": {
   "myparameter": {
     "type": "object",
     "value": "[reference('Microsoft.Compute/availabilitySets/availabilitySet1')]"
   }
 }
}
`

func unmarshalTemplate(data string) (map[string]interface{},error) {
	result := make(map[string]interface{})

	err := json.Unmarshal([]byte(data),&result)
	if err != nil {
		return nil,fmt.Errorf("unmarshal template error:%v",err)
	}
	return result,nil
}

func TestDeploymentsClient_BeginCreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	temp,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)
}

func TestDeploymentsClient_ListByResourceGroup(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"list","westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	template,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,deploymentName,*resp.Name)

	listPager := deploymentsClient.ListByResourceGroup(rgName,nil)
	require.True(t, listPager.NextPage(ctx))
}

func TestDeploymentsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"get","westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	template,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,deploymentName,*resp.Name)

	// get
	getResp,err := deploymentsClient.Get(ctx,rgName,deploymentName,nil)
	require.NoError(t, err)
	require.Equal(t, deploymentName,*getResp.Name)
}

func TestDeploymentsClient_BeginWhatIf(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"what","westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	template,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,deploymentName,*resp.Name)

	// what if
	whatPoller,err := deploymentsClient.BeginWhatIf(
		ctx,
		rgName,
		deploymentName,
		armresources.DeploymentWhatIf{
			Properties: &armresources.DeploymentWhatIfProperties{
				DeploymentProperties: armresources.DeploymentProperties{
					Mode: armresources.DeploymentModeIncremental.ToPtr(),
					Template: template,
				},
			},
		},
		nil,
		)
	require.NoError(t, err)
	whatResp,err := whatPoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, deploymentName,*whatResp.Status)
}

func TestDeploymentsClient_BeginValidate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"validate","westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	template,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,deploymentName,*resp.Name)

	vPoller,err := deploymentsClient.BeginValidate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	vResp,err := vPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, "Incremental",string(*vResp.Properties.Mode))
}

func TestDeploymentsClient_ExportTemplate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"export","westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	template,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,deploymentName,*resp.Name)

	// export template
	exportTemplate,err := deploymentsClient.ExportTemplate(ctx,rgName,deploymentName,nil)
	require.NoError(t, err)
	require.NotNil(t,exportTemplate)
}

func TestDeploymentsClient_BeginDelete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "delete", "westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID, cred, opt)
	deploymentName, err := createRandomName(t, "rs")
	template, err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:     armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	require.Equal(t, deploymentName, *resp.Name)

	// delete deployment
	delPoller,err := deploymentsClient.BeginDelete(ctx,rgName,deploymentName,nil)
	require.NoError(t, err)
	delResp,err := delPoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, 200,delResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_CheckExistenceAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "delete", "westus")
	defer clean()
	scope := *rg.ID

	// check deployment existence
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"dep")
	require.NoError(t, err)
	resp,err := deploymentsClient.CheckExistenceAtScope(ctx,scope,deploymentName,nil)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestDeploymentsClient_BeginCreateOrUpdateAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
	defer clean()
	scope := *rg.ID

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	temp,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtScope(
		ctx,
		scope,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)
}

func TestDeploymentsClient_ListAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"list","westus")
	defer clean()
	scope := *rg.ID

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	temp,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtScope(
		ctx,
		scope,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// list deployment at scope
	listPager := deploymentsClient.ListAtScope(scope,nil)
	require.True(t, listPager.NextPage(ctx))
}

func TestDeploymentsClient_GetAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"get","westus")
	defer clean()
	scope := *rg.ID

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	temp,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtScope(
		ctx,
		scope,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// get deployment at scope
	getResp,err := deploymentsClient.GetAtScope(ctx,scope,deploymentName,nil)
	require.NoError(t, err)
	require.Equal(t, deploymentName,*getResp.Name)
}

func TestDeploymentsClient_BeginValidateAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"validate","westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	template,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,deploymentName,*resp.Name)

	vPoller,err := deploymentsClient.BeginValidateAtScope(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	vResp,err := vPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, "Incremental",string(*vResp.Properties.Mode))
}

func TestDeploymentsClient_ExportTemplateAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"export","westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	template,err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,deploymentName,*resp.Name)

	// export template
	exportTemplate,err := deploymentsClient.ExportTemplateAtScope(ctx,rgName,deploymentName,nil)
	require.NoError(t, err)
	require.NotNil(t,exportTemplate)
}

func TestDeploymentsClient_BeginDeleteAtScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "delete", "westus")
	defer clean()
	rgName := *rg.Name

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID, cred, opt)
	deploymentName, err := createRandomName(t, "rs")
	template, err := unmarshalTemplate(template)
	require.NoError(t, err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:     armresources.DeploymentModeIncremental.ToPtr(),
				Template: template,
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
	require.Equal(t, deploymentName, *resp.Name)

	// delete deployment
	delPoller,err := deploymentsClient.BeginDeleteAtScope(ctx,rgName,deploymentName,nil)
	require.NoError(t, err)
	delResp,err := delPoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, 200,delResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_CheckExistenceAtManagementGroupScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// management group
	groupName := "20000000-0001-0000-0000-000000000123456"

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	resp,err := deploymentsClient.CheckExistenceAtManagementGroupScope(ctx,groupName,deploymentName,nil)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestDeploymentsClient_BeginCreateOrUpdateAtManagementGroupScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create management group
	managementGroupsClient := armmanagementgroups.NewManagementGroupsClient(cred,opt)
	groupName := "20000000-0001-0000-0000-000000000123456"
	mgPoller,err := managementGroupsClient.BeginCreateOrUpdate(
		ctx,
		groupName,
		armmanagementgroups.CreateManagementGroupRequest{
			Name: to.StringPtr(groupName),
		},
		nil,
		)
	require.NoError(t, err)
	mgResp,err := mgPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, groupName,*mgResp.Name)
	defer cleanupManagement(t,managementGroupsClient,groupName)

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(
		ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)
}

func TestDeploymentsClient_ListAtManagementGroupScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create management group
	managementGroupsClient := armmanagementgroups.NewManagementGroupsClient(cred,opt)
	groupName := "20000000-0001-0000-0000-000000000123456"
	mgPoller,err := managementGroupsClient.BeginCreateOrUpdate(
		ctx,
		groupName,
		armmanagementgroups.CreateManagementGroupRequest{
			Name: to.StringPtr(groupName),
		},
		nil,
	)
	require.NoError(t, err)
	mgResp,err := mgPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, groupName,*mgResp.Name)
	defer cleanupManagement(t,managementGroupsClient,groupName)

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(
		ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// list deployment
	listPager := deploymentsClient.ListAtManagementGroupScope(groupName,nil)
	require.True(t, listPager.NextPage(ctx))
}

func TestDeploymentsClient_GetAtManagementGroupScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create management group
	managementGroupsClient := armmanagementgroups.NewManagementGroupsClient(cred,opt)
	groupName := "20000000-0001-0000-0000-000000000123456"
	mgPoller,err := managementGroupsClient.BeginCreateOrUpdate(
		ctx,
		groupName,
		armmanagementgroups.CreateManagementGroupRequest{
			Name: to.StringPtr(groupName),
		},
		nil,
	)
	require.NoError(t, err)
	mgResp,err := mgPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, groupName,*mgResp.Name)
	defer cleanupManagement(t,managementGroupsClient,groupName)

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(
		ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// get deployment
	getResp,err := deploymentsClient.GetAtManagementGroupScope(ctx,groupName,deploymentName,nil)
	require.NoError(t, err)
	require.Equal(t, deploymentName,*getResp.Name)
}

func TestDeploymentsClient_CancelAtManagementGroupScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create management group
	managementGroupsClient := armmanagementgroups.NewManagementGroupsClient(cred,opt)
	groupName := "20000000-0001-0000-0000-000000000123456"
	mgPoller,err := managementGroupsClient.BeginCreateOrUpdate(
		ctx,
		groupName,
		armmanagementgroups.CreateManagementGroupRequest{
			Name: to.StringPtr(groupName),
		},
		nil,
	)
	require.NoError(t, err)
	mgResp,err := mgPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, groupName,*mgResp.Name)
	defer cleanupManagement(t,managementGroupsClient,groupName)

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(
		ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)


	// cancel deployment
	cancelResp,err := deploymentsClient.CancelAtManagementGroupScope(ctx,groupName,deploymentName,nil)
	require.NoError(t, err)
	require.Equal(t, 200,cancelResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_BeginValidateAtManagementGroupScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create management group
	managementGroupsClient := armmanagementgroups.NewManagementGroupsClient(cred,opt)
	groupName := "20000000-0001-0000-0000-000000000123456"
	mgPoller,err := managementGroupsClient.BeginCreateOrUpdate(
		ctx,
		groupName,
		armmanagementgroups.CreateManagementGroupRequest{
			Name: to.StringPtr(groupName),
		},
		nil,
	)
	require.NoError(t, err)
	mgResp,err := mgPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, groupName,*mgResp.Name)

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(
		ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)
	defer cleanupManagement(t,managementGroupsClient,groupName)

	// validate deployment
	validatePoller,err := deploymentsClient.BeginValidateAtManagementGroupScope(
		ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
		)
	require.NoError(t, err)
	validateResp,err := validatePoller.PollUntilDone(ctx,10*time.Second)
	require.NotNil(t,validateResp.Properties)
}

func TestDeploymentsClient_ExportTemplateAtManagementGroupScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create management group
	managementGroupsClient := armmanagementgroups.NewManagementGroupsClient(cred,opt)
	groupName := "20000000-0001-0000-0000-000000000123456"
	mgPoller,err := managementGroupsClient.BeginCreateOrUpdate(
		ctx,
		groupName,
		armmanagementgroups.CreateManagementGroupRequest{
			Name: to.StringPtr(groupName),
		},
		nil,
	)
	require.NoError(t, err)
	mgResp,err := mgPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, groupName,*mgResp.Name)
	defer cleanupManagement(t,managementGroupsClient,groupName)

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(
		ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// export template deployment
	exportResp,err := deploymentsClient.ExportTemplateAtManagementGroupScope(ctx, groupName, deploymentName, nil)
	require.NoError(t, err)
	require.NotNil(t, exportResp.Template)
}

func TestDeploymentsClient_BeginDeleteAtManagementGroupScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create management group
	managementGroupsClient := armmanagementgroups.NewManagementGroupsClient(cred,opt)
	groupName := "20000000-0001-0000-0000-000000000123456"
	mgPoller,err := managementGroupsClient.BeginCreateOrUpdate(
		ctx,
		groupName,
		armmanagementgroups.CreateManagementGroupRequest{
			Name: to.StringPtr(groupName),
		},
		nil,
	)
	require.NoError(t, err)
	mgResp,err := mgPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, groupName,*mgResp.Name)
	defer cleanupManagement(t,managementGroupsClient,groupName)

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(
		ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// delete template deployment
	delPoller,err := deploymentsClient.BeginDeleteAtManagementGroupScope(ctx, groupName, deploymentName, nil)
	require.NoError(t, err)
	delResp,err := delPoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, 200,delResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_CheckExistenceAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// check deployment existence
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"dep")
	require.NoError(t, err)
	resp,err := deploymentsClient.CheckExistenceAtTenantScope(ctx,deploymentName,nil)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestDeploymentsClient_BeginCreateOrUpdateAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtTenantScope(
		ctx,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)
}

func TestDeploymentsClient_ListAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// list deployment
	listResp := deploymentsClient.ListAtTenantScope(nil)
	require.True(t, listResp.NextPage(ctx))
}

func TestDeploymentsClient_GetAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// get deployment
	getResp,err := deploymentsClient.GetAtTenantScope(ctx,deploymentName,nil)
	require.NoError(t, err)
	require.Equal(t, deploymentName,*getResp.Name)
}

func TestDeploymentsClient_BeginWhatIfAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// what if deployment
	whatIfPoller,err := deploymentsClient.BeginWhatIfAtTenantScope(
		ctx,
		deploymentName,
		armresources.ScopedDeploymentWhatIf{
			Properties: &armresources.DeploymentWhatIfProperties{
				DeploymentProperties: armresources.DeploymentProperties{
					Mode: armresources.DeploymentModeIncremental.ToPtr(),
					Template: map[string]interface{}{
						"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
					},
					Parameters: map[string]interface{}{
						"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
					},
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
		)
	require.NoError(t, err)
	whatIfResp,err := whatIfPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, 200,whatIfResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_CancelAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// cancel deployment
	cancelResp,err := deploymentsClient.CancelAtTenantScope(ctx, deploymentName, nil)
	require.NoError(t, err)
	require.Equal(t, 200,cancelResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_BeginValidateAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// validate deployment
	validatePoller,err := deploymentsClient.BeginValidateAtTenantScope(
		ctx,
		deploymentName,
		armresources.ScopedDeployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
		)
	require.NoError(t, err)
	validateResp,err := validatePoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, 200,validateResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_ExportTemplateAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID, cred, opt)
	deploymentName, err := createRandomName(t, "rs")
	require.NoError(t, err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := pollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, deploymentName, *resp.Name)

	// export template deployment
	exportResp, err := deploymentsClient.ExportTemplateAtTenantScope(ctx, deploymentName, nil)
	require.NoError(t, err)
	require.NotNil(t, exportResp.Template)
}

func TestDeploymentsClient_BeginDeleteAtTenantScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID, cred, opt)
	deploymentName, err := createRandomName(t, "rs")
	require.NoError(t, err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := pollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, deploymentName, *resp.Name)

	// export template deployment
	delPoller, err := deploymentsClient.BeginDeleteAtTenantScope(ctx, deploymentName, nil)
	require.NoError(t, err)
	delResp,err := delPoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,200,delResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_CheckExistenceAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// check deployment existence
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"dep")
	require.NoError(t, err)
	resp,err := deploymentsClient.CheckExistenceAtSubscriptionScope(ctx,deploymentName,nil)
	require.NoError(t, err)
	require.False(t, resp.Success)
}

func TestDeploymentsClient_BeginCreateOrUpdateAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)
}

func TestDeploymentsClient_ListAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// list deployment
	listResp := deploymentsClient.ListAtSubscriptionScope(nil)
	require.True(t, listResp.NextPage(ctx))
}

func TestDeploymentsClient_GetAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// get deployment
	getResp,err := deploymentsClient.GetAtSubscriptionScope(ctx,deploymentName,nil)
	require.NoError(t, err)
	require.Equal(t, deploymentName,*getResp.Name)
}

func TestDeploymentsClient_BeginWhatIfAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// what if deployment
	whatIfPoller,err := deploymentsClient.BeginWhatIfAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.DeploymentWhatIf{
			Properties: &armresources.DeploymentWhatIfProperties{
				DeploymentProperties: armresources.DeploymentProperties{
					Mode: armresources.DeploymentModeIncremental.ToPtr(),
					Template: map[string]interface{}{
						"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
					},
					Parameters: map[string]interface{}{
						"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
					},
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	whatIfResp,err := whatIfPoller.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, 200,whatIfResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_CancelAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// cancel deployment
	cancelResp,err := deploymentsClient.CancelAtSubscriptionScope(ctx, deploymentName, nil)
	require.NoError(t, err)
	require.Equal(t, 200,cancelResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_BeginValidateAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID,cred,opt)
	deploymentName,err := createRandomName(t,"rs")
	require.NoError(t, err)
	pollerResp,err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t,deploymentName,*resp.Name)

	// validate deployment
	validatePoller,err := deploymentsClient.BeginValidateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	validateResp,err := validatePoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t, 200,validateResp.RawResponse.StatusCode)
}

func TestDeploymentsClient_ExportTemplateAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID, cred, opt)
	deploymentName, err := createRandomName(t, "rs")
	require.NoError(t, err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := pollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, deploymentName, *resp.Name)

	// export template deployment
	exportResp, err := deploymentsClient.ExportTemplateAtSubscriptionScope(ctx, deploymentName, nil)
	require.NoError(t, err)
	require.NotNil(t, exportResp.Template)
}

func TestDeploymentsClient_BeginDeleteAtSubscriptionScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create deployment
	deploymentsClient := armresources.NewDeploymentsClient(subscriptionID, cred, opt)
	deploymentName, err := createRandomName(t, "rs")
	require.NoError(t, err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode: armresources.DeploymentModeIncremental.ToPtr(),
				Template: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
				Parameters: map[string]interface{}{
					"uri": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json",
				},
			},
			Location: to.StringPtr("West US"),
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := pollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, deploymentName, *resp.Name)

	// delete deployment
	delPoller, err := deploymentsClient.BeginDeleteAtSubscriptionScope(ctx, deploymentName, nil)
	require.NoError(t, err)
	delResp,err := delPoller.PollUntilDone(ctx,10*time.Second)
	require.Equal(t,200,delResp.RawResponse.StatusCode)
}
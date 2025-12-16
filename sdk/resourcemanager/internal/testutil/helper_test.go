// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
)

func TestCreateDeleteResourceGroup(t *testing.T) {
	ctx := context.Background()
	cred, options := GetCredAndClientOptions(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	stop := StartRecording(t, pathToPackage)
	defer stop()
	resourceGroup, _, err := CreateResourceGroup(ctx, subscriptionID, cred, options, "eastus")
	require.NoError(t, err)
	_, err = DeleteResourceGroup(ctx, subscriptionID, cred, options, *resourceGroup.Name)
	require.NoError(t, err)
}

func TestCreateDeployment(t *testing.T) {
	ctx := context.Background()
	cred, options := GetCredAndClientOptions(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	stop := StartRecording(t, pathToPackage)
	defer stop()
	resourceGroup, _, err := CreateResourceGroup(ctx, subscriptionID, cred, options, "eastus")
	require.NoError(t, err)
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"resourceName": map[string]interface{}{
				"type":  "string",
				"value": "[variables('name').value]",
			},
		},
		"resources": []interface{}{},
		"variables": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
				"metadata": map[string]interface{}{
					"description": "Name of the SignalR service.",
				},
				"value": "[concat('sw',uniqueString(resourceGroup().id))]",
			},
		},
	}
	params := map[string]interface{}{}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := CreateDeployment(ctx, subscriptionID, cred, options, *resourceGroup.Name, "Generate_Unique_Name", &deployment)
	require.NoError(t, err)
	require.NotEmpty(t, deploymentExtend.Properties.Outputs.(map[string]interface{})["resourceName"].(map[string]interface{})["value"].(string))
	_, err = DeleteResourceGroup(ctx, subscriptionID, cred, options, *resourceGroup.Name)
	require.NoError(t, err)
}

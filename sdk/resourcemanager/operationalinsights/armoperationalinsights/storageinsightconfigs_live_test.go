// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armoperationalinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type StorageInsightConfigsTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	storageAccountName string
	workspaceName      string
	location           string
	resourceGroupName  string
	storageAccountId   string
	subscriptionId     string
}

func (testsuite *StorageInsightConfigsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.storageAccountName = "oistorageinsightconfigx"
	testsuite.workspaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "oistorageinsightconfig", 28, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.storageAccountId = recording.GetEnvVariable("STORAGE_ACCOUNT_ID", "")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *StorageInsightConfigsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestStorageInsightConfigsTestSuite(t *testing.T) {
	suite.Run(t, new(StorageInsightConfigsTestSuite))
}

func (testsuite *StorageInsightConfigsTestSuite) Prepare() {
	var err error
	// From step StorageAccount_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"storageAccountId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName'))]",
			},
		},
		"parameters": map[string]interface{}{
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"storageAccountName": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.storageAccountName,
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('storageAccountName')]",
				"type":       "Microsoft.Storage/storageAccounts",
				"apiVersion": "2022-05-01",
				"kind":       "StorageV2",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"accessTier":                   "Hot",
					"allowBlobPublicAccess":        true,
					"allowCrossTenantReplication":  true,
					"allowSharedKeyAccess":         true,
					"defaultToOAuthAuthentication": false,
					"dnsEndpointType":              "Standard",
					"encryption": map[string]interface{}{
						"keySource":                       "Microsoft.Storage",
						"requireInfrastructureEncryption": false,
						"services": map[string]interface{}{
							"blob": map[string]interface{}{
								"enabled": true,
								"keyType": "Account",
							},
							"file": map[string]interface{}{
								"enabled": true,
								"keyType": "Account",
							},
						},
					},
					"minimumTlsVersion": "TLS1_2",
					"networkAcls": map[string]interface{}{
						"bypass":              "AzureServices",
						"defaultAction":       "Allow",
						"ipRules":             []interface{}{},
						"virtualNetworkRules": []interface{}{},
					},
					"publicNetworkAccess":      "Enabled",
					"supportsHttpsTrafficOnly": true,
				},
				"sku": map[string]interface{}{
					"name": "Standard_RAGRS",
					"tier": "Standard",
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "StorageAccount_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.storageAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountId"].(map[string]interface{})["value"].(string)

	// From step Workspaces_Create
	fmt.Println("Call operation: Workspaces_CreateOrUpdate")
	workspacesClient, err := armoperationalinsights.NewWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	workspacesClientCreateOrUpdateResponsePoller, err := workspacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, armoperationalinsights.Workspace{
		Location: to.Ptr(testsuite.location),
		Properties: &armoperationalinsights.WorkspaceProperties{
			RetentionInDays: to.Ptr[int32](30),
			SKU: &armoperationalinsights.WorkspaceSKU{
				Name: to.Ptr(armoperationalinsights.WorkspaceSKUNameEnumPerGB2018),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, workspacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.OperationalInsights/workspaces/storageInsightConfigs
func (testsuite *StorageInsightConfigsTestSuite) TestStorageInsightConfig() {
	var err error
	// From step StorageInsightConfigs_CreateOrUpdate
	fmt.Println("Call operation: StorageInsightConfigs_CreateOrUpdate")
	storageInsightConfigsClient, err := armoperationalinsights.NewStorageInsightConfigsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = storageInsightConfigsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "AzTestSI1110", armoperationalinsights.StorageInsight{
		Properties: &armoperationalinsights.StorageInsightProperties{
			Containers: []*string{
				to.Ptr("wad-iis-logfiles")},
			StorageAccount: &armoperationalinsights.StorageAccount{
				ID:  to.Ptr(testsuite.storageAccountId),
				Key: to.Ptr("1234"),
			},
			Tables: []*string{
				to.Ptr("WADWindowsEventLogsTable"),
				to.Ptr("LinuxSyslogVer2v0")},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step StorageInsightConfigs_ListByWorkspace
	fmt.Println("Call operation: StorageInsightConfigs_ListByWorkspace")
	storageInsightConfigsClientNewListByWorkspacePager := storageInsightConfigsClient.NewListByWorkspacePager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for storageInsightConfigsClientNewListByWorkspacePager.More() {
		_, err := storageInsightConfigsClientNewListByWorkspacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StorageInsightConfigs_Get
	fmt.Println("Call operation: StorageInsightConfigs_Get")
	_, err = storageInsightConfigsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "AzTestSI1110", nil)
	testsuite.Require().NoError(err)

	// From step StorageInsightConfigs_Delete
	fmt.Println("Call operation: StorageInsightConfigs_Delete")
	_, err = storageInsightConfigsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "AzTestSI1110", nil)
	testsuite.Require().NoError(err)
}

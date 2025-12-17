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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type DataExportsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	namespacesName    string
	workspaceName     string
	eventhubId        string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *DataExportsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.namespacesName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "eventhubnamespace", 23, false)
	testsuite.workspaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "oidataexport", 18, false)
	testsuite.eventhubId = recording.GetEnvVariable("EVENTHUB_ID", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DataExportsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDataExportsTestSuite(t *testing.T) {
	suite.Run(t, new(DataExportsTestSuite))
}

func (testsuite *DataExportsTestSuite) Prepare() {
	var err error
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

	// From step EventhubNamespace_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"eventhubId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.EventHub/namespaces', parameters('namespacesName'))]",
			},
		},
		"parameters": map[string]interface{}{
			"namespacesName": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.namespacesName,
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('namespacesName')]",
				"type":       "Microsoft.EventHub/namespaces",
				"apiVersion": "2022-01-01-preview",
				"location":   "East US",
				"properties": map[string]interface{}{
					"disableLocalAuth":       false,
					"isAutoInflateEnabled":   false,
					"kafkaEnabled":           true,
					"maximumThroughputUnits": float64(0),
					"minimumTlsVersion":      "1.2",
					"publicNetworkAccess":    "Enabled",
					"zoneRedundant":          true,
				},
				"sku": map[string]interface{}{
					"name":     "Standard",
					"capacity": float64(1),
					"tier":     "Standard",
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "EventhubNamespace_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.eventhubId = deploymentExtend.Properties.Outputs.(map[string]interface{})["eventhubId"].(map[string]interface{})["value"].(string)
}

// Microsoft.OperationalInsights/workspaces/dataExports
func (testsuite *DataExportsTestSuite) TestDataExport() {
	var err error
	// From step DataExports_CreateOrUpdate
	fmt.Println("Call operation: DataExports_CreateOrUpdate")
	dataExportsClient, err := armoperationalinsights.NewDataExportsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = dataExportsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "export1", armoperationalinsights.DataExport{
		Properties: &armoperationalinsights.DataExportProperties{
			Destination: &armoperationalinsights.Destination{
				ResourceID: to.Ptr(testsuite.eventhubId),
			},
			TableNames: []*string{
				to.Ptr("Heartbeat")},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step DataExports_ListByWorkspace
	fmt.Println("Call operation: DataExports_ListByWorkspace")
	dataExportsClientNewListByWorkspacePager := dataExportsClient.NewListByWorkspacePager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for dataExportsClientNewListByWorkspacePager.More() {
		_, err := dataExportsClientNewListByWorkspacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataExports_Get
	fmt.Println("Call operation: DataExports_Get")
	_, err = dataExportsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "export1", nil)
	testsuite.Require().NoError(err)

	// From step DataExports_Delete
	fmt.Println("Call operation: DataExports_Delete")
	_, err = dataExportsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "export1", nil)
	testsuite.Require().NoError(err)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdatabricks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databricks/armdatabricks"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type VnetpeeringTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	peeringName       string
	virtaulNetworkId  string
	workspaceName     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *VnetpeeringTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.peeringName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "peeringn", 14, false)
	testsuite.workspaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "workspac", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *VnetpeeringTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestVnetpeeringTestSuite(t *testing.T) {
	suite.Run(t, new(VnetpeeringTestSuite))
}

func (testsuite *VnetpeeringTestSuite) Prepare() {
	var err error
	// From step Workspaces_CreateOrUpdate
	fmt.Println("Call operation: Workspaces_CreateOrUpdate")
	workspacesClient, err := armdatabricks.NewWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	workspacesClientCreateOrUpdateResponsePoller, err := workspacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, armdatabricks.Workspace{
		Location: to.Ptr(testsuite.location),
		Properties: &armdatabricks.WorkspaceProperties{
			ManagedResourceGroupID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/myManaged" + testsuite.resourceGroupName),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, workspacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Create_VirtualNetwork
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"virtaulNetworkId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "databricksvnet",
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"subnets": []any{
						map[string]any{
							"name": "default",
							"properties": map[string]any{
								"addressPrefix": "10.0.0.0/24",
							},
						},
					},
				},
				"tags": map[string]any{},
			},
		},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_VirtualNetwork", &deployment)
	testsuite.Require().NoError(err)
	testsuite.virtaulNetworkId = deploymentExtend.Properties.Outputs.(map[string]interface{})["virtaulNetworkId"].(map[string]interface{})["value"].(string)
}

// Microsoft.Databricks/workspaces/{workspaceName}/virtualNetworkPeerings/{peeringName}
func (testsuite *VnetpeeringTestSuite) TestVNetPeering() {
	var err error
	// From step vNetPeering_CreateOrUpdate
	fmt.Println("Call operation: vNetPeering_CreateOrUpdate")
	vNetPeeringClient, err := armdatabricks.NewVNetPeeringClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vNetPeeringClientCreateOrUpdateResponsePoller, err := vNetPeeringClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, testsuite.peeringName, armdatabricks.VirtualNetworkPeering{
		Properties: &armdatabricks.VirtualNetworkPeeringPropertiesFormat{
			AllowForwardedTraffic:     to.Ptr(false),
			AllowGatewayTransit:       to.Ptr(false),
			AllowVirtualNetworkAccess: to.Ptr(true),
			RemoteVirtualNetwork: &armdatabricks.VirtualNetworkPeeringPropertiesFormatRemoteVirtualNetwork{
				ID: to.Ptr(testsuite.virtaulNetworkId),
			},
			UseRemoteGateways: to.Ptr(false),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, vNetPeeringClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step vNetPeering_ListByWorkspace
	fmt.Println("Call operation: vNetPeering_ListByWorkspace")
	vNetPeeringClientNewListByWorkspacePager := vNetPeeringClient.NewListByWorkspacePager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for vNetPeeringClientNewListByWorkspacePager.More() {
		_, err := vNetPeeringClientNewListByWorkspacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step vNetPeering_Get
	fmt.Println("Call operation: vNetPeering_Get")
	_, err = vNetPeeringClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, testsuite.peeringName, nil)
	testsuite.Require().NoError(err)

	// From step vNetPeering_Delete
	fmt.Println("Call operation: vNetPeering_Delete")
	vNetPeeringClientDeleteResponsePoller, err := vNetPeeringClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, testsuite.peeringName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, vNetPeeringClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

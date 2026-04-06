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

type DatabricksTestSuite struct {
	suite.Suite

	ctx                      context.Context
	cred                     azcore.TokenCredential
	options                  *arm.ClientOptions
	managedResourceGroupName string
	peeringName              string
	virtaulNetworkId         string
	workspaceId              string
	workspaceName            string
	location                 string
	resourceGroupName        string
	subscriptionId           string
}

func (testsuite *DatabricksTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.managedResourceGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "managedrg", 15, false)
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

func (testsuite *DatabricksTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDatabricksTestSuite(t *testing.T) {
	suite.Run(t, new(DatabricksTestSuite))
}

func (testsuite *DatabricksTestSuite) Prepare() {
	var err error
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
			"networkSecurityGroupName": map[string]any{
				"type":         "string",
				"defaultValue": "databricks-nsg",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "databricksvnet2",
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('networkSecurityGroupName')]",
				"type":       "Microsoft.Network/networkSecurityGroups",
				"apiVersion": "2022-11-01",
				"location":   "[parameters('location')]",
			},
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/networkSecurityGroups', parameters('networkSecurityGroupName'))]",
				},
				"location": "[parameters('location')]",
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
						map[string]any{
							"name": "myPublicSubnet",
							"properties": map[string]any{
								"addressPrefix": "10.0.1.0/24",
								"delegations": []any{
									map[string]any{
										"name": "Microsoft.Databricks.workspaces",
										"type": "Microsoft.Network/virtualNetworks/subnets/delegations",
										"properties": map[string]any{
											"serviceName": "Microsoft.Databricks/workspaces",
										},
									},
								},
								"networkSecurityGroup": map[string]any{
									"id": "[resourceId('Microsoft.Network/networkSecurityGroups', parameters('networkSecurityGroupName'))]",
								},
							},
						},
						map[string]any{
							"name": "myPrivateSubnet",
							"properties": map[string]any{
								"addressPrefix": "10.0.2.0/24",
								"delegations": []any{
									map[string]any{
										"name": "Microsoft.Databricks.workspaces",
										"type": "Microsoft.Network/virtualNetworks/subnets/delegations",
										"properties": map[string]any{
											"serviceName": "Microsoft.Databricks/workspaces",
										},
									},
								},
								"networkSecurityGroup": map[string]any{
									"id": "[resourceId('Microsoft.Network/networkSecurityGroups', parameters('networkSecurityGroupName'))]",
								},
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

	// From step Workspaces_CreateOrUpdate
	fmt.Println("Call operation: Workspaces_CreateOrUpdate")
	workspacesClient, err := armdatabricks.NewWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	workspacesClientCreateOrUpdateResponsePoller, err := workspacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, armdatabricks.Workspace{
		Location: to.Ptr(testsuite.location),
		Properties: &armdatabricks.WorkspaceProperties{
			ManagedResourceGroupID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.managedResourceGroupName),
			Parameters: &armdatabricks.WorkspaceCustomParameters{
				CustomPrivateSubnetName: &armdatabricks.WorkspaceCustomStringParameter{
					Value: to.Ptr("myPrivateSubnet"),
				},
				CustomPublicSubnetName: &armdatabricks.WorkspaceCustomStringParameter{
					Value: to.Ptr("myPublicSubnet"),
				},
				CustomVirtualNetworkID: &armdatabricks.WorkspaceCustomStringParameter{
					Value: to.Ptr(testsuite.virtaulNetworkId),
				},
			},
		},
		SKU: &armdatabricks.SKU{
			Name: to.Ptr("premium"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var workspacesClientCreateOrUpdateResponse *armdatabricks.WorkspacesClientCreateOrUpdateResponse
	workspacesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, workspacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.workspaceId = *workspacesClientCreateOrUpdateResponse.ID

	// From step Create_PrivateEndpoint
	template = map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "endpointdatabricks-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "endpointdatabricks",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "endpointvnet",
			},
			"workspaceId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.workspaceId,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2020-11-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"enableDdosProtection": false,
					"subnets": []any{
						map[string]any{
							"name": "pesubnet",
							"properties": map[string]any{
								"addressPrefix":                     "10.0.3.0/24",
								"delegations":                       []any{},
								"privateEndpointNetworkPolicies":    "Disabled",
								"privateLinkServiceNetworkPolicies": "Enabled",
							},
						},
					},
					"virtualNetworkPeerings": []any{},
				},
			},
			map[string]any{
				"name":       "[parameters('networkInterfaceName')]",
				"type":       "Microsoft.Network/networkInterfaces",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'pesubnet')]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"dnsSettings": map[string]any{
						"dnsServers": []any{},
					},
					"enableIPForwarding": false,
					"ipConfigurations": []any{
						map[string]any{
							"name": "privateEndpointIpConfig",
							"properties": map[string]any{
								"primary":                   true,
								"privateIPAddress":          "10.0.0.4",
								"privateIPAddressVersion":   "IPv4",
								"privateIPAllocationMethod": "Dynamic",
								"subnet": map[string]any{
									"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'pesubnet')]",
								},
							},
						},
					},
				},
			},
			map[string]any{
				"name":       "[parameters('privateEndpointName')]",
				"type":       "Microsoft.Network/privateEndpoints",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'pesubnet')]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"customDnsConfigs":                    []any{},
					"manualPrivateLinkServiceConnections": []any{},
					"privateLinkServiceConnections": []any{
						map[string]any{
							"name": "[parameters('privateEndpointName')]",
							"properties": map[string]any{
								"groupIds": []any{
									"databricks_ui_api",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('workspaceId')]",
							},
						},
					},
					"subnet": map[string]any{
						"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'pesubnet')]",
					},
				},
			},
			map[string]any{
				"name":       "[concat(parameters('virtualNetworksName'), '/pesubnet')]",
				"type":       "Microsoft.Network/virtualNetworks/subnets",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
				},
				"properties": map[string]any{
					"addressPrefix":                     "10.0.3.0/24",
					"delegations":                       []any{},
					"privateEndpointNetworkPolicies":    "Disabled",
					"privateLinkServiceNetworkPolicies": "Enabled",
				},
			},
		},
		"variables": map[string]any{},
	}
	deployment = armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)
}

// Microsoft.Databricks/workspaces/{workspaceName}
func (testsuite *DatabricksTestSuite) TestWorkspaces() {
	var err error
	// From step Workspaces_ListByResourceGroup
	fmt.Println("Call operation: Workspaces_ListByResourceGroup")
	workspacesClient, err := armdatabricks.NewWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	workspacesClientNewListByResourceGroupPager := workspacesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for workspacesClientNewListByResourceGroupPager.More() {
		_, err := workspacesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Workspaces_Get
	fmt.Println("Call operation: Workspaces_Get")
	_, err = workspacesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, nil)
	testsuite.Require().NoError(err)

	// From step Workspaces_Update
	fmt.Println("Call operation: Workspaces_Update")
	workspacesClientUpdateResponsePoller, err := workspacesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, armdatabricks.WorkspaceUpdate{
		Tags: map[string]*string{
			"mytag1": to.Ptr("myvalue1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, workspacesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Databricks/workspaces/{workspaceName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *DatabricksTestSuite) TestPrivateEndpointConnections() {
	var privateEndpointConnectionName string
	var err error
	// From step PrivateEndpointConnections_List
	fmt.Println("Call operation: PrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armdatabricks.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Create
	fmt.Println("Call operation: PrivateEndpointConnections_Create")
	privateEndpointConnectionsClientCreateResponsePoller, err := privateEndpointConnectionsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, privateEndpointConnectionName, armdatabricks.PrivateEndpointConnection{
		Properties: &armdatabricks.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armdatabricks.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Approved by databricksadmin@contoso.com"),
				Status:      to.Ptr(armdatabricks.PrivateLinkServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_List
	fmt.Println("Call operation: PrivateLinkResources_List")
	privateLinkResourcesClient, err := armdatabricks.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListPager := privateLinkResourcesClient.NewListPager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for privateLinkResourcesClientNewListPager.More() {
		_, err := privateLinkResourcesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateLinkResources_Get
	fmt.Println("Call operation: PrivateLinkResources_Get")
	_, err = privateLinkResourcesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "databricks_ui_api", nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Databricks/workspaces/{workspaceName}/outboundNetworkDependenciesEndpoints
func (testsuite *DatabricksTestSuite) TestOutboundNetworkDependenciesEndpoints() {
	var err error
	// From step OutboundNetworkDependenciesEndpoints_List
	fmt.Println("Call operation: OutboundNetworkDependenciesEndpoints_List")
	outboundNetworkDependenciesEndpointsClient, err := armdatabricks.NewOutboundNetworkDependenciesEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = outboundNetworkDependenciesEndpointsClient.List(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *DatabricksTestSuite) Cleanup() {
	var err error
	// From step Workspaces_Delete
	fmt.Println("Call operation: Workspaces_Delete")
	workspacesClient, err := armdatabricks.NewWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	workspacesClientDeleteResponsePoller, err := workspacesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, workspacesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

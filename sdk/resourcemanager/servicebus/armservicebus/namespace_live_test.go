//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armservicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
	"github.com/stretchr/testify/suite"
)

type NamespaceTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	namespaceId           string
	namespaceName         string
	ruleSubnetId          string
	authorizationRuleName string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *NamespaceTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespaceauthoriz", 23, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NamespaceTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(NamespaceTestSuite))
}

func (testsuite *NamespaceTestSuite) Prepare() {
	var err error
	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClient, err := armservicebus.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientCreateOrUpdateResponsePoller, err := namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.SBNamespace{
		Location: to.Ptr(testsuite.location),
		SKU: &armservicebus.SBSKU{
			Name: to.Ptr(armservicebus.SKUNamePremium),
			Tier: to.Ptr(armservicebus.SKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var namespacesClientCreateOrUpdateResponse *armservicebus.NamespacesClientCreateOrUpdateResponse
	namespacesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.namespaceId = *namespacesClientCreateOrUpdateResponse.ID

	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"ruleSubnetId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'networkruleset')]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"namespaceId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.namespaceId,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "peservicebus-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "peservicebus",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "peservicebusvnet",
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
							"name": "default",
							"properties": map[string]any{
								"addressPrefix":                     "10.0.0.0/24",
								"delegations":                       []any{},
								"privateEndpointNetworkPolicies":    "Enabled",
								"privateLinkServiceNetworkPolicies": "Enabled",
							},
						},
						map[string]any{
							"name": "networkruleset",
							"properties": map[string]any{
								"addressPrefix":                     "10.0.1.0/24",
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
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
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
									"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
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
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
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
									"namespace",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('namespaceId')]",
							},
						},
					},
					"subnet": map[string]any{
						"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
					},
				},
			},
		},
		"variables": map[string]any{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)
	testsuite.ruleSubnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["ruleSubnetId"].(map[string]interface{})["value"].(string)
}

// Microsoft.ServiceBus/namespaces/{namespaceName}
func (testsuite *NamespaceTestSuite) TestNamespaces() {
	var err error
	// From step Namespaces_CheckNameAvailability
	fmt.Println("Call operation: Namespaces_CheckNameAvailability")
	namespacesClient, err := armservicebus.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = namespacesClient.CheckNameAvailability(testsuite.ctx, armservicebus.CheckNameAvailability{
		Name: to.Ptr(testsuite.namespaceName),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_ListByResourceGroup
	fmt.Println("Call operation: Namespaces_ListByResourceGroup")
	namespacesClientNewListByResourceGroupPager := namespacesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for namespacesClientNewListByResourceGroupPager.More() {
		_, err := namespacesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_List
	fmt.Println("Call operation: Namespaces_List")
	namespacesClientNewListPager := namespacesClient.NewListPager(nil)
	for namespacesClientNewListPager.More() {
		_, err := namespacesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_Get
	fmt.Println("Call operation: Namespaces_Get")
	_, err = namespacesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_Update
	fmt.Println("Call operation: Namespaces_Update")
	_, err = namespacesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.SBNamespaceUpdateParameters{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag3": to.Ptr("value3"),
			"tag4": to.Ptr("value4"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/networkRuleSets/default
func (testsuite *NamespaceTestSuite) TestNamespacesNetworkRuleSet() {
	var err error
	// From step Namespaces_CreateOrUpdateNetworkRuleSet
	fmt.Println("Call operation: Namespaces_CreateOrUpdateNetworkRuleSet")
	namespacesClient, err := armservicebus.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = namespacesClient.CreateOrUpdateNetworkRuleSet(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.NetworkRuleSet{
		Properties: &armservicebus.NetworkRuleSetProperties{
			DefaultAction: to.Ptr(armservicebus.DefaultActionDeny),
			IPRules: []*armservicebus.NWRuleSetIPRules{
				{
					Action: to.Ptr(armservicebus.NetworkRuleIPActionAllow),
					IPMask: to.Ptr("1.1.1.1"),
				},
				{
					Action: to.Ptr(armservicebus.NetworkRuleIPActionAllow),
					IPMask: to.Ptr("1.1.1.2"),
				}},
			VirtualNetworkRules: []*armservicebus.NWRuleSetVirtualNetworkRules{
				{
					IgnoreMissingVnetServiceEndpoint: to.Ptr(true),
					Subnet: &armservicebus.Subnet{
						ID: to.Ptr(testsuite.ruleSubnetId),
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_GetNetworkRuleSet
	fmt.Println("Call operation: Namespaces_GetNetworkRuleSet")
	_, err = namespacesClient.GetNetworkRuleSet(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_ListNetworkRuleSets
	fmt.Println("Call operation: Namespaces_ListNetworkRuleSets")
	namespacesClientNewListNetworkRuleSetsPager := namespacesClient.NewListNetworkRuleSetsPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for namespacesClientNewListNetworkRuleSetsPager.More() {
		_, err := namespacesClientNewListNetworkRuleSetsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/AuthorizationRules/{authorizationRuleName}
func (testsuite *NamespaceTestSuite) TestNamespacesAuthorization() {
	var err error
	// From step Namespaces_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: Namespaces_CreateOrUpdateAuthorizationRule")
	namespacesClient, err := armservicebus.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = namespacesClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, armservicebus.SBAuthorizationRule{
		Properties: &armservicebus.SBAuthorizationRuleProperties{
			Rights: []*armservicebus.AccessRights{
				to.Ptr(armservicebus.AccessRightsListen),
				to.Ptr(armservicebus.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_ListAuthorizationRules
	fmt.Println("Call operation: Namespaces_ListAuthorizationRules")
	namespacesClientNewListAuthorizationRulesPager := namespacesClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for namespacesClientNewListAuthorizationRulesPager.More() {
		_, err := namespacesClientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_GetAuthorizationRule
	fmt.Println("Call operation: Namespaces_GetAuthorizationRule")
	_, err = namespacesClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_ListKeys
	fmt.Println("Call operation: Namespaces_ListKeys")
	_, err = namespacesClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_RegenerateKeys
	fmt.Println("Call operation: Namespaces_RegenerateKeys")
	_, err = namespacesClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, armservicebus.RegenerateAccessKeyParameters{
		KeyType: to.Ptr(armservicebus.KeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_DeleteAuthorizationRule
	fmt.Println("Call operation: Namespaces_DeleteAuthorizationRule")
	_, err = namespacesClient.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *NamespaceTestSuite) TestPrivateEndpointConnections() {
	var privateEndpointConnectionName string
	var err error
	// From step PrivateEndpointConnections_List
	fmt.Println("Call operation: PrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armservicebus.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_CreateOrUpdate
	fmt.Println("Call operation: PrivateEndpointConnections_CreateOrUpdate")
	_, err = privateEndpointConnectionsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, privateEndpointConnectionName, armservicebus.PrivateEndpointConnection{
		Properties: &armservicebus.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armservicebus.ConnectionState{
				Description: to.Ptr("testing"),
				Status:      to.Ptr(armservicebus.PrivateLinkConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_Get
	fmt.Println("Call operation: PrivateLinkResources_Get")
	privateLinkResourcesClient, err := armservicebus.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = privateLinkResourcesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *NamespaceTestSuite) Cleanup() {
	var err error
	// From step Namespaces_Delete
	fmt.Println("Call operation: Namespaces_Delete")
	namespacesClient, err := armservicebus.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientDeleteResponsePoller, err := namespacesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

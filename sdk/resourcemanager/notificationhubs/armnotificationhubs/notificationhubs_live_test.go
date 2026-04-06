// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnotificationhubs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/notificationhubs/armnotificationhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type NotificationhubsTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	armEndpoint           string
	authorizationRuleName string
	namespaceName         string
	notificationHubName   string
	notificationhubsId    string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *NotificationhubsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "authoriz", 14, false)
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.notificationHubName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "notifica", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NotificationhubsTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestNotificationhubsTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationhubsTestSuite))
}

func (testsuite *NotificationhubsTestSuite) Prepare() {
	var err error
	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClient, err := armnotificationhubs.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientCreateOrUpdateResponsePoller, err := namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armnotificationhubs.NamespaceResource{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
		Properties: &armnotificationhubs.NamespaceProperties{
			NetworkACLs: &armnotificationhubs.NetworkACLs{
				IPRules: []*armnotificationhubs.IPRule{
					{
						IPMask: to.Ptr("185.48.100.00/24"),
						Rights: []*armnotificationhubs.AccessRights{
							to.Ptr(armnotificationhubs.AccessRightsManage),
							to.Ptr(armnotificationhubs.AccessRightsSend),
							to.Ptr(armnotificationhubs.AccessRightsListen),
						},
					},
				},
				PublicNetworkRule: &armnotificationhubs.PublicInternetAuthorizationRule{
					Rights: []*armnotificationhubs.AccessRights{
						to.Ptr(armnotificationhubs.AccessRightsListen),
					},
				},
			},
			ZoneRedundancy: to.Ptr(armnotificationhubs.ZoneRedundancyPreferenceEnabled),
		},
		SKU: &armnotificationhubs.SKU{
			Name: to.Ptr(armnotificationhubs.SKUNameStandard),
			Tier: to.Ptr("Standard"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var namespacesClientCreateOrUpdateResponse *armnotificationhubs.NamespacesClientCreateOrUpdateResponse
	namespacesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.notificationhubsId = *namespacesClientCreateOrUpdateResponse.ID

	// From step NotificationHubs_CreateOrUpdate
	fmt.Println("Call operation: NotificationHubs_CreateOrUpdate")
	client, err := armnotificationhubs.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, armnotificationhubs.NotificationHubResource{
		Location:   to.Ptr(testsuite.location),
		Properties: &armnotificationhubs.NotificationHubProperties{},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.NotificationHubs/namespaces/{namespaceName}
func (testsuite *NotificationhubsTestSuite) TestNamespaces() {
	var err error
	// From step Namespaces_ListAll
	fmt.Println("Call operation: Namespaces_ListAll")
	namespacesClient, err := armnotificationhubs.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientNewListAllPager := namespacesClient.NewListAllPager(&armnotificationhubs.NamespacesClientListAllOptions{
		SkipToken: nil,
		Top:       nil,
	})
	for namespacesClientNewListAllPager.More() {
		_, err := namespacesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_List
	fmt.Println("Call operation: Namespaces_List")
	namespacesClientNewListPager := namespacesClient.NewListPager(testsuite.resourceGroupName, &armnotificationhubs.NamespacesClientListOptions{
		SkipToken: nil,
		Top:       nil,
	})
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
	_, err = namespacesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armnotificationhubs.NamespacePatchParameters{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_CheckAvailability
	fmt.Println("Call operation: Namespaces_CheckAvailability")
	_, err = namespacesClient.CheckAvailability(testsuite.ctx, armnotificationhubs.CheckAvailabilityParameters{
		Name: to.Ptr("sdk-Namespace-2924"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_GetPnsCredentials
	fmt.Println("Call operation: Namespaces_GetPnsCredentials")
	_, err = namespacesClient.GetPnsCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: Namespaces_CreateOrUpdateAuthorizationRule")
	_, err = namespacesClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, armnotificationhubs.SharedAccessAuthorizationRuleResource{
		Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
			Rights: []*armnotificationhubs.AccessRights{
				to.Ptr(armnotificationhubs.AccessRightsListen),
				to.Ptr(armnotificationhubs.AccessRightsSend),
			},
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

	// From step Namespaces_RegenerateKeys
	fmt.Println("Call operation: Namespaces_RegenerateKeys")
	_, err = namespacesClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, armnotificationhubs.PolicyKeyResource{
		PolicyKey: to.Ptr(armnotificationhubs.PolicyKeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_ListKeys
	fmt.Println("Call operation: Namespaces_ListKeys")
	_, err = namespacesClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_DeleteAuthorizationRule
	fmt.Println("Call operation: Namespaces_DeleteAuthorizationRule")
	_, err = namespacesClient.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.NotificationHubs/namespaces/{namespaceName}/notificationHubs/{notificationHubName}
func (testsuite *NotificationhubsTestSuite) TestNotificationHubs() {
	var err error
	// From step NotificationHubs_List
	fmt.Println("Call operation: NotificationHubs_List")
	client, err := armnotificationhubs.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientNewListPager := client.NewListPager(testsuite.resourceGroupName, testsuite.namespaceName, &armnotificationhubs.ClientListOptions{
		SkipToken: nil,
		Top:       nil,
	})
	for clientNewListPager.More() {
		_, err := clientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NotificationHubs_Get
	fmt.Println("Call operation: NotificationHubs_Get")
	_, err = client.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_Update
	fmt.Println("Call operation: NotificationHubs_Update")
	_, err = client.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, armnotificationhubs.NotificationHubPatchParameters{
		Tags: map[string]*string{
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_GetPnsCredentials
	fmt.Println("Call operation: NotificationHubs_GetPnsCredentials")
	_, err = client.GetPnsCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_CheckNotificationHubAvailability
	fmt.Println("Call operation: NotificationHubs_CheckNotificationHubAvailability")
	_, err = client.CheckNotificationHubAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armnotificationhubs.CheckAvailabilityParameters{
		Name:     to.Ptr("sdktest"),
		Location: to.Ptr("West Europe"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: NotificationHubs_CreateOrUpdateAuthorizationRule")
	_, err = client.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, testsuite.authorizationRuleName, armnotificationhubs.SharedAccessAuthorizationRuleResource{
		Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
			Rights: []*armnotificationhubs.AccessRights{
				to.Ptr(armnotificationhubs.AccessRightsListen),
				to.Ptr(armnotificationhubs.AccessRightsSend),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_ListAuthorizationRules
	fmt.Println("Call operation: NotificationHubs_ListAuthorizationRules")
	clientNewListAuthorizationRulesPager := client.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, nil)
	for clientNewListAuthorizationRulesPager.More() {
		_, err := clientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NotificationHubs_GetAuthorizationRule
	fmt.Println("Call operation: NotificationHubs_GetAuthorizationRule")
	_, err = client.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_RegenerateKeys
	fmt.Println("Call operation: NotificationHubs_RegenerateKeys")
	_, err = client.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, testsuite.authorizationRuleName, armnotificationhubs.PolicyKeyResource{
		PolicyKey: to.Ptr(armnotificationhubs.PolicyKeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_ListKeys
	fmt.Println("Call operation: NotificationHubs_ListKeys")
	_, err = client.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_DeleteAuthorizationRule
	fmt.Println("Call operation: NotificationHubs_DeleteAuthorizationRule")
	_, err = client.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.NotificationHubs/namespaces/{namespaceName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *NotificationhubsTestSuite) TestPrivateEndpointConnections() {
	var privateEndpointConnectionName string
	var err error
	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "eepnotificationhubs-nic",
			},
			"notificationhubsId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.notificationhubsId,
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "epnotificationhubs",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "epnotificationhubsvnet",
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
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
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
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
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
								"privateLinkServiceId": "[parameters('notificationhubsId')]",
							},
						},
					},
					"subnet": map[string]any{
						"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
					},
				},
			},
			map[string]any{
				"name":       "[concat(parameters('virtualNetworksName'), '/default')]",
				"type":       "Microsoft.Network/virtualNetworks/subnets",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
				},
				"properties": map[string]any{
					"addressPrefix":                     "10.0.0.0/24",
					"delegations":                       []any{},
					"privateEndpointNetworkPolicies":    "Disabled",
					"privateLinkServiceNetworkPolicies": "Enabled",
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
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_List
	fmt.Println("Call operation: PrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armnotificationhubs.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Update
	fmt.Println("Call operation: PrivateEndpointConnections_Update")
	privateEndpointConnectionsClientUpdateResponsePoller, err := privateEndpointConnectionsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, privateEndpointConnectionName, armnotificationhubs.PrivateEndpointConnectionResource{
		Properties: &armnotificationhubs.PrivateEndpointConnectionProperties{
			PrivateEndpoint: &armnotificationhubs.RemotePrivateEndpointConnection{},
			PrivateLinkServiceConnectionState: &armnotificationhubs.RemotePrivateLinkServiceConnectionState{
				Status: to.Ptr(armnotificationhubs.PrivateLinkConnectionStatusApproved),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_ListGroupIds
	fmt.Println("Call operation: PrivateEndpointConnections_ListGroupIds")
	privateEndpointConnectionsClientNewListGroupIDsPager := privateEndpointConnectionsClient.NewListGroupIDsPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for privateEndpointConnectionsClientNewListGroupIDsPager.More() {
		_, err := privateEndpointConnectionsClientNewListGroupIDsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_GetGroupId
	fmt.Println("Call operation: PrivateEndpointConnections_GetGroupId")
	_, err = privateEndpointConnectionsClient.GetGroupID(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, "namespace", nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.NotificationHubs/operations
func (testsuite *NotificationhubsTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armnotificationhubs.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *NotificationhubsTestSuite) Cleanup() {
	var err error
	// From step NotificationHubs_Delete
	fmt.Println("Call operation: NotificationHubs_Delete")
	client, err := armnotificationhubs.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_Delete
	fmt.Println("Call operation: Namespaces_Delete")
	namespacesClient, err := armnotificationhubs.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = namespacesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)
}

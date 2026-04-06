// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armsignalr_test

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/signalr/armsignalr"
	"github.com/stretchr/testify/suite"
)

type SignalrTestSuite struct {
	suite.Suite

	ctx                           context.Context
	cred                          azcore.TokenCredential
	options                       *arm.ClientOptions
	name                          string
	armEndpoint                   string
	certificateName               string
	resourceName                  string
	sharedPrivateLinkResourceName string
	signalrId                     string
	storageAccountName            string
	location                      string
	resourceGroupName             string
	subscriptionId                string
}

func (testsuite *SignalrTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.name, _ = recording.GenerateAlphaNumericID(testsuite.T(), "name", 10, false)
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.certificateName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "certific", 14, false)
	testsuite.resourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "resource", 14, false)
	testsuite.sharedPrivateLinkResourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sharedpr", 14, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "signalrsc", 15, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SignalrTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSignalrTestSuite(t *testing.T) {
	suite.Run(t, new(SignalrTestSuite))
}

func (testsuite *SignalrTestSuite) Prepare() {
	var err error
	// From step SignalR_CheckNameAvailability
	fmt.Println("Call operation: SignalR_CheckNameAvailability")
	client, err := armsignalr.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.CheckNameAvailability(testsuite.ctx, testsuite.location, armsignalr.NameAvailabilityParameters{
		Name: to.Ptr(testsuite.resourceName),
		Type: to.Ptr("Microsoft.SignalRService/SignalR"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step SignalR_CreateOrUpdate
	fmt.Println("Call operation: SignalR_CreateOrUpdate")
	clientCreateOrUpdateResponsePoller, err := client.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armsignalr.ResourceInfo{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Identity: &armsignalr.ManagedIdentity{
			Type: to.Ptr(armsignalr.ManagedIdentityTypeSystemAssigned),
		},
		Kind: to.Ptr(armsignalr.ServiceKindSignalR),
		Properties: &armsignalr.Properties{
			Cors: &armsignalr.CorsSettings{
				AllowedOrigins: []*string{
					to.Ptr("https://foo.com"),
					to.Ptr("https://bar.com")},
			},
			DisableAADAuth:   to.Ptr(false),
			DisableLocalAuth: to.Ptr(false),
			Features: []*armsignalr.Feature{
				{
					Flag:       to.Ptr(armsignalr.FeatureFlagsServiceMode),
					Properties: map[string]*string{},
					Value:      to.Ptr("Serverless"),
				},
				{
					Flag:       to.Ptr(armsignalr.FeatureFlagsEnableConnectivityLogs),
					Properties: map[string]*string{},
					Value:      to.Ptr("True"),
				},
				{
					Flag:       to.Ptr(armsignalr.FeatureFlagsEnableMessagingLogs),
					Properties: map[string]*string{},
					Value:      to.Ptr("False"),
				},
				{
					Flag:       to.Ptr(armsignalr.FeatureFlagsEnableLiveTrace),
					Properties: map[string]*string{},
					Value:      to.Ptr("False"),
				}},
			LiveTraceConfiguration: &armsignalr.LiveTraceConfiguration{
				Categories: []*armsignalr.LiveTraceCategory{
					{
						Name:    to.Ptr("ConnectivityLogs"),
						Enabled: to.Ptr("true"),
					}},
				Enabled: to.Ptr("false"),
			},
			NetworkACLs: &armsignalr.NetworkACLs{
				DefaultAction: to.Ptr(armsignalr.ACLActionDeny),
				PrivateEndpoints: []*armsignalr.PrivateEndpointACL{
					{
						Allow: []*armsignalr.SignalRRequestType{
							to.Ptr(armsignalr.SignalRRequestTypeServerConnection)},
						Name: to.Ptr(testsuite.resourceName + ".1fa229cd-bf3f-47f0-8c49-afb36723997e"),
					}},
				PublicNetwork: &armsignalr.NetworkACL{
					Allow: []*armsignalr.SignalRRequestType{
						to.Ptr(armsignalr.SignalRRequestTypeClientConnection)},
				},
			},
			PublicNetworkAccess: to.Ptr("Enabled"),
			Serverless: &armsignalr.ServerlessSettings{
				ConnectionTimeoutInSeconds: to.Ptr[int32](5),
			},
			TLS: &armsignalr.TLSSettings{
				ClientCertEnabled: to.Ptr(false),
			},
			Upstream: &armsignalr.ServerlessUpstreamSettings{
				Templates: []*armsignalr.UpstreamTemplate{
					{
						Auth: &armsignalr.UpstreamAuthSettings{
							Type: to.Ptr(armsignalr.UpstreamAuthTypeManagedIdentity),
							ManagedIdentity: &armsignalr.ManagedIdentitySettings{
								Resource: to.Ptr("api://example"),
							},
						},
						CategoryPattern: to.Ptr("*"),
						EventPattern:    to.Ptr("connect,disconnect"),
						HubPattern:      to.Ptr("*"),
						URLTemplate:     to.Ptr("https://example.com/chat/api/connect"),
					}},
			},
		},
		SKU: &armsignalr.ResourceSKU{
			Name:     to.Ptr("Premium_P1"),
			Capacity: to.Ptr[int32](1),
			Tier:     to.Ptr(armsignalr.SignalRSKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var clientCreateOrUpdateResponse *armsignalr.ClientCreateOrUpdateResponse
	clientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, clientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.signalrId = *clientCreateOrUpdateResponse.ID
}

// Microsoft.SignalRService/signalR/{resourceName}
func (testsuite *SignalrTestSuite) TestSignalR() {
	var err error
	// From step SignalR_ListBySubscription
	fmt.Println("Call operation: SignalR_ListBySubscription")
	client, err := armsignalr.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientNewListBySubscriptionPager := client.NewListBySubscriptionPager(nil)
	for clientNewListBySubscriptionPager.More() {
		_, err := clientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SignalR_ListByResourceGroup
	fmt.Println("Call operation: SignalR_ListByResourceGroup")
	clientNewListByResourceGroupPager := client.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for clientNewListByResourceGroupPager.More() {
		_, err := clientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SignalR_Get
	fmt.Println("Call operation: SignalR_Get")
	_, err = client.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step SignalR_ListSkus
	fmt.Println("Call operation: SignalR_ListSkus")
	_, err = client.ListSKUs(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step SignalR_Update
	fmt.Println("Call operation: SignalR_Update")
	clientUpdateResponsePoller, err := client.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armsignalr.ResourceInfo{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Identity: &armsignalr.ManagedIdentity{
			Type: to.Ptr(armsignalr.ManagedIdentityTypeSystemAssigned),
		},
		Kind: to.Ptr(armsignalr.ServiceKindSignalR),
		Properties: &armsignalr.Properties{
			Cors: &armsignalr.CorsSettings{
				AllowedOrigins: []*string{
					to.Ptr("https://foo.com"),
					to.Ptr("https://bar.com")},
			},
			DisableAADAuth:   to.Ptr(false),
			DisableLocalAuth: to.Ptr(false),
			Features: []*armsignalr.Feature{
				{
					Flag:       to.Ptr(armsignalr.FeatureFlagsServiceMode),
					Properties: map[string]*string{},
					Value:      to.Ptr("Serverless"),
				},
				{
					Flag:       to.Ptr(armsignalr.FeatureFlagsEnableConnectivityLogs),
					Properties: map[string]*string{},
					Value:      to.Ptr("True"),
				},
				{
					Flag:       to.Ptr(armsignalr.FeatureFlagsEnableMessagingLogs),
					Properties: map[string]*string{},
					Value:      to.Ptr("False"),
				},
				{
					Flag:       to.Ptr(armsignalr.FeatureFlagsEnableLiveTrace),
					Properties: map[string]*string{},
					Value:      to.Ptr("False"),
				}},
			LiveTraceConfiguration: &armsignalr.LiveTraceConfiguration{
				Categories: []*armsignalr.LiveTraceCategory{
					{
						Name:    to.Ptr("ConnectivityLogs"),
						Enabled: to.Ptr("true"),
					}},
				Enabled: to.Ptr("false"),
			},
			NetworkACLs: &armsignalr.NetworkACLs{
				DefaultAction: to.Ptr(armsignalr.ACLActionDeny),
				PrivateEndpoints: []*armsignalr.PrivateEndpointACL{
					{
						Allow: []*armsignalr.SignalRRequestType{
							to.Ptr(armsignalr.SignalRRequestTypeServerConnection)},
						Name: to.Ptr(testsuite.resourceName + ".1fa229cd-bf3f-47f0-8c49-afb36723997e"),
					}},
				PublicNetwork: &armsignalr.NetworkACL{
					Allow: []*armsignalr.SignalRRequestType{
						to.Ptr(armsignalr.SignalRRequestTypeClientConnection)},
				},
			},
			PublicNetworkAccess: to.Ptr("Enabled"),
			Serverless: &armsignalr.ServerlessSettings{
				ConnectionTimeoutInSeconds: to.Ptr[int32](5),
			},
			TLS: &armsignalr.TLSSettings{
				ClientCertEnabled: to.Ptr(false),
			},
			Upstream: &armsignalr.ServerlessUpstreamSettings{
				Templates: []*armsignalr.UpstreamTemplate{
					{
						Auth: &armsignalr.UpstreamAuthSettings{
							Type: to.Ptr(armsignalr.UpstreamAuthTypeManagedIdentity),
							ManagedIdentity: &armsignalr.ManagedIdentitySettings{
								Resource: to.Ptr("api://example"),
							},
						},
						CategoryPattern: to.Ptr("*"),
						EventPattern:    to.Ptr("connect,disconnect"),
						HubPattern:      to.Ptr("*"),
						URLTemplate:     to.Ptr("https://example.com/chat/api/connect"),
					}},
			},
		},
		SKU: &armsignalr.ResourceSKU{
			Name:     to.Ptr("Premium_P1"),
			Capacity: to.Ptr[int32](1),
			Tier:     to.Ptr(armsignalr.SignalRSKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step SignalR_Restart
	fmt.Println("Call operation: SignalR_Restart")
	clientRestartResponsePoller, err := client.BeginRestart(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientRestartResponsePoller)
	testsuite.Require().NoError(err)

	// From step SignalR_RegenerateKey
	fmt.Println("Call operation: SignalR_RegenerateKey")
	clientRegenerateKeyResponsePoller, err := client.BeginRegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armsignalr.RegenerateKeyParameters{
		KeyType: to.Ptr(armsignalr.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientRegenerateKeyResponsePoller)
	testsuite.Require().NoError(err)

	// From step SignalR_ListKeys
	fmt.Println("Call operation: SignalR_ListKeys")
	_, err = client.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.SignalRService/signalR/{resourceName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *SignalrTestSuite) TestSignalRPrivateEndpointConnections() {
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
				"defaultValue": "epsignalr-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "epsignalr",
			},
			"signalrId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.signalrId,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "epsignalrvnet",
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
									"signalr",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('signalrId')]",
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

	// From step SignalRPrivateEndpointConnections_List
	fmt.Println("Call operation: SignalRPrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armsignalr.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step SignalRPrivateEndpointConnections_Update
	fmt.Println("Call operation: SignalRPrivateEndpointConnections_Update")
	_, err = privateEndpointConnectionsClient.Update(testsuite.ctx, privateEndpointConnectionName, testsuite.resourceGroupName, testsuite.resourceName, armsignalr.PrivateEndpointConnection{
		Properties: &armsignalr.PrivateEndpointConnectionProperties{
			PrivateEndpoint: &armsignalr.PrivateEndpoint{
				ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourcegroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/privateEndpoints/myPrivateEndpoint"),
			},
			PrivateLinkServiceConnectionState: &armsignalr.PrivateLinkServiceConnectionState{
				ActionsRequired: to.Ptr("None"),
				Status:          to.Ptr(armsignalr.PrivateLinkServiceConnectionStatusApproved),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step SignalRPrivateEndpointConnections_Get
	fmt.Println("Call operation: SignalRPrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, privateEndpointConnectionName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step SignalRPrivateEndpointConnections_Delete
	fmt.Println("Call operation: SignalRPrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, privateEndpointConnectionName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.SignalRService/signalR/{resourceName}/sharedPrivateLinkResources/{sharedPrivateLinkResourceName}
func (testsuite *SignalrTestSuite) TestSignalRSharedPrivateLinkResources() {
	var storageAccountId string
	var err error
	// From step Create_StorageAccount
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"storageAccountId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"storageAccountName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.storageAccountName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('storageAccountName')]",
				"type":       "Microsoft.Storage/storageAccounts",
				"apiVersion": "2022-05-01",
				"kind":       "StorageV2",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"accessTier":                   "Hot",
					"allowBlobPublicAccess":        true,
					"allowCrossTenantReplication":  true,
					"allowSharedKeyAccess":         true,
					"defaultToOAuthAuthentication": false,
					"dnsEndpointType":              "Standard",
					"encryption": map[string]any{
						"keySource":                       "Microsoft.Storage",
						"requireInfrastructureEncryption": false,
						"services": map[string]any{
							"blob": map[string]any{
								"enabled": true,
								"keyType": "Account",
							},
							"file": map[string]any{
								"enabled": true,
								"keyType": "Account",
							},
						},
					},
					"minimumTlsVersion": "TLS1_2",
					"networkAcls": map[string]any{
						"bypass":              "AzureServices",
						"defaultAction":       "Allow",
						"ipRules":             []any{},
						"virtualNetworkRules": []any{},
					},
					"publicNetworkAccess":      "Enabled",
					"supportsHttpsTrafficOnly": true,
				},
				"sku": map[string]any{
					"name": "Standard_RAGRS",
					"tier": "Standard",
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_StorageAccount", &deployment)
	testsuite.Require().NoError(err)
	storageAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountId"].(map[string]interface{})["value"].(string)

	// From step SignalRSharedPrivateLinkResources_CreateOrUpdate
	fmt.Println("Call operation: SignalRSharedPrivateLinkResources_CreateOrUpdate")
	sharedPrivateLinkResourcesClient, err := armsignalr.NewSharedPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sharedPrivateLinkResourcesClientCreateOrUpdateResponsePoller, err := sharedPrivateLinkResourcesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.sharedPrivateLinkResourceName, testsuite.resourceGroupName, testsuite.resourceName, armsignalr.SharedPrivateLinkResource{
		Properties: &armsignalr.SharedPrivateLinkResourceProperties{
			GroupID:               to.Ptr("table"),
			PrivateLinkResourceID: to.Ptr(storageAccountId),
			RequestMessage:        to.Ptr("Please approve"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sharedPrivateLinkResourcesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step SignalRSharedPrivateLinkResources_List
	fmt.Println("Call operation: SignalRSharedPrivateLinkResources_List")
	sharedPrivateLinkResourcesClientNewListPager := sharedPrivateLinkResourcesClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for sharedPrivateLinkResourcesClientNewListPager.More() {
		_, err := sharedPrivateLinkResourcesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SignalRSharedPrivateLinkResources_Get
	fmt.Println("Call operation: SignalRSharedPrivateLinkResources_Get")
	_, err = sharedPrivateLinkResourcesClient.Get(testsuite.ctx, testsuite.sharedPrivateLinkResourceName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step SignalRSharedPrivateLinkResources_Delete
	fmt.Println("Call operation: SignalRSharedPrivateLinkResources_Delete")
	sharedPrivateLinkResourcesClientDeleteResponsePoller, err := sharedPrivateLinkResourcesClient.BeginDelete(testsuite.ctx, testsuite.sharedPrivateLinkResourceName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sharedPrivateLinkResourcesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.SignalRService/operations
func (testsuite *SignalrTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armsignalr.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.SignalRService/locations/{location}/usages
func (testsuite *SignalrTestSuite) TestUsages() {
	var err error
	// From step Usages_List
	fmt.Println("Call operation: Usages_List")
	usagesClient, err := armsignalr.NewUsagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	usagesClientNewListPager := usagesClient.NewListPager(testsuite.location, nil)
	for usagesClientNewListPager.More() {
		_, err := usagesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.SignalRService/signalR/{resourceName}/privateLinkResources
func (testsuite *SignalrTestSuite) TestSignalRPrivateLinkResources() {
	var err error
	// From step SignalRPrivateLinkResources_List
	fmt.Println("Call operation: SignalRPrivateLinkResources_List")
	privateLinkResourcesClient, err := armsignalr.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListPager := privateLinkResourcesClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for privateLinkResourcesClientNewListPager.More() {
		_, err := privateLinkResourcesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *SignalrTestSuite) Cleanup() {
	var err error
	// From step SignalR_Delete
	fmt.Println("Call operation: SignalR_Delete")
	client, err := armsignalr.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientDeleteResponsePoller, err := client.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

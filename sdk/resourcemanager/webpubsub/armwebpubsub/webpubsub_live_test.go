// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armwebpubsub_test

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/webpubsub/armwebpubsub"
	"github.com/stretchr/testify/suite"
)

type WebpubsubTestSuite struct {
	suite.Suite

	ctx                           context.Context
	cred                          azcore.TokenCredential
	options                       *arm.ClientOptions
	name                          string
	armEndpoint                   string
	certificateName               string
	hubName                       string
	resourceName                  string
	serverfarmsName               string
	sharedPrivateLinkResourceName string
	sitesName                     string
	webpubsubId                   string
	location                      string
	resourceGroupName             string
	subscriptionId                string
}

func (testsuite *WebpubsubTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.name, _ = recording.GenerateAlphaNumericID(testsuite.T(), "customdomain", 18, false)
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.certificateName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "certific", 14, false)
	testsuite.hubName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "hubname", 13, false)
	testsuite.resourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "resource", 14, false)
	testsuite.serverfarmsName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serverfarm", 16, false)
	testsuite.sharedPrivateLinkResourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sharedpr", 14, false)
	testsuite.sitesName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sitena", 12, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *WebpubsubTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestWebpubsubTestSuite(t *testing.T) {
	suite.Run(t, new(WebpubsubTestSuite))
}

func (testsuite *WebpubsubTestSuite) Prepare() {
	var err error
	// From step WebPubSub_CreateOrUpdate
	fmt.Println("Call operation: WebPubSub_CreateOrUpdate")
	client, err := armwebpubsub.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientCreateOrUpdateResponsePoller, err := client.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armwebpubsub.ResourceInfo{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Identity: &armwebpubsub.ManagedIdentity{
			Type: to.Ptr(armwebpubsub.ManagedIdentityTypeSystemAssigned),
		},
		Properties: &armwebpubsub.Properties{
			DisableAADAuth:   to.Ptr(false),
			DisableLocalAuth: to.Ptr(false),
			LiveTraceConfiguration: &armwebpubsub.LiveTraceConfiguration{
				Categories: []*armwebpubsub.LiveTraceCategory{
					{
						Name:    to.Ptr("ConnectivityLogs"),
						Enabled: to.Ptr("true"),
					}},
				Enabled: to.Ptr("false"),
			},
			NetworkACLs: &armwebpubsub.NetworkACLs{
				DefaultAction: to.Ptr(armwebpubsub.ACLActionDeny),
				PrivateEndpoints: []*armwebpubsub.PrivateEndpointACL{
					{
						Allow: []*armwebpubsub.WebPubSubRequestType{
							to.Ptr(armwebpubsub.WebPubSubRequestTypeServerConnection)},
						Name: to.Ptr(testsuite.resourceName + ".00000000-0000-0000-0000-000000000000"),
					}},
				PublicNetwork: &armwebpubsub.NetworkACL{
					Allow: []*armwebpubsub.WebPubSubRequestType{
						to.Ptr(armwebpubsub.WebPubSubRequestTypeClientConnection)},
				},
			},
			PublicNetworkAccess: to.Ptr("Enabled"),
			TLS: &armwebpubsub.TLSSettings{
				ClientCertEnabled: to.Ptr(false),
			},
		},
		SKU: &armwebpubsub.ResourceSKU{
			Name:     to.Ptr("Premium_P1"),
			Capacity: to.Ptr[int32](1),
			Tier:     to.Ptr(armwebpubsub.WebPubSubSKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var clientCreateOrUpdateResponse *armwebpubsub.ClientCreateOrUpdateResponse
	clientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, clientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.webpubsubId = *clientCreateOrUpdateResponse.ID
}

// Microsoft.SignalRService/webPubSub/{resourceName}
func (testsuite *WebpubsubTestSuite) TestWebPubSub() {
	var err error
	// From step WebPubSub_CheckNameAvailability
	fmt.Println("Call operation: WebPubSub_CheckNameAvailability")
	client, err := armwebpubsub.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.CheckNameAvailability(testsuite.ctx, testsuite.location, armwebpubsub.NameAvailabilityParameters{
		Name: to.Ptr("myWebPubSubService"),
		Type: to.Ptr("Microsoft.SignalRService/WebPubSub"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step WebPubSub_ListBySubscription
	fmt.Println("Call operation: WebPubSub_ListBySubscription")
	clientNewListBySubscriptionPager := client.NewListBySubscriptionPager(nil)
	for clientNewListBySubscriptionPager.More() {
		_, err := clientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WebPubSub_ListByResourceGroup
	fmt.Println("Call operation: WebPubSub_ListByResourceGroup")
	clientNewListByResourceGroupPager := client.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for clientNewListByResourceGroupPager.More() {
		_, err := clientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WebPubSub_ListSkus
	fmt.Println("Call operation: WebPubSub_ListSkus")
	_, err = client.ListSKUs(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step WebPubSub_Get
	fmt.Println("Call operation: WebPubSub_Get")
	_, err = client.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step WebPubSub_Update
	fmt.Println("Call operation: WebPubSub_Update")
	clientUpdateResponsePoller, err := client.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armwebpubsub.ResourceInfo{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Identity: &armwebpubsub.ManagedIdentity{
			Type: to.Ptr(armwebpubsub.ManagedIdentityTypeSystemAssigned),
		},
		Properties: &armwebpubsub.Properties{
			DisableAADAuth:   to.Ptr(false),
			DisableLocalAuth: to.Ptr(false),
			LiveTraceConfiguration: &armwebpubsub.LiveTraceConfiguration{
				Categories: []*armwebpubsub.LiveTraceCategory{
					{
						Name:    to.Ptr("ConnectivityLogs"),
						Enabled: to.Ptr("true"),
					}},
				Enabled: to.Ptr("false"),
			},
			NetworkACLs: &armwebpubsub.NetworkACLs{
				DefaultAction: to.Ptr(armwebpubsub.ACLActionDeny),
				PrivateEndpoints: []*armwebpubsub.PrivateEndpointACL{
					{
						Allow: []*armwebpubsub.WebPubSubRequestType{
							to.Ptr(armwebpubsub.WebPubSubRequestTypeServerConnection)},
						Name: to.Ptr(testsuite.resourceName + ".1fa229cd-bf3f-47f0-8c49-afb36723997e"),
					}},
				PublicNetwork: &armwebpubsub.NetworkACL{
					Allow: []*armwebpubsub.WebPubSubRequestType{
						to.Ptr(armwebpubsub.WebPubSubRequestTypeClientConnection)},
				},
			},
			PublicNetworkAccess: to.Ptr("Enabled"),
			TLS: &armwebpubsub.TLSSettings{
				ClientCertEnabled: to.Ptr(false),
			},
		},
		SKU: &armwebpubsub.ResourceSKU{
			Name:     to.Ptr("Premium_P1"),
			Capacity: to.Ptr[int32](1),
			Tier:     to.Ptr(armwebpubsub.WebPubSubSKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step WebPubSub_Restart
	fmt.Println("Call operation: WebPubSub_Restart")
	clientRestartResponsePoller, err := client.BeginRestart(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientRestartResponsePoller)
	testsuite.Require().NoError(err)

	// From step WebPubSub_RegenerateKey
	fmt.Println("Call operation: WebPubSub_RegenerateKey")
	clientRegenerateKeyResponsePoller, err := client.BeginRegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armwebpubsub.RegenerateKeyParameters{
		KeyType: to.Ptr(armwebpubsub.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientRegenerateKeyResponsePoller)
	testsuite.Require().NoError(err)

	// From step WebPubSub_ListKeys
	fmt.Println("Call operation: WebPubSub_ListKeys")
	_, err = client.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.SignalRService/webPubSub/{resourceName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *WebpubsubTestSuite) TestWebPubSubPrivateEndpointConnections() {
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
				"defaultValue": "epwebpubsub-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "epwebpubsub",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "epwebpubsubvnet",
			},
			"webpubsubId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.webpubsubId,
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
									"webpubsub",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('webpubsubId')]",
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

	// From step WebPubSubPrivateEndpointConnections_List
	fmt.Println("Call operation: WebPubSubPrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armwebpubsub.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step WebPubSubPrivateEndpointConnections_Update
	fmt.Println("Call operation: WebPubSubPrivateEndpointConnections_Update")
	_, err = privateEndpointConnectionsClient.Update(testsuite.ctx, privateEndpointConnectionName, testsuite.resourceGroupName, testsuite.resourceName, armwebpubsub.PrivateEndpointConnection{
		Properties: &armwebpubsub.PrivateEndpointConnectionProperties{
			PrivateEndpoint: &armwebpubsub.PrivateEndpoint{
				ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourcegroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/privateEndpoints/myPrivateEndpoint"),
			},
			PrivateLinkServiceConnectionState: &armwebpubsub.PrivateLinkServiceConnectionState{
				ActionsRequired: to.Ptr("None"),
				Status:          to.Ptr(armwebpubsub.PrivateLinkServiceConnectionStatusApproved),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step WebPubSubPrivateEndpointConnections_Get
	fmt.Println("Call operation: WebPubSubPrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, privateEndpointConnectionName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step WebPubSubPrivateEndpointConnections_Delete
	fmt.Println("Call operation: WebPubSubPrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, privateEndpointConnectionName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.SignalRService/webPubSub/{resourceName}/sharedPrivateLinkResources/{sharedPrivateLinkResourceName}
func (testsuite *WebpubsubTestSuite) TestWebPubSubSharedPrivateLinkResources() {
	var webAppId string
	var err error
	// From step Create_WebApp
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"webAppId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Web/sites', parameters('sitesName'))]",
			},
		},
		"parameters": map[string]any{
			"serverfarmsName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.serverfarmsName,
			},
			"sitesName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.sitesName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('serverfarmsName')]",
				"type":       "Microsoft.Web/serverfarms",
				"apiVersion": "2022-09-01",
				"kind":       "linux",
				"location":   "East US",
				"properties": map[string]any{
					"elasticScaleEnabled":       false,
					"hyperV":                    false,
					"isSpot":                    false,
					"isXenon":                   false,
					"maximumElasticWorkerCount": float64(1),
					"perSiteScaling":            false,
					"reserved":                  true,
					"targetWorkerCount":         float64(0),
					"targetWorkerSizeId":        float64(0),
					"zoneRedundant":             false,
				},
				"sku": map[string]any{
					"name":     "P1v3",
					"capacity": float64(1),
					"family":   "Pv3",
					"size":     "P1v3",
					"tier":     "PremiumV3",
				},
			},
			map[string]any{
				"name":       "[parameters('sitesName')]",
				"type":       "Microsoft.Web/sites",
				"apiVersion": "2022-09-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Web/serverfarms', parameters('serverfarmsName'))]",
				},
				"kind":     "app",
				"location": "East US",
				"properties": map[string]any{
					"enabled":      true,
					"serverFarmId": "[resourceId('Microsoft.Web/serverfarms', parameters('serverfarmsName'))]",
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_WebApp", &deployment)
	testsuite.Require().NoError(err)
	webAppId = deploymentExtend.Properties.Outputs.(map[string]interface{})["webAppId"].(map[string]interface{})["value"].(string)

	// From step WebPubSubSharedPrivateLinkResources_CreateOrUpdate
	fmt.Println("Call operation: WebPubSubSharedPrivateLinkResources_CreateOrUpdate")
	sharedPrivateLinkResourcesClient, err := armwebpubsub.NewSharedPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sharedPrivateLinkResourcesClientCreateOrUpdateResponsePoller, err := sharedPrivateLinkResourcesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.sharedPrivateLinkResourceName, testsuite.resourceGroupName, testsuite.resourceName, armwebpubsub.SharedPrivateLinkResource{
		Properties: &armwebpubsub.SharedPrivateLinkResourceProperties{
			GroupID:               to.Ptr("sites"),
			PrivateLinkResourceID: to.Ptr(webAppId),
			RequestMessage:        to.Ptr("Please approve"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sharedPrivateLinkResourcesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step WebPubSubSharedPrivateLinkResources_List
	fmt.Println("Call operation: WebPubSubSharedPrivateLinkResources_List")
	sharedPrivateLinkResourcesClientNewListPager := sharedPrivateLinkResourcesClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for sharedPrivateLinkResourcesClientNewListPager.More() {
		_, err := sharedPrivateLinkResourcesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WebPubSubSharedPrivateLinkResources_Get
	fmt.Println("Call operation: WebPubSubSharedPrivateLinkResources_Get")
	_, err = sharedPrivateLinkResourcesClient.Get(testsuite.ctx, testsuite.sharedPrivateLinkResourceName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step WebPubSubSharedPrivateLinkResources_Delete
	fmt.Println("Call operation: WebPubSubSharedPrivateLinkResources_Delete")
	sharedPrivateLinkResourcesClientDeleteResponsePoller, err := sharedPrivateLinkResourcesClient.BeginDelete(testsuite.ctx, testsuite.sharedPrivateLinkResourceName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sharedPrivateLinkResourcesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.SignalRService/webPubSub/{resourceName}/hubs/{hubName}
func (testsuite *WebpubsubTestSuite) TestWebPubSubHubs() {
	var err error
	// From step WebPubSubHubs_CreateOrUpdate
	fmt.Println("Call operation: WebPubSubHubs_CreateOrUpdate")
	hubsClient, err := armwebpubsub.NewHubsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	hubsClientCreateOrUpdateResponsePoller, err := hubsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.hubName, testsuite.resourceGroupName, testsuite.resourceName, armwebpubsub.Hub{
		Properties: &armwebpubsub.HubProperties{
			AnonymousConnectPolicy: to.Ptr("allow"),
			EventHandlers: []*armwebpubsub.EventHandler{
				{
					Auth: &armwebpubsub.UpstreamAuthSettings{
						Type: to.Ptr(armwebpubsub.UpstreamAuthTypeManagedIdentity),
						ManagedIdentity: &armwebpubsub.ManagedIdentitySettings{
							Resource: to.Ptr("abc"),
						},
					},
					SystemEvents: []*string{
						to.Ptr("connect"),
						to.Ptr("connected")},
					URLTemplate:      to.Ptr("http://host.com"),
					UserEventPattern: to.Ptr("*"),
				}},
			EventListeners: []*armwebpubsub.EventListener{
				{
					Endpoint: &armwebpubsub.EventHubEndpoint{
						Type:                    to.Ptr(armwebpubsub.EventListenerEndpointDiscriminatorEventHub),
						EventHubName:            to.Ptr("eventHubName1"),
						FullyQualifiedNamespace: to.Ptr("example.servicebus.windows.net"),
					},
					Filter: &armwebpubsub.EventNameFilter{
						Type: to.Ptr(armwebpubsub.EventListenerFilterDiscriminatorEventName),
						SystemEvents: []*string{
							to.Ptr("connected"),
							to.Ptr("disconnected")},
						UserEventPattern: to.Ptr("*"),
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, hubsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step WebPubSubHubs_List
	fmt.Println("Call operation: WebPubSubHubs_List")
	hubsClientNewListPager := hubsClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for hubsClientNewListPager.More() {
		_, err := hubsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WebPubSubHubs_Get
	fmt.Println("Call operation: WebPubSubHubs_Get")
	_, err = hubsClient.Get(testsuite.ctx, testsuite.hubName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step WebPubSubHubs_Delete
	fmt.Println("Call operation: WebPubSubHubs_Delete")
	hubsClientDeleteResponsePoller, err := hubsClient.BeginDelete(testsuite.ctx, testsuite.hubName, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, hubsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.SignalRService/operations
func (testsuite *WebpubsubTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armwebpubsub.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Usages_List
	fmt.Println("Call operation: Usages_List")
	usagesClient, err := armwebpubsub.NewUsagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	usagesClientNewListPager := usagesClient.NewListPager(testsuite.location, nil)
	for usagesClientNewListPager.More() {
		_, err := usagesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WebPubSubPrivateLinkResources_List
	fmt.Println("Call operation: WebPubSubPrivateLinkResources_List")
	privateLinkResourcesClient, err := armwebpubsub.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListPager := privateLinkResourcesClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for privateLinkResourcesClientNewListPager.More() {
		_, err := privateLinkResourcesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *WebpubsubTestSuite) Cleanup() {
	var err error
	// From step WebPubSub_Delete
	fmt.Println("Call operation: WebPubSub_Delete")
	client, err := armwebpubsub.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientDeleteResponsePoller, err := client.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

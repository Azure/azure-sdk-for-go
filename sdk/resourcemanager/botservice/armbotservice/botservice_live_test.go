// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armbotservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/botservice/armbotservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type BotserviceTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	botServiceId      string
	connectionName    string
	resourceName      string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *BotserviceTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.connectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "connecti", 14, false)
	testsuite.resourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "resource", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *BotserviceTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestBotserviceTestSuite(t *testing.T) {
	suite.Run(t, new(BotserviceTestSuite))
}

func (testsuite *BotserviceTestSuite) Prepare() {
	var err error
	// From step Bots_GetCheckNameAvailability
	fmt.Println("Call operation: Bots_GetCheckNameAvailability")
	botsClient, err := armbotservice.NewBotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = botsClient.GetCheckNameAvailability(testsuite.ctx, armbotservice.CheckNameAvailabilityRequestBody{
		Name: to.Ptr("testbotname"),
		Type: to.Ptr("string"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Bots_Create
	fmt.Println("Call operation: Bots_Create")
	botsClientCreateResponse, err := botsClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armbotservice.Bot{
		Etag:     to.Ptr("etag1"),
		Kind:     to.Ptr(armbotservice.KindSdk),
		Location: to.Ptr("global"),
		SKU: &armbotservice.SKU{
			Name: to.Ptr(armbotservice.SKUNameS1),
		},
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
		Properties: &armbotservice.BotProperties{
			Description: to.Ptr("The description of the bot"),
			DisplayName: to.Ptr("The Name of the bot"),
			Endpoint:    to.Ptr("https://bing.com/messages/"),
			MsaAppID:    to.Ptr("00000000-0000-0000-0000-000000000001"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.botServiceId = *botsClientCreateResponse.ID
}

// Microsoft.BotService/botServices/{resourceName}
func (testsuite *BotserviceTestSuite) TestBots() {
	var err error
	// From step Bots_List
	fmt.Println("Call operation: Bots_List")
	botsClient, err := armbotservice.NewBotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	botsClientNewListPager := botsClient.NewListPager(nil)
	for botsClientNewListPager.More() {
		_, err := botsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Bots_Get
	fmt.Println("Call operation: Bots_Get")
	_, err = botsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step Bots_ListByResourceGroup
	fmt.Println("Call operation: Bots_ListByResourceGroup")
	botsClientNewListByResourceGroupPager := botsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for botsClientNewListByResourceGroupPager.More() {
		_, err := botsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Bots_Update
	fmt.Println("Call operation: Bots_Update")
	_, err = botsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armbotservice.Bot{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.BotService/botServices/{resourceName}/channels/{channelName}
func (testsuite *BotserviceTestSuite) TestChannels() {
	var err error
	// From step Channels_Create
	fmt.Println("Call operation: Channels_Create")
	channelsClient, err := armbotservice.NewChannelsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = channelsClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armbotservice.ChannelNameMsTeamsChannel, armbotservice.BotChannel{
		Location: to.Ptr("global"),
		Properties: &armbotservice.MsTeamsChannel{
			ChannelName: to.Ptr("MsTeamsChannel"),
			Properties: &armbotservice.MsTeamsChannelProperties{
				IsEnabled: to.Ptr(true),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Channels_ListByResourceGroup
	fmt.Println("Call operation: Channels_ListByResourceGroup")
	channelsClientNewListByResourceGroupPager := channelsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for channelsClientNewListByResourceGroupPager.More() {
		_, err := channelsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Channels_Get
	fmt.Println("Call operation: Channels_Get")
	_, err = channelsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, "MsTeamsChannel", nil)
	testsuite.Require().NoError(err)

	// From step Channels_Update
	fmt.Println("Call operation: Channels_Update")
	_, err = channelsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armbotservice.ChannelNameMsTeamsChannel, armbotservice.BotChannel{
		Location: to.Ptr("global"),
		Properties: &armbotservice.MsTeamsChannel{
			ChannelName: to.Ptr("MsTeamsChannel"),
			Properties: &armbotservice.MsTeamsChannelProperties{
				IsEnabled: to.Ptr(true),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Channels_ListWithKeys
	fmt.Println("Call operation: Channels_ListWithKeys")
	_, err = channelsClient.ListWithKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armbotservice.ChannelNameMsTeamsChannel, nil)
	testsuite.Require().NoError(err)

	// From step DirectLine_RegenerateKeys
	fmt.Println("Call operation: DirectLine_RegenerateKeys")
	directLineClient, err := armbotservice.NewDirectLineClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = directLineClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armbotservice.RegenerateKeysChannelNameDirectLineChannel, armbotservice.SiteInfo{
		Key:      to.Ptr(armbotservice.KeyKey1),
		SiteName: to.Ptr("testSiteName"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Channels_Delete
	fmt.Println("Call operation: Channels_Delete")
	_, err = channelsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, "MsTeamsChannel", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.BotService/botServices/{resourceName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *BotserviceTestSuite) TestPrivateEndpointConnections() {
	var privateEndpointConnectionName string
	var err error
	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"privateEndpointConnectionName": map[string]any{
				"type":  "string",
				"value": "[parameters('privateEndpointName')]",
			},
		},
		"parameters": map[string]any{
			"botServiceId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.botServiceId,
			},
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "epbotservice-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "epbotservice",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "botservicevnet",
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
									"bot",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('botServiceId')]",
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)
	privateEndpointConnectionName = deploymentExtend.Properties.Outputs.(map[string]interface{})["privateEndpointConnectionName"].(map[string]interface{})["value"].(string)

	// From step PrivateEndpointConnections_List
	fmt.Println("Call operation: PrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armbotservice.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		_, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateEndpointConnections_Create
	fmt.Println("Call operation: PrivateEndpointConnections_Create")
	_, err = privateEndpointConnectionsClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, privateEndpointConnectionName, armbotservice.PrivateEndpointConnection{
		Properties: &armbotservice.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armbotservice.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Auto-Approved"),
				Status:      to.Ptr(armbotservice.PrivateEndpointServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_ListByBotResource
	fmt.Println("Call operation: PrivateLinkResources_ListByBotResource")
	privateLinkResourcesClient, err := armbotservice.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = privateLinkResourcesClient.ListByBotResource(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	_, err = privateEndpointConnectionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.BotService/operations
func (testsuite *BotserviceTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armbotservice.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.BotService/hostSettings
func (testsuite *BotserviceTestSuite) TestHostSettings() {
	var err error
	// From step HostSettings_Get
	fmt.Println("Call operation: HostSettings_Get")
	hostSettingsClient, err := armbotservice.NewHostSettingsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = hostSettingsClient.Get(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.BotService/botServices/{resourceName}/createEmailSignInUrl
func (testsuite *BotserviceTestSuite) TestEmail() {
	var err error
	// From step Email_CreateSignInUrl
	fmt.Println("Call operation: Email_CreateSignInUrl")
	emailClient, err := armbotservice.NewEmailClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = emailClient.CreateSignInURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *BotserviceTestSuite) Cleanup() {
	var err error
	// From step Bots_Delete
	fmt.Println("Call operation: Bots_Delete")
	botsClient, err := armbotservice.NewBotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = botsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armpostgresql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type PrivateEndpointConnectionsTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	postgresqlserverId string
	serverName         string
	adminPassword      string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *PrivateEndpointConnectionsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serverna", 14, true)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *PrivateEndpointConnectionsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPrivateEndpointConnectionsTestSuite(t *testing.T) {
	suite.Run(t, new(PrivateEndpointConnectionsTestSuite))
}

func (testsuite *PrivateEndpointConnectionsTestSuite) Prepare() {
	var err error
	// From step Servers_Create
	fmt.Println("Call operation: Servers_Create")
	serversClient, err := armpostgresql.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClientCreateResponsePoller, err := serversClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armpostgresql.ServerForCreate{
		Location: to.Ptr(testsuite.location),
		Properties: &armpostgresql.ServerPropertiesForDefaultCreate{
			CreateMode:        to.Ptr(armpostgresql.CreateModeDefault),
			MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS12),
			SSLEnforcement:    to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
			StorageProfile: &armpostgresql.StorageProfile{
				BackupRetentionDays: to.Ptr[int32](7),
				GeoRedundantBackup:  to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
				StorageMB:           to.Ptr[int32](128000),
			},
			AdministratorLogin:         to.Ptr("cloudsa"),
			AdministratorLoginPassword: to.Ptr(testsuite.adminPassword),
		},
		SKU: &armpostgresql.SKU{
			Name:   to.Ptr("GP_Gen5_8"),
			Family: to.Ptr("Gen5"),
			Tier:   to.Ptr(armpostgresql.SKUTierGeneralPurpose),
		},
		Tags: map[string]*string{
			"ElasticServer": to.Ptr("1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var serversClientCreateResponse *armpostgresql.ServersClientCreateResponse
	serversClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, serversClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.postgresqlserverId = *serversClientCreateResponse.ID
}

// Microsoft.DBforPostgreSQL/servers/{serverName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *PrivateEndpointConnectionsTestSuite) TestPrivateEndpointConnections() {
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
				"defaultValue": "pepostgresql-nic",
			},
			"postgresqlserverId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.postgresqlserverId,
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "pepostgresql",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "pepostgresqlvnet",
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
									"postgresqlServer",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('postgresqlserverId')]",
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
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)

	var privateEndpointConnectionName string
	// From step PrivateEndpointConnections_ListByServer
	fmt.Println("Call operation: PrivateEndpointConnections_ListByServer")
	privateEndpointConnectionsClient, err := armpostgresql.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListByServerPager := privateEndpointConnectionsClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for privateEndpointConnectionsClientNewListByServerPager.More() {
		result, err := privateEndpointConnectionsClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		privateEndpointConnectionName = *result.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	var privateLinkResourceName string
	// from step PrivateLinkResources_ListByServer
	fmt.Println("Call operation: PrivateLinkResources_ListByServer")
	privateLinkResourcesClient, err := armpostgresql.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListByServerPager := privateLinkResourcesClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for privateLinkResourcesClientNewListByServerPager.More() {
		result, err := privateLinkResourcesClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		privateLinkResourceName = *result.Value[0].Name
		break
	}

	// From step PrivateLinkResources_Get
	fmt.Println("Call operation: PrivateLinkResources_Get")
	_, err = privateLinkResourcesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, privateLinkResourceName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

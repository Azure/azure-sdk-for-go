//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armattestation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/attestation/armattestation/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type AttestationTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	armEndpoint		string
	attestationId		string
	providerName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *AttestationTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.providerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "provider", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AttestationTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAttestationTestSuite(t *testing.T) {
	suite.Run(t, new(AttestationTestSuite))
}

func (testsuite *AttestationTestSuite) Prepare() {
	var err error
	// From step AttestationProviders_Create
	fmt.Println("Call operation: AttestationProviders_Create")
	providersClient, err := armattestation.NewProvidersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	providersClientCreateResponse, err := providersClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.providerName, armattestation.ServiceCreationParams{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.attestationId = *providersClientCreateResponse.ID
}

// Microsoft.Attestation/attestationProviders/{providerName}
func (testsuite *AttestationTestSuite) TestAttestationProviders() {
	var err error
	// From step AttestationProviders_List
	fmt.Println("Call operation: AttestationProviders_List")
	providersClient, err := armattestation.NewProvidersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = providersClient.List(testsuite.ctx, nil)
	testsuite.Require().NoError(err)

	// From step AttestationProviders_ListDefault
	fmt.Println("Call operation: AttestationProviders_ListDefault")
	_, err = providersClient.ListDefault(testsuite.ctx, nil)
	testsuite.Require().NoError(err)

	// From step AttestationProviders_Get
	fmt.Println("Call operation: AttestationProviders_Get")
	_, err = providersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.providerName, nil)
	testsuite.Require().NoError(err)

	// From step AttestationProviders_GetDefaultByLocation
	fmt.Println("Call operation: AttestationProviders_GetDefaultByLocation")
	_, err = providersClient.GetDefaultByLocation(testsuite.ctx, "Central US", nil)
	testsuite.Require().NoError(err)

	// From step AttestationProviders_ListByResourceGroup
	fmt.Println("Call operation: AttestationProviders_ListByResourceGroup")
	_, err = providersClient.ListByResourceGroup(testsuite.ctx, testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)

	// From step AttestationProviders_Update
	fmt.Println("Call operation: AttestationProviders_Update")
	_, err = providersClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.providerName, armattestation.ServicePatchParams{
		Tags: map[string]*string{
			"Property1":	to.Ptr("Value1"),
			"Property2":	to.Ptr("Value2"),
			"Property3":	to.Ptr("Value3"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Attestation/attestationProviders/{providerName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *AttestationTestSuite) TestPrivateEndpointConnections() {
	var privateEndpointConnectionId string
	var privateEndpointConnectionName string
	var err error
	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":		"https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion":	"1.0.0.0",
		"parameters": map[string]any{
			"attestationId": map[string]any{
				"type":		"string",
				"defaultValue":	testsuite.attestationId,
			},
			"location": map[string]any{
				"type":		"string",
				"defaultValue":	testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":		"string",
				"defaultValue":	"epattestation-nic",
			},
			"privateEndpointName": map[string]any{
				"type":		"string",
				"defaultValue":	"epattestation",
			},
			"virtualNetworksName": map[string]any{
				"type":		"string",
				"defaultValue":	"epattestationvnet",
			},
		},
		"resources": []any{
			map[string]any{
				"name":		"[parameters('virtualNetworksName')]",
				"type":		"Microsoft.Network/virtualNetworks",
				"apiVersion":	"2020-11-01",
				"location":	"[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"enableDdosProtection":	false,
					"subnets": []any{
						map[string]any{
							"name":	"default",
							"properties": map[string]any{
								"addressPrefix":			"10.0.0.0/24",
								"delegations":				[]any{},
								"privateEndpointNetworkPolicies":	"Disabled",
								"privateLinkServiceNetworkPolicies":	"Enabled",
							},
						},
					},
					"virtualNetworkPeerings":	[]any{},
				},
			},
			map[string]any{
				"name":		"[parameters('networkInterfaceName')]",
				"type":		"Microsoft.Network/networkInterfaces",
				"apiVersion":	"2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location":	"[parameters('location')]",
				"properties": map[string]any{
					"dnsSettings": map[string]any{
						"dnsServers": []any{},
					},
					"enableIPForwarding":	false,
					"ipConfigurations": []any{
						map[string]any{
							"name":	"privateEndpointIpConfig",
							"properties": map[string]any{
								"primary":			true,
								"privateIPAddress":		"10.0.0.4",
								"privateIPAddressVersion":	"IPv4",
								"privateIPAllocationMethod":	"Dynamic",
								"subnet": map[string]any{
									"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
								},
							},
						},
					},
				},
			},
			map[string]any{
				"name":		"[parameters('privateEndpointName')]",
				"type":		"Microsoft.Network/privateEndpoints",
				"apiVersion":	"2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location":	"[parameters('location')]",
				"properties": map[string]any{
					"customDnsConfigs":			[]any{},
					"manualPrivateLinkServiceConnections":	[]any{},
					"privateLinkServiceConnections": []any{
						map[string]any{
							"name":	"[parameters('privateEndpointName')]",
							"properties": map[string]any{
								"groupIds": []any{
									"standard",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":		"Auto-Approved",
									"actionsRequired":	"None",
									"status":		"Approved",
								},
								"privateLinkServiceId":	"[parameters('attestationId')]",
							},
						},
					},
					"subnet": map[string]any{
						"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
					},
				},
			},
			map[string]any{
				"name":		"[concat(parameters('virtualNetworksName'), '/default')]",
				"type":		"Microsoft.Network/virtualNetworks/subnets",
				"apiVersion":	"2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
				},
				"properties": map[string]any{
					"addressPrefix":			"10.0.0.0/24",
					"delegations":				[]any{},
					"privateEndpointNetworkPolicies":	"Disabled",
					"privateLinkServiceNetworkPolicies":	"Enabled",
				},
			},
		},
		"variables":	map[string]any{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:	template,
			Mode:		to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_List
	fmt.Println("Call operation: PrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armattestation.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.providerName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionId = *nextResult.Value[0].ID

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Create
	fmt.Println("Call operation: PrivateEndpointConnections_Create")
	_, err = privateEndpointConnectionsClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.providerName, privateEndpointConnectionName, armattestation.PrivateEndpointConnection{
		ID:	to.Ptr(privateEndpointConnectionId),
		Properties: &armattestation.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armattestation.PrivateLinkServiceConnectionState{
				Description:	to.Ptr("rejection connection"),
				Status:		to.Ptr(armattestation.PrivateEndpointServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.providerName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Attestation/operations
func (testsuite *AttestationTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armattestation.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = operationsClient.List(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *AttestationTestSuite) Cleanup() {
	var err error
	// From step AttestationProviders_Delete
	fmt.Println("Call operation: AttestationProviders_Delete")
	providersClient, err := armattestation.NewProvidersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = providersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.providerName, nil)
	testsuite.Require().NoError(err)
}

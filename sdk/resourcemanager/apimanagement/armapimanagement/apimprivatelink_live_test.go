// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armapimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type ApimprivatelinkTestSuite struct {
	suite.Suite

	ctx                           context.Context
	cred                          azcore.TokenCredential
	options                       *arm.ClientOptions
	apimId                        string
	privateEndpointConnectionName string
	serviceName                   string
	virtualNetworksName           string
	location                      string
	resourceGroupName             string
	subscriptionId                string
}

func (testsuite *ApimprivatelinkTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.privateEndpointConnectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "apimprivateendpoint", 25, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceprivate", 20, false)
	testsuite.virtualNetworksName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "apimvnet", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimprivatelinkTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimprivatelinkTestSuite(t *testing.T) {
	suite.Run(t, new(ApimprivatelinkTestSuite))
}

func (testsuite *ApimprivatelinkTestSuite) Prepare() {
	var err error
	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceClientCreateOrUpdateResponsePoller, err := serviceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceResource{
		Tags: map[string]*string{
			"Name": to.Ptr("Contoso"),
			"Test": to.Ptr("User"),
		},
		Location: to.Ptr(testsuite.location),
		Properties: &armapimanagement.ServiceProperties{
			PublisherEmail: to.Ptr("foo@contoso.com"),
			PublisherName:  to.Ptr("foo"),
		},
		SKU: &armapimanagement.ServiceSKUProperties{
			Name:     to.Ptr(armapimanagement.SKUTypeStandard),
			Capacity: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var serviceClientCreateOrUpdateResponse *armapimanagement.ServiceClientCreateOrUpdateResponse
	serviceClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.apimId = *serviceClientCreateOrUpdateResponse.ID

	// From step PrivateEndpoint_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]interface{}{
			"apimId": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(apimId)",
			},
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(location)",
			},
			"privateEndpointConnectionName": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(privateEndpointConnectionName)",
			},
			"virtualNetworksName": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(virtualNetworksName)",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2020-11-01",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"addressSpace": map[string]interface{}{
						"addressPrefixes": []interface{}{
							"10.0.0.0/16",
						},
					},
					"enableDdosProtection": false,
					"subnets": []interface{}{
						map[string]interface{}{
							"name": "default",
							"properties": map[string]interface{}{
								"addressPrefix":                     "10.0.0.0/24",
								"delegations":                       []interface{}{},
								"privateEndpointNetworkPolicies":    "Disabled",
								"privateLinkServiceNetworkPolicies": "Enabled",
							},
						},
					},
					"virtualNetworkPeerings": []interface{}{},
				},
			},
			map[string]interface{}{
				"name":       "[concat(parameters('privateEndpointConnectionName'), '-nic')]",
				"type":       "Microsoft.Network/networkInterfaces",
				"apiVersion": "2020-11-01",
				"dependsOn": []interface{}{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]interface{}{
					"dnsSettings": map[string]interface{}{
						"dnsServers": []interface{}{},
					},
					"enableIPForwarding": false,
					"ipConfigurations": []interface{}{
						map[string]interface{}{
							"name": "privateEndpointIpConfig",
							"properties": map[string]interface{}{
								"primary":                   true,
								"privateIPAddress":          "10.0.0.4",
								"privateIPAddressVersion":   "IPv4",
								"privateIPAllocationMethod": "Dynamic",
								"subnet": map[string]interface{}{
									"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
								},
							},
						},
					},
				},
			},
			map[string]interface{}{
				"name":       "[parameters('privateEndpointConnectionName')]",
				"type":       "Microsoft.Network/privateEndpoints",
				"apiVersion": "2020-11-01",
				"dependsOn": []interface{}{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]interface{}{
					"customDnsConfigs":                    []interface{}{},
					"manualPrivateLinkServiceConnections": []interface{}{},
					"privateLinkServiceConnections": []interface{}{
						map[string]interface{}{
							"name": "[parameters('privateEndpointConnectionName')]",
							"properties": map[string]interface{}{
								"groupIds": []interface{}{
									"Gateway",
								},
								"privateLinkServiceConnectionState": map[string]interface{}{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('apimId')]",
							},
						},
					},
					"subnet": map[string]interface{}{
						"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
					},
				},
			},
			map[string]interface{}{
				"name":       "[concat(parameters('virtualNetworksName'), '/default')]",
				"type":       "Microsoft.Network/virtualNetworks/subnets",
				"apiVersion": "2020-11-01",
				"dependsOn": []interface{}{
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
				},
				"properties": map[string]interface{}{
					"addressPrefix":                     "10.0.0.0/24",
					"delegations":                       []interface{}{},
					"privateEndpointNetworkPolicies":    "Disabled",
					"privateLinkServiceNetworkPolicies": "Enabled",
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	params := map[string]interface{}{
		"apimId":                        map[string]interface{}{"value": testsuite.apimId},
		"location":                      map[string]interface{}{"value": testsuite.location},
		"privateEndpointConnectionName": map[string]interface{}{"value": testsuite.privateEndpointConnectionName},
		"virtualNetworksName":           map[string]interface{}{"value": testsuite.virtualNetworksName},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "PrivateEndpoint_Create", &deployment)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/privateEndpointConnections
func (testsuite *ApimprivatelinkTestSuite) TestPrivateendpointconnection() {
	var err error
	// From step PrivateEndpointConnection_CreateOrUpdate
	fmt.Println("Call operation: PrivateEndpointConnection_CreateOrUpdate")
	privateEndpointConnectionClient, err := armapimanagement.NewPrivateEndpointConnectionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionClientCreateOrUpdateResponsePoller, err := privateEndpointConnectionClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.privateEndpointConnectionName, armapimanagement.PrivateEndpointConnectionRequest{
		ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.ApiManagement/service/" + testsuite.serviceName + "/privateEndpointConnections/" + testsuite.privateEndpointConnectionName),
		Properties: &armapimanagement.PrivateEndpointConnectionRequestProperties{
			PrivateLinkServiceConnectionState: &armapimanagement.PrivateLinkServiceConnectionState{
				Description: to.Ptr("The Private Endpoint Connection is approved."),
				Status:      to.Ptr(armapimanagement.PrivateEndpointServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnection_ListByService
	fmt.Println("Call operation: PrivateEndpointConnection_ListByService")
	privateEndpointConnectionClientNewListByServicePager := privateEndpointConnectionClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for privateEndpointConnectionClientNewListByServicePager.More() {
		_, err := privateEndpointConnectionClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateEndpointConnection_ListPrivateLinkResources
	fmt.Println("Call operation: PrivateEndpointConnection_ListPrivateLinkResources")
	_, err = privateEndpointConnectionClient.ListPrivateLinkResources(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnection_GetByName
	fmt.Println("Call operation: PrivateEndpointConnection_GetByName")
	_, err = privateEndpointConnectionClient.GetByName(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnection_Delete
	fmt.Println("Call operation: PrivateEndpointConnection_Delete")
	privateEndpointConnectionClientDeleteResponsePoller, err := privateEndpointConnectionClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

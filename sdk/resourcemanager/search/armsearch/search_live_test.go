// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armsearch_test

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/search/armsearch"
	"github.com/stretchr/testify/suite"
)

type SearchTestSuite struct {
	suite.Suite

	ctx                           context.Context
	cred                          azcore.TokenCredential
	options                       *arm.ClientOptions
	searchServiceName             string
	serviceId                     string
	sharedPrivateLinkResourceName string
	storageAccountName            string
	location                      string
	resourceGroupName             string
	subscriptionId                string
}

func (testsuite *SearchTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.searchServiceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "searchse", 14, true)
	testsuite.sharedPrivateLinkResourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sharedpr", 14, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "searchstorageac", 21, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SearchTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestSearchTestSuite(t *testing.T) {
	// The resource type 'checkNameAvailability' could not be found in the namespace 'Microsoft.Search' for api version '2025-05-01'
	suite.Run(t, new(SearchTestSuite))
}

func (testsuite *SearchTestSuite) Prepare() {
	var err error
	// From step Services_CheckNameAvailability
	fmt.Println("Call operation: Services_CheckNameAvailability")
	servicesClient, err := armsearch.NewServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = servicesClient.CheckNameAvailability(testsuite.ctx, armsearch.CheckNameAvailabilityInput{
		Name: to.Ptr(testsuite.searchServiceName),
		Type: to.Ptr("Microsoft.Search/searchServices"),
	}, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)

	// From step Services_CreateOrUpdate
	fmt.Println("Call operation: Services_CreateOrUpdate")
	servicesClientCreateOrUpdateResponsePoller, err := servicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, armsearch.Service{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"app-name": to.Ptr("My e-commerce app"),
		},
		Properties: &armsearch.ServiceProperties{
			HostingMode:    to.Ptr(armsearch.HostingModeDefault),
			PartitionCount: to.Ptr[int32](1),
			ReplicaCount:   to.Ptr[int32](3),
		},
		SKU: &armsearch.SKU{
			Name: to.Ptr(armsearch.SKUNameStandard),
		},
	}, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
	var servicesClientCreateOrUpdateResponse *armsearch.ServicesClientCreateOrUpdateResponse
	servicesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, servicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.serviceId = *servicesClientCreateOrUpdateResponse.ID
}

// Microsoft.Search/searchServices/{searchServiceName}
func (testsuite *SearchTestSuite) TestServices() {
	var err error
	// From step Services_ListBySubscription
	fmt.Println("Call operation: Services_ListBySubscription")
	servicesClient, err := armsearch.NewServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	servicesClientNewListBySubscriptionPager := servicesClient.NewListBySubscriptionPager(&armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	for servicesClientNewListBySubscriptionPager.More() {
		_, err := servicesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Services_ListByResourceGroup
	fmt.Println("Call operation: Services_ListByResourceGroup")
	servicesClientNewListByResourceGroupPager := servicesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	for servicesClientNewListByResourceGroupPager.More() {
		_, err := servicesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Services_Get
	fmt.Println("Call operation: Services_Get")
	_, err = servicesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)

	// From step Services_Update
	fmt.Println("Call operation: Services_Update")
	_, err = servicesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, armsearch.ServiceUpdate{
		Properties: &armsearch.ServiceProperties{
			ReplicaCount: to.Ptr[int32](2),
		},
		Tags: map[string]*string{
			"app-name": to.Ptr("My e-commerce app"),
			"new-tag":  to.Ptr("Adding a new tag"),
		},
	}, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Search/searchServices/{searchServiceName}/sharedPrivateLinkResources/{sharedPrivateLinkResourceName}
func (testsuite *SearchTestSuite) TestSharedPrivateLinkResources() {
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

	// From step SharedPrivateLinkResources_CreateOrUpdate
	fmt.Println("Call operation: SharedPrivateLinkResources_CreateOrUpdate")
	sharedPrivateLinkResourcesClient, err := armsearch.NewSharedPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sharedPrivateLinkResourcesClientCreateOrUpdateResponsePoller, err := sharedPrivateLinkResourcesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, testsuite.sharedPrivateLinkResourceName, armsearch.SharedPrivateLinkResource{
		Properties: &armsearch.SharedPrivateLinkResourceProperties{
			GroupID:               to.Ptr("blob"),
			PrivateLinkResourceID: to.Ptr(storageAccountId),
			RequestMessage:        to.Ptr("please approve"),
		},
	}, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sharedPrivateLinkResourcesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step SharedPrivateLinkResources_ListByService
	fmt.Println("Call operation: SharedPrivateLinkResources_ListByService")
	sharedPrivateLinkResourcesClientNewListByServicePager := sharedPrivateLinkResourcesClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.searchServiceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	for sharedPrivateLinkResourcesClientNewListByServicePager.More() {
		_, err := sharedPrivateLinkResourcesClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SharedPrivateLinkResources_Get
	fmt.Println("Call operation: SharedPrivateLinkResources_Get")
	_, err = sharedPrivateLinkResourcesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, testsuite.sharedPrivateLinkResourceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)

	// From step SharedPrivateLinkResources_Delete
	fmt.Println("Call operation: SharedPrivateLinkResources_Delete")
	sharedPrivateLinkResourcesClientDeleteResponsePoller, err := sharedPrivateLinkResourcesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, testsuite.sharedPrivateLinkResourceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, sharedPrivateLinkResourcesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Search/searchServices/{searchServiceName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *SearchTestSuite) TestPrivateEndpointConnections() {
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
				"defaultValue": "endpointsearch-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "endpointsearch",
			},
			"serviceId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.serviceId,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "searchvnet",
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
									"searchService",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('serviceId')]",
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

	// From step PrivateEndpointConnections_ListByService
	fmt.Println("Call operation: PrivateEndpointConnections_ListByService")
	privateEndpointConnectionsClient, err := armsearch.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListByServicePager := privateEndpointConnectionsClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.searchServiceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	for privateEndpointConnectionsClientNewListByServicePager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Update
	fmt.Println("Call operation: PrivateEndpointConnections_Update")
	_, err = privateEndpointConnectionsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, privateEndpointConnectionName, armsearch.PrivateEndpointConnection{
		Properties: &armsearch.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armsearch.PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState{
				Description: to.Ptr("Rejected for some reason"),
				Status:      to.Ptr(armsearch.PrivateLinkServiceConnectionStatusRejected),
			},
		},
	}, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, privateEndpointConnectionName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_ListSupported
	fmt.Println("Call operation: PrivateLinkResources_ListSupported")
	privateLinkResourcesClient, err := armsearch.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListSupportedPager := privateLinkResourcesClient.NewListSupportedPager(testsuite.resourceGroupName, testsuite.searchServiceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	for privateLinkResourcesClientNewListSupportedPager.More() {
		_, err := privateLinkResourcesClientNewListSupportedPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	_, err = privateEndpointConnectionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, privateEndpointConnectionName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Search/operations
func (testsuite *SearchTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armsearch.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Search/searchServices/{searchServiceName}/createQueryKey/{name}
func (testsuite *SearchTestSuite) TestQueryKeys() {
	name := "myquerykey"
	var queryKey string
	var err error
	// From step QueryKeys_Create
	fmt.Println("Call operation: QueryKeys_Create")
	queryKeysClient, err := armsearch.NewQueryKeysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	queryKeysClientCreateResponse, err := queryKeysClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, name, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
	queryKey = *queryKeysClientCreateResponse.Key

	// From step QueryKeys_ListBySearchService
	fmt.Println("Call operation: QueryKeys_ListBySearchService")
	queryKeysClientNewListBySearchServicePager := queryKeysClient.NewListBySearchServicePager(testsuite.resourceGroupName, testsuite.searchServiceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	for queryKeysClientNewListBySearchServicePager.More() {
		_, err := queryKeysClientNewListBySearchServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step QueryKeys_Delete
	fmt.Println("Call operation: QueryKeys_Delete")
	_, err = queryKeysClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, queryKey, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Search/searchServices/{searchServiceName}/regenerateAdminKey/{keyKind}
func (testsuite *SearchTestSuite) TestAdminKeys() {
	var err error
	// From step AdminKeys_Get
	fmt.Println("Call operation: AdminKeys_Get")
	adminKeysClient, err := armsearch.NewAdminKeysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = adminKeysClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)

	// From step AdminKeys_Regenerate
	fmt.Println("Call operation: AdminKeys_Regenerate")
	_, err = adminKeysClient.Regenerate(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, armsearch.AdminKeyKindPrimary, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *SearchTestSuite) Cleanup() {
	var err error
	// From step Services_Delete
	fmt.Println("Call operation: Services_Delete")
	servicesClient, err := armsearch.NewServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = servicesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.searchServiceName, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	testsuite.Require().NoError(err)
}

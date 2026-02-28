// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armbatch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/batch/armbatch/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type BatchManagementTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	accountName        string
	applicationName    string
	batchAccountId     string
	poolName           string
	storageAccountId   string
	storageAccountName string
	versionName        string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *BatchManagementTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.applicationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "applicat", 14, false)
	testsuite.poolName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "poolname", 14, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storagea", 14, true)
	testsuite.versionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "versionn", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *BatchManagementTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestBatchManagementTestSuite(t *testing.T) {
	suite.Run(t, new(BatchManagementTestSuite))
}

func (testsuite *BatchManagementTestSuite) Prepare() {
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
	testsuite.storageAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountId"].(map[string]interface{})["value"].(string)

	// From step BatchAccount_Create
	fmt.Println("Call operation: BatchAccount_Create")
	accountClient, err := armbatch.NewAccountClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	accountClientCreateResponsePoller, err := accountClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armbatch.AccountCreateParameters{
		Location: to.Ptr(testsuite.location),
		Properties: &armbatch.AccountCreateProperties{
			AutoStorage: &armbatch.AutoStorageBaseProperties{
				StorageAccountID: to.Ptr(testsuite.storageAccountId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var accountClientCreateResponse *armbatch.AccountClientCreateResponse
	accountClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, accountClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.batchAccountId = *accountClientCreateResponse.ID

	// From step Application_Create
	fmt.Println("Call operation: Application_Create")
	applicationClient, err := armbatch.NewApplicationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = applicationClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.applicationName, armbatch.Application{
		Properties: &armbatch.ApplicationProperties{
			AllowUpdates: to.Ptr(false),
			DisplayName:  to.Ptr("myAppName"),
		},
	},
		&armbatch.ApplicationClientCreateOptions{},
	)
	testsuite.Require().NoError(err)
}

// Microsoft.Batch/locations/{locationName}
func (testsuite *BatchManagementTestSuite) TestLocation() {
	locationName := "westus"
	var err error
	// From step Location_CheckNameAvailability
	fmt.Println("Call operation: Location_CheckNameAvailability")
	locationClient, err := armbatch.NewLocationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = locationClient.CheckNameAvailability(testsuite.ctx, locationName, armbatch.CheckNameAvailabilityParameters{
		Name: to.Ptr("newaccountname"),
		Type: to.Ptr("Microsoft.Batch/batchAccounts"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Location_GetQuotas
	fmt.Println("Call operation: Location_GetQuotas")
	_, err = locationClient.GetQuotas(testsuite.ctx, locationName, nil)
	testsuite.Require().NoError(err)

	// From step Location_ListSupportedVirtualMachineSkus
	fmt.Println("Call operation: Location_ListSupportedVirtualMachineSkus")
	locationClientNewListSupportedVirtualMachineSKUsPager := locationClient.NewListSupportedVirtualMachineSKUsPager(locationName, &armbatch.LocationClientListSupportedVirtualMachineSKUsOptions{Maxresults: nil,
		Filter: nil,
	})
	for locationClientNewListSupportedVirtualMachineSKUsPager.More() {
		_, err := locationClientNewListSupportedVirtualMachineSKUsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Batch/batchAccounts/{accountName}
func (testsuite *BatchManagementTestSuite) TestBatchAccount() {
	var err error
	// From step BatchAccount_List
	fmt.Println("Call operation: BatchAccount_List")
	accountClient, err := armbatch.NewAccountClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	accountClientNewListPager := accountClient.NewListPager(nil)
	for accountClientNewListPager.More() {
		_, err := accountClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BatchAccount_Get
	fmt.Println("Call operation: BatchAccount_Get")
	_, err = accountClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step BatchAccount_ListByResourceGroup
	fmt.Println("Call operation: BatchAccount_ListByResourceGroup")
	accountClientNewListByResourceGroupPager := accountClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for accountClientNewListByResourceGroupPager.More() {
		_, err := accountClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BatchAccount_ListOutboundNetworkDependenciesEndpoints
	fmt.Println("Call operation: BatchAccount_ListOutboundNetworkDependenciesEndpoints")
	accountClientNewListOutboundNetworkDependenciesEndpointsPager := accountClient.NewListOutboundNetworkDependenciesEndpointsPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for accountClientNewListOutboundNetworkDependenciesEndpointsPager.More() {
		_, err := accountClientNewListOutboundNetworkDependenciesEndpointsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BatchAccount_ListDetectors
	fmt.Println("Call operation: BatchAccount_ListDetectors")
	accountClientNewListDetectorsPager := accountClient.NewListDetectorsPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for accountClientNewListDetectorsPager.More() {
		_, err := accountClientNewListDetectorsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BatchAccount_GetDetector
	fmt.Println("Call operation: BatchAccount_GetDetector")
	_, err = accountClient.GetDetector(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, "poolsAndNodes", nil)
	testsuite.Require().NoError(err)

	// From step BatchAccount_Update
	fmt.Println("Call operation: BatchAccount_Update")
	_, err = accountClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armbatch.AccountUpdateParameters{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step BatchAccount_SynchronizeAutoStorageKeys
	fmt.Println("Call operation: BatchAccount_SynchronizeAutoStorageKeys")
	_, err = accountClient.SynchronizeAutoStorageKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step BatchAccount_RegenerateKey
	fmt.Println("Call operation: BatchAccount_RegenerateKey")
	_, err = accountClient.RegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armbatch.AccountRegenerateKeyParameters{
		KeyName: to.Ptr(armbatch.AccountKeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step BatchAccount_GetKeys
	fmt.Println("Call operation: BatchAccount_GetKeys")
	_, err = accountClient.GetKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Batch/batchAccounts/{accountName}/pools/{poolName}
func (testsuite *BatchManagementTestSuite) TestPool() {
	var err error
	// From step Pool_Create
	fmt.Println("Call operation: Pool_Create")
	poolClient, err := armbatch.NewPoolClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = poolClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.poolName, armbatch.Pool{
		Properties: &armbatch.PoolProperties{
			DeploymentConfiguration: &armbatch.DeploymentConfiguration{
				VirtualMachineConfiguration: &armbatch.VirtualMachineConfiguration{},
			},
			ScaleSettings: &armbatch.ScaleSettings{
				FixedScale: &armbatch.FixedScaleSettings{
					TargetDedicatedNodes: to.Ptr[int32](3),
				},
			},
			VMSize: to.Ptr("STANDARD_D4"),
		},
	}, &armbatch.PoolClientCreateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)

	// From step Pool_ListByBatchAccount
	fmt.Println("Call operation: Pool_ListByBatchAccount")
	poolClientNewListByBatchAccountPager := poolClient.NewListByBatchAccountPager(testsuite.resourceGroupName, testsuite.accountName, &armbatch.PoolClientListByBatchAccountOptions{Maxresults: nil,
		Select: nil,
		Filter: nil,
	})
	for poolClientNewListByBatchAccountPager.More() {
		_, err := poolClientNewListByBatchAccountPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Pool_Get
	fmt.Println("Call operation: Pool_Get")
	_, err = poolClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.poolName, nil)
	testsuite.Require().NoError(err)

	// From step Pool_Update
	fmt.Println("Call operation: Pool_Update")
	_, err = poolClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.poolName, armbatch.Pool{
		Properties: &armbatch.PoolProperties{
			ScaleSettings: &armbatch.ScaleSettings{
				FixedScale: &armbatch.FixedScaleSettings{
					NodeDeallocationOption: to.Ptr(armbatch.ComputeNodeDeallocationOptionTaskCompletion),
					ResizeTimeout:          to.Ptr("PT8M"),
					TargetDedicatedNodes:   to.Ptr[int32](5),
					TargetLowPriorityNodes: to.Ptr[int32](0),
				},
			},
		},
	}, &armbatch.PoolClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Pool_DisableAutoScale
	fmt.Println("Call operation: Pool_DisableAutoScale")
	_, err = poolClient.DisableAutoScale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.poolName, nil)
	testsuite.Require().NoError(err)

	// From step Pool_StopResize
	fmt.Println("Call operation: Pool_StopResize")
	_, err = poolClient.StopResize(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.poolName, nil)
	testsuite.Require().NoError(err)

	// From step Pool_Delete
	fmt.Println("Call operation: Pool_Delete")
	poolClientDeleteResponsePoller, err := poolClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.poolName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, poolClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Batch/batchAccounts/{accountName}/applications/{applicationName}
func (testsuite *BatchManagementTestSuite) TestApplication() {
	var err error
	// From step Application_List
	fmt.Println("Call operation: Application_List")
	applicationClient, err := armbatch.NewApplicationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	applicationClientNewListPager := applicationClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armbatch.ApplicationClientListOptions{Maxresults: nil})
	for applicationClientNewListPager.More() {
		_, err := applicationClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Application_Get
	fmt.Println("Call operation: Application_Get")
	_, err = applicationClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.applicationName, nil)
	testsuite.Require().NoError(err)

	// From step Application_Update
	fmt.Println("Call operation: Application_Update")
	_, err = applicationClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.applicationName, armbatch.Application{
		Properties: &armbatch.ApplicationProperties{
			AllowUpdates: to.Ptr(true),
			DisplayName:  to.Ptr("myAppName"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Batch/batchAccounts/{accountName}/applications/{applicationName}/versions/{versionName}
func (testsuite *BatchManagementTestSuite) TestApplicationPackage() {
	var err error
	// From step ApplicationPackage_Create
	fmt.Println("Call operation: ApplicationPackage_Create")
	applicationPackageClient, err := armbatch.NewApplicationPackageClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = applicationPackageClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.applicationName, testsuite.versionName, armbatch.ApplicationPackage{}, &armbatch.ApplicationPackageClientCreateOptions{})
	testsuite.Require().NoError(err)

	// From step ApplicationPackage_List
	fmt.Println("Call operation: ApplicationPackage_List")
	applicationPackageClientNewListPager := applicationPackageClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.applicationName, &armbatch.ApplicationPackageClientListOptions{Maxresults: nil})
	for applicationPackageClientNewListPager.More() {
		_, err := applicationPackageClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApplicationPackage_Get
	fmt.Println("Call operation: ApplicationPackage_Get")
	_, err = applicationPackageClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.applicationName, testsuite.versionName, nil)
	testsuite.Require().NoError(err)

	// From step ApplicationPackage_Delete
	fmt.Println("Call operation: ApplicationPackage_Delete")
	_, err = applicationPackageClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.applicationName, testsuite.versionName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Batch/operations
func (testsuite *BatchManagementTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armbatch.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Batch/batchAccounts/{accountName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *BatchManagementTestSuite) TestPrivateEndpointConnection() {
	var privateEndpointConnectionName string
	var privateLinkResourceName string
	var err error
	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]any{
			"batchAccountId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.batchAccountId,
			},
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "epbatch-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "epbatch",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "batchvnet",
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
									"batchAccount",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('batchAccountId')]",
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

	// From step PrivateEndpointConnection_ListByBatchAccount
	fmt.Println("Call operation: PrivateEndpointConnection_ListByBatchAccount")
	privateEndpointConnectionClient, err := armbatch.NewPrivateEndpointConnectionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionClientNewListByBatchAccountPager := privateEndpointConnectionClient.NewListByBatchAccountPager(testsuite.resourceGroupName, testsuite.accountName, &armbatch.PrivateEndpointConnectionClientListByBatchAccountOptions{Maxresults: nil})
	for privateEndpointConnectionClientNewListByBatchAccountPager.More() {
		nextResult, err := privateEndpointConnectionClientNewListByBatchAccountPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnection_Update
	fmt.Println("Call operation: PrivateEndpointConnection_Update")
	privateEndpointConnectionClientUpdateResponsePoller, err := privateEndpointConnectionClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, privateEndpointConnectionName, armbatch.PrivateEndpointConnection{
		Properties: &armbatch.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armbatch.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Approved by xyz.abc@company.com"),
				Status:      to.Ptr(armbatch.PrivateLinkServiceConnectionStatusRejected),
			},
		},
	}, &armbatch.PrivateEndpointConnectionClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnection_Get
	fmt.Println("Call operation: PrivateEndpointConnection_Get")
	_, err = privateEndpointConnectionClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResource_ListByBatchAccount
	fmt.Println("Call operation: PrivateLinkResource_ListByBatchAccount")
	privateLinkResourceClient, err := armbatch.NewPrivateLinkResourceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourceClientNewListByBatchAccountPager := privateLinkResourceClient.NewListByBatchAccountPager(testsuite.resourceGroupName, testsuite.accountName, &armbatch.PrivateLinkResourceClientListByBatchAccountOptions{Maxresults: nil})
	for privateLinkResourceClientNewListByBatchAccountPager.More() {
		nextResult, err := privateLinkResourceClientNewListByBatchAccountPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateLinkResourceName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateLinkResource_Get
	fmt.Println("Call operation: PrivateLinkResource_Get")
	_, err = privateLinkResourceClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, privateLinkResourceName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnection_Delete
	fmt.Println("Call operation: PrivateEndpointConnection_Delete")
	privateEndpointConnectionClientDeleteResponsePoller, err := privateEndpointConnectionClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *BatchManagementTestSuite) Cleanup() {
	var err error
	// From step Application_Delete
	fmt.Println("Call operation: Application_Delete")
	applicationClient, err := armbatch.NewApplicationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = applicationClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.applicationName, nil)
	testsuite.Require().NoError(err)

	// From step BatchAccount_Delete
	fmt.Println("Call operation: BatchAccount_Delete")
	accountClient, err := armbatch.NewAccountClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	accountClientDeleteResponsePoller, err := accountClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, accountClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

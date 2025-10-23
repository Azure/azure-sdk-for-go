//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v3"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	accountName         string
	encryptionScopeName string
	storageAccountId    string
	location            string
	resourceGroupName   string
	subscriptionId      string
}

func (testsuite *StorageTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName = "accountnam"
	testsuite.encryptionScopeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "encryption", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *StorageTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (testsuite *StorageTestSuite) Prepare() {
	var err error
	// From step StorageAccounts_CheckNameAvailability
	fmt.Println("Call operation: StorageAccounts_CheckNameAvailability")
	accountsClient, err := armstorage.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = accountsClient.CheckNameAvailability(testsuite.ctx, armstorage.AccountCheckNameAvailabilityParameters{
		Name: to.Ptr("sto3363"),
		Type: to.Ptr("Microsoft.Storage/storageAccounts"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step StorageAccount_Create
	fmt.Println("Call operation: StorageAccounts_Create")
	accountsClientCreateResponsePoller, err := accountsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.AccountCreateParameters{
		Kind:     to.Ptr(armstorage.KindStorageV2),
		Location: to.Ptr(testsuite.location),
		Properties: &armstorage.AccountPropertiesCreateParameters{
			AllowBlobPublicAccess:        to.Ptr(false),
			AllowSharedKeyAccess:         to.Ptr(true),
			DefaultToOAuthAuthentication: to.Ptr(false),
			Encryption: &armstorage.Encryption{
				KeySource:                       to.Ptr(armstorage.KeySourceMicrosoftStorage),
				RequireInfrastructureEncryption: to.Ptr(false),
				Services: &armstorage.EncryptionServices{
					Blob: &armstorage.EncryptionService{
						Enabled: to.Ptr(true),
						KeyType: to.Ptr(armstorage.KeyTypeAccount),
					},
					File: &armstorage.EncryptionService{
						Enabled: to.Ptr(true),
						KeyType: to.Ptr(armstorage.KeyTypeAccount),
					},
				},
			},
			IsHnsEnabled:  to.Ptr(true),
			IsSftpEnabled: to.Ptr(true),
			KeyPolicy: &armstorage.KeyPolicy{
				KeyExpirationPeriodInDays: to.Ptr[int32](20),
			},
			MinimumTLSVersion: to.Ptr(armstorage.MinimumTLSVersionTLS12),
			RoutingPreference: &armstorage.RoutingPreference{
				PublishInternetEndpoints:  to.Ptr(true),
				PublishMicrosoftEndpoints: to.Ptr(true),
				RoutingChoice:             to.Ptr(armstorage.RoutingChoiceMicrosoftRouting),
			},
			SasPolicy: &armstorage.SasPolicy{
				ExpirationAction:    to.Ptr(armstorage.ExpirationActionLog),
				SasExpirationPeriod: to.Ptr("1.15:59:59"),
			},
		},
		SKU: &armstorage.SKU{
			Name: to.Ptr(armstorage.SKUNameStandardGRS),
		},
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var accountsClientCreateResponse *armstorage.AccountsClientCreateResponse
	accountsClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, accountsClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.storageAccountId = *accountsClientCreateResponse.ID
}

// Microsoft.Storage/storageAccounts/{accountName}
func (testsuite *StorageTestSuite) TestStorageAccounts() {
	var err error
	// From step StorageAccounts_List
	fmt.Println("Call operation: StorageAccounts_List")
	accountsClient, err := armstorage.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	accountsClientNewListPager := accountsClient.NewListPager(nil)
	for accountsClientNewListPager.More() {
		_, err := accountsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StorageAccounts_ListByResourceGroup
	fmt.Println("Call operation: StorageAccounts_ListByResourceGroup")
	accountsClientNewListByResourceGroupPager := accountsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for accountsClientNewListByResourceGroupPager.More() {
		_, err := accountsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StorageAccounts_GetProperties
	fmt.Println("Call operation: StorageAccounts_GetProperties")
	_, err = accountsClient.GetProperties(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, &armstorage.AccountsClientGetPropertiesOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step StorageAccounts_Update
	fmt.Println("Call operation: StorageAccounts_Update")
	_, err = accountsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.AccountUpdateParameters{
		Tags: map[string]*string{
			"storageKey": to.Ptr("storageValue"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step StorageAccounts_RevokeUserDelegationKeys
	fmt.Println("Call operation: StorageAccounts_RevokeUserDelegationKeys")
	_, err = accountsClient.RevokeUserDelegationKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step StorageAccounts_ListServiceSAS
	fmt.Println("Call operation: StorageAccounts_ListServiceSAS")
	_, err = accountsClient.ListServiceSAS(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.ServiceSasParameters{
		CanonicalizedResource:  to.Ptr("/blob/" + testsuite.accountName + "/music"),
		SharedAccessExpiryTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-05-24T11:32:48.8457197Z"); return t }()),
		Permissions:            to.Ptr(armstorage.PermissionsL),
		Resource:               to.Ptr(armstorage.SignedResourceC),
	}, nil)
	testsuite.Require().NoError(err)

	// From step StorageAccounts_ListKeys
	fmt.Println("Call operation: StorageAccounts_ListKeys")
	_, err = accountsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, &armstorage.AccountsClientListKeysOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step StorageAccounts_RegenerateKey
	fmt.Println("Call operation: StorageAccounts_RegenerateKey")
	_, err = accountsClient.RegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.AccountRegenerateKeyParameters{
		KeyName: to.Ptr("key2"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step StorageAccounts_ListAccountSAS
	fmt.Println("Call operation: StorageAccounts_ListAccountSAS")
	_, err = accountsClient.ListAccountSAS(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.AccountSasParameters{
		KeyToSign:              to.Ptr("key1"),
		SharedAccessExpiryTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-05-24T11:42:03.1567373Z"); return t }()),
		Permissions:            to.Ptr(armstorage.PermissionsR),
		Protocols:              to.Ptr(armstorage.HTTPProtocolHTTPSHTTP),
		ResourceTypes:          to.Ptr(armstorage.SignedResourceTypesS),
		Services:               to.Ptr(armstorage.ServicesB),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/operations
func (testsuite *StorageTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armstorage.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Storage/skus
func (testsuite *StorageTestSuite) TestSkus() {
	var err error
	// From step Skus_List
	fmt.Println("Call operation: SKUs_List")
	sKUsClient, err := armstorage.NewSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sKUsClientNewListPager := sKUsClient.NewListPager(nil)
	for sKUsClientNewListPager.More() {
		_, err := sKUsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Storage/locations/{location}/usages
func (testsuite *StorageTestSuite) TestUsages() {
	var err error
	// From step Usages_ListByLocation
	fmt.Println("Call operation: Usages_ListByLocation")
	usagesClient, err := armstorage.NewUsagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	usagesClientNewListByLocationPager := usagesClient.NewListByLocationPager(testsuite.location, nil)
	for usagesClientNewListByLocationPager.More() {
		_, err := usagesClientNewListByLocationPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Storage/storageAccounts/{accountName}/managementPolicies/{managementPolicyName}
func (testsuite *StorageTestSuite) TestManagementPolicies() {
	var err error
	// From step ManagementPolicies_CreateOrUpdate
	fmt.Println("Call operation: ManagementPolicies_CreateOrUpdate")
	managementPoliciesClient, err := armstorage.NewManagementPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managementPoliciesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.ManagementPolicyNameDefault, armstorage.ManagementPolicy{
		Properties: &armstorage.ManagementPolicyProperties{
			Policy: &armstorage.ManagementPolicySchema{
				Rules: []*armstorage.ManagementPolicyRule{
					{
						Name: to.Ptr("olcmtest1"),
						Type: to.Ptr(armstorage.RuleTypeLifecycle),
						Definition: &armstorage.ManagementPolicyDefinition{
							Actions: &armstorage.ManagementPolicyAction{
								BaseBlob: &armstorage.ManagementPolicyBaseBlob{
									Delete: &armstorage.DateAfterModification{
										DaysAfterModificationGreaterThan: to.Ptr[float32](1000),
									},
									TierToArchive: &armstorage.DateAfterModification{
										DaysAfterModificationGreaterThan: to.Ptr[float32](90),
									},
									TierToCool: &armstorage.DateAfterModification{
										DaysAfterModificationGreaterThan: to.Ptr[float32](30),
									},
								},
								Snapshot: &armstorage.ManagementPolicySnapShot{
									Delete: &armstorage.DateAfterCreation{
										DaysAfterCreationGreaterThan: to.Ptr[float32](30),
									},
								},
							},
							Filters: &armstorage.ManagementPolicyFilter{
								BlobTypes: []*string{
									to.Ptr("blockBlob")},
								PrefixMatch: []*string{
									to.Ptr("olcmtestcontainer1")},
							},
						},
						Enabled: to.Ptr(true),
					},
					{
						Name: to.Ptr("olcmtest2"),
						Type: to.Ptr(armstorage.RuleTypeLifecycle),
						Definition: &armstorage.ManagementPolicyDefinition{
							Actions: &armstorage.ManagementPolicyAction{
								BaseBlob: &armstorage.ManagementPolicyBaseBlob{
									Delete: &armstorage.DateAfterModification{
										DaysAfterModificationGreaterThan: to.Ptr[float32](1000),
									},
									TierToArchive: &armstorage.DateAfterModification{
										DaysAfterModificationGreaterThan: to.Ptr[float32](90),
									},
									TierToCool: &armstorage.DateAfterModification{
										DaysAfterModificationGreaterThan: to.Ptr[float32](30),
									},
								},
							},
							Filters: &armstorage.ManagementPolicyFilter{
								BlobIndexMatch: []*armstorage.TagFilter{
									{
										Name:  to.Ptr("tag1"),
										Op:    to.Ptr("=="),
										Value: to.Ptr("val1"),
									},
									{
										Name:  to.Ptr("tag2"),
										Op:    to.Ptr("=="),
										Value: to.Ptr("val2"),
									}},
								BlobTypes: []*string{
									to.Ptr("blockBlob")},
								PrefixMatch: []*string{
									to.Ptr("olcmtestcontainer2")},
							},
						},
						Enabled: to.Ptr(true),
					}},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ManagementPolicies_Get
	fmt.Println("Call operation: ManagementPolicies_Get")
	_, err = managementPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.ManagementPolicyNameDefault, nil)
	testsuite.Require().NoError(err)

	// From step ManagementPolicies_Delete
	fmt.Println("Call operation: ManagementPolicies_Delete")
	_, err = managementPoliciesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.ManagementPolicyNameDefault, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/inventoryPolicies/{blobInventoryPolicyName}
func (testsuite *StorageTestSuite) TestBlobInventoryPolicies() {
	containerName := "myblob"
	var err error
	// From step BlobContainers_Create
	fmt.Println("Call operation: BlobContainers_Create")
	blobContainersClient, err := armstorage.NewBlobContainersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = blobContainersClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, armstorage.BlobContainer{}, nil)
	testsuite.Require().NoError(err)

	// From step BlobInventoryPolicies_CreateOrUpdate
	fmt.Println("Call operation: BlobInventoryPolicies_CreateOrUpdate")
	blobInventoryPoliciesClient, err := armstorage.NewBlobInventoryPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = blobInventoryPoliciesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.BlobInventoryPolicyNameDefault, armstorage.BlobInventoryPolicy{
		Properties: &armstorage.BlobInventoryPolicyProperties{
			Policy: &armstorage.BlobInventoryPolicySchema{
				Type:    to.Ptr(armstorage.InventoryRuleTypeInventory),
				Enabled: to.Ptr(true),
				Rules: []*armstorage.BlobInventoryPolicyRule{
					{
						Name: to.Ptr("inventoryPolicyRule1"),
						Definition: &armstorage.BlobInventoryPolicyDefinition{
							Format: to.Ptr(armstorage.FormatCSV),
							Filters: &armstorage.BlobInventoryPolicyFilter{
								BlobTypes: []*string{
									to.Ptr("blockBlob")},
							},
							ObjectType: to.Ptr(armstorage.ObjectTypeBlob),
							Schedule:   to.Ptr(armstorage.ScheduleDaily),
							SchemaFields: []*string{
								to.Ptr("Name"),
								to.Ptr("Last-Modified"),
								to.Ptr("ETag"),
								to.Ptr("LeaseStatus"),
								to.Ptr("LeaseState"),
								to.Ptr("LeaseDuration"),
								to.Ptr("Metadata")},
						},
						Destination: to.Ptr(containerName),
						Enabled:     to.Ptr(true),
					}},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step BlobInventoryPolicies_List
	fmt.Println("Call operation: BlobInventoryPolicies_List")
	blobInventoryPoliciesClientNewListPager := blobInventoryPoliciesClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for blobInventoryPoliciesClientNewListPager.More() {
		_, err := blobInventoryPoliciesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BlobInventoryPolicies_Get
	fmt.Println("Call operation: BlobInventoryPolicies_Get")
	_, err = blobInventoryPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.BlobInventoryPolicyNameDefault, nil)
	testsuite.Require().NoError(err)

	// From step BlobInventoryPolicies_Delete
	fmt.Println("Call operation: BlobInventoryPolicies_Delete")
	_, err = blobInventoryPoliciesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.BlobInventoryPolicyNameDefault, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *StorageTestSuite) TestPrivateEndpointConnections() {
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
				"defaultValue": "pestorage-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "pestorage",
			},
			"storageAccountId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.storageAccountId,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "pestoragevnet",
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
							"name": "privateEndpointIpConfig.ab24488f-044e-43f0-b9d1-af1f04071719",
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
									"blob",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('storageAccountId')]",
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

	// From step PrivateEndpointConnections_List
	fmt.Println("Call operation: PrivateEndpointConnections_List")
	privateEndpointConnectionsClient, err := armstorage.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListPager := privateEndpointConnectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for privateEndpointConnectionsClientNewListPager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Put
	fmt.Println("Call operation: PrivateEndpointConnections_Put")
	_, err = privateEndpointConnectionsClient.Put(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, privateEndpointConnectionName, armstorage.PrivateEndpointConnection{
		Properties: &armstorage.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armstorage.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Auto-Approved"),
				Status:      to.Ptr(armstorage.PrivateEndpointServiceConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	_, err = privateEndpointConnectionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/privateLinkResources
func (testsuite *StorageTestSuite) TestPrivateLinkResources() {
	var err error
	// From step PrivateLinkResources_ListByStorageAccount
	fmt.Println("Call operation: PrivateLinkResources_ListByStorageAccount")
	privateLinkResourcesClient, err := armstorage.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = privateLinkResourcesClient.ListByStorageAccount(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/localUsers/{username}
func (testsuite *StorageTestSuite) TestLocalUsers() {
	username := "storageuser"
	var err error
	// From step LocalUsers_CreateOrUpdate
	fmt.Println("Call operation: LocalUsers_CreateOrUpdate")
	localUsersClient, err := armstorage.NewLocalUsersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = localUsersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, username, armstorage.LocalUser{
		Properties: &armstorage.LocalUserProperties{
			HasSSHPassword: to.Ptr(true),
			HomeDirectory:  to.Ptr("homedirectory"),
			PermissionScopes: []*armstorage.PermissionScope{
				{
					Permissions:  to.Ptr("rwd"),
					ResourceName: to.Ptr("share1"),
					Service:      to.Ptr("file"),
				},
				{
					Permissions:  to.Ptr("rw"),
					ResourceName: to.Ptr("share2"),
					Service:      to.Ptr("file"),
				}},
			SSHAuthorizedKeys: []*armstorage.SSHPublicKey{
				{
					Description: to.Ptr("key name"),
					Key:         to.Ptr("ssh-rsa keykeykeykeykey="),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step LocalUsers_List
	fmt.Println("Call operation: LocalUsers_List")
	localUsersClientNewListPager := localUsersClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for localUsersClientNewListPager.More() {
		_, err := localUsersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LocalUsers_Get
	fmt.Println("Call operation: LocalUsers_Get")
	_, err = localUsersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, username, nil)
	testsuite.Require().NoError(err)

	// From step LocalUsers_ListKeys
	fmt.Println("Call operation: LocalUsers_ListKeys")
	_, err = localUsersClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, username, nil)
	testsuite.Require().NoError(err)

	// From step LocalUsers_RegeneratePassword
	fmt.Println("Call operation: LocalUsers_RegeneratePassword")
	_, err = localUsersClient.RegeneratePassword(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, username, nil)
	testsuite.Require().NoError(err)

	// From step LocalUsers_Delete
	fmt.Println("Call operation: LocalUsers_Delete")
	_, err = localUsersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, username, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/encryptionScopes/{encryptionScopeName}
func (testsuite *StorageTestSuite) TestEncryptionScopes() {
	var err error
	// From step EncryptionScopes_Put
	fmt.Println("Call operation: EncryptionScopes_Put")
	encryptionScopesClient, err := armstorage.NewEncryptionScopesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = encryptionScopesClient.Put(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.encryptionScopeName, armstorage.EncryptionScope{}, nil)
	testsuite.Require().NoError(err)

	// From step EncryptionScopes_List
	fmt.Println("Call operation: EncryptionScopes_List")
	encryptionScopesClientNewListPager := encryptionScopesClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armstorage.EncryptionScopesClientListOptions{Maxpagesize: nil,
		Filter:  nil,
		Include: nil,
	})
	for encryptionScopesClientNewListPager.More() {
		_, err := encryptionScopesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EncryptionScopes_Get
	fmt.Println("Call operation: EncryptionScopes_Get")
	_, err = encryptionScopesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.encryptionScopeName, nil)
	testsuite.Require().NoError(err)

	// From step EncryptionScopes_Patch
	fmt.Println("Call operation: EncryptionScopes_Patch")
	_, err = encryptionScopesClient.Patch(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.encryptionScopeName, armstorage.EncryptionScope{}, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *StorageTestSuite) Cleanup() {
	var err error
	// From step StorageAccounts_Delete
	fmt.Println("Call operation: StorageAccounts_Delete")
	accountsClient, err := armstorage.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = accountsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step DeletedAccounts_List
	fmt.Println("Call operation: DeletedAccounts_List")
	deletedAccountsClient, err := armstorage.NewDeletedAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	deletedAccountsClientNewListPager := deletedAccountsClient.NewListPager(nil)
	for deletedAccountsClientNewListPager.More() {
		_, err := deletedAccountsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DeletedAccounts_Get
	fmt.Println("Call operation: DeletedAccounts_Get")
	_, err = deletedAccountsClient.Get(testsuite.ctx, testsuite.accountName, testsuite.location, nil)
	testsuite.Require().NoError(err)
}

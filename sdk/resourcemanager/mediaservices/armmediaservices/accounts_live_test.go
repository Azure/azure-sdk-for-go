//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmediaservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mediaservices/armmediaservices/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type AccountsTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	accountName        string
	storageAccountId   string
	storageAccountName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *AccountsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/mediaservices/armmediaservices/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storageaccount", 20, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AccountsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAccountsTestSuite(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	suite.Run(t, new(AccountsTestSuite))
}

func (testsuite *AccountsTestSuite) Prepare() {
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
					"name": "Standard_LRS",
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
}

func (testsuite *AccountsTestSuite) TestAccountTests() {
	var err error
	// From step Locations_CheckNameAvailability
	fmt.Println("Call operation: Locations_CheckNameAvailability")
	locationsClient, err := armmediaservices.NewLocationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = locationsClient.CheckNameAvailability(testsuite.ctx, "japanwest", armmediaservices.CheckNameAvailabilityInput{
		Name: to.Ptr(testsuite.accountName),
		Type: to.Ptr("Microsoft.Media/mediaservices"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Mediaservices_CreateOrUpdate
	fmt.Println("Call operation: Mediaservices_CreateOrUpdate")
	client, err := armmediaservices.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientCreateOrUpdateResponsePoller, err := client.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmediaservices.MediaService{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
		Properties: &armmediaservices.MediaServiceProperties{
			StorageAccounts: []*armmediaservices.StorageAccount{
				{
					Type: to.Ptr(armmediaservices.StorageAccountTypePrimary),
					ID:   to.Ptr(testsuite.storageAccountId),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Mediaservices_ListBySubscription
	fmt.Println("Call operation: Mediaservices_ListBySubscription")
	clientNewListBySubscriptionPager := client.NewListBySubscriptionPager(nil)
	for clientNewListBySubscriptionPager.More() {
		_, err := clientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Mediaservices_Get
	fmt.Println("Call operation: Mediaservices_Get")
	_, err = client.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step Mediaservices_List
	fmt.Println("Call operation: Mediaservices_List")
	clientNewListPager := client.NewListPager(testsuite.resourceGroupName, nil)
	for clientNewListPager.More() {
		_, err := clientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Mediaservices_Update
	fmt.Println("Call operation: Mediaservices_Update")
	clientUpdateResponsePoller, err := client.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmediaservices.MediaServiceUpdate{
		Tags: map[string]*string{
			"key1": to.Ptr("value3"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Mediaservices_SyncStorageKeys
	fmt.Println("Call operation: Mediaservices_SyncStorageKeys")
	_, err = client.SyncStorageKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmediaservices.SyncStorageKeysInput{
		ID: to.Ptr(testsuite.storageAccountId),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armmediaservices.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = operationsClient.List(testsuite.ctx, nil)
	testsuite.Require().NoError(err)

	// From step Mediaservices_Delete
	fmt.Println("Call operation: Mediaservices_Delete")
	_, err = client.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
}

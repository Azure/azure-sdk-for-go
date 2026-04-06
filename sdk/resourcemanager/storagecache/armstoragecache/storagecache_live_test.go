// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragecache_test

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache/v4"
	"github.com/stretchr/testify/suite"
)

type StoragecacheTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	cacheName         string
	storageTargetName string
	subnetId          string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *StoragecacheTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.cacheName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "cachenam", 14, false)
	testsuite.storageTargetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storaget", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *StoragecacheTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestStoragecacheTestSuite(t *testing.T) {
	suite.Run(t, new(StoragecacheTestSuite))
}

func (testsuite *StoragecacheTestSuite) Prepare() {
	var err error
	// From step Create_VirtualNetwork
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"subnetId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "storagecachevnet",
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"subnets": []any{
						map[string]any{
							"name": "default",
							"properties": map[string]any{
								"addressPrefix": "10.0.0.0/24",
							},
						},
					},
				},
				"tags": map[string]any{},
			},
		},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_VirtualNetwork", &deployment)
	testsuite.Require().NoError(err)
	testsuite.subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step Caches_CreateOrUpdate
	fmt.Println("Call operation: Caches_CreateOrUpdate")
	cachesClient, err := armstoragecache.NewCachesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	cachesClientCreateOrUpdateResponsePoller, err := cachesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, armstoragecache.Cache{
		Location: to.Ptr(testsuite.location),
		Properties: &armstoragecache.CacheProperties{
			CacheSizeGB: to.Ptr[int32](3072),
			Subnet:      to.Ptr(testsuite.subnetId),
		},
		SKU: &armstoragecache.CacheSKU{
			Name: to.Ptr("Standard_2G"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cachesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageCache/caches/{cacheName}
func (testsuite *StoragecacheTestSuite) TestCaches() {
	var err error
	// From step Caches_List
	fmt.Println("Call operation: Caches_List")
	cachesClient, err := armstoragecache.NewCachesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	cachesClientNewListPager := cachesClient.NewListPager(nil)
	for cachesClientNewListPager.More() {
		_, err := cachesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Caches_ListByResourceGroup
	fmt.Println("Call operation: Caches_ListByResourceGroup")
	cachesClientNewListByResourceGroupPager := cachesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for cachesClientNewListByResourceGroupPager.More() {
		_, err := cachesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Caches_Get
	fmt.Println("Call operation: Caches_Get")
	_, err = cachesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, nil)
	testsuite.Require().NoError(err)

	// From step Caches_Update
	fmt.Println("Call operation: Caches_Update")
	cachesClientUpdateResponsePoller, err := cachesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, armstoragecache.Cache{
		Tags: map[string]*string{
			"Dept": to.Ptr("Contoso"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cachesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Caches_DebugInfo
	fmt.Println("Call operation: Caches_DebugInfo")
	cachesClientDebugInfoResponsePoller, err := cachesClient.BeginDebugInfo(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cachesClientDebugInfoResponsePoller)
	testsuite.Require().NoError(err)

	// From step Caches_Start
	fmt.Println("Call operation: Caches_Start")
	cachesClientStartResponsePoller, err := cachesClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cachesClientStartResponsePoller)
	testsuite.Require().NoError(err)

	// From step Caches_UpgradeFirmware
	fmt.Println("Call operation: Caches_UpgradeFirmware")
	cachesClientUpgradeFirmwareResponsePoller, err := cachesClient.BeginUpgradeFirmware(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cachesClientUpgradeFirmwareResponsePoller)
	testsuite.Require().NoError(err)

	// From step Caches_Flush
	fmt.Println("Call operation: Caches_Flush")
	cachesClientFlushResponsePoller, err := cachesClient.BeginFlush(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cachesClientFlushResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageCache/caches/{cacheName}/storageTargets/{storageTargetName}
func (testsuite *StoragecacheTestSuite) TestStorageTargets() {
	var err error
	// From step StorageTargets_CreateOrUpdate
	fmt.Println("Call operation: StorageTargets_CreateOrUpdate")
	storageTargetsClient, err := armstoragecache.NewStorageTargetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	storageTargetsClientCreateOrUpdateResponsePoller, err := storageTargetsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, testsuite.storageTargetName, armstoragecache.StorageTarget{
		Properties: &armstoragecache.StorageTargetProperties{
			Nfs3: &armstoragecache.Nfs3Target{
				Target:            to.Ptr("10.0.44.44"),
				UsageModel:        to.Ptr("READ_ONLY"),
				VerificationTimer: to.Ptr[int32](30),
			},
			TargetType: to.Ptr(armstoragecache.StorageTargetTypeNfs3),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageTargetsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step StorageTargets_ListByCache
	fmt.Println("Call operation: StorageTargets_ListByCache")
	storageTargetsClientNewListByCachePager := storageTargetsClient.NewListByCachePager(testsuite.resourceGroupName, testsuite.cacheName, nil)
	for storageTargetsClientNewListByCachePager.More() {
		_, err := storageTargetsClientNewListByCachePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StorageTargets_Get
	fmt.Println("Call operation: StorageTargets_Get")
	_, err = storageTargetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, testsuite.storageTargetName, nil)
	testsuite.Require().NoError(err)

	// From step StorageTargets_DnsRefresh
	fmt.Println("Call operation: StorageTargets_DnsRefresh")
	storageTargetsClientDNSRefreshResponsePoller, err := storageTargetsClient.BeginDNSRefresh(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, testsuite.storageTargetName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageTargetsClientDNSRefreshResponsePoller)
	testsuite.Require().NoError(err)

	// From step StorageTarget_Flush
	fmt.Println("Call operation: StorageTarget_Flush")
	storageTargetClient, err := armstoragecache.NewStorageTargetClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	storageTargetClientFlushResponsePoller, err := storageTargetClient.BeginFlush(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, testsuite.storageTargetName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageTargetClientFlushResponsePoller)
	testsuite.Require().NoError(err)

	// From step StorageTarget_Invalidate
	fmt.Println("Call operation: StorageTarget_Invalidate")
	storageTargetClientInvalidateResponsePoller, err := storageTargetClient.BeginInvalidate(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, testsuite.storageTargetName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageTargetClientInvalidateResponsePoller)
	testsuite.Require().NoError(err)

	// From step StorageTarget_Resume
	fmt.Println("Call operation: StorageTarget_Resume")
	storageTargetClientResumeResponsePoller, err := storageTargetClient.BeginResume(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, testsuite.storageTargetName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageTargetClientResumeResponsePoller)
	testsuite.Require().NoError(err)

	// From step StorageTarget_Suspend
	fmt.Println("Call operation: StorageTarget_Suspend")
	storageTargetClientSuspendResponsePoller, err := storageTargetClient.BeginSuspend(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, testsuite.storageTargetName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageTargetClientSuspendResponsePoller)
	testsuite.Require().NoError(err)

	// From step StorageTargets_Delete
	fmt.Println("Call operation: StorageTargets_Delete")
	storageTargetsClientDeleteResponsePoller, err := storageTargetsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, testsuite.storageTargetName, &armstoragecache.StorageTargetsClientBeginDeleteOptions{Force: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageTargetsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageCache/operations
func (testsuite *StoragecacheTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armstoragecache.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.StorageCache/skus
func (testsuite *StoragecacheTestSuite) TestSkus() {
	var err error
	// From step Skus_List
	fmt.Println("Call operation: Skus_List")
	sKUsClient, err := armstoragecache.NewSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sKUsClientNewListPager := sKUsClient.NewListPager(nil)
	for sKUsClientNewListPager.More() {
		_, err := sKUsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.StorageCache/usageModels
func (testsuite *StoragecacheTestSuite) TestUsageModels() {
	var err error
	// From step UsageModels_List
	fmt.Println("Call operation: UsageModels_List")
	usageModelsClient, err := armstoragecache.NewUsageModelsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	usageModelsClientNewListPager := usageModelsClient.NewListPager(nil)
	for usageModelsClientNewListPager.More() {
		_, err := usageModelsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.StorageCache/locations/{location}/usages
func (testsuite *StoragecacheTestSuite) TestAscUsages() {
	var err error
	// From step AscUsages_List
	fmt.Println("Call operation: AscUsages_List")
	ascUsagesClient, err := armstoragecache.NewAscUsagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	ascUsagesClientNewListPager := ascUsagesClient.NewListPager(testsuite.location, nil)
	for ascUsagesClientNewListPager.More() {
		_, err := ascUsagesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *StoragecacheTestSuite) Cleanup() {
	var err error
	// From step Caches_Stop
	fmt.Println("Call operation: Caches_Stop")
	cachesClient, err := armstoragecache.NewCachesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	cachesClientStopResponsePoller, err := cachesClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cachesClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step Caches_Delete
	fmt.Println("Call operation: Caches_Delete")
	cachesClientDeleteResponsePoller, err := cachesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.cacheName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, cachesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

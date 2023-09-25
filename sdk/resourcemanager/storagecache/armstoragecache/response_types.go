//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstoragecache

// AmlFilesystemsClientArchiveResponse contains the response from method AmlFilesystemsClient.Archive.
type AmlFilesystemsClientArchiveResponse struct {
	// placeholder for future response values
}

// AmlFilesystemsClientCancelArchiveResponse contains the response from method AmlFilesystemsClient.CancelArchive.
type AmlFilesystemsClientCancelArchiveResponse struct {
	// placeholder for future response values
}

// AmlFilesystemsClientCreateOrUpdateResponse contains the response from method AmlFilesystemsClient.BeginCreateOrUpdate.
type AmlFilesystemsClientCreateOrUpdateResponse struct {
	// An AML file system instance. Follows Azure Resource Manager standards: https://github.com/Azure/azure-resource-manager-rpc/blob/master/v1.0/resource-api-reference.md
	AmlFilesystem
}

// AmlFilesystemsClientDeleteResponse contains the response from method AmlFilesystemsClient.BeginDelete.
type AmlFilesystemsClientDeleteResponse struct {
	// placeholder for future response values
}

// AmlFilesystemsClientGetResponse contains the response from method AmlFilesystemsClient.Get.
type AmlFilesystemsClientGetResponse struct {
	// An AML file system instance. Follows Azure Resource Manager standards: https://github.com/Azure/azure-resource-manager-rpc/blob/master/v1.0/resource-api-reference.md
	AmlFilesystem
}

// AmlFilesystemsClientListByResourceGroupResponse contains the response from method AmlFilesystemsClient.NewListByResourceGroupPager.
type AmlFilesystemsClientListByResourceGroupResponse struct {
	// Result of the request to list AML file systems. It contains a list of AML file systems and a URL link to get the next set
// of results.
	AmlFilesystemsListResult
}

// AmlFilesystemsClientListResponse contains the response from method AmlFilesystemsClient.NewListPager.
type AmlFilesystemsClientListResponse struct {
	// Result of the request to list AML file systems. It contains a list of AML file systems and a URL link to get the next set
// of results.
	AmlFilesystemsListResult
}

// AmlFilesystemsClientUpdateResponse contains the response from method AmlFilesystemsClient.BeginUpdate.
type AmlFilesystemsClientUpdateResponse struct {
	// An AML file system instance. Follows Azure Resource Manager standards: https://github.com/Azure/azure-resource-manager-rpc/blob/master/v1.0/resource-api-reference.md
	AmlFilesystem
}

// AscOperationsClientGetResponse contains the response from method AscOperationsClient.Get.
type AscOperationsClientGetResponse struct {
	// The status of operation.
	AscOperation
}

// AscUsagesClientListResponse contains the response from method AscUsagesClient.NewListPager.
type AscUsagesClientListResponse struct {
	// Result of the request to list resource usages. It contains a list of resource usages & limits and a URL link to get the
// next set of results.
	ResourceUsagesListResult
}

// CachesClientCreateOrUpdateResponse contains the response from method CachesClient.BeginCreateOrUpdate.
type CachesClientCreateOrUpdateResponse struct {
	// A cache instance. Follows Azure Resource Manager standards: https://github.com/Azure/azure-resource-manager-rpc/blob/master/v1.0/resource-api-reference.md
	Cache
}

// CachesClientDebugInfoResponse contains the response from method CachesClient.BeginDebugInfo.
type CachesClientDebugInfoResponse struct {
	// placeholder for future response values
}

// CachesClientDeleteResponse contains the response from method CachesClient.BeginDelete.
type CachesClientDeleteResponse struct {
	// placeholder for future response values
}

// CachesClientFlushResponse contains the response from method CachesClient.BeginFlush.
type CachesClientFlushResponse struct {
	// placeholder for future response values
}

// CachesClientGetResponse contains the response from method CachesClient.Get.
type CachesClientGetResponse struct {
	// A cache instance. Follows Azure Resource Manager standards: https://github.com/Azure/azure-resource-manager-rpc/blob/master/v1.0/resource-api-reference.md
	Cache
}

// CachesClientListByResourceGroupResponse contains the response from method CachesClient.NewListByResourceGroupPager.
type CachesClientListByResourceGroupResponse struct {
	// Result of the request to list caches. It contains a list of caches and a URL link to get the next set of results.
	CachesListResult
}

// CachesClientListResponse contains the response from method CachesClient.NewListPager.
type CachesClientListResponse struct {
	// Result of the request to list caches. It contains a list of caches and a URL link to get the next set of results.
	CachesListResult
}

// CachesClientPausePrimingJobResponse contains the response from method CachesClient.BeginPausePrimingJob.
type CachesClientPausePrimingJobResponse struct {
	// placeholder for future response values
}

// CachesClientResumePrimingJobResponse contains the response from method CachesClient.BeginResumePrimingJob.
type CachesClientResumePrimingJobResponse struct {
	// placeholder for future response values
}

// CachesClientSpaceAllocationResponse contains the response from method CachesClient.BeginSpaceAllocation.
type CachesClientSpaceAllocationResponse struct {
	// placeholder for future response values
}

// CachesClientStartPrimingJobResponse contains the response from method CachesClient.BeginStartPrimingJob.
type CachesClientStartPrimingJobResponse struct {
	// placeholder for future response values
}

// CachesClientStartResponse contains the response from method CachesClient.BeginStart.
type CachesClientStartResponse struct {
	// placeholder for future response values
}

// CachesClientStopPrimingJobResponse contains the response from method CachesClient.BeginStopPrimingJob.
type CachesClientStopPrimingJobResponse struct {
	// placeholder for future response values
}

// CachesClientStopResponse contains the response from method CachesClient.BeginStop.
type CachesClientStopResponse struct {
	// placeholder for future response values
}

// CachesClientUpdateResponse contains the response from method CachesClient.BeginUpdate.
type CachesClientUpdateResponse struct {
	// A cache instance. Follows Azure Resource Manager standards: https://github.com/Azure/azure-resource-manager-rpc/blob/master/v1.0/resource-api-reference.md
	Cache
}

// CachesClientUpgradeFirmwareResponse contains the response from method CachesClient.BeginUpgradeFirmware.
type CachesClientUpgradeFirmwareResponse struct {
	// placeholder for future response values
}

// ManagementClientCheckAmlFSSubnetsResponse contains the response from method ManagementClient.CheckAmlFSSubnets.
type ManagementClientCheckAmlFSSubnetsResponse struct {
	// placeholder for future response values
}

// ManagementClientGetRequiredAmlFSSubnetsSizeResponse contains the response from method ManagementClient.GetRequiredAmlFSSubnetsSize.
type ManagementClientGetRequiredAmlFSSubnetsSizeResponse struct {
	// Information about the number of available IP addresses that are required for the AML file system.
	RequiredAmlFilesystemSubnetsSize
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// Result of the request to list Resource Provider operations. It contains a list of operations and a URL link to get the
// next set of results.
	APIOperationListResult
}

// SKUsClientListResponse contains the response from method SKUsClient.NewListPager.
type SKUsClientListResponse struct {
	// The response from the List Cache SKUs operation.
	ResourceSKUsResult
}

// StorageTargetClientFlushResponse contains the response from method StorageTargetClient.BeginFlush.
type StorageTargetClientFlushResponse struct {
	// placeholder for future response values
}

// StorageTargetClientInvalidateResponse contains the response from method StorageTargetClient.BeginInvalidate.
type StorageTargetClientInvalidateResponse struct {
	// placeholder for future response values
}

// StorageTargetClientResumeResponse contains the response from method StorageTargetClient.BeginResume.
type StorageTargetClientResumeResponse struct {
	// placeholder for future response values
}

// StorageTargetClientSuspendResponse contains the response from method StorageTargetClient.BeginSuspend.
type StorageTargetClientSuspendResponse struct {
	// placeholder for future response values
}

// StorageTargetsClientCreateOrUpdateResponse contains the response from method StorageTargetsClient.BeginCreateOrUpdate.
type StorageTargetsClientCreateOrUpdateResponse struct {
	// Type of the Storage Target.
	StorageTarget
}

// StorageTargetsClientDNSRefreshResponse contains the response from method StorageTargetsClient.BeginDNSRefresh.
type StorageTargetsClientDNSRefreshResponse struct {
	// placeholder for future response values
}

// StorageTargetsClientDeleteResponse contains the response from method StorageTargetsClient.BeginDelete.
type StorageTargetsClientDeleteResponse struct {
	// placeholder for future response values
}

// StorageTargetsClientGetResponse contains the response from method StorageTargetsClient.Get.
type StorageTargetsClientGetResponse struct {
	// Type of the Storage Target.
	StorageTarget
}

// StorageTargetsClientListByCacheResponse contains the response from method StorageTargetsClient.NewListByCachePager.
type StorageTargetsClientListByCacheResponse struct {
	// A list of Storage Targets.
	StorageTargetsResult
}

// StorageTargetsClientRestoreDefaultsResponse contains the response from method StorageTargetsClient.BeginRestoreDefaults.
type StorageTargetsClientRestoreDefaultsResponse struct {
	// placeholder for future response values
}

// UsageModelsClientListResponse contains the response from method UsageModelsClient.NewListPager.
type UsageModelsClientListResponse struct {
	// A list of cache usage models.
	UsageModelsResult
}


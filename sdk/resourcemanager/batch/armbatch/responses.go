//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armbatch

// AccountClientCreateResponse contains the response from method AccountClient.BeginCreate.
type AccountClientCreateResponse struct {
	// Contains information about an Azure Batch account.
	Account
}

// AccountClientDeleteResponse contains the response from method AccountClient.BeginDelete.
type AccountClientDeleteResponse struct {
	// placeholder for future response values
}

// AccountClientGetDetectorResponse contains the response from method AccountClient.GetDetector.
type AccountClientGetDetectorResponse struct {
	// Contains the information for a detector.
	DetectorResponse
}

// AccountClientGetKeysResponse contains the response from method AccountClient.GetKeys.
type AccountClientGetKeysResponse struct {
	// A set of Azure Batch account keys.
	AccountKeys
}

// AccountClientGetResponse contains the response from method AccountClient.Get.
type AccountClientGetResponse struct {
	// Contains information about an Azure Batch account.
	Account
}

// AccountClientListByResourceGroupResponse contains the response from method AccountClient.NewListByResourceGroupPager.
type AccountClientListByResourceGroupResponse struct {
	// Values returned by the List operation.
	AccountListResult
}

// AccountClientListDetectorsResponse contains the response from method AccountClient.NewListDetectorsPager.
type AccountClientListDetectorsResponse struct {
	// Values returned by the List operation.
	DetectorListResult
}

// AccountClientListOutboundNetworkDependenciesEndpointsResponse contains the response from method AccountClient.NewListOutboundNetworkDependenciesEndpointsPager.
type AccountClientListOutboundNetworkDependenciesEndpointsResponse struct {
	// Values returned by the List operation.
	OutboundEnvironmentEndpointCollection
}

// AccountClientListResponse contains the response from method AccountClient.NewListPager.
type AccountClientListResponse struct {
	// Values returned by the List operation.
	AccountListResult
}

// AccountClientRegenerateKeyResponse contains the response from method AccountClient.RegenerateKey.
type AccountClientRegenerateKeyResponse struct {
	// A set of Azure Batch account keys.
	AccountKeys
}

// AccountClientSynchronizeAutoStorageKeysResponse contains the response from method AccountClient.SynchronizeAutoStorageKeys.
type AccountClientSynchronizeAutoStorageKeysResponse struct {
	// placeholder for future response values
}

// AccountClientUpdateResponse contains the response from method AccountClient.Update.
type AccountClientUpdateResponse struct {
	// Contains information about an Azure Batch account.
	Account
}

// ApplicationClientCreateResponse contains the response from method ApplicationClient.Create.
type ApplicationClientCreateResponse struct {
	// Contains information about an application in a Batch account.
	Application
}

// ApplicationClientDeleteResponse contains the response from method ApplicationClient.Delete.
type ApplicationClientDeleteResponse struct {
	// placeholder for future response values
}

// ApplicationClientGetResponse contains the response from method ApplicationClient.Get.
type ApplicationClientGetResponse struct {
	// Contains information about an application in a Batch account.
	Application
}

// ApplicationClientListResponse contains the response from method ApplicationClient.NewListPager.
type ApplicationClientListResponse struct {
	// The result of performing list applications.
	ListApplicationsResult
}

// ApplicationClientUpdateResponse contains the response from method ApplicationClient.Update.
type ApplicationClientUpdateResponse struct {
	// Contains information about an application in a Batch account.
	Application
}

// ApplicationPackageClientActivateResponse contains the response from method ApplicationPackageClient.Activate.
type ApplicationPackageClientActivateResponse struct {
	// An application package which represents a particular version of an application.
	ApplicationPackage
}

// ApplicationPackageClientCreateResponse contains the response from method ApplicationPackageClient.Create.
type ApplicationPackageClientCreateResponse struct {
	// An application package which represents a particular version of an application.
	ApplicationPackage
}

// ApplicationPackageClientDeleteResponse contains the response from method ApplicationPackageClient.Delete.
type ApplicationPackageClientDeleteResponse struct {
	// placeholder for future response values
}

// ApplicationPackageClientGetResponse contains the response from method ApplicationPackageClient.Get.
type ApplicationPackageClientGetResponse struct {
	// An application package which represents a particular version of an application.
	ApplicationPackage
}

// ApplicationPackageClientListResponse contains the response from method ApplicationPackageClient.NewListPager.
type ApplicationPackageClientListResponse struct {
	// The result of performing list application packages.
	ListApplicationPackagesResult
}

// CertificateClientCancelDeletionResponse contains the response from method CertificateClient.CancelDeletion.
type CertificateClientCancelDeletionResponse struct {
	// Contains information about a certificate.
	Certificate

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// CertificateClientCreateResponse contains the response from method CertificateClient.Create.
type CertificateClientCreateResponse struct {
	// Contains information about a certificate.
	Certificate

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// CertificateClientDeleteResponse contains the response from method CertificateClient.BeginDelete.
type CertificateClientDeleteResponse struct {
	// placeholder for future response values
}

// CertificateClientGetResponse contains the response from method CertificateClient.Get.
type CertificateClientGetResponse struct {
	// Contains information about a certificate.
	Certificate

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// CertificateClientListByBatchAccountResponse contains the response from method CertificateClient.NewListByBatchAccountPager.
type CertificateClientListByBatchAccountResponse struct {
	// Values returned by the List operation.
	ListCertificatesResult
}

// CertificateClientUpdateResponse contains the response from method CertificateClient.Update.
type CertificateClientUpdateResponse struct {
	// Contains information about a certificate.
	Certificate

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// LocationClientCheckNameAvailabilityResponse contains the response from method LocationClient.CheckNameAvailability.
type LocationClientCheckNameAvailabilityResponse struct {
	// The CheckNameAvailability operation response.
	CheckNameAvailabilityResult
}

// LocationClientGetQuotasResponse contains the response from method LocationClient.GetQuotas.
type LocationClientGetQuotasResponse struct {
	// Quotas associated with a Batch region for a particular subscription.
	LocationQuota
}

// LocationClientListSupportedVirtualMachineSKUsResponse contains the response from method LocationClient.NewListSupportedVirtualMachineSKUsPager.
type LocationClientListSupportedVirtualMachineSKUsResponse struct {
	// The Batch List supported SKUs operation response.
	SupportedSKUsResult
}

// NetworkSecurityPerimeterClientGetConfigurationResponse contains the response from method NetworkSecurityPerimeterClient.GetConfiguration.
type NetworkSecurityPerimeterClientGetConfigurationResponse struct {
	// Network security perimeter (NSP) configuration resource
	NetworkSecurityPerimeterConfiguration
}

// NetworkSecurityPerimeterClientListConfigurationsResponse contains the response from method NetworkSecurityPerimeterClient.NewListConfigurationsPager.
type NetworkSecurityPerimeterClientListConfigurationsResponse struct {
	// Result of a list NSP (network security perimeter) configurations request.
	NetworkSecurityPerimeterConfigurationListResult
}

// NetworkSecurityPerimeterClientReconcileConfigurationResponse contains the response from method NetworkSecurityPerimeterClient.BeginReconcileConfiguration.
type NetworkSecurityPerimeterClientReconcileConfigurationResponse struct {
	// placeholder for future response values
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// Result of the request to list REST API operations. It contains a list of operations and a URL nextLink to get the next
	// set of results.
	OperationListResult
}

// PoolClientCreateResponse contains the response from method PoolClient.Create.
type PoolClientCreateResponse struct {
	// Contains information about a pool.
	Pool

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// PoolClientDeleteResponse contains the response from method PoolClient.BeginDelete.
type PoolClientDeleteResponse struct {
	// placeholder for future response values
}

// PoolClientDisableAutoScaleResponse contains the response from method PoolClient.DisableAutoScale.
type PoolClientDisableAutoScaleResponse struct {
	// Contains information about a pool.
	Pool

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// PoolClientGetResponse contains the response from method PoolClient.Get.
type PoolClientGetResponse struct {
	// Contains information about a pool.
	Pool

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// PoolClientListByBatchAccountResponse contains the response from method PoolClient.NewListByBatchAccountPager.
type PoolClientListByBatchAccountResponse struct {
	// Values returned by the List operation.
	ListPoolsResult
}

// PoolClientStopResizeResponse contains the response from method PoolClient.StopResize.
type PoolClientStopResizeResponse struct {
	// Contains information about a pool.
	Pool

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// PoolClientUpdateResponse contains the response from method PoolClient.Update.
type PoolClientUpdateResponse struct {
	// Contains information about a pool.
	Pool

	// ETag contains the information returned from the ETag header response.
	ETag *string
}

// PrivateEndpointConnectionClientDeleteResponse contains the response from method PrivateEndpointConnectionClient.BeginDelete.
type PrivateEndpointConnectionClientDeleteResponse struct {
	// placeholder for future response values
}

// PrivateEndpointConnectionClientGetResponse contains the response from method PrivateEndpointConnectionClient.Get.
type PrivateEndpointConnectionClientGetResponse struct {
	// Contains information about a private link resource.
	PrivateEndpointConnection
}

// PrivateEndpointConnectionClientListByBatchAccountResponse contains the response from method PrivateEndpointConnectionClient.NewListByBatchAccountPager.
type PrivateEndpointConnectionClientListByBatchAccountResponse struct {
	// Values returned by the List operation.
	ListPrivateEndpointConnectionsResult
}

// PrivateEndpointConnectionClientUpdateResponse contains the response from method PrivateEndpointConnectionClient.BeginUpdate.
type PrivateEndpointConnectionClientUpdateResponse struct {
	// Contains information about a private link resource.
	PrivateEndpointConnection
}

// PrivateLinkResourceClientGetResponse contains the response from method PrivateLinkResourceClient.Get.
type PrivateLinkResourceClientGetResponse struct {
	// Contains information about a private link resource.
	PrivateLinkResource
}

// PrivateLinkResourceClientListByBatchAccountResponse contains the response from method PrivateLinkResourceClient.NewListByBatchAccountPager.
type PrivateLinkResourceClientListByBatchAccountResponse struct {
	// Values returned by the List operation.
	ListPrivateLinkResourcesResult
}

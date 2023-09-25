//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkeyvault

// KeysClientCreateIfNotExistResponse contains the response from method KeysClient.CreateIfNotExist.
type KeysClientCreateIfNotExistResponse struct {
	// The key resource.
	Key
}

// KeysClientGetResponse contains the response from method KeysClient.Get.
type KeysClientGetResponse struct {
	// The key resource.
	Key
}

// KeysClientGetVersionResponse contains the response from method KeysClient.GetVersion.
type KeysClientGetVersionResponse struct {
	// The key resource.
	Key
}

// KeysClientListResponse contains the response from method KeysClient.NewListPager.
type KeysClientListResponse struct {
	// The page of keys.
	KeyListResult
}

// KeysClientListVersionsResponse contains the response from method KeysClient.NewListVersionsPager.
type KeysClientListVersionsResponse struct {
	// The page of keys.
	KeyListResult
}

// MHSMPrivateEndpointConnectionsClientDeleteResponse contains the response from method MHSMPrivateEndpointConnectionsClient.BeginDelete.
type MHSMPrivateEndpointConnectionsClientDeleteResponse struct {
	// Private endpoint connection resource.
	MHSMPrivateEndpointConnection
}

// MHSMPrivateEndpointConnectionsClientGetResponse contains the response from method MHSMPrivateEndpointConnectionsClient.Get.
type MHSMPrivateEndpointConnectionsClientGetResponse struct {
	// Private endpoint connection resource.
	MHSMPrivateEndpointConnection
}

// MHSMPrivateEndpointConnectionsClientListByResourceResponse contains the response from method MHSMPrivateEndpointConnectionsClient.NewListByResourcePager.
type MHSMPrivateEndpointConnectionsClientListByResourceResponse struct {
	// List of private endpoint connections associated with a managed HSM Pools
	MHSMPrivateEndpointConnectionsListResult
}

// MHSMPrivateEndpointConnectionsClientPutResponse contains the response from method MHSMPrivateEndpointConnectionsClient.Put.
type MHSMPrivateEndpointConnectionsClientPutResponse struct {
	// Private endpoint connection resource.
	MHSMPrivateEndpointConnection

	// AzureAsyncOperation contains the information returned from the Azure-AsyncOperation header response.
	AzureAsyncOperation *string

	// RetryAfter contains the information returned from the Retry-After header response.
	RetryAfter *int32
}

// MHSMPrivateLinkResourcesClientListByMHSMResourceResponse contains the response from method MHSMPrivateLinkResourcesClient.ListByMHSMResource.
type MHSMPrivateLinkResourcesClientListByMHSMResourceResponse struct {
	// A list of private link resources
	MHSMPrivateLinkResourceListResult
}

// MHSMRegionsClientListByResourceResponse contains the response from method MHSMRegionsClient.NewListByResourcePager.
type MHSMRegionsClientListByResourceResponse struct {
	// List of regions associated with a managed HSM Pools
	MHSMRegionsListResult
}

// ManagedHsmKeysClientCreateIfNotExistResponse contains the response from method ManagedHsmKeysClient.CreateIfNotExist.
type ManagedHsmKeysClientCreateIfNotExistResponse struct {
	// The key resource.
	ManagedHsmKey
}

// ManagedHsmKeysClientGetResponse contains the response from method ManagedHsmKeysClient.Get.
type ManagedHsmKeysClientGetResponse struct {
	// The key resource.
	ManagedHsmKey
}

// ManagedHsmKeysClientGetVersionResponse contains the response from method ManagedHsmKeysClient.GetVersion.
type ManagedHsmKeysClientGetVersionResponse struct {
	// The key resource.
	ManagedHsmKey
}

// ManagedHsmKeysClientListResponse contains the response from method ManagedHsmKeysClient.NewListPager.
type ManagedHsmKeysClientListResponse struct {
	// The page of keys.
	ManagedHsmKeyListResult
}

// ManagedHsmKeysClientListVersionsResponse contains the response from method ManagedHsmKeysClient.NewListVersionsPager.
type ManagedHsmKeysClientListVersionsResponse struct {
	// The page of keys.
	ManagedHsmKeyListResult
}

// ManagedHsmsClientCheckMhsmNameAvailabilityResponse contains the response from method ManagedHsmsClient.CheckMhsmNameAvailability.
type ManagedHsmsClientCheckMhsmNameAvailabilityResponse struct {
	// The CheckMhsmNameAvailability operation response.
	CheckMhsmNameAvailabilityResult
}

// ManagedHsmsClientCreateOrUpdateResponse contains the response from method ManagedHsmsClient.BeginCreateOrUpdate.
type ManagedHsmsClientCreateOrUpdateResponse struct {
	// Resource information with extended details.
	ManagedHsm
}

// ManagedHsmsClientDeleteResponse contains the response from method ManagedHsmsClient.BeginDelete.
type ManagedHsmsClientDeleteResponse struct {
	// placeholder for future response values
}

// ManagedHsmsClientGetDeletedResponse contains the response from method ManagedHsmsClient.GetDeleted.
type ManagedHsmsClientGetDeletedResponse struct {
	DeletedManagedHsm
}

// ManagedHsmsClientGetResponse contains the response from method ManagedHsmsClient.Get.
type ManagedHsmsClientGetResponse struct {
	// Resource information with extended details.
	ManagedHsm
}

// ManagedHsmsClientListByResourceGroupResponse contains the response from method ManagedHsmsClient.NewListByResourceGroupPager.
type ManagedHsmsClientListByResourceGroupResponse struct {
	// List of managed HSM Pools
	ManagedHsmListResult
}

// ManagedHsmsClientListBySubscriptionResponse contains the response from method ManagedHsmsClient.NewListBySubscriptionPager.
type ManagedHsmsClientListBySubscriptionResponse struct {
	// List of managed HSM Pools
	ManagedHsmListResult
}

// ManagedHsmsClientListDeletedResponse contains the response from method ManagedHsmsClient.NewListDeletedPager.
type ManagedHsmsClientListDeletedResponse struct {
	// List of deleted managed HSM Pools
	DeletedManagedHsmListResult
}

// ManagedHsmsClientPurgeDeletedResponse contains the response from method ManagedHsmsClient.BeginPurgeDeleted.
type ManagedHsmsClientPurgeDeletedResponse struct {
	// placeholder for future response values
}

// ManagedHsmsClientUpdateResponse contains the response from method ManagedHsmsClient.BeginUpdate.
type ManagedHsmsClientUpdateResponse struct {
	// Resource information with extended details.
	ManagedHsm
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// Result of the request to list Storage operations. It contains a list of operations and a URL link to get the next set of
// results.
	OperationListResult
}

// PrivateEndpointConnectionsClientDeleteResponse contains the response from method PrivateEndpointConnectionsClient.BeginDelete.
type PrivateEndpointConnectionsClientDeleteResponse struct {
	// Private endpoint connection resource.
	PrivateEndpointConnection
}

// PrivateEndpointConnectionsClientGetResponse contains the response from method PrivateEndpointConnectionsClient.Get.
type PrivateEndpointConnectionsClientGetResponse struct {
	// Private endpoint connection resource.
	PrivateEndpointConnection
}

// PrivateEndpointConnectionsClientListByResourceResponse contains the response from method PrivateEndpointConnectionsClient.NewListByResourcePager.
type PrivateEndpointConnectionsClientListByResourceResponse struct {
	// List of private endpoint connections.
	PrivateEndpointConnectionListResult
}

// PrivateEndpointConnectionsClientPutResponse contains the response from method PrivateEndpointConnectionsClient.Put.
type PrivateEndpointConnectionsClientPutResponse struct {
	// Private endpoint connection resource.
	PrivateEndpointConnection

	// AzureAsyncOperation contains the information returned from the Azure-AsyncOperation header response.
	AzureAsyncOperation *string

	// RetryAfter contains the information returned from the Retry-After header response.
	RetryAfter *int32
}

// PrivateLinkResourcesClientListByVaultResponse contains the response from method PrivateLinkResourcesClient.ListByVault.
type PrivateLinkResourcesClientListByVaultResponse struct {
	// A list of private link resources
	PrivateLinkResourceListResult
}

// SecretsClientCreateOrUpdateResponse contains the response from method SecretsClient.CreateOrUpdate.
type SecretsClientCreateOrUpdateResponse struct {
	// Resource information with extended details.
	Secret
}

// SecretsClientGetResponse contains the response from method SecretsClient.Get.
type SecretsClientGetResponse struct {
	// Resource information with extended details.
	Secret
}

// SecretsClientListResponse contains the response from method SecretsClient.NewListPager.
type SecretsClientListResponse struct {
	// List of secrets
	SecretListResult
}

// SecretsClientUpdateResponse contains the response from method SecretsClient.Update.
type SecretsClientUpdateResponse struct {
	// Resource information with extended details.
	Secret
}

// VaultsClientCheckNameAvailabilityResponse contains the response from method VaultsClient.CheckNameAvailability.
type VaultsClientCheckNameAvailabilityResponse struct {
	// The CheckNameAvailability operation response.
	CheckNameAvailabilityResult
}

// VaultsClientCreateOrUpdateResponse contains the response from method VaultsClient.BeginCreateOrUpdate.
type VaultsClientCreateOrUpdateResponse struct {
	// Resource information with extended details.
	Vault
}

// VaultsClientDeleteResponse contains the response from method VaultsClient.Delete.
type VaultsClientDeleteResponse struct {
	// placeholder for future response values
}

// VaultsClientGetDeletedResponse contains the response from method VaultsClient.GetDeleted.
type VaultsClientGetDeletedResponse struct {
	// Deleted vault information with extended details.
	DeletedVault
}

// VaultsClientGetResponse contains the response from method VaultsClient.Get.
type VaultsClientGetResponse struct {
	// Resource information with extended details.
	Vault
}

// VaultsClientListByResourceGroupResponse contains the response from method VaultsClient.NewListByResourceGroupPager.
type VaultsClientListByResourceGroupResponse struct {
	// List of vaults
	VaultListResult
}

// VaultsClientListBySubscriptionResponse contains the response from method VaultsClient.NewListBySubscriptionPager.
type VaultsClientListBySubscriptionResponse struct {
	// List of vaults
	VaultListResult
}

// VaultsClientListDeletedResponse contains the response from method VaultsClient.NewListDeletedPager.
type VaultsClientListDeletedResponse struct {
	// List of vaults
	DeletedVaultListResult
}

// VaultsClientListResponse contains the response from method VaultsClient.NewListPager.
type VaultsClientListResponse struct {
	// List of vault resources.
	ResourceListResult
}

// VaultsClientPurgeDeletedResponse contains the response from method VaultsClient.BeginPurgeDeleted.
type VaultsClientPurgeDeletedResponse struct {
	// placeholder for future response values
}

// VaultsClientUpdateAccessPolicyResponse contains the response from method VaultsClient.UpdateAccessPolicy.
type VaultsClientUpdateAccessPolicyResponse struct {
	// Parameters for updating the access policy in a vault
	VaultAccessPolicyParameters
}

// VaultsClientUpdateResponse contains the response from method VaultsClient.Update.
type VaultsClientUpdateResponse struct {
	// Resource information with extended details.
	Vault
}


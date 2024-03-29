//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridcontainerservice

// AgentPoolClientBeginCreateOrUpdateOptions contains the optional parameters for the AgentPoolClient.BeginCreateOrUpdate
// method.
type AgentPoolClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AgentPoolClientBeginDeleteOptions contains the optional parameters for the AgentPoolClient.BeginDelete method.
type AgentPoolClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AgentPoolClientGetOptions contains the optional parameters for the AgentPoolClient.Get method.
type AgentPoolClientGetOptions struct {
	// placeholder for future optional parameters
}

// AgentPoolClientListByProvisionedClusterOptions contains the optional parameters for the AgentPoolClient.NewListByProvisionedClusterPager
// method.
type AgentPoolClientListByProvisionedClusterOptions struct {
	// placeholder for future optional parameters
}

// ClientBeginDeleteKubernetesVersionsOptions contains the optional parameters for the Client.BeginDeleteKubernetesVersions
// method.
type ClientBeginDeleteKubernetesVersionsOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ClientBeginDeleteVMSKUsOptions contains the optional parameters for the Client.BeginDeleteVMSKUs method.
type ClientBeginDeleteVMSKUsOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ClientBeginPutKubernetesVersionsOptions contains the optional parameters for the Client.BeginPutKubernetesVersions method.
type ClientBeginPutKubernetesVersionsOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ClientBeginPutVMSKUsOptions contains the optional parameters for the Client.BeginPutVMSKUs method.
type ClientBeginPutVMSKUsOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ClientGetKubernetesVersionsOptions contains the optional parameters for the Client.GetKubernetesVersions method.
type ClientGetKubernetesVersionsOptions struct {
	// placeholder for future optional parameters
}

// ClientGetVMSKUsOptions contains the optional parameters for the Client.GetVMSKUs method.
type ClientGetVMSKUsOptions struct {
	// placeholder for future optional parameters
}

// HybridIdentityMetadataClientBeginDeleteOptions contains the optional parameters for the HybridIdentityMetadataClient.BeginDelete
// method.
type HybridIdentityMetadataClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// HybridIdentityMetadataClientGetOptions contains the optional parameters for the HybridIdentityMetadataClient.Get method.
type HybridIdentityMetadataClientGetOptions struct {
	// placeholder for future optional parameters
}

// HybridIdentityMetadataClientListByClusterOptions contains the optional parameters for the HybridIdentityMetadataClient.NewListByClusterPager
// method.
type HybridIdentityMetadataClientListByClusterOptions struct {
	// placeholder for future optional parameters
}

// HybridIdentityMetadataClientPutOptions contains the optional parameters for the HybridIdentityMetadataClient.Put method.
type HybridIdentityMetadataClientPutOptions struct {
	// placeholder for future optional parameters
}

// KubernetesVersionsClientListOptions contains the optional parameters for the KubernetesVersionsClient.NewListPager method.
type KubernetesVersionsClientListOptions struct {
	// placeholder for future optional parameters
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.NewListPager method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// ProvisionedClusterInstancesClientBeginCreateOrUpdateOptions contains the optional parameters for the ProvisionedClusterInstancesClient.BeginCreateOrUpdate
// method.
type ProvisionedClusterInstancesClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ProvisionedClusterInstancesClientBeginDeleteOptions contains the optional parameters for the ProvisionedClusterInstancesClient.BeginDelete
// method.
type ProvisionedClusterInstancesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ProvisionedClusterInstancesClientBeginListAdminKubeconfigOptions contains the optional parameters for the ProvisionedClusterInstancesClient.BeginListAdminKubeconfig
// method.
type ProvisionedClusterInstancesClientBeginListAdminKubeconfigOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ProvisionedClusterInstancesClientBeginListUserKubeconfigOptions contains the optional parameters for the ProvisionedClusterInstancesClient.BeginListUserKubeconfig
// method.
type ProvisionedClusterInstancesClientBeginListUserKubeconfigOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ProvisionedClusterInstancesClientGetOptions contains the optional parameters for the ProvisionedClusterInstancesClient.Get
// method.
type ProvisionedClusterInstancesClientGetOptions struct {
	// placeholder for future optional parameters
}

// ProvisionedClusterInstancesClientGetUpgradeProfileOptions contains the optional parameters for the ProvisionedClusterInstancesClient.GetUpgradeProfile
// method.
type ProvisionedClusterInstancesClientGetUpgradeProfileOptions struct {
	// placeholder for future optional parameters
}

// ProvisionedClusterInstancesClientListOptions contains the optional parameters for the ProvisionedClusterInstancesClient.NewListPager
// method.
type ProvisionedClusterInstancesClientListOptions struct {
	// placeholder for future optional parameters
}

// VMSKUsClientListOptions contains the optional parameters for the VMSKUsClient.NewListPager method.
type VMSKUsClientListOptions struct {
	// placeholder for future optional parameters
}

// VirtualNetworksClientBeginCreateOrUpdateOptions contains the optional parameters for the VirtualNetworksClient.BeginCreateOrUpdate
// method.
type VirtualNetworksClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VirtualNetworksClientBeginDeleteOptions contains the optional parameters for the VirtualNetworksClient.BeginDelete method.
type VirtualNetworksClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VirtualNetworksClientBeginUpdateOptions contains the optional parameters for the VirtualNetworksClient.BeginUpdate method.
type VirtualNetworksClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VirtualNetworksClientListByResourceGroupOptions contains the optional parameters for the VirtualNetworksClient.NewListByResourceGroupPager
// method.
type VirtualNetworksClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// VirtualNetworksClientListBySubscriptionOptions contains the optional parameters for the VirtualNetworksClient.NewListBySubscriptionPager
// method.
type VirtualNetworksClientListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// VirtualNetworksClientRetrieveOptions contains the optional parameters for the VirtualNetworksClient.Retrieve method.
type VirtualNetworksClientRetrieveOptions struct {
	// placeholder for future optional parameters
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armscvmm

// AvailabilitySetsClientCreateOrUpdateResponse contains the response from method AvailabilitySetsClient.BeginCreateOrUpdate.
type AvailabilitySetsClientCreateOrUpdateResponse struct {
	// The AvailabilitySets resource definition.
	AvailabilitySet
}

// AvailabilitySetsClientDeleteResponse contains the response from method AvailabilitySetsClient.BeginDelete.
type AvailabilitySetsClientDeleteResponse struct {
	// placeholder for future response values
}

// AvailabilitySetsClientGetResponse contains the response from method AvailabilitySetsClient.Get.
type AvailabilitySetsClientGetResponse struct {
	// The AvailabilitySets resource definition.
	AvailabilitySet
}

// AvailabilitySetsClientListByResourceGroupResponse contains the response from method AvailabilitySetsClient.NewListByResourceGroupPager.
type AvailabilitySetsClientListByResourceGroupResponse struct {
	// List of AvailabilitySets.
	AvailabilitySetListResult
}

// AvailabilitySetsClientListBySubscriptionResponse contains the response from method AvailabilitySetsClient.NewListBySubscriptionPager.
type AvailabilitySetsClientListBySubscriptionResponse struct {
	// List of AvailabilitySets.
	AvailabilitySetListResult
}

// AvailabilitySetsClientUpdateResponse contains the response from method AvailabilitySetsClient.BeginUpdate.
type AvailabilitySetsClientUpdateResponse struct {
	// The AvailabilitySets resource definition.
	AvailabilitySet
}

// CloudsClientCreateOrUpdateResponse contains the response from method CloudsClient.BeginCreateOrUpdate.
type CloudsClientCreateOrUpdateResponse struct {
	// The Clouds resource definition.
	Cloud
}

// CloudsClientDeleteResponse contains the response from method CloudsClient.BeginDelete.
type CloudsClientDeleteResponse struct {
	// placeholder for future response values
}

// CloudsClientGetResponse contains the response from method CloudsClient.Get.
type CloudsClientGetResponse struct {
	// The Clouds resource definition.
	Cloud
}

// CloudsClientListByResourceGroupResponse contains the response from method CloudsClient.NewListByResourceGroupPager.
type CloudsClientListByResourceGroupResponse struct {
	// List of Clouds.
	CloudListResult
}

// CloudsClientListBySubscriptionResponse contains the response from method CloudsClient.NewListBySubscriptionPager.
type CloudsClientListBySubscriptionResponse struct {
	// List of Clouds.
	CloudListResult
}

// CloudsClientUpdateResponse contains the response from method CloudsClient.BeginUpdate.
type CloudsClientUpdateResponse struct {
	// The Clouds resource definition.
	Cloud
}

// InventoryItemsClientCreateResponse contains the response from method InventoryItemsClient.Create.
type InventoryItemsClientCreateResponse struct {
	// Defines the inventory item.
	InventoryItem
}

// InventoryItemsClientDeleteResponse contains the response from method InventoryItemsClient.Delete.
type InventoryItemsClientDeleteResponse struct {
	// placeholder for future response values
}

// InventoryItemsClientGetResponse contains the response from method InventoryItemsClient.Get.
type InventoryItemsClientGetResponse struct {
	// Defines the inventory item.
	InventoryItem
}

// InventoryItemsClientListByVMMServerResponse contains the response from method InventoryItemsClient.NewListByVMMServerPager.
type InventoryItemsClientListByVMMServerResponse struct {
	// List of InventoryItems.
	InventoryItemsList
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// Results of the request to list operations.
	ResourceProviderOperationList
}

// VirtualMachineTemplatesClientCreateOrUpdateResponse contains the response from method VirtualMachineTemplatesClient.BeginCreateOrUpdate.
type VirtualMachineTemplatesClientCreateOrUpdateResponse struct {
	// The VirtualMachineTemplates resource definition.
	VirtualMachineTemplate
}

// VirtualMachineTemplatesClientDeleteResponse contains the response from method VirtualMachineTemplatesClient.BeginDelete.
type VirtualMachineTemplatesClientDeleteResponse struct {
	// placeholder for future response values
}

// VirtualMachineTemplatesClientGetResponse contains the response from method VirtualMachineTemplatesClient.Get.
type VirtualMachineTemplatesClientGetResponse struct {
	// The VirtualMachineTemplates resource definition.
	VirtualMachineTemplate
}

// VirtualMachineTemplatesClientListByResourceGroupResponse contains the response from method VirtualMachineTemplatesClient.NewListByResourceGroupPager.
type VirtualMachineTemplatesClientListByResourceGroupResponse struct {
	// List of VirtualMachineTemplates.
	VirtualMachineTemplateListResult
}

// VirtualMachineTemplatesClientListBySubscriptionResponse contains the response from method VirtualMachineTemplatesClient.NewListBySubscriptionPager.
type VirtualMachineTemplatesClientListBySubscriptionResponse struct {
	// List of VirtualMachineTemplates.
	VirtualMachineTemplateListResult
}

// VirtualMachineTemplatesClientUpdateResponse contains the response from method VirtualMachineTemplatesClient.BeginUpdate.
type VirtualMachineTemplatesClientUpdateResponse struct {
	// The VirtualMachineTemplates resource definition.
	VirtualMachineTemplate
}

// VirtualMachinesClientCreateCheckpointResponse contains the response from method VirtualMachinesClient.BeginCreateCheckpoint.
type VirtualMachinesClientCreateCheckpointResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientCreateOrUpdateResponse contains the response from method VirtualMachinesClient.BeginCreateOrUpdate.
type VirtualMachinesClientCreateOrUpdateResponse struct {
	// The VirtualMachines resource definition.
	VirtualMachine
}

// VirtualMachinesClientDeleteCheckpointResponse contains the response from method VirtualMachinesClient.BeginDeleteCheckpoint.
type VirtualMachinesClientDeleteCheckpointResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientDeleteResponse contains the response from method VirtualMachinesClient.BeginDelete.
type VirtualMachinesClientDeleteResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientGetResponse contains the response from method VirtualMachinesClient.Get.
type VirtualMachinesClientGetResponse struct {
	// The VirtualMachines resource definition.
	VirtualMachine
}

// VirtualMachinesClientListByResourceGroupResponse contains the response from method VirtualMachinesClient.NewListByResourceGroupPager.
type VirtualMachinesClientListByResourceGroupResponse struct {
	// List of VirtualMachines.
	VirtualMachineListResult
}

// VirtualMachinesClientListBySubscriptionResponse contains the response from method VirtualMachinesClient.NewListBySubscriptionPager.
type VirtualMachinesClientListBySubscriptionResponse struct {
	// List of VirtualMachines.
	VirtualMachineListResult
}

// VirtualMachinesClientRestartResponse contains the response from method VirtualMachinesClient.BeginRestart.
type VirtualMachinesClientRestartResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientRestoreCheckpointResponse contains the response from method VirtualMachinesClient.BeginRestoreCheckpoint.
type VirtualMachinesClientRestoreCheckpointResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientStartResponse contains the response from method VirtualMachinesClient.BeginStart.
type VirtualMachinesClientStartResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientStopResponse contains the response from method VirtualMachinesClient.BeginStop.
type VirtualMachinesClientStopResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientUpdateResponse contains the response from method VirtualMachinesClient.BeginUpdate.
type VirtualMachinesClientUpdateResponse struct {
	// The VirtualMachines resource definition.
	VirtualMachine
}

// VirtualNetworksClientCreateOrUpdateResponse contains the response from method VirtualNetworksClient.BeginCreateOrUpdate.
type VirtualNetworksClientCreateOrUpdateResponse struct {
	// The VirtualNetworks resource definition.
	VirtualNetwork
}

// VirtualNetworksClientDeleteResponse contains the response from method VirtualNetworksClient.BeginDelete.
type VirtualNetworksClientDeleteResponse struct {
	// placeholder for future response values
}

// VirtualNetworksClientGetResponse contains the response from method VirtualNetworksClient.Get.
type VirtualNetworksClientGetResponse struct {
	// The VirtualNetworks resource definition.
	VirtualNetwork
}

// VirtualNetworksClientListByResourceGroupResponse contains the response from method VirtualNetworksClient.NewListByResourceGroupPager.
type VirtualNetworksClientListByResourceGroupResponse struct {
	// List of VirtualNetworks.
	VirtualNetworkListResult
}

// VirtualNetworksClientListBySubscriptionResponse contains the response from method VirtualNetworksClient.NewListBySubscriptionPager.
type VirtualNetworksClientListBySubscriptionResponse struct {
	// List of VirtualNetworks.
	VirtualNetworkListResult
}

// VirtualNetworksClientUpdateResponse contains the response from method VirtualNetworksClient.BeginUpdate.
type VirtualNetworksClientUpdateResponse struct {
	// The VirtualNetworks resource definition.
	VirtualNetwork
}

// VmmServersClientCreateOrUpdateResponse contains the response from method VmmServersClient.BeginCreateOrUpdate.
type VmmServersClientCreateOrUpdateResponse struct {
	// The VmmServers resource definition.
	VMMServer
}

// VmmServersClientDeleteResponse contains the response from method VmmServersClient.BeginDelete.
type VmmServersClientDeleteResponse struct {
	// placeholder for future response values
}

// VmmServersClientGetResponse contains the response from method VmmServersClient.Get.
type VmmServersClientGetResponse struct {
	// The VmmServers resource definition.
	VMMServer
}

// VmmServersClientListByResourceGroupResponse contains the response from method VmmServersClient.NewListByResourceGroupPager.
type VmmServersClientListByResourceGroupResponse struct {
	// List of VmmServers.
	VMMServerListResult
}

// VmmServersClientListBySubscriptionResponse contains the response from method VmmServersClient.NewListBySubscriptionPager.
type VmmServersClientListBySubscriptionResponse struct {
	// List of VmmServers.
	VMMServerListResult
}

// VmmServersClientUpdateResponse contains the response from method VmmServersClient.BeginUpdate.
type VmmServersClientUpdateResponse struct {
	// The VmmServers resource definition.
	VMMServer
}


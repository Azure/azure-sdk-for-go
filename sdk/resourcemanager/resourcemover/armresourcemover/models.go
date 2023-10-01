//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armresourcemover

import "time"

// AffectedMoveResource - The RP custom operation error info.
type AffectedMoveResource struct {
	// READ-ONLY; The affected move resource id.
	ID *string

	// READ-ONLY; The affected move resources.
	MoveResources []*AffectedMoveResource

	// READ-ONLY; The affected move resource source id.
	SourceID *string
}

// AutomaticResolutionProperties - Defines the properties for automatic resolution.
type AutomaticResolutionProperties struct {
	// Gets the MoveResource ARM ID of the dependent resource if the resolution type is Automatic.
	MoveResourceID *string
}

// AvailabilitySetResourceSettings - Gets or sets the availability set resource settings.
type AvailabilitySetResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the target fault domain.
	FaultDomain *int32

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string

	// Gets or sets the target update domain.
	UpdateDomain *int32
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type AvailabilitySetResourceSettings.
func (a *AvailabilitySetResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            a.ResourceType,
		TargetResourceGroupName: a.TargetResourceGroupName,
		TargetResourceName:      a.TargetResourceName,
	}
}

// AzureResourceReference - Defines reference to an Azure resource.
type AzureResourceReference struct {
	// REQUIRED; Gets the ARM resource ID of the tracked resource being referenced.
	SourceArmResourceID *string
}

// BulkRemoveRequest - Defines the request body for bulk remove of move resources operation.
type BulkRemoveRequest struct {
	// Defines the move resource input type.
	MoveResourceInputType *MoveResourceInputType

	// Gets or sets the list of resource Id's, by default it accepts move resource id's unless the input type is switched via
	// moveResourceInputType property.
	MoveResources []*string

	// Gets or sets a value indicating whether the operation needs to only run pre-requisite.
	ValidateOnly *bool
}

// CommitRequest - Defines the request body for commit operation.
type CommitRequest struct {
	// REQUIRED; Gets or sets the list of resource Id's, by default it accepts move resource id's unless the input type is switched
	// via moveResourceInputType property.
	MoveResources []*string

	// Defines the move resource input type.
	MoveResourceInputType *MoveResourceInputType

	// Gets or sets a value indicating whether the operation needs to only run pre-requisite.
	ValidateOnly *bool
}

// DiscardRequest - Defines the request body for discard operation.
type DiscardRequest struct {
	// REQUIRED; Gets or sets the list of resource Id's, by default it accepts move resource id's unless the input type is switched
	// via moveResourceInputType property.
	MoveResources []*string

	// Defines the move resource input type.
	MoveResourceInputType *MoveResourceInputType

	// Gets or sets a value indicating whether the operation needs to only run pre-requisite.
	ValidateOnly *bool
}

// DiskEncryptionSetResourceSettings - Defines the disk encryption set resource settings.
type DiskEncryptionSetResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type DiskEncryptionSetResourceSettings.
func (d *DiskEncryptionSetResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            d.ResourceType,
		TargetResourceGroupName: d.TargetResourceGroupName,
		TargetResourceName:      d.TargetResourceName,
	}
}

// Display - Contains the localized display information for this particular operation / action. These value will be used by
// several clients for (1) custom role definitions for RBAC; (2) complex query filters for
// the event service; and (3) audit history / records for management operations.
type Display struct {
	// Gets or sets the description. The localized friendly description for the operation, as it should be shown to the user.
	// It should be thorough, yet concise – it will be used in tool tips and detailed
	// views. Prescriptive guidance for namespace: Read any 'display.provider' resource Create or Update any 'display.provider'
	// resource Delete any 'display.provider' resource Perform any other action on any
	// 'display.provider' resource Prescriptive guidance for namespace: Read any 'display.resource' Create or Update any 'display.resource'
	// Delete any 'display.resource' 'ActionName' any 'display.resources'.
	Description *string

	// Gets or sets the operation. The localized friendly name for the operation, as it should be shown to the user. It should
	// be concise (to fit in drop downs) but clear (i.e. self-documenting). It should
	// use Title Casing. Prescriptive guidance: Read Create or Update Delete 'ActionName'.
	Operation *string

	// Gets or sets the provider. The localized friendly form of the resource provider name – it is expected to also include the
	// publisher/company responsible. It should use Title Casing and begin with
	// "Microsoft" for 1st party services. e.g. "Microsoft Monitoring Insights" or "Microsoft Compute.".
	Provider *string

	// Gets or sets the resource. The localized friendly form of the resource related to this action/operation – it should match
	// the public documentation for the resource provider. It should use Title
	// Casing. This value should be unique for a particular URL type (e.g. nested types should notreuse their parent’s display.resource
	// field) e.g. "Virtual Machines" or "Scheduler Job Collections", or
	// "Virtual Machine VM Sizes" or "Scheduler Jobs".
	Resource *string
}

// Identity - Defines the MSI properties of the Move Collection.
type Identity struct {
	// Gets or sets the principal id.
	PrincipalID *string

	// Gets or sets the tenant id.
	TenantID *string

	// The type of identity used for the resource mover service.
	Type *ResourceIdentityType
}

// JobStatus - Defines the job status.
type JobStatus struct {
	// READ-ONLY; Defines the job name.
	JobName *JobName

	// READ-ONLY; Gets or sets the monitoring job percentage.
	JobProgress *string
}

// KeyVaultResourceSettings - Defines the key vault resource settings.
type KeyVaultResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type KeyVaultResourceSettings.
func (k *KeyVaultResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            k.ResourceType,
		TargetResourceGroupName: k.TargetResourceGroupName,
		TargetResourceName:      k.TargetResourceName,
	}
}

// LBBackendAddressPoolResourceSettings - Defines load balancer backend address pool properties.
type LBBackendAddressPoolResourceSettings struct {
	// Gets or sets the backend address pool name.
	Name *string
}

// LBFrontendIPConfigurationResourceSettings - Defines load balancer frontend IP configuration properties.
type LBFrontendIPConfigurationResourceSettings struct {
	// Gets or sets the frontend IP configuration name.
	Name *string

	// Gets or sets the IP address of the Load Balancer.This is only specified if a specific private IP address shall be allocated
	// from the subnet specified in subnetRef.
	PrivateIPAddress *string

	// Gets or sets PrivateIP allocation method (Static/Dynamic).
	PrivateIPAllocationMethod *string

	// Defines reference to subnet.
	Subnet *SubnetReference

	// Gets or sets the csv list of zones.
	Zones *string
}

// LoadBalancerBackendAddressPoolReference - Defines reference to load balancer backend address pools.
type LoadBalancerBackendAddressPoolReference struct {
	// REQUIRED; Gets the ARM resource ID of the tracked resource being referenced.
	SourceArmResourceID *string

	// Gets the name of the proxy resource on the target side.
	Name *string
}

// LoadBalancerNatRuleReference - Defines reference to load balancer NAT rules.
type LoadBalancerNatRuleReference struct {
	// REQUIRED; Gets the ARM resource ID of the tracked resource being referenced.
	SourceArmResourceID *string

	// Gets the name of the proxy resource on the target side.
	Name *string
}

// LoadBalancerResourceSettings - Defines the load balancer resource settings.
type LoadBalancerResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the backend address pools of the load balancer.
	BackendAddressPools []*LBBackendAddressPoolResourceSettings

	// Gets or sets the frontend IP configurations of the load balancer.
	FrontendIPConfigurations []*LBFrontendIPConfigurationResourceSettings

	// Gets or sets load balancer sku (Basic/Standard).
	SKU *string

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string

	// Gets or sets the csv list of zones common for all frontend IP configurations. Note this is given precedence only if frontend
	// IP configurations settings are not present.
	Zones *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type LoadBalancerResourceSettings.
func (l *LoadBalancerResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            l.ResourceType,
		TargetResourceGroupName: l.TargetResourceGroupName,
		TargetResourceName:      l.TargetResourceName,
	}
}

// ManualResolutionProperties - Defines the properties for manual resolution.
type ManualResolutionProperties struct {
	// Gets or sets the target resource ARM ID of the dependent resource if the resource type is Manual.
	TargetID *string
}

// MoveCollection - Define the move collection.
type MoveCollection struct {
	// Defines the MSI properties of the Move Collection.
	Identity *Identity

	// The geo-location where the resource lives.
	Location *string

	// Defines the move collection properties.
	Properties *MoveCollectionProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; The etag of the resource.
	Etag *string

	// READ-ONLY; Fully qualified resource Id for the resource.
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData

	// READ-ONLY; The type of the resource.
	Type *string
}

// MoveCollectionProperties - Defines the move collection properties.
type MoveCollectionProperties struct {
	// Gets or sets the move region which indicates the region where the VM Regional to Zonal move will be conducted.
	MoveRegion *string

	// Defines the MoveType.
	MoveType *MoveType

	// Gets or sets the source region.
	SourceRegion *string

	// Gets or sets the target region.
	TargetRegion *string

	// Gets or sets the version of move collection.
	Version *string

	// READ-ONLY; Defines the move collection errors.
	Errors *MoveCollectionPropertiesErrors

	// READ-ONLY; Defines the provisioning states.
	ProvisioningState *ProvisioningState
}

// MoveCollectionPropertiesErrors - Defines the move collection errors.
type MoveCollectionPropertiesErrors struct {
	// The move resource error body.
	Properties *MoveResourceErrorBody
}

// MoveCollectionResultList - Defines the collection of move collections.
type MoveCollectionResultList struct {
	// Gets the value of next link.
	NextLink *string

	// Gets the list of move collections.
	Value []*MoveCollection
}

// MoveErrorInfo - The move custom error info.
type MoveErrorInfo struct {
	// READ-ONLY; The affected move resources.
	MoveResources []*AffectedMoveResource
}

// MoveResource - Defines the move resource.
type MoveResource struct {
	// Defines the move resource properties.
	Properties *MoveResourceProperties

	// READ-ONLY; Fully qualified resource Id for the resource.
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData

	// READ-ONLY; The type of the resource.
	Type *string
}

// MoveResourceCollection - Defines the collection of move resources.
type MoveResourceCollection struct {
	// Gets the value of next link.
	NextLink *string

	// Gets or sets the list of summary items and the field on which summary is done.
	SummaryCollection *SummaryCollection

	// Gets the list of move resources.
	Value []*MoveResource

	// READ-ONLY; Gets the total count.
	TotalCount *int64
}

// MoveResourceDependency - Defines the dependency of the move resource.
type MoveResourceDependency struct {
	// Defines the properties for automatic resolution.
	AutomaticResolution *AutomaticResolutionProperties

	// Defines the dependency type.
	DependencyType *DependencyType

	// Gets the source ARM ID of the dependent resource.
	ID *string

	// Gets or sets a value indicating whether the dependency is optional.
	IsOptional *string

	// Defines the properties for manual resolution.
	ManualResolution *ManualResolutionProperties

	// Gets the dependency resolution status.
	ResolutionStatus *string

	// Defines the resolution type.
	ResolutionType *ResolutionType
}

// MoveResourceDependencyOverride - Defines the dependency override of the move resource.
type MoveResourceDependencyOverride struct {
	// Gets or sets the ARM ID of the dependent resource.
	ID *string

	// Gets or sets the resource ARM id of either the MoveResource or the resource ARM ID of the dependent resource.
	TargetID *string
}

// MoveResourceError - An error response from the azure resource mover service.
type MoveResourceError struct {
	// The move resource error body.
	Properties *MoveResourceErrorBody
}

// MoveResourceErrorBody - An error response from the Azure Migrate service.
type MoveResourceErrorBody struct {
	// READ-ONLY; An identifier for the error. Codes are invariant and are intended to be consumed programmatically.
	Code *string

	// READ-ONLY; A list of additional details about the error.
	Details []*MoveResourceErrorBody

	// READ-ONLY; A message describing the error, intended to be suitable for display in a user interface.
	Message *string

	// READ-ONLY; The target of the particular error. For example, the name of the property in error.
	Target *string
}

// MoveResourceFilter - Move resource filter.
type MoveResourceFilter struct {
	Properties *MoveResourceFilterProperties
}

type MoveResourceFilterProperties struct {
	// The provisioning state.
	ProvisioningState *string
}

// MoveResourceProperties - Defines the move resource properties.
type MoveResourceProperties struct {
	// REQUIRED; Gets or sets the Source ARM Id of the resource.
	SourceID *string

	// Gets or sets the move resource dependencies overrides.
	DependsOnOverrides []*MoveResourceDependencyOverride

	// Gets or sets the existing target ARM Id of the resource.
	ExistingTargetID *string

	// Gets or sets the resource settings.
	ResourceSettings ResourceSettingsClassification

	// READ-ONLY; Gets or sets the move resource dependencies.
	DependsOn []*MoveResourceDependency

	// READ-ONLY; Defines the move resource errors.
	Errors *MoveResourcePropertiesErrors

	// READ-ONLY; Gets a value indicating whether the resolve action is required over the move collection.
	IsResolveRequired *bool

	// READ-ONLY; Defines the move resource status.
	MoveStatus *MoveResourcePropertiesMoveStatus

	// READ-ONLY; Defines the provisioning states.
	ProvisioningState *ProvisioningState

	// READ-ONLY; Gets or sets the source resource settings.
	SourceResourceSettings ResourceSettingsClassification

	// READ-ONLY; Gets or sets the Target ARM Id of the resource.
	TargetID *string
}

// MoveResourcePropertiesErrors - Defines the move resource errors.
type MoveResourcePropertiesErrors struct {
	// The move resource error body.
	Properties *MoveResourceErrorBody
}

// MoveResourcePropertiesMoveStatus - Defines the move resource status.
type MoveResourcePropertiesMoveStatus struct {
	// An error response from the azure resource mover service.
	Errors *MoveResourceError

	// Defines the job status.
	JobStatus *JobStatus

	// READ-ONLY; Defines the MoveResource states.
	MoveState *MoveState
}

// MoveResourceStatus - Defines the move resource status.
type MoveResourceStatus struct {
	// An error response from the azure resource mover service.
	Errors *MoveResourceError

	// Defines the job status.
	JobStatus *JobStatus

	// READ-ONLY; Defines the MoveResource states.
	MoveState *MoveState
}

// NetworkInterfaceResourceSettings - Defines the network interface resource settings.
type NetworkInterfaceResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets a value indicating whether accelerated networking is enabled.
	EnableAcceleratedNetworking *bool

	// Gets or sets the IP configurations of the NIC.
	IPConfigurations []*NicIPConfigurationResourceSettings

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type NetworkInterfaceResourceSettings.
func (n *NetworkInterfaceResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            n.ResourceType,
		TargetResourceGroupName: n.TargetResourceGroupName,
		TargetResourceName:      n.TargetResourceName,
	}
}

// NetworkSecurityGroupResourceSettings - Defines the NSG resource settings.
type NetworkSecurityGroupResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets Security rules of network security group.
	SecurityRules []*NsgSecurityRule

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type NetworkSecurityGroupResourceSettings.
func (n *NetworkSecurityGroupResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            n.ResourceType,
		TargetResourceGroupName: n.TargetResourceGroupName,
		TargetResourceName:      n.TargetResourceName,
	}
}

// NicIPConfigurationResourceSettings - Defines NIC IP configuration properties.
type NicIPConfigurationResourceSettings struct {
	// Gets or sets the references of the load balancer backend address pools.
	LoadBalancerBackendAddressPools []*LoadBalancerBackendAddressPoolReference

	// Gets or sets the references of the load balancer NAT rules.
	LoadBalancerNatRules []*LoadBalancerNatRuleReference

	// Gets or sets the IP configuration name.
	Name *string

	// Gets or sets a value indicating whether this IP configuration is the primary.
	Primary *bool

	// Gets or sets the private IP address of the network interface IP Configuration.
	PrivateIPAddress *string

	// Gets or sets the private IP address allocation method.
	PrivateIPAllocationMethod *string

	// Defines reference to a public IP.
	PublicIP *PublicIPReference

	// Defines reference to subnet.
	Subnet *SubnetReference
}

// NsgReference - Defines reference to NSG.
type NsgReference struct {
	// REQUIRED; Gets the ARM resource ID of the tracked resource being referenced.
	SourceArmResourceID *string
}

// NsgSecurityRule - Security Rule data model for Network Security Groups.
type NsgSecurityRule struct {
	// Gets or sets whether network traffic is allowed or denied. Possible values are “Allow” and “Deny”.
	Access *string

	// Gets or sets a description for this rule. Restricted to 140 chars.
	Description *string

	// Gets or sets destination address prefix. CIDR or source IP range. A “*” can also be used to match all source IPs. Default
	// tags such as ‘VirtualNetwork’, ‘AzureLoadBalancer’ and ‘Internet’ can also be
	// used.
	DestinationAddressPrefix *string

	// Gets or sets Destination Port or Range. Integer or range between 0 and 65535. A “*” can also be used to match all ports.
	DestinationPortRange *string

	// Gets or sets the direction of the rule.InBound or Outbound. The direction specifies if rule will be evaluated on incoming
	// or outgoing traffic.
	Direction *string

	// Gets or sets the Security rule name.
	Name *string

	// Gets or sets the priority of the rule. The value can be between 100 and 4096. The priority number must be unique for each
	// rule in the collection. The lower the priority number, the higher the priority
	// of the rule.
	Priority *int32

	// Gets or sets Network protocol this rule applies to. Can be Tcp, Udp or All(*).
	Protocol *string

	// Gets or sets source address prefix. CIDR or source IP range. A “*” can also be used to match all source IPs. Default tags
	// such as ‘VirtualNetwork’, ‘AzureLoadBalancer’ and ‘Internet’ can also be used.
	// If this is an ingress rule, specifies where network traffic originates from.
	SourceAddressPrefix *string

	// Gets or sets Source Port or Range. Integer or range between 0 and
	// 65535. A “*” can also be used to match all ports.
	SourcePortRange *string
}

// OperationErrorAdditionalInfo - The operation error info.
type OperationErrorAdditionalInfo struct {
	// READ-ONLY; The operation error info.
	Info *MoveErrorInfo

	// READ-ONLY; The error type.
	Type *string
}

// OperationStatus - Operation status REST resource.
type OperationStatus struct {
	// READ-ONLY; End time.
	EndTime *string

	// READ-ONLY; Error stating all error details for the operation.
	Error *OperationStatusError

	// READ-ONLY; Resource Id.
	ID *string

	// READ-ONLY; Operation name.
	Name *string

	// READ-ONLY; Custom data.
	Properties any

	// READ-ONLY; Start time.
	StartTime *string

	// READ-ONLY; Status of the operation. ARM expects the terminal status to be one of Succeeded/ Failed/ Canceled. All other
	// values imply that the operation is still running.
	Status *string
}

// OperationStatusError - Class for operation status errors.
type OperationStatusError struct {
	// READ-ONLY; The additional info.
	AdditionalInfo []*OperationErrorAdditionalInfo

	// READ-ONLY; The error code.
	Code *string

	// READ-ONLY; The error details.
	Details []*OperationStatusError

	// READ-ONLY; The error message.
	Message *string
}

// OperationsDiscovery - Operations discovery class.
type OperationsDiscovery struct {
	// Contains the localized display information for this particular operation / action. These value will be used by several
	// clients for (1) custom role definitions for RBAC; (2) complex query filters for
	// the event service; and (3) audit history / records for management operations.
	Display *Display

	// Indicates whether the operation is a data action
	IsDataAction *bool

	// Gets or sets Name of the API. The name of the operation being performed on this particular object. It should match the
	// action name that appears in RBAC / the event service. Examples of operations
	// include:
	// * Microsoft.Compute/virtualMachine/capture/action
	// * Microsoft.Compute/virtualMachine/restart/action
	// * Microsoft.Compute/virtualMachine/write
	// * Microsoft.Compute/virtualMachine/read
	// * Microsoft.Compute/virtualMachine/delete Each action should include, in order: (1) Resource Provider Namespace (2) Type
	// hierarchy for which the action applies (e.g. server/databases for a SQL Azure
	// database) (3) Read, Write, Action or Delete indicating which type applies. If it is a PUT/PATCH on a collection or named
	// value, Write should be used. If it is a GET, Read should be used. If it is a
	// DELETE, Delete should be used. If it is a POST, Action should be used. As a note: all resource providers would need to
	// include the "{Resource Provider Namespace}/register/action" operation in their
	// response. This API is used to register for their service, and should include details about the operation (e.g. a localized
	// name for the resource provider + any special considerations like PII
	// release).
	Name *string

	// Gets or sets Origin. The intended executor of the operation; governs the display of the operation in the RBAC UX and the
	// audit logs UX. Default value is "user,system".
	Origin *string

	// ClientDiscovery properties.
	Properties any
}

// OperationsDiscoveryCollection - Collection of ClientDiscovery details.
type OperationsDiscoveryCollection struct {
	// Gets or sets the value of next link.
	NextLink *string

	// Gets or sets the ClientDiscovery details.
	Value []*OperationsDiscovery
}

// PrepareRequest - Defines the request body for initiate prepare operation.
type PrepareRequest struct {
	// REQUIRED; Gets or sets the list of resource Id's, by default it accepts move resource id's unless the input type is switched
	// via moveResourceInputType property.
	MoveResources []*string

	// Defines the move resource input type.
	MoveResourceInputType *MoveResourceInputType

	// Gets or sets a value indicating whether the operation needs to only run pre-requisite.
	ValidateOnly *bool
}

// ProxyResourceReference - Defines reference to a proxy resource.
type ProxyResourceReference struct {
	// REQUIRED; Gets the ARM resource ID of the tracked resource being referenced.
	SourceArmResourceID *string

	// Gets the name of the proxy resource on the target side.
	Name *string
}

// PublicIPAddressResourceSettings - Defines the public IP address resource settings.
type PublicIPAddressResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the domain name label.
	DomainNameLabel *string

	// Gets or sets the fully qualified domain name.
	Fqdn *string

	// Gets or sets public IP allocation method.
	PublicIPAllocationMethod *string

	// Gets or sets public IP sku.
	SKU *string

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string

	// Gets or sets public IP zones.
	Zones *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type PublicIPAddressResourceSettings.
func (p *PublicIPAddressResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            p.ResourceType,
		TargetResourceGroupName: p.TargetResourceGroupName,
		TargetResourceName:      p.TargetResourceName,
	}
}

// PublicIPReference - Defines reference to a public IP.
type PublicIPReference struct {
	// REQUIRED; Gets the ARM resource ID of the tracked resource being referenced.
	SourceArmResourceID *string
}

// RequiredForResourcesCollection - Required for resources collection.
type RequiredForResourcesCollection struct {
	// Gets or sets the list of source Ids for which the input resource is required.
	SourceIDs []*string
}

// ResourceGroupResourceSettings - Defines the resource group resource settings.
type ResourceGroupResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type ResourceGroupResourceSettings.
func (r *ResourceGroupResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            r.ResourceType,
		TargetResourceGroupName: r.TargetResourceGroupName,
		TargetResourceName:      r.TargetResourceName,
	}
}

// ResourceMoveRequest - Defines the request body for resource move operation.
type ResourceMoveRequest struct {
	// REQUIRED; Gets or sets the list of resource Id's, by default it accepts move resource id's unless the input type is switched
	// via moveResourceInputType property.
	MoveResources []*string

	// Defines the move resource input type.
	MoveResourceInputType *MoveResourceInputType

	// Gets or sets a value indicating whether the operation needs to only run pre-requisite.
	ValidateOnly *bool
}

// ResourceSettings - Gets or sets the resource settings.
type ResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type ResourceSettings.
func (r *ResourceSettings) GetResourceSettings() *ResourceSettings { return r }

// SQLDatabaseResourceSettings - Defines the Sql Database resource settings.
type SQLDatabaseResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string

	// Defines the zone redundant resource setting.
	ZoneRedundant *ZoneRedundant
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type SQLDatabaseResourceSettings.
func (s *SQLDatabaseResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            s.ResourceType,
		TargetResourceGroupName: s.TargetResourceGroupName,
		TargetResourceName:      s.TargetResourceName,
	}
}

// SQLElasticPoolResourceSettings - Defines the Sql ElasticPool resource settings.
type SQLElasticPoolResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string

	// Defines the zone redundant resource setting.
	ZoneRedundant *ZoneRedundant
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type SQLElasticPoolResourceSettings.
func (s *SQLElasticPoolResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            s.ResourceType,
		TargetResourceGroupName: s.TargetResourceGroupName,
		TargetResourceName:      s.TargetResourceName,
	}
}

// SQLServerResourceSettings - Defines the SQL Server resource settings.
type SQLServerResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type SQLServerResourceSettings.
func (s *SQLServerResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            s.ResourceType,
		TargetResourceGroupName: s.TargetResourceGroupName,
		TargetResourceName:      s.TargetResourceName,
	}
}

// SubnetReference - Defines reference to subnet.
type SubnetReference struct {
	// REQUIRED; Gets the ARM resource ID of the tracked resource being referenced.
	SourceArmResourceID *string

	// Gets the name of the proxy resource on the target side.
	Name *string
}

// SubnetResourceSettings - Defines the virtual network subnets resource settings.
type SubnetResourceSettings struct {
	// Gets or sets address prefix for the subnet.
	AddressPrefix *string

	// Gets or sets the Subnet name.
	Name *string

	// Defines reference to NSG.
	NetworkSecurityGroup *NsgReference
}

// Summary item.
type Summary struct {
	// Gets the count.
	Count *int32

	// Gets the item.
	Item *string
}

// SummaryCollection - Summary Collection.
type SummaryCollection struct {
	// Gets or sets the field name on which summary is done.
	FieldName *string

	// Gets or sets the list of summary items.
	Summary []*Summary
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time

	// The identity that created the resource.
	CreatedBy *string

	// The type of identity that created the resource.
	CreatedByType *CreatedByType

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time

	// The identity that last modified the resource.
	LastModifiedBy *string

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType
}

// UnresolvedDependenciesFilter - Unresolved dependencies contract.
type UnresolvedDependenciesFilter struct {
	Properties *UnresolvedDependenciesFilterProperties
}

type UnresolvedDependenciesFilterProperties struct {
	// The count of the resource.
	Count *int32
}

// UnresolvedDependency - Unresolved dependency.
type UnresolvedDependency struct {
	// Gets or sets the count.
	Count *int32

	// Gets or sets the arm id of the dependency.
	ID *string
}

// UnresolvedDependencyCollection - Unresolved dependency collection.
type UnresolvedDependencyCollection struct {
	// Gets or sets the value of next link.
	NextLink *string

	// Gets or sets the list of unresolved dependencies.
	Value []*UnresolvedDependency

	// READ-ONLY; Gets or sets the list of summary items and the field on which summary is done.
	SummaryCollection *SummaryCollection

	// READ-ONLY; Gets the total count.
	TotalCount *int64
}

// UpdateMoveCollectionRequest - Defines the request body for updating move collection.
type UpdateMoveCollectionRequest struct {
	// Defines the MSI properties of the Move Collection.
	Identity *Identity

	// Gets or sets the Resource tags.
	Tags map[string]*string
}

// VirtualMachineResourceSettings - Gets or sets the virtual machine resource settings.
type VirtualMachineResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target availability set id for virtual machines not in an availability set at source.
	TargetAvailabilitySetID *string

	// Gets or sets the target availability zone.
	TargetAvailabilityZone *TargetAvailabilityZone

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string

	// Gets or sets the target virtual machine size.
	TargetVMSize *string

	// Gets or sets user-managed identities
	UserManagedIdentities []*string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type VirtualMachineResourceSettings.
func (v *VirtualMachineResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            v.ResourceType,
		TargetResourceGroupName: v.TargetResourceGroupName,
		TargetResourceName:      v.TargetResourceName,
	}
}

// VirtualNetworkResourceSettings - Defines the virtual network resource settings.
type VirtualNetworkResourceSettings struct {
	// REQUIRED; The resource type. For example, the value can be Microsoft.Compute/virtualMachines.
	ResourceType *string

	// Gets or sets the address prefixes for the virtual network.
	AddressSpace []*string

	// Gets or sets DHCPOptions that contains an array of DNS servers available to VMs deployed in the virtual network.
	DNSServers []*string

	// Gets or sets a value indicating whether gets or sets whether the DDOS protection should be switched on.
	EnableDdosProtection *bool

	// Gets or sets List of subnets in a VirtualNetwork.
	Subnets []*SubnetResourceSettings

	// Gets or sets the Resource tags.
	Tags map[string]*string

	// Gets or sets the target resource group name.
	TargetResourceGroupName *string

	// Gets or sets the target Resource name.
	TargetResourceName *string
}

// GetResourceSettings implements the ResourceSettingsClassification interface for type VirtualNetworkResourceSettings.
func (v *VirtualNetworkResourceSettings) GetResourceSettings() *ResourceSettings {
	return &ResourceSettings{
		ResourceType:            v.ResourceType,
		TargetResourceGroupName: v.TargetResourceGroupName,
		TargetResourceName:      v.TargetResourceName,
	}
}

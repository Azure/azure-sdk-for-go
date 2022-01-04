//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstoragepool

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
	"time"
)

// ACL - Access Control List (ACL) for an iSCSI Target; defines LUN masking policy
type ACL struct {
	// REQUIRED; iSCSI initiator IQN (iSCSI Qualified Name); example: "iqn.2005-03.org.iscsi:client".
	InitiatorIqn *string `json:"initiatorIqn,omitempty"`

	// REQUIRED; List of LUN names mapped to the ACL.
	MappedLuns []*string `json:"mappedLuns,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type ACL.
func (a ACL) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "initiatorIqn", a.InitiatorIqn)
	populate(objectMap, "mappedLuns", a.MappedLuns)
	return json.Marshal(objectMap)
}

// Disk - Azure Managed Disk to attach to the Disk Pool.
type Disk struct {
	// REQUIRED; Unique Azure Resource ID of the Managed Disk.
	ID *string `json:"id,omitempty"`
}

// DiskPool - Response for Disk Pool request.
type DiskPool struct {
	TrackedResource
	// REQUIRED; Properties of Disk Pool.
	Properties *DiskPoolProperties `json:"properties,omitempty"`

	// Determines the SKU of the Disk pool
	SKU *SKU `json:"sku,omitempty"`

	// READ-ONLY; Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedBy *string `json:"managedBy,omitempty" azure:"ro"`

	// READ-ONLY; List of Azure resource ids that manage this resource.
	ManagedByExtended []*string `json:"managedByExtended,omitempty" azure:"ro"`

	// READ-ONLY; Resource metadata required by ARM RPC
	SystemData *SystemMetadata `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPool.
func (d DiskPool) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	d.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "managedBy", d.ManagedBy)
	populate(objectMap, "managedByExtended", d.ManagedByExtended)
	populate(objectMap, "properties", d.Properties)
	populate(objectMap, "sku", d.SKU)
	populate(objectMap, "systemData", d.SystemData)
	return json.Marshal(objectMap)
}

// DiskPoolCreate - Request payload for create or update Disk Pool request.
type DiskPoolCreate struct {
	// REQUIRED; The geo-location where the resource lives.
	Location *string `json:"location,omitempty"`

	// REQUIRED; Properties for Disk Pool create request.
	Properties *DiskPoolCreateProperties `json:"properties,omitempty"`

	// REQUIRED; Determines the SKU of the Disk Pool
	SKU *SKU `json:"sku,omitempty"`

	// Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedBy *string `json:"managedBy,omitempty"`

	// List of Azure resource ids that manage this resource.
	ManagedByExtended []*string `json:"managedByExtended,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. Ex- Microsoft.Compute/virtualMachines or Microsoft.Storage/storageAccounts.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPoolCreate.
func (d DiskPoolCreate) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", d.ID)
	populate(objectMap, "location", d.Location)
	populate(objectMap, "managedBy", d.ManagedBy)
	populate(objectMap, "managedByExtended", d.ManagedByExtended)
	populate(objectMap, "name", d.Name)
	populate(objectMap, "properties", d.Properties)
	populate(objectMap, "sku", d.SKU)
	populate(objectMap, "tags", d.Tags)
	populate(objectMap, "type", d.Type)
	return json.Marshal(objectMap)
}

// DiskPoolCreateProperties - Properties for Disk Pool create or update request.
type DiskPoolCreateProperties struct {
	// REQUIRED; Azure Resource ID of a Subnet for the Disk Pool.
	SubnetID *string `json:"subnetId,omitempty"`

	// List of additional capabilities for a Disk Pool.
	AdditionalCapabilities []*string `json:"additionalCapabilities,omitempty"`

	// Logical zone for Disk Pool resource; example: ["1"].
	AvailabilityZones []*string `json:"availabilityZones,omitempty"`

	// List of Azure Managed Disks to attach to a Disk Pool.
	Disks []*Disk `json:"disks,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPoolCreateProperties.
func (d DiskPoolCreateProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalCapabilities", d.AdditionalCapabilities)
	populate(objectMap, "availabilityZones", d.AvailabilityZones)
	populate(objectMap, "disks", d.Disks)
	populate(objectMap, "subnetId", d.SubnetID)
	return json.Marshal(objectMap)
}

// DiskPoolListResult - List of Disk Pools
type DiskPoolListResult struct {
	// REQUIRED; An array of Disk pool objects.
	Value []*DiskPool `json:"value,omitempty"`

	// READ-ONLY; URI to fetch the next section of the paginated response.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPoolListResult.
func (d DiskPoolListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", d.NextLink)
	populate(objectMap, "value", d.Value)
	return json.Marshal(objectMap)
}

// DiskPoolProperties - Disk Pool response properties.
type DiskPoolProperties struct {
	// REQUIRED; Logical zone for Disk Pool resource; example: ["1"].
	AvailabilityZones []*string `json:"availabilityZones,omitempty"`

	// REQUIRED; Operational status of the Disk Pool.
	Status *OperationalStatus `json:"status,omitempty"`

	// REQUIRED; Azure Resource ID of a Subnet for the Disk Pool.
	SubnetID *string `json:"subnetId,omitempty"`

	// READ-ONLY; State of the operation on the resource.
	ProvisioningState *ProvisioningStates `json:"provisioningState,omitempty" azure:"ro"`

	// List of additional capabilities for Disk Pool.
	AdditionalCapabilities []*string `json:"additionalCapabilities,omitempty"`

	// List of Azure Managed Disks to attach to a Disk Pool.
	Disks []*Disk `json:"disks,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPoolProperties.
func (d DiskPoolProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalCapabilities", d.AdditionalCapabilities)
	populate(objectMap, "availabilityZones", d.AvailabilityZones)
	populate(objectMap, "disks", d.Disks)
	populate(objectMap, "provisioningState", d.ProvisioningState)
	populate(objectMap, "status", d.Status)
	populate(objectMap, "subnetId", d.SubnetID)
	return json.Marshal(objectMap)
}

// DiskPoolUpdate - Request payload for Update Disk Pool request.
type DiskPoolUpdate struct {
	// REQUIRED; Properties for Disk Pool update request.
	Properties *DiskPoolUpdateProperties `json:"properties,omitempty"`

	// Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedBy *string `json:"managedBy,omitempty"`

	// List of Azure resource ids that manage this resource.
	ManagedByExtended []*string `json:"managedByExtended,omitempty"`

	// Determines the SKU of the Disk Pool
	SKU *SKU `json:"sku,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPoolUpdate.
func (d DiskPoolUpdate) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "managedBy", d.ManagedBy)
	populate(objectMap, "managedByExtended", d.ManagedByExtended)
	populate(objectMap, "properties", d.Properties)
	populate(objectMap, "sku", d.SKU)
	populate(objectMap, "tags", d.Tags)
	return json.Marshal(objectMap)
}

// DiskPoolUpdateProperties - Properties for Disk Pool update request.
type DiskPoolUpdateProperties struct {
	// List of Azure Managed Disks to attach to a Disk Pool.
	Disks []*Disk `json:"disks,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPoolUpdateProperties.
func (d DiskPoolUpdateProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "disks", d.Disks)
	return json.Marshal(objectMap)
}

// DiskPoolZoneInfo - Disk Pool SKU Details
type DiskPoolZoneInfo struct {
	// READ-ONLY; List of additional capabilities for Disk Pool.
	AdditionalCapabilities []*string `json:"additionalCapabilities,omitempty" azure:"ro"`

	// READ-ONLY; Logical zone for Disk Pool resource; example: ["1"].
	AvailabilityZones []*string `json:"availabilityZones,omitempty" azure:"ro"`

	// READ-ONLY; Determines the SKU of VM deployed for Disk Pool
	SKU *SKU `json:"sku,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPoolZoneInfo.
func (d DiskPoolZoneInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalCapabilities", d.AdditionalCapabilities)
	populate(objectMap, "availabilityZones", d.AvailabilityZones)
	populate(objectMap, "sku", d.SKU)
	return json.Marshal(objectMap)
}

// DiskPoolZoneListResult - List Disk Pool skus operation response.
type DiskPoolZoneListResult struct {
	// READ-ONLY; URI to fetch the next section of the paginated response.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; The list of Disk Pool Skus.
	Value []*DiskPoolZoneInfo `json:"value,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type DiskPoolZoneListResult.
func (d DiskPoolZoneListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", d.NextLink)
	populate(objectMap, "value", d.Value)
	return json.Marshal(objectMap)
}

// DiskPoolZonesListOptions contains the optional parameters for the DiskPoolZones.List method.
type DiskPoolZonesListOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsBeginCreateOrUpdateOptions contains the optional parameters for the DiskPools.BeginCreateOrUpdate method.
type DiskPoolsBeginCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsBeginDeallocateOptions contains the optional parameters for the DiskPools.BeginDeallocate method.
type DiskPoolsBeginDeallocateOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsBeginDeleteOptions contains the optional parameters for the DiskPools.BeginDelete method.
type DiskPoolsBeginDeleteOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsBeginStartOptions contains the optional parameters for the DiskPools.BeginStart method.
type DiskPoolsBeginStartOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsBeginUpdateOptions contains the optional parameters for the DiskPools.BeginUpdate method.
type DiskPoolsBeginUpdateOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsBeginUpgradeOptions contains the optional parameters for the DiskPools.BeginUpgrade method.
type DiskPoolsBeginUpgradeOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsGetOptions contains the optional parameters for the DiskPools.Get method.
type DiskPoolsGetOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsListByResourceGroupOptions contains the optional parameters for the DiskPools.ListByResourceGroup method.
type DiskPoolsListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsListBySubscriptionOptions contains the optional parameters for the DiskPools.ListBySubscription method.
type DiskPoolsListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// DiskPoolsListOutboundNetworkDependenciesEndpointsOptions contains the optional parameters for the DiskPools.ListOutboundNetworkDependenciesEndpoints
// method.
type DiskPoolsListOutboundNetworkDependenciesEndpointsOptions struct {
	// placeholder for future optional parameters
}

// EndpointDependency - A domain name that a service is reached at, including details of the current connection status.
type EndpointDependency struct {
	// The domain name of the dependency.
	DomainName *string `json:"domainName,omitempty"`

	// The IP Addresses and Ports used when connecting to DomainName.
	EndpointDetails []*EndpointDetail `json:"endpointDetails,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type EndpointDependency.
func (e EndpointDependency) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "domainName", e.DomainName)
	populate(objectMap, "endpointDetails", e.EndpointDetails)
	return json.Marshal(objectMap)
}

// EndpointDetail - Current TCP connectivity information from the App Service Environment to a single endpoint.
type EndpointDetail struct {
	// An IP Address that Domain Name currently resolves to.
	IPAddress *string `json:"ipAddress,omitempty"`

	// Whether it is possible to create a TCP connection from the App Service Environment to this IpAddress at this Port.
	IsAccessible *bool `json:"isAccessible,omitempty"`

	// The time in milliseconds it takes for a TCP connection to be created from the App Service Environment to this IpAddress at this Port.
	Latency *float64 `json:"latency,omitempty"`

	// The port an endpoint is connected to.
	Port *int32 `json:"port,omitempty"`
}

// Error - The resource management error response.
// Implements the error and azcore.HTTPResponse interfaces.
type Error struct {
	raw string
	// RP error response.
	InnerError *ErrorResponse `json:"error,omitempty"`
}

// Error implements the error interface for type Error.
// The contents of the error text are not contractual and subject to change.
func (e Error) Error() string {
	return e.raw
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info map[string]interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ErrorResponse - The resource management error response.
type ErrorResponse struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*ErrorResponse `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ErrorResponse.
func (e ErrorResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalInfo", e.AdditionalInfo)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// IscsiLun - LUN to expose the Azure Managed Disk.
type IscsiLun struct {
	// REQUIRED; Azure Resource ID of the Managed Disk.
	ManagedDiskAzureResourceID *string `json:"managedDiskAzureResourceId,omitempty"`

	// REQUIRED; User defined name for iSCSI LUN; example: "lun0"
	Name *string `json:"name,omitempty"`

	// READ-ONLY; Specifies the Logical Unit Number of the iSCSI LUN.
	Lun *int32 `json:"lun,omitempty" azure:"ro"`
}

// IscsiTarget - Response for iSCSI Target requests.
type IscsiTarget struct {
	ProxyResource
	// REQUIRED; Properties for iSCSI Target operations.
	Properties *IscsiTargetProperties `json:"properties,omitempty"`

	// READ-ONLY; Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedBy *string `json:"managedBy,omitempty" azure:"ro"`

	// READ-ONLY; List of Azure resource ids that manage this resource.
	ManagedByExtended []*string `json:"managedByExtended,omitempty" azure:"ro"`

	// READ-ONLY; Resource metadata required by ARM RPC
	SystemData *SystemMetadata `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type IscsiTarget.
func (i IscsiTarget) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	i.ProxyResource.marshalInternal(objectMap)
	populate(objectMap, "managedBy", i.ManagedBy)
	populate(objectMap, "managedByExtended", i.ManagedByExtended)
	populate(objectMap, "properties", i.Properties)
	populate(objectMap, "systemData", i.SystemData)
	return json.Marshal(objectMap)
}

// IscsiTargetCreate - Payload for iSCSI Target create or update requests.
type IscsiTargetCreate struct {
	ProxyResource
	// REQUIRED; Properties for iSCSI Target create request.
	Properties *IscsiTargetCreateProperties `json:"properties,omitempty"`

	// Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedBy *string `json:"managedBy,omitempty"`

	// List of Azure resource ids that manage this resource.
	ManagedByExtended []*string `json:"managedByExtended,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type IscsiTargetCreate.
func (i IscsiTargetCreate) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	i.ProxyResource.marshalInternal(objectMap)
	populate(objectMap, "managedBy", i.ManagedBy)
	populate(objectMap, "managedByExtended", i.ManagedByExtended)
	populate(objectMap, "properties", i.Properties)
	return json.Marshal(objectMap)
}

// IscsiTargetCreateProperties - Properties for iSCSI Target create or update request.
type IscsiTargetCreateProperties struct {
	// REQUIRED; Mode for Target connectivity.
	ACLMode *IscsiTargetACLMode `json:"aclMode,omitempty"`

	// List of LUNs to be exposed through iSCSI Target.
	Luns []*IscsiLun `json:"luns,omitempty"`

	// Access Control List (ACL) for an iSCSI Target; defines LUN masking policy
	StaticACLs []*ACL `json:"staticAcls,omitempty"`

	// iSCSI Target IQN (iSCSI Qualified Name); example: "iqn.2005-03.org.iscsi:server".
	TargetIqn *string `json:"targetIqn,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type IscsiTargetCreateProperties.
func (i IscsiTargetCreateProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "aclMode", i.ACLMode)
	populate(objectMap, "luns", i.Luns)
	populate(objectMap, "staticAcls", i.StaticACLs)
	populate(objectMap, "targetIqn", i.TargetIqn)
	return json.Marshal(objectMap)
}

// IscsiTargetList - List of iSCSI Targets.
type IscsiTargetList struct {
	// REQUIRED; An array of iSCSI Targets in a Disk Pool.
	Value []*IscsiTarget `json:"value,omitempty"`

	// READ-ONLY; URI to fetch the next section of the paginated response.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type IscsiTargetList.
func (i IscsiTargetList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", i.NextLink)
	populate(objectMap, "value", i.Value)
	return json.Marshal(objectMap)
}

// IscsiTargetProperties - Response properties for iSCSI Target operations.
type IscsiTargetProperties struct {
	// REQUIRED; Mode for Target connectivity.
	ACLMode *IscsiTargetACLMode `json:"aclMode,omitempty"`

	// REQUIRED; Operational status of the iSCSI Target.
	Status *OperationalStatus `json:"status,omitempty"`

	// REQUIRED; iSCSI Target IQN (iSCSI Qualified Name); example: "iqn.2005-03.org.iscsi:server".
	TargetIqn *string `json:"targetIqn,omitempty"`

	// READ-ONLY; State of the operation on the resource.
	ProvisioningState *ProvisioningStates `json:"provisioningState,omitempty" azure:"ro"`

	// List of private IPv4 addresses to connect to the iSCSI Target.
	Endpoints []*string `json:"endpoints,omitempty"`

	// List of LUNs to be exposed through iSCSI Target.
	Luns []*IscsiLun `json:"luns,omitempty"`

	// The port used by iSCSI Target portal group.
	Port *int32 `json:"port,omitempty"`

	// Access Control List (ACL) for an iSCSI Target; defines LUN masking policy
	StaticACLs []*ACL `json:"staticAcls,omitempty"`

	// READ-ONLY; List of identifiers for active sessions on the iSCSI target
	Sessions []*string `json:"sessions,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type IscsiTargetProperties.
func (i IscsiTargetProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "aclMode", i.ACLMode)
	populate(objectMap, "endpoints", i.Endpoints)
	populate(objectMap, "luns", i.Luns)
	populate(objectMap, "port", i.Port)
	populate(objectMap, "provisioningState", i.ProvisioningState)
	populate(objectMap, "sessions", i.Sessions)
	populate(objectMap, "staticAcls", i.StaticACLs)
	populate(objectMap, "status", i.Status)
	populate(objectMap, "targetIqn", i.TargetIqn)
	return json.Marshal(objectMap)
}

// IscsiTargetUpdate - Payload for iSCSI Target update requests.
type IscsiTargetUpdate struct {
	ProxyResource
	// REQUIRED; Properties for iSCSI Target update request.
	Properties *IscsiTargetUpdateProperties `json:"properties,omitempty"`

	// Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedBy *string `json:"managedBy,omitempty"`

	// List of Azure resource ids that manage this resource.
	ManagedByExtended []*string `json:"managedByExtended,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type IscsiTargetUpdate.
func (i IscsiTargetUpdate) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	i.ProxyResource.marshalInternal(objectMap)
	populate(objectMap, "managedBy", i.ManagedBy)
	populate(objectMap, "managedByExtended", i.ManagedByExtended)
	populate(objectMap, "properties", i.Properties)
	return json.Marshal(objectMap)
}

// IscsiTargetUpdateProperties - Properties for iSCSI Target update request.
type IscsiTargetUpdateProperties struct {
	// List of LUNs to be exposed through iSCSI Target.
	Luns []*IscsiLun `json:"luns,omitempty"`

	// Access Control List (ACL) for an iSCSI Target; defines LUN masking policy
	StaticACLs []*ACL `json:"staticAcls,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type IscsiTargetUpdateProperties.
func (i IscsiTargetUpdateProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "luns", i.Luns)
	populate(objectMap, "staticAcls", i.StaticACLs)
	return json.Marshal(objectMap)
}

// IscsiTargetsBeginCreateOrUpdateOptions contains the optional parameters for the IscsiTargets.BeginCreateOrUpdate method.
type IscsiTargetsBeginCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// IscsiTargetsBeginDeleteOptions contains the optional parameters for the IscsiTargets.BeginDelete method.
type IscsiTargetsBeginDeleteOptions struct {
	// placeholder for future optional parameters
}

// IscsiTargetsBeginUpdateOptions contains the optional parameters for the IscsiTargets.BeginUpdate method.
type IscsiTargetsBeginUpdateOptions struct {
	// placeholder for future optional parameters
}

// IscsiTargetsGetOptions contains the optional parameters for the IscsiTargets.Get method.
type IscsiTargetsGetOptions struct {
	// placeholder for future optional parameters
}

// IscsiTargetsListByDiskPoolOptions contains the optional parameters for the IscsiTargets.ListByDiskPool method.
type IscsiTargetsListByDiskPoolOptions struct {
	// placeholder for future optional parameters
}

// OperationsListOptions contains the optional parameters for the Operations.List method.
type OperationsListOptions struct {
	// placeholder for future optional parameters
}

// OutboundEnvironmentEndpoint - Endpoints accessed for a common purpose that the App Service Environment requires outbound network access to.
type OutboundEnvironmentEndpoint struct {
	// The type of service accessed by the App Service Environment, e.g., Azure Storage, Azure SQL Database, and Azure Active Directory.
	Category *string `json:"category,omitempty"`

	// The endpoints that the App Service Environment reaches the service at.
	Endpoints []*EndpointDependency `json:"endpoints,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type OutboundEnvironmentEndpoint.
func (o OutboundEnvironmentEndpoint) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "category", o.Category)
	populate(objectMap, "endpoints", o.Endpoints)
	return json.Marshal(objectMap)
}

// OutboundEnvironmentEndpointList - Collection of Outbound Environment Endpoints
type OutboundEnvironmentEndpointList struct {
	// REQUIRED; Collection of resources.
	Value []*OutboundEnvironmentEndpoint `json:"value,omitempty"`

	// READ-ONLY; Link to next page of resources.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type OutboundEnvironmentEndpointList.
func (o OutboundEnvironmentEndpointList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", o.NextLink)
	populate(objectMap, "value", o.Value)
	return json.Marshal(objectMap)
}

// ProxyResource - The resource model definition for a ARM proxy resource. It will have everything other than required location and tags
type ProxyResource struct {
	Resource
}

func (p ProxyResource) marshalInternal(objectMap map[string]interface{}) {
	p.Resource.marshalInternal(objectMap)
}

// Resource - ARM resource model definition.
type Resource struct {
	// READ-ONLY; Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. Ex- Microsoft.Compute/virtualMachines or Microsoft.Storage/storageAccounts.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	r.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (r Resource) marshalInternal(objectMap map[string]interface{}) {
	populate(objectMap, "id", r.ID)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "type", r.Type)
}

// ResourceSKUCapability - Capability a resource SKU has.
type ResourceSKUCapability struct {
	// READ-ONLY; Capability name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Capability value
	Value *string `json:"value,omitempty" azure:"ro"`
}

// ResourceSKUInfo - Resource SKU Details
type ResourceSKUInfo struct {
	// READ-ONLY; StoragePool RP API version
	APIVersion *string `json:"apiVersion,omitempty" azure:"ro"`

	// READ-ONLY; List of additional capabilities for StoragePool resource.
	Capabilities []*ResourceSKUCapability `json:"capabilities,omitempty" azure:"ro"`

	// READ-ONLY; Zones and zone capabilities in those locations where the SKU is available.
	LocationInfo *ResourceSKULocationInfo `json:"locationInfo,omitempty" azure:"ro"`

	// READ-ONLY; Sku name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; StoragePool resource type
	ResourceType *string `json:"resourceType,omitempty" azure:"ro"`

	// READ-ONLY; The restrictions because of which SKU cannot be used. This is empty if there are no restrictions.
	Restrictions []*ResourceSKURestrictions `json:"restrictions,omitempty" azure:"ro"`

	// READ-ONLY; Sku tier
	Tier *string `json:"tier,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceSKUInfo.
func (r ResourceSKUInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "apiVersion", r.APIVersion)
	populate(objectMap, "capabilities", r.Capabilities)
	populate(objectMap, "locationInfo", r.LocationInfo)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "resourceType", r.ResourceType)
	populate(objectMap, "restrictions", r.Restrictions)
	populate(objectMap, "tier", r.Tier)
	return json.Marshal(objectMap)
}

// ResourceSKUListResult - List Disk Pool skus operation response.
type ResourceSKUListResult struct {
	// URI to fetch the next section of the paginated response.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of StoragePool resource skus.
	Value []*ResourceSKUInfo `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceSKUListResult.
func (r ResourceSKUListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", r.NextLink)
	populate(objectMap, "value", r.Value)
	return json.Marshal(objectMap)
}

// ResourceSKULocationInfo - Zone and capability info for resource sku
type ResourceSKULocationInfo struct {
	// READ-ONLY; Location of the SKU
	Location *string `json:"location,omitempty" azure:"ro"`

	// READ-ONLY; Details of capabilities available to a SKU in specific zones.
	ZoneDetails []*ResourceSKUZoneDetails `json:"zoneDetails,omitempty" azure:"ro"`

	// READ-ONLY; List of availability zones where the SKU is supported.
	Zones []*string `json:"zones,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceSKULocationInfo.
func (r ResourceSKULocationInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "location", r.Location)
	populate(objectMap, "zoneDetails", r.ZoneDetails)
	populate(objectMap, "zones", r.Zones)
	return json.Marshal(objectMap)
}

// ResourceSKURestrictionInfo - Describes an available Compute SKU Restriction Information.
type ResourceSKURestrictionInfo struct {
	// READ-ONLY; Locations where the SKU is restricted
	Locations []*string `json:"locations,omitempty" azure:"ro"`

	// READ-ONLY; List of availability zones where the SKU is restricted.
	Zones []*string `json:"zones,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceSKURestrictionInfo.
func (r ResourceSKURestrictionInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "locations", r.Locations)
	populate(objectMap, "zones", r.Zones)
	return json.Marshal(objectMap)
}

// ResourceSKURestrictions - Describes scaling information of a SKU.
type ResourceSKURestrictions struct {
	// READ-ONLY; The reason for restriction.
	ReasonCode *ResourceSKURestrictionsReasonCode `json:"reasonCode,omitempty" azure:"ro"`

	// READ-ONLY; The information about the restriction where the SKU cannot be used.
	RestrictionInfo *ResourceSKURestrictionInfo `json:"restrictionInfo,omitempty" azure:"ro"`

	// READ-ONLY; The type of restrictions.
	Type *ResourceSKURestrictionsType `json:"type,omitempty" azure:"ro"`

	// READ-ONLY; The value of restrictions. If the restriction type is set to location. This would be different locations where the SKU is restricted.
	Values []*string `json:"values,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceSKURestrictions.
func (r ResourceSKURestrictions) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "reasonCode", r.ReasonCode)
	populate(objectMap, "restrictionInfo", r.RestrictionInfo)
	populate(objectMap, "type", r.Type)
	populate(objectMap, "values", r.Values)
	return json.Marshal(objectMap)
}

// ResourceSKUZoneDetails - Describes The zonal capabilities of a SKU.
type ResourceSKUZoneDetails struct {
	// READ-ONLY; A list of capabilities that are available for the SKU in the specified list of zones.
	Capabilities []*ResourceSKUCapability `json:"capabilities,omitempty" azure:"ro"`

	// READ-ONLY; The set of zones that the SKU is available in with the specified capabilities.
	Name []*string `json:"name,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceSKUZoneDetails.
func (r ResourceSKUZoneDetails) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "capabilities", r.Capabilities)
	populate(objectMap, "name", r.Name)
	return json.Marshal(objectMap)
}

// ResourceSKUsListOptions contains the optional parameters for the ResourceSKUs.List method.
type ResourceSKUsListOptions struct {
	// placeholder for future optional parameters
}

// SKU - Sku for ARM resource
type SKU struct {
	// REQUIRED; Sku name
	Name *string `json:"name,omitempty"`

	// Sku tier
	Tier *string `json:"tier,omitempty"`
}

// StoragePoolOperationDisplay - Metadata about an operation.
type StoragePoolOperationDisplay struct {
	// REQUIRED; Localized friendly description for the operation, as it should be shown to the user.
	Description *string `json:"description,omitempty"`

	// REQUIRED; Localized friendly name for the operation, as it should be shown to the user.
	Operation *string `json:"operation,omitempty"`

	// REQUIRED; Localized friendly form of the resource provider name.
	Provider *string `json:"provider,omitempty"`

	// REQUIRED; Localized friendly form of the resource type related to this action/operation.
	Resource *string `json:"resource,omitempty"`
}

// StoragePoolOperationListResult - List of operations supported by the RP.
type StoragePoolOperationListResult struct {
	// REQUIRED; An array of operations supported by the StoragePool RP.
	Value []*StoragePoolRPOperation `json:"value,omitempty"`

	// URI to fetch the next section of the paginated response.
	NextLink *string `json:"nextLink,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type StoragePoolOperationListResult.
func (s StoragePoolOperationListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", s.NextLink)
	populate(objectMap, "value", s.Value)
	return json.Marshal(objectMap)
}

// StoragePoolRPOperation - Description of a StoragePool RP Operation
type StoragePoolRPOperation struct {
	// REQUIRED; Additional metadata about RP operation.
	Display *StoragePoolOperationDisplay `json:"display,omitempty"`

	// REQUIRED; Indicates whether the operation applies to data-plane.
	IsDataAction *bool `json:"isDataAction,omitempty"`

	// REQUIRED; The name of the operation being performed on this particular object
	Name *string `json:"name,omitempty"`

	// Indicates the action type.
	ActionType *string `json:"actionType,omitempty"`

	// The intended executor of the operation; governs the display of the operation in the RBAC UX and the audit logs UX.
	Origin *string `json:"origin,omitempty"`
}

// SystemMetadata - Metadata pertaining to creation and last modification of the resource.
type SystemMetadata struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type SystemMetadata.
func (s SystemMetadata) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populateTimeRFC3339(objectMap, "createdAt", s.CreatedAt)
	populate(objectMap, "createdBy", s.CreatedBy)
	populate(objectMap, "createdByType", s.CreatedByType)
	populateTimeRFC3339(objectMap, "lastModifiedAt", s.LastModifiedAt)
	populate(objectMap, "lastModifiedBy", s.LastModifiedBy)
	populate(objectMap, "lastModifiedByType", s.LastModifiedByType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type SystemMetadata.
func (s *SystemMetadata) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "createdAt":
			err = unpopulateTimeRFC3339(val, &s.CreatedAt)
			delete(rawMsg, key)
		case "createdBy":
			err = unpopulate(val, &s.CreatedBy)
			delete(rawMsg, key)
		case "createdByType":
			err = unpopulate(val, &s.CreatedByType)
			delete(rawMsg, key)
		case "lastModifiedAt":
			err = unpopulateTimeRFC3339(val, &s.LastModifiedAt)
			delete(rawMsg, key)
		case "lastModifiedBy":
			err = unpopulate(val, &s.LastModifiedBy)
			delete(rawMsg, key)
		case "lastModifiedByType":
			err = unpopulate(val, &s.LastModifiedByType)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TrackedResource - The resource model definition for a ARM tracked top level resource.
type TrackedResource struct {
	Resource
	// REQUIRED; The geo-location where the resource lives.
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type TrackedResource.
func (t TrackedResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	t.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (t TrackedResource) marshalInternal(objectMap map[string]interface{}) {
	t.Resource.marshalInternal(objectMap)
	populate(objectMap, "location", t.Location)
	populate(objectMap, "tags", t.Tags)
}

func populate(m map[string]interface{}, k string, v interface{}) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}

func unpopulate(data json.RawMessage, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}

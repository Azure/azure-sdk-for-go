//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armtimeseriesinsights

const (
	moduleName    = "armtimeseriesinsights"
	moduleVersion = "v0.1.0"
)

// AccessPolicyRole - A role defining the data plane operations that a principal can perform on a Time Series Insights client.
type AccessPolicyRole string

const (
	AccessPolicyRoleContributor AccessPolicyRole = "Contributor"
	AccessPolicyRoleReader      AccessPolicyRole = "Reader"
)

// PossibleAccessPolicyRoleValues returns the possible values for the AccessPolicyRole const type.
func PossibleAccessPolicyRoleValues() []AccessPolicyRole {
	return []AccessPolicyRole{
		AccessPolicyRoleContributor,
		AccessPolicyRoleReader,
	}
}

// DataStringComparisonBehavior - The reference data set key comparison behavior can be set using this property. By default,
// the value is 'Ordinal' - which means case sensitive key comparison will be performed while joining reference
// data with events or while adding new reference data. When 'OrdinalIgnoreCase' is set, case insensitive comparison will
// be used.
type DataStringComparisonBehavior string

const (
	DataStringComparisonBehaviorOrdinal           DataStringComparisonBehavior = "Ordinal"
	DataStringComparisonBehaviorOrdinalIgnoreCase DataStringComparisonBehavior = "OrdinalIgnoreCase"
)

// PossibleDataStringComparisonBehaviorValues returns the possible values for the DataStringComparisonBehavior const type.
func PossibleDataStringComparisonBehaviorValues() []DataStringComparisonBehavior {
	return []DataStringComparisonBehavior{
		DataStringComparisonBehaviorOrdinal,
		DataStringComparisonBehaviorOrdinalIgnoreCase,
	}
}

// EnvironmentKind - The kind of the environment.
type EnvironmentKind string

const (
	EnvironmentKindGen1 EnvironmentKind = "Gen1"
	EnvironmentKindGen2 EnvironmentKind = "Gen2"
)

// PossibleEnvironmentKindValues returns the possible values for the EnvironmentKind const type.
func PossibleEnvironmentKindValues() []EnvironmentKind {
	return []EnvironmentKind{
		EnvironmentKindGen1,
		EnvironmentKindGen2,
	}
}

// EnvironmentResourceKind - The kind of the environment.
type EnvironmentResourceKind string

const (
	EnvironmentResourceKindGen1 EnvironmentResourceKind = "Gen1"
	EnvironmentResourceKindGen2 EnvironmentResourceKind = "Gen2"
)

// PossibleEnvironmentResourceKindValues returns the possible values for the EnvironmentResourceKind const type.
func PossibleEnvironmentResourceKindValues() []EnvironmentResourceKind {
	return []EnvironmentResourceKind{
		EnvironmentResourceKindGen1,
		EnvironmentResourceKindGen2,
	}
}

// EventSourceKind - The kind of the event source.
type EventSourceKind string

const (
	EventSourceKindMicrosoftEventHub EventSourceKind = "Microsoft.EventHub"
	EventSourceKindMicrosoftIoTHub   EventSourceKind = "Microsoft.IoTHub"
)

// PossibleEventSourceKindValues returns the possible values for the EventSourceKind const type.
func PossibleEventSourceKindValues() []EventSourceKind {
	return []EventSourceKind{
		EventSourceKindMicrosoftEventHub,
		EventSourceKindMicrosoftIoTHub,
	}
}

// EventSourceResourceKind - The kind of the event source.
type EventSourceResourceKind string

const (
	EventSourceResourceKindMicrosoftEventHub EventSourceResourceKind = "Microsoft.EventHub"
	EventSourceResourceKindMicrosoftIoTHub   EventSourceResourceKind = "Microsoft.IoTHub"
)

// PossibleEventSourceResourceKindValues returns the possible values for the EventSourceResourceKind const type.
func PossibleEventSourceResourceKindValues() []EventSourceResourceKind {
	return []EventSourceResourceKind{
		EventSourceResourceKindMicrosoftEventHub,
		EventSourceResourceKindMicrosoftIoTHub,
	}
}

// IngressStartAtType - The type of the ingressStartAt, It can be "EarliestAvailable", "EventSourceCreationTime", "CustomEnqueuedTime".
type IngressStartAtType string

const (
	IngressStartAtTypeCustomEnqueuedTime      IngressStartAtType = "CustomEnqueuedTime"
	IngressStartAtTypeEarliestAvailable       IngressStartAtType = "EarliestAvailable"
	IngressStartAtTypeEventSourceCreationTime IngressStartAtType = "EventSourceCreationTime"
)

// PossibleIngressStartAtTypeValues returns the possible values for the IngressStartAtType const type.
func PossibleIngressStartAtTypeValues() []IngressStartAtType {
	return []IngressStartAtType{
		IngressStartAtTypeCustomEnqueuedTime,
		IngressStartAtTypeEarliestAvailable,
		IngressStartAtTypeEventSourceCreationTime,
	}
}

// IngressState - This string represents the state of ingress operations on an environment. It can be "Disabled", "Ready",
// "Running", "Paused" or "Unknown"
type IngressState string

const (
	IngressStateDisabled IngressState = "Disabled"
	IngressStatePaused   IngressState = "Paused"
	IngressStateReady    IngressState = "Ready"
	IngressStateRunning  IngressState = "Running"
	IngressStateUnknown  IngressState = "Unknown"
)

// PossibleIngressStateValues returns the possible values for the IngressState const type.
func PossibleIngressStateValues() []IngressState {
	return []IngressState{
		IngressStateDisabled,
		IngressStatePaused,
		IngressStateReady,
		IngressStateRunning,
		IngressStateUnknown,
	}
}

// LocalTimestampFormat - An enum that represents the format of the local timestamp property that needs to be set.
type LocalTimestampFormat string

const (
	LocalTimestampFormatEmbedded LocalTimestampFormat = "Embedded"
)

// PossibleLocalTimestampFormatValues returns the possible values for the LocalTimestampFormat const type.
func PossibleLocalTimestampFormatValues() []LocalTimestampFormat {
	return []LocalTimestampFormat{
		LocalTimestampFormatEmbedded,
	}
}

// PrivateEndpointConnectionProvisioningState - The current provisioning state.
type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

// PossiblePrivateEndpointConnectionProvisioningStateValues returns the possible values for the PrivateEndpointConnectionProvisioningState const type.
func PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState {
	return []PrivateEndpointConnectionProvisioningState{
		PrivateEndpointConnectionProvisioningStateCreating,
		PrivateEndpointConnectionProvisioningStateDeleting,
		PrivateEndpointConnectionProvisioningStateFailed,
		PrivateEndpointConnectionProvisioningStateSucceeded,
	}
}

// PrivateEndpointServiceConnectionStatus - The private endpoint connection status.
type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

// PossiblePrivateEndpointServiceConnectionStatusValues returns the possible values for the PrivateEndpointServiceConnectionStatus const type.
func PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus {
	return []PrivateEndpointServiceConnectionStatus{
		PrivateEndpointServiceConnectionStatusApproved,
		PrivateEndpointServiceConnectionStatusPending,
		PrivateEndpointServiceConnectionStatusRejected,
	}
}

// PropertyType - The type of the property.
type PropertyType string

const (
	PropertyTypeString PropertyType = "String"
)

// PossiblePropertyTypeValues returns the possible values for the PropertyType const type.
func PossiblePropertyTypeValues() []PropertyType {
	return []PropertyType{
		PropertyTypeString,
	}
}

// ProvisioningState - Provisioning state of the resource.
type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{
		ProvisioningStateAccepted,
		ProvisioningStateCreating,
		ProvisioningStateDeleting,
		ProvisioningStateFailed,
		ProvisioningStateSucceeded,
		ProvisioningStateUpdating,
	}
}

// PublicNetworkAccess - This value can be set to 'enabled' to avoid breaking changes on existing customer resources and templates.
// If set to 'disabled', traffic over public interface is not allowed, and private endpoint
// connections would be the exclusive access method.
type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "enabled"
)

// PossiblePublicNetworkAccessValues returns the possible values for the PublicNetworkAccess const type.
func PossiblePublicNetworkAccessValues() []PublicNetworkAccess {
	return []PublicNetworkAccess{
		PublicNetworkAccessDisabled,
		PublicNetworkAccessEnabled,
	}
}

// ReferenceDataKeyPropertyType - The type of the key property.
type ReferenceDataKeyPropertyType string

const (
	ReferenceDataKeyPropertyTypeBool     ReferenceDataKeyPropertyType = "Bool"
	ReferenceDataKeyPropertyTypeDateTime ReferenceDataKeyPropertyType = "DateTime"
	ReferenceDataKeyPropertyTypeDouble   ReferenceDataKeyPropertyType = "Double"
	ReferenceDataKeyPropertyTypeString   ReferenceDataKeyPropertyType = "String"
)

// PossibleReferenceDataKeyPropertyTypeValues returns the possible values for the ReferenceDataKeyPropertyType const type.
func PossibleReferenceDataKeyPropertyTypeValues() []ReferenceDataKeyPropertyType {
	return []ReferenceDataKeyPropertyType{
		ReferenceDataKeyPropertyTypeBool,
		ReferenceDataKeyPropertyTypeDateTime,
		ReferenceDataKeyPropertyTypeDouble,
		ReferenceDataKeyPropertyTypeString,
	}
}

// SKUName - The name of this SKU.
type SKUName string

const (
	SKUNameL1 SKUName = "L1"
	SKUNameP1 SKUName = "P1"
	SKUNameS1 SKUName = "S1"
	SKUNameS2 SKUName = "S2"
)

// PossibleSKUNameValues returns the possible values for the SKUName const type.
func PossibleSKUNameValues() []SKUName {
	return []SKUName{
		SKUNameL1,
		SKUNameP1,
		SKUNameS1,
		SKUNameS2,
	}
}

// StorageLimitExceededBehavior - The behavior the Time Series Insights service should take when the environment's capacity
// has been exceeded. If "PauseIngress" is specified, new events will not be read from the event source. If
// "PurgeOldData" is specified, new events will continue to be read and old events will be deleted from the environment. The
// default behavior is PurgeOldData.
type StorageLimitExceededBehavior string

const (
	StorageLimitExceededBehaviorPauseIngress StorageLimitExceededBehavior = "PauseIngress"
	StorageLimitExceededBehaviorPurgeOldData StorageLimitExceededBehavior = "PurgeOldData"
)

// PossibleStorageLimitExceededBehaviorValues returns the possible values for the StorageLimitExceededBehavior const type.
func PossibleStorageLimitExceededBehaviorValues() []StorageLimitExceededBehavior {
	return []StorageLimitExceededBehavior{
		StorageLimitExceededBehaviorPauseIngress,
		StorageLimitExceededBehaviorPurgeOldData,
	}
}

// WarmStoragePropertiesState - This string represents the state of warm storage properties usage. It can be "Ok", "Error",
// "Unknown".
type WarmStoragePropertiesState string

const (
	WarmStoragePropertiesStateError   WarmStoragePropertiesState = "Error"
	WarmStoragePropertiesStateOk      WarmStoragePropertiesState = "Ok"
	WarmStoragePropertiesStateUnknown WarmStoragePropertiesState = "Unknown"
)

// PossibleWarmStoragePropertiesStateValues returns the possible values for the WarmStoragePropertiesState const type.
func PossibleWarmStoragePropertiesStateValues() []WarmStoragePropertiesState {
	return []WarmStoragePropertiesState{
		WarmStoragePropertiesStateError,
		WarmStoragePropertiesStateOk,
		WarmStoragePropertiesStateUnknown,
	}
}

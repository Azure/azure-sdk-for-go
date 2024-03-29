//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnotificationhubs

const (
	moduleName    = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/notificationhubs/armnotificationhubs"
	moduleVersion = "v2.0.0-beta.1"
)

// AccessRights - Defines values for AccessRights.
type AccessRights string

const (
	AccessRightsListen AccessRights = "Listen"
	AccessRightsManage AccessRights = "Manage"
	AccessRightsSend   AccessRights = "Send"
)

// PossibleAccessRightsValues returns the possible values for the AccessRights const type.
func PossibleAccessRightsValues() []AccessRights {
	return []AccessRights{
		AccessRightsListen,
		AccessRightsManage,
		AccessRightsSend,
	}
}

// CreatedByType - The type of identity that created the resource.
type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

// PossibleCreatedByTypeValues returns the possible values for the CreatedByType const type.
func PossibleCreatedByTypeValues() []CreatedByType {
	return []CreatedByType{
		CreatedByTypeApplication,
		CreatedByTypeKey,
		CreatedByTypeManagedIdentity,
		CreatedByTypeUser,
	}
}

// NamespaceStatus - Namespace status.
type NamespaceStatus string

const (
	NamespaceStatusCreated   NamespaceStatus = "Created"
	NamespaceStatusCreating  NamespaceStatus = "Creating"
	NamespaceStatusDeleting  NamespaceStatus = "Deleting"
	NamespaceStatusSuspended NamespaceStatus = "Suspended"
)

// PossibleNamespaceStatusValues returns the possible values for the NamespaceStatus const type.
func PossibleNamespaceStatusValues() []NamespaceStatus {
	return []NamespaceStatus{
		NamespaceStatusCreated,
		NamespaceStatusCreating,
		NamespaceStatusDeleting,
		NamespaceStatusSuspended,
	}
}

// NamespaceType - Defines values for NamespaceType.
type NamespaceType string

const (
	NamespaceTypeMessaging       NamespaceType = "Messaging"
	NamespaceTypeNotificationHub NamespaceType = "NotificationHub"
)

// PossibleNamespaceTypeValues returns the possible values for the NamespaceType const type.
func PossibleNamespaceTypeValues() []NamespaceType {
	return []NamespaceType{
		NamespaceTypeMessaging,
		NamespaceTypeNotificationHub,
	}
}

// OperationProvisioningState - Defines values for OperationProvisioningState.
type OperationProvisioningState string

const (
	OperationProvisioningStateCanceled   OperationProvisioningState = "Canceled"
	OperationProvisioningStateDisabled   OperationProvisioningState = "Disabled"
	OperationProvisioningStateFailed     OperationProvisioningState = "Failed"
	OperationProvisioningStateInProgress OperationProvisioningState = "InProgress"
	OperationProvisioningStatePending    OperationProvisioningState = "Pending"
	OperationProvisioningStateSucceeded  OperationProvisioningState = "Succeeded"
	OperationProvisioningStateUnknown    OperationProvisioningState = "Unknown"
)

// PossibleOperationProvisioningStateValues returns the possible values for the OperationProvisioningState const type.
func PossibleOperationProvisioningStateValues() []OperationProvisioningState {
	return []OperationProvisioningState{
		OperationProvisioningStateCanceled,
		OperationProvisioningStateDisabled,
		OperationProvisioningStateFailed,
		OperationProvisioningStateInProgress,
		OperationProvisioningStatePending,
		OperationProvisioningStateSucceeded,
		OperationProvisioningStateUnknown,
	}
}

// PolicyKeyType - Type of Shared Access Policy Key (primary or secondary).
type PolicyKeyType string

const (
	PolicyKeyTypePrimaryKey   PolicyKeyType = "PrimaryKey"
	PolicyKeyTypeSecondaryKey PolicyKeyType = "SecondaryKey"
)

// PossiblePolicyKeyTypeValues returns the possible values for the PolicyKeyType const type.
func PossiblePolicyKeyTypeValues() []PolicyKeyType {
	return []PolicyKeyType{
		PolicyKeyTypePrimaryKey,
		PolicyKeyTypeSecondaryKey,
	}
}

// PrivateEndpointConnectionProvisioningState - State of Private Endpoint Connection.
type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating        PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleted         PrivateEndpointConnectionProvisioningState = "Deleted"
	PrivateEndpointConnectionProvisioningStateDeleting        PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateDeletingByProxy PrivateEndpointConnectionProvisioningState = "DeletingByProxy"
	PrivateEndpointConnectionProvisioningStateSucceeded       PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUnknown         PrivateEndpointConnectionProvisioningState = "Unknown"
	PrivateEndpointConnectionProvisioningStateUpdating        PrivateEndpointConnectionProvisioningState = "Updating"
	PrivateEndpointConnectionProvisioningStateUpdatingByProxy PrivateEndpointConnectionProvisioningState = "UpdatingByProxy"
)

// PossiblePrivateEndpointConnectionProvisioningStateValues returns the possible values for the PrivateEndpointConnectionProvisioningState const type.
func PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState {
	return []PrivateEndpointConnectionProvisioningState{
		PrivateEndpointConnectionProvisioningStateCreating,
		PrivateEndpointConnectionProvisioningStateDeleted,
		PrivateEndpointConnectionProvisioningStateDeleting,
		PrivateEndpointConnectionProvisioningStateDeletingByProxy,
		PrivateEndpointConnectionProvisioningStateSucceeded,
		PrivateEndpointConnectionProvisioningStateUnknown,
		PrivateEndpointConnectionProvisioningStateUpdating,
		PrivateEndpointConnectionProvisioningStateUpdatingByProxy,
	}
}

// PrivateLinkConnectionStatus - State of Private Link Connection.
type PrivateLinkConnectionStatus string

const (
	PrivateLinkConnectionStatusApproved     PrivateLinkConnectionStatus = "Approved"
	PrivateLinkConnectionStatusDisconnected PrivateLinkConnectionStatus = "Disconnected"
	PrivateLinkConnectionStatusPending      PrivateLinkConnectionStatus = "Pending"
	PrivateLinkConnectionStatusRejected     PrivateLinkConnectionStatus = "Rejected"
)

// PossiblePrivateLinkConnectionStatusValues returns the possible values for the PrivateLinkConnectionStatus const type.
func PossiblePrivateLinkConnectionStatusValues() []PrivateLinkConnectionStatus {
	return []PrivateLinkConnectionStatus{
		PrivateLinkConnectionStatusApproved,
		PrivateLinkConnectionStatusDisconnected,
		PrivateLinkConnectionStatusPending,
		PrivateLinkConnectionStatusRejected,
	}
}

// PublicNetworkAccess - Type of public network access.
type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

// PossiblePublicNetworkAccessValues returns the possible values for the PublicNetworkAccess const type.
func PossiblePublicNetworkAccessValues() []PublicNetworkAccess {
	return []PublicNetworkAccess{
		PublicNetworkAccessDisabled,
		PublicNetworkAccessEnabled,
	}
}

// ReplicationRegion - Allowed replication region
type ReplicationRegion string

const (
	ReplicationRegionAustraliaEast    ReplicationRegion = "AustraliaEast"
	ReplicationRegionBrazilSouth      ReplicationRegion = "BrazilSouth"
	ReplicationRegionDefault          ReplicationRegion = "Default"
	ReplicationRegionNone             ReplicationRegion = "None"
	ReplicationRegionNorthEurope      ReplicationRegion = "NorthEurope"
	ReplicationRegionSouthAfricaNorth ReplicationRegion = "SouthAfricaNorth"
	ReplicationRegionSouthEastAsia    ReplicationRegion = "SouthEastAsia"
	ReplicationRegionWestUs2          ReplicationRegion = "WestUs2"
)

// PossibleReplicationRegionValues returns the possible values for the ReplicationRegion const type.
func PossibleReplicationRegionValues() []ReplicationRegion {
	return []ReplicationRegion{
		ReplicationRegionAustraliaEast,
		ReplicationRegionBrazilSouth,
		ReplicationRegionDefault,
		ReplicationRegionNone,
		ReplicationRegionNorthEurope,
		ReplicationRegionSouthAfricaNorth,
		ReplicationRegionSouthEastAsia,
		ReplicationRegionWestUs2,
	}
}

// SKUName - Namespace SKU name.
type SKUName string

const (
	SKUNameBasic    SKUName = "Basic"
	SKUNameFree     SKUName = "Free"
	SKUNameStandard SKUName = "Standard"
)

// PossibleSKUNameValues returns the possible values for the SKUName const type.
func PossibleSKUNameValues() []SKUName {
	return []SKUName{
		SKUNameBasic,
		SKUNameFree,
		SKUNameStandard,
	}
}

// ZoneRedundancyPreference - Namespace SKU name.
type ZoneRedundancyPreference string

const (
	ZoneRedundancyPreferenceDisabled ZoneRedundancyPreference = "Disabled"
	ZoneRedundancyPreferenceEnabled  ZoneRedundancyPreference = "Enabled"
)

// PossibleZoneRedundancyPreferenceValues returns the possible values for the ZoneRedundancyPreference const type.
func PossibleZoneRedundancyPreferenceValues() []ZoneRedundancyPreference {
	return []ZoneRedundancyPreference{
		ZoneRedundancyPreferenceDisabled,
		ZoneRedundancyPreferenceEnabled,
	}
}

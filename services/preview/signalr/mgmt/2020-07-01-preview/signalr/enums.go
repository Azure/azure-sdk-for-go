package signalr

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// ACLAction enumerates the values for acl action.
type ACLAction string

const (
	// Allow ...
	Allow ACLAction = "Allow"
	// Deny ...
	Deny ACLAction = "Deny"
)

// PossibleACLActionValues returns an array of possible values for the ACLAction const type.
func PossibleACLActionValues() []ACLAction {
	return []ACLAction{Allow, Deny}
}

// FeatureFlags enumerates the values for feature flags.
type FeatureFlags string

const (
	// EnableConnectivityLogs ...
	EnableConnectivityLogs FeatureFlags = "EnableConnectivityLogs"
	// EnableMessagingLogs ...
	EnableMessagingLogs FeatureFlags = "EnableMessagingLogs"
	// ServiceMode ...
	ServiceMode FeatureFlags = "ServiceMode"
)

// PossibleFeatureFlagsValues returns an array of possible values for the FeatureFlags const type.
func PossibleFeatureFlagsValues() []FeatureFlags {
	return []FeatureFlags{EnableConnectivityLogs, EnableMessagingLogs, ServiceMode}
}

// KeyType enumerates the values for key type.
type KeyType string

const (
	// Primary ...
	Primary KeyType = "Primary"
	// Secondary ...
	Secondary KeyType = "Secondary"
)

// PossibleKeyTypeValues returns an array of possible values for the KeyType const type.
func PossibleKeyTypeValues() []KeyType {
	return []KeyType{Primary, Secondary}
}

// ManagedIdentityType enumerates the values for managed identity type.
type ManagedIdentityType string

const (
	// None ...
	None ManagedIdentityType = "None"
	// SystemAssigned ...
	SystemAssigned ManagedIdentityType = "SystemAssigned"
	// UserAssigned ...
	UserAssigned ManagedIdentityType = "UserAssigned"
)

// PossibleManagedIdentityTypeValues returns an array of possible values for the ManagedIdentityType const type.
func PossibleManagedIdentityTypeValues() []ManagedIdentityType {
	return []ManagedIdentityType{None, SystemAssigned, UserAssigned}
}

// PrivateLinkServiceConnectionStatus enumerates the values for private link service connection status.
type PrivateLinkServiceConnectionStatus string

const (
	// Approved ...
	Approved PrivateLinkServiceConnectionStatus = "Approved"
	// Disconnected ...
	Disconnected PrivateLinkServiceConnectionStatus = "Disconnected"
	// Pending ...
	Pending PrivateLinkServiceConnectionStatus = "Pending"
	// Rejected ...
	Rejected PrivateLinkServiceConnectionStatus = "Rejected"
)

// PossiblePrivateLinkServiceConnectionStatusValues returns an array of possible values for the PrivateLinkServiceConnectionStatus const type.
func PossiblePrivateLinkServiceConnectionStatusValues() []PrivateLinkServiceConnectionStatus {
	return []PrivateLinkServiceConnectionStatus{Approved, Disconnected, Pending, Rejected}
}

// ProvisioningState enumerates the values for provisioning state.
type ProvisioningState string

const (
	// Canceled ...
	Canceled ProvisioningState = "Canceled"
	// Creating ...
	Creating ProvisioningState = "Creating"
	// Deleting ...
	Deleting ProvisioningState = "Deleting"
	// Failed ...
	Failed ProvisioningState = "Failed"
	// Moving ...
	Moving ProvisioningState = "Moving"
	// Running ...
	Running ProvisioningState = "Running"
	// Succeeded ...
	Succeeded ProvisioningState = "Succeeded"
	// Unknown ...
	Unknown ProvisioningState = "Unknown"
	// Updating ...
	Updating ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns an array of possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{Canceled, Creating, Deleting, Failed, Moving, Running, Succeeded, Unknown, Updating}
}

// RequestType enumerates the values for request type.
type RequestType string

const (
	// ClientConnection ...
	ClientConnection RequestType = "ClientConnection"
	// RESTAPI ...
	RESTAPI RequestType = "RESTAPI"
	// ServerConnection ...
	ServerConnection RequestType = "ServerConnection"
)

// PossibleRequestTypeValues returns an array of possible values for the RequestType const type.
func PossibleRequestTypeValues() []RequestType {
	return []RequestType{ClientConnection, RESTAPI, ServerConnection}
}

// ServiceKind enumerates the values for service kind.
type ServiceKind string

const (
	// RawWebSockets ...
	RawWebSockets ServiceKind = "RawWebSockets"
	// SignalR ...
	SignalR ServiceKind = "SignalR"
)

// PossibleServiceKindValues returns an array of possible values for the ServiceKind const type.
func PossibleServiceKindValues() []ServiceKind {
	return []ServiceKind{RawWebSockets, SignalR}
}

// SkuTier enumerates the values for sku tier.
type SkuTier string

const (
	// Basic ...
	Basic SkuTier = "Basic"
	// Free ...
	Free SkuTier = "Free"
	// Premium ...
	Premium SkuTier = "Premium"
	// Standard ...
	Standard SkuTier = "Standard"
)

// PossibleSkuTierValues returns an array of possible values for the SkuTier const type.
func PossibleSkuTierValues() []SkuTier {
	return []SkuTier{Basic, Free, Premium, Standard}
}

// UpstreamAuthType enumerates the values for upstream auth type.
type UpstreamAuthType string

const (
	// UpstreamAuthTypeManagedIdentity ...
	UpstreamAuthTypeManagedIdentity UpstreamAuthType = "ManagedIdentity"
	// UpstreamAuthTypeNone ...
	UpstreamAuthTypeNone UpstreamAuthType = "None"
)

// PossibleUpstreamAuthTypeValues returns an array of possible values for the UpstreamAuthType const type.
func PossibleUpstreamAuthTypeValues() []UpstreamAuthType {
	return []UpstreamAuthType{UpstreamAuthTypeManagedIdentity, UpstreamAuthTypeNone}
}

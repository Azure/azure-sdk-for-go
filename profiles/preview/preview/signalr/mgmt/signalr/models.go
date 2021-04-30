// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package signalr

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/preview/signalr/mgmt/2021-04-01-preview/signalr"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type ACLAction = original.ACLAction

const (
	ACLActionAllow ACLAction = original.ACLActionAllow
	ACLActionDeny  ACLAction = original.ACLActionDeny
)

type CreatedByType = original.CreatedByType

const (
	CreatedByTypeApplication     CreatedByType = original.CreatedByTypeApplication
	CreatedByTypeKey             CreatedByType = original.CreatedByTypeKey
	CreatedByTypeManagedIdentity CreatedByType = original.CreatedByTypeManagedIdentity
	CreatedByTypeUser            CreatedByType = original.CreatedByTypeUser
)

type FeatureFlags = original.FeatureFlags

const (
	FeatureFlagsEnableConnectivityLogs FeatureFlags = original.FeatureFlagsEnableConnectivityLogs
	FeatureFlagsEnableLiveTrace        FeatureFlags = original.FeatureFlagsEnableLiveTrace
	FeatureFlagsEnableMessagingLogs    FeatureFlags = original.FeatureFlagsEnableMessagingLogs
	FeatureFlagsServiceMode            FeatureFlags = original.FeatureFlagsServiceMode
)

type KeyType = original.KeyType

const (
	KeyTypePrimary   KeyType = original.KeyTypePrimary
	KeyTypeSecondary KeyType = original.KeyTypeSecondary
)

type ManagedIdentityType = original.ManagedIdentityType

const (
	ManagedIdentityTypeNone           ManagedIdentityType = original.ManagedIdentityTypeNone
	ManagedIdentityTypeSystemAssigned ManagedIdentityType = original.ManagedIdentityTypeSystemAssigned
	ManagedIdentityTypeUserAssigned   ManagedIdentityType = original.ManagedIdentityTypeUserAssigned
)

type PrivateLinkServiceConnectionStatus = original.PrivateLinkServiceConnectionStatus

const (
	PrivateLinkServiceConnectionStatusApproved     PrivateLinkServiceConnectionStatus = original.PrivateLinkServiceConnectionStatusApproved
	PrivateLinkServiceConnectionStatusDisconnected PrivateLinkServiceConnectionStatus = original.PrivateLinkServiceConnectionStatusDisconnected
	PrivateLinkServiceConnectionStatusPending      PrivateLinkServiceConnectionStatus = original.PrivateLinkServiceConnectionStatusPending
	PrivateLinkServiceConnectionStatusRejected     PrivateLinkServiceConnectionStatus = original.PrivateLinkServiceConnectionStatusRejected
)

type ProvisioningState = original.ProvisioningState

const (
	ProvisioningStateCanceled  ProvisioningState = original.ProvisioningStateCanceled
	ProvisioningStateCreating  ProvisioningState = original.ProvisioningStateCreating
	ProvisioningStateDeleting  ProvisioningState = original.ProvisioningStateDeleting
	ProvisioningStateFailed    ProvisioningState = original.ProvisioningStateFailed
	ProvisioningStateMoving    ProvisioningState = original.ProvisioningStateMoving
	ProvisioningStateRunning   ProvisioningState = original.ProvisioningStateRunning
	ProvisioningStateSucceeded ProvisioningState = original.ProvisioningStateSucceeded
	ProvisioningStateUnknown   ProvisioningState = original.ProvisioningStateUnknown
	ProvisioningStateUpdating  ProvisioningState = original.ProvisioningStateUpdating
)

type RequestType = original.RequestType

const (
	RequestTypeClientConnection RequestType = original.RequestTypeClientConnection
	RequestTypeRESTAPI          RequestType = original.RequestTypeRESTAPI
	RequestTypeServerConnection RequestType = original.RequestTypeServerConnection
	RequestTypeTrace            RequestType = original.RequestTypeTrace
)

type ServiceKind = original.ServiceKind

const (
	ServiceKindRawWebSockets ServiceKind = original.ServiceKindRawWebSockets
	ServiceKindSignalR       ServiceKind = original.ServiceKindSignalR
)

type SharedPrivateLinkResourceStatus = original.SharedPrivateLinkResourceStatus

const (
	SharedPrivateLinkResourceStatusApproved     SharedPrivateLinkResourceStatus = original.SharedPrivateLinkResourceStatusApproved
	SharedPrivateLinkResourceStatusDisconnected SharedPrivateLinkResourceStatus = original.SharedPrivateLinkResourceStatusDisconnected
	SharedPrivateLinkResourceStatusPending      SharedPrivateLinkResourceStatus = original.SharedPrivateLinkResourceStatusPending
	SharedPrivateLinkResourceStatusRejected     SharedPrivateLinkResourceStatus = original.SharedPrivateLinkResourceStatusRejected
	SharedPrivateLinkResourceStatusTimeout      SharedPrivateLinkResourceStatus = original.SharedPrivateLinkResourceStatusTimeout
)

type SkuTier = original.SkuTier

const (
	SkuTierBasic    SkuTier = original.SkuTierBasic
	SkuTierFree     SkuTier = original.SkuTierFree
	SkuTierPremium  SkuTier = original.SkuTierPremium
	SkuTierStandard SkuTier = original.SkuTierStandard
)

type UpstreamAuthType = original.UpstreamAuthType

const (
	UpstreamAuthTypeManagedIdentity UpstreamAuthType = original.UpstreamAuthTypeManagedIdentity
	UpstreamAuthTypeNone            UpstreamAuthType = original.UpstreamAuthTypeNone
)

type BaseClient = original.BaseClient
type Client = original.Client
type CorsSettings = original.CorsSettings
type CreateOrUpdateFuture = original.CreateOrUpdateFuture
type DeleteFuture = original.DeleteFuture
type Dimension = original.Dimension
type ErrorAdditionalInfo = original.ErrorAdditionalInfo
type ErrorDetail = original.ErrorDetail
type ErrorResponse = original.ErrorResponse
type Feature = original.Feature
type Keys = original.Keys
type LogSpecification = original.LogSpecification
type ManagedIdentity = original.ManagedIdentity
type ManagedIdentitySettings = original.ManagedIdentitySettings
type MetricSpecification = original.MetricSpecification
type NameAvailability = original.NameAvailability
type NameAvailabilityParameters = original.NameAvailabilityParameters
type NetworkACL = original.NetworkACL
type NetworkACLs = original.NetworkACLs
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationList = original.OperationList
type OperationListIterator = original.OperationListIterator
type OperationListPage = original.OperationListPage
type OperationProperties = original.OperationProperties
type OperationsClient = original.OperationsClient
type PrivateEndpoint = original.PrivateEndpoint
type PrivateEndpointACL = original.PrivateEndpointACL
type PrivateEndpointConnection = original.PrivateEndpointConnection
type PrivateEndpointConnectionList = original.PrivateEndpointConnectionList
type PrivateEndpointConnectionListIterator = original.PrivateEndpointConnectionListIterator
type PrivateEndpointConnectionListPage = original.PrivateEndpointConnectionListPage
type PrivateEndpointConnectionProperties = original.PrivateEndpointConnectionProperties
type PrivateEndpointConnectionsClient = original.PrivateEndpointConnectionsClient
type PrivateEndpointConnectionsDeleteFuture = original.PrivateEndpointConnectionsDeleteFuture
type PrivateLinkResource = original.PrivateLinkResource
type PrivateLinkResourceList = original.PrivateLinkResourceList
type PrivateLinkResourceListIterator = original.PrivateLinkResourceListIterator
type PrivateLinkResourceListPage = original.PrivateLinkResourceListPage
type PrivateLinkResourceProperties = original.PrivateLinkResourceProperties
type PrivateLinkResourcesClient = original.PrivateLinkResourcesClient
type PrivateLinkServiceConnectionState = original.PrivateLinkServiceConnectionState
type Properties = original.Properties
type ProxyResource = original.ProxyResource
type RegenerateKeyFuture = original.RegenerateKeyFuture
type RegenerateKeyParameters = original.RegenerateKeyParameters
type Resource = original.Resource
type ResourceList = original.ResourceList
type ResourceListIterator = original.ResourceListIterator
type ResourceListPage = original.ResourceListPage
type ResourceSku = original.ResourceSku
type ResourceType = original.ResourceType
type RestartFuture = original.RestartFuture
type ServerlessUpstreamSettings = original.ServerlessUpstreamSettings
type ServiceSpecification = original.ServiceSpecification
type ShareablePrivateLinkResourceProperties = original.ShareablePrivateLinkResourceProperties
type ShareablePrivateLinkResourceType = original.ShareablePrivateLinkResourceType
type SharedPrivateLinkResource = original.SharedPrivateLinkResource
type SharedPrivateLinkResourceList = original.SharedPrivateLinkResourceList
type SharedPrivateLinkResourceListIterator = original.SharedPrivateLinkResourceListIterator
type SharedPrivateLinkResourceListPage = original.SharedPrivateLinkResourceListPage
type SharedPrivateLinkResourceProperties = original.SharedPrivateLinkResourceProperties
type SharedPrivateLinkResourcesClient = original.SharedPrivateLinkResourcesClient
type SharedPrivateLinkResourcesCreateOrUpdateFuture = original.SharedPrivateLinkResourcesCreateOrUpdateFuture
type SharedPrivateLinkResourcesDeleteFuture = original.SharedPrivateLinkResourcesDeleteFuture
type SystemData = original.SystemData
type TLSSettings = original.TLSSettings
type TrackedResource = original.TrackedResource
type UpdateFuture = original.UpdateFuture
type UpstreamAuthSettings = original.UpstreamAuthSettings
type UpstreamTemplate = original.UpstreamTemplate
type Usage = original.Usage
type UsageList = original.UsageList
type UsageListIterator = original.UsageListIterator
type UsageListPage = original.UsageListPage
type UsageName = original.UsageName
type UsagesClient = original.UsagesClient
type UserAssignedIdentityProperty = original.UserAssignedIdentityProperty

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewClient(subscriptionID string) Client {
	return original.NewClient(subscriptionID)
}
func NewClientWithBaseURI(baseURI string, subscriptionID string) Client {
	return original.NewClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationListIterator(page OperationListPage) OperationListIterator {
	return original.NewOperationListIterator(page)
}
func NewOperationListPage(cur OperationList, getNextPage func(context.Context, OperationList) (OperationList, error)) OperationListPage {
	return original.NewOperationListPage(cur, getNextPage)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewPrivateEndpointConnectionListIterator(page PrivateEndpointConnectionListPage) PrivateEndpointConnectionListIterator {
	return original.NewPrivateEndpointConnectionListIterator(page)
}
func NewPrivateEndpointConnectionListPage(cur PrivateEndpointConnectionList, getNextPage func(context.Context, PrivateEndpointConnectionList) (PrivateEndpointConnectionList, error)) PrivateEndpointConnectionListPage {
	return original.NewPrivateEndpointConnectionListPage(cur, getNextPage)
}
func NewPrivateEndpointConnectionsClient(subscriptionID string) PrivateEndpointConnectionsClient {
	return original.NewPrivateEndpointConnectionsClient(subscriptionID)
}
func NewPrivateEndpointConnectionsClientWithBaseURI(baseURI string, subscriptionID string) PrivateEndpointConnectionsClient {
	return original.NewPrivateEndpointConnectionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewPrivateLinkResourceListIterator(page PrivateLinkResourceListPage) PrivateLinkResourceListIterator {
	return original.NewPrivateLinkResourceListIterator(page)
}
func NewPrivateLinkResourceListPage(cur PrivateLinkResourceList, getNextPage func(context.Context, PrivateLinkResourceList) (PrivateLinkResourceList, error)) PrivateLinkResourceListPage {
	return original.NewPrivateLinkResourceListPage(cur, getNextPage)
}
func NewPrivateLinkResourcesClient(subscriptionID string) PrivateLinkResourcesClient {
	return original.NewPrivateLinkResourcesClient(subscriptionID)
}
func NewPrivateLinkResourcesClientWithBaseURI(baseURI string, subscriptionID string) PrivateLinkResourcesClient {
	return original.NewPrivateLinkResourcesClientWithBaseURI(baseURI, subscriptionID)
}
func NewResourceListIterator(page ResourceListPage) ResourceListIterator {
	return original.NewResourceListIterator(page)
}
func NewResourceListPage(cur ResourceList, getNextPage func(context.Context, ResourceList) (ResourceList, error)) ResourceListPage {
	return original.NewResourceListPage(cur, getNextPage)
}
func NewSharedPrivateLinkResourceListIterator(page SharedPrivateLinkResourceListPage) SharedPrivateLinkResourceListIterator {
	return original.NewSharedPrivateLinkResourceListIterator(page)
}
func NewSharedPrivateLinkResourceListPage(cur SharedPrivateLinkResourceList, getNextPage func(context.Context, SharedPrivateLinkResourceList) (SharedPrivateLinkResourceList, error)) SharedPrivateLinkResourceListPage {
	return original.NewSharedPrivateLinkResourceListPage(cur, getNextPage)
}
func NewSharedPrivateLinkResourcesClient(subscriptionID string) SharedPrivateLinkResourcesClient {
	return original.NewSharedPrivateLinkResourcesClient(subscriptionID)
}
func NewSharedPrivateLinkResourcesClientWithBaseURI(baseURI string, subscriptionID string) SharedPrivateLinkResourcesClient {
	return original.NewSharedPrivateLinkResourcesClientWithBaseURI(baseURI, subscriptionID)
}
func NewUsageListIterator(page UsageListPage) UsageListIterator {
	return original.NewUsageListIterator(page)
}
func NewUsageListPage(cur UsageList, getNextPage func(context.Context, UsageList) (UsageList, error)) UsageListPage {
	return original.NewUsageListPage(cur, getNextPage)
}
func NewUsagesClient(subscriptionID string) UsagesClient {
	return original.NewUsagesClient(subscriptionID)
}
func NewUsagesClientWithBaseURI(baseURI string, subscriptionID string) UsagesClient {
	return original.NewUsagesClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleACLActionValues() []ACLAction {
	return original.PossibleACLActionValues()
}
func PossibleCreatedByTypeValues() []CreatedByType {
	return original.PossibleCreatedByTypeValues()
}
func PossibleFeatureFlagsValues() []FeatureFlags {
	return original.PossibleFeatureFlagsValues()
}
func PossibleKeyTypeValues() []KeyType {
	return original.PossibleKeyTypeValues()
}
func PossibleManagedIdentityTypeValues() []ManagedIdentityType {
	return original.PossibleManagedIdentityTypeValues()
}
func PossiblePrivateLinkServiceConnectionStatusValues() []PrivateLinkServiceConnectionStatus {
	return original.PossiblePrivateLinkServiceConnectionStatusValues()
}
func PossibleProvisioningStateValues() []ProvisioningState {
	return original.PossibleProvisioningStateValues()
}
func PossibleRequestTypeValues() []RequestType {
	return original.PossibleRequestTypeValues()
}
func PossibleServiceKindValues() []ServiceKind {
	return original.PossibleServiceKindValues()
}
func PossibleSharedPrivateLinkResourceStatusValues() []SharedPrivateLinkResourceStatus {
	return original.PossibleSharedPrivateLinkResourceStatusValues()
}
func PossibleSkuTierValues() []SkuTier {
	return original.PossibleSkuTierValues()
}
func PossibleUpstreamAuthTypeValues() []UpstreamAuthType {
	return original.PossibleUpstreamAuthTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

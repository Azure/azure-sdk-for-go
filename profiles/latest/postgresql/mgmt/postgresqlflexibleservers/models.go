//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package postgresqlflexibleservers

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2021-06-01/postgresqlflexibleservers"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type ConfigurationDataType = original.ConfigurationDataType

const (
	ConfigurationDataTypeBoolean     ConfigurationDataType = original.ConfigurationDataTypeBoolean
	ConfigurationDataTypeEnumeration ConfigurationDataType = original.ConfigurationDataTypeEnumeration
	ConfigurationDataTypeInteger     ConfigurationDataType = original.ConfigurationDataTypeInteger
	ConfigurationDataTypeNumeric     ConfigurationDataType = original.ConfigurationDataTypeNumeric
)

type CreateMode = original.CreateMode

const (
	CreateModeCreate             CreateMode = original.CreateModeCreate
	CreateModeDefault            CreateMode = original.CreateModeDefault
	CreateModePointInTimeRestore CreateMode = original.CreateModePointInTimeRestore
	CreateModeUpdate             CreateMode = original.CreateModeUpdate
)

type CreateModeForUpdate = original.CreateModeForUpdate

const (
	CreateModeForUpdateDefault CreateModeForUpdate = original.CreateModeForUpdateDefault
	CreateModeForUpdateUpdate  CreateModeForUpdate = original.CreateModeForUpdateUpdate
)

type CreatedByType = original.CreatedByType

const (
	CreatedByTypeApplication     CreatedByType = original.CreatedByTypeApplication
	CreatedByTypeKey             CreatedByType = original.CreatedByTypeKey
	CreatedByTypeManagedIdentity CreatedByType = original.CreatedByTypeManagedIdentity
	CreatedByTypeUser            CreatedByType = original.CreatedByTypeUser
)

type FailoverMode = original.FailoverMode

const (
	FailoverModeForcedFailover    FailoverMode = original.FailoverModeForcedFailover
	FailoverModeForcedSwitchover  FailoverMode = original.FailoverModeForcedSwitchover
	FailoverModePlannedFailover   FailoverMode = original.FailoverModePlannedFailover
	FailoverModePlannedSwitchover FailoverMode = original.FailoverModePlannedSwitchover
)

type GeoRedundantBackupEnum = original.GeoRedundantBackupEnum

const (
	GeoRedundantBackupEnumDisabled GeoRedundantBackupEnum = original.GeoRedundantBackupEnumDisabled
	GeoRedundantBackupEnumEnabled  GeoRedundantBackupEnum = original.GeoRedundantBackupEnumEnabled
)

type HighAvailabilityMode = original.HighAvailabilityMode

const (
	HighAvailabilityModeDisabled      HighAvailabilityMode = original.HighAvailabilityModeDisabled
	HighAvailabilityModeZoneRedundant HighAvailabilityMode = original.HighAvailabilityModeZoneRedundant
)

type OperationOrigin = original.OperationOrigin

const (
	OperationOriginNotSpecified OperationOrigin = original.OperationOriginNotSpecified
	OperationOriginSystem       OperationOrigin = original.OperationOriginSystem
	OperationOriginUser         OperationOrigin = original.OperationOriginUser
)

type ServerHAState = original.ServerHAState

const (
	ServerHAStateCreatingStandby ServerHAState = original.ServerHAStateCreatingStandby
	ServerHAStateFailingOver     ServerHAState = original.ServerHAStateFailingOver
	ServerHAStateHealthy         ServerHAState = original.ServerHAStateHealthy
	ServerHAStateNotEnabled      ServerHAState = original.ServerHAStateNotEnabled
	ServerHAStateRemovingStandby ServerHAState = original.ServerHAStateRemovingStandby
	ServerHAStateReplicatingData ServerHAState = original.ServerHAStateReplicatingData
)

type ServerPublicNetworkAccessState = original.ServerPublicNetworkAccessState

const (
	ServerPublicNetworkAccessStateDisabled ServerPublicNetworkAccessState = original.ServerPublicNetworkAccessStateDisabled
	ServerPublicNetworkAccessStateEnabled  ServerPublicNetworkAccessState = original.ServerPublicNetworkAccessStateEnabled
)

type ServerState = original.ServerState

const (
	ServerStateDisabled ServerState = original.ServerStateDisabled
	ServerStateDropping ServerState = original.ServerStateDropping
	ServerStateReady    ServerState = original.ServerStateReady
	ServerStateStarting ServerState = original.ServerStateStarting
	ServerStateStopped  ServerState = original.ServerStateStopped
	ServerStateStopping ServerState = original.ServerStateStopping
	ServerStateUpdating ServerState = original.ServerStateUpdating
)

type ServerVersion = original.ServerVersion

const (
	ServerVersionOneOne   ServerVersion = original.ServerVersionOneOne
	ServerVersionOneThree ServerVersion = original.ServerVersionOneThree
	ServerVersionOneTwo   ServerVersion = original.ServerVersionOneTwo
)

type SkuTier = original.SkuTier

const (
	SkuTierBurstable       SkuTier = original.SkuTierBurstable
	SkuTierGeneralPurpose  SkuTier = original.SkuTierGeneralPurpose
	SkuTierMemoryOptimized SkuTier = original.SkuTierMemoryOptimized
)

type AzureEntityResource = original.AzureEntityResource
type Backup = original.Backup
type BaseClient = original.BaseClient
type CapabilitiesListResult = original.CapabilitiesListResult
type CapabilitiesListResultIterator = original.CapabilitiesListResultIterator
type CapabilitiesListResultPage = original.CapabilitiesListResultPage
type CapabilityProperties = original.CapabilityProperties
type CheckNameAvailabilityClient = original.CheckNameAvailabilityClient
type CloudError = original.CloudError
type Configuration = original.Configuration
type ConfigurationListResult = original.ConfigurationListResult
type ConfigurationListResultIterator = original.ConfigurationListResultIterator
type ConfigurationListResultPage = original.ConfigurationListResultPage
type ConfigurationProperties = original.ConfigurationProperties
type ConfigurationsClient = original.ConfigurationsClient
type ConfigurationsPutFuture = original.ConfigurationsPutFuture
type ConfigurationsUpdateFuture = original.ConfigurationsUpdateFuture
type Database = original.Database
type DatabaseListResult = original.DatabaseListResult
type DatabaseListResultIterator = original.DatabaseListResultIterator
type DatabaseListResultPage = original.DatabaseListResultPage
type DatabaseProperties = original.DatabaseProperties
type DatabasesClient = original.DatabasesClient
type DatabasesCreateFuture = original.DatabasesCreateFuture
type DatabasesDeleteFuture = original.DatabasesDeleteFuture
type DelegatedSubnetUsage = original.DelegatedSubnetUsage
type ErrorAdditionalInfo = original.ErrorAdditionalInfo
type ErrorResponse = original.ErrorResponse
type FirewallRule = original.FirewallRule
type FirewallRuleListResult = original.FirewallRuleListResult
type FirewallRuleListResultIterator = original.FirewallRuleListResultIterator
type FirewallRuleListResultPage = original.FirewallRuleListResultPage
type FirewallRuleProperties = original.FirewallRuleProperties
type FirewallRulesClient = original.FirewallRulesClient
type FirewallRulesCreateOrUpdateFuture = original.FirewallRulesCreateOrUpdateFuture
type FirewallRulesDeleteFuture = original.FirewallRulesDeleteFuture
type FlexibleServerEditionCapability = original.FlexibleServerEditionCapability
type GetPrivateDNSZoneSuffixClient = original.GetPrivateDNSZoneSuffixClient
type HighAvailability = original.HighAvailability
type HyperscaleNodeEditionCapability = original.HyperscaleNodeEditionCapability
type LocationBasedCapabilitiesClient = original.LocationBasedCapabilitiesClient
type MaintenanceWindow = original.MaintenanceWindow
type NameAvailability = original.NameAvailability
type NameAvailabilityRequest = original.NameAvailabilityRequest
type Network = original.Network
type NodeTypeCapability = original.NodeTypeCapability
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationsClient = original.OperationsClient
type ProxyResource = original.ProxyResource
type Resource = original.Resource
type RestartParameter = original.RestartParameter
type Server = original.Server
type ServerForUpdate = original.ServerForUpdate
type ServerListResult = original.ServerListResult
type ServerListResultIterator = original.ServerListResultIterator
type ServerListResultPage = original.ServerListResultPage
type ServerProperties = original.ServerProperties
type ServerPropertiesForUpdate = original.ServerPropertiesForUpdate
type ServerVersionCapability = original.ServerVersionCapability
type ServersClient = original.ServersClient
type ServersCreateFuture = original.ServersCreateFuture
type ServersDeleteFuture = original.ServersDeleteFuture
type ServersRestartFuture = original.ServersRestartFuture
type ServersStartFuture = original.ServersStartFuture
type ServersStopFuture = original.ServersStopFuture
type ServersUpdateFuture = original.ServersUpdateFuture
type Sku = original.Sku
type Storage = original.Storage
type StorageEditionCapability = original.StorageEditionCapability
type StorageMBCapability = original.StorageMBCapability
type String = original.String
type SystemData = original.SystemData
type TrackedResource = original.TrackedResource
type VcoreCapability = original.VcoreCapability
type VirtualNetworkSubnetUsageClient = original.VirtualNetworkSubnetUsageClient
type VirtualNetworkSubnetUsageParameter = original.VirtualNetworkSubnetUsageParameter
type VirtualNetworkSubnetUsageResult = original.VirtualNetworkSubnetUsageResult

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewCapabilitiesListResultIterator(page CapabilitiesListResultPage) CapabilitiesListResultIterator {
	return original.NewCapabilitiesListResultIterator(page)
}
func NewCapabilitiesListResultPage(cur CapabilitiesListResult, getNextPage func(context.Context, CapabilitiesListResult) (CapabilitiesListResult, error)) CapabilitiesListResultPage {
	return original.NewCapabilitiesListResultPage(cur, getNextPage)
}
func NewCheckNameAvailabilityClient(subscriptionID string) CheckNameAvailabilityClient {
	return original.NewCheckNameAvailabilityClient(subscriptionID)
}
func NewCheckNameAvailabilityClientWithBaseURI(baseURI string, subscriptionID string) CheckNameAvailabilityClient {
	return original.NewCheckNameAvailabilityClientWithBaseURI(baseURI, subscriptionID)
}
func NewConfigurationListResultIterator(page ConfigurationListResultPage) ConfigurationListResultIterator {
	return original.NewConfigurationListResultIterator(page)
}
func NewConfigurationListResultPage(cur ConfigurationListResult, getNextPage func(context.Context, ConfigurationListResult) (ConfigurationListResult, error)) ConfigurationListResultPage {
	return original.NewConfigurationListResultPage(cur, getNextPage)
}
func NewConfigurationsClient(subscriptionID string) ConfigurationsClient {
	return original.NewConfigurationsClient(subscriptionID)
}
func NewConfigurationsClientWithBaseURI(baseURI string, subscriptionID string) ConfigurationsClient {
	return original.NewConfigurationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewDatabaseListResultIterator(page DatabaseListResultPage) DatabaseListResultIterator {
	return original.NewDatabaseListResultIterator(page)
}
func NewDatabaseListResultPage(cur DatabaseListResult, getNextPage func(context.Context, DatabaseListResult) (DatabaseListResult, error)) DatabaseListResultPage {
	return original.NewDatabaseListResultPage(cur, getNextPage)
}
func NewDatabasesClient(subscriptionID string) DatabasesClient {
	return original.NewDatabasesClient(subscriptionID)
}
func NewDatabasesClientWithBaseURI(baseURI string, subscriptionID string) DatabasesClient {
	return original.NewDatabasesClientWithBaseURI(baseURI, subscriptionID)
}
func NewFirewallRuleListResultIterator(page FirewallRuleListResultPage) FirewallRuleListResultIterator {
	return original.NewFirewallRuleListResultIterator(page)
}
func NewFirewallRuleListResultPage(cur FirewallRuleListResult, getNextPage func(context.Context, FirewallRuleListResult) (FirewallRuleListResult, error)) FirewallRuleListResultPage {
	return original.NewFirewallRuleListResultPage(cur, getNextPage)
}
func NewFirewallRulesClient(subscriptionID string) FirewallRulesClient {
	return original.NewFirewallRulesClient(subscriptionID)
}
func NewFirewallRulesClientWithBaseURI(baseURI string, subscriptionID string) FirewallRulesClient {
	return original.NewFirewallRulesClientWithBaseURI(baseURI, subscriptionID)
}
func NewGetPrivateDNSZoneSuffixClient(subscriptionID string) GetPrivateDNSZoneSuffixClient {
	return original.NewGetPrivateDNSZoneSuffixClient(subscriptionID)
}
func NewGetPrivateDNSZoneSuffixClientWithBaseURI(baseURI string, subscriptionID string) GetPrivateDNSZoneSuffixClient {
	return original.NewGetPrivateDNSZoneSuffixClientWithBaseURI(baseURI, subscriptionID)
}
func NewLocationBasedCapabilitiesClient(subscriptionID string) LocationBasedCapabilitiesClient {
	return original.NewLocationBasedCapabilitiesClient(subscriptionID)
}
func NewLocationBasedCapabilitiesClientWithBaseURI(baseURI string, subscriptionID string) LocationBasedCapabilitiesClient {
	return original.NewLocationBasedCapabilitiesClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewServerListResultIterator(page ServerListResultPage) ServerListResultIterator {
	return original.NewServerListResultIterator(page)
}
func NewServerListResultPage(cur ServerListResult, getNextPage func(context.Context, ServerListResult) (ServerListResult, error)) ServerListResultPage {
	return original.NewServerListResultPage(cur, getNextPage)
}
func NewServersClient(subscriptionID string) ServersClient {
	return original.NewServersClient(subscriptionID)
}
func NewServersClientWithBaseURI(baseURI string, subscriptionID string) ServersClient {
	return original.NewServersClientWithBaseURI(baseURI, subscriptionID)
}
func NewVirtualNetworkSubnetUsageClient(subscriptionID string) VirtualNetworkSubnetUsageClient {
	return original.NewVirtualNetworkSubnetUsageClient(subscriptionID)
}
func NewVirtualNetworkSubnetUsageClientWithBaseURI(baseURI string, subscriptionID string) VirtualNetworkSubnetUsageClient {
	return original.NewVirtualNetworkSubnetUsageClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleConfigurationDataTypeValues() []ConfigurationDataType {
	return original.PossibleConfigurationDataTypeValues()
}
func PossibleCreateModeForUpdateValues() []CreateModeForUpdate {
	return original.PossibleCreateModeForUpdateValues()
}
func PossibleCreateModeValues() []CreateMode {
	return original.PossibleCreateModeValues()
}
func PossibleCreatedByTypeValues() []CreatedByType {
	return original.PossibleCreatedByTypeValues()
}
func PossibleFailoverModeValues() []FailoverMode {
	return original.PossibleFailoverModeValues()
}
func PossibleGeoRedundantBackupEnumValues() []GeoRedundantBackupEnum {
	return original.PossibleGeoRedundantBackupEnumValues()
}
func PossibleHighAvailabilityModeValues() []HighAvailabilityMode {
	return original.PossibleHighAvailabilityModeValues()
}
func PossibleOperationOriginValues() []OperationOrigin {
	return original.PossibleOperationOriginValues()
}
func PossibleServerHAStateValues() []ServerHAState {
	return original.PossibleServerHAStateValues()
}
func PossibleServerPublicNetworkAccessStateValues() []ServerPublicNetworkAccessState {
	return original.PossibleServerPublicNetworkAccessStateValues()
}
func PossibleServerStateValues() []ServerState {
	return original.PossibleServerStateValues()
}
func PossibleServerVersionValues() []ServerVersion {
	return original.PossibleServerVersionValues()
}
func PossibleSkuTierValues() []SkuTier {
	return original.PossibleSkuTierValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}

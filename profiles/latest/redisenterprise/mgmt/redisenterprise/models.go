// +build go1.9

// Copyright 2021 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package redisenterprise

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/redisenterprise/mgmt/2021-03-01/redisenterprise"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type AccessKeyType = original.AccessKeyType

const (
	Primary   AccessKeyType = original.Primary
	Secondary AccessKeyType = original.Secondary
)

type ActionType = original.ActionType

const (
	Internal ActionType = original.Internal
)

type AofFrequency = original.AofFrequency

const (
	Always AofFrequency = original.Always
	Ones   AofFrequency = original.Ones
)

type ClusteringPolicy = original.ClusteringPolicy

const (
	EnterpriseCluster ClusteringPolicy = original.EnterpriseCluster
	OSSCluster        ClusteringPolicy = original.OSSCluster
)

type EvictionPolicy = original.EvictionPolicy

const (
	AllKeysLFU     EvictionPolicy = original.AllKeysLFU
	AllKeysLRU     EvictionPolicy = original.AllKeysLRU
	AllKeysRandom  EvictionPolicy = original.AllKeysRandom
	NoEviction     EvictionPolicy = original.NoEviction
	VolatileLFU    EvictionPolicy = original.VolatileLFU
	VolatileLRU    EvictionPolicy = original.VolatileLRU
	VolatileRandom EvictionPolicy = original.VolatileRandom
	VolatileTTL    EvictionPolicy = original.VolatileTTL
)

type Origin = original.Origin

const (
	System     Origin = original.System
	User       Origin = original.User
	Usersystem Origin = original.Usersystem
)

type PrivateEndpointConnectionProvisioningState = original.PrivateEndpointConnectionProvisioningState

const (
	Creating  PrivateEndpointConnectionProvisioningState = original.Creating
	Deleting  PrivateEndpointConnectionProvisioningState = original.Deleting
	Failed    PrivateEndpointConnectionProvisioningState = original.Failed
	Succeeded PrivateEndpointConnectionProvisioningState = original.Succeeded
)

type PrivateEndpointServiceConnectionStatus = original.PrivateEndpointServiceConnectionStatus

const (
	Approved PrivateEndpointServiceConnectionStatus = original.Approved
	Pending  PrivateEndpointServiceConnectionStatus = original.Pending
	Rejected PrivateEndpointServiceConnectionStatus = original.Rejected
)

type Protocol = original.Protocol

const (
	Encrypted Protocol = original.Encrypted
	Plaintext Protocol = original.Plaintext
)

type ProvisioningState = original.ProvisioningState

const (
	ProvisioningStateCanceled  ProvisioningState = original.ProvisioningStateCanceled
	ProvisioningStateCreating  ProvisioningState = original.ProvisioningStateCreating
	ProvisioningStateDeleting  ProvisioningState = original.ProvisioningStateDeleting
	ProvisioningStateFailed    ProvisioningState = original.ProvisioningStateFailed
	ProvisioningStateSucceeded ProvisioningState = original.ProvisioningStateSucceeded
	ProvisioningStateUpdating  ProvisioningState = original.ProvisioningStateUpdating
)

type RdbFrequency = original.RdbFrequency

const (
	Oneh    RdbFrequency = original.Oneh
	OneTwoh RdbFrequency = original.OneTwoh
	Sixh    RdbFrequency = original.Sixh
)

type ResourceState = original.ResourceState

const (
	ResourceStateCreateFailed  ResourceState = original.ResourceStateCreateFailed
	ResourceStateCreating      ResourceState = original.ResourceStateCreating
	ResourceStateDeleteFailed  ResourceState = original.ResourceStateDeleteFailed
	ResourceStateDeleting      ResourceState = original.ResourceStateDeleting
	ResourceStateDisabled      ResourceState = original.ResourceStateDisabled
	ResourceStateDisableFailed ResourceState = original.ResourceStateDisableFailed
	ResourceStateDisabling     ResourceState = original.ResourceStateDisabling
	ResourceStateEnableFailed  ResourceState = original.ResourceStateEnableFailed
	ResourceStateEnabling      ResourceState = original.ResourceStateEnabling
	ResourceStateRunning       ResourceState = original.ResourceStateRunning
	ResourceStateUpdateFailed  ResourceState = original.ResourceStateUpdateFailed
	ResourceStateUpdating      ResourceState = original.ResourceStateUpdating
)

type SkuName = original.SkuName

const (
	EnterpriseE10        SkuName = original.EnterpriseE10
	EnterpriseE100       SkuName = original.EnterpriseE100
	EnterpriseE20        SkuName = original.EnterpriseE20
	EnterpriseE50        SkuName = original.EnterpriseE50
	EnterpriseFlashF1500 SkuName = original.EnterpriseFlashF1500
	EnterpriseFlashF300  SkuName = original.EnterpriseFlashF300
	EnterpriseFlashF700  SkuName = original.EnterpriseFlashF700
)

type TLSVersion = original.TLSVersion

const (
	OneFullStopOne  TLSVersion = original.OneFullStopOne
	OneFullStopTwo  TLSVersion = original.OneFullStopTwo
	OneFullStopZero TLSVersion = original.OneFullStopZero
)

type AccessKeys = original.AccessKeys
type AzureEntityResource = original.AzureEntityResource
type BaseClient = original.BaseClient
type Client = original.Client
type Cluster = original.Cluster
type ClusterList = original.ClusterList
type ClusterListIterator = original.ClusterListIterator
type ClusterListPage = original.ClusterListPage
type ClusterProperties = original.ClusterProperties
type ClusterUpdate = original.ClusterUpdate
type CreateFuture = original.CreateFuture
type Database = original.Database
type DatabaseList = original.DatabaseList
type DatabaseListIterator = original.DatabaseListIterator
type DatabaseListPage = original.DatabaseListPage
type DatabaseProperties = original.DatabaseProperties
type DatabaseUpdate = original.DatabaseUpdate
type DatabasesClient = original.DatabasesClient
type DatabasesCreateFuture = original.DatabasesCreateFuture
type DatabasesDeleteFuture = original.DatabasesDeleteFuture
type DatabasesExportFuture = original.DatabasesExportFuture
type DatabasesImportFuture = original.DatabasesImportFuture
type DatabasesRegenerateKeyFuture = original.DatabasesRegenerateKeyFuture
type DatabasesUpdateFuture = original.DatabasesUpdateFuture
type DeleteFuture = original.DeleteFuture
type ErrorAdditionalInfo = original.ErrorAdditionalInfo
type ErrorDetail = original.ErrorDetail
type ErrorResponse = original.ErrorResponse
type ExportClusterParameters = original.ExportClusterParameters
type ImportClusterParameters = original.ImportClusterParameters
type Module = original.Module
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationStatus = original.OperationStatus
type OperationsClient = original.OperationsClient
type OperationsStatusClient = original.OperationsStatusClient
type Persistence = original.Persistence
type PrivateEndpoint = original.PrivateEndpoint
type PrivateEndpointConnection = original.PrivateEndpointConnection
type PrivateEndpointConnectionListResult = original.PrivateEndpointConnectionListResult
type PrivateEndpointConnectionProperties = original.PrivateEndpointConnectionProperties
type PrivateEndpointConnectionsClient = original.PrivateEndpointConnectionsClient
type PrivateEndpointConnectionsPutFuture = original.PrivateEndpointConnectionsPutFuture
type PrivateLinkResource = original.PrivateLinkResource
type PrivateLinkResourceListResult = original.PrivateLinkResourceListResult
type PrivateLinkResourceProperties = original.PrivateLinkResourceProperties
type PrivateLinkResourcesClient = original.PrivateLinkResourcesClient
type PrivateLinkServiceConnectionState = original.PrivateLinkServiceConnectionState
type ProxyResource = original.ProxyResource
type RegenerateKeyParameters = original.RegenerateKeyParameters
type Resource = original.Resource
type Sku = original.Sku
type TrackedResource = original.TrackedResource
type UpdateFuture = original.UpdateFuture

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewClient(subscriptionID string) Client {
	return original.NewClient(subscriptionID)
}
func NewClientWithBaseURI(baseURI string, subscriptionID string) Client {
	return original.NewClientWithBaseURI(baseURI, subscriptionID)
}
func NewClusterListIterator(page ClusterListPage) ClusterListIterator {
	return original.NewClusterListIterator(page)
}
func NewClusterListPage(cur ClusterList, getNextPage func(context.Context, ClusterList) (ClusterList, error)) ClusterListPage {
	return original.NewClusterListPage(cur, getNextPage)
}
func NewDatabaseListIterator(page DatabaseListPage) DatabaseListIterator {
	return original.NewDatabaseListIterator(page)
}
func NewDatabaseListPage(cur DatabaseList, getNextPage func(context.Context, DatabaseList) (DatabaseList, error)) DatabaseListPage {
	return original.NewDatabaseListPage(cur, getNextPage)
}
func NewDatabasesClient(subscriptionID string) DatabasesClient {
	return original.NewDatabasesClient(subscriptionID)
}
func NewDatabasesClientWithBaseURI(baseURI string, subscriptionID string) DatabasesClient {
	return original.NewDatabasesClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationListResultIterator(page OperationListResultPage) OperationListResultIterator {
	return original.NewOperationListResultIterator(page)
}
func NewOperationListResultPage(cur OperationListResult, getNextPage func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage {
	return original.NewOperationListResultPage(cur, getNextPage)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationsStatusClient(subscriptionID string) OperationsStatusClient {
	return original.NewOperationsStatusClient(subscriptionID)
}
func NewOperationsStatusClientWithBaseURI(baseURI string, subscriptionID string) OperationsStatusClient {
	return original.NewOperationsStatusClientWithBaseURI(baseURI, subscriptionID)
}
func NewPrivateEndpointConnectionsClient(subscriptionID string) PrivateEndpointConnectionsClient {
	return original.NewPrivateEndpointConnectionsClient(subscriptionID)
}
func NewPrivateEndpointConnectionsClientWithBaseURI(baseURI string, subscriptionID string) PrivateEndpointConnectionsClient {
	return original.NewPrivateEndpointConnectionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewPrivateLinkResourcesClient(subscriptionID string) PrivateLinkResourcesClient {
	return original.NewPrivateLinkResourcesClient(subscriptionID)
}
func NewPrivateLinkResourcesClientWithBaseURI(baseURI string, subscriptionID string) PrivateLinkResourcesClient {
	return original.NewPrivateLinkResourcesClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleAccessKeyTypeValues() []AccessKeyType {
	return original.PossibleAccessKeyTypeValues()
}
func PossibleActionTypeValues() []ActionType {
	return original.PossibleActionTypeValues()
}
func PossibleAofFrequencyValues() []AofFrequency {
	return original.PossibleAofFrequencyValues()
}
func PossibleClusteringPolicyValues() []ClusteringPolicy {
	return original.PossibleClusteringPolicyValues()
}
func PossibleEvictionPolicyValues() []EvictionPolicy {
	return original.PossibleEvictionPolicyValues()
}
func PossibleOriginValues() []Origin {
	return original.PossibleOriginValues()
}
func PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState {
	return original.PossiblePrivateEndpointConnectionProvisioningStateValues()
}
func PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus {
	return original.PossiblePrivateEndpointServiceConnectionStatusValues()
}
func PossibleProtocolValues() []Protocol {
	return original.PossibleProtocolValues()
}
func PossibleProvisioningStateValues() []ProvisioningState {
	return original.PossibleProvisioningStateValues()
}
func PossibleRdbFrequencyValues() []RdbFrequency {
	return original.PossibleRdbFrequencyValues()
}
func PossibleResourceStateValues() []ResourceState {
	return original.PossibleResourceStateValues()
}
func PossibleSkuNameValues() []SkuName {
	return original.PossibleSkuNameValues()
}
func PossibleTLSVersionValues() []TLSVersion {
	return original.PossibleTLSVersionValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}

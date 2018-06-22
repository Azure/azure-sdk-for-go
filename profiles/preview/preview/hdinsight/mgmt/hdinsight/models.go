// +build go1.9

// Copyright 2018 Microsoft Corporation
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

package hdinsight

import original "github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"

type ApplicationsClient = original.ApplicationsClient

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type ClustersClient = original.ClustersClient
type ConfigurationsClient = original.ConfigurationsClient
type ExtensionsClient = original.ExtensionsClient
type LocationClient = original.LocationClient
type AsyncOperationState = original.AsyncOperationState

const (
	Failed     AsyncOperationState = original.Failed
	InProgress AsyncOperationState = original.InProgress
	Succeeded  AsyncOperationState = original.Succeeded
)

type ClusterProvisioningState = original.ClusterProvisioningState

const (
	ClusterProvisioningStateCanceled   ClusterProvisioningState = original.ClusterProvisioningStateCanceled
	ClusterProvisioningStateDeleting   ClusterProvisioningState = original.ClusterProvisioningStateDeleting
	ClusterProvisioningStateFailed     ClusterProvisioningState = original.ClusterProvisioningStateFailed
	ClusterProvisioningStateInProgress ClusterProvisioningState = original.ClusterProvisioningStateInProgress
	ClusterProvisioningStateSucceeded  ClusterProvisioningState = original.ClusterProvisioningStateSucceeded
)

type DirectoryType = original.DirectoryType

const (
	ActiveDirectory DirectoryType = original.ActiveDirectory
)

type OSType = original.OSType

const (
	Linux   OSType = original.Linux
	Windows OSType = original.Windows
)

type Tier = original.Tier

const (
	Premium  Tier = original.Premium
	Standard Tier = original.Standard
)

type Application = original.Application
type ApplicationGetEndpoint = original.ApplicationGetEndpoint
type ApplicationGetHTTPSEndpoint = original.ApplicationGetHTTPSEndpoint
type ApplicationListResult = original.ApplicationListResult
type ApplicationListResultIterator = original.ApplicationListResultIterator
type ApplicationListResultPage = original.ApplicationListResultPage
type ApplicationProperties = original.ApplicationProperties
type ApplicationsCreateFuture = original.ApplicationsCreateFuture
type ApplicationsDeleteFuture = original.ApplicationsDeleteFuture
type Cluster = original.Cluster
type ClusterCreateParametersExtended = original.ClusterCreateParametersExtended
type ClusterCreateProperties = original.ClusterCreateProperties
type ClusterDefinition = original.ClusterDefinition
type ClusterGetProperties = original.ClusterGetProperties
type ClusterListPersistedScriptActionsResult = original.ClusterListPersistedScriptActionsResult
type ClusterListResult = original.ClusterListResult
type ClusterListResultIterator = original.ClusterListResultIterator
type ClusterListResultPage = original.ClusterListResultPage
type ClusterListRuntimeScriptActionDetailResult = original.ClusterListRuntimeScriptActionDetailResult
type ClusterMonitoringRequest = original.ClusterMonitoringRequest
type ClusterMonitoringResponse = original.ClusterMonitoringResponse
type ClusterPatchParameters = original.ClusterPatchParameters
type ClusterResizeParameters = original.ClusterResizeParameters
type ClustersCreateFuture = original.ClustersCreateFuture
type ClustersDeleteFuture = original.ClustersDeleteFuture
type ClustersExecuteScriptActionsFuture = original.ClustersExecuteScriptActionsFuture
type ClustersResizeFuture = original.ClustersResizeFuture
type ComputeProfile = original.ComputeProfile
type ConfigurationsUpdateHTTPSettingsFuture = original.ConfigurationsUpdateHTTPSettingsFuture
type ConnectivityEndpoint = original.ConnectivityEndpoint
type DataDisksGroups = original.DataDisksGroups
type ErrorResponse = original.ErrorResponse
type Errors = original.Errors
type ExecuteScriptActionParameters = original.ExecuteScriptActionParameters
type Extension = original.Extension
type ExtensionsCreateFuture = original.ExtensionsCreateFuture
type ExtensionsDeleteFuture = original.ExtensionsDeleteFuture
type ExtensionsDisableMonitoringFuture = original.ExtensionsDisableMonitoringFuture
type ExtensionsEnableMonitoringFuture = original.ExtensionsEnableMonitoringFuture
type HardwareProfile = original.HardwareProfile
type LinuxOperatingSystemProfile = original.LinuxOperatingSystemProfile
type LocalizedName = original.LocalizedName
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationResource = original.OperationResource
type OsProfile = original.OsProfile
type ProxyResource = original.ProxyResource
type QuotaInfo = original.QuotaInfo
type Resource = original.Resource
type Role = original.Role
type RuntimeScriptAction = original.RuntimeScriptAction
type RuntimeScriptActionDetail = original.RuntimeScriptActionDetail
type ScriptAction = original.ScriptAction
type ScriptActionExecutionHistoryList = original.ScriptActionExecutionHistoryList
type ScriptActionExecutionHistoryListIterator = original.ScriptActionExecutionHistoryListIterator
type ScriptActionExecutionHistoryListPage = original.ScriptActionExecutionHistoryListPage
type ScriptActionExecutionSummary = original.ScriptActionExecutionSummary
type ScriptActionPersistedGetResponseSpec = original.ScriptActionPersistedGetResponseSpec
type ScriptActionsList = original.ScriptActionsList
type ScriptActionsListIterator = original.ScriptActionsListIterator
type ScriptActionsListPage = original.ScriptActionsListPage
type SecurityProfile = original.SecurityProfile
type SetString = original.SetString
type SSHProfile = original.SSHProfile
type SSHPublicKey = original.SSHPublicKey
type StorageAccount = original.StorageAccount
type StorageProfile = original.StorageProfile
type TrackedResource = original.TrackedResource
type Usage = original.Usage
type UsagesListResult = original.UsagesListResult
type VirtualNetworkProfile = original.VirtualNetworkProfile
type OperationsClient = original.OperationsClient
type ScriptActionsClient = original.ScriptActionsClient
type ScriptExecutionHistoryClient = original.ScriptExecutionHistoryClient

func NewApplicationsClient(subscriptionID string) ApplicationsClient {
	return original.NewApplicationsClient(subscriptionID)
}
func NewApplicationsClientWithBaseURI(baseURI string, subscriptionID string) ApplicationsClient {
	return original.NewApplicationsClientWithBaseURI(baseURI, subscriptionID)
}
func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func NewClustersClient(subscriptionID string) ClustersClient {
	return original.NewClustersClient(subscriptionID)
}
func NewClustersClientWithBaseURI(baseURI string, subscriptionID string) ClustersClient {
	return original.NewClustersClientWithBaseURI(baseURI, subscriptionID)
}
func NewConfigurationsClient(subscriptionID string) ConfigurationsClient {
	return original.NewConfigurationsClient(subscriptionID)
}
func NewConfigurationsClientWithBaseURI(baseURI string, subscriptionID string) ConfigurationsClient {
	return original.NewConfigurationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewExtensionsClient(subscriptionID string) ExtensionsClient {
	return original.NewExtensionsClient(subscriptionID)
}
func NewExtensionsClientWithBaseURI(baseURI string, subscriptionID string) ExtensionsClient {
	return original.NewExtensionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewLocationClient(subscriptionID string) LocationClient {
	return original.NewLocationClient(subscriptionID)
}
func NewLocationClientWithBaseURI(baseURI string, subscriptionID string) LocationClient {
	return original.NewLocationClientWithBaseURI(baseURI, subscriptionID)
}
func PossibleAsyncOperationStateValues() []AsyncOperationState {
	return original.PossibleAsyncOperationStateValues()
}
func PossibleClusterProvisioningStateValues() []ClusterProvisioningState {
	return original.PossibleClusterProvisioningStateValues()
}
func PossibleDirectoryTypeValues() []DirectoryType {
	return original.PossibleDirectoryTypeValues()
}
func PossibleOSTypeValues() []OSType {
	return original.PossibleOSTypeValues()
}
func PossibleTierValues() []Tier {
	return original.PossibleTierValues()
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewScriptActionsClient(subscriptionID string) ScriptActionsClient {
	return original.NewScriptActionsClient(subscriptionID)
}
func NewScriptActionsClientWithBaseURI(baseURI string, subscriptionID string) ScriptActionsClient {
	return original.NewScriptActionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewScriptExecutionHistoryClient(subscriptionID string) ScriptExecutionHistoryClient {
	return original.NewScriptExecutionHistoryClient(subscriptionID)
}
func NewScriptExecutionHistoryClientWithBaseURI(baseURI string, subscriptionID string) ScriptExecutionHistoryClient {
	return original.NewScriptExecutionHistoryClientWithBaseURI(baseURI, subscriptionID)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

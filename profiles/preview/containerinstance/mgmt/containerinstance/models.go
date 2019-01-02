// +build go1.9

// Copyright 2019 Microsoft Corporation
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

package containerinstance

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type ContainerGroupIPAddressType = original.ContainerGroupIPAddressType

const (
	Private ContainerGroupIPAddressType = original.Private
	Public  ContainerGroupIPAddressType = original.Public
)

type ContainerGroupNetworkProtocol = original.ContainerGroupNetworkProtocol

const (
	TCP ContainerGroupNetworkProtocol = original.TCP
	UDP ContainerGroupNetworkProtocol = original.UDP
)

type ContainerGroupRestartPolicy = original.ContainerGroupRestartPolicy

const (
	Always    ContainerGroupRestartPolicy = original.Always
	Never     ContainerGroupRestartPolicy = original.Never
	OnFailure ContainerGroupRestartPolicy = original.OnFailure
)

type ContainerNetworkProtocol = original.ContainerNetworkProtocol

const (
	ContainerNetworkProtocolTCP ContainerNetworkProtocol = original.ContainerNetworkProtocolTCP
	ContainerNetworkProtocolUDP ContainerNetworkProtocol = original.ContainerNetworkProtocolUDP
)

type GpuSku = original.GpuSku

const (
	K80  GpuSku = original.K80
	P100 GpuSku = original.P100
	V100 GpuSku = original.V100
)

type LogAnalyticsLogType = original.LogAnalyticsLogType

const (
	ContainerInsights     LogAnalyticsLogType = original.ContainerInsights
	ContainerInstanceLogs LogAnalyticsLogType = original.ContainerInstanceLogs
)

type OperatingSystemTypes = original.OperatingSystemTypes

const (
	Linux   OperatingSystemTypes = original.Linux
	Windows OperatingSystemTypes = original.Windows
)

type OperationsOrigin = original.OperationsOrigin

const (
	System OperationsOrigin = original.System
	User   OperationsOrigin = original.User
)

type ResourceIdentityType = original.ResourceIdentityType

const (
	None                       ResourceIdentityType = original.None
	SystemAssigned             ResourceIdentityType = original.SystemAssigned
	SystemAssignedUserAssigned ResourceIdentityType = original.SystemAssignedUserAssigned
	UserAssigned               ResourceIdentityType = original.UserAssigned
)

type Scheme = original.Scheme

const (
	HTTP  Scheme = original.HTTP
	HTTPS Scheme = original.HTTPS
)

type AzureFileVolume = original.AzureFileVolume
type BaseClient = original.BaseClient
type Container = original.Container
type ContainerClient = original.ContainerClient
type ContainerExec = original.ContainerExec
type ContainerExecRequest = original.ContainerExecRequest
type ContainerExecRequestTerminalSize = original.ContainerExecRequestTerminalSize
type ContainerExecResponse = original.ContainerExecResponse
type ContainerGroup = original.ContainerGroup
type ContainerGroupDiagnostics = original.ContainerGroupDiagnostics
type ContainerGroupIdentity = original.ContainerGroupIdentity
type ContainerGroupIdentityUserAssignedIdentitiesValue = original.ContainerGroupIdentityUserAssignedIdentitiesValue
type ContainerGroupListResult = original.ContainerGroupListResult
type ContainerGroupListResultIterator = original.ContainerGroupListResultIterator
type ContainerGroupListResultPage = original.ContainerGroupListResultPage
type ContainerGroupNetworkProfile = original.ContainerGroupNetworkProfile
type ContainerGroupProperties = original.ContainerGroupProperties
type ContainerGroupPropertiesInstanceView = original.ContainerGroupPropertiesInstanceView
type ContainerGroupUsageClient = original.ContainerGroupUsageClient
type ContainerGroupsClient = original.ContainerGroupsClient
type ContainerGroupsCreateOrUpdateFuture = original.ContainerGroupsCreateOrUpdateFuture
type ContainerGroupsRestartFuture = original.ContainerGroupsRestartFuture
type ContainerGroupsStartFuture = original.ContainerGroupsStartFuture
type ContainerHTTPGet = original.ContainerHTTPGet
type ContainerPort = original.ContainerPort
type ContainerProbe = original.ContainerProbe
type ContainerProperties = original.ContainerProperties
type ContainerPropertiesInstanceView = original.ContainerPropertiesInstanceView
type ContainerState = original.ContainerState
type DNSConfiguration = original.DNSConfiguration
type EnvironmentVariable = original.EnvironmentVariable
type Event = original.Event
type GitRepoVolume = original.GitRepoVolume
type GpuResource = original.GpuResource
type IPAddress = original.IPAddress
type ImageRegistryCredential = original.ImageRegistryCredential
type LogAnalytics = original.LogAnalytics
type Logs = original.Logs
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationsClient = original.OperationsClient
type Port = original.Port
type Resource = original.Resource
type ResourceLimits = original.ResourceLimits
type ResourceRequests = original.ResourceRequests
type ResourceRequirements = original.ResourceRequirements
type ServiceAssociationLinkClient = original.ServiceAssociationLinkClient
type Usage = original.Usage
type UsageListResult = original.UsageListResult
type UsageName = original.UsageName
type Volume = original.Volume
type VolumeMount = original.VolumeMount

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewContainerClient(subscriptionID string) ContainerClient {
	return original.NewContainerClient(subscriptionID)
}
func NewContainerClientWithBaseURI(baseURI string, subscriptionID string) ContainerClient {
	return original.NewContainerClientWithBaseURI(baseURI, subscriptionID)
}
func NewContainerGroupListResultIterator(page ContainerGroupListResultPage) ContainerGroupListResultIterator {
	return original.NewContainerGroupListResultIterator(page)
}
func NewContainerGroupListResultPage(getNextPage func(context.Context, ContainerGroupListResult) (ContainerGroupListResult, error)) ContainerGroupListResultPage {
	return original.NewContainerGroupListResultPage(getNextPage)
}
func NewContainerGroupUsageClient(subscriptionID string) ContainerGroupUsageClient {
	return original.NewContainerGroupUsageClient(subscriptionID)
}
func NewContainerGroupUsageClientWithBaseURI(baseURI string, subscriptionID string) ContainerGroupUsageClient {
	return original.NewContainerGroupUsageClientWithBaseURI(baseURI, subscriptionID)
}
func NewContainerGroupsClient(subscriptionID string) ContainerGroupsClient {
	return original.NewContainerGroupsClient(subscriptionID)
}
func NewContainerGroupsClientWithBaseURI(baseURI string, subscriptionID string) ContainerGroupsClient {
	return original.NewContainerGroupsClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewServiceAssociationLinkClient(subscriptionID string) ServiceAssociationLinkClient {
	return original.NewServiceAssociationLinkClient(subscriptionID)
}
func NewServiceAssociationLinkClientWithBaseURI(baseURI string, subscriptionID string) ServiceAssociationLinkClient {
	return original.NewServiceAssociationLinkClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleContainerGroupIPAddressTypeValues() []ContainerGroupIPAddressType {
	return original.PossibleContainerGroupIPAddressTypeValues()
}
func PossibleContainerGroupNetworkProtocolValues() []ContainerGroupNetworkProtocol {
	return original.PossibleContainerGroupNetworkProtocolValues()
}
func PossibleContainerGroupRestartPolicyValues() []ContainerGroupRestartPolicy {
	return original.PossibleContainerGroupRestartPolicyValues()
}
func PossibleContainerNetworkProtocolValues() []ContainerNetworkProtocol {
	return original.PossibleContainerNetworkProtocolValues()
}
func PossibleGpuSkuValues() []GpuSku {
	return original.PossibleGpuSkuValues()
}
func PossibleLogAnalyticsLogTypeValues() []LogAnalyticsLogType {
	return original.PossibleLogAnalyticsLogTypeValues()
}
func PossibleOperatingSystemTypesValues() []OperatingSystemTypes {
	return original.PossibleOperatingSystemTypesValues()
}
func PossibleOperationsOriginValues() []OperationsOrigin {
	return original.PossibleOperationsOriginValues()
}
func PossibleResourceIdentityTypeValues() []ResourceIdentityType {
	return original.PossibleResourceIdentityTypeValues()
}
func PossibleSchemeValues() []Scheme {
	return original.PossibleSchemeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

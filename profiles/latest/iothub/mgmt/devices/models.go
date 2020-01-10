// +build go1.9

// Copyright 2020 Microsoft Corporation
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

package devices

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2018-04-01/devices"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type AccessRights = original.AccessRights

const (
	DeviceConnect                                        AccessRights = original.DeviceConnect
	RegistryRead                                         AccessRights = original.RegistryRead
	RegistryReadDeviceConnect                            AccessRights = original.RegistryReadDeviceConnect
	RegistryReadRegistryWrite                            AccessRights = original.RegistryReadRegistryWrite
	RegistryReadRegistryWriteDeviceConnect               AccessRights = original.RegistryReadRegistryWriteDeviceConnect
	RegistryReadRegistryWriteServiceConnect              AccessRights = original.RegistryReadRegistryWriteServiceConnect
	RegistryReadRegistryWriteServiceConnectDeviceConnect AccessRights = original.RegistryReadRegistryWriteServiceConnectDeviceConnect
	RegistryReadServiceConnect                           AccessRights = original.RegistryReadServiceConnect
	RegistryReadServiceConnectDeviceConnect              AccessRights = original.RegistryReadServiceConnectDeviceConnect
	RegistryWrite                                        AccessRights = original.RegistryWrite
	RegistryWriteDeviceConnect                           AccessRights = original.RegistryWriteDeviceConnect
	RegistryWriteServiceConnect                          AccessRights = original.RegistryWriteServiceConnect
	RegistryWriteServiceConnectDeviceConnect             AccessRights = original.RegistryWriteServiceConnectDeviceConnect
	ServiceConnect                                       AccessRights = original.ServiceConnect
	ServiceConnectDeviceConnect                          AccessRights = original.ServiceConnectDeviceConnect
)

type Capabilities = original.Capabilities

const (
	DeviceManagement Capabilities = original.DeviceManagement
	None             Capabilities = original.None
)

type EndpointHealthStatus = original.EndpointHealthStatus

const (
	Dead      EndpointHealthStatus = original.Dead
	Healthy   EndpointHealthStatus = original.Healthy
	Unhealthy EndpointHealthStatus = original.Unhealthy
	Unknown   EndpointHealthStatus = original.Unknown
)

type IPFilterActionType = original.IPFilterActionType

const (
	Accept IPFilterActionType = original.Accept
	Reject IPFilterActionType = original.Reject
)

type IotHubNameUnavailabilityReason = original.IotHubNameUnavailabilityReason

const (
	AlreadyExists IotHubNameUnavailabilityReason = original.AlreadyExists
	Invalid       IotHubNameUnavailabilityReason = original.Invalid
)

type IotHubScaleType = original.IotHubScaleType

const (
	IotHubScaleTypeAutomatic IotHubScaleType = original.IotHubScaleTypeAutomatic
	IotHubScaleTypeManual    IotHubScaleType = original.IotHubScaleTypeManual
	IotHubScaleTypeNone      IotHubScaleType = original.IotHubScaleTypeNone
)

type IotHubSku = original.IotHubSku

const (
	B1 IotHubSku = original.B1
	B2 IotHubSku = original.B2
	B3 IotHubSku = original.B3
	F1 IotHubSku = original.F1
	S1 IotHubSku = original.S1
	S2 IotHubSku = original.S2
	S3 IotHubSku = original.S3
)

type IotHubSkuTier = original.IotHubSkuTier

const (
	Basic    IotHubSkuTier = original.Basic
	Free     IotHubSkuTier = original.Free
	Standard IotHubSkuTier = original.Standard
)

type JobStatus = original.JobStatus

const (
	JobStatusCancelled JobStatus = original.JobStatusCancelled
	JobStatusCompleted JobStatus = original.JobStatusCompleted
	JobStatusEnqueued  JobStatus = original.JobStatusEnqueued
	JobStatusFailed    JobStatus = original.JobStatusFailed
	JobStatusRunning   JobStatus = original.JobStatusRunning
	JobStatusUnknown   JobStatus = original.JobStatusUnknown
)

type JobType = original.JobType

const (
	JobTypeBackup                    JobType = original.JobTypeBackup
	JobTypeExport                    JobType = original.JobTypeExport
	JobTypeFactoryResetDevice        JobType = original.JobTypeFactoryResetDevice
	JobTypeFirmwareUpdate            JobType = original.JobTypeFirmwareUpdate
	JobTypeImport                    JobType = original.JobTypeImport
	JobTypeReadDeviceProperties      JobType = original.JobTypeReadDeviceProperties
	JobTypeRebootDevice              JobType = original.JobTypeRebootDevice
	JobTypeUnknown                   JobType = original.JobTypeUnknown
	JobTypeUpdateDeviceConfiguration JobType = original.JobTypeUpdateDeviceConfiguration
	JobTypeWriteDeviceProperties     JobType = original.JobTypeWriteDeviceProperties
)

type OperationMonitoringLevel = original.OperationMonitoringLevel

const (
	OperationMonitoringLevelError            OperationMonitoringLevel = original.OperationMonitoringLevelError
	OperationMonitoringLevelErrorInformation OperationMonitoringLevel = original.OperationMonitoringLevelErrorInformation
	OperationMonitoringLevelInformation      OperationMonitoringLevel = original.OperationMonitoringLevelInformation
	OperationMonitoringLevelNone             OperationMonitoringLevel = original.OperationMonitoringLevelNone
)

type RouteErrorSeverity = original.RouteErrorSeverity

const (
	Error   RouteErrorSeverity = original.Error
	Warning RouteErrorSeverity = original.Warning
)

type RoutingSource = original.RoutingSource

const (
	RoutingSourceDeviceJobLifecycleEvents RoutingSource = original.RoutingSourceDeviceJobLifecycleEvents
	RoutingSourceDeviceLifecycleEvents    RoutingSource = original.RoutingSourceDeviceLifecycleEvents
	RoutingSourceDeviceMessages           RoutingSource = original.RoutingSourceDeviceMessages
	RoutingSourceInvalid                  RoutingSource = original.RoutingSourceInvalid
	RoutingSourceTwinChangeEvents         RoutingSource = original.RoutingSourceTwinChangeEvents
)

type TestResultStatus = original.TestResultStatus

const (
	False     TestResultStatus = original.False
	True      TestResultStatus = original.True
	Undefined TestResultStatus = original.Undefined
)

type BaseClient = original.BaseClient
type CertificateBodyDescription = original.CertificateBodyDescription
type CertificateDescription = original.CertificateDescription
type CertificateListDescription = original.CertificateListDescription
type CertificateProperties = original.CertificateProperties
type CertificatePropertiesWithNonce = original.CertificatePropertiesWithNonce
type CertificateVerificationDescription = original.CertificateVerificationDescription
type CertificateWithNonceDescription = original.CertificateWithNonceDescription
type CertificatesClient = original.CertificatesClient
type CloudToDeviceProperties = original.CloudToDeviceProperties
type EndpointHealthData = original.EndpointHealthData
type EndpointHealthDataListResult = original.EndpointHealthDataListResult
type EndpointHealthDataListResultIterator = original.EndpointHealthDataListResultIterator
type EndpointHealthDataListResultPage = original.EndpointHealthDataListResultPage
type ErrorDetails = original.ErrorDetails
type EventHubConsumerGroupInfo = original.EventHubConsumerGroupInfo
type EventHubConsumerGroupsListResult = original.EventHubConsumerGroupsListResult
type EventHubConsumerGroupsListResultIterator = original.EventHubConsumerGroupsListResultIterator
type EventHubConsumerGroupsListResultPage = original.EventHubConsumerGroupsListResultPage
type EventHubProperties = original.EventHubProperties
type ExportDevicesRequest = original.ExportDevicesRequest
type FallbackRouteProperties = original.FallbackRouteProperties
type FeedbackProperties = original.FeedbackProperties
type IPFilterRule = original.IPFilterRule
type ImportDevicesRequest = original.ImportDevicesRequest
type IotHubCapacity = original.IotHubCapacity
type IotHubDescription = original.IotHubDescription
type IotHubDescriptionListResult = original.IotHubDescriptionListResult
type IotHubDescriptionListResultIterator = original.IotHubDescriptionListResultIterator
type IotHubDescriptionListResultPage = original.IotHubDescriptionListResultPage
type IotHubNameAvailabilityInfo = original.IotHubNameAvailabilityInfo
type IotHubProperties = original.IotHubProperties
type IotHubQuotaMetricInfo = original.IotHubQuotaMetricInfo
type IotHubQuotaMetricInfoListResult = original.IotHubQuotaMetricInfoListResult
type IotHubQuotaMetricInfoListResultIterator = original.IotHubQuotaMetricInfoListResultIterator
type IotHubQuotaMetricInfoListResultPage = original.IotHubQuotaMetricInfoListResultPage
type IotHubResourceClient = original.IotHubResourceClient
type IotHubResourceCreateOrUpdateFuture = original.IotHubResourceCreateOrUpdateFuture
type IotHubResourceDeleteFuture = original.IotHubResourceDeleteFuture
type IotHubResourceUpdateFuture = original.IotHubResourceUpdateFuture
type IotHubSkuDescription = original.IotHubSkuDescription
type IotHubSkuDescriptionListResult = original.IotHubSkuDescriptionListResult
type IotHubSkuDescriptionListResultIterator = original.IotHubSkuDescriptionListResultIterator
type IotHubSkuDescriptionListResultPage = original.IotHubSkuDescriptionListResultPage
type IotHubSkuInfo = original.IotHubSkuInfo
type JobResponse = original.JobResponse
type JobResponseListResult = original.JobResponseListResult
type JobResponseListResultIterator = original.JobResponseListResultIterator
type JobResponseListResultPage = original.JobResponseListResultPage
type MatchedRoute = original.MatchedRoute
type MessagingEndpointProperties = original.MessagingEndpointProperties
type Name = original.Name
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationInputs = original.OperationInputs
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationsClient = original.OperationsClient
type OperationsMonitoringProperties = original.OperationsMonitoringProperties
type RegistryStatistics = original.RegistryStatistics
type Resource = original.Resource
type ResourceProviderCommonClient = original.ResourceProviderCommonClient
type RouteCompilationError = original.RouteCompilationError
type RouteErrorPosition = original.RouteErrorPosition
type RouteErrorRange = original.RouteErrorRange
type RouteProperties = original.RouteProperties
type RoutingEndpoints = original.RoutingEndpoints
type RoutingEventHubProperties = original.RoutingEventHubProperties
type RoutingMessage = original.RoutingMessage
type RoutingProperties = original.RoutingProperties
type RoutingServiceBusQueueEndpointProperties = original.RoutingServiceBusQueueEndpointProperties
type RoutingServiceBusTopicEndpointProperties = original.RoutingServiceBusTopicEndpointProperties
type RoutingStorageContainerProperties = original.RoutingStorageContainerProperties
type RoutingTwin = original.RoutingTwin
type RoutingTwinProperties = original.RoutingTwinProperties
type SetObject = original.SetObject
type SharedAccessSignatureAuthorizationRule = original.SharedAccessSignatureAuthorizationRule
type SharedAccessSignatureAuthorizationRuleListResult = original.SharedAccessSignatureAuthorizationRuleListResult
type SharedAccessSignatureAuthorizationRuleListResultIterator = original.SharedAccessSignatureAuthorizationRuleListResultIterator
type SharedAccessSignatureAuthorizationRuleListResultPage = original.SharedAccessSignatureAuthorizationRuleListResultPage
type StorageEndpointProperties = original.StorageEndpointProperties
type TagsResource = original.TagsResource
type TestAllRoutesInput = original.TestAllRoutesInput
type TestAllRoutesResult = original.TestAllRoutesResult
type TestRouteInput = original.TestRouteInput
type TestRouteResult = original.TestRouteResult
type TestRouteResultDetails = original.TestRouteResultDetails
type UserSubscriptionQuota = original.UserSubscriptionQuota
type UserSubscriptionQuotaListResult = original.UserSubscriptionQuotaListResult

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewCertificatesClient(subscriptionID string) CertificatesClient {
	return original.NewCertificatesClient(subscriptionID)
}
func NewCertificatesClientWithBaseURI(baseURI string, subscriptionID string) CertificatesClient {
	return original.NewCertificatesClientWithBaseURI(baseURI, subscriptionID)
}
func NewEndpointHealthDataListResultIterator(page EndpointHealthDataListResultPage) EndpointHealthDataListResultIterator {
	return original.NewEndpointHealthDataListResultIterator(page)
}
func NewEndpointHealthDataListResultPage(getNextPage func(context.Context, EndpointHealthDataListResult) (EndpointHealthDataListResult, error)) EndpointHealthDataListResultPage {
	return original.NewEndpointHealthDataListResultPage(getNextPage)
}
func NewEventHubConsumerGroupsListResultIterator(page EventHubConsumerGroupsListResultPage) EventHubConsumerGroupsListResultIterator {
	return original.NewEventHubConsumerGroupsListResultIterator(page)
}
func NewEventHubConsumerGroupsListResultPage(getNextPage func(context.Context, EventHubConsumerGroupsListResult) (EventHubConsumerGroupsListResult, error)) EventHubConsumerGroupsListResultPage {
	return original.NewEventHubConsumerGroupsListResultPage(getNextPage)
}
func NewIotHubDescriptionListResultIterator(page IotHubDescriptionListResultPage) IotHubDescriptionListResultIterator {
	return original.NewIotHubDescriptionListResultIterator(page)
}
func NewIotHubDescriptionListResultPage(getNextPage func(context.Context, IotHubDescriptionListResult) (IotHubDescriptionListResult, error)) IotHubDescriptionListResultPage {
	return original.NewIotHubDescriptionListResultPage(getNextPage)
}
func NewIotHubQuotaMetricInfoListResultIterator(page IotHubQuotaMetricInfoListResultPage) IotHubQuotaMetricInfoListResultIterator {
	return original.NewIotHubQuotaMetricInfoListResultIterator(page)
}
func NewIotHubQuotaMetricInfoListResultPage(getNextPage func(context.Context, IotHubQuotaMetricInfoListResult) (IotHubQuotaMetricInfoListResult, error)) IotHubQuotaMetricInfoListResultPage {
	return original.NewIotHubQuotaMetricInfoListResultPage(getNextPage)
}
func NewIotHubResourceClient(subscriptionID string) IotHubResourceClient {
	return original.NewIotHubResourceClient(subscriptionID)
}
func NewIotHubResourceClientWithBaseURI(baseURI string, subscriptionID string) IotHubResourceClient {
	return original.NewIotHubResourceClientWithBaseURI(baseURI, subscriptionID)
}
func NewIotHubSkuDescriptionListResultIterator(page IotHubSkuDescriptionListResultPage) IotHubSkuDescriptionListResultIterator {
	return original.NewIotHubSkuDescriptionListResultIterator(page)
}
func NewIotHubSkuDescriptionListResultPage(getNextPage func(context.Context, IotHubSkuDescriptionListResult) (IotHubSkuDescriptionListResult, error)) IotHubSkuDescriptionListResultPage {
	return original.NewIotHubSkuDescriptionListResultPage(getNextPage)
}
func NewJobResponseListResultIterator(page JobResponseListResultPage) JobResponseListResultIterator {
	return original.NewJobResponseListResultIterator(page)
}
func NewJobResponseListResultPage(getNextPage func(context.Context, JobResponseListResult) (JobResponseListResult, error)) JobResponseListResultPage {
	return original.NewJobResponseListResultPage(getNextPage)
}
func NewOperationListResultIterator(page OperationListResultPage) OperationListResultIterator {
	return original.NewOperationListResultIterator(page)
}
func NewOperationListResultPage(getNextPage func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage {
	return original.NewOperationListResultPage(getNextPage)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewResourceProviderCommonClient(subscriptionID string) ResourceProviderCommonClient {
	return original.NewResourceProviderCommonClient(subscriptionID)
}
func NewResourceProviderCommonClientWithBaseURI(baseURI string, subscriptionID string) ResourceProviderCommonClient {
	return original.NewResourceProviderCommonClientWithBaseURI(baseURI, subscriptionID)
}
func NewSharedAccessSignatureAuthorizationRuleListResultIterator(page SharedAccessSignatureAuthorizationRuleListResultPage) SharedAccessSignatureAuthorizationRuleListResultIterator {
	return original.NewSharedAccessSignatureAuthorizationRuleListResultIterator(page)
}
func NewSharedAccessSignatureAuthorizationRuleListResultPage(getNextPage func(context.Context, SharedAccessSignatureAuthorizationRuleListResult) (SharedAccessSignatureAuthorizationRuleListResult, error)) SharedAccessSignatureAuthorizationRuleListResultPage {
	return original.NewSharedAccessSignatureAuthorizationRuleListResultPage(getNextPage)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleAccessRightsValues() []AccessRights {
	return original.PossibleAccessRightsValues()
}
func PossibleCapabilitiesValues() []Capabilities {
	return original.PossibleCapabilitiesValues()
}
func PossibleEndpointHealthStatusValues() []EndpointHealthStatus {
	return original.PossibleEndpointHealthStatusValues()
}
func PossibleIPFilterActionTypeValues() []IPFilterActionType {
	return original.PossibleIPFilterActionTypeValues()
}
func PossibleIotHubNameUnavailabilityReasonValues() []IotHubNameUnavailabilityReason {
	return original.PossibleIotHubNameUnavailabilityReasonValues()
}
func PossibleIotHubScaleTypeValues() []IotHubScaleType {
	return original.PossibleIotHubScaleTypeValues()
}
func PossibleIotHubSkuTierValues() []IotHubSkuTier {
	return original.PossibleIotHubSkuTierValues()
}
func PossibleIotHubSkuValues() []IotHubSku {
	return original.PossibleIotHubSkuValues()
}
func PossibleJobStatusValues() []JobStatus {
	return original.PossibleJobStatusValues()
}
func PossibleJobTypeValues() []JobType {
	return original.PossibleJobTypeValues()
}
func PossibleOperationMonitoringLevelValues() []OperationMonitoringLevel {
	return original.PossibleOperationMonitoringLevelValues()
}
func PossibleRouteErrorSeverityValues() []RouteErrorSeverity {
	return original.PossibleRouteErrorSeverityValues()
}
func PossibleRoutingSourceValues() []RoutingSource {
	return original.PossibleRoutingSourceValues()
}
func PossibleTestResultStatusValues() []TestResultStatus {
	return original.PossibleTestResultStatusValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}

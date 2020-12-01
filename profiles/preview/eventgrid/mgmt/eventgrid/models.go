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

package eventgrid

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2020-06-01/eventgrid"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type DomainProvisioningState = original.DomainProvisioningState

const (
	Canceled  DomainProvisioningState = original.Canceled
	Creating  DomainProvisioningState = original.Creating
	Deleting  DomainProvisioningState = original.Deleting
	Failed    DomainProvisioningState = original.Failed
	Succeeded DomainProvisioningState = original.Succeeded
	Updating  DomainProvisioningState = original.Updating
)

type DomainTopicProvisioningState = original.DomainTopicProvisioningState

const (
	DomainTopicProvisioningStateCanceled  DomainTopicProvisioningState = original.DomainTopicProvisioningStateCanceled
	DomainTopicProvisioningStateCreating  DomainTopicProvisioningState = original.DomainTopicProvisioningStateCreating
	DomainTopicProvisioningStateDeleting  DomainTopicProvisioningState = original.DomainTopicProvisioningStateDeleting
	DomainTopicProvisioningStateFailed    DomainTopicProvisioningState = original.DomainTopicProvisioningStateFailed
	DomainTopicProvisioningStateSucceeded DomainTopicProvisioningState = original.DomainTopicProvisioningStateSucceeded
	DomainTopicProvisioningStateUpdating  DomainTopicProvisioningState = original.DomainTopicProvisioningStateUpdating
)

type EndpointType = original.EndpointType

const (
	EndpointTypeAzureFunction                EndpointType = original.EndpointTypeAzureFunction
	EndpointTypeEventHub                     EndpointType = original.EndpointTypeEventHub
	EndpointTypeEventSubscriptionDestination EndpointType = original.EndpointTypeEventSubscriptionDestination
	EndpointTypeHybridConnection             EndpointType = original.EndpointTypeHybridConnection
	EndpointTypeServiceBusQueue              EndpointType = original.EndpointTypeServiceBusQueue
	EndpointTypeServiceBusTopic              EndpointType = original.EndpointTypeServiceBusTopic
	EndpointTypeStorageQueue                 EndpointType = original.EndpointTypeStorageQueue
	EndpointTypeWebHook                      EndpointType = original.EndpointTypeWebHook
)

type EndpointTypeBasicDeadLetterDestination = original.EndpointTypeBasicDeadLetterDestination

const (
	EndpointTypeDeadLetterDestination EndpointTypeBasicDeadLetterDestination = original.EndpointTypeDeadLetterDestination
	EndpointTypeStorageBlob           EndpointTypeBasicDeadLetterDestination = original.EndpointTypeStorageBlob
)

type EventDeliverySchema = original.EventDeliverySchema

const (
	CloudEventSchemaV10 EventDeliverySchema = original.CloudEventSchemaV10
	CustomInputSchema   EventDeliverySchema = original.CustomInputSchema
	EventGridSchema     EventDeliverySchema = original.EventGridSchema
)

type EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningState

const (
	EventSubscriptionProvisioningStateAwaitingManualAction EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningStateAwaitingManualAction
	EventSubscriptionProvisioningStateCanceled             EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningStateCanceled
	EventSubscriptionProvisioningStateCreating             EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningStateCreating
	EventSubscriptionProvisioningStateDeleting             EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningStateDeleting
	EventSubscriptionProvisioningStateFailed               EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningStateFailed
	EventSubscriptionProvisioningStateSucceeded            EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningStateSucceeded
	EventSubscriptionProvisioningStateUpdating             EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningStateUpdating
)

type IPActionType = original.IPActionType

const (
	Allow IPActionType = original.Allow
)

type InputSchema = original.InputSchema

const (
	InputSchemaCloudEventSchemaV10 InputSchema = original.InputSchemaCloudEventSchemaV10
	InputSchemaCustomEventSchema   InputSchema = original.InputSchemaCustomEventSchema
	InputSchemaEventGridSchema     InputSchema = original.InputSchemaEventGridSchema
)

type InputSchemaMappingType = original.InputSchemaMappingType

const (
	InputSchemaMappingTypeInputSchemaMapping InputSchemaMappingType = original.InputSchemaMappingTypeInputSchemaMapping
	InputSchemaMappingTypeJSON               InputSchemaMappingType = original.InputSchemaMappingTypeJSON
)

type OperatorType = original.OperatorType

const (
	OperatorTypeAdvancedFilter            OperatorType = original.OperatorTypeAdvancedFilter
	OperatorTypeBoolEquals                OperatorType = original.OperatorTypeBoolEquals
	OperatorTypeNumberGreaterThan         OperatorType = original.OperatorTypeNumberGreaterThan
	OperatorTypeNumberGreaterThanOrEquals OperatorType = original.OperatorTypeNumberGreaterThanOrEquals
	OperatorTypeNumberIn                  OperatorType = original.OperatorTypeNumberIn
	OperatorTypeNumberLessThan            OperatorType = original.OperatorTypeNumberLessThan
	OperatorTypeNumberLessThanOrEquals    OperatorType = original.OperatorTypeNumberLessThanOrEquals
	OperatorTypeNumberNotIn               OperatorType = original.OperatorTypeNumberNotIn
	OperatorTypeStringBeginsWith          OperatorType = original.OperatorTypeStringBeginsWith
	OperatorTypeStringContains            OperatorType = original.OperatorTypeStringContains
	OperatorTypeStringEndsWith            OperatorType = original.OperatorTypeStringEndsWith
	OperatorTypeStringIn                  OperatorType = original.OperatorTypeStringIn
	OperatorTypeStringNotIn               OperatorType = original.OperatorTypeStringNotIn
)

type PersistedConnectionStatus = original.PersistedConnectionStatus

const (
	Approved     PersistedConnectionStatus = original.Approved
	Disconnected PersistedConnectionStatus = original.Disconnected
	Pending      PersistedConnectionStatus = original.Pending
	Rejected     PersistedConnectionStatus = original.Rejected
)

type PublicNetworkAccess = original.PublicNetworkAccess

const (
	Disabled PublicNetworkAccess = original.Disabled
	Enabled  PublicNetworkAccess = original.Enabled
)

type ResourceProvisioningState = original.ResourceProvisioningState

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = original.ResourceProvisioningStateCanceled
	ResourceProvisioningStateCreating  ResourceProvisioningState = original.ResourceProvisioningStateCreating
	ResourceProvisioningStateDeleting  ResourceProvisioningState = original.ResourceProvisioningStateDeleting
	ResourceProvisioningStateFailed    ResourceProvisioningState = original.ResourceProvisioningStateFailed
	ResourceProvisioningStateSucceeded ResourceProvisioningState = original.ResourceProvisioningStateSucceeded
	ResourceProvisioningStateUpdating  ResourceProvisioningState = original.ResourceProvisioningStateUpdating
)

type ResourceRegionType = original.ResourceRegionType

const (
	GlobalResource   ResourceRegionType = original.GlobalResource
	RegionalResource ResourceRegionType = original.RegionalResource
)

type TopicProvisioningState = original.TopicProvisioningState

const (
	TopicProvisioningStateCanceled  TopicProvisioningState = original.TopicProvisioningStateCanceled
	TopicProvisioningStateCreating  TopicProvisioningState = original.TopicProvisioningStateCreating
	TopicProvisioningStateDeleting  TopicProvisioningState = original.TopicProvisioningStateDeleting
	TopicProvisioningStateFailed    TopicProvisioningState = original.TopicProvisioningStateFailed
	TopicProvisioningStateSucceeded TopicProvisioningState = original.TopicProvisioningStateSucceeded
	TopicProvisioningStateUpdating  TopicProvisioningState = original.TopicProvisioningStateUpdating
)

type TopicTypeProvisioningState = original.TopicTypeProvisioningState

const (
	TopicTypeProvisioningStateCanceled  TopicTypeProvisioningState = original.TopicTypeProvisioningStateCanceled
	TopicTypeProvisioningStateCreating  TopicTypeProvisioningState = original.TopicTypeProvisioningStateCreating
	TopicTypeProvisioningStateDeleting  TopicTypeProvisioningState = original.TopicTypeProvisioningStateDeleting
	TopicTypeProvisioningStateFailed    TopicTypeProvisioningState = original.TopicTypeProvisioningStateFailed
	TopicTypeProvisioningStateSucceeded TopicTypeProvisioningState = original.TopicTypeProvisioningStateSucceeded
	TopicTypeProvisioningStateUpdating  TopicTypeProvisioningState = original.TopicTypeProvisioningStateUpdating
)

type AdvancedFilter = original.AdvancedFilter
type AzureFunctionEventSubscriptionDestination = original.AzureFunctionEventSubscriptionDestination
type AzureFunctionEventSubscriptionDestinationProperties = original.AzureFunctionEventSubscriptionDestinationProperties
type BaseClient = original.BaseClient
type BasicAdvancedFilter = original.BasicAdvancedFilter
type BasicDeadLetterDestination = original.BasicDeadLetterDestination
type BasicEventSubscriptionDestination = original.BasicEventSubscriptionDestination
type BasicInputSchemaMapping = original.BasicInputSchemaMapping
type BoolEqualsAdvancedFilter = original.BoolEqualsAdvancedFilter
type ConnectionState = original.ConnectionState
type DeadLetterDestination = original.DeadLetterDestination
type Domain = original.Domain
type DomainProperties = original.DomainProperties
type DomainRegenerateKeyRequest = original.DomainRegenerateKeyRequest
type DomainSharedAccessKeys = original.DomainSharedAccessKeys
type DomainTopic = original.DomainTopic
type DomainTopicProperties = original.DomainTopicProperties
type DomainTopicsClient = original.DomainTopicsClient
type DomainTopicsCreateOrUpdateFuture = original.DomainTopicsCreateOrUpdateFuture
type DomainTopicsDeleteFuture = original.DomainTopicsDeleteFuture
type DomainTopicsListResult = original.DomainTopicsListResult
type DomainTopicsListResultIterator = original.DomainTopicsListResultIterator
type DomainTopicsListResultPage = original.DomainTopicsListResultPage
type DomainUpdateParameterProperties = original.DomainUpdateParameterProperties
type DomainUpdateParameters = original.DomainUpdateParameters
type DomainsClient = original.DomainsClient
type DomainsCreateOrUpdateFuture = original.DomainsCreateOrUpdateFuture
type DomainsDeleteFuture = original.DomainsDeleteFuture
type DomainsListResult = original.DomainsListResult
type DomainsListResultIterator = original.DomainsListResultIterator
type DomainsListResultPage = original.DomainsListResultPage
type DomainsUpdateFuture = original.DomainsUpdateFuture
type EventHubEventSubscriptionDestination = original.EventHubEventSubscriptionDestination
type EventHubEventSubscriptionDestinationProperties = original.EventHubEventSubscriptionDestinationProperties
type EventSubscription = original.EventSubscription
type EventSubscriptionDestination = original.EventSubscriptionDestination
type EventSubscriptionFilter = original.EventSubscriptionFilter
type EventSubscriptionFullURL = original.EventSubscriptionFullURL
type EventSubscriptionProperties = original.EventSubscriptionProperties
type EventSubscriptionUpdateParameters = original.EventSubscriptionUpdateParameters
type EventSubscriptionsClient = original.EventSubscriptionsClient
type EventSubscriptionsCreateOrUpdateFuture = original.EventSubscriptionsCreateOrUpdateFuture
type EventSubscriptionsDeleteFuture = original.EventSubscriptionsDeleteFuture
type EventSubscriptionsListResult = original.EventSubscriptionsListResult
type EventSubscriptionsListResultIterator = original.EventSubscriptionsListResultIterator
type EventSubscriptionsListResultPage = original.EventSubscriptionsListResultPage
type EventSubscriptionsUpdateFuture = original.EventSubscriptionsUpdateFuture
type EventType = original.EventType
type EventTypeProperties = original.EventTypeProperties
type EventTypesListResult = original.EventTypesListResult
type HybridConnectionEventSubscriptionDestination = original.HybridConnectionEventSubscriptionDestination
type HybridConnectionEventSubscriptionDestinationProperties = original.HybridConnectionEventSubscriptionDestinationProperties
type InboundIPRule = original.InboundIPRule
type InputSchemaMapping = original.InputSchemaMapping
type JSONField = original.JSONField
type JSONFieldWithDefault = original.JSONFieldWithDefault
type JSONInputSchemaMapping = original.JSONInputSchemaMapping
type JSONInputSchemaMappingProperties = original.JSONInputSchemaMappingProperties
type NumberGreaterThanAdvancedFilter = original.NumberGreaterThanAdvancedFilter
type NumberGreaterThanOrEqualsAdvancedFilter = original.NumberGreaterThanOrEqualsAdvancedFilter
type NumberInAdvancedFilter = original.NumberInAdvancedFilter
type NumberLessThanAdvancedFilter = original.NumberLessThanAdvancedFilter
type NumberLessThanOrEqualsAdvancedFilter = original.NumberLessThanOrEqualsAdvancedFilter
type NumberNotInAdvancedFilter = original.NumberNotInAdvancedFilter
type Operation = original.Operation
type OperationInfo = original.OperationInfo
type OperationsClient = original.OperationsClient
type OperationsListResult = original.OperationsListResult
type PrivateEndpoint = original.PrivateEndpoint
type PrivateEndpointConnection = original.PrivateEndpointConnection
type PrivateEndpointConnectionListResult = original.PrivateEndpointConnectionListResult
type PrivateEndpointConnectionListResultIterator = original.PrivateEndpointConnectionListResultIterator
type PrivateEndpointConnectionListResultPage = original.PrivateEndpointConnectionListResultPage
type PrivateEndpointConnectionProperties = original.PrivateEndpointConnectionProperties
type PrivateEndpointConnectionsClient = original.PrivateEndpointConnectionsClient
type PrivateEndpointConnectionsDeleteFuture = original.PrivateEndpointConnectionsDeleteFuture
type PrivateEndpointConnectionsUpdateFuture = original.PrivateEndpointConnectionsUpdateFuture
type PrivateLinkResource = original.PrivateLinkResource
type PrivateLinkResourceProperties = original.PrivateLinkResourceProperties
type PrivateLinkResourcesClient = original.PrivateLinkResourcesClient
type PrivateLinkResourcesListResult = original.PrivateLinkResourcesListResult
type PrivateLinkResourcesListResultIterator = original.PrivateLinkResourcesListResultIterator
type PrivateLinkResourcesListResultPage = original.PrivateLinkResourcesListResultPage
type Resource = original.Resource
type RetryPolicy = original.RetryPolicy
type ServiceBusQueueEventSubscriptionDestination = original.ServiceBusQueueEventSubscriptionDestination
type ServiceBusQueueEventSubscriptionDestinationProperties = original.ServiceBusQueueEventSubscriptionDestinationProperties
type ServiceBusTopicEventSubscriptionDestination = original.ServiceBusTopicEventSubscriptionDestination
type ServiceBusTopicEventSubscriptionDestinationProperties = original.ServiceBusTopicEventSubscriptionDestinationProperties
type StorageBlobDeadLetterDestination = original.StorageBlobDeadLetterDestination
type StorageBlobDeadLetterDestinationProperties = original.StorageBlobDeadLetterDestinationProperties
type StorageQueueEventSubscriptionDestination = original.StorageQueueEventSubscriptionDestination
type StorageQueueEventSubscriptionDestinationProperties = original.StorageQueueEventSubscriptionDestinationProperties
type StringBeginsWithAdvancedFilter = original.StringBeginsWithAdvancedFilter
type StringContainsAdvancedFilter = original.StringContainsAdvancedFilter
type StringEndsWithAdvancedFilter = original.StringEndsWithAdvancedFilter
type StringInAdvancedFilter = original.StringInAdvancedFilter
type StringNotInAdvancedFilter = original.StringNotInAdvancedFilter
type Topic = original.Topic
type TopicProperties = original.TopicProperties
type TopicRegenerateKeyRequest = original.TopicRegenerateKeyRequest
type TopicSharedAccessKeys = original.TopicSharedAccessKeys
type TopicTypeInfo = original.TopicTypeInfo
type TopicTypeProperties = original.TopicTypeProperties
type TopicTypesClient = original.TopicTypesClient
type TopicTypesListResult = original.TopicTypesListResult
type TopicUpdateParameterProperties = original.TopicUpdateParameterProperties
type TopicUpdateParameters = original.TopicUpdateParameters
type TopicsClient = original.TopicsClient
type TopicsCreateOrUpdateFuture = original.TopicsCreateOrUpdateFuture
type TopicsDeleteFuture = original.TopicsDeleteFuture
type TopicsListResult = original.TopicsListResult
type TopicsListResultIterator = original.TopicsListResultIterator
type TopicsListResultPage = original.TopicsListResultPage
type TopicsUpdateFuture = original.TopicsUpdateFuture
type TrackedResource = original.TrackedResource
type WebHookEventSubscriptionDestination = original.WebHookEventSubscriptionDestination
type WebHookEventSubscriptionDestinationProperties = original.WebHookEventSubscriptionDestinationProperties

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewDomainTopicsClient(subscriptionID string) DomainTopicsClient {
	return original.NewDomainTopicsClient(subscriptionID)
}
func NewDomainTopicsClientWithBaseURI(baseURI string, subscriptionID string) DomainTopicsClient {
	return original.NewDomainTopicsClientWithBaseURI(baseURI, subscriptionID)
}
func NewDomainTopicsListResultIterator(page DomainTopicsListResultPage) DomainTopicsListResultIterator {
	return original.NewDomainTopicsListResultIterator(page)
}
func NewDomainTopicsListResultPage(cur DomainTopicsListResult, getNextPage func(context.Context, DomainTopicsListResult) (DomainTopicsListResult, error)) DomainTopicsListResultPage {
	return original.NewDomainTopicsListResultPage(cur, getNextPage)
}
func NewDomainsClient(subscriptionID string) DomainsClient {
	return original.NewDomainsClient(subscriptionID)
}
func NewDomainsClientWithBaseURI(baseURI string, subscriptionID string) DomainsClient {
	return original.NewDomainsClientWithBaseURI(baseURI, subscriptionID)
}
func NewDomainsListResultIterator(page DomainsListResultPage) DomainsListResultIterator {
	return original.NewDomainsListResultIterator(page)
}
func NewDomainsListResultPage(cur DomainsListResult, getNextPage func(context.Context, DomainsListResult) (DomainsListResult, error)) DomainsListResultPage {
	return original.NewDomainsListResultPage(cur, getNextPage)
}
func NewEventSubscriptionsClient(subscriptionID string) EventSubscriptionsClient {
	return original.NewEventSubscriptionsClient(subscriptionID)
}
func NewEventSubscriptionsClientWithBaseURI(baseURI string, subscriptionID string) EventSubscriptionsClient {
	return original.NewEventSubscriptionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewEventSubscriptionsListResultIterator(page EventSubscriptionsListResultPage) EventSubscriptionsListResultIterator {
	return original.NewEventSubscriptionsListResultIterator(page)
}
func NewEventSubscriptionsListResultPage(cur EventSubscriptionsListResult, getNextPage func(context.Context, EventSubscriptionsListResult) (EventSubscriptionsListResult, error)) EventSubscriptionsListResultPage {
	return original.NewEventSubscriptionsListResultPage(cur, getNextPage)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewPrivateEndpointConnectionListResultIterator(page PrivateEndpointConnectionListResultPage) PrivateEndpointConnectionListResultIterator {
	return original.NewPrivateEndpointConnectionListResultIterator(page)
}
func NewPrivateEndpointConnectionListResultPage(cur PrivateEndpointConnectionListResult, getNextPage func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error)) PrivateEndpointConnectionListResultPage {
	return original.NewPrivateEndpointConnectionListResultPage(cur, getNextPage)
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
func NewPrivateLinkResourcesListResultIterator(page PrivateLinkResourcesListResultPage) PrivateLinkResourcesListResultIterator {
	return original.NewPrivateLinkResourcesListResultIterator(page)
}
func NewPrivateLinkResourcesListResultPage(cur PrivateLinkResourcesListResult, getNextPage func(context.Context, PrivateLinkResourcesListResult) (PrivateLinkResourcesListResult, error)) PrivateLinkResourcesListResultPage {
	return original.NewPrivateLinkResourcesListResultPage(cur, getNextPage)
}
func NewTopicTypesClient(subscriptionID string) TopicTypesClient {
	return original.NewTopicTypesClient(subscriptionID)
}
func NewTopicTypesClientWithBaseURI(baseURI string, subscriptionID string) TopicTypesClient {
	return original.NewTopicTypesClientWithBaseURI(baseURI, subscriptionID)
}
func NewTopicsClient(subscriptionID string) TopicsClient {
	return original.NewTopicsClient(subscriptionID)
}
func NewTopicsClientWithBaseURI(baseURI string, subscriptionID string) TopicsClient {
	return original.NewTopicsClientWithBaseURI(baseURI, subscriptionID)
}
func NewTopicsListResultIterator(page TopicsListResultPage) TopicsListResultIterator {
	return original.NewTopicsListResultIterator(page)
}
func NewTopicsListResultPage(cur TopicsListResult, getNextPage func(context.Context, TopicsListResult) (TopicsListResult, error)) TopicsListResultPage {
	return original.NewTopicsListResultPage(cur, getNextPage)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleDomainProvisioningStateValues() []DomainProvisioningState {
	return original.PossibleDomainProvisioningStateValues()
}
func PossibleDomainTopicProvisioningStateValues() []DomainTopicProvisioningState {
	return original.PossibleDomainTopicProvisioningStateValues()
}
func PossibleEndpointTypeBasicDeadLetterDestinationValues() []EndpointTypeBasicDeadLetterDestination {
	return original.PossibleEndpointTypeBasicDeadLetterDestinationValues()
}
func PossibleEndpointTypeValues() []EndpointType {
	return original.PossibleEndpointTypeValues()
}
func PossibleEventDeliverySchemaValues() []EventDeliverySchema {
	return original.PossibleEventDeliverySchemaValues()
}
func PossibleEventSubscriptionProvisioningStateValues() []EventSubscriptionProvisioningState {
	return original.PossibleEventSubscriptionProvisioningStateValues()
}
func PossibleIPActionTypeValues() []IPActionType {
	return original.PossibleIPActionTypeValues()
}
func PossibleInputSchemaMappingTypeValues() []InputSchemaMappingType {
	return original.PossibleInputSchemaMappingTypeValues()
}
func PossibleInputSchemaValues() []InputSchema {
	return original.PossibleInputSchemaValues()
}
func PossibleOperatorTypeValues() []OperatorType {
	return original.PossibleOperatorTypeValues()
}
func PossiblePersistedConnectionStatusValues() []PersistedConnectionStatus {
	return original.PossiblePersistedConnectionStatusValues()
}
func PossiblePublicNetworkAccessValues() []PublicNetworkAccess {
	return original.PossiblePublicNetworkAccessValues()
}
func PossibleResourceProvisioningStateValues() []ResourceProvisioningState {
	return original.PossibleResourceProvisioningStateValues()
}
func PossibleResourceRegionTypeValues() []ResourceRegionType {
	return original.PossibleResourceRegionTypeValues()
}
func PossibleTopicProvisioningStateValues() []TopicProvisioningState {
	return original.PossibleTopicProvisioningStateValues()
}
func PossibleTopicTypeProvisioningStateValues() []TopicTypeProvisioningState {
	return original.PossibleTopicTypeProvisioningStateValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

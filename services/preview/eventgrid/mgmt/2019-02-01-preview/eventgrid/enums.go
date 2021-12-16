package eventgrid

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// DomainProvisioningState enumerates the values for domain provisioning state.
type DomainProvisioningState string

const (
	// Canceled ...
	Canceled DomainProvisioningState = "Canceled"
	// Creating ...
	Creating DomainProvisioningState = "Creating"
	// Deleting ...
	Deleting DomainProvisioningState = "Deleting"
	// Failed ...
	Failed DomainProvisioningState = "Failed"
	// Succeeded ...
	Succeeded DomainProvisioningState = "Succeeded"
	// Updating ...
	Updating DomainProvisioningState = "Updating"
)

// PossibleDomainProvisioningStateValues returns an array of possible values for the DomainProvisioningState const type.
func PossibleDomainProvisioningStateValues() []DomainProvisioningState {
	return []DomainProvisioningState{Canceled, Creating, Deleting, Failed, Succeeded, Updating}
}

// DomainTopicProvisioningState enumerates the values for domain topic provisioning state.
type DomainTopicProvisioningState string

const (
	// DomainTopicProvisioningStateCanceled ...
	DomainTopicProvisioningStateCanceled DomainTopicProvisioningState = "Canceled"
	// DomainTopicProvisioningStateCreating ...
	DomainTopicProvisioningStateCreating DomainTopicProvisioningState = "Creating"
	// DomainTopicProvisioningStateDeleting ...
	DomainTopicProvisioningStateDeleting DomainTopicProvisioningState = "Deleting"
	// DomainTopicProvisioningStateFailed ...
	DomainTopicProvisioningStateFailed DomainTopicProvisioningState = "Failed"
	// DomainTopicProvisioningStateSucceeded ...
	DomainTopicProvisioningStateSucceeded DomainTopicProvisioningState = "Succeeded"
	// DomainTopicProvisioningStateUpdating ...
	DomainTopicProvisioningStateUpdating DomainTopicProvisioningState = "Updating"
)

// PossibleDomainTopicProvisioningStateValues returns an array of possible values for the DomainTopicProvisioningState const type.
func PossibleDomainTopicProvisioningStateValues() []DomainTopicProvisioningState {
	return []DomainTopicProvisioningState{DomainTopicProvisioningStateCanceled, DomainTopicProvisioningStateCreating, DomainTopicProvisioningStateDeleting, DomainTopicProvisioningStateFailed, DomainTopicProvisioningStateSucceeded, DomainTopicProvisioningStateUpdating}
}

// EndpointType enumerates the values for endpoint type.
type EndpointType string

const (
	// EndpointTypeEventHub ...
	EndpointTypeEventHub EndpointType = "EventHub"
	// EndpointTypeEventSubscriptionDestination ...
	EndpointTypeEventSubscriptionDestination EndpointType = "EventSubscriptionDestination"
	// EndpointTypeHybridConnection ...
	EndpointTypeHybridConnection EndpointType = "HybridConnection"
	// EndpointTypeServiceBusQueue ...
	EndpointTypeServiceBusQueue EndpointType = "ServiceBusQueue"
	// EndpointTypeStorageQueue ...
	EndpointTypeStorageQueue EndpointType = "StorageQueue"
	// EndpointTypeWebHook ...
	EndpointTypeWebHook EndpointType = "WebHook"
)

// PossibleEndpointTypeValues returns an array of possible values for the EndpointType const type.
func PossibleEndpointTypeValues() []EndpointType {
	return []EndpointType{EndpointTypeEventHub, EndpointTypeEventSubscriptionDestination, EndpointTypeHybridConnection, EndpointTypeServiceBusQueue, EndpointTypeStorageQueue, EndpointTypeWebHook}
}

// EndpointTypeBasicDeadLetterDestination enumerates the values for endpoint type basic dead letter
// destination.
type EndpointTypeBasicDeadLetterDestination string

const (
	// EndpointTypeDeadLetterDestination ...
	EndpointTypeDeadLetterDestination EndpointTypeBasicDeadLetterDestination = "DeadLetterDestination"
	// EndpointTypeStorageBlob ...
	EndpointTypeStorageBlob EndpointTypeBasicDeadLetterDestination = "StorageBlob"
)

// PossibleEndpointTypeBasicDeadLetterDestinationValues returns an array of possible values for the EndpointTypeBasicDeadLetterDestination const type.
func PossibleEndpointTypeBasicDeadLetterDestinationValues() []EndpointTypeBasicDeadLetterDestination {
	return []EndpointTypeBasicDeadLetterDestination{EndpointTypeDeadLetterDestination, EndpointTypeStorageBlob}
}

// EventDeliverySchema enumerates the values for event delivery schema.
type EventDeliverySchema string

const (
	// CloudEventV01Schema ...
	CloudEventV01Schema EventDeliverySchema = "CloudEventV01Schema"
	// CustomInputSchema ...
	CustomInputSchema EventDeliverySchema = "CustomInputSchema"
	// EventGridSchema ...
	EventGridSchema EventDeliverySchema = "EventGridSchema"
)

// PossibleEventDeliverySchemaValues returns an array of possible values for the EventDeliverySchema const type.
func PossibleEventDeliverySchemaValues() []EventDeliverySchema {
	return []EventDeliverySchema{CloudEventV01Schema, CustomInputSchema, EventGridSchema}
}

// EventSubscriptionProvisioningState enumerates the values for event subscription provisioning state.
type EventSubscriptionProvisioningState string

const (
	// EventSubscriptionProvisioningStateAwaitingManualAction ...
	EventSubscriptionProvisioningStateAwaitingManualAction EventSubscriptionProvisioningState = "AwaitingManualAction"
	// EventSubscriptionProvisioningStateCanceled ...
	EventSubscriptionProvisioningStateCanceled EventSubscriptionProvisioningState = "Canceled"
	// EventSubscriptionProvisioningStateCreating ...
	EventSubscriptionProvisioningStateCreating EventSubscriptionProvisioningState = "Creating"
	// EventSubscriptionProvisioningStateDeleting ...
	EventSubscriptionProvisioningStateDeleting EventSubscriptionProvisioningState = "Deleting"
	// EventSubscriptionProvisioningStateFailed ...
	EventSubscriptionProvisioningStateFailed EventSubscriptionProvisioningState = "Failed"
	// EventSubscriptionProvisioningStateSucceeded ...
	EventSubscriptionProvisioningStateSucceeded EventSubscriptionProvisioningState = "Succeeded"
	// EventSubscriptionProvisioningStateUpdating ...
	EventSubscriptionProvisioningStateUpdating EventSubscriptionProvisioningState = "Updating"
)

// PossibleEventSubscriptionProvisioningStateValues returns an array of possible values for the EventSubscriptionProvisioningState const type.
func PossibleEventSubscriptionProvisioningStateValues() []EventSubscriptionProvisioningState {
	return []EventSubscriptionProvisioningState{EventSubscriptionProvisioningStateAwaitingManualAction, EventSubscriptionProvisioningStateCanceled, EventSubscriptionProvisioningStateCreating, EventSubscriptionProvisioningStateDeleting, EventSubscriptionProvisioningStateFailed, EventSubscriptionProvisioningStateSucceeded, EventSubscriptionProvisioningStateUpdating}
}

// InputSchema enumerates the values for input schema.
type InputSchema string

const (
	// InputSchemaCloudEventV01Schema ...
	InputSchemaCloudEventV01Schema InputSchema = "CloudEventV01Schema"
	// InputSchemaCustomEventSchema ...
	InputSchemaCustomEventSchema InputSchema = "CustomEventSchema"
	// InputSchemaEventGridSchema ...
	InputSchemaEventGridSchema InputSchema = "EventGridSchema"
)

// PossibleInputSchemaValues returns an array of possible values for the InputSchema const type.
func PossibleInputSchemaValues() []InputSchema {
	return []InputSchema{InputSchemaCloudEventV01Schema, InputSchemaCustomEventSchema, InputSchemaEventGridSchema}
}

// InputSchemaMappingType enumerates the values for input schema mapping type.
type InputSchemaMappingType string

const (
	// InputSchemaMappingTypeInputSchemaMapping ...
	InputSchemaMappingTypeInputSchemaMapping InputSchemaMappingType = "InputSchemaMapping"
	// InputSchemaMappingTypeJSON ...
	InputSchemaMappingTypeJSON InputSchemaMappingType = "Json"
)

// PossibleInputSchemaMappingTypeValues returns an array of possible values for the InputSchemaMappingType const type.
func PossibleInputSchemaMappingTypeValues() []InputSchemaMappingType {
	return []InputSchemaMappingType{InputSchemaMappingTypeInputSchemaMapping, InputSchemaMappingTypeJSON}
}

// OperatorType enumerates the values for operator type.
type OperatorType string

const (
	// OperatorTypeAdvancedFilter ...
	OperatorTypeAdvancedFilter OperatorType = "AdvancedFilter"
	// OperatorTypeBoolEquals ...
	OperatorTypeBoolEquals OperatorType = "BoolEquals"
	// OperatorTypeNumberGreaterThan ...
	OperatorTypeNumberGreaterThan OperatorType = "NumberGreaterThan"
	// OperatorTypeNumberGreaterThanOrEquals ...
	OperatorTypeNumberGreaterThanOrEquals OperatorType = "NumberGreaterThanOrEquals"
	// OperatorTypeNumberIn ...
	OperatorTypeNumberIn OperatorType = "NumberIn"
	// OperatorTypeNumberLessThan ...
	OperatorTypeNumberLessThan OperatorType = "NumberLessThan"
	// OperatorTypeNumberLessThanOrEquals ...
	OperatorTypeNumberLessThanOrEquals OperatorType = "NumberLessThanOrEquals"
	// OperatorTypeNumberNotIn ...
	OperatorTypeNumberNotIn OperatorType = "NumberNotIn"
	// OperatorTypeStringBeginsWith ...
	OperatorTypeStringBeginsWith OperatorType = "StringBeginsWith"
	// OperatorTypeStringContains ...
	OperatorTypeStringContains OperatorType = "StringContains"
	// OperatorTypeStringEndsWith ...
	OperatorTypeStringEndsWith OperatorType = "StringEndsWith"
	// OperatorTypeStringIn ...
	OperatorTypeStringIn OperatorType = "StringIn"
	// OperatorTypeStringNotIn ...
	OperatorTypeStringNotIn OperatorType = "StringNotIn"
)

// PossibleOperatorTypeValues returns an array of possible values for the OperatorType const type.
func PossibleOperatorTypeValues() []OperatorType {
	return []OperatorType{OperatorTypeAdvancedFilter, OperatorTypeBoolEquals, OperatorTypeNumberGreaterThan, OperatorTypeNumberGreaterThanOrEquals, OperatorTypeNumberIn, OperatorTypeNumberLessThan, OperatorTypeNumberLessThanOrEquals, OperatorTypeNumberNotIn, OperatorTypeStringBeginsWith, OperatorTypeStringContains, OperatorTypeStringEndsWith, OperatorTypeStringIn, OperatorTypeStringNotIn}
}

// ResourceRegionType enumerates the values for resource region type.
type ResourceRegionType string

const (
	// GlobalResource ...
	GlobalResource ResourceRegionType = "GlobalResource"
	// RegionalResource ...
	RegionalResource ResourceRegionType = "RegionalResource"
)

// PossibleResourceRegionTypeValues returns an array of possible values for the ResourceRegionType const type.
func PossibleResourceRegionTypeValues() []ResourceRegionType {
	return []ResourceRegionType{GlobalResource, RegionalResource}
}

// TopicProvisioningState enumerates the values for topic provisioning state.
type TopicProvisioningState string

const (
	// TopicProvisioningStateCanceled ...
	TopicProvisioningStateCanceled TopicProvisioningState = "Canceled"
	// TopicProvisioningStateCreating ...
	TopicProvisioningStateCreating TopicProvisioningState = "Creating"
	// TopicProvisioningStateDeleting ...
	TopicProvisioningStateDeleting TopicProvisioningState = "Deleting"
	// TopicProvisioningStateFailed ...
	TopicProvisioningStateFailed TopicProvisioningState = "Failed"
	// TopicProvisioningStateSucceeded ...
	TopicProvisioningStateSucceeded TopicProvisioningState = "Succeeded"
	// TopicProvisioningStateUpdating ...
	TopicProvisioningStateUpdating TopicProvisioningState = "Updating"
)

// PossibleTopicProvisioningStateValues returns an array of possible values for the TopicProvisioningState const type.
func PossibleTopicProvisioningStateValues() []TopicProvisioningState {
	return []TopicProvisioningState{TopicProvisioningStateCanceled, TopicProvisioningStateCreating, TopicProvisioningStateDeleting, TopicProvisioningStateFailed, TopicProvisioningStateSucceeded, TopicProvisioningStateUpdating}
}

// TopicTypeProvisioningState enumerates the values for topic type provisioning state.
type TopicTypeProvisioningState string

const (
	// TopicTypeProvisioningStateCanceled ...
	TopicTypeProvisioningStateCanceled TopicTypeProvisioningState = "Canceled"
	// TopicTypeProvisioningStateCreating ...
	TopicTypeProvisioningStateCreating TopicTypeProvisioningState = "Creating"
	// TopicTypeProvisioningStateDeleting ...
	TopicTypeProvisioningStateDeleting TopicTypeProvisioningState = "Deleting"
	// TopicTypeProvisioningStateFailed ...
	TopicTypeProvisioningStateFailed TopicTypeProvisioningState = "Failed"
	// TopicTypeProvisioningStateSucceeded ...
	TopicTypeProvisioningStateSucceeded TopicTypeProvisioningState = "Succeeded"
	// TopicTypeProvisioningStateUpdating ...
	TopicTypeProvisioningStateUpdating TopicTypeProvisioningState = "Updating"
)

// PossibleTopicTypeProvisioningStateValues returns an array of possible values for the TopicTypeProvisioningState const type.
func PossibleTopicTypeProvisioningStateValues() []TopicTypeProvisioningState {
	return []TopicTypeProvisioningState{TopicTypeProvisioningStateCanceled, TopicTypeProvisioningStateCreating, TopicTypeProvisioningStateDeleting, TopicTypeProvisioningStateFailed, TopicTypeProvisioningStateSucceeded, TopicTypeProvisioningStateUpdating}
}

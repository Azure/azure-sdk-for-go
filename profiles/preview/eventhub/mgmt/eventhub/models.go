//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package eventhub

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type AccessRights = original.AccessRights

const (
	Listen        AccessRights = original.Listen
	Manage        AccessRights = original.Manage
	SendEnumValue AccessRights = original.SendEnumValue
)

type DefaultAction = original.DefaultAction

const (
	Allow DefaultAction = original.Allow
	Deny  DefaultAction = original.Deny
)

type EncodingCaptureDescription = original.EncodingCaptureDescription

const (
	Avro        EncodingCaptureDescription = original.Avro
	AvroDeflate EncodingCaptureDescription = original.AvroDeflate
)

type EntityStatus = original.EntityStatus

const (
	Active          EntityStatus = original.Active
	Creating        EntityStatus = original.Creating
	Deleting        EntityStatus = original.Deleting
	Disabled        EntityStatus = original.Disabled
	ReceiveDisabled EntityStatus = original.ReceiveDisabled
	Renaming        EntityStatus = original.Renaming
	Restoring       EntityStatus = original.Restoring
	SendDisabled    EntityStatus = original.SendDisabled
	Unknown         EntityStatus = original.Unknown
)

type KeyType = original.KeyType

const (
	PrimaryKey   KeyType = original.PrimaryKey
	SecondaryKey KeyType = original.SecondaryKey
)

type NetworkRuleIPAction = original.NetworkRuleIPAction

const (
	NetworkRuleIPActionAllow NetworkRuleIPAction = original.NetworkRuleIPActionAllow
)

type ProvisioningStateDR = original.ProvisioningStateDR

const (
	Accepted  ProvisioningStateDR = original.Accepted
	Failed    ProvisioningStateDR = original.Failed
	Succeeded ProvisioningStateDR = original.Succeeded
)

type RoleDisasterRecovery = original.RoleDisasterRecovery

const (
	Primary               RoleDisasterRecovery = original.Primary
	PrimaryNotReplicating RoleDisasterRecovery = original.PrimaryNotReplicating
	Secondary             RoleDisasterRecovery = original.Secondary
)

type SkuName = original.SkuName

const (
	Basic    SkuName = original.Basic
	Standard SkuName = original.Standard
)

type SkuTier = original.SkuTier

const (
	SkuTierBasic    SkuTier = original.SkuTierBasic
	SkuTierStandard SkuTier = original.SkuTierStandard
)

type UnavailableReason = original.UnavailableReason

const (
	InvalidName                           UnavailableReason = original.InvalidName
	NameInLockdown                        UnavailableReason = original.NameInLockdown
	NameInUse                             UnavailableReason = original.NameInUse
	None                                  UnavailableReason = original.None
	SubscriptionIsDisabled                UnavailableReason = original.SubscriptionIsDisabled
	TooManyNamespaceInCurrentSubscription UnavailableReason = original.TooManyNamespaceInCurrentSubscription
)

type AccessKeys = original.AccessKeys
type ArmDisasterRecovery = original.ArmDisasterRecovery
type ArmDisasterRecoveryListResult = original.ArmDisasterRecoveryListResult
type ArmDisasterRecoveryListResultIterator = original.ArmDisasterRecoveryListResultIterator
type ArmDisasterRecoveryListResultPage = original.ArmDisasterRecoveryListResultPage
type ArmDisasterRecoveryProperties = original.ArmDisasterRecoveryProperties
type AuthorizationRule = original.AuthorizationRule
type AuthorizationRuleListResult = original.AuthorizationRuleListResult
type AuthorizationRuleListResultIterator = original.AuthorizationRuleListResultIterator
type AuthorizationRuleListResultPage = original.AuthorizationRuleListResultPage
type AuthorizationRuleProperties = original.AuthorizationRuleProperties
type BaseClient = original.BaseClient
type CaptureDescription = original.CaptureDescription
type CheckNameAvailabilityParameter = original.CheckNameAvailabilityParameter
type CheckNameAvailabilityResult = original.CheckNameAvailabilityResult
type ConsumerGroup = original.ConsumerGroup
type ConsumerGroupListResult = original.ConsumerGroupListResult
type ConsumerGroupListResultIterator = original.ConsumerGroupListResultIterator
type ConsumerGroupListResultPage = original.ConsumerGroupListResultPage
type ConsumerGroupProperties = original.ConsumerGroupProperties
type ConsumerGroupsClient = original.ConsumerGroupsClient
type Destination = original.Destination
type DestinationProperties = original.DestinationProperties
type DisasterRecoveryConfigsClient = original.DisasterRecoveryConfigsClient
type EHNamespace = original.EHNamespace
type EHNamespaceListResult = original.EHNamespaceListResult
type EHNamespaceListResultIterator = original.EHNamespaceListResultIterator
type EHNamespaceListResultPage = original.EHNamespaceListResultPage
type EHNamespaceProperties = original.EHNamespaceProperties
type ErrorResponse = original.ErrorResponse
type EventHubsClient = original.EventHubsClient
type ListResult = original.ListResult
type ListResultIterator = original.ListResultIterator
type ListResultPage = original.ListResultPage
type MessagingPlan = original.MessagingPlan
type MessagingPlanProperties = original.MessagingPlanProperties
type MessagingRegions = original.MessagingRegions
type MessagingRegionsListResult = original.MessagingRegionsListResult
type MessagingRegionsListResultIterator = original.MessagingRegionsListResultIterator
type MessagingRegionsListResultPage = original.MessagingRegionsListResultPage
type MessagingRegionsProperties = original.MessagingRegionsProperties
type Model = original.Model
type NWRuleSetIPRules = original.NWRuleSetIPRules
type NWRuleSetVirtualNetworkRules = original.NWRuleSetVirtualNetworkRules
type NamespacesClient = original.NamespacesClient
type NamespacesCreateOrUpdateFuture = original.NamespacesCreateOrUpdateFuture
type NamespacesDeleteFuture = original.NamespacesDeleteFuture
type NetworkRuleSet = original.NetworkRuleSet
type NetworkRuleSetListResult = original.NetworkRuleSetListResult
type NetworkRuleSetListResultIterator = original.NetworkRuleSetListResultIterator
type NetworkRuleSetListResultPage = original.NetworkRuleSetListResultPage
type NetworkRuleSetProperties = original.NetworkRuleSetProperties
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationsClient = original.OperationsClient
type Properties = original.Properties
type RegenerateAccessKeyParameters = original.RegenerateAccessKeyParameters
type RegionsClient = original.RegionsClient
type Resource = original.Resource
type Sku = original.Sku
type Subnet = original.Subnet
type TrackedResource = original.TrackedResource

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewArmDisasterRecoveryListResultIterator(page ArmDisasterRecoveryListResultPage) ArmDisasterRecoveryListResultIterator {
	return original.NewArmDisasterRecoveryListResultIterator(page)
}
func NewArmDisasterRecoveryListResultPage(cur ArmDisasterRecoveryListResult, getNextPage func(context.Context, ArmDisasterRecoveryListResult) (ArmDisasterRecoveryListResult, error)) ArmDisasterRecoveryListResultPage {
	return original.NewArmDisasterRecoveryListResultPage(cur, getNextPage)
}
func NewAuthorizationRuleListResultIterator(page AuthorizationRuleListResultPage) AuthorizationRuleListResultIterator {
	return original.NewAuthorizationRuleListResultIterator(page)
}
func NewAuthorizationRuleListResultPage(cur AuthorizationRuleListResult, getNextPage func(context.Context, AuthorizationRuleListResult) (AuthorizationRuleListResult, error)) AuthorizationRuleListResultPage {
	return original.NewAuthorizationRuleListResultPage(cur, getNextPage)
}
func NewConsumerGroupListResultIterator(page ConsumerGroupListResultPage) ConsumerGroupListResultIterator {
	return original.NewConsumerGroupListResultIterator(page)
}
func NewConsumerGroupListResultPage(cur ConsumerGroupListResult, getNextPage func(context.Context, ConsumerGroupListResult) (ConsumerGroupListResult, error)) ConsumerGroupListResultPage {
	return original.NewConsumerGroupListResultPage(cur, getNextPage)
}
func NewConsumerGroupsClient(subscriptionID string) ConsumerGroupsClient {
	return original.NewConsumerGroupsClient(subscriptionID)
}
func NewConsumerGroupsClientWithBaseURI(baseURI string, subscriptionID string) ConsumerGroupsClient {
	return original.NewConsumerGroupsClientWithBaseURI(baseURI, subscriptionID)
}
func NewDisasterRecoveryConfigsClient(subscriptionID string) DisasterRecoveryConfigsClient {
	return original.NewDisasterRecoveryConfigsClient(subscriptionID)
}
func NewDisasterRecoveryConfigsClientWithBaseURI(baseURI string, subscriptionID string) DisasterRecoveryConfigsClient {
	return original.NewDisasterRecoveryConfigsClientWithBaseURI(baseURI, subscriptionID)
}
func NewEHNamespaceListResultIterator(page EHNamespaceListResultPage) EHNamespaceListResultIterator {
	return original.NewEHNamespaceListResultIterator(page)
}
func NewEHNamespaceListResultPage(cur EHNamespaceListResult, getNextPage func(context.Context, EHNamespaceListResult) (EHNamespaceListResult, error)) EHNamespaceListResultPage {
	return original.NewEHNamespaceListResultPage(cur, getNextPage)
}
func NewEventHubsClient(subscriptionID string) EventHubsClient {
	return original.NewEventHubsClient(subscriptionID)
}
func NewEventHubsClientWithBaseURI(baseURI string, subscriptionID string) EventHubsClient {
	return original.NewEventHubsClientWithBaseURI(baseURI, subscriptionID)
}
func NewListResultIterator(page ListResultPage) ListResultIterator {
	return original.NewListResultIterator(page)
}
func NewListResultPage(cur ListResult, getNextPage func(context.Context, ListResult) (ListResult, error)) ListResultPage {
	return original.NewListResultPage(cur, getNextPage)
}
func NewMessagingRegionsListResultIterator(page MessagingRegionsListResultPage) MessagingRegionsListResultIterator {
	return original.NewMessagingRegionsListResultIterator(page)
}
func NewMessagingRegionsListResultPage(cur MessagingRegionsListResult, getNextPage func(context.Context, MessagingRegionsListResult) (MessagingRegionsListResult, error)) MessagingRegionsListResultPage {
	return original.NewMessagingRegionsListResultPage(cur, getNextPage)
}
func NewNamespacesClient(subscriptionID string) NamespacesClient {
	return original.NewNamespacesClient(subscriptionID)
}
func NewNamespacesClientWithBaseURI(baseURI string, subscriptionID string) NamespacesClient {
	return original.NewNamespacesClientWithBaseURI(baseURI, subscriptionID)
}
func NewNetworkRuleSetListResultIterator(page NetworkRuleSetListResultPage) NetworkRuleSetListResultIterator {
	return original.NewNetworkRuleSetListResultIterator(page)
}
func NewNetworkRuleSetListResultPage(cur NetworkRuleSetListResult, getNextPage func(context.Context, NetworkRuleSetListResult) (NetworkRuleSetListResult, error)) NetworkRuleSetListResultPage {
	return original.NewNetworkRuleSetListResultPage(cur, getNextPage)
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
func NewRegionsClient(subscriptionID string) RegionsClient {
	return original.NewRegionsClient(subscriptionID)
}
func NewRegionsClientWithBaseURI(baseURI string, subscriptionID string) RegionsClient {
	return original.NewRegionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleAccessRightsValues() []AccessRights {
	return original.PossibleAccessRightsValues()
}
func PossibleDefaultActionValues() []DefaultAction {
	return original.PossibleDefaultActionValues()
}
func PossibleEncodingCaptureDescriptionValues() []EncodingCaptureDescription {
	return original.PossibleEncodingCaptureDescriptionValues()
}
func PossibleEntityStatusValues() []EntityStatus {
	return original.PossibleEntityStatusValues()
}
func PossibleKeyTypeValues() []KeyType {
	return original.PossibleKeyTypeValues()
}
func PossibleNetworkRuleIPActionValues() []NetworkRuleIPAction {
	return original.PossibleNetworkRuleIPActionValues()
}
func PossibleProvisioningStateDRValues() []ProvisioningStateDR {
	return original.PossibleProvisioningStateDRValues()
}
func PossibleRoleDisasterRecoveryValues() []RoleDisasterRecovery {
	return original.PossibleRoleDisasterRecoveryValues()
}
func PossibleSkuNameValues() []SkuName {
	return original.PossibleSkuNameValues()
}
func PossibleSkuTierValues() []SkuTier {
	return original.PossibleSkuTierValues()
}
func PossibleUnavailableReasonValues() []UnavailableReason {
	return original.PossibleUnavailableReasonValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

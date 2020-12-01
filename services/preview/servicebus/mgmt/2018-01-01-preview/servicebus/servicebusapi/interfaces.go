package servicebusapi

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

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2018-01-01-preview/servicebus"
	"github.com/Azure/go-autorest/autorest"
)

// NamespacesClientAPI contains the set of methods on the NamespacesClient type.
type NamespacesClientAPI interface {
	CheckNameAvailabilityMethod(ctx context.Context, parameters servicebus.CheckNameAvailability) (result servicebus.CheckNameAvailabilityResult, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, parameters servicebus.SBNamespace) (result servicebus.NamespacesCreateOrUpdateFuture, err error)
	CreateOrUpdateAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, authorizationRuleName string, parameters servicebus.SBAuthorizationRule) (result servicebus.SBAuthorizationRule, err error)
	CreateOrUpdateIPFilterRule(ctx context.Context, resourceGroupName string, namespaceName string, IPFilterRuleName string, parameters servicebus.IPFilterRule) (result servicebus.IPFilterRule, err error)
	CreateOrUpdateNetworkRuleSet(ctx context.Context, resourceGroupName string, namespaceName string, parameters servicebus.NetworkRuleSet) (result servicebus.NetworkRuleSet, err error)
	CreateOrUpdateVirtualNetworkRule(ctx context.Context, resourceGroupName string, namespaceName string, virtualNetworkRuleName string, parameters servicebus.VirtualNetworkRule) (result servicebus.VirtualNetworkRule, err error)
	Delete(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.NamespacesDeleteFuture, err error)
	DeleteAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, authorizationRuleName string) (result autorest.Response, err error)
	DeleteIPFilterRule(ctx context.Context, resourceGroupName string, namespaceName string, IPFilterRuleName string) (result autorest.Response, err error)
	DeleteVirtualNetworkRule(ctx context.Context, resourceGroupName string, namespaceName string, virtualNetworkRuleName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.SBNamespace, err error)
	GetAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, authorizationRuleName string) (result servicebus.SBAuthorizationRule, err error)
	GetIPFilterRule(ctx context.Context, resourceGroupName string, namespaceName string, IPFilterRuleName string) (result servicebus.IPFilterRule, err error)
	GetNetworkRuleSet(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.NetworkRuleSet, err error)
	GetVirtualNetworkRule(ctx context.Context, resourceGroupName string, namespaceName string, virtualNetworkRuleName string) (result servicebus.VirtualNetworkRule, err error)
	List(ctx context.Context) (result servicebus.SBNamespaceListResultPage, err error)
	ListComplete(ctx context.Context) (result servicebus.SBNamespaceListResultIterator, err error)
	ListAuthorizationRules(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.SBAuthorizationRuleListResultPage, err error)
	ListAuthorizationRulesComplete(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.SBAuthorizationRuleListResultIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result servicebus.SBNamespaceListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result servicebus.SBNamespaceListResultIterator, err error)
	ListIPFilterRules(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.IPFilterRuleListResultPage, err error)
	ListIPFilterRulesComplete(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.IPFilterRuleListResultIterator, err error)
	ListKeys(ctx context.Context, resourceGroupName string, namespaceName string, authorizationRuleName string) (result servicebus.AccessKeys, err error)
	ListVirtualNetworkRules(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.VirtualNetworkRuleListResultPage, err error)
	ListVirtualNetworkRulesComplete(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.VirtualNetworkRuleListResultIterator, err error)
	Migrate(ctx context.Context, resourceGroupName string, namespaceName string, parameters servicebus.SBNamespaceMigrate) (result autorest.Response, err error)
	RegenerateKeys(ctx context.Context, resourceGroupName string, namespaceName string, authorizationRuleName string, parameters servicebus.RegenerateAccessKeyParameters) (result servicebus.AccessKeys, err error)
	Update(ctx context.Context, resourceGroupName string, namespaceName string, parameters servicebus.SBNamespaceUpdateParameters) (result servicebus.SBNamespace, err error)
}

var _ NamespacesClientAPI = (*servicebus.NamespacesClient)(nil)

// PrivateEndpointConnectionsClientAPI contains the set of methods on the PrivateEndpointConnectionsClient type.
type PrivateEndpointConnectionsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, privateEndpointConnectionName string, parameters servicebus.PrivateEndpointConnection) (result servicebus.PrivateEndpointConnection, err error)
	Delete(ctx context.Context, resourceGroupName string, namespaceName string, privateEndpointConnectionName string) (result servicebus.PrivateEndpointConnectionsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, namespaceName string, privateEndpointConnectionName string) (result servicebus.PrivateEndpointConnection, err error)
	List(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.PrivateEndpointConnectionListResultPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.PrivateEndpointConnectionListResultIterator, err error)
}

var _ PrivateEndpointConnectionsClientAPI = (*servicebus.PrivateEndpointConnectionsClient)(nil)

// PrivateLinkResourcesClientAPI contains the set of methods on the PrivateLinkResourcesClient type.
type PrivateLinkResourcesClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.PrivateLinkResourcesListResult, err error)
}

var _ PrivateLinkResourcesClientAPI = (*servicebus.PrivateLinkResourcesClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result servicebus.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result servicebus.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*servicebus.OperationsClient)(nil)

// DisasterRecoveryConfigsClientAPI contains the set of methods on the DisasterRecoveryConfigsClient type.
type DisasterRecoveryConfigsClientAPI interface {
	BreakPairing(ctx context.Context, resourceGroupName string, namespaceName string, alias string) (result autorest.Response, err error)
	CheckNameAvailabilityMethod(ctx context.Context, resourceGroupName string, namespaceName string, parameters servicebus.CheckNameAvailability) (result servicebus.CheckNameAvailabilityResult, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, alias string, parameters servicebus.ArmDisasterRecovery) (result servicebus.ArmDisasterRecovery, err error)
	Delete(ctx context.Context, resourceGroupName string, namespaceName string, alias string) (result autorest.Response, err error)
	FailOver(ctx context.Context, resourceGroupName string, namespaceName string, alias string, parameters *servicebus.FailoverProperties) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, namespaceName string, alias string) (result servicebus.ArmDisasterRecovery, err error)
	GetAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, alias string, authorizationRuleName string) (result servicebus.SBAuthorizationRule, err error)
	List(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.ArmDisasterRecoveryListResultPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.ArmDisasterRecoveryListResultIterator, err error)
	ListAuthorizationRules(ctx context.Context, resourceGroupName string, namespaceName string, alias string) (result servicebus.SBAuthorizationRuleListResultPage, err error)
	ListAuthorizationRulesComplete(ctx context.Context, resourceGroupName string, namespaceName string, alias string) (result servicebus.SBAuthorizationRuleListResultIterator, err error)
	ListKeys(ctx context.Context, resourceGroupName string, namespaceName string, alias string, authorizationRuleName string) (result servicebus.AccessKeys, err error)
}

var _ DisasterRecoveryConfigsClientAPI = (*servicebus.DisasterRecoveryConfigsClient)(nil)

// QueuesClientAPI contains the set of methods on the QueuesClient type.
type QueuesClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, queueName string, parameters servicebus.SBQueue) (result servicebus.SBQueue, err error)
	CreateOrUpdateAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, queueName string, authorizationRuleName string, parameters servicebus.SBAuthorizationRule) (result servicebus.SBAuthorizationRule, err error)
	Delete(ctx context.Context, resourceGroupName string, namespaceName string, queueName string) (result autorest.Response, err error)
	DeleteAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, queueName string, authorizationRuleName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, namespaceName string, queueName string) (result servicebus.SBQueue, err error)
	GetAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, queueName string, authorizationRuleName string) (result servicebus.SBAuthorizationRule, err error)
	ListAuthorizationRules(ctx context.Context, resourceGroupName string, namespaceName string, queueName string) (result servicebus.SBAuthorizationRuleListResultPage, err error)
	ListAuthorizationRulesComplete(ctx context.Context, resourceGroupName string, namespaceName string, queueName string) (result servicebus.SBAuthorizationRuleListResultIterator, err error)
	ListByNamespace(ctx context.Context, resourceGroupName string, namespaceName string, skip *int32, top *int32) (result servicebus.SBQueueListResultPage, err error)
	ListByNamespaceComplete(ctx context.Context, resourceGroupName string, namespaceName string, skip *int32, top *int32) (result servicebus.SBQueueListResultIterator, err error)
	ListKeys(ctx context.Context, resourceGroupName string, namespaceName string, queueName string, authorizationRuleName string) (result servicebus.AccessKeys, err error)
	RegenerateKeys(ctx context.Context, resourceGroupName string, namespaceName string, queueName string, authorizationRuleName string, parameters servicebus.RegenerateAccessKeyParameters) (result servicebus.AccessKeys, err error)
}

var _ QueuesClientAPI = (*servicebus.QueuesClient)(nil)

// TopicsClientAPI contains the set of methods on the TopicsClient type.
type TopicsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, parameters servicebus.SBTopic) (result servicebus.SBTopic, err error)
	CreateOrUpdateAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, authorizationRuleName string, parameters servicebus.SBAuthorizationRule) (result servicebus.SBAuthorizationRule, err error)
	Delete(ctx context.Context, resourceGroupName string, namespaceName string, topicName string) (result autorest.Response, err error)
	DeleteAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, authorizationRuleName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, namespaceName string, topicName string) (result servicebus.SBTopic, err error)
	GetAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, authorizationRuleName string) (result servicebus.SBAuthorizationRule, err error)
	ListAuthorizationRules(ctx context.Context, resourceGroupName string, namespaceName string, topicName string) (result servicebus.SBAuthorizationRuleListResultPage, err error)
	ListAuthorizationRulesComplete(ctx context.Context, resourceGroupName string, namespaceName string, topicName string) (result servicebus.SBAuthorizationRuleListResultIterator, err error)
	ListByNamespace(ctx context.Context, resourceGroupName string, namespaceName string, skip *int32, top *int32) (result servicebus.SBTopicListResultPage, err error)
	ListByNamespaceComplete(ctx context.Context, resourceGroupName string, namespaceName string, skip *int32, top *int32) (result servicebus.SBTopicListResultIterator, err error)
	ListKeys(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, authorizationRuleName string) (result servicebus.AccessKeys, err error)
	RegenerateKeys(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, authorizationRuleName string, parameters servicebus.RegenerateAccessKeyParameters) (result servicebus.AccessKeys, err error)
}

var _ TopicsClientAPI = (*servicebus.TopicsClient)(nil)

// EventHubsClientAPI contains the set of methods on the EventHubsClient type.
type EventHubsClientAPI interface {
	ListByNamespace(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.EventHubListResultPage, err error)
	ListByNamespaceComplete(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.EventHubListResultIterator, err error)
}

var _ EventHubsClientAPI = (*servicebus.EventHubsClient)(nil)

// MigrationConfigsClientAPI contains the set of methods on the MigrationConfigsClient type.
type MigrationConfigsClientAPI interface {
	CompleteMigration(ctx context.Context, resourceGroupName string, namespaceName string) (result autorest.Response, err error)
	CreateAndStartMigration(ctx context.Context, resourceGroupName string, namespaceName string, parameters servicebus.MigrationConfigProperties) (result servicebus.MigrationConfigsCreateAndStartMigrationFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, namespaceName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.MigrationConfigProperties, err error)
	List(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.MigrationConfigListResultPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, namespaceName string) (result servicebus.MigrationConfigListResultIterator, err error)
	Revert(ctx context.Context, resourceGroupName string, namespaceName string) (result autorest.Response, err error)
}

var _ MigrationConfigsClientAPI = (*servicebus.MigrationConfigsClient)(nil)

// PremiumMessagingRegionsClientAPI contains the set of methods on the PremiumMessagingRegionsClient type.
type PremiumMessagingRegionsClientAPI interface {
	List(ctx context.Context) (result servicebus.PremiumMessagingRegionsListResultPage, err error)
	ListComplete(ctx context.Context) (result servicebus.PremiumMessagingRegionsListResultIterator, err error)
}

var _ PremiumMessagingRegionsClientAPI = (*servicebus.PremiumMessagingRegionsClient)(nil)

// RegionsClientAPI contains the set of methods on the RegionsClient type.
type RegionsClientAPI interface {
	ListBySku(ctx context.Context, sku string) (result servicebus.PremiumMessagingRegionsListResultPage, err error)
	ListBySkuComplete(ctx context.Context, sku string) (result servicebus.PremiumMessagingRegionsListResultIterator, err error)
}

var _ RegionsClientAPI = (*servicebus.RegionsClient)(nil)

// SubscriptionsClientAPI contains the set of methods on the SubscriptionsClient type.
type SubscriptionsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, parameters servicebus.SBSubscription) (result servicebus.SBSubscription, err error)
	Delete(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string) (result servicebus.SBSubscription, err error)
	ListByTopic(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, skip *int32, top *int32) (result servicebus.SBSubscriptionListResultPage, err error)
	ListByTopicComplete(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, skip *int32, top *int32) (result servicebus.SBSubscriptionListResultIterator, err error)
}

var _ SubscriptionsClientAPI = (*servicebus.SubscriptionsClient)(nil)

// RulesClientAPI contains the set of methods on the RulesClient type.
type RulesClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, ruleName string, parameters servicebus.Rule) (result servicebus.Rule, err error)
	Delete(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, ruleName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, ruleName string) (result servicebus.Rule, err error)
	ListBySubscriptions(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, skip *int32, top *int32) (result servicebus.RuleListResultPage, err error)
	ListBySubscriptionsComplete(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, skip *int32, top *int32) (result servicebus.RuleListResultIterator, err error)
}

var _ RulesClientAPI = (*servicebus.RulesClient)(nil)

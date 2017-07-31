package servicebus

import (
	 original "github.com/Azure/azure-sdk-for-go/service/servicebus/management/2017-04-01/servicebus"
)

type (
	 PremiumMessagingRegionsClient = original.PremiumMessagingRegionsClient
	 QueuesClient = original.QueuesClient
	 RegionsClient = original.RegionsClient
	 RulesClient = original.RulesClient
	 SubscriptionsClient = original.SubscriptionsClient
	 TopicsClient = original.TopicsClient
	 ManagementClient = original.ManagementClient
	 OperationsClient = original.OperationsClient
	 NamespacesClient = original.NamespacesClient
	 EventHubsClient = original.EventHubsClient
	 AccessRights = original.AccessRights
	 EncodingCaptureDescription = original.EncodingCaptureDescription
	 EntityStatus = original.EntityStatus
	 FilterType = original.FilterType
	 KeyType = original.KeyType
	 SkuName = original.SkuName
	 SkuTier = original.SkuTier
	 UnavailableReason = original.UnavailableReason
	 AccessKeys = original.AccessKeys
	 Action = original.Action
	 AuthorizationRuleProperties = original.AuthorizationRuleProperties
	 CaptureDescription = original.CaptureDescription
	 CheckNameAvailability = original.CheckNameAvailability
	 CheckNameAvailabilityResult = original.CheckNameAvailabilityResult
	 CorrelationFilter = original.CorrelationFilter
	 Destination = original.Destination
	 DestinationProperties = original.DestinationProperties
	 ErrorResponse = original.ErrorResponse
	 Eventhub = original.Eventhub
	 EventhubProperties = original.EventhubProperties
	 EventHubListResult = original.EventHubListResult
	 MessageCountDetails = original.MessageCountDetails
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 PremiumMessagingRegions = original.PremiumMessagingRegions
	 PremiumMessagingRegionsProperties = original.PremiumMessagingRegionsProperties
	 PremiumMessagingRegionsListResult = original.PremiumMessagingRegionsListResult
	 RegenerateAccessKeyParameters = original.RegenerateAccessKeyParameters
	 Resource = original.Resource
	 ResourceNamespacePatch = original.ResourceNamespacePatch
	 Rule = original.Rule
	 RuleListResult = original.RuleListResult
	 Ruleproperties = original.Ruleproperties
	 SBAuthorizationRule = original.SBAuthorizationRule
	 SBAuthorizationRuleProperties = original.SBAuthorizationRuleProperties
	 SBAuthorizationRuleListResult = original.SBAuthorizationRuleListResult
	 SBNamespace = original.SBNamespace
	 SBNamespaceListResult = original.SBNamespaceListResult
	 SBNamespaceProperties = original.SBNamespaceProperties
	 SBNamespaceUpdateParameters = original.SBNamespaceUpdateParameters
	 SBQueue = original.SBQueue
	 SBQueueListResult = original.SBQueueListResult
	 SBQueueProperties = original.SBQueueProperties
	 SBSku = original.SBSku
	 SBSubscription = original.SBSubscription
	 SBSubscriptionListResult = original.SBSubscriptionListResult
	 SBSubscriptionProperties = original.SBSubscriptionProperties
	 SBTopic = original.SBTopic
	 SBTopicListResult = original.SBTopicListResult
	 SBTopicProperties = original.SBTopicProperties
	 SQLFilter = original.SQLFilter
	 SQLRuleAction = original.SQLRuleAction
	 TrackedResource = original.TrackedResource
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Listen = original.Listen
	 Manage = original.Manage
	 Send = original.Send
	 Avro = original.Avro
	 AvroDeflate = original.AvroDeflate
	 Active = original.Active
	 Creating = original.Creating
	 Deleting = original.Deleting
	 Disabled = original.Disabled
	 ReceiveDisabled = original.ReceiveDisabled
	 Renaming = original.Renaming
	 Restoring = original.Restoring
	 SendDisabled = original.SendDisabled
	 Unknown = original.Unknown
	 FilterTypeCorrelationFilter = original.FilterTypeCorrelationFilter
	 FilterTypeSQLFilter = original.FilterTypeSQLFilter
	 PrimaryKey = original.PrimaryKey
	 SecondaryKey = original.SecondaryKey
	 Basic = original.Basic
	 Premium = original.Premium
	 Standard = original.Standard
	 SkuTierBasic = original.SkuTierBasic
	 SkuTierPremium = original.SkuTierPremium
	 SkuTierStandard = original.SkuTierStandard
	 InvalidName = original.InvalidName
	 NameInLockdown = original.NameInLockdown
	 NameInUse = original.NameInUse
	 None = original.None
	 SubscriptionIsDisabled = original.SubscriptionIsDisabled
	 TooManyNamespaceInCurrentSubscription = original.TooManyNamespaceInCurrentSubscription
)

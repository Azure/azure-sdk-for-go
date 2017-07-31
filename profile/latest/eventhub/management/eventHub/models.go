package eventhub

import (
	 original "github.com/Azure/azure-sdk-for-go/service/eventhub/management/2017-04-01/eventHub"
)

type (
	 ManagementClient = original.ManagementClient
	 ConsumerGroupsClient = original.ConsumerGroupsClient
	 EventHubsClient = original.EventHubsClient
	 AccessRights = original.AccessRights
	 EncodingCaptureDescription = original.EncodingCaptureDescription
	 EntityStatus = original.EntityStatus
	 KeyType = original.KeyType
	 SkuName = original.SkuName
	 SkuTier = original.SkuTier
	 UnavailableReason = original.UnavailableReason
	 AccessKeys = original.AccessKeys
	 AuthorizationRule = original.AuthorizationRule
	 AuthorizationRuleProperties = original.AuthorizationRuleProperties
	 AuthorizationRuleListResult = original.AuthorizationRuleListResult
	 CaptureDescription = original.CaptureDescription
	 CheckNameAvailabilityParameter = original.CheckNameAvailabilityParameter
	 CheckNameAvailabilityResult = original.CheckNameAvailabilityResult
	 ConsumerGroup = original.ConsumerGroup
	 ConsumerGroupProperties = original.ConsumerGroupProperties
	 ConsumerGroupListResult = original.ConsumerGroupListResult
	 Destination = original.Destination
	 DestinationProperties = original.DestinationProperties
	 EHNamespace = original.EHNamespace
	 EHNamespaceProperties = original.EHNamespaceProperties
	 EHNamespaceListResult = original.EHNamespaceListResult
	 ErrorResponse = original.ErrorResponse
	 Eventhub = original.Eventhub
	 EventhubProperties = original.EventhubProperties
	 ListResult = original.ListResult
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 RegenerateAccessKeyParameters = original.RegenerateAccessKeyParameters
	 Resource = original.Resource
	 Sku = original.Sku
	 TrackedResource = original.TrackedResource
	 NamespacesClient = original.NamespacesClient
	 OperationsClient = original.OperationsClient
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
	 PrimaryKey = original.PrimaryKey
	 SecondaryKey = original.SecondaryKey
	 Basic = original.Basic
	 Standard = original.Standard
	 SkuTierBasic = original.SkuTierBasic
	 SkuTierStandard = original.SkuTierStandard
	 InvalidName = original.InvalidName
	 NameInLockdown = original.NameInLockdown
	 NameInUse = original.NameInUse
	 None = original.None
	 SubscriptionIsDisabled = original.SubscriptionIsDisabled
	 TooManyNamespaceInCurrentSubscription = original.TooManyNamespaceInCurrentSubscription
)

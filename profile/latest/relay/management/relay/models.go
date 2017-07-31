package relay

import (
	 original "github.com/Azure/azure-sdk-for-go/service/relay/management/2017-04-01/relay"
)

type (
	 NamespacesClient = original.NamespacesClient
	 OperationsClient = original.OperationsClient
	 WCFRelaysClient = original.WCFRelaysClient
	 ManagementClient = original.ManagementClient
	 HybridConnectionsClient = original.HybridConnectionsClient
	 AccessRights = original.AccessRights
	 KeyType = original.KeyType
	 ProvisioningStateEnum = original.ProvisioningStateEnum
	 RelaytypeEnum = original.RelaytypeEnum
	 SkuTier = original.SkuTier
	 UnavailableReason = original.UnavailableReason
	 AccessKeys = original.AccessKeys
	 AuthorizationRule = original.AuthorizationRule
	 AuthorizationRuleProperties = original.AuthorizationRuleProperties
	 AuthorizationRuleListResult = original.AuthorizationRuleListResult
	 CheckNameAvailability = original.CheckNameAvailability
	 CheckNameAvailabilityResult = original.CheckNameAvailabilityResult
	 ErrorResponse = original.ErrorResponse
	 HybridConnection = original.HybridConnection
	 HybridConnectionProperties = original.HybridConnectionProperties
	 HybridConnectionListResult = original.HybridConnectionListResult
	 Namespace = original.Namespace
	 NamespaceListResult = original.NamespaceListResult
	 NamespaceProperties = original.NamespaceProperties
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 RegenerateAccessKeyParameters = original.RegenerateAccessKeyParameters
	 Resource = original.Resource
	 ResourceNamespacePatch = original.ResourceNamespacePatch
	 Sku = original.Sku
	 TrackedResource = original.TrackedResource
	 UpdateParameters = original.UpdateParameters
	 WcfRelay = original.WcfRelay
	 WcfRelayProperties = original.WcfRelayProperties
	 WcfRelaysListResult = original.WcfRelaysListResult
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Listen = original.Listen
	 Manage = original.Manage
	 Send = original.Send
	 PrimaryKey = original.PrimaryKey
	 SecondaryKey = original.SecondaryKey
	 Created = original.Created
	 Deleted = original.Deleted
	 Failed = original.Failed
	 Succeeded = original.Succeeded
	 Unknown = original.Unknown
	 Updating = original.Updating
	 HTTP = original.HTTP
	 NetTCP = original.NetTCP
	 Standard = original.Standard
	 InvalidName = original.InvalidName
	 NameInLockdown = original.NameInLockdown
	 NameInUse = original.NameInUse
	 None = original.None
	 SubscriptionIsDisabled = original.SubscriptionIsDisabled
	 TooManyNamespaceInCurrentSubscription = original.TooManyNamespaceInCurrentSubscription
)

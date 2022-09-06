# Change History

## Breaking Changes

### Removed Constants

1. IdentityType.None
1. IdentityType.SystemAssigned
1. IdentityType.UserAssigned
1. KeyName.Key1
1. KeyName.Key2
1. KeySource.MicrosoftCognitiveServices
1. KeySource.MicrosoftKeyVault
1. NetworkRuleAction.Allow
1. NetworkRuleAction.Deny
1. PrivateEndpointServiceConnectionStatus.Approved
1. PrivateEndpointServiceConnectionStatus.Disconnected
1. PrivateEndpointServiceConnectionStatus.Pending
1. PrivateEndpointServiceConnectionStatus.Rejected
1. ProvisioningState.Creating
1. ProvisioningState.Deleting
1. ProvisioningState.Failed
1. ProvisioningState.Moving
1. ProvisioningState.ResolvingDNS
1. ProvisioningState.Succeeded
1. PublicNetworkAccess.Disabled
1. PublicNetworkAccess.Enabled
1. QuotaUsageStatus.Blocked
1. QuotaUsageStatus.InOverage
1. QuotaUsageStatus.Included
1. QuotaUsageStatus.Unknown
1. ResourceSkuRestrictionsReasonCode.NotAvailableForSubscription
1. ResourceSkuRestrictionsReasonCode.QuotaID
1. ResourceSkuRestrictionsType.Location
1. ResourceSkuRestrictionsType.Zone
1. SkuTier.Enterprise
1. SkuTier.Free
1. SkuTier.Premium
1. SkuTier.Standard
1. UnitType.Bytes
1. UnitType.BytesPerSecond
1. UnitType.Count
1. UnitType.CountPerSecond
1. UnitType.Milliseconds
1. UnitType.Percent
1. UnitType.Seconds

### Struct Changes

#### Removed Struct Fields

1. PrivateLinkServiceConnectionState.ActionRequired

## Additive Changes

### New Constants

1. IdentityType.IdentityTypeNone
1. IdentityType.IdentityTypeSystemAssigned
1. IdentityType.IdentityTypeUserAssigned
1. KeyName.KeyNameKey1
1. KeyName.KeyNameKey2
1. KeySource.KeySourceMicrosoftCognitiveServices
1. KeySource.KeySourceMicrosoftKeyVault
1. NetworkRuleAction.NetworkRuleActionAllow
1. NetworkRuleAction.NetworkRuleActionDeny
1. PrivateEndpointServiceConnectionStatus.PrivateEndpointServiceConnectionStatusApproved
1. PrivateEndpointServiceConnectionStatus.PrivateEndpointServiceConnectionStatusDisconnected
1. PrivateEndpointServiceConnectionStatus.PrivateEndpointServiceConnectionStatusPending
1. PrivateEndpointServiceConnectionStatus.PrivateEndpointServiceConnectionStatusRejected
1. ProvisioningState.ProvisioningStateCreating
1. ProvisioningState.ProvisioningStateDeleting
1. ProvisioningState.ProvisioningStateFailed
1. ProvisioningState.ProvisioningStateMoving
1. ProvisioningState.ProvisioningStateResolvingDNS
1. ProvisioningState.ProvisioningStateSucceeded
1. PublicNetworkAccess.PublicNetworkAccessDisabled
1. PublicNetworkAccess.PublicNetworkAccessEnabled
1. QuotaUsageStatus.QuotaUsageStatusBlocked
1. QuotaUsageStatus.QuotaUsageStatusInOverage
1. QuotaUsageStatus.QuotaUsageStatusIncluded
1. QuotaUsageStatus.QuotaUsageStatusUnknown
1. ResourceSkuRestrictionsReasonCode.ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription
1. ResourceSkuRestrictionsReasonCode.ResourceSkuRestrictionsReasonCodeQuotaID
1. ResourceSkuRestrictionsType.ResourceSkuRestrictionsTypeLocation
1. ResourceSkuRestrictionsType.ResourceSkuRestrictionsTypeZone
1. SkuTier.SkuTierEnterprise
1. SkuTier.SkuTierFree
1. SkuTier.SkuTierPremium
1. SkuTier.SkuTierStandard
1. UnitType.UnitTypeBytes
1. UnitType.UnitTypeBytesPerSecond
1. UnitType.UnitTypeCount
1. UnitType.UnitTypeCountPerSecond
1. UnitType.UnitTypeMilliseconds
1. UnitType.UnitTypePercent
1. UnitType.UnitTypeSeconds

### New Funcs

1. AccountSkuChangeInfo.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. AccountSkuChangeInfo

#### New Struct Fields

1. AccountAPIProperties.QnaAzureSearchEndpointID
1. AccountAPIProperties.QnaAzureSearchEndpointKey
1. AccountProperties.IsMigrated
1. AccountProperties.SkuChangeInfo
1. PrivateEndpointConnection.Etag
1. PrivateEndpointConnection.Location
1. PrivateLinkServiceConnectionState.ActionsRequired

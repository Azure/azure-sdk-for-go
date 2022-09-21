# Change History

## Breaking Changes

### Removed Constants

1. ProvisioningState1.ProvisioningState1BillingFailed
1. ProvisioningState1.ProvisioningState1Cancelled
1. ProvisioningState1.ProvisioningState1ConfirmedBilling
1. ProvisioningState1.ProvisioningState1ConfirmedResourceHold
1. ProvisioningState1.ProvisioningState1Created
1. ProvisioningState1.ProvisioningState1Creating
1. ProvisioningState1.ProvisioningState1Expired
1. ProvisioningState1.ProvisioningState1Failed
1. ProvisioningState1.ProvisioningState1Merged
1. ProvisioningState1.ProvisioningState1PendingBilling
1. ProvisioningState1.ProvisioningState1PendingResourceHold
1. ProvisioningState1.ProvisioningState1Split
1. ProvisioningState1.ProvisioningState1Succeeded

### Removed Type Aliases

1. string.ProvisioningState1

### Removed Funcs

1. PossibleProvisioningState1Values() []ProvisioningState1

### Signature Changes

#### Struct Fields

1. OrderProperties.ProvisioningState changed type from ProvisioningState1 to ProvisioningState

## Additive Changes

### New Constants

1. DisplayProvisioningState.DisplayProvisioningStateCancelled
1. DisplayProvisioningState.DisplayProvisioningStateExpired
1. DisplayProvisioningState.DisplayProvisioningStateExpiring
1. DisplayProvisioningState.DisplayProvisioningStateFailed
1. DisplayProvisioningState.DisplayProvisioningStatePending
1. DisplayProvisioningState.DisplayProvisioningStateProcessing
1. DisplayProvisioningState.DisplayProvisioningStateSucceeded
1. ErrorResponseCode.ErrorResponseCodeRefundLimitExceeded
1. ErrorResponseCode.ErrorResponseCodeSelfServiceRefundNotSupported
1. Location.LocationAustraliaeast
1. Location.LocationAustraliasoutheast
1. Location.LocationBrazilsouth
1. Location.LocationCanadacentral
1. Location.LocationCanadaeast
1. Location.LocationCentralindia
1. Location.LocationCentralus
1. Location.LocationEastasia
1. Location.LocationEastus
1. Location.LocationEastus2
1. Location.LocationJapaneast
1. Location.LocationJapanwest
1. Location.LocationNorthcentralus
1. Location.LocationNortheurope
1. Location.LocationSouthcentralus
1. Location.LocationSoutheastasia
1. Location.LocationSouthindia
1. Location.LocationUksouth
1. Location.LocationUkwest
1. Location.LocationWestcentralus
1. Location.LocationWesteurope
1. Location.LocationWestindia
1. Location.LocationWestus
1. Location.LocationWestus2
1. UserFriendlyAppliedScopeType.UserFriendlyAppliedScopeTypeManagementGroup
1. UserFriendlyAppliedScopeType.UserFriendlyAppliedScopeTypeNone
1. UserFriendlyAppliedScopeType.UserFriendlyAppliedScopeTypeResourceGroup
1. UserFriendlyAppliedScopeType.UserFriendlyAppliedScopeTypeShared
1. UserFriendlyAppliedScopeType.UserFriendlyAppliedScopeTypeSingle
1. UserFriendlyRenewState.UserFriendlyRenewStateNotApplicable
1. UserFriendlyRenewState.UserFriendlyRenewStateNotRenewed
1. UserFriendlyRenewState.UserFriendlyRenewStateOff
1. UserFriendlyRenewState.UserFriendlyRenewStateOn
1. UserFriendlyRenewState.UserFriendlyRenewStateRenewed

### New Type Aliases

1. string.DisplayProvisioningState
1. string.Location
1. string.UserFriendlyAppliedScopeType
1. string.UserFriendlyRenewState

### New Funcs

1. CalculateRefundClient.Post(context.Context, string, CalculateRefundRequest) (CalculateRefundResponse, error)
1. CalculateRefundClient.PostPreparer(context.Context, string, CalculateRefundRequest) (*http.Request, error)
1. CalculateRefundClient.PostResponder(*http.Response) (CalculateRefundResponse, error)
1. CalculateRefundClient.PostSender(*http.Request) (*http.Response, error)
1. Client.Archive(context.Context, string, string) (autorest.Response, error)
1. Client.ArchivePreparer(context.Context, string, string) (*http.Request, error)
1. Client.ArchiveResponder(*http.Response) (autorest.Response, error)
1. Client.ArchiveSender(*http.Request) (*http.Response, error)
1. Client.Unarchive(context.Context, string, string) (autorest.Response, error)
1. Client.UnarchivePreparer(context.Context, string, string) (*http.Request, error)
1. Client.UnarchiveResponder(*http.Response) (autorest.Response, error)
1. Client.UnarchiveSender(*http.Request) (*http.Response, error)
1. NewCalculateRefundClient() CalculateRefundClient
1. NewCalculateRefundClientWithBaseURI(string) CalculateRefundClient
1. NewReturnClient() ReturnClient
1. NewReturnClientWithBaseURI(string) ReturnClient
1. PossibleDisplayProvisioningStateValues() []DisplayProvisioningState
1. PossibleLocationValues() []Location
1. PossibleUserFriendlyAppliedScopeTypeValues() []UserFriendlyAppliedScopeType
1. PossibleUserFriendlyRenewStateValues() []UserFriendlyRenewState
1. ReturnClient.Post(context.Context, string, RefundRequest) (RefundResponse, error)
1. ReturnClient.PostPreparer(context.Context, string, RefundRequest) (*http.Request, error)
1. ReturnClient.PostResponder(*http.Response) (RefundResponse, error)
1. ReturnClient.PostSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. CalculateRefundClient
1. CalculateRefundRequest
1. CalculateRefundRequestProperties
1. CalculateRefundResponse
1. RefundBillingInformation
1. RefundPolicyError
1. RefundPolicyResult
1. RefundPolicyResultProperty
1. RefundRequest
1. RefundRequestProperties
1. RefundResponse
1. RefundResponseProperties
1. ReturnClient

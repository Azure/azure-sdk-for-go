# Unreleased

## Breaking Changes

### Removed Constants

1. TransactionTypeKind.All
1. TransactionTypeKind.Reservation

### Struct Changes

#### Removed Struct Fields

1. EnrollmentPolicies.MarketplacesEnabled

## Additive Changes

### New Constants

1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeBillingAccountInactive
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeDestinationBillingProfileInactive
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeDestinationBillingProfileNotFound
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeDestinationInvoiceSectionInactive
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeDestinationInvoiceSectionNotFound
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeInvalidDestination
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeMarketplaceNotEnabledOnDestination
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeProductInactive
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeProductNotFound
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeProductTypeNotSupported
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeSourceBillingProfilePastDue
1. SubscriptionTransferValidationErrorCode.SubscriptionTransferValidationErrorCodeSourceInvoiceSectionInactive
1. TransactionTypeKind.TransactionTypeKindAll
1. TransactionTypeKind.TransactionTypeKindReservation

### New Funcs

1. *Reservation.UnmarshalJSON([]byte) error
1. *ReservationsListResultIterator.Next() error
1. *ReservationsListResultIterator.NextWithContext(context.Context) error
1. *ReservationsListResultPage.Next() error
1. *ReservationsListResultPage.NextWithContext(context.Context) error
1. NewReservationsClient(string) ReservationsClient
1. NewReservationsClientWithBaseURI(string, string) ReservationsClient
1. NewReservationsListResultIterator(ReservationsListResultPage) ReservationsListResultIterator
1. NewReservationsListResultPage(ReservationsListResult, func(context.Context, ReservationsListResult) (ReservationsListResult, error)) ReservationsListResultPage
1. Reservation.MarshalJSON() ([]byte, error)
1. ReservationProperty.MarshalJSON() ([]byte, error)
1. ReservationPropertyUtilization.MarshalJSON() ([]byte, error)
1. ReservationSkuProperty.MarshalJSON() ([]byte, error)
1. ReservationSummary.MarshalJSON() ([]byte, error)
1. ReservationUtilizationAggregates.MarshalJSON() ([]byte, error)
1. ReservationsClient.ListByBillingAccount(context.Context, string, string, string, string, string) (ReservationsListResultPage, error)
1. ReservationsClient.ListByBillingAccountComplete(context.Context, string, string, string, string, string) (ReservationsListResultIterator, error)
1. ReservationsClient.ListByBillingAccountPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. ReservationsClient.ListByBillingAccountResponder(*http.Response) (ReservationsListResult, error)
1. ReservationsClient.ListByBillingAccountSender(*http.Request) (*http.Response, error)
1. ReservationsClient.ListByBillingProfile(context.Context, string, string, string, string, string, string) (ReservationsListResultPage, error)
1. ReservationsClient.ListByBillingProfileComplete(context.Context, string, string, string, string, string, string) (ReservationsListResultIterator, error)
1. ReservationsClient.ListByBillingProfilePreparer(context.Context, string, string, string, string, string, string) (*http.Request, error)
1. ReservationsClient.ListByBillingProfileResponder(*http.Response) (ReservationsListResult, error)
1. ReservationsClient.ListByBillingProfileSender(*http.Request) (*http.Response, error)
1. ReservationsListResult.IsEmpty() bool
1. ReservationsListResult.MarshalJSON() ([]byte, error)
1. ReservationsListResultIterator.NotDone() bool
1. ReservationsListResultIterator.Response() ReservationsListResult
1. ReservationsListResultIterator.Value() Reservation
1. ReservationsListResultPage.NotDone() bool
1. ReservationsListResultPage.Response() ReservationsListResult
1. ReservationsListResultPage.Values() []Reservation

### Struct Changes

#### New Structs

1. Reservation
1. ReservationProperty
1. ReservationPropertyUtilization
1. ReservationSkuProperty
1. ReservationSummary
1. ReservationUtilizationAggregates
1. ReservationsClient
1. ReservationsListResult
1. ReservationsListResultIterator
1. ReservationsListResultPage

#### New Struct Fields

1. AccountProperties.NotificationEmailAddress
1. AddressDetails.MiddleName
1. CustomerListResult.TotalCount
1. EnrollmentAccountProperties.AccountOwnerEmail
1. EnrollmentPolicies.MarketplaceEnabled
1. InvoiceSectionListResult.TotalCount
1. Operation.IsDataAction
1. OperationDisplay.Description
1. ProfileProperties.Tags
1. SubscriptionsListResult.TotalCount
1. TransactionListResult.TotalCount

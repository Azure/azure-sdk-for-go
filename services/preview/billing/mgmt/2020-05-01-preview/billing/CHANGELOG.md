
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewRoleAssignmentsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewCustomersClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewInstructionListResultPage` signature has been changed from `(func(context.Context, InstructionListResult) (InstructionListResult, error))` to `(InstructionListResult,func(context.Context, InstructionListResult) (InstructionListResult, error))`
- Function `NewAvailableBalancesClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewProductsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewPropertyClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewSubscriptionsListResultPage` signature has been changed from `(func(context.Context, SubscriptionsListResult) (SubscriptionsListResult, error))` to `(SubscriptionsListResult,func(context.Context, SubscriptionsListResult) (SubscriptionsListResult, error))`
- Function `NewInstructionsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewSubscriptionsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewPoliciesClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewAgreementsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewTransactionListResultPage` signature has been changed from `(func(context.Context, TransactionListResult) (TransactionListResult, error))` to `(TransactionListResult,func(context.Context, TransactionListResult) (TransactionListResult, error))`
- Function `NewTransactionsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewOperationsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewOperationsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewRoleDefinitionsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewPropertyClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewEnrollmentAccountListResultPage` signature has been changed from `(func(context.Context, EnrollmentAccountListResult) (EnrollmentAccountListResult, error))` to `(EnrollmentAccountListResult,func(context.Context, EnrollmentAccountListResult) (EnrollmentAccountListResult, error))`
- Function `NewTransactionsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewProfilesClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewInvoiceSectionsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewProductsListResultPage` signature has been changed from `(func(context.Context, ProductsListResult) (ProductsListResult, error))` to `(ProductsListResult,func(context.Context, ProductsListResult) (ProductsListResult, error))`
- Function `NewInvoiceListResultPage` signature has been changed from `(func(context.Context, InvoiceListResult) (InvoiceListResult, error))` to `(InvoiceListResult,func(context.Context, InvoiceListResult) (InvoiceListResult, error))`
- Function `NewAvailableBalancesClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewAddressClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewPeriodsListResultPage` signature has been changed from `(func(context.Context, PeriodsListResult) (PeriodsListResult, error))` to `(PeriodsListResult,func(context.Context, PeriodsListResult) (PeriodsListResult, error))`
- Function `NewInvoiceSectionsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewPermissionsListResultPage` signature has been changed from `(func(context.Context, PermissionsListResult) (PermissionsListResult, error))` to `(PermissionsListResult,func(context.Context, PermissionsListResult) (PermissionsListResult, error))`
- Function `NewRoleDefinitionsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewEnrollmentAccountsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewProfilesClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewCustomerListResultPage` signature has been changed from `(func(context.Context, CustomerListResult) (CustomerListResult, error))` to `(CustomerListResult,func(context.Context, CustomerListResult) (CustomerListResult, error))`
- Function `NewEnrollmentAccountsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewPeriodsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewCustomersClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewAddressClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewAccountsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewInstructionsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewAgreementListResultPage` signature has been changed from `(func(context.Context, AgreementListResult) (AgreementListResult, error))` to `(AgreementListResult,func(context.Context, AgreementListResult) (AgreementListResult, error))`
- Function `NewInvoiceSectionListResultPage` signature has been changed from `(func(context.Context, InvoiceSectionListResult) (InvoiceSectionListResult, error))` to `(InvoiceSectionListResult,func(context.Context, InvoiceSectionListResult) (InvoiceSectionListResult, error))`
- Function `NewProductsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewAccountsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewSubscriptionsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewAccountListResultPage` signature has been changed from `(func(context.Context, AccountListResult) (AccountListResult, error))` to `(AccountListResult,func(context.Context, AccountListResult) (AccountListResult, error))`
- Function `NewRoleAssignmentsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewInvoicesClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewRoleAssignmentListResultPage` signature has been changed from `(func(context.Context, RoleAssignmentListResult) (RoleAssignmentListResult, error))` to `(RoleAssignmentListResult,func(context.Context, RoleAssignmentListResult) (RoleAssignmentListResult, error))`
- Function `NewPeriodsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewPoliciesClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewAgreementsClient` signature has been changed from `(string,string)` to `(string)`
- Function `NewPermissionsClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewInvoicesClientWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewProfileListResultPage` signature has been changed from `(func(context.Context, ProfileListResult) (ProfileListResult, error))` to `(ProfileListResult,func(context.Context, ProfileListResult) (ProfileListResult, error))`
- Function `NewPermissionsClient` signature has been changed from `(string,string)` to `(string)`
- Function `New` signature has been changed from `(string,string)` to `(string)`
- Function `NewInvoiceSectionListWithCreateSubPermissionResultPage` signature has been changed from `(func(context.Context, InvoiceSectionListWithCreateSubPermissionResult) (InvoiceSectionListWithCreateSubPermissionResult, error))` to `(InvoiceSectionListWithCreateSubPermissionResult,func(context.Context, InvoiceSectionListWithCreateSubPermissionResult) (InvoiceSectionListWithCreateSubPermissionResult, error))`
- Function `NewWithBaseURI` signature has been changed from `(string,string,string)` to `(string,string)`
- Function `NewRoleDefinitionListResultPage` signature has been changed from `(func(context.Context, RoleDefinitionListResult) (RoleDefinitionListResult, error))` to `(RoleDefinitionListResult,func(context.Context, RoleDefinitionListResult) (RoleDefinitionListResult, error))`
- Field `SubscriptionID1` of struct `BaseClient` has been removed

## New Content

- Const `InvoiceDocumentTypeCreditNote` is added
- Const `Void` is added
- Const `InvoiceDocumentTypeInvoice` is added
- Function `InvoicesClient.DownloadMultipleBillingSubscriptionInvoicesResponder(*http.Response) (DownloadURL,error)` is added
- Function `InvoicesClient.DownloadMultipleBillingProfileInvoices(context.Context,string,[]string) (InvoicesDownloadMultipleBillingProfileInvoicesFuture,error)` is added
- Function `InvoicesClient.DownloadMultipleBillingSubscriptionInvoices(context.Context,[]string) (InvoicesDownloadMultipleBillingSubscriptionInvoicesFuture,error)` is added
- Function `InvoicesClient.DownloadMultipleBillingProfileInvoicesResponder(*http.Response) (DownloadURL,error)` is added
- Function `InvoicesClient.DownloadMultipleBillingProfileInvoicesSender(*http.Request) (InvoicesDownloadMultipleBillingProfileInvoicesFuture,error)` is added
- Function `InvoicesClient.DownloadMultipleBillingProfileInvoicesPreparer(context.Context,string,[]string) (*http.Request,error)` is added
- Function `*InvoicesDownloadMultipleBillingSubscriptionInvoicesFuture.Result(InvoicesClient) (DownloadURL,error)` is added
- Function `RebillDetails.MarshalJSON() ([]byte,error)` is added
- Function `PossibleInvoiceDocumentTypeValues() []InvoiceDocumentType` is added
- Function `InvoicesClient.DownloadMultipleBillingSubscriptionInvoicesPreparer(context.Context,[]string) (*http.Request,error)` is added
- Function `InvoicesClient.DownloadMultipleBillingSubscriptionInvoicesSender(*http.Request) (InvoicesDownloadMultipleBillingSubscriptionInvoicesFuture,error)` is added
- Function `*InvoicesDownloadMultipleBillingProfileInvoicesFuture.Result(InvoicesClient) (DownloadURL,error)` is added
- Function `InvoiceProperties.MarshalJSON() ([]byte,error)` is added
- Struct `ErrorSubDetailsItem` is added
- Struct `InvoicesDownloadMultipleBillingProfileInvoicesFuture` is added
- Struct `InvoicesDownloadMultipleBillingSubscriptionInvoicesFuture` is added
- Struct `RebillDetails` is added
- Field `Details` is added to struct `ErrorDetails`
- Field `CreditForDocumentID` is added to struct `InvoiceProperties`
- Field `BilledDocumentID` is added to struct `InvoiceProperties`
- Field `RebillDetails` is added to struct `InvoiceProperties`
- Field `DocumentType` is added to struct `InvoiceProperties`


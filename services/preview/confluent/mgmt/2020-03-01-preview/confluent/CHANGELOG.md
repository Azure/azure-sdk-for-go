
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewOrganizationResourceListResultPage` signature has been changed from `(func(context.Context, OrganizationResourceListResult) (OrganizationResourceListResult, error))` to `(OrganizationResourceListResult,func(context.Context, OrganizationResourceListResult) (OrganizationResourceListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`

## New Content

- Function `NewAgreementResourceListResponsePage(AgreementResourceListResponse,func(context.Context, AgreementResourceListResponse) (AgreementResourceListResponse, error)) AgreementResourceListResponsePage` is added
- Function `MarketplaceAgreementsClient.CreateResponder(*http.Response) (AgreementResource,error)` is added
- Function `*AgreementResourceListResponseIterator.Next() error` is added
- Function `MarketplaceAgreementsClient.List(context.Context) (AgreementResourceListResponsePage,error)` is added
- Function `MarketplaceAgreementsClient.CreatePreparer(context.Context,*AgreementResource) (*http.Request,error)` is added
- Function `AgreementResourceListResponsePage.Values() []AgreementResource` is added
- Function `MarketplaceAgreementsClient.ListPreparer(context.Context) (*http.Request,error)` is added
- Function `MarketplaceAgreementsClient.ListResponder(*http.Response) (AgreementResourceListResponse,error)` is added
- Function `NewMarketplaceAgreementsClientWithBaseURI(string,string) MarketplaceAgreementsClient` is added
- Function `MarketplaceAgreementsClient.CreateSender(*http.Request) (*http.Response,error)` is added
- Function `AgreementResourceListResponsePage.Response() AgreementResourceListResponse` is added
- Function `AgreementResource.MarshalJSON() ([]byte,error)` is added
- Function `*AgreementResourceListResponsePage.Next() error` is added
- Function `AgreementResourceListResponse.IsEmpty() bool` is added
- Function `MarketplaceAgreementsClient.Create(context.Context,*AgreementResource) (AgreementResource,error)` is added
- Function `*AgreementResourceListResponseIterator.NextWithContext(context.Context) error` is added
- Function `AgreementResourceListResponsePage.NotDone() bool` is added
- Function `NewAgreementResourceListResponseIterator(AgreementResourceListResponsePage) AgreementResourceListResponseIterator` is added
- Function `AgreementResourceListResponseIterator.Value() AgreementResource` is added
- Function `AgreementResourceListResponseIterator.Response() AgreementResourceListResponse` is added
- Function `NewMarketplaceAgreementsClient(string) MarketplaceAgreementsClient` is added
- Function `MarketplaceAgreementsClient.ListSender(*http.Request) (*http.Response,error)` is added
- Function `MarketplaceAgreementsClient.ListComplete(context.Context) (AgreementResourceListResponseIterator,error)` is added
- Function `AgreementResourceListResponseIterator.NotDone() bool` is added
- Function `*AgreementResourceListResponsePage.NextWithContext(context.Context) error` is added
- Struct `AgreementProperties` is added
- Struct `AgreementResource` is added
- Struct `AgreementResourceListResponse` is added
- Struct `AgreementResourceListResponseIterator` is added
- Struct `AgreementResourceListResponsePage` is added
- Struct `MarketplaceAgreementsClient` is added


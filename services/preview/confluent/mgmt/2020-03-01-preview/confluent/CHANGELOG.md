Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewOrganizationResourceListResultPage` parameter(s) have been changed from `(func(context.Context, OrganizationResourceListResult) (OrganizationResourceListResult, error))` to `(OrganizationResourceListResult, func(context.Context, OrganizationResourceListResult) (OrganizationResourceListResult, error))`

## New Content

- New function `AgreementResourceListResponse.IsEmpty() bool`
- New function `NewMarketplaceAgreementsClientWithBaseURI(string, string) MarketplaceAgreementsClient`
- New function `MarketplaceAgreementsClient.ListResponder(*http.Response) (AgreementResourceListResponse, error)`
- New function `MarketplaceAgreementsClient.CreateSender(*http.Request) (*http.Response, error)`
- New function `NewAgreementResourceListResponseIterator(AgreementResourceListResponsePage) AgreementResourceListResponseIterator`
- New function `AgreementResourceListResponsePage.NotDone() bool`
- New function `AgreementResourceListResponsePage.Response() AgreementResourceListResponse`
- New function `MarketplaceAgreementsClient.Create(context.Context, *AgreementResource) (AgreementResource, error)`
- New function `MarketplaceAgreementsClient.ListComplete(context.Context) (AgreementResourceListResponseIterator, error)`
- New function `AgreementResourceListResponseIterator.NotDone() bool`
- New function `MarketplaceAgreementsClient.ListSender(*http.Request) (*http.Response, error)`
- New function `*AgreementResourceListResponseIterator.Next() error`
- New function `AgreementResourceListResponseIterator.Response() AgreementResourceListResponse`
- New function `*AgreementResourceListResponsePage.Next() error`
- New function `MarketplaceAgreementsClient.CreateResponder(*http.Response) (AgreementResource, error)`
- New function `AgreementResourceListResponsePage.Values() []AgreementResource`
- New function `*AgreementResourceListResponseIterator.NextWithContext(context.Context) error`
- New function `AgreementResourceListResponseIterator.Value() AgreementResource`
- New function `MarketplaceAgreementsClient.ListPreparer(context.Context) (*http.Request, error)`
- New function `MarketplaceAgreementsClient.CreatePreparer(context.Context, *AgreementResource) (*http.Request, error)`
- New function `NewAgreementResourceListResponsePage(AgreementResourceListResponse, func(context.Context, AgreementResourceListResponse) (AgreementResourceListResponse, error)) AgreementResourceListResponsePage`
- New function `AgreementResource.MarshalJSON() ([]byte, error)`
- New function `*AgreementResourceListResponsePage.NextWithContext(context.Context) error`
- New function `MarketplaceAgreementsClient.List(context.Context) (AgreementResourceListResponsePage, error)`
- New function `NewMarketplaceAgreementsClient(string) MarketplaceAgreementsClient`
- New struct `AgreementProperties`
- New struct `AgreementResource`
- New struct `AgreementResourceListResponse`
- New struct `AgreementResourceListResponseIterator`
- New struct `AgreementResourceListResponsePage`
- New struct `MarketplaceAgreementsClient`

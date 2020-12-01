Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewAPIKeyListResponsePage` parameter(s) have been changed from `(func(context.Context, APIKeyListResponse) (APIKeyListResponse, error))` to `(APIKeyListResponse, func(context.Context, APIKeyListResponse) (APIKeyListResponse, error))`
- Function `NewLinkedResourceListResponsePage` parameter(s) have been changed from `(func(context.Context, LinkedResourceListResponse) (LinkedResourceListResponse, error))` to `(LinkedResourceListResponse, func(context.Context, LinkedResourceListResponse) (LinkedResourceListResponse, error))`
- Function `NewMonitoredResourceListResponsePage` parameter(s) have been changed from `(func(context.Context, MonitoredResourceListResponse) (MonitoredResourceListResponse, error))` to `(MonitoredResourceListResponse, func(context.Context, MonitoredResourceListResponse) (MonitoredResourceListResponse, error))`
- Function `NewHostListResponsePage` parameter(s) have been changed from `(func(context.Context, HostListResponse) (HostListResponse, error))` to `(HostListResponse, func(context.Context, HostListResponse) (HostListResponse, error))`
- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewSingleSignOnResourceListResponsePage` parameter(s) have been changed from `(func(context.Context, SingleSignOnResourceListResponse) (SingleSignOnResourceListResponse, error))` to `(SingleSignOnResourceListResponse, func(context.Context, SingleSignOnResourceListResponse) (SingleSignOnResourceListResponse, error))`
- Function `NewMonitorResourceListResponsePage` parameter(s) have been changed from `(func(context.Context, MonitorResourceListResponse) (MonitorResourceListResponse, error))` to `(MonitorResourceListResponse, func(context.Context, MonitorResourceListResponse) (MonitorResourceListResponse, error))`
- Function `NewMonitoringTagRulesListResponsePage` parameter(s) have been changed from `(func(context.Context, MonitoringTagRulesListResponse) (MonitoringTagRulesListResponse, error))` to `(MonitoringTagRulesListResponse, func(context.Context, MonitoringTagRulesListResponse) (MonitoringTagRulesListResponse, error))`

## New Content

- New function `AgreementResourceListResponsePage.Values() []AgreementResource`
- New function `AgreementResourceListResponsePage.Response() AgreementResourceListResponse`
- New function `AgreementResourceListResponse.IsEmpty() bool`
- New function `NewAgreementResourceListResponseIterator(AgreementResourceListResponsePage) AgreementResourceListResponseIterator`
- New function `*AgreementResourceListResponseIterator.NextWithContext(context.Context) error`
- New function `AgreementResource.MarshalJSON() ([]byte, error)`
- New function `*AgreementResourceListResponseIterator.Next() error`
- New function `MarketplaceAgreementsClient.CreateResponder(*http.Response) (AgreementResource, error)`
- New function `MarketplaceAgreementsClient.ListPreparer(context.Context) (*http.Request, error)`
- New function `*AgreementResourceListResponsePage.Next() error`
- New function `AgreementResourceListResponseIterator.Value() AgreementResource`
- New function `MarketplaceAgreementsClient.ListComplete(context.Context) (AgreementResourceListResponseIterator, error)`
- New function `NewMarketplaceAgreementsClient(string) MarketplaceAgreementsClient`
- New function `AgreementResourceListResponsePage.NotDone() bool`
- New function `NewAgreementResourceListResponsePage(AgreementResourceListResponse, func(context.Context, AgreementResourceListResponse) (AgreementResourceListResponse, error)) AgreementResourceListResponsePage`
- New function `AgreementResourceListResponseIterator.Response() AgreementResourceListResponse`
- New function `MarketplaceAgreementsClient.List(context.Context) (AgreementResourceListResponsePage, error)`
- New function `MarketplaceAgreementsClient.ListSender(*http.Request) (*http.Response, error)`
- New function `MarketplaceAgreementsClient.Create(context.Context, *AgreementResource) (AgreementResource, error)`
- New function `MarketplaceAgreementsClient.ListResponder(*http.Response) (AgreementResourceListResponse, error)`
- New function `NewMarketplaceAgreementsClientWithBaseURI(string, string) MarketplaceAgreementsClient`
- New function `MarketplaceAgreementsClient.CreatePreparer(context.Context, *AgreementResource) (*http.Request, error)`
- New function `MarketplaceAgreementsClient.CreateSender(*http.Request) (*http.Response, error)`
- New function `AgreementResourceListResponseIterator.NotDone() bool`
- New function `*AgreementResourceListResponsePage.NextWithContext(context.Context) error`
- New struct `AgreementProperties`
- New struct `AgreementResource`
- New struct `AgreementResourceListResponse`
- New struct `AgreementResourceListResponseIterator`
- New struct `AgreementResourceListResponsePage`
- New struct `MarketplaceAgreementsClient`

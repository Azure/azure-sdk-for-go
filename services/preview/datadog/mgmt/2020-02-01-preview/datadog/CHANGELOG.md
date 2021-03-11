Generated from https://github.com/Azure/azure-rest-api-specs/tree/d9506558e6389e62926ae385f1d625a1376a0f9d/specification/datadog/resource-manager/readme.md tag: `package-2020-02-preview`

Code generator @microsoft.azure/autorest.go@2.1.175


## Breaking Changes

### Removed Funcs

1. APIKeysClient.GetDefaultKey(context.Context, string, string) (APIKey, error)
1. APIKeysClient.GetDefaultKeyPreparer(context.Context, string, string) (*http.Request, error)
1. APIKeysClient.GetDefaultKeyResponder(*http.Response) (APIKey, error)
1. APIKeysClient.GetDefaultKeySender(*http.Request) (*http.Response, error)
1. APIKeysClient.List(context.Context, string, string) (APIKeyListResponsePage, error)
1. APIKeysClient.ListComplete(context.Context, string, string) (APIKeyListResponseIterator, error)
1. APIKeysClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. APIKeysClient.ListResponder(*http.Response) (APIKeyListResponse, error)
1. APIKeysClient.ListSender(*http.Request) (*http.Response, error)
1. APIKeysClient.SetDefaultKey(context.Context, string, string, *APIKey) (autorest.Response, error)
1. APIKeysClient.SetDefaultKeyPreparer(context.Context, string, string, *APIKey) (*http.Request, error)
1. APIKeysClient.SetDefaultKeyResponder(*http.Response) (autorest.Response, error)
1. APIKeysClient.SetDefaultKeySender(*http.Request) (*http.Response, error)
1. HostsClient.List(context.Context, string, string) (HostListResponsePage, error)
1. HostsClient.ListComplete(context.Context, string, string) (HostListResponseIterator, error)
1. HostsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. HostsClient.ListResponder(*http.Response) (HostListResponse, error)
1. HostsClient.ListSender(*http.Request) (*http.Response, error)
1. LinkedResourcesClient.List(context.Context, string, string) (LinkedResourceListResponsePage, error)
1. LinkedResourcesClient.ListComplete(context.Context, string, string) (LinkedResourceListResponseIterator, error)
1. LinkedResourcesClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. LinkedResourcesClient.ListResponder(*http.Response) (LinkedResourceListResponse, error)
1. LinkedResourcesClient.ListSender(*http.Request) (*http.Response, error)
1. MarketplaceAgreementsClient.Create(context.Context, *AgreementResource) (AgreementResource, error)
1. MarketplaceAgreementsClient.CreatePreparer(context.Context, *AgreementResource) (*http.Request, error)
1. MarketplaceAgreementsClient.CreateResponder(*http.Response) (AgreementResource, error)
1. MarketplaceAgreementsClient.CreateSender(*http.Request) (*http.Response, error)
1. MonitoredResourcesClient.List(context.Context, string, string) (MonitoredResourceListResponsePage, error)
1. MonitoredResourcesClient.ListComplete(context.Context, string, string) (MonitoredResourceListResponseIterator, error)
1. MonitoredResourcesClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. MonitoredResourcesClient.ListResponder(*http.Response) (MonitoredResourceListResponse, error)
1. MonitoredResourcesClient.ListSender(*http.Request) (*http.Response, error)
1. NewAPIKeysClient(string) APIKeysClient
1. NewAPIKeysClientWithBaseURI(string, string) APIKeysClient
1. NewHostsClient(string) HostsClient
1. NewHostsClientWithBaseURI(string, string) HostsClient
1. NewLinkedResourcesClient(string) LinkedResourcesClient
1. NewLinkedResourcesClientWithBaseURI(string, string) LinkedResourcesClient
1. NewMonitoredResourcesClient(string) MonitoredResourcesClient
1. NewMonitoredResourcesClientWithBaseURI(string, string) MonitoredResourcesClient
1. NewRefreshSetPasswordClient(string) RefreshSetPasswordClient
1. NewRefreshSetPasswordClientWithBaseURI(string, string) RefreshSetPasswordClient
1. RefreshSetPasswordClient.Get(context.Context, string, string) (SetPasswordLink, error)
1. RefreshSetPasswordClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. RefreshSetPasswordClient.GetResponder(*http.Response) (SetPasswordLink, error)
1. RefreshSetPasswordClient.GetSender(*http.Request) (*http.Response, error)

## Struct Changes

### Removed Structs

1. APIKeysClient
1. HostsClient
1. LinkedResourcesClient
1. MonitoredResourcesClient
1. RefreshSetPasswordClient

### New Constants

1. MarketplaceSubscriptionStatus.Provisioning
1. MarketplaceSubscriptionStatus.Unsubscribed

### New Funcs

1. MarketplaceAgreementsClient.CreateOrUpdate(context.Context, *AgreementResource) (AgreementResource, error)
1. MarketplaceAgreementsClient.CreateOrUpdatePreparer(context.Context, *AgreementResource) (*http.Request, error)
1. MarketplaceAgreementsClient.CreateOrUpdateResponder(*http.Response) (AgreementResource, error)
1. MarketplaceAgreementsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. MonitorsClient.GetDefaultKey(context.Context, string, string) (APIKey, error)
1. MonitorsClient.GetDefaultKeyPreparer(context.Context, string, string) (*http.Request, error)
1. MonitorsClient.GetDefaultKeyResponder(*http.Response) (APIKey, error)
1. MonitorsClient.GetDefaultKeySender(*http.Request) (*http.Response, error)
1. MonitorsClient.ListAPIKeys(context.Context, string, string) (APIKeyListResponsePage, error)
1. MonitorsClient.ListAPIKeysComplete(context.Context, string, string) (APIKeyListResponseIterator, error)
1. MonitorsClient.ListAPIKeysPreparer(context.Context, string, string) (*http.Request, error)
1. MonitorsClient.ListAPIKeysResponder(*http.Response) (APIKeyListResponse, error)
1. MonitorsClient.ListAPIKeysSender(*http.Request) (*http.Response, error)
1. MonitorsClient.ListHosts(context.Context, string, string) (HostListResponsePage, error)
1. MonitorsClient.ListHostsComplete(context.Context, string, string) (HostListResponseIterator, error)
1. MonitorsClient.ListHostsPreparer(context.Context, string, string) (*http.Request, error)
1. MonitorsClient.ListHostsResponder(*http.Response) (HostListResponse, error)
1. MonitorsClient.ListHostsSender(*http.Request) (*http.Response, error)
1. MonitorsClient.ListLinkedResources(context.Context, string, string) (LinkedResourceListResponsePage, error)
1. MonitorsClient.ListLinkedResourcesComplete(context.Context, string, string) (LinkedResourceListResponseIterator, error)
1. MonitorsClient.ListLinkedResourcesPreparer(context.Context, string, string) (*http.Request, error)
1. MonitorsClient.ListLinkedResourcesResponder(*http.Response) (LinkedResourceListResponse, error)
1. MonitorsClient.ListLinkedResourcesSender(*http.Request) (*http.Response, error)
1. MonitorsClient.ListMonitoredResources(context.Context, string, string) (MonitoredResourceListResponsePage, error)
1. MonitorsClient.ListMonitoredResourcesComplete(context.Context, string, string) (MonitoredResourceListResponseIterator, error)
1. MonitorsClient.ListMonitoredResourcesPreparer(context.Context, string, string) (*http.Request, error)
1. MonitorsClient.ListMonitoredResourcesResponder(*http.Response) (MonitoredResourceListResponse, error)
1. MonitorsClient.ListMonitoredResourcesSender(*http.Request) (*http.Response, error)
1. MonitorsClient.RefreshSetPasswordLink(context.Context, string, string) (SetPasswordLink, error)
1. MonitorsClient.RefreshSetPasswordLinkPreparer(context.Context, string, string) (*http.Request, error)
1. MonitorsClient.RefreshSetPasswordLinkResponder(*http.Response) (SetPasswordLink, error)
1. MonitorsClient.RefreshSetPasswordLinkSender(*http.Request) (*http.Response, error)
1. MonitorsClient.SetDefaultKey(context.Context, string, string, *APIKey) (autorest.Response, error)
1. MonitorsClient.SetDefaultKeyPreparer(context.Context, string, string, *APIKey) (*http.Request, error)
1. MonitorsClient.SetDefaultKeyResponder(*http.Response) (autorest.Response, error)
1. MonitorsClient.SetDefaultKeySender(*http.Request) (*http.Response, error)
1. SingleSignOnProperties.MarshalJSON() ([]byte, error)

## Struct Changes

### New Struct Fields

1. MonitoringTagRulesProperties.ProvisioningState
1. OrganizationProperties.APIKey
1. OrganizationProperties.ApplicationKey
1. OrganizationProperties.RedirectURI
1. SingleSignOnProperties.ProvisioningState

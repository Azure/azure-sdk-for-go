# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. ProductsClient.GetProducts
	- Params
		- From: context.Context, string, string, *DeviceConfiguration
		- To: context.Context, string, string, string, *DeviceConfiguration
1. ProductsClient.GetProductsPreparer
	- Params
		- From: context.Context, string, string, *DeviceConfiguration
		- To: context.Context, string, string, string, *DeviceConfiguration

## Additive Changes

### New Funcs

1. RegistrationsClient.EnableRemoteManagement(context.Context, string, string) (autorest.Response, error)
1. RegistrationsClient.EnableRemoteManagementPreparer(context.Context, string, string) (*http.Request, error)
1. RegistrationsClient.EnableRemoteManagementResponder(*http.Response) (autorest.Response, error)
1. RegistrationsClient.EnableRemoteManagementSender(*http.Request) (*http.Response, error)
1. RegistrationsClient.ListBySubscription(context.Context) (RegistrationListPage, error)
1. RegistrationsClient.ListBySubscriptionComplete(context.Context) (RegistrationListIterator, error)
1. RegistrationsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. RegistrationsClient.ListBySubscriptionResponder(*http.Response) (RegistrationList, error)
1. RegistrationsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)

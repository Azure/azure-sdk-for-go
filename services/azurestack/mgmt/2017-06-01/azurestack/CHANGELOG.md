Generated from https://github.com/Azure/azure-rest-api-specs/tree/../../../../../azure-rest-api-specs/specification/azurestack/resource-manager/readme.md tag: `package-2017-06-01`

Code generator 


### Breaking Changes

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

#### New Funcs

1. RegistrationsClient.EnableRemoteManagement(context.Context, string, string) (autorest.Response, error)
1. RegistrationsClient.EnableRemoteManagementPreparer(context.Context, string, string) (*http.Request, error)
1. RegistrationsClient.EnableRemoteManagementResponder(*http.Response) (autorest.Response, error)
1. RegistrationsClient.EnableRemoteManagementSender(*http.Request) (*http.Response, error)

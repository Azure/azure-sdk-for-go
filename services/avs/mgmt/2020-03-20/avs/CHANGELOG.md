
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewHcxEnterpriseSiteListPage` signature has been changed from `(func(context.Context, HcxEnterpriseSiteList) (HcxEnterpriseSiteList, error))` to `(HcxEnterpriseSiteList,func(context.Context, HcxEnterpriseSiteList) (HcxEnterpriseSiteList, error))`
- Function `NewPrivateCloudListPage` signature has been changed from `(func(context.Context, PrivateCloudList) (PrivateCloudList, error))` to `(PrivateCloudList,func(context.Context, PrivateCloudList) (PrivateCloudList, error))`
- Function `NewClusterListPage` signature has been changed from `(func(context.Context, ClusterList) (ClusterList, error))` to `(ClusterList,func(context.Context, ClusterList) (ClusterList, error))`
- Function `AuthorizationsClient.CreateOrUpdatePreparer` signature has been changed from `(context.Context,string,string,string,interface{})` to `(context.Context,string,string,string,ExpressRouteAuthorization)`
- Function `HcxEnterpriseSitesClient.CreateOrUpdatePreparer` signature has been changed from `(context.Context,string,string,string,interface{})` to `(context.Context,string,string,string,HcxEnterpriseSite)`
- Function `AuthorizationsClient.CreateOrUpdate` signature has been changed from `(context.Context,string,string,string,interface{})` to `(context.Context,string,string,string,ExpressRouteAuthorization)`
- Function `NewOperationListPage` signature has been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList,func(context.Context, OperationList) (OperationList, error))`
- Function `HcxEnterpriseSitesClient.CreateOrUpdate` signature has been changed from `(context.Context,string,string,string,interface{})` to `(context.Context,string,string,string,HcxEnterpriseSite)`
- Function `NewExpressRouteAuthorizationListPage` signature has been changed from `(func(context.Context, ExpressRouteAuthorizationList) (ExpressRouteAuthorizationList, error))` to `(ExpressRouteAuthorizationList,func(context.Context, ExpressRouteAuthorizationList) (ExpressRouteAuthorizationList, error))`

## New Content

- Function `Operation.MarshalJSON() ([]byte,error)` is added
- Struct `LogSpecification` is added
- Struct `MetricDimension` is added
- Struct `MetricSpecification` is added
- Struct `OperationProperties` is added
- Struct `ServiceSpecification` is added
- Field `ProvisioningState` is added to struct `ManagementCluster`
- Field `Properties` is added to struct `Operation`
- Field `IsDataAction` is added to struct `Operation`
- Field `Origin` is added to struct `Operation`


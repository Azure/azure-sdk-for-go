# Release History

## 0.3.0 (2021-11-18)
### Breaking Changes

- Function `*SyncGroupsClient.ListLogs` parameter(s) have been changed from `(string, string, string, string, string, string, Enum75, *SyncGroupsListLogsOptions)` to `(string, string, string, string, string, string, Enum74, *SyncGroupsListLogsOptions)`
- Const `Enum75Success` has been removed
- Const `Enum75All` has been removed
- Const `Enum75Error` has been removed
- Const `Enum75Warning` has been removed
- Function `Enum75.ToPtr` has been removed
- Function `*ServerConnectionPoliciesClient.CreateOrUpdate` has been removed
- Function `PossibleEnum75Values` has been removed
- Struct `ServerConnectionPoliciesCreateOrUpdateOptions` has been removed

### New Content

- New const `Enum74Warning`
- New const `Enum74All`
- New const `Enum74Error`
- New const `Enum74Success`
- New function `ServerConnectionPoliciesCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (ServerConnectionPoliciesCreateOrUpdateResponse, error)`
- New function `ServerConnectionPolicyListResult.MarshalJSON() ([]byte, error)`
- New function `*ServerConnectionPoliciesCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ServerConnectionPoliciesListByServerPager.NextPage(context.Context) bool`
- New function `*ServerConnectionPoliciesClient.ListByServer(string, string, *ServerConnectionPoliciesListByServerOptions) *ServerConnectionPoliciesListByServerPager`
- New function `Enum74.ToPtr() *Enum74`
- New function `*ServerConnectionPoliciesListByServerPager.Err() error`
- New function `*ServerConnectionPoliciesCreateOrUpdatePollerResponse.Resume(context.Context, *ServerConnectionPoliciesClient, string) error`
- New function `PossibleEnum74Values() []Enum74`
- New function `*ServerConnectionPoliciesListByServerPager.PageResponse() ServerConnectionPoliciesListByServerResponse`
- New function `*ServerConnectionPoliciesCreateOrUpdatePoller.Done() bool`
- New function `*ServerConnectionPoliciesClient.BeginCreateOrUpdate(context.Context, string, string, ConnectionPolicyName, ServerConnectionPolicy, *ServerConnectionPoliciesBeginCreateOrUpdateOptions) (ServerConnectionPoliciesCreateOrUpdatePollerResponse, error)`
- New function `*ServerConnectionPoliciesCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*ServerConnectionPoliciesCreateOrUpdatePoller.FinalResponse(context.Context) (ServerConnectionPoliciesCreateOrUpdateResponse, error)`
- New struct `ServerConnectionPoliciesBeginCreateOrUpdateOptions`
- New struct `ServerConnectionPoliciesCreateOrUpdatePoller`
- New struct `ServerConnectionPoliciesCreateOrUpdatePollerResponse`
- New struct `ServerConnectionPoliciesListByServerOptions`
- New struct `ServerConnectionPoliciesListByServerPager`
- New struct `ServerConnectionPoliciesListByServerResponse`
- New struct `ServerConnectionPoliciesListByServerResult`
- New struct `ServerConnectionPolicyListResult`

Total 10 breaking change(s), 34 additive change(s).


## 0.2.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.1 (2021-10-14)
- fix wrong module path in go.mod

## 0.1.0 (2021-10-08)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/sql/armsql") to avoid confusion.

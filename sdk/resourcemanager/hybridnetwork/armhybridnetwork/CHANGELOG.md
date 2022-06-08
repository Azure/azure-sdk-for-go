# Release History

## 2.0.0-beta.1 (2022-05-24)
### Breaking Changes

- Struct `OperationList` has been removed
- Field `OperationList` of struct `OperationsClientListResponse` has been removed

### Features Added

- New const `HTTPMethodDelete`
- New const `HTTPMethodPost`
- New const `HTTPMethodPatch`
- New const `HTTPMethodUnknown`
- New const `OriginSystem`
- New const `ActionTypeInternal`
- New const `HTTPMethodPut`
- New const `HTTPMethodGet`
- New const `OriginUser`
- New const `OriginUserSystem`
- New function `PossibleHTTPMethodValues() []HTTPMethod`
- New function `*NetworkFunctionsClient.BeginExecuteRequest(context.Context, string, string, ExecuteRequestParameters, *NetworkFunctionsClientBeginExecuteRequestOptions) (*runtime.Poller[NetworkFunctionsClientExecuteRequestResponse], error)`
- New function `*SKUCredential.UnmarshalJSON([]byte) error`
- New function `PossibleOriginValues() []Origin`
- New function `*VendorSKUsClient.ListCredential(context.Context, string, string, *VendorSKUsClientListCredentialOptions) (VendorSKUsClientListCredentialResponse, error)`
- New function `PossibleActionTypeValues() []ActionType`
- New struct `ExecuteRequestParameters`
- New struct `NetworkFunctionsClientBeginExecuteRequestOptions`
- New struct `NetworkFunctionsClientExecuteRequestResponse`
- New struct `OperationListResult`
- New struct `RequestMetadata`
- New struct `SKUCredential`
- New struct `VendorSKUsClientListCredentialOptions`
- New struct `VendorSKUsClientListCredentialResponse`
- New field `Origin` in struct `Operation`
- New field `ActionType` in struct `Operation`
- New field `IsDataAction` in struct `Operation`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridnetwork/armhybridnetwork` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
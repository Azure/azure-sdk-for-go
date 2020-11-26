
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewResourceListResultPage` signature has been changed from `(func(context.Context, ResourceListResult) (ResourceListResult, error))` to `(ResourceListResult,func(context.Context, ResourceListResult) (ResourceListResult, error))`
- Function `NewDeletedVaultListResultPage` signature has been changed from `(func(context.Context, DeletedVaultListResult) (DeletedVaultListResult, error))` to `(DeletedVaultListResult,func(context.Context, DeletedVaultListResult) (DeletedVaultListResult, error))`
- Function `NewVaultListResultPage` signature has been changed from `(func(context.Context, VaultListResult) (VaultListResult, error))` to `(VaultListResult,func(context.Context, VaultListResult) (VaultListResult, error))`

## New Content

- Const `P384` is added
- Const `P256K` is added
- Const `RSAHSM` is added
- Const `RecoverablePurgeable` is added
- Const `JSONWebKeyOperationDecrypt` is added
- Const `StoragePermissionsAll` is added
- Const `ECHSM` is added
- Const `RSA` is added
- Const `SecretPermissionsAll` is added
- Const `Purgeable` is added
- Const `JSONWebKeyOperationVerify` is added
- Const `JSONWebKeyOperationWrapKey` is added
- Const `JSONWebKeyOperationEncrypt` is added
- Const `P521` is added
- Const `JSONWebKeyOperationImport` is added
- Const `P256` is added
- Const `JSONWebKeyOperationSign` is added
- Const `All` is added
- Const `JSONWebKeyOperationUnwrapKey` is added
- Const `EC` is added
- Const `Recoverable` is added
- Const `RecoverableProtectedSubscription` is added
- Const `KeyPermissionsAll` is added
- Function `KeyCreateParameters.MarshalJSON() ([]byte,error)` is added
- Function `KeysClient.GetVersionResponder(*http.Response) (Key,error)` is added
- Function `DeletedVault.MarshalJSON() ([]byte,error)` is added
- Function `KeyProperties.MarshalJSON() ([]byte,error)` is added
- Function `KeysClient.ListResponder(*http.Response) (KeyListResult,error)` is added
- Function `Attributes.MarshalJSON() ([]byte,error)` is added
- Function `KeysClient.ListPreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `KeysClient.ListVersionsSender(*http.Request) (*http.Response,error)` is added
- Function `KeysClient.GetSender(*http.Request) (*http.Response,error)` is added
- Function `KeysClient.GetVersionPreparer(context.Context,string,string,string,string) (*http.Request,error)` is added
- Function `KeyAttributes.MarshalJSON() ([]byte,error)` is added
- Function `NewKeysClientWithBaseURI(string,string) KeysClient` is added
- Function `Key.MarshalJSON() ([]byte,error)` is added
- Function `PossibleJSONWebKeyOperationValues() []JSONWebKeyOperation` is added
- Function `KeyListResultIterator.Response() KeyListResult` is added
- Function `KeysClient.ListVersionsResponder(*http.Response) (KeyListResult,error)` is added
- Function `KeysClient.Get(context.Context,string,string,string) (Key,error)` is added
- Function `*KeyListResultPage.NextWithContext(context.Context) error` is added
- Function `VaultAccessPolicyParameters.MarshalJSON() ([]byte,error)` is added
- Function `PrivateLinkResourceProperties.MarshalJSON() ([]byte,error)` is added
- Function `KeysClient.List(context.Context,string,string) (KeyListResultPage,error)` is added
- Function `VaultProperties.MarshalJSON() ([]byte,error)` is added
- Function `KeysClient.ListVersions(context.Context,string,string,string) (KeyListResultPage,error)` is added
- Function `KeysClient.GetVersion(context.Context,string,string,string,string) (Key,error)` is added
- Function `*KeyListResultIterator.Next() error` is added
- Function `*KeyListResultPage.Next() error` is added
- Function `KeysClient.GetVersionSender(*http.Request) (*http.Response,error)` is added
- Function `PossibleDeletionRecoveryLevelValues() []DeletionRecoveryLevel` is added
- Function `KeyListResultPage.Response() KeyListResult` is added
- Function `KeysClient.ListVersionsComplete(context.Context,string,string,string) (KeyListResultIterator,error)` is added
- Function `KeysClient.ListComplete(context.Context,string,string) (KeyListResultIterator,error)` is added
- Function `KeysClient.CreateIfNotExist(context.Context,string,string,string,KeyCreateParameters) (Key,error)` is added
- Function `*Key.UnmarshalJSON([]byte) error` is added
- Function `KeysClient.GetResponder(*http.Response) (Key,error)` is added
- Function `KeysClient.CreateIfNotExistPreparer(context.Context,string,string,string,KeyCreateParameters) (*http.Request,error)` is added
- Function `KeyListResultIterator.Value() Key` is added
- Function `KeysClient.ListVersionsPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `KeyListResultPage.Values() []Key` is added
- Function `NewKeyListResultIterator(KeyListResultPage) KeyListResultIterator` is added
- Function `PossibleJSONWebKeyTypeValues() []JSONWebKeyType` is added
- Function `KeyListResultIterator.NotDone() bool` is added
- Function `KeysClient.GetPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `PossibleJSONWebKeyCurveNameValues() []JSONWebKeyCurveName` is added
- Function `KeysClient.CreateIfNotExistSender(*http.Request) (*http.Response,error)` is added
- Function `KeysClient.CreateIfNotExistResponder(*http.Response) (Key,error)` is added
- Function `*KeyListResultIterator.NextWithContext(context.Context) error` is added
- Function `NewKeysClient(string) KeysClient` is added
- Function `NewKeyListResultPage(KeyListResult,func(context.Context, KeyListResult) (KeyListResult, error)) KeyListResultPage` is added
- Function `KeyListResult.IsEmpty() bool` is added
- Function `KeysClient.ListSender(*http.Request) (*http.Response,error)` is added
- Function `KeyListResultPage.NotDone() bool` is added
- Struct `Attributes` is added
- Struct `Key` is added
- Struct `KeyAttributes` is added
- Struct `KeyCreateParameters` is added
- Struct `KeyListResult` is added
- Struct `KeyListResultIterator` is added
- Struct `KeyListResultPage` is added
- Struct `KeyProperties` is added
- Struct `KeysClient` is added


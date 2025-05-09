# Release History

## 1.3.0-beta.3 (2025-05-23)
### Features Added

- New enum type `ApplicationType` with values `ApplicationTypeCodeTransparency`, `ApplicationTypeConfidentialLedger`
- New enum type `EnclavePlatform` with values `EnclavePlatformAmdSevSnp`, `EnclavePlatformIntelSgx`
- New field `ApplicationType`, `EnclavePlatform`, `HostLevel`, `MaxBodySizeInMb`, `NodeCount`, `SubjectName`, `WorkerThreads`, `WriteLBAddressPrefix` in struct `LedgerProperties`
- New field `EnclavePlatform` in struct `ManagedCCFProperties`


## 1.3.0-beta.2 (2024-04-26)
### Features Added

- New enum type `LedgerSKU` with values `LedgerSKUBasic`, `LedgerSKUStandard`, `LedgerSKUUnknown`
- New function `*LedgerClient.BeginRestore(context.Context, string, string, Restore, *LedgerClientBeginRestoreOptions) (*runtime.Poller[LedgerClientRestoreResponse], error)`
- New function `*LedgerClient.BeginBackup(context.Context, string, string, Backup, *LedgerClientBeginBackupOptions) (*runtime.Poller[LedgerClientBackupResponse], error)`
- New function `*ManagedCCFClient.BeginRestore(context.Context, string, string, ManagedCCFRestore, *ManagedCCFClientBeginRestoreOptions) (*runtime.Poller[ManagedCCFClientRestoreResponse], error)`
- New function `*ManagedCCFClient.BeginBackup(context.Context, string, string, ManagedCCFBackup, *ManagedCCFClientBeginBackupOptions) (*runtime.Poller[ManagedCCFClientBackupResponse], error)`
- New struct `Backup`
- New struct `BackupResponse`
- New struct `ManagedCCFBackup`
- New struct `ManagedCCFBackupResponse`
- New struct `ManagedCCFRestore`
- New struct `ManagedCCFRestoreResponse`
- New struct `Restore`
- New struct `RestoreResponse`
- New field `LedgerSKU` in struct `LedgerProperties`
- New anonymous field `ManagedCCF` in struct `ManagedCCFClientUpdateResponse`
- New field `RunningState` in struct `ManagedCCFProperties`


## 1.3.0-beta.1 (2023-11-30)
### Features Added

- New enum type `LanguageRuntime` with values `LanguageRuntimeCPP`, `LanguageRuntimeJS`
- New enum type `RunningState` with values `RunningStateActive`, `RunningStatePaused`, `RunningStatePausing`, `RunningStateResuming`, `RunningStateUnknown`
- New function `*ClientFactory.NewManagedCCFClient() *ManagedCCFClient`
- New function `NewManagedCCFClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedCCFClient, error)`
- New function `*ManagedCCFClient.BeginCreate(context.Context, string, string, ManagedCCF, *ManagedCCFClientBeginCreateOptions) (*runtime.Poller[ManagedCCFClientCreateResponse], error)`
- New function `*ManagedCCFClient.BeginDelete(context.Context, string, string, *ManagedCCFClientBeginDeleteOptions) (*runtime.Poller[ManagedCCFClientDeleteResponse], error)`
- New function `*ManagedCCFClient.Get(context.Context, string, string, *ManagedCCFClientGetOptions) (ManagedCCFClientGetResponse, error)`
- New function `*ManagedCCFClient.NewListByResourceGroupPager(string, *ManagedCCFClientListByResourceGroupOptions) *runtime.Pager[ManagedCCFClientListByResourceGroupResponse]`
- New function `*ManagedCCFClient.NewListBySubscriptionPager(*ManagedCCFClientListBySubscriptionOptions) *runtime.Pager[ManagedCCFClientListBySubscriptionResponse]`
- New function `*ManagedCCFClient.BeginUpdate(context.Context, string, string, ManagedCCF, *ManagedCCFClientBeginUpdateOptions) (*runtime.Poller[ManagedCCFClientUpdateResponse], error)`
- New struct `DeploymentType`
- New struct `ManagedCCF`
- New struct `ManagedCCFList`
- New struct `ManagedCCFProperties`
- New struct `MemberIdentityCertificate`
- New field `RunningState` in struct `LedgerProperties`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0-beta.1 (2023-04-28)

### Features Added

- New enum type `LanguageRuntime` with values `LanguageRuntimeCPP`, `LanguageRuntimeJS`
- New enum type `RunningState` with values `RunningStateActive`, `RunningStatePaused`, `RunningStatePausing`, `RunningStateResuming`, `RunningStateUnknown`
- New function `*ClientFactory.NewManagedCCFClient() *ManagedCCFClient`
- New function `NewManagedCCFClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedCCFClient, error)`
- New function `*ManagedCCFClient.BeginCreate(context.Context, string, string, ManagedCCF, *ManagedCCFClientBeginCreateOptions) (*runtime.Poller[ManagedCCFClientCreateResponse], error)`
- New function `*ManagedCCFClient.BeginDelete(context.Context, string, string, *ManagedCCFClientBeginDeleteOptions) (*runtime.Poller[ManagedCCFClientDeleteResponse], error)`
- New function `*ManagedCCFClient.Get(context.Context, string, string, *ManagedCCFClientGetOptions) (ManagedCCFClientGetResponse, error)`
- New function `*ManagedCCFClient.NewListByResourceGroupPager(string, *ManagedCCFClientListByResourceGroupOptions) *runtime.Pager[ManagedCCFClientListByResourceGroupResponse]`
- New function `*ManagedCCFClient.NewListBySubscriptionPager(*ManagedCCFClientListBySubscriptionOptions) *runtime.Pager[ManagedCCFClientListBySubscriptionResponse]`
- New function `*ManagedCCFClient.BeginUpdate(context.Context, string, string, ManagedCCF, *ManagedCCFClientBeginUpdateOptions) (*runtime.Poller[ManagedCCFClientUpdateResponse], error)`
- New struct `DeploymentType`
- New struct `ManagedCCF`
- New struct `ManagedCCFList`
- New struct `ManagedCCFProperties`
- New struct `MemberIdentityCertificate`
- New field `RunningState` in struct `LedgerProperties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/confidentialledger/armconfidentialledger` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
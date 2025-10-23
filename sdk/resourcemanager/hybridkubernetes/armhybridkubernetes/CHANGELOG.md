# Release History

## 2.0.0-beta.2 (2025-10-23)
### Breaking Changes

- Type of `SystemData.LastModifiedByType` has been changed from `*LastModifiedByType` to `*CreatedByType`
- Enum `LastModifiedByType` has been removed
- Function `*ConnectedClusterClient.Update` has been removed
- Struct `ErrorAdditionalInfo` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorResponse` has been removed
- Struct `Resource` has been removed
- Struct `TrackedResource` has been removed
- Field `ResourceID` of struct `Gateway` has been removed

### Features Added

- New value `ConnectedClusterKindGCP` added to enum type `ConnectedClusterKind`
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New function `*ConnectedClusterClient.BeginUpdateAsync(context.Context, string, string, ConnectedClusterPatch, *ConnectedClusterClientBeginUpdateAsyncOptions) (*runtime.Poller[ConnectedClusterClientUpdateAsyncResponse], error)`
- New field `Gateway` in struct `ConnectedClusterPatchProperties`
- New field `ActionType`, `IsDataAction`, `Origin` in struct `Operation`


## 2.0.0-beta.1 (2025-03-17)
### Breaking Changes

- Type of `ConnectedClusterPatch.Properties` has been changed from `any` to `*ConnectedClusterPatchProperties`
- function  `*ConnectedClusterClient.BeginCreate` has been renamed to `*ConnectedClusterClient.BeginCreateOrReplace`

### Features Added

- New value `ConnectivityStatusAgentNotInstalled` added to enum type `ConnectivityStatus`
- New enum type `AutoUpgradeOptions` with values `AutoUpgradeOptionsDisabled`, `AutoUpgradeOptionsEnabled`
- New enum type `AzureHybridBenefit` with values `AzureHybridBenefitFalse`, `AzureHybridBenefitNotApplicable`, `AzureHybridBenefitTrue`
- New enum type `ConnectedClusterKind` with values `ConnectedClusterKindAWS`, `ConnectedClusterKindProvisionedCluster`
- New enum type `PrivateLinkState` with values `PrivateLinkStateDisabled`, `PrivateLinkStateEnabled`
- New struct `AADProfile`
- New struct `AgentError`
- New struct `ArcAgentProfile`
- New struct `ArcAgentryConfigurations`
- New struct `ConnectedClusterPatchProperties`
- New struct `Gateway`
- New struct `OidcIssuerProfile`
- New struct `SecurityProfile`
- New struct `SecurityProfileWorkloadIdentity`
- New struct `SystemComponent`
- New field `Kind` in struct `ConnectedCluster`
- New field `AADProfile`, `ArcAgentProfile`, `ArcAgentryConfigurations`, `AzureHybridBenefit`, `DistributionVersion`, `Gateway`, `MiscellaneousProperties`, `OidcIssuerProfile`, `PrivateLinkScopeResourceID`, `PrivateLinkState`, `SecurityProfile` in struct `ConnectedClusterProperties`
- New field `RelayTid`, `RelayType` in struct `HybridConnectionConfig`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridkubernetes/armhybridkubernetes` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
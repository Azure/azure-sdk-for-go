# Release History

## 2.0.0-beta.1 (2026-11-20)
### Breaking Changes

- Struct `AddressSpacesModel` has been removed
- Struct `VirtualNetworkModel` has been removed
- Field `EnclaveVirtualNetwork` of struct `CheckAddressSpaceAvailabilityRequest` has been removed
- Field `EnclaveVirtualNetwork` of struct `VirtualEnclavePatchProperties` has been removed
- Field `EnclaveAddressSpaces`, `EnclaveVirtualNetwork` of struct `VirtualEnclaveProperties` has been removed

### Features Added

- New value `ConnectionStateDegraded`, `ConnectionStateReconciling` added to enum type `ConnectionState`
- New value `DiagnosticDestinationUnspecified` added to enum type `DiagnosticDestination`
- New value `ServiceIdentifierAzureLocal` added to enum type `ServiceIdentifier`
- New enum type `AddressSpaceOwnership` with values `AddressSpaceOwnershipCustomerManaged`, `AddressSpaceOwnershipProductManaged`
- New enum type `RuleDirection` with values `RuleDirectionInbound`, `RuleDirectionOutbound`
- New enum type `SubnetJoinPolicyMode` with values `SubnetJoinPolicyModeAudit`, `SubnetJoinPolicyModeDeny`
- New enum type `SubnetPurpose` with values `SubnetPurposeAKS`, `SubnetPurposeAKSAPIServer`, `SubnetPurposeGeneral`
- New enum type `SyncStatus` with values `SyncStatusFailed`, `SyncStatusInProgress`, `SyncStatusNotConfigured`, `SyncStatusPartiallySucceeded`, `SyncStatusSucceeded`
- New function `*CommunityEndpointsClient.BeginSync(ctx context.Context, resourceGroupName string, communityName string, communityEndpointName string, options *CommunityEndpointsClientBeginSyncOptions) (*runtime.Poller[CommunityEndpointsClientSyncResponse], error)`
- New struct `AddressSpaceEntry`
- New struct `ComputedAddressSpaceEntry`
- New struct `SyncDestinationStatus`
- New struct `SyncPerformedBy`
- New field `AddressSpace`, `PrefixLength` in struct `CheckAddressSpaceAvailabilityRequest`
- New field `Direction` in struct `CommunityEndpointDestinationRule`
- New field `SyncTo` in struct `CommunityEndpointPatchProperties`
- New field `SyncDestinationStatuses`, `SyncFrom`, `SyncPerformedBy`, `SyncStatus`, `SyncTo` in struct `CommunityEndpointProperties`
- New field `SubnetJoinPolicyMode` in struct `CommunityPatchProperties`
- New field `SubnetJoinPolicyMode` in struct `CommunityProperties`
- New field `SubnetPurpose` in struct `SubnetConfiguration`
- New field `AddressSpaces`, `SubnetConfigurations` in struct `VirtualEnclavePatchProperties`
- New field `AddressSpaces`, `AllowSubnetCommunication`, `ComputedAddressSpaces`, `NetworkName`, `SubnetConfigurations` in struct `VirtualEnclaveProperties`


## 1.0.0 (2026-08-27)
### Other Changes

- First stable (GA) release of the `armenclave` module, targeting API version `2026-04-01`. The API surface is unchanged from `0.1.0`; this release promotes the module to a stable, generally available version.

## 0.1.0 (2026-07-17)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/enclave/armenclave` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
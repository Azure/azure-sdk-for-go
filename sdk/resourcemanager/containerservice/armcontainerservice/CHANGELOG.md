# Release History

## 0.2.0 (2021-10-29)
### Breaking Changes

- Function `NewPrivateLinkResourcesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateEndpointConnectionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewMaintenanceConfigurationsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAgentPoolsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewResolvePrivateLinkServiceIDClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewManagedClustersClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSnapshotsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`

### New Content

- New struct `WindowsGmsaProfile`
- New field `EnableFIPS` in struct `SnapshotProperties`
- New field `KubernetesVersion` in struct `SnapshotProperties`
- New field `NodeImageVersion` in struct `SnapshotProperties`
- New field `OSSKU` in struct `SnapshotProperties`
- New field `OSType` in struct `SnapshotProperties`
- New field `VMSize` in struct `SnapshotProperties`
- New field `GmsaProfile` in struct `ManagedClusterWindowsProfile`

Total 8 breaking change(s), 4 additive change(s).


## 0.1.0 (2021-10-08)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/containerservice/armcontainerservice") to avoid confusion.
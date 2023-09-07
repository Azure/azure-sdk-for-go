# Release History

## 0.2.0 (2023-09-22)
### Breaking Changes

- Operation `*FleetMembersClient.Update` has been changed to LRO, use `*FleetMembersClient.BeginUpdate` instead.
- Operation `*FleetsClient.Update` has been changed to LRO, use `*FleetsClient.BeginUpdate` instead.

### Features Added

- New value `UpdateStateSkipped` added to enum type `UpdateState`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `NodeImageSelectionType` with values `NodeImageSelectionTypeConsistent`, `NodeImageSelectionTypeLatest`
- New struct `APIServerAccessProfile`
- New struct `AgentProfile`
- New struct `ManagedServiceIdentity`
- New struct `NodeImageSelection`
- New struct `NodeImageSelectionStatus`
- New struct `NodeImageVersion`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Fleet`
- New field `APIServerAccessProfile`, `AgentProfile` in struct `FleetHubProfile`
- New field `Identity` in struct `FleetPatch`
- New field `NodeImageSelection` in struct `ManagedClusterUpdate`
- New field `Message` in struct `MemberUpdateStatus`
- New field `NodeImageSelection` in struct `UpdateRunStatus`


## 0.1.0 (2023-06-23)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservicefleet/armcontainerservicefleet` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
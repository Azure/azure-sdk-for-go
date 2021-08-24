# Release History

## v0.2.0 (released)
### Breaking Changes

- Type of `ClusterSKU.Capacity` has been changed from `*int64` to `*Capacity`
- Type of `WorkspaceSKU.CapacityReservationLevel` has been changed from `*int32` to `*CapacityReservationLevel`
- Function `Table.MarshalJSON` has been removed
- Function `NewTablesClient` has been removed
- Function `*TablesClient.Get` has been removed
- Function `*TablesClient.ListByWorkspace` has been removed
- Function `TablesListResult.MarshalJSON` has been removed
- Function `*TablesClient.Update` has been removed
- Struct `Table` has been removed
- Struct `TableProperties` has been removed
- Struct `TablesClient` has been removed
- Struct `TablesGetOptions` has been removed
- Struct `TablesGetResponse` has been removed
- Struct `TablesGetResult` has been removed
- Struct `TablesListByWorkspaceOptions` has been removed
- Struct `TablesListByWorkspaceResponse` has been removed
- Struct `TablesListByWorkspaceResult` has been removed
- Struct `TablesListResult` has been removed
- Struct `TablesUpdateOptions` has been removed
- Struct `TablesUpdateResponse` has been removed
- Struct `TablesUpdateResult` has been removed

### New Content

- New const `CapacityReservationLevelFourHundred`
- New const `CapacityReservationLevelFiveThousand`
- New const `CapacityTwoThousand`
- New const `CapacityReservationLevelTenHundred`
- New const `CapacityTenHundred`
- New const `CapacityFiveThousand`
- New const `CapacityReservationLevelThreeHundred`
- New const `CapacityReservationLevelTwoHundred`
- New const `CapacityReservationLevelFiveHundred`
- New const `CapacityReservationLevelOneHundred`
- New const `CapacityReservationLevelTwoThousand`
- New const `CapacityFiveHundred`
- New function `PossibleCapacityReservationLevelValues() []CapacityReservationLevel`
- New function `Capacity.ToPtr() *Capacity`
- New function `CapacityReservationLevel.ToPtr() *CapacityReservationLevel`
- New function `PossibleCapacityValues() []Capacity`

Total 34 breaking change(s), 16 additive change(s).


## v0.1.0 (released)

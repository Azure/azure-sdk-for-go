# Release History

## 0.7.0 (2022-10-18)
### Features Added

- New const `ChangeTypeCreate`
- New const `PropertyChangeTypeInsert`
- New const `ChangeTypeUpdate`
- New const `PropertyChangeTypeUpdate`
- New const `ChangeTypeDelete`
- New const `PropertyChangeTypeRemove`
- New const `ChangeCategorySystem`
- New const `ChangeCategoryUser`
- New type alias `ChangeType`
- New type alias `PropertyChangeType`
- New type alias `ChangeCategory`
- New function `*Client.ResourceChanges(context.Context, ResourceChangesRequestParameters, *ClientResourceChangesOptions) (ClientResourceChangesResponse, error)`
- New function `PossibleChangeCategoryValues() []ChangeCategory`
- New function `PossiblePropertyChangeTypeValues() []PropertyChangeType`
- New function `*Client.ResourceChangeDetails(context.Context, ResourceChangeDetailsRequestParameters, *ClientResourceChangeDetailsOptions) (ClientResourceChangeDetailsResponse, error)`
- New function `PossibleChangeTypeValues() []ChangeType`
- New struct `ClientResourceChangeDetailsOptions`
- New struct `ClientResourceChangeDetailsResponse`
- New struct `ClientResourceChangesOptions`
- New struct `ClientResourceChangesResponse`
- New struct `ResourceChangeData`
- New struct `ResourceChangeDataAfterSnapshot`
- New struct `ResourceChangeDataBeforeSnapshot`
- New struct `ResourceChangeDetailsRequestParameters`
- New struct `ResourceChangeList`
- New struct `ResourceChangesRequestParameters`
- New struct `ResourceChangesRequestParametersInterval`
- New struct `ResourcePropertyChange`
- New struct `ResourceSnapshotData`


## 0.6.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
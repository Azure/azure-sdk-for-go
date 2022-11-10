# Release History

## 0.6.0 (2022-11-10)
### Features Added

- New const `TypeAzurePrivateEndpoint`
- New const `TypeIP`
- New type alias `Type`
- New function `*TrafficFiltersClient.Delete(context.Context, string, string, *TrafficFiltersClientDeleteOptions) (TrafficFiltersClientDeleteResponse, error)`
- New function `*AssociateTrafficFilterClient.BeginAssociate(context.Context, string, string, *AssociateTrafficFilterClientBeginAssociateOptions) (*runtime.Poller[AssociateTrafficFilterClientAssociateResponse], error)`
- New function `NewCreateAndAssociateIPFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CreateAndAssociateIPFilterClient, error)`
- New function `NewUpgradableVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UpgradableVersionsClient, error)`
- New function `NewTrafficFiltersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TrafficFiltersClient, error)`
- New function `NewDetachTrafficFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DetachTrafficFilterClient, error)`
- New function `NewAllTrafficFiltersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AllTrafficFiltersClient, error)`
- New function `NewDetachAndDeleteTrafficFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DetachAndDeleteTrafficFilterClient, error)`
- New function `*MonitorClient.BeginUpgrade(context.Context, string, string, *MonitorClientBeginUpgradeOptions) (*runtime.Poller[MonitorClientUpgradeResponse], error)`
- New function `*AllTrafficFiltersClient.List(context.Context, string, string, *AllTrafficFiltersClientListOptions) (AllTrafficFiltersClientListResponse, error)`
- New function `*ExternalUserClient.CreateOrUpdate(context.Context, string, string, *ExternalUserClientCreateOrUpdateOptions) (ExternalUserClientCreateOrUpdateResponse, error)`
- New function `NewListAssociatedTrafficFiltersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ListAssociatedTrafficFiltersClient, error)`
- New function `*CreateAndAssociateIPFilterClient.BeginCreate(context.Context, string, string, *CreateAndAssociateIPFilterClientBeginCreateOptions) (*runtime.Poller[CreateAndAssociateIPFilterClientCreateResponse], error)`
- New function `*DetachAndDeleteTrafficFilterClient.Delete(context.Context, string, string, *DetachAndDeleteTrafficFilterClientDeleteOptions) (DetachAndDeleteTrafficFilterClientDeleteResponse, error)`
- New function `NewAssociateTrafficFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AssociateTrafficFilterClient, error)`
- New function `*DetachTrafficFilterClient.BeginUpdate(context.Context, string, string, *DetachTrafficFilterClientBeginUpdateOptions) (*runtime.Poller[DetachTrafficFilterClientUpdateResponse], error)`
- New function `*CreateAndAssociatePLFilterClient.BeginCreate(context.Context, string, string, *CreateAndAssociatePLFilterClientBeginCreateOptions) (*runtime.Poller[CreateAndAssociatePLFilterClientCreateResponse], error)`
- New function `*ListAssociatedTrafficFiltersClient.List(context.Context, string, string, *ListAssociatedTrafficFiltersClientListOptions) (ListAssociatedTrafficFiltersClientListResponse, error)`
- New function `PossibleTypeValues() []Type`
- New function `NewCreateAndAssociatePLFilterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CreateAndAssociatePLFilterClient, error)`
- New function `NewMonitorClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MonitorClient, error)`
- New function `NewExternalUserClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExternalUserClient, error)`
- New function `*UpgradableVersionsClient.Details(context.Context, string, string, *UpgradableVersionsClientDetailsOptions) (UpgradableVersionsClientDetailsResponse, error)`
- New struct `AllTrafficFiltersClient`
- New struct `AllTrafficFiltersClientListOptions`
- New struct `AllTrafficFiltersClientListResponse`
- New struct `AssociateTrafficFilterClient`
- New struct `AssociateTrafficFilterClientAssociateResponse`
- New struct `AssociateTrafficFilterClientBeginAssociateOptions`
- New struct `CreateAndAssociateIPFilterClient`
- New struct `CreateAndAssociateIPFilterClientBeginCreateOptions`
- New struct `CreateAndAssociateIPFilterClientCreateResponse`
- New struct `CreateAndAssociatePLFilterClient`
- New struct `CreateAndAssociatePLFilterClientBeginCreateOptions`
- New struct `CreateAndAssociatePLFilterClientCreateResponse`
- New struct `DetachAndDeleteTrafficFilterClient`
- New struct `DetachAndDeleteTrafficFilterClientDeleteOptions`
- New struct `DetachAndDeleteTrafficFilterClientDeleteResponse`
- New struct `DetachTrafficFilterClient`
- New struct `DetachTrafficFilterClientBeginUpdateOptions`
- New struct `DetachTrafficFilterClientUpdateResponse`
- New struct `ExternalUserClient`
- New struct `ExternalUserClientCreateOrUpdateOptions`
- New struct `ExternalUserClientCreateOrUpdateResponse`
- New struct `ExternalUserCreationResponse`
- New struct `ExternalUserInfo`
- New struct `ListAssociatedTrafficFiltersClient`
- New struct `ListAssociatedTrafficFiltersClientListOptions`
- New struct `ListAssociatedTrafficFiltersClientListResponse`
- New struct `MonitorClient`
- New struct `MonitorClientBeginUpgradeOptions`
- New struct `MonitorClientUpgradeResponse`
- New struct `MonitorUpgrade`
- New struct `TrafficFilter`
- New struct `TrafficFilterResponse`
- New struct `TrafficFilterRule`
- New struct `TrafficFiltersClient`
- New struct `TrafficFiltersClientDeleteOptions`
- New struct `TrafficFiltersClientDeleteResponse`
- New struct `UpgradableVersionsClient`
- New struct `UpgradableVersionsClientDetailsOptions`
- New struct `UpgradableVersionsClientDetailsResponse`
- New struct `UpgradableVersionsList`
- New field `Version` in struct `MonitorProperties`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/elastic/armelastic` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
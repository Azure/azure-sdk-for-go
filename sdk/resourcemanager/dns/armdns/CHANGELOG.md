# Release History

## 1.3.0-beta.1 (2024-10-25)
### Features Added

- New value `RecordTypeDS`, `RecordTypeNAPTR`, `RecordTypeTLSA` added to enum type `RecordType`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New function `*ClientFactory.NewDnssecConfigsClient() *DnssecConfigsClient`
- New function `NewDnssecConfigsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DnssecConfigsClient, error)`
- New function `*DnssecConfigsClient.BeginCreateOrUpdate(context.Context, string, string, *DnssecConfigsClientBeginCreateOrUpdateOptions) (*runtime.Poller[DnssecConfigsClientCreateOrUpdateResponse], error)`
- New function `*DnssecConfigsClient.BeginDelete(context.Context, string, string, *DnssecConfigsClientBeginDeleteOptions) (*runtime.Poller[DnssecConfigsClientDeleteResponse], error)`
- New function `*DnssecConfigsClient.Get(context.Context, string, string, *DnssecConfigsClientGetOptions) (DnssecConfigsClientGetResponse, error)`
- New function `*DnssecConfigsClient.NewListByDNSZonePager(string, string, *DnssecConfigsClientListByDNSZoneOptions) *runtime.Pager[DnssecConfigsClientListByDNSZoneResponse]`
- New struct `DelegationSignerInfo`
- New struct `Digest`
- New struct `DnssecConfig`
- New struct `DnssecConfigListResult`
- New struct `DnssecProperties`
- New struct `DsRecord`
- New struct `NaptrRecord`
- New struct `SigningKey`
- New struct `SystemData`
- New struct `TlsaRecord`
- New field `DsRecords`, `NaptrRecords`, `TlsaRecords`, `TrafficManagementProfile` in struct `RecordSetProperties`
- New field `SystemData` in struct `Zone`
- New field `SigningKeys` in struct `ZoneProperties`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
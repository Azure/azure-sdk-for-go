# Unreleased

## Additive Changes

### New Funcs

1. *ConfigurationsBatchUpdateFuture.UnmarshalJSON([]byte) error
1. ConfigurationsClient.BatchUpdate(context.Context, string, string, ConfigurationListResult) (ConfigurationsBatchUpdateFuture, error)
1. ConfigurationsClient.BatchUpdatePreparer(context.Context, string, string, ConfigurationListResult) (*http.Request, error)
1. ConfigurationsClient.BatchUpdateResponder(*http.Response) (ConfigurationListResult, error)
1. ConfigurationsClient.BatchUpdateSender(*http.Request) (ConfigurationsBatchUpdateFuture, error)
1. StorageProfile.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ConfigurationsBatchUpdateFuture
1. PrivateDNSZoneArguments

#### New Struct Fields

1. CapabilityProperties.Status
1. ServerEditionCapability.Status
1. ServerProperties.PrivateDNSZoneArguments
1. ServerVersionCapability.Status
1. StorageEditionCapability.Status
1. StorageProfile.FileStorageSkuName
1. VcoreCapability.Status

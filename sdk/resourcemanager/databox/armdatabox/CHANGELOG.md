# Release History

## 2.2.0 (2025-03-27)
### Features Added

- New value `DataCenterCodeAMS25`, `DataCenterCodeBL24`, `DataCenterCodeCPQ21`, `DataCenterCodeDSM11`, `DataCenterCodeDXB23`, `DataCenterCodeIDC5`, `DataCenterCodeNTG20`, `DataCenterCodeOSA23`, `DataCenterCodeTYO23` added to enum type `DataCenterCode`
- New enum type `DelayNotificationStatus` with values `DelayNotificationStatusActive`, `DelayNotificationStatusResolved`
- New enum type `ModelName` with values `ModelNameAzureDataBox120`, `ModelNameAzureDataBox525`, `ModelNameDataBox`, `ModelNameDataBoxCustomerDisk`, `ModelNameDataBoxDisk`, `ModelNameDataBoxHeavy`
- New enum type `PortalDelayErrorCode` with values `PortalDelayErrorCodeActiveOrderLimitBreachedDelay`, `PortalDelayErrorCodeHighDemandDelay`, `PortalDelayErrorCodeInternalIssueDelay`, `PortalDelayErrorCodeLargeNumberOfFilesDelay`
- New struct `DeviceCapabilityDetails`
- New struct `DeviceCapabilityRequest`
- New struct `DeviceCapabilityResponse`
- New struct `JobDelayDetails`
- New field `Model` in struct `CommonScheduleAvailabilityRequest`
- New field `Model` in struct `CreateOrderLimitForSubscriptionValidationRequest`
- New field `Model` in struct `DataTransferDetailsValidationRequest`
- New field `Model` in struct `DatacenterAddressRequest`
- New field `Model` in struct `DiskScheduleAvailabilityRequest`
- New field `Model` in struct `HeavyScheduleAvailabilityRequest`
- New field `AllDevicesLost`, `DelayedStage` in struct `JobProperties`
- New field `DelayInformation` in struct `JobStages`
- New field `Model` in struct `PreferencesValidationRequest`
- New field `DeviceCapabilityRequest` in struct `RegionConfigurationRequest`
- New field `DeviceCapabilityResponse` in struct `RegionConfigurationResponse`
- New field `Model` in struct `SKU`
- New field `Model` in struct `SKUAvailabilityValidationRequest`
- New field `IndividualSKUUsable` in struct `SKUCapacity`
- New field `Model` in struct `ScheduleAvailabilityRequest`
- New field `Model` in struct `TransportAvailabilityRequest`
- New field `Model` in struct `ValidateAddress`


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-05-26)
### Breaking Changes

- Field `AccountName` of struct `DiskGranularCopyLogDetails` has been removed

### Features Added

- New enum type `HardwareEncryption` with values `HardwareEncryptionDisabled`, `HardwareEncryptionEnabled`
- New enum type `ReverseShippingDetailsEditStatus` with values `ReverseShippingDetailsEditStatusDisabled`, `ReverseShippingDetailsEditStatusEnabled`, `ReverseShippingDetailsEditStatusNotSupported`
- New enum type `ReverseTransportPreferenceEditStatus` with values `ReverseTransportPreferenceEditStatusDisabled`, `ReverseTransportPreferenceEditStatusEnabled`, `ReverseTransportPreferenceEditStatusNotSupported`
- New struct `ContactInfo`
- New struct `ReverseShippingDetails`
- New field `ReverseShippingDetails` in struct `CommonJobDetails`
- New field `Actions`, `Error` in struct `CopyProgress`
- New field `Actions`, `Error` in struct `CustomerDiskCopyProgress`
- New field `ReverseShippingDetails` in struct `CustomerDiskJobDetails`
- New field `Actions`, `Error` in struct `DiskCopyProgress`
- New field `AccountID` in struct `DiskGranularCopyLogDetails`
- New field `Actions`, `Error` in struct `DiskGranularCopyProgress`
- New field `GranularCopyLogDetails`, `ReverseShippingDetails` in struct `DiskJobDetails`
- New field `HardwareEncryption` in struct `EncryptionPreferences`
- New field `Actions`, `Error` in struct `GranularCopyProgress`
- New field `ReverseShippingDetails` in struct `HeavyJobDetails`
- New field `ReverseShippingDetails` in struct `JobDetails`
- New field `ReverseShippingDetailsUpdate`, `ReverseTransportPreferenceUpdate` in struct `JobProperties`
- New field `SerialNumberCustomerResolutionMap` in struct `MitigateJobRequest`
- New field `ReverseTransportPreferences` in struct `Preferences`
- New field `CountriesWithinCommerceBoundary` in struct `SKUProperties`
- New field `SkipAddressValidation`, `TaxIdentificationNumber` in struct `ShippingAddress`
- New field `IsUpdated` in struct `TransportPreferences`
- New field `Preferences`, `ReverseShippingDetails` in struct `UpdateJobDetails`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databox/armdatabox` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
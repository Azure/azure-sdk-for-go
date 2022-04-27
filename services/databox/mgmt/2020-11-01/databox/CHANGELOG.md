# Unreleased

## Breaking Changes

### Removed Constants

1. ShareDestinationFormatType.ShareDestinationFormatTypeAzurePremiumFiles

### Struct Changes

#### Removed Struct Fields

1. DcAccessSecurityCode.ForwardDcAccessCode
1. DcAccessSecurityCode.ReverseDcAccessCode
1. DiskJobDetails.ExpectedDataSizeInTerabytes
1. DiskScheduleAvailabilityRequest.ExpectedDataSizeInTerabytes
1. HeavyJobDetails.ExpectedDataSizeInTerabytes
1. JobDetails.ExpectedDataSizeInTerabytes
1. JobDetailsType.ExpectedDataSizeInTerabytes

## Additive Changes

### New Funcs

1. SystemData.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. SystemData

#### New Struct Fields

1. DcAccessSecurityCode.ForwardDCAccessCode
1. DcAccessSecurityCode.ReverseDCAccessCode
1. DiskJobDetails.ExpectedDataSizeInTeraBytes
1. DiskScheduleAvailabilityRequest.ExpectedDataSizeInTeraBytes
1. HeavyJobDetails.ExpectedDataSizeInTeraBytes
1. JobDetails.ExpectedDataSizeInTeraBytes
1. JobDetailsType.ExpectedDataSizeInTeraBytes
1. JobResource.SystemData

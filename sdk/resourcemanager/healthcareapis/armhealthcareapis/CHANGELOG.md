# Release History

## 1.1.0 (2022-11-11)
### Features Added

- New struct `CorsConfiguration`
- New struct `FhirServiceImportConfiguration`
- New struct `ServiceImportConfigurationInfo`
- New field `ImportConfiguration` in struct `FhirServiceProperties`
- New field `ResourceIDDimensionNameOverride` in struct `MetricSpecification`
- New field `EnableRegionalMdmAccount` in struct `MetricSpecification`
- New field `MetricFilterPattern` in struct `MetricSpecification`
- New field `IsInternal` in struct `MetricSpecification`
- New field `SourceMdmAccount` in struct `MetricSpecification`
- New field `CorsConfiguration` in struct `DicomServiceProperties`
- New field `ImportConfiguration` in struct `ServicesProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/healthcareapis/armhealthcareapis` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
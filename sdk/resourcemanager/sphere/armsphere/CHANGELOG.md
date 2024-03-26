# Release History

## 1.0.0 (2024-03-26)
### Breaking Changes

- Struct `CountDeviceResponse` has been renamed `CountDevicesResponse`

### Features Added

- New function `*CatalogsClient.BeginUploadImage(context.Context, string, string, Image, *CatalogsClientBeginUploadImageOptions) (*runtime.Poller[CatalogsClientUploadImageResponse], error)`
- New field `TenantID` in struct `CatalogProperties`


## 0.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.1.0 (2023-07-28)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sphere/armsphere` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
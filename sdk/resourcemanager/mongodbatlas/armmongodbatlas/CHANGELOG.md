# Release History

## 1.0.0 (2025-07-02)
### Breaking Changes

- Function `*OrganizationsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, OrganizationResource, *OrganizationsClientBeginUpdateOptions)` to `(context.Context, string, string, OrganizationResourceUpdate, *OrganizationsClientBeginUpdateOptions)`

### Features Added

- New struct `OrganizationResourceUpdate`
- New struct `OrganizationResourceUpdateProperties`


## 0.1.0 (2025-05-07)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mongodbatlas/armmongodbatlas` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
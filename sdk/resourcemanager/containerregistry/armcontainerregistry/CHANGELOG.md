# Release History

## 0.1.1 (2021-10-26)
### Breaking Changes

### New Content

- New const `ConnectedRegistryModeReadOnly`
- New const `ConnectedRegistryModeReadWrite`
- New field `NotificationsList` in struct `ConnectedRegistryUpdateProperties`
- New field `NotificationsList` in struct `ConnectedRegistryProperties`

Total 0 breaking change(s), 4 additive change(s).


## 0.1.0 (2021-10-08)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/containerregistry/armcontainerregistry") to avoid confusion.
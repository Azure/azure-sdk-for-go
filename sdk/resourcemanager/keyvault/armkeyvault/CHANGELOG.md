# Release History

## 1.1.0-beta.1 (2022-05-19)
### Features Added

- New const `KeyRotationPolicyActionTypeNotify`
- New const `JSONWebKeyOperationRelease`
- New const `KeyRotationPolicyActionTypeRotate`
- New const `KeyPermissionsRotate`
- New const `KeyPermissionsRelease`
- New const `KeyPermissionsSetrotationpolicy`
- New const `KeyPermissionsGetrotationpolicy`
- New function `PossibleKeyRotationPolicyActionTypeValues() []KeyRotationPolicyActionType`
- New function `*KeyReleasePolicy.UnmarshalJSON([]byte) error`
- New function `KeyReleasePolicy.MarshalJSON() ([]byte, error)`
- New function `RotationPolicy.MarshalJSON() ([]byte, error)`
- New struct `Action`
- New struct `KeyReleasePolicy`
- New struct `KeyRotationPolicyAttributes`
- New struct `LifetimeAction`
- New struct `RotationPolicy`
- New struct `Trigger`
- New field `ReleasePolicy` in struct `KeyProperties`
- New field `RotationPolicy` in struct `KeyProperties`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
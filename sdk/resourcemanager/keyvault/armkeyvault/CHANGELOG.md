# Release History

## 1.1.0 (2022-09-02)
### Features Added

- New const `KeyPermissionsRelease`
- New const `KeyRotationPolicyActionTypeNotify`
- New const `KeyPermissionsSetrotationpolicy`
- New const `ActivationStatusFailed`
- New const `KeyPermissionsGetrotationpolicy`
- New const `ActivationStatusNotActivated`
- New const `JSONWebKeyOperationRelease`
- New const `ActivationStatusActive`
- New const `KeyPermissionsRotate`
- New const `ActivationStatusUnknown`
- New const `KeyRotationPolicyActionTypeRotate`
- New type alias `KeyRotationPolicyActionType`
- New type alias `ActivationStatus`
- New function `*ManagedHsmsClient.CheckMhsmNameAvailability(context.Context, CheckMhsmNameAvailabilityParameters, *ManagedHsmsClientCheckMhsmNameAvailabilityOptions) (ManagedHsmsClientCheckMhsmNameAvailabilityResponse, error)`
- New function `PossibleKeyRotationPolicyActionTypeValues() []KeyRotationPolicyActionType`
- New function `PossibleActivationStatusValues() []ActivationStatus`
- New struct `Action`
- New struct `CheckMhsmNameAvailabilityParameters`
- New struct `CheckMhsmNameAvailabilityResult`
- New struct `KeyReleasePolicy`
- New struct `KeyRotationPolicyAttributes`
- New struct `LifetimeAction`
- New struct `ManagedHSMSecurityDomainProperties`
- New struct `ManagedHsmsClientCheckMhsmNameAvailabilityOptions`
- New struct `ManagedHsmsClientCheckMhsmNameAvailabilityResponse`
- New struct `RotationPolicy`
- New struct `Trigger`
- New field `ReleasePolicy` in struct `KeyProperties`
- New field `RotationPolicy` in struct `KeyProperties`
- New field `Etag` in struct `MHSMPrivateEndpointConnectionItem`
- New field `ID` in struct `MHSMPrivateEndpointConnectionItem`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
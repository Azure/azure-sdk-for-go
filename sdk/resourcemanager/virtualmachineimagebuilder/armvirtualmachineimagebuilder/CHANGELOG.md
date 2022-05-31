# Release History

## 1.1.0 (2022-05-30)
### Features Added

- New const `ProvisioningErrorCodeBadValidatorType`
- New const `ProvisioningErrorCodeUnsupportedValidatorType`
- New const `RunSubStateValidating`
- New const `ProvisioningErrorCodeNoValidatorScript`
- New const `ProvisioningErrorCodeBadStagingResourceGroup`
- New function `ImageTemplatePowerShellValidator.MarshalJSON() ([]byte, error)`
- New function `*ImageTemplatePropertiesValidate.UnmarshalJSON([]byte) error`
- New function `*ImageTemplateShellValidator.UnmarshalJSON([]byte) error`
- New function `*ImageTemplatePowerShellValidator.GetImageTemplateInVMValidator() *ImageTemplateInVMValidator`
- New function `ImageTemplateShellValidator.MarshalJSON() ([]byte, error)`
- New function `*ImageTemplateShellValidator.GetImageTemplateInVMValidator() *ImageTemplateInVMValidator`
- New function `ImageTemplatePropertiesValidate.MarshalJSON() ([]byte, error)`
- New function `*ImageTemplateInVMValidator.GetImageTemplateInVMValidator() *ImageTemplateInVMValidator`
- New function `*ImageTemplatePowerShellValidator.UnmarshalJSON([]byte) error`
- New struct `ImageTemplateInVMValidator`
- New struct `ImageTemplatePowerShellValidator`
- New struct `ImageTemplatePropertiesValidate`
- New struct `ImageTemplateShellValidator`
- New struct `ProxyResource`
- New field `SystemData` in struct `TrackedResource`
- New field `SystemData` in struct `RunOutput`
- New field `ExactStagingResourceGroup` in struct `ImageTemplateProperties`
- New field `StagingResourceGroup` in struct `ImageTemplateProperties`
- New field `Validate` in struct `ImageTemplateProperties`
- New field `SystemData` in struct `Resource`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/virtualmachineimagebuilder/armvirtualmachineimagebuilder` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
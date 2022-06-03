# Release History

## 1.1.0-beta.1 (2022-05-19)
### Features Added

- New const `TaskScopeResource`
- New const `RebootOptionsNever`
- New const `RebootOptionsAlways`
- New const `TaskScopeGlobal`
- New const `RebootOptionsIfRequired`
- New function `InputLinuxParameters.MarshalJSON() ([]byte, error)`
- New function `*ConfigurationAssignmentsClient.Get(context.Context, string, string, string, string, string, *ConfigurationAssignmentsClientGetOptions) (ConfigurationAssignmentsClientGetResponse, error)`
- New function `InputWindowsParameters.MarshalJSON() ([]byte, error)`
- New function `TaskProperties.MarshalJSON() ([]byte, error)`
- New function `PossibleRebootOptionsValues() []RebootOptions`
- New function `SoftwareUpdateConfigurationTasks.MarshalJSON() ([]byte, error)`
- New function `PossibleTaskScopeValues() []TaskScope`
- New function `*ConfigurationAssignmentsClient.GetParent(context.Context, string, string, string, string, string, string, string, *ConfigurationAssignmentsClientGetParentOptions) (ConfigurationAssignmentsClientGetParentResponse, error)`
- New struct `ConfigurationAssignmentsClientGetOptions`
- New struct `ConfigurationAssignmentsClientGetParentOptions`
- New struct `ConfigurationAssignmentsClientGetParentResponse`
- New struct `ConfigurationAssignmentsClientGetResponse`
- New struct `ConfigurationAssignmentsWithinSubscriptionClientListOptions`
- New struct `ConfigurationAssignmentsWithinSubscriptionClientListResponse`
- New struct `InputLinuxParameters`
- New struct `InputPatchConfiguration`
- New struct `InputWindowsParameters`
- New struct `SoftwareUpdateConfigurationTasks`
- New struct `TaskProperties`
- New field `InstallPatches` in struct `ConfigurationProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/maintenance/armmaintenance` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
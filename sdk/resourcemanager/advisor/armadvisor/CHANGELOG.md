# Release History

## 1.1.0 (2022-09-14)
### Features Added

- New const `DurationSixty`
- New const `DurationNinety`
- New const `DurationFourteen`
- New const `DurationSeven`
- New const `DurationThirty`
- New const `DurationTwentyOne`
- New const `PredictionTypePredictiveRightsizing`
- New type alias `PredictionType`
- New type alias `Duration`
- New function `NewManagementClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagementClient, error)`
- New function `PossiblePredictionTypeValues() []PredictionType`
- New function `*ManagementClient.Predict(context.Context, PredictionRequest, *ManagementClientPredictOptions) (ManagementClientPredictResponse, error)`
- New function `PossibleDurationValues() []Duration`
- New struct `ManagementClient`
- New struct `ManagementClientPredictOptions`
- New struct `ManagementClientPredictResponse`
- New struct `PredictionRequest`
- New struct `PredictionRequestProperties`
- New struct `PredictionResponse`
- New struct `PredictionResponseProperties`
- New field `Duration` in struct `ConfigDataProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/advisor/armadvisor` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
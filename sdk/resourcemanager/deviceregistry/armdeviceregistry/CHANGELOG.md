# Release History

## 1.0.0 (2024-12-09)
### Breaking Changes

- Type of `AssetProperties.Version` has been changed from `*int32` to `*int64`
- Type of `AssetStatus.Version` has been changed from `*int32` to `*int64`
- Type of `DataPoint.ObservabilityMode` has been changed from `*DataPointsObservabilityMode` to `*DataPointObservabilityMode`
- Type of `Event.ObservabilityMode` has been changed from `*EventsObservabilityMode` to `*EventObservabilityMode`
- Enum `DataPointsObservabilityMode` has been removed
- Enum `EventsObservabilityMode` has been removed
- Enum `UserAuthenticationMode` has been removed
- Struct `OwnCertificate` has been removed
- Struct `TransportAuthentication` has been removed
- Struct `TransportAuthenticationUpdate` has been removed
- Struct `UserAuthentication` has been removed
- Struct `UserAuthenticationUpdate` has been removed
- Field `TransportAuthentication`, `UserAuthentication` of struct `AssetEndpointProfileProperties` has been removed
- Field `TransportAuthentication`, `UserAuthentication` of struct `AssetEndpointProfileUpdateProperties` has been removed
- Field `AssetEndpointProfileURI`, `AssetType`, `DataPoints`, `DefaultDataPointsConfiguration` of struct `AssetProperties` has been removed
- Field `AssetType`, `DataPoints`, `DefaultDataPointsConfiguration` of struct `AssetUpdateProperties` has been removed
- Field `CapabilityID` of struct `DataPoint` has been removed
- Field `CapabilityID` of struct `Event` has been removed
- Field `PasswordReference`, `UsernameReference` of struct `UsernamePasswordCredentials` has been removed
- Field `PasswordReference`, `UsernameReference` of struct `UsernamePasswordCredentialsUpdate` has been removed
- Field `CertificateReference` of struct `X509Credentials` has been removed
- Field `CertificateReference` of struct `X509CredentialsUpdate` has been removed

### Features Added

- New value `ProvisioningStateDeleting` added to enum type `ProvisioningState`
- New enum type `AuthenticationMethod` with values `AuthenticationMethodAnonymous`, `AuthenticationMethodCertificate`, `AuthenticationMethodUsernamePassword`
- New enum type `DataPointObservabilityMode` with values `DataPointObservabilityModeCounter`, `DataPointObservabilityModeGauge`, `DataPointObservabilityModeHistogram`, `DataPointObservabilityModeLog`, `DataPointObservabilityModeNone`
- New enum type `EventObservabilityMode` with values `EventObservabilityModeLog`, `EventObservabilityModeNone`
- New enum type `TopicRetainType` with values `TopicRetainTypeKeep`, `TopicRetainTypeNever`
- New function `NewBillingContainersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BillingContainersClient, error)`
- New function `*BillingContainersClient.Get(context.Context, string, *BillingContainersClientGetOptions) (BillingContainersClientGetResponse, error)`
- New function `*BillingContainersClient.NewListBySubscriptionPager(*BillingContainersClientListBySubscriptionOptions) *runtime.Pager[BillingContainersClientListBySubscriptionResponse]`
- New function `*ClientFactory.NewBillingContainersClient() *BillingContainersClient`
- New struct `AssetEndpointProfileStatus`
- New struct `AssetEndpointProfileStatusError`
- New struct `AssetStatusDataset`
- New struct `AssetStatusEvent`
- New struct `Authentication`
- New struct `AuthenticationUpdate`
- New struct `BillingContainer`
- New struct `BillingContainerListResult`
- New struct `BillingContainerProperties`
- New struct `Dataset`
- New struct `MessageSchemaReference`
- New struct `Topic`
- New struct `TopicUpdate`
- New field `Authentication`, `DiscoveredAssetEndpointProfileRef`, `EndpointProfileType`, `Status` in struct `AssetEndpointProfileProperties`
- New field `Authentication`, `EndpointProfileType` in struct `AssetEndpointProfileUpdateProperties`
- New field `AssetEndpointProfileRef`, `Datasets`, `DefaultDatasetsConfiguration`, `DefaultTopic`, `DiscoveredAssetRefs` in struct `AssetProperties`
- New field `Datasets`, `Events` in struct `AssetStatus`
- New field `Datasets`, `DefaultDatasetsConfiguration`, `DefaultTopic` in struct `AssetUpdateProperties`
- New field `Topic` in struct `Event`
- New field `ResourceID` in struct `OperationStatusResult`
- New field `PasswordSecretName`, `UsernameSecretName` in struct `UsernamePasswordCredentials`
- New field `PasswordSecretName`, `UsernameSecretName` in struct `UsernamePasswordCredentialsUpdate`
- New field `CertificateSecretName` in struct `X509Credentials`
- New field `CertificateSecretName` in struct `X509CredentialsUpdate`


## 0.1.0 (2024-04-26)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deviceregistry/armdeviceregistry` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
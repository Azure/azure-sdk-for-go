# Release History

## 2.0.0-beta.2 (2022-08-14)
### Other Changes
- Replace API version `2021-10-01-preview` with `2022-07-01-preview`

## 2.0.0-beta.1 (2022-06-09)
### Breaking Changes

- Type of `ServiceProperties.ProvisioningState` has been changed from `*ProvisioningState` to `*CommunicationServicesProvisioningState`
- Const `ProvisioningStateCreating` has been removed
- Const `ProvisioningStateSucceeded` has been removed
- Const `ProvisioningStateFailed` has been removed
- Const `ProvisioningStateUpdating` has been removed
- Const `ProvisioningStateUnknown` has been removed
- Const `ProvisioningStateRunning` has been removed
- Const `ProvisioningStateDeleting` has been removed
- Const `ProvisioningStateMoving` has been removed
- Const `ProvisioningStateCanceled` has been removed
- Function `*ServiceClient.RegenerateKey` has been removed
- Function `*ServiceClient.Update` has been removed
- Function `*ServiceClient.NewListByResourceGroupPager` has been removed
- Function `*ServiceClient.NewListBySubscriptionPager` has been removed
- Function `*ServiceClient.ListKeys` has been removed
- Function `*ServiceClient.Get` has been removed
- Function `PossibleProvisioningStateValues` has been removed
- Function `*ServiceClient.BeginCreateOrUpdate` has been removed
- Function `*ServiceClient.CheckNameAvailability` has been removed
- Function `*ServiceClient.BeginDelete` has been removed
- Function `NewServiceClient` has been removed
- Function `*ServiceClient.LinkNotificationHub` has been removed
- Struct `LocationResource` has been removed
- Struct `NameAvailability` has been removed
- Struct `ServiceClient` has been removed
- Struct `ServiceClientBeginCreateOrUpdateOptions` has been removed
- Struct `ServiceClientBeginDeleteOptions` has been removed
- Struct `ServiceClientCheckNameAvailabilityOptions` has been removed
- Struct `ServiceClientCheckNameAvailabilityResponse` has been removed
- Struct `ServiceClientCreateOrUpdateResponse` has been removed
- Struct `ServiceClientDeleteResponse` has been removed
- Struct `ServiceClientGetOptions` has been removed
- Struct `ServiceClientGetResponse` has been removed
- Struct `ServiceClientLinkNotificationHubOptions` has been removed
- Struct `ServiceClientLinkNotificationHubResponse` has been removed
- Struct `ServiceClientListByResourceGroupOptions` has been removed
- Struct `ServiceClientListByResourceGroupResponse` has been removed
- Struct `ServiceClientListBySubscriptionOptions` has been removed
- Struct `ServiceClientListBySubscriptionResponse` has been removed
- Struct `ServiceClientListKeysOptions` has been removed
- Struct `ServiceClientListKeysResponse` has been removed
- Struct `ServiceClientRegenerateKeyOptions` has been removed
- Struct `ServiceClientRegenerateKeyResponse` has been removed
- Struct `ServiceClientUpdateOptions` has been removed
- Struct `ServiceClientUpdateResponse` has been removed

### Features Added

- New const `DomainsProvisioningStateUnknown`
- New const `CommunicationServicesProvisioningStateUnknown`
- New const `CommunicationServicesProvisioningStateMoving`
- New const `CommunicationServicesProvisioningStateCreating`
- New const `VerificationStatusVerificationRequested`
- New const `VerificationStatusVerificationFailed`
- New const `DomainsProvisioningStateRunning`
- New const `DomainsProvisioningStateCreating`
- New const `CheckNameAvailabilityReasonInvalid`
- New const `UserEngagementTrackingEnabled`
- New const `DomainManagementAzureManaged`
- New const `DomainsProvisioningStateDeleting`
- New const `DomainManagementCustomerManagedInExchangeOnline`
- New const `UserEngagementTrackingDisabled`
- New const `DomainsProvisioningStateSucceeded`
- New const `CommunicationServicesProvisioningStateDeleting`
- New const `EmailServicesProvisioningStateDeleting`
- New const `CommunicationServicesProvisioningStateFailed`
- New const `EmailServicesProvisioningStateMoving`
- New const `CommunicationServicesProvisioningStateSucceeded`
- New const `DomainsProvisioningStateFailed`
- New const `VerificationStatusNotStarted`
- New const `EmailServicesProvisioningStateFailed`
- New const `CommunicationServicesProvisioningStateCanceled`
- New const `VerificationTypeDomain`
- New const `CheckNameAvailabilityReasonAlreadyExists`
- New const `CommunicationServicesProvisioningStateUpdating`
- New const `VerificationStatusVerified`
- New const `DomainsProvisioningStateMoving`
- New const `VerificationTypeDKIM`
- New const `VerificationStatusCancellationRequested`
- New const `EmailServicesProvisioningStateCreating`
- New const `DomainManagementCustomerManaged`
- New const `VerificationTypeDMARC`
- New const `DomainsProvisioningStateCanceled`
- New const `EmailServicesProvisioningStateUnknown`
- New const `EmailServicesProvisioningStateCanceled`
- New const `EmailServicesProvisioningStateSucceeded`
- New const `DomainsProvisioningStateUpdating`
- New const `EmailServicesProvisioningStateUpdating`
- New const `VerificationTypeSPF`
- New const `VerificationTypeDKIM2`
- New const `EmailServicesProvisioningStateRunning`
- New const `VerificationStatusVerificationInProgress`
- New const `CommunicationServicesProvisioningStateRunning`
- New function `PossibleVerificationStatusValues() []VerificationStatus`
- New function `UpdateDomainProperties.MarshalJSON() ([]byte, error)`
- New function `UpdateDomainRequestParameters.MarshalJSON() ([]byte, error)`
- New function `TrackedResource.MarshalJSON() ([]byte, error)`
- New function `ServiceResourceUpdate.MarshalJSON() ([]byte, error)`
- New function `PossibleDomainsProvisioningStateValues() []DomainsProvisioningState`
- New function `DomainProperties.MarshalJSON() ([]byte, error)`
- New function `PossibleCheckNameAvailabilityReasonValues() []CheckNameAvailabilityReason`
- New function `PossibleCommunicationServicesProvisioningStateValues() []CommunicationServicesProvisioningState`
- New function `DomainResource.MarshalJSON() ([]byte, error)`
- New function `PossibleEmailServicesProvisioningStateValues() []EmailServicesProvisioningState`
- New function `EmailServiceResource.MarshalJSON() ([]byte, error)`
- New function `PossibleDomainManagementValues() []DomainManagement`
- New function `EmailServiceResourceUpdate.MarshalJSON() ([]byte, error)`
- New function `PossibleUserEngagementTrackingValues() []UserEngagementTracking`
- New function `ServiceUpdateProperties.MarshalJSON() ([]byte, error)`
- New function `PossibleVerificationTypeValues() []VerificationType`
- New function `ServiceProperties.MarshalJSON() ([]byte, error)`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityResponse`
- New struct `DNSRecord`
- New struct `DomainProperties`
- New struct `DomainPropertiesVerificationRecords`
- New struct `DomainPropertiesVerificationStates`
- New struct `DomainResource`
- New struct `DomainResourceList`
- New struct `DomainsClientBeginCancelVerificationOptions`
- New struct `DomainsClientBeginCreateOrUpdateOptions`
- New struct `DomainsClientBeginDeleteOptions`
- New struct `DomainsClientBeginInitiateVerificationOptions`
- New struct `DomainsClientBeginUpdateOptions`
- New struct `DomainsClientCancelVerificationResponse`
- New struct `DomainsClientCreateOrUpdateResponse`
- New struct `DomainsClientDeleteResponse`
- New struct `DomainsClientGetOptions`
- New struct `DomainsClientGetResponse`
- New struct `DomainsClientInitiateVerificationResponse`
- New struct `DomainsClientListByEmailServiceResourceOptions`
- New struct `DomainsClientListByEmailServiceResourceResponse`
- New struct `DomainsClientUpdateResponse`
- New struct `EmailServiceProperties`
- New struct `EmailServiceResource`
- New struct `EmailServiceResourceList`
- New struct `EmailServiceResourceUpdate`
- New struct `EmailServicesClientBeginCreateOrUpdateOptions`
- New struct `EmailServicesClientBeginDeleteOptions`
- New struct `EmailServicesClientBeginUpdateOptions`
- New struct `EmailServicesClientCreateOrUpdateResponse`
- New struct `EmailServicesClientDeleteResponse`
- New struct `EmailServicesClientGetOptions`
- New struct `EmailServicesClientGetResponse`
- New struct `EmailServicesClientListByResourceGroupOptions`
- New struct `EmailServicesClientListByResourceGroupResponse`
- New struct `EmailServicesClientListBySubscriptionOptions`
- New struct `EmailServicesClientListBySubscriptionResponse`
- New struct `EmailServicesClientListVerifiedExchangeOnlineDomainsOptions`
- New struct `EmailServicesClientListVerifiedExchangeOnlineDomainsResponse`
- New struct `EmailServicesClientUpdateResponse`
- New struct `ServiceResourceUpdate`
- New struct `ServiceUpdateProperties`
- New struct `ServicesClientBeginCreateOrUpdateOptions`
- New struct `ServicesClientBeginDeleteOptions`
- New struct `ServicesClientBeginRegenerateKeyOptions`
- New struct `ServicesClientBeginUpdateOptions`
- New struct `ServicesClientCheckNameAvailabilityOptions`
- New struct `ServicesClientCheckNameAvailabilityResponse`
- New struct `ServicesClientCreateOrUpdateResponse`
- New struct `ServicesClientDeleteResponse`
- New struct `ServicesClientGetOptions`
- New struct `ServicesClientGetResponse`
- New struct `ServicesClientLinkNotificationHubOptions`
- New struct `ServicesClientLinkNotificationHubResponse`
- New struct `ServicesClientListByResourceGroupOptions`
- New struct `ServicesClientListByResourceGroupResponse`
- New struct `ServicesClientListBySubscriptionOptions`
- New struct `ServicesClientListBySubscriptionResponse`
- New struct `ServicesClientListKeysOptions`
- New struct `ServicesClientListKeysResponse`
- New struct `ServicesClientRegenerateKeyResponse`
- New struct `ServicesClientUpdateResponse`
- New struct `TrackedResource`
- New struct `UpdateDomainProperties`
- New struct `UpdateDomainRequestParameters`
- New struct `VerificationParameter`
- New struct `VerificationStatusRecord`
- New field `SystemData` in struct `Resource`
- New field `LinkedDomains` in struct `ServiceProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/communication/armcommunication` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).

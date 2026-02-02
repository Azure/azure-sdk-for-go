# Release History

## 2.0.0 (2026-02-02)
### Breaking Changes

- Function `*NetworkExperimentProfilesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, profileName string, resourceGroupName string, parameters Profile, options *NetworkExperimentProfilesClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, profileName string, parameters Profile, options *NetworkExperimentProfilesClientBeginCreateOrUpdateOptions)`
- Type of `BackendPoolProperties.ResourceState` has been changed from `*FrontDoorResourceState` to `*ResourceState`
- Type of `CacheConfiguration.QueryParameterStripDirective` has been changed from `*FrontDoorQuery` to `*Query`
- Type of `CertificateSourceParameters.CertificateType` has been changed from `*FrontDoorCertificateType` to `*CertificateType`
- Type of `CustomHTTPSConfiguration.CertificateSource` has been changed from `*FrontDoorCertificateSource` to `*CertificateSource`
- Type of `CustomHTTPSConfiguration.ProtocolType` has been changed from `*FrontDoorTLSProtocolType` to `*TLSProtocolType`
- Type of `ForwardingConfiguration.ForwardingProtocol` has been changed from `*FrontDoorForwardingProtocol` to `*ForwardingProtocol`
- Type of `FrontendEndpointProperties.ResourceState` has been changed from `*FrontDoorResourceState` to `*ResourceState`
- Type of `HealthProbeSettingsProperties.HealthProbeMethod` has been changed from `*FrontDoorHealthProbeMethod` to `*HealthProbeMethod`
- Type of `HealthProbeSettingsProperties.Protocol` has been changed from `*FrontDoorProtocol` to `*Protocol`
- Type of `HealthProbeSettingsProperties.ResourceState` has been changed from `*FrontDoorResourceState` to `*ResourceState`
- Type of `LoadBalancingSettingsProperties.ResourceState` has been changed from `*FrontDoorResourceState` to `*ResourceState`
- Type of `Properties.EnabledState` has been changed from `*FrontDoorEnabledState` to `*EnabledState`
- Type of `Properties.ResourceState` has been changed from `*FrontDoorResourceState` to `*ResourceState`
- Type of `RedirectConfiguration.RedirectProtocol` has been changed from `*FrontDoorRedirectProtocol` to `*RedirectProtocol`
- Type of `RedirectConfiguration.RedirectType` has been changed from `*FrontDoorRedirectType` to `*RedirectType`
- Type of `RoutingRuleProperties.AcceptedProtocols` has been changed from `[]*FrontDoorProtocol` to `[]*Protocol`
- Type of `RoutingRuleProperties.ResourceState` has been changed from `*FrontDoorResourceState` to `*ResourceState`
- Type of `RulesEngineProperties.ResourceState` has been changed from `*FrontDoorResourceState` to `*ResourceState`
- `MinimumTLSVersionOne0`, `MinimumTLSVersionOne2` from enum `MinimumTLSVersion` has been removed
- Enum `FrontDoorCertificateSource` has been removed
- Enum `FrontDoorCertificateType` has been removed
- Enum `FrontDoorEnabledState` has been removed
- Enum `FrontDoorForwardingProtocol` has been removed
- Enum `FrontDoorHealthProbeMethod` has been removed
- Enum `FrontDoorProtocol` has been removed
- Enum `FrontDoorQuery` has been removed
- Enum `FrontDoorRedirectProtocol` has been removed
- Enum `FrontDoorRedirectType` has been removed
- Enum `FrontDoorResourceState` has been removed
- Enum `FrontDoorTLSProtocolType` has been removed
- Enum `NetworkOperationStatus` has been removed
- Struct `AzureAsyncOperationResult` has been removed
- Struct `BackendPoolListResult` has been removed
- Struct `BackendPoolUpdateParameters` has been removed
- Struct `DefaultErrorResponse` has been removed
- Struct `DefaultErrorResponseError` has been removed
- Struct `Error` has been removed
- Struct `ErrorDetails` has been removed
- Struct `ErrorResponse` has been removed
- Struct `FrontendEndpointUpdateParameters` has been removed
- Struct `HealthProbeSettingsListResult` has been removed
- Struct `HealthProbeSettingsUpdateParameters` has been removed
- Struct `LoadBalancingSettingsListResult` has been removed
- Struct `LoadBalancingSettingsUpdateParameters` has been removed
- Struct `Resource` has been removed
- Struct `RoutingRuleListResult` has been removed
- Struct `RoutingRuleUpdateParameters` has been removed
- Struct `RulesEngineUpdateParameters` has been removed
- Struct `UpdateParameters` has been removed

### Features Added

- New value `ActionTypeCAPTCHA` added to enum type `ActionType`
- New value `MinimumTLSVersion10`, `MinimumTLSVersion12` added to enum type `MinimumTLSVersion`
- New value `OperatorServiceTagMatch` added to enum type `Operator`
- New enum type `CertificateSource` with values `CertificateSourceAzureKeyVault`, `CertificateSourceFrontDoor`
- New enum type `CertificateType` with values `CertificateTypeDedicated`
- New enum type `EnabledState` with values `EnabledStateDisabled`, `EnabledStateEnabled`
- New enum type `ForwardingProtocol` with values `ForwardingProtocolHTTPOnly`, `ForwardingProtocolHTTPSOnly`, `ForwardingProtocolMatchRequest`
- New enum type `HealthProbeMethod` with values `HealthProbeMethodGET`, `HealthProbeMethodHEAD`
- New enum type `Protocol` with values `ProtocolHTTP`, `ProtocolHTTPS`
- New enum type `Query` with values `QueryStripAll`, `QueryStripAllExcept`, `QueryStripNone`, `QueryStripOnly`
- New enum type `RedirectProtocol` with values `RedirectProtocolHTTPOnly`, `RedirectProtocolHTTPSOnly`, `RedirectProtocolMatchRequest`
- New enum type `RedirectType` with values `RedirectTypeFound`, `RedirectTypeMoved`, `RedirectTypePermanentRedirect`, `RedirectTypeTemporaryRedirect`
- New enum type `ResourceState` with values `ResourceStateCreating`, `ResourceStateDeleting`, `ResourceStateDisabled`, `ResourceStateDisabling`, `ResourceStateEnabled`, `ResourceStateEnabling`, `ResourceStateMigrated`, `ResourceStateMigrating`
- New enum type `SensitivityType` with values `SensitivityTypeHigh`, `SensitivityTypeLow`, `SensitivityTypeMedium`
- New enum type `TLSProtocolType` with values `TLSProtocolTypeServerNameIndication`
- New field `DefaultSensitivity` in struct `ManagedRuleDefinition`
- New field `Sensitivity` in struct `ManagedRuleOverride`
- New field `CaptchaExpirationInMinutes` in struct `PolicySettings`


## 1.4.0 (2024-04-26)
### Features Added

- New value `ActionTypeJSChallenge` added to enum type `ActionType`
- New enum type `ScrubbingRuleEntryMatchOperator` with values `ScrubbingRuleEntryMatchOperatorEquals`, `ScrubbingRuleEntryMatchOperatorEqualsAny`
- New enum type `ScrubbingRuleEntryMatchVariable` with values `ScrubbingRuleEntryMatchVariableQueryStringArgNames`, `ScrubbingRuleEntryMatchVariableRequestBodyJSONArgNames`, `ScrubbingRuleEntryMatchVariableRequestBodyPostArgNames`, `ScrubbingRuleEntryMatchVariableRequestCookieNames`, `ScrubbingRuleEntryMatchVariableRequestHeaderNames`, `ScrubbingRuleEntryMatchVariableRequestIPAddress`, `ScrubbingRuleEntryMatchVariableRequestURI`
- New enum type `ScrubbingRuleEntryState` with values `ScrubbingRuleEntryStateDisabled`, `ScrubbingRuleEntryStateEnabled`
- New enum type `VariableName` with values `VariableNameGeoLocation`, `VariableNameNone`, `VariableNameSocketAddr`
- New enum type `WebApplicationFirewallScrubbingState` with values `WebApplicationFirewallScrubbingStateDisabled`, `WebApplicationFirewallScrubbingStateEnabled`
- New struct `GroupByVariable`
- New struct `PolicySettingsLogScrubbing`
- New struct `WebApplicationFirewallScrubbingRules`
- New field `GroupBy` in struct `CustomRule`
- New field `JavascriptChallengeExpirationInMinutes`, `LogScrubbing` in struct `PolicySettings`


## 1.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-05-26)
### Features Added

- New value `ActionTypeAnomalyScoring` added to enum type `ActionType`
- New value `FrontDoorResourceStateMigrated`, `FrontDoorResourceStateMigrating` added to enum type `FrontDoorResourceState`
- New function `*PoliciesClient.NewListBySubscriptionPager(*PoliciesClientListBySubscriptionOptions) *runtime.Pager[PoliciesClientListBySubscriptionResponse]`
- New function `*PoliciesClient.BeginUpdate(context.Context, string, string, TagsObject, *PoliciesClientBeginUpdateOptions) (*runtime.Poller[PoliciesClientUpdateResponse], error)`
- New struct `DefaultErrorResponse`
- New struct `DefaultErrorResponseError`
- New field `ExtendedProperties` in struct `Properties`


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
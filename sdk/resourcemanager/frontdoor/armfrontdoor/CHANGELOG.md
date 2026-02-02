# Release History

## 1.5.0 (2026-02-02)
### Features Added

- New value `ActionTypeCAPTCHA` added to enum type `ActionType`
- New value `OperatorServiceTagMatch` added to enum type `Operator`
- New enum type `SensitivityType` with values `SensitivityTypeHigh`, `SensitivityTypeLow`, `SensitivityTypeMedium`
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
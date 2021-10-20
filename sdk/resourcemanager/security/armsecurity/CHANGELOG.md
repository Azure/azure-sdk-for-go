# Release History

## 0.2.0 (2021-10-20)
### Breaking Changes

- Function `*AdvancedThreatProtectionClient.Get` parameter(s) have been changed from `(context.Context, string, Enum6, *AdvancedThreatProtectionGetOptions)` to `(context.Context, string, Enum7, *AdvancedThreatProtectionGetOptions)`
- Function `*AdvancedThreatProtectionClient.Create` parameter(s) have been changed from `(context.Context, string, Enum6, AdvancedThreatProtectionSetting, *AdvancedThreatProtectionCreateOptions)` to `(context.Context, string, Enum7, AdvancedThreatProtectionSetting, *AdvancedThreatProtectionCreateOptions)`
- Function `*CustomAssessmentAutomationsClient.Create` parameter(s) have been changed from `(context.Context, string, string, CustomAssessmentAutomation, *CustomAssessmentAutomationsCreateOptions)` to `(context.Context, string, string, CustomAssessmentAutomationRequest, *CustomAssessmentAutomationsCreateOptions)`
- Const `Enum6Current` has been removed
- Function `PossibleEnum6Values` has been removed
- Function `Enum6.ToPtr` has been removed

### New Content

- New const `Enum7Current`
- New function `CustomAssessmentAutomationRequest.MarshalJSON() ([]byte, error)`
- New function `*CustomAssessmentAutomationRequest.UnmarshalJSON([]byte) error`
- New function `Enum7.ToPtr() *Enum7`
- New function `PossibleEnum7Values() []Enum7`
- New struct `CustomAssessmentAutomationRequest`
- New struct `CustomAssessmentAutomationRequestProperties`
- New field `AssessmentKey` in struct `CustomAssessmentAutomationProperties`
- New field `SystemData` in struct `CustomAssessmentAutomation`
- New field `SystemData` in struct `CustomEntityStoreAssignment`

Total 6 breaking change(s), 12 additive change(s).


## 0.1.0 (2021-10-15)

- Initial preview release.

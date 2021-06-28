# Unreleased

## Additive Changes

### New Funcs

1. *AlertRulePatchObject.UnmarshalJSON([]byte) error
1. AlertRulePatchObject.MarshalJSON() ([]byte, error)
1. SmartDetectorAlertRulesClient.Patch(context.Context, string, string, AlertRulePatchObject) (AlertRule, error)
1. SmartDetectorAlertRulesClient.PatchPreparer(context.Context, string, string, AlertRulePatchObject) (*http.Request, error)
1. SmartDetectorAlertRulesClient.PatchResponder(*http.Response) (AlertRule, error)
1. SmartDetectorAlertRulesClient.PatchSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. AlertRulePatchObject
1. AlertRulePatchProperties

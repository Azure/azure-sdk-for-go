# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. AlertsClient.ChangeState
	- Params
		- From: context.Context, string, AlertState
		- To: context.Context, string, AlertState, string
1. AlertsClient.ChangeStatePreparer
	- Params
		- From: context.Context, string, AlertState
		- To: context.Context, string, AlertState, string

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

1. ActionStatus
1. AlertRulePatchObject
1. AlertRulePatchProperties

#### New Struct Fields

1. Essentials.ActionStatus
1. Essentials.Description

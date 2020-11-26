
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewSBQueueListResultPage` signature has been changed from `(func(context.Context, SBQueueListResult) (SBQueueListResult, error))` to `(SBQueueListResult,func(context.Context, SBQueueListResult) (SBQueueListResult, error))`
- Function `NewSBNamespaceListResultPage` signature has been changed from `(func(context.Context, SBNamespaceListResult) (SBNamespaceListResult, error))` to `(SBNamespaceListResult,func(context.Context, SBNamespaceListResult) (SBNamespaceListResult, error))`
- Function `NewPremiumMessagingRegionsListResultPage` signature has been changed from `(func(context.Context, PremiumMessagingRegionsListResult) (PremiumMessagingRegionsListResult, error))` to `(PremiumMessagingRegionsListResult,func(context.Context, PremiumMessagingRegionsListResult) (PremiumMessagingRegionsListResult, error))`
- Function `NewSBSubscriptionListResultPage` signature has been changed from `(func(context.Context, SBSubscriptionListResult) (SBSubscriptionListResult, error))` to `(SBSubscriptionListResult,func(context.Context, SBSubscriptionListResult) (SBSubscriptionListResult, error))`
- Function `NewRuleListResultPage` signature has been changed from `(func(context.Context, RuleListResult) (RuleListResult, error))` to `(RuleListResult,func(context.Context, RuleListResult) (RuleListResult, error))`
- Function `NewSBAuthorizationRuleListResultPage` signature has been changed from `(func(context.Context, SBAuthorizationRuleListResult) (SBAuthorizationRuleListResult, error))` to `(SBAuthorizationRuleListResult,func(context.Context, SBAuthorizationRuleListResult) (SBAuthorizationRuleListResult, error))`
- Function `NewEventHubListResultPage` signature has been changed from `(func(context.Context, EventHubListResult) (EventHubListResult, error))` to `(EventHubListResult,func(context.Context, EventHubListResult) (EventHubListResult, error))`
- Function `NewArmDisasterRecoveryListResultPage` signature has been changed from `(func(context.Context, ArmDisasterRecoveryListResult) (ArmDisasterRecoveryListResult, error))` to `(ArmDisasterRecoveryListResult,func(context.Context, ArmDisasterRecoveryListResult) (ArmDisasterRecoveryListResult, error))`
- Function `NewNetworkRuleSetListResultPage` signature has been changed from `(func(context.Context, NetworkRuleSetListResult) (NetworkRuleSetListResult, error))` to `(NetworkRuleSetListResult,func(context.Context, NetworkRuleSetListResult) (NetworkRuleSetListResult, error))`
- Function `NewSBTopicListResultPage` signature has been changed from `(func(context.Context, SBTopicListResult) (SBTopicListResult, error))` to `(SBTopicListResult,func(context.Context, SBTopicListResult) (SBTopicListResult, error))`
- Function `NewMigrationConfigListResultPage` signature has been changed from `(func(context.Context, MigrationConfigListResult) (MigrationConfigListResult, error))` to `(MigrationConfigListResult,func(context.Context, MigrationConfigListResult) (MigrationConfigListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Struct `AuthorizationRuleProperties` has been removed
- Field `Code` of struct `ErrorResponse` has been removed
- Field `Message` of struct `ErrorResponse` has been removed

## New Content

- Struct `ErrorAdditionalInfo` is added
- Struct `ErrorResponseError` is added
- Field `Error` is added to struct `ErrorResponse`


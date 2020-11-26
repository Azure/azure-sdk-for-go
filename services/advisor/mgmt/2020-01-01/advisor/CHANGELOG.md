
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewConfigurationListResultPage` signature has been changed from `(func(context.Context, ConfigurationListResult) (ConfigurationListResult, error))` to `(ConfigurationListResult,func(context.Context, ConfigurationListResult) (ConfigurationListResult, error))`
- Function `NewSuppressionContractListResultPage` signature has been changed from `(func(context.Context, SuppressionContractListResult) (SuppressionContractListResult, error))` to `(SuppressionContractListResult,func(context.Context, SuppressionContractListResult) (SuppressionContractListResult, error))`
- Function `SuppressionsClient.GetResponder` return values have been changed from `(SuppressionContract,error)` to `(SetObject,error)`
- Function `NewResourceRecommendationBaseListResultPage` signature has been changed from `(func(context.Context, ResourceRecommendationBaseListResult) (ResourceRecommendationBaseListResult, error))` to `(ResourceRecommendationBaseListResult,func(context.Context, ResourceRecommendationBaseListResult) (ResourceRecommendationBaseListResult, error))`
- Function `NewOperationEntityListResultPage` signature has been changed from `(func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))` to `(OperationEntityListResult,func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))`
- Function `NewMetadataEntityListResultPage` signature has been changed from `(func(context.Context, MetadataEntityListResult) (MetadataEntityListResult, error))` to `(MetadataEntityListResult,func(context.Context, MetadataEntityListResult) (MetadataEntityListResult, error))`
- Function `SuppressionsClient.Get` return values have been changed from `(SuppressionContract,error)` to `(SetObject,error)`

## New Content

- Function `SuppressionProperties.MarshalJSON() ([]byte,error)` is added
- Field `ExpirationTimeStamp` is added to struct `SuppressionProperties`


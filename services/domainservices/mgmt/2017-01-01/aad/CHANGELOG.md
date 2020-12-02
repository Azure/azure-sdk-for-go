Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewDomainServiceListResultPage` parameter(s) have been changed from `(func(context.Context, DomainServiceListResult) (DomainServiceListResult, error))` to `(DomainServiceListResult, func(context.Context, DomainServiceListResult) (DomainServiceListResult, error))`
- Function `NewOperationEntityListResultPage` parameter(s) have been changed from `(func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))` to `(OperationEntityListResult, func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))`
- Type of `DomainServiceProperties.HealthLastEvaluated` has been changed from `*date.Time` to `*date.TimeRFC1123`

## New Content

- New field `DeploymentID` in struct `DomainServiceProperties`

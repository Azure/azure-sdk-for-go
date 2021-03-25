Generated from https://github.com/Azure/azure-rest-api-specs/tree/0f0e41fa4e3679510fcf03ecd60084f1cdbd5805/specification/deviceupdate/resource-manager/readme.md tag: `package-2020-03-01-preview`

Code generator @microsoft.azure/autorest.go@2.1.175


## Breaking Changes

### Removed Funcs

1. ErrorDefinition.MarshalJSON() ([]byte, error)
1. InstancesClient.ListBySubscription(context.Context, string) (InstanceListPage, error)
1. InstancesClient.ListBySubscriptionComplete(context.Context, string) (InstanceListIterator, error)
1. InstancesClient.ListBySubscriptionPreparer(context.Context, string) (*http.Request, error)
1. InstancesClient.ListBySubscriptionResponder(*http.Response) (InstanceList, error)
1. InstancesClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)

## Struct Changes

### Removed Struct Fields

1. ErrorDefinition.Code
1. ErrorDefinition.Details
1. ErrorDefinition.Message
1. ErrorResponse.Error

## Struct Changes

### New Structs

1. ErrorAdditionalInfo

### New Struct Fields

1. ErrorDefinition.Error
1. ErrorResponse.AdditionalInfo
1. ErrorResponse.Details
1. ErrorResponse.Message
1. ErrorResponse.Target

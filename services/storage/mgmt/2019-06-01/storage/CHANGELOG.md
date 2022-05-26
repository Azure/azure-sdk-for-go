# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. ErrorResponse.Code
1. ErrorResponse.Message

## Additive Changes

### New Constants

1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User
1. Name.AccessTimeTracking

### New Funcs

1. *BlobInventoryPolicy.UnmarshalJSON([]byte) error
1. BlobInventoryPoliciesClient.CreateOrUpdate(context.Context, string, string, BlobInventoryPolicy) (BlobInventoryPolicy, error)
1. BlobInventoryPoliciesClient.CreateOrUpdatePreparer(context.Context, string, string, BlobInventoryPolicy) (*http.Request, error)
1. BlobInventoryPoliciesClient.CreateOrUpdateResponder(*http.Response) (BlobInventoryPolicy, error)
1. BlobInventoryPoliciesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. BlobInventoryPoliciesClient.Delete(context.Context, string, string) (autorest.Response, error)
1. BlobInventoryPoliciesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. BlobInventoryPoliciesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. BlobInventoryPoliciesClient.DeleteSender(*http.Request) (*http.Response, error)
1. BlobInventoryPoliciesClient.Get(context.Context, string, string) (BlobInventoryPolicy, error)
1. BlobInventoryPoliciesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. BlobInventoryPoliciesClient.GetResponder(*http.Response) (BlobInventoryPolicy, error)
1. BlobInventoryPoliciesClient.GetSender(*http.Request) (*http.Response, error)
1. BlobInventoryPoliciesClient.List(context.Context, string, string) (ListBlobInventoryPolicy, error)
1. BlobInventoryPoliciesClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. BlobInventoryPoliciesClient.ListResponder(*http.Response) (ListBlobInventoryPolicy, error)
1. BlobInventoryPoliciesClient.ListSender(*http.Request) (*http.Response, error)
1. BlobInventoryPolicy.MarshalJSON() ([]byte, error)
1. BlobInventoryPolicyProperties.MarshalJSON() ([]byte, error)
1. ListBlobInventoryPolicy.MarshalJSON() ([]byte, error)
1. NewBlobInventoryPoliciesClient(string) BlobInventoryPoliciesClient
1. NewBlobInventoryPoliciesClientWithBaseURI(string, string) BlobInventoryPoliciesClient
1. PossibleCreatedByTypeValues() []CreatedByType
1. PossibleNameValues() []Name

### Struct Changes

#### New Structs

1. BlobInventoryPoliciesClient
1. BlobInventoryPolicy
1. BlobInventoryPolicyDefinition
1. BlobInventoryPolicyFilter
1. BlobInventoryPolicyProperties
1. BlobInventoryPolicyRule
1. BlobInventoryPolicySchema
1. ErrorResponseBody
1. LastAccessTimeTrackingPolicy
1. ListBlobInventoryPolicy
1. ManagementPolicyVersion
1. SystemData

#### New Struct Fields

1. AccountProperties.AllowSharedKeyAccess
1. AccountPropertiesCreateParameters.AllowSharedKeyAccess
1. AccountPropertiesUpdateParameters.AllowSharedKeyAccess
1. BlobServicePropertiesProperties.LastAccessTimeTrackingPolicy
1. ChangeFeed.RetentionInDays
1. DateAfterModification.DaysAfterLastAccessTimeGreaterThan
1. ErrorResponse.Error
1. ManagementPolicyAction.Version
1. ManagementPolicyBaseBlob.EnableAutoTierToHotFromCool
1. ManagementPolicySnapShot.TierToArchive
1. ManagementPolicySnapShot.TierToCool

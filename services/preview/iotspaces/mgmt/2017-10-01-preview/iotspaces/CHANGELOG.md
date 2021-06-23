# Unreleased

## Breaking Changes

### Removed Constants

1. NameUnavailabilityReason.AlreadyExists
1. NameUnavailabilityReason.Invalid
1. ProvisioningState.Canceled
1. ProvisioningState.Deleting
1. ProvisioningState.Failed
1. ProvisioningState.Provisioning
1. ProvisioningState.Succeeded
1. Sku.F1
1. Sku.S1
1. Sku.S2
1. Sku.S3

### Removed Funcs

1. *CreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *DeleteFuture.UnmarshalJSON([]byte) error
1. *DescriptionListResultIterator.Next() error
1. *DescriptionListResultIterator.NextWithContext(context.Context) error
1. *DescriptionListResultPage.Next() error
1. *DescriptionListResultPage.NextWithContext(context.Context) error
1. *OperationListResultIterator.Next() error
1. *OperationListResultIterator.NextWithContext(context.Context) error
1. *OperationListResultPage.Next() error
1. *OperationListResultPage.NextWithContext(context.Context) error
1. *UpdateFuture.UnmarshalJSON([]byte) error
1. Client.CheckNameAvailability(context.Context, OperationInputs) (NameAvailabilityInfo, error)
1. Client.CheckNameAvailabilityPreparer(context.Context, OperationInputs) (*http.Request, error)
1. Client.CheckNameAvailabilityResponder(*http.Response) (NameAvailabilityInfo, error)
1. Client.CheckNameAvailabilitySender(*http.Request) (*http.Response, error)
1. Client.CreateOrUpdate(context.Context, string, string, Description) (CreateOrUpdateFuture, error)
1. Client.CreateOrUpdatePreparer(context.Context, string, string, Description) (*http.Request, error)
1. Client.CreateOrUpdateResponder(*http.Response) (Description, error)
1. Client.CreateOrUpdateSender(*http.Request) (CreateOrUpdateFuture, error)
1. Client.Delete(context.Context, string, string) (DeleteFuture, error)
1. Client.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. Client.DeleteResponder(*http.Response) (Description, error)
1. Client.DeleteSender(*http.Request) (DeleteFuture, error)
1. Client.Get(context.Context, string, string) (Description, error)
1. Client.GetPreparer(context.Context, string, string) (*http.Request, error)
1. Client.GetResponder(*http.Response) (Description, error)
1. Client.GetSender(*http.Request) (*http.Response, error)
1. Client.List(context.Context) (DescriptionListResultPage, error)
1. Client.ListByResourceGroup(context.Context, string) (DescriptionListResultPage, error)
1. Client.ListByResourceGroupComplete(context.Context, string) (DescriptionListResultIterator, error)
1. Client.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. Client.ListByResourceGroupResponder(*http.Response) (DescriptionListResult, error)
1. Client.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. Client.ListComplete(context.Context) (DescriptionListResultIterator, error)
1. Client.ListPreparer(context.Context) (*http.Request, error)
1. Client.ListResponder(*http.Response) (DescriptionListResult, error)
1. Client.ListSender(*http.Request) (*http.Response, error)
1. Client.Update(context.Context, string, string, PatchDescription) (UpdateFuture, error)
1. Client.UpdatePreparer(context.Context, string, string, PatchDescription) (*http.Request, error)
1. Client.UpdateResponder(*http.Response) (Description, error)
1. Client.UpdateSender(*http.Request) (UpdateFuture, error)
1. Description.MarshalJSON() ([]byte, error)
1. DescriptionListResult.IsEmpty() bool
1. DescriptionListResultIterator.NotDone() bool
1. DescriptionListResultIterator.Response() DescriptionListResult
1. DescriptionListResultIterator.Value() Description
1. DescriptionListResultPage.NotDone() bool
1. DescriptionListResultPage.Response() DescriptionListResult
1. DescriptionListResultPage.Values() []Description
1. ErrorDetails.MarshalJSON() ([]byte, error)
1. NameAvailabilityInfo.MarshalJSON() ([]byte, error)
1. NewClient(uuid.UUID) Client
1. NewClientWithBaseURI(string, uuid.UUID) Client
1. NewDescriptionListResultIterator(DescriptionListResultPage) DescriptionListResultIterator
1. NewDescriptionListResultPage(DescriptionListResult, func(context.Context, DescriptionListResult) (DescriptionListResult, error)) DescriptionListResultPage
1. NewOperationListResultIterator(OperationListResultPage) OperationListResultIterator
1. NewOperationListResultPage(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage
1. NewOperationsClient(uuid.UUID) OperationsClient
1. NewOperationsClientWithBaseURI(string, uuid.UUID) OperationsClient
1. Operation.MarshalJSON() ([]byte, error)
1. OperationDisplay.MarshalJSON() ([]byte, error)
1. OperationListResult.IsEmpty() bool
1. OperationListResult.MarshalJSON() ([]byte, error)
1. OperationListResultIterator.NotDone() bool
1. OperationListResultIterator.Response() OperationListResult
1. OperationListResultIterator.Value() Operation
1. OperationListResultPage.NotDone() bool
1. OperationListResultPage.Response() OperationListResult
1. OperationListResultPage.Values() []Operation
1. OperationsClient.List(context.Context) (OperationListResultPage, error)
1. OperationsClient.ListComplete(context.Context) (OperationListResultIterator, error)
1. OperationsClient.ListPreparer(context.Context) (*http.Request, error)
1. OperationsClient.ListResponder(*http.Response) (OperationListResult, error)
1. OperationsClient.ListSender(*http.Request) (*http.Response, error)
1. PatchDescription.MarshalJSON() ([]byte, error)
1. PossibleNameUnavailabilityReasonValues() []NameUnavailabilityReason
1. PossibleProvisioningStateValues() []ProvisioningState
1. PossibleSkuValues() []Sku
1. Properties.MarshalJSON() ([]byte, error)
1. Resource.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. Client
1. CreateOrUpdateFuture
1. DeleteFuture
1. Description
1. DescriptionListResult
1. DescriptionListResultIterator
1. DescriptionListResultPage
1. ErrorDetails
1. NameAvailabilityInfo
1. Operation
1. OperationDisplay
1. OperationInputs
1. OperationListResult
1. OperationListResultIterator
1. OperationListResultPage
1. OperationsClient
1. PatchDescription
1. Properties
1. Resource
1. SkuInfo
1. StorageContainerProperties
1. UpdateFuture


Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## New Content

- Function `BaseClient.PublishCustomEventEvents(context.Context,string,[]interface{}) (autorest.Response,error)` is added
- Function `BaseClient.PublishCloudEventEventsResponder(*http.Response) (autorest.Response,error)` is added
- Function `CloudEventEvent.MarshalJSON() ([]byte,error)` is added
- Function `BaseClient.PublishCustomEventEventsSender(*http.Request) (*http.Response,error)` is added
- Function `BaseClient.PublishCloudEventEventsSender(*http.Request) (*http.Response,error)` is added
- Function `BaseClient.PublishCloudEventEvents(context.Context,string,[]CloudEventEvent) (autorest.Response,error)` is added
- Function `*CloudEventEvent.UnmarshalJSON([]byte) error` is added
- Function `BaseClient.PublishCloudEventEventsPreparer(context.Context,string,[]CloudEventEvent) (*http.Request,error)` is added
- Function `BaseClient.PublishCustomEventEventsResponder(*http.Response) (autorest.Response,error)` is added
- Function `ACSChatThreadPropertiesUpdatedPerUserEventData.MarshalJSON() ([]byte,error)` is added
- Function `ACSChatThreadCreatedWithUserEventData.MarshalJSON() ([]byte,error)` is added
- Function `BaseClient.PublishCustomEventEventsPreparer(context.Context,string,[]interface{}) (*http.Request,error)` is added
- Struct `ACSChatEventBaseProperties` is added
- Struct `ACSChatMemberAddedToThreadWithUserEventData` is added
- Struct `ACSChatMemberRemovedFromThreadWithUserEventData` is added
- Struct `ACSChatMessageDeletedEventData` is added
- Struct `ACSChatMessageEditedEventData` is added
- Struct `ACSChatMessageEventBaseProperties` is added
- Struct `ACSChatMessageReceivedEventData` is added
- Struct `ACSChatThreadCreatedWithUserEventData` is added
- Struct `ACSChatThreadEventBaseProperties` is added
- Struct `ACSChatThreadMemberProperties` is added
- Struct `ACSChatThreadPropertiesUpdatedPerUserEventData` is added
- Struct `ACSChatThreadWithUserDeletedEventData` is added
- Struct `AcsSmsDeliveryAttemptProperties` is added
- Struct `AcsSmsDeliveryReportReceivedEventData` is added
- Struct `AcsSmsEventBaseProperties` is added
- Struct `AcsSmsReceivedEventData` is added
- Struct `CloudEventEvent` is added
- Struct `KeyVaultVaultAccessPolicyChangedEventData` is added


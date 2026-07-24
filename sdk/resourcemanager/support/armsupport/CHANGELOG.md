# Release History

## 2.1.0 (2026-07-06)
### Features Added

- New enum type `Channel` with values `ChannelChat`, `ChannelWeb`
- New enum type `ChatConversationStatus` with values `ChatConversationStatusActive`, `ChatConversationStatusClosed`
- New enum type `EscalationStatus` with values `EscalationStatusEscalationAvailable`, `EscalationStatusEscalationInitiated`, `EscalationStatusEscalationProcessed`, `EscalationStatusEscalationUnavailable`, `EscalationStatusEscalationUnsupported`
- New function `NewClassifyProblemsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClassifyProblemsClient, error)`
- New function `*ClassifyProblemsClient.ClassifyProblems(ctx context.Context, problemServiceName string, problemClassificationsClassificationInput ProblemClassificationsClassificationInput, options *ClassifyProblemsClientClassifyProblemsOptions) (ClassifyProblemsClientClassifyProblemsResponse, error)`
- New function `NewClassifyProblemsNoSubscriptionClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ClassifyProblemsNoSubscriptionClient, error)`
- New function `*ClassifyProblemsNoSubscriptionClient.ClassifyProblems(ctx context.Context, problemServiceName string, problemClassificationsClassificationInput ProblemClassificationsClassificationInput, options *ClassifyProblemsNoSubscriptionClientClassifyProblemsOptions) (ClassifyProblemsNoSubscriptionClientClassifyProblemsResponse, error)`
- New function `NewClassifyServicesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClassifyServicesClient, error)`
- New function `*ClassifyServicesClient.ClassifyServices(ctx context.Context, serviceClassificationRequest ServiceClassificationRequest, options *ClassifyServicesClientClassifyServicesOptions) (ClassifyServicesClientClassifyServicesResponse, error)`
- New function `NewClassifyServicesNoSubscriptionClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ClassifyServicesNoSubscriptionClient, error)`
- New function `*ClassifyServicesNoSubscriptionClient.ClassifyServices(ctx context.Context, serviceClassificationRequest ServiceClassificationRequest, options *ClassifyServicesNoSubscriptionClientClassifyServicesOptions) (ClassifyServicesNoSubscriptionClientClassifyServicesResponse, error)`
- New function `*ClientFactory.NewClassifyProblemsClient() *ClassifyProblemsClient`
- New function `*ClientFactory.NewClassifyProblemsNoSubscriptionClient() *ClassifyProblemsNoSubscriptionClient`
- New function `*ClientFactory.NewClassifyServicesClient() *ClassifyServicesClient`
- New function `*ClientFactory.NewClassifyServicesNoSubscriptionClient() *ClassifyServicesNoSubscriptionClient`
- New function `*TicketsClient.LookUpResourceID(ctx context.Context, lookUpResourceIDRequest LookUpResourceIDRequest, options *TicketsClientLookUpResourceIDOptions) (TicketsClientLookUpResourceIDResponse, error)`
- New struct `ClassificationService`
- New struct `DirectConnectEscalation`
- New struct `LookUpResourceIDRequest`
- New struct `LookUpResourceIDResponse`
- New struct `ProblemClassificationsClassificationInput`
- New struct `ProblemClassificationsClassificationOutput`
- New struct `ProblemClassificationsClassificationResult`
- New struct `ServiceClassificationAnswer`
- New struct `ServiceClassificationOutput`
- New struct `ServiceClassificationRequest`
- New field `ChatConversationStatus`, `CommunityForumPost`, `DirectConnectEscalation`, `SupportChannel` in struct `TicketDetailsProperties`
- New field `DirectConnectEscalation` in struct `UpdateSupportTicket`


## 2.0.0 (2026-06-24)
### Breaking Changes

- Type of `MessageProperties.ContentType` has been changed from `*TranscriptContentType` to `*string`
- Enum `TranscriptContentType` has been removed
- Function `PossibleTranscriptContentTypeValues` has been removed
- Struct `OperationsListResult` has been removed
- Field `OperationsListResult` of struct `OperationsClientListResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New struct `OperationListResult`
- New field `SystemData` in struct `CommunicationDetails`
- New field `ActionType`, `IsDataAction`, `Origin` in struct `Operation`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New field `SystemData` in struct `ProblemClassification`
- New field `NextLink` in struct `ProblemClassificationsListResult`
- New field `SystemData` in struct `Service`
- New field `NextLink` in struct `ServicesListResult`
- New field `SystemData` in struct `TicketDetails`


## 1.3.0 (2024-04-26)
### Features Added

- New enum type `Consent` with values `ConsentNo`, `ConsentYes`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `IsTemporaryTicket` with values `IsTemporaryTicketNo`, `IsTemporaryTicketYes`
- New enum type `UserConsent` with values `UserConsentNo`, `UserConsentYes`
- New function `NewChatTranscriptsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ChatTranscriptsClient, error)`
- New function `*ChatTranscriptsClient.Get(context.Context, string, string, *ChatTranscriptsClientGetOptions) (ChatTranscriptsClientGetResponse, error)`
- New function `*ChatTranscriptsClient.NewListPager(string, *ChatTranscriptsClientListOptions) *runtime.Pager[ChatTranscriptsClientListResponse]`
- New function `NewChatTranscriptsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*ChatTranscriptsNoSubscriptionClient, error)`
- New function `*ChatTranscriptsNoSubscriptionClient.Get(context.Context, string, string, *ChatTranscriptsNoSubscriptionClientGetOptions) (ChatTranscriptsNoSubscriptionClientGetResponse, error)`
- New function `*ChatTranscriptsNoSubscriptionClient.NewListPager(string, *ChatTranscriptsNoSubscriptionClientListOptions) *runtime.Pager[ChatTranscriptsNoSubscriptionClientListResponse]`
- New function `*ClientFactory.NewChatTranscriptsClient() *ChatTranscriptsClient`
- New function `*ClientFactory.NewChatTranscriptsNoSubscriptionClient() *ChatTranscriptsNoSubscriptionClient`
- New function `*ClientFactory.NewCommunicationsNoSubscriptionClient() *CommunicationsNoSubscriptionClient`
- New function `*ClientFactory.NewFileWorkspacesClient() *FileWorkspacesClient`
- New function `*ClientFactory.NewFileWorkspacesNoSubscriptionClient() *FileWorkspacesNoSubscriptionClient`
- New function `*ClientFactory.NewFilesClient() *FilesClient`
- New function `*ClientFactory.NewFilesNoSubscriptionClient() *FilesNoSubscriptionClient`
- New function `*ClientFactory.NewTicketsNoSubscriptionClient() *TicketsNoSubscriptionClient`
- New function `NewCommunicationsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*CommunicationsNoSubscriptionClient, error)`
- New function `*CommunicationsNoSubscriptionClient.CheckNameAvailability(context.Context, string, CheckNameAvailabilityInput, *CommunicationsNoSubscriptionClientCheckNameAvailabilityOptions) (CommunicationsNoSubscriptionClientCheckNameAvailabilityResponse, error)`
- New function `*CommunicationsNoSubscriptionClient.BeginCreate(context.Context, string, string, CommunicationDetails, *CommunicationsNoSubscriptionClientBeginCreateOptions) (*runtime.Poller[CommunicationsNoSubscriptionClientCreateResponse], error)`
- New function `*CommunicationsNoSubscriptionClient.Get(context.Context, string, string, *CommunicationsNoSubscriptionClientGetOptions) (CommunicationsNoSubscriptionClientGetResponse, error)`
- New function `*CommunicationsNoSubscriptionClient.NewListPager(string, *CommunicationsNoSubscriptionClientListOptions) *runtime.Pager[CommunicationsNoSubscriptionClientListResponse]`
- New function `NewFileWorkspacesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FileWorkspacesClient, error)`
- New function `*FileWorkspacesClient.Create(context.Context, string, *FileWorkspacesClientCreateOptions) (FileWorkspacesClientCreateResponse, error)`
- New function `*FileWorkspacesClient.Get(context.Context, string, *FileWorkspacesClientGetOptions) (FileWorkspacesClientGetResponse, error)`
- New function `NewFileWorkspacesNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*FileWorkspacesNoSubscriptionClient, error)`
- New function `*FileWorkspacesNoSubscriptionClient.Create(context.Context, string, *FileWorkspacesNoSubscriptionClientCreateOptions) (FileWorkspacesNoSubscriptionClientCreateResponse, error)`
- New function `*FileWorkspacesNoSubscriptionClient.Get(context.Context, string, *FileWorkspacesNoSubscriptionClientGetOptions) (FileWorkspacesNoSubscriptionClientGetResponse, error)`
- New function `NewFilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FilesClient, error)`
- New function `*FilesClient.Create(context.Context, string, string, FileDetails, *FilesClientCreateOptions) (FilesClientCreateResponse, error)`
- New function `*FilesClient.Get(context.Context, string, string, *FilesClientGetOptions) (FilesClientGetResponse, error)`
- New function `*FilesClient.NewListPager(string, *FilesClientListOptions) *runtime.Pager[FilesClientListResponse]`
- New function `*FilesClient.Upload(context.Context, string, string, UploadFile, *FilesClientUploadOptions) (FilesClientUploadResponse, error)`
- New function `NewFilesNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*FilesNoSubscriptionClient, error)`
- New function `*FilesNoSubscriptionClient.Create(context.Context, string, string, FileDetails, *FilesNoSubscriptionClientCreateOptions) (FilesNoSubscriptionClientCreateResponse, error)`
- New function `*FilesNoSubscriptionClient.Get(context.Context, string, string, *FilesNoSubscriptionClientGetOptions) (FilesNoSubscriptionClientGetResponse, error)`
- New function `*FilesNoSubscriptionClient.NewListPager(string, *FilesNoSubscriptionClientListOptions) *runtime.Pager[FilesNoSubscriptionClientListResponse]`
- New function `*FilesNoSubscriptionClient.Upload(context.Context, string, string, UploadFile, *FilesNoSubscriptionClientUploadOptions) (FilesNoSubscriptionClientUploadResponse, error)`
- New function `NewTicketsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*TicketsNoSubscriptionClient, error)`
- New function `*TicketsNoSubscriptionClient.CheckNameAvailability(context.Context, CheckNameAvailabilityInput, *TicketsNoSubscriptionClientCheckNameAvailabilityOptions) (TicketsNoSubscriptionClientCheckNameAvailabilityResponse, error)`
- New function `*TicketsNoSubscriptionClient.BeginCreate(context.Context, string, TicketDetails, *TicketsNoSubscriptionClientBeginCreateOptions) (*runtime.Poller[TicketsNoSubscriptionClientCreateResponse], error)`
- New function `*TicketsNoSubscriptionClient.Get(context.Context, string, *TicketsNoSubscriptionClientGetOptions) (TicketsNoSubscriptionClientGetResponse, error)`
- New function `*TicketsNoSubscriptionClient.NewListPager(*TicketsNoSubscriptionClientListOptions) *runtime.Pager[TicketsNoSubscriptionClientListResponse]`
- New function `*TicketsNoSubscriptionClient.Update(context.Context, string, UpdateSupportTicket, *TicketsNoSubscriptionClientUpdateOptions) (TicketsNoSubscriptionClientUpdateResponse, error)`
- New struct `ChatTranscriptDetails`
- New struct `ChatTranscriptDetailsProperties`
- New struct `ChatTranscriptsListResult`
- New struct `FileDetails`
- New struct `FileDetailsProperties`
- New struct `FileWorkspaceDetails`
- New struct `FileWorkspaceDetailsProperties`
- New struct `FilesListResult`
- New struct `MessageProperties`
- New struct `SecondaryConsent`
- New struct `SecondaryConsentEnabled`
- New struct `SystemData`
- New struct `UploadFile`
- New field `SecondaryConsentEnabled` in struct `ProblemClassificationProperties`
- New field `AdvancedDiagnosticConsent`, `FileWorkspaceName`, `IsTemporaryTicket`, `ProblemScopingQuestions`, `SecondaryConsent`, `SupportPlanDisplayName`, `SupportPlanID` in struct `TicketDetailsProperties`
- New field `AdvancedDiagnosticConsent`, `SecondaryConsent` in struct `UpdateSupportTicket`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/support/armsupport` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
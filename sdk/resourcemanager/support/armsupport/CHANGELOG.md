# Release History

## 2.0.0-beta.3 (2024-03-22)
### Breaking Changes

- Type of `FileDetailsProperties.ChunkSize` has been changed from `*float32` to `*int32`
- Type of `FileDetailsProperties.FileSize` has been changed from `*float32` to `*int32`
- Type of `FileDetailsProperties.NumberOfChunks` has been changed from `*float32` to `*int32`
- Type of `UploadFile.ChunkIndex` has been changed from `*float32` to `*int32`
- Function `*ClientFactory.NewTicketChatTranscriptsNoSubscriptionClient` has been removed
- Function `*ClientFactory.NewTicketCommunicationsNoSubscriptionClient` has been removed
- Function `NewTicketChatTranscriptsNoSubscriptionClient` has been removed
- Function `*TicketChatTranscriptsNoSubscriptionClient.NewListPager` has been removed
- Function `NewTicketCommunicationsNoSubscriptionClient` has been removed
- Function `*TicketCommunicationsNoSubscriptionClient.NewListPager` has been removed

### Features Added

- New enum type `IsTemporaryTicket` with values `IsTemporaryTicketNo`, `IsTemporaryTicketYes`
- New function `*ChatTranscriptsNoSubscriptionClient.NewListPager(string, *ChatTranscriptsNoSubscriptionClientListOptions) *runtime.Pager[ChatTranscriptsNoSubscriptionClientListResponse]`
- New function `*ClientFactory.NewLookUpResourceIDClient() *LookUpResourceIDClient`
- New function `*ClientFactory.NewProblemClassificationsNoSubscriptionClient() *ProblemClassificationsNoSubscriptionClient`
- New function `*ClientFactory.NewServiceClassificationsClient() *ServiceClassificationsClient`
- New function `*ClientFactory.NewServiceClassificationsNoSubscriptionClient() *ServiceClassificationsNoSubscriptionClient`
- New function `*CommunicationsNoSubscriptionClient.NewListPager(string, *CommunicationsNoSubscriptionClientListOptions) *runtime.Pager[CommunicationsNoSubscriptionClientListResponse]`
- New function `NewLookUpResourceIDClient(azcore.TokenCredential, *arm.ClientOptions) (*LookUpResourceIDClient, error)`
- New function `*LookUpResourceIDClient.Post(context.Context, LookUpResourceIDRequest, *LookUpResourceIDClientPostOptions) (LookUpResourceIDClientPostResponse, error)`
- New function `*ProblemClassificationsClient.ClassifyProblems(context.Context, string, ProblemClassificationsClassificationInput, *ProblemClassificationsClientClassifyProblemsOptions) (ProblemClassificationsClientClassifyProblemsResponse, error)`
- New function `NewProblemClassificationsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*ProblemClassificationsNoSubscriptionClient, error)`
- New function `*ProblemClassificationsNoSubscriptionClient.ClassifyProblems(context.Context, string, ProblemClassificationsClassificationInput, *ProblemClassificationsNoSubscriptionClientClassifyProblemsOptions) (ProblemClassificationsNoSubscriptionClientClassifyProblemsResponse, error)`
- New function `NewServiceClassificationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServiceClassificationsClient, error)`
- New function `*ServiceClassificationsClient.ClassifyServices(context.Context, ServiceClassificationRequest, *ServiceClassificationsClientClassifyServicesOptions) (ServiceClassificationsClientClassifyServicesResponse, error)`
- New function `NewServiceClassificationsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*ServiceClassificationsNoSubscriptionClient, error)`
- New function `*ServiceClassificationsNoSubscriptionClient.ClassifyServices(context.Context, ServiceClassificationRequest, *ServiceClassificationsNoSubscriptionClientClassifyServicesOptions) (ServiceClassificationsNoSubscriptionClientClassifyServicesResponse, error)`
- New struct `ClassificationService`
- New struct `LookUpResourceIDRequest`
- New struct `LookUpResourceIDResponse`
- New struct `ProblemClassificationsClassificationInput`
- New struct `ProblemClassificationsClassificationOutput`
- New struct `ProblemClassificationsClassificationResult`
- New struct `ServiceClassificationAnswer`
- New struct `ServiceClassificationOutput`
- New struct `ServiceClassificationRequest`
- New field `Metadata`, `ParentProblemClassification` in struct `ProblemClassificationProperties`
- New field `Metadata` in struct `ServiceProperties`
- New field `IsTemporaryTicket` in struct `TicketDetailsProperties`


## 2.0.0-beta.2 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0-beta.1 (2023-10-27)
### Breaking Changes

- Struct `ExceptionResponse` has been removed
- Struct `ServiceError` has been removed
- Struct `ServiceErrorDetail` has been removed

### Features Added

- New enum type `Consent` with values `ConsentNo`, `ConsentYes`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `UserConsent` with values `UserConsentNo`, `UserConsentYes`
- New function `NewChatTranscriptsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ChatTranscriptsClient, error)`
- New function `*ChatTranscriptsClient.Get(context.Context, string, string, *ChatTranscriptsClientGetOptions) (ChatTranscriptsClientGetResponse, error)`
- New function `*ChatTranscriptsClient.NewListPager(string, *ChatTranscriptsClientListOptions) *runtime.Pager[ChatTranscriptsClientListResponse]`
- New function `NewChatTranscriptsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*ChatTranscriptsNoSubscriptionClient, error)`
- New function `*ChatTranscriptsNoSubscriptionClient.Get(context.Context, string, string, *ChatTranscriptsNoSubscriptionClientGetOptions) (ChatTranscriptsNoSubscriptionClientGetResponse, error)`
- New function `*ClientFactory.NewChatTranscriptsClient() *ChatTranscriptsClient`
- New function `*ClientFactory.NewChatTranscriptsNoSubscriptionClient() *ChatTranscriptsNoSubscriptionClient`
- New function `*ClientFactory.NewCommunicationsNoSubscriptionClient() *CommunicationsNoSubscriptionClient`
- New function `*ClientFactory.NewFileWorkspacesClient() *FileWorkspacesClient`
- New function `*ClientFactory.NewFileWorkspacesNoSubscriptionClient() *FileWorkspacesNoSubscriptionClient`
- New function `*ClientFactory.NewFilesClient() *FilesClient`
- New function `*ClientFactory.NewFilesNoSubscriptionClient() *FilesNoSubscriptionClient`
- New function `*ClientFactory.NewTicketChatTranscriptsNoSubscriptionClient() *TicketChatTranscriptsNoSubscriptionClient`
- New function `*ClientFactory.NewTicketCommunicationsNoSubscriptionClient() *TicketCommunicationsNoSubscriptionClient`
- New function `*ClientFactory.NewTicketsNoSubscriptionClient() *TicketsNoSubscriptionClient`
- New function `NewCommunicationsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*CommunicationsNoSubscriptionClient, error)`
- New function `*CommunicationsNoSubscriptionClient.CheckNameAvailability(context.Context, string, CheckNameAvailabilityInput, *CommunicationsNoSubscriptionClientCheckNameAvailabilityOptions) (CommunicationsNoSubscriptionClientCheckNameAvailabilityResponse, error)`
- New function `*CommunicationsNoSubscriptionClient.BeginCreate(context.Context, string, string, CommunicationDetails, *CommunicationsNoSubscriptionClientBeginCreateOptions) (*runtime.Poller[CommunicationsNoSubscriptionClientCreateResponse], error)`
- New function `*CommunicationsNoSubscriptionClient.Get(context.Context, string, string, *CommunicationsNoSubscriptionClientGetOptions) (CommunicationsNoSubscriptionClientGetResponse, error)`
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
- New function `NewTicketChatTranscriptsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*TicketChatTranscriptsNoSubscriptionClient, error)`
- New function `*TicketChatTranscriptsNoSubscriptionClient.NewListPager(string, *TicketChatTranscriptsNoSubscriptionClientListOptions) *runtime.Pager[TicketChatTranscriptsNoSubscriptionClientListResponse]`
- New function `NewTicketCommunicationsNoSubscriptionClient(azcore.TokenCredential, *arm.ClientOptions) (*TicketCommunicationsNoSubscriptionClient, error)`
- New function `*TicketCommunicationsNoSubscriptionClient.NewListPager(string, *TicketCommunicationsNoSubscriptionClientListOptions) *runtime.Pager[TicketCommunicationsNoSubscriptionClientListResponse]`
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
- New field `AdvancedDiagnosticConsent`, `FileWorkspaceName`, `ProblemScopingQuestions`, `SecondaryConsent`, `SupportPlanDisplayName`, `SupportPlanID` in struct `TicketDetailsProperties`
- New field `AdvancedDiagnosticConsent`, `SecondaryConsent` in struct `UpdateSupportTicket`


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/support/armsupport` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
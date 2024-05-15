# Release History

## 0.4.0 (Unreleased)

### Breaking Changes

- `Type` has been removed, making it simpler to compare the EventGridEvent.Type and CloudEvent.Type values against
our provided constants.

- The following models have had 'Advanced' removed from their name:
  - ACSMessageButtonContent
  - ACSMessageContext
  - ACSMessageDeliveryStatusUpdatedEventData
  - ACSMessageEventData
  - ACSMessageInteractiveButtonReplyContent
  - ACSMessageInteractiveContent
  - ACSMessageInteractiveListReplyContent
  - ACSMessageMediaContent
  - ACSMessageReceivedEventData

- Several models' fields have had renames to improve consistency with other languages:
  - ACSEmailDeliveryReportReceivedEventData:
    - DeliveryAttemptTimestamp -> DeliveryAttemptTimeStamp
  - ACSEmailEngagementTrackingReportReceivedEventData:
    - Engagement -> EngagementType
    - UserActionTimestamp -> UserActionTimeStamp
  - ACSIncomingCallEventData:
    - FromCommunicationIdentifier -> From
    - ToCommunicationIdentifier -> To
  - ACSMessageContext:
    - MessageID -> ID
  - ACSMessageDeliveryStatusUpdatedEventData:
    - ChannelKind -> ChannelType
    - ReceivedTimestamp -> ReceivedTimeStamp
  - ACSMessageEventData:
    - ReceivedTimestamp -> ReceivedTimeStamp
  - ACSMessageInteractiveButtonReplyContent:
    - ButtonID -> ID
  - ACSMessageInteractiveContent:
    - ReplyKind -> Type
  - ACSMessageInteractiveListReplyContent:
    - ListItemID -> ID
  - ACSMessageMediaContent:
    - MediaID -> ID
  - ACSMessageReceivedEventData:
    - ChannelKind -> ChannelType
    - InteractiveContent -> Interactive
    - MediaContent -> Media
    - ReceivedTimestamp -> ReceivedTimeStamp
  - ACSRouterWorkerSelector:
    - LabelOperator -> Operator
    - State -> SelectorState
    - TimeToLive -> TTLSeconds
  - EventHubCaptureFileCreatedEventData:
    - Fileurl -> FileURL
  - HealthcareFhirResourceCreatedEventData
    - FhirResourceID -> ResourceFhirID
    - FhirResourceType -> ResourceType
    - FhirResourceVersionID -> ResourceVersionID
    - FhirServiceHostName -> ResourceFhirAccount
  - HealthcareFhirResourceDeletedEventData
    - FhirResourceID -> ResourceFhirID
    - FhirResourceType -> ResourceType
    - FhirResourceVersionID -> ResourceVersionID
    - FhirServiceHostName -> ResourceFhirAccount
  - HealthcareFhirResourceUpdatedEventData
    - FhirResourceID -> ResourceFhirID
    - FhirResourceType -> ResourceType
    - FhirResourceVersionID -> ResourceVersionID
    - FhirServiceHostName -> ResourceFhirAccount
  - ResourceNotificationsHealthResourcesAnnotatedEventData
    - OperationalDetails -> OperationalInfo
    - ResourceDetails -> ResourceInfo
  - ResourceNotificationsHealthResourcesAvailabilityStatusChangedEventData
    - OperationalDetails -> OperationalInfo
    - ResourceDetails -> ResourceInfo
  - ResourceNotificationsResourceDeletedEventData
    - OperationalDetails -> OperationalInfo
    - ResourceDetails -> ResourceInfo
  - ResourceNotificationsResourceManagementCreatedOrUpdatedEventData
    - OperationalDetails -> OperationalInfo
    - ResourceDetails -> ResourceInfo
  - ResourceNotificationsResourceManagementDeletedEventData
    - OperationalDetails -> OperationalInfo
    - ResourceDetails -> ResourceInfo
  - ResourceNotificationsResourceUpdatedEventData
    - OperationalDetails -> OperationalInfo
    - ResourceDetails -> ResourceInfo
  - StorageTaskAssignmentCompletedEventData:
    - CompletedOn -> CompletedDateTime
    - SummaryReportBlobURI -> SummaryReportBlobURL
  - StorageTaskAssignmentQueuedEventData:
    - QueuedOn -> QueuedDateTime

## 0.3.0 (2024-04-03)

### Features Added

- Added events ACSRouterWorkerUpdatedEventData and ACSAdvancedMessageDeliveryStatusUpdatedEventData. (PR#22638)

### Breaking Changes

Field and type renames:
- Globally, types and fields named ChannelType has been renamed to ChannelKind
- ACS events and constants have been changed to use an all-caps name (ex: AcsEmailDeliveryReportStatusDetails -> ACSEmailDeliveryReportStatusDetails).
- ACSAdvancedMessageContext.ID -> MessageID
- ACSAdvancedMessageReceivedEventData
  - .Media -> MediaContent
  - .Interactive -> InteractiveContent

## 0.2.0 (2024-03-14)

### Features Added

- Added API Center system events under their official names.

### Breaking Changes

- Events have been renamed:
  - APIDefinitionAddedEventData renamed to APICenterAPIDefinitionAddedEventData
  - APIDefinitionUpdatedEventData renamed to APICenterAPIDefinitionUpdatedEventData

## 0.1.0 (2024-03-05)

### Features Added

- Initial preview for Event Grid system events.

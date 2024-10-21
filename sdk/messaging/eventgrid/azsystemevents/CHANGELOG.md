# Release History

## 0.4.4 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.4.3 (2024-10-14)

### Features Added

- New field has been added to ACSIncomingCallEventData: OnBehalfOfCallee.

## 0.4.2 (2024-09-19)

### Features Added

- A new field has been added to StorageLifecyclePolicyCompletedEventData:
  - TierToColdSummary

## 0.4.1 (2024-08-20)

### Features Added

- New fields have been added:
  - StorageBlobCreatedEventData: AccessTier
  - StorageBlobTierChangedEventData: AccessTier and PreviousTier

### Breaking Changes

- Models that were not system events (ex: ACSChatMessageEventInThreadBaseProperties), or referenced by system events, have been removed.

## 0.4.0 (2024-06-11)

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

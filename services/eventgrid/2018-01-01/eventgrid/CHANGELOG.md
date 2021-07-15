# Unreleased

## Breaking Changes

### Removed Constants

1. StampKind.AseV1
1. StampKind.AseV2

### Removed Funcs

1. ACSChatThreadCreatedWithUserEventData.MarshalJSON() ([]byte, error)
1. ACSChatThreadPropertiesUpdatedPerUserEventData.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. ACSChatEventBaseProperties
1. ACSChatMemberAddedToThreadWithUserEventData
1. ACSChatMemberRemovedFromThreadWithUserEventData
1. ACSChatMessageDeletedEventData
1. ACSChatMessageEditedEventData
1. ACSChatMessageEventBaseProperties
1. ACSChatMessageReceivedEventData
1. ACSChatThreadCreatedWithUserEventData
1. ACSChatThreadEventBaseProperties
1. ACSChatThreadMemberProperties
1. ACSChatThreadPropertiesUpdatedPerUserEventData
1. ACSChatThreadWithUserDeletedEventData

### Signature Changes

#### Const Types

1. Public changed type from StampKind to CommunicationCloudEnvironmentModel

## Additive Changes

### New Constants

1. CommunicationCloudEnvironmentModel.Dod
1. CommunicationCloudEnvironmentModel.Gcch
1. StampKind.StampKindAseV1
1. StampKind.StampKindAseV2
1. StampKind.StampKindPublic

### New Funcs

1. AcsChatMessageEditedEventData.MarshalJSON() ([]byte, error)
1. AcsChatMessageEditedInThreadEventData.MarshalJSON() ([]byte, error)
1. AcsChatMessageReceivedEventData.MarshalJSON() ([]byte, error)
1. AcsChatMessageReceivedInThreadEventData.MarshalJSON() ([]byte, error)
1. AcsChatThreadCreatedEventData.MarshalJSON() ([]byte, error)
1. AcsChatThreadCreatedWithUserEventData.MarshalJSON() ([]byte, error)
1. AcsChatThreadPropertiesUpdatedEventData.MarshalJSON() ([]byte, error)
1. AcsChatThreadPropertiesUpdatedPerUserEventData.MarshalJSON() ([]byte, error)
1. PossibleCommunicationCloudEnvironmentModelValues() []CommunicationCloudEnvironmentModel

### Struct Changes

#### New Structs

1. AcsChatEventBaseProperties
1. AcsChatEventInThreadBaseProperties
1. AcsChatMessageDeletedEventData
1. AcsChatMessageDeletedInThreadEventData
1. AcsChatMessageEditedEventData
1. AcsChatMessageEditedInThreadEventData
1. AcsChatMessageEventBaseProperties
1. AcsChatMessageEventInThreadBaseProperties
1. AcsChatMessageReceivedEventData
1. AcsChatMessageReceivedInThreadEventData
1. AcsChatParticipantAddedToThreadEventData
1. AcsChatParticipantAddedToThreadWithUserEventData
1. AcsChatParticipantRemovedFromThreadEventData
1. AcsChatParticipantRemovedFromThreadWithUserEventData
1. AcsChatThreadCreatedEventData
1. AcsChatThreadCreatedWithUserEventData
1. AcsChatThreadDeletedEventData
1. AcsChatThreadEventBaseProperties
1. AcsChatThreadEventInThreadBaseProperties
1. AcsChatThreadParticipantProperties
1. AcsChatThreadPropertiesUpdatedEventData
1. AcsChatThreadPropertiesUpdatedPerUserEventData
1. AcsChatThreadWithUserDeletedEventData
1. AcsRecordingChunkInfoProperties
1. AcsRecordingFileStatusUpdatedEventData
1. AcsRecordingStorageInfoProperties
1. CommunicationIdentifierModel
1. CommunicationUserIdentifierModel
1. ContainerServiceNewKubernetesVersionAvailableEventData
1. MicrosoftTeamsUserIdentifierModel
1. PhoneNumberIdentifierModel
1. PolicyInsightsPolicyStateChangedEventData
1. PolicyInsightsPolicyStateCreatedEventData
1. PolicyInsightsPolicyStateDeletedEventData
1. ServiceBusActiveMessagesAvailablePeriodicNotificationsEventData
1. ServiceBusDeadletterMessagesAvailablePeriodicNotificationsEventData
1. StorageAsyncOperationInitiatedEventData
1. StorageBlobInventoryPolicyCompletedEventData
1. StorageBlobTierChangedEventData

#### New Struct Fields

1. AcsSmsDeliveryReportReceivedEventData.Tag
1. AppConfigurationKeyValueDeletedEventData.SyncToken
1. AppConfigurationKeyValueModifiedEventData.SyncToken

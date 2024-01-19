//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package systemevents

// Type represents the value set in EventData.EventType or messaging.CloudEvent.Type
// for system events.
type Type string

const (
	TypeAPIManagementAPICreated                                       Type = "Microsoft.ApiManagement.APICreated"                                        // maps to APIManagementAPICreatedEventData
	TypeAPIManagementAPIDeleted                                       Type = "Microsoft.ApiManagement.APIDeleted"                                        // maps to APIManagementAPIDeletedEventData
	TypeAPIManagementAPIReleaseCreated                                Type = "Microsoft.ApiManagement.APIReleaseCreated"                                 // maps to APIManagementAPIReleaseCreatedEventData
	TypeAPIManagementAPIReleaseDeleted                                Type = "Microsoft.ApiManagement.APIReleaseDeleted"                                 // maps to APIManagementAPIReleaseDeletedEventData
	TypeAPIManagementAPIReleaseUpdated                                Type = "Microsoft.ApiManagement.APIReleaseUpdated"                                 // maps to APIManagementAPIReleaseUpdatedEventData
	TypeAPIManagementAPIUpdated                                       Type = "Microsoft.ApiManagement.APIUpdated"                                        // maps to APIManagementAPIUpdatedEventData
	TypeAPIManagementGatewayAPIAdded                                  Type = "Microsoft.ApiManagement.GatewayAPIAdded"                                   // maps to APIManagementGatewayAPIAddedEventData
	TypeAPIManagementGatewayAPIRemoved                                Type = "Microsoft.ApiManagement.GatewayAPIRemoved"                                 // maps to APIManagementGatewayAPIRemovedEventData
	TypeAPIManagementGatewayCertificateAuthorityCreated               Type = "Microsoft.ApiManagement.GatewayCertificateAuthorityCreated"                // maps to APIManagementGatewayCertificateAuthorityCreatedEventData
	TypeAPIManagementGatewayCertificateAuthorityDeleted               Type = "Microsoft.ApiManagement.GatewayCertificateAuthorityDeleted"                // maps to APIManagementGatewayCertificateAuthorityDeletedEventData
	TypeAPIManagementGatewayCertificateAuthorityUpdated               Type = "Microsoft.ApiManagement.GatewayCertificateAuthorityUpdated"                // maps to APIManagementGatewayCertificateAuthorityUpdatedEventData
	TypeAPIManagementGatewayCreated                                   Type = "Microsoft.ApiManagement.GatewayCreated"                                    // maps to APIManagementGatewayCreatedEventData
	TypeAPIManagementGatewayDeleted                                   Type = "Microsoft.ApiManagement.GatewayDeleted"                                    // maps to APIManagementGatewayDeletedEventData
	TypeAPIManagementGatewayHostnameConfigurationCreated              Type = "Microsoft.ApiManagement.GatewayHostnameConfigurationCreated"               // maps to APIManagementGatewayHostnameConfigurationCreatedEventData
	TypeAPIManagementGatewayHostnameConfigurationDeleted              Type = "Microsoft.ApiManagement.GatewayHostnameConfigurationDeleted"               // maps to APIManagementGatewayHostnameConfigurationDeletedEventData
	TypeAPIManagementGatewayHostnameConfigurationUpdated              Type = "Microsoft.ApiManagement.GatewayHostnameConfigurationUpdated"               // maps to APIManagementGatewayHostnameConfigurationUpdatedEventData
	TypeAPIManagementGatewayUpdated                                   Type = "Microsoft.ApiManagement.GatewayUpdated"                                    // maps to APIManagementGatewayUpdatedEventData
	TypeAPIManagementProductCreated                                   Type = "Microsoft.ApiManagement.ProductCreated"                                    // maps to APIManagementProductCreatedEventData
	TypeAPIManagementProductDeleted                                   Type = "Microsoft.ApiManagement.ProductDeleted"                                    // maps to APIManagementProductDeletedEventData
	TypeAPIManagementProductUpdated                                   Type = "Microsoft.ApiManagement.ProductUpdated"                                    // maps to APIManagementProductUpdatedEventData
	TypeAPIManagementSubscriptionCreated                              Type = "Microsoft.ApiManagement.SubscriptionCreated"                               // maps to APIManagementSubscriptionCreatedEventData
	TypeAPIManagementSubscriptionDeleted                              Type = "Microsoft.ApiManagement.SubscriptionDeleted"                               // maps to APIManagementSubscriptionDeletedEventData
	TypeAPIManagementSubscriptionUpdated                              Type = "Microsoft.ApiManagement.SubscriptionUpdated"                               // maps to APIManagementSubscriptionUpdatedEventData
	TypeAPIManagementUserCreated                                      Type = "Microsoft.ApiManagement.UserCreated"                                       // maps to APIManagementUserCreatedEventData
	TypeAPIManagementUserDeleted                                      Type = "Microsoft.ApiManagement.UserDeleted"                                       // maps to APIManagementUserDeletedEventData
	TypeAPIManagementUserUpdated                                      Type = "Microsoft.ApiManagement.UserUpdated"                                       // maps to APIManagementUserUpdatedEventData
	TypeAppConfigurationKeyValueDeleted                               Type = "Microsoft.AppConfiguration.KeyValueDeleted"                                // maps to AppConfigurationKeyValueDeletedEventData
	TypeAppConfigurationKeyValueModified                              Type = "Microsoft.AppConfiguration.KeyValueModified"                               // maps to AppConfigurationKeyValueModifiedEventData
	TypeAppConfigurationSnapshotCreated                               Type = "Microsoft.AppConfiguration.SnapshotCreated"                                // maps to AppConfigurationSnapshotCreatedEventData
	TypeAppConfigurationSnapshotModified                              Type = "Microsoft.AppConfiguration.SnapshotModified"                               // maps to AppConfigurationSnapshotModifiedEventData
	TypeRedisExportRDBCompleted                                       Type = "Microsoft.Cache.ExportRDBCompleted"                                        // maps to RedisExportRDBCompletedEventData
	TypeRedisImportRDBCompleted                                       Type = "Microsoft.Cache.ImportRDBCompleted"                                        // maps to RedisImportRDBCompletedEventData
	TypeRedisPatchingCompleted                                        Type = "Microsoft.Cache.PatchingCompleted"                                         // maps to RedisPatchingCompletedEventData
	TypeRedisScalingCompleted                                         Type = "Microsoft.Cache.ScalingCompleted"                                          // maps to RedisScalingCompletedEventData
	TypeAcsChatMessageDeleted                                         Type = "Microsoft.Communication.ChatMessageDeleted"                                // maps to AcsChatMessageDeletedEventData
	TypeAcsChatMessageDeletedInThread                                 Type = "Microsoft.Communication.ChatMessageDeletedInThread"                        // maps to AcsChatMessageDeletedInThreadEventData
	TypeAcsChatMessageEdited                                          Type = "Microsoft.Communication.ChatMessageEdited"                                 // maps to AcsChatMessageEditedEventData
	TypeAcsChatMessageEditedInThread                                  Type = "Microsoft.Communication.ChatMessageEditedInThread"                         // maps to AcsChatMessageEditedInThreadEventData
	TypeAcsChatMessageReceived                                        Type = "Microsoft.Communication.ChatMessageReceived"                               // maps to AcsChatMessageReceivedEventData
	TypeAcsChatMessageReceivedInThread                                Type = "Microsoft.Communication.ChatMessageReceivedInThread"                       // maps to AcsChatMessageReceivedInThreadEventData
	TypeAcsChatParticipantAddedToThreadWithUser                       Type = "Microsoft.Communication.ChatParticipantAddedToThreadWithUser"              // maps to AcsChatParticipantAddedToThreadWithUserEventData
	TypeAcsChatParticipantRemovedFromThreadWithUser                   Type = "Microsoft.Communication.ChatParticipantRemovedFromThreadWithUser"          // maps to AcsChatParticipantRemovedFromThreadWithUserEventData
	TypeAcsChatThreadCreated                                          Type = "Microsoft.Communication.ChatThreadCreated"                                 // maps to AcsChatThreadCreatedEventData
	TypeAcsChatThreadCreatedWithUser                                  Type = "Microsoft.Communication.ChatThreadCreatedWithUser"                         // maps to AcsChatThreadCreatedWithUserEventData
	TypeAcsChatThreadDeleted                                          Type = "Microsoft.Communication.ChatThreadDeleted"                                 // maps to AcsChatThreadDeletedEventData
	TypeAcsChatParticipantAddedToThread                               Type = "Microsoft.Communication.ChatThreadParticipantAdded"                        // maps to AcsChatParticipantAddedToThreadEventData
	TypeAcsChatParticipantRemovedFromThread                           Type = "Microsoft.Communication.ChatThreadParticipantRemoved"                      // maps to AcsChatParticipantRemovedFromThreadEventData
	TypeAcsChatThreadPropertiesUpdated                                Type = "Microsoft.Communication.ChatThreadPropertiesUpdated"                       // maps to AcsChatThreadPropertiesUpdatedEventData
	TypeAcsChatThreadPropertiesUpdatedPerUser                         Type = "Microsoft.Communication.ChatThreadPropertiesUpdatedPerUser"                // maps to AcsChatThreadPropertiesUpdatedPerUserEventData
	TypeAcsChatThreadWithUserDeleted                                  Type = "Microsoft.Communication.ChatThreadWithUserDeleted"                         // maps to AcsChatThreadWithUserDeletedEventData
	TypeAcsEmailDeliveryReportReceived                                Type = "Microsoft.Communication.EmailDeliveryReportReceived"                       // maps to AcsEmailDeliveryReportReceivedEventData
	TypeAcsEmailEngagementTrackingReportReceived                      Type = "Microsoft.Communication.EmailEngagementTrackingReportReceived"             // maps to AcsEmailEngagementTrackingReportReceivedEventData
	TypeAcsIncomingCall                                               Type = "Microsoft.Communication.IncomingCall"                                      // maps to AcsIncomingCallEventData
	TypeAcsRecordingFileStatusUpdated                                 Type = "Microsoft.Communication.RecordingFileStatusUpdated"                        // maps to AcsRecordingFileStatusUpdatedEventData
	TypeAcsRouterJobCancelled                                         Type = "Microsoft.Communication.RouterJobCancelled"                                // maps to AcsRouterJobCancelledEventData
	TypeAcsRouterJobClassificationFailed                              Type = "Microsoft.Communication.RouterJobClassificationFailed"                     // maps to AcsRouterJobClassificationFailedEventData
	TypeAcsRouterJobClassified                                        Type = "Microsoft.Communication.RouterJobClassified"                               // maps to AcsRouterJobClassifiedEventData
	TypeAcsRouterJobClosed                                            Type = "Microsoft.Communication.RouterJobClosed"                                   // maps to AcsRouterJobClosedEventData
	TypeAcsRouterJobCompleted                                         Type = "Microsoft.Communication.RouterJobCompleted"                                // maps to AcsRouterJobCompletedEventData
	TypeAcsRouterJobDeleted                                           Type = "Microsoft.Communication.RouterJobDeleted"                                  // maps to AcsRouterJobDeletedEventData
	TypeAcsRouterJobExceptionTriggered                                Type = "Microsoft.Communication.RouterJobExceptionTriggered"                       // maps to AcsRouterJobExceptionTriggeredEventData
	TypeAcsRouterJobQueued                                            Type = "Microsoft.Communication.RouterJobQueued"                                   // maps to AcsRouterJobQueuedEventData
	TypeAcsRouterJobReceived                                          Type = "Microsoft.Communication.RouterJobReceived"                                 // maps to AcsRouterJobReceivedEventData
	TypeAcsRouterJobSchedulingFailed                                  Type = "Microsoft.Communication.RouterJobSchedulingFailed"                         // maps to AcsRouterJobSchedulingFailedEventData
	TypeAcsRouterJobUnassigned                                        Type = "Microsoft.Communication.RouterJobUnassigned"                               // maps to AcsRouterJobUnassignedEventData
	TypeAcsRouterJobWaitingForActivation                              Type = "Microsoft.Communication.RouterJobWaitingForActivation"                     // maps to AcsRouterJobWaitingForActivationEventData
	TypeAcsRouterJobWorkerSelectorsExpired                            Type = "Microsoft.Communication.RouterJobWorkerSelectorsExpired"                   // maps to AcsRouterJobWorkerSelectorsExpiredEventData
	TypeAcsRouterWorkerDeleted                                        Type = "Microsoft.Communication.RouterWorkerDeleted"                               // maps to AcsRouterWorkerDeletedEventData
	TypeAcsRouterWorkerDeregistered                                   Type = "Microsoft.Communication.RouterWorkerDeregistered"                          // maps to AcsRouterWorkerDeregisteredEventData
	TypeAcsRouterWorkerOfferAccepted                                  Type = "Microsoft.Communication.RouterWorkerOfferAccepted"                         // maps to AcsRouterWorkerOfferAcceptedEventData
	TypeAcsRouterWorkerOfferDeclined                                  Type = "Microsoft.Communication.RouterWorkerOfferDeclined"                         // maps to AcsRouterWorkerOfferDeclinedEventData
	TypeAcsRouterWorkerOfferExpired                                   Type = "Microsoft.Communication.RouterWorkerOfferExpired"                          // maps to AcsRouterWorkerOfferExpiredEventData
	TypeAcsRouterWorkerOfferIssued                                    Type = "Microsoft.Communication.RouterWorkerOfferIssued"                           // maps to AcsRouterWorkerOfferIssuedEventData
	TypeAcsRouterWorkerOfferRevoked                                   Type = "Microsoft.Communication.RouterWorkerOfferRevoked"                          // maps to AcsRouterWorkerOfferRevokedEventData
	TypeAcsRouterWorkerRegistered                                     Type = "Microsoft.Communication.RouterWorkerRegistered"                            // maps to AcsRouterWorkerRegisteredEventData
	TypeAcsSmsDeliveryReportReceived                                  Type = "Microsoft.Communication.SMSDeliveryReportReceived"                         // maps to AcsSmsDeliveryReportReceivedEventData
	TypeAcsSmsReceived                                                Type = "Microsoft.Communication.SMSReceived"                                       // maps to AcsSmsReceivedEventData
	TypeAcsUserDisconnected                                           Type = "Microsoft.Communication.UserDisconnected"                                  // maps to AcsUserDisconnectedEventData
	TypeContainerRegistryChartDeleted                                 Type = "Microsoft.ContainerRegistry.ChartDeleted"                                  // maps to ContainerRegistryChartDeletedEventData
	TypeContainerRegistryChartPushed                                  Type = "Microsoft.ContainerRegistry.ChartPushed"                                   // maps to ContainerRegistryChartPushedEventData
	TypeContainerRegistryImageDeleted                                 Type = "Microsoft.ContainerRegistry.ImageDeleted"                                  // maps to ContainerRegistryImageDeletedEventData
	TypeContainerRegistryImagePushed                                  Type = "Microsoft.ContainerRegistry.ImagePushed"                                   // maps to ContainerRegistryImagePushedEventData
	TypeContainerServiceClusterSupportEnded                           Type = "Microsoft.ContainerService.ClusterSupportEnded"                            // maps to ContainerServiceClusterSupportEndedEventData
	TypeContainerServiceClusterSupportEnding                          Type = "Microsoft.ContainerService.ClusterSupportEnding"                           // maps to ContainerServiceClusterSupportEndingEventData
	TypeContainerServiceNewKubernetesVersionAvailable                 Type = "Microsoft.ContainerService.NewKubernetesVersionAvailable"                  // maps to ContainerServiceNewKubernetesVersionAvailableEventData
	TypeContainerServiceNodePoolRollingFailed                         Type = "Microsoft.ContainerService.NodePoolRollingFailed"                          // maps to ContainerServiceNodePoolRollingFailedEventData
	TypeContainerServiceNodePoolRollingStarted                        Type = "Microsoft.ContainerService.NodePoolRollingStarted"                         // maps to ContainerServiceNodePoolRollingStartedEventData
	TypeContainerServiceNodePoolRollingSucceeded                      Type = "Microsoft.ContainerService.NodePoolRollingSucceeded"                       // maps to ContainerServiceNodePoolRollingSucceededEventData
	TypeDataBoxCopyCompleted                                          Type = "Microsoft.DataBox.CopyCompleted"                                           // maps to DataBoxCopyCompletedEventData
	TypeDataBoxCopyStarted                                            Type = "Microsoft.DataBox.CopyStarted"                                             // maps to DataBoxCopyStartedEventData
	TypeDataBoxOrderCompleted                                         Type = "Microsoft.DataBox.OrderCompleted"                                          // maps to DataBoxOrderCompletedEventData
	TypeIotHubDeviceConnected                                         Type = "Microsoft.Devices.DeviceConnected"                                         // maps to IotHubDeviceConnectedEventData
	TypeIotHubDeviceCreated                                           Type = "Microsoft.Devices.DeviceCreated"                                           // maps to IotHubDeviceCreatedEventData
	TypeIotHubDeviceDeleted                                           Type = "Microsoft.Devices.DeviceDeleted"                                           // maps to IotHubDeviceDeletedEventData
	TypeIotHubDeviceDisconnected                                      Type = "Microsoft.Devices.DeviceDisconnected"                                      // maps to IotHubDeviceDisconnectedEventData
	TypeIotHubDeviceTelemetry                                         Type = "Microsoft.Devices.DeviceTelemetry"                                         // maps to IotHubDeviceTelemetryEventData
	TypeEventGridMQTTClientCreatedOrUpdated                           Type = "Microsoft.EventGrid.MQTTClientCreatedOrUpdated"                            // maps to EventGridMQTTClientCreatedOrUpdatedEventData
	TypeEventGridMQTTClientDeleted                                    Type = "Microsoft.EventGrid.MQTTClientDeleted"                                     // maps to EventGridMQTTClientDeletedEventData
	TypeEventGridMQTTClientSessionConnected                           Type = "Microsoft.EventGrid.MQTTClientSessionConnected"                            // maps to EventGridMQTTClientSessionConnectedEventData
	TypeEventGridMQTTClientSessionDisconnected                        Type = "Microsoft.EventGrid.MQTTClientSessionDisconnected"                         // maps to EventGridMQTTClientSessionDisconnectedEventData
	TypeSubscriptionDeleted                                           Type = "Microsoft.EventGrid.SubscriptionDeletedEvent"                              // maps to SubscriptionDeletedEventData
	TypeSubscriptionValidation                                        Type = "Microsoft.EventGrid.SubscriptionValidationEvent"                           // maps to SubscriptionValidationEventData
	TypeEventHubCaptureFileCreated                                    Type = "Microsoft.EventHub.CaptureFileCreated"                                     // maps to EventHubCaptureFileCreatedEventData
	TypeHealthcareDicomImageCreated                                   Type = "Microsoft.HealthcareApis.DicomImageCreated"                                // maps to HealthcareDicomImageCreatedEventData
	TypeHealthcareDicomImageDeleted                                   Type = "Microsoft.HealthcareApis.DicomImageDeleted"                                // maps to HealthcareDicomImageDeletedEventData
	TypeHealthcareDicomImageUpdated                                   Type = "Microsoft.HealthcareApis.DicomImageUpdated"                                // maps to HealthcareDicomImageUpdatedEventData
	TypeHealthcareFhirResourceCreated                                 Type = "Microsoft.HealthcareApis.FhirResourceCreated"                              // maps to HealthcareFhirResourceCreatedEventData
	TypeHealthcareFhirResourceDeleted                                 Type = "Microsoft.HealthcareApis.FhirResourceDeleted"                              // maps to HealthcareFhirResourceDeletedEventData
	TypeHealthcareFhirResourceUpdated                                 Type = "Microsoft.HealthcareApis.FhirResourceUpdated"                              // maps to HealthcareFhirResourceUpdatedEventData
	TypeKeyVaultCertificateExpired                                    Type = "Microsoft.KeyVault.CertificateExpired"                                     // maps to KeyVaultCertificateExpiredEventData
	TypeKeyVaultCertificateNearExpiry                                 Type = "Microsoft.KeyVault.CertificateNearExpiry"                                  // maps to KeyVaultCertificateNearExpiryEventData
	TypeKeyVaultCertificateNewVersionCreated                          Type = "Microsoft.KeyVault.CertificateNewVersionCreated"                           // maps to KeyVaultCertificateNewVersionCreatedEventData
	TypeKeyVaultKeyExpired                                            Type = "Microsoft.KeyVault.KeyExpired"                                             // maps to KeyVaultKeyExpiredEventData
	TypeKeyVaultKeyNearExpiry                                         Type = "Microsoft.KeyVault.KeyNearExpiry"                                          // maps to KeyVaultKeyNearExpiryEventData
	TypeKeyVaultKeyNewVersionCreated                                  Type = "Microsoft.KeyVault.KeyNewVersionCreated"                                   // maps to KeyVaultKeyNewVersionCreatedEventData
	TypeKeyVaultSecretExpired                                         Type = "Microsoft.KeyVault.SecretExpired"                                          // maps to KeyVaultSecretExpiredEventData
	TypeKeyVaultSecretNearExpiry                                      Type = "Microsoft.KeyVault.SecretNearExpiry"                                       // maps to KeyVaultSecretNearExpiryEventData
	TypeKeyVaultSecretNewVersionCreated                               Type = "Microsoft.KeyVault.SecretNewVersionCreated"                                // maps to KeyVaultSecretNewVersionCreatedEventData
	TypeKeyVaultAccessPolicyChanged                                   Type = "Microsoft.KeyVault.VaultAccessPolicyChanged"                               // maps to KeyVaultAccessPolicyChangedEventData
	TypeMachineLearningServicesDatasetDriftDetected                   Type = "Microsoft.MachineLearningServices.DatasetDriftDetected"                    // maps to MachineLearningServicesDatasetDriftDetectedEventData
	TypeMachineLearningServicesModelDeployed                          Type = "Microsoft.MachineLearningServices.ModelDeployed"                           // maps to MachineLearningServicesModelDeployedEventData
	TypeMachineLearningServicesModelRegistered                        Type = "Microsoft.MachineLearningServices.ModelRegistered"                         // maps to MachineLearningServicesModelRegisteredEventData
	TypeMachineLearningServicesRunCompleted                           Type = "Microsoft.MachineLearningServices.RunCompleted"                            // maps to MachineLearningServicesRunCompletedEventData
	TypeMachineLearningServicesRunStatusChanged                       Type = "Microsoft.MachineLearningServices.RunStatusChanged"                        // maps to MachineLearningServicesRunStatusChangedEventData
	TypeMapsGeofenceEntered                                           Type = "Microsoft.Maps.GeofenceEntered"                                            // maps to MapsGeofenceEnteredEventData
	TypeMapsGeofenceExited                                            Type = "Microsoft.Maps.GeofenceExited"                                             // maps to MapsGeofenceExitedEventData
	TypeMapsGeofenceResult                                            Type = "Microsoft.Maps.GeofenceResult"                                             // maps to MapsGeofenceResultEventData
	TypeMediaJobCanceled                                              Type = "Microsoft.Media.JobCanceled"                                               // maps to MediaJobCanceledEventData
	TypeMediaJobCanceling                                             Type = "Microsoft.Media.JobCanceling"                                              // maps to MediaJobCancelingEventData
	TypeMediaJobErrored                                               Type = "Microsoft.Media.JobErrored"                                                // maps to MediaJobErroredEventData
	TypeMediaJobFinished                                              Type = "Microsoft.Media.JobFinished"                                               // maps to MediaJobFinishedEventData
	TypeMediaJobOutputCanceled                                        Type = "Microsoft.Media.JobOutputCanceled"                                         // maps to MediaJobOutputCanceledEventData
	TypeMediaJobOutputCanceling                                       Type = "Microsoft.Media.JobOutputCanceling"                                        // maps to MediaJobOutputCancelingEventData
	TypeMediaJobOutputErrored                                         Type = "Microsoft.Media.JobOutputErrored"                                          // maps to MediaJobOutputErroredEventData
	TypeMediaJobOutputFinished                                        Type = "Microsoft.Media.JobOutputFinished"                                         // maps to MediaJobOutputFinishedEventData
	TypeMediaJobOutputProcessing                                      Type = "Microsoft.Media.JobOutputProcessing"                                       // maps to MediaJobOutputProcessingEventData
	TypeMediaJobOutputProgress                                        Type = "Microsoft.Media.JobOutputProgress"                                         // maps to MediaJobOutputProgressEventData
	TypeMediaJobOutputScheduled                                       Type = "Microsoft.Media.JobOutputScheduled"                                        // maps to MediaJobOutputScheduledEventData
	TypeMediaJobOutputStateChange                                     Type = "Microsoft.Media.JobOutputStateChange"                                      // maps to MediaJobOutputStateChangeEventData
	TypeMediaJobProcessing                                            Type = "Microsoft.Media.JobProcessing"                                             // maps to MediaJobProcessingEventData
	TypeMediaJobScheduled                                             Type = "Microsoft.Media.JobScheduled"                                              // maps to MediaJobScheduledEventData
	TypeMediaJobStateChange                                           Type = "Microsoft.Media.JobStateChange"                                            // maps to MediaJobStateChangeEventData
	TypeMediaLiveEventChannelArchiveHeartbeat                         Type = "Microsoft.Media.LiveEventChannelArchiveHeartbeat"                          // maps to MediaLiveEventChannelArchiveHeartbeatEventData
	TypeMediaLiveEventConnectionRejected                              Type = "Microsoft.Media.LiveEventConnectionRejected"                               // maps to MediaLiveEventConnectionRejectedEventData
	TypeMediaLiveEventEncoderConnected                                Type = "Microsoft.Media.LiveEventEncoderConnected"                                 // maps to MediaLiveEventEncoderConnectedEventData
	TypeMediaLiveEventEncoderDisconnected                             Type = "Microsoft.Media.LiveEventEncoderDisconnected"                              // maps to MediaLiveEventEncoderDisconnectedEventData
	TypeMediaLiveEventIncomingDataChunkDropped                        Type = "Microsoft.Media.LiveEventIncomingDataChunkDropped"                         // maps to MediaLiveEventIncomingDataChunkDroppedEventData
	TypeMediaLiveEventIncomingStreamReceived                          Type = "Microsoft.Media.LiveEventIncomingStreamReceived"                           // maps to MediaLiveEventIncomingStreamReceivedEventData
	TypeMediaLiveEventIncomingStreamsOutOfSync                        Type = "Microsoft.Media.LiveEventIncomingStreamsOutOfSync"                         // maps to MediaLiveEventIncomingStreamsOutOfSyncEventData
	TypeMediaLiveEventIncomingVideoStreamsOutOfSync                   Type = "Microsoft.Media.LiveEventIncomingVideoStreamsOutOfSync"                    // maps to MediaLiveEventIncomingVideoStreamsOutOfSyncEventData
	TypeMediaLiveEventIngestHeartbeat                                 Type = "Microsoft.Media.LiveEventIngestHeartbeat"                                  // maps to MediaLiveEventIngestHeartbeatEventData
	TypeMediaLiveEventTrackDiscontinuityDetected                      Type = "Microsoft.Media.LiveEventTrackDiscontinuityDetected"                       // maps to MediaLiveEventTrackDiscontinuityDetectedEventData
	TypePolicyInsightsPolicyStateChanged                              Type = "Microsoft.PolicyInsights.PolicyStateChanged"                               // maps to PolicyInsightsPolicyStateChangedEventData
	TypePolicyInsightsPolicyStateCreated                              Type = "Microsoft.PolicyInsights.PolicyStateCreated"                               // maps to PolicyInsightsPolicyStateCreatedEventData
	TypePolicyInsightsPolicyStateDeleted                              Type = "Microsoft.PolicyInsights.PolicyStateDeleted"                               // maps to PolicyInsightsPolicyStateDeletedEventData
	TypeResourceNotificationsHealthResourcesAvailabilityStatusChanged Type = "Microsoft.ResourceNotifications.HealthResources.AvailabilityStatusChanged" // maps to ResourceNotificationsHealthResourcesAvailabilityStatusChangedEventData
	TypeResourceNotificationsHealthResourcesAnnotated                 Type = "Microsoft.ResourceNotifications.HealthResources.ResourceAnnotated"         // maps to ResourceNotificationsHealthResourcesAnnotatedEventData
	TypeResourceNotificationsResourceManagementCreatedOrUpdated       Type = "Microsoft.ResourceNotifications.Resources.CreatedOrUpdated"                // maps to ResourceNotificationsResourceManagementCreatedOrUpdatedEventData
	TypeResourceNotificationsResourceManagementDeleted                Type = "Microsoft.ResourceNotifications.Resources.Deleted"                         // maps to ResourceNotificationsResourceManagementDeletedEventData
	TypeResourceActionCancel                                          Type = "Microsoft.Resources.ResourceActionCancel"                                  // maps to ResourceActionCancelEventData
	TypeResourceActionFailure                                         Type = "Microsoft.Resources.ResourceActionFailure"                                 // maps to ResourceActionFailureEventData
	TypeResourceActionSuccess                                         Type = "Microsoft.Resources.ResourceActionSuccess"                                 // maps to ResourceActionSuccessEventData
	TypeResourceDeleteCancel                                          Type = "Microsoft.Resources.ResourceDeleteCancel"                                  // maps to ResourceDeleteCancelEventData
	TypeResourceDeleteFailure                                         Type = "Microsoft.Resources.ResourceDeleteFailure"                                 // maps to ResourceDeleteFailureEventData
	TypeResourceDeleteSuccess                                         Type = "Microsoft.Resources.ResourceDeleteSuccess"                                 // maps to ResourceDeleteSuccessEventData
	TypeResourceWriteCancel                                           Type = "Microsoft.Resources.ResourceWriteCancel"                                   // maps to ResourceWriteCancelEventData
	TypeResourceWriteFailure                                          Type = "Microsoft.Resources.ResourceWriteFailure"                                  // maps to ResourceWriteFailureEventData
	TypeResourceWriteSuccess                                          Type = "Microsoft.Resources.ResourceWriteSuccess"                                  // maps to ResourceWriteSuccessEventData
	TypeServiceBusActiveMessagesAvailablePeriodicNotifications        Type = "Microsoft.ServiceBus.ActiveMessagesAvailablePeriodicNotifications"         // maps to ServiceBusActiveMessagesAvailablePeriodicNotificationsEventData
	TypeServiceBusActiveMessagesAvailableWithNoListeners              Type = "Microsoft.ServiceBus.ActiveMessagesAvailableWithNoListeners"               // maps to ServiceBusActiveMessagesAvailableWithNoListenersEventData
	TypeServiceBusDeadletterMessagesAvailablePeriodicNotifications    Type = "Microsoft.ServiceBus.DeadletterMessagesAvailablePeriodicNotifications"     // maps to ServiceBusDeadletterMessagesAvailablePeriodicNotificationsEventData
	TypeServiceBusDeadletterMessagesAvailableWithNoListeners          Type = "Microsoft.ServiceBus.DeadletterMessagesAvailableWithNoListeners"           // maps to ServiceBusDeadletterMessagesAvailableWithNoListenersEventData
	TypeSignalRServiceClientConnectionConnected                       Type = "Microsoft.SignalRService.ClientConnectionConnected"                        // maps to SignalRServiceClientConnectionConnectedEventData
	TypeSignalRServiceClientConnectionDisconnected                    Type = "Microsoft.SignalRService.ClientConnectionDisconnected"                     // maps to SignalRServiceClientConnectionDisconnectedEventData
	TypeStorageAsyncOperationInitiated                                Type = "Microsoft.Storage.AsyncOperationInitiated"                                 // maps to StorageAsyncOperationInitiatedEventData
	TypeStorageBlobCreated                                            Type = "Microsoft.Storage.BlobCreated"                                             // maps to StorageBlobCreatedEventData
	TypeStorageBlobDeleted                                            Type = "Microsoft.Storage.BlobDeleted"                                             // maps to StorageBlobDeletedEventData
	TypeStorageBlobInventoryPolicyCompleted                           Type = "Microsoft.Storage.BlobInventoryPolicyCompleted"                            // maps to StorageBlobInventoryPolicyCompletedEventData
	TypeStorageBlobRenamed                                            Type = "Microsoft.Storage.BlobRenamed"                                             // maps to StorageBlobRenamedEventData
	TypeStorageBlobTierChanged                                        Type = "Microsoft.Storage.BlobTierChanged"                                         // maps to StorageBlobTierChangedEventData
	TypeStorageDirectoryCreated                                       Type = "Microsoft.Storage.DirectoryCreated"                                        // maps to StorageDirectoryCreatedEventData
	TypeStorageDirectoryDeleted                                       Type = "Microsoft.Storage.DirectoryDeleted"                                        // maps to StorageDirectoryDeletedEventData
	TypeStorageDirectoryRenamed                                       Type = "Microsoft.Storage.DirectoryRenamed"                                        // maps to StorageDirectoryRenamedEventData
	TypeStorageLifecyclePolicyCompleted                               Type = "Microsoft.Storage.LifecyclePolicyCompleted"                                // maps to StorageLifecyclePolicyCompletedEventData
	TypeStorageTaskCompleted                                          Type = "Microsoft.Storage.StorageTaskCompleted"                                    // maps to StorageTaskCompletedEventData
	TypeStorageTaskQueued                                             Type = "Microsoft.Storage.StorageTaskQueued"                                       // maps to StorageTaskQueuedEventData
	TypeWebAppServicePlanUpdated                                      Type = "Microsoft.Web.AppServicePlanUpdated"                                       // maps to WebAppServicePlanUpdatedEventData
	TypeWebAppUpdated                                                 Type = "Microsoft.Web.AppUpdated"                                                  // maps to WebAppUpdatedEventData
	TypeWebBackupOperationCompleted                                   Type = "Microsoft.Web.BackupOperationCompleted"                                    // maps to WebBackupOperationCompletedEventData
	TypeWebBackupOperationFailed                                      Type = "Microsoft.Web.BackupOperationFailed"                                       // maps to WebBackupOperationFailedEventData
	TypeWebBackupOperationStarted                                     Type = "Microsoft.Web.BackupOperationStarted"                                      // maps to WebBackupOperationStartedEventData
	TypeWebRestoreOperationCompleted                                  Type = "Microsoft.Web.RestoreOperationCompleted"                                   // maps to WebRestoreOperationCompletedEventData
	TypeWebRestoreOperationFailed                                     Type = "Microsoft.Web.RestoreOperationFailed"                                      // maps to WebRestoreOperationFailedEventData
	TypeWebRestoreOperationStarted                                    Type = "Microsoft.Web.RestoreOperationStarted"                                     // maps to WebRestoreOperationStartedEventData
	TypeWebSlotSwapCompleted                                          Type = "Microsoft.Web.SlotSwapCompleted"                                           // maps to WebSlotSwapCompletedEventData
	TypeWebSlotSwapFailed                                             Type = "Microsoft.Web.SlotSwapFailed"                                              // maps to WebSlotSwapFailedEventData
	TypeWebSlotSwapStarted                                            Type = "Microsoft.Web.SlotSwapStarted"                                             // maps to WebSlotSwapStartedEventData
	TypeWebSlotSwapWithPreviewCancelled                               Type = "Microsoft.Web.SlotSwapWithPreviewCancelled"                                // maps to WebSlotSwapWithPreviewCancelledEventData
	TypeWebSlotSwapWithPreviewStarted                                 Type = "Microsoft.Web.SlotSwapWithPreviewStarted"                                  // maps to WebSlotSwapWithPreviewStartedEventData
)

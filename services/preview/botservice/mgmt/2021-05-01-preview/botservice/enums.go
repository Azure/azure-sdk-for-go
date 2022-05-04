package botservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// ChannelName enumerates the values for channel name.
type ChannelName string

const (
	// ChannelNameAlexaChannel ...
	ChannelNameAlexaChannel ChannelName = "AlexaChannel"
	// ChannelNameDirectLineChannel ...
	ChannelNameDirectLineChannel ChannelName = "DirectLineChannel"
	// ChannelNameDirectLineSpeechChannel ...
	ChannelNameDirectLineSpeechChannel ChannelName = "DirectLineSpeechChannel"
	// ChannelNameEmailChannel ...
	ChannelNameEmailChannel ChannelName = "EmailChannel"
	// ChannelNameFacebookChannel ...
	ChannelNameFacebookChannel ChannelName = "FacebookChannel"
	// ChannelNameKikChannel ...
	ChannelNameKikChannel ChannelName = "KikChannel"
	// ChannelNameLineChannel ...
	ChannelNameLineChannel ChannelName = "LineChannel"
	// ChannelNameMsTeamsChannel ...
	ChannelNameMsTeamsChannel ChannelName = "MsTeamsChannel"
	// ChannelNameOutlookChannel ...
	ChannelNameOutlookChannel ChannelName = "OutlookChannel"
	// ChannelNameSkypeChannel ...
	ChannelNameSkypeChannel ChannelName = "SkypeChannel"
	// ChannelNameSlackChannel ...
	ChannelNameSlackChannel ChannelName = "SlackChannel"
	// ChannelNameSmsChannel ...
	ChannelNameSmsChannel ChannelName = "SmsChannel"
	// ChannelNameTelegramChannel ...
	ChannelNameTelegramChannel ChannelName = "TelegramChannel"
	// ChannelNameWebChatChannel ...
	ChannelNameWebChatChannel ChannelName = "WebChatChannel"
)

// PossibleChannelNameValues returns an array of possible values for the ChannelName const type.
func PossibleChannelNameValues() []ChannelName {
	return []ChannelName{ChannelNameAlexaChannel, ChannelNameDirectLineChannel, ChannelNameDirectLineSpeechChannel, ChannelNameEmailChannel, ChannelNameFacebookChannel, ChannelNameKikChannel, ChannelNameLineChannel, ChannelNameMsTeamsChannel, ChannelNameOutlookChannel, ChannelNameSkypeChannel, ChannelNameSlackChannel, ChannelNameSmsChannel, ChannelNameTelegramChannel, ChannelNameWebChatChannel}
}

// ChannelNameBasicChannel enumerates the values for channel name basic channel.
type ChannelNameBasicChannel string

const (
	// ChannelNameBasicChannelChannelNameAlexaChannel ...
	ChannelNameBasicChannelChannelNameAlexaChannel ChannelNameBasicChannel = "AlexaChannel"
	// ChannelNameBasicChannelChannelNameChannel ...
	ChannelNameBasicChannelChannelNameChannel ChannelNameBasicChannel = "Channel"
	// ChannelNameBasicChannelChannelNameDirectLineChannel ...
	ChannelNameBasicChannelChannelNameDirectLineChannel ChannelNameBasicChannel = "DirectLineChannel"
	// ChannelNameBasicChannelChannelNameDirectLineSpeechChannel ...
	ChannelNameBasicChannelChannelNameDirectLineSpeechChannel ChannelNameBasicChannel = "DirectLineSpeechChannel"
	// ChannelNameBasicChannelChannelNameEmailChannel ...
	ChannelNameBasicChannelChannelNameEmailChannel ChannelNameBasicChannel = "EmailChannel"
	// ChannelNameBasicChannelChannelNameFacebookChannel ...
	ChannelNameBasicChannelChannelNameFacebookChannel ChannelNameBasicChannel = "FacebookChannel"
	// ChannelNameBasicChannelChannelNameKikChannel ...
	ChannelNameBasicChannelChannelNameKikChannel ChannelNameBasicChannel = "KikChannel"
	// ChannelNameBasicChannelChannelNameLineChannel ...
	ChannelNameBasicChannelChannelNameLineChannel ChannelNameBasicChannel = "LineChannel"
	// ChannelNameBasicChannelChannelNameMsTeamsChannel ...
	ChannelNameBasicChannelChannelNameMsTeamsChannel ChannelNameBasicChannel = "MsTeamsChannel"
	// ChannelNameBasicChannelChannelNameSkypeChannel ...
	ChannelNameBasicChannelChannelNameSkypeChannel ChannelNameBasicChannel = "SkypeChannel"
	// ChannelNameBasicChannelChannelNameSlackChannel ...
	ChannelNameBasicChannelChannelNameSlackChannel ChannelNameBasicChannel = "SlackChannel"
	// ChannelNameBasicChannelChannelNameSmsChannel ...
	ChannelNameBasicChannelChannelNameSmsChannel ChannelNameBasicChannel = "SmsChannel"
	// ChannelNameBasicChannelChannelNameTelegramChannel ...
	ChannelNameBasicChannelChannelNameTelegramChannel ChannelNameBasicChannel = "TelegramChannel"
	// ChannelNameBasicChannelChannelNameWebChatChannel ...
	ChannelNameBasicChannelChannelNameWebChatChannel ChannelNameBasicChannel = "WebChatChannel"
)

// PossibleChannelNameBasicChannelValues returns an array of possible values for the ChannelNameBasicChannel const type.
func PossibleChannelNameBasicChannelValues() []ChannelNameBasicChannel {
	return []ChannelNameBasicChannel{ChannelNameBasicChannelChannelNameAlexaChannel, ChannelNameBasicChannelChannelNameChannel, ChannelNameBasicChannelChannelNameDirectLineChannel, ChannelNameBasicChannelChannelNameDirectLineSpeechChannel, ChannelNameBasicChannelChannelNameEmailChannel, ChannelNameBasicChannelChannelNameFacebookChannel, ChannelNameBasicChannelChannelNameKikChannel, ChannelNameBasicChannelChannelNameLineChannel, ChannelNameBasicChannelChannelNameMsTeamsChannel, ChannelNameBasicChannelChannelNameSkypeChannel, ChannelNameBasicChannelChannelNameSlackChannel, ChannelNameBasicChannelChannelNameSmsChannel, ChannelNameBasicChannelChannelNameTelegramChannel, ChannelNameBasicChannelChannelNameWebChatChannel}
}

// Key enumerates the values for key.
type Key string

const (
	// Key1 ...
	Key1 Key = "key1"
	// Key2 ...
	Key2 Key = "key2"
)

// PossibleKeyValues returns an array of possible values for the Key const type.
func PossibleKeyValues() []Key {
	return []Key{Key1, Key2}
}

// Kind enumerates the values for kind.
type Kind string

const (
	// KindAzurebot ...
	KindAzurebot Kind = "azurebot"
	// KindBot ...
	KindBot Kind = "bot"
	// KindDesigner ...
	KindDesigner Kind = "designer"
	// KindFunction ...
	KindFunction Kind = "function"
	// KindSdk ...
	KindSdk Kind = "sdk"
)

// PossibleKindValues returns an array of possible values for the Kind const type.
func PossibleKindValues() []Kind {
	return []Kind{KindAzurebot, KindBot, KindDesigner, KindFunction, KindSdk}
}

// MsaAppType enumerates the values for msa app type.
type MsaAppType string

const (
	// MsaAppTypeMultiTenant ...
	MsaAppTypeMultiTenant MsaAppType = "MultiTenant"
	// MsaAppTypeSingleTenant ...
	MsaAppTypeSingleTenant MsaAppType = "SingleTenant"
	// MsaAppTypeUserAssignedMSI ...
	MsaAppTypeUserAssignedMSI MsaAppType = "UserAssignedMSI"
)

// PossibleMsaAppTypeValues returns an array of possible values for the MsaAppType const type.
func PossibleMsaAppTypeValues() []MsaAppType {
	return []MsaAppType{MsaAppTypeMultiTenant, MsaAppTypeSingleTenant, MsaAppTypeUserAssignedMSI}
}

// OperationResultStatus enumerates the values for operation result status.
type OperationResultStatus string

const (
	// OperationResultStatusCanceled ...
	OperationResultStatusCanceled OperationResultStatus = "Canceled"
	// OperationResultStatusFailed ...
	OperationResultStatusFailed OperationResultStatus = "Failed"
	// OperationResultStatusRequested ...
	OperationResultStatusRequested OperationResultStatus = "Requested"
	// OperationResultStatusRunning ...
	OperationResultStatusRunning OperationResultStatus = "Running"
	// OperationResultStatusSucceeded ...
	OperationResultStatusSucceeded OperationResultStatus = "Succeeded"
)

// PossibleOperationResultStatusValues returns an array of possible values for the OperationResultStatus const type.
func PossibleOperationResultStatusValues() []OperationResultStatus {
	return []OperationResultStatus{OperationResultStatusCanceled, OperationResultStatusFailed, OperationResultStatusRequested, OperationResultStatusRunning, OperationResultStatusSucceeded}
}

// PrivateEndpointConnectionProvisioningState enumerates the values for private endpoint connection
// provisioning state.
type PrivateEndpointConnectionProvisioningState string

const (
	// PrivateEndpointConnectionProvisioningStateCreating ...
	PrivateEndpointConnectionProvisioningStateCreating PrivateEndpointConnectionProvisioningState = "Creating"
	// PrivateEndpointConnectionProvisioningStateDeleting ...
	PrivateEndpointConnectionProvisioningStateDeleting PrivateEndpointConnectionProvisioningState = "Deleting"
	// PrivateEndpointConnectionProvisioningStateFailed ...
	PrivateEndpointConnectionProvisioningStateFailed PrivateEndpointConnectionProvisioningState = "Failed"
	// PrivateEndpointConnectionProvisioningStateSucceeded ...
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

// PossiblePrivateEndpointConnectionProvisioningStateValues returns an array of possible values for the PrivateEndpointConnectionProvisioningState const type.
func PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState {
	return []PrivateEndpointConnectionProvisioningState{PrivateEndpointConnectionProvisioningStateCreating, PrivateEndpointConnectionProvisioningStateDeleting, PrivateEndpointConnectionProvisioningStateFailed, PrivateEndpointConnectionProvisioningStateSucceeded}
}

// PrivateEndpointServiceConnectionStatus enumerates the values for private endpoint service connection status.
type PrivateEndpointServiceConnectionStatus string

const (
	// PrivateEndpointServiceConnectionStatusApproved ...
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	// PrivateEndpointServiceConnectionStatusPending ...
	PrivateEndpointServiceConnectionStatusPending PrivateEndpointServiceConnectionStatus = "Pending"
	// PrivateEndpointServiceConnectionStatusRejected ...
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

// PossiblePrivateEndpointServiceConnectionStatusValues returns an array of possible values for the PrivateEndpointServiceConnectionStatus const type.
func PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus {
	return []PrivateEndpointServiceConnectionStatus{PrivateEndpointServiceConnectionStatusApproved, PrivateEndpointServiceConnectionStatusPending, PrivateEndpointServiceConnectionStatusRejected}
}

// PublicNetworkAccess enumerates the values for public network access.
type PublicNetworkAccess string

const (
	// PublicNetworkAccessDisabled ...
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	// PublicNetworkAccessEnabled ...
	PublicNetworkAccessEnabled PublicNetworkAccess = "Enabled"
)

// PossiblePublicNetworkAccessValues returns an array of possible values for the PublicNetworkAccess const type.
func PossiblePublicNetworkAccessValues() []PublicNetworkAccess {
	return []PublicNetworkAccess{PublicNetworkAccessDisabled, PublicNetworkAccessEnabled}
}

// RegenerateKeysChannelName enumerates the values for regenerate keys channel name.
type RegenerateKeysChannelName string

const (
	// RegenerateKeysChannelNameDirectLineChannel ...
	RegenerateKeysChannelNameDirectLineChannel RegenerateKeysChannelName = "DirectLineChannel"
	// RegenerateKeysChannelNameWebChatChannel ...
	RegenerateKeysChannelNameWebChatChannel RegenerateKeysChannelName = "WebChatChannel"
)

// PossibleRegenerateKeysChannelNameValues returns an array of possible values for the RegenerateKeysChannelName const type.
func PossibleRegenerateKeysChannelNameValues() []RegenerateKeysChannelName {
	return []RegenerateKeysChannelName{RegenerateKeysChannelNameDirectLineChannel, RegenerateKeysChannelNameWebChatChannel}
}

// SkuName enumerates the values for sku name.
type SkuName string

const (
	// SkuNameF0 ...
	SkuNameF0 SkuName = "F0"
	// SkuNameS1 ...
	SkuNameS1 SkuName = "S1"
)

// PossibleSkuNameValues returns an array of possible values for the SkuName const type.
func PossibleSkuNameValues() []SkuName {
	return []SkuName{SkuNameF0, SkuNameS1}
}

// SkuTier enumerates the values for sku tier.
type SkuTier string

const (
	// SkuTierFree ...
	SkuTierFree SkuTier = "Free"
	// SkuTierStandard ...
	SkuTierStandard SkuTier = "Standard"
)

// PossibleSkuTierValues returns an array of possible values for the SkuTier const type.
func PossibleSkuTierValues() []SkuTier {
	return []SkuTier{SkuTierFree, SkuTierStandard}
}

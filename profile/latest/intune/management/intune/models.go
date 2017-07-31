package intune

import (
	 original "github.com/Azure/azure-sdk-for-go/service/intune/management/2015-01-14-privatepreview/intune"
)

type (
	 IosClient = original.IosClient
	 AppSharingFromLevel = original.AppSharingFromLevel
	 AppSharingToLevel = original.AppSharingToLevel
	 Authentication = original.Authentication
	 ClipboardSharingLevel = original.ClipboardSharingLevel
	 DataBackup = original.DataBackup
	 DeviceCompliance = original.DeviceCompliance
	 FileEncryption = original.FileEncryption
	 FileEncryptionLevel = original.FileEncryptionLevel
	 FileSharingSaveAs = original.FileSharingSaveAs
	 GroupStatus = original.GroupStatus
	 ManagedBrowser = original.ManagedBrowser
	 Pin = original.Pin
	 Platform = original.Platform
	 ScreenCapture = original.ScreenCapture
	 TouchID = original.TouchID
	 AndroidMAMPolicy = original.AndroidMAMPolicy
	 AndroidMAMPolicyCollection = original.AndroidMAMPolicyCollection
	 AndroidMAMPolicyProperties = original.AndroidMAMPolicyProperties
	 Application = original.Application
	 ApplicationCollection = original.ApplicationCollection
	 ApplicationProperties = original.ApplicationProperties
	 Device = original.Device
	 DeviceCollection = original.DeviceCollection
	 DeviceProperties = original.DeviceProperties
	 Error = original.Error
	 FlaggedEnrolledApp = original.FlaggedEnrolledApp
	 FlaggedEnrolledAppCollection = original.FlaggedEnrolledAppCollection
	 FlaggedEnrolledAppError = original.FlaggedEnrolledAppError
	 FlaggedEnrolledAppProperties = original.FlaggedEnrolledAppProperties
	 FlaggedUser = original.FlaggedUser
	 FlaggedUserCollection = original.FlaggedUserCollection
	 FlaggedUserProperties = original.FlaggedUserProperties
	 GroupItem = original.GroupItem
	 GroupProperties = original.GroupProperties
	 GroupsCollection = original.GroupsCollection
	 IOSMAMPolicy = original.IOSMAMPolicy
	 IOSMAMPolicyCollection = original.IOSMAMPolicyCollection
	 IOSMAMPolicyProperties = original.IOSMAMPolicyProperties
	 Location = original.Location
	 LocationCollection = original.LocationCollection
	 LocationProperties = original.LocationProperties
	 MAMPolicyAppIDOrGroupIDPayload = original.MAMPolicyAppIDOrGroupIDPayload
	 MAMPolicyAppOrGroupIDProperties = original.MAMPolicyAppOrGroupIDProperties
	 MAMPolicyProperties = original.MAMPolicyProperties
	 OperationMetadataProperties = original.OperationMetadataProperties
	 OperationResult = original.OperationResult
	 OperationResultCollection = original.OperationResultCollection
	 OperationResultProperties = original.OperationResultProperties
	 Resource = original.Resource
	 StatusesDefault = original.StatusesDefault
	 StatusesProperties = original.StatusesProperties
	 WipeDeviceOperationResult = original.WipeDeviceOperationResult
	 WipeDeviceOperationResultProperties = original.WipeDeviceOperationResultProperties
	 AndroidClient = original.AndroidClient
	 ManagementClient = original.ManagementClient
)

const (
	 AllApps = original.AllApps
	 None = original.None
	 PolicyManagedApps = original.PolicyManagedApps
	 AppSharingToLevelAllApps = original.AppSharingToLevelAllApps
	 AppSharingToLevelNone = original.AppSharingToLevelNone
	 AppSharingToLevelPolicyManagedApps = original.AppSharingToLevelPolicyManagedApps
	 NotRequired = original.NotRequired
	 Required = original.Required
	 ClipboardSharingLevelAllApps = original.ClipboardSharingLevelAllApps
	 ClipboardSharingLevelBlocked = original.ClipboardSharingLevelBlocked
	 ClipboardSharingLevelPolicyManagedApps = original.ClipboardSharingLevelPolicyManagedApps
	 ClipboardSharingLevelPolicyManagedAppsWithPasteIn = original.ClipboardSharingLevelPolicyManagedAppsWithPasteIn
	 Allow = original.Allow
	 Block = original.Block
	 Disable = original.Disable
	 Enable = original.Enable
	 FileEncryptionNotRequired = original.FileEncryptionNotRequired
	 FileEncryptionRequired = original.FileEncryptionRequired
	 AfterDeviceRestart = original.AfterDeviceRestart
	 DeviceLocked = original.DeviceLocked
	 DeviceLockedExceptFilesOpen = original.DeviceLockedExceptFilesOpen
	 UseDeviceSettings = original.UseDeviceSettings
	 FileSharingSaveAsAllow = original.FileSharingSaveAsAllow
	 FileSharingSaveAsBlock = original.FileSharingSaveAsBlock
	 NotTargeted = original.NotTargeted
	 Targeted = original.Targeted
	 ManagedBrowserNotRequired = original.ManagedBrowserNotRequired
	 ManagedBrowserRequired = original.ManagedBrowserRequired
	 PinNotRequired = original.PinNotRequired
	 PinRequired = original.PinRequired
	 Android = original.Android
	 Ios = original.Ios
	 Windows = original.Windows
	 ScreenCaptureAllow = original.ScreenCaptureAllow
	 ScreenCaptureBlock = original.ScreenCaptureBlock
	 TouchIDDisable = original.TouchIDDisable
	 TouchIDEnable = original.TouchIDEnable
	 DefaultBaseURI = original.DefaultBaseURI
)

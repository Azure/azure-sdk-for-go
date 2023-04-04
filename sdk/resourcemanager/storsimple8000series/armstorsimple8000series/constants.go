//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armstorsimple8000series

const (
	moduleName    = "armstorsimple8000series"
	moduleVersion = "v1.1.0"
)

// AlertEmailNotificationStatus - Indicates whether email notification enabled or not.
type AlertEmailNotificationStatus string

const (
	AlertEmailNotificationStatusEnabled  AlertEmailNotificationStatus = "Enabled"
	AlertEmailNotificationStatusDisabled AlertEmailNotificationStatus = "Disabled"
)

// PossibleAlertEmailNotificationStatusValues returns the possible values for the AlertEmailNotificationStatus const type.
func PossibleAlertEmailNotificationStatusValues() []AlertEmailNotificationStatus {
	return []AlertEmailNotificationStatus{
		AlertEmailNotificationStatusEnabled,
		AlertEmailNotificationStatusDisabled,
	}
}

// AlertScope - The scope of the alert
type AlertScope string

const (
	AlertScopeResource AlertScope = "Resource"
	AlertScopeDevice   AlertScope = "Device"
)

// PossibleAlertScopeValues returns the possible values for the AlertScope const type.
func PossibleAlertScopeValues() []AlertScope {
	return []AlertScope{
		AlertScopeResource,
		AlertScopeDevice,
	}
}

// AlertSeverity - Specifies the severity of the alerts to be filtered. Only 'Equality' operator is supported for this property.
type AlertSeverity string

const (
	AlertSeverityInformational AlertSeverity = "Informational"
	AlertSeverityWarning       AlertSeverity = "Warning"
	AlertSeverityCritical      AlertSeverity = "Critical"
)

// PossibleAlertSeverityValues returns the possible values for the AlertSeverity const type.
func PossibleAlertSeverityValues() []AlertSeverity {
	return []AlertSeverity{
		AlertSeverityInformational,
		AlertSeverityWarning,
		AlertSeverityCritical,
	}
}

// AlertSourceType - Specifies the source type of the alerts to be filtered. Only 'Equality' operator is supported for this
// property.
type AlertSourceType string

const (
	AlertSourceTypeResource AlertSourceType = "Resource"
	AlertSourceTypeDevice   AlertSourceType = "Device"
)

// PossibleAlertSourceTypeValues returns the possible values for the AlertSourceType const type.
func PossibleAlertSourceTypeValues() []AlertSourceType {
	return []AlertSourceType{
		AlertSourceTypeResource,
		AlertSourceTypeDevice,
	}
}

// AlertStatus - Specifies the status of the alerts to be filtered. Only 'Equality' operator is supported for this property.
type AlertStatus string

const (
	AlertStatusActive  AlertStatus = "Active"
	AlertStatusCleared AlertStatus = "Cleared"
)

// PossibleAlertStatusValues returns the possible values for the AlertStatus const type.
func PossibleAlertStatusValues() []AlertStatus {
	return []AlertStatus{
		AlertStatusActive,
		AlertStatusCleared,
	}
}

// AuthenticationType - The authentication type.
type AuthenticationType string

const (
	AuthenticationTypeInvalid AuthenticationType = "Invalid"
	AuthenticationTypeNone    AuthenticationType = "None"
	AuthenticationTypeBasic   AuthenticationType = "Basic"
	AuthenticationTypeNTLM    AuthenticationType = "NTLM"
)

// PossibleAuthenticationTypeValues returns the possible values for the AuthenticationType const type.
func PossibleAuthenticationTypeValues() []AuthenticationType {
	return []AuthenticationType{
		AuthenticationTypeInvalid,
		AuthenticationTypeNone,
		AuthenticationTypeBasic,
		AuthenticationTypeNTLM,
	}
}

// AuthorizationEligibility - The eligibility status of device for service data encryption key rollover.
type AuthorizationEligibility string

const (
	AuthorizationEligibilityInEligible AuthorizationEligibility = "InEligible"
	AuthorizationEligibilityEligible   AuthorizationEligibility = "Eligible"
)

// PossibleAuthorizationEligibilityValues returns the possible values for the AuthorizationEligibility const type.
func PossibleAuthorizationEligibilityValues() []AuthorizationEligibility {
	return []AuthorizationEligibility{
		AuthorizationEligibilityInEligible,
		AuthorizationEligibilityEligible,
	}
}

// AuthorizationStatus - The authorization status of the device for service data encryption key rollover.
type AuthorizationStatus string

const (
	AuthorizationStatusDisabled AuthorizationStatus = "Disabled"
	AuthorizationStatusEnabled  AuthorizationStatus = "Enabled"
)

// PossibleAuthorizationStatusValues returns the possible values for the AuthorizationStatus const type.
func PossibleAuthorizationStatusValues() []AuthorizationStatus {
	return []AuthorizationStatus{
		AuthorizationStatusDisabled,
		AuthorizationStatusEnabled,
	}
}

// BackupJobCreationType - The backup job creation type.
type BackupJobCreationType string

const (
	BackupJobCreationTypeAdhoc      BackupJobCreationType = "Adhoc"
	BackupJobCreationTypeBySchedule BackupJobCreationType = "BySchedule"
	BackupJobCreationTypeBySSM      BackupJobCreationType = "BySSM"
)

// PossibleBackupJobCreationTypeValues returns the possible values for the BackupJobCreationType const type.
func PossibleBackupJobCreationTypeValues() []BackupJobCreationType {
	return []BackupJobCreationType{
		BackupJobCreationTypeAdhoc,
		BackupJobCreationTypeBySchedule,
		BackupJobCreationTypeBySSM,
	}
}

// BackupPolicyCreationType - The backup policy creation type. Indicates whether this was created through SaaS or through
// StorSimple Snapshot Manager.
type BackupPolicyCreationType string

const (
	BackupPolicyCreationTypeBySaaS BackupPolicyCreationType = "BySaaS"
	BackupPolicyCreationTypeBySSM  BackupPolicyCreationType = "BySSM"
)

// PossibleBackupPolicyCreationTypeValues returns the possible values for the BackupPolicyCreationType const type.
func PossibleBackupPolicyCreationTypeValues() []BackupPolicyCreationType {
	return []BackupPolicyCreationType{
		BackupPolicyCreationTypeBySaaS,
		BackupPolicyCreationTypeBySSM,
	}
}

// BackupStatus - The backup status of the volume.
type BackupStatus string

const (
	BackupStatusEnabled  BackupStatus = "Enabled"
	BackupStatusDisabled BackupStatus = "Disabled"
)

// PossibleBackupStatusValues returns the possible values for the BackupStatus const type.
func PossibleBackupStatusValues() []BackupStatus {
	return []BackupStatus{
		BackupStatusEnabled,
		BackupStatusDisabled,
	}
}

// BackupType - The type of the backup.
type BackupType string

const (
	BackupTypeLocalSnapshot BackupType = "LocalSnapshot"
	BackupTypeCloudSnapshot BackupType = "CloudSnapshot"
)

// PossibleBackupTypeValues returns the possible values for the BackupType const type.
func PossibleBackupTypeValues() []BackupType {
	return []BackupType{
		BackupTypeLocalSnapshot,
		BackupTypeCloudSnapshot,
	}
}

// ControllerID - The active controller that the request is expecting on the device.
type ControllerID string

const (
	ControllerIDUnknown     ControllerID = "Unknown"
	ControllerIDNone        ControllerID = "None"
	ControllerIDController0 ControllerID = "Controller0"
	ControllerIDController1 ControllerID = "Controller1"
)

// PossibleControllerIDValues returns the possible values for the ControllerID const type.
func PossibleControllerIDValues() []ControllerID {
	return []ControllerID{
		ControllerIDUnknown,
		ControllerIDNone,
		ControllerIDController0,
		ControllerIDController1,
	}
}

// ControllerPowerStateAction - The power state that the request is expecting for the controller of the device.
type ControllerPowerStateAction string

const (
	ControllerPowerStateActionStart    ControllerPowerStateAction = "Start"
	ControllerPowerStateActionRestart  ControllerPowerStateAction = "Restart"
	ControllerPowerStateActionShutdown ControllerPowerStateAction = "Shutdown"
)

// PossibleControllerPowerStateActionValues returns the possible values for the ControllerPowerStateAction const type.
func PossibleControllerPowerStateActionValues() []ControllerPowerStateAction {
	return []ControllerPowerStateAction{
		ControllerPowerStateActionStart,
		ControllerPowerStateActionRestart,
		ControllerPowerStateActionShutdown,
	}
}

// ControllerStatus - The controller 0's status that the request is expecting on the device.
type ControllerStatus string

const (
	ControllerStatusNotPresent ControllerStatus = "NotPresent"
	ControllerStatusPoweredOff ControllerStatus = "PoweredOff"
	ControllerStatusOk         ControllerStatus = "Ok"
	ControllerStatusRecovering ControllerStatus = "Recovering"
	ControllerStatusWarning    ControllerStatus = "Warning"
	ControllerStatusFailure    ControllerStatus = "Failure"
)

// PossibleControllerStatusValues returns the possible values for the ControllerStatus const type.
func PossibleControllerStatusValues() []ControllerStatus {
	return []ControllerStatus{
		ControllerStatusNotPresent,
		ControllerStatusPoweredOff,
		ControllerStatusOk,
		ControllerStatusRecovering,
		ControllerStatusWarning,
		ControllerStatusFailure,
	}
}

type DayOfWeek string

const (
	DayOfWeekSunday    DayOfWeek = "Sunday"
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
)

// PossibleDayOfWeekValues returns the possible values for the DayOfWeek const type.
func PossibleDayOfWeekValues() []DayOfWeek {
	return []DayOfWeek{
		DayOfWeekSunday,
		DayOfWeekMonday,
		DayOfWeekTuesday,
		DayOfWeekWednesday,
		DayOfWeekThursday,
		DayOfWeekFriday,
		DayOfWeekSaturday,
	}
}

// DeviceConfigurationStatus - The current configuration status of the device.
type DeviceConfigurationStatus string

const (
	DeviceConfigurationStatusComplete DeviceConfigurationStatus = "Complete"
	DeviceConfigurationStatusPending  DeviceConfigurationStatus = "Pending"
)

// PossibleDeviceConfigurationStatusValues returns the possible values for the DeviceConfigurationStatus const type.
func PossibleDeviceConfigurationStatusValues() []DeviceConfigurationStatus {
	return []DeviceConfigurationStatus{
		DeviceConfigurationStatusComplete,
		DeviceConfigurationStatusPending,
	}
}

// DeviceStatus - The current status of the device.
type DeviceStatus string

const (
	DeviceStatusUnknown           DeviceStatus = "Unknown"
	DeviceStatusOnline            DeviceStatus = "Online"
	DeviceStatusOffline           DeviceStatus = "Offline"
	DeviceStatusDeactivated       DeviceStatus = "Deactivated"
	DeviceStatusRequiresAttention DeviceStatus = "RequiresAttention"
	DeviceStatusMaintenanceMode   DeviceStatus = "MaintenanceMode"
	DeviceStatusCreating          DeviceStatus = "Creating"
	DeviceStatusProvisioning      DeviceStatus = "Provisioning"
	DeviceStatusDeactivating      DeviceStatus = "Deactivating"
	DeviceStatusDeleted           DeviceStatus = "Deleted"
	DeviceStatusReadyToSetup      DeviceStatus = "ReadyToSetup"
)

// PossibleDeviceStatusValues returns the possible values for the DeviceStatus const type.
func PossibleDeviceStatusValues() []DeviceStatus {
	return []DeviceStatus{
		DeviceStatusUnknown,
		DeviceStatusOnline,
		DeviceStatusOffline,
		DeviceStatusDeactivated,
		DeviceStatusRequiresAttention,
		DeviceStatusMaintenanceMode,
		DeviceStatusCreating,
		DeviceStatusProvisioning,
		DeviceStatusDeactivating,
		DeviceStatusDeleted,
		DeviceStatusReadyToSetup,
	}
}

// DeviceType - The type of the device.
type DeviceType string

const (
	DeviceTypeInvalid                     DeviceType = "Invalid"
	DeviceTypeSeries8000VirtualAppliance  DeviceType = "Series8000VirtualAppliance"
	DeviceTypeSeries8000PhysicalAppliance DeviceType = "Series8000PhysicalAppliance"
)

// PossibleDeviceTypeValues returns the possible values for the DeviceType const type.
func PossibleDeviceTypeValues() []DeviceType {
	return []DeviceType{
		DeviceTypeInvalid,
		DeviceTypeSeries8000VirtualAppliance,
		DeviceTypeSeries8000PhysicalAppliance,
	}
}

// EncryptionAlgorithm - The algorithm used to encrypt "Value".
type EncryptionAlgorithm string

const (
	EncryptionAlgorithmNone          EncryptionAlgorithm = "None"
	EncryptionAlgorithmAES256        EncryptionAlgorithm = "AES256"
	EncryptionAlgorithmRSAESPKCS1V15 EncryptionAlgorithm = "RSAES_PKCS1_v_1_5"
)

// PossibleEncryptionAlgorithmValues returns the possible values for the EncryptionAlgorithm const type.
func PossibleEncryptionAlgorithmValues() []EncryptionAlgorithm {
	return []EncryptionAlgorithm{
		EncryptionAlgorithmNone,
		EncryptionAlgorithmAES256,
		EncryptionAlgorithmRSAESPKCS1V15,
	}
}

// EncryptionStatus - The encryption status to indicates if encryption is enabled or not.
type EncryptionStatus string

const (
	EncryptionStatusEnabled  EncryptionStatus = "Enabled"
	EncryptionStatusDisabled EncryptionStatus = "Disabled"
)

// PossibleEncryptionStatusValues returns the possible values for the EncryptionStatus const type.
func PossibleEncryptionStatusValues() []EncryptionStatus {
	return []EncryptionStatus{
		EncryptionStatusEnabled,
		EncryptionStatusDisabled,
	}
}

// FeatureSupportStatus - The feature support status.
type FeatureSupportStatus string

const (
	FeatureSupportStatusNotAvailable             FeatureSupportStatus = "NotAvailable"
	FeatureSupportStatusUnsupportedDeviceVersion FeatureSupportStatus = "UnsupportedDeviceVersion"
	FeatureSupportStatusSupported                FeatureSupportStatus = "Supported"
)

// PossibleFeatureSupportStatusValues returns the possible values for the FeatureSupportStatus const type.
func PossibleFeatureSupportStatusValues() []FeatureSupportStatus {
	return []FeatureSupportStatus{
		FeatureSupportStatusNotAvailable,
		FeatureSupportStatusUnsupportedDeviceVersion,
		FeatureSupportStatusSupported,
	}
}

// HardwareComponentStatus - The status of the hardware component.
type HardwareComponentStatus string

const (
	HardwareComponentStatusUnknown    HardwareComponentStatus = "Unknown"
	HardwareComponentStatusNotPresent HardwareComponentStatus = "NotPresent"
	HardwareComponentStatusPoweredOff HardwareComponentStatus = "PoweredOff"
	HardwareComponentStatusOk         HardwareComponentStatus = "Ok"
	HardwareComponentStatusRecovering HardwareComponentStatus = "Recovering"
	HardwareComponentStatusWarning    HardwareComponentStatus = "Warning"
	HardwareComponentStatusFailure    HardwareComponentStatus = "Failure"
)

// PossibleHardwareComponentStatusValues returns the possible values for the HardwareComponentStatus const type.
func PossibleHardwareComponentStatusValues() []HardwareComponentStatus {
	return []HardwareComponentStatus{
		HardwareComponentStatusUnknown,
		HardwareComponentStatusNotPresent,
		HardwareComponentStatusPoweredOff,
		HardwareComponentStatusOk,
		HardwareComponentStatusRecovering,
		HardwareComponentStatusWarning,
		HardwareComponentStatusFailure,
	}
}

// ISCSIAndCloudStatus - Value indicating cloud and ISCSI status of network adapter.
type ISCSIAndCloudStatus string

const (
	ISCSIAndCloudStatusDisabled             ISCSIAndCloudStatus = "Disabled"
	ISCSIAndCloudStatusIscsiEnabled         ISCSIAndCloudStatus = "IscsiEnabled"
	ISCSIAndCloudStatusCloudEnabled         ISCSIAndCloudStatus = "CloudEnabled"
	ISCSIAndCloudStatusIscsiAndCloudEnabled ISCSIAndCloudStatus = "IscsiAndCloudEnabled"
)

// PossibleISCSIAndCloudStatusValues returns the possible values for the ISCSIAndCloudStatus const type.
func PossibleISCSIAndCloudStatusValues() []ISCSIAndCloudStatus {
	return []ISCSIAndCloudStatus{
		ISCSIAndCloudStatusDisabled,
		ISCSIAndCloudStatusIscsiEnabled,
		ISCSIAndCloudStatusCloudEnabled,
		ISCSIAndCloudStatusIscsiAndCloudEnabled,
	}
}

// InEligibilityCategory - The reason for inEligibility of device, in case it's not eligible for service data encryption key
// rollover.
type InEligibilityCategory string

const (
	InEligibilityCategoryDeviceNotOnline       InEligibilityCategory = "DeviceNotOnline"
	InEligibilityCategoryNotSupportedAppliance InEligibilityCategory = "NotSupportedAppliance"
	InEligibilityCategoryRolloverPending       InEligibilityCategory = "RolloverPending"
)

// PossibleInEligibilityCategoryValues returns the possible values for the InEligibilityCategory const type.
func PossibleInEligibilityCategoryValues() []InEligibilityCategory {
	return []InEligibilityCategory{
		InEligibilityCategoryDeviceNotOnline,
		InEligibilityCategoryNotSupportedAppliance,
		InEligibilityCategoryRolloverPending,
	}
}

// JobStatus - The current status of the job.
type JobStatus string

const (
	JobStatusRunning   JobStatus = "Running"
	JobStatusSucceeded JobStatus = "Succeeded"
	JobStatusFailed    JobStatus = "Failed"
	JobStatusCanceled  JobStatus = "Canceled"
)

// PossibleJobStatusValues returns the possible values for the JobStatus const type.
func PossibleJobStatusValues() []JobStatus {
	return []JobStatus{
		JobStatusRunning,
		JobStatusSucceeded,
		JobStatusFailed,
		JobStatusCanceled,
	}
}

// JobType - The type of the job.
type JobType string

const (
	JobTypeScheduledBackup           JobType = "ScheduledBackup"
	JobTypeManualBackup              JobType = "ManualBackup"
	JobTypeRestoreBackup             JobType = "RestoreBackup"
	JobTypeCloneVolume               JobType = "CloneVolume"
	JobTypeFailoverVolumeContainers  JobType = "FailoverVolumeContainers"
	JobTypeCreateLocallyPinnedVolume JobType = "CreateLocallyPinnedVolume"
	JobTypeModifyVolume              JobType = "ModifyVolume"
	JobTypeInstallUpdates            JobType = "InstallUpdates"
	JobTypeSupportPackageLogs        JobType = "SupportPackageLogs"
	JobTypeCreateCloudAppliance      JobType = "CreateCloudAppliance"
)

// PossibleJobTypeValues returns the possible values for the JobType const type.
func PossibleJobTypeValues() []JobType {
	return []JobType{
		JobTypeScheduledBackup,
		JobTypeManualBackup,
		JobTypeRestoreBackup,
		JobTypeCloneVolume,
		JobTypeFailoverVolumeContainers,
		JobTypeCreateLocallyPinnedVolume,
		JobTypeModifyVolume,
		JobTypeInstallUpdates,
		JobTypeSupportPackageLogs,
		JobTypeCreateCloudAppliance,
	}
}

// KeyRolloverStatus - The key rollover status to indicates if key rollover is required or not. If secret's encryption has
// been upgraded, then it requires key rollover.
type KeyRolloverStatus string

const (
	KeyRolloverStatusRequired    KeyRolloverStatus = "Required"
	KeyRolloverStatusNotRequired KeyRolloverStatus = "NotRequired"
)

// PossibleKeyRolloverStatusValues returns the possible values for the KeyRolloverStatus const type.
func PossibleKeyRolloverStatusValues() []KeyRolloverStatus {
	return []KeyRolloverStatus{
		KeyRolloverStatusRequired,
		KeyRolloverStatusNotRequired,
	}
}

// ManagerType - The type of StorSimple Manager.
type ManagerType string

const (
	ManagerTypeGardaV1    ManagerType = "GardaV1"
	ManagerTypeHelsinkiV1 ManagerType = "HelsinkiV1"
)

// PossibleManagerTypeValues returns the possible values for the ManagerType const type.
func PossibleManagerTypeValues() []ManagerType {
	return []ManagerType{
		ManagerTypeGardaV1,
		ManagerTypeHelsinkiV1,
	}
}

// MetricAggregationType - The metric aggregation type.
type MetricAggregationType string

const (
	MetricAggregationTypeAverage MetricAggregationType = "Average"
	MetricAggregationTypeLast    MetricAggregationType = "Last"
	MetricAggregationTypeMaximum MetricAggregationType = "Maximum"
	MetricAggregationTypeMinimum MetricAggregationType = "Minimum"
	MetricAggregationTypeNone    MetricAggregationType = "None"
	MetricAggregationTypeTotal   MetricAggregationType = "Total"
)

// PossibleMetricAggregationTypeValues returns the possible values for the MetricAggregationType const type.
func PossibleMetricAggregationTypeValues() []MetricAggregationType {
	return []MetricAggregationType{
		MetricAggregationTypeAverage,
		MetricAggregationTypeLast,
		MetricAggregationTypeMaximum,
		MetricAggregationTypeMinimum,
		MetricAggregationTypeNone,
		MetricAggregationTypeTotal,
	}
}

// MetricUnit - The metric unit.
type MetricUnit string

const (
	MetricUnitBytes          MetricUnit = "Bytes"
	MetricUnitBytesPerSecond MetricUnit = "BytesPerSecond"
	MetricUnitCount          MetricUnit = "Count"
	MetricUnitCountPerSecond MetricUnit = "CountPerSecond"
	MetricUnitPercent        MetricUnit = "Percent"
	MetricUnitSeconds        MetricUnit = "Seconds"
)

// PossibleMetricUnitValues returns the possible values for the MetricUnit const type.
func PossibleMetricUnitValues() []MetricUnit {
	return []MetricUnit{
		MetricUnitBytes,
		MetricUnitBytesPerSecond,
		MetricUnitCount,
		MetricUnitCountPerSecond,
		MetricUnitPercent,
		MetricUnitSeconds,
	}
}

// MonitoringStatus - The monitoring status of the volume.
type MonitoringStatus string

const (
	MonitoringStatusEnabled  MonitoringStatus = "Enabled"
	MonitoringStatusDisabled MonitoringStatus = "Disabled"
)

// PossibleMonitoringStatusValues returns the possible values for the MonitoringStatus const type.
func PossibleMonitoringStatusValues() []MonitoringStatus {
	return []MonitoringStatus{
		MonitoringStatusEnabled,
		MonitoringStatusDisabled,
	}
}

// NetInterfaceID - The ID of the network adapter.
type NetInterfaceID string

const (
	NetInterfaceIDInvalid NetInterfaceID = "Invalid"
	NetInterfaceIDData0   NetInterfaceID = "Data0"
	NetInterfaceIDData1   NetInterfaceID = "Data1"
	NetInterfaceIDData2   NetInterfaceID = "Data2"
	NetInterfaceIDData3   NetInterfaceID = "Data3"
	NetInterfaceIDData4   NetInterfaceID = "Data4"
	NetInterfaceIDData5   NetInterfaceID = "Data5"
)

// PossibleNetInterfaceIDValues returns the possible values for the NetInterfaceID const type.
func PossibleNetInterfaceIDValues() []NetInterfaceID {
	return []NetInterfaceID{
		NetInterfaceIDInvalid,
		NetInterfaceIDData0,
		NetInterfaceIDData1,
		NetInterfaceIDData2,
		NetInterfaceIDData3,
		NetInterfaceIDData4,
		NetInterfaceIDData5,
	}
}

// NetInterfaceStatus - Value indicating status of network adapter.
type NetInterfaceStatus string

const (
	NetInterfaceStatusEnabled  NetInterfaceStatus = "Enabled"
	NetInterfaceStatusDisabled NetInterfaceStatus = "Disabled"
)

// PossibleNetInterfaceStatusValues returns the possible values for the NetInterfaceStatus const type.
func PossibleNetInterfaceStatusValues() []NetInterfaceStatus {
	return []NetInterfaceStatus{
		NetInterfaceStatusEnabled,
		NetInterfaceStatusDisabled,
	}
}

// NetworkMode - The mode of network adapter, either IPv4, IPv6 or both.
type NetworkMode string

const (
	NetworkModeInvalid NetworkMode = "Invalid"
	NetworkModeIPV4    NetworkMode = "IPV4"
	NetworkModeIPV6    NetworkMode = "IPV6"
	NetworkModeBOTH    NetworkMode = "BOTH"
)

// PossibleNetworkModeValues returns the possible values for the NetworkMode const type.
func PossibleNetworkModeValues() []NetworkMode {
	return []NetworkMode{
		NetworkModeInvalid,
		NetworkModeIPV4,
		NetworkModeIPV6,
		NetworkModeBOTH,
	}
}

// OperationStatus - The operation status on the volume.
type OperationStatus string

const (
	OperationStatusNone      OperationStatus = "None"
	OperationStatusUpdating  OperationStatus = "Updating"
	OperationStatusDeleting  OperationStatus = "Deleting"
	OperationStatusRestoring OperationStatus = "Restoring"
)

// PossibleOperationStatusValues returns the possible values for the OperationStatus const type.
func PossibleOperationStatusValues() []OperationStatus {
	return []OperationStatus{
		OperationStatusNone,
		OperationStatusUpdating,
		OperationStatusDeleting,
		OperationStatusRestoring,
	}
}

// OwnerShipStatus - The owner ship status of the volume container. Only when the status is "NotOwned", the delete operation
// on the volume container is permitted.
type OwnerShipStatus string

const (
	OwnerShipStatusOwned    OwnerShipStatus = "Owned"
	OwnerShipStatusNotOwned OwnerShipStatus = "NotOwned"
)

// PossibleOwnerShipStatusValues returns the possible values for the OwnerShipStatus const type.
func PossibleOwnerShipStatusValues() []OwnerShipStatus {
	return []OwnerShipStatus{
		OwnerShipStatusOwned,
		OwnerShipStatusNotOwned,
	}
}

// RecurrenceType - The recurrence type.
type RecurrenceType string

const (
	RecurrenceTypeMinutes RecurrenceType = "Minutes"
	RecurrenceTypeHourly  RecurrenceType = "Hourly"
	RecurrenceTypeDaily   RecurrenceType = "Daily"
	RecurrenceTypeWeekly  RecurrenceType = "Weekly"
)

// PossibleRecurrenceTypeValues returns the possible values for the RecurrenceType const type.
func PossibleRecurrenceTypeValues() []RecurrenceType {
	return []RecurrenceType{
		RecurrenceTypeMinutes,
		RecurrenceTypeHourly,
		RecurrenceTypeDaily,
		RecurrenceTypeWeekly,
	}
}

// RemoteManagementModeConfiguration - The remote management mode.
type RemoteManagementModeConfiguration string

const (
	RemoteManagementModeConfigurationUnknown             RemoteManagementModeConfiguration = "Unknown"
	RemoteManagementModeConfigurationDisabled            RemoteManagementModeConfiguration = "Disabled"
	RemoteManagementModeConfigurationHTTPSEnabled        RemoteManagementModeConfiguration = "HttpsEnabled"
	RemoteManagementModeConfigurationHTTPSAndHTTPEnabled RemoteManagementModeConfiguration = "HttpsAndHttpEnabled"
)

// PossibleRemoteManagementModeConfigurationValues returns the possible values for the RemoteManagementModeConfiguration const type.
func PossibleRemoteManagementModeConfigurationValues() []RemoteManagementModeConfiguration {
	return []RemoteManagementModeConfiguration{
		RemoteManagementModeConfigurationUnknown,
		RemoteManagementModeConfigurationDisabled,
		RemoteManagementModeConfigurationHTTPSEnabled,
		RemoteManagementModeConfigurationHTTPSAndHTTPEnabled,
	}
}

// SSLStatus - Signifies whether SSL needs to be enabled or not.
type SSLStatus string

const (
	SSLStatusEnabled  SSLStatus = "Enabled"
	SSLStatusDisabled SSLStatus = "Disabled"
)

// PossibleSSLStatusValues returns the possible values for the SSLStatus const type.
func PossibleSSLStatusValues() []SSLStatus {
	return []SSLStatus{
		SSLStatusEnabled,
		SSLStatusDisabled,
	}
}

// ScheduleStatus - The schedule status.
type ScheduleStatus string

const (
	ScheduleStatusEnabled  ScheduleStatus = "Enabled"
	ScheduleStatusDisabled ScheduleStatus = "Disabled"
)

// PossibleScheduleStatusValues returns the possible values for the ScheduleStatus const type.
func PossibleScheduleStatusValues() []ScheduleStatus {
	return []ScheduleStatus{
		ScheduleStatusEnabled,
		ScheduleStatusDisabled,
	}
}

// ScheduledBackupStatus - Indicates whether at least one of the schedules in the backup policy is active or not.
type ScheduledBackupStatus string

const (
	ScheduledBackupStatusDisabled ScheduledBackupStatus = "Disabled"
	ScheduledBackupStatusEnabled  ScheduledBackupStatus = "Enabled"
)

// PossibleScheduledBackupStatusValues returns the possible values for the ScheduledBackupStatus const type.
func PossibleScheduledBackupStatusValues() []ScheduledBackupStatus {
	return []ScheduledBackupStatus{
		ScheduledBackupStatusDisabled,
		ScheduledBackupStatusEnabled,
	}
}

// TargetEligibilityResultCode - The result code for the error, due to which the device does not qualify as a failover target
// device.
type TargetEligibilityResultCode string

const (
	TargetEligibilityResultCodeTargetAndSourceCannotBeSameError          TargetEligibilityResultCode = "TargetAndSourceCannotBeSameError"
	TargetEligibilityResultCodeTargetIsNotOnlineError                    TargetEligibilityResultCode = "TargetIsNotOnlineError"
	TargetEligibilityResultCodeTargetSourceIncompatibleVersionError      TargetEligibilityResultCode = "TargetSourceIncompatibleVersionError"
	TargetEligibilityResultCodeLocalToTieredVolumesConversionWarning     TargetEligibilityResultCode = "LocalToTieredVolumesConversionWarning"
	TargetEligibilityResultCodeTargetInsufficientCapacityError           TargetEligibilityResultCode = "TargetInsufficientCapacityError"
	TargetEligibilityResultCodeTargetInsufficientLocalVolumeMemoryError  TargetEligibilityResultCode = "TargetInsufficientLocalVolumeMemoryError"
	TargetEligibilityResultCodeTargetInsufficientTieredVolumeMemoryError TargetEligibilityResultCode = "TargetInsufficientTieredVolumeMemoryError"
)

// PossibleTargetEligibilityResultCodeValues returns the possible values for the TargetEligibilityResultCode const type.
func PossibleTargetEligibilityResultCodeValues() []TargetEligibilityResultCode {
	return []TargetEligibilityResultCode{
		TargetEligibilityResultCodeTargetAndSourceCannotBeSameError,
		TargetEligibilityResultCodeTargetIsNotOnlineError,
		TargetEligibilityResultCodeTargetSourceIncompatibleVersionError,
		TargetEligibilityResultCodeLocalToTieredVolumesConversionWarning,
		TargetEligibilityResultCodeTargetInsufficientCapacityError,
		TargetEligibilityResultCodeTargetInsufficientLocalVolumeMemoryError,
		TargetEligibilityResultCodeTargetInsufficientTieredVolumeMemoryError,
	}
}

// TargetEligibilityStatus - The eligibility status of device, as a failover target device.
type TargetEligibilityStatus string

const (
	TargetEligibilityStatusNotEligible TargetEligibilityStatus = "NotEligible"
	TargetEligibilityStatusEligible    TargetEligibilityStatus = "Eligible"
)

// PossibleTargetEligibilityStatusValues returns the possible values for the TargetEligibilityStatus const type.
func PossibleTargetEligibilityStatusValues() []TargetEligibilityStatus {
	return []TargetEligibilityStatus{
		TargetEligibilityStatusNotEligible,
		TargetEligibilityStatusEligible,
	}
}

// VirtualMachineAPIType - The virtual machine API type.
type VirtualMachineAPIType string

const (
	VirtualMachineAPITypeClassic VirtualMachineAPIType = "Classic"
	VirtualMachineAPITypeArm     VirtualMachineAPIType = "Arm"
)

// PossibleVirtualMachineAPITypeValues returns the possible values for the VirtualMachineAPIType const type.
func PossibleVirtualMachineAPITypeValues() []VirtualMachineAPIType {
	return []VirtualMachineAPIType{
		VirtualMachineAPITypeClassic,
		VirtualMachineAPITypeArm,
	}
}

// VolumeStatus - The volume status.
type VolumeStatus string

const (
	VolumeStatusOnline  VolumeStatus = "Online"
	VolumeStatusOffline VolumeStatus = "Offline"
)

// PossibleVolumeStatusValues returns the possible values for the VolumeStatus const type.
func PossibleVolumeStatusValues() []VolumeStatus {
	return []VolumeStatus{
		VolumeStatusOnline,
		VolumeStatusOffline,
	}
}

// VolumeType - The volume type.
type VolumeType string

const (
	VolumeTypeTiered        VolumeType = "Tiered"
	VolumeTypeArchival      VolumeType = "Archival"
	VolumeTypeLocallyPinned VolumeType = "LocallyPinned"
)

// PossibleVolumeTypeValues returns the possible values for the VolumeType const type.
func PossibleVolumeTypeValues() []VolumeType {
	return []VolumeType{
		VolumeTypeTiered,
		VolumeTypeArchival,
		VolumeTypeLocallyPinned,
	}
}

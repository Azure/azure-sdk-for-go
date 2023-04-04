//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdevtestlabs

const (
	moduleName    = "armdevtestlabs"
	moduleVersion = "v1.1.0"
)

// CostThresholdStatus - Indicates whether this threshold will be displayed on cost charts.
type CostThresholdStatus string

const (
	CostThresholdStatusDisabled CostThresholdStatus = "Disabled"
	CostThresholdStatusEnabled  CostThresholdStatus = "Enabled"
)

// PossibleCostThresholdStatusValues returns the possible values for the CostThresholdStatus const type.
func PossibleCostThresholdStatusValues() []CostThresholdStatus {
	return []CostThresholdStatus{
		CostThresholdStatusDisabled,
		CostThresholdStatusEnabled,
	}
}

// CostType - The type of the cost.
type CostType string

const (
	CostTypeProjected   CostType = "Projected"
	CostTypeReported    CostType = "Reported"
	CostTypeUnavailable CostType = "Unavailable"
)

// PossibleCostTypeValues returns the possible values for the CostType const type.
func PossibleCostTypeValues() []CostType {
	return []CostType{
		CostTypeProjected,
		CostTypeReported,
		CostTypeUnavailable,
	}
}

// CustomImageOsType - The OS type of the custom image (i.e. Windows, Linux)
type CustomImageOsType string

const (
	CustomImageOsTypeLinux   CustomImageOsType = "Linux"
	CustomImageOsTypeNone    CustomImageOsType = "None"
	CustomImageOsTypeWindows CustomImageOsType = "Windows"
)

// PossibleCustomImageOsTypeValues returns the possible values for the CustomImageOsType const type.
func PossibleCustomImageOsTypeValues() []CustomImageOsType {
	return []CustomImageOsType{
		CustomImageOsTypeLinux,
		CustomImageOsTypeNone,
		CustomImageOsTypeWindows,
	}
}

// EnableStatus - Indicates if the artifact source is enabled (values: Enabled, Disabled).
type EnableStatus string

const (
	EnableStatusDisabled EnableStatus = "Disabled"
	EnableStatusEnabled  EnableStatus = "Enabled"
)

// PossibleEnableStatusValues returns the possible values for the EnableStatus const type.
func PossibleEnableStatusValues() []EnableStatus {
	return []EnableStatus{
		EnableStatusDisabled,
		EnableStatusEnabled,
	}
}

// EnvironmentPermission - The access rights to be granted to the user when provisioning an environment
type EnvironmentPermission string

const (
	EnvironmentPermissionContributor EnvironmentPermission = "Contributor"
	EnvironmentPermissionReader      EnvironmentPermission = "Reader"
)

// PossibleEnvironmentPermissionValues returns the possible values for the EnvironmentPermission const type.
func PossibleEnvironmentPermissionValues() []EnvironmentPermission {
	return []EnvironmentPermission{
		EnvironmentPermissionContributor,
		EnvironmentPermissionReader,
	}
}

// FileUploadOptions - Options for uploading the files for the artifact. UploadFilesAndGenerateSasTokens is the default value.
type FileUploadOptions string

const (
	FileUploadOptionsNone                            FileUploadOptions = "None"
	FileUploadOptionsUploadFilesAndGenerateSasTokens FileUploadOptions = "UploadFilesAndGenerateSasTokens"
)

// PossibleFileUploadOptionsValues returns the possible values for the FileUploadOptions const type.
func PossibleFileUploadOptionsValues() []FileUploadOptions {
	return []FileUploadOptions{
		FileUploadOptionsNone,
		FileUploadOptionsUploadFilesAndGenerateSasTokens,
	}
}

// HTTPStatusCode - The status code for the operation.
type HTTPStatusCode string

const (
	HTTPStatusCodeAccepted                     HTTPStatusCode = "Accepted"
	HTTPStatusCodeAmbiguous                    HTTPStatusCode = "Ambiguous"
	HTTPStatusCodeBadGateway                   HTTPStatusCode = "BadGateway"
	HTTPStatusCodeBadRequest                   HTTPStatusCode = "BadRequest"
	HTTPStatusCodeConflict                     HTTPStatusCode = "Conflict"
	HTTPStatusCodeContinue                     HTTPStatusCode = "Continue"
	HTTPStatusCodeCreated                      HTTPStatusCode = "Created"
	HTTPStatusCodeExpectationFailed            HTTPStatusCode = "ExpectationFailed"
	HTTPStatusCodeForbidden                    HTTPStatusCode = "Forbidden"
	HTTPStatusCodeFound                        HTTPStatusCode = "Found"
	HTTPStatusCodeGatewayTimeout               HTTPStatusCode = "GatewayTimeout"
	HTTPStatusCodeGone                         HTTPStatusCode = "Gone"
	HTTPStatusCodeHTTPVersionNotSupported      HTTPStatusCode = "HttpVersionNotSupported"
	HTTPStatusCodeInternalServerError          HTTPStatusCode = "InternalServerError"
	HTTPStatusCodeLengthRequired               HTTPStatusCode = "LengthRequired"
	HTTPStatusCodeMethodNotAllowed             HTTPStatusCode = "MethodNotAllowed"
	HTTPStatusCodeMoved                        HTTPStatusCode = "Moved"
	HTTPStatusCodeMovedPermanently             HTTPStatusCode = "MovedPermanently"
	HTTPStatusCodeMultipleChoices              HTTPStatusCode = "MultipleChoices"
	HTTPStatusCodeNoContent                    HTTPStatusCode = "NoContent"
	HTTPStatusCodeNonAuthoritativeInformation  HTTPStatusCode = "NonAuthoritativeInformation"
	HTTPStatusCodeNotAcceptable                HTTPStatusCode = "NotAcceptable"
	HTTPStatusCodeNotFound                     HTTPStatusCode = "NotFound"
	HTTPStatusCodeNotImplemented               HTTPStatusCode = "NotImplemented"
	HTTPStatusCodeNotModified                  HTTPStatusCode = "NotModified"
	HTTPStatusCodeOK                           HTTPStatusCode = "OK"
	HTTPStatusCodePartialContent               HTTPStatusCode = "PartialContent"
	HTTPStatusCodePaymentRequired              HTTPStatusCode = "PaymentRequired"
	HTTPStatusCodePreconditionFailed           HTTPStatusCode = "PreconditionFailed"
	HTTPStatusCodeProxyAuthenticationRequired  HTTPStatusCode = "ProxyAuthenticationRequired"
	HTTPStatusCodeRedirect                     HTTPStatusCode = "Redirect"
	HTTPStatusCodeRedirectKeepVerb             HTTPStatusCode = "RedirectKeepVerb"
	HTTPStatusCodeRedirectMethod               HTTPStatusCode = "RedirectMethod"
	HTTPStatusCodeRequestEntityTooLarge        HTTPStatusCode = "RequestEntityTooLarge"
	HTTPStatusCodeRequestTimeout               HTTPStatusCode = "RequestTimeout"
	HTTPStatusCodeRequestURITooLong            HTTPStatusCode = "RequestUriTooLong"
	HTTPStatusCodeRequestedRangeNotSatisfiable HTTPStatusCode = "RequestedRangeNotSatisfiable"
	HTTPStatusCodeResetContent                 HTTPStatusCode = "ResetContent"
	HTTPStatusCodeSeeOther                     HTTPStatusCode = "SeeOther"
	HTTPStatusCodeServiceUnavailable           HTTPStatusCode = "ServiceUnavailable"
	HTTPStatusCodeSwitchingProtocols           HTTPStatusCode = "SwitchingProtocols"
	HTTPStatusCodeTemporaryRedirect            HTTPStatusCode = "TemporaryRedirect"
	HTTPStatusCodeUnauthorized                 HTTPStatusCode = "Unauthorized"
	HTTPStatusCodeUnsupportedMediaType         HTTPStatusCode = "UnsupportedMediaType"
	HTTPStatusCodeUnused                       HTTPStatusCode = "Unused"
	HTTPStatusCodeUpgradeRequired              HTTPStatusCode = "UpgradeRequired"
	HTTPStatusCodeUseProxy                     HTTPStatusCode = "UseProxy"
)

// PossibleHTTPStatusCodeValues returns the possible values for the HTTPStatusCode const type.
func PossibleHTTPStatusCodeValues() []HTTPStatusCode {
	return []HTTPStatusCode{
		HTTPStatusCodeAccepted,
		HTTPStatusCodeAmbiguous,
		HTTPStatusCodeBadGateway,
		HTTPStatusCodeBadRequest,
		HTTPStatusCodeConflict,
		HTTPStatusCodeContinue,
		HTTPStatusCodeCreated,
		HTTPStatusCodeExpectationFailed,
		HTTPStatusCodeForbidden,
		HTTPStatusCodeFound,
		HTTPStatusCodeGatewayTimeout,
		HTTPStatusCodeGone,
		HTTPStatusCodeHTTPVersionNotSupported,
		HTTPStatusCodeInternalServerError,
		HTTPStatusCodeLengthRequired,
		HTTPStatusCodeMethodNotAllowed,
		HTTPStatusCodeMoved,
		HTTPStatusCodeMovedPermanently,
		HTTPStatusCodeMultipleChoices,
		HTTPStatusCodeNoContent,
		HTTPStatusCodeNonAuthoritativeInformation,
		HTTPStatusCodeNotAcceptable,
		HTTPStatusCodeNotFound,
		HTTPStatusCodeNotImplemented,
		HTTPStatusCodeNotModified,
		HTTPStatusCodeOK,
		HTTPStatusCodePartialContent,
		HTTPStatusCodePaymentRequired,
		HTTPStatusCodePreconditionFailed,
		HTTPStatusCodeProxyAuthenticationRequired,
		HTTPStatusCodeRedirect,
		HTTPStatusCodeRedirectKeepVerb,
		HTTPStatusCodeRedirectMethod,
		HTTPStatusCodeRequestEntityTooLarge,
		HTTPStatusCodeRequestTimeout,
		HTTPStatusCodeRequestURITooLong,
		HTTPStatusCodeRequestedRangeNotSatisfiable,
		HTTPStatusCodeResetContent,
		HTTPStatusCodeSeeOther,
		HTTPStatusCodeServiceUnavailable,
		HTTPStatusCodeSwitchingProtocols,
		HTTPStatusCodeTemporaryRedirect,
		HTTPStatusCodeUnauthorized,
		HTTPStatusCodeUnsupportedMediaType,
		HTTPStatusCodeUnused,
		HTTPStatusCodeUpgradeRequired,
		HTTPStatusCodeUseProxy,
	}
}

// HostCachingOptions - Caching option for a data disk (i.e. None, ReadOnly, ReadWrite).
type HostCachingOptions string

const (
	HostCachingOptionsNone      HostCachingOptions = "None"
	HostCachingOptionsReadOnly  HostCachingOptions = "ReadOnly"
	HostCachingOptionsReadWrite HostCachingOptions = "ReadWrite"
)

// PossibleHostCachingOptionsValues returns the possible values for the HostCachingOptions const type.
func PossibleHostCachingOptionsValues() []HostCachingOptions {
	return []HostCachingOptions{
		HostCachingOptionsNone,
		HostCachingOptionsReadOnly,
		HostCachingOptionsReadWrite,
	}
}

// LinuxOsState - The state of the Linux OS (i.e. NonDeprovisioned, DeprovisionRequested, DeprovisionApplied).
type LinuxOsState string

const (
	LinuxOsStateDeprovisionApplied   LinuxOsState = "DeprovisionApplied"
	LinuxOsStateDeprovisionRequested LinuxOsState = "DeprovisionRequested"
	LinuxOsStateNonDeprovisioned     LinuxOsState = "NonDeprovisioned"
)

// PossibleLinuxOsStateValues returns the possible values for the LinuxOsState const type.
func PossibleLinuxOsStateValues() []LinuxOsState {
	return []LinuxOsState{
		LinuxOsStateDeprovisionApplied,
		LinuxOsStateDeprovisionRequested,
		LinuxOsStateNonDeprovisioned,
	}
}

// ManagedIdentityType - Managed identity.
type ManagedIdentityType string

const (
	ManagedIdentityTypeNone                       ManagedIdentityType = "None"
	ManagedIdentityTypeSystemAssigned             ManagedIdentityType = "SystemAssigned"
	ManagedIdentityTypeSystemAssignedUserAssigned ManagedIdentityType = "SystemAssigned,UserAssigned"
	ManagedIdentityTypeUserAssigned               ManagedIdentityType = "UserAssigned"
)

// PossibleManagedIdentityTypeValues returns the possible values for the ManagedIdentityType const type.
func PossibleManagedIdentityTypeValues() []ManagedIdentityType {
	return []ManagedIdentityType{
		ManagedIdentityTypeNone,
		ManagedIdentityTypeSystemAssigned,
		ManagedIdentityTypeSystemAssignedUserAssigned,
		ManagedIdentityTypeUserAssigned,
	}
}

// NotificationChannelEventType - The event type for which this notification is enabled (i.e. AutoShutdown, Cost)
type NotificationChannelEventType string

const (
	NotificationChannelEventTypeAutoShutdown NotificationChannelEventType = "AutoShutdown"
	NotificationChannelEventTypeCost         NotificationChannelEventType = "Cost"
)

// PossibleNotificationChannelEventTypeValues returns the possible values for the NotificationChannelEventType const type.
func PossibleNotificationChannelEventTypeValues() []NotificationChannelEventType {
	return []NotificationChannelEventType{
		NotificationChannelEventTypeAutoShutdown,
		NotificationChannelEventTypeCost,
	}
}

// PolicyEvaluatorType - The evaluator type of the policy (i.e. AllowedValuesPolicy, MaxValuePolicy).
type PolicyEvaluatorType string

const (
	PolicyEvaluatorTypeAllowedValuesPolicy PolicyEvaluatorType = "AllowedValuesPolicy"
	PolicyEvaluatorTypeMaxValuePolicy      PolicyEvaluatorType = "MaxValuePolicy"
)

// PossiblePolicyEvaluatorTypeValues returns the possible values for the PolicyEvaluatorType const type.
func PossiblePolicyEvaluatorTypeValues() []PolicyEvaluatorType {
	return []PolicyEvaluatorType{
		PolicyEvaluatorTypeAllowedValuesPolicy,
		PolicyEvaluatorTypeMaxValuePolicy,
	}
}

// PolicyFactName - The fact name of the policy (e.g. LabVmCount, LabVmSize, MaxVmsAllowedPerLab, etc.
type PolicyFactName string

const (
	PolicyFactNameEnvironmentTemplate         PolicyFactName = "EnvironmentTemplate"
	PolicyFactNameGalleryImage                PolicyFactName = "GalleryImage"
	PolicyFactNameLabPremiumVMCount           PolicyFactName = "LabPremiumVmCount"
	PolicyFactNameLabTargetCost               PolicyFactName = "LabTargetCost"
	PolicyFactNameLabVMCount                  PolicyFactName = "LabVmCount"
	PolicyFactNameLabVMSize                   PolicyFactName = "LabVmSize"
	PolicyFactNameScheduleEditPermission      PolicyFactName = "ScheduleEditPermission"
	PolicyFactNameUserOwnedLabPremiumVMCount  PolicyFactName = "UserOwnedLabPremiumVmCount"
	PolicyFactNameUserOwnedLabVMCount         PolicyFactName = "UserOwnedLabVmCount"
	PolicyFactNameUserOwnedLabVMCountInSubnet PolicyFactName = "UserOwnedLabVmCountInSubnet"
)

// PossiblePolicyFactNameValues returns the possible values for the PolicyFactName const type.
func PossiblePolicyFactNameValues() []PolicyFactName {
	return []PolicyFactName{
		PolicyFactNameEnvironmentTemplate,
		PolicyFactNameGalleryImage,
		PolicyFactNameLabPremiumVMCount,
		PolicyFactNameLabTargetCost,
		PolicyFactNameLabVMCount,
		PolicyFactNameLabVMSize,
		PolicyFactNameScheduleEditPermission,
		PolicyFactNameUserOwnedLabPremiumVMCount,
		PolicyFactNameUserOwnedLabVMCount,
		PolicyFactNameUserOwnedLabVMCountInSubnet,
	}
}

// PolicyStatus - The status of the policy.
type PolicyStatus string

const (
	PolicyStatusDisabled PolicyStatus = "Disabled"
	PolicyStatusEnabled  PolicyStatus = "Enabled"
)

// PossiblePolicyStatusValues returns the possible values for the PolicyStatus const type.
func PossiblePolicyStatusValues() []PolicyStatus {
	return []PolicyStatus{
		PolicyStatusDisabled,
		PolicyStatusEnabled,
	}
}

// PremiumDataDisk - The setting to enable usage of premium data disks. When its value is 'Enabled', creation of standard
// or premium data disks is allowed. When its value is 'Disabled', only creation of standard data
// disks is allowed.
type PremiumDataDisk string

const (
	PremiumDataDiskDisabled PremiumDataDisk = "Disabled"
	PremiumDataDiskEnabled  PremiumDataDisk = "Enabled"
)

// PossiblePremiumDataDiskValues returns the possible values for the PremiumDataDisk const type.
func PossiblePremiumDataDiskValues() []PremiumDataDisk {
	return []PremiumDataDisk{
		PremiumDataDiskDisabled,
		PremiumDataDiskEnabled,
	}
}

// ReportingCycleType - Reporting cycle type.
type ReportingCycleType string

const (
	ReportingCycleTypeCalendarMonth ReportingCycleType = "CalendarMonth"
	ReportingCycleTypeCustom        ReportingCycleType = "Custom"
)

// PossibleReportingCycleTypeValues returns the possible values for the ReportingCycleType const type.
func PossibleReportingCycleTypeValues() []ReportingCycleType {
	return []ReportingCycleType{
		ReportingCycleTypeCalendarMonth,
		ReportingCycleTypeCustom,
	}
}

// SourceControlType - The artifact source's type.
type SourceControlType string

const (
	SourceControlTypeGitHub         SourceControlType = "GitHub"
	SourceControlTypeStorageAccount SourceControlType = "StorageAccount"
	SourceControlTypeVsoGit         SourceControlType = "VsoGit"
)

// PossibleSourceControlTypeValues returns the possible values for the SourceControlType const type.
func PossibleSourceControlTypeValues() []SourceControlType {
	return []SourceControlType{
		SourceControlTypeGitHub,
		SourceControlTypeStorageAccount,
		SourceControlTypeVsoGit,
	}
}

// StorageType - The storage type for the disk (i.e. Standard, Premium).
type StorageType string

const (
	StorageTypePremium     StorageType = "Premium"
	StorageTypeStandard    StorageType = "Standard"
	StorageTypeStandardSSD StorageType = "StandardSSD"
)

// PossibleStorageTypeValues returns the possible values for the StorageType const type.
func PossibleStorageTypeValues() []StorageType {
	return []StorageType{
		StorageTypePremium,
		StorageTypeStandard,
		StorageTypeStandardSSD,
	}
}

// TargetCostStatus - Target cost status
type TargetCostStatus string

const (
	TargetCostStatusDisabled TargetCostStatus = "Disabled"
	TargetCostStatusEnabled  TargetCostStatus = "Enabled"
)

// PossibleTargetCostStatusValues returns the possible values for the TargetCostStatus const type.
func PossibleTargetCostStatusValues() []TargetCostStatus {
	return []TargetCostStatus{
		TargetCostStatusDisabled,
		TargetCostStatusEnabled,
	}
}

// TransportProtocol - The transport protocol for the endpoint.
type TransportProtocol string

const (
	TransportProtocolTCP TransportProtocol = "Tcp"
	TransportProtocolUDP TransportProtocol = "Udp"
)

// PossibleTransportProtocolValues returns the possible values for the TransportProtocol const type.
func PossibleTransportProtocolValues() []TransportProtocol {
	return []TransportProtocol{
		TransportProtocolTCP,
		TransportProtocolUDP,
	}
}

// UsagePermissionType - The permission policy of the subnet for allowing public IP addresses (i.e. Allow, Deny)).
type UsagePermissionType string

const (
	UsagePermissionTypeAllow   UsagePermissionType = "Allow"
	UsagePermissionTypeDefault UsagePermissionType = "Default"
	UsagePermissionTypeDeny    UsagePermissionType = "Deny"
)

// PossibleUsagePermissionTypeValues returns the possible values for the UsagePermissionType const type.
func PossibleUsagePermissionTypeValues() []UsagePermissionType {
	return []UsagePermissionType{
		UsagePermissionTypeAllow,
		UsagePermissionTypeDefault,
		UsagePermissionTypeDeny,
	}
}

// VirtualMachineCreationSource - Tells source of creation of lab virtual machine. Output property only.
type VirtualMachineCreationSource string

const (
	VirtualMachineCreationSourceFromCustomImage        VirtualMachineCreationSource = "FromCustomImage"
	VirtualMachineCreationSourceFromGalleryImage       VirtualMachineCreationSource = "FromGalleryImage"
	VirtualMachineCreationSourceFromSharedGalleryImage VirtualMachineCreationSource = "FromSharedGalleryImage"
)

// PossibleVirtualMachineCreationSourceValues returns the possible values for the VirtualMachineCreationSource const type.
func PossibleVirtualMachineCreationSourceValues() []VirtualMachineCreationSource {
	return []VirtualMachineCreationSource{
		VirtualMachineCreationSourceFromCustomImage,
		VirtualMachineCreationSourceFromGalleryImage,
		VirtualMachineCreationSourceFromSharedGalleryImage,
	}
}

// WindowsOsState - The state of the Windows OS (i.e. NonSysprepped, SysprepRequested, SysprepApplied).
type WindowsOsState string

const (
	WindowsOsStateNonSysprepped    WindowsOsState = "NonSysprepped"
	WindowsOsStateSysprepApplied   WindowsOsState = "SysprepApplied"
	WindowsOsStateSysprepRequested WindowsOsState = "SysprepRequested"
)

// PossibleWindowsOsStateValues returns the possible values for the WindowsOsState const type.
func PossibleWindowsOsStateValues() []WindowsOsState {
	return []WindowsOsState{
		WindowsOsStateNonSysprepped,
		WindowsOsStateSysprepApplied,
		WindowsOsStateSysprepRequested,
	}
}

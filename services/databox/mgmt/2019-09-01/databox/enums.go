package databox

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// AccessProtocol enumerates the values for access protocol.
type AccessProtocol string

const (
	// NFS Network File System protocol(NFS).
	NFS AccessProtocol = "NFS"
	// SMB Server Message Block protocol(SMB).
	SMB AccessProtocol = "SMB"
)

// PossibleAccessProtocolValues returns an array of possible values for the AccessProtocol const type.
func PossibleAccessProtocolValues() []AccessProtocol {
	return []AccessProtocol{NFS, SMB}
}

// AddressType enumerates the values for address type.
type AddressType string

const (
	// Commercial Commercial Address.
	Commercial AddressType = "Commercial"
	// None Address type not known.
	None AddressType = "None"
	// Residential Residential Address.
	Residential AddressType = "Residential"
)

// PossibleAddressTypeValues returns an array of possible values for the AddressType const type.
func PossibleAddressTypeValues() []AddressType {
	return []AddressType{Commercial, None, Residential}
}

// AddressValidationStatus enumerates the values for address validation status.
type AddressValidationStatus string

const (
	// Ambiguous Address provided is ambiguous, please choose one of the alternate addresses returned.
	Ambiguous AddressValidationStatus = "Ambiguous"
	// Invalid Address provided is invalid or not supported.
	Invalid AddressValidationStatus = "Invalid"
	// Valid Address provided is valid.
	Valid AddressValidationStatus = "Valid"
)

// PossibleAddressValidationStatusValues returns an array of possible values for the AddressValidationStatus const type.
func PossibleAddressValidationStatusValues() []AddressValidationStatus {
	return []AddressValidationStatus{Ambiguous, Invalid, Valid}
}

// CopyLogDetailsType enumerates the values for copy log details type.
type CopyLogDetailsType string

const (
	// CopyLogDetailsTypeCopyLogDetails ...
	CopyLogDetailsTypeCopyLogDetails CopyLogDetailsType = "CopyLogDetails"
	// CopyLogDetailsTypeDataBox ...
	CopyLogDetailsTypeDataBox CopyLogDetailsType = "DataBox"
	// CopyLogDetailsTypeDataBoxDisk ...
	CopyLogDetailsTypeDataBoxDisk CopyLogDetailsType = "DataBoxDisk"
	// CopyLogDetailsTypeDataBoxHeavy ...
	CopyLogDetailsTypeDataBoxHeavy CopyLogDetailsType = "DataBoxHeavy"
)

// PossibleCopyLogDetailsTypeValues returns an array of possible values for the CopyLogDetailsType const type.
func PossibleCopyLogDetailsTypeValues() []CopyLogDetailsType {
	return []CopyLogDetailsType{CopyLogDetailsTypeCopyLogDetails, CopyLogDetailsTypeDataBox, CopyLogDetailsTypeDataBoxDisk, CopyLogDetailsTypeDataBoxHeavy}
}

// CopyStatus enumerates the values for copy status.
type CopyStatus string

const (
	// Completed Data copy completed.
	Completed CopyStatus = "Completed"
	// CompletedWithErrors Data copy completed with errors.
	CompletedWithErrors CopyStatus = "CompletedWithErrors"
	// DeviceFormatted Data copy failed. The Device was formatted by user.
	DeviceFormatted CopyStatus = "DeviceFormatted"
	// DeviceMetadataModified Data copy failed. Device metadata was modified by user.
	DeviceMetadataModified CopyStatus = "DeviceMetadataModified"
	// Failed Data copy failed. No data was copied.
	Failed CopyStatus = "Failed"
	// HardwareError The Device has hit hardware issues.
	HardwareError CopyStatus = "HardwareError"
	// InProgress Data copy is in progress.
	InProgress CopyStatus = "InProgress"
	// NotReturned No copy triggered as device was not returned.
	NotReturned CopyStatus = "NotReturned"
	// NotStarted Data copy hasn't started yet.
	NotStarted CopyStatus = "NotStarted"
	// StorageAccountNotAccessible Data copy failed. Storage Account was not accessible during copy.
	StorageAccountNotAccessible CopyStatus = "StorageAccountNotAccessible"
	// UnsupportedData Data copy failed. The Device data content is not supported.
	UnsupportedData CopyStatus = "UnsupportedData"
)

// PossibleCopyStatusValues returns an array of possible values for the CopyStatus const type.
func PossibleCopyStatusValues() []CopyStatus {
	return []CopyStatus{Completed, CompletedWithErrors, DeviceFormatted, DeviceMetadataModified, Failed, HardwareError, InProgress, NotReturned, NotStarted, StorageAccountNotAccessible, UnsupportedData}
}

// DataDestinationType enumerates the values for data destination type.
type DataDestinationType string

const (
	// ManagedDisk Azure Managed disk storage.
	ManagedDisk DataDestinationType = "ManagedDisk"
	// StorageAccount Storage Accounts .
	StorageAccount DataDestinationType = "StorageAccount"
)

// PossibleDataDestinationTypeValues returns an array of possible values for the DataDestinationType const type.
func PossibleDataDestinationTypeValues() []DataDestinationType {
	return []DataDestinationType{ManagedDisk, StorageAccount}
}

// DataDestinationTypeBasicDestinationAccountDetails enumerates the values for data destination type basic
// destination account details.
type DataDestinationTypeBasicDestinationAccountDetails string

const (
	// DataDestinationTypeDestinationAccountDetails ...
	DataDestinationTypeDestinationAccountDetails DataDestinationTypeBasicDestinationAccountDetails = "DestinationAccountDetails"
	// DataDestinationTypeManagedDisk ...
	DataDestinationTypeManagedDisk DataDestinationTypeBasicDestinationAccountDetails = "ManagedDisk"
	// DataDestinationTypeStorageAccount ...
	DataDestinationTypeStorageAccount DataDestinationTypeBasicDestinationAccountDetails = "StorageAccount"
)

// PossibleDataDestinationTypeBasicDestinationAccountDetailsValues returns an array of possible values for the DataDestinationTypeBasicDestinationAccountDetails const type.
func PossibleDataDestinationTypeBasicDestinationAccountDetailsValues() []DataDestinationTypeBasicDestinationAccountDetails {
	return []DataDestinationTypeBasicDestinationAccountDetails{DataDestinationTypeDestinationAccountDetails, DataDestinationTypeManagedDisk, DataDestinationTypeStorageAccount}
}

// JobDeliveryType enumerates the values for job delivery type.
type JobDeliveryType string

const (
	// NonScheduled Non Scheduled job.
	NonScheduled JobDeliveryType = "NonScheduled"
	// Scheduled Scheduled job.
	Scheduled JobDeliveryType = "Scheduled"
)

// PossibleJobDeliveryTypeValues returns an array of possible values for the JobDeliveryType const type.
func PossibleJobDeliveryTypeValues() []JobDeliveryType {
	return []JobDeliveryType{NonScheduled, Scheduled}
}

// JobDetailsTypeEnum enumerates the values for job details type enum.
type JobDetailsTypeEnum string

const (
	// JobDetailsTypeDataBox ...
	JobDetailsTypeDataBox JobDetailsTypeEnum = "DataBox"
	// JobDetailsTypeDataBoxDisk ...
	JobDetailsTypeDataBoxDisk JobDetailsTypeEnum = "DataBoxDisk"
	// JobDetailsTypeDataBoxHeavy ...
	JobDetailsTypeDataBoxHeavy JobDetailsTypeEnum = "DataBoxHeavy"
	// JobDetailsTypeJobDetails ...
	JobDetailsTypeJobDetails JobDetailsTypeEnum = "JobDetails"
)

// PossibleJobDetailsTypeEnumValues returns an array of possible values for the JobDetailsTypeEnum const type.
func PossibleJobDetailsTypeEnumValues() []JobDetailsTypeEnum {
	return []JobDetailsTypeEnum{JobDetailsTypeDataBox, JobDetailsTypeDataBoxDisk, JobDetailsTypeDataBoxHeavy, JobDetailsTypeJobDetails}
}

// JobSecretsTypeEnum enumerates the values for job secrets type enum.
type JobSecretsTypeEnum string

const (
	// JobSecretsTypeDataBox ...
	JobSecretsTypeDataBox JobSecretsTypeEnum = "DataBox"
	// JobSecretsTypeDataBoxDisk ...
	JobSecretsTypeDataBoxDisk JobSecretsTypeEnum = "DataBoxDisk"
	// JobSecretsTypeDataBoxHeavy ...
	JobSecretsTypeDataBoxHeavy JobSecretsTypeEnum = "DataBoxHeavy"
	// JobSecretsTypeJobSecrets ...
	JobSecretsTypeJobSecrets JobSecretsTypeEnum = "JobSecrets"
)

// PossibleJobSecretsTypeEnumValues returns an array of possible values for the JobSecretsTypeEnum const type.
func PossibleJobSecretsTypeEnumValues() []JobSecretsTypeEnum {
	return []JobSecretsTypeEnum{JobSecretsTypeDataBox, JobSecretsTypeDataBoxDisk, JobSecretsTypeDataBoxHeavy, JobSecretsTypeJobSecrets}
}

// NotificationStageName enumerates the values for notification stage name.
type NotificationStageName string

const (
	// AtAzureDC Notification at device received at azure datacenter stage.
	AtAzureDC NotificationStageName = "AtAzureDC"
	// DataCopy Notification at data copy started stage.
	DataCopy NotificationStageName = "DataCopy"
	// Delivered Notification at device delivered stage.
	Delivered NotificationStageName = "Delivered"
	// DevicePrepared Notification at device prepared stage.
	DevicePrepared NotificationStageName = "DevicePrepared"
	// Dispatched Notification at device dispatched stage.
	Dispatched NotificationStageName = "Dispatched"
	// PickedUp Notification at device picked up from user stage.
	PickedUp NotificationStageName = "PickedUp"
)

// PossibleNotificationStageNameValues returns an array of possible values for the NotificationStageName const type.
func PossibleNotificationStageNameValues() []NotificationStageName {
	return []NotificationStageName{AtAzureDC, DataCopy, Delivered, DevicePrepared, Dispatched, PickedUp}
}

// OverallValidationStatus enumerates the values for overall validation status.
type OverallValidationStatus string

const (
	// AllValidToProceed Every input request is valid.
	AllValidToProceed OverallValidationStatus = "AllValidToProceed"
	// CertainInputValidationsSkipped Certain input validations skipped.
	CertainInputValidationsSkipped OverallValidationStatus = "CertainInputValidationsSkipped"
	// InputsRevisitRequired Some input requests are not valid.
	InputsRevisitRequired OverallValidationStatus = "InputsRevisitRequired"
)

// PossibleOverallValidationStatusValues returns an array of possible values for the OverallValidationStatus const type.
func PossibleOverallValidationStatusValues() []OverallValidationStatus {
	return []OverallValidationStatus{AllValidToProceed, CertainInputValidationsSkipped, InputsRevisitRequired}
}

// ShareDestinationFormatType enumerates the values for share destination format type.
type ShareDestinationFormatType string

const (
	// ShareDestinationFormatTypeAzureFile Azure storage file format.
	ShareDestinationFormatTypeAzureFile ShareDestinationFormatType = "AzureFile"
	// ShareDestinationFormatTypeBlockBlob Azure storage block blob format.
	ShareDestinationFormatTypeBlockBlob ShareDestinationFormatType = "BlockBlob"
	// ShareDestinationFormatTypeHCS Storsimple data format.
	ShareDestinationFormatTypeHCS ShareDestinationFormatType = "HCS"
	// ShareDestinationFormatTypeManagedDisk Azure Compute Disk.
	ShareDestinationFormatTypeManagedDisk ShareDestinationFormatType = "ManagedDisk"
	// ShareDestinationFormatTypePageBlob Azure storage page blob format.
	ShareDestinationFormatTypePageBlob ShareDestinationFormatType = "PageBlob"
	// ShareDestinationFormatTypeUnknownType Unknown format.
	ShareDestinationFormatTypeUnknownType ShareDestinationFormatType = "UnknownType"
)

// PossibleShareDestinationFormatTypeValues returns an array of possible values for the ShareDestinationFormatType const type.
func PossibleShareDestinationFormatTypeValues() []ShareDestinationFormatType {
	return []ShareDestinationFormatType{ShareDestinationFormatTypeAzureFile, ShareDestinationFormatTypeBlockBlob, ShareDestinationFormatTypeHCS, ShareDestinationFormatTypeManagedDisk, ShareDestinationFormatTypePageBlob, ShareDestinationFormatTypeUnknownType}
}

// SkuDisabledReason enumerates the values for sku disabled reason.
type SkuDisabledReason string

const (
	// SkuDisabledReasonCountry SKU is not available in the requested country.
	SkuDisabledReasonCountry SkuDisabledReason = "Country"
	// SkuDisabledReasonFeature Required features are not enabled for the SKU.
	SkuDisabledReasonFeature SkuDisabledReason = "Feature"
	// SkuDisabledReasonNone SKU is not disabled.
	SkuDisabledReasonNone SkuDisabledReason = "None"
	// SkuDisabledReasonNoSubscriptionInfo Subscription has not registered to Microsoft.DataBox and Service
	// does not have the subscription notification.
	SkuDisabledReasonNoSubscriptionInfo SkuDisabledReason = "NoSubscriptionInfo"
	// SkuDisabledReasonOfferType Subscription does not have required offer types for the SKU.
	SkuDisabledReasonOfferType SkuDisabledReason = "OfferType"
	// SkuDisabledReasonRegion SKU is not available to push data to the requested Azure region.
	SkuDisabledReasonRegion SkuDisabledReason = "Region"
)

// PossibleSkuDisabledReasonValues returns an array of possible values for the SkuDisabledReason const type.
func PossibleSkuDisabledReasonValues() []SkuDisabledReason {
	return []SkuDisabledReason{SkuDisabledReasonCountry, SkuDisabledReasonFeature, SkuDisabledReasonNone, SkuDisabledReasonNoSubscriptionInfo, SkuDisabledReasonOfferType, SkuDisabledReasonRegion}
}

// SkuName enumerates the values for sku name.
type SkuName string

const (
	// DataBox Databox.
	DataBox SkuName = "DataBox"
	// DataBoxDisk DataboxDisk.
	DataBoxDisk SkuName = "DataBoxDisk"
	// DataBoxHeavy DataboxHeavy.
	DataBoxHeavy SkuName = "DataBoxHeavy"
)

// PossibleSkuNameValues returns an array of possible values for the SkuName const type.
func PossibleSkuNameValues() []SkuName {
	return []SkuName{DataBox, DataBoxDisk, DataBoxHeavy}
}

// SkuNameBasicScheduleAvailabilityRequest enumerates the values for sku name basic schedule availability
// request.
type SkuNameBasicScheduleAvailabilityRequest string

const (
	// SkuNameDataBox ...
	SkuNameDataBox SkuNameBasicScheduleAvailabilityRequest = "DataBox"
	// SkuNameDataBoxDisk ...
	SkuNameDataBoxDisk SkuNameBasicScheduleAvailabilityRequest = "DataBoxDisk"
	// SkuNameDataBoxHeavy ...
	SkuNameDataBoxHeavy SkuNameBasicScheduleAvailabilityRequest = "DataBoxHeavy"
	// SkuNameScheduleAvailabilityRequest ...
	SkuNameScheduleAvailabilityRequest SkuNameBasicScheduleAvailabilityRequest = "ScheduleAvailabilityRequest"
)

// PossibleSkuNameBasicScheduleAvailabilityRequestValues returns an array of possible values for the SkuNameBasicScheduleAvailabilityRequest const type.
func PossibleSkuNameBasicScheduleAvailabilityRequestValues() []SkuNameBasicScheduleAvailabilityRequest {
	return []SkuNameBasicScheduleAvailabilityRequest{SkuNameDataBox, SkuNameDataBoxDisk, SkuNameDataBoxHeavy, SkuNameScheduleAvailabilityRequest}
}

// StageName enumerates the values for stage name.
type StageName string

const (
	// StageNameAborted Order has been aborted.
	StageNameAborted StageName = "Aborted"
	// StageNameAtAzureDC Device has been received at azure datacenter from the user.
	StageNameAtAzureDC StageName = "AtAzureDC"
	// StageNameCancelled Order has been cancelled.
	StageNameCancelled StageName = "Cancelled"
	// StageNameCompleted Order has completed.
	StageNameCompleted StageName = "Completed"
	// StageNameCompletedWithErrors Order has completed with errors.
	StageNameCompletedWithErrors StageName = "CompletedWithErrors"
	// StageNameCompletedWithWarnings Order has completed with warnings.
	StageNameCompletedWithWarnings StageName = "CompletedWithWarnings"
	// StageNameDataCopy Data copy from the device at azure datacenter.
	StageNameDataCopy StageName = "DataCopy"
	// StageNameDelivered Device has been delivered to the user of the order.
	StageNameDelivered StageName = "Delivered"
	// StageNameDeviceOrdered An order has been created.
	StageNameDeviceOrdered StageName = "DeviceOrdered"
	// StageNameDevicePrepared A device has been prepared for the order.
	StageNameDevicePrepared StageName = "DevicePrepared"
	// StageNameDispatched Device has been dispatched to the user of the order.
	StageNameDispatched StageName = "Dispatched"
	// StageNameFailedIssueDetectedAtAzureDC Order has failed due to issue detected at azure datacenter.
	StageNameFailedIssueDetectedAtAzureDC StageName = "Failed_IssueDetectedAtAzureDC"
	// StageNameFailedIssueReportedAtCustomer Order has failed due to issue reported by user.
	StageNameFailedIssueReportedAtCustomer StageName = "Failed_IssueReportedAtCustomer"
	// StageNamePickedUp Device has been picked up from user and in transit to azure datacenter.
	StageNamePickedUp StageName = "PickedUp"
	// StageNameReadyToDispatchFromAzureDC Device is ready to be handed to customer from Azure DC.
	StageNameReadyToDispatchFromAzureDC StageName = "ReadyToDispatchFromAzureDC"
	// StageNameReadyToReceiveAtAzureDC Device can be dropped off at Azure DC.
	StageNameReadyToReceiveAtAzureDC StageName = "ReadyToReceiveAtAzureDC"
)

// PossibleStageNameValues returns an array of possible values for the StageName const type.
func PossibleStageNameValues() []StageName {
	return []StageName{StageNameAborted, StageNameAtAzureDC, StageNameCancelled, StageNameCompleted, StageNameCompletedWithErrors, StageNameCompletedWithWarnings, StageNameDataCopy, StageNameDelivered, StageNameDeviceOrdered, StageNameDevicePrepared, StageNameDispatched, StageNameFailedIssueDetectedAtAzureDC, StageNameFailedIssueReportedAtCustomer, StageNamePickedUp, StageNameReadyToDispatchFromAzureDC, StageNameReadyToReceiveAtAzureDC}
}

// StageStatus enumerates the values for stage status.
type StageStatus string

const (
	// StageStatusCancelled Stage has been cancelled.
	StageStatusCancelled StageStatus = "Cancelled"
	// StageStatusCancelling Stage is cancelling.
	StageStatusCancelling StageStatus = "Cancelling"
	// StageStatusFailed Stage has failed.
	StageStatusFailed StageStatus = "Failed"
	// StageStatusInProgress Stage is in progress.
	StageStatusInProgress StageStatus = "InProgress"
	// StageStatusNone No status available yet.
	StageStatusNone StageStatus = "None"
	// StageStatusSucceeded Stage has succeeded.
	StageStatusSucceeded StageStatus = "Succeeded"
	// StageStatusSucceededWithErrors Stage has succeeded with errors.
	StageStatusSucceededWithErrors StageStatus = "SucceededWithErrors"
)

// PossibleStageStatusValues returns an array of possible values for the StageStatus const type.
func PossibleStageStatusValues() []StageStatus {
	return []StageStatus{StageStatusCancelled, StageStatusCancelling, StageStatusFailed, StageStatusInProgress, StageStatusNone, StageStatusSucceeded, StageStatusSucceededWithErrors}
}

// TransportShipmentTypes enumerates the values for transport shipment types.
type TransportShipmentTypes string

const (
	// CustomerManaged Shipment Logistics is handled by the customer.
	CustomerManaged TransportShipmentTypes = "CustomerManaged"
	// MicrosoftManaged Shipment Logistics is handled by Microsoft.
	MicrosoftManaged TransportShipmentTypes = "MicrosoftManaged"
)

// PossibleTransportShipmentTypesValues returns an array of possible values for the TransportShipmentTypes const type.
func PossibleTransportShipmentTypesValues() []TransportShipmentTypes {
	return []TransportShipmentTypes{CustomerManaged, MicrosoftManaged}
}

// ValidationCategory enumerates the values for validation category.
type ValidationCategory string

const (
	// ValidationCategoryJobCreationValidation ...
	ValidationCategoryJobCreationValidation ValidationCategory = "JobCreationValidation"
	// ValidationCategoryValidationRequest ...
	ValidationCategoryValidationRequest ValidationCategory = "ValidationRequest"
)

// PossibleValidationCategoryValues returns an array of possible values for the ValidationCategory const type.
func PossibleValidationCategoryValues() []ValidationCategory {
	return []ValidationCategory{ValidationCategoryJobCreationValidation, ValidationCategoryValidationRequest}
}

// ValidationStatus enumerates the values for validation status.
type ValidationStatus string

const (
	// ValidationStatusInvalid Validation is not successful
	ValidationStatusInvalid ValidationStatus = "Invalid"
	// ValidationStatusSkipped Validation is skipped
	ValidationStatusSkipped ValidationStatus = "Skipped"
	// ValidationStatusValid Validation is successful
	ValidationStatusValid ValidationStatus = "Valid"
)

// PossibleValidationStatusValues returns an array of possible values for the ValidationStatus const type.
func PossibleValidationStatusValues() []ValidationStatus {
	return []ValidationStatus{ValidationStatusInvalid, ValidationStatusSkipped, ValidationStatusValid}
}

// ValidationType enumerates the values for validation type.
type ValidationType string

const (
	// ValidationTypeValidateAddress ...
	ValidationTypeValidateAddress ValidationType = "ValidateAddress"
	// ValidationTypeValidateCreateOrderLimit ...
	ValidationTypeValidateCreateOrderLimit ValidationType = "ValidateCreateOrderLimit"
	// ValidationTypeValidateDataDestinationDetails ...
	ValidationTypeValidateDataDestinationDetails ValidationType = "ValidateDataDestinationDetails"
	// ValidationTypeValidatePreferences ...
	ValidationTypeValidatePreferences ValidationType = "ValidatePreferences"
	// ValidationTypeValidateSkuAvailability ...
	ValidationTypeValidateSkuAvailability ValidationType = "ValidateSkuAvailability"
	// ValidationTypeValidateSubscriptionIsAllowedToCreateJob ...
	ValidationTypeValidateSubscriptionIsAllowedToCreateJob ValidationType = "ValidateSubscriptionIsAllowedToCreateJob"
	// ValidationTypeValidationInputRequest ...
	ValidationTypeValidationInputRequest ValidationType = "ValidationInputRequest"
)

// PossibleValidationTypeValues returns an array of possible values for the ValidationType const type.
func PossibleValidationTypeValues() []ValidationType {
	return []ValidationType{ValidationTypeValidateAddress, ValidationTypeValidateCreateOrderLimit, ValidationTypeValidateDataDestinationDetails, ValidationTypeValidatePreferences, ValidationTypeValidateSkuAvailability, ValidationTypeValidateSubscriptionIsAllowedToCreateJob, ValidationTypeValidationInputRequest}
}

// ValidationTypeBasicValidationInputResponse enumerates the values for validation type basic validation input
// response.
type ValidationTypeBasicValidationInputResponse string

const (
	// ValidationTypeBasicValidationInputResponseValidationTypeValidateAddress ...
	ValidationTypeBasicValidationInputResponseValidationTypeValidateAddress ValidationTypeBasicValidationInputResponse = "ValidateAddress"
	// ValidationTypeBasicValidationInputResponseValidationTypeValidateCreateOrderLimit ...
	ValidationTypeBasicValidationInputResponseValidationTypeValidateCreateOrderLimit ValidationTypeBasicValidationInputResponse = "ValidateCreateOrderLimit"
	// ValidationTypeBasicValidationInputResponseValidationTypeValidateDataDestinationDetails ...
	ValidationTypeBasicValidationInputResponseValidationTypeValidateDataDestinationDetails ValidationTypeBasicValidationInputResponse = "ValidateDataDestinationDetails"
	// ValidationTypeBasicValidationInputResponseValidationTypeValidatePreferences ...
	ValidationTypeBasicValidationInputResponseValidationTypeValidatePreferences ValidationTypeBasicValidationInputResponse = "ValidatePreferences"
	// ValidationTypeBasicValidationInputResponseValidationTypeValidateSkuAvailability ...
	ValidationTypeBasicValidationInputResponseValidationTypeValidateSkuAvailability ValidationTypeBasicValidationInputResponse = "ValidateSkuAvailability"
	// ValidationTypeBasicValidationInputResponseValidationTypeValidateSubscriptionIsAllowedToCreateJob ...
	ValidationTypeBasicValidationInputResponseValidationTypeValidateSubscriptionIsAllowedToCreateJob ValidationTypeBasicValidationInputResponse = "ValidateSubscriptionIsAllowedToCreateJob"
	// ValidationTypeBasicValidationInputResponseValidationTypeValidationInputResponse ...
	ValidationTypeBasicValidationInputResponseValidationTypeValidationInputResponse ValidationTypeBasicValidationInputResponse = "ValidationInputResponse"
)

// PossibleValidationTypeBasicValidationInputResponseValues returns an array of possible values for the ValidationTypeBasicValidationInputResponse const type.
func PossibleValidationTypeBasicValidationInputResponseValues() []ValidationTypeBasicValidationInputResponse {
	return []ValidationTypeBasicValidationInputResponse{ValidationTypeBasicValidationInputResponseValidationTypeValidateAddress, ValidationTypeBasicValidationInputResponseValidationTypeValidateCreateOrderLimit, ValidationTypeBasicValidationInputResponseValidationTypeValidateDataDestinationDetails, ValidationTypeBasicValidationInputResponseValidationTypeValidatePreferences, ValidationTypeBasicValidationInputResponseValidationTypeValidateSkuAvailability, ValidationTypeBasicValidationInputResponseValidationTypeValidateSubscriptionIsAllowedToCreateJob, ValidationTypeBasicValidationInputResponseValidationTypeValidationInputResponse}
}

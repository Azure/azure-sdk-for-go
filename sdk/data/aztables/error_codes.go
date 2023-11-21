// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

// TableErrorCode is the error code returned by the service on failed operations. For more
// information about Table service error codes: https://learn.microsoft.com/rest/api/storageservices/table-service-error-codes
type TableErrorCode string

const (
	DuplicatePropertiesSpecified TableErrorCode = "DuplicatePropertiesSpecified"
	EntityAlreadyExists          TableErrorCode = "EntityAlreadyExists"
	EntityNotFound               TableErrorCode = "EntityNotFound"
	EntityTooLarge               TableErrorCode = "EntityTooLarge"
	HostInformationNotPresent    TableErrorCode = "HostInformationNotPresent"
	InvalidDuplicateRow          TableErrorCode = "InvalidDuplicateRow"
	InvalidInput                 TableErrorCode = "InvalidInput"
	InvalidValueType             TableErrorCode = "InvalidValueType"
	InvalidXmlDocument           TableErrorCode = "InvalidXmlDocument"
	JSONFormatNotSupported       TableErrorCode = "JsonFormatNotSupported"
	MethodNotAllowed             TableErrorCode = "MethodNotAllowed"
	NotImplemented               TableErrorCode = "NotImplemented"
	OutOfRangeInput              TableErrorCode = "OutOfRangeInput"
	PropertiesNeedValue          TableErrorCode = "PropertiesNeedValue"
	PropertyNameInvalid          TableErrorCode = "PropertyNameInvalid"
	PropertyNameTooLong          TableErrorCode = "PropertyNameTooLong"
	PropertyValueTooLarge        TableErrorCode = "PropertyValueTooLarge"
	ResourceNotFound             TableErrorCode = "ResourceNotFound"
	TableAlreadyExists           TableErrorCode = "TableAlreadyExists"
	TableBeingDeleted            TableErrorCode = "TableBeingDeleted"
	TableNotFound                TableErrorCode = "TableNotFound"
	TooManyProperties            TableErrorCode = "TooManyProperties"
	UpdateConditionNotSatisfied  TableErrorCode = "UpdateConditionNotSatisfied"
	XMethodIncorrectCount        TableErrorCode = "XMethodIncorrectCount"
	XMethodIncorrectValue        TableErrorCode = "XMethodIncorrectValue"
	XMethodNotUsingPost          TableErrorCode = "XMethodNotUsingPost"
)

// PossibleTableErrorCodeValues returns a slice of all possible TableErrorCode values
func PossibleTableErrorCodeValues() []TableErrorCode {
	return []TableErrorCode{
		DuplicatePropertiesSpecified,
		EntityAlreadyExists,
		EntityNotFound,
		EntityTooLarge,
		HostInformationNotPresent,
		InvalidDuplicateRow,
		InvalidInput,
		InvalidValueType,
		InvalidXmlDocument,
		JSONFormatNotSupported,
		MethodNotAllowed,
		NotImplemented,
		OutOfRangeInput,
		PropertiesNeedValue,
		PropertyNameInvalid,
		PropertyNameTooLong,
		PropertyValueTooLarge,
		ResourceNotFound,
		TableAlreadyExists,
		TableBeingDeleted,
		TableNotFound,
		TooManyProperties,
		UpdateConditionNotSatisfied,
		XMethodIncorrectCount,
		XMethodIncorrectValue,
		XMethodNotUsingPost,
	}
}

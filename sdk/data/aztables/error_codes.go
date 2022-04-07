// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

// TableErrorCode is the error code returned by the service on failed operations
type TableErrorCode string

const (
	AccountAlreadyExists                 TableErrorCode = "AccountAlreadyExists"
	AccountBeingCreated                  TableErrorCode = "AccountBeingCreated"
	AccountIsDisabled                    TableErrorCode = "AccountIsDisabled"
	AuthenticationFailed                 TableErrorCode = "AuthenticationFailed"
	AuthorizationFailure                 TableErrorCode = "AuthorizationFailure"
	ConditionHeadersNotSupported         TableErrorCode = "ConditionHeadersNotSupported"
	ConditionNotMet                      TableErrorCode = "ConditionNotMet"
	DuplicatePropertiesSpecified         TableErrorCode = "DuplicatePropertiesSpecified"
	EmptyMetadataKey                     TableErrorCode = "EmptyMetadataKey"
	EntityAlreadyExists                  TableErrorCode = "EntityAlreadyExists"
	EntityNotFound                       TableErrorCode = "EntityNotFound"
	EntityTooLarge                       TableErrorCode = "EntityTooLarge"
	HostInformationNotPresent            TableErrorCode = "HostInformationNotPresent"
	InsufficientAccountPermissions       TableErrorCode = "InsufficientAccountPermissions"
	InternalError                        TableErrorCode = "InternalError"
	InvalidAuthenticationInfo            TableErrorCode = "InvalidAuthenticationInfo"
	InvalidDuplicateRow                  TableErrorCode = "InvalidDuplicateRow"
	InvalidHeaderValue                   TableErrorCode = "InvalidHeaderValue"
	InvalidHTTPVerb                      TableErrorCode = "InvalidHttpVerb"
	InvalidInput                         TableErrorCode = "InvalidInput"
	InvalidMd5                           TableErrorCode = "InvalidMd5"
	InvalidMetadata                      TableErrorCode = "InvalidMetadata"
	InvalidQueryParameterValue           TableErrorCode = "InvalidQueryParameterValue"
	InvalidRange                         TableErrorCode = "InvalidRange"
	InvalidResourceName                  TableErrorCode = "InvalidResourceName"
	InvalidURI                           TableErrorCode = "InvalidUri"
	InvalidValueType                     TableErrorCode = "InvalidValueType"
	InvalidXMLDocument                   TableErrorCode = "InvalidXmlDocument"
	InvalidXMLNodeValue                  TableErrorCode = "InvalidXmlNodeValue"
	JSONFormatNotSupported               TableErrorCode = "JsonFormatNotSupported"
	Md5Mismatch                          TableErrorCode = "Md5Mismatch"
	MetadataTooLarge                     TableErrorCode = "MetadataTooLarge"
	MethodNotAllowed                     TableErrorCode = "MethodNotAllowed"
	MissingContentLengthHeader           TableErrorCode = "MissingContentLengthHeader"
	MissingRequiredHeader                TableErrorCode = "MissingRequiredHeader"
	MissingRequiredQueryParameter        TableErrorCode = "MissingRequiredQueryParameter"
	MissingRequiredXMLNode               TableErrorCode = "MissingRequiredXmlNode"
	MultipleConditionHeadersNotSupported TableErrorCode = "MultipleConditionHeadersNotSupported"
	NoAuthenticationInformation          TableErrorCode = "NoAuthenticationInformation"
	NotImplemented                       TableErrorCode = "NotImplemented"
	OperationTimedOut                    TableErrorCode = "OperationTimedOut"
	OutOfRangeInput                      TableErrorCode = "OutOfRangeInput"
	OutOfRangeQueryParameterValue        TableErrorCode = "OutOfRangeQueryParameterValue"
	PropertiesNeedValue                  TableErrorCode = "PropertiesNeedValue"
	PropertyNameInvalid                  TableErrorCode = "PropertyNameInvalid"
	PropertyNameTooLong                  TableErrorCode = "PropertyNameTooLong"
	PropertyValueTooLarge                TableErrorCode = "PropertyValueTooLarge"
	RequestBodyTooLarge                  TableErrorCode = "RequestBodyTooLarge"
	RequestUrlFailedToParse              TableErrorCode = "RequestUrlFailedToParse"
	ResourceAlreadyExists                TableErrorCode = "ResourceAlreadyExists"
	ResourceNotFound                     TableErrorCode = "ResourceNotFound"
	ResourceTypeMismatch                 TableErrorCode = "ResourceTypeMismatch"
	ServerBusy                           TableErrorCode = "ServerBusy"
	TableAlreadyExists                   TableErrorCode = "TableAlreadyExists"
	TableBeingDeleted                    TableErrorCode = "TableBeingDeleted"
	TableNotFound                        TableErrorCode = "TableNotFound"
	TooManyProperties                    TableErrorCode = "TooManyProperties"
	UnsupportedHeader                    TableErrorCode = "UnsupportedHeader"
	UnsupportedHTTPVerb                  TableErrorCode = "UnsupportedHttpVerb"
	UnsupportedQueryParameter            TableErrorCode = "UnsupportedQueryParameter"
	UnsupportedXMLNode                   TableErrorCode = "UnsupportedXmlNode"
	UpdateConditionNotSatisfied          TableErrorCode = "UpdateConditionNotSatisfied"
	XMethodIncorrectCount                TableErrorCode = "XMethodIncorrectCount"
	XMethodIncorrectValue                TableErrorCode = "XMethodIncorrectValue"
	XMethodNotUsingPost                  TableErrorCode = "XMethodNotUsingPost"
)

// PossibleTableErrorCodeValues returns a slice of all possible TableErrorCode values
func PossibleTableErrorCodeValues() []TableErrorCode {
	return []TableErrorCode{
		AccountAlreadyExists,
		AccountBeingCreated,
		AccountIsDisabled,
		AuthenticationFailed,
		AuthorizationFailure,
		ConditionHeadersNotSupported,
		ConditionNotMet,
		DuplicatePropertiesSpecified,
		EmptyMetadataKey,
		EntityAlreadyExists,
		EntityNotFound,
		EntityTooLarge,
		HostInformationNotPresent,
		InsufficientAccountPermissions,
		InternalError,
		InvalidAuthenticationInfo,
		InvalidDuplicateRow,
		InvalidHeaderValue,
		InvalidHTTPVerb,
		InvalidInput,
		InvalidMd5,
		InvalidMetadata,
		InvalidQueryParameterValue,
		InvalidRange,
		InvalidResourceName,
		InvalidURI,
		InvalidValueType,
		InvalidXMLDocument,
		InvalidXMLNodeValue,
		JSONFormatNotSupported,
		Md5Mismatch,
		MetadataTooLarge,
		MethodNotAllowed,
		MissingContentLengthHeader,
		MissingRequiredHeader,
		MissingRequiredQueryParameter,
		MissingRequiredXMLNode,
		MultipleConditionHeadersNotSupported,
		NoAuthenticationInformation,
		NotImplemented,
		OperationTimedOut,
		OutOfRangeInput,
		OutOfRangeQueryParameterValue,
		PropertiesNeedValue,
		PropertyNameInvalid,
		PropertyNameTooLong,
		PropertyValueTooLarge,
		RequestBodyTooLarge,
		RequestUrlFailedToParse,
		ResourceAlreadyExists,
		ResourceNotFound,
		ResourceTypeMismatch,
		ServerBusy,
		TableAlreadyExists,
		TableBeingDeleted,
		TableNotFound,
		TooManyProperties,
		UnsupportedHeader,
		UnsupportedHTTPVerb,
		UnsupportedQueryParameter,
		UnsupportedXMLNode,
		UpdateConditionNotSatisfied,
		XMethodIncorrectCount,
		XMethodIncorrectValue,
		XMethodNotUsingPost,
	}
}

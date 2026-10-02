// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package queueerror

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/internal/generated"
)

// HasCode returns true if the provided error is an *azcore.ResponseError
// with its ErrorCode field equal to one of the specified Codes.
func HasCode(err error, codes ...Code) bool {
	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return false
	}

	for _, code := range codes {
		if respErr.ErrorCode == string(code) {
			return true
		}
	}

	return false
}

// Code - Error codes returned by the service
type Code = generated.StorageErrorCode

const (
	AccountAlreadyExists                 Code = generated.StorageErrorCodeAccountAlreadyExists
	AccountBeingCreated                  Code = generated.StorageErrorCodeAccountBeingCreated
	AccountIsDisabled                    Code = generated.StorageErrorCodeAccountIsDisabled
	AuthenticationFailed                 Code = generated.StorageErrorCodeAuthenticationFailed
	AuthorizationFailure                 Code = generated.StorageErrorCodeAuthorizationFailure
	AuthorizationPermissionMismatch      Code = generated.StorageErrorCodeAuthorizationPermissionMismatch
	AuthorizationProtocolMismatch        Code = generated.StorageErrorCodeAuthorizationProtocolMismatch
	AuthorizationResourceTypeMismatch    Code = generated.StorageErrorCodeAuthorizationResourceTypeMismatch
	AuthorizationServiceMismatch         Code = generated.StorageErrorCodeAuthorizationServiceMismatch
	AuthorizationSourceIPMismatch        Code = generated.StorageErrorCodeAuthorizationSourceIPMismatch
	ConditionHeadersNotSupported         Code = generated.StorageErrorCodeConditionHeadersNotSupported
	ConditionNotMet                      Code = generated.StorageErrorCodeConditionNotMet
	EmptyMetadataKey                     Code = generated.StorageErrorCodeEmptyMetadataKey
	FeatureVersionMismatch               Code = generated.StorageErrorCodeFeatureVersionMismatch
	InsufficientAccountPermissions       Code = generated.StorageErrorCodeInsufficientAccountPermissions
	InternalError                        Code = generated.StorageErrorCodeInternalError
	InvalidAuthenticationInfo            Code = generated.StorageErrorCodeInvalidAuthenticationInfo
	InvalidHTTPVerb                      Code = generated.StorageErrorCodeInvalidHTTPVerb
	InvalidHeaderValue                   Code = generated.StorageErrorCodeInvalidHeaderValue
	InvalidInput                         Code = generated.StorageErrorCodeInvalidInput
	InvalidMD5                           Code = generated.StorageErrorCodeInvalidMD5
	InvalidMarker                        Code = generated.StorageErrorCodeInvalidMarker
	InvalidMetadata                      Code = generated.StorageErrorCodeInvalidMetadata
	InvalidQueryParameterValue           Code = generated.StorageErrorCodeInvalidQueryParameterValue
	InvalidRange                         Code = generated.StorageErrorCodeInvalidRange
	InvalidResourceName                  Code = generated.StorageErrorCodeInvalidResourceName
	InvalidURI                           Code = generated.StorageErrorCodeInvalidURI
	InvalidXMLDocument                   Code = generated.StorageErrorCodeInvalidXMLDocument
	InvalidXMLNodeValue                  Code = generated.StorageErrorCodeInvalidXMLNodeValue
	MD5Mismatch                          Code = generated.StorageErrorCodeMD5Mismatch
	MessageNotFound                      Code = generated.StorageErrorCodeMessageNotFound
	MessageTooLarge                      Code = generated.StorageErrorCodeMessageTooLarge
	MetadataTooLarge                     Code = generated.StorageErrorCodeMetadataTooLarge
	MissingContentLengthHeader           Code = generated.StorageErrorCodeMissingContentLengthHeader
	MissingRequiredHeader                Code = generated.StorageErrorCodeMissingRequiredHeader
	MissingRequiredQueryParameter        Code = generated.StorageErrorCodeMissingRequiredQueryParameter
	MissingRequiredXMLNode               Code = generated.StorageErrorCodeMissingRequiredXMLNode
	MultipleConditionHeadersNotSupported Code = generated.StorageErrorCodeMultipleConditionHeadersNotSupported
	OperationTimedOut                    Code = generated.StorageErrorCodeOperationTimedOut
	OutOfRangeInput                      Code = generated.StorageErrorCodeOutOfRangeInput
	OutOfRangeQueryParameterValue        Code = generated.StorageErrorCodeOutOfRangeQueryParameterValue
	PopReceiptMismatch                   Code = generated.StorageErrorCodePopReceiptMismatch
	QueueAlreadyExists                   Code = generated.StorageErrorCodeQueueAlreadyExists
	QueueBeingDeleted                    Code = generated.StorageErrorCodeQueueBeingDeleted
	QueueDisabled                        Code = generated.StorageErrorCodeQueueDisabled
	QueueNotEmpty                        Code = generated.StorageErrorCodeQueueNotEmpty
	QueueNotFound                        Code = generated.StorageErrorCodeQueueNotFound
	RequestBodyTooLarge                  Code = generated.StorageErrorCodeRequestBodyTooLarge
	RequestURLFailedToParse              Code = generated.StorageErrorCodeRequestURLFailedToParse
	ResourceAlreadyExists                Code = generated.StorageErrorCodeResourceAlreadyExists
	ResourceNotFound                     Code = generated.StorageErrorCodeResourceNotFound
	ResourceTypeMismatch                 Code = generated.StorageErrorCodeResourceTypeMismatch
	ServerBusy                           Code = generated.StorageErrorCodeServerBusy
	UnsupportedHTTPVerb                  Code = generated.StorageErrorCodeUnsupportedHTTPVerb
	UnsupportedHeader                    Code = generated.StorageErrorCodeUnsupportedHeader
	UnsupportedQueryParameter            Code = generated.StorageErrorCodeUnsupportedQueryParameter
	UnsupportedXMLNode                   Code = generated.StorageErrorCodeUnsupportedXMLNode
)

var (
	// MissingSharedKeyCredential - Error is returned when SAS URL is being created without SharedKeyCredential.
	//nolint:staticcheck // ST1012: Renaming these errors would be a breaking change, so suppressing linter warning.
	MissingSharedKeyCredential = errors.New("SAS can only be signed with a SharedKeyCredential")
)

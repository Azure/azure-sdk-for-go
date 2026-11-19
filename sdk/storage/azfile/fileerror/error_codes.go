// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fileerror

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
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
	AccountAlreadyExists                               Code = "AccountAlreadyExists"
	AccountBeingCreated                                Code = "AccountBeingCreated"
	AccountIsDisabled                                  Code = "AccountIsDisabled"
	AuthenticationFailed                               Code = "AuthenticationFailed"
	AuthorizationFailure                               Code = "AuthorizationFailure"
	AuthorizationPermissionMismatch                    Code = "AuthorizationPermissionMismatch"
	AuthorizationProtocolMismatch                      Code = "AuthorizationProtocolMismatch"
	AuthorizationResourceTypeMismatch                  Code = "AuthorizationResourceTypeMismatch"
	AuthorizationServiceMismatch                       Code = "AuthorizationServiceMismatch"
	AuthorizationSourceIPMismatch                      Code = "AuthorizationSourceIPMismatch"
	CannotDeleteFileOrDirectory                        Code = "CannotDeleteFileOrDirectory"
	ClientCacheFlushDelay                              Code = "ClientCacheFlushDelay"
	ConditionHeadersNotSupported                       Code = "ConditionHeadersNotSupported"
	ConditionNotMet                                    Code = "ConditionNotMet"
	DeletePending                                      Code = "DeletePending"
	DirectoryNotEmpty                                  Code = "DirectoryNotEmpty"
	EmptyMetadataKey                                   Code = "EmptyMetadataKey"
	FeatureVersionMismatch                             Code = "FeatureVersionMismatch"
	FileLockConflict                                   Code = "FileLockConflict"
	InsufficientAccountPermissions                     Code = "InsufficientAccountPermissions"
	InternalError                                      Code = "InternalError"
	InvalidAuthenticationInfo                          Code = "InvalidAuthenticationInfo"
	InvalidFileOrDirectoryPathName                     Code = "InvalidFileOrDirectoryPathName"
	InvalidHTTPVerb                                    Code = "InvalidHttpVerb"
	InvalidHeaderValue                                 Code = "InvalidHeaderValue"
	InvalidInput                                       Code = "InvalidInput"
	InvalidMD5                                         Code = "InvalidMd5"
	InvalidMetadata                                    Code = "InvalidMetadata"
	InvalidQueryParameterValue                         Code = "InvalidQueryParameterValue"
	InvalidRange                                       Code = "InvalidRange"
	InvalidResourceName                                Code = "InvalidResourceName"
	InvalidURI                                         Code = "InvalidUri"
	InvalidXMLDocument                                 Code = "InvalidXmlDocument"
	InvalidXMLNodeValue                                Code = "InvalidXmlNodeValue"
	MD5Mismatch                                        Code = "Md5Mismatch"
	MetadataTooLarge                                   Code = "MetadataTooLarge"
	MissingContentLengthHeader                         Code = "MissingContentLengthHeader"
	MissingRequiredHeader                              Code = "MissingRequiredHeader"
	MissingRequiredQueryParameter                      Code = "MissingRequiredQueryParameter"
	MissingRequiredXMLNode                             Code = "MissingRequiredXmlNode"
	MultipleConditionHeadersNotSupported               Code = "MultipleConditionHeadersNotSupported"
	OperationTimedOut                                  Code = "OperationTimedOut"
	OutOfRangeInput                                    Code = "OutOfRangeInput"
	OutOfRangeQueryParameterValue                      Code = "OutOfRangeQueryParameterValue"
	ParentNotFound                                     Code = "ParentNotFound"
	ReadOnlyAttribute                                  Code = "ReadOnlyAttribute"
	RequestBodyTooLarge                                Code = "RequestBodyTooLarge"
	RequestURLFailedToParse                            Code = "RequestUrlFailedToParse"
	ResourceAlreadyExists                              Code = "ResourceAlreadyExists"
	ResourceNotFound                                   Code = "ResourceNotFound"
	ResourceTypeMismatch                               Code = "ResourceTypeMismatch"
	ServerBusy                                         Code = "ServerBusy"
	ShareAlreadyExists                                 Code = "ShareAlreadyExists"
	ShareBeingDeleted                                  Code = "ShareBeingDeleted"
	ShareDisabled                                      Code = "ShareDisabled"
	ShareHasSnapshots                                  Code = "ShareHasSnapshots"
	ShareNotFound                                      Code = "ShareNotFound"
	ShareSizeLimitReached                              Code = "ShareSizeLimitReached"
	ShareSnapshotCountExceeded                         Code = "ShareSnapshotCountExceeded"
	ShareSnapshotInProgress                            Code = "ShareSnapshotInProgress"
	ShareSnapshotOperationNotSupported                 Code = "ShareSnapshotOperationNotSupported"
	SharingViolation                                   Code = "SharingViolation"
	UnsupportedHTTPVerb                                Code = "UnsupportedHttpVerb"
	UnsupportedHeader                                  Code = "UnsupportedHeader"
	UnsupportedQueryParameter                          Code = "UnsupportedQueryParameter"
	UnsupportedXMLNode                                 Code = "UnsupportedXmlNode"
	FileOAuthManagementAPIRestrictedToSRP              Code = "FileOAuthManagementApiRestrictedToSrp"
	FileShareProvisionedBandwidthInvalid               Code = "FileShareProvisionedBandwidthInvalid"
	FileShareProvisionedIopsInvalid                    Code = "FileShareProvisionedIopsInvalid"
	FileShareProvisionedStorageInvalid                 Code = "FileShareProvisionedStorageInvalid"
	TotalSharesCountExceedsAccountLimit                Code = "TotalSharesCountExceedsAccountLimit"
	TotalSharesProvisionedBandwidthExceedsAccountLimit Code = "TotalSharesProvisionedBandwidthExceedsAccountLimit"
	TotalSharesProvisionedCapacityExceedsAccountLimit  Code = "TotalSharesProvisionedCapacityExceedsAccountLimit"
	TotalSharesProvisionedIopsExceedsAccountLimit      Code = "TotalSharesProvisionedIopsExceedsAccountLimit"
)

// nolint:staticcheck // ST1012: Renaming would be a breaking change, so suppressing linter warning.
var (
	// MissingSharedKeyCredential - Error is returned when SAS URL is being created without SharedKeyCredential.
	MissingSharedKeyCredential = errors.New("SAS can only be signed with a SharedKeyCredential")

	// UnsupportedChecksum - Error is returned when a precomputed checksum is used with multi-chunk uploads.
	UnsupportedChecksum = errors.New("for multi-part uploads, user generated checksums cannot be validated")
)

// nolint:staticcheck // ST1012: Named without the Err prefix to match the other errors in this package.
var (
	// FileIDNotSupported is returned when an operation that requires the path of a file or directory is called
	// on a client that addresses the resource by its file ID, e.g. one created by share.Client.NewFileClientByFileID.
	// The request is not sent to the service. Use a client constructed with the path of the resource instead.
	FileIDNotSupported = errors.New("operation is not supported when the client addresses the resource by file ID; use a client constructed with the path of the resource instead")

	// FileIDRequired is returned when an operation that is only supported on a client that addresses the
	// resource by its file ID is called on a client that addresses it by path. The request is not sent to the
	// service. Use share.Client.NewFileClientByFileID to create a client that addresses the file by its file ID.
	FileIDRequired = errors.New("operation is only supported when the client addresses the resource by file ID; use share.Client.NewFileClientByFileID to create one")
)

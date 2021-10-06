//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package filesystem

import original "github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"

const (
	DefaultAdlsFileSystemDNSSuffix = original.DefaultAdlsFileSystemDNSSuffix
)

type AppendModeType = original.AppendModeType

const (
	Autocreate AppendModeType = original.Autocreate
)

type Exception = original.Exception

const (
	ExceptionAccessControlException        Exception = original.ExceptionAccessControlException
	ExceptionAdlsRemoteException           Exception = original.ExceptionAdlsRemoteException
	ExceptionBadOffsetException            Exception = original.ExceptionBadOffsetException
	ExceptionFileAlreadyExistsException    Exception = original.ExceptionFileAlreadyExistsException
	ExceptionFileNotFoundException         Exception = original.ExceptionFileNotFoundException
	ExceptionIllegalArgumentException      Exception = original.ExceptionIllegalArgumentException
	ExceptionIOException                   Exception = original.ExceptionIOException
	ExceptionRuntimeException              Exception = original.ExceptionRuntimeException
	ExceptionSecurityException             Exception = original.ExceptionSecurityException
	ExceptionThrottledException            Exception = original.ExceptionThrottledException
	ExceptionUnsupportedOperationException Exception = original.ExceptionUnsupportedOperationException
)

type ExpiryOptionType = original.ExpiryOptionType

const (
	Absolute               ExpiryOptionType = original.Absolute
	NeverExpire            ExpiryOptionType = original.NeverExpire
	RelativeToCreationDate ExpiryOptionType = original.RelativeToCreationDate
	RelativeToNow          ExpiryOptionType = original.RelativeToNow
)

type FileType = original.FileType

const (
	DIRECTORY FileType = original.DIRECTORY
	FILE      FileType = original.FILE
)

type SyncFlag = original.SyncFlag

const (
	CLOSE    SyncFlag = original.CLOSE
	DATA     SyncFlag = original.DATA
	METADATA SyncFlag = original.METADATA
)

type ACLStatus = original.ACLStatus
type ACLStatusResult = original.ACLStatusResult
type AdlsAccessControlException = original.AdlsAccessControlException
type AdlsBadOffsetException = original.AdlsBadOffsetException
type AdlsError = original.AdlsError
type AdlsFileAlreadyExistsException = original.AdlsFileAlreadyExistsException
type AdlsFileNotFoundException = original.AdlsFileNotFoundException
type AdlsIOException = original.AdlsIOException
type AdlsIllegalArgumentException = original.AdlsIllegalArgumentException
type AdlsRemoteException = original.AdlsRemoteException
type AdlsRuntimeException = original.AdlsRuntimeException
type AdlsSecurityException = original.AdlsSecurityException
type AdlsThrottledException = original.AdlsThrottledException
type AdlsUnsupportedOperationException = original.AdlsUnsupportedOperationException
type BaseClient = original.BaseClient
type BasicAdlsRemoteException = original.BasicAdlsRemoteException
type Client = original.Client
type ContentSummary = original.ContentSummary
type ContentSummaryResult = original.ContentSummaryResult
type FileOperationResult = original.FileOperationResult
type FileStatusProperties = original.FileStatusProperties
type FileStatusResult = original.FileStatusResult
type FileStatuses = original.FileStatuses
type FileStatusesResult = original.FileStatusesResult
type ReadCloser = original.ReadCloser

func New() BaseClient {
	return original.New()
}
func NewClient() Client {
	return original.NewClient()
}
func NewWithoutDefaults(adlsFileSystemDNSSuffix string) BaseClient {
	return original.NewWithoutDefaults(adlsFileSystemDNSSuffix)
}
func PossibleAppendModeTypeValues() []AppendModeType {
	return original.PossibleAppendModeTypeValues()
}
func PossibleExceptionValues() []Exception {
	return original.PossibleExceptionValues()
}
func PossibleExpiryOptionTypeValues() []ExpiryOptionType {
	return original.PossibleExpiryOptionTypeValues()
}
func PossibleFileTypeValues() []FileType {
	return original.PossibleFileTypeValues()
}
func PossibleSyncFlagValues() []SyncFlag {
	return original.PossibleSyncFlagValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}

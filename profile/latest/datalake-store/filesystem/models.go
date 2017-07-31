package filesystem

import (
	 original "github.com/Azure/azure-sdk-for-go/service/datalake-store/2016-11-01/filesystem"
)

type (
	 GroupClient = original.GroupClient
	 AppendModeType = original.AppendModeType
	 ExpiryOptionType = original.ExpiryOptionType
	 FileType = original.FileType
	 SyncFlag = original.SyncFlag
	 ACLStatus = original.ACLStatus
	 ACLStatusResult = original.ACLStatusResult
	 AdlsAccessControlException = original.AdlsAccessControlException
	 AdlsBadOffsetException = original.AdlsBadOffsetException
	 AdlsError = original.AdlsError
	 AdlsFileAlreadyExistsException = original.AdlsFileAlreadyExistsException
	 AdlsFileNotFoundException = original.AdlsFileNotFoundException
	 AdlsIllegalArgumentException = original.AdlsIllegalArgumentException
	 AdlsIOException = original.AdlsIOException
	 AdlsRemoteException = original.AdlsRemoteException
	 AdlsRuntimeException = original.AdlsRuntimeException
	 AdlsSecurityException = original.AdlsSecurityException
	 AdlsThrottledException = original.AdlsThrottledException
	 AdlsUnsupportedOperationException = original.AdlsUnsupportedOperationException
	 ContentSummary = original.ContentSummary
	 ContentSummaryResult = original.ContentSummaryResult
	 FileOperationResult = original.FileOperationResult
	 FileStatuses = original.FileStatuses
	 FileStatusesResult = original.FileStatusesResult
	 FileStatusProperties = original.FileStatusProperties
	 FileStatusResult = original.FileStatusResult
	 ReadCloser = original.ReadCloser
	 ManagementClient = original.ManagementClient
)

const (
	 Autocreate = original.Autocreate
	 Absolute = original.Absolute
	 NeverExpire = original.NeverExpire
	 RelativeToCreationDate = original.RelativeToCreationDate
	 RelativeToNow = original.RelativeToNow
	 DIRECTORY = original.DIRECTORY
	 FILE = original.FILE
	 CLOSE = original.CLOSE
	 DATA = original.DATA
	 METADATA = original.METADATA
	 DefaultAdlsFileSystemDNSSuffix = original.DefaultAdlsFileSystemDNSSuffix
)

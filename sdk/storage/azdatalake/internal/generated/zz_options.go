// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package generated

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"time"
)

// CPKInfo contains a group of parameters for the PathClient.Create method.
type CPKInfo struct {
	// The algorithm used to produce the encryption key hash. Currently, the only accepted value is "AES256". Must be provided
	// if the x-ms-encryption-key header is provided.
	EncryptionAlgorithm *EncryptionAlgorithmType

	// Optional. Specifies the encryption key to use to encrypt the data provided in the request. If not specified, encryption
	// is performed with the root account encryption key. For more information, see
	// Encryption at Rest for Azure Storage Services.
	EncryptionKey *string

	// The SHA-256 hash of the provided encryption key. Must be provided if the x-ms-encryption-key header is provided.
	EncryptionKeySHA256 *string
}

// FileSystemClientCreateOptions contains the optional parameters for the FileSystemClient.Create method.
type FileSystemClientCreateOptions struct {
	// Optional. User-defined properties to be stored with the filesystem, in the format of a comma-separated list of name and
	// value pairs "n1=v1, n2=v2, …", where each value is a base64 encoded string. Note
	// that the string may only contain ASCII characters in the ISO-8859-1 character set. If the filesystem exists, any properties
	// not included in the list will be removed. All properties are removed if the
	// header is omitted. To merge new and existing properties, first get all existing properties and the current E-Tag, then
	// make a conditional request with the E-Tag and include values for all properties.
	Properties *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// FileSystemClientDeleteOptions contains the optional parameters for the FileSystemClient.Delete method.
type FileSystemClientDeleteOptions struct {
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// FileSystemClientGetPropertiesOptions contains the optional parameters for the FileSystemClient.GetProperties method.
type FileSystemClientGetPropertiesOptions struct {
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// FileSystemClientListBlobHierarchySegmentOptions contains the optional parameters for the FileSystemClient.NewListBlobHierarchySegmentPager
// method.
type FileSystemClientListBlobHierarchySegmentOptions struct {
	// When the request includes this parameter, the operation returns a PathPrefix element in the response body that acts as
	// a placeholder for all blobs whose names begin with the same substring up to the
	// appearance of the delimiter character. The delimiter may be a single character or a string.
	Delimiter *string

	// Include this parameter to specify one or more datasets to include in the response.
	Include []ListBlobsIncludeItem

	// A string value that identifies the portion of the list of containers to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing
	// operation did not return all containers remaining to be listed with the current page. The NextMarker value can be used
	// as the value for the marker parameter in a subsequent call to request the next
	// page of list items. The marker value is opaque to the client.
	Marker *string

	// An optional value that specifies the maximum number of items to return. If omitted or greater than 5,000, the response
	// will include up to 5,000 items.
	MaxResults *int32

	// Filters results to filesystems within the specified prefix.
	Prefix *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// Include this parameter to specify one or more datasets to include in the response.. Specifying any value will set the value
	// to deleted.
	Showonly *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// FileSystemClientListPathsOptions contains the optional parameters for the FileSystemClient.NewListPathsPager method.
type FileSystemClientListPathsOptions struct {
	// Optional. When deleting a directory, the number of paths that are deleted with each invocation is limited. If the number
	// of paths to be deleted exceeds this limit, a continuation token is returned in
	// this response header. When a continuation token is returned in the response, it must be specified in a subsequent invocation
	// of the delete operation to continue deleting the directory.
	Continuation *string

	// An optional value that specifies the maximum number of items to return. If omitted or greater than 5,000, the response
	// will include up to 5,000 items.
	MaxResults *int32

	// Optional. Filters results to paths within the specified directory. An error occurs if the directory does not exist.
	Path *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32

	// Optional. Valid only when Hierarchical Namespace is enabled for the account. If "true", the user identity values returned
	// in the x-ms-owner, x-ms-group, and x-ms-acl response headers will be
	// transformed from Azure Active Directory Object IDs to User Principal Names. If "false", the values will be returned as
	// Azure Active Directory Object IDs. The default value is false. Note that group
	// and application Object IDs are not translated because they do not have unique friendly names.
	Upn *bool
}

// FileSystemClientSetPropertiesOptions contains the optional parameters for the FileSystemClient.SetProperties method.
type FileSystemClientSetPropertiesOptions struct {
	// Optional. User-defined properties to be stored with the filesystem, in the format of a comma-separated list of name and
	// value pairs "n1=v1, n2=v2, …", where each value is a base64 encoded string. Note
	// that the string may only contain ASCII characters in the ISO-8859-1 character set. If the filesystem exists, any properties
	// not included in the list will be removed. All properties are removed if the
	// header is omitted. To merge new and existing properties, first get all existing properties and the current E-Tag, then
	// make a conditional request with the E-Tag and include values for all properties.
	Properties *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// LeaseAccessConditions contains a group of parameters for the PathClient.Create method.
type LeaseAccessConditions struct {
	// If specified, the operation only succeeds if the resource's lease is active and matches this ID.
	LeaseID *string
}

// ModifiedAccessConditions contains a group of parameters for the FileSystemClient.SetProperties method.
type ModifiedAccessConditions struct {
	// Specify an ETag value to operate only on blobs with a matching value.
	IfMatch *azcore.ETag

	// Specify this header value to operate only on a blob if it has been modified since the specified date/time.
	IfModifiedSince *time.Time

	// Specify an ETag value to operate only on blobs without a matching value.
	IfNoneMatch *azcore.ETag

	// Specify this header value to operate only on a blob if it has not been modified since the specified date/time.
	IfUnmodifiedSince *time.Time
}

// PathClientAppendDataOptions contains the optional parameters for the PathClient.AppendData method.
type PathClientAppendDataOptions struct {
	// Required for "Append Data" and "Flush Data". Must be 0 for "Flush Data". Must be the length of the request content in bytes
	// for "Append Data".
	ContentLength *int64

	// If file should be flushed after the append
	Flush *bool

	// Optional. If "acquire" it will acquire the lease. If "auto-renew" it will renew the lease. If "release" it will release
	// the lease only on flush. If "acquire-release" it will acquire & complete the
	// operation & release the lease once operation is done.
	LeaseAction *LeaseAction

	// The lease duration is required to acquire a lease, and specifies the duration of the lease in seconds. The lease duration
	// must be between 15 and 60 seconds or -1 for infinite lease.
	LeaseDuration *int64

	// This parameter allows the caller to upload data in parallel and control the order in which it is appended to the file.
	// It is required when uploading data to be appended to the file and when flushing
	// previously uploaded data to the file. The value must be the position where the data is to be appended. Uploaded data is
	// not immediately flushed, or written, to the file. To flush, the previously
	// uploaded data must be contiguous, the position parameter must be specified and equal to the length of the file after all
	// data has been written, and there must not be a request entity body included
	// with the request.
	Position *int64

	// Proposed lease ID, in a GUID string format. The Blob service returns 400 (Invalid request) if the proposed lease ID is
	// not in the correct format. See Guid Constructor (String) for a list of valid GUID
	// string formats.
	ProposedLeaseID *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// Required if the request body is a structured message. Specifies the message schema version and properties.
	StructuredBodyType *string

	// Required if the request body is a structured message. Specifies the length of the blob/file content inside the message
	// body. Will always be smaller than Content-Length.
	StructuredContentLength *int64

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32

	// Specify the transactional crc64 for the body, to be validated by the service.
	TransactionalContentCRC64 []byte
}

// PathClientCreateOptions contains the optional parameters for the PathClient.Create method.
type PathClientCreateOptions struct {
	// Sets POSIX access control rights on files and directories. The value is a comma-separated list of access control entries.
	// Each access control entry (ACE) consists of a scope, a type, a user or group
	// identifier, and permissions in the format "[scope:][type]:[id]:[permissions]".
	ACL *string

	// Optional. When deleting a directory, the number of paths that are deleted with each invocation is limited. If the number
	// of paths to be deleted exceeds this limit, a continuation token is returned in
	// this response header. When a continuation token is returned in the response, it must be specified in a subsequent invocation
	// of the delete operation to continue deleting the directory.
	Continuation *string

	// Specifies the encryption context to set on the file.
	EncryptionContext *string

	// The time to set the blob to expiry
	ExpiresOn *string

	// Required. Indicates mode of the expiry time
	ExpiryOptions *PathExpiryOptions

	// Optional. The owning group of the blob or directory.
	Group *string

	// The lease duration is required to acquire a lease, and specifies the duration of the lease in seconds. The lease duration
	// must be between 15 and 60 seconds or -1 for infinite lease.
	LeaseDuration *int64

	// Optional. Valid only when namespace is enabled. This parameter determines the behavior of the rename operation. The value
	// must be "legacy" or "posix", and the default value will be "posix".
	Mode *PathRenameMode

	// Optional. The owner of the blob or directory.
	Owner *string

	// Optional and only valid if Hierarchical Namespace is enabled for the account. Sets POSIX access permissions for the file
	// owner, the file owning group, and others. Each class may be granted read,
	// write, or execute permission. The sticky bit is also supported. Both symbolic (rwxrw-rw-) and 4-digit octal notation (e.g.
	// 0766) are supported.
	Permissions *string

	// Optional. User-defined properties to be stored with the filesystem, in the format of a comma-separated list of name and
	// value pairs "n1=v1, n2=v2, …", where each value is a base64 encoded string. Note
	// that the string may only contain ASCII characters in the ISO-8859-1 character set. If the filesystem exists, any properties
	// not included in the list will be removed. All properties are removed if the
	// header is omitted. To merge new and existing properties, first get all existing properties and the current E-Tag, then
	// make a conditional request with the E-Tag and include values for all properties.
	Properties *string

	// Proposed lease ID, in a GUID string format. The Blob service returns 400 (Invalid request) if the proposed lease ID is
	// not in the correct format. See Guid Constructor (String) for a list of valid GUID
	// string formats.
	ProposedLeaseID *string

	// An optional file or directory to be renamed. The value must have the following format: "/{filesystem}/{path}". If "x-ms-properties"
	// is specified, the properties will overwrite the existing properties;
	// otherwise, the existing properties will be preserved. This value must be a URL percent-encoded string. Note that the string
	// may only contain ASCII characters in the ISO-8859-1 character set.
	RenameSource *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// Required only for Create File and Create Directory. The value must be "file" or "directory".
	Resource *PathResourceType

	// A lease ID for the source path. If specified, the source path must have an active lease and the lease ID must match.
	SourceLeaseID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32

	// Optional and only valid if Hierarchical Namespace is enabled for the account. When creating a file or directory and the
	// parent folder does not have a default ACL, the umask restricts the permissions
	// of the file or directory to be created. The resulting permission is given by p bitwise and not u, where p is the permission
	// and u is the umask. For example, if p is 0777 and u is 0057, then the
	// resulting permission is 0720. The default permission is 0777 for a directory and 0666 for a file. The default umask is
	// 0027. The umask must be specified in 4-digit octal notation (e.g. 0766).
	Umask *string
}

// PathClientDeleteOptions contains the optional parameters for the PathClient.Delete method.
type PathClientDeleteOptions struct {
	// Optional. When deleting a directory, the number of paths that are deleted with each invocation is limited. If the number
	// of paths to be deleted exceeds this limit, a continuation token is returned in
	// this response header. When a continuation token is returned in the response, it must be specified in a subsequent invocation
	// of the delete operation to continue deleting the directory.
	Continuation *string

	// If true, paginated behavior will be seen. Pagination is for the recursive ACL checks as a POSIX requirement in the server
	// and Delete in an atomic operation once the ACL checks are completed. If false
	// or missing, normal default behavior will kick in, which may timeout in case of very large directories due to recursive
	// ACL checks. This new parameter is introduced for backward compatibility.
	Paginated *bool

	// Required
	Recursive *bool

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// PathClientFlushDataOptions contains the optional parameters for the PathClient.FlushData method.
type PathClientFlushDataOptions struct {
	// Azure Storage Events allow applications to receive notifications when files change. When Azure Storage Events are enabled,
	// a file changed event is raised. This event has a property indicating whether
	// this is the final change to distinguish the difference between an intermediate flush to a file stream and the final close
	// of a file stream. The close query parameter is valid only when the action is
	// "flush" and change notifications are enabled. If the value of close is "true" and the flush operation completes successfully,
	// the service raises a file change notification with a property indicating
	// that this is the final update (the file stream has been closed). If "false" a change notification is raised indicating
	// the file has changed. The default is false. This query parameter is set to true
	// by the Hadoop ABFS driver to indicate that the file stream has been closed."
	Close *bool

	// Required for "Append Data" and "Flush Data". Must be 0 for "Flush Data". Must be the length of the request content in bytes
	// for "Append Data".
	ContentLength *int64

	// Optional. If "acquire" it will acquire the lease. If "auto-renew" it will renew the lease. If "release" it will release
	// the lease only on flush. If "acquire-release" it will acquire & complete the
	// operation & release the lease once operation is done.
	LeaseAction *LeaseAction

	// The lease duration is required to acquire a lease, and specifies the duration of the lease in seconds. The lease duration
	// must be between 15 and 60 seconds or -1 for infinite lease.
	LeaseDuration *int64

	// This parameter allows the caller to upload data in parallel and control the order in which it is appended to the file.
	// It is required when uploading data to be appended to the file and when flushing
	// previously uploaded data to the file. The value must be the position where the data is to be appended. Uploaded data is
	// not immediately flushed, or written, to the file. To flush, the previously
	// uploaded data must be contiguous, the position parameter must be specified and equal to the length of the file after all
	// data has been written, and there must not be a request entity body included
	// with the request.
	Position *int64

	// Proposed lease ID, in a GUID string format. The Blob service returns 400 (Invalid request) if the proposed lease ID is
	// not in the correct format. See Guid Constructor (String) for a list of valid GUID
	// string formats.
	ProposedLeaseID *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// Valid only for flush operations. If "true", uncommitted data is retained after the flush operation completes; otherwise,
	// the uncommitted data is deleted after the flush operation. The default is
	// false. Data at offsets less than the specified position are written to the file when flush succeeds, but this optional
	// parameter allows data after the flush position to be retained for a future flush
	// operation.
	RetainUncommittedData *bool

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// PathClientGetPropertiesOptions contains the optional parameters for the PathClient.GetProperties method.
type PathClientGetPropertiesOptions struct {
	// Optional. If the value is "getStatus" only the system defined properties for the path are returned. If the value is "getAccessControl"
	// the access control list is returned in the response headers
	// (Hierarchical Namespace must be enabled for the account), otherwise the properties are returned.
	Action *PathGetPropertiesAction

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32

	// Optional. Valid only when Hierarchical Namespace is enabled for the account. If "true", the user identity values returned
	// in the x-ms-owner, x-ms-group, and x-ms-acl response headers will be
	// transformed from Azure Active Directory Object IDs to User Principal Names. If "false", the values will be returned as
	// Azure Active Directory Object IDs. The default value is false. Note that group
	// and application Object IDs are not translated because they do not have unique friendly names.
	Upn *bool
}

// PathClientLeaseOptions contains the optional parameters for the PathClient.Lease method.
type PathClientLeaseOptions struct {
	// Proposed lease ID, in a GUID string format. The Blob service returns 400 (Invalid request) if the proposed lease ID is
	// not in the correct format. See Guid Constructor (String) for a list of valid GUID
	// string formats.
	ProposedLeaseID *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32

	// The lease break period duration is optional to break a lease, and specifies the break period of the lease in seconds. The
	// lease break duration must be between 0 and 60 seconds.
	XMSLeaseBreakPeriod *int32
}

// PathClientReadOptions contains the optional parameters for the PathClient.Read method.
type PathClientReadOptions struct {
	// The HTTP Range request header specifies one or more byte ranges of the resource to be retrieved.
	Range *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32

	// Optional. When this header is set to "true" and specified together with the Range header, the service returns the MD5 hash
	// for the range, as long as the range is less than or equal to 4MB in size. If
	// this header is specified without the Range header, the service returns status code 400 (Bad Request). If this header is
	// set to true when the range exceeds 4 MB in size, the service returns status code
	// 400 (Bad Request).
	XMSRangeGetContentMD5 *bool
}

// PathClientSetAccessControlOptions contains the optional parameters for the PathClient.SetAccessControl method.
type PathClientSetAccessControlOptions struct {
	// Sets POSIX access control rights on files and directories. The value is a comma-separated list of access control entries.
	// Each access control entry (ACE) consists of a scope, a type, a user or group
	// identifier, and permissions in the format "[scope:][type]:[id]:[permissions]".
	ACL *string

	// Optional. The owning group of the blob or directory.
	Group *string

	// Optional. The owner of the blob or directory.
	Owner *string

	// Optional and only valid if Hierarchical Namespace is enabled for the account. Sets POSIX access permissions for the file
	// owner, the file owning group, and others. Each class may be granted read,
	// write, or execute permission. The sticky bit is also supported. Both symbolic (rwxrw-rw-) and 4-digit octal notation (e.g.
	// 0766) are supported.
	Permissions *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// PathClientSetAccessControlRecursiveOptions contains the optional parameters for the PathClient.SetAccessControlRecursive
// method.
type PathClientSetAccessControlRecursiveOptions struct {
	// Sets POSIX access control rights on files and directories. The value is a comma-separated list of access control entries.
	// Each access control entry (ACE) consists of a scope, a type, a user or group
	// identifier, and permissions in the format "[scope:][type]:[id]:[permissions]".
	ACL *string

	// Optional. When deleting a directory, the number of paths that are deleted with each invocation is limited. If the number
	// of paths to be deleted exceeds this limit, a continuation token is returned in
	// this response header. When a continuation token is returned in the response, it must be specified in a subsequent invocation
	// of the delete operation to continue deleting the directory.
	Continuation *string

	// Optional. Valid for "SetAccessControlRecursive" operation. If set to false, the operation will terminate quickly on encountering
	// user errors (4XX). If true, the operation will ignore user errors and
	// proceed with the operation on other sub-entities of the directory. Continuation token will only be returned when forceFlag
	// is true in case of user errors. If not set the default value is false for
	// this.
	ForceFlag *bool

	// Optional. It specifies the maximum number of files or directories on which the acl change will be applied. If omitted or
	// greater than 2,000, the request will process up to 2,000 items
	MaxRecords *int32

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// PathClientSetExpiryOptions contains the optional parameters for the PathClient.SetExpiry method.
type PathClientSetExpiryOptions struct {
	// The time to set the blob to expiry
	ExpiresOn *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// PathClientUndeleteOptions contains the optional parameters for the PathClient.Undelete method.
type PathClientUndeleteOptions struct {
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32

	// Only for hierarchical namespace enabled accounts. Optional. The path of the soft deleted blob to undelete.
	UndeleteSource *string
}

// PathClientUpdateOptions contains the optional parameters for the PathClient.Update method.
type PathClientUpdateOptions struct {
	// Sets POSIX access control rights on files and directories. The value is a comma-separated list of access control entries.
	// Each access control entry (ACE) consists of a scope, a type, a user or group
	// identifier, and permissions in the format "[scope:][type]:[id]:[permissions]".
	ACL *string

	// Azure Storage Events allow applications to receive notifications when files change. When Azure Storage Events are enabled,
	// a file changed event is raised. This event has a property indicating whether
	// this is the final change to distinguish the difference between an intermediate flush to a file stream and the final close
	// of a file stream. The close query parameter is valid only when the action is
	// "flush" and change notifications are enabled. If the value of close is "true" and the flush operation completes successfully,
	// the service raises a file change notification with a property indicating
	// that this is the final update (the file stream has been closed). If "false" a change notification is raised indicating
	// the file has changed. The default is false. This query parameter is set to true
	// by the Hadoop ABFS driver to indicate that the file stream has been closed."
	Close *bool

	// Required for "Append Data" and "Flush Data". Must be 0 for "Flush Data". Must be the length of the request content in bytes
	// for "Append Data".
	ContentLength *int64

	// Optional. The number of paths processed with each invocation is limited. If the number of paths to be processed exceeds
	// this limit, a continuation token is returned in the response header
	// x-ms-continuation. When a continuation token is returned in the response, it must be percent-encoded and specified in a
	// subsequent invocation of setAccessControlRecursive operation.
	Continuation *string

	// Optional. Valid for "SetAccessControlRecursive" operation. If set to false, the operation will terminate quickly on encountering
	// user errors (4XX). If true, the operation will ignore user errors and
	// proceed with the operation on other sub-entities of the directory. Continuation token will only be returned when forceFlag
	// is true in case of user errors. If not set the default value is false for
	// this.
	ForceFlag *bool

	// Optional. The owning group of the blob or directory.
	Group *string

	// Optional. Valid for "SetAccessControlRecursive" operation. It specifies the maximum number of files or directories on which
	// the acl change will be applied. If omitted or greater than 2,000, the
	// request will process up to 2,000 items
	MaxRecords *int32

	// Optional. The owner of the blob or directory.
	Owner *string

	// Optional and only valid if Hierarchical Namespace is enabled for the account. Sets POSIX access permissions for the file
	// owner, the file owning group, and others. Each class may be granted read,
	// write, or execute permission. The sticky bit is also supported. Both symbolic (rwxrw-rw-) and 4-digit octal notation (e.g.
	// 0766) are supported.
	Permissions *string

	// This parameter allows the caller to upload data in parallel and control the order in which it is appended to the file.
	// It is required when uploading data to be appended to the file and when flushing
	// previously uploaded data to the file. The value must be the position where the data is to be appended. Uploaded data is
	// not immediately flushed, or written, to the file. To flush, the previously
	// uploaded data must be contiguous, the position parameter must be specified and equal to the length of the file after all
	// data has been written, and there must not be a request entity body included
	// with the request.
	Position *int64

	// Optional. User-defined properties to be stored with the filesystem, in the format of a comma-separated list of name and
	// value pairs "n1=v1, n2=v2, …", where each value is a base64 encoded string. Note
	// that the string may only contain ASCII characters in the ISO-8859-1 character set. If the filesystem exists, any properties
	// not included in the list will be removed. All properties are removed if the
	// header is omitted. To merge new and existing properties, first get all existing properties and the current E-Tag, then
	// make a conditional request with the E-Tag and include values for all properties.
	Properties *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// Valid only for flush operations. If "true", uncommitted data is retained after the flush operation completes; otherwise,
	// the uncommitted data is deleted after the flush operation. The default is
	// false. Data at offsets less than the specified position are written to the file when flush succeeds, but this optional
	// parameter allows data after the flush position to be retained for a future flush
	// operation.
	RetainUncommittedData *bool

	// Required if the request body is a structured message. Specifies the message schema version and properties.
	StructuredBodyType *string

	// Required if the request body is a structured message. Specifies the length of the blob/file content inside the message
	// body. Will always be smaller than Content-Length.
	StructuredContentLength *int64

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// PathHTTPHeaders contains a group of parameters for the PathClient.Create method.
type PathHTTPHeaders struct {
	// Optional. Sets the blob's cache control. If specified, this property is stored with the blob and returned with a read request.
	CacheControl *string

	// Optional. Sets the blob's Content-Disposition header.
	ContentDisposition *string

	// Optional. Sets the blob's content encoding. If specified, this property is stored with the blob and returned with a read
	// request.
	ContentEncoding *string

	// Optional. Set the blob's content language. If specified, this property is stored with the blob and returned with a read
	// request.
	ContentLanguage *string

	// Specify the transactional md5 for the body, to be validated by the service.
	ContentMD5 []byte

	// Optional. Sets the blob's content type. If specified, this property is stored with the blob and returned with a read request.
	ContentType *string

	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentHash []byte
}

// ServiceClientListFileSystemsOptions contains the optional parameters for the ServiceClient.NewListFileSystemsPager method.
type ServiceClientListFileSystemsOptions struct {
	// Optional. When deleting a directory, the number of paths that are deleted with each invocation is limited. If the number
	// of paths to be deleted exceeds this limit, a continuation token is returned in
	// this response header. When a continuation token is returned in the response, it must be specified in a subsequent invocation
	// of the delete operation to continue deleting the directory.
	Continuation *string

	// An optional value that specifies the maximum number of items to return. If omitted or greater than 5,000, the response
	// will include up to 5,000 items.
	MaxResults *int32

	// Filters results to filesystems within the specified prefix.
	Prefix *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// SourceModifiedAccessConditions contains a group of parameters for the PathClient.Create method.
type SourceModifiedAccessConditions struct {
	// Specify an ETag value to operate only on blobs with a matching value.
	SourceIfMatch *azcore.ETag

	// Specify this header value to operate only on a blob if it has been modified since the specified date/time.
	SourceIfModifiedSince *time.Time

	// Specify an ETag value to operate only on blobs without a matching value.
	SourceIfNoneMatch *azcore.ETag

	// Specify this header value to operate only on a blob if it has not been modified since the specified date/time.
	SourceIfUnmodifiedSince *time.Time
}

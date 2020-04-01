// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azblob

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"strconv"
	"time"
)

// DirectoryOperations contains the methods for the Directory group.
type DirectoryOperations interface {
	// Create - Create a directory. By default, the destination is overwritten and if the destination already exists and has a lease the lease is broken. This operation supports conditional HTTP requests.  For more information, see [Specifying Conditional Headers for Blob Service Operations](https://docs.microsoft.com/en-us/rest/api/storageservices/specifying-conditional-headers-for-blob-service-operations).  To fail if the destination already exists, use a conditional request with If-None-Match: "*".
	Create(ctx context.Context, options *DirectoryCreateOptions) (*DirectoryCreateResponse, error)
	// Delete - Deletes the directory
	Delete(ctx context.Context, recursiveDirectoryDelete bool, options *DirectoryDeleteOptions) (*DirectoryDeleteResponse, error)
	// GetAccessControl - Get the owner, group, permissions, or access control list for a directory.
	GetAccessControl(ctx context.Context, options *DirectoryGetAccessControlOptions) (*DirectoryGetAccessControlResponse, error)
	// Rename - Rename a directory. By default, the destination is overwritten and if the destination already exists and has a lease the lease is broken. This operation supports conditional HTTP requests. For more information, see [Specifying Conditional Headers for Blob Service Operations](https://docs.microsoft.com/en-us/rest/api/storageservices/specifying-conditional-headers-for-blob-service-operations). To fail if the destination already exists, use a conditional request with If-None-Match: "*".
	Rename(ctx context.Context, renameSource string, options *DirectoryRenameOptions) (*DirectoryRenameResponse, error)
	// SetAccessControl - Set the owner, group, permissions, or access control list for a directory.
	SetAccessControl(ctx context.Context, options *DirectorySetAccessControlOptions) (*DirectorySetAccessControlResponse, error)
}

// directoryOperations implements the DirectoryOperations interface.
type directoryOperations struct {
	*Client
	pathRenameMode *PathRenameMode
}

// Create - Create a directory. By default, the destination is overwritten and if the destination already exists and has a lease the lease is broken. This operation supports conditional HTTP requests.  For more information, see [Specifying Conditional Headers for Blob Service Operations](https://docs.microsoft.com/en-us/rest/api/storageservices/specifying-conditional-headers-for-blob-service-operations).  To fail if the destination already exists, use a conditional request with If-None-Match: "*".
func (client *directoryOperations) Create(ctx context.Context, options *DirectoryCreateOptions) (*DirectoryCreateResponse, error) {
	req, err := client.createCreateRequest(options)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.createHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// createCreateRequest creates the Create request.
func (client *directoryOperations) createCreateRequest(options *DirectoryCreateOptions) (*azcore.Request, error) {
	u := client.u
	query := u.Query()
	query.Set("resource", "directory")
	if options != nil && options.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodPut, *u)
	if options != nil && options.DirectoryProperties != nil {
		req.Header.Set("x-ms-properties", *options.DirectoryProperties)
	}
	if options != nil && options.PosixPermissions != nil {
		req.Header.Set("x-ms-permissions", *options.PosixPermissions)
	}
	if options != nil && options.PosixUmask != nil {
		req.Header.Set("x-ms-umask", *options.PosixUmask)
	}
	if options != nil && options.CacheControl != nil {
		req.Header.Set("x-ms-cache-control", *options.CacheControl)
	}
	if options != nil && options.ContentType != nil {
		req.Header.Set("x-ms-content-type", *options.ContentType)
	}
	if options != nil && options.ContentEncoding != nil {
		req.Header.Set("x-ms-content-encoding", *options.ContentEncoding)
	}
	if options != nil && options.ContentLanguage != nil {
		req.Header.Set("x-ms-content-language", *options.ContentLanguage)
	}
	if options != nil && options.ContentDisposition != nil {
		req.Header.Set("x-ms-content-disposition", *options.ContentDisposition)
	}
	if options != nil && options.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *options.LeaseId)
	}
	if options != nil && options.IfModifiedSince != nil {
		req.Header.Set("If-Modified-Since", options.IfModifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.IfUnmodifiedSince != nil {
		req.Header.Set("If-Unmodified-Since", options.IfUnmodifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.IfMatch != nil {
		req.Header.Set("If-Match", *options.IfMatch)
	}
	if options != nil && options.IfNoneMatch != nil {
		req.Header.Set("If-None-Match", *options.IfNoneMatch)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	if options != nil && options.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *options.RequestId)
	}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *directoryOperations) createHandleResponse(resp *azcore.Response) (*DirectoryCreateResponse, error) {
	if !resp.HasStatusCode(http.StatusCreated) {
		return nil, newDataLakeStorageError(resp)
	}
	result := DirectoryCreateResponse{RawResponse: resp.Response}
	eTag := resp.Header.Get("ETag")
	result.ETag = &eTag
	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	result.LastModified = &lastModified
	clientRequestId := resp.Header.Get("x-ms-client-request-id")
	result.ClientRequestId = &clientRequestId
	requestId := resp.Header.Get("x-ms-request-id")
	result.RequestId = &requestId
	version := resp.Header.Get("x-ms-version")
	result.Version = &version
	contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, err
	}
	result.ContentLength = &contentLength
	date, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return nil, err
	}
	result.Date = &date
	return &result, nil
}

// Delete - Deletes the directory
func (client *directoryOperations) Delete(ctx context.Context, recursiveDirectoryDelete bool, options *DirectoryDeleteOptions) (*DirectoryDeleteResponse, error) {
	req, err := client.deleteCreateRequest(recursiveDirectoryDelete, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.deleteHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// deleteCreateRequest creates the Delete request.
func (client *directoryOperations) deleteCreateRequest(recursiveDirectoryDelete bool, options *DirectoryDeleteOptions) (*azcore.Request, error) {
	u := client.u
	query := u.Query()
	if options != nil && options.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	query.Set("recursive", strconv.FormatBool(recursiveDirectoryDelete))
	if options != nil && options.Marker != nil {
		query.Set("continuation", *options.Marker)
	}
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodDelete, *u)
	if options != nil && options.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *options.LeaseId)
	}
	if options != nil && options.IfModifiedSince != nil {
		req.Header.Set("If-Modified-Since", options.IfModifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.IfUnmodifiedSince != nil {
		req.Header.Set("If-Unmodified-Since", options.IfUnmodifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.IfMatch != nil {
		req.Header.Set("If-Match", *options.IfMatch)
	}
	if options != nil && options.IfNoneMatch != nil {
		req.Header.Set("If-None-Match", *options.IfNoneMatch)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	if options != nil && options.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *options.RequestId)
	}
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *directoryOperations) deleteHandleResponse(resp *azcore.Response) (*DirectoryDeleteResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, newDataLakeStorageError(resp)
	}
	result := DirectoryDeleteResponse{RawResponse: resp.Response}
	continuation := resp.Header.Get("x-ms-continuation")
	result.Continuation = &continuation
	clientRequestId := resp.Header.Get("x-ms-client-request-id")
	result.ClientRequestId = &clientRequestId
	requestId := resp.Header.Get("x-ms-request-id")
	result.RequestId = &requestId
	version := resp.Header.Get("x-ms-version")
	result.Version = &version
	date, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return nil, err
	}
	result.Date = &date
	return &result, nil
}

// GetAccessControl - Get the owner, group, permissions, or access control list for a directory.
func (client *directoryOperations) GetAccessControl(ctx context.Context, options *DirectoryGetAccessControlOptions) (*DirectoryGetAccessControlResponse, error) {
	req, err := client.getAccessControlCreateRequest(options)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.getAccessControlHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getAccessControlCreateRequest creates the GetAccessControl request.
func (client *directoryOperations) getAccessControlCreateRequest(options *DirectoryGetAccessControlOptions) (*azcore.Request, error) {
	u := client.u
	query := u.Query()
	query.Set("action", "getAccessControl")
	if options != nil && options.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	if options != nil && options.Upn != nil {
		query.Set("upn", strconv.FormatBool(*options.Upn))
	}
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodHead, *u)
	if options != nil && options.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *options.LeaseId)
	}
	if options != nil && options.IfMatch != nil {
		req.Header.Set("If-Match", *options.IfMatch)
	}
	if options != nil && options.IfNoneMatch != nil {
		req.Header.Set("If-None-Match", *options.IfNoneMatch)
	}
	if options != nil && options.IfModifiedSince != nil {
		req.Header.Set("If-Modified-Since", options.IfModifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.IfUnmodifiedSince != nil {
		req.Header.Set("If-Unmodified-Since", options.IfUnmodifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *options.RequestId)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	return req, nil
}

// getAccessControlHandleResponse handles the GetAccessControl response.
func (client *directoryOperations) getAccessControlHandleResponse(resp *azcore.Response) (*DirectoryGetAccessControlResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, newDataLakeStorageError(resp)
	}
	result := DirectoryGetAccessControlResponse{RawResponse: resp.Response}
	date, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return nil, err
	}
	result.Date = &date
	eTag := resp.Header.Get("ETag")
	result.ETag = &eTag
	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	result.LastModified = &lastModified
	owner := resp.Header.Get("x-ms-owner")
	result.Owner = &owner
	group := resp.Header.Get("x-ms-group")
	result.Group = &group
	permissions := resp.Header.Get("x-ms-permissions")
	result.Permissions = &permissions
	acl := resp.Header.Get("x-ms-acl")
	result.Acl = &acl
	requestId := resp.Header.Get("x-ms-request-id")
	result.RequestId = &requestId
	version := resp.Header.Get("x-ms-version")
	result.Version = &version
	return &result, nil
}

// Rename - Rename a directory. By default, the destination is overwritten and if the destination already exists and has a lease the lease is broken. This operation supports conditional HTTP requests. For more information, see [Specifying Conditional Headers for Blob Service Operations](https://docs.microsoft.com/en-us/rest/api/storageservices/specifying-conditional-headers-for-blob-service-operations). To fail if the destination already exists, use a conditional request with If-None-Match: "*".
func (client *directoryOperations) Rename(ctx context.Context, renameSource string, options *DirectoryRenameOptions) (*DirectoryRenameResponse, error) {
	req, err := client.renameCreateRequest(renameSource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.renameHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// renameCreateRequest creates the Rename request.
func (client *directoryOperations) renameCreateRequest(renameSource string, options *DirectoryRenameOptions) (*azcore.Request, error) {
	u := client.u
	query := u.Query()
	if options != nil && options.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	if options != nil && options.Marker != nil {
		query.Set("continuation", *options.Marker)
	}
	if client.pathRenameMode != nil {
		query.Set("mode", string(*client.pathRenameMode))
	}
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodPut, *u)
	req.Header.Set("x-ms-rename-source", renameSource)
	if options != nil && options.DirectoryProperties != nil {
		req.Header.Set("x-ms-properties", *options.DirectoryProperties)
	}
	if options != nil && options.PosixPermissions != nil {
		req.Header.Set("x-ms-permissions", *options.PosixPermissions)
	}
	if options != nil && options.PosixUmask != nil {
		req.Header.Set("x-ms-umask", *options.PosixUmask)
	}
	if options != nil && options.CacheControl != nil {
		req.Header.Set("x-ms-cache-control", *options.CacheControl)
	}
	if options != nil && options.ContentType != nil {
		req.Header.Set("x-ms-content-type", *options.ContentType)
	}
	if options != nil && options.ContentEncoding != nil {
		req.Header.Set("x-ms-content-encoding", *options.ContentEncoding)
	}
	if options != nil && options.ContentLanguage != nil {
		req.Header.Set("x-ms-content-language", *options.ContentLanguage)
	}
	if options != nil && options.ContentDisposition != nil {
		req.Header.Set("x-ms-content-disposition", *options.ContentDisposition)
	}
	if options != nil && options.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *options.LeaseId)
	}
	if options != nil && options.SourceLeaseId != nil {
		req.Header.Set("x-ms-source-lease-id", *options.SourceLeaseId)
	}
	if options != nil && options.IfModifiedSince != nil {
		req.Header.Set("If-Modified-Since", options.IfModifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.IfUnmodifiedSince != nil {
		req.Header.Set("If-Unmodified-Since", options.IfUnmodifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.IfMatch != nil {
		req.Header.Set("If-Match", *options.IfMatch)
	}
	if options != nil && options.IfNoneMatch != nil {
		req.Header.Set("If-None-Match", *options.IfNoneMatch)
	}
	if options != nil && options.SourceIfModifiedSince != nil {
		req.Header.Set("x-ms-source-if-modified-since", options.SourceIfModifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.SourceIfUnmodifiedSince != nil {
		req.Header.Set("x-ms-source-if-unmodified-since", options.SourceIfUnmodifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.SourceIfMatch != nil {
		req.Header.Set("x-ms-source-if-match", *options.SourceIfMatch)
	}
	if options != nil && options.SourceIfNoneMatch != nil {
		req.Header.Set("x-ms-source-if-none-match", *options.SourceIfNoneMatch)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	if options != nil && options.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *options.RequestId)
	}
	return req, nil
}

// renameHandleResponse handles the Rename response.
func (client *directoryOperations) renameHandleResponse(resp *azcore.Response) (*DirectoryRenameResponse, error) {
	if !resp.HasStatusCode(http.StatusCreated) {
		return nil, newDataLakeStorageError(resp)
	}
	result := DirectoryRenameResponse{RawResponse: resp.Response}
	continuation := resp.Header.Get("x-ms-continuation")
	result.Continuation = &continuation
	eTag := resp.Header.Get("ETag")
	result.ETag = &eTag
	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	result.LastModified = &lastModified
	clientRequestId := resp.Header.Get("x-ms-client-request-id")
	result.ClientRequestId = &clientRequestId
	requestId := resp.Header.Get("x-ms-request-id")
	result.RequestId = &requestId
	version := resp.Header.Get("x-ms-version")
	result.Version = &version
	contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, err
	}
	result.ContentLength = &contentLength
	date, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return nil, err
	}
	result.Date = &date
	return &result, nil
}

// SetAccessControl - Set the owner, group, permissions, or access control list for a directory.
func (client *directoryOperations) SetAccessControl(ctx context.Context, options *DirectorySetAccessControlOptions) (*DirectorySetAccessControlResponse, error) {
	req, err := client.setAccessControlCreateRequest(options)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.setAccessControlHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// setAccessControlCreateRequest creates the SetAccessControl request.
func (client *directoryOperations) setAccessControlCreateRequest(options *DirectorySetAccessControlOptions) (*azcore.Request, error) {
	u := client.u
	query := u.Query()
	query.Set("action", "setAccessControl")
	if options != nil && options.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodPatch, *u)
	if options != nil && options.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *options.LeaseId)
	}
	if options != nil && options.Owner != nil {
		req.Header.Set("x-ms-owner", *options.Owner)
	}
	if options != nil && options.Group != nil {
		req.Header.Set("x-ms-group", *options.Group)
	}
	if options != nil && options.PosixPermissions != nil {
		req.Header.Set("x-ms-permissions", *options.PosixPermissions)
	}
	if options != nil && options.PosixAcl != nil {
		req.Header.Set("x-ms-acl", *options.PosixAcl)
	}
	if options != nil && options.IfMatch != nil {
		req.Header.Set("If-Match", *options.IfMatch)
	}
	if options != nil && options.IfNoneMatch != nil {
		req.Header.Set("If-None-Match", *options.IfNoneMatch)
	}
	if options != nil && options.IfModifiedSince != nil {
		req.Header.Set("If-Modified-Since", options.IfModifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.IfUnmodifiedSince != nil {
		req.Header.Set("If-Unmodified-Since", options.IfUnmodifiedSince.Format(time.RFC1123))
	}
	if options != nil && options.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *options.RequestId)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	return req, nil
}

// setAccessControlHandleResponse handles the SetAccessControl response.
func (client *directoryOperations) setAccessControlHandleResponse(resp *azcore.Response) (*DirectorySetAccessControlResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, newDataLakeStorageError(resp)
	}
	result := DirectorySetAccessControlResponse{RawResponse: resp.Response}
	date, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return nil, err
	}
	result.Date = &date
	eTag := resp.Header.Get("ETag")
	result.ETag = &eTag
	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	result.LastModified = &lastModified
	requestId := resp.Header.Get("x-ms-request-id")
	result.RequestId = &requestId
	version := resp.Header.Get("x-ms-version")
	result.Version = &version
	return &result, nil
}

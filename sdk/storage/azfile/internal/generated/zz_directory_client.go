//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package generated

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DirectoryClient contains the methods for the Directory group.
// Don't use this type directly, use a constructor function instead.
type DirectoryClient struct {
	internal               *azcore.Client
	endpoint               string
	allowTrailingDot       *bool
	fileRequestIntent      *ShareTokenIntent
	allowSourceTrailingDot *bool
}

// Create - Creates a new directory under the specified share or parent directory.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-05
//   - options - DirectoryClientCreateOptions contains the optional parameters for the DirectoryClient.Create method.
func (client *DirectoryClient) Create(ctx context.Context, options *DirectoryClientCreateOptions) (DirectoryClientCreateResponse, error) {
	var err error
	req, err := client.createCreateRequest(ctx, options)
	if err != nil {
		return DirectoryClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DirectoryClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return DirectoryClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *DirectoryClient) createCreateRequest(ctx context.Context, options *DirectoryClientCreateOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPut, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "directory")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	if options != nil && options.Metadata != nil {
		for k, v := range options.Metadata {
			if v != nil {
				req.Raw().Header["x-ms-meta-"+k] = []string{*v}
			}
		}
	}
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	if options != nil && options.FilePermission != nil {
		req.Raw().Header["x-ms-file-permission"] = []string{*options.FilePermission}
	}
	if options != nil && options.FilePermissionFormat != nil {
		req.Raw().Header["x-ms-file-permission-format"] = []string{string(*options.FilePermissionFormat)}
	}
	if options != nil && options.FilePermissionKey != nil {
		req.Raw().Header["x-ms-file-permission-key"] = []string{*options.FilePermissionKey}
	}
	if options != nil && options.FileAttributes != nil {
		req.Raw().Header["x-ms-file-attributes"] = []string{*options.FileAttributes}
	}
	if options != nil && options.FileCreationTime != nil {
		req.Raw().Header["x-ms-file-creation-time"] = []string{*options.FileCreationTime}
	}
	if options != nil && options.FileLastWriteTime != nil {
		req.Raw().Header["x-ms-file-last-write-time"] = []string{*options.FileLastWriteTime}
	}
	if options != nil && options.FileChangeTime != nil {
		req.Raw().Header["x-ms-file-change-time"] = []string{*options.FileChangeTime}
	}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *DirectoryClient) createHandleResponse(resp *http.Response) (DirectoryClientCreateResponse, error) {
	result := DirectoryClientCreateResponse{}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientCreateResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = (*azcore.ETag)(&val)
	}
	if val := resp.Header.Get("x-ms-file-attributes"); val != "" {
		result.FileAttributes = &val
	}
	if val := resp.Header.Get("x-ms-file-change-time"); val != "" {
		fileChangeTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientCreateResponse{}, err
		}
		result.FileChangeTime = &fileChangeTime
	}
	if val := resp.Header.Get("x-ms-file-creation-time"); val != "" {
		fileCreationTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientCreateResponse{}, err
		}
		result.FileCreationTime = &fileCreationTime
	}
	if val := resp.Header.Get("x-ms-file-id"); val != "" {
		result.ID = &val
	}
	if val := resp.Header.Get("x-ms-file-last-write-time"); val != "" {
		fileLastWriteTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientCreateResponse{}, err
		}
		result.FileLastWriteTime = &fileLastWriteTime
	}
	if val := resp.Header.Get("x-ms-file-parent-id"); val != "" {
		result.ParentID = &val
	}
	if val := resp.Header.Get("x-ms-file-permission-key"); val != "" {
		result.FilePermissionKey = &val
	}
	if val := resp.Header.Get("x-ms-request-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return DirectoryClientCreateResponse{}, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("Last-Modified"); val != "" {
		lastModified, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientCreateResponse{}, err
		}
		result.LastModified = &lastModified
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// Delete - Removes the specified empty directory. Note that the directory must be empty before it can be deleted.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-05
//   - options - DirectoryClientDeleteOptions contains the optional parameters for the DirectoryClient.Delete method.
func (client *DirectoryClient) Delete(ctx context.Context, options *DirectoryClientDeleteOptions) (DirectoryClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, options)
	if err != nil {
		return DirectoryClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DirectoryClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return DirectoryClientDeleteResponse{}, err
	}
	resp, err := client.deleteHandleResponse(httpResp)
	return resp, err
}

// deleteCreateRequest creates the Delete request.
func (client *DirectoryClient) deleteCreateRequest(ctx context.Context, options *DirectoryClientDeleteOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodDelete, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "directory")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *DirectoryClient) deleteHandleResponse(resp *http.Response) (DirectoryClientDeleteResponse, error) {
	result := DirectoryClientDeleteResponse{}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientDeleteResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// ForceCloseHandles - Closes all handles open for given directory.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-05
//   - handleID - Specifies handle ID opened on the file or directory to be closed. Asterisk (‘*’) is a wildcard that specifies
//     all handles.
//   - options - DirectoryClientForceCloseHandlesOptions contains the optional parameters for the DirectoryClient.ForceCloseHandles
//     method.
func (client *DirectoryClient) ForceCloseHandles(ctx context.Context, handleID string, options *DirectoryClientForceCloseHandlesOptions) (DirectoryClientForceCloseHandlesResponse, error) {
	var err error
	req, err := client.forceCloseHandlesCreateRequest(ctx, handleID, options)
	if err != nil {
		return DirectoryClientForceCloseHandlesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DirectoryClientForceCloseHandlesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DirectoryClientForceCloseHandlesResponse{}, err
	}
	resp, err := client.forceCloseHandlesHandleResponse(httpResp)
	return resp, err
}

// forceCloseHandlesCreateRequest creates the ForceCloseHandles request.
func (client *DirectoryClient) forceCloseHandlesCreateRequest(ctx context.Context, handleID string, options *DirectoryClientForceCloseHandlesOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPut, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("comp", "forceclosehandles")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	if options != nil && options.Marker != nil {
		reqQP.Set("marker", *options.Marker)
	}
	if options != nil && options.Sharesnapshot != nil {
		reqQP.Set("sharesnapshot", *options.Sharesnapshot)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-handle-id"] = []string{handleID}
	if options != nil && options.Recursive != nil {
		req.Raw().Header["x-ms-recursive"] = []string{strconv.FormatBool(*options.Recursive)}
	}
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// forceCloseHandlesHandleResponse handles the ForceCloseHandles response.
func (client *DirectoryClient) forceCloseHandlesHandleResponse(resp *http.Response) (DirectoryClientForceCloseHandlesResponse, error) {
	result := DirectoryClientForceCloseHandlesResponse{}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientForceCloseHandlesResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-marker"); val != "" {
		result.Marker = &val
	}
	if val := resp.Header.Get("x-ms-number-of-handles-closed"); val != "" {
		numberOfHandlesClosed32, err := strconv.ParseInt(val, 10, 32)
		numberOfHandlesClosed := int32(numberOfHandlesClosed32)
		if err != nil {
			return DirectoryClientForceCloseHandlesResponse{}, err
		}
		result.NumberOfHandlesClosed = &numberOfHandlesClosed
	}
	if val := resp.Header.Get("x-ms-number-of-handles-failed"); val != "" {
		numberOfHandlesFailedToClose32, err := strconv.ParseInt(val, 10, 32)
		numberOfHandlesFailedToClose := int32(numberOfHandlesFailedToClose32)
		if err != nil {
			return DirectoryClientForceCloseHandlesResponse{}, err
		}
		result.NumberOfHandlesFailedToClose = &numberOfHandlesFailedToClose
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// GetProperties - Returns all system properties for the specified directory, and can also be used to check the existence
// of a directory. The data returned does not include the files in the directory or any
// subdirectories.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-05
//   - options - DirectoryClientGetPropertiesOptions contains the optional parameters for the DirectoryClient.GetProperties method.
func (client *DirectoryClient) GetProperties(ctx context.Context, options *DirectoryClientGetPropertiesOptions) (DirectoryClientGetPropertiesResponse, error) {
	var err error
	req, err := client.getPropertiesCreateRequest(ctx, options)
	if err != nil {
		return DirectoryClientGetPropertiesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DirectoryClientGetPropertiesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DirectoryClientGetPropertiesResponse{}, err
	}
	resp, err := client.getPropertiesHandleResponse(httpResp)
	return resp, err
}

// getPropertiesCreateRequest creates the GetProperties request.
func (client *DirectoryClient) getPropertiesCreateRequest(ctx context.Context, options *DirectoryClientGetPropertiesOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "directory")
	if options != nil && options.Sharesnapshot != nil {
		reqQP.Set("sharesnapshot", *options.Sharesnapshot)
	}
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// getPropertiesHandleResponse handles the GetProperties response.
func (client *DirectoryClient) getPropertiesHandleResponse(resp *http.Response) (DirectoryClientGetPropertiesResponse, error) {
	result := DirectoryClientGetPropertiesResponse{}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientGetPropertiesResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = (*azcore.ETag)(&val)
	}
	if val := resp.Header.Get("x-ms-file-attributes"); val != "" {
		result.FileAttributes = &val
	}
	if val := resp.Header.Get("x-ms-file-change-time"); val != "" {
		fileChangeTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientGetPropertiesResponse{}, err
		}
		result.FileChangeTime = &fileChangeTime
	}
	if val := resp.Header.Get("x-ms-file-creation-time"); val != "" {
		fileCreationTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientGetPropertiesResponse{}, err
		}
		result.FileCreationTime = &fileCreationTime
	}
	if val := resp.Header.Get("x-ms-file-id"); val != "" {
		result.ID = &val
	}
	if val := resp.Header.Get("x-ms-file-last-write-time"); val != "" {
		fileLastWriteTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientGetPropertiesResponse{}, err
		}
		result.FileLastWriteTime = &fileLastWriteTime
	}
	if val := resp.Header.Get("x-ms-file-parent-id"); val != "" {
		result.ParentID = &val
	}
	if val := resp.Header.Get("x-ms-file-permission-key"); val != "" {
		result.FilePermissionKey = &val
	}
	if val := resp.Header.Get("x-ms-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return DirectoryClientGetPropertiesResponse{}, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("Last-Modified"); val != "" {
		lastModified, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientGetPropertiesResponse{}, err
		}
		result.LastModified = &lastModified
	}
	for hh := range resp.Header {
		if len(hh) > len("x-ms-meta-") && strings.EqualFold(hh[:len("x-ms-meta-")], "x-ms-meta-") {
			if result.Metadata == nil {
				result.Metadata = map[string]*string{}
			}
			result.Metadata[hh[len("x-ms-meta-"):]] = to.Ptr(resp.Header.Get(hh))
		}
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// NewListFilesAndDirectoriesSegmentPager - Returns a list of files or directories under the specified share or directory.
// It lists the contents only for a single level of the directory hierarchy.
//
// Generated from API version 2025-01-05
//   - options - DirectoryClientListFilesAndDirectoriesSegmentOptions contains the optional parameters for the DirectoryClient.NewListFilesAndDirectoriesSegmentPager
//     method.
//
// listFilesAndDirectoriesSegmentCreateRequest creates the ListFilesAndDirectoriesSegment request.
func (client *DirectoryClient) ListFilesAndDirectoriesSegmentCreateRequest(ctx context.Context, options *DirectoryClientListFilesAndDirectoriesSegmentOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "directory")
	reqQP.Set("comp", "list")
	if options != nil && options.Prefix != nil {
		reqQP.Set("prefix", *options.Prefix)
	}
	if options != nil && options.Sharesnapshot != nil {
		reqQP.Set("sharesnapshot", *options.Sharesnapshot)
	}
	if options != nil && options.Marker != nil {
		reqQP.Set("marker", *options.Marker)
	}
	if options != nil && options.Maxresults != nil {
		reqQP.Set("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
	}
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	if options != nil && options.Include != nil {
		reqQP.Set("include", strings.Join(strings.Fields(strings.Trim(fmt.Sprint(options.Include), "[]")), ","))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	if options != nil && options.IncludeExtendedInfo != nil {
		req.Raw().Header["x-ms-file-extended-info"] = []string{strconv.FormatBool(*options.IncludeExtendedInfo)}
	}
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// listFilesAndDirectoriesSegmentHandleResponse handles the ListFilesAndDirectoriesSegment response.
func (client *DirectoryClient) ListFilesAndDirectoriesSegmentHandleResponse(resp *http.Response) (DirectoryClientListFilesAndDirectoriesSegmentResponse, error) {
	result := DirectoryClientListFilesAndDirectoriesSegmentResponse{}
	if val := resp.Header.Get("Content-Type"); val != "" {
		result.ContentType = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientListFilesAndDirectoriesSegmentResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if err := runtime.UnmarshalAsXML(resp, &result.ListFilesAndDirectoriesSegmentResponse); err != nil {
		return DirectoryClientListFilesAndDirectoriesSegmentResponse{}, err
	}
	return result, nil
}

// ListHandles - Lists handles for directory.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-05
//   - options - DirectoryClientListHandlesOptions contains the optional parameters for the DirectoryClient.ListHandles method.
func (client *DirectoryClient) ListHandles(ctx context.Context, options *DirectoryClientListHandlesOptions) (DirectoryClientListHandlesResponse, error) {
	var err error
	req, err := client.listHandlesCreateRequest(ctx, options)
	if err != nil {
		return DirectoryClientListHandlesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DirectoryClientListHandlesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DirectoryClientListHandlesResponse{}, err
	}
	resp, err := client.listHandlesHandleResponse(httpResp)
	return resp, err
}

// listHandlesCreateRequest creates the ListHandles request.
func (client *DirectoryClient) listHandlesCreateRequest(ctx context.Context, options *DirectoryClientListHandlesOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("comp", "listhandles")
	if options != nil && options.Marker != nil {
		reqQP.Set("marker", *options.Marker)
	}
	if options != nil && options.Maxresults != nil {
		reqQP.Set("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
	}
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	if options != nil && options.Sharesnapshot != nil {
		reqQP.Set("sharesnapshot", *options.Sharesnapshot)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.Recursive != nil {
		req.Raw().Header["x-ms-recursive"] = []string{strconv.FormatBool(*options.Recursive)}
	}
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// listHandlesHandleResponse handles the ListHandles response.
func (client *DirectoryClient) listHandlesHandleResponse(resp *http.Response) (DirectoryClientListHandlesResponse, error) {
	result := DirectoryClientListHandlesResponse{}
	if val := resp.Header.Get("Content-Type"); val != "" {
		result.ContentType = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientListHandlesResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if err := runtime.UnmarshalAsXML(resp, &result.ListHandlesResponse); err != nil {
		return DirectoryClientListHandlesResponse{}, err
	}
	return result, nil
}

// Rename - Renames a directory
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-05
//   - renameSource - Required. Specifies the URI-style path of the source file, up to 2 KB in length.
//   - options - DirectoryClientRenameOptions contains the optional parameters for the DirectoryClient.Rename method.
//   - SourceLeaseAccessConditions - SourceLeaseAccessConditions contains a group of parameters for the DirectoryClient.Rename
//     method.
//   - DestinationLeaseAccessConditions - DestinationLeaseAccessConditions contains a group of parameters for the DirectoryClient.Rename
//     method.
//   - CopyFileSMBInfo - CopyFileSMBInfo contains a group of parameters for the DirectoryClient.Rename method.
func (client *DirectoryClient) Rename(ctx context.Context, renameSource string, options *DirectoryClientRenameOptions, sourceLeaseAccessConditions *SourceLeaseAccessConditions, destinationLeaseAccessConditions *DestinationLeaseAccessConditions, copyFileSMBInfo *CopyFileSMBInfo) (DirectoryClientRenameResponse, error) {
	var err error
	req, err := client.renameCreateRequest(ctx, renameSource, options, sourceLeaseAccessConditions, destinationLeaseAccessConditions, copyFileSMBInfo)
	if err != nil {
		return DirectoryClientRenameResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DirectoryClientRenameResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DirectoryClientRenameResponse{}, err
	}
	resp, err := client.renameHandleResponse(httpResp)
	return resp, err
}

// renameCreateRequest creates the Rename request.
func (client *DirectoryClient) renameCreateRequest(ctx context.Context, renameSource string, options *DirectoryClientRenameOptions, sourceLeaseAccessConditions *SourceLeaseAccessConditions, destinationLeaseAccessConditions *DestinationLeaseAccessConditions, copyFileSMBInfo *CopyFileSMBInfo) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPut, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "directory")
	reqQP.Set("comp", "rename")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	req.Raw().Header["x-ms-file-rename-source"] = []string{renameSource}
	if options != nil && options.ReplaceIfExists != nil {
		req.Raw().Header["x-ms-file-rename-replace-if-exists"] = []string{strconv.FormatBool(*options.ReplaceIfExists)}
	}
	if options != nil && options.IgnoreReadOnly != nil {
		req.Raw().Header["x-ms-file-rename-ignore-readonly"] = []string{strconv.FormatBool(*options.IgnoreReadOnly)}
	}
	if sourceLeaseAccessConditions != nil && sourceLeaseAccessConditions.SourceLeaseID != nil {
		req.Raw().Header["x-ms-source-lease-id"] = []string{*sourceLeaseAccessConditions.SourceLeaseID}
	}
	if destinationLeaseAccessConditions != nil && destinationLeaseAccessConditions.DestinationLeaseID != nil {
		req.Raw().Header["x-ms-destination-lease-id"] = []string{*destinationLeaseAccessConditions.DestinationLeaseID}
	}
	if copyFileSMBInfo != nil && copyFileSMBInfo.FileAttributes != nil {
		req.Raw().Header["x-ms-file-attributes"] = []string{*copyFileSMBInfo.FileAttributes}
	}
	if copyFileSMBInfo != nil && copyFileSMBInfo.FileCreationTime != nil {
		req.Raw().Header["x-ms-file-creation-time"] = []string{*copyFileSMBInfo.FileCreationTime}
	}
	if copyFileSMBInfo != nil && copyFileSMBInfo.FileLastWriteTime != nil {
		req.Raw().Header["x-ms-file-last-write-time"] = []string{*copyFileSMBInfo.FileLastWriteTime}
	}
	if copyFileSMBInfo != nil && copyFileSMBInfo.FileChangeTime != nil {
		req.Raw().Header["x-ms-file-change-time"] = []string{*copyFileSMBInfo.FileChangeTime}
	}
	if options != nil && options.FilePermission != nil {
		req.Raw().Header["x-ms-file-permission"] = []string{*options.FilePermission}
	}
	if options != nil && options.FilePermissionFormat != nil {
		req.Raw().Header["x-ms-file-permission-format"] = []string{string(*options.FilePermissionFormat)}
	}
	if options != nil && options.FilePermissionKey != nil {
		req.Raw().Header["x-ms-file-permission-key"] = []string{*options.FilePermissionKey}
	}
	if options != nil && options.Metadata != nil {
		for k, v := range options.Metadata {
			if v != nil {
				req.Raw().Header["x-ms-meta-"+k] = []string{*v}
			}
		}
	}
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	if client.allowSourceTrailingDot != nil {
		req.Raw().Header["x-ms-source-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowSourceTrailingDot)}
	}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// renameHandleResponse handles the Rename response.
func (client *DirectoryClient) renameHandleResponse(resp *http.Response) (DirectoryClientRenameResponse, error) {
	result := DirectoryClientRenameResponse{}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientRenameResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = (*azcore.ETag)(&val)
	}
	if val := resp.Header.Get("x-ms-file-attributes"); val != "" {
		result.FileAttributes = &val
	}
	if val := resp.Header.Get("x-ms-file-change-time"); val != "" {
		fileChangeTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientRenameResponse{}, err
		}
		result.FileChangeTime = &fileChangeTime
	}
	if val := resp.Header.Get("x-ms-file-creation-time"); val != "" {
		fileCreationTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientRenameResponse{}, err
		}
		result.FileCreationTime = &fileCreationTime
	}
	if val := resp.Header.Get("x-ms-file-id"); val != "" {
		result.ID = &val
	}
	if val := resp.Header.Get("x-ms-file-last-write-time"); val != "" {
		fileLastWriteTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientRenameResponse{}, err
		}
		result.FileLastWriteTime = &fileLastWriteTime
	}
	if val := resp.Header.Get("x-ms-file-parent-id"); val != "" {
		result.ParentID = &val
	}
	if val := resp.Header.Get("x-ms-file-permission-key"); val != "" {
		result.FilePermissionKey = &val
	}
	if val := resp.Header.Get("x-ms-request-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return DirectoryClientRenameResponse{}, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("Last-Modified"); val != "" {
		lastModified, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientRenameResponse{}, err
		}
		result.LastModified = &lastModified
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// SetMetadata - Updates user defined metadata for the specified directory.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-05
//   - options - DirectoryClientSetMetadataOptions contains the optional parameters for the DirectoryClient.SetMetadata method.
func (client *DirectoryClient) SetMetadata(ctx context.Context, options *DirectoryClientSetMetadataOptions) (DirectoryClientSetMetadataResponse, error) {
	var err error
	req, err := client.setMetadataCreateRequest(ctx, options)
	if err != nil {
		return DirectoryClientSetMetadataResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DirectoryClientSetMetadataResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DirectoryClientSetMetadataResponse{}, err
	}
	resp, err := client.setMetadataHandleResponse(httpResp)
	return resp, err
}

// setMetadataCreateRequest creates the SetMetadata request.
func (client *DirectoryClient) setMetadataCreateRequest(ctx context.Context, options *DirectoryClientSetMetadataOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPut, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "directory")
	reqQP.Set("comp", "metadata")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.Metadata != nil {
		for k, v := range options.Metadata {
			if v != nil {
				req.Raw().Header["x-ms-meta-"+k] = []string{*v}
			}
		}
	}
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// setMetadataHandleResponse handles the SetMetadata response.
func (client *DirectoryClient) setMetadataHandleResponse(resp *http.Response) (DirectoryClientSetMetadataResponse, error) {
	result := DirectoryClientSetMetadataResponse{}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientSetMetadataResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = (*azcore.ETag)(&val)
	}
	if val := resp.Header.Get("x-ms-request-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return DirectoryClientSetMetadataResponse{}, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// SetProperties - Sets properties on the directory.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-05
//   - options - DirectoryClientSetPropertiesOptions contains the optional parameters for the DirectoryClient.SetProperties method.
func (client *DirectoryClient) SetProperties(ctx context.Context, options *DirectoryClientSetPropertiesOptions) (DirectoryClientSetPropertiesResponse, error) {
	var err error
	req, err := client.setPropertiesCreateRequest(ctx, options)
	if err != nil {
		return DirectoryClientSetPropertiesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DirectoryClientSetPropertiesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DirectoryClientSetPropertiesResponse{}, err
	}
	resp, err := client.setPropertiesHandleResponse(httpResp)
	return resp, err
}

// setPropertiesCreateRequest creates the SetProperties request.
func (client *DirectoryClient) setPropertiesCreateRequest(ctx context.Context, options *DirectoryClientSetPropertiesOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPut, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "directory")
	reqQP.Set("comp", "properties")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-version"] = []string{"2025-01-05"}
	if options != nil && options.FilePermission != nil {
		req.Raw().Header["x-ms-file-permission"] = []string{*options.FilePermission}
	}
	if options != nil && options.FilePermissionFormat != nil {
		req.Raw().Header["x-ms-file-permission-format"] = []string{string(*options.FilePermissionFormat)}
	}
	if options != nil && options.FilePermissionKey != nil {
		req.Raw().Header["x-ms-file-permission-key"] = []string{*options.FilePermissionKey}
	}
	if options != nil && options.FileAttributes != nil {
		req.Raw().Header["x-ms-file-attributes"] = []string{*options.FileAttributes}
	}
	if options != nil && options.FileCreationTime != nil {
		req.Raw().Header["x-ms-file-creation-time"] = []string{*options.FileCreationTime}
	}
	if options != nil && options.FileLastWriteTime != nil {
		req.Raw().Header["x-ms-file-last-write-time"] = []string{*options.FileLastWriteTime}
	}
	if options != nil && options.FileChangeTime != nil {
		req.Raw().Header["x-ms-file-change-time"] = []string{*options.FileChangeTime}
	}
	if client.allowTrailingDot != nil {
		req.Raw().Header["x-ms-allow-trailing-dot"] = []string{strconv.FormatBool(*client.allowTrailingDot)}
	}
	if client.fileRequestIntent != nil {
		req.Raw().Header["x-ms-file-request-intent"] = []string{string(*client.fileRequestIntent)}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// setPropertiesHandleResponse handles the SetProperties response.
func (client *DirectoryClient) setPropertiesHandleResponse(resp *http.Response) (DirectoryClientSetPropertiesResponse, error) {
	result := DirectoryClientSetPropertiesResponse{}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientSetPropertiesResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = (*azcore.ETag)(&val)
	}
	if val := resp.Header.Get("x-ms-file-attributes"); val != "" {
		result.FileAttributes = &val
	}
	if val := resp.Header.Get("x-ms-file-change-time"); val != "" {
		fileChangeTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientSetPropertiesResponse{}, err
		}
		result.FileChangeTime = &fileChangeTime
	}
	if val := resp.Header.Get("x-ms-file-creation-time"); val != "" {
		fileCreationTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientSetPropertiesResponse{}, err
		}
		result.FileCreationTime = &fileCreationTime
	}
	if val := resp.Header.Get("x-ms-file-id"); val != "" {
		result.ID = &val
	}
	if val := resp.Header.Get("x-ms-file-last-write-time"); val != "" {
		fileLastWriteTime, err := time.Parse(ISO8601, val)
		if err != nil {
			return DirectoryClientSetPropertiesResponse{}, err
		}
		result.FileLastWriteTime = &fileLastWriteTime
	}
	if val := resp.Header.Get("x-ms-file-parent-id"); val != "" {
		result.ParentID = &val
	}
	if val := resp.Header.Get("x-ms-file-permission-key"); val != "" {
		result.FilePermissionKey = &val
	}
	if val := resp.Header.Get("x-ms-request-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return DirectoryClientSetPropertiesResponse{}, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("Last-Modified"); val != "" {
		lastModified, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return DirectoryClientSetPropertiesResponse{}, err
		}
		result.LastModified = &lastModified
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

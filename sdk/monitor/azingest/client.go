//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package azingest

import (
	"bytes"
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"net/http"
	"net/url"
	"strings"
)

// Client contains the methods for the Client group.
// Don't use this type directly, use a constructor function instead.
type Client struct {
	internal *azcore.Client
	endpoint string
}

// Upload - See error response code and error response message for more detail.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
//   - ruleID - The immutable Id of the Data Collection Rule resource.
//   - stream - The streamDeclaration name as defined in the Data Collection Rule.
//   - body - An array of objects matching the schema defined by the provided stream.
//   - options - UploadOptions contains the optional parameters for the Client.Upload method.
func (client *Client) Upload(ctx context.Context, ruleID string, stream string, body []byte, options *UploadOptions) (UploadResponse, error) {
	req, err := client.uploadCreateRequest(ctx, ruleID, stream, body, options)
	if err != nil {
		return UploadResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return UploadResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusNoContent) {
		return UploadResponse{}, runtime.NewResponseError(resp)
	}
	return UploadResponse{}, nil
}

// uploadCreateRequest creates the Upload request.
func (client *Client) uploadCreateRequest(ctx context.Context, ruleID string, stream string, body []byte, options *UploadOptions) (*policy.Request, error) {
	urlPath := "/dataCollectionRules/{ruleId}/streams/{stream}"
	if ruleID == "" {
		return nil, errors.New("parameter ruleID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ruleId}", url.PathEscape(ruleID))
	if stream == "" {
		return nil, errors.New("parameter stream cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{stream}", url.PathEscape(stream))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.ContentEncoding != nil {
		req.Raw().Header["Content-Encoding"] = []string{*options.ContentEncoding}
	}
	if options != nil && options.XMSClientRequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.XMSClientRequestID}
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, req.SetBody(streaming.NopCloser(bytes.NewReader(body)), "application/json")
}

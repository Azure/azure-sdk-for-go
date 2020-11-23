// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcompute

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GalleryImagesClient contains the methods for the GalleryImages group.
// Don't use this type directly, use NewGalleryImagesClient() instead.
type GalleryImagesClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewGalleryImagesClient creates a new instance of GalleryImagesClient with the specified values.
func NewGalleryImagesClient(con *armcore.Connection, subscriptionID string) GalleryImagesClient {
	return GalleryImagesClient{con: con, subscriptionID: subscriptionID}
}

// Pipeline returns the pipeline associated with this client.
func (client GalleryImagesClient) Pipeline() azcore.Pipeline {
	return client.con.Pipeline()
}

// BeginCreateOrUpdate - Create or update a gallery image definition.
func (client GalleryImagesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, galleryImage GalleryImage, options *GalleryImagesCreateOrUpdateOptions) (*GalleryImagePollerResponse, error) {
	resp, err := client.CreateOrUpdate(ctx, resourceGroupName, galleryName, galleryImageName, galleryImage, options)
	if err != nil {
		return nil, err
	}
	result := &GalleryImagePollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("GalleryImagesClient.CreateOrUpdate", "", resp, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	poller := &galleryImagePoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*GalleryImageResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeCreateOrUpdate creates a new GalleryImagePoller from the specified resume token.
// token - The value must come from a previous call to GalleryImagePoller.ResumeToken().
func (client GalleryImagesClient) ResumeCreateOrUpdate(token string) (GalleryImagePoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("GalleryImagesClient.CreateOrUpdate", token, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	return &galleryImagePoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// CreateOrUpdate - Create or update a gallery image definition.
func (client GalleryImagesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, galleryImage GalleryImage, options *GalleryImagesCreateOrUpdateOptions) (*azcore.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, galleryName, galleryImageName, galleryImage, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client GalleryImagesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, galleryImage GalleryImage, options *GalleryImagesCreateOrUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{galleryImageName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryImageName}", url.PathEscape(galleryImageName))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-09-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(galleryImage)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client GalleryImagesClient) createOrUpdateHandleResponse(resp *azcore.Response) (*GalleryImageResponse, error) {
	result := GalleryImageResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.GalleryImage)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client GalleryImagesClient) createOrUpdateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// BeginDelete - Delete a gallery image.
func (client GalleryImagesClient) BeginDelete(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, options *GalleryImagesDeleteOptions) (*HTTPPollerResponse, error) {
	resp, err := client.Delete(ctx, resourceGroupName, galleryName, galleryImageName, options)
	if err != nil {
		return nil, err
	}
	result := &HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("GalleryImagesClient.Delete", "", resp, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	poller := &httpPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeDelete creates a new HTTPPoller from the specified resume token.
// token - The value must come from a previous call to HTTPPoller.ResumeToken().
func (client GalleryImagesClient) ResumeDelete(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("GalleryImagesClient.Delete", token, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// Delete - Delete a gallery image.
func (client GalleryImagesClient) Delete(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, options *GalleryImagesDeleteOptions) (*azcore.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, galleryName, galleryImageName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client GalleryImagesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, options *GalleryImagesDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{galleryImageName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryImageName}", url.PathEscape(galleryImageName))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-09-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client GalleryImagesClient) deleteHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Get - Retrieves information about a gallery image definition.
func (client GalleryImagesClient) Get(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, options *GalleryImagesGetOptions) (*GalleryImageResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, galleryName, galleryImageName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.getHandleError(resp)
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client GalleryImagesClient) getCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, options *GalleryImagesGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{galleryImageName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryImageName}", url.PathEscape(galleryImageName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-09-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client GalleryImagesClient) getHandleResponse(resp *azcore.Response) (*GalleryImageResponse, error) {
	result := GalleryImageResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.GalleryImage)
}

// getHandleError handles the Get error response.
func (client GalleryImagesClient) getHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// ListByGallery - List gallery image definitions in a gallery.
func (client GalleryImagesClient) ListByGallery(resourceGroupName string, galleryName string, options *GalleryImagesListByGalleryOptions) GalleryImageListPager {
	return &galleryImageListPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listByGalleryCreateRequest(ctx, resourceGroupName, galleryName, options)
		},
		responder: client.listByGalleryHandleResponse,
		errorer:   client.listByGalleryHandleError,
		advancer: func(ctx context.Context, resp *GalleryImageListResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.GalleryImageList.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listByGalleryCreateRequest creates the ListByGallery request.
func (client GalleryImagesClient) listByGalleryCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, options *GalleryImagesListByGalleryOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-09-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listByGalleryHandleResponse handles the ListByGallery response.
func (client GalleryImagesClient) listByGalleryHandleResponse(resp *azcore.Response) (*GalleryImageListResponse, error) {
	result := GalleryImageListResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.GalleryImageList)
}

// listByGalleryHandleError handles the ListByGallery error response.
func (client GalleryImagesClient) listByGalleryHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// BeginUpdate - Update a gallery image definition.
func (client GalleryImagesClient) BeginUpdate(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, galleryImage GalleryImageUpdate, options *GalleryImagesUpdateOptions) (*GalleryImagePollerResponse, error) {
	resp, err := client.Update(ctx, resourceGroupName, galleryName, galleryImageName, galleryImage, options)
	if err != nil {
		return nil, err
	}
	result := &GalleryImagePollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("GalleryImagesClient.Update", "", resp, client.updateHandleError)
	if err != nil {
		return nil, err
	}
	poller := &galleryImagePoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*GalleryImageResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeUpdate creates a new GalleryImagePoller from the specified resume token.
// token - The value must come from a previous call to GalleryImagePoller.ResumeToken().
func (client GalleryImagesClient) ResumeUpdate(token string) (GalleryImagePoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("GalleryImagesClient.Update", token, client.updateHandleError)
	if err != nil {
		return nil, err
	}
	return &galleryImagePoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// Update - Update a gallery image definition.
func (client GalleryImagesClient) Update(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, galleryImage GalleryImageUpdate, options *GalleryImagesUpdateOptions) (*azcore.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, galleryName, galleryImageName, galleryImage, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.updateHandleError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client GalleryImagesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, galleryImage GalleryImageUpdate, options *GalleryImagesUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{galleryImageName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	urlPath = strings.ReplaceAll(urlPath, "{galleryImageName}", url.PathEscape(galleryImageName))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-09-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(galleryImage)
}

// updateHandleResponse handles the Update response.
func (client GalleryImagesClient) updateHandleResponse(resp *azcore.Response) (*GalleryImageResponse, error) {
	result := GalleryImageResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.GalleryImage)
}

// updateHandleError handles the Update error response.
func (client GalleryImagesClient) updateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

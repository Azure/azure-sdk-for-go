// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcompute

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// VirtualMachineImagesOperations contains the methods for the VirtualMachineImages group.
type VirtualMachineImagesOperations interface {
	// Get - Gets a virtual machine image.
	Get(ctx context.Context, location string, publisherName string, offer string, skus string, version string, options *VirtualMachineImagesGetOptions) (*VirtualMachineImageResponse, error)
	// List - Gets a list of all virtual machine image versions for the specified location, publisher, offer, and SKU.
	List(ctx context.Context, location string, publisherName string, offer string, skus string, options *VirtualMachineImagesListOptions) (*VirtualMachineImageResourceArrayResponse, error)
	// ListOffers - Gets a list of virtual machine image offers for the specified location and publisher.
	ListOffers(ctx context.Context, location string, publisherName string, options *VirtualMachineImagesListOffersOptions) (*VirtualMachineImageResourceArrayResponse, error)
	// ListPublishers - Gets a list of virtual machine image publishers for the specified Azure location.
	ListPublishers(ctx context.Context, location string, options *VirtualMachineImagesListPublishersOptions) (*VirtualMachineImageResourceArrayResponse, error)
	// ListSKUs - Gets a list of virtual machine image SKUs for the specified location, publisher, and offer.
	ListSKUs(ctx context.Context, location string, publisherName string, offer string, options *VirtualMachineImagesListSKUsOptions) (*VirtualMachineImageResourceArrayResponse, error)
}

// VirtualMachineImagesClient implements the VirtualMachineImagesOperations interface.
// Don't use this type directly, use NewVirtualMachineImagesClient() instead.
type VirtualMachineImagesClient struct {
	*Client
	subscriptionID string
}

// NewVirtualMachineImagesClient creates a new instance of VirtualMachineImagesClient with the specified values.
func NewVirtualMachineImagesClient(c *Client, subscriptionID string) VirtualMachineImagesOperations {
	return &VirtualMachineImagesClient{Client: c, subscriptionID: subscriptionID}
}

// Do invokes the Do() method on the pipeline associated with this client.
func (client *VirtualMachineImagesClient) Do(req *azcore.Request) (*azcore.Response, error) {
	return client.p.Do(req)
}

// Get - Gets a virtual machine image.
func (client *VirtualMachineImagesClient) Get(ctx context.Context, location string, publisherName string, offer string, skus string, version string, options *VirtualMachineImagesGetOptions) (*VirtualMachineImageResponse, error) {
	req, err := client.GetCreateRequest(ctx, location, publisherName, offer, skus, version, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.GetHandleError(resp)
	}
	result, err := client.GetHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCreateRequest creates the Get request.
func (client *VirtualMachineImagesClient) GetCreateRequest(ctx context.Context, location string, publisherName string, offer string, skus string, version string, options *VirtualMachineImagesGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus/{skus}/versions/{version}"
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	urlPath = strings.ReplaceAll(urlPath, "{offer}", url.PathEscape(offer))
	urlPath = strings.ReplaceAll(urlPath, "{skus}", url.PathEscape(skus))
	urlPath = strings.ReplaceAll(urlPath, "{version}", url.PathEscape(version))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2019-12-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetHandleResponse handles the Get response.
func (client *VirtualMachineImagesClient) GetHandleResponse(resp *azcore.Response) (*VirtualMachineImageResponse, error) {
	result := VirtualMachineImageResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VirtualMachineImage)
}

// GetHandleError handles the Get error response.
func (client *VirtualMachineImagesClient) GetHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// List - Gets a list of all virtual machine image versions for the specified location, publisher, offer, and SKU.
func (client *VirtualMachineImagesClient) List(ctx context.Context, location string, publisherName string, offer string, skus string, options *VirtualMachineImagesListOptions) (*VirtualMachineImageResourceArrayResponse, error) {
	req, err := client.ListCreateRequest(ctx, location, publisherName, offer, skus, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.ListHandleError(resp)
	}
	result, err := client.ListHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListCreateRequest creates the List request.
func (client *VirtualMachineImagesClient) ListCreateRequest(ctx context.Context, location string, publisherName string, offer string, skus string, options *VirtualMachineImagesListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus/{skus}/versions"
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	urlPath = strings.ReplaceAll(urlPath, "{offer}", url.PathEscape(offer))
	urlPath = strings.ReplaceAll(urlPath, "{skus}", url.PathEscape(skus))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if options != nil && options.Expand != nil {
		query.Set("$expand", *options.Expand)
	}
	if options != nil && options.Top != nil {
		query.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Orderby != nil {
		query.Set("$orderby", *options.Orderby)
	}
	query.Set("api-version", "2019-12-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListHandleResponse handles the List response.
func (client *VirtualMachineImagesClient) ListHandleResponse(resp *azcore.Response) (*VirtualMachineImageResourceArrayResponse, error) {
	result := VirtualMachineImageResourceArrayResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VirtualMachineImageResourceArray)
}

// ListHandleError handles the List error response.
func (client *VirtualMachineImagesClient) ListHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListOffers - Gets a list of virtual machine image offers for the specified location and publisher.
func (client *VirtualMachineImagesClient) ListOffers(ctx context.Context, location string, publisherName string, options *VirtualMachineImagesListOffersOptions) (*VirtualMachineImageResourceArrayResponse, error) {
	req, err := client.ListOffersCreateRequest(ctx, location, publisherName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.ListOffersHandleError(resp)
	}
	result, err := client.ListOffersHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListOffersCreateRequest creates the ListOffers request.
func (client *VirtualMachineImagesClient) ListOffersCreateRequest(ctx context.Context, location string, publisherName string, options *VirtualMachineImagesListOffersOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmimage/offers"
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2019-12-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListOffersHandleResponse handles the ListOffers response.
func (client *VirtualMachineImagesClient) ListOffersHandleResponse(resp *azcore.Response) (*VirtualMachineImageResourceArrayResponse, error) {
	result := VirtualMachineImageResourceArrayResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VirtualMachineImageResourceArray)
}

// ListOffersHandleError handles the ListOffers error response.
func (client *VirtualMachineImagesClient) ListOffersHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListPublishers - Gets a list of virtual machine image publishers for the specified Azure location.
func (client *VirtualMachineImagesClient) ListPublishers(ctx context.Context, location string, options *VirtualMachineImagesListPublishersOptions) (*VirtualMachineImageResourceArrayResponse, error) {
	req, err := client.ListPublishersCreateRequest(ctx, location, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.ListPublishersHandleError(resp)
	}
	result, err := client.ListPublishersHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListPublishersCreateRequest creates the ListPublishers request.
func (client *VirtualMachineImagesClient) ListPublishersCreateRequest(ctx context.Context, location string, options *VirtualMachineImagesListPublishersOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers"
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2019-12-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListPublishersHandleResponse handles the ListPublishers response.
func (client *VirtualMachineImagesClient) ListPublishersHandleResponse(resp *azcore.Response) (*VirtualMachineImageResourceArrayResponse, error) {
	result := VirtualMachineImageResourceArrayResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VirtualMachineImageResourceArray)
}

// ListPublishersHandleError handles the ListPublishers error response.
func (client *VirtualMachineImagesClient) ListPublishersHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListSKUs - Gets a list of virtual machine image SKUs for the specified location, publisher, and offer.
func (client *VirtualMachineImagesClient) ListSKUs(ctx context.Context, location string, publisherName string, offer string, options *VirtualMachineImagesListSKUsOptions) (*VirtualMachineImageResourceArrayResponse, error) {
	req, err := client.ListSKUsCreateRequest(ctx, location, publisherName, offer, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.ListSKUsHandleError(resp)
	}
	result, err := client.ListSKUsHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListSKUsCreateRequest creates the ListSKUs request.
func (client *VirtualMachineImagesClient) ListSKUsCreateRequest(ctx context.Context, location string, publisherName string, offer string, options *VirtualMachineImagesListSKUsOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus"
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	urlPath = strings.ReplaceAll(urlPath, "{offer}", url.PathEscape(offer))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2019-12-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListSKUsHandleResponse handles the ListSKUs response.
func (client *VirtualMachineImagesClient) ListSKUsHandleResponse(resp *azcore.Response) (*VirtualMachineImageResourceArrayResponse, error) {
	result := VirtualMachineImageResourceArrayResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VirtualMachineImageResourceArray)
}

// ListSKUsHandleError handles the ListSKUs error response.
func (client *VirtualMachineImagesClient) ListSKUsHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcompute

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// VirtualMachineImagesEdgeZoneClient contains the methods for the VirtualMachineImagesEdgeZone group.
// Don't use this type directly, use NewVirtualMachineImagesEdgeZoneClient() instead.
type VirtualMachineImagesEdgeZoneClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewVirtualMachineImagesEdgeZoneClient creates a new instance of VirtualMachineImagesEdgeZoneClient with the specified values.
// subscriptionID - Subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID forms
// part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewVirtualMachineImagesEdgeZoneClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VirtualMachineImagesEdgeZoneClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &VirtualMachineImagesEdgeZoneClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Get - Gets a virtual machine image in an edge zone.
// If the operation fails it returns an *azcore.ResponseError type.
// location - The name of a supported Azure region.
// edgeZone - The name of the edge zone.
// publisherName - A valid image publisher.
// offer - A valid image publisher offer.
// skus - A valid image SKU.
// version - A valid image SKU version.
// options - VirtualMachineImagesEdgeZoneClientGetOptions contains the optional parameters for the VirtualMachineImagesEdgeZoneClient.Get
// method.
func (client *VirtualMachineImagesEdgeZoneClient) Get(ctx context.Context, location string, edgeZone string, publisherName string, offer string, skus string, version string, options *VirtualMachineImagesEdgeZoneClientGetOptions) (VirtualMachineImagesEdgeZoneClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, location, edgeZone, publisherName, offer, skus, version, options)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VirtualMachineImagesEdgeZoneClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VirtualMachineImagesEdgeZoneClient) getCreateRequest(ctx context.Context, location string, edgeZone string, publisherName string, offer string, skus string, version string, options *VirtualMachineImagesEdgeZoneClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/edgeZones/{edgeZone}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus/{skus}/versions/{version}"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if edgeZone == "" {
		return nil, errors.New("parameter edgeZone cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeZone}", url.PathEscape(edgeZone))
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	if offer == "" {
		return nil, errors.New("parameter offer cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{offer}", url.PathEscape(offer))
	if skus == "" {
		return nil, errors.New("parameter skus cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{skus}", url.PathEscape(skus))
	if version == "" {
		return nil, errors.New("parameter version cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{version}", url.PathEscape(version))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VirtualMachineImagesEdgeZoneClient) getHandleResponse(resp *http.Response) (VirtualMachineImagesEdgeZoneClientGetResponse, error) {
	result := VirtualMachineImagesEdgeZoneClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualMachineImage); err != nil {
		return VirtualMachineImagesEdgeZoneClientGetResponse{}, err
	}
	return result, nil
}

// List - Gets a list of all virtual machine image versions for the specified location, edge zone, publisher, offer, and SKU.
// If the operation fails it returns an *azcore.ResponseError type.
// location - The name of a supported Azure region.
// edgeZone - The name of the edge zone.
// publisherName - A valid image publisher.
// offer - A valid image publisher offer.
// skus - A valid image SKU.
// options - VirtualMachineImagesEdgeZoneClientListOptions contains the optional parameters for the VirtualMachineImagesEdgeZoneClient.List
// method.
func (client *VirtualMachineImagesEdgeZoneClient) List(ctx context.Context, location string, edgeZone string, publisherName string, offer string, skus string, options *VirtualMachineImagesEdgeZoneClientListOptions) (VirtualMachineImagesEdgeZoneClientListResponse, error) {
	req, err := client.listCreateRequest(ctx, location, edgeZone, publisherName, offer, skus, options)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientListResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VirtualMachineImagesEdgeZoneClientListResponse{}, runtime.NewResponseError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *VirtualMachineImagesEdgeZoneClient) listCreateRequest(ctx context.Context, location string, edgeZone string, publisherName string, offer string, skus string, options *VirtualMachineImagesEdgeZoneClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/edgeZones/{edgeZone}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus/{skus}/versions"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if edgeZone == "" {
		return nil, errors.New("parameter edgeZone cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeZone}", url.PathEscape(edgeZone))
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	if offer == "" {
		return nil, errors.New("parameter offer cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{offer}", url.PathEscape(offer))
	if skus == "" {
		return nil, errors.New("parameter skus cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{skus}", url.PathEscape(skus))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *VirtualMachineImagesEdgeZoneClient) listHandleResponse(resp *http.Response) (VirtualMachineImagesEdgeZoneClientListResponse, error) {
	result := VirtualMachineImagesEdgeZoneClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualMachineImageResourceArray); err != nil {
		return VirtualMachineImagesEdgeZoneClientListResponse{}, err
	}
	return result, nil
}

// ListOffers - Gets a list of virtual machine image offers for the specified location, edge zone and publisher.
// If the operation fails it returns an *azcore.ResponseError type.
// location - The name of a supported Azure region.
// edgeZone - The name of the edge zone.
// publisherName - A valid image publisher.
// options - VirtualMachineImagesEdgeZoneClientListOffersOptions contains the optional parameters for the VirtualMachineImagesEdgeZoneClient.ListOffers
// method.
func (client *VirtualMachineImagesEdgeZoneClient) ListOffers(ctx context.Context, location string, edgeZone string, publisherName string, options *VirtualMachineImagesEdgeZoneClientListOffersOptions) (VirtualMachineImagesEdgeZoneClientListOffersResponse, error) {
	req, err := client.listOffersCreateRequest(ctx, location, edgeZone, publisherName, options)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientListOffersResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientListOffersResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VirtualMachineImagesEdgeZoneClientListOffersResponse{}, runtime.NewResponseError(resp)
	}
	return client.listOffersHandleResponse(resp)
}

// listOffersCreateRequest creates the ListOffers request.
func (client *VirtualMachineImagesEdgeZoneClient) listOffersCreateRequest(ctx context.Context, location string, edgeZone string, publisherName string, options *VirtualMachineImagesEdgeZoneClientListOffersOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/edgeZones/{edgeZone}/publishers/{publisherName}/artifacttypes/vmimage/offers"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if edgeZone == "" {
		return nil, errors.New("parameter edgeZone cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeZone}", url.PathEscape(edgeZone))
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listOffersHandleResponse handles the ListOffers response.
func (client *VirtualMachineImagesEdgeZoneClient) listOffersHandleResponse(resp *http.Response) (VirtualMachineImagesEdgeZoneClientListOffersResponse, error) {
	result := VirtualMachineImagesEdgeZoneClientListOffersResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualMachineImageResourceArray); err != nil {
		return VirtualMachineImagesEdgeZoneClientListOffersResponse{}, err
	}
	return result, nil
}

// ListPublishers - Gets a list of virtual machine image publishers for the specified Azure location and edge zone.
// If the operation fails it returns an *azcore.ResponseError type.
// location - The name of a supported Azure region.
// edgeZone - The name of the edge zone.
// options - VirtualMachineImagesEdgeZoneClientListPublishersOptions contains the optional parameters for the VirtualMachineImagesEdgeZoneClient.ListPublishers
// method.
func (client *VirtualMachineImagesEdgeZoneClient) ListPublishers(ctx context.Context, location string, edgeZone string, options *VirtualMachineImagesEdgeZoneClientListPublishersOptions) (VirtualMachineImagesEdgeZoneClientListPublishersResponse, error) {
	req, err := client.listPublishersCreateRequest(ctx, location, edgeZone, options)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientListPublishersResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientListPublishersResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VirtualMachineImagesEdgeZoneClientListPublishersResponse{}, runtime.NewResponseError(resp)
	}
	return client.listPublishersHandleResponse(resp)
}

// listPublishersCreateRequest creates the ListPublishers request.
func (client *VirtualMachineImagesEdgeZoneClient) listPublishersCreateRequest(ctx context.Context, location string, edgeZone string, options *VirtualMachineImagesEdgeZoneClientListPublishersOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/edgeZones/{edgeZone}/publishers"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if edgeZone == "" {
		return nil, errors.New("parameter edgeZone cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeZone}", url.PathEscape(edgeZone))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listPublishersHandleResponse handles the ListPublishers response.
func (client *VirtualMachineImagesEdgeZoneClient) listPublishersHandleResponse(resp *http.Response) (VirtualMachineImagesEdgeZoneClientListPublishersResponse, error) {
	result := VirtualMachineImagesEdgeZoneClientListPublishersResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualMachineImageResourceArray); err != nil {
		return VirtualMachineImagesEdgeZoneClientListPublishersResponse{}, err
	}
	return result, nil
}

// ListSKUs - Gets a list of virtual machine image SKUs for the specified location, edge zone, publisher, and offer.
// If the operation fails it returns an *azcore.ResponseError type.
// location - The name of a supported Azure region.
// edgeZone - The name of the edge zone.
// publisherName - A valid image publisher.
// offer - A valid image publisher offer.
// options - VirtualMachineImagesEdgeZoneClientListSKUsOptions contains the optional parameters for the VirtualMachineImagesEdgeZoneClient.ListSKUs
// method.
func (client *VirtualMachineImagesEdgeZoneClient) ListSKUs(ctx context.Context, location string, edgeZone string, publisherName string, offer string, options *VirtualMachineImagesEdgeZoneClientListSKUsOptions) (VirtualMachineImagesEdgeZoneClientListSKUsResponse, error) {
	req, err := client.listSKUsCreateRequest(ctx, location, edgeZone, publisherName, offer, options)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientListSKUsResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineImagesEdgeZoneClientListSKUsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VirtualMachineImagesEdgeZoneClientListSKUsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listSKUsHandleResponse(resp)
}

// listSKUsCreateRequest creates the ListSKUs request.
func (client *VirtualMachineImagesEdgeZoneClient) listSKUsCreateRequest(ctx context.Context, location string, edgeZone string, publisherName string, offer string, options *VirtualMachineImagesEdgeZoneClientListSKUsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/edgeZones/{edgeZone}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if edgeZone == "" {
		return nil, errors.New("parameter edgeZone cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeZone}", url.PathEscape(edgeZone))
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	if offer == "" {
		return nil, errors.New("parameter offer cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{offer}", url.PathEscape(offer))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listSKUsHandleResponse handles the ListSKUs response.
func (client *VirtualMachineImagesEdgeZoneClient) listSKUsHandleResponse(resp *http.Response) (VirtualMachineImagesEdgeZoneClientListSKUsResponse, error) {
	result := VirtualMachineImagesEdgeZoneClientListSKUsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualMachineImageResourceArray); err != nil {
		return VirtualMachineImagesEdgeZoneClientListSKUsResponse{}, err
	}
	return result, nil
}

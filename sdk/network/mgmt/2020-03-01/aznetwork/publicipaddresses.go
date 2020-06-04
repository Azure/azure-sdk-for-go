// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package aznetwork

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PublicIPAddressesOperations contains the methods for the PublicIPAddresses group.
type PublicIPAddressesOperations interface {
	// BeginCreateOrUpdate - Creates or updates a static or dynamic public IP address.
	BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, publicIPAddressName string, parameters PublicIPAddress) (*PublicIPAddressResponse, error)
	// ResumeCreateOrUpdate - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeCreateOrUpdate(token string) (PublicIPAddressPoller, error)
	// BeginDelete - Deletes the specified public IP address.
	BeginDelete(ctx context.Context, resourceGroupName string, publicIPAddressName string) (*HTTPResponse, error)
	// ResumeDelete - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeDelete(token string) (HTTPPoller, error)
	// Get - Gets the specified public IP address in a specified resource group.
	Get(ctx context.Context, resourceGroupName string, publicIPAddressName string, publicIPAddressesGetOptions *PublicIPAddressesGetOptions) (*PublicIPAddressResponse, error)
	// GetVirtualMachineScaleSetPublicIPAddress - Get the specified public IP address in a virtual machine scale set.
	GetVirtualMachineScaleSetPublicIPAddress(ctx context.Context, resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, networkInterfaceName string, ipConfigurationName string, publicIPAddressName string, publicIPAddressesGetVirtualMachineScaleSetPublicIpaddressOptions *PublicIPAddressesGetVirtualMachineScaleSetPublicIPAddressOptions) (*PublicIPAddressResponse, error)
	// List - Gets all public IP addresses in a resource group.
	List(resourceGroupName string) (PublicIPAddressListResultPager, error)
	// ListAll - Gets all the public IP addresses in a subscription.
	ListAll() (PublicIPAddressListResultPager, error)
	// ListVirtualMachineScaleSetPublicIPAddresses - Gets information about all public IP addresses on a virtual machine scale set level.
	ListVirtualMachineScaleSetPublicIPAddresses(resourceGroupName string, virtualMachineScaleSetName string) (PublicIPAddressListResultPager, error)
	// ListVirtualMachineScaleSetVMPublicIPaddresses - Gets information about all public IP addresses in a virtual machine IP configuration in a virtual machine scale set.
	ListVirtualMachineScaleSetVMPublicIPaddresses(resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, networkInterfaceName string, ipConfigurationName string) (PublicIPAddressListResultPager, error)
	// UpdateTags - Updates public IP address tags.
	UpdateTags(ctx context.Context, resourceGroupName string, publicIPAddressName string, parameters TagsObject) (*PublicIPAddressResponse, error)
}

// publicIPAddressesOperations implements the PublicIPAddressesOperations interface.
type publicIPAddressesOperations struct {
	*Client
	subscriptionID string
}

// CreateOrUpdate - Creates or updates a static or dynamic public IP address.
func (client *publicIPAddressesOperations) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, publicIPAddressName string, parameters PublicIPAddress) (*PublicIPAddressResponse, error) {
	req, err := client.createOrUpdateCreateRequest(resourceGroupName, publicIPAddressName, parameters)
	if err != nil {
		return nil, err
	}
	// send the first request to initialize the poller
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.createOrUpdateHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	pt, err := createPollingTracker("publicIPAddressesOperations.CreateOrUpdate", "azure-async-operation", resp, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	poller := &publicIPAddressPoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*PublicIPAddressResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *publicIPAddressesOperations) ResumeCreateOrUpdate(token string) (PublicIPAddressPoller, error) {
	pt, err := resumePollingTracker("publicIPAddressesOperations.CreateOrUpdate", token, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	return &publicIPAddressPoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *publicIPAddressesOperations) createOrUpdateCreateRequest(resourceGroupName string, publicIPAddressName string, parameters PublicIPAddress) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPAddresses/{publicIpAddressName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{publicIpAddressName}", url.PathEscape(publicIPAddressName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodPut, *u)
	return req, req.MarshalAsJSON(parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *publicIPAddressesOperations) createOrUpdateHandleResponse(resp *azcore.Response) (*PublicIPAddressResponse, error) {
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated, http.StatusNoContent) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	result := PublicIPAddressResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PublicIPAddress)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *publicIPAddressesOperations) createOrUpdateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// Delete - Deletes the specified public IP address.
func (client *publicIPAddressesOperations) BeginDelete(ctx context.Context, resourceGroupName string, publicIPAddressName string) (*HTTPResponse, error) {
	req, err := client.deleteCreateRequest(resourceGroupName, publicIPAddressName)
	if err != nil {
		return nil, err
	}
	// send the first request to initialize the poller
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.deleteHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	pt, err := createPollingTracker("publicIPAddressesOperations.Delete", "location", resp, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	poller := &httpPoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *publicIPAddressesOperations) ResumeDelete(token string) (HTTPPoller, error) {
	pt, err := resumePollingTracker("publicIPAddressesOperations.Delete", token, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *publicIPAddressesOperations) deleteCreateRequest(resourceGroupName string, publicIPAddressName string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPAddresses/{publicIpAddressName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{publicIpAddressName}", url.PathEscape(publicIPAddressName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodDelete, *u)
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *publicIPAddressesOperations) deleteHandleResponse(resp *azcore.Response) (*HTTPResponse, error) {
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return &HTTPResponse{RawResponse: resp.Response}, nil
}

// deleteHandleError handles the Delete error response.
func (client *publicIPAddressesOperations) deleteHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// Get - Gets the specified public IP address in a specified resource group.
func (client *publicIPAddressesOperations) Get(ctx context.Context, resourceGroupName string, publicIPAddressName string, publicIPAddressesGetOptions *PublicIPAddressesGetOptions) (*PublicIPAddressResponse, error) {
	req, err := client.getCreateRequest(resourceGroupName, publicIPAddressName, publicIPAddressesGetOptions)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client *publicIPAddressesOperations) getCreateRequest(resourceGroupName string, publicIPAddressName string, publicIPAddressesGetOptions *PublicIPAddressesGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPAddresses/{publicIpAddressName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{publicIpAddressName}", url.PathEscape(publicIPAddressName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	if publicIPAddressesGetOptions != nil && publicIPAddressesGetOptions.Expand != nil {
		query.Set("$expand", *publicIPAddressesGetOptions.Expand)
	}
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *publicIPAddressesOperations) getHandleResponse(resp *azcore.Response) (*PublicIPAddressResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.getHandleError(resp)
	}
	result := PublicIPAddressResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PublicIPAddress)
}

// getHandleError handles the Get error response.
func (client *publicIPAddressesOperations) getHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// GetVirtualMachineScaleSetPublicIPAddress - Get the specified public IP address in a virtual machine scale set.
func (client *publicIPAddressesOperations) GetVirtualMachineScaleSetPublicIPAddress(ctx context.Context, resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, networkInterfaceName string, ipConfigurationName string, publicIPAddressName string, publicIPAddressesGetVirtualMachineScaleSetPublicIpaddressOptions *PublicIPAddressesGetVirtualMachineScaleSetPublicIPAddressOptions) (*PublicIPAddressResponse, error) {
	req, err := client.getVirtualMachineScaleSetPublicIPAddressCreateRequest(resourceGroupName, virtualMachineScaleSetName, virtualmachineIndex, networkInterfaceName, ipConfigurationName, publicIPAddressName, publicIPAddressesGetVirtualMachineScaleSetPublicIpaddressOptions)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.getVirtualMachineScaleSetPublicIPAddressHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getVirtualMachineScaleSetPublicIPAddressCreateRequest creates the GetVirtualMachineScaleSetPublicIPAddress request.
func (client *publicIPAddressesOperations) getVirtualMachineScaleSetPublicIPAddressCreateRequest(resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, networkInterfaceName string, ipConfigurationName string, publicIPAddressName string, publicIPAddressesGetVirtualMachineScaleSetPublicIpaddressOptions *PublicIPAddressesGetVirtualMachineScaleSetPublicIPAddressOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{virtualMachineScaleSetName}/virtualMachines/{virtualmachineIndex}/networkInterfaces/{networkInterfaceName}/ipconfigurations/{ipConfigurationName}/publicipaddresses/{publicIpAddressName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineScaleSetName}", url.PathEscape(virtualMachineScaleSetName))
	urlPath = strings.ReplaceAll(urlPath, "{virtualmachineIndex}", url.PathEscape(virtualmachineIndex))
	urlPath = strings.ReplaceAll(urlPath, "{networkInterfaceName}", url.PathEscape(networkInterfaceName))
	urlPath = strings.ReplaceAll(urlPath, "{ipConfigurationName}", url.PathEscape(ipConfigurationName))
	urlPath = strings.ReplaceAll(urlPath, "{publicIpAddressName}", url.PathEscape(publicIPAddressName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2018-10-01")
	if publicIPAddressesGetVirtualMachineScaleSetPublicIpaddressOptions != nil && publicIPAddressesGetVirtualMachineScaleSetPublicIpaddressOptions.Expand != nil {
		query.Set("$expand", *publicIPAddressesGetVirtualMachineScaleSetPublicIpaddressOptions.Expand)
	}
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// getVirtualMachineScaleSetPublicIPAddressHandleResponse handles the GetVirtualMachineScaleSetPublicIPAddress response.
func (client *publicIPAddressesOperations) getVirtualMachineScaleSetPublicIPAddressHandleResponse(resp *azcore.Response) (*PublicIPAddressResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.getVirtualMachineScaleSetPublicIPAddressHandleError(resp)
	}
	result := PublicIPAddressResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PublicIPAddress)
}

// getVirtualMachineScaleSetPublicIPAddressHandleError handles the GetVirtualMachineScaleSetPublicIPAddress error response.
func (client *publicIPAddressesOperations) getVirtualMachineScaleSetPublicIPAddressHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// List - Gets all public IP addresses in a resource group.
func (client *publicIPAddressesOperations) List(resourceGroupName string) (PublicIPAddressListResultPager, error) {
	req, err := client.listCreateRequest(resourceGroupName)
	if err != nil {
		return nil, err
	}
	return &publicIPAddressListResultPager{
		pipeline:  client.p,
		request:   req,
		responder: client.listHandleResponse,
		advancer: func(resp *PublicIPAddressListResultResponse) (*azcore.Request, error) {
			u, err := url.Parse(*resp.PublicIPAddressListResult.NextLink)
			if err != nil {
				return nil, fmt.Errorf("invalid NextLink: %w", err)
			}
			if u.Scheme == "" {
				return nil, fmt.Errorf("no scheme detected in NextLink %s", *resp.PublicIPAddressListResult.NextLink)
			}
			return azcore.NewRequest(http.MethodGet, *u), nil
		},
	}, nil
}

// listCreateRequest creates the List request.
func (client *publicIPAddressesOperations) listCreateRequest(resourceGroupName string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPAddresses"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// listHandleResponse handles the List response.
func (client *publicIPAddressesOperations) listHandleResponse(resp *azcore.Response) (*PublicIPAddressListResultResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.listHandleError(resp)
	}
	result := PublicIPAddressListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PublicIPAddressListResult)
}

// listHandleError handles the List error response.
func (client *publicIPAddressesOperations) listHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// ListAll - Gets all the public IP addresses in a subscription.
func (client *publicIPAddressesOperations) ListAll() (PublicIPAddressListResultPager, error) {
	req, err := client.listAllCreateRequest()
	if err != nil {
		return nil, err
	}
	return &publicIPAddressListResultPager{
		pipeline:  client.p,
		request:   req,
		responder: client.listAllHandleResponse,
		advancer: func(resp *PublicIPAddressListResultResponse) (*azcore.Request, error) {
			u, err := url.Parse(*resp.PublicIPAddressListResult.NextLink)
			if err != nil {
				return nil, fmt.Errorf("invalid NextLink: %w", err)
			}
			if u.Scheme == "" {
				return nil, fmt.Errorf("no scheme detected in NextLink %s", *resp.PublicIPAddressListResult.NextLink)
			}
			return azcore.NewRequest(http.MethodGet, *u), nil
		},
	}, nil
}

// listAllCreateRequest creates the ListAll request.
func (client *publicIPAddressesOperations) listAllCreateRequest() (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/publicIPAddresses"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// listAllHandleResponse handles the ListAll response.
func (client *publicIPAddressesOperations) listAllHandleResponse(resp *azcore.Response) (*PublicIPAddressListResultResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.listAllHandleError(resp)
	}
	result := PublicIPAddressListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PublicIPAddressListResult)
}

// listAllHandleError handles the ListAll error response.
func (client *publicIPAddressesOperations) listAllHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// ListVirtualMachineScaleSetPublicIPAddresses - Gets information about all public IP addresses on a virtual machine scale set level.
func (client *publicIPAddressesOperations) ListVirtualMachineScaleSetPublicIPAddresses(resourceGroupName string, virtualMachineScaleSetName string) (PublicIPAddressListResultPager, error) {
	req, err := client.listVirtualMachineScaleSetPublicIPAddressesCreateRequest(resourceGroupName, virtualMachineScaleSetName)
	if err != nil {
		return nil, err
	}
	return &publicIPAddressListResultPager{
		pipeline:  client.p,
		request:   req,
		responder: client.listVirtualMachineScaleSetPublicIPAddressesHandleResponse,
		advancer: func(resp *PublicIPAddressListResultResponse) (*azcore.Request, error) {
			u, err := url.Parse(*resp.PublicIPAddressListResult.NextLink)
			if err != nil {
				return nil, fmt.Errorf("invalid NextLink: %w", err)
			}
			if u.Scheme == "" {
				return nil, fmt.Errorf("no scheme detected in NextLink %s", *resp.PublicIPAddressListResult.NextLink)
			}
			return azcore.NewRequest(http.MethodGet, *u), nil
		},
	}, nil
}

// listVirtualMachineScaleSetPublicIPAddressesCreateRequest creates the ListVirtualMachineScaleSetPublicIPAddresses request.
func (client *publicIPAddressesOperations) listVirtualMachineScaleSetPublicIPAddressesCreateRequest(resourceGroupName string, virtualMachineScaleSetName string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{virtualMachineScaleSetName}/publicipaddresses"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineScaleSetName}", url.PathEscape(virtualMachineScaleSetName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2018-10-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// listVirtualMachineScaleSetPublicIPAddressesHandleResponse handles the ListVirtualMachineScaleSetPublicIPAddresses response.
func (client *publicIPAddressesOperations) listVirtualMachineScaleSetPublicIPAddressesHandleResponse(resp *azcore.Response) (*PublicIPAddressListResultResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.listVirtualMachineScaleSetPublicIPAddressesHandleError(resp)
	}
	result := PublicIPAddressListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PublicIPAddressListResult)
}

// listVirtualMachineScaleSetPublicIPAddressesHandleError handles the ListVirtualMachineScaleSetPublicIPAddresses error response.
func (client *publicIPAddressesOperations) listVirtualMachineScaleSetPublicIPAddressesHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// ListVirtualMachineScaleSetVMPublicIPaddresses - Gets information about all public IP addresses in a virtual machine IP configuration in a virtual machine scale set.
func (client *publicIPAddressesOperations) ListVirtualMachineScaleSetVMPublicIPaddresses(resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, networkInterfaceName string, ipConfigurationName string) (PublicIPAddressListResultPager, error) {
	req, err := client.listVirtualMachineScaleSetVMPublicIpaddressesCreateRequest(resourceGroupName, virtualMachineScaleSetName, virtualmachineIndex, networkInterfaceName, ipConfigurationName)
	if err != nil {
		return nil, err
	}
	return &publicIPAddressListResultPager{
		pipeline:  client.p,
		request:   req,
		responder: client.listVirtualMachineScaleSetVMPublicIpaddressesHandleResponse,
		advancer: func(resp *PublicIPAddressListResultResponse) (*azcore.Request, error) {
			u, err := url.Parse(*resp.PublicIPAddressListResult.NextLink)
			if err != nil {
				return nil, fmt.Errorf("invalid NextLink: %w", err)
			}
			if u.Scheme == "" {
				return nil, fmt.Errorf("no scheme detected in NextLink %s", *resp.PublicIPAddressListResult.NextLink)
			}
			return azcore.NewRequest(http.MethodGet, *u), nil
		},
	}, nil
}

// listVirtualMachineScaleSetVMPublicIpaddressesCreateRequest creates the ListVirtualMachineScaleSetVMPublicIPaddresses request.
func (client *publicIPAddressesOperations) listVirtualMachineScaleSetVMPublicIpaddressesCreateRequest(resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, networkInterfaceName string, ipConfigurationName string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{virtualMachineScaleSetName}/virtualMachines/{virtualmachineIndex}/networkInterfaces/{networkInterfaceName}/ipconfigurations/{ipConfigurationName}/publicipaddresses"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineScaleSetName}", url.PathEscape(virtualMachineScaleSetName))
	urlPath = strings.ReplaceAll(urlPath, "{virtualmachineIndex}", url.PathEscape(virtualmachineIndex))
	urlPath = strings.ReplaceAll(urlPath, "{networkInterfaceName}", url.PathEscape(networkInterfaceName))
	urlPath = strings.ReplaceAll(urlPath, "{ipConfigurationName}", url.PathEscape(ipConfigurationName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2018-10-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// listVirtualMachineScaleSetVMPublicIpaddressesHandleResponse handles the ListVirtualMachineScaleSetVMPublicIPaddresses response.
func (client *publicIPAddressesOperations) listVirtualMachineScaleSetVMPublicIpaddressesHandleResponse(resp *azcore.Response) (*PublicIPAddressListResultResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.listVirtualMachineScaleSetVMPublicIpaddressesHandleError(resp)
	}
	result := PublicIPAddressListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PublicIPAddressListResult)
}

// listVirtualMachineScaleSetVMPublicIpaddressesHandleError handles the ListVirtualMachineScaleSetVMPublicIPaddresses error response.
func (client *publicIPAddressesOperations) listVirtualMachineScaleSetVMPublicIpaddressesHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// UpdateTags - Updates public IP address tags.
func (client *publicIPAddressesOperations) UpdateTags(ctx context.Context, resourceGroupName string, publicIPAddressName string, parameters TagsObject) (*PublicIPAddressResponse, error) {
	req, err := client.updateTagsCreateRequest(resourceGroupName, publicIPAddressName, parameters)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.updateTagsHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *publicIPAddressesOperations) updateTagsCreateRequest(resourceGroupName string, publicIPAddressName string, parameters TagsObject) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPAddresses/{publicIpAddressName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{publicIpAddressName}", url.PathEscape(publicIPAddressName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodPatch, *u)
	return req, req.MarshalAsJSON(parameters)
}

// updateTagsHandleResponse handles the UpdateTags response.
func (client *publicIPAddressesOperations) updateTagsHandleResponse(resp *azcore.Response) (*PublicIPAddressResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.updateTagsHandleError(resp)
	}
	result := PublicIPAddressResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PublicIPAddress)
}

// updateTagsHandleError handles the UpdateTags error response.
func (client *publicIPAddressesOperations) updateTagsHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

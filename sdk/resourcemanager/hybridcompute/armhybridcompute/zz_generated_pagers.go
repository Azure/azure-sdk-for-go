//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridcompute

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"reflect"
)

// MachineExtensionsListPager provides operations for iterating over paged responses.
type MachineExtensionsListPager struct {
	client    *MachineExtensionsClient
	current   MachineExtensionsListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, MachineExtensionsListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *MachineExtensionsListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *MachineExtensionsListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.MachineExtensionsListResult.NextLink == nil || len(*p.current.MachineExtensionsListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.pl.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listHandleError(resp)
		return false
	}
	result, err := p.client.listHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current MachineExtensionsListResponse page.
func (p *MachineExtensionsListPager) PageResponse() MachineExtensionsListResponse {
	return p.current
}

// MachinesListByResourceGroupPager provides operations for iterating over paged responses.
type MachinesListByResourceGroupPager struct {
	client    *MachinesClient
	current   MachinesListByResourceGroupResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, MachinesListByResourceGroupResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *MachinesListByResourceGroupPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *MachinesListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.MachineListResult.NextLink == nil || len(*p.current.MachineListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.pl.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listByResourceGroupHandleError(resp)
		return false
	}
	result, err := p.client.listByResourceGroupHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current MachinesListByResourceGroupResponse page.
func (p *MachinesListByResourceGroupPager) PageResponse() MachinesListByResourceGroupResponse {
	return p.current
}

// MachinesListBySubscriptionPager provides operations for iterating over paged responses.
type MachinesListBySubscriptionPager struct {
	client    *MachinesClient
	current   MachinesListBySubscriptionResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, MachinesListBySubscriptionResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *MachinesListBySubscriptionPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *MachinesListBySubscriptionPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.MachineListResult.NextLink == nil || len(*p.current.MachineListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.pl.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listBySubscriptionHandleError(resp)
		return false
	}
	result, err := p.client.listBySubscriptionHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current MachinesListBySubscriptionResponse page.
func (p *MachinesListBySubscriptionPager) PageResponse() MachinesListBySubscriptionResponse {
	return p.current
}

// PrivateEndpointConnectionsListByPrivateLinkScopePager provides operations for iterating over paged responses.
type PrivateEndpointConnectionsListByPrivateLinkScopePager struct {
	client    *PrivateEndpointConnectionsClient
	current   PrivateEndpointConnectionsListByPrivateLinkScopeResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateEndpointConnectionsListByPrivateLinkScopeResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateEndpointConnectionsListByPrivateLinkScopePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateEndpointConnectionsListByPrivateLinkScopePager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateEndpointConnectionListResult.NextLink == nil || len(*p.current.PrivateEndpointConnectionListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.pl.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listByPrivateLinkScopeHandleError(resp)
		return false
	}
	result, err := p.client.listByPrivateLinkScopeHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateEndpointConnectionsListByPrivateLinkScopeResponse page.
func (p *PrivateEndpointConnectionsListByPrivateLinkScopePager) PageResponse() PrivateEndpointConnectionsListByPrivateLinkScopeResponse {
	return p.current
}

// PrivateLinkResourcesListByPrivateLinkScopePager provides operations for iterating over paged responses.
type PrivateLinkResourcesListByPrivateLinkScopePager struct {
	client    *PrivateLinkResourcesClient
	current   PrivateLinkResourcesListByPrivateLinkScopeResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkResourcesListByPrivateLinkScopeResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkResourcesListByPrivateLinkScopePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkResourcesListByPrivateLinkScopePager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkResourceListResult.NextLink == nil || len(*p.current.PrivateLinkResourceListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.pl.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listByPrivateLinkScopeHandleError(resp)
		return false
	}
	result, err := p.client.listByPrivateLinkScopeHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateLinkResourcesListByPrivateLinkScopeResponse page.
func (p *PrivateLinkResourcesListByPrivateLinkScopePager) PageResponse() PrivateLinkResourcesListByPrivateLinkScopeResponse {
	return p.current
}

// PrivateLinkScopesListByResourceGroupPager provides operations for iterating over paged responses.
type PrivateLinkScopesListByResourceGroupPager struct {
	client    *PrivateLinkScopesClient
	current   PrivateLinkScopesListByResourceGroupResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkScopesListByResourceGroupResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkScopesListByResourceGroupPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkScopesListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.HybridComputePrivateLinkScopeListResult.NextLink == nil || len(*p.current.HybridComputePrivateLinkScopeListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.pl.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listByResourceGroupHandleError(resp)
		return false
	}
	result, err := p.client.listByResourceGroupHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateLinkScopesListByResourceGroupResponse page.
func (p *PrivateLinkScopesListByResourceGroupPager) PageResponse() PrivateLinkScopesListByResourceGroupResponse {
	return p.current
}

// PrivateLinkScopesListPager provides operations for iterating over paged responses.
type PrivateLinkScopesListPager struct {
	client    *PrivateLinkScopesClient
	current   PrivateLinkScopesListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkScopesListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkScopesListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkScopesListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.HybridComputePrivateLinkScopeListResult.NextLink == nil || len(*p.current.HybridComputePrivateLinkScopeListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.pl.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listHandleError(resp)
		return false
	}
	result, err := p.client.listHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateLinkScopesListResponse page.
func (p *PrivateLinkScopesListPager) PageResponse() PrivateLinkScopesListResponse {
	return p.current
}

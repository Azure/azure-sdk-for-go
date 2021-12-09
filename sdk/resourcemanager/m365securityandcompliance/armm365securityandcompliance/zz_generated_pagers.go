//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armm365securityandcompliance

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"reflect"
)

// OperationsListPager provides operations for iterating over paged responses.
type OperationsListPager struct {
	client    *OperationsClient
	current   OperationsListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, OperationsListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *OperationsListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *OperationsListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.OperationListResult.NextLink == nil || len(*p.current.OperationListResult.NextLink) == 0 {
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

// PageResponse returns the current OperationsListResponse page.
func (p *OperationsListPager) PageResponse() OperationsListResponse {
	return p.current
}

// PrivateEndpointConnectionsAdtAPIListByServicePager provides operations for iterating over paged responses.
type PrivateEndpointConnectionsAdtAPIListByServicePager struct {
	client    *PrivateEndpointConnectionsAdtAPIClient
	current   PrivateEndpointConnectionsAdtAPIListByServiceResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateEndpointConnectionsAdtAPIListByServiceResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateEndpointConnectionsAdtAPIListByServicePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateEndpointConnectionsAdtAPIListByServicePager) NextPage(ctx context.Context) bool {
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
		p.err = p.client.listByServiceHandleError(resp)
		return false
	}
	result, err := p.client.listByServiceHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateEndpointConnectionsAdtAPIListByServiceResponse page.
func (p *PrivateEndpointConnectionsAdtAPIListByServicePager) PageResponse() PrivateEndpointConnectionsAdtAPIListByServiceResponse {
	return p.current
}

// PrivateEndpointConnectionsCompListByServicePager provides operations for iterating over paged responses.
type PrivateEndpointConnectionsCompListByServicePager struct {
	client    *PrivateEndpointConnectionsCompClient
	current   PrivateEndpointConnectionsCompListByServiceResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateEndpointConnectionsCompListByServiceResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateEndpointConnectionsCompListByServicePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateEndpointConnectionsCompListByServicePager) NextPage(ctx context.Context) bool {
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
		p.err = p.client.listByServiceHandleError(resp)
		return false
	}
	result, err := p.client.listByServiceHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateEndpointConnectionsCompListByServiceResponse page.
func (p *PrivateEndpointConnectionsCompListByServicePager) PageResponse() PrivateEndpointConnectionsCompListByServiceResponse {
	return p.current
}

// PrivateEndpointConnectionsForEDMListByServicePager provides operations for iterating over paged responses.
type PrivateEndpointConnectionsForEDMListByServicePager struct {
	client    *PrivateEndpointConnectionsForEDMClient
	current   PrivateEndpointConnectionsForEDMListByServiceResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateEndpointConnectionsForEDMListByServiceResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateEndpointConnectionsForEDMListByServicePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateEndpointConnectionsForEDMListByServicePager) NextPage(ctx context.Context) bool {
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
		p.err = p.client.listByServiceHandleError(resp)
		return false
	}
	result, err := p.client.listByServiceHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateEndpointConnectionsForEDMListByServiceResponse page.
func (p *PrivateEndpointConnectionsForEDMListByServicePager) PageResponse() PrivateEndpointConnectionsForEDMListByServiceResponse {
	return p.current
}

// PrivateEndpointConnectionsForMIPPolicySyncListByServicePager provides operations for iterating over paged responses.
type PrivateEndpointConnectionsForMIPPolicySyncListByServicePager struct {
	client    *PrivateEndpointConnectionsForMIPPolicySyncClient
	current   PrivateEndpointConnectionsForMIPPolicySyncListByServiceResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateEndpointConnectionsForMIPPolicySyncListByServiceResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateEndpointConnectionsForMIPPolicySyncListByServicePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateEndpointConnectionsForMIPPolicySyncListByServicePager) NextPage(ctx context.Context) bool {
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
		p.err = p.client.listByServiceHandleError(resp)
		return false
	}
	result, err := p.client.listByServiceHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateEndpointConnectionsForMIPPolicySyncListByServiceResponse page.
func (p *PrivateEndpointConnectionsForMIPPolicySyncListByServicePager) PageResponse() PrivateEndpointConnectionsForMIPPolicySyncListByServiceResponse {
	return p.current
}

// PrivateEndpointConnectionsForSCCPowershellListByServicePager provides operations for iterating over paged responses.
type PrivateEndpointConnectionsForSCCPowershellListByServicePager struct {
	client    *PrivateEndpointConnectionsForSCCPowershellClient
	current   PrivateEndpointConnectionsForSCCPowershellListByServiceResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateEndpointConnectionsForSCCPowershellListByServiceResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateEndpointConnectionsForSCCPowershellListByServicePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateEndpointConnectionsForSCCPowershellListByServicePager) NextPage(ctx context.Context) bool {
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
		p.err = p.client.listByServiceHandleError(resp)
		return false
	}
	result, err := p.client.listByServiceHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateEndpointConnectionsForSCCPowershellListByServiceResponse page.
func (p *PrivateEndpointConnectionsForSCCPowershellListByServicePager) PageResponse() PrivateEndpointConnectionsForSCCPowershellListByServiceResponse {
	return p.current
}

// PrivateEndpointConnectionsSecListByServicePager provides operations for iterating over paged responses.
type PrivateEndpointConnectionsSecListByServicePager struct {
	client    *PrivateEndpointConnectionsSecClient
	current   PrivateEndpointConnectionsSecListByServiceResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateEndpointConnectionsSecListByServiceResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateEndpointConnectionsSecListByServicePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateEndpointConnectionsSecListByServicePager) NextPage(ctx context.Context) bool {
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
		p.err = p.client.listByServiceHandleError(resp)
		return false
	}
	result, err := p.client.listByServiceHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current PrivateEndpointConnectionsSecListByServiceResponse page.
func (p *PrivateEndpointConnectionsSecListByServicePager) PageResponse() PrivateEndpointConnectionsSecListByServiceResponse {
	return p.current
}

// PrivateLinkServicesForEDMUploadListByResourceGroupPager provides operations for iterating over paged responses.
type PrivateLinkServicesForEDMUploadListByResourceGroupPager struct {
	client    *PrivateLinkServicesForEDMUploadClient
	current   PrivateLinkServicesForEDMUploadListByResourceGroupResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForEDMUploadListByResourceGroupResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForEDMUploadListByResourceGroupPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForEDMUploadListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForEDMUploadDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForEDMUploadDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForEDMUploadListByResourceGroupResponse page.
func (p *PrivateLinkServicesForEDMUploadListByResourceGroupPager) PageResponse() PrivateLinkServicesForEDMUploadListByResourceGroupResponse {
	return p.current
}

// PrivateLinkServicesForEDMUploadListPager provides operations for iterating over paged responses.
type PrivateLinkServicesForEDMUploadListPager struct {
	client    *PrivateLinkServicesForEDMUploadClient
	current   PrivateLinkServicesForEDMUploadListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForEDMUploadListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForEDMUploadListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForEDMUploadListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForEDMUploadDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForEDMUploadDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForEDMUploadListResponse page.
func (p *PrivateLinkServicesForEDMUploadListPager) PageResponse() PrivateLinkServicesForEDMUploadListResponse {
	return p.current
}

// PrivateLinkServicesForM365ComplianceCenterListByResourceGroupPager provides operations for iterating over paged responses.
type PrivateLinkServicesForM365ComplianceCenterListByResourceGroupPager struct {
	client    *PrivateLinkServicesForM365ComplianceCenterClient
	current   PrivateLinkServicesForM365ComplianceCenterListByResourceGroupResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForM365ComplianceCenterListByResourceGroupResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForM365ComplianceCenterListByResourceGroupPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForM365ComplianceCenterListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForM365ComplianceCenterDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForM365ComplianceCenterDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForM365ComplianceCenterListByResourceGroupResponse page.
func (p *PrivateLinkServicesForM365ComplianceCenterListByResourceGroupPager) PageResponse() PrivateLinkServicesForM365ComplianceCenterListByResourceGroupResponse {
	return p.current
}

// PrivateLinkServicesForM365ComplianceCenterListPager provides operations for iterating over paged responses.
type PrivateLinkServicesForM365ComplianceCenterListPager struct {
	client    *PrivateLinkServicesForM365ComplianceCenterClient
	current   PrivateLinkServicesForM365ComplianceCenterListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForM365ComplianceCenterListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForM365ComplianceCenterListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForM365ComplianceCenterListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForM365ComplianceCenterDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForM365ComplianceCenterDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForM365ComplianceCenterListResponse page.
func (p *PrivateLinkServicesForM365ComplianceCenterListPager) PageResponse() PrivateLinkServicesForM365ComplianceCenterListResponse {
	return p.current
}

// PrivateLinkServicesForM365SecurityCenterListByResourceGroupPager provides operations for iterating over paged responses.
type PrivateLinkServicesForM365SecurityCenterListByResourceGroupPager struct {
	client    *PrivateLinkServicesForM365SecurityCenterClient
	current   PrivateLinkServicesForM365SecurityCenterListByResourceGroupResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForM365SecurityCenterListByResourceGroupResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForM365SecurityCenterListByResourceGroupPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForM365SecurityCenterListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForM365SecurityCenterDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForM365SecurityCenterDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForM365SecurityCenterListByResourceGroupResponse page.
func (p *PrivateLinkServicesForM365SecurityCenterListByResourceGroupPager) PageResponse() PrivateLinkServicesForM365SecurityCenterListByResourceGroupResponse {
	return p.current
}

// PrivateLinkServicesForM365SecurityCenterListPager provides operations for iterating over paged responses.
type PrivateLinkServicesForM365SecurityCenterListPager struct {
	client    *PrivateLinkServicesForM365SecurityCenterClient
	current   PrivateLinkServicesForM365SecurityCenterListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForM365SecurityCenterListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForM365SecurityCenterListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForM365SecurityCenterListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForM365SecurityCenterDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForM365SecurityCenterDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForM365SecurityCenterListResponse page.
func (p *PrivateLinkServicesForM365SecurityCenterListPager) PageResponse() PrivateLinkServicesForM365SecurityCenterListResponse {
	return p.current
}

// PrivateLinkServicesForMIPPolicySyncListByResourceGroupPager provides operations for iterating over paged responses.
type PrivateLinkServicesForMIPPolicySyncListByResourceGroupPager struct {
	client    *PrivateLinkServicesForMIPPolicySyncClient
	current   PrivateLinkServicesForMIPPolicySyncListByResourceGroupResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForMIPPolicySyncListByResourceGroupResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForMIPPolicySyncListByResourceGroupPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForMIPPolicySyncListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForMIPPolicySyncDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForMIPPolicySyncDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForMIPPolicySyncListByResourceGroupResponse page.
func (p *PrivateLinkServicesForMIPPolicySyncListByResourceGroupPager) PageResponse() PrivateLinkServicesForMIPPolicySyncListByResourceGroupResponse {
	return p.current
}

// PrivateLinkServicesForMIPPolicySyncListPager provides operations for iterating over paged responses.
type PrivateLinkServicesForMIPPolicySyncListPager struct {
	client    *PrivateLinkServicesForMIPPolicySyncClient
	current   PrivateLinkServicesForMIPPolicySyncListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForMIPPolicySyncListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForMIPPolicySyncListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForMIPPolicySyncListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForMIPPolicySyncDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForMIPPolicySyncDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForMIPPolicySyncListResponse page.
func (p *PrivateLinkServicesForMIPPolicySyncListPager) PageResponse() PrivateLinkServicesForMIPPolicySyncListResponse {
	return p.current
}

// PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupPager provides operations for iterating over paged responses.
type PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupPager struct {
	client    *PrivateLinkServicesForO365ManagementActivityAPIClient
	current   PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForO365ManagementActivityAPIDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForO365ManagementActivityAPIDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupResponse page.
func (p *PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupPager) PageResponse() PrivateLinkServicesForO365ManagementActivityAPIListByResourceGroupResponse {
	return p.current
}

// PrivateLinkServicesForO365ManagementActivityAPIListPager provides operations for iterating over paged responses.
type PrivateLinkServicesForO365ManagementActivityAPIListPager struct {
	client    *PrivateLinkServicesForO365ManagementActivityAPIClient
	current   PrivateLinkServicesForO365ManagementActivityAPIListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForO365ManagementActivityAPIListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForO365ManagementActivityAPIListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForO365ManagementActivityAPIListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForO365ManagementActivityAPIDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForO365ManagementActivityAPIDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForO365ManagementActivityAPIListResponse page.
func (p *PrivateLinkServicesForO365ManagementActivityAPIListPager) PageResponse() PrivateLinkServicesForO365ManagementActivityAPIListResponse {
	return p.current
}

// PrivateLinkServicesForSCCPowershellListByResourceGroupPager provides operations for iterating over paged responses.
type PrivateLinkServicesForSCCPowershellListByResourceGroupPager struct {
	client    *PrivateLinkServicesForSCCPowershellClient
	current   PrivateLinkServicesForSCCPowershellListByResourceGroupResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForSCCPowershellListByResourceGroupResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForSCCPowershellListByResourceGroupPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForSCCPowershellListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForSCCPowershellDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForSCCPowershellDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForSCCPowershellListByResourceGroupResponse page.
func (p *PrivateLinkServicesForSCCPowershellListByResourceGroupPager) PageResponse() PrivateLinkServicesForSCCPowershellListByResourceGroupResponse {
	return p.current
}

// PrivateLinkServicesForSCCPowershellListPager provides operations for iterating over paged responses.
type PrivateLinkServicesForSCCPowershellListPager struct {
	client    *PrivateLinkServicesForSCCPowershellClient
	current   PrivateLinkServicesForSCCPowershellListResponse
	err       error
	requester func(context.Context) (*policy.Request, error)
	advancer  func(context.Context, PrivateLinkServicesForSCCPowershellListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *PrivateLinkServicesForSCCPowershellListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *PrivateLinkServicesForSCCPowershellListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PrivateLinkServicesForSCCPowershellDescriptionListResult.NextLink == nil || len(*p.current.PrivateLinkServicesForSCCPowershellDescriptionListResult.NextLink) == 0 {
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

// PageResponse returns the current PrivateLinkServicesForSCCPowershellListResponse page.
func (p *PrivateLinkServicesForSCCPowershellListPager) PageResponse() PrivateLinkServicesForSCCPowershellListResponse {
	return p.current
}

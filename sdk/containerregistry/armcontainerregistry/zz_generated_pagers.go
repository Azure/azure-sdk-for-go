// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcontainerregistry

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// AgentPoolListResultPager provides iteration over AgentPoolListResult pages.
type AgentPoolListResultPager interface {
	azcore.Pager

	// PageResponse returns the current AgentPoolListResultResponse.
	PageResponse() AgentPoolListResultResponse
}

type agentPoolListResultCreateRequest func(context.Context) (*azcore.Request, error)

type agentPoolListResultHandleError func(*azcore.Response) error

type agentPoolListResultHandleResponse func(*azcore.Response) (AgentPoolListResultResponse, error)

type agentPoolListResultAdvancePage func(context.Context, AgentPoolListResultResponse) (*azcore.Request, error)

type agentPoolListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester agentPoolListResultCreateRequest
	// callback for handling response errors
	errorer agentPoolListResultHandleError
	// callback for handling the HTTP response
	responder agentPoolListResultHandleResponse
	// callback for advancing to the next page
	advancer agentPoolListResultAdvancePage
	// contains the current response
	current AgentPoolListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *agentPoolListResultPager) Err() error {
	return p.err
}

func (p *agentPoolListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.AgentPoolListResult.NextLink == nil || len(*p.current.AgentPoolListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *agentPoolListResultPager) PageResponse() AgentPoolListResultResponse {
	return p.current
}

// ConnectedRegistryListResultPager provides iteration over ConnectedRegistryListResult pages.
type ConnectedRegistryListResultPager interface {
	azcore.Pager

	// PageResponse returns the current ConnectedRegistryListResultResponse.
	PageResponse() ConnectedRegistryListResultResponse
}

type connectedRegistryListResultCreateRequest func(context.Context) (*azcore.Request, error)

type connectedRegistryListResultHandleError func(*azcore.Response) error

type connectedRegistryListResultHandleResponse func(*azcore.Response) (ConnectedRegistryListResultResponse, error)

type connectedRegistryListResultAdvancePage func(context.Context, ConnectedRegistryListResultResponse) (*azcore.Request, error)

type connectedRegistryListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester connectedRegistryListResultCreateRequest
	// callback for handling response errors
	errorer connectedRegistryListResultHandleError
	// callback for handling the HTTP response
	responder connectedRegistryListResultHandleResponse
	// callback for advancing to the next page
	advancer connectedRegistryListResultAdvancePage
	// contains the current response
	current ConnectedRegistryListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *connectedRegistryListResultPager) Err() error {
	return p.err
}

func (p *connectedRegistryListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.ConnectedRegistryListResult.NextLink == nil || len(*p.current.ConnectedRegistryListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *connectedRegistryListResultPager) PageResponse() ConnectedRegistryListResultResponse {
	return p.current
}

// EventListResultPager provides iteration over EventListResult pages.
type EventListResultPager interface {
	azcore.Pager

	// PageResponse returns the current EventListResultResponse.
	PageResponse() EventListResultResponse
}

type eventListResultCreateRequest func(context.Context) (*azcore.Request, error)

type eventListResultHandleError func(*azcore.Response) error

type eventListResultHandleResponse func(*azcore.Response) (EventListResultResponse, error)

type eventListResultAdvancePage func(context.Context, EventListResultResponse) (*azcore.Request, error)

type eventListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester eventListResultCreateRequest
	// callback for handling response errors
	errorer eventListResultHandleError
	// callback for handling the HTTP response
	responder eventListResultHandleResponse
	// callback for advancing to the next page
	advancer eventListResultAdvancePage
	// contains the current response
	current EventListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *eventListResultPager) Err() error {
	return p.err
}

func (p *eventListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.EventListResult.NextLink == nil || len(*p.current.EventListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *eventListResultPager) PageResponse() EventListResultResponse {
	return p.current
}

// ExportPipelineListResultPager provides iteration over ExportPipelineListResult pages.
type ExportPipelineListResultPager interface {
	azcore.Pager

	// PageResponse returns the current ExportPipelineListResultResponse.
	PageResponse() ExportPipelineListResultResponse
}

type exportPipelineListResultCreateRequest func(context.Context) (*azcore.Request, error)

type exportPipelineListResultHandleError func(*azcore.Response) error

type exportPipelineListResultHandleResponse func(*azcore.Response) (ExportPipelineListResultResponse, error)

type exportPipelineListResultAdvancePage func(context.Context, ExportPipelineListResultResponse) (*azcore.Request, error)

type exportPipelineListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester exportPipelineListResultCreateRequest
	// callback for handling response errors
	errorer exportPipelineListResultHandleError
	// callback for handling the HTTP response
	responder exportPipelineListResultHandleResponse
	// callback for advancing to the next page
	advancer exportPipelineListResultAdvancePage
	// contains the current response
	current ExportPipelineListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *exportPipelineListResultPager) Err() error {
	return p.err
}

func (p *exportPipelineListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.ExportPipelineListResult.NextLink == nil || len(*p.current.ExportPipelineListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *exportPipelineListResultPager) PageResponse() ExportPipelineListResultResponse {
	return p.current
}

// ImportPipelineListResultPager provides iteration over ImportPipelineListResult pages.
type ImportPipelineListResultPager interface {
	azcore.Pager

	// PageResponse returns the current ImportPipelineListResultResponse.
	PageResponse() ImportPipelineListResultResponse
}

type importPipelineListResultCreateRequest func(context.Context) (*azcore.Request, error)

type importPipelineListResultHandleError func(*azcore.Response) error

type importPipelineListResultHandleResponse func(*azcore.Response) (ImportPipelineListResultResponse, error)

type importPipelineListResultAdvancePage func(context.Context, ImportPipelineListResultResponse) (*azcore.Request, error)

type importPipelineListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester importPipelineListResultCreateRequest
	// callback for handling response errors
	errorer importPipelineListResultHandleError
	// callback for handling the HTTP response
	responder importPipelineListResultHandleResponse
	// callback for advancing to the next page
	advancer importPipelineListResultAdvancePage
	// contains the current response
	current ImportPipelineListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *importPipelineListResultPager) Err() error {
	return p.err
}

func (p *importPipelineListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.ImportPipelineListResult.NextLink == nil || len(*p.current.ImportPipelineListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *importPipelineListResultPager) PageResponse() ImportPipelineListResultResponse {
	return p.current
}

// OperationListResultPager provides iteration over OperationListResult pages.
type OperationListResultPager interface {
	azcore.Pager

	// PageResponse returns the current OperationListResultResponse.
	PageResponse() OperationListResultResponse
}

type operationListResultCreateRequest func(context.Context) (*azcore.Request, error)

type operationListResultHandleError func(*azcore.Response) error

type operationListResultHandleResponse func(*azcore.Response) (OperationListResultResponse, error)

type operationListResultAdvancePage func(context.Context, OperationListResultResponse) (*azcore.Request, error)

type operationListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester operationListResultCreateRequest
	// callback for handling response errors
	errorer operationListResultHandleError
	// callback for handling the HTTP response
	responder operationListResultHandleResponse
	// callback for advancing to the next page
	advancer operationListResultAdvancePage
	// contains the current response
	current OperationListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *operationListResultPager) Err() error {
	return p.err
}

func (p *operationListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *operationListResultPager) PageResponse() OperationListResultResponse {
	return p.current
}

// PipelineRunListResultPager provides iteration over PipelineRunListResult pages.
type PipelineRunListResultPager interface {
	azcore.Pager

	// PageResponse returns the current PipelineRunListResultResponse.
	PageResponse() PipelineRunListResultResponse
}

type pipelineRunListResultCreateRequest func(context.Context) (*azcore.Request, error)

type pipelineRunListResultHandleError func(*azcore.Response) error

type pipelineRunListResultHandleResponse func(*azcore.Response) (PipelineRunListResultResponse, error)

type pipelineRunListResultAdvancePage func(context.Context, PipelineRunListResultResponse) (*azcore.Request, error)

type pipelineRunListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester pipelineRunListResultCreateRequest
	// callback for handling response errors
	errorer pipelineRunListResultHandleError
	// callback for handling the HTTP response
	responder pipelineRunListResultHandleResponse
	// callback for advancing to the next page
	advancer pipelineRunListResultAdvancePage
	// contains the current response
	current PipelineRunListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *pipelineRunListResultPager) Err() error {
	return p.err
}

func (p *pipelineRunListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.PipelineRunListResult.NextLink == nil || len(*p.current.PipelineRunListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *pipelineRunListResultPager) PageResponse() PipelineRunListResultResponse {
	return p.current
}

// PrivateEndpointConnectionListResultPager provides iteration over PrivateEndpointConnectionListResult pages.
type PrivateEndpointConnectionListResultPager interface {
	azcore.Pager

	// PageResponse returns the current PrivateEndpointConnectionListResultResponse.
	PageResponse() PrivateEndpointConnectionListResultResponse
}

type privateEndpointConnectionListResultCreateRequest func(context.Context) (*azcore.Request, error)

type privateEndpointConnectionListResultHandleError func(*azcore.Response) error

type privateEndpointConnectionListResultHandleResponse func(*azcore.Response) (PrivateEndpointConnectionListResultResponse, error)

type privateEndpointConnectionListResultAdvancePage func(context.Context, PrivateEndpointConnectionListResultResponse) (*azcore.Request, error)

type privateEndpointConnectionListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester privateEndpointConnectionListResultCreateRequest
	// callback for handling response errors
	errorer privateEndpointConnectionListResultHandleError
	// callback for handling the HTTP response
	responder privateEndpointConnectionListResultHandleResponse
	// callback for advancing to the next page
	advancer privateEndpointConnectionListResultAdvancePage
	// contains the current response
	current PrivateEndpointConnectionListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *privateEndpointConnectionListResultPager) Err() error {
	return p.err
}

func (p *privateEndpointConnectionListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *privateEndpointConnectionListResultPager) PageResponse() PrivateEndpointConnectionListResultResponse {
	return p.current
}

// PrivateLinkResourceListResultPager provides iteration over PrivateLinkResourceListResult pages.
type PrivateLinkResourceListResultPager interface {
	azcore.Pager

	// PageResponse returns the current PrivateLinkResourceListResultResponse.
	PageResponse() PrivateLinkResourceListResultResponse
}

type privateLinkResourceListResultCreateRequest func(context.Context) (*azcore.Request, error)

type privateLinkResourceListResultHandleError func(*azcore.Response) error

type privateLinkResourceListResultHandleResponse func(*azcore.Response) (PrivateLinkResourceListResultResponse, error)

type privateLinkResourceListResultAdvancePage func(context.Context, PrivateLinkResourceListResultResponse) (*azcore.Request, error)

type privateLinkResourceListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester privateLinkResourceListResultCreateRequest
	// callback for handling response errors
	errorer privateLinkResourceListResultHandleError
	// callback for handling the HTTP response
	responder privateLinkResourceListResultHandleResponse
	// callback for advancing to the next page
	advancer privateLinkResourceListResultAdvancePage
	// contains the current response
	current PrivateLinkResourceListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *privateLinkResourceListResultPager) Err() error {
	return p.err
}

func (p *privateLinkResourceListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *privateLinkResourceListResultPager) PageResponse() PrivateLinkResourceListResultResponse {
	return p.current
}

// RegistryListResultPager provides iteration over RegistryListResult pages.
type RegistryListResultPager interface {
	azcore.Pager

	// PageResponse returns the current RegistryListResultResponse.
	PageResponse() RegistryListResultResponse
}

type registryListResultCreateRequest func(context.Context) (*azcore.Request, error)

type registryListResultHandleError func(*azcore.Response) error

type registryListResultHandleResponse func(*azcore.Response) (RegistryListResultResponse, error)

type registryListResultAdvancePage func(context.Context, RegistryListResultResponse) (*azcore.Request, error)

type registryListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester registryListResultCreateRequest
	// callback for handling response errors
	errorer registryListResultHandleError
	// callback for handling the HTTP response
	responder registryListResultHandleResponse
	// callback for advancing to the next page
	advancer registryListResultAdvancePage
	// contains the current response
	current RegistryListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *registryListResultPager) Err() error {
	return p.err
}

func (p *registryListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RegistryListResult.NextLink == nil || len(*p.current.RegistryListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *registryListResultPager) PageResponse() RegistryListResultResponse {
	return p.current
}

// ReplicationListResultPager provides iteration over ReplicationListResult pages.
type ReplicationListResultPager interface {
	azcore.Pager

	// PageResponse returns the current ReplicationListResultResponse.
	PageResponse() ReplicationListResultResponse
}

type replicationListResultCreateRequest func(context.Context) (*azcore.Request, error)

type replicationListResultHandleError func(*azcore.Response) error

type replicationListResultHandleResponse func(*azcore.Response) (ReplicationListResultResponse, error)

type replicationListResultAdvancePage func(context.Context, ReplicationListResultResponse) (*azcore.Request, error)

type replicationListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester replicationListResultCreateRequest
	// callback for handling response errors
	errorer replicationListResultHandleError
	// callback for handling the HTTP response
	responder replicationListResultHandleResponse
	// callback for advancing to the next page
	advancer replicationListResultAdvancePage
	// contains the current response
	current ReplicationListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *replicationListResultPager) Err() error {
	return p.err
}

func (p *replicationListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.ReplicationListResult.NextLink == nil || len(*p.current.ReplicationListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *replicationListResultPager) PageResponse() ReplicationListResultResponse {
	return p.current
}

// RunListResultPager provides iteration over RunListResult pages.
type RunListResultPager interface {
	azcore.Pager

	// PageResponse returns the current RunListResultResponse.
	PageResponse() RunListResultResponse
}

type runListResultCreateRequest func(context.Context) (*azcore.Request, error)

type runListResultHandleError func(*azcore.Response) error

type runListResultHandleResponse func(*azcore.Response) (RunListResultResponse, error)

type runListResultAdvancePage func(context.Context, RunListResultResponse) (*azcore.Request, error)

type runListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester runListResultCreateRequest
	// callback for handling response errors
	errorer runListResultHandleError
	// callback for handling the HTTP response
	responder runListResultHandleResponse
	// callback for advancing to the next page
	advancer runListResultAdvancePage
	// contains the current response
	current RunListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *runListResultPager) Err() error {
	return p.err
}

func (p *runListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RunListResult.NextLink == nil || len(*p.current.RunListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *runListResultPager) PageResponse() RunListResultResponse {
	return p.current
}

// ScopeMapListResultPager provides iteration over ScopeMapListResult pages.
type ScopeMapListResultPager interface {
	azcore.Pager

	// PageResponse returns the current ScopeMapListResultResponse.
	PageResponse() ScopeMapListResultResponse
}

type scopeMapListResultCreateRequest func(context.Context) (*azcore.Request, error)

type scopeMapListResultHandleError func(*azcore.Response) error

type scopeMapListResultHandleResponse func(*azcore.Response) (ScopeMapListResultResponse, error)

type scopeMapListResultAdvancePage func(context.Context, ScopeMapListResultResponse) (*azcore.Request, error)

type scopeMapListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester scopeMapListResultCreateRequest
	// callback for handling response errors
	errorer scopeMapListResultHandleError
	// callback for handling the HTTP response
	responder scopeMapListResultHandleResponse
	// callback for advancing to the next page
	advancer scopeMapListResultAdvancePage
	// contains the current response
	current ScopeMapListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *scopeMapListResultPager) Err() error {
	return p.err
}

func (p *scopeMapListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.ScopeMapListResult.NextLink == nil || len(*p.current.ScopeMapListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *scopeMapListResultPager) PageResponse() ScopeMapListResultResponse {
	return p.current
}

// TaskListResultPager provides iteration over TaskListResult pages.
type TaskListResultPager interface {
	azcore.Pager

	// PageResponse returns the current TaskListResultResponse.
	PageResponse() TaskListResultResponse
}

type taskListResultCreateRequest func(context.Context) (*azcore.Request, error)

type taskListResultHandleError func(*azcore.Response) error

type taskListResultHandleResponse func(*azcore.Response) (TaskListResultResponse, error)

type taskListResultAdvancePage func(context.Context, TaskListResultResponse) (*azcore.Request, error)

type taskListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester taskListResultCreateRequest
	// callback for handling response errors
	errorer taskListResultHandleError
	// callback for handling the HTTP response
	responder taskListResultHandleResponse
	// callback for advancing to the next page
	advancer taskListResultAdvancePage
	// contains the current response
	current TaskListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *taskListResultPager) Err() error {
	return p.err
}

func (p *taskListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.TaskListResult.NextLink == nil || len(*p.current.TaskListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *taskListResultPager) PageResponse() TaskListResultResponse {
	return p.current
}

// TaskRunListResultPager provides iteration over TaskRunListResult pages.
type TaskRunListResultPager interface {
	azcore.Pager

	// PageResponse returns the current TaskRunListResultResponse.
	PageResponse() TaskRunListResultResponse
}

type taskRunListResultCreateRequest func(context.Context) (*azcore.Request, error)

type taskRunListResultHandleError func(*azcore.Response) error

type taskRunListResultHandleResponse func(*azcore.Response) (TaskRunListResultResponse, error)

type taskRunListResultAdvancePage func(context.Context, TaskRunListResultResponse) (*azcore.Request, error)

type taskRunListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester taskRunListResultCreateRequest
	// callback for handling response errors
	errorer taskRunListResultHandleError
	// callback for handling the HTTP response
	responder taskRunListResultHandleResponse
	// callback for advancing to the next page
	advancer taskRunListResultAdvancePage
	// contains the current response
	current TaskRunListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *taskRunListResultPager) Err() error {
	return p.err
}

func (p *taskRunListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.TaskRunListResult.NextLink == nil || len(*p.current.TaskRunListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *taskRunListResultPager) PageResponse() TaskRunListResultResponse {
	return p.current
}

// TokenListResultPager provides iteration over TokenListResult pages.
type TokenListResultPager interface {
	azcore.Pager

	// PageResponse returns the current TokenListResultResponse.
	PageResponse() TokenListResultResponse
}

type tokenListResultCreateRequest func(context.Context) (*azcore.Request, error)

type tokenListResultHandleError func(*azcore.Response) error

type tokenListResultHandleResponse func(*azcore.Response) (TokenListResultResponse, error)

type tokenListResultAdvancePage func(context.Context, TokenListResultResponse) (*azcore.Request, error)

type tokenListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester tokenListResultCreateRequest
	// callback for handling response errors
	errorer tokenListResultHandleError
	// callback for handling the HTTP response
	responder tokenListResultHandleResponse
	// callback for advancing to the next page
	advancer tokenListResultAdvancePage
	// contains the current response
	current TokenListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *tokenListResultPager) Err() error {
	return p.err
}

func (p *tokenListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.TokenListResult.NextLink == nil || len(*p.current.TokenListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *tokenListResultPager) PageResponse() TokenListResultResponse {
	return p.current
}

// WebhookListResultPager provides iteration over WebhookListResult pages.
type WebhookListResultPager interface {
	azcore.Pager

	// PageResponse returns the current WebhookListResultResponse.
	PageResponse() WebhookListResultResponse
}

type webhookListResultCreateRequest func(context.Context) (*azcore.Request, error)

type webhookListResultHandleError func(*azcore.Response) error

type webhookListResultHandleResponse func(*azcore.Response) (WebhookListResultResponse, error)

type webhookListResultAdvancePage func(context.Context, WebhookListResultResponse) (*azcore.Request, error)

type webhookListResultPager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// creates the initial request (non-LRO case)
	requester webhookListResultCreateRequest
	// callback for handling response errors
	errorer webhookListResultHandleError
	// callback for handling the HTTP response
	responder webhookListResultHandleResponse
	// callback for advancing to the next page
	advancer webhookListResultAdvancePage
	// contains the current response
	current WebhookListResultResponse
	// status codes for successful retrieval
	statusCodes []int
	// any error encountered
	err error
}

func (p *webhookListResultPager) Err() error {
	return p.err
}

func (p *webhookListResultPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.WebhookListResult.NextLink == nil || len(*p.current.WebhookListResult.NextLink) == 0 {
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
	resp, err := p.pipeline.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(p.statusCodes...) {
		p.err = p.errorer(resp)
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *webhookListResultPager) PageResponse() WebhookListResultResponse {
	return p.current
}

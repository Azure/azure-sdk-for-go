//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armedgeorder

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
)

// EdgeOrderManagementClientCreateAddressPoller provides polling facilities until the operation reaches a terminal state.
type EdgeOrderManagementClientCreateAddressPoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *EdgeOrderManagementClientCreateAddressPoller) Done() bool {
	return p.pt.Done()
}

// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP
// response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is
// updated and the error is returned.
// If the LRO has not reached a terminal state, the poller's state is updated and
// the latest HTTP response is returned.
// If Poll fails, the poller's state is unmodified and the error is returned.
// Calling Poll on an LRO that has reached a terminal state will return the final
// HTTP response or error.
func (p *EdgeOrderManagementClientCreateAddressPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final EdgeOrderManagementClientCreateAddressResponse will be returned.
func (p *EdgeOrderManagementClientCreateAddressPoller) FinalResponse(ctx context.Context) (EdgeOrderManagementClientCreateAddressResponse, error) {
	respType := EdgeOrderManagementClientCreateAddressResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.AddressResource)
	if err != nil {
		return EdgeOrderManagementClientCreateAddressResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *EdgeOrderManagementClientCreateAddressPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// EdgeOrderManagementClientCreateOrderItemPoller provides polling facilities until the operation reaches a terminal state.
type EdgeOrderManagementClientCreateOrderItemPoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *EdgeOrderManagementClientCreateOrderItemPoller) Done() bool {
	return p.pt.Done()
}

// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP
// response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is
// updated and the error is returned.
// If the LRO has not reached a terminal state, the poller's state is updated and
// the latest HTTP response is returned.
// If Poll fails, the poller's state is unmodified and the error is returned.
// Calling Poll on an LRO that has reached a terminal state will return the final
// HTTP response or error.
func (p *EdgeOrderManagementClientCreateOrderItemPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final EdgeOrderManagementClientCreateOrderItemResponse will be returned.
func (p *EdgeOrderManagementClientCreateOrderItemPoller) FinalResponse(ctx context.Context) (EdgeOrderManagementClientCreateOrderItemResponse, error) {
	respType := EdgeOrderManagementClientCreateOrderItemResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.OrderItemResource)
	if err != nil {
		return EdgeOrderManagementClientCreateOrderItemResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *EdgeOrderManagementClientCreateOrderItemPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// EdgeOrderManagementClientDeleteAddressByNamePoller provides polling facilities until the operation reaches a terminal state.
type EdgeOrderManagementClientDeleteAddressByNamePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *EdgeOrderManagementClientDeleteAddressByNamePoller) Done() bool {
	return p.pt.Done()
}

// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP
// response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is
// updated and the error is returned.
// If the LRO has not reached a terminal state, the poller's state is updated and
// the latest HTTP response is returned.
// If Poll fails, the poller's state is unmodified and the error is returned.
// Calling Poll on an LRO that has reached a terminal state will return the final
// HTTP response or error.
func (p *EdgeOrderManagementClientDeleteAddressByNamePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final EdgeOrderManagementClientDeleteAddressByNameResponse will be returned.
func (p *EdgeOrderManagementClientDeleteAddressByNamePoller) FinalResponse(ctx context.Context) (EdgeOrderManagementClientDeleteAddressByNameResponse, error) {
	respType := EdgeOrderManagementClientDeleteAddressByNameResponse{}
	resp, err := p.pt.FinalResponse(ctx, nil)
	if err != nil {
		return EdgeOrderManagementClientDeleteAddressByNameResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *EdgeOrderManagementClientDeleteAddressByNamePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// EdgeOrderManagementClientDeleteOrderItemByNamePoller provides polling facilities until the operation reaches a terminal state.
type EdgeOrderManagementClientDeleteOrderItemByNamePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *EdgeOrderManagementClientDeleteOrderItemByNamePoller) Done() bool {
	return p.pt.Done()
}

// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP
// response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is
// updated and the error is returned.
// If the LRO has not reached a terminal state, the poller's state is updated and
// the latest HTTP response is returned.
// If Poll fails, the poller's state is unmodified and the error is returned.
// Calling Poll on an LRO that has reached a terminal state will return the final
// HTTP response or error.
func (p *EdgeOrderManagementClientDeleteOrderItemByNamePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final EdgeOrderManagementClientDeleteOrderItemByNameResponse will be returned.
func (p *EdgeOrderManagementClientDeleteOrderItemByNamePoller) FinalResponse(ctx context.Context) (EdgeOrderManagementClientDeleteOrderItemByNameResponse, error) {
	respType := EdgeOrderManagementClientDeleteOrderItemByNameResponse{}
	resp, err := p.pt.FinalResponse(ctx, nil)
	if err != nil {
		return EdgeOrderManagementClientDeleteOrderItemByNameResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *EdgeOrderManagementClientDeleteOrderItemByNamePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// EdgeOrderManagementClientReturnOrderItemPoller provides polling facilities until the operation reaches a terminal state.
type EdgeOrderManagementClientReturnOrderItemPoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *EdgeOrderManagementClientReturnOrderItemPoller) Done() bool {
	return p.pt.Done()
}

// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP
// response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is
// updated and the error is returned.
// If the LRO has not reached a terminal state, the poller's state is updated and
// the latest HTTP response is returned.
// If Poll fails, the poller's state is unmodified and the error is returned.
// Calling Poll on an LRO that has reached a terminal state will return the final
// HTTP response or error.
func (p *EdgeOrderManagementClientReturnOrderItemPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final EdgeOrderManagementClientReturnOrderItemResponse will be returned.
func (p *EdgeOrderManagementClientReturnOrderItemPoller) FinalResponse(ctx context.Context) (EdgeOrderManagementClientReturnOrderItemResponse, error) {
	respType := EdgeOrderManagementClientReturnOrderItemResponse{}
	resp, err := p.pt.FinalResponse(ctx, nil)
	if err != nil {
		return EdgeOrderManagementClientReturnOrderItemResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *EdgeOrderManagementClientReturnOrderItemPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// EdgeOrderManagementClientUpdateAddressPoller provides polling facilities until the operation reaches a terminal state.
type EdgeOrderManagementClientUpdateAddressPoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *EdgeOrderManagementClientUpdateAddressPoller) Done() bool {
	return p.pt.Done()
}

// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP
// response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is
// updated and the error is returned.
// If the LRO has not reached a terminal state, the poller's state is updated and
// the latest HTTP response is returned.
// If Poll fails, the poller's state is unmodified and the error is returned.
// Calling Poll on an LRO that has reached a terminal state will return the final
// HTTP response or error.
func (p *EdgeOrderManagementClientUpdateAddressPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final EdgeOrderManagementClientUpdateAddressResponse will be returned.
func (p *EdgeOrderManagementClientUpdateAddressPoller) FinalResponse(ctx context.Context) (EdgeOrderManagementClientUpdateAddressResponse, error) {
	respType := EdgeOrderManagementClientUpdateAddressResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.AddressResource)
	if err != nil {
		return EdgeOrderManagementClientUpdateAddressResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *EdgeOrderManagementClientUpdateAddressPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// EdgeOrderManagementClientUpdateOrderItemPoller provides polling facilities until the operation reaches a terminal state.
type EdgeOrderManagementClientUpdateOrderItemPoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *EdgeOrderManagementClientUpdateOrderItemPoller) Done() bool {
	return p.pt.Done()
}

// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP
// response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is
// updated and the error is returned.
// If the LRO has not reached a terminal state, the poller's state is updated and
// the latest HTTP response is returned.
// If Poll fails, the poller's state is unmodified and the error is returned.
// Calling Poll on an LRO that has reached a terminal state will return the final
// HTTP response or error.
func (p *EdgeOrderManagementClientUpdateOrderItemPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final EdgeOrderManagementClientUpdateOrderItemResponse will be returned.
func (p *EdgeOrderManagementClientUpdateOrderItemPoller) FinalResponse(ctx context.Context) (EdgeOrderManagementClientUpdateOrderItemResponse, error) {
	respType := EdgeOrderManagementClientUpdateOrderItemResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.OrderItemResource)
	if err != nil {
		return EdgeOrderManagementClientUpdateOrderItemResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *EdgeOrderManagementClientUpdateOrderItemPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

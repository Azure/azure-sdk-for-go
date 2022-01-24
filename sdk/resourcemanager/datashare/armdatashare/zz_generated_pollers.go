//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatashare

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
)

// AccountsClientCreatePoller provides polling facilities until the operation reaches a terminal state.
type AccountsClientCreatePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *AccountsClientCreatePoller) Done() bool {
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
func (p *AccountsClientCreatePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final AccountsClientCreateResponse will be returned.
func (p *AccountsClientCreatePoller) FinalResponse(ctx context.Context) (AccountsClientCreateResponse, error) {
	respType := AccountsClientCreateResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.Account)
	if err != nil {
		return AccountsClientCreateResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *AccountsClientCreatePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// AccountsClientDeletePoller provides polling facilities until the operation reaches a terminal state.
type AccountsClientDeletePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *AccountsClientDeletePoller) Done() bool {
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
func (p *AccountsClientDeletePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final AccountsClientDeleteResponse will be returned.
func (p *AccountsClientDeletePoller) FinalResponse(ctx context.Context) (AccountsClientDeleteResponse, error) {
	respType := AccountsClientDeleteResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.OperationResponse)
	if err != nil {
		return AccountsClientDeleteResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *AccountsClientDeletePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// DataSetsClientDeletePoller provides polling facilities until the operation reaches a terminal state.
type DataSetsClientDeletePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *DataSetsClientDeletePoller) Done() bool {
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
func (p *DataSetsClientDeletePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final DataSetsClientDeleteResponse will be returned.
func (p *DataSetsClientDeletePoller) FinalResponse(ctx context.Context) (DataSetsClientDeleteResponse, error) {
	respType := DataSetsClientDeleteResponse{}
	resp, err := p.pt.FinalResponse(ctx, nil)
	if err != nil {
		return DataSetsClientDeleteResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *DataSetsClientDeletePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// ProviderShareSubscriptionsClientRevokePoller provides polling facilities until the operation reaches a terminal state.
type ProviderShareSubscriptionsClientRevokePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *ProviderShareSubscriptionsClientRevokePoller) Done() bool {
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
func (p *ProviderShareSubscriptionsClientRevokePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final ProviderShareSubscriptionsClientRevokeResponse will be returned.
func (p *ProviderShareSubscriptionsClientRevokePoller) FinalResponse(ctx context.Context) (ProviderShareSubscriptionsClientRevokeResponse, error) {
	respType := ProviderShareSubscriptionsClientRevokeResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.ProviderShareSubscription)
	if err != nil {
		return ProviderShareSubscriptionsClientRevokeResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *ProviderShareSubscriptionsClientRevokePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// ShareSubscriptionsClientCancelSynchronizationPoller provides polling facilities until the operation reaches a terminal state.
type ShareSubscriptionsClientCancelSynchronizationPoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *ShareSubscriptionsClientCancelSynchronizationPoller) Done() bool {
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
func (p *ShareSubscriptionsClientCancelSynchronizationPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final ShareSubscriptionsClientCancelSynchronizationResponse will be returned.
func (p *ShareSubscriptionsClientCancelSynchronizationPoller) FinalResponse(ctx context.Context) (ShareSubscriptionsClientCancelSynchronizationResponse, error) {
	respType := ShareSubscriptionsClientCancelSynchronizationResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.ShareSubscriptionSynchronization)
	if err != nil {
		return ShareSubscriptionsClientCancelSynchronizationResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *ShareSubscriptionsClientCancelSynchronizationPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// ShareSubscriptionsClientDeletePoller provides polling facilities until the operation reaches a terminal state.
type ShareSubscriptionsClientDeletePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *ShareSubscriptionsClientDeletePoller) Done() bool {
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
func (p *ShareSubscriptionsClientDeletePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final ShareSubscriptionsClientDeleteResponse will be returned.
func (p *ShareSubscriptionsClientDeletePoller) FinalResponse(ctx context.Context) (ShareSubscriptionsClientDeleteResponse, error) {
	respType := ShareSubscriptionsClientDeleteResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.OperationResponse)
	if err != nil {
		return ShareSubscriptionsClientDeleteResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *ShareSubscriptionsClientDeletePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// ShareSubscriptionsClientSynchronizePoller provides polling facilities until the operation reaches a terminal state.
type ShareSubscriptionsClientSynchronizePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *ShareSubscriptionsClientSynchronizePoller) Done() bool {
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
func (p *ShareSubscriptionsClientSynchronizePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final ShareSubscriptionsClientSynchronizeResponse will be returned.
func (p *ShareSubscriptionsClientSynchronizePoller) FinalResponse(ctx context.Context) (ShareSubscriptionsClientSynchronizeResponse, error) {
	respType := ShareSubscriptionsClientSynchronizeResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.ShareSubscriptionSynchronization)
	if err != nil {
		return ShareSubscriptionsClientSynchronizeResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *ShareSubscriptionsClientSynchronizePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// SharesClientDeletePoller provides polling facilities until the operation reaches a terminal state.
type SharesClientDeletePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *SharesClientDeletePoller) Done() bool {
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
func (p *SharesClientDeletePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final SharesClientDeleteResponse will be returned.
func (p *SharesClientDeletePoller) FinalResponse(ctx context.Context) (SharesClientDeleteResponse, error) {
	respType := SharesClientDeleteResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.OperationResponse)
	if err != nil {
		return SharesClientDeleteResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *SharesClientDeletePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// SynchronizationSettingsClientDeletePoller provides polling facilities until the operation reaches a terminal state.
type SynchronizationSettingsClientDeletePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *SynchronizationSettingsClientDeletePoller) Done() bool {
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
func (p *SynchronizationSettingsClientDeletePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final SynchronizationSettingsClientDeleteResponse will be returned.
func (p *SynchronizationSettingsClientDeletePoller) FinalResponse(ctx context.Context) (SynchronizationSettingsClientDeleteResponse, error) {
	respType := SynchronizationSettingsClientDeleteResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.OperationResponse)
	if err != nil {
		return SynchronizationSettingsClientDeleteResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *SynchronizationSettingsClientDeletePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// TriggersClientCreatePoller provides polling facilities until the operation reaches a terminal state.
type TriggersClientCreatePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *TriggersClientCreatePoller) Done() bool {
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
func (p *TriggersClientCreatePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final TriggersClientCreateResponse will be returned.
func (p *TriggersClientCreatePoller) FinalResponse(ctx context.Context) (TriggersClientCreateResponse, error) {
	respType := TriggersClientCreateResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.TriggersClientCreateResult)
	if err != nil {
		return TriggersClientCreateResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *TriggersClientCreatePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

// TriggersClientDeletePoller provides polling facilities until the operation reaches a terminal state.
type TriggersClientDeletePoller struct {
	pt *azcore.Poller
}

// Done returns true if the LRO has reached a terminal state.
func (p *TriggersClientDeletePoller) Done() bool {
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
func (p *TriggersClientDeletePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// FinalResponse performs a final GET to the service and returns the final response
// for the polling operation. If there is an error performing the final GET then an error is returned.
// If the final GET succeeded then the final TriggersClientDeleteResponse will be returned.
func (p *TriggersClientDeletePoller) FinalResponse(ctx context.Context) (TriggersClientDeleteResponse, error) {
	respType := TriggersClientDeleteResponse{}
	resp, err := p.pt.FinalResponse(ctx, &respType.OperationResponse)
	if err != nil {
		return TriggersClientDeleteResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *TriggersClientDeletePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhanaonazure

import (
	"context"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"net/http"
	"time"
)

// OperationsListResponse contains the response from method Operations.List.
type OperationsListResponse struct {
	OperationsListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// OperationsListResult contains the result from method Operations.List.
type OperationsListResult struct {
	OperationList
}

// ProviderInstancesCreatePollerResponse contains the response from method ProviderInstances.Create.
type ProviderInstancesCreatePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *ProviderInstancesCreatePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l ProviderInstancesCreatePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (ProviderInstancesCreateResponse, error) {
	respType := ProviderInstancesCreateResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, &respType.ProviderInstance)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a ProviderInstancesCreatePollerResponse from the provided client and resume token.
func (l *ProviderInstancesCreatePollerResponse) Resume(ctx context.Context, client *ProviderInstancesClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("ProviderInstancesClient.Create", token, client.pl, client.createHandleError)
	if err != nil {
		return err
	}
	poller := &ProviderInstancesCreatePoller{
		pt: pt,
	}
	resp, err := poller.Poll(ctx)
	if err != nil {
		return err
	}
	l.Poller = poller
	l.RawResponse = resp
	return nil
}

// ProviderInstancesCreateResponse contains the response from method ProviderInstances.Create.
type ProviderInstancesCreateResponse struct {
	ProviderInstancesCreateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ProviderInstancesCreateResult contains the result from method ProviderInstances.Create.
type ProviderInstancesCreateResult struct {
	ProviderInstance
}

// ProviderInstancesDeletePollerResponse contains the response from method ProviderInstances.Delete.
type ProviderInstancesDeletePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *ProviderInstancesDeletePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l ProviderInstancesDeletePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (ProviderInstancesDeleteResponse, error) {
	respType := ProviderInstancesDeleteResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, nil)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a ProviderInstancesDeletePollerResponse from the provided client and resume token.
func (l *ProviderInstancesDeletePollerResponse) Resume(ctx context.Context, client *ProviderInstancesClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("ProviderInstancesClient.Delete", token, client.pl, client.deleteHandleError)
	if err != nil {
		return err
	}
	poller := &ProviderInstancesDeletePoller{
		pt: pt,
	}
	resp, err := poller.Poll(ctx)
	if err != nil {
		return err
	}
	l.Poller = poller
	l.RawResponse = resp
	return nil
}

// ProviderInstancesDeleteResponse contains the response from method ProviderInstances.Delete.
type ProviderInstancesDeleteResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ProviderInstancesGetResponse contains the response from method ProviderInstances.Get.
type ProviderInstancesGetResponse struct {
	ProviderInstancesGetResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ProviderInstancesGetResult contains the result from method ProviderInstances.Get.
type ProviderInstancesGetResult struct {
	ProviderInstance
}

// ProviderInstancesListResponse contains the response from method ProviderInstances.List.
type ProviderInstancesListResponse struct {
	ProviderInstancesListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ProviderInstancesListResult contains the result from method ProviderInstances.List.
type ProviderInstancesListResult struct {
	ProviderInstanceListResult
}

// SapMonitorsCreatePollerResponse contains the response from method SapMonitors.Create.
type SapMonitorsCreatePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SapMonitorsCreatePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SapMonitorsCreatePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SapMonitorsCreateResponse, error) {
	respType := SapMonitorsCreateResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, &respType.SapMonitor)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SapMonitorsCreatePollerResponse from the provided client and resume token.
func (l *SapMonitorsCreatePollerResponse) Resume(ctx context.Context, client *SapMonitorsClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SapMonitorsClient.Create", token, client.pl, client.createHandleError)
	if err != nil {
		return err
	}
	poller := &SapMonitorsCreatePoller{
		pt: pt,
	}
	resp, err := poller.Poll(ctx)
	if err != nil {
		return err
	}
	l.Poller = poller
	l.RawResponse = resp
	return nil
}

// SapMonitorsCreateResponse contains the response from method SapMonitors.Create.
type SapMonitorsCreateResponse struct {
	SapMonitorsCreateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SapMonitorsCreateResult contains the result from method SapMonitors.Create.
type SapMonitorsCreateResult struct {
	SapMonitor
}

// SapMonitorsDeletePollerResponse contains the response from method SapMonitors.Delete.
type SapMonitorsDeletePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SapMonitorsDeletePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SapMonitorsDeletePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SapMonitorsDeleteResponse, error) {
	respType := SapMonitorsDeleteResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, nil)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SapMonitorsDeletePollerResponse from the provided client and resume token.
func (l *SapMonitorsDeletePollerResponse) Resume(ctx context.Context, client *SapMonitorsClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SapMonitorsClient.Delete", token, client.pl, client.deleteHandleError)
	if err != nil {
		return err
	}
	poller := &SapMonitorsDeletePoller{
		pt: pt,
	}
	resp, err := poller.Poll(ctx)
	if err != nil {
		return err
	}
	l.Poller = poller
	l.RawResponse = resp
	return nil
}

// SapMonitorsDeleteResponse contains the response from method SapMonitors.Delete.
type SapMonitorsDeleteResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SapMonitorsGetResponse contains the response from method SapMonitors.Get.
type SapMonitorsGetResponse struct {
	SapMonitorsGetResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SapMonitorsGetResult contains the result from method SapMonitors.Get.
type SapMonitorsGetResult struct {
	SapMonitor
}

// SapMonitorsListResponse contains the response from method SapMonitors.List.
type SapMonitorsListResponse struct {
	SapMonitorsListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SapMonitorsListResult contains the result from method SapMonitors.List.
type SapMonitorsListResult struct {
	SapMonitorListResult
}

// SapMonitorsUpdateResponse contains the response from method SapMonitors.Update.
type SapMonitorsUpdateResponse struct {
	SapMonitorsUpdateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SapMonitorsUpdateResult contains the result from method SapMonitors.Update.
type SapMonitorsUpdateResult struct {
	SapMonitor
}

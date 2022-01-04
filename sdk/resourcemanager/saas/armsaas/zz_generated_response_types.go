//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsaas

import (
	"context"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"net/http"
	"time"
)

// ApplicationsListResponse contains the response from method Applications.List.
type ApplicationsListResponse struct {
	ApplicationsListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ApplicationsListResult contains the result from method Applications.List.
type ApplicationsListResult struct {
	SaasAppResponseWithContinuation
}

// OperationsListResponse contains the response from method Operations.List.
type OperationsListResponse struct {
	OperationsListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// OperationsListResult contains the result from method Operations.List.
type OperationsListResult struct {
	SaasAppOperationsResponseWithContinuation
}

// SaaSCreateResourcePollerResponse contains the response from method SaaS.CreateResource.
type SaaSCreateResourcePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaaSCreateResourcePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaaSCreateResourcePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaaSCreateResourceResponse, error) {
	respType := SaaSCreateResourceResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, &respType.SaasResource)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaaSCreateResourcePollerResponse from the provided client and resume token.
func (l *SaaSCreateResourcePollerResponse) Resume(ctx context.Context, client *SaaSClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaaSClient.CreateResource", token, client.pl, client.createResourceHandleError)
	if err != nil {
		return err
	}
	poller := &SaaSCreateResourcePoller{
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

// SaaSCreateResourceResponse contains the response from method SaaS.CreateResource.
type SaaSCreateResourceResponse struct {
	SaaSCreateResourceResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaaSCreateResourceResult contains the result from method SaaS.CreateResource.
type SaaSCreateResourceResult struct {
	SaasResource
}

// SaaSDeletePollerResponse contains the response from method SaaS.Delete.
type SaaSDeletePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaaSDeletePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaaSDeletePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaaSDeleteResponse, error) {
	respType := SaaSDeleteResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, nil)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaaSDeletePollerResponse from the provided client and resume token.
func (l *SaaSDeletePollerResponse) Resume(ctx context.Context, client *SaaSClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaaSClient.Delete", token, client.pl, client.deleteHandleError)
	if err != nil {
		return err
	}
	poller := &SaaSDeletePoller{
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

// SaaSDeleteResponse contains the response from method SaaS.Delete.
type SaaSDeleteResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaaSGetResourceResponse contains the response from method SaaS.GetResource.
type SaaSGetResourceResponse struct {
	SaaSGetResourceResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaaSGetResourceResult contains the result from method SaaS.GetResource.
type SaaSGetResourceResult struct {
	SaasResource
}

// SaaSOperationGetPollerResponse contains the response from method SaaSOperation.Get.
type SaaSOperationGetPollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaaSOperationGetPoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaaSOperationGetPollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaaSOperationGetResponse, error) {
	respType := SaaSOperationGetResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, &respType.SaasResource)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaaSOperationGetPollerResponse from the provided client and resume token.
func (l *SaaSOperationGetPollerResponse) Resume(ctx context.Context, client *SaaSOperationClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaaSOperationClient.Get", token, client.pl, client.getHandleError)
	if err != nil {
		return err
	}
	poller := &SaaSOperationGetPoller{
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

// SaaSOperationGetResponse contains the response from method SaaSOperation.Get.
type SaaSOperationGetResponse struct {
	SaaSOperationGetResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaaSOperationGetResult contains the result from method SaaSOperation.Get.
type SaaSOperationGetResult struct {
	SaasResource
}

// SaaSUpdateResourcePollerResponse contains the response from method SaaS.UpdateResource.
type SaaSUpdateResourcePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaaSUpdateResourcePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaaSUpdateResourcePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaaSUpdateResourceResponse, error) {
	respType := SaaSUpdateResourceResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, &respType.SaasResource)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaaSUpdateResourcePollerResponse from the provided client and resume token.
func (l *SaaSUpdateResourcePollerResponse) Resume(ctx context.Context, client *SaaSClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaaSClient.UpdateResource", token, client.pl, client.updateResourceHandleError)
	if err != nil {
		return err
	}
	poller := &SaaSUpdateResourcePoller{
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

// SaaSUpdateResourceResponse contains the response from method SaaS.UpdateResource.
type SaaSUpdateResourceResponse struct {
	SaaSUpdateResourceResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaaSUpdateResourceResult contains the result from method SaaS.UpdateResource.
type SaaSUpdateResourceResult struct {
	SaasResource
}

// SaasResourcesListAccessTokenResponse contains the response from method SaasResources.ListAccessToken.
type SaasResourcesListAccessTokenResponse struct {
	SaasResourcesListAccessTokenResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasResourcesListAccessTokenResult contains the result from method SaasResources.ListAccessToken.
type SaasResourcesListAccessTokenResult struct {
	AccessTokenResult
}

// SaasResourcesListResponse contains the response from method SaasResources.List.
type SaasResourcesListResponse struct {
	SaasResourcesListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasResourcesListResult contains the result from method SaasResources.List.
type SaasResourcesListResult struct {
	SaasResourceResponseWithContinuation
}

// SaasSubscriptionLevelCreateOrUpdatePollerResponse contains the response from method SaasSubscriptionLevel.CreateOrUpdate.
type SaasSubscriptionLevelCreateOrUpdatePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaasSubscriptionLevelCreateOrUpdatePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaasSubscriptionLevelCreateOrUpdatePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaasSubscriptionLevelCreateOrUpdateResponse, error) {
	respType := SaasSubscriptionLevelCreateOrUpdateResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, &respType.SaasResource)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaasSubscriptionLevelCreateOrUpdatePollerResponse from the provided client and resume token.
func (l *SaasSubscriptionLevelCreateOrUpdatePollerResponse) Resume(ctx context.Context, client *SaasSubscriptionLevelClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaasSubscriptionLevelClient.CreateOrUpdate", token, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return err
	}
	poller := &SaasSubscriptionLevelCreateOrUpdatePoller{
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

// SaasSubscriptionLevelCreateOrUpdateResponse contains the response from method SaasSubscriptionLevel.CreateOrUpdate.
type SaasSubscriptionLevelCreateOrUpdateResponse struct {
	SaasSubscriptionLevelCreateOrUpdateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelCreateOrUpdateResult contains the result from method SaasSubscriptionLevel.CreateOrUpdate.
type SaasSubscriptionLevelCreateOrUpdateResult struct {
	SaasResource
}

// SaasSubscriptionLevelDeletePollerResponse contains the response from method SaasSubscriptionLevel.Delete.
type SaasSubscriptionLevelDeletePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaasSubscriptionLevelDeletePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaasSubscriptionLevelDeletePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaasSubscriptionLevelDeleteResponse, error) {
	respType := SaasSubscriptionLevelDeleteResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, nil)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaasSubscriptionLevelDeletePollerResponse from the provided client and resume token.
func (l *SaasSubscriptionLevelDeletePollerResponse) Resume(ctx context.Context, client *SaasSubscriptionLevelClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaasSubscriptionLevelClient.Delete", token, client.pl, client.deleteHandleError)
	if err != nil {
		return err
	}
	poller := &SaasSubscriptionLevelDeletePoller{
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

// SaasSubscriptionLevelDeleteResponse contains the response from method SaasSubscriptionLevel.Delete.
type SaasSubscriptionLevelDeleteResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelGetResponse contains the response from method SaasSubscriptionLevel.Get.
type SaasSubscriptionLevelGetResponse struct {
	SaasSubscriptionLevelGetResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelGetResult contains the result from method SaasSubscriptionLevel.Get.
type SaasSubscriptionLevelGetResult struct {
	SaasResource
}

// SaasSubscriptionLevelListAccessTokenResponse contains the response from method SaasSubscriptionLevel.ListAccessToken.
type SaasSubscriptionLevelListAccessTokenResponse struct {
	SaasSubscriptionLevelListAccessTokenResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelListAccessTokenResult contains the result from method SaasSubscriptionLevel.ListAccessToken.
type SaasSubscriptionLevelListAccessTokenResult struct {
	AccessTokenResult
}

// SaasSubscriptionLevelListByAzureSubscriptionResponse contains the response from method SaasSubscriptionLevel.ListByAzureSubscription.
type SaasSubscriptionLevelListByAzureSubscriptionResponse struct {
	SaasSubscriptionLevelListByAzureSubscriptionResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelListByAzureSubscriptionResult contains the result from method SaasSubscriptionLevel.ListByAzureSubscription.
type SaasSubscriptionLevelListByAzureSubscriptionResult struct {
	SaasResourceResponseWithContinuation
}

// SaasSubscriptionLevelListByResourceGroupResponse contains the response from method SaasSubscriptionLevel.ListByResourceGroup.
type SaasSubscriptionLevelListByResourceGroupResponse struct {
	SaasSubscriptionLevelListByResourceGroupResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelListByResourceGroupResult contains the result from method SaasSubscriptionLevel.ListByResourceGroup.
type SaasSubscriptionLevelListByResourceGroupResult struct {
	SaasResourceResponseWithContinuation
}

// SaasSubscriptionLevelMoveResourcesPollerResponse contains the response from method SaasSubscriptionLevel.MoveResources.
type SaasSubscriptionLevelMoveResourcesPollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaasSubscriptionLevelMoveResourcesPoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaasSubscriptionLevelMoveResourcesPollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaasSubscriptionLevelMoveResourcesResponse, error) {
	respType := SaasSubscriptionLevelMoveResourcesResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, nil)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaasSubscriptionLevelMoveResourcesPollerResponse from the provided client and resume token.
func (l *SaasSubscriptionLevelMoveResourcesPollerResponse) Resume(ctx context.Context, client *SaasSubscriptionLevelClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaasSubscriptionLevelClient.MoveResources", token, client.pl, client.moveResourcesHandleError)
	if err != nil {
		return err
	}
	poller := &SaasSubscriptionLevelMoveResourcesPoller{
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

// SaasSubscriptionLevelMoveResourcesResponse contains the response from method SaasSubscriptionLevel.MoveResources.
type SaasSubscriptionLevelMoveResourcesResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelUpdatePollerResponse contains the response from method SaasSubscriptionLevel.Update.
type SaasSubscriptionLevelUpdatePollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaasSubscriptionLevelUpdatePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaasSubscriptionLevelUpdatePollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaasSubscriptionLevelUpdateResponse, error) {
	respType := SaasSubscriptionLevelUpdateResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, &respType.SaasResource)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaasSubscriptionLevelUpdatePollerResponse from the provided client and resume token.
func (l *SaasSubscriptionLevelUpdatePollerResponse) Resume(ctx context.Context, client *SaasSubscriptionLevelClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaasSubscriptionLevelClient.Update", token, client.pl, client.updateHandleError)
	if err != nil {
		return err
	}
	poller := &SaasSubscriptionLevelUpdatePoller{
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

// SaasSubscriptionLevelUpdateResponse contains the response from method SaasSubscriptionLevel.Update.
type SaasSubscriptionLevelUpdateResponse struct {
	SaasSubscriptionLevelUpdateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelUpdateResult contains the result from method SaasSubscriptionLevel.Update.
type SaasSubscriptionLevelUpdateResult struct {
	SaasResource
}

// SaasSubscriptionLevelUpdateToUnsubscribedPollerResponse contains the response from method SaasSubscriptionLevel.UpdateToUnsubscribed.
type SaasSubscriptionLevelUpdateToUnsubscribedPollerResponse struct {
	// Poller contains an initialized poller.
	Poller *SaasSubscriptionLevelUpdateToUnsubscribedPoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error is received.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
// A good starting value is 30 seconds. Note that some resources might benefit from a different value.
func (l SaasSubscriptionLevelUpdateToUnsubscribedPollerResponse) PollUntilDone(ctx context.Context, freq time.Duration) (SaasSubscriptionLevelUpdateToUnsubscribedResponse, error) {
	respType := SaasSubscriptionLevelUpdateToUnsubscribedResponse{}
	resp, err := l.Poller.pt.PollUntilDone(ctx, freq, nil)
	if err != nil {
		return respType, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// Resume rehydrates a SaasSubscriptionLevelUpdateToUnsubscribedPollerResponse from the provided client and resume token.
func (l *SaasSubscriptionLevelUpdateToUnsubscribedPollerResponse) Resume(ctx context.Context, client *SaasSubscriptionLevelClient, token string) error {
	pt, err := armruntime.NewPollerFromResumeToken("SaasSubscriptionLevelClient.UpdateToUnsubscribed", token, client.pl, client.updateToUnsubscribedHandleError)
	if err != nil {
		return err
	}
	poller := &SaasSubscriptionLevelUpdateToUnsubscribedPoller{
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

// SaasSubscriptionLevelUpdateToUnsubscribedResponse contains the response from method SaasSubscriptionLevel.UpdateToUnsubscribed.
type SaasSubscriptionLevelUpdateToUnsubscribedResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SaasSubscriptionLevelValidateMoveResourcesResponse contains the response from method SaasSubscriptionLevel.ValidateMoveResources.
type SaasSubscriptionLevelValidateMoveResourcesResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

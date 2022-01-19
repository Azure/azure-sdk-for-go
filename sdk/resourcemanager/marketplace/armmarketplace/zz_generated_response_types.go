//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmarketplace

import "net/http"

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	OperationsClientListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// OperationsClientListResult contains the result from method OperationsClient.List.
type OperationsClientListResult struct {
	OperationListResult
}

// PrivateStoreClientAcknowledgeOfferNotificationResponse contains the response from method PrivateStoreClient.AcknowledgeOfferNotification.
type PrivateStoreClientAcknowledgeOfferNotificationResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientAdminRequestApprovalsListResponse contains the response from method PrivateStoreClient.AdminRequestApprovalsList.
type PrivateStoreClientAdminRequestApprovalsListResponse struct {
	PrivateStoreClientAdminRequestApprovalsListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientAdminRequestApprovalsListResult contains the result from method PrivateStoreClient.AdminRequestApprovalsList.
type PrivateStoreClientAdminRequestApprovalsListResult struct {
	AdminRequestApprovalsList
}

// PrivateStoreClientBillingAccountsResponse contains the response from method PrivateStoreClient.BillingAccounts.
type PrivateStoreClientBillingAccountsResponse struct {
	PrivateStoreClientBillingAccountsResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientBillingAccountsResult contains the result from method PrivateStoreClient.BillingAccounts.
type PrivateStoreClientBillingAccountsResult struct {
	BillingAccountsResponse
}

// PrivateStoreClientBulkCollectionsActionResponse contains the response from method PrivateStoreClient.BulkCollectionsAction.
type PrivateStoreClientBulkCollectionsActionResponse struct {
	PrivateStoreClientBulkCollectionsActionResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientBulkCollectionsActionResult contains the result from method PrivateStoreClient.BulkCollectionsAction.
type PrivateStoreClientBulkCollectionsActionResult struct {
	BulkCollectionsResponse
}

// PrivateStoreClientCollectionsToSubscriptionsMappingResponse contains the response from method PrivateStoreClient.CollectionsToSubscriptionsMapping.
type PrivateStoreClientCollectionsToSubscriptionsMappingResponse struct {
	PrivateStoreClientCollectionsToSubscriptionsMappingResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientCollectionsToSubscriptionsMappingResult contains the result from method PrivateStoreClient.CollectionsToSubscriptionsMapping.
type PrivateStoreClientCollectionsToSubscriptionsMappingResult struct {
	CollectionsToSubscriptionsMappingResponse
}

// PrivateStoreClientCreateApprovalRequestResponse contains the response from method PrivateStoreClient.CreateApprovalRequest.
type PrivateStoreClientCreateApprovalRequestResponse struct {
	PrivateStoreClientCreateApprovalRequestResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientCreateApprovalRequestResult contains the result from method PrivateStoreClient.CreateApprovalRequest.
type PrivateStoreClientCreateApprovalRequestResult struct {
	RequestApprovalResource
}

// PrivateStoreClientCreateOrUpdateResponse contains the response from method PrivateStoreClient.CreateOrUpdate.
type PrivateStoreClientCreateOrUpdateResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientDeleteResponse contains the response from method PrivateStoreClient.Delete.
type PrivateStoreClientDeleteResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientGetAdminRequestApprovalResponse contains the response from method PrivateStoreClient.GetAdminRequestApproval.
type PrivateStoreClientGetAdminRequestApprovalResponse struct {
	PrivateStoreClientGetAdminRequestApprovalResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientGetAdminRequestApprovalResult contains the result from method PrivateStoreClient.GetAdminRequestApproval.
type PrivateStoreClientGetAdminRequestApprovalResult struct {
	AdminRequestApprovalsResource
}

// PrivateStoreClientGetApprovalRequestsListResponse contains the response from method PrivateStoreClient.GetApprovalRequestsList.
type PrivateStoreClientGetApprovalRequestsListResponse struct {
	PrivateStoreClientGetApprovalRequestsListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientGetApprovalRequestsListResult contains the result from method PrivateStoreClient.GetApprovalRequestsList.
type PrivateStoreClientGetApprovalRequestsListResult struct {
	RequestApprovalsList
}

// PrivateStoreClientGetRequestApprovalResponse contains the response from method PrivateStoreClient.GetRequestApproval.
type PrivateStoreClientGetRequestApprovalResponse struct {
	PrivateStoreClientGetRequestApprovalResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientGetRequestApprovalResult contains the result from method PrivateStoreClient.GetRequestApproval.
type PrivateStoreClientGetRequestApprovalResult struct {
	RequestApprovalResource
}

// PrivateStoreClientGetResponse contains the response from method PrivateStoreClient.Get.
type PrivateStoreClientGetResponse struct {
	PrivateStoreClientGetResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientGetResult contains the result from method PrivateStoreClient.Get.
type PrivateStoreClientGetResult struct {
	PrivateStore
}

// PrivateStoreClientListResponse contains the response from method PrivateStoreClient.List.
type PrivateStoreClientListResponse struct {
	PrivateStoreClientListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientListResult contains the result from method PrivateStoreClient.List.
type PrivateStoreClientListResult struct {
	PrivateStoreList
}

// PrivateStoreClientQueryApprovedPlansResponse contains the response from method PrivateStoreClient.QueryApprovedPlans.
type PrivateStoreClientQueryApprovedPlansResponse struct {
	PrivateStoreClientQueryApprovedPlansResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientQueryApprovedPlansResult contains the result from method PrivateStoreClient.QueryApprovedPlans.
type PrivateStoreClientQueryApprovedPlansResult struct {
	QueryApprovedPlansResponse
}

// PrivateStoreClientQueryNotificationsStateResponse contains the response from method PrivateStoreClient.QueryNotificationsState.
type PrivateStoreClientQueryNotificationsStateResponse struct {
	PrivateStoreClientQueryNotificationsStateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientQueryNotificationsStateResult contains the result from method PrivateStoreClient.QueryNotificationsState.
type PrivateStoreClientQueryNotificationsStateResult struct {
	PrivateStoreNotificationsState
}

// PrivateStoreClientQueryOffersResponse contains the response from method PrivateStoreClient.QueryOffers.
type PrivateStoreClientQueryOffersResponse struct {
	PrivateStoreClientQueryOffersResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientQueryOffersResult contains the result from method PrivateStoreClient.QueryOffers.
type PrivateStoreClientQueryOffersResult struct {
	QueryOffers
}

// PrivateStoreClientQueryRequestApprovalResponse contains the response from method PrivateStoreClient.QueryRequestApproval.
type PrivateStoreClientQueryRequestApprovalResponse struct {
	PrivateStoreClientQueryRequestApprovalResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientQueryRequestApprovalResult contains the result from method PrivateStoreClient.QueryRequestApproval.
type PrivateStoreClientQueryRequestApprovalResult struct {
	QueryRequestApproval
}

// PrivateStoreClientUpdateAdminRequestApprovalResponse contains the response from method PrivateStoreClient.UpdateAdminRequestApproval.
type PrivateStoreClientUpdateAdminRequestApprovalResponse struct {
	PrivateStoreClientUpdateAdminRequestApprovalResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreClientUpdateAdminRequestApprovalResult contains the result from method PrivateStoreClient.UpdateAdminRequestApproval.
type PrivateStoreClientUpdateAdminRequestApprovalResult struct {
	AdminRequestApprovalsResource
}

// PrivateStoreClientWithdrawPlanResponse contains the response from method PrivateStoreClient.WithdrawPlan.
type PrivateStoreClientWithdrawPlanResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionClientCreateOrUpdateResponse contains the response from method PrivateStoreCollectionClient.CreateOrUpdate.
type PrivateStoreCollectionClientCreateOrUpdateResponse struct {
	PrivateStoreCollectionClientCreateOrUpdateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionClientCreateOrUpdateResult contains the result from method PrivateStoreCollectionClient.CreateOrUpdate.
type PrivateStoreCollectionClientCreateOrUpdateResult struct {
	Collection
}

// PrivateStoreCollectionClientDeleteResponse contains the response from method PrivateStoreCollectionClient.Delete.
type PrivateStoreCollectionClientDeleteResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionClientGetResponse contains the response from method PrivateStoreCollectionClient.Get.
type PrivateStoreCollectionClientGetResponse struct {
	PrivateStoreCollectionClientGetResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionClientGetResult contains the result from method PrivateStoreCollectionClient.Get.
type PrivateStoreCollectionClientGetResult struct {
	Collection
}

// PrivateStoreCollectionClientListResponse contains the response from method PrivateStoreCollectionClient.List.
type PrivateStoreCollectionClientListResponse struct {
	PrivateStoreCollectionClientListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionClientListResult contains the result from method PrivateStoreCollectionClient.List.
type PrivateStoreCollectionClientListResult struct {
	CollectionsList
}

// PrivateStoreCollectionClientPostResponse contains the response from method PrivateStoreCollectionClient.Post.
type PrivateStoreCollectionClientPostResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionClientTransferOffersResponse contains the response from method PrivateStoreCollectionClient.TransferOffers.
type PrivateStoreCollectionClientTransferOffersResponse struct {
	PrivateStoreCollectionClientTransferOffersResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionClientTransferOffersResult contains the result from method PrivateStoreCollectionClient.TransferOffers.
type PrivateStoreCollectionClientTransferOffersResult struct {
	TransferOffersResponse
}

// PrivateStoreCollectionOfferClientCreateOrUpdateResponse contains the response from method PrivateStoreCollectionOfferClient.CreateOrUpdate.
type PrivateStoreCollectionOfferClientCreateOrUpdateResponse struct {
	PrivateStoreCollectionOfferClientCreateOrUpdateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionOfferClientCreateOrUpdateResult contains the result from method PrivateStoreCollectionOfferClient.CreateOrUpdate.
type PrivateStoreCollectionOfferClientCreateOrUpdateResult struct {
	Offer
}

// PrivateStoreCollectionOfferClientDeleteResponse contains the response from method PrivateStoreCollectionOfferClient.Delete.
type PrivateStoreCollectionOfferClientDeleteResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionOfferClientGetResponse contains the response from method PrivateStoreCollectionOfferClient.Get.
type PrivateStoreCollectionOfferClientGetResponse struct {
	PrivateStoreCollectionOfferClientGetResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionOfferClientGetResult contains the result from method PrivateStoreCollectionOfferClient.Get.
type PrivateStoreCollectionOfferClientGetResult struct {
	Offer
}

// PrivateStoreCollectionOfferClientListResponse contains the response from method PrivateStoreCollectionOfferClient.List.
type PrivateStoreCollectionOfferClientListResponse struct {
	PrivateStoreCollectionOfferClientListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PrivateStoreCollectionOfferClientListResult contains the result from method PrivateStoreCollectionOfferClient.List.
type PrivateStoreCollectionOfferClientListResult struct {
	OfferListResponse
}

// PrivateStoreCollectionOfferClientPostResponse contains the response from method PrivateStoreCollectionOfferClient.Post.
type PrivateStoreCollectionOfferClientPostResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

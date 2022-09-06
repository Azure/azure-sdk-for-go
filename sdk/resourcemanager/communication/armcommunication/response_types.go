//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcommunication

// DomainsClientCancelVerificationResponse contains the response from method DomainsClient.CancelVerification.
type DomainsClientCancelVerificationResponse struct {
	// placeholder for future response values
}

// DomainsClientCreateOrUpdateResponse contains the response from method DomainsClient.CreateOrUpdate.
type DomainsClientCreateOrUpdateResponse struct {
	DomainResource
}

// DomainsClientDeleteResponse contains the response from method DomainsClient.Delete.
type DomainsClientDeleteResponse struct {
	// placeholder for future response values
}

// DomainsClientGetResponse contains the response from method DomainsClient.Get.
type DomainsClientGetResponse struct {
	DomainResource
}

// DomainsClientInitiateVerificationResponse contains the response from method DomainsClient.InitiateVerification.
type DomainsClientInitiateVerificationResponse struct {
	// placeholder for future response values
}

// DomainsClientListByEmailServiceResourceResponse contains the response from method DomainsClient.ListByEmailServiceResource.
type DomainsClientListByEmailServiceResourceResponse struct {
	DomainResourceList
}

// DomainsClientUpdateResponse contains the response from method DomainsClient.Update.
type DomainsClientUpdateResponse struct {
	DomainResource
}

// EmailServicesClientCreateOrUpdateResponse contains the response from method EmailServicesClient.CreateOrUpdate.
type EmailServicesClientCreateOrUpdateResponse struct {
	EmailServiceResource
}

// EmailServicesClientDeleteResponse contains the response from method EmailServicesClient.Delete.
type EmailServicesClientDeleteResponse struct {
	// placeholder for future response values
}

// EmailServicesClientGetResponse contains the response from method EmailServicesClient.Get.
type EmailServicesClientGetResponse struct {
	EmailServiceResource
}

// EmailServicesClientListByResourceGroupResponse contains the response from method EmailServicesClient.ListByResourceGroup.
type EmailServicesClientListByResourceGroupResponse struct {
	EmailServiceResourceList
}

// EmailServicesClientListBySubscriptionResponse contains the response from method EmailServicesClient.ListBySubscription.
type EmailServicesClientListBySubscriptionResponse struct {
	EmailServiceResourceList
}

// EmailServicesClientListVerifiedExchangeOnlineDomainsResponse contains the response from method EmailServicesClient.ListVerifiedExchangeOnlineDomains.
type EmailServicesClientListVerifiedExchangeOnlineDomainsResponse struct {
	// List of FQDNs of verified domains in Exchange Online.
	StringArray []*string
}

// EmailServicesClientUpdateResponse contains the response from method EmailServicesClient.Update.
type EmailServicesClientUpdateResponse struct {
	EmailServiceResource
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	OperationListResult
}

// ServicesClientCheckNameAvailabilityResponse contains the response from method ServicesClient.CheckNameAvailability.
type ServicesClientCheckNameAvailabilityResponse struct {
	CheckNameAvailabilityResponse
}

// ServicesClientCreateOrUpdateResponse contains the response from method ServicesClient.CreateOrUpdate.
type ServicesClientCreateOrUpdateResponse struct {
	ServiceResource
}

// ServicesClientDeleteResponse contains the response from method ServicesClient.Delete.
type ServicesClientDeleteResponse struct {
	// placeholder for future response values
}

// ServicesClientGetResponse contains the response from method ServicesClient.Get.
type ServicesClientGetResponse struct {
	ServiceResource
}

// ServicesClientLinkNotificationHubResponse contains the response from method ServicesClient.LinkNotificationHub.
type ServicesClientLinkNotificationHubResponse struct {
	LinkedNotificationHub
}

// ServicesClientListByResourceGroupResponse contains the response from method ServicesClient.ListByResourceGroup.
type ServicesClientListByResourceGroupResponse struct {
	ServiceResourceList
}

// ServicesClientListBySubscriptionResponse contains the response from method ServicesClient.ListBySubscription.
type ServicesClientListBySubscriptionResponse struct {
	ServiceResourceList
}

// ServicesClientListKeysResponse contains the response from method ServicesClient.ListKeys.
type ServicesClientListKeysResponse struct {
	ServiceKeys
}

// ServicesClientRegenerateKeyResponse contains the response from method ServicesClient.RegenerateKey.
type ServicesClientRegenerateKeyResponse struct {
	ServiceKeys
}

// ServicesClientUpdateResponse contains the response from method ServicesClient.Update.
type ServicesClientUpdateResponse struct {
	ServiceResource
}

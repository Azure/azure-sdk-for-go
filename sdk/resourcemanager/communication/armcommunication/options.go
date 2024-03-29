//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcommunication

// DomainsClientBeginCancelVerificationOptions contains the optional parameters for the DomainsClient.BeginCancelVerification
// method.
type DomainsClientBeginCancelVerificationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// DomainsClientBeginCreateOrUpdateOptions contains the optional parameters for the DomainsClient.BeginCreateOrUpdate method.
type DomainsClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// DomainsClientBeginDeleteOptions contains the optional parameters for the DomainsClient.BeginDelete method.
type DomainsClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// DomainsClientBeginInitiateVerificationOptions contains the optional parameters for the DomainsClient.BeginInitiateVerification
// method.
type DomainsClientBeginInitiateVerificationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// DomainsClientBeginUpdateOptions contains the optional parameters for the DomainsClient.BeginUpdate method.
type DomainsClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// DomainsClientGetOptions contains the optional parameters for the DomainsClient.Get method.
type DomainsClientGetOptions struct {
	// placeholder for future optional parameters
}

// DomainsClientListByEmailServiceResourceOptions contains the optional parameters for the DomainsClient.NewListByEmailServiceResourcePager
// method.
type DomainsClientListByEmailServiceResourceOptions struct {
	// placeholder for future optional parameters
}

// EmailServicesClientBeginCreateOrUpdateOptions contains the optional parameters for the EmailServicesClient.BeginCreateOrUpdate
// method.
type EmailServicesClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// EmailServicesClientBeginDeleteOptions contains the optional parameters for the EmailServicesClient.BeginDelete method.
type EmailServicesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// EmailServicesClientBeginUpdateOptions contains the optional parameters for the EmailServicesClient.BeginUpdate method.
type EmailServicesClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// EmailServicesClientGetOptions contains the optional parameters for the EmailServicesClient.Get method.
type EmailServicesClientGetOptions struct {
	// placeholder for future optional parameters
}

// EmailServicesClientListByResourceGroupOptions contains the optional parameters for the EmailServicesClient.NewListByResourceGroupPager
// method.
type EmailServicesClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// EmailServicesClientListBySubscriptionOptions contains the optional parameters for the EmailServicesClient.NewListBySubscriptionPager
// method.
type EmailServicesClientListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// EmailServicesClientListVerifiedExchangeOnlineDomainsOptions contains the optional parameters for the EmailServicesClient.ListVerifiedExchangeOnlineDomains
// method.
type EmailServicesClientListVerifiedExchangeOnlineDomainsOptions struct {
	// placeholder for future optional parameters
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.NewListPager method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// SenderUsernamesClientCreateOrUpdateOptions contains the optional parameters for the SenderUsernamesClient.CreateOrUpdate
// method.
type SenderUsernamesClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// SenderUsernamesClientDeleteOptions contains the optional parameters for the SenderUsernamesClient.Delete method.
type SenderUsernamesClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// SenderUsernamesClientGetOptions contains the optional parameters for the SenderUsernamesClient.Get method.
type SenderUsernamesClientGetOptions struct {
	// placeholder for future optional parameters
}

// SenderUsernamesClientListByDomainsOptions contains the optional parameters for the SenderUsernamesClient.NewListByDomainsPager
// method.
type SenderUsernamesClientListByDomainsOptions struct {
	// placeholder for future optional parameters
}

// ServicesClientBeginCreateOrUpdateOptions contains the optional parameters for the ServicesClient.BeginCreateOrUpdate method.
type ServicesClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ServicesClientBeginDeleteOptions contains the optional parameters for the ServicesClient.BeginDelete method.
type ServicesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ServicesClientCheckNameAvailabilityOptions contains the optional parameters for the ServicesClient.CheckNameAvailability
// method.
type ServicesClientCheckNameAvailabilityOptions struct {
	// placeholder for future optional parameters
}

// ServicesClientGetOptions contains the optional parameters for the ServicesClient.Get method.
type ServicesClientGetOptions struct {
	// placeholder for future optional parameters
}

// ServicesClientLinkNotificationHubOptions contains the optional parameters for the ServicesClient.LinkNotificationHub method.
type ServicesClientLinkNotificationHubOptions struct {
	// Parameters supplied to the operation.
	LinkNotificationHubParameters *LinkNotificationHubParameters
}

// ServicesClientListByResourceGroupOptions contains the optional parameters for the ServicesClient.NewListByResourceGroupPager
// method.
type ServicesClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// ServicesClientListBySubscriptionOptions contains the optional parameters for the ServicesClient.NewListBySubscriptionPager
// method.
type ServicesClientListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// ServicesClientListKeysOptions contains the optional parameters for the ServicesClient.ListKeys method.
type ServicesClientListKeysOptions struct {
	// placeholder for future optional parameters
}

// ServicesClientRegenerateKeyOptions contains the optional parameters for the ServicesClient.RegenerateKey method.
type ServicesClientRegenerateKeyOptions struct {
	// placeholder for future optional parameters
}

// ServicesClientUpdateOptions contains the optional parameters for the ServicesClient.Update method.
type ServicesClientUpdateOptions struct {
	// placeholder for future optional parameters
}

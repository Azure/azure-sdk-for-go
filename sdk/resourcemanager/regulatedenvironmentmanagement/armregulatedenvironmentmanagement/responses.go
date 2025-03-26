// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armregulatedenvironmentmanagement

// LandingZoneAccountOperationsClientCreateResponse contains the response from method LandingZoneAccountOperationsClient.BeginCreate.
type LandingZoneAccountOperationsClientCreateResponse struct {
	// The Landing zone account resource type. A Landing zone account is the container for configuring, deploying and managing
	// multiple landing zones.
	LZAccount
}

// LandingZoneAccountOperationsClientDeleteResponse contains the response from method LandingZoneAccountOperationsClient.BeginDelete.
type LandingZoneAccountOperationsClientDeleteResponse struct {
	// placeholder for future response values
}

// LandingZoneAccountOperationsClientGetResponse contains the response from method LandingZoneAccountOperationsClient.Get.
type LandingZoneAccountOperationsClientGetResponse struct {
	// The Landing zone account resource type. A Landing zone account is the container for configuring, deploying and managing
	// multiple landing zones.
	LZAccount
}

// LandingZoneAccountOperationsClientListByResourceGroupResponse contains the response from method LandingZoneAccountOperationsClient.NewListByResourceGroupPager.
type LandingZoneAccountOperationsClientListByResourceGroupResponse struct {
	// The response of a LandingZoneAccountResource list operation.
	LandingZoneAccountResourceListResult
}

// LandingZoneAccountOperationsClientListBySubscriptionResponse contains the response from method LandingZoneAccountOperationsClient.NewListBySubscriptionPager.
type LandingZoneAccountOperationsClientListBySubscriptionResponse struct {
	// The response of a LandingZoneAccountResource list operation.
	LandingZoneAccountResourceListResult
}

// LandingZoneAccountOperationsClientUpdateResponse contains the response from method LandingZoneAccountOperationsClient.BeginUpdate.
type LandingZoneAccountOperationsClientUpdateResponse struct {
	// The Landing zone account resource type. A Landing zone account is the container for configuring, deploying and managing
	// multiple landing zones.
	LZAccount
}

// LandingZoneConfigurationOperationsClientCreateCopyResponse contains the response from method LandingZoneConfigurationOperationsClient.BeginCreateCopy.
type LandingZoneConfigurationOperationsClientCreateCopyResponse struct {
	// The response of the create duplicate landing zone configuration.
	CreateLZConfigurationCopyResult
}

// LandingZoneConfigurationOperationsClientCreateResponse contains the response from method LandingZoneConfigurationOperationsClient.BeginCreate.
type LandingZoneConfigurationOperationsClientCreateResponse struct {
	// Concrete proxy resource types can be created by aliasing this type using a specific property type.
	LZConfiguration
}

// LandingZoneConfigurationOperationsClientDeleteResponse contains the response from method LandingZoneConfigurationOperationsClient.BeginDelete.
type LandingZoneConfigurationOperationsClientDeleteResponse struct {
	// placeholder for future response values
}

// LandingZoneConfigurationOperationsClientGenerateLandingZoneResponse contains the response from method LandingZoneConfigurationOperationsClient.BeginGenerateLandingZone.
type LandingZoneConfigurationOperationsClientGenerateLandingZoneResponse struct {
	// The response payload for generating infrastructure-as-code for the landing zone.
	GenerateLandingZoneResult
}

// LandingZoneConfigurationOperationsClientGetResponse contains the response from method LandingZoneConfigurationOperationsClient.Get.
type LandingZoneConfigurationOperationsClientGetResponse struct {
	// Concrete proxy resource types can be created by aliasing this type using a specific property type.
	LZConfiguration
}

// LandingZoneConfigurationOperationsClientListByResourceGroupResponse contains the response from method LandingZoneConfigurationOperationsClient.NewListByResourceGroupPager.
type LandingZoneConfigurationOperationsClientListByResourceGroupResponse struct {
	// The response of a LandingZoneConfigurationResource list operation.
	LandingZoneConfigurationResourceListResult
}

// LandingZoneConfigurationOperationsClientListBySubscriptionResponse contains the response from method LandingZoneConfigurationOperationsClient.NewListBySubscriptionPager.
type LandingZoneConfigurationOperationsClientListBySubscriptionResponse struct {
	// The response of a LandingZoneConfigurationResource list operation.
	LandingZoneConfigurationResourceListResult
}

// LandingZoneConfigurationOperationsClientUpdateAuthoringStatusResponse contains the response from method LandingZoneConfigurationOperationsClient.BeginUpdateAuthoringStatus.
type LandingZoneConfigurationOperationsClientUpdateAuthoringStatusResponse struct {
	// The response for authoring status update request.
	UpdateAuthoringStatusResult
}

// LandingZoneConfigurationOperationsClientUpdateResponse contains the response from method LandingZoneConfigurationOperationsClient.BeginUpdate.
type LandingZoneConfigurationOperationsClientUpdateResponse struct {
	// Concrete proxy resource types can be created by aliasing this type using a specific property type.
	LZConfiguration
}

// LandingZoneRegistrationOperationsClientCreateResponse contains the response from method LandingZoneRegistrationOperationsClient.BeginCreate.
type LandingZoneRegistrationOperationsClientCreateResponse struct {
	// The Landing zone registration resource type.
	LZRegistration
}

// LandingZoneRegistrationOperationsClientDeleteResponse contains the response from method LandingZoneRegistrationOperationsClient.Delete.
type LandingZoneRegistrationOperationsClientDeleteResponse struct {
	// placeholder for future response values
}

// LandingZoneRegistrationOperationsClientGetResponse contains the response from method LandingZoneRegistrationOperationsClient.Get.
type LandingZoneRegistrationOperationsClientGetResponse struct {
	// The Landing zone registration resource type.
	LZRegistration
}

// LandingZoneRegistrationOperationsClientListByResourceGroupResponse contains the response from method LandingZoneRegistrationOperationsClient.NewListByResourceGroupPager.
type LandingZoneRegistrationOperationsClientListByResourceGroupResponse struct {
	// The response of a LandingZoneRegistrationResource list operation.
	LandingZoneRegistrationResourceListResult
}

// LandingZoneRegistrationOperationsClientListBySubscriptionResponse contains the response from method LandingZoneRegistrationOperationsClient.NewListBySubscriptionPager.
type LandingZoneRegistrationOperationsClientListBySubscriptionResponse struct {
	// The response of a LandingZoneRegistrationResource list operation.
	LandingZoneRegistrationResourceListResult
}

// LandingZoneRegistrationOperationsClientUpdateResponse contains the response from method LandingZoneRegistrationOperationsClient.BeginUpdate.
type LandingZoneRegistrationOperationsClientUpdateResponse struct {
	// The Landing zone registration resource type.
	LZRegistration
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to get the next set of results.
	OperationListResult
}

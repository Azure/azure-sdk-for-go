// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armcontainerservicefleet

// AutoUpgradeProfileOperationsClientGenerateUpdateRunResponse contains the response from method AutoUpgradeProfileOperationsClient.BeginGenerateUpdateRun.
type AutoUpgradeProfileOperationsClientGenerateUpdateRunResponse struct {
	// The request should only proceed if an entity matches this string.
	IfMatch *string
}

// AutoUpgradeProfilesClientCreateOrUpdateResponse contains the response from method AutoUpgradeProfilesClient.BeginCreateOrUpdate.
type AutoUpgradeProfilesClientCreateOrUpdateResponse struct {
	// The AutoUpgradeProfile resource.
	AutoUpgradeProfile
}

// AutoUpgradeProfilesClientDeleteResponse contains the response from method AutoUpgradeProfilesClient.BeginDelete.
type AutoUpgradeProfilesClientDeleteResponse struct {
	// placeholder for future response values
}

// AutoUpgradeProfilesClientGetResponse contains the response from method AutoUpgradeProfilesClient.Get.
type AutoUpgradeProfilesClientGetResponse struct {
	// The AutoUpgradeProfile resource.
	AutoUpgradeProfile
}

// AutoUpgradeProfilesClientListByFleetResponse contains the response from method AutoUpgradeProfilesClient.NewListByFleetPager.
type AutoUpgradeProfilesClientListByFleetResponse struct {
	// The response of a AutoUpgradeProfile list operation.
	AutoUpgradeProfileListResult
}

// FleetMembersClientCreateResponse contains the response from method FleetMembersClient.BeginCreate.
type FleetMembersClientCreateResponse struct {
	// A member of the Fleet. It contains a reference to an existing Kubernetes cluster on Azure.
	FleetMember
}

// FleetMembersClientDeleteResponse contains the response from method FleetMembersClient.BeginDelete.
type FleetMembersClientDeleteResponse struct {
	// placeholder for future response values
}

// FleetMembersClientGetResponse contains the response from method FleetMembersClient.Get.
type FleetMembersClientGetResponse struct {
	// A member of the Fleet. It contains a reference to an existing Kubernetes cluster on Azure.
	FleetMember
}

// FleetMembersClientListByFleetResponse contains the response from method FleetMembersClient.NewListByFleetPager.
type FleetMembersClientListByFleetResponse struct {
	// The response of a FleetMember list operation.
	FleetMemberListResult
}

// FleetMembersClientUpdateAsyncResponse contains the response from method FleetMembersClient.BeginUpdateAsync.
type FleetMembersClientUpdateAsyncResponse struct {
	// A member of the Fleet. It contains a reference to an existing Kubernetes cluster on Azure.
	FleetMember
}

// FleetUpdateStrategiesClientCreateOrUpdateResponse contains the response from method FleetUpdateStrategiesClient.BeginCreateOrUpdate.
type FleetUpdateStrategiesClientCreateOrUpdateResponse struct {
	// Defines a multi-stage process to perform update operations across members of a Fleet.
	FleetUpdateStrategy
}

// FleetUpdateStrategiesClientDeleteResponse contains the response from method FleetUpdateStrategiesClient.BeginDelete.
type FleetUpdateStrategiesClientDeleteResponse struct {
	// placeholder for future response values
}

// FleetUpdateStrategiesClientGetResponse contains the response from method FleetUpdateStrategiesClient.Get.
type FleetUpdateStrategiesClientGetResponse struct {
	// Defines a multi-stage process to perform update operations across members of a Fleet.
	FleetUpdateStrategy
}

// FleetUpdateStrategiesClientListByFleetResponse contains the response from method FleetUpdateStrategiesClient.NewListByFleetPager.
type FleetUpdateStrategiesClientListByFleetResponse struct {
	// The response of a FleetUpdateStrategy list operation.
	FleetUpdateStrategyListResult
}

// FleetsClientCreateResponse contains the response from method FleetsClient.BeginCreate.
type FleetsClientCreateResponse struct {
	// The Fleet resource.
	Fleet
}

// FleetsClientDeleteResponse contains the response from method FleetsClient.BeginDelete.
type FleetsClientDeleteResponse struct {
	// placeholder for future response values
}

// FleetsClientGetResponse contains the response from method FleetsClient.Get.
type FleetsClientGetResponse struct {
	// The Fleet resource.
	Fleet
}

// FleetsClientListByResourceGroupResponse contains the response from method FleetsClient.NewListByResourceGroupPager.
type FleetsClientListByResourceGroupResponse struct {
	// The response of a Fleet list operation.
	FleetListResult
}

// FleetsClientListBySubscriptionResponse contains the response from method FleetsClient.NewListBySubscriptionPager.
type FleetsClientListBySubscriptionResponse struct {
	// The response of a Fleet list operation.
	FleetListResult
}

// FleetsClientListCredentialsResponse contains the response from method FleetsClient.ListCredentials.
type FleetsClientListCredentialsResponse struct {
	// The Credential results response.
	FleetCredentialResults
}

// FleetsClientUpdateAsyncResponse contains the response from method FleetsClient.BeginUpdateAsync.
type FleetsClientUpdateAsyncResponse struct {
	// The Fleet resource.
	Fleet
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to get the next set of results.
	OperationListResult
}

// UpdateRunsClientCreateOrUpdateResponse contains the response from method UpdateRunsClient.BeginCreateOrUpdate.
type UpdateRunsClientCreateOrUpdateResponse struct {
	// A multi-stage process to perform update operations across members of a Fleet.
	UpdateRun
}

// UpdateRunsClientDeleteResponse contains the response from method UpdateRunsClient.BeginDelete.
type UpdateRunsClientDeleteResponse struct {
	// placeholder for future response values
}

// UpdateRunsClientGetResponse contains the response from method UpdateRunsClient.Get.
type UpdateRunsClientGetResponse struct {
	// A multi-stage process to perform update operations across members of a Fleet.
	UpdateRun
}

// UpdateRunsClientListByFleetResponse contains the response from method UpdateRunsClient.NewListByFleetPager.
type UpdateRunsClientListByFleetResponse struct {
	// The response of a UpdateRun list operation.
	UpdateRunListResult
}

// UpdateRunsClientSkipResponse contains the response from method UpdateRunsClient.BeginSkip.
type UpdateRunsClientSkipResponse struct {
	// A multi-stage process to perform update operations across members of a Fleet.
	UpdateRun
}

// UpdateRunsClientStartResponse contains the response from method UpdateRunsClient.BeginStart.
type UpdateRunsClientStartResponse struct {
	// A multi-stage process to perform update operations across members of a Fleet.
	UpdateRun
}

// UpdateRunsClientStopResponse contains the response from method UpdateRunsClient.BeginStop.
type UpdateRunsClientStopResponse struct {
	// A multi-stage process to perform update operations across members of a Fleet.
	UpdateRun
}

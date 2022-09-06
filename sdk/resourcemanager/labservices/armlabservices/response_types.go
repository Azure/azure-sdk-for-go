//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armlabservices

// ImagesClientCreateOrUpdateResponse contains the response from method ImagesClient.CreateOrUpdate.
type ImagesClientCreateOrUpdateResponse struct {
	Image
}

// ImagesClientGetResponse contains the response from method ImagesClient.Get.
type ImagesClientGetResponse struct {
	Image
}

// ImagesClientListByLabPlanResponse contains the response from method ImagesClient.ListByLabPlan.
type ImagesClientListByLabPlanResponse struct {
	PagedImages
}

// ImagesClientUpdateResponse contains the response from method ImagesClient.Update.
type ImagesClientUpdateResponse struct {
	Image
}

// LabPlansClientCreateOrUpdateResponse contains the response from method LabPlansClient.CreateOrUpdate.
type LabPlansClientCreateOrUpdateResponse struct {
	LabPlan
}

// LabPlansClientDeleteResponse contains the response from method LabPlansClient.Delete.
type LabPlansClientDeleteResponse struct {
	// placeholder for future response values
}

// LabPlansClientGetResponse contains the response from method LabPlansClient.Get.
type LabPlansClientGetResponse struct {
	LabPlan
}

// LabPlansClientListByResourceGroupResponse contains the response from method LabPlansClient.ListByResourceGroup.
type LabPlansClientListByResourceGroupResponse struct {
	PagedLabPlans
}

// LabPlansClientListBySubscriptionResponse contains the response from method LabPlansClient.ListBySubscription.
type LabPlansClientListBySubscriptionResponse struct {
	PagedLabPlans
}

// LabPlansClientSaveImageResponse contains the response from method LabPlansClient.SaveImage.
type LabPlansClientSaveImageResponse struct {
	// placeholder for future response values
}

// LabPlansClientUpdateResponse contains the response from method LabPlansClient.Update.
type LabPlansClientUpdateResponse struct {
	LabPlan
}

// LabsClientCreateOrUpdateResponse contains the response from method LabsClient.CreateOrUpdate.
type LabsClientCreateOrUpdateResponse struct {
	Lab
}

// LabsClientDeleteResponse contains the response from method LabsClient.Delete.
type LabsClientDeleteResponse struct {
	// placeholder for future response values
}

// LabsClientGetResponse contains the response from method LabsClient.Get.
type LabsClientGetResponse struct {
	Lab
}

// LabsClientListByResourceGroupResponse contains the response from method LabsClient.ListByResourceGroup.
type LabsClientListByResourceGroupResponse struct {
	PagedLabs
}

// LabsClientListBySubscriptionResponse contains the response from method LabsClient.ListBySubscription.
type LabsClientListBySubscriptionResponse struct {
	PagedLabs
}

// LabsClientPublishResponse contains the response from method LabsClient.Publish.
type LabsClientPublishResponse struct {
	// placeholder for future response values
}

// LabsClientSyncGroupResponse contains the response from method LabsClient.SyncGroup.
type LabsClientSyncGroupResponse struct {
	// placeholder for future response values
}

// LabsClientUpdateResponse contains the response from method LabsClient.Update.
type LabsClientUpdateResponse struct {
	Lab
}

// OperationResultsClientGetResponse contains the response from method OperationResultsClient.Get.
type OperationResultsClientGetResponse struct {
	OperationResult
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	OperationListResult
}

// SKUsClientListResponse contains the response from method SKUsClient.List.
type SKUsClientListResponse struct {
	PagedSKUInfos
}

// SchedulesClientCreateOrUpdateResponse contains the response from method SchedulesClient.CreateOrUpdate.
type SchedulesClientCreateOrUpdateResponse struct {
	Schedule
}

// SchedulesClientDeleteResponse contains the response from method SchedulesClient.Delete.
type SchedulesClientDeleteResponse struct {
	// placeholder for future response values
}

// SchedulesClientGetResponse contains the response from method SchedulesClient.Get.
type SchedulesClientGetResponse struct {
	Schedule
}

// SchedulesClientListByLabResponse contains the response from method SchedulesClient.ListByLab.
type SchedulesClientListByLabResponse struct {
	PagedSchedules
}

// SchedulesClientUpdateResponse contains the response from method SchedulesClient.Update.
type SchedulesClientUpdateResponse struct {
	Schedule
}

// UsagesClientListByLocationResponse contains the response from method UsagesClient.ListByLocation.
type UsagesClientListByLocationResponse struct {
	ListUsagesResult
}

// UsersClientCreateOrUpdateResponse contains the response from method UsersClient.CreateOrUpdate.
type UsersClientCreateOrUpdateResponse struct {
	User
}

// UsersClientDeleteResponse contains the response from method UsersClient.Delete.
type UsersClientDeleteResponse struct {
	// placeholder for future response values
}

// UsersClientGetResponse contains the response from method UsersClient.Get.
type UsersClientGetResponse struct {
	User
}

// UsersClientInviteResponse contains the response from method UsersClient.Invite.
type UsersClientInviteResponse struct {
	// placeholder for future response values
}

// UsersClientListByLabResponse contains the response from method UsersClient.ListByLab.
type UsersClientListByLabResponse struct {
	PagedUsers
}

// UsersClientUpdateResponse contains the response from method UsersClient.Update.
type UsersClientUpdateResponse struct {
	User
}

// VirtualMachinesClientGetResponse contains the response from method VirtualMachinesClient.Get.
type VirtualMachinesClientGetResponse struct {
	VirtualMachine
}

// VirtualMachinesClientListByLabResponse contains the response from method VirtualMachinesClient.ListByLab.
type VirtualMachinesClientListByLabResponse struct {
	PagedVirtualMachines
}

// VirtualMachinesClientRedeployResponse contains the response from method VirtualMachinesClient.Redeploy.
type VirtualMachinesClientRedeployResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientReimageResponse contains the response from method VirtualMachinesClient.Reimage.
type VirtualMachinesClientReimageResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientResetPasswordResponse contains the response from method VirtualMachinesClient.ResetPassword.
type VirtualMachinesClientResetPasswordResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientStartResponse contains the response from method VirtualMachinesClient.Start.
type VirtualMachinesClientStartResponse struct {
	// placeholder for future response values
}

// VirtualMachinesClientStopResponse contains the response from method VirtualMachinesClient.Stop.
type VirtualMachinesClientStopResponse struct {
	// placeholder for future response values
}

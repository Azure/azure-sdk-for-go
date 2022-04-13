package authorizationapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2020-10-01/authorization"
	"github.com/Azure/go-autorest/autorest"
)

// PermissionsClientAPI contains the set of methods on the PermissionsClient type.
type PermissionsClientAPI interface {
	ListForResource(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string) (result authorization.PermissionGetResultPage, err error)
	ListForResourceComplete(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string) (result authorization.PermissionGetResultIterator, err error)
	ListForResourceGroup(ctx context.Context, resourceGroupName string) (result authorization.PermissionGetResultPage, err error)
	ListForResourceGroupComplete(ctx context.Context, resourceGroupName string) (result authorization.PermissionGetResultIterator, err error)
}

var _ PermissionsClientAPI = (*authorization.PermissionsClient)(nil)

// RoleDefinitionsClientAPI contains the set of methods on the RoleDefinitionsClient type.
type RoleDefinitionsClientAPI interface {
	CreateOrUpdate(ctx context.Context, scope string, roleDefinitionID string, roleDefinition authorization.RoleDefinition) (result authorization.RoleDefinition, err error)
	Delete(ctx context.Context, scope string, roleDefinitionID string) (result authorization.RoleDefinition, err error)
	Get(ctx context.Context, scope string, roleDefinitionID string) (result authorization.RoleDefinition, err error)
	GetByID(ctx context.Context, roleDefinitionID string) (result authorization.RoleDefinition, err error)
	List(ctx context.Context, scope string, filter string) (result authorization.RoleDefinitionListResultPage, err error)
	ListComplete(ctx context.Context, scope string, filter string) (result authorization.RoleDefinitionListResultIterator, err error)
}

var _ RoleDefinitionsClientAPI = (*authorization.RoleDefinitionsClient)(nil)

// ProviderOperationsMetadataClientAPI contains the set of methods on the ProviderOperationsMetadataClient type.
type ProviderOperationsMetadataClientAPI interface {
	Get(ctx context.Context, resourceProviderNamespace string, expand string) (result authorization.ProviderOperationsMetadata, err error)
	List(ctx context.Context, expand string) (result authorization.ProviderOperationsMetadataListResultPage, err error)
	ListComplete(ctx context.Context, expand string) (result authorization.ProviderOperationsMetadataListResultIterator, err error)
}

var _ ProviderOperationsMetadataClientAPI = (*authorization.ProviderOperationsMetadataClient)(nil)

// GlobalAdministratorClientAPI contains the set of methods on the GlobalAdministratorClient type.
type GlobalAdministratorClientAPI interface {
	ElevateAccess(ctx context.Context) (result autorest.Response, err error)
}

var _ GlobalAdministratorClientAPI = (*authorization.GlobalAdministratorClient)(nil)

// RoleAssignmentsClientAPI contains the set of methods on the RoleAssignmentsClient type.
type RoleAssignmentsClientAPI interface {
	Create(ctx context.Context, scope string, roleAssignmentName string, parameters authorization.RoleAssignmentCreateParameters) (result authorization.RoleAssignment, err error)
	CreateByID(ctx context.Context, roleAssignmentID string, parameters authorization.RoleAssignmentCreateParameters) (result authorization.RoleAssignment, err error)
	Delete(ctx context.Context, scope string, roleAssignmentName string) (result authorization.RoleAssignment, err error)
	DeleteByID(ctx context.Context, roleAssignmentID string) (result authorization.RoleAssignment, err error)
	Get(ctx context.Context, scope string, roleAssignmentName string) (result authorization.RoleAssignment, err error)
	GetByID(ctx context.Context, roleAssignmentID string) (result authorization.RoleAssignment, err error)
	List(ctx context.Context, filter string) (result authorization.RoleAssignmentListResultPage, err error)
	ListComplete(ctx context.Context, filter string) (result authorization.RoleAssignmentListResultIterator, err error)
	ListForResource(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, filter string) (result authorization.RoleAssignmentListResultPage, err error)
	ListForResourceComplete(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, filter string) (result authorization.RoleAssignmentListResultIterator, err error)
	ListForResourceGroup(ctx context.Context, resourceGroupName string, filter string) (result authorization.RoleAssignmentListResultPage, err error)
	ListForResourceGroupComplete(ctx context.Context, resourceGroupName string, filter string) (result authorization.RoleAssignmentListResultIterator, err error)
	ListForScope(ctx context.Context, scope string, filter string) (result authorization.RoleAssignmentListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string, filter string) (result authorization.RoleAssignmentListResultIterator, err error)
}

var _ RoleAssignmentsClientAPI = (*authorization.RoleAssignmentsClient)(nil)

// ClassicAdministratorsClientAPI contains the set of methods on the ClassicAdministratorsClient type.
type ClassicAdministratorsClientAPI interface {
	List(ctx context.Context) (result authorization.ClassicAdministratorListResultPage, err error)
	ListComplete(ctx context.Context) (result authorization.ClassicAdministratorListResultIterator, err error)
}

var _ ClassicAdministratorsClientAPI = (*authorization.ClassicAdministratorsClient)(nil)

// EligibleChildResourcesClientAPI contains the set of methods on the EligibleChildResourcesClient type.
type EligibleChildResourcesClientAPI interface {
	Get(ctx context.Context, scope string, filter string) (result authorization.EligibleChildResourcesListResultPage, err error)
	GetComplete(ctx context.Context, scope string, filter string) (result authorization.EligibleChildResourcesListResultIterator, err error)
}

var _ EligibleChildResourcesClientAPI = (*authorization.EligibleChildResourcesClient)(nil)

// RoleAssignmentSchedulesClientAPI contains the set of methods on the RoleAssignmentSchedulesClient type.
type RoleAssignmentSchedulesClientAPI interface {
	Get(ctx context.Context, scope string, roleAssignmentScheduleName string) (result authorization.RoleAssignmentSchedule, err error)
	ListForScope(ctx context.Context, scope string, filter string) (result authorization.RoleAssignmentScheduleListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string, filter string) (result authorization.RoleAssignmentScheduleListResultIterator, err error)
}

var _ RoleAssignmentSchedulesClientAPI = (*authorization.RoleAssignmentSchedulesClient)(nil)

// RoleAssignmentScheduleInstancesClientAPI contains the set of methods on the RoleAssignmentScheduleInstancesClient type.
type RoleAssignmentScheduleInstancesClientAPI interface {
	Get(ctx context.Context, scope string, roleAssignmentScheduleInstanceName string) (result authorization.RoleAssignmentScheduleInstance, err error)
	ListForScope(ctx context.Context, scope string, filter string) (result authorization.RoleAssignmentScheduleInstanceListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string, filter string) (result authorization.RoleAssignmentScheduleInstanceListResultIterator, err error)
}

var _ RoleAssignmentScheduleInstancesClientAPI = (*authorization.RoleAssignmentScheduleInstancesClient)(nil)

// RoleAssignmentScheduleRequestsClientAPI contains the set of methods on the RoleAssignmentScheduleRequestsClient type.
type RoleAssignmentScheduleRequestsClientAPI interface {
	Cancel(ctx context.Context, scope string, roleAssignmentScheduleRequestName string) (result autorest.Response, err error)
	Create(ctx context.Context, scope string, roleAssignmentScheduleRequestName string, parameters authorization.RoleAssignmentScheduleRequest) (result authorization.RoleAssignmentScheduleRequest, err error)
	Get(ctx context.Context, scope string, roleAssignmentScheduleRequestName string) (result authorization.RoleAssignmentScheduleRequest, err error)
	ListForScope(ctx context.Context, scope string, filter string) (result authorization.RoleAssignmentScheduleRequestListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string, filter string) (result authorization.RoleAssignmentScheduleRequestListResultIterator, err error)
	Validate(ctx context.Context, scope string, roleAssignmentScheduleRequestName string, parameters authorization.RoleAssignmentScheduleRequest) (result authorization.RoleAssignmentScheduleRequest, err error)
}

var _ RoleAssignmentScheduleRequestsClientAPI = (*authorization.RoleAssignmentScheduleRequestsClient)(nil)

// RoleEligibilitySchedulesClientAPI contains the set of methods on the RoleEligibilitySchedulesClient type.
type RoleEligibilitySchedulesClientAPI interface {
	Get(ctx context.Context, scope string, roleEligibilityScheduleName string) (result authorization.RoleEligibilitySchedule, err error)
	ListForScope(ctx context.Context, scope string, filter string) (result authorization.RoleEligibilityScheduleListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string, filter string) (result authorization.RoleEligibilityScheduleListResultIterator, err error)
}

var _ RoleEligibilitySchedulesClientAPI = (*authorization.RoleEligibilitySchedulesClient)(nil)

// RoleEligibilityScheduleInstancesClientAPI contains the set of methods on the RoleEligibilityScheduleInstancesClient type.
type RoleEligibilityScheduleInstancesClientAPI interface {
	Get(ctx context.Context, scope string, roleEligibilityScheduleInstanceName string) (result authorization.RoleEligibilityScheduleInstance, err error)
	ListForScope(ctx context.Context, scope string, filter string) (result authorization.RoleEligibilityScheduleInstanceListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string, filter string) (result authorization.RoleEligibilityScheduleInstanceListResultIterator, err error)
}

var _ RoleEligibilityScheduleInstancesClientAPI = (*authorization.RoleEligibilityScheduleInstancesClient)(nil)

// RoleEligibilityScheduleRequestsClientAPI contains the set of methods on the RoleEligibilityScheduleRequestsClient type.
type RoleEligibilityScheduleRequestsClientAPI interface {
	Cancel(ctx context.Context, scope string, roleEligibilityScheduleRequestName string) (result autorest.Response, err error)
	Create(ctx context.Context, scope string, roleEligibilityScheduleRequestName string, parameters authorization.RoleEligibilityScheduleRequest) (result authorization.RoleEligibilityScheduleRequest, err error)
	Get(ctx context.Context, scope string, roleEligibilityScheduleRequestName string) (result authorization.RoleEligibilityScheduleRequest, err error)
	ListForScope(ctx context.Context, scope string, filter string) (result authorization.RoleEligibilityScheduleRequestListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string, filter string) (result authorization.RoleEligibilityScheduleRequestListResultIterator, err error)
	Validate(ctx context.Context, scope string, roleEligibilityScheduleRequestName string, parameters authorization.RoleEligibilityScheduleRequest) (result authorization.RoleEligibilityScheduleRequest, err error)
}

var _ RoleEligibilityScheduleRequestsClientAPI = (*authorization.RoleEligibilityScheduleRequestsClient)(nil)

// RoleManagementPoliciesClientAPI contains the set of methods on the RoleManagementPoliciesClient type.
type RoleManagementPoliciesClientAPI interface {
	Delete(ctx context.Context, scope string, roleManagementPolicyName string) (result autorest.Response, err error)
	Get(ctx context.Context, scope string, roleManagementPolicyName string) (result authorization.RoleManagementPolicy, err error)
	ListForScope(ctx context.Context, scope string) (result authorization.RoleManagementPolicyListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string) (result authorization.RoleManagementPolicyListResultIterator, err error)
	Update(ctx context.Context, scope string, roleManagementPolicyName string, parameters authorization.RoleManagementPolicy) (result authorization.RoleManagementPolicy, err error)
}

var _ RoleManagementPoliciesClientAPI = (*authorization.RoleManagementPoliciesClient)(nil)

// RoleManagementPolicyAssignmentsClientAPI contains the set of methods on the RoleManagementPolicyAssignmentsClient type.
type RoleManagementPolicyAssignmentsClientAPI interface {
	Create(ctx context.Context, scope string, roleManagementPolicyAssignmentName string, parameters authorization.RoleManagementPolicyAssignment) (result authorization.RoleManagementPolicyAssignment, err error)
	Delete(ctx context.Context, scope string, roleManagementPolicyAssignmentName string) (result autorest.Response, err error)
	Get(ctx context.Context, scope string, roleManagementPolicyAssignmentName string) (result authorization.RoleManagementPolicyAssignment, err error)
	ListForScope(ctx context.Context, scope string) (result authorization.RoleManagementPolicyAssignmentListResultPage, err error)
	ListForScopeComplete(ctx context.Context, scope string) (result authorization.RoleManagementPolicyAssignmentListResultIterator, err error)
}

var _ RoleManagementPolicyAssignmentsClientAPI = (*authorization.RoleManagementPolicyAssignmentsClient)(nil)

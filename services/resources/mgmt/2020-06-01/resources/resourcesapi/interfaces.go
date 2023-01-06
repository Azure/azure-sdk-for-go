// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package resourcesapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/Azure/go-autorest/autorest"
)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result resources.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result resources.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*resources.OperationsClient)(nil)

// DeploymentsClientAPI contains the set of methods on the DeploymentsClient type.
type DeploymentsClientAPI interface {
	CalculateTemplateHash(ctx context.Context, templateParameter interface{}) (result resources.TemplateHashResult, err error)
	Cancel(ctx context.Context, resourceGroupName string, deploymentName string) (result autorest.Response, err error)
	CancelAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string) (result autorest.Response, err error)
	CancelAtScope(ctx context.Context, scope string, deploymentName string) (result autorest.Response, err error)
	CancelAtSubscriptionScope(ctx context.Context, deploymentName string) (result autorest.Response, err error)
	CancelAtTenantScope(ctx context.Context, deploymentName string) (result autorest.Response, err error)
	CheckExistence(ctx context.Context, resourceGroupName string, deploymentName string) (result autorest.Response, err error)
	CheckExistenceAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string) (result autorest.Response, err error)
	CheckExistenceAtScope(ctx context.Context, scope string, deploymentName string) (result autorest.Response, err error)
	CheckExistenceAtSubscriptionScope(ctx context.Context, deploymentName string) (result autorest.Response, err error)
	CheckExistenceAtTenantScope(ctx context.Context, deploymentName string) (result autorest.Response, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters resources.Deployment) (result resources.DeploymentsCreateOrUpdateFuture, err error)
	CreateOrUpdateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters resources.ScopedDeployment) (result resources.DeploymentsCreateOrUpdateAtManagementGroupScopeFuture, err error)
	CreateOrUpdateAtScope(ctx context.Context, scope string, deploymentName string, parameters resources.Deployment) (result resources.DeploymentsCreateOrUpdateAtScopeFuture, err error)
	CreateOrUpdateAtSubscriptionScope(ctx context.Context, deploymentName string, parameters resources.Deployment) (result resources.DeploymentsCreateOrUpdateAtSubscriptionScopeFuture, err error)
	CreateOrUpdateAtTenantScope(ctx context.Context, deploymentName string, parameters resources.ScopedDeployment) (result resources.DeploymentsCreateOrUpdateAtTenantScopeFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentsDeleteFuture, err error)
	DeleteAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string) (result resources.DeploymentsDeleteAtManagementGroupScopeFuture, err error)
	DeleteAtScope(ctx context.Context, scope string, deploymentName string) (result resources.DeploymentsDeleteAtScopeFuture, err error)
	DeleteAtSubscriptionScope(ctx context.Context, deploymentName string) (result resources.DeploymentsDeleteAtSubscriptionScopeFuture, err error)
	DeleteAtTenantScope(ctx context.Context, deploymentName string) (result resources.DeploymentsDeleteAtTenantScopeFuture, err error)
	ExportTemplate(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentExportResult, err error)
	ExportTemplateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string) (result resources.DeploymentExportResult, err error)
	ExportTemplateAtScope(ctx context.Context, scope string, deploymentName string) (result resources.DeploymentExportResult, err error)
	ExportTemplateAtSubscriptionScope(ctx context.Context, deploymentName string) (result resources.DeploymentExportResult, err error)
	ExportTemplateAtTenantScope(ctx context.Context, deploymentName string) (result resources.DeploymentExportResult, err error)
	Get(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentExtended, err error)
	GetAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string) (result resources.DeploymentExtended, err error)
	GetAtScope(ctx context.Context, scope string, deploymentName string) (result resources.DeploymentExtended, err error)
	GetAtSubscriptionScope(ctx context.Context, deploymentName string) (result resources.DeploymentExtended, err error)
	GetAtTenantScope(ctx context.Context, deploymentName string) (result resources.DeploymentExtended, err error)
	ListAtManagementGroupScope(ctx context.Context, groupID string, filter string, top *int32) (result resources.DeploymentListResultPage, err error)
	ListAtManagementGroupScopeComplete(ctx context.Context, groupID string, filter string, top *int32) (result resources.DeploymentListResultIterator, err error)
	ListAtScope(ctx context.Context, scope string, filter string, top *int32) (result resources.DeploymentListResultPage, err error)
	ListAtScopeComplete(ctx context.Context, scope string, filter string, top *int32) (result resources.DeploymentListResultIterator, err error)
	ListAtSubscriptionScope(ctx context.Context, filter string, top *int32) (result resources.DeploymentListResultPage, err error)
	ListAtSubscriptionScopeComplete(ctx context.Context, filter string, top *int32) (result resources.DeploymentListResultIterator, err error)
	ListAtTenantScope(ctx context.Context, filter string, top *int32) (result resources.DeploymentListResultPage, err error)
	ListAtTenantScopeComplete(ctx context.Context, filter string, top *int32) (result resources.DeploymentListResultIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string, filter string, top *int32) (result resources.DeploymentListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string, filter string, top *int32) (result resources.DeploymentListResultIterator, err error)
	Validate(ctx context.Context, resourceGroupName string, deploymentName string, parameters resources.Deployment) (result resources.DeploymentsValidateFuture, err error)
	ValidateAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters resources.ScopedDeployment) (result resources.DeploymentsValidateAtManagementGroupScopeFuture, err error)
	ValidateAtScope(ctx context.Context, scope string, deploymentName string, parameters resources.Deployment) (result resources.DeploymentsValidateAtScopeFuture, err error)
	ValidateAtSubscriptionScope(ctx context.Context, deploymentName string, parameters resources.Deployment) (result resources.DeploymentsValidateAtSubscriptionScopeFuture, err error)
	ValidateAtTenantScope(ctx context.Context, deploymentName string, parameters resources.ScopedDeployment) (result resources.DeploymentsValidateAtTenantScopeFuture, err error)
	WhatIf(ctx context.Context, resourceGroupName string, deploymentName string, parameters resources.DeploymentWhatIf) (result resources.DeploymentsWhatIfFuture, err error)
	WhatIfAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, parameters resources.ScopedDeploymentWhatIf) (result resources.DeploymentsWhatIfAtManagementGroupScopeFuture, err error)
	WhatIfAtSubscriptionScope(ctx context.Context, deploymentName string, parameters resources.DeploymentWhatIf) (result resources.DeploymentsWhatIfAtSubscriptionScopeFuture, err error)
	WhatIfAtTenantScope(ctx context.Context, deploymentName string, parameters resources.ScopedDeploymentWhatIf) (result resources.DeploymentsWhatIfAtTenantScopeFuture, err error)
}

var _ DeploymentsClientAPI = (*resources.DeploymentsClient)(nil)

// ProvidersClientAPI contains the set of methods on the ProvidersClient type.
type ProvidersClientAPI interface {
	Get(ctx context.Context, resourceProviderNamespace string, expand string) (result resources.Provider, err error)
	GetAtTenantScope(ctx context.Context, resourceProviderNamespace string, expand string) (result resources.Provider, err error)
	List(ctx context.Context, top *int32, expand string) (result resources.ProviderListResultPage, err error)
	ListComplete(ctx context.Context, top *int32, expand string) (result resources.ProviderListResultIterator, err error)
	ListAtTenantScope(ctx context.Context, top *int32, expand string) (result resources.ProviderListResultPage, err error)
	ListAtTenantScopeComplete(ctx context.Context, top *int32, expand string) (result resources.ProviderListResultIterator, err error)
	Register(ctx context.Context, resourceProviderNamespace string) (result resources.Provider, err error)
	RegisterAtManagementGroupScope(ctx context.Context, resourceProviderNamespace string, groupID string) (result autorest.Response, err error)
	Unregister(ctx context.Context, resourceProviderNamespace string) (result resources.Provider, err error)
}

var _ ProvidersClientAPI = (*resources.ProvidersClient)(nil)

// ClientAPI contains the set of methods on the Client type.
type ClientAPI interface {
	CheckExistence(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, APIVersion string) (result autorest.Response, err error)
	CheckExistenceByID(ctx context.Context, resourceID string, APIVersion string) (result autorest.Response, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, APIVersion string, parameters resources.GenericResource) (result resources.CreateOrUpdateFuture, err error)
	CreateOrUpdateByID(ctx context.Context, resourceID string, APIVersion string, parameters resources.GenericResource) (result resources.CreateOrUpdateByIDFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, APIVersion string) (result resources.DeleteFuture, err error)
	DeleteByID(ctx context.Context, resourceID string, APIVersion string) (result resources.DeleteByIDFuture, err error)
	Get(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, APIVersion string) (result resources.GenericResource, err error)
	GetByID(ctx context.Context, resourceID string, APIVersion string) (result resources.GenericResource, err error)
	List(ctx context.Context, filter string, expand string, top *int32) (result resources.ListResultPage, err error)
	ListComplete(ctx context.Context, filter string, expand string, top *int32) (result resources.ListResultIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string, filter string, expand string, top *int32) (result resources.ListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string, filter string, expand string, top *int32) (result resources.ListResultIterator, err error)
	MoveResources(ctx context.Context, sourceResourceGroupName string, parameters resources.MoveInfo) (result resources.MoveResourcesFuture, err error)
	Update(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, APIVersion string, parameters resources.GenericResource) (result resources.UpdateFuture, err error)
	UpdateByID(ctx context.Context, resourceID string, APIVersion string, parameters resources.GenericResource) (result resources.UpdateByIDFuture, err error)
	ValidateMoveResources(ctx context.Context, sourceResourceGroupName string, parameters resources.MoveInfo) (result resources.ValidateMoveResourcesFuture, err error)
}

var _ ClientAPI = (*resources.Client)(nil)

// GroupsClientAPI contains the set of methods on the GroupsClient type.
type GroupsClientAPI interface {
	CheckExistence(ctx context.Context, resourceGroupName string) (result autorest.Response, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, parameters resources.Group) (result resources.Group, err error)
	Delete(ctx context.Context, resourceGroupName string, forceDeletionResourceTypes string) (result resources.GroupsDeleteFuture, err error)
	ExportTemplate(ctx context.Context, resourceGroupName string, parameters resources.ExportTemplateRequest) (result resources.GroupsExportTemplateFuture, err error)
	Get(ctx context.Context, resourceGroupName string) (result resources.Group, err error)
	List(ctx context.Context, filter string, top *int32) (result resources.GroupListResultPage, err error)
	ListComplete(ctx context.Context, filter string, top *int32) (result resources.GroupListResultIterator, err error)
	Update(ctx context.Context, resourceGroupName string, parameters resources.GroupPatchable) (result resources.Group, err error)
}

var _ GroupsClientAPI = (*resources.GroupsClient)(nil)

// TagsClientAPI contains the set of methods on the TagsClient type.
type TagsClientAPI interface {
	CreateOrUpdate(ctx context.Context, tagName string) (result resources.TagDetails, err error)
	CreateOrUpdateAtScope(ctx context.Context, scope string, parameters resources.TagsResource) (result resources.TagsResource, err error)
	CreateOrUpdateValue(ctx context.Context, tagName string, tagValue string) (result resources.TagValue, err error)
	Delete(ctx context.Context, tagName string) (result autorest.Response, err error)
	DeleteAtScope(ctx context.Context, scope string) (result autorest.Response, err error)
	DeleteValue(ctx context.Context, tagName string, tagValue string) (result autorest.Response, err error)
	GetAtScope(ctx context.Context, scope string) (result resources.TagsResource, err error)
	List(ctx context.Context) (result resources.TagsListResultPage, err error)
	ListComplete(ctx context.Context) (result resources.TagsListResultIterator, err error)
	UpdateAtScope(ctx context.Context, scope string, parameters resources.TagsPatchResource) (result resources.TagsResource, err error)
}

var _ TagsClientAPI = (*resources.TagsClient)(nil)

// DeploymentOperationsClientAPI contains the set of methods on the DeploymentOperationsClient type.
type DeploymentOperationsClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, deploymentName string, operationID string) (result resources.DeploymentOperation, err error)
	GetAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, operationID string) (result resources.DeploymentOperation, err error)
	GetAtScope(ctx context.Context, scope string, deploymentName string, operationID string) (result resources.DeploymentOperation, err error)
	GetAtSubscriptionScope(ctx context.Context, deploymentName string, operationID string) (result resources.DeploymentOperation, err error)
	GetAtTenantScope(ctx context.Context, deploymentName string, operationID string) (result resources.DeploymentOperation, err error)
	List(ctx context.Context, resourceGroupName string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultIterator, err error)
	ListAtManagementGroupScope(ctx context.Context, groupID string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultPage, err error)
	ListAtManagementGroupScopeComplete(ctx context.Context, groupID string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultIterator, err error)
	ListAtScope(ctx context.Context, scope string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultPage, err error)
	ListAtScopeComplete(ctx context.Context, scope string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultIterator, err error)
	ListAtSubscriptionScope(ctx context.Context, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultPage, err error)
	ListAtSubscriptionScopeComplete(ctx context.Context, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultIterator, err error)
	ListAtTenantScope(ctx context.Context, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultPage, err error)
	ListAtTenantScopeComplete(ctx context.Context, deploymentName string, top *int32) (result resources.DeploymentOperationsListResultIterator, err error)
}

var _ DeploymentOperationsClientAPI = (*resources.DeploymentOperationsClient)(nil)

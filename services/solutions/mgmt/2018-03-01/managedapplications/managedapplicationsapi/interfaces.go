// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armmanagedapplications](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armmanagedapplications). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package managedapplicationsapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/solutions/mgmt/2018-03-01/managedapplications"
	"github.com/Azure/go-autorest/autorest"
)

// BaseClientAPI contains the set of methods on the BaseClient type.
type BaseClientAPI interface {
	ListOperations(ctx context.Context) (result managedapplications.OperationListResultPage, err error)
	ListOperationsComplete(ctx context.Context) (result managedapplications.OperationListResultIterator, err error)
}

var _ BaseClientAPI = (*managedapplications.BaseClient)(nil)

// ApplicationsClientAPI contains the set of methods on the ApplicationsClient type.
type ApplicationsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, applicationName string, parameters managedapplications.Application) (result managedapplications.ApplicationsCreateOrUpdateFuture, err error)
	CreateOrUpdateByID(ctx context.Context, applicationID string, parameters managedapplications.Application) (result managedapplications.ApplicationsCreateOrUpdateByIDFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, applicationName string) (result managedapplications.ApplicationsDeleteFuture, err error)
	DeleteByID(ctx context.Context, applicationID string) (result managedapplications.ApplicationsDeleteByIDFuture, err error)
	Get(ctx context.Context, resourceGroupName string, applicationName string) (result managedapplications.Application, err error)
	GetByID(ctx context.Context, applicationID string) (result managedapplications.Application, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result managedapplications.ApplicationListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result managedapplications.ApplicationListResultIterator, err error)
	ListBySubscription(ctx context.Context) (result managedapplications.ApplicationListResultPage, err error)
	ListBySubscriptionComplete(ctx context.Context) (result managedapplications.ApplicationListResultIterator, err error)
	ListTokens(ctx context.Context, resourceGroupName string, applicationName string, parameters managedapplications.ListTokenRequest) (result managedapplications.ManagedIdentityTokenResult, err error)
	RefreshPermissions(ctx context.Context, resourceGroupName string, applicationName string) (result managedapplications.ApplicationsRefreshPermissionsFuture, err error)
	Update(ctx context.Context, resourceGroupName string, applicationName string, parameters *managedapplications.ApplicationPatchable) (result managedapplications.ApplicationsUpdateFuture, err error)
	UpdateAccess(ctx context.Context, resourceGroupName string, applicationName string, parameters managedapplications.UpdateAccessDefinition) (result managedapplications.ApplicationsUpdateAccessFuture, err error)
	UpdateByID(ctx context.Context, applicationID string, parameters *managedapplications.ApplicationPatchable) (result managedapplications.ApplicationsUpdateByIDFuture, err error)
}

var _ ApplicationsClientAPI = (*managedapplications.ApplicationsClient)(nil)

// ApplicationDefinitionsClientAPI contains the set of methods on the ApplicationDefinitionsClient type.
type ApplicationDefinitionsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, applicationDefinitionName string, parameters managedapplications.ApplicationDefinition) (result managedapplications.ApplicationDefinition, err error)
	CreateOrUpdateByID(ctx context.Context, resourceGroupName string, applicationDefinitionName string, parameters managedapplications.ApplicationDefinition) (result managedapplications.ApplicationDefinition, err error)
	Delete(ctx context.Context, resourceGroupName string, applicationDefinitionName string) (result autorest.Response, err error)
	DeleteByID(ctx context.Context, resourceGroupName string, applicationDefinitionName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, applicationDefinitionName string) (result managedapplications.ApplicationDefinition, err error)
	GetByID(ctx context.Context, resourceGroupName string, applicationDefinitionName string) (result managedapplications.ApplicationDefinition, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result managedapplications.ApplicationDefinitionListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result managedapplications.ApplicationDefinitionListResultIterator, err error)
	ListBySubscription(ctx context.Context) (result managedapplications.ApplicationDefinitionListResultPage, err error)
	ListBySubscriptionComplete(ctx context.Context) (result managedapplications.ApplicationDefinitionListResultIterator, err error)
	Update(ctx context.Context, resourceGroupName string, applicationDefinitionName string, parameters managedapplications.ApplicationDefinitionPatchable) (result managedapplications.ApplicationDefinition, err error)
	UpdateByID(ctx context.Context, resourceGroupName string, applicationDefinitionName string, parameters managedapplications.ApplicationDefinitionPatchable) (result managedapplications.ApplicationDefinition, err error)
}

var _ ApplicationDefinitionsClientAPI = (*managedapplications.ApplicationDefinitionsClient)(nil)

// JitRequestsClientAPI contains the set of methods on the JitRequestsClient type.
type JitRequestsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, jitRequestName string, parameters managedapplications.JitRequestDefinition) (result managedapplications.JitRequestsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, jitRequestName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, jitRequestName string) (result managedapplications.JitRequestDefinition, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result managedapplications.JitRequestDefinitionListResult, err error)
	ListBySubscription(ctx context.Context) (result managedapplications.JitRequestDefinitionListResult, err error)
	Update(ctx context.Context, resourceGroupName string, jitRequestName string, parameters managedapplications.JitRequestPatchable) (result managedapplications.JitRequestDefinition, err error)
}

var _ JitRequestsClientAPI = (*managedapplications.JitRequestsClient)(nil)

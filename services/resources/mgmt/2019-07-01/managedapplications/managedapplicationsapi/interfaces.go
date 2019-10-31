package managedapplicationsapi

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/managedapplications"
	"github.com/Azure/go-autorest/autorest"
)

// BaseClientAPI contains the set of methods on the BaseClient type.
type BaseClientAPI interface {
	ListOperations(ctx context.Context) (result managedapplications.OperationListResultPage, err error)
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
	ListBySubscription(ctx context.Context) (result managedapplications.ApplicationListResultPage, err error)
	Update(ctx context.Context, resourceGroupName string, applicationName string, parameters *managedapplications.Application) (result managedapplications.Application, err error)
	UpdateByID(ctx context.Context, applicationID string, parameters *managedapplications.Application) (result managedapplications.Application, err error)
}

var _ ApplicationsClientAPI = (*managedapplications.ApplicationsClient)(nil)

// ApplicationDefinitionsClientAPI contains the set of methods on the ApplicationDefinitionsClient type.
type ApplicationDefinitionsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, applicationDefinitionName string, parameters managedapplications.ApplicationDefinition) (result managedapplications.ApplicationDefinitionsCreateOrUpdateFuture, err error)
	CreateOrUpdateByID(ctx context.Context, applicationDefinitionID string, parameters managedapplications.ApplicationDefinition) (result managedapplications.ApplicationDefinitionsCreateOrUpdateByIDFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, applicationDefinitionName string) (result managedapplications.ApplicationDefinitionsDeleteFuture, err error)
	DeleteByID(ctx context.Context, applicationDefinitionID string) (result managedapplications.ApplicationDefinitionsDeleteByIDFuture, err error)
	Get(ctx context.Context, resourceGroupName string, applicationDefinitionName string) (result managedapplications.ApplicationDefinition, err error)
	GetByID(ctx context.Context, applicationDefinitionID string) (result managedapplications.ApplicationDefinition, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result managedapplications.ApplicationDefinitionListResultPage, err error)
}

var _ ApplicationDefinitionsClientAPI = (*managedapplications.ApplicationDefinitionsClient)(nil)

// ApplicationDefinitionArtifactsClientAPI contains the set of methods on the ApplicationDefinitionArtifactsClient type.
type ApplicationDefinitionArtifactsClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, applicationDefinitionName string, applicationArtifactName string) (result managedapplications.String, err error)
}

var _ ApplicationDefinitionArtifactsClientAPI = (*managedapplications.ApplicationDefinitionArtifactsClient)(nil)

// UpdateAccessClientAPI contains the set of methods on the UpdateAccessClient type.
type UpdateAccessClientAPI interface {
	Post(ctx context.Context, resourceGroupName string, applicationName string, parameters managedapplications.JitUpdateAccessDefinition) (result managedapplications.UpdateAccessPostFuture, err error)
}

var _ UpdateAccessClientAPI = (*managedapplications.UpdateAccessClient)(nil)

// JitRequestsClientAPI contains the set of methods on the JitRequestsClient type.
type JitRequestsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, jitRequestName string, parameters managedapplications.JitRequestDefinition) (result managedapplications.JitRequestsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, jitRequestName string) (result autorest.Response, err error)
	GetByName(ctx context.Context, resourceGroupName string, jitRequestName string) (result managedapplications.JitRequestDefinition, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result managedapplications.JitRequestDefinitionListResult, err error)
	ListBySubscription(ctx context.Context) (result managedapplications.JitRequestDefinitionListResult, err error)
	Patch(ctx context.Context, resourceGroupName string, jitRequestName string, parameters managedapplications.JitRequestPatchable) (result managedapplications.JitRequestDefinition, err error)
}

var _ JitRequestsClientAPI = (*managedapplications.JitRequestsClient)(nil)

// RefreshPermissionsClientAPI contains the set of methods on the RefreshPermissionsClient type.
type RefreshPermissionsClientAPI interface {
	Post(ctx context.Context, resourceGroupName string, applicationName string) (result managedapplications.RefreshPermissionsPostFuture, err error)
}

var _ RefreshPermissionsClientAPI = (*managedapplications.RefreshPermissionsClient)(nil)

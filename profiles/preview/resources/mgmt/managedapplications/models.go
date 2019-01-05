// +build go1.9

// Copyright 2019 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package managedapplications

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-06-01/managedapplications"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type ApplicationArtifactType = original.ApplicationArtifactType

const (
	Custom   ApplicationArtifactType = original.Custom
	Template ApplicationArtifactType = original.Template
)

type ApplicationLockLevel = original.ApplicationLockLevel

const (
	CanNotDelete ApplicationLockLevel = original.CanNotDelete
	None         ApplicationLockLevel = original.None
	ReadOnly     ApplicationLockLevel = original.ReadOnly
)

type ProvisioningState = original.ProvisioningState

const (
	Accepted  ProvisioningState = original.Accepted
	Canceled  ProvisioningState = original.Canceled
	Created   ProvisioningState = original.Created
	Creating  ProvisioningState = original.Creating
	Deleted   ProvisioningState = original.Deleted
	Deleting  ProvisioningState = original.Deleting
	Failed    ProvisioningState = original.Failed
	Ready     ProvisioningState = original.Ready
	Running   ProvisioningState = original.Running
	Succeeded ProvisioningState = original.Succeeded
	Updating  ProvisioningState = original.Updating
)

type ResourceIdentityType = original.ResourceIdentityType

const (
	SystemAssigned ResourceIdentityType = original.SystemAssigned
)

type Application = original.Application
type ApplicationArtifact = original.ApplicationArtifact
type ApplicationDefinition = original.ApplicationDefinition
type ApplicationDefinitionListResult = original.ApplicationDefinitionListResult
type ApplicationDefinitionListResultIterator = original.ApplicationDefinitionListResultIterator
type ApplicationDefinitionListResultPage = original.ApplicationDefinitionListResultPage
type ApplicationDefinitionProperties = original.ApplicationDefinitionProperties
type ApplicationDefinitionsClient = original.ApplicationDefinitionsClient
type ApplicationDefinitionsCreateOrUpdateByIDFuture = original.ApplicationDefinitionsCreateOrUpdateByIDFuture
type ApplicationDefinitionsCreateOrUpdateFuture = original.ApplicationDefinitionsCreateOrUpdateFuture
type ApplicationDefinitionsDeleteByIDFuture = original.ApplicationDefinitionsDeleteByIDFuture
type ApplicationDefinitionsDeleteFuture = original.ApplicationDefinitionsDeleteFuture
type ApplicationListResult = original.ApplicationListResult
type ApplicationListResultIterator = original.ApplicationListResultIterator
type ApplicationListResultPage = original.ApplicationListResultPage
type ApplicationPatchable = original.ApplicationPatchable
type ApplicationProperties = original.ApplicationProperties
type ApplicationPropertiesPatchable = original.ApplicationPropertiesPatchable
type ApplicationProviderAuthorization = original.ApplicationProviderAuthorization
type ApplicationsClient = original.ApplicationsClient
type ApplicationsCreateOrUpdateByIDFuture = original.ApplicationsCreateOrUpdateByIDFuture
type ApplicationsCreateOrUpdateFuture = original.ApplicationsCreateOrUpdateFuture
type ApplicationsDeleteByIDFuture = original.ApplicationsDeleteByIDFuture
type ApplicationsDeleteFuture = original.ApplicationsDeleteFuture
type BaseClient = original.BaseClient
type ErrorResponse = original.ErrorResponse
type GenericResource = original.GenericResource
type Identity = original.Identity
type Plan = original.Plan
type PlanPatchable = original.PlanPatchable
type Resource = original.Resource
type Sku = original.Sku

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewApplicationDefinitionListResultIterator(page ApplicationDefinitionListResultPage) ApplicationDefinitionListResultIterator {
	return original.NewApplicationDefinitionListResultIterator(page)
}
func NewApplicationDefinitionListResultPage(getNextPage func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error)) ApplicationDefinitionListResultPage {
	return original.NewApplicationDefinitionListResultPage(getNextPage)
}
func NewApplicationDefinitionsClient(subscriptionID string) ApplicationDefinitionsClient {
	return original.NewApplicationDefinitionsClient(subscriptionID)
}
func NewApplicationDefinitionsClientWithBaseURI(baseURI string, subscriptionID string) ApplicationDefinitionsClient {
	return original.NewApplicationDefinitionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewApplicationListResultIterator(page ApplicationListResultPage) ApplicationListResultIterator {
	return original.NewApplicationListResultIterator(page)
}
func NewApplicationListResultPage(getNextPage func(context.Context, ApplicationListResult) (ApplicationListResult, error)) ApplicationListResultPage {
	return original.NewApplicationListResultPage(getNextPage)
}
func NewApplicationsClient(subscriptionID string) ApplicationsClient {
	return original.NewApplicationsClient(subscriptionID)
}
func NewApplicationsClientWithBaseURI(baseURI string, subscriptionID string) ApplicationsClient {
	return original.NewApplicationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleApplicationArtifactTypeValues() []ApplicationArtifactType {
	return original.PossibleApplicationArtifactTypeValues()
}
func PossibleApplicationLockLevelValues() []ApplicationLockLevel {
	return original.PossibleApplicationLockLevelValues()
}
func PossibleProvisioningStateValues() []ProvisioningState {
	return original.PossibleProvisioningStateValues()
}
func PossibleResourceIdentityTypeValues() []ResourceIdentityType {
	return original.PossibleResourceIdentityTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

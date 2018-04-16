// +build go1.9

// Copyright 2018 Microsoft Corporation
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

package managementgroups

import original "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-01-01-preview/management"

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type OperationsClient = original.OperationsClient
type Permissions = original.Permissions

const (
	Delete   Permissions = original.Delete
	Edit     Permissions = original.Edit
	Noaccess Permissions = original.Noaccess
	View     Permissions = original.View
)

type Permissions1 = original.Permissions1

const (
	Permissions1Delete   Permissions1 = original.Permissions1Delete
	Permissions1Edit     Permissions1 = original.Permissions1Edit
	Permissions1Noaccess Permissions1 = original.Permissions1Noaccess
	Permissions1View     Permissions1 = original.Permissions1View
)

type ProvisioningState = original.ProvisioningState

const (
	Updating ProvisioningState = original.Updating
)

type Reason = original.Reason

const (
	AlreadyExists Reason = original.AlreadyExists
	Invalid       Reason = original.Invalid
)

type Type = original.Type

const (
	ProvidersMicrosoftManagementmanagementGroup Type = original.ProvidersMicrosoftManagementmanagementGroup
)

type Type1 = original.Type1

const (
	ProvidersMicrosoftManagementmanagementGroups Type1 = original.ProvidersMicrosoftManagementmanagementGroups
	Subscriptions                                Type1 = original.Subscriptions
)

type Type2 = original.Type2

const (
	Type2ProvidersMicrosoftManagementmanagementGroups Type2 = original.Type2ProvidersMicrosoftManagementmanagementGroups
	Type2Subscriptions                                Type2 = original.Type2Subscriptions
)

type CheckNameAvailabilityRequest = original.CheckNameAvailabilityRequest
type CheckNameAvailabilityResult = original.CheckNameAvailabilityResult
type ChildInfo = original.ChildInfo
type CreateManagementGroupChildInfo = original.CreateManagementGroupChildInfo
type CreateManagementGroupDetails = original.CreateManagementGroupDetails
type CreateManagementGroupProperties = original.CreateManagementGroupProperties
type CreateManagementGroupRequest = original.CreateManagementGroupRequest
type CreateOrUpdateFuture = original.CreateOrUpdateFuture
type CreateParentGroupInfo = original.CreateParentGroupInfo
type DeleteFuture = original.DeleteFuture
type Details = original.Details
type EntityHierarchyItem = original.EntityHierarchyItem
type EntityHierarchyItemProperties = original.EntityHierarchyItemProperties
type EntityInfo = original.EntityInfo
type EntityInfoProperties = original.EntityInfoProperties
type EntityListResult = original.EntityListResult
type EntityListResultIterator = original.EntityListResultIterator
type EntityListResultPage = original.EntityListResultPage
type EntityParentGroupInfo = original.EntityParentGroupInfo
type ErrorDetails = original.ErrorDetails
type ErrorResponse = original.ErrorResponse
type Info = original.Info
type InfoProperties = original.InfoProperties
type ListResult = original.ListResult
type ListResultIterator = original.ListResultIterator
type ListResultPage = original.ListResultPage
type ManagementGroup = original.ManagementGroup
type Operation = original.Operation
type OperationDisplayProperties = original.OperationDisplayProperties
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationResults = original.OperationResults
type OperationResultsProperties = original.OperationResultsProperties
type ParentGroupInfo = original.ParentGroupInfo
type PatchManagementGroupRequest = original.PatchManagementGroupRequest
type Properties = original.Properties
type SetObject = original.SetObject
type SubscriptionsClient = original.SubscriptionsClient
type EntitiesClient = original.EntitiesClient
type Client = original.Client

func PossiblePermissionsValues() []Permissions {
	return original.PossiblePermissionsValues()
}
func PossiblePermissions1Values() []Permissions1 {
	return original.PossiblePermissions1Values()
}
func PossibleProvisioningStateValues() []ProvisioningState {
	return original.PossibleProvisioningStateValues()
}
func PossibleReasonValues() []Reason {
	return original.PossibleReasonValues()
}
func PossibleTypeValues() []Type {
	return original.PossibleTypeValues()
}
func PossibleType1Values() []Type1 {
	return original.PossibleType1Values()
}
func PossibleType2Values() []Type2 {
	return original.PossibleType2Values()
}
func NewSubscriptionsClient(operationResultID string, skiptoken string) SubscriptionsClient {
	return original.NewSubscriptionsClient(operationResultID, skiptoken)
}
func NewSubscriptionsClientWithBaseURI(baseURI string, operationResultID string, skiptoken string) SubscriptionsClient {
	return original.NewSubscriptionsClientWithBaseURI(baseURI, operationResultID, skiptoken)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
func NewEntitiesClient(operationResultID string, skiptoken string) EntitiesClient {
	return original.NewEntitiesClient(operationResultID, skiptoken)
}
func NewEntitiesClientWithBaseURI(baseURI string, operationResultID string, skiptoken string) EntitiesClient {
	return original.NewEntitiesClientWithBaseURI(baseURI, operationResultID, skiptoken)
}
func NewClient(operationResultID string, skiptoken string) Client {
	return original.NewClient(operationResultID, skiptoken)
}
func NewClientWithBaseURI(baseURI string, operationResultID string, skiptoken string) Client {
	return original.NewClientWithBaseURI(baseURI, operationResultID, skiptoken)
}
func New(operationResultID string, skiptoken string) BaseClient {
	return original.New(operationResultID, skiptoken)
}
func NewWithBaseURI(baseURI string, operationResultID string, skiptoken string) BaseClient {
	return original.NewWithBaseURI(baseURI, operationResultID, skiptoken)
}
func NewOperationsClient(operationResultID string, skiptoken string) OperationsClient {
	return original.NewOperationsClient(operationResultID, skiptoken)
}
func NewOperationsClientWithBaseURI(baseURI string, operationResultID string, skiptoken string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, operationResultID, skiptoken)
}

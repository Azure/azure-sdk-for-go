//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package kubernetesconfiguration

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/kubernetesconfiguration/mgmt/2021-09-01/kubernetesconfiguration"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type CreatedByType = original.CreatedByType

const (
	CreatedByTypeApplication     CreatedByType = original.CreatedByTypeApplication
	CreatedByTypeKey             CreatedByType = original.CreatedByTypeKey
	CreatedByTypeManagedIdentity CreatedByType = original.CreatedByTypeManagedIdentity
	CreatedByTypeUser            CreatedByType = original.CreatedByTypeUser
)

type LevelType = original.LevelType

const (
	LevelTypeError       LevelType = original.LevelTypeError
	LevelTypeInformation LevelType = original.LevelTypeInformation
	LevelTypeWarning     LevelType = original.LevelTypeWarning
)

type ProvisioningState = original.ProvisioningState

const (
	ProvisioningStateCanceled  ProvisioningState = original.ProvisioningStateCanceled
	ProvisioningStateCreating  ProvisioningState = original.ProvisioningStateCreating
	ProvisioningStateDeleting  ProvisioningState = original.ProvisioningStateDeleting
	ProvisioningStateFailed    ProvisioningState = original.ProvisioningStateFailed
	ProvisioningStateSucceeded ProvisioningState = original.ProvisioningStateSucceeded
	ProvisioningStateUpdating  ProvisioningState = original.ProvisioningStateUpdating
)

type ResourceIdentityType = original.ResourceIdentityType

const (
	ResourceIdentityTypeSystemAssigned ResourceIdentityType = original.ResourceIdentityTypeSystemAssigned
)

type SkuTier = original.SkuTier

const (
	SkuTierBasic    SkuTier = original.SkuTierBasic
	SkuTierFree     SkuTier = original.SkuTierFree
	SkuTierPremium  SkuTier = original.SkuTierPremium
	SkuTierStandard SkuTier = original.SkuTierStandard
)

type AzureEntityResource = original.AzureEntityResource
type BaseClient = original.BaseClient
type ErrorAdditionalInfo = original.ErrorAdditionalInfo
type ErrorDetail = original.ErrorDetail
type ErrorResponse = original.ErrorResponse
type Extension = original.Extension
type ExtensionProperties = original.ExtensionProperties
type ExtensionPropertiesAksAssignedIdentity = original.ExtensionPropertiesAksAssignedIdentity
type ExtensionStatus = original.ExtensionStatus
type ExtensionsClient = original.ExtensionsClient
type ExtensionsCreateFuture = original.ExtensionsCreateFuture
type ExtensionsDeleteFuture = original.ExtensionsDeleteFuture
type ExtensionsList = original.ExtensionsList
type ExtensionsListIterator = original.ExtensionsListIterator
type ExtensionsListPage = original.ExtensionsListPage
type ExtensionsUpdateFuture = original.ExtensionsUpdateFuture
type Identity = original.Identity
type OperationStatusClient = original.OperationStatusClient
type OperationStatusList = original.OperationStatusList
type OperationStatusListIterator = original.OperationStatusListIterator
type OperationStatusListPage = original.OperationStatusListPage
type OperationStatusResult = original.OperationStatusResult
type OperationsClient = original.OperationsClient
type PatchExtension = original.PatchExtension
type PatchExtensionProperties = original.PatchExtensionProperties
type Plan = original.Plan
type ProxyResource = original.ProxyResource
type Resource = original.Resource
type ResourceModelWithAllowedPropertySet = original.ResourceModelWithAllowedPropertySet
type ResourceModelWithAllowedPropertySetIdentity = original.ResourceModelWithAllowedPropertySetIdentity
type ResourceModelWithAllowedPropertySetPlan = original.ResourceModelWithAllowedPropertySetPlan
type ResourceModelWithAllowedPropertySetSku = original.ResourceModelWithAllowedPropertySetSku
type ResourceProviderOperation = original.ResourceProviderOperation
type ResourceProviderOperationDisplay = original.ResourceProviderOperationDisplay
type ResourceProviderOperationList = original.ResourceProviderOperationList
type ResourceProviderOperationListIterator = original.ResourceProviderOperationListIterator
type ResourceProviderOperationListPage = original.ResourceProviderOperationListPage
type Scope = original.Scope
type ScopeCluster = original.ScopeCluster
type ScopeNamespace = original.ScopeNamespace
type Sku = original.Sku
type SystemData = original.SystemData
type TrackedResource = original.TrackedResource

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewExtensionsClient(subscriptionID string) ExtensionsClient {
	return original.NewExtensionsClient(subscriptionID)
}
func NewExtensionsClientWithBaseURI(baseURI string, subscriptionID string) ExtensionsClient {
	return original.NewExtensionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewExtensionsListIterator(page ExtensionsListPage) ExtensionsListIterator {
	return original.NewExtensionsListIterator(page)
}
func NewExtensionsListPage(cur ExtensionsList, getNextPage func(context.Context, ExtensionsList) (ExtensionsList, error)) ExtensionsListPage {
	return original.NewExtensionsListPage(cur, getNextPage)
}
func NewOperationStatusClient(subscriptionID string) OperationStatusClient {
	return original.NewOperationStatusClient(subscriptionID)
}
func NewOperationStatusClientWithBaseURI(baseURI string, subscriptionID string) OperationStatusClient {
	return original.NewOperationStatusClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationStatusListIterator(page OperationStatusListPage) OperationStatusListIterator {
	return original.NewOperationStatusListIterator(page)
}
func NewOperationStatusListPage(cur OperationStatusList, getNextPage func(context.Context, OperationStatusList) (OperationStatusList, error)) OperationStatusListPage {
	return original.NewOperationStatusListPage(cur, getNextPage)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewResourceProviderOperationListIterator(page ResourceProviderOperationListPage) ResourceProviderOperationListIterator {
	return original.NewResourceProviderOperationListIterator(page)
}
func NewResourceProviderOperationListPage(cur ResourceProviderOperationList, getNextPage func(context.Context, ResourceProviderOperationList) (ResourceProviderOperationList, error)) ResourceProviderOperationListPage {
	return original.NewResourceProviderOperationListPage(cur, getNextPage)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleCreatedByTypeValues() []CreatedByType {
	return original.PossibleCreatedByTypeValues()
}
func PossibleLevelTypeValues() []LevelType {
	return original.PossibleLevelTypeValues()
}
func PossibleProvisioningStateValues() []ProvisioningState {
	return original.PossibleProvisioningStateValues()
}
func PossibleResourceIdentityTypeValues() []ResourceIdentityType {
	return original.PossibleResourceIdentityTypeValues()
}
func PossibleSkuTierValues() []SkuTier {
	return original.PossibleSkuTierValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

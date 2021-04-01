// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package kubernetesconfiguration

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/preview/kubernetesconfiguration/mgmt/2020-07-01-preview/kubernetesconfiguration"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type ComplianceStateType = original.ComplianceStateType

const (
	Compliant    ComplianceStateType = original.Compliant
	Failed       ComplianceStateType = original.Failed
	Installed    ComplianceStateType = original.Installed
	Noncompliant ComplianceStateType = original.Noncompliant
	Pending      ComplianceStateType = original.Pending
)

type InstallStateType = original.InstallStateType

const (
	InstallStateTypeFailed    InstallStateType = original.InstallStateTypeFailed
	InstallStateTypeInstalled InstallStateType = original.InstallStateTypeInstalled
	InstallStateTypePending   InstallStateType = original.InstallStateTypePending
)

type LevelType = original.LevelType

const (
	Error       LevelType = original.Error
	Information LevelType = original.Information
	Warning     LevelType = original.Warning
)

type MessageLevelType = original.MessageLevelType

const (
	MessageLevelTypeError       MessageLevelType = original.MessageLevelTypeError
	MessageLevelTypeInformation MessageLevelType = original.MessageLevelTypeInformation
	MessageLevelTypeWarning     MessageLevelType = original.MessageLevelTypeWarning
)

type OperatorScopeType = original.OperatorScopeType

const (
	Cluster   OperatorScopeType = original.Cluster
	Namespace OperatorScopeType = original.Namespace
)

type OperatorType = original.OperatorType

const (
	Flux OperatorType = original.Flux
)

type ProvisioningStateType = original.ProvisioningStateType

const (
	ProvisioningStateTypeAccepted  ProvisioningStateType = original.ProvisioningStateTypeAccepted
	ProvisioningStateTypeDeleting  ProvisioningStateType = original.ProvisioningStateTypeDeleting
	ProvisioningStateTypeFailed    ProvisioningStateType = original.ProvisioningStateTypeFailed
	ProvisioningStateTypeRunning   ProvisioningStateType = original.ProvisioningStateTypeRunning
	ProvisioningStateTypeSucceeded ProvisioningStateType = original.ProvisioningStateTypeSucceeded
)

type ResourceIdentityType = original.ResourceIdentityType

const (
	None           ResourceIdentityType = original.None
	SystemAssigned ResourceIdentityType = original.SystemAssigned
)

type BaseClient = original.BaseClient
type ComplianceStatus = original.ComplianceStatus
type ConfigurationIdentity = original.ConfigurationIdentity
type ErrorDefinition = original.ErrorDefinition
type ErrorResponse = original.ErrorResponse
type ExtensionInstance = original.ExtensionInstance
type ExtensionInstanceProperties = original.ExtensionInstanceProperties
type ExtensionInstanceUpdate = original.ExtensionInstanceUpdate
type ExtensionInstanceUpdateProperties = original.ExtensionInstanceUpdateProperties
type ExtensionInstancesList = original.ExtensionInstancesList
type ExtensionInstancesListIterator = original.ExtensionInstancesListIterator
type ExtensionInstancesListPage = original.ExtensionInstancesListPage
type ExtensionStatus = original.ExtensionStatus
type ExtensionsClient = original.ExtensionsClient
type HelmOperatorProperties = original.HelmOperatorProperties
type OperationsClient = original.OperationsClient
type ProxyResource = original.ProxyResource
type Resource = original.Resource
type ResourceProviderOperation = original.ResourceProviderOperation
type ResourceProviderOperationDisplay = original.ResourceProviderOperationDisplay
type ResourceProviderOperationList = original.ResourceProviderOperationList
type ResourceProviderOperationListIterator = original.ResourceProviderOperationListIterator
type ResourceProviderOperationListPage = original.ResourceProviderOperationListPage
type Result = original.Result
type Scope = original.Scope
type ScopeCluster = original.ScopeCluster
type ScopeNamespace = original.ScopeNamespace
type SourceControlConfiguration = original.SourceControlConfiguration
type SourceControlConfigurationList = original.SourceControlConfigurationList
type SourceControlConfigurationListIterator = original.SourceControlConfigurationListIterator
type SourceControlConfigurationListPage = original.SourceControlConfigurationListPage
type SourceControlConfigurationProperties = original.SourceControlConfigurationProperties
type SourceControlConfigurationsClient = original.SourceControlConfigurationsClient
type SourceControlConfigurationsDeleteFuture = original.SourceControlConfigurationsDeleteFuture
type SystemData = original.SystemData

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewExtensionInstancesListIterator(page ExtensionInstancesListPage) ExtensionInstancesListIterator {
	return original.NewExtensionInstancesListIterator(page)
}
func NewExtensionInstancesListPage(cur ExtensionInstancesList, getNextPage func(context.Context, ExtensionInstancesList) (ExtensionInstancesList, error)) ExtensionInstancesListPage {
	return original.NewExtensionInstancesListPage(cur, getNextPage)
}
func NewExtensionsClient(subscriptionID string) ExtensionsClient {
	return original.NewExtensionsClient(subscriptionID)
}
func NewExtensionsClientWithBaseURI(baseURI string, subscriptionID string) ExtensionsClient {
	return original.NewExtensionsClientWithBaseURI(baseURI, subscriptionID)
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
func NewSourceControlConfigurationListIterator(page SourceControlConfigurationListPage) SourceControlConfigurationListIterator {
	return original.NewSourceControlConfigurationListIterator(page)
}
func NewSourceControlConfigurationListPage(cur SourceControlConfigurationList, getNextPage func(context.Context, SourceControlConfigurationList) (SourceControlConfigurationList, error)) SourceControlConfigurationListPage {
	return original.NewSourceControlConfigurationListPage(cur, getNextPage)
}
func NewSourceControlConfigurationsClient(subscriptionID string) SourceControlConfigurationsClient {
	return original.NewSourceControlConfigurationsClient(subscriptionID)
}
func NewSourceControlConfigurationsClientWithBaseURI(baseURI string, subscriptionID string) SourceControlConfigurationsClient {
	return original.NewSourceControlConfigurationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleComplianceStateTypeValues() []ComplianceStateType {
	return original.PossibleComplianceStateTypeValues()
}
func PossibleInstallStateTypeValues() []InstallStateType {
	return original.PossibleInstallStateTypeValues()
}
func PossibleLevelTypeValues() []LevelType {
	return original.PossibleLevelTypeValues()
}
func PossibleMessageLevelTypeValues() []MessageLevelType {
	return original.PossibleMessageLevelTypeValues()
}
func PossibleOperatorScopeTypeValues() []OperatorScopeType {
	return original.PossibleOperatorScopeTypeValues()
}
func PossibleOperatorTypeValues() []OperatorType {
	return original.PossibleOperatorTypeValues()
}
func PossibleProvisioningStateTypeValues() []ProvisioningStateType {
	return original.PossibleProvisioningStateTypeValues()
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

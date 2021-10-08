//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package hybridcompute

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/hybridcompute/mgmt/2020-08-02/hybridcompute"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type InstanceViewTypes = original.InstanceViewTypes

const (
	InstanceView InstanceViewTypes = original.InstanceView
)

type StatusLevelTypes = original.StatusLevelTypes

const (
	Error   StatusLevelTypes = original.Error
	Info    StatusLevelTypes = original.Info
	Warning StatusLevelTypes = original.Warning
)

type StatusTypes = original.StatusTypes

const (
	StatusTypesConnected    StatusTypes = original.StatusTypesConnected
	StatusTypesDisconnected StatusTypes = original.StatusTypesDisconnected
	StatusTypesError        StatusTypes = original.StatusTypesError
)

type AzureEntityResource = original.AzureEntityResource
type BaseClient = original.BaseClient
type ErrorDetail = original.ErrorDetail
type ErrorResponse = original.ErrorResponse
type Identity = original.Identity
type LocationData = original.LocationData
type Machine = original.Machine
type MachineExtension = original.MachineExtension
type MachineExtensionInstanceView = original.MachineExtensionInstanceView
type MachineExtensionInstanceViewStatus = original.MachineExtensionInstanceViewStatus
type MachineExtensionProperties = original.MachineExtensionProperties
type MachineExtensionPropertiesInstanceView = original.MachineExtensionPropertiesInstanceView
type MachineExtensionPropertiesModel = original.MachineExtensionPropertiesModel
type MachineExtensionUpdate = original.MachineExtensionUpdate
type MachineExtensionUpdateProperties = original.MachineExtensionUpdateProperties
type MachineExtensionUpdatePropertiesModel = original.MachineExtensionUpdatePropertiesModel
type MachineExtensionsClient = original.MachineExtensionsClient
type MachineExtensionsCreateOrUpdateFuture = original.MachineExtensionsCreateOrUpdateFuture
type MachineExtensionsDeleteFuture = original.MachineExtensionsDeleteFuture
type MachineExtensionsListResult = original.MachineExtensionsListResult
type MachineExtensionsListResultIterator = original.MachineExtensionsListResultIterator
type MachineExtensionsListResultPage = original.MachineExtensionsListResultPage
type MachineExtensionsUpdateFuture = original.MachineExtensionsUpdateFuture
type MachineIdentity = original.MachineIdentity
type MachineListResult = original.MachineListResult
type MachineListResultIterator = original.MachineListResultIterator
type MachineListResultPage = original.MachineListResultPage
type MachineProperties = original.MachineProperties
type MachinePropertiesModel = original.MachinePropertiesModel
type MachinePropertiesOsProfile = original.MachinePropertiesOsProfile
type MachineUpdate = original.MachineUpdate
type MachineUpdateIdentity = original.MachineUpdateIdentity
type MachineUpdateProperties = original.MachineUpdateProperties
type MachineUpdatePropertiesModel = original.MachineUpdatePropertiesModel
type MachinesClient = original.MachinesClient
type OSProfile = original.OSProfile
type OperationListResult = original.OperationListResult
type OperationValue = original.OperationValue
type OperationValueDisplay = original.OperationValueDisplay
type OperationValueDisplayModel = original.OperationValueDisplayModel
type OperationsClient = original.OperationsClient
type ProxyResource = original.ProxyResource
type Resource = original.Resource
type TrackedResource = original.TrackedResource
type UpdateResource = original.UpdateResource

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewMachineExtensionsClient(subscriptionID string) MachineExtensionsClient {
	return original.NewMachineExtensionsClient(subscriptionID)
}
func NewMachineExtensionsClientWithBaseURI(baseURI string, subscriptionID string) MachineExtensionsClient {
	return original.NewMachineExtensionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewMachineExtensionsListResultIterator(page MachineExtensionsListResultPage) MachineExtensionsListResultIterator {
	return original.NewMachineExtensionsListResultIterator(page)
}
func NewMachineExtensionsListResultPage(cur MachineExtensionsListResult, getNextPage func(context.Context, MachineExtensionsListResult) (MachineExtensionsListResult, error)) MachineExtensionsListResultPage {
	return original.NewMachineExtensionsListResultPage(cur, getNextPage)
}
func NewMachineListResultIterator(page MachineListResultPage) MachineListResultIterator {
	return original.NewMachineListResultIterator(page)
}
func NewMachineListResultPage(cur MachineListResult, getNextPage func(context.Context, MachineListResult) (MachineListResult, error)) MachineListResultPage {
	return original.NewMachineListResultPage(cur, getNextPage)
}
func NewMachinesClient(subscriptionID string) MachinesClient {
	return original.NewMachinesClient(subscriptionID)
}
func NewMachinesClientWithBaseURI(baseURI string, subscriptionID string) MachinesClient {
	return original.NewMachinesClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleInstanceViewTypesValues() []InstanceViewTypes {
	return original.PossibleInstanceViewTypesValues()
}
func PossibleStatusLevelTypesValues() []StatusLevelTypes {
	return original.PossibleStatusLevelTypesValues()
}
func PossibleStatusTypesValues() []StatusTypes {
	return original.PossibleStatusTypesValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

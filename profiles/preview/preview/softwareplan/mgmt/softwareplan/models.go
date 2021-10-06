//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package softwareplan

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/preview/softwareplan/mgmt/2019-06-01-preview/softwareplan"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type ErrorCode = original.ErrorCode

const (
	InvalidRequestParameter ErrorCode = original.InvalidRequestParameter
	MissingRequestParameter ErrorCode = original.MissingRequestParameter
)

type ProvisioningState = original.ProvisioningState

const (
	Cancelled ProvisioningState = original.Cancelled
	Failed    ProvisioningState = original.Failed
	Succeeded ProvisioningState = original.Succeeded
)

type AzureEntityResource = original.AzureEntityResource
type BaseClient = original.BaseClient
type Client = original.Client
type Error = original.Error
type HybridUseBenefitClient = original.HybridUseBenefitClient
type HybridUseBenefitListResult = original.HybridUseBenefitListResult
type HybridUseBenefitListResultIterator = original.HybridUseBenefitListResultIterator
type HybridUseBenefitListResultPage = original.HybridUseBenefitListResultPage
type HybridUseBenefitModel = original.HybridUseBenefitModel
type HybridUseBenefitProperties = original.HybridUseBenefitProperties
type HybridUseBenefitRevisionClient = original.HybridUseBenefitRevisionClient
type OperationDisplay = original.OperationDisplay
type OperationList = original.OperationList
type OperationListIterator = original.OperationListIterator
type OperationListPage = original.OperationListPage
type OperationResponse = original.OperationResponse
type OperationsClient = original.OperationsClient
type ProxyResource = original.ProxyResource
type Resource = original.Resource
type Sku = original.Sku
type TrackedResource = original.TrackedResource

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewClient(subscriptionID string) Client {
	return original.NewClient(subscriptionID)
}
func NewClientWithBaseURI(baseURI string, subscriptionID string) Client {
	return original.NewClientWithBaseURI(baseURI, subscriptionID)
}
func NewHybridUseBenefitClient(subscriptionID string) HybridUseBenefitClient {
	return original.NewHybridUseBenefitClient(subscriptionID)
}
func NewHybridUseBenefitClientWithBaseURI(baseURI string, subscriptionID string) HybridUseBenefitClient {
	return original.NewHybridUseBenefitClientWithBaseURI(baseURI, subscriptionID)
}
func NewHybridUseBenefitListResultIterator(page HybridUseBenefitListResultPage) HybridUseBenefitListResultIterator {
	return original.NewHybridUseBenefitListResultIterator(page)
}
func NewHybridUseBenefitListResultPage(cur HybridUseBenefitListResult, getNextPage func(context.Context, HybridUseBenefitListResult) (HybridUseBenefitListResult, error)) HybridUseBenefitListResultPage {
	return original.NewHybridUseBenefitListResultPage(cur, getNextPage)
}
func NewHybridUseBenefitRevisionClient(subscriptionID string) HybridUseBenefitRevisionClient {
	return original.NewHybridUseBenefitRevisionClient(subscriptionID)
}
func NewHybridUseBenefitRevisionClientWithBaseURI(baseURI string, subscriptionID string) HybridUseBenefitRevisionClient {
	return original.NewHybridUseBenefitRevisionClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationListIterator(page OperationListPage) OperationListIterator {
	return original.NewOperationListIterator(page)
}
func NewOperationListPage(cur OperationList, getNextPage func(context.Context, OperationList) (OperationList, error)) OperationListPage {
	return original.NewOperationListPage(cur, getNextPage)
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
func PossibleErrorCodeValues() []ErrorCode {
	return original.PossibleErrorCodeValues()
}
func PossibleProvisioningStateValues() []ProvisioningState {
	return original.PossibleProvisioningStateValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

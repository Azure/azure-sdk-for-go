//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package datacatalog

import original "github.com/Azure/azure-sdk-for-go/services/datacatalog/mgmt/2016-03-30/datacatalog"

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type SkuType = original.SkuType

const (
	Free     SkuType = original.Free
	Standard SkuType = original.Standard
)

type ADCCatalog = original.ADCCatalog
type ADCCatalogProperties = original.ADCCatalogProperties
type ADCCatalogsClient = original.ADCCatalogsClient
type ADCCatalogsDeleteFuture = original.ADCCatalogsDeleteFuture
type ADCCatalogsListResult = original.ADCCatalogsListResult
type ADCOperationsClient = original.ADCOperationsClient
type BaseClient = original.BaseClient
type OperationDisplayInfo = original.OperationDisplayInfo
type OperationEntity = original.OperationEntity
type OperationEntityListResult = original.OperationEntityListResult
type Principals = original.Principals
type Resource = original.Resource

func New(subscriptionID string, catalogName string) BaseClient {
	return original.New(subscriptionID, catalogName)
}
func NewADCCatalogsClient(subscriptionID string, catalogName string) ADCCatalogsClient {
	return original.NewADCCatalogsClient(subscriptionID, catalogName)
}
func NewADCCatalogsClientWithBaseURI(baseURI string, subscriptionID string, catalogName string) ADCCatalogsClient {
	return original.NewADCCatalogsClientWithBaseURI(baseURI, subscriptionID, catalogName)
}
func NewADCOperationsClient(subscriptionID string, catalogName string) ADCOperationsClient {
	return original.NewADCOperationsClient(subscriptionID, catalogName)
}
func NewADCOperationsClientWithBaseURI(baseURI string, subscriptionID string, catalogName string) ADCOperationsClient {
	return original.NewADCOperationsClientWithBaseURI(baseURI, subscriptionID, catalogName)
}
func NewWithBaseURI(baseURI string, subscriptionID string, catalogName string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID, catalogName)
}
func PossibleSkuTypeValues() []SkuType {
	return original.PossibleSkuTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}

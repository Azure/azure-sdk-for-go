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

package iotcentral

import original "github.com/Azure/azure-sdk-for-go/services/iotcentral/mgmt/2018-09-01/iotcentral"

type AppsClient = original.AppsClient

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type AppNameUnavailabilityReason = original.AppNameUnavailabilityReason

const (
	AlreadyExists AppNameUnavailabilityReason = original.AlreadyExists
	Invalid       AppNameUnavailabilityReason = original.Invalid
)

type AppSku = original.AppSku

const (
	F1 AppSku = original.F1
	S1 AppSku = original.S1
)

type App = original.App
type AppAvailabilityInfo = original.AppAvailabilityInfo
type AppListResult = original.AppListResult
type AppListResultIterator = original.AppListResultIterator
type AppListResultPage = original.AppListResultPage
type AppPatch = original.AppPatch
type AppProperties = original.AppProperties
type AppsCreateOrUpdateFuture = original.AppsCreateOrUpdateFuture
type AppsDeleteFuture = original.AppsDeleteFuture
type AppSkuInfo = original.AppSkuInfo
type AppsUpdateFuture = original.AppsUpdateFuture
type ErrorDetails = original.ErrorDetails
type ErrorResponseBody = original.ErrorResponseBody
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationInputs = original.OperationInputs
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type Resource = original.Resource
type OperationsClient = original.OperationsClient

func NewAppsClient(subscriptionID string) AppsClient {
	return original.NewAppsClient(subscriptionID)
}
func NewAppsClientWithBaseURI(baseURI string, subscriptionID string) AppsClient {
	return original.NewAppsClientWithBaseURI(baseURI, subscriptionID)
}
func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleAppNameUnavailabilityReasonValues() []AppNameUnavailabilityReason {
	return original.PossibleAppNameUnavailabilityReasonValues()
}
func PossibleAppSkuValues() []AppSku {
	return original.PossibleAppSkuValues()
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

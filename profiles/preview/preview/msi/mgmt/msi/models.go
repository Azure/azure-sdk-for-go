// +build go1.9

// Copyright 2020 Microsoft Corporation
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

package msi

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type AzureEntityResource = original.AzureEntityResource
type BaseClient = original.BaseClient
type CloudError = original.CloudError
type CloudErrorBody = original.CloudErrorBody
type Identity = original.Identity
type IdentityProperties = original.IdentityProperties
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationsClient = original.OperationsClient
type ProxyResource = original.ProxyResource
type Resource = original.Resource
type SystemAssignedIdentitiesClient = original.SystemAssignedIdentitiesClient
type SystemAssignedIdentity = original.SystemAssignedIdentity
type TrackedResource = original.TrackedResource
type UserAssignedIdentitiesClient = original.UserAssignedIdentitiesClient
type UserAssignedIdentitiesListResult = original.UserAssignedIdentitiesListResult
type UserAssignedIdentitiesListResultIterator = original.UserAssignedIdentitiesListResultIterator
type UserAssignedIdentitiesListResultPage = original.UserAssignedIdentitiesListResultPage

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewOperationListResultIterator(page OperationListResultPage) OperationListResultIterator {
	return original.NewOperationListResultIterator(page)
}
func NewOperationListResultPage(cur OperationListResult, getNextPage func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage {
	return original.NewOperationListResultPage(cur, getNextPage)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewSystemAssignedIdentitiesClient(subscriptionID string) SystemAssignedIdentitiesClient {
	return original.NewSystemAssignedIdentitiesClient(subscriptionID)
}
func NewSystemAssignedIdentitiesClientWithBaseURI(baseURI string, subscriptionID string) SystemAssignedIdentitiesClient {
	return original.NewSystemAssignedIdentitiesClientWithBaseURI(baseURI, subscriptionID)
}
func NewUserAssignedIdentitiesClient(subscriptionID string) UserAssignedIdentitiesClient {
	return original.NewUserAssignedIdentitiesClient(subscriptionID)
}
func NewUserAssignedIdentitiesClientWithBaseURI(baseURI string, subscriptionID string) UserAssignedIdentitiesClient {
	return original.NewUserAssignedIdentitiesClientWithBaseURI(baseURI, subscriptionID)
}
func NewUserAssignedIdentitiesListResultIterator(page UserAssignedIdentitiesListResultPage) UserAssignedIdentitiesListResultIterator {
	return original.NewUserAssignedIdentitiesListResultIterator(page)
}
func NewUserAssignedIdentitiesListResultPage(cur UserAssignedIdentitiesListResult, getNextPage func(context.Context, UserAssignedIdentitiesListResult) (UserAssignedIdentitiesListResult, error)) UserAssignedIdentitiesListResultPage {
	return original.NewUserAssignedIdentitiesListResultPage(cur, getNextPage)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

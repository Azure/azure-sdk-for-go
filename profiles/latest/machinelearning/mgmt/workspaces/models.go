//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package workspaces

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/machinelearning/mgmt/2019-10-01/workspaces"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type State = original.State

const (
	Deleted      State = original.Deleted
	Disabled     State = original.Disabled
	Enabled      State = original.Enabled
	Migrated     State = original.Migrated
	Registered   State = original.Registered
	Unregistered State = original.Unregistered
	Updated      State = original.Updated
)

type WorkspaceType = original.WorkspaceType

const (
	Anonymous    WorkspaceType = original.Anonymous
	Free         WorkspaceType = original.Free
	PaidPremium  WorkspaceType = original.PaidPremium
	PaidStandard WorkspaceType = original.PaidStandard
	Production   WorkspaceType = original.Production
)

type BaseClient = original.BaseClient
type Client = original.Client
type ErrorResponse = original.ErrorResponse
type KeysResponse = original.KeysResponse
type ListResult = original.ListResult
type ListResultIterator = original.ListResultIterator
type ListResultPage = original.ListResultPage
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationsClient = original.OperationsClient
type Properties = original.Properties
type PropertiesUpdateParameters = original.PropertiesUpdateParameters
type Resource = original.Resource
type Sku = original.Sku
type UpdateParameters = original.UpdateParameters
type Workspace = original.Workspace

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewClient(subscriptionID string) Client {
	return original.NewClient(subscriptionID)
}
func NewClientWithBaseURI(baseURI string, subscriptionID string) Client {
	return original.NewClientWithBaseURI(baseURI, subscriptionID)
}
func NewListResultIterator(page ListResultPage) ListResultIterator {
	return original.NewListResultIterator(page)
}
func NewListResultPage(cur ListResult, getNextPage func(context.Context, ListResult) (ListResult, error)) ListResultPage {
	return original.NewListResultPage(cur, getNextPage)
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
func PossibleStateValues() []State {
	return original.PossibleStateValues()
}
func PossibleWorkspaceTypeValues() []WorkspaceType {
	return original.PossibleWorkspaceTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}

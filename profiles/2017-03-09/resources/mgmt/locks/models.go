//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package locks

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2015-01-01/locks"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type LockLevel = original.LockLevel

const (
	CanNotDelete LockLevel = original.CanNotDelete
	NotSpecified LockLevel = original.NotSpecified
	ReadOnly     LockLevel = original.ReadOnly
)

type BaseClient = original.BaseClient
type ManagementLockListResult = original.ManagementLockListResult
type ManagementLockListResultIterator = original.ManagementLockListResultIterator
type ManagementLockListResultPage = original.ManagementLockListResultPage
type ManagementLockObject = original.ManagementLockObject
type ManagementLockProperties = original.ManagementLockProperties
type ManagementLocksClient = original.ManagementLocksClient

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewManagementLockListResultIterator(page ManagementLockListResultPage) ManagementLockListResultIterator {
	return original.NewManagementLockListResultIterator(page)
}
func NewManagementLockListResultPage(cur ManagementLockListResult, getNextPage func(context.Context, ManagementLockListResult) (ManagementLockListResult, error)) ManagementLockListResultPage {
	return original.NewManagementLockListResultPage(cur, getNextPage)
}
func NewManagementLocksClient(subscriptionID string) ManagementLocksClient {
	return original.NewManagementLocksClient(subscriptionID)
}
func NewManagementLocksClientWithBaseURI(baseURI string, subscriptionID string) ManagementLocksClient {
	return original.NewManagementLocksClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleLockLevelValues() []LockLevel {
	return original.PossibleLockLevelValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/2017-03-09"
}
func Version() string {
	return original.Version()
}

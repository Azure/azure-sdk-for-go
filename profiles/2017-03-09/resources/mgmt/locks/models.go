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

package locks

import original "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2015-01-01/locks"

type ManagementLocksClient = original.ManagementLocksClient
type LockLevel = original.LockLevel

const (
	CanNotDelete LockLevel = original.CanNotDelete
	NotSpecified LockLevel = original.NotSpecified
	ReadOnly     LockLevel = original.ReadOnly
)

type ManagementLockListResult = original.ManagementLockListResult
type ManagementLockListResultIterator = original.ManagementLockListResultIterator
type ManagementLockListResultPage = original.ManagementLockListResultPage
type ManagementLockObject = original.ManagementLockObject
type ManagementLockProperties = original.ManagementLockProperties

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient

func NewManagementLocksClient(subscriptionID string) ManagementLocksClient {
	return original.NewManagementLocksClient(subscriptionID)
}
func NewManagementLocksClientWithBaseURI(baseURI string, subscriptionID string) ManagementLocksClient {
	return original.NewManagementLocksClientWithBaseURI(baseURI, subscriptionID)
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
func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}

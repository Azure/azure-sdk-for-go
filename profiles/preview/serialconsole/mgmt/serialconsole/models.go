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

package serialconsole

import original "github.com/Azure/azure-sdk-for-go/services/serialconsole/mgmt/2018-05-01/serialconsole"

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type DisableSerialConsoleResult = original.DisableSerialConsoleResult
type EnableSerialConsoleResult = original.EnableSerialConsoleResult
type GetSerialConsoleSubscriptionNotFound = original.GetSerialConsoleSubscriptionNotFound
type Operations = original.Operations
type OperationsValueItem = original.OperationsValueItem
type OperationsValueItemDisplay = original.OperationsValueItemDisplay
type SetObject = original.SetObject
type Status = original.Status

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
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

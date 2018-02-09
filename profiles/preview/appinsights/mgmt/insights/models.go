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

package insights

import original "github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"

type APIKeysClient = original.APIKeysClient

func NewAPIKeysClient(subscriptionID string) APIKeysClient {
	return original.NewAPIKeysClient(subscriptionID)
}
func NewAPIKeysClientWithBaseURI(baseURI string, subscriptionID string) APIKeysClient {
	return original.NewAPIKeysClientWithBaseURI(baseURI, subscriptionID)
}

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}

type ComponentCurrentBillingFeaturesClient = original.ComponentCurrentBillingFeaturesClient

func NewComponentCurrentBillingFeaturesClient(subscriptionID string) ComponentCurrentBillingFeaturesClient {
	return original.NewComponentCurrentBillingFeaturesClient(subscriptionID)
}
func NewComponentCurrentBillingFeaturesClientWithBaseURI(baseURI string, subscriptionID string) ComponentCurrentBillingFeaturesClient {
	return original.NewComponentCurrentBillingFeaturesClientWithBaseURI(baseURI, subscriptionID)
}

type ComponentQuotaStatusClient = original.ComponentQuotaStatusClient

func NewComponentQuotaStatusClient(subscriptionID string) ComponentQuotaStatusClient {
	return original.NewComponentQuotaStatusClient(subscriptionID)
}
func NewComponentQuotaStatusClientWithBaseURI(baseURI string, subscriptionID string) ComponentQuotaStatusClient {
	return original.NewComponentQuotaStatusClientWithBaseURI(baseURI, subscriptionID)
}

type ComponentsClient = original.ComponentsClient

func NewComponentsClient(subscriptionID string) ComponentsClient {
	return original.NewComponentsClient(subscriptionID)
}
func NewComponentsClientWithBaseURI(baseURI string, subscriptionID string) ComponentsClient {
	return original.NewComponentsClientWithBaseURI(baseURI, subscriptionID)
}

type ExportConfigurationsClient = original.ExportConfigurationsClient

func NewExportConfigurationsClient(subscriptionID string) ExportConfigurationsClient {
	return original.NewExportConfigurationsClient(subscriptionID)
}
func NewExportConfigurationsClientWithBaseURI(baseURI string, subscriptionID string) ExportConfigurationsClient {
	return original.NewExportConfigurationsClientWithBaseURI(baseURI, subscriptionID)
}

type ApplicationType = original.ApplicationType

const (
	Other ApplicationType = original.Other
	Web   ApplicationType = original.Web
)

type FlowType = original.FlowType

const (
	Bluefield FlowType = original.Bluefield
)

type RequestSource = original.RequestSource

const (
	Rest RequestSource = original.Rest
)

type WebTestKind = original.WebTestKind

const (
	Multistep WebTestKind = original.Multistep
	Ping      WebTestKind = original.Ping
)

type APIKeyRequest = original.APIKeyRequest
type ApplicationInsightsComponent = original.ApplicationInsightsComponent
type ApplicationInsightsComponentAPIKey = original.ApplicationInsightsComponentAPIKey
type ApplicationInsightsComponentAPIKeyListResult = original.ApplicationInsightsComponentAPIKeyListResult
type ApplicationInsightsComponentBillingFeatures = original.ApplicationInsightsComponentBillingFeatures
type ApplicationInsightsComponentDataVolumeCap = original.ApplicationInsightsComponentDataVolumeCap
type ApplicationInsightsComponentExportConfiguration = original.ApplicationInsightsComponentExportConfiguration
type ApplicationInsightsComponentExportRequest = original.ApplicationInsightsComponentExportRequest
type ApplicationInsightsComponentListResult = original.ApplicationInsightsComponentListResult
type ApplicationInsightsComponentListResultIterator = original.ApplicationInsightsComponentListResultIterator
type ApplicationInsightsComponentListResultPage = original.ApplicationInsightsComponentListResultPage
type ApplicationInsightsComponentProperties = original.ApplicationInsightsComponentProperties
type ApplicationInsightsComponentQuotaStatus = original.ApplicationInsightsComponentQuotaStatus
type ErrorResponse = original.ErrorResponse
type ListApplicationInsightsComponentExportConfiguration = original.ListApplicationInsightsComponentExportConfiguration
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type Resource = original.Resource
type TagsResource = original.TagsResource
type WebTest = original.WebTest
type WebTestGeolocation = original.WebTestGeolocation
type WebTestListResult = original.WebTestListResult
type WebTestListResultIterator = original.WebTestListResultIterator
type WebTestListResultPage = original.WebTestListResultPage
type WebTestProperties = original.WebTestProperties
type WebTestPropertiesConfiguration = original.WebTestPropertiesConfiguration
type OperationsClient = original.OperationsClient

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

type WebTestsClient = original.WebTestsClient

func NewWebTestsClient(subscriptionID string) WebTestsClient {
	return original.NewWebTestsClient(subscriptionID)
}
func NewWebTestsClientWithBaseURI(baseURI string, subscriptionID string) WebTestsClient {
	return original.NewWebTestsClientWithBaseURI(baseURI, subscriptionID)
}

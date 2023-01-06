// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package insightsapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/appinsights/mgmt/2020-02-02-preview/insights"
	"github.com/Azure/go-autorest/autorest"
)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result insights.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result insights.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*insights.OperationsClient)(nil)

// AnnotationsClientAPI contains the set of methods on the AnnotationsClient type.
type AnnotationsClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, resourceName string, annotationProperties insights.Annotation) (result insights.ListAnnotation, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string, annotationID string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string, annotationID string) (result insights.ListAnnotation, err error)
	List(ctx context.Context, resourceGroupName string, resourceName string, start string, end string) (result insights.AnnotationsListResult, err error)
}

var _ AnnotationsClientAPI = (*insights.AnnotationsClient)(nil)

// APIKeysClientAPI contains the set of methods on the APIKeysClient type.
type APIKeysClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, resourceName string, APIKeyProperties insights.APIKeyRequest) (result insights.ApplicationInsightsComponentAPIKey, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string, keyID string) (result insights.ApplicationInsightsComponentAPIKey, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string, keyID string) (result insights.ApplicationInsightsComponentAPIKey, err error)
	List(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ApplicationInsightsComponentAPIKeyListResult, err error)
}

var _ APIKeysClientAPI = (*insights.APIKeysClient)(nil)

// ExportConfigurationsClientAPI contains the set of methods on the ExportConfigurationsClient type.
type ExportConfigurationsClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, resourceName string, exportProperties insights.ApplicationInsightsComponentExportRequest) (result insights.ListApplicationInsightsComponentExportConfiguration, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string, exportID string) (result insights.ApplicationInsightsComponentExportConfiguration, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string, exportID string) (result insights.ApplicationInsightsComponentExportConfiguration, err error)
	List(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ListApplicationInsightsComponentExportConfiguration, err error)
	Update(ctx context.Context, resourceGroupName string, resourceName string, exportID string, exportProperties insights.ApplicationInsightsComponentExportRequest) (result insights.ApplicationInsightsComponentExportConfiguration, err error)
}

var _ ExportConfigurationsClientAPI = (*insights.ExportConfigurationsClient)(nil)

// ComponentCurrentBillingFeaturesClientAPI contains the set of methods on the ComponentCurrentBillingFeaturesClient type.
type ComponentCurrentBillingFeaturesClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ApplicationInsightsComponentBillingFeatures, err error)
	Update(ctx context.Context, resourceGroupName string, resourceName string, billingFeaturesProperties insights.ApplicationInsightsComponentBillingFeatures) (result insights.ApplicationInsightsComponentBillingFeatures, err error)
}

var _ ComponentCurrentBillingFeaturesClientAPI = (*insights.ComponentCurrentBillingFeaturesClient)(nil)

// ComponentQuotaStatusClientAPI contains the set of methods on the ComponentQuotaStatusClient type.
type ComponentQuotaStatusClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ApplicationInsightsComponentQuotaStatus, err error)
}

var _ ComponentQuotaStatusClientAPI = (*insights.ComponentQuotaStatusClient)(nil)

// ComponentFeatureCapabilitiesClientAPI contains the set of methods on the ComponentFeatureCapabilitiesClient type.
type ComponentFeatureCapabilitiesClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ApplicationInsightsComponentFeatureCapabilities, err error)
}

var _ ComponentFeatureCapabilitiesClientAPI = (*insights.ComponentFeatureCapabilitiesClient)(nil)

// ComponentAvailableFeaturesClientAPI contains the set of methods on the ComponentAvailableFeaturesClient type.
type ComponentAvailableFeaturesClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ApplicationInsightsComponentAvailableFeatures, err error)
}

var _ ComponentAvailableFeaturesClientAPI = (*insights.ComponentAvailableFeaturesClient)(nil)

// ProactiveDetectionConfigurationsClientAPI contains the set of methods on the ProactiveDetectionConfigurationsClient type.
type ProactiveDetectionConfigurationsClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, resourceName string, configurationID string) (result insights.ApplicationInsightsComponentProactiveDetectionConfiguration, err error)
	List(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ListApplicationInsightsComponentProactiveDetectionConfiguration, err error)
	Update(ctx context.Context, resourceGroupName string, resourceName string, configurationID string, proactiveDetectionProperties insights.ApplicationInsightsComponentProactiveDetectionConfiguration) (result insights.ApplicationInsightsComponentProactiveDetectionConfiguration, err error)
}

var _ ProactiveDetectionConfigurationsClientAPI = (*insights.ProactiveDetectionConfigurationsClient)(nil)

// WorkItemConfigurationsClientAPI contains the set of methods on the WorkItemConfigurationsClient type.
type WorkItemConfigurationsClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, resourceName string, workItemConfigurationProperties insights.WorkItemCreateConfiguration) (result insights.WorkItemConfiguration, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string, workItemConfigID string) (result autorest.Response, err error)
	GetDefault(ctx context.Context, resourceGroupName string, resourceName string) (result insights.WorkItemConfiguration, err error)
	GetItem(ctx context.Context, resourceGroupName string, resourceName string, workItemConfigID string) (result insights.WorkItemConfiguration, err error)
	List(ctx context.Context, resourceGroupName string, resourceName string) (result insights.WorkItemConfigurationsListResult, err error)
	UpdateItem(ctx context.Context, resourceGroupName string, resourceName string, workItemConfigID string, workItemConfigurationProperties insights.WorkItemCreateConfiguration) (result insights.WorkItemConfiguration, err error)
}

var _ WorkItemConfigurationsClientAPI = (*insights.WorkItemConfigurationsClient)(nil)

// FavoritesClientAPI contains the set of methods on the FavoritesClient type.
type FavoritesClientAPI interface {
	Add(ctx context.Context, resourceGroupName string, resourceName string, favoriteID string, favoriteProperties insights.ApplicationInsightsComponentFavorite) (result insights.ApplicationInsightsComponentFavorite, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string, favoriteID string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string, favoriteID string) (result insights.ApplicationInsightsComponentFavorite, err error)
	List(ctx context.Context, resourceGroupName string, resourceName string, favoriteType insights.FavoriteType, sourceType insights.FavoriteSourceType, canFetchContent *bool, tags []string) (result insights.ListApplicationInsightsComponentFavorite, err error)
	Update(ctx context.Context, resourceGroupName string, resourceName string, favoriteID string, favoriteProperties insights.ApplicationInsightsComponentFavorite) (result insights.ApplicationInsightsComponentFavorite, err error)
}

var _ FavoritesClientAPI = (*insights.FavoritesClient)(nil)

// WebTestLocationsClientAPI contains the set of methods on the WebTestLocationsClient type.
type WebTestLocationsClientAPI interface {
	List(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ApplicationInsightsWebTestLocationsListResult, err error)
}

var _ WebTestLocationsClientAPI = (*insights.WebTestLocationsClient)(nil)

// WebTestsClientAPI contains the set of methods on the WebTestsClient type.
type WebTestsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, webTestName string, webTestDefinition insights.WebTest) (result insights.WebTest, err error)
	Delete(ctx context.Context, resourceGroupName string, webTestName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, webTestName string) (result insights.WebTest, err error)
	List(ctx context.Context) (result insights.WebTestListResultPage, err error)
	ListComplete(ctx context.Context) (result insights.WebTestListResultIterator, err error)
	ListByComponent(ctx context.Context, componentName string, resourceGroupName string) (result insights.WebTestListResultPage, err error)
	ListByComponentComplete(ctx context.Context, componentName string, resourceGroupName string) (result insights.WebTestListResultIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result insights.WebTestListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result insights.WebTestListResultIterator, err error)
	UpdateTags(ctx context.Context, resourceGroupName string, webTestName string, webTestTags insights.TagsResource) (result insights.WebTest, err error)
}

var _ WebTestsClientAPI = (*insights.WebTestsClient)(nil)

// AnalyticsItemsClientAPI contains the set of methods on the AnalyticsItemsClient type.
type AnalyticsItemsClientAPI interface {
	Delete(ctx context.Context, resourceGroupName string, resourceName string, scopePath insights.ItemScopePath, ID string, name string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string, scopePath insights.ItemScopePath, ID string, name string) (result insights.ApplicationInsightsComponentAnalyticsItem, err error)
	List(ctx context.Context, resourceGroupName string, resourceName string, scopePath insights.ItemScopePath, scope insights.ItemScope, typeParameter insights.ItemTypeParameter, includeContent *bool) (result insights.ListApplicationInsightsComponentAnalyticsItem, err error)
	Put(ctx context.Context, resourceGroupName string, resourceName string, scopePath insights.ItemScopePath, itemProperties insights.ApplicationInsightsComponentAnalyticsItem, overrideItem *bool) (result insights.ApplicationInsightsComponentAnalyticsItem, err error)
}

var _ AnalyticsItemsClientAPI = (*insights.AnalyticsItemsClient)(nil)

// WorkbooksClientAPI contains the set of methods on the WorkbooksClient type.
type WorkbooksClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties insights.Workbook) (result insights.Workbook, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result insights.Workbook, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string, category insights.CategoryType, tags []string, canFetchContent *bool) (result insights.WorkbooksListResult, err error)
	Update(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties insights.Workbook) (result insights.Workbook, err error)
}

var _ WorkbooksClientAPI = (*insights.WorkbooksClient)(nil)

// MyWorkbooksClientAPI contains the set of methods on the MyWorkbooksClient type.
type MyWorkbooksClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties insights.MyWorkbook) (result insights.MyWorkbook, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result insights.MyWorkbook, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string, category insights.CategoryType, tags []string, canFetchContent *bool) (result insights.MyWorkbooksListResult, err error)
	ListBySubscription(ctx context.Context, category insights.CategoryType, tags []string, canFetchContent *bool) (result insights.MyWorkbooksListResult, err error)
	Update(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties insights.MyWorkbook) (result insights.MyWorkbook, err error)
}

var _ MyWorkbooksClientAPI = (*insights.MyWorkbooksClient)(nil)

// ComponentsClientAPI contains the set of methods on the ComponentsClient type.
type ComponentsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, insightProperties insights.ApplicationInsightsComponent) (result insights.ApplicationInsightsComponent, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ApplicationInsightsComponent, err error)
	GetPurgeStatus(ctx context.Context, resourceGroupName string, resourceName string, purgeID string) (result insights.ComponentPurgeStatusResponse, err error)
	List(ctx context.Context) (result insights.ApplicationInsightsComponentListResultPage, err error)
	ListComplete(ctx context.Context) (result insights.ApplicationInsightsComponentListResultIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result insights.ApplicationInsightsComponentListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result insights.ApplicationInsightsComponentListResultIterator, err error)
	Purge(ctx context.Context, resourceGroupName string, resourceName string, body insights.ComponentPurgeBody) (result insights.ComponentPurgeResponse, err error)
	UpdateTags(ctx context.Context, resourceGroupName string, resourceName string, componentTags insights.TagsResource) (result insights.ApplicationInsightsComponent, err error)
}

var _ ComponentsClientAPI = (*insights.ComponentsClient)(nil)

// ComponentLinkedStorageAccountsClientAPI contains the set of methods on the ComponentLinkedStorageAccountsClient type.
type ComponentLinkedStorageAccountsClientAPI interface {
	CreateAndUpdate(ctx context.Context, resourceGroupName string, resourceName string, linkedStorageAccountsProperties insights.ComponentLinkedStorageAccounts) (result insights.ComponentLinkedStorageAccounts, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result insights.ComponentLinkedStorageAccounts, err error)
	Update(ctx context.Context, resourceGroupName string, resourceName string, linkedStorageAccountsProperties insights.ComponentLinkedStorageAccountsPatch) (result insights.ComponentLinkedStorageAccounts, err error)
}

var _ ComponentLinkedStorageAccountsClientAPI = (*insights.ComponentLinkedStorageAccountsClient)(nil)

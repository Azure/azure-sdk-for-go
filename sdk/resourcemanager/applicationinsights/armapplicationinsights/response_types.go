//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armapplicationinsights

// APIKeysClientCreateResponse contains the response from method APIKeysClient.Create.
type APIKeysClientCreateResponse struct {
	ComponentAPIKey
}

// APIKeysClientDeleteResponse contains the response from method APIKeysClient.Delete.
type APIKeysClientDeleteResponse struct {
	ComponentAPIKey
}

// APIKeysClientGetResponse contains the response from method APIKeysClient.Get.
type APIKeysClientGetResponse struct {
	ComponentAPIKey
}

// APIKeysClientListResponse contains the response from method APIKeysClient.NewListPager.
type APIKeysClientListResponse struct {
	ComponentAPIKeyListResult
}

// AnalyticsItemsClientDeleteResponse contains the response from method AnalyticsItemsClient.Delete.
type AnalyticsItemsClientDeleteResponse struct {
	// placeholder for future response values
}

// AnalyticsItemsClientGetResponse contains the response from method AnalyticsItemsClient.Get.
type AnalyticsItemsClientGetResponse struct {
	ComponentAnalyticsItem
}

// AnalyticsItemsClientListResponse contains the response from method AnalyticsItemsClient.List.
type AnalyticsItemsClientListResponse struct {
	// Array of ApplicationInsightsComponentAnalyticsItem
	ComponentAnalyticsItemArray []*ComponentAnalyticsItem
}

// AnalyticsItemsClientPutResponse contains the response from method AnalyticsItemsClient.Put.
type AnalyticsItemsClientPutResponse struct {
	ComponentAnalyticsItem
}

// AnnotationsClientCreateResponse contains the response from method AnnotationsClient.Create.
type AnnotationsClientCreateResponse struct {
	// Array of Annotation
	AnnotationArray []*Annotation
}

// AnnotationsClientDeleteResponse contains the response from method AnnotationsClient.Delete.
type AnnotationsClientDeleteResponse struct {
	// placeholder for future response values
}

// AnnotationsClientGetResponse contains the response from method AnnotationsClient.Get.
type AnnotationsClientGetResponse struct {
	// Array of Annotation
	AnnotationArray []*Annotation
}

// AnnotationsClientListResponse contains the response from method AnnotationsClient.NewListPager.
type AnnotationsClientListResponse struct {
	AnnotationsListResult
}

// ComponentAvailableFeaturesClientGetResponse contains the response from method ComponentAvailableFeaturesClient.Get.
type ComponentAvailableFeaturesClientGetResponse struct {
	ComponentAvailableFeatures
}

// ComponentCurrentBillingFeaturesClientGetResponse contains the response from method ComponentCurrentBillingFeaturesClient.Get.
type ComponentCurrentBillingFeaturesClientGetResponse struct {
	ComponentBillingFeatures
}

// ComponentCurrentBillingFeaturesClientUpdateResponse contains the response from method ComponentCurrentBillingFeaturesClient.Update.
type ComponentCurrentBillingFeaturesClientUpdateResponse struct {
	ComponentBillingFeatures
}

// ComponentFeatureCapabilitiesClientGetResponse contains the response from method ComponentFeatureCapabilitiesClient.Get.
type ComponentFeatureCapabilitiesClientGetResponse struct {
	ComponentFeatureCapabilities
}

// ComponentLinkedStorageAccountsClientCreateAndUpdateResponse contains the response from method ComponentLinkedStorageAccountsClient.CreateAndUpdate.
type ComponentLinkedStorageAccountsClientCreateAndUpdateResponse struct {
	ComponentLinkedStorageAccounts
}

// ComponentLinkedStorageAccountsClientDeleteResponse contains the response from method ComponentLinkedStorageAccountsClient.Delete.
type ComponentLinkedStorageAccountsClientDeleteResponse struct {
	// placeholder for future response values
}

// ComponentLinkedStorageAccountsClientGetResponse contains the response from method ComponentLinkedStorageAccountsClient.Get.
type ComponentLinkedStorageAccountsClientGetResponse struct {
	ComponentLinkedStorageAccounts
}

// ComponentLinkedStorageAccountsClientUpdateResponse contains the response from method ComponentLinkedStorageAccountsClient.Update.
type ComponentLinkedStorageAccountsClientUpdateResponse struct {
	ComponentLinkedStorageAccounts
}

// ComponentQuotaStatusClientGetResponse contains the response from method ComponentQuotaStatusClient.Get.
type ComponentQuotaStatusClientGetResponse struct {
	ComponentQuotaStatus
}

// ComponentsClientCreateOrUpdateResponse contains the response from method ComponentsClient.CreateOrUpdate.
type ComponentsClientCreateOrUpdateResponse struct {
	Component
}

// ComponentsClientDeleteResponse contains the response from method ComponentsClient.Delete.
type ComponentsClientDeleteResponse struct {
	// placeholder for future response values
}

// ComponentsClientGetPurgeStatusResponse contains the response from method ComponentsClient.GetPurgeStatus.
type ComponentsClientGetPurgeStatusResponse struct {
	ComponentPurgeStatusResponse
}

// ComponentsClientGetResponse contains the response from method ComponentsClient.Get.
type ComponentsClientGetResponse struct {
	Component
}

// ComponentsClientListByResourceGroupResponse contains the response from method ComponentsClient.NewListByResourceGroupPager.
type ComponentsClientListByResourceGroupResponse struct {
	ComponentListResult
}

// ComponentsClientListResponse contains the response from method ComponentsClient.NewListPager.
type ComponentsClientListResponse struct {
	ComponentListResult
}

// ComponentsClientPurgeResponse contains the response from method ComponentsClient.Purge.
type ComponentsClientPurgeResponse struct {
	ComponentPurgeResponse
}

// ComponentsClientUpdateTagsResponse contains the response from method ComponentsClient.UpdateTags.
type ComponentsClientUpdateTagsResponse struct {
	Component
}

// ExportConfigurationsClientCreateResponse contains the response from method ExportConfigurationsClient.Create.
type ExportConfigurationsClientCreateResponse struct {
	// A list of Continuous Export configurations.
	ComponentExportConfigurationArray []*ComponentExportConfiguration
}

// ExportConfigurationsClientDeleteResponse contains the response from method ExportConfigurationsClient.Delete.
type ExportConfigurationsClientDeleteResponse struct {
	ComponentExportConfiguration
}

// ExportConfigurationsClientGetResponse contains the response from method ExportConfigurationsClient.Get.
type ExportConfigurationsClientGetResponse struct {
	ComponentExportConfiguration
}

// ExportConfigurationsClientListResponse contains the response from method ExportConfigurationsClient.List.
type ExportConfigurationsClientListResponse struct {
	// A list of Continuous Export configurations.
	ComponentExportConfigurationArray []*ComponentExportConfiguration
}

// ExportConfigurationsClientUpdateResponse contains the response from method ExportConfigurationsClient.Update.
type ExportConfigurationsClientUpdateResponse struct {
	ComponentExportConfiguration
}

// FavoritesClientAddResponse contains the response from method FavoritesClient.Add.
type FavoritesClientAddResponse struct {
	ComponentFavorite
}

// FavoritesClientDeleteResponse contains the response from method FavoritesClient.Delete.
type FavoritesClientDeleteResponse struct {
	// placeholder for future response values
}

// FavoritesClientGetResponse contains the response from method FavoritesClient.Get.
type FavoritesClientGetResponse struct {
	ComponentFavorite
}

// FavoritesClientListResponse contains the response from method FavoritesClient.List.
type FavoritesClientListResponse struct {
	// Array of ApplicationInsightsComponentFavorite
	ComponentFavoriteArray []*ComponentFavorite
}

// FavoritesClientUpdateResponse contains the response from method FavoritesClient.Update.
type FavoritesClientUpdateResponse struct {
	ComponentFavorite
}

// LiveTokenClientGetResponse contains the response from method LiveTokenClient.Get.
type LiveTokenClientGetResponse struct {
	LiveTokenResponse
}

// MyWorkbooksClientCreateOrUpdateResponse contains the response from method MyWorkbooksClient.CreateOrUpdate.
type MyWorkbooksClientCreateOrUpdateResponse struct {
	MyWorkbook
}

// MyWorkbooksClientDeleteResponse contains the response from method MyWorkbooksClient.Delete.
type MyWorkbooksClientDeleteResponse struct {
	// placeholder for future response values
}

// MyWorkbooksClientGetResponse contains the response from method MyWorkbooksClient.Get.
type MyWorkbooksClientGetResponse struct {
	MyWorkbook
}

// MyWorkbooksClientListByResourceGroupResponse contains the response from method MyWorkbooksClient.NewListByResourceGroupPager.
type MyWorkbooksClientListByResourceGroupResponse struct {
	MyWorkbooksListResult
}

// MyWorkbooksClientListBySubscriptionResponse contains the response from method MyWorkbooksClient.NewListBySubscriptionPager.
type MyWorkbooksClientListBySubscriptionResponse struct {
	MyWorkbooksListResult
}

// MyWorkbooksClientUpdateResponse contains the response from method MyWorkbooksClient.Update.
type MyWorkbooksClientUpdateResponse struct {
	MyWorkbook
}

// ProactiveDetectionConfigurationsClientGetResponse contains the response from method ProactiveDetectionConfigurationsClient.Get.
type ProactiveDetectionConfigurationsClientGetResponse struct {
	ComponentProactiveDetectionConfiguration
}

// ProactiveDetectionConfigurationsClientListResponse contains the response from method ProactiveDetectionConfigurationsClient.List.
type ProactiveDetectionConfigurationsClientListResponse struct {
	// A list of ProactiveDetection configurations.
	ComponentProactiveDetectionConfigurationArray []*ComponentProactiveDetectionConfiguration
}

// ProactiveDetectionConfigurationsClientUpdateResponse contains the response from method ProactiveDetectionConfigurationsClient.Update.
type ProactiveDetectionConfigurationsClientUpdateResponse struct {
	ComponentProactiveDetectionConfiguration
}

// WebTestLocationsClientListResponse contains the response from method WebTestLocationsClient.NewListPager.
type WebTestLocationsClientListResponse struct {
	WebTestLocationsListResult
}

// WebTestsClientCreateOrUpdateResponse contains the response from method WebTestsClient.CreateOrUpdate.
type WebTestsClientCreateOrUpdateResponse struct {
	WebTest
}

// WebTestsClientDeleteResponse contains the response from method WebTestsClient.Delete.
type WebTestsClientDeleteResponse struct {
	// placeholder for future response values
}

// WebTestsClientGetResponse contains the response from method WebTestsClient.Get.
type WebTestsClientGetResponse struct {
	WebTest
}

// WebTestsClientListByComponentResponse contains the response from method WebTestsClient.NewListByComponentPager.
type WebTestsClientListByComponentResponse struct {
	WebTestListResult
}

// WebTestsClientListByResourceGroupResponse contains the response from method WebTestsClient.NewListByResourceGroupPager.
type WebTestsClientListByResourceGroupResponse struct {
	WebTestListResult
}

// WebTestsClientListResponse contains the response from method WebTestsClient.NewListPager.
type WebTestsClientListResponse struct {
	WebTestListResult
}

// WebTestsClientUpdateTagsResponse contains the response from method WebTestsClient.UpdateTags.
type WebTestsClientUpdateTagsResponse struct {
	WebTest
}

// WorkItemConfigurationsClientCreateResponse contains the response from method WorkItemConfigurationsClient.Create.
type WorkItemConfigurationsClientCreateResponse struct {
	WorkItemConfiguration
}

// WorkItemConfigurationsClientDeleteResponse contains the response from method WorkItemConfigurationsClient.Delete.
type WorkItemConfigurationsClientDeleteResponse struct {
	// placeholder for future response values
}

// WorkItemConfigurationsClientGetDefaultResponse contains the response from method WorkItemConfigurationsClient.GetDefault.
type WorkItemConfigurationsClientGetDefaultResponse struct {
	WorkItemConfiguration
}

// WorkItemConfigurationsClientGetItemResponse contains the response from method WorkItemConfigurationsClient.GetItem.
type WorkItemConfigurationsClientGetItemResponse struct {
	WorkItemConfiguration
}

// WorkItemConfigurationsClientListResponse contains the response from method WorkItemConfigurationsClient.NewListPager.
type WorkItemConfigurationsClientListResponse struct {
	WorkItemConfigurationsListResult
}

// WorkItemConfigurationsClientUpdateItemResponse contains the response from method WorkItemConfigurationsClient.UpdateItem.
type WorkItemConfigurationsClientUpdateItemResponse struct {
	WorkItemConfiguration
}

// WorkbookTemplatesClientCreateOrUpdateResponse contains the response from method WorkbookTemplatesClient.CreateOrUpdate.
type WorkbookTemplatesClientCreateOrUpdateResponse struct {
	WorkbookTemplate
}

// WorkbookTemplatesClientDeleteResponse contains the response from method WorkbookTemplatesClient.Delete.
type WorkbookTemplatesClientDeleteResponse struct {
	// placeholder for future response values
}

// WorkbookTemplatesClientGetResponse contains the response from method WorkbookTemplatesClient.Get.
type WorkbookTemplatesClientGetResponse struct {
	WorkbookTemplate
}

// WorkbookTemplatesClientListByResourceGroupResponse contains the response from method WorkbookTemplatesClient.NewListByResourceGroupPager.
type WorkbookTemplatesClientListByResourceGroupResponse struct {
	WorkbookTemplatesListResult
}

// WorkbookTemplatesClientUpdateResponse contains the response from method WorkbookTemplatesClient.Update.
type WorkbookTemplatesClientUpdateResponse struct {
	WorkbookTemplate
}

// WorkbooksClientCreateOrUpdateResponse contains the response from method WorkbooksClient.CreateOrUpdate.
type WorkbooksClientCreateOrUpdateResponse struct {
	Workbook
}

// WorkbooksClientDeleteResponse contains the response from method WorkbooksClient.Delete.
type WorkbooksClientDeleteResponse struct {
	// placeholder for future response values
}

// WorkbooksClientGetResponse contains the response from method WorkbooksClient.Get.
type WorkbooksClientGetResponse struct {
	Workbook
}

// WorkbooksClientListByResourceGroupResponse contains the response from method WorkbooksClient.NewListByResourceGroupPager.
type WorkbooksClientListByResourceGroupResponse struct {
	WorkbooksListResult
}

// WorkbooksClientListBySubscriptionResponse contains the response from method WorkbooksClient.NewListBySubscriptionPager.
type WorkbooksClientListBySubscriptionResponse struct {
	WorkbooksListResult
}

// WorkbooksClientRevisionGetResponse contains the response from method WorkbooksClient.RevisionGet.
type WorkbooksClientRevisionGetResponse struct {
	Workbook
}

// WorkbooksClientRevisionsListResponse contains the response from method WorkbooksClient.NewRevisionsListPager.
type WorkbooksClientRevisionsListResponse struct {
	WorkbooksListResult
}

// WorkbooksClientUpdateResponse contains the response from method WorkbooksClient.Update.
type WorkbooksClientUpdateResponse struct {
	Workbook
}

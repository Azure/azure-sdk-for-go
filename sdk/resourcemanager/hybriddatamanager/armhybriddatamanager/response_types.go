//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybriddatamanager

// DataManagersClientCreateResponse contains the response from method DataManagersClient.BeginCreate.
type DataManagersClientCreateResponse struct {
	// The DataManager resource.
	DataManager
}

// DataManagersClientDeleteResponse contains the response from method DataManagersClient.BeginDelete.
type DataManagersClientDeleteResponse struct {
	// placeholder for future response values
}

// DataManagersClientGetResponse contains the response from method DataManagersClient.Get.
type DataManagersClientGetResponse struct {
	// The DataManager resource.
	DataManager
}

// DataManagersClientListByResourceGroupResponse contains the response from method DataManagersClient.NewListByResourceGroupPager.
type DataManagersClientListByResourceGroupResponse struct {
	// DataManager resources Collection.
	DataManagerList
}

// DataManagersClientListResponse contains the response from method DataManagersClient.NewListPager.
type DataManagersClientListResponse struct {
	// DataManager resources Collection.
	DataManagerList
}

// DataManagersClientUpdateResponse contains the response from method DataManagersClient.BeginUpdate.
type DataManagersClientUpdateResponse struct {
	// The DataManager resource.
	DataManager
}

// DataServicesClientGetResponse contains the response from method DataServicesClient.Get.
type DataServicesClientGetResponse struct {
	// Data Service.
	DataService
}

// DataServicesClientListByDataManagerResponse contains the response from method DataServicesClient.NewListByDataManagerPager.
type DataServicesClientListByDataManagerResponse struct {
	// Data Service Collection.
	DataServiceList
}

// DataStoreTypesClientGetResponse contains the response from method DataStoreTypesClient.Get.
type DataStoreTypesClientGetResponse struct {
	// Data Store Type.
	DataStoreType
}

// DataStoreTypesClientListByDataManagerResponse contains the response from method DataStoreTypesClient.NewListByDataManagerPager.
type DataStoreTypesClientListByDataManagerResponse struct {
	// Data Store Type Collection.
	DataStoreTypeList
}

// DataStoresClientCreateOrUpdateResponse contains the response from method DataStoresClient.BeginCreateOrUpdate.
type DataStoresClientCreateOrUpdateResponse struct {
	// Data store.
	DataStore
}

// DataStoresClientDeleteResponse contains the response from method DataStoresClient.BeginDelete.
type DataStoresClientDeleteResponse struct {
	// placeholder for future response values
}

// DataStoresClientGetResponse contains the response from method DataStoresClient.Get.
type DataStoresClientGetResponse struct {
	// Data store.
	DataStore
}

// DataStoresClientListByDataManagerResponse contains the response from method DataStoresClient.NewListByDataManagerPager.
type DataStoresClientListByDataManagerResponse struct {
	// Data Store Collection.
	DataStoreList
}

// JobDefinitionsClientCreateOrUpdateResponse contains the response from method JobDefinitionsClient.BeginCreateOrUpdate.
type JobDefinitionsClientCreateOrUpdateResponse struct {
	// Job Definition.
	JobDefinition
}

// JobDefinitionsClientDeleteResponse contains the response from method JobDefinitionsClient.BeginDelete.
type JobDefinitionsClientDeleteResponse struct {
	// placeholder for future response values
}

// JobDefinitionsClientGetResponse contains the response from method JobDefinitionsClient.Get.
type JobDefinitionsClientGetResponse struct {
	// Job Definition.
	JobDefinition
}

// JobDefinitionsClientListByDataManagerResponse contains the response from method JobDefinitionsClient.NewListByDataManagerPager.
type JobDefinitionsClientListByDataManagerResponse struct {
	// Job Definition Collection.
	JobDefinitionList
}

// JobDefinitionsClientListByDataServiceResponse contains the response from method JobDefinitionsClient.NewListByDataServicePager.
type JobDefinitionsClientListByDataServiceResponse struct {
	// Job Definition Collection.
	JobDefinitionList
}

// JobDefinitionsClientRunResponse contains the response from method JobDefinitionsClient.BeginRun.
type JobDefinitionsClientRunResponse struct {
	// placeholder for future response values
}

// JobsClientCancelResponse contains the response from method JobsClient.BeginCancel.
type JobsClientCancelResponse struct {
	// placeholder for future response values
}

// JobsClientGetResponse contains the response from method JobsClient.Get.
type JobsClientGetResponse struct {
	// Data service job.
	Job
}

// JobsClientListByDataManagerResponse contains the response from method JobsClient.NewListByDataManagerPager.
type JobsClientListByDataManagerResponse struct {
	// Job Collection.
	JobList
}

// JobsClientListByDataServiceResponse contains the response from method JobsClient.NewListByDataServicePager.
type JobsClientListByDataServiceResponse struct {
	// Job Collection.
	JobList
}

// JobsClientListByJobDefinitionResponse contains the response from method JobsClient.NewListByJobDefinitionPager.
type JobsClientListByJobDefinitionResponse struct {
	// Job Collection.
	JobList
}

// JobsClientResumeResponse contains the response from method JobsClient.BeginResume.
type JobsClientResumeResponse struct {
	// placeholder for future response values
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// Class for set of operations used for discovery of available provider operations.
	AvailableProviderOperations
}

// PublicKeysClientGetResponse contains the response from method PublicKeysClient.Get.
type PublicKeysClientGetResponse struct {
	// Public key
	PublicKey
}

// PublicKeysClientListByDataManagerResponse contains the response from method PublicKeysClient.NewListByDataManagerPager.
type PublicKeysClientListByDataManagerResponse struct {
	// PublicKey Collection
	PublicKeyList
}

package search

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 0.11.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

type ProvisioningState string

const (
	ProvisioningStateFailed       ProvisioningState = "failed"
	ProvisioningStateProvisioning ProvisioningState = "provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "succeeded"
)

type SearchServiceStatus string

const (
	SearchServiceStatusDegraded     SearchServiceStatus = "degraded"
	SearchServiceStatusDeleting     SearchServiceStatus = "deleting"
	SearchServiceStatusDisabled     SearchServiceStatus = "disabled"
	SearchServiceStatusError        SearchServiceStatus = "error"
	SearchServiceStatusProvisioning SearchServiceStatus = "provisioning"
	SearchServiceStatusRunning      SearchServiceStatus = "running"
)

type SkuType string

const (
	Free      SkuType = "free"
	Standard  SkuType = "standard"
	Standard2 SkuType = "standard2"
)

// Response containing the primary and secondary API keys for a given Azure
// Search service.
type AdminKeyResult struct {
	PrimaryKey   string `json:"primaryKey,omitempty"`
	SecondaryKey string `json:"secondaryKey,omitempty"`
}

// Response containing the query API keys for a given Azure Search service.
type ListQueryKeysResult struct {
	Value []QueryKey `json:"value,omitempty"`
}

// Describes an API key for a given Azure Search service that has permissions
// for query operations only.
type QueryKey struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Properties that describe an Azure Search service.
type SearchServiceCreateOrUpdateParameters struct {
	Location   string            `json:"location,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
	Properties struct {
		Sku struct {
			Name SkuType `json:"name,omitempty"`
		} `json:"sku,omitempty"`
		ReplicaCount   int `json:"replicaCount,omitempty"`
		PartitionCount int `json:"partitionCount,omitempty"`
	} `json:"properties,omitempty"`
}

// Response containing a list of Azure Search services for a given resource
// group.
type SearchServiceListResult struct {
	Value []SearchServiceResource `json:"value,omitempty"`
}

// Describes an Azure Search service and its current state.
type SearchServiceResource struct {
	Name       string            `json:"name,omitempty"`
	Location   string            `json:"location,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
	Properties struct {
		Status            SearchServiceStatus `json:"status,omitempty"`
		StatusDetails     string              `json:"statusDetails,omitempty"`
		ProvisioningState ProvisioningState   `json:"provisioningState,omitempty"`
		Sku               struct {
			Name SkuType `json:"name,omitempty"`
		} `json:"sku,omitempty"`
		ReplicaCount   int `json:"replicaCount,omitempty"`
		PartitionCount int `json:"partitionCount,omitempty"`
	} `json:"properties,omitempty"`
}

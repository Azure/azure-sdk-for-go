// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// Headers
const (
	cosmosHeaderRequestCharge                      string = "x-ms-request-charge"
	cosmosHeaderActivityId                         string = "x-ms-activity-id"
	cosmosHeaderEtag                               string = "etag"
	cosmosHeaderSubstatus                          string = "x-ms-substatus"
	cosmosHeaderPopulateQuotaInfo                  string = "x-ms-documentdb-populatequotainfo"
	cosmosHeaderPreTriggerInclude                  string = "x-ms-documentdb-pre-trigger-include"
	cosmosHeaderPostTriggerInclude                 string = "x-ms-documentdb-post-trigger-include"
	cosmosHeaderIndexingDirective                  string = "x-ms-indexing-directive"
	cosmosHeaderSessionToken                       string = "x-ms-session-token"
	cosmosHeaderConsistencyLevel                   string = "x-ms-consistency-level"
	cosmosHeaderPartitionKey                       string = "x-ms-documentdb-partitionkey"
	cosmosHeaderPrefer                             string = "Prefer"
	cosmosHeaderIsUpsert                           string = "x-ms-documentdb-is-upsert"
	cosmosHeaderOfferThroughput                    string = "x-ms-offer-throughput"
	cosmosHeaderOfferAutoscale                     string = "x-ms-cosmos-offer-autopilot-settings"
	cosmosHeaderQuery                              string = "x-ms-documentdb-query"
	cosmosHeaderOfferReplacePending                string = "x-ms-offer-replace-pending"
	cosmosHeaderOfferMinimumThroughput             string = "x-ms-cosmos-min-throughput"
	cosmosHeaderResponseContinuationTokenLimitInKb string = "x-ms-documentdb-responsecontinuationtokenlimitinkb"
	cosmosHeaderEnableScanInQuery                  string = "x-ms-documentdb-force-query-scan"
	cosmosHeaderMaxItemCount                       string = "x-ms-max-item-count"
	cosmosHeaderContinuationToken                  string = "x-ms-continuation"
	cosmosHeaderPopulateIndexMetrics               string = "x-ms-cosmos-populateindexmetrics"
	cosmosHeaderPopulateQueryMetrics               string = "x-ms-documentdb-populatequerymetrics"
	cosmosHeaderQueryMetrics                       string = "x-ms-documentdb-query-metrics"
	cosmosHeaderIndexUtilization                   string = "x-ms-cosmos-index-utilization"
	cosmosHeaderCorrelatedActivityId               string = "x-ms-cosmos-correlated-activityid"
	cosmosHeaderIsBatchRequest                     string = "x-ms-cosmos-is-batch-request"
	cosmosHeaderIsBatchAtomic                      string = "x-ms-cosmos-batch-atomic"
	cosmosHeaderIsBatchOrdered                     string = "x-ms-cosmos-batch-ordered"
	cosmosHeaderSDKSupportedCapabilities           string = "x-ms-cosmos-sdk-supportedcapabilities"
	headerXmsDate                                  string = "x-ms-date"
	headerAuthorization                            string = "Authorization"
	headerContentType                              string = "Content-Type"
	headerIfMatch                                  string = "If-Match"
	headerIfNoneMatch                              string = "If-None-Match"
	headerXmsVersion                               string = "x-ms-version"
)

const (
	cosmosHeaderValuesPreferMinimal string = "return=minimal"
	cosmosHeaderValuesQuery         string = "application/query+json"
)

// Substatus Codes
const (
	subStatusWriteForbidden          string = "3"
	subStatusDatabaseAccountNotFound string = "1008"
	subStatusReadSessionNotAvailable string = "1002"
)

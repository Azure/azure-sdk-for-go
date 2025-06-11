// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestItemRequestOptionsToHeaders(t *testing.T) {
	options := &ItemOptions{}
	options.PreTriggers = []string{"preTrigger1", "preTrigger2"}
	options.PostTriggers = []string{"postTrigger1", "postTrigger2"}
	options.ConsistencyLevel = ConsistencyLevelSession.ToPtr()
	sessionToken := "sessionToken"
	options.SessionToken = &sessionToken
	options.IndexingDirective = IndexingDirectiveInclude.ToPtr()
	etagValue := azcore.ETag("someEtag")
	options.IfMatchEtag = &etagValue
	maxIntegratedCacheStalenessDuration := time.Duration(5 * time.Minute)
	options.DedicatedGatewayRequestOptions = &DedicatedGatewayRequestOptions{
		MaxIntegratedCacheStaleness: &maxIntegratedCacheStalenessDuration,
	}
	options.DedicatedGatewayRequestOptions.BypassIntegratedCache = true
	header := options.toHeaders()
	if header == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	headers := *header
	if headers[cosmosHeaderPreTriggerInclude] != "preTrigger1,preTrigger2" {
		t.Errorf("PreTriggerInclude should be preTrigger1,preTrigger2 but got %v", headers[cosmosHeaderPreTriggerInclude])
	}
	if headers[cosmosHeaderPostTriggerInclude] != "postTrigger1,postTrigger2" {
		t.Errorf("PostTriggerInclude should be postTrigger1,postTrigger2 but got %v", headers[cosmosHeaderPostTriggerInclude])
	}
	if headers[cosmosHeaderConsistencyLevel] != "Session" {
		t.Errorf("ConsistencyLevel should be Session but got %v", headers[cosmosHeaderConsistencyLevel])
	}
	if headers[cosmosHeaderIndexingDirective] != "Include" {
		t.Errorf("IndexingDirective should be Include but got %v", headers[cosmosHeaderIndexingDirective])
	}
	if headers[cosmosHeaderSessionToken] != "sessionToken" {
		t.Errorf("SessionToken should be sessionToken but got %v", headers[cosmosHeaderSessionToken])
	}
	if headers[headerIfMatch] != string(*options.IfMatchEtag) {
		t.Errorf("IfMatchEtag should be someEtag but got %v", headers[headerIfMatch])
	}
	if headers[headerDedicatedGatewayMaxAge] != strconv.FormatInt(300000, 10) {
		t.Errorf("headerDedicatedGatewayMaxAge should be 300000 but got %v", headers[headerDedicatedGatewayMaxAge])
	}
	if headers[headerDedicatedGatewayBypassCache] != "true" {
		t.Errorf("headerDedicatedGatewayBypassCache should be true but got %v", headers[headerDedicatedGatewayBypassCache])
	}
}

func TestItemRequestOptionsToHeaders_bypassIntegratedCacheNotSet(t *testing.T) {
	options := &ItemOptions{}
	header := options.toHeaders()
	if header == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	headers := *header
	if _, exists := headers[headerDedicatedGatewayBypassCache]; exists {
		t.Errorf("headerDedicatedGatewayBypassCache should not exist when BypassIntegratedCache is not set")
	}
}

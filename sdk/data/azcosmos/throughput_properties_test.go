// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestThroughputPropertiesManualRawSerialization(t *testing.T) {
	nowAsUnix := time.Unix(1630100602, 0)

	jsonString := []byte("{\"offerType\":\"Invalid\",\"offerResourceId\":\"4SRTANCD3Dw=\",\"resource\":\"dbs/dbid/colls/collid/\", \"offerVersion\":\"V2\",\"content\":{\"offerThroughput\":400},\"_rid\":\"HFln\",\"_etag\":\"\\\"00000000-0000-0000-9b8c-8ea3e19601d7\\\"\",\"_ts\":1630100602}")

	otherProperties := &ThroughputProperties{}
	err := json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if otherProperties.offerType != "Invalid" {
		t.Errorf("OfferType mismatch %v", otherProperties.offerType)
	}

	if otherProperties.offerResourceId != "4SRTANCD3Dw=" {
		t.Errorf("OfferResourceId mismatch %v", otherProperties.offerResourceId)
	}

	if otherProperties.version != offerVersion2 {
		t.Errorf("OfferVersion mismatch %v", otherProperties.version)
	}

	if otherProperties.offerId != "HFln" {
		t.Errorf("OfferId mismatch %v", otherProperties.offerId)
	}

	if *otherProperties.ETag != "\"00000000-0000-0000-9b8c-8ea3e19601d7\"" {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified != nowAsUnix {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified)
	}

	if otherProperties.resource != "dbs/dbid/colls/collid/" {
		t.Errorf("resource mismatch %v", otherProperties.resource)
	}

	mt, isManual := otherProperties.ManualThroughput()
	if !isManual {
		t.Fatal("Expected to have manual throughput available")
	}

	if mt != 400 {
		t.Errorf("ManualThroughput mismatch %v", mt)
	}
}

func TestThroughputPropertiesManualE2ESerialization(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("\"00000000-0000-0000-9b8c-8ea3e19601d7\"")
	properties := NewManualThroughputProperties(400)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.resource = "dbs/dbid/colls/collid/"
	properties.ETag = &etag
	properties.LastModified = nowAsUnix
	jsonString, err := json.Marshal(&properties)
	if err != nil {
		t.Fatal(err)
	}

	otherProperties := &ThroughputProperties{}
	err = json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if otherProperties.offerType != "" {
		t.Errorf("OfferType mismatch %v", otherProperties.offerType)
	}

	if otherProperties.offerResourceId != "4SRTANCD3Dw=" {
		t.Errorf("OfferResourceId mismatch %v", otherProperties.offerResourceId)
	}

	if otherProperties.version != offerVersion2 {
		t.Errorf("OfferVersion mismatch %v", otherProperties.version)
	}

	if otherProperties.offerId != "HFln" {
		t.Errorf("OfferId mismatch %v", otherProperties.offerId)
	}

	if *otherProperties.ETag != etag {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified != properties.LastModified {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified)
	}

	if otherProperties.resource != properties.resource {
		t.Errorf("resource mismatch %v", otherProperties.resource)
	}

	mt, isManual := otherProperties.ManualThroughput()
	if !isManual {
		t.Fatal("Expected to have manual throughput available")
	}

	if mt != 400 {
		t.Errorf("ManualThroughput mismatch %v", mt)
	}
}

func TestThroughputPropertiesAutoscaleWithIncrementE2ESerialization(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("\"00000000-0000-0000-9b8c-8ea3e19601d7\"")
	properties := NewAutoscaleThroughputPropertiesWithIncrement(400, 500)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = &etag
	properties.LastModified = nowAsUnix
	jsonString, err := json.Marshal(&properties)
	if err != nil {
		t.Fatal(err)
	}

	otherProperties := &ThroughputProperties{}
	err = json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if otherProperties.offerType != "" {
		t.Errorf("OfferType mismatch %v", otherProperties.offerType)
	}

	if otherProperties.offerResourceId != "4SRTANCD3Dw=" {
		t.Errorf("OfferResourceId mismatch %v", otherProperties.offerResourceId)
	}

	if otherProperties.version != offerVersion2 {
		t.Errorf("OfferVersion mismatch %v", otherProperties.version)
	}

	if otherProperties.offerId != "HFln" {
		t.Errorf("OfferId mismatch %v", otherProperties.offerId)
	}

	if *otherProperties.ETag != etag {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified != properties.LastModified {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified)
	}

	at, hasAutoscale := otherProperties.AutoscaleMaxThroughput()
	if !hasAutoscale {
		t.Errorf("Expected to have autoscale")
	}

	inc, hasAutoscale := otherProperties.AutoscaleIncrement()
	if !hasAutoscale {
		t.Errorf("Expected to have autoscale")
	}

	if at != 400 {
		t.Errorf("MaxThroughput mismatch %v", at)
	}

	if inc != 500 {
		t.Errorf("Increment mismatch %v", inc)
	}

	if otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties == nil {
		t.Errorf("AutoscaleAutoUpgradeProperties mismatch %v", *otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties)
	}
}

func TestThroughputPropertiesAutoscaleE2ESerialization(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("\"00000000-0000-0000-9b8c-8ea3e19601d7\"")
	properties := NewAutoscaleThroughputProperties(400)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = &etag
	properties.LastModified = nowAsUnix
	jsonString, err := json.Marshal(&properties)
	if err != nil {
		t.Fatal(err)
	}

	otherProperties := &ThroughputProperties{}
	err = json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if otherProperties.offerType != "" {
		t.Errorf("OfferType mismatch %v", otherProperties.offerType)
	}

	if otherProperties.offerResourceId != "4SRTANCD3Dw=" {
		t.Errorf("OfferResourceId mismatch %v", otherProperties.offerResourceId)
	}

	if otherProperties.version != offerVersion2 {
		t.Errorf("OfferVersion mismatch %v", otherProperties.version)
	}

	if otherProperties.offerId != "HFln" {
		t.Errorf("OfferId mismatch %v", otherProperties.offerId)
	}

	if *otherProperties.ETag != etag {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified != properties.LastModified {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified)
	}

	at, hasAutoscale := otherProperties.AutoscaleMaxThroughput()
	if !hasAutoscale {
		t.Errorf("Expected to have autoscale")
	}

	_, hasAutoscaleIncrement := otherProperties.AutoscaleIncrement()
	if hasAutoscaleIncrement {
		t.Errorf("Expected not to have autoscale increment")
	}

	if at != 400 {
		t.Errorf("MaxThroughput mismatch %v", at)
	}

	if otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties != nil {
		t.Errorf("AutoscaleAutoUpgradeProperties mismatch %v", *otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties)
	}
}

func TestThroughputPropertiesAutoscaleIncrementE2ESerialization(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("\"00000000-0000-0000-9b8c-8ea3e19601d7\"")
	properties := NewAutoscaleThroughputPropertiesWithIncrement(400, 10)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = &etag
	properties.LastModified = nowAsUnix
	jsonString, err := json.Marshal(&properties)
	if err != nil {
		t.Fatal(err)
	}

	otherProperties := &ThroughputProperties{}
	err = json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if otherProperties.offerType != "" {
		t.Errorf("OfferType mismatch %v", otherProperties.offerType)
	}

	if otherProperties.offerResourceId != "4SRTANCD3Dw=" {
		t.Errorf("OfferResourceId mismatch %v", otherProperties.offerResourceId)
	}

	if otherProperties.version != offerVersion2 {
		t.Errorf("OfferVersion mismatch %v", otherProperties.version)
	}

	if otherProperties.offerId != "HFln" {
		t.Errorf("OfferId mismatch %v", otherProperties.offerId)
	}

	if *otherProperties.ETag != etag {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified != properties.LastModified {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified)
	}

	at, hasAutoscale := otherProperties.AutoscaleMaxThroughput()
	if !hasAutoscale {
		t.Errorf("Expected to have autoscale")
	}

	if at != 400 {
		t.Errorf("MaxThroughput mismatch %v", at)
	}

	if otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties.ThroughputPolicy.IncrementPercent != 10 {
		t.Errorf("IncrementPercent mismatch %v", otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties.ThroughputPolicy.IncrementPercent)
	}
}

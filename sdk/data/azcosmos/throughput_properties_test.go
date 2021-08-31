// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"
)

func TestThroughputPropertiesManualRawSerialization(t *testing.T) {
	nowAsUnix := time.Unix(1630100602, 0)

	jsonString := []byte("{\"offerType\":\"Invalid\",\"offerResourceId\":\"4SRTANCD3Dw=\",\"offerVersion\":\"V2\",\"content\":{\"offerThroughput\":400},\"id\":\"HFln\",\"_etag\":\"\\\"00000000-0000-0000-9b8c-8ea3e19601d7\\\"\",\"_ts\":1630100602}")

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

	if otherProperties.ETag != "\"00000000-0000-0000-9b8c-8ea3e19601d7\"" {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified.Time != nowAsUnix {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified.Time)
	}

	mt, err := otherProperties.ManualThroughput()
	if err != nil {
		t.Fatal(err)
	}

	if mt != 400 {
		t.Errorf("ManualThroughput mismatch %v", mt)
	}
}

func TestThroughputPropertiesManualE2ESerialization(t *testing.T) {
	nowAsUnix := time.Now().Unix()

	now := UnixTime{
		Time: time.Unix(nowAsUnix, 0),
	}

	properties := NewManualThroughputProperties(400)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = "\"00000000-0000-0000-9b8c-8ea3e19601d7\""
	properties.LastModified = &now
	jsonString, err := json.Marshal(properties)
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

	if otherProperties.ETag != "\"00000000-0000-0000-9b8c-8ea3e19601d7\"" {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified.Time != properties.LastModified.Time {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified.Time)
	}

	mt, err := otherProperties.ManualThroughput()
	if err != nil {
		t.Fatal(err)
	}

	if mt != 400 {
		t.Errorf("ManualThroughput mismatch %v", mt)
	}
}

func TestThroughputPropertiesAutoscaleE2ESerialization(t *testing.T) {
	nowAsUnix := time.Now().Unix()

	now := UnixTime{
		Time: time.Unix(nowAsUnix, 0),
	}

	properties := NewAutoscaleThroughputProperties(400)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = "\"00000000-0000-0000-9b8c-8ea3e19601d7\""
	properties.LastModified = &now
	jsonString, err := json.Marshal(properties)
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

	if otherProperties.ETag != "\"00000000-0000-0000-9b8c-8ea3e19601d7\"" {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified.Time != properties.LastModified.Time {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified.Time)
	}

	at, err := otherProperties.AutoscaleMaxThroughput()
	if err != nil {
		t.Fatal(err)
	}

	if at != 400 {
		t.Errorf("MaxThroughput mismatch %v", at)
	}

	if otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties != nil {
		t.Errorf("AutoscaleAutoUpgradeProperties mismatch %v", otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties)
	}
}

func TestThroughputPropertiesAutoscaleIncrementE2ESerialization(t *testing.T) {
	nowAsUnix := time.Now().Unix()

	now := UnixTime{
		Time: time.Unix(nowAsUnix, 0),
	}

	properties := NewAutoscaleThroughputPropertiesWithIncrement(400, 10)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = "\"00000000-0000-0000-9b8c-8ea3e19601d7\""
	properties.LastModified = &now
	jsonString, err := json.Marshal(properties)
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

	if otherProperties.ETag != "\"00000000-0000-0000-9b8c-8ea3e19601d7\"" {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified.Time != properties.LastModified.Time {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified.Time)
	}

	at, err := otherProperties.AutoscaleMaxThroughput()
	if err != nil {
		t.Fatal(err)
	}

	if at != 400 {
		t.Errorf("MaxThroughput mismatch %v", at)
	}

	if otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties.ThroughputPolicy.IncrementPercent != 10 {
		t.Errorf("IncrementPercent mismatch %v", otherProperties.offer.AutoScale.AutoscaleAutoUpgradeProperties.ThroughputPolicy.IncrementPercent)
	}
}

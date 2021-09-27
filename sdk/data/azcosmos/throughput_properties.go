// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	offerVersion2 string = "V2"
)

// ThroughputProperties describes the throughput configuration of a resource.
type ThroughputProperties struct {
	ETag         string
	LastModified *UnixTime

	version         string
	offerType       string
	offer           *offer
	offerResourceId string
	offerId         string
	selfLink        string
}

// NewManualThroughputProperties returns a ThroughputProperties object with the given throughput in manual mode.
// throughput - the throughput in RU/s
func NewManualThroughputProperties(throughput int) *ThroughputProperties {
	return &ThroughputProperties{
		version: offerVersion2,
		offer:   newManualOffer(throughput),
	}
}

// NewAutoscaleThroughputPropertiesWithIncrement returns a ThroughputProperties object with the given max throughput on autoscale mode.
// maxThroughput - the max throughput in RU/s
// incrementPercentage - the auto upgrade max throughput increment percentage
func NewAutoscaleThroughputPropertiesWithIncrement(startingMaxThroughput int, incrementPercentage int) *ThroughputProperties {
	return &ThroughputProperties{
		version: offerVersion2,
		offer:   newAutoscaleOfferWithIncrement(startingMaxThroughput, incrementPercentage),
	}
}

// NewAutoscaleThroughputProperties returns a ThroughputProperties object with the given max throughput on autoscale mode.
// maxThroughput - the max throughput in RU/s
func NewAutoscaleThroughputProperties(startingMaxThroughput int) *ThroughputProperties {
	return &ThroughputProperties{
		version: offerVersion2,
		offer:   newAutoscaleOffer(startingMaxThroughput),
	}
}

func (tp *ThroughputProperties) MarshalJSON() ([]byte, error) {
	offer, err := json.Marshal(tp.offer)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"content\":")
	buffer.Write(offer)

	if tp.offerResourceId != "" {
		buffer.WriteString(fmt.Sprintf(",\"offerResourceId\":\"%s\"", tp.offerResourceId))
	}

	if tp.offerId != "" {
		buffer.WriteString(fmt.Sprintf(",\"id\":\"%s\"", tp.offerId))
	}

	buffer.WriteString(fmt.Sprintf(",\"offerType\":\"%s\"", tp.offerType))
	buffer.WriteString(fmt.Sprintf(",\"offerVersion\":\"%s\"", tp.version))

	if tp.ETag != "" {
		buffer.WriteString(",\"_etag\":")
		etag, err := json.Marshal(tp.ETag)
		if err != nil {
			return nil, err
		}
		buffer.Write(etag)
	}

	if tp.selfLink != "" {
		buffer.WriteString(fmt.Sprintf(",\"_self\":\"%s\"", tp.selfLink))
	}

	if tp.LastModified != nil {
		buffer.WriteString(",\"_ts\":")
		ts, err := json.Marshal(tp.LastModified)
		if err != nil {
			return nil, err
		}
		buffer.Write(ts)
	}

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func (tp *ThroughputProperties) UnmarshalJSON(b []byte) error {
	var attributes map[string]json.RawMessage
	err := json.Unmarshal(b, &attributes)
	if err != nil {
		return err
	}

	if content, ok := attributes["content"]; ok {
		if err := json.Unmarshal(content, &tp.offer); err != nil {
			return err
		}
	}

	if offerType, ok := attributes["offerType"]; ok {
		if err := json.Unmarshal(offerType, &tp.offerType); err != nil {
			return err
		}
	}

	if version, ok := attributes["offerVersion"]; ok {
		if err := json.Unmarshal(version, &tp.version); err != nil {
			return err
		}
	}

	if offerResourceId, ok := attributes["offerResourceId"]; ok {
		if err := json.Unmarshal(offerResourceId, &tp.offerResourceId); err != nil {
			return err
		}
	}

	if etag, ok := attributes["_etag"]; ok {
		if err := json.Unmarshal(etag, &tp.ETag); err != nil {
			return err
		}
	}

	if ts, ok := attributes["_ts"]; ok {
		if err := json.Unmarshal(ts, &tp.LastModified); err != nil {
			return err
		}
	}

	if id, ok := attributes["id"]; ok {
		if err := json.Unmarshal(id, &tp.offerId); err != nil {
			return err
		}
	}

	if self, ok := attributes["_self"]; ok {
		if err := json.Unmarshal(self, &tp.selfLink); err != nil {
			return err
		}
	}

	return nil
}

// ManualThroughput returns the provisioned throughput in manual mode.
func (tp *ThroughputProperties) ManualThroughput() (int, error) {
	if tp.offer.Throughput == nil {
		return 0, fmt.Errorf("offer is not a manual offer")
	}

	return *tp.offer.Throughput, nil
}

// AutoscaleMaxThroughput returns the configured max throughput on autoscale mode.
func (tp *ThroughputProperties) AutoscaleMaxThroughput() (int, error) {
	if tp.offer.AutoScale == nil {
		return 0, fmt.Errorf("offer is not an autoscale offer")
	}

	return tp.offer.AutoScale.MaxThroughput, nil
}

func (tp *ThroughputProperties) addHeadersToRequest(req *policy.Request) {
	if tp == nil {
		return
	}

	if tp.offer.Throughput != nil {
		req.Raw().Header.Add(cosmosHeaderOfferThroughput, strconv.Itoa(*tp.offer.Throughput))
	} else {
		req.Raw().Header.Add(cosmosHeaderOfferAutoscale, tp.offer.AutoScale.ToJsonString())
	}
}

type offer struct {
	Throughput *int               `json:"offerThroughput,omitempty"`
	AutoScale  *autoscaleSettings `json:"offerAutopilotSettings,omitempty"`
}

func newManualOffer(throughput int) *offer {
	return &offer{
		Throughput: &throughput,
	}
}

func newAutoscaleOfferWithIncrement(startingMaxThroughput int, incrementPercentage int) *offer {
	return &offer{
		AutoScale: &autoscaleSettings{
			MaxThroughput: startingMaxThroughput,
			AutoscaleAutoUpgradeProperties: &autoscaleAutoUpgradeProperties{
				ThroughputPolicy: &autoscaleThroughputPolicy{
					IncrementPercent: incrementPercentage,
				},
			},
		},
	}
}

func newAutoscaleOffer(startingMaxThroughput int) *offer {
	return &offer{
		AutoScale: &autoscaleSettings{
			MaxThroughput: startingMaxThroughput,
		},
	}
}

type autoscaleSettings struct {
	MaxThroughput                  int                             `json:"maxThroughput,omitempty"`
	AutoscaleAutoUpgradeProperties *autoscaleAutoUpgradeProperties `json:"autoUpgradePolicy,omitempty"`
}

func (as *autoscaleSettings) ToJsonString() string {
	if as == nil {
		return ""
	}

	jsonString, _ := json.Marshal(as)

	return string(jsonString)
}

type autoscaleAutoUpgradeProperties struct {
	ThroughputPolicy *autoscaleThroughputPolicy `json:"throughputPolicy,omitempty"`
}

type autoscaleThroughputPolicy struct {
	IncrementPercent int `json:"incrementPercent,omitempty"`
}

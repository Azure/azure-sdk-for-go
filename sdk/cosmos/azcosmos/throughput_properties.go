// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	offerVersion1 string = "V1"
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

	return nil
}

func (tp *ThroughputProperties) ManualThroughput() (int, error) {
	if tp.offer.Throughput == nil {
		return 0, fmt.Errorf("offer is not a manual offer")
	}

	return *tp.offer.Throughput, nil
}

func (tp *ThroughputProperties) AutoscaleMaxThroughput() (int, error) {
	if tp.offer.AutoScale == nil {
		return 0, fmt.Errorf("offer is not an autoscale offer")
	}

	return tp.offer.AutoScale.MaxThroughput, nil
}

func NewManualThroughputProperties(throughput int) *ThroughputProperties {
	return &ThroughputProperties{
		version: offerVersion2,
		offer:   newManualOffer(throughput),
	}
}

func NewAutoscaleThroughputProperties(startingMaxThroughput int, incrementPercentage int) *ThroughputProperties {
	return &ThroughputProperties{
		version: offerVersion2,
		offer:   newAutoscaleOfferWithIncrement(startingMaxThroughput, incrementPercentage),
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
				ThroughputProperties: &autoscaleThroughputProperties{
					IncrementPercent: incrementPercentage,
				},
			},
		},
	}
}

type autoscaleSettings struct {
	MaxThroughput                  int                             `json:"maxThroughput,omitempty"`
	AutoscaleAutoUpgradeProperties *autoscaleAutoUpgradeProperties `json:"autoUpgradePolicy,omitempty"`
}

type autoscaleAutoUpgradeProperties struct {
	ThroughputProperties *autoscaleThroughputProperties `json:"throughputPolicy,omitempty"`
}

type autoscaleThroughputProperties struct {
	IncrementPercent int `json:"incrementPercent,omitempty"`
}

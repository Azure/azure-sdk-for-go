// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	offerVersion2 string = "V2"
)

// ThroughputProperties describes the throughput configuration of a resource.
// It must be initialized through the available constructors.
type ThroughputProperties struct {
	// ETag contains the entity etag of the throughput information.
	ETag *azcore.ETag
	// LastModified contains the last modified time of the throughput information.
	LastModified time.Time

	resource        string
	version         string
	offerType       string
	offer           *offer
	offerResourceId string
	offerId         string
	selfLink        string
}

// NewManualThroughputProperties returns a ThroughputProperties object with the given throughput in manual mode.
// throughput - the throughput in RU/s
func NewManualThroughputProperties(throughput int32) ThroughputProperties {
	return ThroughputProperties{
		version: offerVersion2,
		offer:   newManualOffer(throughput),
	}
}

// NewAutoscaleThroughputPropertiesWithIncrement returns a ThroughputProperties object with the given max throughput on autoscale mode.
// maxThroughput - the max throughput in RU/s
// incrementPercentage - the auto upgrade max throughput increment percentage
func NewAutoscaleThroughputPropertiesWithIncrement(startingMaxThroughput int32, incrementPercentage int32) ThroughputProperties {
	return ThroughputProperties{
		version: offerVersion2,
		offer:   newAutoscaleOfferWithIncrement(startingMaxThroughput, incrementPercentage),
	}
}

// NewAutoscaleThroughputProperties returns a ThroughputProperties object with the given max throughput on autoscale mode.
// maxThroughput - the max throughput in RU/s
func NewAutoscaleThroughputProperties(startingMaxThroughput int32) ThroughputProperties {
	return ThroughputProperties{
		version: offerVersion2,
		offer:   newAutoscaleOffer(startingMaxThroughput),
	}
}

// MarshalJSON implements the json.Marshaler interface
func (tp *ThroughputProperties) MarshalJSON() ([]byte, error) {
	offer, err := json.Marshal(tp.offer)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBufferString("{")
	fmt.Fprint(buffer, "\"content\":")
	buffer.Write(offer)

	if tp.offerResourceId != "" {
		fmt.Fprintf(buffer, ",\"offerResourceId\":\"%s\"", tp.offerResourceId)
	}

	if tp.offerId != "" {
		fmt.Fprintf(buffer, ",\"id\":\"%s\"", tp.offerId)
		fmt.Fprintf(buffer, ",\"_rid\":\"%s\"", tp.offerId)
	}

	fmt.Fprintf(buffer, ",\"offerType\":\"%s\"", tp.offerType)
	fmt.Fprintf(buffer, ",\"offerVersion\":\"%s\"", tp.version)

	if tp.ETag != nil {
		fmt.Fprint(buffer, ",\"_etag\":")
		etag, err := json.Marshal(tp.ETag)
		if err != nil {
			return nil, err
		}
		buffer.Write(etag)
	}

	if tp.selfLink != "" {
		fmt.Fprintf(buffer, ",\"_self\":\"%s\"", tp.selfLink)
	}

	if tp.resource != "" {
		fmt.Fprintf(buffer, ",\"resource\":\"%s\"", tp.resource)
	}

	if !tp.LastModified.IsZero() {
		fmt.Fprintf(buffer, ",\"_ts\":%v", strconv.FormatInt(tp.LastModified.Unix(), 10))
	}

	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
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
		var timestamp int64
		if err := json.Unmarshal(ts, &timestamp); err != nil {
			return err
		}
		tp.LastModified = time.Unix(timestamp, 0)
	}

	if id, ok := attributes["id"]; ok {
		if err := json.Unmarshal(id, &tp.offerId); err != nil {
			return err
		}
	}

	if resource, ok := attributes["resource"]; ok {
		if err := json.Unmarshal(resource, &tp.resource); err != nil {
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
func (tp *ThroughputProperties) ManualThroughput() (int32, bool) {
	if tp.offer.Throughput == nil {
		return 0, false
	}

	return *tp.offer.Throughput, true
}

// AutoscaleMaxThroughput returns the configured max throughput on autoscale mode.
func (tp *ThroughputProperties) AutoscaleMaxThroughput() (int32, bool) {
	if tp.offer.AutoScale == nil {
		return 0, false
	}

	return tp.offer.AutoScale.MaxThroughput, true
}

// AutoscaleIncrement returns the configured percent increment on autoscale mode.
func (tp *ThroughputProperties) AutoscaleIncrement() (int32, bool) {
	if tp.offer.AutoScale == nil ||
		tp.offer.AutoScale.AutoscaleAutoUpgradeProperties == nil ||
		tp.offer.AutoScale.AutoscaleAutoUpgradeProperties.ThroughputPolicy == nil {
		return 0, false
	}

	return tp.offer.AutoScale.AutoscaleAutoUpgradeProperties.ThroughputPolicy.IncrementPercent, true
}

func (tp *ThroughputProperties) addHeadersToRequest(req *policy.Request) {
	if tp == nil {
		return
	}

	if tp.offer.Throughput != nil {
		req.Raw().Header.Add(cosmosHeaderOfferThroughput, strconv.Itoa(int(*tp.offer.Throughput)))
	} else {
		req.Raw().Header.Add(cosmosHeaderOfferAutoscale, tp.offer.AutoScale.ToJsonString())
	}
}

type offer struct {
	Throughput *int32             `json:"offerThroughput,omitempty"`
	AutoScale  *autoscaleSettings `json:"offerAutopilotSettings,omitempty"`
}

func newManualOffer(throughput int32) *offer {
	return &offer{
		Throughput: &throughput,
	}
}

func newAutoscaleOfferWithIncrement(startingMaxThroughput int32, incrementPercentage int32) *offer {
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

func newAutoscaleOffer(startingMaxThroughput int32) *offer {
	return &offer{
		AutoScale: &autoscaleSettings{
			MaxThroughput: startingMaxThroughput,
		},
	}
}

type autoscaleSettings struct {
	MaxThroughput                  int32                           `json:"maxThroughput,omitempty"`
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
	IncrementPercent int32 `json:"incrementPercent,omitempty"`
}

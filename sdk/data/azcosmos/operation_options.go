// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// OperationOptions configures requests, client defaults, and runtime defaults.
// Unset fields inherit; resolution and environment overrides are owned by the Rust driver.
type OperationOptions struct {
	// ConsistencyStrategy controls read consistency. Zero inherits; explicit Default uses the account level.
	ConsistencyStrategy ReadConsistencyStrategy
	// EnableContentResponseOnWrite overrides whether writes return item content. Nil inherits.
	EnableContentResponseOnWrite *bool
	// ExcludedRegions replaces inherited exclusions. Nil inherits; an empty non-nil slice clears them.
	ExcludedRegions []Region
	// EndToEndTimeout bounds initialization, metadata lookup, and execution. Zero inherits.
	// Context deadlines can stop the Go wait sooner than the driver's one-second minimum,
	// but do not cancel submitted native work.
	EndToEndTimeout time.Duration
	// PatchStrategy selects patch execution. A patch-specific Strategy takes precedence.
	PatchStrategy PatchStrategy
	// SessionCapturingDisabled disables automatic session capture/resolution, not explicit tokens.
	SessionCapturingDisabled *bool
	// MaxFailoverRetryCount overrides the failover retry limit, including explicit zero.
	MaxFailoverRetryCount *uint32
	// MaxSessionRetryCount overrides the session retry limit, including explicit zero.
	MaxSessionRetryCount *uint32
	// EndpointUnavailabilityTTL controls how long unavailable endpoints are avoided. Nil inherits.
	EndpointUnavailabilityTTL *time.Duration
	// CustomHeaders replaces inherited custom headers. Nil inherits; an empty map clears them.
	CustomHeaders map[string]string
	// BinaryEncoding overrides the entire encoding group. Nil inherits the driver's binary default.
	BinaryEncoding *BinaryEncodingOptions
	// ThroughputControl configures independently inherited throughput controls.
	ThroughputControl ThroughputControlOptions
	// ThrottlingRetry configures each transport invocation's 429 retry budget, not the whole operation.
	ThrottlingRetry ThrottlingRetryOptions
	// HedgingEnabled controls the hedging master switch; the driver's environment override wins.
	HedgingEnabled *bool
	// AvailabilityStrategy selects disabled or threshold-based hedging. Zero inherits.
	AvailabilityStrategy AvailabilityStrategy
}

// PriorityLevel selects the service request priority.
type PriorityLevel string

const (
	// PriorityLevelUnset inherits the priority.
	PriorityLevelUnset PriorityLevel = ""
	// PriorityLevelHigh selects high priority.
	PriorityLevelHigh PriorityLevel = "High"
	// PriorityLevelLow selects low priority.
	PriorityLevelLow PriorityLevel = "Low"
)

// ThroughputControlOptions configures service throughput controls. Fields inherit independently.
type ThroughputControlOptions struct {
	// ThroughputBucket selects a bucket, including explicit zero. Nil inherits.
	ThroughputBucket *uint32
	// PriorityLevel selects a service priority. Zero inherits.
	PriorityLevel PriorityLevel
}

// ThrottlingRetryOptions configures retries for one transport-pipeline invocation.
type ThrottlingRetryOptions struct {
	// MaxRetryCount allows this many retries after the initial attempt. Zero disables retries.
	MaxRetryCount *uint32
	// MaxRetryWaitTime bounds cumulative throttling delay. Nil inherits; zero is explicit.
	MaxRetryWaitTime *time.Duration
}

// AvailabilityStrategy is an immutable whole-value override. Its zero value inherits.
type AvailabilityStrategy struct {
	kind      uint8
	threshold time.Duration
}

// DisabledAvailability disables availability hedging unless the driver's master switch enables it.
func DisabledAvailability() AvailabilityStrategy {
	return AvailabilityStrategy{kind: 1}
}

// HedgingAvailability selects hedging with a strictly positive threshold.
// An invalid threshold is rejected when the strategy is used.
func HedgingAvailability(threshold time.Duration) AvailabilityStrategy {
	return AvailabilityStrategy{kind: 2, threshold: threshold}
}

// BinaryEncodingOptions controls Cosmos binary JSON wire encoding and response conversion.
// A supplied zero value disables binary encoding. Raw responses may be binary when enabled.
type BinaryEncodingOptions struct {
	// Enabled permits binary JSON on the wire.
	Enabled bool
	// RequestTextResponse converts binary responses to text JSON without disabling binary wire encoding.
	RequestTextResponse bool
}

func clonePointer[T any](value *T) *T {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func (o OperationOptions) clone() OperationOptions {
	o.EnableContentResponseOnWrite = clonePointer(o.EnableContentResponseOnWrite)
	o.ExcludedRegions = slices.Clone(o.ExcludedRegions)
	o.SessionCapturingDisabled = clonePointer(o.SessionCapturingDisabled)
	o.MaxFailoverRetryCount = clonePointer(o.MaxFailoverRetryCount)
	o.MaxSessionRetryCount = clonePointer(o.MaxSessionRetryCount)
	o.EndpointUnavailabilityTTL = clonePointer(o.EndpointUnavailabilityTTL)
	o.BinaryEncoding = clonePointer(o.BinaryEncoding)
	o.CustomHeaders = maps.Clone(o.CustomHeaders)
	o.ThroughputControl.ThroughputBucket = clonePointer(o.ThroughputControl.ThroughputBucket)
	o.ThrottlingRetry.MaxRetryCount = clonePointer(o.ThrottlingRetry.MaxRetryCount)
	o.ThrottlingRetry.MaxRetryWaitTime = clonePointer(o.ThrottlingRetry.MaxRetryWaitTime)
	o.HedgingEnabled = clonePointer(o.HedgingEnabled)
	return o
}

func (o OperationOptions) validate() error {
	if err := o.ConsistencyStrategy.validate(); err != nil {
		return err
	}
	if err := o.PatchStrategy.validate(); err != nil {
		return err
	}
	if o.EndToEndTimeout < 0 {
		return errors.New("azcosmos: EndToEndTimeout must not be negative")
	}
	if o.EndpointUnavailabilityTTL != nil && *o.EndpointUnavailabilityTTL < 0 {
		return errors.New("azcosmos: EndpointUnavailabilityTTL must not be negative")
	}
	if o.ThrottlingRetry.MaxRetryWaitTime != nil && *o.ThrottlingRetry.MaxRetryWaitTime < 0 {
		return errors.New("azcosmos: MaxRetryWaitTime must not be negative")
	}
	switch o.ThroughputControl.PriorityLevel {
	case PriorityLevelUnset, PriorityLevelHigh, PriorityLevelLow:
	default:
		return fmt.Errorf("azcosmos: unknown priority level %q", o.ThroughputControl.PriorityLevel)
	}
	if o.AvailabilityStrategy.kind > 2 || o.AvailabilityStrategy.kind == 2 && o.AvailabilityStrategy.threshold <= 0 {
		return errors.New("azcosmos: hedging availability requires a positive threshold")
	}
	for _, region := range o.ExcludedRegions {
		if region == "" || strings.ContainsRune(string(region), 0) {
			return errors.New("azcosmos: ExcludedRegions must contain nonempty regions without NUL bytes")
		}
	}
	seen := make(map[string]bool, len(o.CustomHeaders))
	for _, name := range slices.Sorted(maps.Keys(o.CustomHeaders)) {
		value := o.CustomHeaders[name]
		if !utf8.ValidString(value) {
			return fmt.Errorf("azcosmos: invalid UTF-8 value for custom header %q", name)
		}
		if name == "" {
			return errors.New("azcosmos: custom header name must not be empty")
		}
		for i := range len(name) {
			c := name[i]
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
				strings.ContainsRune("!#$%&'*+-.^_`|~", rune(c))) {
				return fmt.Errorf("azcosmos: invalid custom header name %q", name)
			}
		}
		normalized := strings.ToLower(name)
		if seen[normalized] {
			return fmt.Errorf("azcosmos: duplicate case-insensitive custom header %q", name)
		}
		seen[normalized] = true
		for i := range len(value) {
			if value[i] < 32 && value[i] != '\t' || value[i] == 127 {
				return fmt.Errorf("azcosmos: invalid value for custom header %q", name)
			}
		}
	}
	return nil
}

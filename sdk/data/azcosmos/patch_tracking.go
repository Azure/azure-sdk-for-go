// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/hex"
	"errors"
	"math"
	"strings"
	"time"
)

// PatchTrackingID is a UUID identifying a logical client-side patch across application retries.
// Tracking modifies reserved item metadata; it is bounded by retention and capacity, not permanent.
type PatchTrackingID string

func (o PatchItemOptions) validateTracking() error {
	if o.MaxAttempts != nil && *o.MaxAttempts == 0 {
		return errors.New("azcosmos: MaxAttempts must be between 1 and 255")
	}
	if o.TrackingCapacity != nil && *o.TrackingCapacity == 0 {
		return errors.New("azcosmos: TrackingCapacity must be between 1 and 65535")
	}
	if o.TrackingRetention != nil && (*o.TrackingRetention < 0 || *o.TrackingRetention/time.Second > math.MaxUint32) {
		return errors.New("azcosmos: TrackingRetention must be nonnegative and fit uint32 seconds")
	}
	if o.TrackingID != "" {
		id := string(o.TrackingID)
		if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
			return errors.New("azcosmos: TrackingID must be a hyphenated UUID")
		}
		if decoded, err := hex.DecodeString(strings.ReplaceAll(id, "-", "")); err != nil || len(decoded) != 16 {
			return errors.New("azcosmos: TrackingID must be a hyphenated UUID")
		}
	}
	return nil
}

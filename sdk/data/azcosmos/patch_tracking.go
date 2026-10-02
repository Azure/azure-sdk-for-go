// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/hex"
	"errors"
	"strings"
)

// PatchTrackingID is a UUID identifying a logical client-side patch across application retries.
// Tracking modifies reserved item metadata; it is bounded by retention and capacity, not permanent.
type PatchTrackingID string

func (id PatchTrackingID) normalized() (PatchTrackingID, error) {
	if id == "" {
		return "", nil
	}
	value := string(id)
	if strings.HasPrefix(value, "urn:uuid:") {
		value = strings.TrimPrefix(value, "urn:uuid:")
	} else if len(value) == 38 && value[0] == '{' && value[37] == '}' {
		value = value[1:37]
	}
	if len(value) == 36 && value[8] == '-' && value[13] == '-' && value[18] == '-' && value[23] == '-' {
		value = strings.ReplaceAll(value, "-", "")
	}
	decoded, err := hex.DecodeString(value)
	if err != nil || len(value) != 32 || len(decoded) != 16 {
		return "", errors.New("azcosmos: TrackingID must be a UUID")
	}
	value = strings.ToLower(value)
	return PatchTrackingID(value[:8] + "-" + value[8:12] + "-" + value[12:16] + "-" + value[16:20] + "-" + value[20:]), nil
}

func (o PatchItemOptions) validateTracking() error {
	if o.MaxAttempts != nil && *o.MaxAttempts == 0 {
		return errors.New("azcosmos: MaxAttempts must be between 1 and 255")
	}
	if o.TrackingCapacity != nil && *o.TrackingCapacity == 0 {
		return errors.New("azcosmos: TrackingCapacity must be between 1 and 65535")
	}
	if o.TrackingRetention != nil && *o.TrackingRetention < 0 {
		return errors.New("azcosmos: TrackingRetention must not be negative")
	}
	_, err := o.TrackingID.normalized()
	return err
}

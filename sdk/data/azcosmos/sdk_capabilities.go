// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "strconv"

type supportedCapabilities uint64

const (
	supportedCapabilitiesNone           supportedCapabilities = 0
	supportedCapabilitiesPartitionMerge supportedCapabilities = 1 << 0
)

var supportedCapabilitiesHeaderValue = supportedCapabilitiesAsString()

func supportedCapabilitiesAsString() string {
	supported := supportedCapabilitiesNone
	supported |= supportedCapabilitiesPartitionMerge
	return strconv.FormatUint(uint64(supported), 10)
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

//type IotHubDeviceConnectedEventData DeviceConnectionStateEventProperties

func fixNAValue(s **string) {
	if *s != nil && **s == "n/a" {
		*s = nil
	}
}

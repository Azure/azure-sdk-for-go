//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"time"
)

func (client *FileSystemClient) Endpoint() string {
	return client.endpoint
}

func (client *FileSystemClient) Pipeline() runtime.Pipeline {
	return client.internal.Pipeline()
}

// used to convert times from UTC to GMT before sending across the wire
var gmt = time.FixedZone("GMT", 0)

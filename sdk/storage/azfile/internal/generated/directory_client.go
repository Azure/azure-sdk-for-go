//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

const (
	// ISO8601 is used for formatting file creation, last write and change time.
	ISO8601 = "2006-01-02T15:04:05.0000000Z07:00"
)

func (client *DirectoryClient) Endpoint() string {
	return client.endpoint
}

func (client *DirectoryClient) Pipeline() runtime.Pipeline {
	return client.pl
}

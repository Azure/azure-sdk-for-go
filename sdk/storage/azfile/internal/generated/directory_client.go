//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	// ISO8601 is used for formatting file creation, last write and change time.
	ISO8601 = "2006-01-02T15:04:05.0000000Z07:00"
)

func (client *DirectoryClient) Endpoint() string {
	return client.endpoint
}

func (client *DirectoryClient) InternalClient() *azcore.Client {
	return client.internal
}

// NewDirectoryClient creates a new instance of DirectoryClient with the specified values.
//   - endpoint - The URL of the service account, share, directory or file that is the target of the desired operation.
//   - allowTrailingDot - If true, the trailing dot will not be trimmed from the target URI.
//   - fileRequestIntent - Valid value is backup
//   - allowSourceTrailingDot - If true, the trailing dot will not be trimmed from the source URI.
//   - azClient - azcore.Client is a basic HTTP client.  It consists of a pipeline and tracing provider.
func NewDirectoryClient(endpoint string, allowTrailingDot *bool, fileRequestIntent *ShareTokenIntent, allowSourceTrailingDot *bool, azClient *azcore.Client) *DirectoryClient {
	client := &DirectoryClient{
		internal:               azClient,
		endpoint:               endpoint,
		allowTrailingDot:       allowTrailingDot,
		fileRequestIntent:      fileRequestIntent,
		allowSourceTrailingDot: allowSourceTrailingDot,
	}
	return client
}

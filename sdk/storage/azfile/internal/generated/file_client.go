// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func (client *FileClient) Endpoint() string {
	return client.endpoint
}

func (client *FileClient) InternalClient() *azcore.Client {
	return client.internal
}

// NewFileClient creates a new instance of FileClient with the specified values.
//   - endpoint - The URL of the service account, share, directory or file that is the target of the desired operation.
//   - allowTrailingDot - If true, the trailing dot will not be trimmed from the target URI.
//   - fileRequestIntent - Valid value is backup
//   - allowSourceTrailingDot - If true, the trailing dot will not be trimmed from the source URI.
//   - azClient - azcore.Client is a basic HTTP client.  It consists of a pipeline and tracing provider.
func NewFileClient(endpoint string, allowTrailingDot *bool, fileRequestIntent *ShareTokenIntent, allowSourceTrailingDot *bool, azClient *azcore.Client) *FileClient {
	client := &FileClient{
		internal:               azClient,
		endpoint:               endpoint,
		allowTrailingDot:       allowTrailingDot,
		version:                ServiceVersion,
		fileRequestIntent:      fileRequestIntent,
		allowSourceTrailingDot: allowSourceTrailingDot,
	}
	return client
}
